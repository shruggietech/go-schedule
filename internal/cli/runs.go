package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func newRunsCmd() *cobra.Command {
	var task string
	var limit int
	cmd := &cobra.Command{
		Use:   "runs",
		Short: "Show run history",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			runs, err := newClient().ListRuns(ctx, task, limit)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(runs)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "SCHEDULED\tOUTCOME\tTRIGGER\tSOURCE TASK\tSOURCE RUN\tEXIT")
			for _, r := range runs {
				exit := "-"
				if r.ExitCode != nil {
					exit = fmt.Sprintf("%d", *r.ExitCode)
				}
				sourceTask, sourceRun := runSourceColumns(r)
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", r.ScheduledFor.Format(time.RFC3339), r.Outcome, r.Trigger, sourceTask, sourceRun, exit)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "filter by task ID")
	cmd.Flags().IntVar(&limit, "limit", 50, "maximum rows")
	return cmd
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
		Use:        "alerts",
		Short:      "Show alerts (deprecated: see `logs`)",
		Deprecated: "use `gosched logs` for the unified log view; `alerts` will be removed in a future release.",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			alerts, err := newClient().ListAlerts(ctx, unacked)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(alerts)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
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
		RunE: func(_ *cobra.Command, a []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			if err := newClient().AckAlert(ctx, a[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "acknowledged alert %s\n", a[0])
			return nil
		},
	}
	cmd.AddCommand(ack)
	return cmd
}
