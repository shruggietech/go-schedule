package clievaluate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

type workerLease struct {
	PID      int    `json:"pid"`
	Creation string `json:"creation"`
	RunID    string `json:"run_id"`
}

func resumeCleanup(prior Report) (map[string]Cleanup, []Artifact) {
	cleanup := map[string]Cleanup{}
	required := map[string]bool{}
	for _, item := range prior.Cleanup {
		cleanup[item.EnvironmentID] = item
		for _, action := range item.Actions {
			for _, path := range action.Evidence {
				required[path] = true
			}
		}
	}
	artifacts := []Artifact{}
	for _, artifact := range prior.Artifacts {
		if required[artifact.Path] || strings.HasSuffix(artifact.Path, "/ownership.json") {
			artifacts = append(artifacts, artifact)
		}
	}
	return cleanup, artifacts
}

func compareResumeInputs(prior, current captureIdentity) error {
	current.RunID = prior.RunID
	current.Attempt = prior.Attempt
	a, err := Encode(prior)
	if err != nil {
		return err
	}
	b, err := Encode(current)
	if err != nil {
		return err
	}
	if !bytes.Equal(a, b) {
		return errors.New("resume inputs changed: source, binaries, budget, manifest, distro and user must match")
	}
	return nil
}

func confirmPreviousWorkersStopped(ctx context.Context, store *ArtifactStore, prior captureIdentity) error {
	for _, env := range []string{"windows", "wsl"} {
		prefix := fmt.Sprintf("attempt-%03d/%s", prior.Attempt, env)
		data, err := store.readRegular(prefix+"/lease.json", 4096)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		var lease workerLease
		if err := Decode(bytes.NewReader(data), &lease); err != nil {
			return err
		}
		if lease.RunID != prior.RunID || lease.PID <= 0 {
			return errors.New("resume worker lease identity mismatch")
		}
		actual, err := processIdentity(lease.PID)
		if err != nil {
			return err
		}
		if actual != "" && actual == lease.Creation {
			return errors.New("the previous owned launcher is still active; resume after its bounded shutdown finishes")
		}
		if env == "wsl" {
			command := clifixture.Command{Executable: filepath.Join(os.Getenv("WINDIR"), "System32", "wsl.exe"), Arguments: []string{"--distribution", prior.Distro, "--user", prior.User, "--exec", "systemctl", "--user", "show", "gosched-qa-" + prior.RunID + ".service", "--property=LoadState", "--property=ActiveState"}, Directory: filepath.Dir(os.Args[0]), Environment: quietEnvironment()}
			probeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			capture := captureCommand(probeCtx, command, nil, 4096)
			cancel()
			if capture.Failure != nil {
				return capture.Failure
			}
			if err := validateStoppedUnit(capture.Code, strings.TrimSpace(string(capture.Out))); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateStoppedUnit(code int, output string) error {
	fields := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		key, value, ok := strings.Cut(line, "=")
		_, duplicate := fields[key]
		if !ok || duplicate || key != "LoadState" && key != "ActiveState" {
			return errors.New("malformed prior WSL unit state; check the selected distribution and user, then retry resume")
		}
		fields[key] = value
	}
	loadedStopped := code == 0 && fields["LoadState"] == "loaded" && (fields["ActiveState"] == "inactive" || fields["ActiveState"] == "failed")
	collected := (code == 0 || code == 1 || code == 4) && fields["LoadState"] == "not-found" && fields["ActiveState"] == "inactive"
	if len(fields) != 2 || !loadedStopped && !collected {
		return errors.New("could not establish the prior WSL unit's stopped state; check the selected distribution and user, then retry resume")
	}
	return nil
}

func validateResumeEvidence(store *ArtifactStore, prior captureIdentity, manifest Manifest) error {
	if _, err := store.root.Stat("report.json"); err == nil {
		r, _, _, loadErr := loadRun(store)
		if loadErr != nil {
			return loadErr
		}
		if r.RunID == prior.RunID && r.CompletionState != "running" {
			return nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// An interrupted coordinator may leave a valid worker checkpoint before aggregation.
	inspected := false
	for _, env := range []string{"windows", "wsl"} {
		prefix := "attempt-" + fmt.Sprintf("%03d", prior.Attempt) + "/" + env
		r, err := loadWorkerCheckpoint(store, prefix)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		inspected = true
		if r.RunID != prior.RunID || r.Runner.ManifestSHA256 != prior.ManifestSHA256 {
			return errors.New("interrupted checkpoint input identity mismatch")
		}
		for _, artifact := range r.Artifacts {
			if _, err := store.Read(artifact); err != nil {
				return err
			}
		}
		if len(r.Cases) != len(manifest.Cases) {
			return errors.New("interrupted worker ledger is incomplete")
		}
		var corpus Corpus
		corpus.SchemaVersion = Version
		corpus.Kind = "go-schedule.cli-qa-corpus"
		corpus.RunID = r.RunID
		for i := range r.Cases {
			c := &r.Cases[i]
			if c.Observed.ExitCode == nil {
				continue
			}
			message := false
			for _, channel := range []struct {
				name string
				path *string
			}{{"stdout", c.StdoutArtifact}, {"stderr", c.StderrArtifact}} {
				if channel.path == nil {
					return errors.New("checkpoint executed case lacks captured channel")
				}
				for _, artifact := range r.Artifacts {
					if artifact.Path == *channel.path {
						data, err := store.Read(artifact)
						if err != nil {
							return err
						}
						message = AddMessage(&corpus, *c, channel.name, string(data)) != "" || message
					}
				}
			}
			for _, path := range c.LogArtifacts {
				for _, artifact := range r.Artifacts {
					if artifact.Path == path {
						data, err := store.Read(artifact)
						if err != nil {
							return err
						}
						message = AddMessage(&corpus, *c, "log", string(data)) != "" || message
					}
				}
			}
			if !message {
				corpus.NoMessages = append(corpus.NoMessages, CaseRef{c.EnvironmentID, c.ID})
			}
		}
		SortCorpus(&corpus)
		data, err := Encode(corpus)
		if err != nil {
			return err
		}
		r.CorpusSHA256 = Digest(data)
		r.CompletionState = "interrupted"
		r.FinishedAt = stamp()
		r.Summary = DeriveSummary(r)
		if err := ValidateReport(store, r, manifest, corpus); err != nil {
			return err
		}
	}
	if !inspected {
		return errors.New("resume has no durable worker ledger; start a new capture directory")
	}
	return nil
}
