//go:build windows

package ipc

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/internal/config"
)

func defaultEndpoint(_ config.Config) string {
	return `\\.\pipe\goschedd`
}

const compatibilityPipeSDDL = "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;AU)"

type windowsIPCOperations struct {
	lookupGroup      func(string) (string, uint32, error)
	listGroupMembers func(string) ([]windowsGroupMember, error)
	listenPipe       func(string, string) (net.Listener, error)
}

type windowsGroupMember struct {
	sid         string
	accountType uint32
}

type localGroupMembersInfo1 struct {
	sid         *windows.SID
	accountType uint32
	name        *uint16
}

const (
	maxPreferredLength  = ^uint32(0)
	maxGroupMemberPages = 1024
)

var (
	netapi32                    = windows.NewLazySystemDLL("netapi32.dll")
	procNetLocalGroupGetMembers = netapi32.NewProc("NetLocalGroupGetMembers")
	procNetApiBufferFree        = netapi32.NewProc("NetApiBufferFree")
)

var productionWindowsIPCOperations = windowsIPCOperations{
	lookupGroup: func(name string) (string, uint32, error) {
		sid, _, accountType, err := windows.LookupSID("", name)
		if err != nil {
			return "", 0, err
		}
		return sid.String(), accountType, nil
	},
	listGroupMembers: listLocalGroupMembers,
	listenPipe: func(endpoint, descriptor string) (net.Listener, error) {
		return winio.ListenPipe(endpoint, &winio.PipeConfig{SecurityDescriptor: descriptor})
	},
}

func listLocalGroupMembers(name string) ([]windowsGroupMember, error) {
	groupName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf("encode group name: %w", err)
	}

	var members []windowsGroupMember
	var resume uintptr
	for page := 0; page < maxGroupMemberPages; page++ {
		var buffer *localGroupMembersInfo1
		var entriesRead uint32
		var totalEntries uint32
		status, _, _ := procNetLocalGroupGetMembers.Call(
			0,
			uintptr(unsafe.Pointer(groupName)),
			1,
			uintptr(unsafe.Pointer(&buffer)),
			uintptr(maxPreferredLength),
			uintptr(unsafe.Pointer(&entriesRead)),
			uintptr(unsafe.Pointer(&totalEntries)),
			uintptr(unsafe.Pointer(&resume)),
		)

		if buffer != nil {
			entries := unsafe.Slice(buffer, int(entriesRead))
			for _, entry := range entries {
				member := windowsGroupMember{accountType: entry.accountType}
				if entry.sid != nil && entry.sid.IsValid() {
					member.sid = entry.sid.String()
				}
				members = append(members, member)
			}
			freeStatus, _, _ := procNetApiBufferFree.Call(uintptr(unsafe.Pointer(buffer)))
			if freeStatus != 0 {
				return nil, fmt.Errorf("free group member buffer: %w", syscall.Errno(freeStatus))
			}
		}

		switch status {
		case 0:
			return members, nil
		case uintptr(windows.ERROR_MORE_DATA):
			if entriesRead == 0 {
				return nil, fmt.Errorf("enumeration returned more data without entries")
			}
			continue
		default:
			return nil, syscall.Errno(status)
		}
	}
	return nil, fmt.Errorf("enumeration exceeded %d pages", maxGroupMemberPages)
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
		var members []windowsGroupMember
		if accountType == windows.SidTypeAlias {
			members, err = ops.listGroupMembers(cfg.AdminGroup)
			if err != nil {
				return nil, AccessInfo{}, fmt.Errorf("ipc: enumerate admin_group %q members: %w", cfg.AdminGroup, err)
			}
		}
		descriptor, err = restrictedPipeDescriptor(sid, members)
		if err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: admin_group %q descriptor: %w", cfg.AdminGroup, err)
		}
		access = AccessInfo{Mode: AccessModeRestricted, AdminGroup: cfg.AdminGroup}
	}

	endpoint := Endpoint(cfg)
	l, err := ops.listenPipe(endpoint, descriptor)
	if err != nil {
		return nil, AccessInfo{}, fmt.Errorf("ipc: listen pipe %s: %w", endpoint, err)
	}
	return l, access, nil
}

func restrictedPipeDescriptor(groupSID string, members []windowsGroupMember) (string, error) {
	normalizedGroupSID, err := normalizeWindowsSID(groupSID)
	if err != nil {
		return "", fmt.Errorf("invalid group SID %q: %w", groupSID, err)
	}

	seen := map[string]struct{}{strings.ToLower(normalizedGroupSID): {}}
	memberSIDs := make([]string, 0, len(members))
	for _, member := range members {
		if member.accountType != windows.SidTypeUser {
			continue
		}
		normalizedSID, err := normalizeWindowsSID(member.sid)
		if err != nil {
			return "", fmt.Errorf("invalid direct member SID %q: %w", member.sid, err)
		}
		key := strings.ToLower(normalizedSID)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		memberSIDs = append(memberSIDs, normalizedSID)
	}
	sort.Strings(memberSIDs)

	var descriptor strings.Builder
	descriptor.WriteString("D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;")
	descriptor.WriteString(normalizedGroupSID)
	descriptor.WriteByte(')')
	for _, memberSID := range memberSIDs {
		descriptor.WriteString("(A;;GRGW;;;")
		descriptor.WriteString(memberSID)
		descriptor.WriteByte(')')
	}
	return descriptor.String(), nil
}

func normalizeWindowsSID(value string) (string, error) {
	sid, err := windows.StringToSid(value)
	if err != nil {
		return "", err
	}
	if !sid.IsValid() {
		return "", fmt.Errorf("SID is not valid")
	}
	return sid.String(), nil
}

// DialContext connects to the named-pipe endpoint.
func DialContext(ctx context.Context, endpoint string) (net.Conn, error) {
	return winio.DialPipeContext(ctx, endpoint)
}
