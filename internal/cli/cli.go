// Package cli implements the gosched command-line interface (cobra). Commands
// operate on the daemon through the shared API client, so the CLI and GUI act on
// identical state. Results go to stdout; diagnostics/errors go to stderr; exit
// codes follow the contract (0 ok, 1 runtime error, 2 usage/validation).
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

// errUsage marks validation/usage failures so Execute can return exit code 2.
var errUsage = errors.New("usage")

var jsonOut bool
var selectedClient *client.Client
var targetFlags remoteTargetFlags

// Execute runs the root command and returns a process exit code.
func Execute() int {
	root := newRoot()
	return handleExecuteError(os.Stderr, root.Execute())
}

func handleExecuteError(stderr io.Writer, err error) int {
	if err == nil {
		return 0
	}
	if !isReported(err) {
		fmt.Fprintln(stderr, "gosched: "+err.Error())
	}
	if errors.Is(err, errUsage) {
		return 2
	}
	// Server-side validation failures map to the usage/validation exit code.
	var se *client.StatusError
	if errors.As(err, &se) && se.Code == server.CodeValidation {
		return 2
	}
	return 1
}

type reportedError struct{ err error }

func (e *reportedError) Error() string { return e.err.Error() }
func (e *reportedError) Unwrap() error { return e.err }

func reported(err error) error { return &reportedError{err: err} }

func isReported(err error) bool {
	var target *reportedError
	return errors.As(err, &target)
}

func fmtUsage(message string) error { return fmt.Errorf("%w: %s", errUsage, message) }

func newRoot() *cobra.Command {
	jsonOut = false
	cfg, _ := config.Load("")
	selectedClient = client.NewSwitchable(client.New(ipc.Endpoint(cfg)))
	targetFlags = remoteTargetFlags{}
	root := &cobra.Command{
		Use:           "gosched",
		Short:         "Cross-platform scheduler with readable schedules and supported cron",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       buildinfo.Version,
	}
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		resolved, description, err := resolveTarget(cmd.Context(), targetFlags)
		if err != nil {
			return err
		}
		selectedClient.Use(resolved)
		if description != "" && !jsonOut {
			fmt.Fprintln(cmd.ErrOrStderr(), description)
		}
		return nil
	}
	root.PersistentFlags().BoolVar(&jsonOut, "json", false, "machine-readable JSON output")
	root.PersistentFlags().StringVar(&targetFlags.Profile, "profile", "", "saved remote profile ID or unambiguous label")
	root.PersistentFlags().StringVar(&targetFlags.Endpoint, "endpoint", "", "explicit remote HTTPS origin")
	root.PersistentFlags().StringVar(&targetFlags.DaemonID, "daemon-id", "", "expected daemon installation ID for an explicit target")
	root.PersistentFlags().StringVar(&targetFlags.CredentialID, "credential-id", "", "native credential reference for an explicit target")
	root.PersistentFlags().StringVar(&targetFlags.CertificateFile, "certificate-file", "", "trusted PEM certificate file for an explicit target")
	root.AddCommand(
		newTaskCmd(),
		newChainCmd(),
		newTriggerCmd(),
		newWatcherCmd(),
		newCronCmd(),
		newGroupCmd(),
		newRunsCmd(),
		newLogsCmd(),
		newAlertsCmd(),
		newNotificationCmd(),
		newServiceCmd(),
		newGUICmd(),
		newHealthCmd(),
		newDaemonCmd(),
		newMCPCmd(),
		newActorCmd(),
		newAuditCmd(),
		newPairingCmd(),
		newCredentialCmd(),
		newProfileCmd(),
	)
	return root
}

func newClient() *client.Client {
	if selectedClient != nil {
		return selectedClient
	}
	cfg, _ := config.Load("")
	return client.New(ipc.Endpoint(cfg))
}

func reqCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// printJSON writes v as indented JSON to stdout.
func printJSON(v any) error {
	return printJSONTo(os.Stdout, v)
}

func printJSONTo(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func newHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check that the daemon is running",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			h, err := newClient().Health(ctx)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(h)
			}
			fmt.Fprintf(os.Stdout, "daemon ok (version %s)\n", h.Version)
			return nil
		},
	}
}
