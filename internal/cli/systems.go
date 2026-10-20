package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientobserve"
)

func newSystemsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "systems", Short: "Read identified operational observations from selected or all saved schedulers"}
	cmd.AddCommand(newSystemsSummaryCmd(), newSearchCmd())
	return cmd
}

func newSystemsSummaryCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{Use: "summary", Short: "Show an operational summary, preserving independent target failures", Args: cobra.NoArgs, Annotations: map[string]string{"cli.aggregate": "true"}, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		targets, err := observationTargets(ctx, cmd, all)
		if err != nil {
			return err
		}
		observer := clientobserve.New()
		observer.Now = stateFor(cmd).deps.Now
		result, observeErr := observer.Summary(ctx, targets)
		if jsonOutput(cmd) {
			err = printJSONTo(cmd.OutOrStdout(), result)
		} else {
			err = writeSystemSummary(cmd.OutOrStdout(), result)
		}
		if err != nil {
			return err
		}
		return observeErr
	}}
	cmd.Flags().BoolVar(&all, "all", false, "observe the local scheduler and every saved profile with independent outcomes")
	return cmd
}

func observationTargets(ctx context.Context, cmd *cobra.Command, all bool) ([]clientobserve.Target, error) {
	state := stateFor(cmd)
	api := clientFor(cmd)
	registration := clientobserve.Registration{Key: "local", Label: "This computer", Endpoint: api.Endpoint()}
	if all {
		flags := state.flags
		for _, name := range []string{"profile", "endpoint", "daemon-id", "credential-id", "certificate-file"} {
			if flag := cmd.Flag(name); flag != nil && flag.Changed {
				return nil, fmtUsage("--all cannot be combined with a selected remote profile or explicit remote target flags")
			}
		}
		if flags.IPCPath != "" {
			if err := validateIPCPath(flags.IPCPath); err != nil {
				return nil, err
			}
			api = client.New(flags.IPCPath)
			registration.Endpoint = flags.IPCPath
		}
		local := clientobserve.Target{Registration: registration, Client: api}
		return clientobserve.Targets(ctx, local, state.deps.Profiles, state.deps.Secrets, nil), nil
	}
	if api.Remote() {
		registration = clientobserve.Registration{Key: "explicit", Label: "Explicit remote target", Endpoint: api.Endpoint(), ExpectedDaemonID: state.flags.DaemonID}
		if state.flags.Profile != "" {
			collection, err := state.deps.Profiles.Load()
			if err != nil {
				return nil, fmt.Errorf("load selected profile: %w", err)
			}
			profile, err := collection.Find(state.flags.Profile)
			if err != nil {
				return nil, fmtUsage("selected remote profile is missing or ambiguous")
			}
			registration.Key = profile.ID
			registration.Label = profile.Label
			registration.ExpectedDaemonID = profile.DaemonID
		}
	}
	return []clientobserve.Target{{Registration: registration, Client: api}}, nil
}

func writeSystemSummary(w io.Writer, snapshot clientobserve.Snapshot) error {
	for index, item := range snapshot.Observations {
		registration := item.Registration
		if index > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "Source: %s\nRegistration: %s\nEndpoint: %s\nExpected daemon: %s\nObserved daemon: %s\nObserved at: %s\n", registration.Label, registration.Key, registration.Endpoint, displayOrDash(registration.ExpectedDaemonID), displayOrDash(registration.ObservedDaemonID), item.ObservedAt.Format(time.RFC3339)); err != nil {
			return err
		}
		if item.Failure != nil {
			if _, err := fmt.Fprintf(w, "Outcome: %s\n%s\n", item.Failure.Code, item.Failure.Message); err != nil {
				return err
			}
			continue
		}
		summary := item.Summary
		if _, err := fmt.Fprintf(w, "Outcome: ok\nActive tasks: %d\nRecent failures: %d\nUnacknowledged alerts: %d\nNotification problems: %d\n", summary.ActiveTaskCount, summary.RecentFailureCount, summary.UnacknowledgedAlertCount, summary.NotificationProblemCount); err != nil {
			return err
		}
		if next := summary.NextOccurrence; next != nil {
			if _, err := fmt.Fprintf(w, "Next occurrence: %s (task %s) at %s\n", next.TaskName, next.TaskID, next.ScheduledFor.Format(time.RFC3339)); err != nil {
				return err
			}
		}
		if failure := summary.RecentFailure; failure != nil {
			if _, err := fmt.Fprintf(w, "Latest failure: run %s, task %s (%s), ended %s\n", failure.RunID, failure.TaskID, failure.TaskName, failure.EndedAt.Format(time.RFC3339)); err != nil {
				return err
			}
		}
		if alert := summary.UnacknowledgedAlert; alert != nil {
			if _, err := fmt.Fprintf(w, "Latest unacknowledged alert: %s, task %s, run %s, %s/%s, created %s\n", alert.AlertID, displayOrDash(alert.TaskID), displayOrDash(alert.RunID), alert.Severity, alert.Kind, alert.CreatedAt.Format(time.RFC3339)); err != nil {
				return err
			}
		}
		if delivery := summary.NotificationProblem; delivery != nil {
			if _, err := fmt.Fprintf(w, "Latest notification problem: delivery %s, task %s, run %s, channel %s, state %s, created %s\n", delivery.DeliveryID, displayOrDash(delivery.TaskID), displayOrDash(delivery.RunID), delivery.ChannelName, delivery.State, delivery.CreatedAt.Format(time.RFC3339)); err != nil {
				return err
			}
		}
	}
	return nil
}
