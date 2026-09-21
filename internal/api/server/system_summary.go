package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/store"
)

const systemSummaryWindow = 24 * time.Hour

func (s *Server) handleSystemSummary(w http.ResponseWriter, _ *http.Request) {
	observedAt := time.Now().UTC()
	summary, err := s.store.SystemSummaryFacts(observedAt)
	if err != nil {
		s.internal(w, err)
		return
	}
	summary.NextOccurrence, err = s.nextSummaryOccurrence(observedAt, observedAt.Add(systemSummaryWindow))
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) nextSummaryOccurrence(from, to time.Time) (*domain.UpcomingSummary, error) {
	tasks, err := s.store.ListTasks("", string(domain.TaskActive))
	if err != nil {
		return nil, err
	}
	var nearest *domain.UpcomingSummary
	for _, task := range tasks {
		if !task.Enabled || task.ScheduleID == "" {
			continue
		}
		sch, err := s.store.GetSchedule(task.ScheduleID)
		if errors.Is(err, store.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read schedule %s for system summary: %w", task.ScheduleID, err)
		}
		upcoming, err := schedule.UpcomingRunsWithPolicy(sch, task.Timezone, task.SchedulePolicy(), from, 1)
		if err != nil || len(upcoming) == 0 || upcoming[0].After(to) {
			continue
		}
		candidate := domain.UpcomingSummary{TaskID: task.ID, TaskName: task.Name, ScheduledFor: upcoming[0]}
		if nearest == nil || candidate.ScheduledFor.Before(nearest.ScheduledFor) || (candidate.ScheduledFor.Equal(nearest.ScheduledFor) && candidate.TaskID < nearest.TaskID) {
			nearest = &candidate
		}
	}
	return nearest, nil
}
