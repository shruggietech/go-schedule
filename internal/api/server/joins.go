package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/events"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

// JoinCreateRequest defines an entry, named invocation slots and joined target.
type JoinCreateRequest struct {
	Name           string                          `json:"name"`
	EntryTaskID    string                          `json:"entry_task_id"`
	Branches       []domain.ParallelJoinBranch     `json:"branches"`
	TargetTaskID   string                          `json:"target_task_id"`
	RecoveryTaskID string                          `json:"recovery_task_id,omitempty"`
	RecoveryOn     *[]domain.JoinRecoveryCondition `json:"recovery_on,omitempty"`
	CleanupTaskID  string                          `json:"cleanup_task_id,omitempty"`
}

// JoinUpdateRequest changes supplied fields, optionally refusing a stale revision.
type JoinUpdateRequest struct {
	Name              *string                         `json:"name,omitempty"`
	EntryTaskID       *string                         `json:"entry_task_id,omitempty"`
	Branches          *[]domain.ParallelJoinBranch    `json:"branches,omitempty"`
	TargetTaskID      *string                         `json:"target_task_id,omitempty"`
	RecoveryTaskID    *string                         `json:"recovery_task_id,omitempty"`
	RecoveryOn        *[]domain.JoinRecoveryCondition `json:"recovery_on,omitempty"`
	CleanupTaskID     *string                         `json:"cleanup_task_id,omitempty"`
	ExpectedUpdatedAt *time.Time                      `json:"expected_updated_at,omitempty"`
}

type joinBranchInput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TaskID string `json:"task_id"`
}

// UnmarshalJSON accepts writable branch fields only while preserving shared Go DTO types.
func (req *JoinCreateRequest) UnmarshalJSON(data []byte) error {
	var wire struct {
		Name           string                          `json:"name"`
		EntryTaskID    string                          `json:"entry_task_id"`
		TargetTaskID   string                          `json:"target_task_id"`
		RecoveryTaskID string                          `json:"recovery_task_id"`
		RecoveryOn     *[]domain.JoinRecoveryCondition `json:"recovery_on"`
		CleanupTaskID  string                          `json:"cleanup_task_id"`
		Branches       []joinBranchInput               `json:"branches"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	req.Name, req.EntryTaskID, req.TargetTaskID = wire.Name, wire.EntryTaskID, wire.TargetTaskID
	req.RecoveryTaskID, req.RecoveryOn, req.CleanupTaskID = wire.RecoveryTaskID, wire.RecoveryOn, wire.CleanupTaskID
	req.Branches = make([]domain.ParallelJoinBranch, len(wire.Branches))
	for i, b := range wire.Branches {
		req.Branches[i] = domain.ParallelJoinBranch{ID: b.ID, Name: b.Name, TaskID: b.TaskID}
	}
	return nil
}

// UnmarshalJSON retains PATCH omission semantics and rejects read-only metadata.
func (req *JoinUpdateRequest) UnmarshalJSON(data []byte) error {
	var wire struct {
		Name              *string                         `json:"name"`
		EntryTaskID       *string                         `json:"entry_task_id"`
		TargetTaskID      *string                         `json:"target_task_id"`
		RecoveryTaskID    *string                         `json:"recovery_task_id"`
		RecoveryOn        *[]domain.JoinRecoveryCondition `json:"recovery_on"`
		CleanupTaskID     *string                         `json:"cleanup_task_id"`
		Branches          *[]joinBranchInput              `json:"branches"`
		ExpectedUpdatedAt *time.Time                      `json:"expected_updated_at"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	req.Name, req.EntryTaskID, req.TargetTaskID, req.ExpectedUpdatedAt = wire.Name, wire.EntryTaskID, wire.TargetTaskID, wire.ExpectedUpdatedAt
	req.RecoveryTaskID, req.RecoveryOn, req.CleanupTaskID = wire.RecoveryTaskID, wire.RecoveryOn, wire.CleanupTaskID
	if wire.Branches != nil {
		branches := make([]domain.ParallelJoinBranch, len(*wire.Branches))
		for i, b := range *wire.Branches {
			branches[i] = domain.ParallelJoinBranch{ID: b.ID, Name: b.Name, TaskID: b.TaskID}
		}
		req.Branches = &branches
	}
	return nil
}

func joinInputBranches(branches []domain.ParallelJoinBranch) []domain.ParallelJoinBranch {
	out := make([]domain.ParallelJoinBranch, len(branches))
	for i, b := range branches {
		out[i] = domain.ParallelJoinBranch{ID: b.ID, Name: b.Name, TaskID: b.TaskID}
	}
	return out
}

func (s *Server) handleCreateJoin(w http.ResponseWriter, r *http.Request) {
	var req JoinCreateRequest
	if !decodeJoinDefinitionRequest(w, r, &req) {
		return
	}
	join := domain.ParallelJoin{Name: req.Name, EntryTaskID: req.EntryTaskID, TargetTaskID: req.TargetTaskID, Branches: joinInputBranches(req.Branches)}
	join.RecoveryTaskID, join.CleanupTaskID = req.RecoveryTaskID, req.CleanupTaskID
	if req.RecoveryOn != nil {
		join.RecoveryOn = append([]domain.JoinRecoveryCondition{}, (*req.RecoveryOn)...)
	}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		s.joinError(w, err)
		return
	}
	s.publishJoin(events.VerbCreated, join)
	writeJSON(w, http.StatusCreated, join)
}

