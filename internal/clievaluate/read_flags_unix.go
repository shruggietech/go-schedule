//go:build !windows

package clievaluate

import "syscall"

const regularReadFlags = syscall.O_NONBLOCK
