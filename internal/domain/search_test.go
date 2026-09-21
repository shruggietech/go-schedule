package domain

import (
	"slices"
	"testing"
	"time"
)

func TestSortSearchMatchesIsDeterministic(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 9, 21, 18, 0, 0, 0, time.UTC)
	values := []SearchMatch{
		{Kind: SearchKindFailure, ObjectID: "run-b", Name: "Archive", OccurredAt: &now},
		{Kind: SearchKindTask, ObjectID: "task-b", Name: "beta"},
		{Kind: SearchKindTask, ObjectID: "task-a", Name: "Alpha"},
		{Kind: SearchKindFailure, ObjectID: "run-a", Name: "Archive", OccurredAt: &now},
	}

	SortSearchMatches(values)

	got := make([]string, 0, len(values))
	for _, value := range values {
		got = append(got, string(value.Kind)+":"+value.ObjectID)
	}
	want := []string{"task:task-a", "task:task-b", "failure:run-a", "failure:run-b"}
	if !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestSearchKindAndActionValidation(t *testing.T) {
	t.Parallel()
	for _, kind := range []SearchKind{SearchKindTask, SearchKindGroup, SearchKindFailure, SearchKindSchedule, SearchKindAlert} {
		if !kind.Valid() {
			t.Fatalf("expected %q to be valid", kind)
		}
	}
	if SearchKind("output").Valid() {
		t.Fatal("secret-bearing output must not become a search kind")
	}
	for _, action := range []SearchAction{SearchActionOpen, SearchActionAcknowledge, SearchActionEnable, SearchActionDisable, SearchActionRunNow} {
		if !action.Valid() {
			t.Fatalf("expected %q to be valid", action)
		}
	}
}
