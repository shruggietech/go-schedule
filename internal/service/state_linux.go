//go:build linux

package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const linuxUnit = "goschedd.service"

func platformQueryState(string) (State, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "/usr/bin/systemctl", "show", "--no-pager", "--property=LoadState,ActiveState,SubState", linuxUnit).Output()
	if err != nil {
		return StateUnknown, fmt.Errorf("query %s: %w", linuxUnit, err)
	}
	return parseSystemdState(string(output))
}

func parseSystemdState(output string) (State, error) {
	properties := make(map[string]string, 3)
	for _, line := range strings.Split(output, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			properties[key] = value
		}
	}
	switch properties["LoadState"] {
	case "not-found":
		return StateNotInstalled, nil
	case "loaded":
	default:
		return StateUnknown, fmt.Errorf("unrecognized %s load state %q", linuxUnit, properties["LoadState"])
	}
	switch properties["ActiveState"] {
	case "active":
		return StateRunning, nil
	case "inactive", "failed":
		return StateStopped, nil
	case "activating", "reloading":
		return StateStarting, nil
	case "deactivating":
		return StateStopping, nil
	default:
		return StateUnknown, fmt.Errorf("unrecognized %s active state %q", linuxUnit, properties["ActiveState"])
	}
}
