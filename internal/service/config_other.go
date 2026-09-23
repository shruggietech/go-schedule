//go:build !linux

package service

import "github.com/shruggietech/go-schedule/internal/config"

// InstalledConfigPath retains the existing non-Linux default.
func InstalledConfigPath() (string, error) { return config.DefaultPath(), nil }
