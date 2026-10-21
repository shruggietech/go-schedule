package clievaluate

import (
	"bytes"
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
		report.Cases[i] = record.Case
		if record.Case.FinishedAt > report.FinishedAt {
			report.FinishedAt = record.Case.FinishedAt
		}
	}
	report.Summary = DeriveSummary(report)
	return report, nil
}

func freezeProgress(store *ArtifactStore, prefix string, report Report) error {
	for _, result := range report.Cases {
		path := prefix + "/progress/" + result.ID + ".json"
		data, err := store.readRegular(path, documentLimit)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if err := store.Replace(path, data); err != nil {
			return err
		}
	}
	return nil
}
