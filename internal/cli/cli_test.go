package cli

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleExecuteError_ReportedUsageDoesNotDuplicateDiagnostic(t *testing.T) {
	var stderr bytes.Buffer
	err := reported(fmtUsage("conversion refused"))
	if got := handleExecuteError(&stderr, err); got != 2 {
		t.Fatalf("exit code = %d, want 2", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("reported error wrote duplicate stderr %q", stderr.String())
	}
	if !errors.Is(err, errUsage) {
		t.Fatal("reported error lost usage classification")
	}
}

func TestHandleExecuteError_UnreportedUsageKeepsExistingDiagnostic(t *testing.T) {
	var stderr bytes.Buffer
	if got := handleExecuteError(&stderr, fmtUsage("bad flag")); got != 2 {
		t.Fatalf("exit code = %d, want 2", got)
	}
	if got := stderr.String(); got != "gosched: usage: bad flag\n" {
		t.Fatalf("stderr = %q", got)
	}
}
