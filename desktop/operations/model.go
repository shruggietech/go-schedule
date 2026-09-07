// Package operations exposes transport-neutral schedule and activity models.
package operations

type ScheduleOccurrence struct {
	ID       string `json:"id"`
	TaskID   string `json:"taskId"`
	TaskName string `json:"taskName"`
	RunID    string `json:"runId,omitempty"`
	Time     string `json:"time"`
	Kind     string `json:"kind"`
	State    string `json:"state"`
	Outcome  string `json:"outcome,omitempty"`
}

type ScheduleSnapshot struct {
	From        string               `json:"from"`
	To          string               `json:"to"`
	LoadedAt    string               `json:"loadedAt"`
	Occurrences []ScheduleOccurrence `json:"occurrences"`
}

type RunRecord struct {
	ID              string `json:"id"`
	TaskID          string `json:"taskId"`
	ScheduledFor    string `json:"scheduledFor"`
	StartedAt       string `json:"startedAt,omitempty"`
	EndedAt         string `json:"endedAt,omitempty"`
	State           string `json:"state"`
	Outcome         string `json:"outcome,omitempty"`
	ExitCode        *int   `json:"exitCode,omitempty"`
	Output          string `json:"output"`
	OutputTruncated bool   `json:"outputTruncated"`
	Trigger         string `json:"trigger"`
	SourceTaskID    string `json:"sourceTaskId,omitempty"`
	SourceRunID     string `json:"sourceRunId,omitempty"`
	SourceTriggerID string `json:"sourceTriggerId,omitempty"`
	SourceWatcherID string `json:"sourceWatcherId,omitempty"`
}

type LogRecord struct {
	ID       string `json:"id"`
	Time     string `json:"time"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
	TaskID   string `json:"taskId,omitempty"`
	RunID    string `json:"runId,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

type AlertRecord struct {
	ID           string `json:"id"`
	TaskID       string `json:"taskId,omitempty"`
	RunID        string `json:"runId,omitempty"`
	Time         string `json:"time"`
	Severity     string `json:"severity"`
	Kind         string `json:"kind"`
	Message      string `json:"message"`
	Acknowledged bool   `json:"acknowledged"`
}

type ActivityWorkspace struct {
	Runs     []RunRecord   `json:"runs"`
	Logs     []LogRecord   `json:"logs"`
	Alerts   []AlertRecord `json:"alerts"`
	LogPath  string        `json:"logPath"`
	LoadedAt string        `json:"loadedAt"`
}

type OperationResult struct {
	Action   string             `json:"action"`
	Outcome  string             `json:"outcome"`
	Message  string             `json:"message"`
	Field    string             `json:"field,omitempty"`
	Schedule *ScheduleSnapshot  `json:"schedule,omitempty"`
	Activity *ActivityWorkspace `json:"activity,omitempty"`
}
