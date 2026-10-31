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

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/buildinfo"
)

// errUsage marks validation/usage failures so Execute can return exit code 2.
var errUsage = errors.New("usage")

// Execute runs the root command and returns a process exit code.
func Execute() int {
	root := newRoot()
	return executeInvocation(context.Background(), root, os.Args[1:])
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

func newRoot() *cobra.Command { return newRootWithDependencies(Dependencies{}) }

func newRootWithDependencies(deps Dependencies) *cobra.Command {
	state := makeInvocation(deps)
	targetFlags := &state.flags
	root := &cobra.Command{
		Use:           "gosched",
		Short:         "Cross-platform scheduler with readable schedules and supported cron",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       buildinfo.Version,
	}
	root.SetContext(context.WithValue(context.Background(), invocationKey{}, state))
	invocationStreams(root, deps)
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.SetContext(context.WithValue(cmd.Context(), invocationKey{}, state))
		return contextualCommandError(cmd, args, state.resolve(cmd))
	}
	root.PersistentFlags().BoolVar(&state.json, "json", false, "machine-readable JSON output")
	root.PersistentFlags().StringVar(&targetFlags.Profile, "profile", "", "saved remote profile ID or unambiguous label")
	root.PersistentFlags().StringVar(&targetFlags.Endpoint, "endpoint", "", "explicit remote HTTPS origin")
	root.PersistentFlags().StringVar(&targetFlags.Endpoint, "daemon-endpoint", "", "explicit remote daemon HTTPS origin (independent of webhook --endpoint)")
	root.PersistentFlags().StringVar(&targetFlags.DaemonID, "daemon-id", "", "expected daemon installation ID for an explicit target")
	root.PersistentFlags().StringVar(&targetFlags.CredentialID, "credential-id", "", "native credential reference for an explicit target")
	root.PersistentFlags().StringVar(&targetFlags.CertificateFile, "certificate-file", "", "trusted PEM certificate file for an explicit target")
	root.PersistentFlags().StringVar(&targetFlags.IPCPath, "ipc-path", "", "explicit absolute local socket path or Windows named pipe")
	root.AddCommand(
		newTaskCmd(),
		newScheduleCmd(),
		newCalendarCmd(),
		newRuntimeCmd(),
		newSystemsCmd(),
		newSearchCmd(),
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
		newDaemonCmdWithClient(state.client),
		newMCPCmd(),
		newActorCmdWithClient(state.client),
		newAuditCmdWithClient(state.client),
		newPairingCmdWithClient(state.client),
		newCredentialCmdWithClient(state.client),
		newProfileCmd(),
		newPreferencesCmd(),
		newBundleCmd(),
	)
	configureCommandContracts(root)
	return root
}

func reqCtx(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	return context.WithTimeout(cmd.Context(), 10*time.Second)
}

// printJSONTo writes v as indented JSON to the invocation's selected output.
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
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			h, err := clientFor(cmd).Health(ctx)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), h)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "daemon ok (version %s)\n", h.Version)
			return nil
		},
	}
}
