//go:build windows

package ipc

import (
	"errors"
	"net"
	"strings"
	"testing"

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
	const sid = "S-1-5-21-1-2-3-1001"
	var gotDescriptor string
	ops := windowsIPCOperations{
		lookupGroup: func(name string) (string, uint32, error) {
			if name != cfg.AdminGroup {
				t.Fatalf("lookup name = %q", name)
			}
			return sid, windows.SidTypeAlias, nil
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
	want := "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;" + sid + ")"
	if gotDescriptor != want {
		t.Fatalf("descriptor = %q, want %q", gotDescriptor, want)
	}
	if access != (AccessInfo{Mode: AccessModeRestricted, AdminGroup: cfg.AdminGroup}) {
		t.Fatalf("access = %+v", access)
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

var _ net.Listener = stubListener{}
var _ net.Addr = stubAddr("")
