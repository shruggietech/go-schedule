package domain

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	DefaultDaemonDisplayName  = "go-schedule daemon"
	MaxDaemonDisplayNameRunes = 80
)

var ErrInvalidDaemonDisplayName = errors.New("invalid daemon display name")

// DaemonIdentity is the durable identity of one logical daemon installation.
type DaemonIdentity struct {
	InstallationID string
	DisplayName    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NormalizeDaemonDisplayName applies the display contract used by every write path.
func NormalizeDaemonDisplayName(name string) (string, error) {
	if !utf8.ValidString(name) {
		return "", ErrInvalidDaemonDisplayName
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return "", ErrInvalidDaemonDisplayName
		}
	}
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > MaxDaemonDisplayNameRunes {
		return "", ErrInvalidDaemonDisplayName
	}
	return name, nil
}
