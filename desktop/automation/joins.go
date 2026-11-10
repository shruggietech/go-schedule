package automation

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// SaveJoin changes future block instances, using atomic revision validation.
func (s *Service) SaveJoin(ctx context.Context, d JoinDraft) OperationResult {
	if strings.TrimSpace(d.Name) == "" {
		return rejected("save_join", "name", "Enter a parallel block name.")
	}
	if d.EntryTaskID == "" {
		return rejected("save_join", "entryTaskId", "Select an entry task.")
	}
	if d.TargetTaskID == "" {
		return rejected("save_join", "targetTaskId", "Select a target task.")
	}
	if len(d.Branches) == 0 {
		return rejected("save_join", "branches", "Add at least one named branch.")
	}
	if d.RecoveryTaskID != "" && d.RecoveryOn == nil {
		d.RecoveryOn = []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure}
	}
	if d.RecoveryTaskID == "" && len(d.RecoveryOn) > 0 {
		return rejected("save_join", "recoveryOn", "Select a recovery task before choosing recovery conditions.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	id := d.ID
	if d.IsNew {
		request := server.JoinCreateRequest{Name: d.Name, EntryTaskID: d.EntryTaskID, TargetTaskID: d.TargetTaskID, Branches: d.Branches, RecoveryTaskID: d.RecoveryTaskID, CleanupTaskID: d.CleanupTaskID}
		if d.RecoveryTaskID != "" {
			request.RecoveryOn = &d.RecoveryOn
		}
		v, err := s.backend.CreateJoin(c, request)
		if err != nil {
			return failure("save_join", err)
		}
		id = v.ID
	} else {
		current, err := s.backend.GetJoinDefinition(c, id)
		if err != nil {
			return failure("save_join", err)
		}
		r := server.JoinUpdateRequest{}
		if !d.OverwriteStale {
			revision, err := time.Parse(time.RFC3339Nano, d.OriginalUpdatedAt)
			if err != nil {
				return rejected("save_join", "", "Reload the definition before saving.")
			}
			if !current.UpdatedAt.Equal(revision) {
				return staleResult("save_join", id)
			}
			r.ExpectedUpdatedAt = &revision
		}
		if current.Name != d.Name {
			r.Name = &d.Name
		}
		if current.EntryTaskID != d.EntryTaskID {
			r.EntryTaskID = &d.EntryTaskID
		}
		if current.TargetTaskID != d.TargetTaskID {
			r.TargetTaskID = &d.TargetTaskID
		}
		if current.RecoveryTaskID != d.RecoveryTaskID {
			r.RecoveryTaskID = &d.RecoveryTaskID
		}
		if !slices.Equal(current.RecoveryOn, d.RecoveryOn) || (d.RecoveryTaskID != "" && current.RecoveryTaskID == "") {
			r.RecoveryOn = &d.RecoveryOn
		}
		if current.CleanupTaskID != d.CleanupTaskID {
			r.CleanupTaskID = &d.CleanupTaskID
		}
		if !sameJoinBranches(current.Branches, d.Branches) {
			r.Branches = &d.Branches
		}
		if r.Name == nil && r.EntryTaskID == nil && r.TargetTaskID == nil && r.Branches == nil && r.RecoveryTaskID == nil && r.RecoveryOn == nil && r.CleanupTaskID == nil {
			return s.refreshed(ctx, "save_join", "No changes were submitted; the loaded definition already matches.", id)
		}
		if _, err := s.backend.UpdateJoin(c, id, r); err != nil {
			return failure("save_join", err)
		}
	}
	return s.refreshed(ctx, "save_join", "Parallel join saved.", id)
}

// DeleteJoin resolves outstanding decisions without canceling unrelated work.
func (s *Service) DeleteJoin(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "delete_join", id, func(c context.Context) error { return s.backend.DeleteJoin(c, id) })
}

// JoinInstances reads a bounded complete instance page for one definition.
func (s *Service) JoinInstances(ctx context.Context, id string, offset int) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	v, err := s.backend.ListJoinInstances(c, domain.ParallelJoinInstanceFilter{DefinitionID: id, Offset: offset, Limit: 100})
	if err != nil {
		return failure("load_join_instances", err)
	}
	return OperationResult{Action: "load_join_instances", Outcome: "accepted", Message: "Join instances are current.", Instances: v}
}

// JoinInstance retrieves exact historical evidence outside the recent page.
func (s *Service) JoinInstance(ctx context.Context, id string) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	v, err := s.backend.GetJoinInstance(c, id)
	if err != nil {
		return failure("load_join_instance", err)
	}
	return OperationResult{Action: "load_join_instance", Outcome: "accepted", Message: "Opened exact join instance.", Instance: &v}
}

// JoinDefinition loads the complete editable definition through its Manage-only read.
func (s *Service) JoinDefinition(ctx context.Context, id string) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	v, err := s.backend.GetJoinDefinition(c, id)
	if err != nil {
		return failure("load_join_definition", err)
	}
	return OperationResult{Action: "load_join_definition", Outcome: "accepted", Message: "The complete join definition is ready to edit.", Join: &v}
}

func sameJoinBranches(current, desired []domain.ParallelJoinBranch) bool {
	if len(current) != len(desired) {
		return false
	}
	for i, v := range current {
		next := desired[i]
		if v.ID != next.ID || v.Name != next.Name || v.TaskID != next.TaskID {
			return false
		}
	}
	return true
}
