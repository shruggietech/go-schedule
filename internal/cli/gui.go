package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

func newGUICmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "gui",
		Short: "Launch the desktop GUI (no console window)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			bin, err := stateFor(cmd).deps.GUIPath()
			if err != nil {
				return err
			}
			if err := stateFor(cmd).deps.LaunchGUI(bin); err != nil {
				return fmt.Errorf("launch gui: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "launched gosched-gui")
			return nil
		},
	}
	MarkLocal(command)
	return command
}

func launchGUI(bin string) error {
	cmd := exec.Command(bin)
	platform.HideConsole(cmd)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
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
		return "", fmt.Errorf("desktop binary not found next to gosched (looked for %s); build it with: cd desktop && go run github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 build", p)
	}
	return p, nil
}
