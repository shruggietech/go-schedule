//go:build !windows

package integration

import "syscall"

func s143ProcessExited(pid int) bool { return syscall.Kill(pid, 0) == syscall.ESRCH }
