package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-scheduler/internal/api/server"
)

func newTriggerCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "trigger", Short: "Chain tasks on completion (event triggers)"}
	cmd.AddCommand(triggerAdd(), triggerList(), triggerRm())
	return cmd
}

func triggerAdd() *cobra.Command {
	var source, target, on, dedupKey, window string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Run --target when --source completes",
		RunE: func(_ *cobra.Command, _ []string) error {
			if source == "" || target == "" {
				return fmt.Errorf("%w: --source and --target are required", errUsage)
			}
			ctx, cancel := reqCtx()
			defer cancel()
			tr, err := newClient().CreateTrigger(ctx, server.TriggerCreateRequest{
				SourceTaskID: source, TargetTaskID: target, OnOutcome: on,
				DedupKey: dedupKey, DedupWindow: window,
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(tr)
			}
			fmt.Fprintf(os.Stdout, "created trigger %s: %s -> %s (on %s)\n", tr.ID, tr.SourceTaskID, tr.TargetTaskID, tr.OnOutcome)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&source, "source", "", "source task ID (its completion fires the trigger)")
	f.StringVar(&target, "target", "", "target task ID (runs when source completes)")
	f.StringVar(&on, "on", "success", "fire on outcome: success|failure|any")
	f.StringVar(&dedupKey, "dedup-key", "", "fixed dedup key (default: per source run)")
	f.StringVar(&window, "dedup-window", "", "dedup window, e.g. 5m")
	return cmd
}

func triggerList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List triggers",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			triggers, err := newClient().ListTriggers(ctx)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(triggers)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tSOURCE\tTARGET\tON\tWINDOW")
			for _, t := range triggers {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.SourceTaskID, t.TargetTaskID, t.OnOutcome, t.DedupWindow)
			}
			return tw.Flush()
		},
	}
}

func triggerRm() *cobra.Command {
	return &cobra.Command{
		Use: "rm <id>", Short: "Delete a trigger", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			if err := newClient().DeleteTrigger(ctx, a[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "deleted trigger %s\n", a[0])
			return nil
		},
	}
}
