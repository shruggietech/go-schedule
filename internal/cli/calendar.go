package cli

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func newCalendarCmd() *cobra.Command {
	var fromText, toText string
	cmd := &cobra.Command{Use: "calendar", Short: "Show recorded runs and predicted occurrences within a bounded time window", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		from, to, err := calendarWindow(stateFor(cmd).deps.Now(), fromText, toText)
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		result, err := clientFor(cmd).GetCalendar(ctx, from, to)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), result)
		}
		return writeCalendar(cmd.OutOrStdout(), result)
	}}
	cmd.Flags().StringVar(&fromText, "from", "", "start timestamp with timezone (RFC3339; default now)")
	cmd.Flags().StringVar(&toText, "to", "", "end timestamp with timezone (RFC3339; default seven days after start; maximum 30 days)")
	return cmd
}

func calendarWindow(now time.Time, fromText, toText string) (time.Time, time.Time, error) {
	from := now.UTC()
	var err error
	if fromText != "" {
		from, err = time.Parse(time.RFC3339, fromText)
		if err != nil {
			return time.Time{}, time.Time{}, fmtUsage("--from must be an RFC3339 timestamp with a timezone, such as 2026-10-19T12:00:00Z")
		}
	}
	to := from.Add(7 * 24 * time.Hour)
	if toText != "" {
		to, err = time.Parse(time.RFC3339, toText)
		if err != nil {
			return time.Time{}, time.Time{}, fmtUsage("--to must be an RFC3339 timestamp with a timezone, such as 2026-10-26T12:00:00Z")
		}
	}
	if !to.After(from) || to.Sub(from) > 30*24*time.Hour {
		return time.Time{}, time.Time{}, fmtUsage("calendar end must follow start by no more than 30 days")
	}
	return from.UTC(), to.UTC(), nil
}

func writeCalendar(w io.Writer, response server.CalendarResponse) error {
	items := append([]server.Occurrence(nil), response.Occurrences...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Time.Before(items[j].Time) })
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tKIND\tTASK ID\tTASK NAME\tRUN ID\tOUTCOME")
	for _, item := range items {
		kind := "predicted"
		if item.Kind == "past" || item.RunID != "" {
			kind = "recorded"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", item.Time.Format(time.RFC3339), kind, item.TaskID, item.TaskName, displayOrDash(item.RunID), displayOrDash(string(item.Outcome)))
	}
	return tw.Flush()
}

func displayOrDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
