package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/cron"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// The cron command group converts strings, reads crontabs at import, and writes
// them at export. Imported expressions use the API's explicit cron input
// boundary so their editable source survives while execution remains RRULE
// based.

func newCronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Convert to and from crontab format",
		Long: "Convert between crontab expressions and this scheduler's schedules.\n\n" +
			"Cron is supported for local string conversion, task input, and interchange:\n" +
			"expressions can be converted, imported, explained, exported, and supplied to\n" +
			"task add or edit when they fit the supported cron subset.",
	}
	cmd.AddCommand(cronConvert(), cronExplain(), cronImport(), cronExport())
	return cmd
}

// ---- convert ------------------------------------------------------------

func cronConvert() *cobra.Command {
	var destination string
	cmd := &cobra.Command{
		Use:           "convert [--to cron|human] <schedule-string>",
		Short:         "Convert one cron or human schedule string locally",
		SilenceErrors: true,
		SilenceUsage:  true,
		Long: "Convert one schedule string to the opposite syntax without contacting the daemon.\n\n" +
			"Automatic mode treats @-prefixed values and five or six cron-shaped fields\n" +
			"as cron. Use --to to select the output syntax explicitly.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return fmtUsage(err.Error())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := cron.Convert(args[0], cron.Syntax(destination))
			if err != nil {
				return fmtUsage(err.Error())
			}
			if result.RefusalReason != "" {
				if jsonOut {
					if err := printJSONTo(cmd.ErrOrStderr(), result); err != nil {
						return fmt.Errorf("write conversion refusal: %w", err)
					}
					return reported(fmtUsage(result.RefusalReason))
				}
				return fmtUsage(result.RefusalReason)
			}
			if jsonOut {
				return printJSONTo(cmd.OutOrStdout(), result)
			}
			fmt.Fprintln(cmd.OutOrStdout(), result.Output)
			return nil
		},
	}
	cmd.Flags().StringVar(&destination, "to", "", "output syntax: cron or human (default: detect input)")
	return cmd
}

// ---- explain ------------------------------------------------------------

func cronExplain() *cobra.Command {
	var tz string
	var count int
	cmd := &cobra.Command{
		Use:   "explain <expression>",
		Short: "Translate one cron expression into plain language",
		Long: "Print the plain-language phrase a cron expression maps to, plus its next\n" +
			"run times. Nothing is created or changed.\n\n" +
			"An expression that cannot be represented is reported by name — that is an\n" +
			"answer, not a failure, so the exit code stays 0.",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			expr := a[0]
			phrase, bad, err := cron.Explain(expr)
			if err != nil {
				return fmt.Errorf("%w: %v", errUsage, err)
			}
			zone := orLocal(tz)
			out := explainResult{Expression: strings.TrimSpace(expr), Timezone: zone}
			if bad.Reason != "" {
				out.Unsupported = bad.Reason
			} else {
				out.Phrase = phrase
				runs, rerr := previewRuns(expr, string(cron.SyntaxCron), zone, count)
				if rerr != nil {
					return rerr
				}
				out.NextRuns = runs
			}
			if jsonOut {
				return printJSON(out)
			}
			printExplain(os.Stdout, out)
			return nil
		},
	}
	cmd.Flags().StringVar(&tz, "timezone", "", "IANA timezone for the displayed run times (default: task default)")
	cmd.Flags().IntVar(&count, "count", 3, "how many upcoming runs to show")
	return cmd
}

// explainResult is the machine-readable shape of an explanation. Exactly one of
// Phrase or Unsupported is set.
type explainResult struct {
	Expression  string      `json:"expression"`
	Phrase      string      `json:"phrase,omitempty"`
	Unsupported string      `json:"unsupported,omitempty"`
	Timezone    string      `json:"timezone"`
	NextRuns    []time.Time `json:"next_runs,omitempty"`
}

func printExplain(w io.Writer, r explainResult) {
	fmt.Fprintln(w, r.Expression)
	if r.Unsupported != "" {
		fmt.Fprintf(w, "  unsupported: %s\n", r.Unsupported)
		return
	}
	fmt.Fprintf(w, "  phrase: %s\n", r.Phrase)
	for i, t := range r.NextRuns {
		label := "  next:  "
		if i > 0 {
			label = "         "
		}
		fmt.Fprintf(w, "%s %s\n", label, t.Format(time.RFC3339))
	}
}

