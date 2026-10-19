package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpmanage"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
	"github.com/shruggietech/go-schedule/internal/mcpoperate"
)

func mcpHTTPStatus(cmd *cobra.Command) (server.MCPHTTPStatusResponse, error) {
	if callback := stateFor(cmd).deps.MCPHTTPStatus; callback != nil {
		return callback(cmd)
	}
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	return clientFor(cmd).MCPHTTPStatus(ctx)
}

func mcpHTTPEnable(cmd *cobra.Command, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	if callback := stateFor(cmd).deps.MCPHTTPEnable; callback != nil {
		return callback(cmd, req)
	}
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	return clientFor(cmd).EnableMCPHTTP(ctx, req)
}

func mcpHTTPRotate(cmd *cobra.Command) (server.MCPHTTPCredentialResponse, error) {
	if callback := stateFor(cmd).deps.MCPHTTPRotate; callback != nil {
		return callback(cmd)
	}
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	return clientFor(cmd).RotateMCPHTTPCredential(ctx)
}

func mcpHTTPDisable(cmd *cobra.Command) (server.MCPHTTPStatusResponse, error) {
	if callback := stateFor(cmd).deps.MCPHTTPDisable; callback != nil {
		return callback(cmd)
	}
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	return clientFor(cmd).DisableMCPHTTP(ctx)
}

func newMCPCmd() *cobra.Command {
	mcpCmd := &cobra.Command{Use: "mcp", Short: "Local Model Context Protocol integration"}
	mcpCmd.AddCommand(newMCPServeCmd())
	mcpCmd.AddCommand(newMCPHTTPCmd())
	return mcpCmd
}

func newMCPServeCmd() *cobra.Command {
	var permission string
	var clientName string
	var requireConfirmation bool
	cmd := &cobra.Command{Use: "serve", Short: "Serve scheduler resources and explicitly authorized tools over stdio", Args: cobra.NoArgs, SilenceUsage: true}
	cmd.Flags().StringVar(&permission, "permission", string(domain.CapabilityObserve), "session permission: observe, operate, or manage")
	cmd.Flags().StringVar(&clientName, "name", "Local MCP stdio client", "audit display name for a mutation-capable session")
	cmd.Flags().BoolVar(&requireConfirmation, "require-confirmation", false, "require confirmed=true on every Manage mutation")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		capability := domain.Capability(permission)
		if capability == domain.CapabilityObserve {
			return mcpobserve.NewServer(clientFor(cmd), buildinfo.Version).Run(cmd.Context(), scopedMCPTransport(cmd))
		}
		if capability != domain.CapabilityOperate && capability != domain.CapabilityManage {
			return fmtUsage("--permission must be observe, operate, or manage")
		}
		if requireConfirmation && capability != domain.CapabilityManage {
			return fmtUsage("--require-confirmation requires --permission manage")
		}
		base := clientFor(cmd)
		if base.Remote() {
			return fmtUsage("MCP Operate stdio is available only for the local daemon")
		}
		createCtx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		session, err := base.CreateMCPSession(createCtx, clientName, capability)
		cancel()
		if err != nil {
			return err
		}
		defer func() {
			revokeCtx, revokeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer revokeCancel()
			_ = base.RevokeMCPSession(revokeCtx, session.ID)
		}()
		mcpServer := mcpobserve.NewServer(base, buildinfo.Version)
		mcpoperate.AddTools(mcpServer, mcpoperate.New(base.WithMCPSession(session.Credential)))
		if capability == domain.CapabilityManage {
			mcpmanage.AddTools(mcpServer, mcpmanage.New(base.WithMCPSession(session.Credential), requireConfirmation))
		}
		return mcpServer.Run(cmd.Context(), scopedMCPTransport(cmd))
	}
	return cmd
}

