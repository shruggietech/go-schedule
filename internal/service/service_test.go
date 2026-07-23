package service

import (
	"testing"

	"github.com/kardianos/service"
)

// The exact wording is part of the CLI's user-visible contract and is asserted
// rather than described: the Windows least-privilege status path added for
// issue #6 must render identically to the library path used elsewhere, so a
// user reading `gosched service status` output cannot tell which one answered.
func TestStatusString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   service.Status
		want string
	}{
		{"running", service.StatusRunning, "running"},
		{"stopped", service.StatusStopped, "stopped"},
		{"unknown", service.StatusUnknown, "unknown (is the service installed?)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := statusString(tc.in); got != tc.want {
				t.Errorf("statusString(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestActionsIncludesEveryControlVerb(t *testing.T) {
	t.Parallel()

	want := map[string]bool{
		"install": false, "uninstall": false, "start": false,
		"stop": false, "restart": false, "status": false,
	}
	for _, a := range Actions() {
		if _, ok := want[a]; !ok {
			t.Errorf("unexpected action %q", a)
			continue
		}
		want[a] = true
	}
	for a, seen := range want {
		if !seen {
			t.Errorf("missing action %q", a)
		}
	}
}
