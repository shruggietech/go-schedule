package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-scheduler/internal/platform"
)

func newGUICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Launch the desktop GUI (no console window)",
		RunE: func(_ *cobra.Command, _ []string) error {
			bin, err := guiPath()
			if err != nil {
				return err
			}
			cmd := exec.Command(bin)
			platform.HideConsole(cmd) // no console window for the GUI process
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("launch gui: %w", err)
			}
			// Detach: don't wait for the GUI to exit.
			_ = cmd.Process.Release()
			fmt.Fprintln(os.Stdout, "launched gosched-gui")
			return nil
		},
	}
}

func guiPath() (string, error) {
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	name := "gosched-gui"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	p := filepath.Join(filepath.Dir(self), name)
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("GUI binary not found next to gosched (looked for %s); build it with: go build ./cmd/gosched-gui", p)
	}
	return p, nil
}
