//go:build windows

package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/internal/config"
)

type stubListener struct{}

func (stubListener) Accept() (net.Conn, error) { return nil, errors.New("not implemented") }
func (stubListener) Close() error              { return nil }
func (stubListener) Addr() net.Addr            { return stubAddr("pipe") }

type stubAddr string

func (a stubAddr) Network() string { return "pipe" }
func (a stubAddr) String() string  { return string(a) }

func TestListenWindowsRestrictedDescriptor(t *testing.T) {
	cfg := config.Default()
	cfg.IPCPath = `\\.\pipe\test`
	const sid = "S-1-5-21-1-2-3-9000"
	var gotDescriptor string
	ops := windowsIPCOperations{
		lookupGroup: func(name string) (string, uint32, error) {
			if name != cfg.AdminGroup {
				t.Fatalf("lookup name = %q", name)
			}
			return sid, windows.SidTypeAlias, nil
		},
		listGroupMembers: func(name string) ([]windowsGroupMember, error) {
			if name != cfg.AdminGroup {
				t.Fatalf("member lookup name = %q", name)
			}
			return []windowsGroupMember{
				{sid: "S-1-5-21-1-2-3-1003", accountType: windows.SidTypeGroup},
				{sid: "S-1-5-21-1-2-3-1002", accountType: windows.SidTypeUser},
				{sid: "S-1-5-21-1-2-3-1001", accountType: windows.SidTypeUser},
				{sid: "s-1-5-21-1-2-3-1002", accountType: windows.SidTypeUser},
			}, nil
		},
		listenPipe: func(endpoint, descriptor string) (net.Listener, error) {
			if endpoint != cfg.IPCPath {
				t.Fatalf("endpoint = %q", endpoint)
			}
			gotDescriptor = descriptor
			return stubListener{}, nil
		},
	}

	ln, access, err := listenWindows(cfg, ops)
	if err != nil {
		t.Fatalf("listenWindows() error = %v", err)
	}
	defer ln.Close()
	want := "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;" + sid + ")" +
		"(A;;GRGW;;;S-1-5-21-1-2-3-1001)" +
		"(A;;GRGW;;;S-1-5-21-1-2-3-1002)"
	if gotDescriptor != want {
		t.Fatalf("descriptor = %q, want %q", gotDescriptor, want)
	}
	if access != (AccessInfo{Mode: AccessModeRestricted, AdminGroup: cfg.AdminGroup}) {
		t.Fatalf("access = %+v", access)
	}
}

func TestListenWindowsRestrictedDescriptorEmptyGroupKeepsGroupACE(t *testing.T) {
	cfg := config.Default()
	const sid = "S-1-5-21-1-2-3-1001"
	var descriptor string
	ln, _, err := listenWindows(cfg, windowsIPCOperations{
		lookupGroup: func(string) (string, uint32, error) {
			return sid, windows.SidTypeAlias, nil
		},
		listGroupMembers: func(string) ([]windowsGroupMember, error) { return nil, nil },
		listenPipe: func(_ string, got string) (net.Listener, error) {
			descriptor = got
			return stubListener{}, nil
		},
	})
	if err != nil {
		t.Fatalf("listenWindows() error = %v", err)
	}
	defer ln.Close()
	want := "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;" + sid + ")"
	if descriptor != want {
		t.Fatalf("descriptor = %q, want %q", descriptor, want)
	}
	if strings.Contains(descriptor, ";;;AU)") || strings.Contains(descriptor, ";;;WD)") {
		t.Fatalf("restricted descriptor broadened access: %q", descriptor)
	}
}

func TestListenWindowsNonAliasGroupKeepsGroupACEWithoutEnumeration(t *testing.T) {
	cfg := config.Default()
	const sid = "S-1-5-21-1-2-3-9000"
	var descriptor string
	ln, _, err := listenWindows(cfg, windowsIPCOperations{
		lookupGroup: func(string) (string, uint32, error) {
			return sid, windows.SidTypeGroup, nil
		},
		listGroupMembers: func(string) ([]windowsGroupMember, error) {
			t.Fatal("non-alias group was passed to local-group enumeration")
			return nil, nil
		},
		listenPipe: func(_ string, got string) (net.Listener, error) {
			descriptor = got
			return stubListener{}, nil
		},
	})
	if err != nil {
		t.Fatalf("listenWindows() error = %v", err)
	}
	defer ln.Close()
	want := "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;" + sid + ")"
	if descriptor != want {
		t.Fatalf("descriptor = %q, want %q", descriptor, want)
	}
}

func TestListenWindowsFailsClosedOnMemberEnumeration(t *testing.T) {
	cfg := config.Default()
	listErr := errors.New("access denied")
	listened := false
	_, _, err := listenWindows(cfg, windowsIPCOperations{
		lookupGroup: func(string) (string, uint32, error) {
			return "S-1-5-21-1-2-3-1001", windows.SidTypeAlias, nil
		},
		listGroupMembers: func(string) ([]windowsGroupMember, error) { return nil, listErr },
		listenPipe: func(string, string) (net.Listener, error) {
			listened = true
			return stubListener{}, nil
		},
	})
	if !errors.Is(err, listErr) || !strings.Contains(err.Error(), "enumerate admin_group") {
		t.Fatalf("error = %v, want wrapped member enumeration failure", err)
	}
	if listened {
		t.Fatal("listener opened after member enumeration failure")
	}
}

