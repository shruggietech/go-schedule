//go:build windows

package clientsecret

import (
	"os"
	"os/exec"

	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

func configureReadCommand(command *exec.Cmd) func() error {
	platform.HideConsole(command)
	// Windows native credential reads do not launch descendants. CommandContext
	// terminates the exact helper handle and Run reaps it before returning.
	return func() error { return nil }
}

func expireReadHelper() { os.Exit(1) }
