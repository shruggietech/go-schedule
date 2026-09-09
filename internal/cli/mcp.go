package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
)

var runMCPStdio = func(cmd *cobra.Command, _ []string) error {
	return mcpobserve.RunStdio(cmd.Context(), newClient(), buildinfo.Version)
}

var mcpHTTPStatus = func(cmd *cobra.Command) (server.MCPHTTPStatusResponse, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return newClient().MCPHTTPStatus(ctx)
}

var mcpHTTPEnable = func(cmd *cobra.Command, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return newClient().EnableMCPHTTP(ctx, req)
}

var mcpHTTPRotate = func(cmd *cobra.Command) (server.MCPHTTPCredentialResponse, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return newClient().RotateMCPHTTPCredential(ctx)
}

var mcpHTTPDisable = func(cmd *cobra.Command) (server.MCPHTTPStatusResponse, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return newClient().DisableMCPHTTP(ctx)
}

func newMCPCmd() *cobra.Command {
	mcpCmd := &cobra.Command{Use: "mcp", Short: "Local Model Context Protocol integration"}
	mcpCmd.AddCommand(&cobra.Command{
		Use:          "serve",
		Short:        "Serve observe-only scheduler resources over stdio",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         runMCPStdio,
	})
	mcpCmd.AddCommand(newMCPHTTPCmd())
	return mcpCmd
}

func newMCPHTTPCmd() *cobra.Command {
	httpCmd := &cobra.Command{Use: "http", Short: "Manage temporary authenticated localhost MCP access"}
	httpCmd.AddCommand(
		&cobra.Command{Use: "status", Short: "Show localhost MCP status without revealing its credential", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := mcpHTTPStatus(cmd)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(status)
			}
			if !status.Enabled {
				fmt.Fprintln(os.Stdout, "localhost MCP disabled")
				return nil
			}
			fmt.Fprintf(os.Stdout, "localhost MCP enabled for %s at %s (credential %s, successful requests %d)\n", status.ClientName, status.Endpoint, status.CredentialFingerprint, status.RequestCount)
			return nil
		}},
		newMCPHTTPEnableCmd(),
		&cobra.Command{Use: "rotate", Short: "Replace the localhost MCP credential", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := mcpHTTPRotate(cmd)
			if err != nil {
				return err
			}
			return printMCPHTTPCredential(result, "localhost MCP credential rotated")
		}},
		&cobra.Command{Use: "disable", Short: "Revoke and stop localhost MCP access", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := mcpHTTPDisable(cmd)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(status)
			}
			fmt.Fprintln(os.Stdout, "localhost MCP disabled")
			return nil
		}},
	)
	return httpCmd
}

func newMCPHTTPEnableCmd() *cobra.Command {
	var port int
	var origins []string
	var clientName string
	cmd := &cobra.Command{Use: "enable", Short: "Start authenticated MCP on numeric IPv4 loopback", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if port < 1 || port > 65535 {
			return fmtUsage("--port must be between 1 and 65535")
		}
		result, err := mcpHTTPEnable(cmd, server.MCPHTTPEnableRequest{Port: port, AllowedOrigins: origins, ClientName: clientName})
		if err != nil {
			return err
		}
		return printMCPHTTPCredential(result, "localhost MCP enabled")
	}}
	cmd.Flags().IntVar(&port, "port", 0, "loopback TCP port (required)")
	cmd.Flags().StringSliceVar(&origins, "origin", nil, "allowed loopback browser origin (repeatable)")
	cmd.Flags().StringVar(&clientName, "name", "", "display name for the active local client")
	_ = cmd.MarkFlagRequired("port")
	return cmd
}

func printMCPHTTPCredential(result server.MCPHTTPCredentialResponse, heading string) error {
	if jsonOut {
		return printJSON(result)
	}
	fmt.Fprintf(os.Stdout, "%s for %s at %s\n", heading, result.ClientName, result.Endpoint)
	fmt.Fprintf(os.Stdout, "credential (shown once): %s\n", result.Credential)
	return nil
}
