//go:build !windows && !linux

package clievaluate

import (
	"errors"
	"runtime"
)

func workerPrivilege() error {
	return errors.New("unattended execution requires Windows or the selected WSL2 Linux environment")
}

func platformDetails() (string, *string) { return runtime.GOOS, nil }
