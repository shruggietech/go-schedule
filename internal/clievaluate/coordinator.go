package clievaluate

import (
	"bytes"
	"context"
	"debug/buildinfo"
	"errors"
	"fmt"
	"net/url"
	"os"
	pathpkg "path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/lock"
)

// CaptureOptions names existing executables and explicitly selected environments.
type CaptureOptions struct {
	Root         string
	Source       string
	BinDirectory string
	WSLDistro    string
	WSLUser      string
	Resume       bool
	Budget       Budget
}

type captureIdentity struct {
	RunID          string  `json:"run_id"`
	Subject        Subject `json:"subject"`
	Budget         Budget  `json:"budget"`
	Distro         string  `json:"distro"`
	User           string  `json:"user"`
	ManifestSHA256 string  `json:"manifest_sha256"`
	Attempt        int     `json:"attempt"`
}

// Capture coordinates sequential contained workers and freezes a review corpus.
func Capture(ctx context.Context, options CaptureOptions) (result Report, resultErr error) {
	source, sourceErr := filepath.Abs(options.Source)
	if sourceErr != nil {
		return result, sourceErr
	}
	options.Source = source
	if runtime.GOOS != "windows" {
		return result, errors.New("the coordinator runs on Windows and selects an existing WSL2 environment explicitly")
	}
	if err := ValidateBudget(options.Budget); err != nil {
		return result, err
	}
	ctx, cancelRun := context.WithTimeout(ctx, time.Duration(options.Budget.TotalTimeoutMS)*time.Millisecond)
	defer cancelRun()
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return result, err
	}
	if options.Resume {
		info, err := os.Lstat(root)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return result, errors.Join(errors.New("resume requires the existing regular run directory"), err)
		}
	} else if err := os.Mkdir(root, 0700); err != nil {
		return result, err
	}
	runLock, err := lock.AcquireFile(filepath.Join(root, "coordinator.lock"), 0600)
	if err != nil {
		return result, err
	}
	defer func() { resultErr = errors.Join(resultErr, runLock.Release()) }()
	store, err := OpenArtifacts(root, options.Budget.ArtifactBytes, NewRedactor())
	if err != nil {
		return result, err
	}
	defer func() {
		if store != nil {
			resultErr = errors.Join(resultErr, store.Close())
		}
	}()
	manifest, _, err := BuildManifest()
	if err != nil {
		return result, err
	}
	manifestBytes, err := Encode(manifest)
	if err != nil {
		return result, err
	}
	subject, subjectErr := captureSubject(ctx, options)
	identity := captureIdentity{RunID: uuid.NewString(), Subject: subject, Budget: options.Budget, Distro: options.WSLDistro, User: options.WSLUser, ManifestSHA256: Digest(manifestBytes), Attempt: 1}
	priorCleanup := map[string]Cleanup{}
	priorOwnership := []Artifact{}
	if options.Resume {
		var prior captureIdentity
		data, err := store.readRegular("identity.json", documentLimit)
		if err != nil {
			return result, err
		}
		if err := Decode(bytes.NewReader(data), &prior); err != nil {
			return result, err
		}
		if err := compareResumeInputs(prior, identity); err != nil {
			return result, err
		}
		if err := confirmPreviousWorkersStopped(ctx, store, prior); err != nil {
			return result, err
		}
		if err := validateResumeEvidence(store, prior, manifest); err != nil {
			return result, fmt.Errorf("resume evidence integrity: %w", err)
		}
		previous, _, _, err := loadRun(store)
		if err != nil {
			return result, err
		}
		if previous.CompletionState == "running" {
			// The root ledger precedes worker startup. Recover durable ownership
			// from interrupted workers as well as the last aggregate checkpoint.
			for _, env := range []string{"windows", "wsl"} {
				var worker Report
				path := fmt.Sprintf("attempt-%03d/%s/report.json", prior.Attempt, env)
				worker, err = loadWorkerCheckpoint(store, strings.TrimSuffix(path, "/report.json"))
				if errors.Is(err, os.ErrNotExist) {
					continue
				} else if err != nil {
					return result, err
				}
				for i, item := range previous.Cleanup {
					if item.EnvironmentID == env {
						previous.Cleanup[i] = worker.Cleanup[0]
					}
				}
				previous.Artifacts = append(previous.Artifacts, worker.Artifacts...)
			}
		}
		priorCleanup, priorOwnership = resumeCleanup(previous)
		identity = prior
		identity.Attempt++
		identity.RunID = uuid.NewString()
	}
	// A fresh attempt replays stateful cases; prior artifacts never imply a new success.
	identityBytes, err := Encode(identity)
	if err != nil {
		return result, err
	}
	if err := store.Replace("identity.json", identityBytes); err != nil {
		return result, err
	}
	if err := store.Replace("manifest.json", manifestBytes); err != nil {
		return result, err
	}
	now := stamp()
	result = Report{SchemaVersion: Version, Kind: "go-schedule.cli-qa-report", RunID: identity.RunID, StartedAt: now, FinishedAt: now, Phase: "capturing", CompletionState: "running", Subject: subject, Runner: Runner{Version: Version, ManifestSHA256: Digest(manifestBytes), Budget: options.Budget, Parameters: []NamedValue{{Name: "attempt", Value: strconv.Itoa(identity.Attempt)}, {Name: "wsl_distro", Value: options.WSLDistro}, {Name: "wsl_user", Value: options.WSLUser}}}}
	result.Artifacts = append(result.Artifacts, priorOwnership...)
	corpus := Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: identity.RunID}
	events := []byte{}
	envs := []Environment{{ID: "windows", OS: "windows", Arch: runtime.GOARCH, Version: runtime.Version(), Status: "blocked", ReasonCode: "not_started"}, {ID: "wsl", OS: "linux", Arch: "amd64", Distro: &options.WSLDistro, User: &options.WSLUser, Version: "preflight pending", Status: "blocked", ReasonCode: "not_started"}}
	checkpoint := result
	for _, env := range envs {
		pending := blockedWorker(identity.RunID, subject, env, manifest, now, "not_started")
		pending.Cleanup[0] = mergeCleanup(priorCleanup[env.ID], pending.Cleanup[0])
		checkpoint.Environments = append(checkpoint.Environments, pending.Environments...)
		checkpoint.Cases = append(checkpoint.Cases, pending.Cases...)
		checkpoint.Cleanup = append(checkpoint.Cleanup, pending.Cleanup...)
	}
	if err := persistCapture(store, &checkpoint, manifest, corpus, nil); err != nil {
		return result, err
	}
	// Each worker writes its complete initial ledger before touching its owned daemon.
	for _, env := range envs {
		if subjectErr != nil {
			worker := blockedWorker(identity.RunID, subject, env, manifest, now, "subject_unavailable")
			worker.Cleanup[0] = mergeCleanup(priorCleanup[env.ID], worker.Cleanup[0])
			worker.Environments[0].Tools = []NamedValue{{Name: "subject_error", Value: subjectErr.Error()}}
			result.Environments = append(result.Environments, worker.Environments...)
			result.Cases = append(result.Cases, worker.Cases...)
			result.Cleanup = append(result.Cleanup, worker.Cleanup...)
			result.CompletionState = "interrupted"
			continue
		}
		prefix := fmt.Sprintf("attempt-%03d/%s", identity.Attempt, env.ID)
		if err := store.root.MkdirAll(prefix, 0700); err != nil {
			return result, err
		}
		workerSubject := subject
		workerSubject.Binaries = append([]Binary{}, subject.Binaries...)
		if env.OS == "linux" {
			for i := range workerSubject.Binaries {
				if workerSubject.Binaries[i].OS == "linux" {
					workerSubject.Binaries[i].Path, err = wslPath(workerSubject.Binaries[i].Path)
					if err != nil {
						return result, err
					}
				}
			}
		}
		config := WorkerConfig{RunID: identity.RunID, Root: root, Prefix: prefix, Subject: workerSubject, Budget: options.Budget, Environment: env}
		configName := prefix + "/worker.json"
		workerExe := binaryPath(subject, "runner", env.OS)
		configArg := filepath.Join(root, filepath.FromSlash(configName))
		command := clifixture.Command{Executable: workerExe, Arguments: []string{"worker", "--config", configArg}, Environment: quietEnvironment(), Directory: options.Source, Limits: &clifixture.ProcessLimits{CPUPercent: options.Budget.CPUPercent, MemoryBytes: uint64(options.Budget.MemoryBytes), MaxProcesses: uint32(options.Budget.WindowsProcesses)}}
		if env.OS == "linux" {
			config.Root, err = wslPath(root)
			if err != nil {
				return result, err
			}
			workerExe, err = wslPath(workerExe)
			if err != nil {
				return result, err
			}
			configArg, err = wslPath(configArg)
			if err != nil {
				return result, err
			}
			args, e := WSLWorkerArguments(identity.RunID, config.Root, workerExe, []string{"worker", "--config", configArg}, options.Budget)
			if e != nil {
				return result, e
			}
			command.Executable = filepath.Join(os.Getenv("WINDIR"), "System32", "wsl.exe")
			command.Arguments = append([]string{"--distribution", options.WSLDistro, "--user", options.WSLUser, "--exec"}, args...)
		}
		data, err := Encode(config)
		if err != nil {
			return result, err
		}
		if err := store.Replace(configName, data); err != nil {
			return result, err
		}
		workerCtx, cancel := context.WithTimeout(ctx, time.Duration(options.Budget.TotalTimeoutMS)*time.Millisecond+45*time.Second)
		capture := captureWithStart(workerCtx, command, nil, options.Budget.ChannelBytes, true, func(pid int) error {
			created, err := processIdentity(pid)
			if err != nil {
				return err
			}
			data, err := Encode(workerLease{pid, created, identity.RunID})
			if err != nil {
				return err
			}
			return store.Replace(prefix+"/lease.json", data)
		})
		cancel()
		// Reopen after the worker's writes so aggregate accounting includes all bytes.
		if err := store.Close(); err != nil {
			return result, err
		}
		store, err = OpenArtifacts(root, options.Budget.ArtifactBytes, NewRedactor())
		if err != nil {
			return result, err
		}
		var worker Report
		var messages Corpus
		worker, loadErr := loadWorkerCheckpoint(store, prefix)
		if loadErr == nil {
			loadErr = store.ReadDocument(prefix+"/corpus.json", "corpus", &messages)
			if errors.Is(loadErr, os.ErrNotExist) {
				messages, loadErr = freezeChannels(store, &worker)
				worker.Phase = "pending_review"
				worker.CompletionState = "interrupted"
				worker.FinishedAt = stamp()
				if loadErr == nil {
					data, e := Encode(messages)
					loadErr = e
					worker.CorpusSHA256 = Digest(data)
					worker.Summary = DeriveSummary(worker)
				}
			}
		}
		if loadErr != nil {
			worker = blockedWorker(identity.RunID, subject, env, manifest, now, "worker_unavailable")
			detail := string(capture.Err)
			if capture.Failure != nil {
				detail += "\n" + capture.Failure.Error()
			}
			worker.Environments[0].Tools = []NamedValue{{Name: "worker_error", Value: detail}}
			messages = Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: identity.RunID}
		} else if worker.CompletionState != "completed" || capture.Code != 0 || capture.Cleanup != nil {
			worker.CompletionState = "interrupted"
			for i := range worker.Cases {
				if worker.Cases[i].ReasonCode == "not_started" {
					worker.Cases[i].ReasonCode = "worker_interrupted"
				}
			}
		}
		for _, message := range messages.Messages {
			for _, ref := range message.Cases {
				AddMessage(&corpus, CaseResult{EnvironmentID: ref.EnvironmentID, ID: ref.CaseID, CommandPath: message.CommandPath, ScenarioClass: message.ScenarioClass}, message.Channel, message.Text)
			}
		}
		corpus.NoMessages = append(corpus.NoMessages, messages.NoMessages...)
		result.Environments = append(result.Environments, worker.Environments...)
		result.Cases = append(result.Cases, worker.Cases...)
		result.Artifacts = append(result.Artifacts, worker.Artifacts...)
		for _, cleanup := range worker.Cleanup {
			result.Cleanup = append(result.Cleanup, mergeCleanup(priorCleanup[cleanup.EnvironmentID], cleanup))
		}
		if worker.CompletionState == "interrupted" {
			result.CompletionState = "interrupted"
		}
		eventData, e := store.readRegular(prefix+"/events.jsonl", documentLimit)
		if e == nil {
			events = append(events, eventData...)
		}
	}
	SortCorpus(&corpus)
	result.Phase = "pending_review"
	if result.CompletionState == "running" {
		result.CompletionState = "completed"
	}
	result.FinishedAt = stamp()
	result.Summary = DeriveSummary(result)
	if err := persistCapture(store, &result, manifest, corpus, events); err != nil {
		return result, err
	}
	return result, nil
}

