//go:build linux

package clievaluate

import (
	"errors"
	"os"
	"strings"
)

func workerPrivilege() error {
	if os.Geteuid() == 0 {
		return errors.New("the unattended QA worker requires an existing nonroot user")
	}
	return nil
}

func platformDetails() (string, *string) {
	version := "Linux"
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				version = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
		}
	}
	kernel := "unavailable"
	if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		kernel = strings.TrimSpace(string(data))
	}
	return version, &kernel
}
