package clievaluate

import "slices"

func persistCapture(store *ArtifactStore, r *Report, m Manifest, c Corpus, events []byte) error {
	r.FinishedAt = stamp()
	SortCorpus(&c)
	manifest, err := Encode(m)
	if err != nil {
		return err
	}
	corpus, err := Encode(c)
	if err != nil {
		return err
	}
	identity, err := Encode(frozenMetadata(*r))
	if err != nil {
		return err
	}
	r.CorpusSHA256 = Digest(corpus)
	r.Runner.ManifestSHA256 = Digest(manifest)
	r.Artifacts = slices.DeleteFunc(r.Artifacts, func(a Artifact) bool {
		return slices.Contains([]string{"manifest.json", "corpus.json", "events.jsonl", "capture-identity.json"}, a.Path)
	})
	for _, document := range []struct {
		name, media string
		data        []byte
	}{{"manifest.json", "application/json", manifest}, {"corpus.json", "application/json", corpus}, {"events.jsonl", "application/x-ndjson", events}, {"capture-identity.json", "application/json", identity}} {
		if err := store.Replace(document.name, document.data); err != nil {
			return err
		}
		r.Artifacts = append(r.Artifacts, Artifact{Path: document.name, MediaType: document.media, Bytes: int64(len(document.data)), SHA256: Digest(document.data), Redacted: true, RedactionPolicy: "fixture-secrets-v1"})
	}
	frozen := EvidenceManifest{Version, "go-schedule.cli-qa-evidence-manifest", r.RunID, r.CorpusSHA256, r.Artifacts}
	data, err := Encode(frozen)
	if err != nil {
		return err
	}
	r.Runner.EvidenceSHA256 = Digest(data)
	if err := store.Replace("evidence-manifest.json", data); err != nil {
		return err
	}
	r.Summary = DeriveSummary(*r)
	if err := ValidateReport(store, *r, m, c); err != nil {
		return err
	}
	return saveReport(store, *r)
}
