// Package taskgroup exposes transport-neutral task and group authoring models.
package taskgroup

type EnvironmentRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GroupSummary struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ParentID         string `json:"parentId"`
	Path             string `json:"path"`
	EffectiveReason  string `json:"effectiveReason"`
	UpdatedAt        string `json:"updatedAt"`
	Depth            int    `json:"depth"`
	ChildCount       int    `json:"childCount"`
	TaskCount        int    `json:"taskCount"`
	DescendantCount  int    `json:"descendantCount"`
	DeclaredEnabled  bool   `json:"declaredEnabled"`
	EffectiveEnabled bool   `json:"effectiveEnabled"`
}

type TaskSummary struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	GroupID           string   `json:"groupId"`
	GroupPath         string   `json:"groupPath"`
	EffectiveState    string   `json:"effectiveState"`
	EffectiveReason   string   `json:"effectiveReason"`
	Lifecycle         string   `json:"lifecycle"`
	Timezone          string   `json:"timezone"`
	ScheduleSummary   string   `json:"scheduleSummary"`
	PolicySummary     string   `json:"policySummary"`
	UpdatedAt         string   `json:"updatedAt"`
	CommandConfigured bool     `json:"commandConfigured"`
	DeclaredEnabled   bool     `json:"declaredEnabled"`
	NextRuns          []string `json:"nextRuns"`
}

type Workspace struct {
	Tasks    []TaskSummary  `json:"tasks"`
	Groups   []GroupSummary `json:"groups"`
	LoadedAt string         `json:"loadedAt"`
}

type TaskDetail struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	GroupID           string           `json:"groupId"`
	CommandLine       string           `json:"commandLine"`
	WorkingDir        string           `json:"workingDir"`
	Stdin             string           `json:"stdin"`
	RunAs             string           `json:"runAs"`
	Timezone          string           `json:"timezone"`
	Mode              string           `json:"mode"`
	Schedule          string           `json:"schedule"`
	ScheduleSyntax    string           `json:"scheduleSyntax"`
	At                string           `json:"at"`
	OverlapPolicy     string           `json:"overlapPolicy"`
	CatchupPolicy     string           `json:"catchupPolicy"`
	MissingDatePolicy string           `json:"missingDatePolicy"`
	TimeBasis         string           `json:"timeBasis"`
	DSTGapPolicy      string           `json:"dstGapPolicy"`
	DSTOverlapPolicy  string           `json:"dstOverlapPolicy"`
	ScheduleSummary   string           `json:"scheduleSummary"`
	PolicySummary     string           `json:"policySummary"`
	Readiness         string           `json:"readiness"`
	UpdatedAt         string           `json:"updatedAt"`
	Enabled           bool             `json:"enabled"`
	Environment       []EnvironmentRow `json:"environment"`
	NextRuns          []string         `json:"nextRuns"`
}

type TaskDraft struct {
	TaskDetail
	IsNew             bool   `json:"isNew"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type GroupDraft struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ParentID          string `json:"parentId"`
	OriginalUpdatedAt string `json:"originalUpdatedAt"`
	Enabled           bool   `json:"enabled"`
	IsNew             bool   `json:"isNew"`
	OverwriteStale    bool   `json:"overwriteStale"`
}

type CommandPreview struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
}

type OperationResult struct {
	Action    string          `json:"action"`
	Outcome   string          `json:"outcome"`
	Message   string          `json:"message"`
	Field     string          `json:"field,omitempty"`
	EntityID  string          `json:"entityId,omitempty"`
	Workspace *Workspace      `json:"workspace,omitempty"`
	Task      *TaskDetail     `json:"task,omitempty"`
	Command   *CommandPreview `json:"command,omitempty"`
}
