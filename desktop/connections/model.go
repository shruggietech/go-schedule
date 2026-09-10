// Package connections owns the desktop profile lifecycle without exposing bearer credentials.
package connections

type Profile struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	Endpoint         string `json:"endpoint"`
	DaemonID         string `json:"daemonId"`
	ShortDaemonID    string `json:"shortDaemonId"`
	Fingerprint      string `json:"fingerprint"`
	Capability       string `json:"capability"`
	Platform         string `json:"platform"`
	Architecture     string `json:"architecture,omitempty"`
	ProductVersion   string `json:"productVersion,omitempty"`
	LastSuccessfulAt string `json:"lastSuccessfulAt,omitempty"`
	Active           bool   `json:"active"`
}

type Workspace struct {
	ActiveProfileID string    `json:"activeProfileId,omitempty"`
	Profiles        []Profile `json:"profiles"`
}

type Result struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Workspace *Workspace `json:"workspace,omitempty"`
}
