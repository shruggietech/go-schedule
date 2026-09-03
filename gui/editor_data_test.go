package gui

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestPolicyLabelRoundTrip(t *testing.T) {
	for _, v := range []domain.OverlapPolicy{domain.OverlapQueueOne, domain.OverlapSkip, domain.OverlapAllowConcurrent} {
		if got := overlapValue(overlapLabel(v)); got != v {
			t.Fatalf("overlap round-trip: %q -> %q -> %q", v, overlapLabel(v), got)
		}
	}
	for _, v := range []domain.CatchupPolicy{domain.CatchupOne, domain.CatchupNone} {
		if got := catchupValue(catchupLabel(v)); got != v {
			t.Fatalf("catchup round-trip: %q -> %q -> %q", v, catchupLabel(v), got)
		}
	}
}

func TestPolicyLabelUnknownFallsBackToDefault(t *testing.T) {
	if got := overlapValue("not a real label"); got != domain.OverlapQueueOne {
		t.Fatalf("unknown overlap label = %q, want default %q", got, domain.OverlapQueueOne)
	}
	if got := catchupValue(""); got != domain.CatchupOne {
		t.Fatalf("unknown catchup label = %q, want default %q", got, domain.CatchupOne)
	}
	if got := overlapLabel(domain.OverlapPolicy("legacy")); got != overlapChoices[0].label {
		t.Fatalf("unknown overlap value label = %q, want default", got)
	}
}

func TestCommandPreviewTextExposesExactBoundaries(t *testing.T) {
	tests := []struct {
		command string
		args    []string
		want    string
	}{
		{"cmd", nil, "Program\n\"cmd\"\nArguments in order (0)\nNone"},
		{"cmd", []string{"/c", "echo hello world"}, "Program\n\"cmd\"\nArguments in order (2)\n1  \"/c\"\n2  \"echo hello world\""},
		{"program", []string{"", "a\tb", "a\r\nb", `C:\x`, `say "hi"`}, "Program\n\"program\"\nArguments in order (5)\n1  \"\"\n2  \"a\\tb\"\n3  \"a\\r\\nb\"\n4  \"C:\\\\x\"\n5  \"say \\\"hi\\\"\""},
	}
	for _, tt := range tests {
		if got := commandPreviewText(tt.command, tt.args); got != tt.want {
			t.Fatalf("commandPreviewText(%q, %v) = %q, want %q", tt.command, tt.args, got, tt.want)
		}
	}
}
