package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeDaemonDisplayNameRejectsControlsBeforeTrimming(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"Workshop\n", "\tWorkshop", "Workshop\r\n"} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := NormalizeDaemonDisplayName(name); !errors.Is(err, ErrInvalidDaemonDisplayName) {
				t.Fatalf("NormalizeDaemonDisplayName(%q) error = %v, want %v", name, err, ErrInvalidDaemonDisplayName)
			}
		})
	}
}

func TestNormalizeDaemonDisplayNameTrimsNonControlSpace(t *testing.T) {
	t.Parallel()

	name, err := NormalizeDaemonDisplayName("\u00a0Workshop\u00a0")
	if err != nil {
		t.Fatalf("NormalizeDaemonDisplayName() error = %v", err)
	}
	if name != "Workshop" {
		t.Fatalf("NormalizeDaemonDisplayName() = %q, want %q", name, "Workshop")
	}
}

func TestNormalizeDaemonDisplayNameRejectsEmptyAndOversizedNames(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"   ", strings.Repeat("x", MaxDaemonDisplayNameRunes+1)} {
		if _, err := NormalizeDaemonDisplayName(name); !errors.Is(err, ErrInvalidDaemonDisplayName) {
			t.Fatalf("NormalizeDaemonDisplayName(%q) error = %v, want %v", name, err, ErrInvalidDaemonDisplayName)
		}
	}
}