func blockedWorker(runID string, subject Subject, env Environment, m Manifest, now, reason string) Report {
	env.Status = "blocked"
	env.ReasonCode = reason
	r := Report{RunID: runID, Subject: subject, Environments: []Environment{env}, CompletionState: "interrupted", Cleanup: []Cleanup{{EnvironmentID: env.ID}}}
	for _, spec := range m.Cases {
		status, why := "blocked", reason
		if !slices.Contains(spec.Platforms, env.OS) {
			status, why = "not_applicable", spec.ApplicabilityReason
		}
		r.Cases = append(r.Cases, CaseResult{ID: spec.ID, EnvironmentID: env.ID, CommandPath: spec.CommandPath, ScenarioClass: spec.ScenarioClass, ExecutionKind: spec.ExecutionKind, Status: status, ExecutionStatus: status, ReasonCode: why, Argv: spec.Argv, ExpectedExit: spec.ExpectedExit, ExpectedChecks: spec.Checks, StartedAt: now, FinishedAt: stamp()})
	}
	return r
}

func binaryPath(subject Subject, role, platform string) string {
	for _, binary := range subject.Binaries {
		if binary.Role == role && binary.OS == platform {
			return binary.Path
		}
	}
	return ""
}
func wslPath(path string) (string, error) {
	if strings.HasPrefix(path, "/mnt/") && len(path) > 7 && path[5] >= 'a' && path[5] <= 'z' && path[6] == '/' {
		if clean := pathpkg.Clean(path); clean == path {
			return path, nil
		}
		return "", errors.New("WSL path must be canonical")
	}
	path = filepath.Clean(path)
	if len(path) < 3 || path[1] != ':' || path[2] != '\\' {
		return "", errors.New("WSL inputs require an absolute local drive path")
	}
	return "/mnt/" + strings.ToLower(path[:1]) + strings.ReplaceAll(path[2:], "\\", "/"), nil
}

