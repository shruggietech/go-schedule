//go:build !windows && !linux

package clievaluate

import "errors"

// CurrentLimits reports unsupported execution controls while contracts remain portable.
func CurrentLimits(string, Budget) ([]NamedValue, error) {
	return nil, errors.New("full CLI QA execution requires Windows 11 or the selected existing Linux/WSL environment")
}
