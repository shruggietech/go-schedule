package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/service"
)

func newServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the system-wide background service (install requires admin)",
	}
	for _, action := range service.Actions() {
		cmd.AddCommand(serviceAction(action))
	}
	return cmd
}

func serviceAction(action string) *cobra.Command {
	return &cobra.Command{
		Use:   action,
		Short: "Service: " + action,
		RunE: func(_ *cobra.Command, _ []string) error {
			exec, err := daemonPath()
			if err != nil {
				return err
			}
			msg, err := service.Control(action, exec, nil)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%s: %s\n", action, msg)
			return nil
		},
	}
}

// daemonPath locates the goschedd binary, assumed to live next to gosched.
func daemonPath() (string, error) {
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	name := "goschedd"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	p := filepath.Join(filepath.Dir(self), name)
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("daemon binary not found next to gosched (looked for %s)", p)
	}
	return p, nil
}
