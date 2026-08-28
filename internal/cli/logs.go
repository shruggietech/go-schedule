package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func newLogsCmd() *cobra.Command {
	var severity string
	var limit int
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show recent daemon logs (info/warning/error)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch severity {
			case "", "info", "warning", "error":
			default:
				return fmt.Errorf("--severity must be one of: info, warning, error")
			}
			ctx, cancel := reqCtx()
			defer cancel()
			response, err := newClient().ListLogs(ctx, severity, limit)
			if err != nil {
				return err
			}
			return writeLogs(cmd.OutOrStdout(), response.Logs, jsonOut)
		},
	}
	cmd.Flags().StringVar(&severity, "severity", "", "filter by severity: info, warning, or error")
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum rows")
	return cmd
}

func writeLogs(w io.Writer, logs []domain.LogRecord, asJSON bool) error {
	if asJSON {
		return printJSONTo(w, logs)
	}
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tSEVERITY\tSOURCE\tMESSAGE")
	for _, l := range logs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			l.Time.Format(time.RFC3339), l.Severity, l.Source, l.Message)
	}
	return tw.Flush()
}
