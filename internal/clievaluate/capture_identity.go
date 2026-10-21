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
	Cleanup         []Cleanup     `json:"cleanup"`
	Cases           []CaseResult  `json:"cases"`
}

func frozenMetadata(r Report) captureMetadata {
	runner := r.Runner
	runner.EvidenceSHA256 = ""
	cases := append([]CaseResult{}, r.Cases...)
	for i := range cases {
		cases[i].Status = cases[i].ExecutionStatus
		if cases[i].ExecutionStatus == "pass" {
			cases[i].ReasonCode = "executed"
		}
		cases[i].CommunicationReviewIDs = nil
		cases[i].FindingIDs = nil
	}
	return captureMetadata{r.Subject, runner, r.Environments, r.StartedAt, r.FinishedAt, r.CompletionState, r.Cleanup, cases}
}
