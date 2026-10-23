package clievaluate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
)

type caseCheckpoint struct {
	Case      CaseResult `json:"case"`
	Artifacts []Artifact `json:"artifacts"`
}

func persistCase(store *ArtifactStore, prefix string, result CaseResult, artifacts []Artifact) error {
	data, err := Encode(caseCheckpoint{result, artifacts})
	if err != nil {
		return err
	}
	_, err = store.Write(prefix+"/progress/"+result.ID+".json", "application/json", data)
	return err
}

func loadWorkerCheckpoint(store *ArtifactStore, prefix string) (Report, error) {
	return loadWorkerCheckpointContext(context.Background(), store, prefix)
}

func loadWorkerCheckpointContext(ctx context.Context, store *ArtifactStore, prefix string) (Report, error) {
	var report Report
	if err := store.ReadDocument(prefix+"/report.json", "report", &report); err != nil {
		return report, err
	}
	if report.Phase != "capturing" {
		return report, nil
	}
	known := map[string]Artifact{}
	for _, a := range report.Artifacts {
		known[a.Path] = a
	}
	for i := range report.Cases {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		path := prefix + "/progress/" + report.Cases[i].ID + ".json"
		data, err := store.readRegular(path, documentLimit)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return report, err
		}
		var record caseCheckpoint
		if err := Decode(bytes.NewReader(data), &record); err != nil {
			return report, err
		}
		if record.Case.ID != report.Cases[i].ID || record.Case.EnvironmentID != report.Cases[i].EnvironmentID {
			return report, errors.New("case checkpoint identity mismatch")
		}
		for _, a := range record.Artifacts {
			if _, err := store.Read(a); err != nil {
				return report, err
			}
			if prior, exists := known[a.Path]; exists {
				if prior != a {
					return report, fmt.Errorf("conflicting case checkpoint artifact: %s", a.Path)
				}
			} else {
				report.Artifacts = append(report.Artifacts, a)
				known[a.Path] = a
			}
		}
		if report.Cases[i].Observed.ExitCode == nil || record.Case.FinishedAt > report.Cases[i].FinishedAt {
			report.Cases[i] = record.Case
		}
		if record.Case.FinishedAt > report.FinishedAt {
			report.FinishedAt = record.Case.FinishedAt
		}
	}
	report.Summary = DeriveSummary(report)
	return report, nil
}

func freezeProgress(store *ArtifactStore, prefix string, report Report) error {
	return freezeProgressContext(context.Background(), store, prefix, report)
}

func freezeProgressContext(ctx context.Context, store *ArtifactStore, prefix string, report Report) error {
	byPath := map[string]Artifact{}
	for _, a := range report.Artifacts {
		byPath[a.Path] = a
	}
	for _, result := range report.Cases {
		if err := ctx.Err(); err != nil {
			return err
		}
		path := prefix + "/progress/" + result.ID + ".json"
		data, err := store.readRegular(path, documentLimit)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		original := data
		var checkpoint caseCheckpoint
		if err := Decode(bytes.NewReader(data), &checkpoint); err != nil {
			return err
		}
		checkpoint.Case = result
		for i, a := range checkpoint.Artifacts {
			current, ok := byPath[a.Path]
			if !ok {
				return fmt.Errorf("checkpoint artifact missing from frozen report: %s", a.Path)
			}
			checkpoint.Artifacts[i] = current
		}
		data, err = Encode(checkpoint)
		if err != nil {
			return err
		}
		data = store.Redactor.Sanitize(data)
		if bytes.Equal(original, data) {
			continue
		}
		if err := store.Replace(path, data); err != nil {
			return err
		}
	}
	return nil
}

func persistWorkerCheckpoint(store *ArtifactStore, prefix string, report *Report) error {
	report.FinishedAt = stamp()
	report.Summary = DeriveSummary(*report)
	data, err := Encode(*report)
	if err != nil {
		return err
	}
	return store.Replace(prefix+"/report.json", data)
}
