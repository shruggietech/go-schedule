package domain

import (
	"sort"
	"strings"
	"time"
)

// DaemonSearchSchema identifies the additive daemon search response contract.
const DaemonSearchSchema = "go-schedule.daemon-search.v1"

// SearchKind identifies one safe searchable projection.
type SearchKind string

const (
	SearchKindTask     SearchKind = "task"
	SearchKindGroup    SearchKind = "group"
	SearchKindFailure  SearchKind = "failure"
	SearchKindSchedule SearchKind = "schedule"
	SearchKindAlert    SearchKind = "alert"
)

// Valid reports whether the kind belongs to the public search contract.
func (k SearchKind) Valid() bool {
	switch k {
	case SearchKindTask, SearchKindGroup, SearchKindFailure, SearchKindSchedule, SearchKindAlert:
		return true
	default:
		return false
	}
}

// SearchAction identifies an action a search result may advertise.
type SearchAction string

const (
	SearchActionOpen        SearchAction = "open"
	SearchActionAcknowledge SearchAction = "acknowledge"
	SearchActionEnable      SearchAction = "enable"
	SearchActionDisable     SearchAction = "disable"
	SearchActionRunNow      SearchAction = "run_now"
)

// Valid reports whether the action belongs to the target-safe search workflow.
func (a SearchAction) Valid() bool {
	switch a {
	case SearchActionOpen, SearchActionAcknowledge, SearchActionEnable, SearchActionDisable, SearchActionRunNow:
		return true
	default:
		return false
	}
}

// SearchMatch is a bounded, secret-free projection of one searchable object.
type SearchMatch struct {
	Kind        SearchKind     `json:"kind"`
	ObjectID    string         `json:"object_id"`
	TaskID      string         `json:"task_id,omitempty"`
	Name        string         `json:"name"`
	Context     string         `json:"context,omitempty"`
	OccurredAt  *time.Time     `json:"occurred_at,omitempty"`
	Enabled     *bool          `json:"enabled,omitempty"`
	ActionHints []SearchAction `json:"action_hints"`
}

// DaemonSearch is one authoritative bounded search observation.
type DaemonSearch struct {
	Schema     string        `json:"schema"`
	Query      string        `json:"query"`
	ObservedAt time.Time     `json:"observed_at"`
	Truncated  bool          `json:"truncated"`
	Results    []SearchMatch `json:"results"`
}

var searchKindOrder = map[SearchKind]int{
	SearchKindTask: 0, SearchKindGroup: 1, SearchKindFailure: 2, SearchKindSchedule: 3, SearchKindAlert: 4,
}

// SortSearchMatches applies the stable public ordering used by every transport.
func SortSearchMatches(values []SearchMatch) {
	sort.SliceStable(values, func(i, j int) bool {
		left, right := values[i], values[j]
		if searchKindOrder[left.Kind] != searchKindOrder[right.Kind] {
			return searchKindOrder[left.Kind] < searchKindOrder[right.Kind]
		}
		leftName, rightName := strings.ToLower(left.Name), strings.ToLower(right.Name)
		if leftName != rightName {
			return leftName < rightName
		}
		if left.OccurredAt != nil && right.OccurredAt != nil && !left.OccurredAt.Equal(*right.OccurredAt) {
			return left.OccurredAt.After(*right.OccurredAt)
		}
		return left.ObjectID < right.ObjectID
	})
}