// previewRuns asks the daemon what an input resolves to, so the times shown come
// from the same evaluator that will run the task rather than a second one here.
func previewRuns(input, syntax, tz string, count int) ([]time.Time, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	resp, err := newClient().Preview(ctx, server.PreviewRequest{
		Schedule: input, ScheduleSyntax: syntax, Timezone: tz,
	})
	if err != nil {
		return nil, err
	}
	if count > 0 && len(resp.NextRuns) > count {
		return resp.NextRuns[:count], nil
	}
	return resp.NextRuns, nil
}

func orLocal(tz string) string {
	if tz == "" {
		return "Local"
	}
	return tz
}

// ---- import -------------------------------------------------------------

// taskCreator is the slice of the API client the import needs, so the reporting
// logic can be exercised without a running daemon.
type taskCreator interface {
	CreateTask(ctx context.Context, req server.TaskCreateRequest) (server.TaskResponse, error)
}

type importOptions struct {
	dryRun   bool
	timezone string
	runAs    string
	group    string
	count    int
	// runs resolves a typed input to its upcoming run times. It is a field so the
	// reporting can be exercised without a daemon, and so an unreachable daemon
	// degrades to a report without run times rather than to no report at all.
	runs func(input, syntax, tz string, count int) ([]time.Time, error)
}

func cronImport() *cobra.Command {
	var opts importOptions
	var file string
	var dialect string
	var system bool
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a crontab, creating one task per schedule line",
		Long: "Read a crontab and create a task for every line that can be represented.\n\n" +
			"Each line is reported with the phrase it maps to; lines that cannot be\n" +
			"represented are reported by name rather than dropped. Use --dry-run to see\n" +
			"the whole report without creating anything.",
		RunE: func(_ *cobra.Command, _ []string) error {
			if file == "" {
				return fmt.Errorf("%w: --file is required (use - for standard input)", errUsage)
			}
			if system && opts.runAs != "" {
				return fmt.Errorf("%w: --run-as cannot be combined with --system", errUsage)
			}
			r, closeFn, err := openInput(file)
			if err != nil {
				return fmt.Errorf("%w: %v", errUsage, err)
			}
			defer closeFn()

			rep, err := cron.ScanCrontabWithOptions(r, cron.ScanOptions{Dialect: cron.Dialect(dialect), System: system})
			if err != nil {
				return err
			}
			var creator taskCreator
			if !opts.dryRun {
				creator = newClient()
			}
			opts.runs = previewRuns
			return runImport(os.Stdout, &rep, opts, creator)
		},
	}
	f := cmd.Flags()
	f.StringVar(&file, "file", "", "crontab file to read, or - for standard input (required)")
	f.StringVar(&dialect, "dialect", string(cron.DialectUnix), "timing layout: unix (five fields) or quartz (six fields)")
	f.BoolVar(&system, "system", false, "read a system crontab with a user field after timing")
	f.StringVar(&opts.runAs, "run-as", "", "run every user-crontab job as this account")
	f.BoolVar(&opts.dryRun, "dry-run", false, "print the report without creating anything")
	f.StringVar(&opts.timezone, "timezone", "", "IANA timezone for the created tasks (default: task default)")
	f.StringVar(&opts.group, "group", "", "group ID to place the imported tasks in")
	f.IntVar(&opts.count, "count", 3, "how many upcoming runs to show per line")
	return cmd
}

func openInput(path string) (io.Reader, func(), error) {
	if path == "-" {
		return os.Stdin, func() {}, nil
	}
	f, err := os.Open(path) //nolint:gosec // the path is the operator's own argument
	if err != nil {
		return nil, nil, err
	}
	return f, func() { _ = f.Close() }, nil
}

