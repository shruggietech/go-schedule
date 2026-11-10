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
	Name         string                      `json:"name"`
	EntryTaskID  string                      `json:"entry_task_id"`
	Branches     []domain.ParallelJoinBranch `json:"branches"`
	TargetTaskID string                      `json:"target_task_id"`
}

// JoinUpdateRequest changes supplied fields, optionally refusing a stale revision.
type JoinUpdateRequest struct {
	Name              *string                      `json:"name,omitempty"`
	EntryTaskID       *string                      `json:"entry_task_id,omitempty"`
	Branches          *[]domain.ParallelJoinBranch `json:"branches,omitempty"`
	TargetTaskID      *string                      `json:"target_task_id,omitempty"`
	ExpectedUpdatedAt *time.Time                   `json:"expected_updated_at,omitempty"`
}

type joinBranchInput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TaskID string `json:"task_id"`
}

// UnmarshalJSON accepts writable branch fields only while preserving shared Go DTO types.
func (req *JoinCreateRequest) UnmarshalJSON(data []byte) error {
	var wire struct {
		Name         string            `json:"name"`
		EntryTaskID  string            `json:"entry_task_id"`
		TargetTaskID string            `json:"target_task_id"`
		Branches     []joinBranchInput `json:"branches"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	req.Name, req.EntryTaskID, req.TargetTaskID = wire.Name, wire.EntryTaskID, wire.TargetTaskID
	req.Branches = make([]domain.ParallelJoinBranch, len(wire.Branches))
	for i, b := range wire.Branches {
		req.Branches[i] = domain.ParallelJoinBranch{ID: b.ID, Name: b.Name, TaskID: b.TaskID}
	}
	return nil
}

// UnmarshalJSON retains PATCH omission semantics and rejects read-only metadata.
func (req *JoinUpdateRequest) UnmarshalJSON(data []byte) error {
	var wire struct {
		Name              *string            `json:"name"`
		EntryTaskID       *string            `json:"entry_task_id"`
		TargetTaskID      *string            `json:"target_task_id"`
		Branches          *[]joinBranchInput `json:"branches"`
		ExpectedUpdatedAt *time.Time         `json:"expected_updated_at"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	req.Name, req.EntryTaskID, req.TargetTaskID, req.ExpectedUpdatedAt = wire.Name, wire.EntryTaskID, wire.TargetTaskID, wire.ExpectedUpdatedAt
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
	if req.Name == nil && req.EntryTaskID == nil && req.TargetTaskID == nil && req.Branches == nil {
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
	case "", domain.JoinInstanceWaiting, domain.JoinInstanceTargetPending, domain.JoinInstanceCompleted, domain.JoinInstanceBlocked, domain.JoinInstanceResolved:
	default:
		writeError(w, http.StatusBadRequest, CodeValidation, "state", "state must be waiting, target_pending, completed, blocked, or resolved")
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
		if !ok || (branch.DeliveryState != domain.DeliveryClaimed && branch.DeliveryState != domain.DeliveryResolved) || branch.Outcome != "" || run.TaskID != branch.TaskID || run.JoinInstanceID != instance.ID || run.JoinBranchID != branch.ID || run.JoinRole != "branch" {
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
	run, ok := active[instance.TargetDeliveryID]
	if ok && (instance.TargetDeliveryState == domain.DeliveryClaimed || instance.TargetDeliveryState == domain.DeliveryResolved) && instance.TargetOutcome == "" && run.TaskID == instance.TargetTaskID && run.JoinInstanceID == instance.ID && run.JoinRole == "target" {
		instance.TargetRunID, instance.TargetStartedAt, instance.TargetEndedAt, instance.TargetOutcome = run.ID, run.StartedAt, run.EndedAt, run.Outcome
		instance.TargetCancellationStatus = run.CancellationStatus
		var truncated bool
		instance.TargetCancellationError, truncated = boundObservationText(run.CancellationError, textLimit)
		if truncated {
			instance.TruncatedFields = append(instance.TruncatedFields, "target_cancellation_error")
		}
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
