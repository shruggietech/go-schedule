// Package settings owns desktop-local preferences, storage visibility, product information, and bounded native actions.
package settings

import "time"

// CurrentPreferenceVersion is the only preference document schema accepted by this build.
const CurrentPreferenceVersion = 1

// Appearance is the complete durable desktop color-mode vocabulary.
type Appearance string

const (
	AppearanceSystem Appearance = "system"
	AppearanceLight  Appearance = "light"
	AppearanceDark   Appearance = "dark"
)

// PreferenceTransition records the one-time legacy preference disposition.
type PreferenceTransition struct {
	Status  string   `json:"status"`
	Retired []string `json:"retired"`
}

// DesktopPreferences is the versioned desktop-local preference document.
type DesktopPreferences struct {
	Version    int                  `json:"version"`
	Appearance Appearance           `json:"appearance"`
	Transition PreferenceTransition `json:"transition"`
}

// StorageRecord describes one exact storage source and its removal semantics.
type StorageRecord struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Path          string `json:"path,omitempty"`
	Owner         string `json:"owner"`
	Scope         string `json:"scope"`
	Existence     string `json:"existence"`
	NormalRemoval string `json:"normalRemoval"`
	ExplicitWipe  string `json:"explicitWipe"`
	Copyable      bool   `json:"copyable"`
}

// ProductLink names one fixed HTTPS destination authorized by the backend.
type ProductLink struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Destination string `json:"destination"`
}

// ProductInformation identifies this build and its trusted support links.
type ProductInformation struct {
	Name      string        `json:"name"`
	Version   string        `json:"version"`
	Publisher string        `json:"publisher"`
	Links     []ProductLink `json:"links"`
}

// Workspace is one complete settings and local-information snapshot.
type Workspace struct {
	Preferences     DesktopPreferences `json:"preferences"`
	PreferencePath  string             `json:"preferencePath"`
	Storage         []StorageRecord    `json:"storage"`
	Product         ProductInformation `json:"product"`
	DaemonAvailable bool               `json:"daemonAvailable"`
	LoadedAt        string             `json:"loadedAt"`
}

// Result gives React a stable settings action outcome without native error details.
type Result struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Workspace *Workspace `json:"workspace,omitempty"`
}

func validAppearance(value Appearance) bool {
	return value == AppearanceSystem || value == AppearanceLight || value == AppearanceDark
}

func loadedAt(now time.Time) string { return now.UTC().Format(time.RFC3339Nano) }
