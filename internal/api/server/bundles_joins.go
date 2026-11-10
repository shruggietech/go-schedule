package server

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/events"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

type bundleRelationship struct {
	item    bundle.Item
	id      string
	created bool
}

// applyBundleRelationships writes the reviewed final mixed graph atomically,
// including creations that could otherwise form a transient cycle during swaps.
func (s *Server) applyBundleRelationships(doc bundle.Document, plan bundle.Plan) []bundle.Item {
	chains := []domain.CompletionChain{}
	joins := []domain.ParallelJoin{}
	objects := []bundleRelationship{}
	var prepareErr error
	for _, source := range doc.Chains {
		if !bundlePlanChanges(plan, "chain", source.PortableID) {
			continue
		}
		object, err := s.bundleRelationshipIdentity("chain", source.PortableID, source.PortableID)
		objects = append(objects, object)
		if err != nil {
			prepareErr = err
			continue
		}
		from, err := s.store.ObjectIDForPortableID("task", source.SourceTaskID)
		if err != nil {
			prepareErr = fmt.Errorf("source task was not applied: %w", err)
			continue
		}
		to, err := s.store.ObjectIDForPortableID("task", source.TargetTaskID)
		if err != nil {
			prepareErr = fmt.Errorf("target task was not applied: %w", err)
			continue
		}
		chains = append(chains, domain.CompletionChain{ID: object.id, SourceTaskID: from, TargetTaskID: to, OnOutcome: domain.CompletionOutcome(source.OnOutcome)})
	}
	for _, source := range doc.Joins {
		if !bundlePlanChanges(plan, "join", source.PortableID) {
			continue
		}
		object, err := s.bundleRelationshipIdentity("join", source.PortableID, source.Name)
		objects = append(objects, object)
		if err != nil {
			prepareErr = err
			continue
		}
		join, err := s.bundleJoin(source, object.id)
		if err != nil {
			prepareErr = err
			continue
		}
		joins = append(joins, join)
	}
	if prepareErr == nil && len(objects) > 0 {
		prepareErr = s.store.ReplaceAutomationRelationships(chains, joins)
	}
	outcomes := make([]bundle.Item, 0, len(objects))
	for _, object := range objects {
		item := object.item
		if prepareErr != nil {
			item.Action, item.Message = bundle.ActionFailed, prepareErr.Error()
		} else {
			item.Action = bundle.ActionApplied
			if object.created {
				if err := s.store.BindPortableID(item.Kind, object.id, item.PortableID); err != nil {
					item.Action, item.Message = bundle.ActionUncertain, err.Error()
				}
			}
			if item.Action == bundle.ActionApplied {
				s.publishBundleRelationship(object)
			}
		}
		outcomes = append(outcomes, item)
	}
	return outcomes
}

func (s *Server) bundleRelationshipIdentity(kind, portableID, name string) (bundleRelationship, error) {
	object := bundleRelationship{item: bundle.Item{Kind: kind, PortableID: portableID, Name: name}}
	id, err := s.store.ObjectIDForPortableID(kind, portableID)
	if errors.Is(err, store.ErrNotFound) {
		object.id, object.created = uuid.NewString(), true
		return object, nil
	}
	object.id = id
	return object, err
}

func (s *Server) bundleJoin(source bundle.Join, id string) (domain.ParallelJoin, error) {
	join := domain.ParallelJoin{ID: id, Name: source.Name, Branches: []domain.ParallelJoinBranch{}}
	entry, err := s.store.ObjectIDForPortableID("task", source.EntryTaskID)
	if err != nil {
		return join, fmt.Errorf("join entry task was not applied: %w", err)
	}
	target, err := s.store.ObjectIDForPortableID("task", source.TargetTaskID)
	if err != nil {
		return join, fmt.Errorf("join target task was not applied: %w", err)
	}
	join.EntryTaskID, join.TargetTaskID = entry, target
	if source.RecoveryTaskID != "" {
		join.RecoveryTaskID, err = s.store.ObjectIDForPortableID("task", source.RecoveryTaskID)
		if err != nil {
			return join, fmt.Errorf("join recovery task was not applied: %w", err)
		}
	}
	if source.CleanupTaskID != "" {
		join.CleanupTaskID, err = s.store.ObjectIDForPortableID("task", source.CleanupTaskID)
		if err != nil {
			return join, fmt.Errorf("join cleanup task was not applied: %w", err)
		}
	}
	if source.RecoveryOn != nil {
		join.RecoveryOn = make([]domain.JoinRecoveryCondition, 0, len(*source.RecoveryOn))
		for _, cause := range *source.RecoveryOn {
			join.RecoveryOn = append(join.RecoveryOn, domain.JoinRecoveryCondition(cause))
		}
	}
	for _, branch := range source.Branches {
		taskID, err := s.store.ObjectIDForPortableID("task", branch.TaskID)
		if err != nil {
			return join, fmt.Errorf("join branch %s task was not applied: %w", branch.Name, err)
		}
		join.Branches = append(join.Branches, domain.ParallelJoinBranch{ID: branch.ID, Name: branch.Name, TaskID: taskID})
	}
	return join, nil
}

func (s *Server) publishBundleRelationship(object bundleRelationship) {
	verb := events.VerbUpdated
	if object.created {
		verb = events.VerbCreated
	}
	if object.item.Kind == "join" {
		if join, err := s.store.GetParallelJoin(object.id); err == nil {
			s.publishJoin(verb, join)
		}
	} else if chain, err := s.store.GetCompletionChain(object.id); err == nil {
		s.publishChain(verb, chain)
	}
}
