package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func newLogsCmd() *cobra.Command {
	var severity string
	var limit int
	var metadata bool
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show recent daemon logs (info/warning/error)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch severity {
			case "", "info", "warning", "error":
			default:
				return fmtUsage("--severity must be one of: info, warning, error")
			}
			if limit < 0 {
				return fmtUsage("--limit must be a non-negative number of log records")
			}
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			response, err := clientFor(cmd).ListLogs(ctx, severity, limit)
			if err != nil {
				return err
			}
			if metadata {
				return writeLogsMetadata(cmd.OutOrStdout(), response, jsonOutput(cmd))
			}
			return writeLogs(cmd.OutOrStdout(), response.Logs, jsonOutput(cmd))
		},
	}
	cmd.Flags().StringVar(&severity, "severity", "", "filter by severity: info, warning, or error")
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum rows (0 shows all retained log records)")
	cmd.Flags().BoolVar(&metadata, "metadata", false, "include the daemon log path with the log records")
	return cmd
}

func writeLogsMetadata(w io.Writer, response server.LogsResponse, asJSON bool) error {
	if asJSON {
		return printJSONTo(w, response)
	}
	if _, err := fmt.Fprintf(w, "Log path: %s\n", displayOrDash(response.LogPath)); err != nil {
		return err
	}
	return writeLogs(w, response.Logs, false)
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