func TestListenWindowsRejectsInvalidMemberSID(t *testing.T) {
	cfg := config.Default()
	listened := false
	_, _, err := listenWindows(cfg, windowsIPCOperations{
		lookupGroup: func(string) (string, uint32, error) {
			return "S-1-5-21-1-2-3-1001", windows.SidTypeAlias, nil
		},
		listGroupMembers: func(string) ([]windowsGroupMember, error) {
			return []windowsGroupMember{{sid: "S-1-5-21-1);(A;;GA;;;WD", accountType: windows.SidTypeUser}}, nil
		},
		listenPipe: func(string, string) (net.Listener, error) {
			listened = true
			return stubListener{}, nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "member SID") {
		t.Fatalf("error = %v, want invalid member SID failure", err)
	}
	if listened {
		t.Fatal("listener opened after invalid member SID")
	}
}

func TestListenWindowsRejectsMissingOrNonGroupAccount(t *testing.T) {
	for _, test := range []struct {
		name        string
		sid         string
		accountType uint32
		lookupErr   error
	}{
		{name: "missing", lookupErr: errors.New("none mapped")},
		{name: "empty SID", accountType: windows.SidTypeGroup},
		{name: "user", sid: "S-1-5-21-1", accountType: windows.SidTypeUser},
		{name: "domain", sid: "S-1-5-21-1", accountType: windows.SidTypeDomain},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.Default()
			listened := false
			_, _, err := listenWindows(cfg, windowsIPCOperations{
				lookupGroup: func(string) (string, uint32, error) {
					return test.sid, test.accountType, test.lookupErr
				},
				listGroupMembers: func(string) ([]windowsGroupMember, error) { return nil, nil },
				listenPipe: func(string, string) (net.Listener, error) {
					listened = true
					return stubListener{}, nil
				},
			})
			if err == nil || !strings.Contains(err.Error(), "admin_group") || !strings.Contains(err.Error(), cfg.AdminGroup) {
				t.Fatalf("error = %v, want actionable admin_group failure", err)
			}
			if listened {
				t.Fatal("listener opened after group resolution failure")
			}
		})
	}
}

func TestListenWindowsCompatibilityMode(t *testing.T) {
	cfg := config.Default()
	cfg.AdminGroup = ""
	lookupCalled := false
	var descriptor string
	ln, access, err := listenWindows(cfg, windowsIPCOperations{
		lookupGroup: func(string) (string, uint32, error) {
			lookupCalled = true
			return "", 0, nil
		},
		listGroupMembers: func(string) ([]windowsGroupMember, error) {
			t.Fatal("compatibility mode enumerated group members")
			return nil, nil
		},
		listenPipe: func(_ string, got string) (net.Listener, error) {
			descriptor = got
			return stubListener{}, nil
		},
	})
	if err != nil {
		t.Fatalf("listenWindows() error = %v", err)
	}
	defer ln.Close()
	if lookupCalled {
		t.Fatal("compatibility mode attempted a group lookup")
	}
	if descriptor != compatibilityPipeSDDL {
		t.Fatalf("descriptor = %q, want %q", descriptor, compatibilityPipeSDDL)
	}
	if access != (AccessInfo{Mode: AccessModeCompatibility}) {
		t.Fatalf("access = %+v", access)
	}
}

func TestListenWindowsCurrentDirectMemberWithoutAliasToken(t *testing.T) {
	cfg := config.Default()
	groupSID, _, err := productionWindowsIPCOperations.lookupGroup(cfg.AdminGroup)
	if err != nil {
		t.Skipf("configured local group unavailable: %v", err)
	}
	members, err := productionWindowsIPCOperations.listGroupMembers(cfg.AdminGroup)
	if err != nil {
		t.Skipf("configured local group cannot be enumerated: %v", err)
	}

	token := windows.GetCurrentProcessToken()
	user, err := token.GetTokenUser()
	if err != nil {
		t.Fatalf("current token user: %v", err)
	}
	userSID := user.User.Sid.String()
	directMember := false
	for _, member := range members {
		if member.accountType == windows.SidTypeUser && strings.EqualFold(member.sid, userSID) {
			directMember = true
			break
		}
	}
	if !directMember {
		t.Skipf("current user %s is not a direct %s member", userSID, cfg.AdminGroup)
	}
	groups, err := token.GetTokenGroups()
	if err != nil {
		t.Fatalf("current token groups: %v", err)
	}
	for _, group := range groups.AllGroups() {
		if group.Sid != nil && strings.EqualFold(group.Sid.String(), groupSID) {
			t.Skipf("current token already contains %s SID", cfg.AdminGroup)
		}
	}

	cfg.IPCPath = `\\.\pipe\goschedd-s038-` + fmt.Sprint(os.Getpid())
	ln, access, err := Listen(cfg)
	if err != nil {
		t.Fatalf("Listen() for direct-member native contract: %v", err)
	}
	defer ln.Close()
	if access.Mode != AccessModeRestricted {
		t.Fatalf("access mode = %q, want restricted", access.Mode)
	}

	acceptErr := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			err = conn.Close()
		}
		acceptErr <- err
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := DialContext(ctx, cfg.IPCPath)
	if err != nil {
		t.Fatalf("standard-token direct member could not dial restricted pipe: %v", err)
	}
	_ = conn.Close()
	if err := <-acceptErr; err != nil {
		t.Fatalf("accept direct-member connection: %v", err)
	}
}

var _ net.Listener = stubListener{}
var _ net.Addr = stubAddr("")
