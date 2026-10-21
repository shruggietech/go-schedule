// Package clievaluate runs owned maintainer evaluations and validates their evidence.
package clievaluate

// Version identifies the compatible artifact contract, independent of product version.
const Version = "1"

// Budget bounds an environment independently of the host's installed scheduler.
type Budget struct {
	CPUPercent       int   `json:"cpu_percent"`
	MemoryBytes      int64 `json:"memory_bytes"`
	WindowsProcesses int   `json:"windows_processes"`
	LinuxTasks       int   `json:"linux_tasks"`
	CaseTimeoutMS    int64 `json:"case_timeout_ms"`
	TotalTimeoutMS   int64 `json:"total_timeout_ms"`
	ChannelBytes     int64 `json:"channel_bytes"`
	ArtifactBytes    int64 `json:"artifact_bytes"`
}

// DefaultBudget is intentionally quiet and leaves reporting/cleanup reserve.
func DefaultBudget() Budget { return Budget{25, 1 << 30, 32, 128, 15000, 1800000, 262144, 64 << 20} }

// Binary binds executable evidence to exact bytes and platform.
type Binary struct {
	Role   string `json:"role"`
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// BuildReceipt binds prepared executable bytes to their build-time source contents.
type BuildReceipt struct {
	SchemaVersion string   `json:"schema_version"`
	Kind          string   `json:"kind"`
	Repository    string   `json:"repository"`
	Commit        string   `json:"commit"`
	Dirty         bool     `json:"dirty"`
	SourceSHA256  string   `json:"source_sha256"`
	Binaries      []Binary `json:"binaries"`
}

// Subject identifies the evaluated source and executable set.
type Subject struct {
	Repository   string   `json:"repository"`
	Commit       string   `json:"commit"`
	Dirty        bool     `json:"dirty"`
	SourceSHA256 string   `json:"source_sha256"`
	Version      string   `json:"version"`
	Binaries     []Binary `json:"binaries"`
}

// Runner identifies the evaluation implementation and frozen inputs.
type Runner struct {
	Version        string       `json:"version"`
	ManifestSHA256 string       `json:"manifest_sha256"`
	EvidenceSHA256 string       `json:"evidence_sha256"`
	Budget         Budget       `json:"budget"`
	Parameters     []NamedValue `json:"parameters"`
}

// NamedValue represents bounded ordered metadata without an open-ended object.
type NamedValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Environment records preflight outcome and actual applied OS controls.
type Environment struct {
	ID         string       `json:"id"`
	OS         string       `json:"os"`
	Version    string       `json:"version"`
	Arch       string       `json:"arch"`
	Distro     *string      `json:"distro"`
	User       *string      `json:"user"`
	Kernel     *string      `json:"kernel"`
	Tools      []NamedValue `json:"tools"`
	Limits     []NamedValue `json:"limits"`
	Status     string       `json:"status"`
	ReasonCode string       `json:"reason_code"`
}

// CheckSpec describes one independently asserted observation.
type CheckSpec struct {
	Name     string `json:"name"`
	Expected string `json:"expected"`
}

// Check is an actual observation with retained evidence references.
type Check struct {
	Name     string   `json:"name"`
	Expected string   `json:"expected"`
	Actual   string   `json:"actual"`
	Passed   bool     `json:"passed"`
	Evidence []string `json:"evidence"`
}

// CaseSpec is a stable manifest scenario, not a generated claim of execution.
type CaseSpec struct {
	ID                  string      `json:"id"`
	CommandPath         string      `json:"command_path"`
	ScenarioClass       string      `json:"scenario_class"`
	ExecutionKind       string      `json:"execution_kind"`
	Platforms           []string    `json:"platforms"`
	ApplicabilityReason string      `json:"applicability_reason"`
	Argv                []string    `json:"argv"`
	ExpectedExit        int         `json:"expected_exit"`
	Checks              []CheckSpec `json:"checks"`
	Coverage            []string    `json:"coverage"`
}

// Flag records the effective command flag, including shadows and inherited scope.
type Flag struct {
	Name       string `json:"name"`
	Shorthand  string `json:"shorthand"`
	Type       string `json:"type"`
	Default    string `json:"default"`
	Usage      string `json:"usage"`
	Scope      string `json:"scope"`
	Hidden     bool   `json:"hidden"`
	Deprecated string `json:"deprecated"`
	Required   bool   `json:"required"`
}

// Command records the actual initialized Cobra command metadata.
type Command struct {
	Path     string   `json:"path"`
	Use      string   `json:"use"`
	Short    string   `json:"short"`
	Long     string   `json:"long"`
	Aliases  []string `json:"aliases"`
	Hidden   bool     `json:"hidden"`
	Runnable bool     `json:"runnable"`
	Flags    []Flag   `json:"flags"`
}

// Manifest freezes metadata and meaningful behavioral cases for both environments.
type Manifest struct {
	SchemaVersion string     `json:"schema_version"`
	Kind          string     `json:"kind"`
	Commands      []Command  `json:"commands"`
	Cases         []CaseSpec `json:"cases"`
}

// Observation preserves raw execution outcome separately from semantic review status.
type Observation struct {
	ExitCode   *int    `json:"exit_code"`
	DurationMS int64   `json:"duration_ms"`
	Checks     []Check `json:"checks"`
}

// CaseResult is present even when prerequisites prevent a planned execution.
type CaseResult struct {
	ID                     string      `json:"id"`
	EnvironmentID          string      `json:"environment_id"`
	CommandPath            string      `json:"command_path"`
	ScenarioClass          string      `json:"scenario_class"`
	ExecutionKind          string      `json:"execution_kind"`
	Status                 string      `json:"status"`
	ExecutionStatus        string      `json:"execution_status"`
	ReasonCode             string      `json:"reason_code"`
	Argv                   []string    `json:"argv"`
	StdinArtifact          *string     `json:"stdin_artifact"`
	ExpectedExit           int         `json:"expected_exit"`
	ExpectedChecks         []CheckSpec `json:"expected_checks"`
	Observed               Observation `json:"observed"`
	StdoutArtifact         *string     `json:"stdout_artifact"`
	StderrArtifact         *string     `json:"stderr_artifact"`
	LogArtifacts           []string    `json:"log_artifacts"`
	SideEffects            []Check     `json:"side_effects"`
	CommunicationReviewIDs []string    `json:"communication_review_ids"`
	FindingIDs             []string    `json:"finding_ids"`
	StartedAt              string      `json:"started_at"`
	FinishedAt             string      `json:"finished_at"`
}

// Artifact describes sanitized regular-file bytes under the run root.
type Artifact struct {
	Path            string `json:"path"`
	MediaType       string `json:"media_type"`
	Bytes           int64  `json:"bytes"`
	SHA256          string `json:"sha256"`
	Redacted        bool   `json:"redacted"`
	RedactionPolicy string `json:"redaction_policy"`
}

// CaseRef identifies an occurrence without ambiguous concatenated identifiers.
type CaseRef struct {
	EnvironmentID string `json:"environment_id"`
	CaseID        string `json:"case_id"`
}

// Message binds exact redacted text to a meaningful scenario/channel context.
type Message struct {
	ID            string    `json:"id"`
	CommandPath   string    `json:"command_path"`
	ScenarioClass string    `json:"scenario_class"`
	Channel       string    `json:"channel"`
	Text          string    `json:"text"`
	Cases         []CaseRef `json:"cases"`
}

// Corpus is immutable after capture and hash-bound to all review input.
type Corpus struct {
	SchemaVersion string    `json:"schema_version"`
	Kind          string    `json:"kind"`
	RunID         string    `json:"run_id"`
	Messages      []Message `json:"messages"`
	NoMessages    []CaseRef `json:"no_messages"`
}

// Rubric records an agent judgment and its explicit rationale.
type Rubric struct {
	Name    string `json:"name"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

// Review contains semantic judgments supplied by the executing agent.
type Review struct {
	ID         string    `json:"id"`
	MessageID  string    `json:"message_id"`
	Cases      []CaseRef `json:"cases"`
	Channel    string    `json:"channel"`
	Excerpt    string    `json:"excerpt"`
	Rubric     []Rubric  `json:"rubric"`
	Verdict    string    `json:"verdict"`
	FindingIDs []string  `json:"finding_ids"`
}

// EvidenceSpan points at verified sanitized evidence, with inclusive line bounds.
type EvidenceSpan struct {
	Artifact  string `json:"artifact"`
	SHA256    string `json:"sha256"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

// Finding provides a directly actionable independently closeable defect.
type Finding struct {
	ID                 string         `json:"id"`
	DeduplicationKey   string         `json:"deduplication_key"`
	Title              string         `json:"title"`
	Severity           string         `json:"severity"`
	Category           string         `json:"category"`
	Target             string         `json:"target"`
	Commit             string         `json:"commit"`
	Version            string         `json:"version"`
	Platforms          []string       `json:"platforms"`
	Cases              []CaseRef      `json:"cases"`
	Preconditions      []string       `json:"preconditions"`
	Argv               []string       `json:"argv"`
	Reproduction       []string       `json:"reproduction"`
	Expected           string         `json:"expected"`
	Actual             string         `json:"actual"`
	Impact             string         `json:"impact"`
	Evidence           []EvidenceSpan `json:"evidence"`
	AcceptanceCriteria []string       `json:"acceptance_criteria"`
	ExistingIssueURL   *string        `json:"existing_issue_url"`
}

// ReviewSet is the agent's complete submission against one frozen corpus.
type ReviewSet struct {
	SchemaVersion string    `json:"schema_version"`
	Kind          string    `json:"kind"`
	RunID         string    `json:"run_id"`
	CorpusSHA256  string    `json:"corpus_sha256"`
	Reviews       []Review  `json:"reviews"`
	Findings      []Finding `json:"findings"`
}

// Cleanup records exact owned resources and any unresolved cleanup requirement.
type Cleanup struct {
	EnvironmentID string       `json:"environment_id"`
	Resources     []NamedValue `json:"resources"`
	Actions       []Check      `json:"actions"`
	Remaining     []string     `json:"remaining"`
	Retained      []string     `json:"retained"`
}

// Summary is derived from cases and reviews and independently checked on input.
type Summary struct {
	Outcome         string `json:"outcome"`
	Planned         int    `json:"planned"`
	Passed          int    `json:"passed"`
	Failed          int    `json:"failed"`
	Blocked         int    `json:"blocked"`
	NotApplicable   int    `json:"not_applicable"`
	Findings        int    `json:"findings"`
	ExecutionPassed int    `json:"execution_passed"`
	ExecutionFailed int    `json:"execution_failed"`
}

// Report is the canonical machine result; Markdown is derived from it.
type Report struct {
	SchemaVersion        string        `json:"schema_version"`
	Kind                 string        `json:"kind"`
	RunID                string        `json:"run_id"`
	StartedAt            string        `json:"started_at"`
	FinishedAt           string        `json:"finished_at"`
	Phase                string        `json:"phase"`
	CompletionState      string        `json:"completion_state"`
	CorpusSHA256         string        `json:"corpus_sha256"`
	Subject              Subject       `json:"subject"`
	Runner               Runner        `json:"runner"`
	Environments         []Environment `json:"environments"`
	Summary              Summary       `json:"summary"`
	Cases                []CaseResult  `json:"cases"`
	CommunicationReviews []Review      `json:"communication_reviews"`
	Findings             []Finding     `json:"findings"`
	Artifacts            []Artifact    `json:"artifacts"`
	Cleanup              []Cleanup     `json:"cleanup"`
}

// Event is an append-only, typed, sequence-numbered sanitized transition.
type Event struct {
	SchemaVersion string  `json:"schema_version"`
	Kind          string  `json:"kind"`
	RunID         string  `json:"run_id"`
	Sequence      int64   `json:"sequence"`
	At            string  `json:"at"`
	Type          string  `json:"type"`
	EnvironmentID *string `json:"environment_id"`
	CaseID        *string `json:"case_id"`
	Detail        string  `json:"detail"`
}

// EvidenceManifest freezes captured content and excludes recursive report hashes.
type EvidenceManifest struct {
	SchemaVersion string     `json:"schema_version"`
	Kind          string     `json:"kind"`
	RunID         string     `json:"run_id"`
	CorpusSHA256  string     `json:"corpus_sha256"`
	Artifacts     []Artifact `json:"artifacts"`
}
