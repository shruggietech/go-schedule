package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type daemonIdentityClient interface {
	Manifest(context.Context) (server.ManifestResponse, error)
	RenameDaemon(context.Context, string) (server.ManifestResponse, error)
	ResetDaemonIdentity(context.Context, string) (server.ManifestResponse, error)
}

func newDaemonCmd() *cobra.Command { return newDaemonCmdWithClient(newClient()) }

func newDaemonCmdWithClient(api daemonIdentityClient) *cobra.Command {
	command := &cobra.Command{Use: "daemon", Short: "Inspect and manage daemon identity"}
	command.AddCommand(
		&cobra.Command{
			Use:   "manifest",
			Short: "Show daemon identity and capabilities",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				ctx, cancel := reqCtx()
				defer cancel()
				manifest, err := api.Manifest(ctx)
				if err != nil {
					return err
				}
				return printManifest(cmd, manifest)
			},
		},
		&cobra.Command{
			Use:   "rename <display-name>",
			Short: "Change the daemon display name",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				name, err := domain.NormalizeDaemonDisplayName(args[0])
				if err != nil {
					return fmtUsage("display name must contain 1 through 80 characters and no control characters")
				}
				ctx, cancel := reqCtx()
				defer cancel()
				manifest, err := api.RenameDaemon(ctx, name)
				if err != nil {
					return err
				}
				if jsonOut {
					return printJSONTo(cmd.OutOrStdout(), manifest)
				}
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Daemon renamed to %s (%s).\n", manifest.DisplayName, manifest.InstallationID)
				return err
			},
		},
		newDaemonResetCmd(api),
	)
	return command
}

func newDaemonResetCmd(api daemonIdentityClient) *cobra.Command {
	var confirmation string
	command := &cobra.Command{
		Use:   "reset-identity",
		Short: "Reset a cloned daemon identity",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if confirmation == "" {
				return fmtUsage("--confirm must contain the current installation ID")
			}
			ctx, cancel := reqCtx()
			defer cancel()
			manifest, err := api.ResetDaemonIdentity(ctx, confirmation)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSONTo(cmd.OutOrStdout(), manifest)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Installation identity reset to %s for %s.\n", manifest.InstallationID, manifest.DisplayName)
			return err
		},
	}
	command.Flags().StringVar(&confirmation, "confirm", "", "exact current installation ID")
	return command
}

func printManifest(cmd *cobra.Command, manifest server.ManifestResponse) error {
	if jsonOut {
		return printJSONTo(cmd.OutOrStdout(), manifest)
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "Daemon: %s\nInstallation ID: %s\nVersion: %s\nMode: %s\nPlatform: %s/%s\nCapabilities: %s\n", manifest.DisplayName, manifest.InstallationID, manifest.ProductVersion, manifest.OperatingMode, manifest.Platform.OS, manifest.Platform.Architecture, strings.Join(manifest.Capabilities, ", "))
	return err
}
