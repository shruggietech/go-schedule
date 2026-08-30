// Package ipc provides the local transport the daemon and its clients use: a
// Unix domain socket on Unix and a named pipe on Windows. HTTP/JSON is served
// over this transport, so neither a TCP port nor network exposure is required.
package ipc

import "github.com/shruggietech/go-schedule/internal/config"

// AccessMode identifies the local authorization policy applied to an IPC
// listener.
type AccessMode string

const (
	// AccessModeRestricted admits the configured administrative group.
	AccessModeRestricted AccessMode = "restricted"
	// AccessModeCompatibility preserves broad local access after an operator
	// explicitly configures an empty admin_group.
	AccessModeCompatibility AccessMode = "compatibility"
)

// AccessInfo is startup evidence describing the policy applied to a listener.
type AccessInfo struct {
	Mode       AccessMode
	AdminGroup string
}

// Endpoint resolves the IPC endpoint (socket path or pipe name) from config,
// falling back to the platform default when unset.
func Endpoint(cfg config.Config) string {
	if cfg.IPCPath != "" {
		return cfg.IPCPath
	}
	return defaultEndpoint(cfg)
}
