package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func newLogsCmd() *cobra.Command {
	var severity string
	var limit int
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show recent daemon logs (info/warning/error)",
		RunE: func(_ *cobra.Command, _ []string) error {
			switch severity {
			case "", "info", "warning", "error":
			default:
				return fmt.Errorf("--severity must be one of: info, warning, error")
			}
			ctx, cancel := reqCtx()
			defer cancel()
			logs, err := newClient().ListLogs(ctx, severity, limit)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(logs)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "TIME\tSEVERITY\tSOURCE\tMESSAGE")
			for _, l := range logs {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
					l.Time.Format(time.RFC3339), l.Severity, l.Source, l.Message)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().StringVar(&severity, "severity", "", "filter by severity: info, warning, or error")
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum rows")
	return cmd
}
