package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/clientobserve"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func newSearchCmd() *cobra.Command {
	var all bool
	var kinds []string
	var limit int
	cmd := &cobra.Command{Use: "search QUERY", Short: "Search the selected scheduler or all registrations without changing the selected target", Args: cobra.ExactArgs(1), Annotations: map[string]string{"cli.aggregate": "true"}, RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(args[0])
		if query == "" || utf8.RuneCountInString(query) > 200 {
			return fmtUsage("search query must contain between 1 and 200 characters")
		}
		if limit < 1 || limit > 50 {
			return fmtUsage("--limit must be between 1 and 50 results per scheduler")
		}
		request := clientobserve.SearchRequest{Query: query, Limit: limit}
		for _, value := range kinds {
			kind := domain.SearchKind(value)
			if !kind.Valid() {
				return fmtUsage("--kind must be task, group, failure, schedule, or alert")
			}
			request.Kinds = append(request.Kinds, kind)
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		targets, err := observationTargets(ctx, cmd, all)
		if err != nil {
			return err
		}
		observer := clientobserve.New()
		observer.Now = stateFor(cmd).deps.Now
		result, observeErr := observer.Search(ctx, targets, request)
		if jsonOutput(cmd) {
			err = printJSONTo(cmd.OutOrStdout(), result)
		} else {
			err = writeSearch(cmd.OutOrStdout(), result)
		}
		if err != nil {
			return err
		}
		return observeErr
	}}
	cmd.Flags().BoolVar(&all, "all", false, "search the local scheduler and every saved profile with independent outcomes")
	cmd.Flags().StringSliceVar(&kinds, "kind", nil, "filter kinds: task, group, failure, schedule, alert (repeat or separate with commas)")
	cmd.Flags().IntVar(&limit, "limit", 50, "maximum results per scheduler (1 through 50)")
	return cmd
}

func writeSearch(w io.Writer, snapshot clientobserve.Snapshot) error {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "SOURCE\tREGISTRATION\tEXPECTED DAEMON\tOBSERVED DAEMON\tENDPOINT\tOBSERVED AT\tKIND\tOBJECT ID\tNAME\tCONTEXT\tTRUNCATED\tOUTCOME")
	for _, item := range snapshot.Observations {
		r := item.Registration
		prefix := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", r.Label, r.Key, displayOrDash(r.ExpectedDaemonID), displayOrDash(r.ObservedDaemonID), r.Endpoint, item.ObservedAt.Format("2006-01-02T15:04:05Z07:00"))
		if item.Failure != nil {
			fmt.Fprintf(tw, "%s\t-\t-\t-\t-\t-\t%s: %s\n", prefix, item.Failure.Code, item.Failure.Message)
			continue
		}
		if len(item.Search.Results) == 0 {
			fmt.Fprintf(tw, "%s\t-\t-\t-\t-\t%t\tno matches\n", prefix, item.Search.Truncated)
			continue
		}
		for _, match := range item.Search.Results {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%t\tok\n", prefix, match.Kind, match.ObjectID, match.Name, match.Context, item.Search.Truncated)
		}
	}
	return tw.Flush()
}
