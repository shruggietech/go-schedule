package domain

import "time"

// CompletionDeliveryRole identifies ordinary chains and exact join-owned work.
type CompletionDeliveryRole string

const (
	DeliveryRoleChain        CompletionDeliveryRole = "chain"
	DeliveryRoleJoinBranch   CompletionDeliveryRole = "join_branch"
	DeliveryRoleJoinTarget   CompletionDeliveryRole = "join_target"
	DeliveryRoleJoinRecovery CompletionDeliveryRole = "join_recovery"
	DeliveryRoleJoinCleanup  CompletionDeliveryRole = "join_cleanup"
)

// JoinInstanceState describes the durable decision, not an inferred native result.
type JoinInstanceState string

const (
	JoinInstanceWaiting         JoinInstanceState = "waiting"
	JoinInstanceTargetPending   JoinInstanceState = "target_pending"
	JoinInstanceCompleted       JoinInstanceState = "completed"
	JoinInstanceBlocked         JoinInstanceState = "blocked"
	JoinInstanceResolved        JoinInstanceState = "resolved"
	JoinInstanceRecoveryPending JoinInstanceState = "recovery_pending"
	JoinInstanceCleanupPending  JoinInstanceState = "cleanup_pending"
)

// ParallelJoinBranch names one stable definition slot referencing an existing task.
type ParallelJoinBranch struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	TaskID          string   `json:"task_id"`
	TaskName        string   `json:"task_name,omitempty"`
	TruncatedFields []string `json:"truncated_fields,omitempty"`
}

