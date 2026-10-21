package clievaluate

import (
	"bytes"
	"errors"
	"reflect"
)

// freezeChannels reapplies the complete fixture registry before corpus identity freezes.
func freezeChannels(store *ArtifactStore, r *Report) (Corpus, error) {
	for i := range r.Artifacts {
		a := &r.Artifacts[i]
		data, err := store.Read(*a)
		if err != nil {
			return Corpus{}, err
		}
		data = store.Redactor.Sanitize(data)
		if err := store.Replace(a.Path, data); err != nil {
			return Corpus{}, err
		}
		a.Bytes = int64(len(data))
		a.SHA256 = Digest(data)
	}
	// Keep argv and all derived strings in the same sanitized form as persisted reports.
	data, err := Encode(*r)
	if err != nil {
		return Corpus{}, err
	}
	if err := Decode(bytes.NewReader(store.Redactor.Sanitize(data)), r); err != nil {
		return Corpus{}, err
	}
	corpus := Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: r.RunID}
	byPath := map[string]Artifact{}
	for _, a := range r.Artifacts {
		byPath[a.Path] = a
	}
	for _, c := range r.Cases {
		if c.Observed.ExitCode == nil {
			continue
		}
		hasMessage := false
		channels := []struct {
			name string
			path *string
		}{{"stdout", c.StdoutArtifact}, {"stderr", c.StderrArtifact}}
		for _, p := range c.LogArtifacts {
			p := p
			channels = append(channels, struct {
				name string
				path *string
			}{"log", &p})
		}
		for _, channel := range channels {
			if channel.path == nil {
				return corpus, errors.New("executed case lacks a frozen channel")
			}
			a, ok := byPath[*channel.path]
			if !ok {
				return corpus, errors.New("frozen channel is absent from the evidence index")
			}
			data, err := store.Read(a)
			if err != nil {
				return corpus, err
			}
			hasMessage = AddMessage(&corpus, c, channel.name, string(data)) != "" || hasMessage
		}
		if !hasMessage {
			corpus.NoMessages = append(corpus.NoMessages, CaseRef{c.EnvironmentID, c.ID})
		}
	}
	SortCorpus(&corpus)
	return corpus, nil
}

func validateFrozen(store *ArtifactStore, r Report) error {
	data, err := store.readRegular("evidence-manifest.json", documentLimit)
	if err != nil {
		return err
	}
	if Digest(data) != r.Runner.EvidenceSHA256 {
		return errors.New("frozen evidence manifest digest mismatch")
	}
	var frozen EvidenceManifest
	if err := Decode(bytes.NewReader(data), &frozen); err != nil {
		return err
	}
	if err := ValidateSchema("evidence-manifest", data); err != nil {
		return err
	}
	if frozen.RunID != r.RunID || frozen.CorpusSHA256 != r.CorpusSHA256 || !reflect.DeepEqual(frozen.Artifacts, r.Artifacts) {
		return errors.New("frozen evidence inventory differs from the report")
	}
	expected, err := Encode(frozenMetadata(r))
	if err != nil {
		return err
	}
	for _, artifact := range r.Artifacts {
		if artifact.Path != "capture-identity.json" {
			continue
		}
		data, err := store.Read(artifact)
		if err != nil {
			return err
		}
		if !bytes.Equal(data, expected) {
			return errors.New("source, budget or environment identity differs from the frozen capture")
		}
		return nil
	}
	return errors.New("frozen capture identity is missing")
}
