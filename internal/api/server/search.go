package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

const (
	defaultSearchLimit = 50
	maximumSearchLimit = 50
	maximumSearchQuery = 200
	searchWindow       = 30 * 24 * time.Hour
)

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "q", "enter a search query")
		return
	}
	if !utf8.ValidString(query) || utf8.RuneCountInString(query) > maximumSearchQuery {
		writeError(w, http.StatusBadRequest, CodeValidation, "q", "must contain at most 200 valid UTF-8 characters")
		return
	}
	limit := defaultSearchLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > maximumSearchLimit {
			writeError(w, http.StatusBadRequest, CodeValidation, "limit", "must be between 1 and 50")
			return
		}
		limit = parsed
	}
	kinds := make(map[domain.SearchKind]bool)
	for _, raw := range r.URL.Query()["kind"] {
		kind := domain.SearchKind(strings.TrimSpace(raw))
		if !kind.Valid() {
			writeError(w, http.StatusBadRequest, CodeValidation, "kind", "must be task, group, failure, schedule, or alert")
			return
		}
		kinds[kind] = true
	}

	observedAt := time.Now().UTC()
	searchKinds := kinds
	if kinds[domain.SearchKindSchedule] || len(kinds) == 0 {
		searchKinds = make(map[domain.SearchKind]bool, len(kinds)+4)
		if len(kinds) == 0 {
			for _, kind := range []domain.SearchKind{domain.SearchKindTask, domain.SearchKindGroup, domain.SearchKindFailure, domain.SearchKindAlert} {
				searchKinds[kind] = true
			}
		} else {
			for kind := range kinds {
				searchKinds[kind] = true
			}
		}
		searchKinds[domain.SearchKindSchedule] = false
	}
	results, err := s.store.SearchFacts(query, searchKinds, limit)
	if err != nil {
		s.internal(w, err)
		return
	}
	if kinds[domain.SearchKindSchedule] || len(kinds) == 0 {
		schedules := make([]domain.SearchMatch, 0, limit+1)
		for offset := 0; len(schedules) <= limit; offset += limit + 1 {
			candidates, searchErr := s.store.SearchScheduleFacts(query, offset, limit+1)
			if searchErr != nil {
				s.internal(w, searchErr)
				return
			}
			if len(candidates) == 0 {
				break
			}
			schedules = append(schedules, s.resolveSearchOccurrences(candidates, observedAt)...)
			if len(candidates) < limit+1 {
				break
			}
		}
		results = append(results, schedules...)
	}
	domain.SortSearchMatches(results)
	truncated := len(results) > limit
	if truncated {
		results = results[:limit]
	}
	writeJSON(w, http.StatusOK, domain.DaemonSearch{Schema: domain.DaemonSearchSchema, Query: query, ObservedAt: observedAt, Truncated: truncated, Results: results})
}

func (s *Server) resolveSearchOccurrences(results []domain.SearchMatch, observedAt time.Time) []domain.SearchMatch {
	resolved := results[:0]
	for _, match := range results {
		if match.Kind != domain.SearchKindSchedule {
			resolved = append(resolved, match)
			continue
		}
		task, err := s.store.GetTask(match.TaskID)
		if err != nil || !task.Enabled || task.State != domain.TaskActive {
			continue
		}
		enabled, err := s.store.GroupChainEnabled(task.GroupID)
		if err != nil || !enabled {
			continue
		}
		stored, err := s.store.GetSchedule(task.ScheduleID)
		if err != nil {
			continue
		}
		upcoming, err := schedule.UpcomingRunsWithPolicy(stored, task.Timezone, task.SchedulePolicy(), observedAt, 1)
		if err != nil || len(upcoming) == 0 || upcoming[0].After(observedAt.Add(searchWindow)) {
			continue
		}
		at := upcoming[0].UTC()
		match.ObjectID = fmt.Sprintf("%s@%s", task.ID, at.Format(time.RFC3339Nano))
		match.OccurredAt = &at
		match.Context = "Next scheduled occurrence"
		resolved = append(resolved, match)
	}
	return resolved
}