func (s *Server) handleUpdateJoin(w http.ResponseWriter, r *http.Request) {
	var req JoinUpdateRequest
	if !decodeJoinDefinitionRequest(w, r, &req) {
		return
	}
	if req.Name == nil && req.EntryTaskID == nil && req.TargetTaskID == nil && req.Branches == nil && req.RecoveryTaskID == nil && req.RecoveryOn == nil && req.CleanupTaskID == nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "provide at least one join field")
		return
	}
	join, err := s.store.GetParallelJoin(r.PathValue("id"))
	if err != nil {
		s.joinError(w, err)
		return
	}
	if req.Name != nil {
		join.Name = *req.Name
	}
	if req.EntryTaskID != nil {
		join.EntryTaskID = *req.EntryTaskID
	}
	if req.TargetTaskID != nil {
		join.TargetTaskID = *req.TargetTaskID
	}
	if req.RecoveryTaskID != nil {
		previous := join.RecoveryTaskID
		join.RecoveryTaskID = *req.RecoveryTaskID
		if (join.RecoveryTaskID == "" || previous == "") && req.RecoveryOn == nil {
			join.RecoveryOn = nil
		}
	}
	if req.RecoveryOn != nil {
		join.RecoveryOn = append([]domain.JoinRecoveryCondition{}, (*req.RecoveryOn)...)
	}
	if req.CleanupTaskID != nil {
		join.CleanupTaskID = *req.CleanupTaskID
	}
	if req.Branches != nil {
		join.Branches = joinInputBranches(*req.Branches)
		join.BranchCount = len(join.Branches)
		join.BranchOffset = 0
		join.HasMoreBranches = false
		join.TruncatedFields = nil
	}
	if err := s.store.UpdateParallelJoinIfUnchanged(&join, req.ExpectedUpdatedAt); err != nil {
		s.joinError(w, err)
		return
	}
	s.publishJoin(events.VerbUpdated, join)
	writeJSON(w, http.StatusOK, join)
}

