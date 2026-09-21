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
	Permission            string   `json:"permission,omitempty"`
	LastAccessedAt        string   `json:"lastAccessedAt,omitempty"`
	RequestCount          uint64   `json:"requestCount"`
	RequireConfirmation   bool     `json:"requireConfirmation"`
}

// Authority describes one MCP permission class without implying unavailable controls.
type Authority struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// Daemon identifies the target of every projected grant.
type Daemon struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Transport describes one independently controlled MCP path.
type Transport struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Description string `json:"description"`
}

// Grant is the secret-free projection of one MCP actor.
type Grant struct {
	ID                    string `json:"id"`
	ClientName            string `json:"clientName"`
	DaemonID              string `json:"daemonId"`
	DaemonName            string `json:"daemonName"`
	Capability            string `json:"capability"`
	CapabilityDescription string `json:"capabilityDescription"`
	Transport             string `json:"transport"`
	CreatedAt             string `json:"createdAt"`
	LastUsedAt            string `json:"lastUsedAt,omitempty"`
	ExpiresAt             string `json:"expiresAt,omitempty"`
	State                 string `json:"state"`
	CredentialFingerprint string `json:"credentialFingerprint,omitempty"`
}

// Action is one bounded shared audit record.
type Action struct {
	DaemonID   string `json:"daemonId"`
	Operation  string `json:"operation"`
	TargetKind string `json:"targetKind"`
	TargetID   string `json:"targetId,omitempty"`
	Result     string `json:"result"`
	OccurredAt string `json:"occurredAt"`
}

// Workspace is the complete non-secret Agent Access projection.
type Workspace struct {
	StdioDescription string      `json:"stdioDescription"`
	HTTP             HTTPStatus  `json:"http"`
	Authorities      []Authority `json:"authorities"`
	Daemon           Daemon      `json:"daemon"`
	MCPState         string      `json:"mcpState"`
	Transports       []Transport `json:"transports"`
	Grants           []Grant     `json:"grants"`
}

// EnableDraft contains user-selectable localhost settings.
type EnableDraft struct {
	ClientName          string   `json:"clientName"`
	Port                int      `json:"port"`
	AllowedOrigins      []string `json:"allowedOrigins"`
	Permission          string   `json:"permission"`
	RequireConfirmation bool     `json:"requireConfirmation"`
}

// GrantDraft contains the deliberate choices for one remote MCP enrollment.
type GrantDraft struct {
	ClientName string `json:"clientName"`
	Capability string `json:"capability"`
	Duration   string `json:"duration"`
}

// GrantEditDraft contains monotonic changes to an existing grant.
type GrantEditDraft struct {
	ActorID    string `json:"actorId"`
	Capability string `json:"capability,omitempty"`
	Duration   string `json:"duration,omitempty"`
}

// Result gives React a stable result without credentials or native errors.
type Result struct {
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Message   string     `json:"message"`
	Workspace *Workspace `json:"workspace,omitempty"`
}

// ActionsResult returns bounded audit evidence without changing the workspace.
type ActionsResult struct {
	Action  string   `json:"action"`
	Outcome string   `json:"outcome"`
	Message string   `json:"message"`
	Actions []Action `json:"actions"`
}
