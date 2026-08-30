package gui

import (
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestChainRowUsesTaskNamesAndPlainOutcome(t *testing.T) {
	row := chainRowText(domain.CompletionChain{
		SourceTaskName: "Fetch", TargetTaskName: "Publish", OnOutcome: domain.CompletionOnSuccess,
	})
	for _, want := range []string{"Fetch", "Publish", "on success"} {
		if !strings.Contains(row, want) {
			t.Fatalf("row %q missing %q", row, want)
		}
	}
}

func TestChainTaskLabelsDisambiguateDuplicateNames(t *testing.T) {
	labels, ids, byID := chainTaskLabels([]domain.Task{{ID: "b", Name: "Build"}, {ID: "a", Name: "Build"}})
	if len(labels) != 2 || labels[0] != "Build (a)" || labels[1] != "Build (b)" {
		t.Fatalf("labels = %+v", labels)
	}
	if ids[labels[0]] != "a" || byID["b"] != "Build (b)" {
		t.Fatalf("identity maps = %+v / %+v", ids, byID)
	}
}

func TestChainOutcomeChoicesRoundTrip(t *testing.T) {
	for _, outcome := range []domain.CompletionOutcome{domain.CompletionOnSuccess, domain.CompletionOnFailure, domain.CompletionOnAny} {
		if got := chainOutcomeForChoice(chainOutcomeChoice(outcome)); got != string(outcome) {
			t.Errorf("%s round trip = %s", outcome, got)
		}
	}
}

func TestUIBuildsChainsTabWithEmptyState(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if len(ui.tabs.Items) < 3 || ui.tabs.Items[2].Text != "Chains" {
		t.Fatalf("chains tab missing: %+v", ui.tabs.Items)
	}
}
