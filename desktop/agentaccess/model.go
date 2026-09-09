// Package agentaccess owns the desktop projection and lifecycle of local MCP access.
package agentaccess

// HTTPStatus is the non-secret state of the optional localhost transport.
type HTTPStatus struct {
	Enabled               bool     `json:"enabled"`
	Endpoint              string   `json:"endpoint,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins"`
	CredentialFingerprint string   `json:"credentialFingerprint,omitempty"`
	EnabledAt             string   `json:"enabledAt,omitempty"`
	ClientName            string   `json:"clientName,omitempty"`
	LastAccessedAt        string   `json:"lastAccessedAt,omitempty"`
	RequestCount          uint64   `json:"requestCount"`
}

// Authority describes one MCP permission class without implying unavailable controls.
type Authority struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// Workspace is the complete non-secret Agent Access projection.
type Workspace struct {
	StdioDescription string      `json:"stdioDescription"`
	HTTP             HTTPStatus  `json:"http"`
	Authorities      []Authority `json:"authorities"`
}

// EnableDraft contains user-selectable localhost settings.
type EnableDraft struct {
	ClientName     string   `json:"clientName"`
	Port           int      `json:"port"`
	AllowedOrigins []string `json:"allowedOrigins"`
}

// Result gives React a stable result without credentials or native errors.
type Result struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Workspace *Workspace `json:"workspace,omitempty"`
}
