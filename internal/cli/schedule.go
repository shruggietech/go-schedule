package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func newScheduleCmd() *cobra.Command {
	root := &cobra.Command{Use: "schedule", Short: "Preview schedules without changing tasks"}
	var syntax, tz, missingDate, timeBasis, dstGap, dstOverlap string
	cmd := &cobra.Command{
		Use: "preview <expression>", Short: "Predict human or cron occurrences using the complete timing policy", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(args[0]) == "" {
				return fmtUsage("provide a human-readable schedule or a supported cron expression")
			}
			hint, err := scheduleSyntax(syntax)
			if err != nil {
				return err
			}
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			response, err := clientFor(cmd).Preview(ctx, server.PreviewRequest{
				Schedule: args[0], ScheduleSyntax: hint, Timezone: tz, MissingDatePolicy: missingDate,
				TimeBasis: timeBasis, DSTGapPolicy: dstGap, DSTOverlapPolicy: dstOverlap,
			})
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), response)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Predicted schedule: %s\nsyntax: %s\ntimezone: %s\n", response.HumanSummary, response.SourceSyntax, tz); err != nil {
				return err
			}
			printPolicySummary(cmd, response.PolicySummary)
			printNextRuns(cmd, response.NextRuns)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&syntax, "syntax", "auto", "schedule interpretation: auto|human|cron")
	f.StringVar(&tz, "tz", "Local", "IANA timezone for the preview")
	f.StringVar(&missingDate, "missing-date", "", missingDateUsage)
	f.StringVar(&timeBasis, "time-basis", "", timeBasisUsage)
	f.StringVar(&dstGap, "dst-gap", "", dstGapUsage)
	f.StringVar(&dstOverlap, "dst-overlap", "", dstOverlapUsage)
	root.AddCommand(cmd)
	return root
}