func captureSubject(ctx context.Context, o CaptureOptions) (Subject, error) {
	s, err := identifySource(ctx, o.Source)
	if err != nil {
		return s, err
	}
	for _, platform := range []string{"windows", "linux"} {
		for _, role := range []string{"cli", "daemon", "runner"} {
			path, err := filepath.Abs(filepath.Join(o.BinDirectory, filepath.FromSlash(preparedName(role, platform))))
			if err != nil {
				return s, err
			}
			data, err := readBinary(path)
			if err != nil {
				return s, err
			}
			info, err := buildinfo.Read(bytes.NewReader(data))
			if err != nil {
				return s, errors.New("evaluation input has no readable Go build identity")
			}
			settings := map[string]string{}
			for _, setting := range info.Settings {
				settings[setting.Key] = setting.Value
			}
			if settings["vcs.revision"] != s.Commit || settings["GOOS"] != platform || settings["GOARCH"] != "amd64" {
				return s, errors.New("evaluation binary source revision or platform differs from the selected checkout")
			}
			if (settings["vcs.modified"] == "true") != s.Dirty {
				return s, errors.New("evaluation binary dirtiness differs from the selected checkout")
			}
			if platform == "linux" {
				path, err = wslPath(path)
				if err != nil {
					return s, err
				}
			}
			s.Binaries = append(s.Binaries, Binary{role, platform, "amd64", path, Digest(data)})
		}
	}
	if err := requireBuildReceipt(o.BinDirectory, s); err != nil {
		return s, err
	}
	version := captureCommand(ctx, clifixture.Command{Executable: binaryPath(s, "cli", "windows"), Arguments: []string{"--version"}, Directory: o.Source, Environment: quietEnvironment()}, nil, 4096)
	if version.Code != 0 || version.Failure != nil {
		return s, errors.New("cannot read the evaluated CLI version")
	}
	s.Version = strings.TrimSpace(string(version.Out))
	return s, nil
}

