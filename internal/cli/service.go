package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/service"
)

type serviceControlFunc func(action, executable string, arguments []string) (string, error)

func newServiceCmd() *cobra.Command {
	return newServiceCmdWith(service.Control, daemonPath)
}

func newServiceCmdWith(control serviceControlFunc, resolveDaemon func() (string, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the system-wide background service (install requires admin)",
	}
	for _, action := range service.Actions() {
		cmd.AddCommand(serviceAction(action, control, resolveDaemon))
	}
	return cmd
}

func serviceAction(action string, control serviceControlFunc, resolveDaemon func() (string, error)) *cobra.Command {
	var configPath string
	command := &cobra.Command{
		Use:   action,
		Short: "Service: " + action,
		RunE: func(_ *cobra.Command, _ []string) error {
			arguments, err := serviceArguments(action, configPath)
			if err != nil {
				return err
			}
			exec, err := resolveDaemon()
			if err != nil {
				return err
			}
			msg, err := control(action, exec, arguments)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%s: %s\n", action, msg)
			return nil
		},
	}
	if action == "install" {
		command.Flags().StringVar(&configPath, "config", "", "absolute or relative daemon configuration file to retain in the service definition")
	}
	return command
}

func serviceArguments(action, configPath string) ([]string, error) {
	if action != "install" || configPath == "" {
		return nil, nil
	}
	abs, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("resolve service configuration: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("read service configuration %s: %w", abs, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("read service configuration %s: path is not a regular file", abs)
	}
	if _, err := config.LoadRequired(abs); err != nil {
		return nil, fmt.Errorf("validate service configuration %s: %w", abs, err)
	}
	return []string{"--config", abs}, nil
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