func (s *Server) handleDeleteJoin(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.DeleteParallelJoin(id); err != nil {
		s.joinError(w, err)
		return
	}
	if s.broker != nil {
		s.broker.PublishJoin(events.VerbDeleted, id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}

type joinPage struct{ offset, limit, text, branchOffset, branchLimit int }

func joinQueryPage(w http.ResponseWriter, r *http.Request) (joinPage, bool) {
	var p joinPage
	var ok bool
	if p.offset, ok = nonNegativeQuery(w, r, "offset", 0); !ok {
		return p, false
	}
	if p.limit, ok = nonNegativeQuery(w, r, "limit", 100); !ok {
		return p, false
	}
	if p.text, ok = nonNegativeQuery(w, r, "text_limit", 4096); !ok {
		return p, false
	}
	if p.branchOffset, ok = nonNegativeQuery(w, r, "branch_offset", 0); !ok {
		return p, false
	}
	if p.branchLimit, ok = nonNegativeQuery(w, r, "branch_limit", 100); !ok {
		return p, false
	}
	if p.limit < 1 || p.limit > 1000 {
		writeError(w, http.StatusBadRequest, CodeValidation, "limit", "limit must be between 1 and 1000")
		return p, false
	}
	if p.branchLimit < 1 || p.branchLimit > 1000 {
		writeError(w, http.StatusBadRequest, CodeValidation, "branch_limit", "branch_limit must be between 1 and 1000")
		return p, false
	}
	if p.text < 1 || p.text > 16384 {
		writeError(w, http.StatusBadRequest, CodeValidation, "text_limit", "text_limit must be between 1 and 16384")
		return p, false
	}
	return p, true
}

func (s *Server) handleListJoins(w http.ResponseWriter, r *http.Request) {
	p, ok := joinQueryPage(w, r)
	if !ok {
		return
	}
	joins, err := s.store.ListParallelJoinsPageContext(r.Context(), p.offset, p.limit+1, p.text, p.branchOffset, p.branchLimit)
	if err != nil {
		s.joinError(w, err)
		return
	}
	more := len(joins) > p.limit
	if more {
		joins = joins[:p.limit]
	}
	if joins == nil {
		joins = []domain.ParallelJoin{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"joins": joins, "has_more": more, "offset": p.offset, "limit": p.limit})
}

func (s *Server) handleGetJoin(w http.ResponseWriter, r *http.Request) {
	p, ok := joinQueryPage(w, r)
	if !ok {
		return
	}
	join, err := s.store.GetParallelJoinPageContext(r.Context(), r.PathValue("id"), p.branchOffset, p.branchLimit, p.text)
	if err != nil {
		s.joinError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, join)
}

func (s *Server) joinActiveRuns() map[string]domain.Run {
	out := map[string]domain.Run{}
	if source, ok := s.sched.(activeRunSource); ok {
		for _, run := range source.ActiveRuns() {
			if run.DeliveryID != "" {
				out[run.DeliveryID] = run
			}
		}
	}
	return out
}

func (s *Server) handleListJoinInstances(w http.ResponseWriter, r *http.Request) {
	p, ok := joinQueryPage(w, r)
	if !ok {
		return
	}
	state := domain.JoinInstanceState(r.URL.Query().Get("state"))
	switch state {
	case "", domain.JoinInstanceWaiting, domain.JoinInstanceTargetPending, domain.JoinInstanceRecoveryPending, domain.JoinInstanceCleanupPending, domain.JoinInstanceCompleted, domain.JoinInstanceBlocked, domain.JoinInstanceResolved:
	default:
		writeError(w, http.StatusBadRequest, CodeValidation, "state", "state must be waiting, target_pending, recovery_pending, cleanup_pending, completed, blocked, or resolved")
		return
	}
	active := s.joinActiveRuns() // Sample first: a subsequent durable terminal result always wins.
	filter := domain.ParallelJoinInstanceFilter{DefinitionID: r.URL.Query().Get("join_id"), EntryRunID: r.URL.Query().Get("entry_run_id"), State: state, Offset: p.offset, Limit: p.limit + 1, TextLimit: p.text, BranchOffset: p.branchOffset, BranchLimit: p.branchLimit}
	instances, err := s.store.ListParallelJoinInstancesContext(r.Context(), filter)
	if err != nil {
		s.joinError(w, err)
		return
	}
	more := len(instances) > p.limit
	if more {
		instances = instances[:p.limit]
	}
	if instances == nil {
		instances = []domain.ParallelJoinInstance{}
	}
	for i := range instances {
		mergeJoinActive(&instances[i], active, p.text)
	}
	writeJSON(w, http.StatusOK, map[string]any{"instances": instances, "has_more": more, "offset": p.offset, "limit": p.limit})
}

func (s *Server) handleGetJoinInstance(w http.ResponseWriter, r *http.Request) {
	p, ok := joinQueryPage(w, r)
	if !ok {
		return
	}
	active := s.joinActiveRuns()
	instance, err := s.store.GetParallelJoinInstancePageContext(r.Context(), r.PathValue("id"), p.branchOffset, p.branchLimit, p.text)
	if err != nil {
		s.joinError(w, err)
		return
	}
	mergeJoinActive(&instance, active, p.text)
	writeJSON(w, http.StatusOK, instance)
}

func mergeJoinActive(instance *domain.ParallelJoinInstance, active map[string]domain.Run, textLimit int) {
	for i := range instance.Branches {
		branch := &instance.Branches[i]
		run, ok := active[branch.DeliveryID]
		if !ok || (branch.DeliveryState != domain.DeliveryClaimed && branch.DeliveryState != domain.DeliveryResolved) || branch.Outcome != "" || branch.Settlement != domain.JoinSettlementPending || run.TaskID != branch.TaskID || run.JoinInstanceID != instance.ID || run.JoinBranchID != branch.ID || run.JoinRole != "branch" {
			continue
		}
		branch.RunID, branch.StartedAt, branch.EndedAt, branch.Outcome = run.ID, run.StartedAt, run.EndedAt, run.Outcome
		branch.CancellationStatus = run.CancellationStatus
		var truncated bool
		branch.CancellationError, truncated = boundObservationText(run.CancellationError, textLimit)
		if truncated {
			branch.TruncatedFields = append(branch.TruncatedFields, "cancellation_error")
		}
	}
	mergeJoinStageActive(instance, active, textLimit, "target", instance.TargetTaskID, instance.TargetDeliveryID, instance.TargetDeliveryState, instance.TargetSettlement, &instance.TargetRunID, &instance.TargetOutcome, &instance.TargetStartedAt, &instance.TargetEndedAt, &instance.TargetCancellationStatus, &instance.TargetCancellationError)
	mergeJoinStageActive(instance, active, textLimit, "recovery", instance.RecoveryTaskID, instance.RecoveryDeliveryID, instance.RecoveryDeliveryState, instance.RecoverySettlement, &instance.RecoveryRunID, &instance.RecoveryOutcome, &instance.RecoveryStartedAt, &instance.RecoveryEndedAt, &instance.RecoveryCancellationStatus, &instance.RecoveryCancellationError)
	mergeJoinStageActive(instance, active, textLimit, "cleanup", instance.CleanupTaskID, instance.CleanupDeliveryID, instance.CleanupDeliveryState, instance.CleanupSettlement, &instance.CleanupRunID, &instance.CleanupOutcome, &instance.CleanupStartedAt, &instance.CleanupEndedAt, &instance.CleanupCancellationStatus, &instance.CleanupCancellationError)
}

func mergeJoinStageActive(instance *domain.ParallelJoinInstance, active map[string]domain.Run, textLimit int, role, taskID, deliveryID string, state domain.DeliveryState, settlement domain.JoinSettlement, runID *string, outcome *domain.RunOutcome, started, ended **time.Time, cancellationStatus, cancellationError *string) {
	run, ok := active[deliveryID]
	if !ok || (state != domain.DeliveryClaimed && state != domain.DeliveryResolved) || *outcome != "" || settlement != domain.JoinSettlementPending || run.TaskID != taskID || run.JoinInstanceID != instance.ID || run.JoinRole != role {
		return
	}
	*runID, *started, *ended, *outcome = run.ID, run.StartedAt, run.EndedAt, run.Outcome
	*cancellationStatus = run.CancellationStatus
	value, truncated := boundObservationText(run.CancellationError, textLimit)
	*cancellationError = value
	if truncated {
		instance.TruncatedFields = append(instance.TruncatedFields, role+"_cancellation_error")
	}
}

func (s *Server) publishJoin(verb events.Verb, join domain.ParallelJoin) {
	if s.broker != nil {
		s.broker.PublishJoin(verb, join.ID, &join)
	}
}

func (s *Server) joinError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "join", "join, instance or referenced task was not found")
	case errors.Is(err, store.ErrStaleJoin):
		writeError(w, http.StatusConflict, CodeConflict, "expected_updated_at", "join changed; refresh before retrying or explicitly overwrite")
	case errors.Is(err, store.ErrInvalidJoin):
		writeError(w, http.StatusBadRequest, CodeValidation, "join", err.Error())
	case errors.Is(err, store.ErrChainCycle):
		writeError(w, http.StatusBadRequest, CodeValidation, "branches", "relationship would create a cycle with completion chains or parallel joins")
	default:
		s.internal(w, err)
	}
}