// runImport prints the per-line report, creates the tasks unless creator is nil,
// and prints the summary. It is the whole of the import's behavior, with the
// daemon behind a one-method interface so the reporting can be tested directly.
//
// A declined line never stops the run: the supported lines are still created and
// the summary accounts for every line (FR-005a, FR-010a).
func runImport(w io.Writer, rep *cron.Report, opts importOptions, creator taskCreator) error {
	for _, line := range rep.Lines {
		for _, warn := range line.Warnings {
			fmt.Fprintf(w, "! %s\n", warn)
		}
		switch line.Kind {
		case cron.LineSkipped:
			continue
		case cron.LineError:
			fmt.Fprintf(w, "line %d: %s\n  error: %s\n", line.Number, line.Raw, line.Reason)
		case cron.LineDeclined:
			fmt.Fprintf(w, "line %d: %s\n  unsupported: %s\n", line.Number, line.Expr, line.Reason)
		case cron.LineJob:
			zone := importZone(line, opts.timezone)
			fmt.Fprintf(w, "line %d: %s\n  phrase:  %s\n  timezone: %s\n  command: %s\n",
				line.Number, line.Expr, line.Phrase, zone, commandLine(line))
			if runAs := importRunAs(line, opts.runAs); runAs != "" {
				fmt.Fprintf(w, "  run as:   %s\n", runAs)
			}
			if len(line.Env) > 0 {
				fmt.Fprintf(w, "  env:      %d variable(s)\n", len(line.Env))
			}
			if line.HasStdin {
				fmt.Fprintf(w, "  stdin:    %d byte(s)\n", len(line.Stdin))
			}
			printLineRuns(w, line.Expr, zone, opts)
			if creator == nil {
				continue
			}
			if err := createFromLine(w, creator, line, zone, opts.group, opts.runAs, rep); err != nil {
				return err
			}
		}
	}

	printImportSummary(w, rep, opts.timezone, opts.dryRun)
	if rep.Failed > 0 {
		return fmt.Errorf("%d of %d task(s) could not be created; the rest were", rep.Failed, rep.Jobs)
	}
	return nil
}

func importRunAs(line cron.Line, fallback string) string {
	if line.RunAs != "" {
		return line.RunAs
	}
	return fallback
}

func importZone(line cron.Line, override string) string {
	if override != "" {
		return override
	}
	return orLocal(line.Timezone)
}

// printLineRuns shows when a line would actually fire, which is the half of the
// report that catches a misreading: a phrase can look right and still mean
// something else. An unreachable daemon costs the run times, not the report —
// the conversion itself is local, and a preview that refused to print because
// the daemon was down would be useless exactly when it is most wanted.
func printLineRuns(w io.Writer, expression, zone string, opts importOptions) {
	if opts.runs == nil || opts.count <= 0 {
		return
	}
	runs, err := opts.runs(expression, string(cron.SyntaxCron), zone, opts.count)
	if err != nil {
		fmt.Fprintf(w, "  next:    (unavailable: %v)\n", err)
		return
	}
	for i, ts := range runs {
		label := "  next:   "
		if i > 0 {
			label = "          "
		}
		fmt.Fprintf(w, "%s%s\n", label, ts.Format(time.RFC3339))
	}
}

func createFromLine(w io.Writer, creator taskCreator, line cron.Line, zone, group, runAs string, rep *cron.Report) error {
	ctx, cancel := reqCtx()
	defer cancel()
	resp, err := creator.CreateTask(ctx, server.TaskCreateRequest{
		Name:           importName(line),
		Command:        line.Command,
		Args:           line.Args,
		Env:            line.Env,
		Stdin:          line.Stdin,
		RunAs:          importRunAs(line, runAs),
		GroupID:        group,
		Timezone:       zone,
		Schedule:       line.Expr,
		ScheduleSyntax: string(cron.SyntaxCron),
	})
	if err != nil {
		rep.Failed++
		fmt.Fprintf(w, "  not created: %v\n", err)
		return nil
	}
	rep.Created++
	fmt.Fprintf(w, "  created: %s (%s)\n", resp.Task.ID, resp.Task.Name)
	return nil
}

// importName derives a task name from the command, because crontabs do not name
// jobs. The base name of the program is what an operator would recognize in a
// task list; the daemon does not require names to be unique.
func importName(line cron.Line) string {
	nameSource := line.CommandText
	if fields := strings.Fields(nameSource); len(fields) > 0 {
		nameSource = fields[0]
	}
	base := filepath.Base(nameSource)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return fmt.Sprintf("imported line %d", line.Number)
	}
	return base
}

func commandLine(line cron.Line) string {
	if len(line.Args) == 0 {
		return line.Command
	}
	return line.Command + " " + strings.Join(line.Args, " ")
}

