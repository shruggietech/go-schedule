package clievaluate

// captureMetadata binds immutable source and environment claims to frozen evidence.
// Review finalization may change verdicts, but cannot rewrite these observations.
type captureMetadata struct {
	Subject         Subject       `json:"subject"`
	Runner          Runner        `json:"runner"`
	Environments    []Environment `json:"environments"`
	StartedAt       string        `json:"started_at"`
	FinishedAt      string        `json:"finished_at"`
	CompletionState string        `json:"completion_state"`
}

func frozenMetadata(r Report) captureMetadata {
	runner := r.Runner
	runner.EvidenceSHA256 = ""
	return captureMetadata{r.Subject, runner, r.Environments, r.StartedAt, r.FinishedAt, r.CompletionState}
}
