package store

import (
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// SearchFacts returns at most limit plus one safe matches so callers can report truncation.
func (s *Store) SearchFacts(query string, kinds map[domain.SearchKind]bool, limit int) ([]domain.SearchMatch, error) {
	query = strings.TrimSpace(query)
	if query == "" || limit < 1 {
		return []domain.SearchMatch{}, nil
	}
	pattern := "%" + escapeLike(strings.ToLower(query)) + "%"
	queryLimit := limit + 1
	allKinds := len(kinds) == 0
	result := make([]domain.SearchMatch, 0, queryLimit)

	if allKinds || kinds[domain.SearchKindTask] || kinds[domain.SearchKindSchedule] {
		rows, err := s.db.Query(`SELECT id,name,enabled,state,COALESCE(schedule_id,''),CASE WHEN length(trim(command))>0 THEN 1 ELSE 0 END
			FROM tasks WHERE lower(id) LIKE ? ESCAPE '\' OR lower(name) LIKE ? ESCAPE '\'
			ORDER BY lower(name),id LIMIT ?`, pattern, pattern, queryLimit)
		if err != nil {
			return nil, fmt.Errorf("store: search tasks: %w", err)
		}
		for rows.Next() {
			var id, name, state, scheduleID string
			var enabled, commandReady int
			if err := rows.Scan(&id, &name, &enabled, &state, &scheduleID, &commandReady); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: scan search task: %w", err)
			}
			if allKinds || kinds[domain.SearchKindTask] {
				enabledValue := enabled != 0
				actions := []domain.SearchAction{domain.SearchActionOpen}
				if enabledValue {
					actions = append(actions, domain.SearchActionDisable)
				} else {
					actions = append(actions, domain.SearchActionEnable)
				}
				if commandReady != 0 {
					actions = append(actions, domain.SearchActionRunNow)
				}
				status := "disabled"
				if enabledValue {
					status = "enabled"
				}
				result = append(result, domain.SearchMatch{Kind: domain.SearchKindTask, ObjectID: id, TaskID: id, Name: name, Context: state + ", " + status, Enabled: &enabledValue, ActionHints: actions})
			}
			if scheduleID != "" && (allKinds || kinds[domain.SearchKindSchedule]) {
				result = append(result, domain.SearchMatch{Kind: domain.SearchKindSchedule, ObjectID: id, TaskID: id, Name: name, Context: "Scheduled task", ActionHints: []domain.SearchAction{domain.SearchActionOpen}})
			}
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("store: close task search: %w", err)
		}
	}

	if allKinds || kinds[domain.SearchKindGroup] {
		rows, err := s.db.Query(`SELECT id,name,enabled FROM groups
			WHERE lower(id) LIKE ? ESCAPE '\' OR lower(name) LIKE ? ESCAPE '\'
			ORDER BY lower(name),id LIMIT ?`, pattern, pattern, queryLimit)
		if err != nil {
			return nil, fmt.Errorf("store: search groups: %w", err)
		}
		for rows.Next() {
			var id, name string
			var enabled int
			if err := rows.Scan(&id, &name, &enabled); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: scan search group: %w", err)
			}
			status := "disabled"
			if enabled != 0 {
				status = "enabled"
			}
			result = append(result, domain.SearchMatch{Kind: domain.SearchKindGroup, ObjectID: id, Name: name, Context: status, ActionHints: []domain.SearchAction{domain.SearchActionOpen}})
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("store: close group search: %w", err)
		}
	}

	if allKinds || kinds[domain.SearchKindFailure] {
		rows, err := s.db.Query(`SELECT r.id,r.task_id,t.name,COALESCE(r.ended_at,r.scheduled_for)
			FROM runs r JOIN tasks t ON t.id=r.task_id WHERE r.outcome=? AND
			(lower(r.id) LIKE ? ESCAPE '\' OR lower(r.task_id) LIKE ? ESCAPE '\' OR lower(t.name) LIKE ? ESCAPE '\')
			ORDER BY COALESCE(r.ended_at,r.scheduled_for) DESC,r.id DESC LIMIT ?`, string(domain.OutcomeFailure), pattern, pattern, pattern, queryLimit)
		if err != nil {
			return nil, fmt.Errorf("store: search failed runs: %w", err)
		}
		for rows.Next() {
			var id, taskID, name, occurred string
			if err := rows.Scan(&id, &taskID, &name, &occurred); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: scan search failure: %w", err)
			}
			at, err := parseTime(occurred)
			if err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: parse search failure time: %w", err)
			}
			result = append(result, domain.SearchMatch{Kind: domain.SearchKindFailure, ObjectID: id, TaskID: taskID, Name: name, Context: "Failed run", OccurredAt: &at, ActionHints: []domain.SearchAction{domain.SearchActionOpen}})
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("store: close failure search: %w", err)
		}
	}

	if allKinds || kinds[domain.SearchKindAlert] {
		rows, err := s.db.Query(`SELECT a.id,COALESCE(a.task_id,''),COALESCE(t.name,''),a.severity,a.kind,a.created_at
			FROM alerts a LEFT JOIN tasks t ON t.id=a.task_id WHERE a.acknowledged=0 AND
			(lower(a.id) LIKE ? ESCAPE '\' OR lower(COALESCE(a.task_id,'')) LIKE ? ESCAPE '\' OR lower(COALESCE(a.run_id,'')) LIKE ? ESCAPE '\' OR lower(COALESCE(t.name,'')) LIKE ? ESCAPE '\' OR lower(a.severity) LIKE ? ESCAPE '\' OR lower(a.kind) LIKE ? ESCAPE '\')
			ORDER BY a.created_at DESC,a.id DESC LIMIT ?`, pattern, pattern, pattern, pattern, pattern, pattern, queryLimit)
		if err != nil {
			return nil, fmt.Errorf("store: search alerts: %w", err)
		}
		for rows.Next() {
			var id, taskID, name, severity, kind, occurred string
			if err := rows.Scan(&id, &taskID, &name, &severity, &kind, &occurred); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: scan search alert: %w", err)
			}
			at, err := parseTime(occurred)
			if err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("store: parse search alert time: %w", err)
			}
			if name == "" {
				name = "Daemon alert"
			}
			result = append(result, domain.SearchMatch{Kind: domain.SearchKindAlert, ObjectID: id, TaskID: taskID, Name: name, Context: severity + " " + kind, OccurredAt: &at, ActionHints: []domain.SearchAction{domain.SearchActionOpen, domain.SearchActionAcknowledge}})
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("store: close alert search: %w", err)
		}
	}

	domain.SortSearchMatches(result)
	if len(result) > queryLimit {
		result = result[:queryLimit]
	}
	return result, nil
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "%", `\%`)
	return strings.ReplaceAll(value, "_", `\_`)
}