func canonicalRemote(value string) (string, error) {
	value = strings.TrimSuffix(strings.TrimSpace(value), ".git")
	if strings.HasPrefix(value, "git@") {
		host, path, ok := strings.Cut(strings.TrimPrefix(value, "git@"), ":")
		if !ok || host == "" || path == "" {
			return "", errors.New("source remote identity is invalid")
		}
		return "https://" + host + "/" + path, nil
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || parsed.Path == "" {
		return "", errors.New("source remote identity is invalid")
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	if parsed.Scheme == "ssh" {
		parsed.Scheme = "https"
	}
	return parsed.String(), nil
}

func loadRun(store *ArtifactStore) (Report, Manifest, Corpus, error) {
	var r Report
	var m Manifest
	var c Corpus
	for _, d := range []struct {
		name, kind string
		target     any
	}{{"report.json", "report", &r}, {"manifest.json", "manifest", &m}, {"corpus.json", "corpus", &c}} {
		if err := store.ReadDocument(d.name, d.kind, d.target); err != nil {
			return r, m, c, err
		}
	}
	if err := validateFrozen(store, r); err != nil {
		return r, m, c, err
	}
	events, err := store.readRegular("events.jsonl", documentLimit)
	if err != nil {
		return r, m, c, err
	}
	if err := validateEvents(events, r); err != nil {
		return r, m, c, err
	}
	return r, m, c, ValidateReport(store, r, m, c)
}
func saveReport(store *ArtifactStore, r Report) error {
	data, err := Encode(r)
	if err != nil {
		return err
	}
	markdown, err := Render(r)
	if err != nil {
		return err
	}
	if err := store.Replace("report.json", data); err != nil {
		return err
	}
	return store.Replace("report.md", markdown)
}

// ProcessReport validates all evidence and optionally applies external agent reviews.
func ProcessReport(root, reviewPath string) (r Report, resultErr error) {
	held, err := lock.AcquireFile(filepath.Join(root, "coordinator.lock"), 0600)
	if err != nil {
		return r, err
	}
	defer func() { resultErr = errors.Join(resultErr, held.Release()) }()
	store, err := OpenArtifacts(root, 256<<20, NewRedactor())
	if err != nil {
		return r, err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	r, m, c, err := loadRun(store)
	if err != nil {
		return r, err
	}
	store.limit = r.Runner.Budget.ArtifactBytes
	if store.used > store.limit {
		return r, errors.New("run directory exceeds its recorded total artifact budget")
	}
	if reviewPath != "" {
		var set ReviewSet
		if err := ReadInput(reviewPath, &set); err != nil {
			return r, err
		}
		if err := ApplyReviews(&r, c, set); err != nil {
			return r, err
		}
		if err := ValidateReport(store, r, m, c); err != nil {
			return r, err
		}
	}
	return r, saveReport(store, r)
}

// ValidateRun checks a frozen report and its evidence without rewriting files.
func ValidateRun(root string) (r Report, resultErr error) {
	held, err := lock.AcquireFile(filepath.Join(root, "coordinator.lock"), 0600)
	if err != nil {
		return r, err
	}
	defer func() { resultErr = errors.Join(resultErr, held.Release()) }()
	store, err := OpenArtifacts(root, 256<<20, NewRedactor())
	if err != nil {
		return r, err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	r, _, _, err = loadRun(store)
	if err != nil {
		return r, err
	}
	if store.used > r.Runner.Budget.ArtifactBytes {
		return r, errors.New("run directory exceeds its recorded total artifact budget")
	}
	return r, nil
}
