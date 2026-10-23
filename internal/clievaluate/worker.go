package clievaluate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

// WorkerConfig binds one contained worker to immutable source and output inputs.
type WorkerConfig struct {
	RunID       string      `json:"run_id"`
	Root        string      `json:"root"`
	Prefix      string      `json:"prefix"`
	Subject     Subject     `json:"subject"`
	Budget      Budget      `json:"budget"`
	Environment Environment `json:"environment"`
}

func stamp() string { return time.Now().UTC().Format(time.RFC3339Nano) }

// RunWorker executes only after verifying its actual OS containment controls.
func RunWorker(ctx context.Context, c WorkerConfig) (resultErr error) {
	if !validRunID(c.RunID) || ValidateArtifactPath(c.Prefix) != nil {
		return errors.New("invalid worker identity or prefix")
	}
	if err := ValidateBudget(c.Budget); err != nil {
		return err
	}
	redactor := NewRedactor()
	for _, value := range []string{"fixture-webhook-secret", "fixture-webhook-rotated"} {
		redactor.Register(value)
	}
	store, err := OpenArtifacts(c.Root, c.Budget.ArtifactBytes, redactor)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	manifest, steps, err := BuildManifest()
	if err != nil {
		return err
	}
	manifestData, err := Encode(manifest)
	if err != nil {
		return err
	}
	now := stamp()
	r := Report{SchemaVersion: Version, Kind: "go-schedule.cli-qa-report", RunID: c.RunID, StartedAt: now, FinishedAt: now, Phase: "capturing", CompletionState: "running", Subject: c.Subject, Runner: Runner{Version: Version, ManifestSHA256: Digest(manifestData), Budget: c.Budget}, Environments: []Environment{c.Environment}}
	corpus := Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: c.RunID}
	r.Cleanup = []Cleanup{{EnvironmentID: c.Environment.ID}}
	r.Environments[0].Version, r.Environments[0].Kernel = platformDetails()
	r.Environments[0].Tools = []NamedValue{{Name: "Go runtime", Value: runtime.Version()}, {Name: "runner", Value: Version}}
	for _, spec := range manifest.Cases {
		status, reason := "blocked", "not_started"
		if !slices.Contains(spec.Platforms, runtime.GOOS) {
			status, reason = "not_applicable", spec.ApplicabilityReason
		}
		r.Cases = append(r.Cases, CaseResult{ID: spec.ID, EnvironmentID: c.Environment.ID, CommandPath: spec.CommandPath, ScenarioClass: spec.ScenarioClass, ExecutionKind: spec.ExecutionKind, Status: status, ExecutionStatus: status, ReasonCode: reason, Argv: spec.Argv, ExpectedExit: spec.ExpectedExit, ExpectedChecks: spec.Checks, StartedAt: now, FinishedAt: now})
	}
	sequence := int64(0)
	event := func(kind string, caseID *string, detail string) error {
		sequence++
		id := c.Environment.ID
		e := Event{Version, "go-schedule.cli-qa-event", c.RunID, sequence, stamp(), kind, &id, caseID, detail}
		data, err := Encode(e)
		if err != nil {
			return err
		}
		var compact []byte
		compact, err = compactJSON(data)
		if err != nil {
			return err
		}
		return store.AppendJournal(c.Prefix+"/events.jsonl", append(compact, '\n'))
	}
	checkpoint := func() error { return persistWorkerCheckpoint(store, c.Prefix, &r) }
	if err := checkpoint(); err != nil {
		return err
	}
	var owned *session
	var limits []NamedValue
	preflightErr := workerPrivilege()
	if preflightErr == nil {
		limits, preflightErr = CurrentLimits(c.RunID, c.Budget)
	}
	binaries := map[string]string{}
	if preflightErr == nil {
		for _, binary := range c.Subject.Binaries {
			if binary.OS != runtime.GOOS {
				continue
			}
			data, e := readBinary(binary.Path)
			if e != nil || Digest(data) != binary.SHA256 {
				preflightErr = errors.Join(errors.New("evaluated executable identity mismatch"), e)
				break
			}
			binaries[binary.Role] = binary.Path
		}
		for _, role := range []string{"cli", "daemon", "runner"} {
			if binaries[role] == "" {
				preflightErr = errors.New("evaluated executable role is missing: " + role)
			}
		}
	}
	if preflightErr == nil {
		if err := CheckInventory(manifest.Commands); err != nil {
			preflightErr = err
		}
	}
	if preflightErr == nil {
		version := captureCommand(ctx, clifixture.Command{Executable: binaries["cli"], Arguments: []string{"--version"}, Directory: filepath.Dir(binaries["cli"]), Environment: quietEnvironment()}, nil, 4096)
		if version.Code != 0 || version.Failure != nil || strings.TrimSpace(string(version.Out)) != c.Subject.Version {
			preflightErr = errors.New("worker CLI version differs from the frozen evaluated subject")
		}
	}
	phase := "json"
	completedCleanup := Cleanup{EnvironmentID: c.Environment.ID}
	startSession := func() error {
		setup, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		var setupErr error
		owned, setupErr = newSession(setup, binaries, c.Budget, redactor, func(w *clifixture.Workspace) error {
			r.Cleanup = []Cleanup{mergeCleanup(completedCleanup, workspaceCleanup(c.Environment.ID, phase, w))}
			if err := retainOwnership(store, &r, c.Prefix+"/"+phase, w); err != nil {
				return err
			}
			return checkpoint()
		})
		return setupErr
	}
	finishSession := func() error {
		if owned == nil {
			return nil
		}
		cleanup, err := closeWorkerSession(store, &r, c.Environment.ID, c.Prefix+"/"+phase, phase, owned)
		completedCleanup = mergeCleanup(completedCleanup, cleanup)
		r.Cleanup = []Cleanup{completedCleanup}
		owned = nil
		return err
	}
	if preflightErr == nil {
		preflightErr = startSession()
	}

	if preflightErr != nil {
		r.Environments[0].Status = "blocked"
		r.Environments[0].ReasonCode = "preflight_failed"
		for i := range r.Cases {
			if r.Cases[i].ExecutionStatus == "blocked" {
				r.Cases[i].ReasonCode = "preflight_failed"
			}
		}
		if err := event("preflight_failed", nil, preflightErr.Error()); err != nil {
			return err
		}
	} else {
		r.Environments[0].Status = "ready"
		r.Environments[0].ReasonCode = "controls_verified"
		r.Environments[0].Limits = limits
		if err := event("environment_ready", nil, "hard limits and exact binaries verified"); err != nil {
			return err
		}
	}
	if owned != nil {
		defer func() {
			if owned != nil {
				resultErr = errors.Join(resultErr, owned.close())
			}
		}()
	}
	for i, definition := range steps {
		result := &r.Cases[i]
		if result.ExecutionStatus == "not_applicable" {
			continue
		}
		if preflightErr != nil {
			if result.ReasonCode == "not_started" {
				result.ReasonCode = "fixture_unavailable"
			}
			continue
		}
		if ctx.Err() != nil {
			result.ReasonCode = "run_deadline"
			continue
		}
		if store.used > store.limit-(8<<20) {
			result.ReasonCode = "artifact_budget_reserve"
			continue
		}
		if definition.Human && phase != "human" {
			if err := finishSession(); err != nil {
				preflightErr = err
				result.ReasonCode = "prior_fixture_cleanup_failed"
				continue
			}
			phase = "human"
			if err := startSession(); err != nil {
				preflightErr = err
				result.ReasonCode = "human_fixture_setup_failed"
				r.Environments[0].Tools = append(r.Environments[0].Tools, NamedValue{Name: "human_fixture_error", Value: err.Error()})
				continue
			}
		}
		result.StartedAt = stamp()
		if err := event("case_started", &result.ID, definition.Spec.ExecutionKind); err != nil {
			return err
		}
		artifactStart := len(r.Artifacts)
		caseCtx, cancel := context.WithTimeout(ctx, time.Duration(c.Budget.CaseTimeoutMS)*time.Millisecond)
		started := time.Now()
		observed := owned.execute(caseCtx, definition)
		cancel()
		result.FinishedAt = stamp()
		result.Observed.DurationMS = time.Since(started).Milliseconds()
		if observed.Blocked != "" {
			result.ReasonCode = observed.Blocked
			artifact, err := store.Write(c.Prefix+"/cases/"+result.ID+"/state.txt", "application/json", observed.State)
			if err != nil {
				return err
			}
			r.Artifacts = append(r.Artifacts, artifact)
			for _, check := range definition.Spec.Checks {
				result.Observed.Checks = append(result.Observed.Checks, Check{Name: check.Name, Expected: check.Expected, Actual: "unavailable", Passed: false, Evidence: []string{artifact.Path}})
			}
			if err := event("case_blocked", &result.ID, observed.Detail); err != nil {
				return err
			}
		} else {
			result.Argv = observed.Args
			result.ExecutionStatus = "pass"
			result.Status = "blocked"
			result.ReasonCode = "pending_review"
			result.Observed.ExitCode = &observed.Code
			if observed.Code != result.ExpectedExit {
				result.ExecutionStatus = "fail"
			}
			prefix := c.Prefix + "/cases/" + result.ID
			channels := []struct {
				name, media string
				data        []byte
			}{{"stdout", "text/plain", observed.Out}, {"stderr", "text/plain", observed.Err}, {"state", "application/json", observed.State}}
			if len(observed.Input) > 0 {
				channels = append(channels, struct {
					name, media string
					data        []byte
				}{"stdin", "text/plain", observed.Input})
			}
			if len(observed.Logs) > 0 {
				channels = append(channels, struct {
					name, media string
					data        []byte
				}{"log", "text/plain", observed.Logs})
			}
			hasMessage := false
			statePath := ""
			for _, channel := range channels {
				artifact, err := store.Write(prefix+"/"+channel.name+".txt", channel.media, channel.data)
				if err != nil {
					return err
				}
				r.Artifacts = append(r.Artifacts, artifact)
				name := artifact.Path
				switch channel.name {
				case "stdout":
					result.StdoutArtifact = &name
				case "stderr":
					result.StderrArtifact = &name
				case "stdin":
					result.StdinArtifact = &name
				case "state":
					statePath = name
				case "log":
					result.LogArtifacts = append(result.LogArtifacts, name)
				}
				if channel.name == "stdout" || channel.name == "stderr" || channel.name == "log" {
					sanitized, err := store.Read(artifact)
					if err != nil {
						return err
					}
					hasMessage = AddMessage(&corpus, *result, channel.name, string(sanitized)) != "" || hasMessage
				}
			}
			if !hasMessage {
				corpus.NoMessages = append(corpus.NoMessages, CaseRef{result.EnvironmentID, result.ID})
			}
			for j, check := range definition.Spec.Checks {
				passed := j < len(observed.Checks) && observed.Checks[j]
				result.Observed.Checks = append(result.Observed.Checks, Check{Name: check.Name, Expected: check.Expected, Actual: strconv.FormatBool(passed), Passed: passed, Evidence: []string{statePath, *result.StdoutArtifact, *result.StderrArtifact}})
				if !passed {
					result.ExecutionStatus = "fail"
				}
			}
			if result.ExecutionStatus == "fail" {
				result.Status = "fail"
				result.ReasonCode = "execution_finding"
			}
			if ctx.Err() != nil {
				result.ExecutionStatus = "blocked"
				result.Status = "blocked"
				result.ReasonCode = "run_interrupted"
			}
			if err := event("case_finished", &result.ID, result.ExecutionStatus); err != nil {
				return err
			}
		}
		if err := persistCase(store, c.Prefix, *result, r.Artifacts[artifactStart:]); err != nil {
			return err
		}
		if i%32 == 31 {
			if err := checkpoint(); err != nil {
				return err
			}
		}
	}
	if err := finishSession(); err != nil {
		if e := event("cleanup_failed", nil, err.Error()); e != nil {
			return e
		}
	}
	if len(r.Cleanup) == 0 {
		r.Cleanup = []Cleanup{completedCleanup}
	}

	finalCtx, cancelFinal := workerFinalizationContext(ctx, 30*time.Second)
	defer cancelFinal()
	store.ctx = finalCtx
	corpus, err = freezeWorkerEvidence(finalCtx, store, c.Prefix, &r)
	if err != nil {
		return err
	}
	r.Phase = "pending_review"
	r.CompletionState = "completed"
	if ctx.Err() != nil || preflightErr != nil && r.Environments[0].Status == "ready" {
		r.CompletionState = "interrupted"
	}
	SortCorpus(&corpus)
	corpusData, err := Encode(corpus)
	if err != nil {
		return err
	}
	r.CorpusSHA256 = Digest(corpusData)
	if err := store.Replace(c.Prefix+"/corpus.json", corpusData); err != nil {
		return err
	}
	if err := event("capture_finished", nil, r.CompletionState); err != nil {
		return err
	}
	return checkpoint()
}

func readBinary(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > 128<<20 {
		return nil, errors.Join(errors.New("binary must be a bounded regular file"), err, file.Close())
	}
	data, err := io.ReadAll(io.LimitReader(file, 128<<20+1))
	if len(data) > 128<<20 {
		err = errors.Join(err, errors.New("binary grew beyond the byte limit"))
	}
	return data, errors.Join(err, file.Close())
}

func compactJSON(data []byte) ([]byte, error) {
	var b bytes.Buffer
	err := json.Compact(&b, data)
	return b.Bytes(), err
}
