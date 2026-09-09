package cli

import (
	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
)

var runMCPStdio = func(cmd *cobra.Command, _ []string) error {
	return mcpobserve.RunStdio(cmd.Context(), newClient(), buildinfo.Version)
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
	return mcpCmd
}
