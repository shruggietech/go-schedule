// Package notifications exposes secret-free desktop notification models.
package notifications

type Channel struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Kind                  string `json:"kind"`
	EndpointSummary       string `json:"endpointSummary"`
	HasAuthorization      bool   `json:"hasAuthorization"`
	Enabled               bool   `json:"enabled"`
	HealthIntervalMinutes int64  `json:"healthIntervalMinutes"`
	UpdatedAt             string `json:"updatedAt"`
}

type ChannelDraft struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Endpoint              string `json:"endpoint"`
	Authorization         string `json:"authorization"`
	ReplaceEndpoint       bool   `json:"replaceEndpoint"`
	ReplaceAuthorization  bool   `json:"replaceAuthorization"`
	IsNew                 bool   `json:"isNew"`
	HealthIntervalMinutes int64  `json:"healthIntervalMinutes"`
}

type Scope struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Context string `json:"context"`
}

type Assignment struct {
	ChannelID                string `json:"channelId"`
	OnSuccess                bool   `json:"onSuccess"`
	OnFailure                bool   `json:"onFailure"`
	FailureThreshold         int    `json:"failureThreshold"`
	OnFailureToStart         bool   `json:"onFailureToStart"`
	DurationThresholdMinutes int64  `json:"durationThresholdMinutes"`
	OnRecovery               bool   `json:"onRecovery"`
	ReminderIntervalMinutes  int64  `json:"reminderIntervalMinutes"`
	QuietPeriodMinutes       int64  `json:"quietPeriodMinutes"`
}

type PolicyDraft struct {
	ScopeType   string       `json:"scopeType"`
	ScopeID     string       `json:"scopeId"`
	Assignments []Assignment `json:"assignments"`
}

type Policy struct {
	Scope                Scope        `json:"scope"`
	DirectAssignments    []Assignment `json:"directAssignments"`
	EffectiveSourceType  string       `json:"effectiveSourceType"`
	EffectiveSourceID    string       `json:"effectiveSourceId,omitempty"`
	EffectiveSourceName  string       `json:"effectiveSourceName"`
	EffectiveAssignments []Assignment `json:"effectiveAssignments"`
}

type ConfiguredScope struct {
	Type                    string `json:"type"`
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Context                 string `json:"context"`
	SourceType              string `json:"sourceType"`
	SourceName              string `json:"sourceName"`
	OnSuccess               bool   `json:"onSuccess"`
	OnFailure               bool   `json:"onFailure"`
	OnProblem               bool   `json:"onProblem"`
	DestinationCount        int    `json:"destinationCount"`
	EnabledDestinationCount int    `json:"enabledDestinationCount"`
	EnabledSuccessCount     int    `json:"enabledSuccessDestinationCount"`
	EnabledFailureCount     int    `json:"enabledFailureDestinationCount"`
	EnabledProblemCount     int    `json:"enabledProblemDestinationCount"`
}

type Delivery struct {
	ID                 string `json:"id"`
	ChannelID          string `json:"channelId,omitempty"`
	ChannelName        string `json:"channelName"`
	DestinationSummary string `json:"destinationSummary"`
	Kind               string `json:"kind"`
	TaskID             string `json:"taskId,omitempty"`
	TaskName           string `json:"taskName,omitempty"`
	GroupID            string `json:"groupId,omitempty"`
	GroupName          string `json:"groupName,omitempty"`
	RunID              string `json:"runId,omitempty"`
	State              string `json:"state"`
	Attempts           int    `json:"attempts"`
	NextAttemptAt      string `json:"nextAttemptAt,omitempty"`
	CreatedAt          string `json:"createdAt"`
	ClaimedAt          string `json:"claimedAt,omitempty"`
	CompletedAt        string `json:"completedAt,omitempty"`
	LastStatus         int    `json:"lastStatus,omitempty"`
	LastError          string `json:"lastError,omitempty"`
	ConditionKind      string `json:"conditionKind,omitempty"`
	ConditionSummary   string `json:"conditionSummary,omitempty"`
}

type Workspace struct {
	Channels         []Channel         `json:"channels"`
	Tasks            []Scope           `json:"tasks"`
	Groups           []Scope           `json:"groups"`
	Coverage         []ConfiguredScope `json:"coverage"`
	CoverageComplete bool              `json:"coverageComplete"`
	Deliveries       []Delivery        `json:"deliveries"`
	LoadedAt         string            `json:"loadedAt"`
}

type Result struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Field     string     `json:"field,omitempty"`
	EntityID  string     `json:"entityId,omitempty"`
	Workspace *Workspace `json:"workspace,omitempty"`
	Policy    *Policy    `json:"policy,omitempty"`
}
