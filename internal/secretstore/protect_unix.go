//go:build !windows

// Package secretstore protects daemon-owned secret values with platform facilities.
package secretstore

// Protect returns the value for storage in the daemon-only Unix database.
func Protect(value string) (string, error) { return value, nil }

// Unprotect returns a value read from the daemon-only Unix database.
func Unprotect(value string) (string, error) { return value, nil }
