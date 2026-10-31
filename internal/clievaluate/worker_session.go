package clievaluate

import (
	"errors"
	"strconv"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func workspaceCleanup(environment, phase string, w *clifixture.Workspace) Cleanup {
	return Cleanup{EnvironmentID: environment, Resources: []NamedValue{{Name: phase + ".workspace", Value: w.Root}, {Name: phase + ".daemon_pid", Value: strconv.Itoa(w.PID)}, {Name: phase + ".daemon_id", Value: w.DaemonID}, {Name: phase + ".ipc", Value: w.IPCPath}}, Remaining: []string{w.Runtime}, Retained: []string{w.Evidence}}
}

func mergeCleanup(first, second Cleanup) Cleanup {
	result := Cleanup{EnvironmentID: first.EnvironmentID}
	if result.EnvironmentID == "" {
		result.EnvironmentID = second.EnvironmentID
	}
	result.Resources = append(append([]NamedValue{}, first.Resources...), second.Resources...)
	result.Actions = append(append([]Check{}, first.Actions...), second.Actions...)
	result.Remaining = append(append([]string{}, first.Remaining...), second.Remaining...)
	result.Retained = append(append([]string{}, first.Retained...), second.Retained...)
	return result
}

func closeWorkerSession(store *ArtifactStore, r *Report, environment, prefix, phase string, owned *session) (Cleanup, error) {
	cleanup := workspaceCleanup(environment, phase, owned.w)
	if err := retainOwnership(store, r, prefix, owned.w); err != nil {
		return cleanup, err
	}
	cleanupErr := owned.close()
	for _, capture := range []*memoryCapture{owned.daemonOut, owned.daemonErr} {
		if capture == nil {
			continue
		}
		data, overflow := capture.Snapshot()
		if len(data) == 0 {
			continue
		}
		name := prefix + "/daemon-" + strconv.Itoa(len(r.Artifacts)) + ".log"
		artifact, err := store.Write(name, "text/plain", data)
		if err != nil {
			return cleanup, errors.Join(cleanupErr, err)
		}
		r.Artifacts = append(r.Artifacts, artifact)
		for i := len(r.Cases) - 1; i >= 0; i-- {
			c := &r.Cases[i]
			if c.Observed.ExitCode == nil {
				continue
			}
			c.LogArtifacts = append(c.LogArtifacts, name)
			if overflow {
				c.SideEffects = append(c.SideEffects, Check{Name: "daemon log capture remained within its channel limit", Expected: "true", Actual: "false", Passed: false, Evidence: []string{name}})
				c.ExecutionStatus = "fail"
				c.Status = "fail"
				c.ReasonCode = "log_budget_exceeded"
			}
			break
		}
	}
	cleanup.Actions = []Check{{Name: phase + " owned process tree drained and runtime removed", Expected: "true", Actual: strconv.FormatBool(cleanupErr == nil), Passed: cleanupErr == nil, Evidence: []string{prefix + "/ownership.json"}}}
	if cleanupErr == nil {
		cleanup.Remaining = []string{}
	}
	return cleanup, cleanupErr
}
