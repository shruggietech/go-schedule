// Package automation exposes transport-neutral automation-source administration models.
package automation

type TaskChoice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Readiness string `json:"readiness"`
	Reason    string `json:"reason"`
}

type ChainSummary struct {
	ID             string `json:"id"`
	SourceTaskID   string `json:"sourceTaskId"`
	SourceTaskName string `json:"sourceTaskName"`
	TargetTaskID   string `json:"targetTaskId"`
	TargetTaskName string `json:"targetTaskName"`
	OnOutcome      string `json:"onOutcome"`
	Readiness      string `json:"readiness"`
	Reason         string `json:"reason"`
	UpdatedAt      string `json:"updatedAt"`
}

type TriggerSummary struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TargetTaskID   string `json:"targetTaskId"`
	TargetTaskName string `json:"targetTaskName"`
	Readiness      string `json:"readiness"`
	Reason         string `json:"reason"`
	UpdatedAt      string `json:"updatedAt"`
	Enabled        bool   `json:"enabled"`
}

type TriggerSetMember struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Readiness string `json:"readiness"`
	Reason    string `json:"reason"`
	Position  int    `json:"position"`
	Enabled   bool   `json:"enabled"`
}

type TriggerSetSummary struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	TargetTaskID   string             `json:"targetTaskId"`
	TargetTaskName string             `json:"targetTaskName"`
	Readiness      string             `json:"readiness"`
	Reason         string             `json:"reason"`
	UpdatedAt      string             `json:"updatedAt"`
	MemberCount    int                `json:"memberCount"`
	EnabledCount   int                `json:"enabledCount"`
	Members        []TriggerSetMember `json:"members"`
}

type WatcherSummary struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	Path           string `json:"path"`
	Pattern        string `json:"pattern"`
	Debounce       string `json:"debounce"`
	Stability      string `json:"stability"`
	TargetTaskID   string `json:"targetTaskId"`
	TargetTaskName string `json:"targetTaskName"`
	Health         string `json:"health"`
	HealthReason   string `json:"healthReason"`
	Readiness      string `json:"readiness"`
	Reason         string `json:"reason"`
	UpdatedAt      string `json:"updatedAt"`
	Recursive      bool   `json:"recursive"`
	Enabled        bool   `json:"enabled"`
}

type Workspace struct {
	Tasks       []TaskChoice        `json:"tasks"`
	Chains      []ChainSummary      `json:"chains"`
	Triggers    []TriggerSummary    `json:"triggers"`
	TriggerSets []TriggerSetSummary `json:"triggerSets"`
	Watchers    []WatcherSummary    `json:"watchers"`
	LoadedAt    string              `json:"loadedAt"`
}

type ChainDraft struct {
	ID                string `json:"id"`
	SourceTaskID      string `json:"sourceTaskId"`
	TargetTaskID      string `json:"targetTaskId"`
	OnOutcome         string `json:"onOutcome"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	IsNew             bool   `json:"isNew"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type TriggerDraft struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	TargetTaskID      string `json:"targetTaskId"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	Enabled           bool   `json:"enabled"`
	IsNew             bool   `json:"isNew"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type TriggerSetDraft struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	TargetTaskID      string `json:"targetTaskId"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	Count             int    `json:"count"`
	Enabled           bool   `json:"enabled"`
	IsNew             bool   `json:"isNew"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type WatcherDraft struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Kind              string `json:"kind"`
	Path              string `json:"path"`
	Pattern           string `json:"pattern"`
	Debounce          string `json:"debounce"`
	Stability         string `json:"stability"`
	TargetTaskID      string `json:"targetTaskId"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	Recursive         bool   `json:"recursive"`
	Enabled           bool   `json:"enabled"`
	IsNew             bool   `json:"isNew"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type OperationResult struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Field     string     `json:"field,omitempty"`
	EntityID  string     `json:"entityId,omitempty"`
	Workspace *Workspace `json:"workspace,omitempty"`
}

type Secret struct {
	Label   string `json:"label"`
	Key     string `json:"key"`
	Command string `json:"command"`
}

type SecretResult struct {
	Action   string   `json:"action"`
	Outcome  string   `json:"outcome"`
	Message  string   `json:"message"`
	Field    string   `json:"field,omitempty"`
	EntityID string   `json:"entityId,omitempty"`
	Title    string   `json:"title,omitempty"`
	Secrets  []Secret `json:"secrets,omitempty"`
}
