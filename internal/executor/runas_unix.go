//go:build !windows

package executor

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// applyRunAs configures cmd to run as the named user (by name or numeric UID).
// An empty runAs runs as the daemon's own account (no-op).
func applyRunAs(cmd *exec.Cmd, runAs string) error {
	if runAs == "" {
		return nil
	}
	u, err := lookupUser(runAs)
	if err != nil {
		return fmt.Errorf("unknown user %q: %w", runAs, err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("invalid uid for %q: %w", runAs, err)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return fmt.Errorf("invalid gid for %q: %w", runAs, err)
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	return nil
}

func lookupUser(s string) (*user.User, error) {
	if u, err := user.Lookup(s); err == nil {
		return u, nil
	}
	return user.LookupId(s)
}

// ValidateRunAs reports whether runAs is acceptable on this platform.
func ValidateRunAs(runAs string) error {
	if runAs == "" {
		return nil
	}
	if _, err := lookupUser(runAs); err != nil {
		return fmt.Errorf("run_as: unknown user %q", runAs)
	}
	return nil
}
