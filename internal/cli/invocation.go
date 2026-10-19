package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

// ProfileStore is the shared client-local connection metadata boundary.
type ProfileStore interface {
	Load() (clientprofile.Collection, error)
	Add(clientprofile.Profile) (clientprofile.Profile, error)
	Rename(string, string) (clientprofile.Profile, error)
	SetActive(string) error
	RemoveWith(string, func(clientprofile.Profile) error) (clientprofile.Profile, error)
}

// SecretStore retains native credential handling unless explicitly injected.
type SecretStore interface {
	Probe() error
	Save(string, string, string) error
	Load(string, string) (string, error)
	LoadContext(context.Context, string, string) (string, error)
	Delete(string, string) error
}

// Dependencies scopes one CLI invocation. Empty fields select production adapters.
// Only explicit callers supply recording adapters; environment variables cannot downgrade secrets.
type Dependencies struct {
	Client                               *client.Client
	Profiles                             ProfileStore
	Secrets                              SecretStore
	PreferencePath, LegacyPreferencePath string
	Now                                  func() time.Time
	ServiceControl                       serviceControlFunc
	DaemonPath, GUIPath                  func() (string, error)
	LaunchGUI                            func(string) error
	MCPHTTPStatus                        func(*cobra.Command) (server.MCPHTTPStatusResponse, error)
	MCPHTTPEnable                        func(*cobra.Command, server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error)
	MCPHTTPRotate                        func(*cobra.Command) (server.MCPHTTPCredentialResponse, error)
	MCPHTTPDisable                       func(*cobra.Command) (server.MCPHTTPStatusResponse, error)
	In                                   io.Reader
	Out, Err                             io.Writer
}

type invocation struct {
	deps   Dependencies
	json   bool
	flags  remoteTargetFlags
	client *client.Client
}
type invocationKey struct{}

func makeInvocation(deps Dependencies) *invocation {
	if deps.Profiles == nil {
		deps.Profiles = clientprofile.NewStore("")
	}
	if deps.Secrets == nil {
		deps.Secrets = clientsecret.New()
	}
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.ServiceControl == nil {
		deps.ServiceControl = defaultServiceControl
	}
	if deps.DaemonPath == nil {
		deps.DaemonPath = daemonPath
	}
	if deps.GUIPath == nil {
		deps.GUIPath = guiPath
	}
	if deps.LaunchGUI == nil {
		deps.LaunchGUI = launchGUI
	}
	initial := deps.Client
	if initial == nil {
		cfg, _ := config.Load("")
		initial = client.New(ipc.Endpoint(cfg))
	}
	return &invocation{deps: deps, client: client.NewSwitchable(initial)}
}

func stateFor(cmd *cobra.Command) *invocation {
	if cmd == nil {
		return makeInvocation(Dependencies{})
	}
	if cmd.Context() == nil {
		cmd.SetContext(context.Background())
	}
	if state, ok := cmd.Context().Value(invocationKey{}).(*invocation); ok {
		return state
	}
	root := cmd.Root()
	if root.Context() == nil {
		root.SetContext(context.Background())
	}
	if state, ok := root.Context().Value(invocationKey{}).(*invocation); ok {
		return state
	}
	state := makeInvocation(Dependencies{})
	root.SetContext(context.WithValue(root.Context(), invocationKey{}, state))
	return state
}
func clientFor(cmd *cobra.Command) *client.Client { return stateFor(cmd).client }
func jsonOutput(cmd *cobra.Command) bool          { return stateFor(cmd).json }

// MarkLocal marks commands that operate on client-local persistence or host adapters.
func MarkLocal(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations["cli.local"] = "true"
}
func isLocal(cmd *cobra.Command) bool {
	for item := cmd; item != nil; item = item.Parent() {
		if item.Name() == "completion" || item.Name() == "__complete" || item.Name() == "__completeNoDesc" {
			return true
		}
		if item.Annotations["cli.local"] == "true" {
			return true
		}
	}
	return false
}

func validateIPCPath(path string) error {
	if path == "" || strings.TrimSpace(path) != path || strings.ContainsAny(path, "\x00\r\n") {
		return fmtUsage("--ipc-path must name a nonempty absolute local socket or named pipe")
	}
	if runtime.GOOS == "windows" {
		if !strings.HasPrefix(strings.ToLower(path), `\\.\pipe\`) || len(path) <= len(`\\.\pipe\`) || strings.ContainsAny(path[len(`\\.\pipe\`):], `/\`) {
			return fmtUsage("--ipc-path must name a local Windows pipe such as \\\\.\\pipe\\goschedd-test")
		}
	} else if !filepath.IsAbs(path) {
		return fmtUsage("--ipc-path must name an absolute Unix socket path")
	}
	return nil
}

func (s *invocation) resolve(cmd *cobra.Command) error {
	if cmd.Flags().Changed("ipc-path") {
		for _, name := range []string{"profile", "endpoint", "daemon-id", "credential-id", "certificate-file"} {
			if cmd.Flags().Changed(name) {
				return fmtUsage("--ipc-path cannot be combined with remote target flags")
			}
		}
	}
	if cmd.Flags().Changed("ipc-path") && s.flags.IPCPath == "" {
		return fmtUsage("--ipc-path must name a nonempty absolute local socket or named pipe")
	}
	if isLocal(cmd) {
		return nil
	}
	if flag := cmd.Flags().Lookup("all"); flag != nil && flag.Value.String() == "true" && cmd.Annotations["cli.aggregate"] == "true" {
		return nil
	}
	if s.flags.IPCPath == "" && s.flags.Profile == "" && s.flags.Endpoint == "" && s.flags.DaemonID == "" && s.flags.CredentialID == "" && s.flags.CertificateFile == "" && s.deps.Client != nil {
		return nil
	}
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	resolved, description, err := resolveTargetWith(ctx, s.flags, s.deps)
	if err != nil {
		return err
	}
	s.client.Use(resolved)
	if description != "" && !s.json {
		_, err = fmt.Fprintln(cmd.ErrOrStderr(), description)
	}
	return err
}

// ExecuteWithDependencies runs the same command handlers with explicit scoped adapters.
// This is the entry point for focused unattended command harnesses.
func ExecuteWithDependencies(ctx context.Context, args []string, deps Dependencies) int {
	root := newRootWithDependencies(deps)
	root.SetArgs(args)
	return handleExecuteError(root.ErrOrStderr(), root.ExecuteContext(ctx))
}

func invocationStreams(root *cobra.Command, deps Dependencies) {
	if deps.In == nil {
		deps.In = os.Stdin
	}
	if deps.Out == nil {
		deps.Out = os.Stdout
	}
	if deps.Err == nil {
		deps.Err = os.Stderr
	}
	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.Err)
}
