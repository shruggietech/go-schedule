package server

import (
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

// The publish* helpers broadcast entity-change events to connected GUI clients so
// every view updates in real time (FR-022). They are nil-broker safe.

func (s *Server) publishTaskCreated(t domain.Task) {
	if s.broker != nil {
		s.broker.PublishTask(events.VerbCreated, t.ID, &t)
	}
}

func (s *Server) publishTaskUpdated(id string) {
	if s.broker == nil {
		return
	}
	// Include the current task so clients can upsert without a refetch.
	if t, err := s.store.GetTask(id); err == nil {
		s.broker.PublishTask(events.VerbUpdated, id, &t)
	} else {
		s.broker.PublishTask(events.VerbUpdated, id, nil)
	}
}

func (s *Server) publishTaskDeleted(id string) {
	if s.broker != nil {
		s.broker.PublishTask(events.VerbDeleted, id, nil)
	}
}

func (s *Server) publishGroupCreated(g domain.Group) {
	if s.broker != nil {
		s.broker.PublishGroup(events.VerbCreated, g.ID, &g)
	}
}

func (s *Server) publishGroupUpdated(id string) {
	if s.broker == nil {
		return
	}
	if g, err := s.store.GetGroup(id); err == nil {
		s.broker.PublishGroup(events.VerbUpdated, id, &g)
	} else {
		s.broker.PublishGroup(events.VerbUpdated, id, nil)
	}
}

func (s *Server) publishGroupDeleted(id string) {
	if s.broker != nil {
		s.broker.PublishGroup(events.VerbDeleted, id, nil)
	}
}
