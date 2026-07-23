//go:build windows

package service

import (
	"testing"

	"github.com/kardianos/service"
)

// The defect behind issue #6 was that a status query asked for start and stop
// rights it did not need, so OpenService failed with "Access is denied" before
// any status was read. The observable consequence is that "the service is
// absent" and "we asked for too much access" arrive as the same message.
//
// Querying a name that cannot exist separates them: a least-privilege query
// reaches OpenService and gets ERROR_SERVICE_DOES_NOT_EXIST, which maps to the
// not-installed result with no error. An access failure here — or on any
// machine where this test runs unelevated — means the wide mask is back.
func TestPlatformStatusMissingServiceIsNotAnAccessError(t *testing.T) {
	t.Parallel()

	st, handled, err := platformStatus("goschedd-nonexistent-test-service")
	if err != nil {
		t.Fatalf("querying a missing service must not error (a least-privilege "+
			"open reaches ERROR_SERVICE_DOES_NOT_EXIST): %v", err)
	}
	if !handled {
		t.Fatal("Windows must handle the status query itself, not fall through")
	}
	if st != service.StatusUnknown {
		t.Errorf("status = %v, want StatusUnknown for a missing service", st)
	}
	if got := statusString(st); got != "unknown (is the service installed?)" {
		t.Errorf("wording = %q, want the existing not-installed wording", got)
	}
}
