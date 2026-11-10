package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func newRunsCmd() *cobra.Command {
	var task string
	var limit int
	var details bool
	cmd := &cobra.Command{
		Use:   "runs",
		Short: "Show run history and control active executions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if limit < 0 {
				return fmtUsage("--limit must be a non-negative number of run records")
			}
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			runs, err := clientFor(cmd).ListRuns(ctx, task, limit)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), runs)
			}
			return writeRunHistory(cmd.OutOrStdout(), runs, details)
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "filter by task ID")
	cmd.Flags().IntVar(&limit, "limit", 50, "maximum rows (0 uses the daemon default)")
	cmd.Flags().BoolVar(&details, "details", false, "include run identities, lifecycle timestamps, captured output and provenance")
	cmd.AddCommand(newActiveRunsCmd(), newGetRunCmd(), newCancelRunCmd())
	return cmd
}

func writeRunHistory(w io.Writer, runs []domain.Run, details bool) error {
	if details {
		for index, run := range runs {
			if index > 0 {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}
			if err := writeRunDetail(w, run); err != nil {
				return err
			}
		}
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "SCHEDULED\tOUTCOME\tTRIGGER\tSOURCE TASK\tSOURCE RUN\tSOURCE TRIGGER\tSOURCE WATCHER\tEXIT")
	for _, r := range runs {
		exit := "-"
		if r.ExitCode != nil {
			exit = fmt.Sprintf("%d", *r.ExitCode)
		}
		sourceTask, sourceRun := runSourceColumns(r)
		sourceTrigger := r.SourceTriggerID
		if sourceTrigger == "" {
			sourceTrigger = "-"
		}
		sourceWatcher := r.SourceWatcherID
		if sourceWatcher == "" {
			sourceWatcher = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", r.ScheduledFor.Format(time.RFC3339), r.Outcome, r.Trigger, sourceTask, sourceRun, sourceTrigger, sourceWatcher, exit)
	}
	return tw.Flush()
}

func newActiveRunsCmd() *cobra.Command {
	return &cobra.Command{Use: "active", Short: "Show authoritative in-memory executions that are still running", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		runs, err := clientFor(cmd).ListActiveRuns(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), runs)
		}
		return writeRunHistory(cmd.OutOrStdout(), runs, true)
	}}
}

func newGetRunCmd() *cobra.Command {
	return &cobra.Command{Use: "get ID", Short: "Show one exact run, including captured output and source correlation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		run, err := clientFor(cmd).GetRun(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), run)
		}
		return writeRunDetail(cmd.OutOrStdout(), run)
	}}
}

func writeRunDetail(w io.Writer, run domain.Run) error {
	state := string(run.Outcome)
	if state == "" {
		state = "unavailable"
		if run.StartedAt != nil && run.EndedAt == nil {
			state = "running"
		}
	}
	started, ended, exit := "-", "-", "-"
	if run.StartedAt != nil {
		started = run.StartedAt.Format(time.RFC3339Nano)
	}
	if run.EndedAt != nil {
		ended = run.EndedAt.Format(time.RFC3339Nano)
	}
	if run.ExitCode != nil {
		exit = fmt.Sprint(*run.ExitCode)
	}
	_, err := fmt.Fprintf(w, "Run ID: %s\nTask ID: %s\nScheduled: %s\nStarted: %s\nEnded: %s\nState: %s\nOutcome: %s\nExit code: %s\nStart failed: %t\nCancellation status: %s\nCancellation error: %s\nTrigger: %s\nSource task: %s\nSource run: %s\nSource trigger: %s\nSource watcher: %s\nOutput truncated: %t\nOutput:\n%s\n", run.ID, run.TaskID, run.ScheduledFor.Format(time.RFC3339Nano), started, ended, state, displayOrDash(string(run.Outcome)), exit, run.StartFailed, displayOrDash(run.CancellationStatus), displayOrDash(run.CancellationError), run.Trigger, displayOrDash(run.SourceTaskID), displayOrDash(run.SourceRunID), displayOrDash(run.SourceTriggerID), displayOrDash(run.SourceWatcherID), run.OutputTruncated, run.Output)
	return err
}

func runSourceColumns(run domain.Run) (string, string) {
	sourceTask, sourceRun := run.SourceTaskID, run.SourceRunID
	if sourceTask == "" {
		sourceTask = "-"
	}
	if sourceRun == "" {
		sourceRun = "-"
	}
	return sourceTask, sourceRun
}

func newAlertsCmd() *cobra.Command {
	var unacked bool
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Show alerts (deprecated: see `logs`)",
		Long:  "Show legacy alerts. This command is deprecated; use gosched logs for the unified log view.",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(cmd.ErrOrStderr(), "gosched: alerts is deprecated; use gosched logs for the unified log view. alerts will be removed in a future release.")
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			alerts, err := clientFor(cmd).ListAlerts(ctx, unacked)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), alerts)
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tSEVERITY\tKIND\tMESSAGE\tACK")
			for _, a := range alerts {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\n", a.ID, a.Severity, a.Kind, a.Message, a.Acknowledged)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().BoolVar(&unacked, "unacked", false, "show only unacknowledged alerts")

	ack := &cobra.Command{
		Use: "ack <id>", Short: "Acknowledge an alert", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			if err := clientFor(cmd).AckAlert(ctx, a[0]); err != nil {
				return err
			}
			return printAcknowledgement(cmd, map[string]string{"id": a[0], "status": "acknowledged"}, "acknowledged alert %s\n", a[0])
		},
	}
	cmd.AddCommand(ack)
	return cmd
}

func newCancelRunCmd() *cobra.Command {
	return &cobra.Command{Use: "cancel RUN_ID", Short: "Request cancellation of one active execution by its exact Run ID", Long: "Request cancellation of one active run on the selected scheduler. Find its Run ID with gosched runs active. A requested acknowledgement does not establish termination; inspect gosched runs get RUN_ID for the final outcome.", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		result, err := clientFor(cmd).CancelRun(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), result)
		}
		if result.Status == "already_finished" {
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "run %s already finished (outcome %s)\n", result.RunID, result.Outcome)
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "cancellation requested for run %s; inspect its final outcome with gosched runs get %s\n", result.RunID, result.RunID)
		return err
	}}
}
