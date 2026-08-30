//go:build windows

package ipc

import (
	"context"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/internal/config"
)

func defaultEndpoint(_ config.Config) string {
	return `\\.\pipe\goschedd`
}

const compatibilityPipeSDDL = "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;AU)"

type windowsIPCOperations struct {
	lookupGroup func(string) (string, uint32, error)
	listenPipe  func(string, string) (net.Listener, error)
}

var productionWindowsIPCOperations = windowsIPCOperations{
	lookupGroup: func(name string) (string, uint32, error) {
		sid, _, accountType, err := windows.LookupSID("", name)
		if err != nil {
			return "", 0, err
		}
		return sid.String(), accountType, nil
	},
	listenPipe: func(endpoint, descriptor string) (net.Listener, error) {
		return winio.ListenPipe(endpoint, &winio.PipeConfig{SecurityDescriptor: descriptor})
	},
}

// Listen creates a Windows named-pipe listener with an ACL that lets the
// interactive user reach a service-hosted (LocalSystem) daemon.
func Listen(cfg config.Config) (net.Listener, AccessInfo, error) {
	return listenWindows(cfg, productionWindowsIPCOperations)
}

func listenWindows(cfg config.Config, ops windowsIPCOperations) (net.Listener, AccessInfo, error) {
	descriptor := compatibilityPipeSDDL
	access := AccessInfo{Mode: AccessModeCompatibility}
	if cfg.AdminGroup != "" {
		sid, accountType, err := ops.lookupGroup(cfg.AdminGroup)
		if err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: resolve admin_group %q: %w", cfg.AdminGroup, err)
		}
		if sid == "" {
			return nil, AccessInfo{}, fmt.Errorf("ipc: admin_group %q resolved to an empty SID", cfg.AdminGroup)
		}
		if accountType != windows.SidTypeGroup && accountType != windows.SidTypeAlias && accountType != windows.SidTypeWellKnownGroup {
			return nil, AccessInfo{}, fmt.Errorf("ipc: admin_group %q resolves to a non-group account", cfg.AdminGroup)
		}
		descriptor = "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;" + sid + ")"
		access = AccessInfo{Mode: AccessModeRestricted, AdminGroup: cfg.AdminGroup}
	}

	endpoint := Endpoint(cfg)
	l, err := ops.listenPipe(endpoint, descriptor)
	if err != nil {
		return nil, AccessInfo{}, fmt.Errorf("ipc: listen pipe %s: %w", endpoint, err)
	}
	return l, access, nil
}

// DialContext connects to the named-pipe endpoint.
func DialContext(ctx context.Context, endpoint string) (net.Conn, error) {
	return winio.DialPipeContext(ctx, endpoint)
}