// ParallelJoin owns explicit branch invocations after each successful entry run.
type ParallelJoin struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	EntryTaskID      string                  `json:"entry_task_id"`
	EntryTaskName    string                  `json:"entry_task_name,omitempty"`
	TargetTaskID     string                  `json:"target_task_id"`
	TargetTaskName   string                  `json:"target_task_name,omitempty"`
	RecoveryTaskID   string                  `json:"recovery_task_id,omitempty"`
	RecoveryTaskName string                  `json:"recovery_task_name,omitempty"`
	RecoveryOn       []JoinRecoveryCondition `json:"recovery_on"`
	CleanupTaskID    string                  `json:"cleanup_task_id,omitempty"`
	CleanupTaskName  string                  `json:"cleanup_task_name,omitempty"`
	Branches         []ParallelJoinBranch    `json:"branches"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	BranchCount      int                     `json:"branch_count"`
	BranchOffset     int                     `json:"branch_offset,omitempty"`
	HasMoreBranches  bool                    `json:"has_more_branches,omitempty"`
	TruncatedFields  []string                `json:"truncated_fields,omitempty"`
}

// ParallelJoinInvocation retains one exact branch contribution, never task latest history.
type ParallelJoinInvocation struct {
	ID                 string         `json:"id"`
	BranchID           string         `json:"branch_id"`
	Name               string         `json:"name"`
	TaskID             string         `json:"task_id"`
	TaskName           string         `json:"task_name,omitempty"`
	DeliveryID         string         `json:"delivery_id"`
	DeliveryState      DeliveryState  `json:"delivery_state"`
	RunID              string         `json:"run_id,omitempty"`
	Outcome            RunOutcome     `json:"outcome,omitempty"`
	Settlement         JoinSettlement `json:"settlement,omitempty"`
	SettledAt          *time.Time     `json:"settled_at,omitempty"`
	Reason             string         `json:"reason,omitempty"`
	StartedAt          *time.Time     `json:"started_at,omitempty"`
	EndedAt            *time.Time     `json:"ended_at,omitempty"`
	CancellationStatus string         `json:"cancellation_status,omitempty"`
	CancellationError  string         `json:"cancellation_error,omitempty"`
	TruncatedFields    []string       `json:"truncated_fields,omitempty"`
}

// ParallelJoinInstance snapshots one entry execution and its durable join decision.
type ParallelJoinInstance struct {
	ID                         string                   `json:"id"`
	DefinitionID               string                   `json:"definition_id"`
	DefinitionName             string                   `json:"definition_name"`
	EntryTaskID                string                   `json:"entry_task_id"`
	EntryTaskName              string                   `json:"entry_task_name,omitempty"`
	EntryRunID                 string                   `json:"entry_run_id"`
	TargetTaskID               string                   `json:"target_task_id"`
	TargetTaskName             string                   `json:"target_task_name,omitempty"`
	RecoveryTaskID             string                   `json:"recovery_task_id,omitempty"`
	RecoveryTaskName           string                   `json:"recovery_task_name,omitempty"`
	RecoveryOn                 []JoinRecoveryCondition  `json:"recovery_on"`
	CleanupTaskID              string                   `json:"cleanup_task_id,omitempty"`
	CleanupTaskName            string                   `json:"cleanup_task_name,omitempty"`
	SelectedRoute              JoinSelectedRoute        `json:"selected_route,omitempty"`
	Branches                   []ParallelJoinInvocation `json:"branches"`
	State                      JoinInstanceState        `json:"state"`
	Reason                     string                   `json:"reason,omitempty"`
	TargetDeliveryID           string                   `json:"target_delivery_id,omitempty"`
	TargetDeliveryState        DeliveryState            `json:"target_delivery_state,omitempty"`
	TargetRunID                string                   `json:"target_run_id,omitempty"`
	TargetOutcome              RunOutcome               `json:"target_outcome,omitempty"`
	TargetStartedAt            *time.Time               `json:"target_started_at,omitempty"`
	TargetEndedAt              *time.Time               `json:"target_ended_at,omitempty"`
	TargetCancellationStatus   string                   `json:"target_cancellation_status,omitempty"`
	TargetCancellationError    string                   `json:"target_cancellation_error,omitempty"`
	TargetSettlement           JoinSettlement           `json:"target_settlement,omitempty"`
	TargetSettledAt            *time.Time               `json:"target_settled_at,omitempty"`
	TargetReason               string                   `json:"target_reason,omitempty"`
	RecoveryDeliveryID         string                   `json:"recovery_delivery_id,omitempty"`
	RecoveryDeliveryState      DeliveryState            `json:"recovery_delivery_state,omitempty"`
	RecoveryRunID              string                   `json:"recovery_run_id,omitempty"`
	RecoveryOutcome            RunOutcome               `json:"recovery_outcome,omitempty"`
	RecoveryStartedAt          *time.Time               `json:"recovery_started_at,omitempty"`
	RecoveryEndedAt            *time.Time               `json:"recovery_ended_at,omitempty"`
	RecoveryCancellationStatus string                   `json:"recovery_cancellation_status,omitempty"`
	RecoveryCancellationError  string                   `json:"recovery_cancellation_error,omitempty"`
	RecoverySettlement         JoinSettlement           `json:"recovery_settlement,omitempty"`
	RecoverySettledAt          *time.Time               `json:"recovery_settled_at,omitempty"`
	RecoveryReason             string                   `json:"recovery_reason,omitempty"`
	CleanupDeliveryID          string                   `json:"cleanup_delivery_id,omitempty"`
	CleanupDeliveryState       DeliveryState            `json:"cleanup_delivery_state,omitempty"`
	CleanupRunID               string                   `json:"cleanup_run_id,omitempty"`
	CleanupOutcome             RunOutcome               `json:"cleanup_outcome,omitempty"`
	CleanupStartedAt           *time.Time               `json:"cleanup_started_at,omitempty"`
	CleanupEndedAt             *time.Time               `json:"cleanup_ended_at,omitempty"`
	CleanupCancellationStatus  string                   `json:"cleanup_cancellation_status,omitempty"`
	CleanupCancellationError   string                   `json:"cleanup_cancellation_error,omitempty"`
	CleanupSettlement          JoinSettlement           `json:"cleanup_settlement,omitempty"`
	CleanupSettledAt           *time.Time               `json:"cleanup_settled_at,omitempty"`
	CleanupReason              string                   `json:"cleanup_reason,omitempty"`
	CreatedAt                  time.Time                `json:"created_at"`
	UpdatedAt                  time.Time                `json:"updated_at"`
	CompletedAt                *time.Time               `json:"completed_at,omitempty"`
	BranchCount                int                      `json:"branch_count"`
	BranchOffset               int                      `json:"branch_offset,omitempty"`
	HasMoreBranches            bool                     `json:"has_more_branches,omitempty"`
	TruncatedFields            []string                 `json:"truncated_fields,omitempty"`
}

// ParallelJoinInstanceFilter bounds authoritative instance and nested-branch reads.
type ParallelJoinInstanceFilter struct {
	DefinitionID string
	EntryRunID   string
	State        JoinInstanceState
	Offset       int
	Limit        int
	TextLimit    int
	BranchOffset int
	BranchLimit  int
}

// JoinRecoveryCondition selects an unsuccessful branch cause for aggregate recovery.
type JoinRecoveryCondition string

const (
	JoinRecoveryOnFailure    JoinRecoveryCondition = "failure"
	JoinRecoveryOnCanceled   JoinRecoveryCondition = "canceled"
	JoinRecoveryOnSkipped    JoinRecoveryCondition = "skipped"
	JoinRecoveryOnNotStarted JoinRecoveryCondition = "not_started"
)

// JoinSettlement records actual terminal or definitive no-start evidence for one owned stage.
type JoinSettlement string

const (
	JoinSettlementPending    JoinSettlement = ""
	JoinSettlementSuccess    JoinSettlement = "success"
	JoinSettlementFailure    JoinSettlement = "failure"
	JoinSettlementCanceled   JoinSettlement = "canceled"
	JoinSettlementSkipped    JoinSettlement = "skipped"
	JoinSettlementNotStarted JoinSettlement = "not_started"
)

// JoinSelectedRoute is the immutable aggregate follow-up decision for an instance.
type JoinSelectedRoute string

const (
	JoinRoutePending  JoinSelectedRoute = ""
	JoinRouteSuccess  JoinSelectedRoute = "success"
	JoinRouteRecovery JoinSelectedRoute = "recovery"
	JoinRouteNone     JoinSelectedRoute = "none"
)
