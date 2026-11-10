package domain

// RunCancellation acknowledges an exact execution request without claiming termination.
// Status is requested or already_finished; Outcome is authoritative only when observed.
type RunCancellation struct {
	RunID              string     `json:"run_id"`
	TaskID             string     `json:"task_id"`
	Status             string     `json:"status"`
	Outcome            RunOutcome `json:"outcome,omitempty"`
	CancellationStatus string     `json:"cancellation_status,omitempty"`
	CancellationError  string     `json:"cancellation_error,omitempty"`
}
