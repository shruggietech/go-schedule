//go:build linux

package desktopcontrol

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var errElevationCancelled = errors.New("authorization cancelled")

const linuxServiceUnit = "goschedd.service"

func linuxControlCommand(action string, root bool) (string, []string, error) {
	if action != "start" && action != "stop" && action != "restart" {
		return "", nil, fmt.Errorf("unsupported service action %q", action)
	}
	if root {
		return "/usr/bin/systemctl", []string{action, linuxServiceUnit}, nil
	}
	return "/usr/bin/pkexec", []string{"/usr/bin/systemctl", action, linuxServiceUnit}, nil
}

func requestElevation(action string) error {
	path, args, err := linuxControlCommand(action, os.Geteuid() == 0)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, args...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return fmt.Errorf("%w after 35 seconds", errHelperTimedOut)
	}
	if err == nil {
		return nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && path == "/usr/bin/pkexec" && exit.ExitCode() == 126 {
		return errElevationCancelled
	}
	if errors.Is(err, os.ErrNotExist) && path == "/usr/bin/pkexec" {
		return fmt.Errorf("graphical authorization is unavailable: install polkit and a session authentication agent: %w", err)
	}
	if errors.As(err, &exit) && path == "/usr/bin/pkexec" && exit.ExitCode() == 127 {
		return fmt.Errorf("authorization was denied or no graphical authentication agent is available: %w", err)
	}
	message := strings.TrimSpace(string(output))
	if message == "" {
		return fmt.Errorf("service manager rejected %s: %w", action, err)
	}
	return fmt.Errorf("service manager rejected %s (%s): %w", action, message, err)
}
