//go:build !windows

package clientsecret

import (
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/shruggietech/go-schedule/internal/processgroup"
)

func configureReadCommand(command *exec.Cmd) func() error {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var once sync.Once
	var cleanupErr error
	cleanup := func() error {
		once.Do(func() {
			if command.Process != nil {
				cleanupErr = processgroup.Terminate(command.Process.Pid, 5*time.Second)
			}
		})
		return cleanupErr
	}
	// macOS native keyring access uses a security child. Own and drain that group too.
	command.Cancel = cleanup
	return cleanup
}
