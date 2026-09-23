//go:build linux

package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/shruggietech/go-schedule/internal/config"
)

// systemdExec is the public org.freedesktop.systemd1.Service.ExecStart tuple.
// Reading its argv through D-Bus avoids re-parsing systemd's human-readable
// show output or unit-file quoting and observes effective drop-in overrides.
type systemdExec struct {
	Path           string
	Argv           []string
	IgnoreErrors   bool
	StartRealtime  uint64
	StartMonotonic uint64
	StopRealtime   uint64
	StopMonotonic  uint64
	PID            uint32
	Code           int32
	Status         int32
}

// InstalledConfigPath returns the configuration retained by the installed
// systemd service. An uninstalled service still uses the standalone default.
func InstalledConfigPath() (string, error) {
	state, err := QueryState()
	if err != nil {
		return "", fmt.Errorf("read installed service state before resolving configuration: %w", err)
	}
	if state == StateNotInstalled {
		return config.DefaultPath(), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	argv, err := systemdExecArgs(ctx, linuxUnit)
	if err != nil {
		return "", fmt.Errorf("discover installed service configuration: %w", err)
	}
	path, err := configArg(argv)
	if err != nil {
		return "", err
	}
	if path == "" {
		return config.DefaultPath(), nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("read installed service configuration %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("installed service configuration %s is not a regular file", path)
	}
	return path, nil
}

func systemdExecArgs(ctx context.Context, unit string) ([]string, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("connect to system service manager: %w", err)
	}
	defer conn.Close() //nolint:errcheck // bounded read-only D-Bus connection
	manager := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var unitPath dbus.ObjectPath
	if err := manager.CallWithContext(ctx, "org.freedesktop.systemd1.Manager.LoadUnit", 0, unit).Store(&unitPath); err != nil {
		return nil, fmt.Errorf("load %s definition: %w", unit, err)
	}
	var property dbus.Variant
	if err := conn.Object("org.freedesktop.systemd1", unitPath).CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Service", "ExecStart").Store(&property); err != nil {
		return nil, fmt.Errorf("read %s ExecStart: %w", unit, err)
	}
	var commands []systemdExec
	if err := property.Store(&commands); err != nil {
		return nil, fmt.Errorf("decode %s ExecStart: %w", unit, err)
	}
	if len(commands) != 1 || len(commands[0].Argv) == 0 {
		return nil, fmt.Errorf("%s must have exactly one effective ExecStart command", unit)
	}
	return commands[0].Argv, nil
}

func configArg(argv []string) (string, error) {
	path := ""
	for i := 0; i < len(argv); i++ {
		if argv[i] != "--config" {
			continue
		}
		if i+1 >= len(argv) || argv[i+1] == "" || path != "" {
			return "", fmt.Errorf("installed service has an invalid or repeated --config argument")
		}
		path = argv[i+1]
		i++
	}
	if path != "" && !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("installed service configuration path must be absolute")
	}
	return path, nil
}
