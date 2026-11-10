package automation

import (
	"context"
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
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	id := d.ID
	if d.IsNew {
		v, err := s.backend.CreateJoin(c, server.JoinCreateRequest{Name: d.Name, EntryTaskID: d.EntryTaskID, TargetTaskID: d.TargetTaskID, Branches: d.Branches})
		if err != nil {
			return failure("save_join", err)
		}
		id = v.ID
	} else {
		r := server.JoinUpdateRequest{Name: &d.Name, EntryTaskID: &d.EntryTaskID, TargetTaskID: &d.TargetTaskID, Branches: &d.Branches}
		if !d.OverwriteStale {
			revision, err := time.Parse(time.RFC3339Nano, d.OriginalUpdatedAt)
			if err != nil {
				return rejected("save_join", "", "Reload the definition before saving.")
			}
			r.ExpectedUpdatedAt = &revision
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