// printImportSummary states the counts and the fidelity facts. The fidelity
// paragraph is not decoration: file timezone context is explicit, while cron
// still carries no catch-up, overlap policy, or restart recovery. The operator
// must be told what the imported jobs gained and what defaults were selected.
func printImportSummary(w io.Writer, rep *cron.Report, timezoneOverride string, dryRun bool) {
	fmt.Fprintln(w)
	verb := "created"
	if dryRun {
		verb = "would create"
	}
	fmt.Fprintf(w, "%d line(s) read: %d %s, %d skipped (comments, blanks, variables), %d unsupported, %d error(s)\n",
		rep.Read, countCreated(rep, dryRun), verb, rep.Skipped, rep.Declined, rep.Errors)
	if rep.Failed > 0 {
		fmt.Fprintf(w, "%d task(s) failed to create; those already created were kept\n", rep.Failed)
	}
	fmt.Fprintln(w)
	if timezoneOverride != "" {
		fmt.Fprintf(w, "The explicit timezone override applies to every imported task: %s.\n", timezoneOverride)
	} else {
		fmt.Fprintln(w, "Each task uses its effective CRON_TZ value, or Local when none is set.")
	}
	fmt.Fprintf(w, "Cron also has no catch-up, overlap, or restart recovery. Imported tasks take\n"+
		"the defaults: catch-up %q (one run after downtime), overlap %q, missing dates %q.\n",
		domain.CatchupOne, domain.OverlapQueueOne, domain.MissingDateSkip)
	if dryRun {
		fmt.Fprintln(w, "\nThis was a preview. Re-run without --dry-run to create these tasks.")
	}
}

func countCreated(rep *cron.Report, dryRun bool) int {
	if dryRun {
		return rep.Jobs
	}
	return rep.Created
}

// ---- export -------------------------------------------------------------

func cronExport() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Emit the task set as crontab lines",
		Long: "Print a crontab line for every task whose schedule cron can carry, and a\n" +
			"commented refusal naming every task it cannot. Nothing is approximated and\n" +
			"no task is silently omitted.",
		RunE: func(_ *cobra.Command, _ []string) error {
			details, err := exportTargets(taskID)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(exportLines(details))
			}
			printExport(os.Stdout, details)
			return nil
		},
	}
	cmd.Flags().StringVar(&taskID, "task", "", "export a single task by ID")
	return cmd
}

func exportTargets(taskID string) ([]server.TaskResponse, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	c := newClient()
	if taskID != "" {
		resp, err := c.GetTask(ctx, taskID)
		if err != nil {
			return nil, err
		}
		return []server.TaskResponse{resp}, nil
	}
	tasks, err := c.ListTasks(ctx, "", "")
	if err != nil {
		return nil, err
	}
	out := make([]server.TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		detail, err := c.GetTask(ctx, t.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, detail)
	}
	return out, nil
}

// exportLine is one task's export outcome; exactly one of Line or Declined is
// set, and every task produces one of these.
type exportLine struct {
	TaskID   string `json:"task_id"`
	Name     string `json:"name"`
	Line     string `json:"line,omitempty"`
	Declined string `json:"declined,omitempty"`
}

func exportLines(details []server.TaskResponse) []exportLine {
	out := make([]exportLine, 0, len(details))
	for _, d := range details {
		e := exportLine{TaskID: d.Task.ID, Name: d.Task.Name}
		expr, bad, ok := cron.Export(d.Task, d.Schedule)
		if ok {
			e.Line = expr + " " + commandOf(d.Task)
		} else {
			e.Declined = bad.Reason
		}
		out = append(out, e)
	}
	return out
}

func printExport(w io.Writer, details []server.TaskResponse) {
	fmt.Fprintf(w, "# gosched cron export — %d task(s)\n", len(details))
	for _, e := range exportLines(details) {
		if e.Declined != "" {
			fmt.Fprintf(w, "# declined: %q — %s\n", e.Name, e.Declined)
			continue
		}
		fmt.Fprintln(w, e.Line)
	}
}

func commandOf(t domain.Task) string {
	if len(t.Args) == 0 {
		return t.Command
	}
	return t.Command + " " + strings.Join(t.Args, " ")
}
