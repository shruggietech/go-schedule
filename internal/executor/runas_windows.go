//go:build windows

package executor

import (
	"fmt"
	"os/exec"
)

// applyRunAs on Windows: per-task user impersonation (LogonUser/CreateProcessAsUser)
// is not supported in this version. An empty runAs runs as the service account
// (no-op); any other value is rejected so behavior is explicit rather than
// silently ignored.
func applyRunAs(_ *exec.Cmd, runAs string) error {
	if runAs == "" {
		return nil
	}
	return fmt.Errorf("per-task run_as is not supported on Windows in this version")
}

// ValidateRunAs reports whether runAs is acceptable on this platform.
func ValidateRunAs(runAs string) error {
	if runAs == "" {
		return nil
	}
	return fmt.Errorf("run_as: not supported on Windows in this version")
}