func newMCPHTTPCmd() *cobra.Command {
	httpCmd := &cobra.Command{Use: "http", Short: "Manage temporary authenticated localhost MCP access"}
	httpCmd.AddCommand(
		&cobra.Command{Use: "status", Short: "Show localhost MCP status without revealing its credential", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := mcpHTTPStatus(cmd)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), status)
			}
			if !status.Enabled {
				fmt.Fprintln(cmd.OutOrStdout(), "localhost MCP disabled")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "localhost MCP %s enabled for %s at %s (credential %s, successful requests %d)\n", status.Permission, status.ClientName, status.Endpoint, status.CredentialFingerprint, status.RequestCount)
			return nil
		}},
		newMCPHTTPEnableCmd(),
		&cobra.Command{Use: "rotate", Short: "Replace the localhost MCP credential", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := mcpHTTPRotate(cmd)
			if err != nil {
				return err
			}
			return printMCPHTTPCredential(cmd, result, "localhost MCP credential rotated")
		}},
		&cobra.Command{Use: "disable", Short: "Revoke and stop localhost MCP access", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := mcpHTTPDisable(cmd)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), status)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "localhost MCP disabled")
			return nil
		}},
	)
	return httpCmd
}

func newMCPHTTPEnableCmd() *cobra.Command {
	var port int
	var origins []string
	var clientName string
	var permission string
	var requireConfirmation bool
	cmd := &cobra.Command{Use: "enable", Short: "Start authenticated MCP on numeric IPv4 loopback", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if port < 1 || port > 65535 {
			return fmtUsage("--port must be between 1 and 65535")
		}
		capability := domain.Capability(permission)
		if capability != domain.CapabilityObserve && capability != domain.CapabilityOperate && capability != domain.CapabilityManage {
			return fmtUsage("--permission must be observe, operate, or manage")
		}
		if requireConfirmation && capability != domain.CapabilityManage {
			return fmtUsage("--require-confirmation requires --permission manage")
		}
		result, err := mcpHTTPEnable(cmd, server.MCPHTTPEnableRequest{Port: port, AllowedOrigins: origins, ClientName: clientName, Permission: capability, RequireConfirmation: requireConfirmation})
		if err != nil {
			return err
		}
		return printMCPHTTPCredential(cmd, result, "localhost MCP enabled")
	}}
	cmd.Flags().IntVar(&port, "port", 0, "loopback TCP port (required)")
	cmd.Flags().StringSliceVar(&origins, "origin", nil, "allowed loopback browser origin (repeatable)")
	cmd.Flags().StringVar(&clientName, "name", "", "display name for the active local client")
	cmd.Flags().StringVar(&permission, "permission", string(domain.CapabilityObserve), "session permission: observe, operate, or manage")
	cmd.Flags().BoolVar(&requireConfirmation, "require-confirmation", false, "require confirmed=true on every Manage mutation")
	_ = cmd.MarkFlagRequired("port")
	return cmd
}

func printMCPHTTPCredential(cmd *cobra.Command, result server.MCPHTTPCredentialResponse, heading string) error {
	if jsonOutput(cmd) {
		return printJSONTo(cmd.OutOrStdout(), result)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s for %s at %s\n", heading, result.ClientName, result.Endpoint)
	fmt.Fprintf(cmd.OutOrStdout(), "credential (shown once): %s\n", result.Credential)
	return nil
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
func scopedMCPTransport(cmd *cobra.Command) *mcp.IOTransport {
	reader, ok := cmd.InOrStdin().(io.ReadCloser)
	if !ok {
		reader = io.NopCloser(cmd.InOrStdin())
	}
	var writer io.WriteCloser = nopWriteCloser{cmd.OutOrStdout()}
	// Explicit owned pipes must close on cancellation. Shared process stdout stays open.
	if out := cmd.OutOrStdout(); out != os.Stdout {
		if closer, ok := out.(io.WriteCloser); ok {
			writer = closer
		}
	}
	return &mcp.IOTransport{Reader: reader, Writer: writer}
}
