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

// The cron command group is the project's only cron surface. Cron is an
// interchange format here — read at import, written at export — and never an
// authoring syntax: nothing in this file feeds an expression anywhere a phrase
// is accepted. Every conversion goes through the human phrase, so what a preview
// prints is literally what gets parsed and stored.

func newCronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Convert to and from crontab format (import, explain, export)",
		Long: "Convert between crontab expressions and this scheduler's schedules.\n\n" +
			"Cron is supported as an interchange format only: expressions can be imported,\n" +
			"explained, and exported, but are never accepted where a schedule phrase is.",
	}
	cmd.AddCommand(cronExplain(), cronImport(), cronExport())
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
				runs, rerr := previewRuns(phrase, zone, count)
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

// previewRuns asks the daemon what a phrase resolves to, so the times shown come
// from the same evaluator that will run the task rather than a second one here.
func previewRuns(phrase, tz string, count int) ([]time.Time, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	resp, err := newClient().Preview(ctx, server.PreviewRequest{Schedule: phrase, Timezone: tz})
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
	group    string
	count    int
	// runs resolves a phrase to its upcoming run times. It is a field so the
	// reporting can be exercised without a daemon, and so an unreachable daemon
	// degrades to a report without run times rather than to no report at all.
	runs func(phrase, tz string, count int) ([]time.Time, error)
}

func cronImport() *cobra.Command {
	var opts importOptions
	var file string
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
			r, closeFn, err := openInput(file)
			if err != nil {
				return fmt.Errorf("%w: %v", errUsage, err)
			}
			defer closeFn()

			rep, err := cron.ScanCrontab(r)
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
	zone := orLocal(opts.timezone)

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
			fmt.Fprintf(w, "line %d: %s\n  phrase:  %s\n  command: %s\n",
				line.Number, line.Expr, line.Phrase, commandLine(line))
			printLineRuns(w, line.Phrase, zone, opts)
			if creator == nil {
				continue
			}
			if err := createFromLine(w, creator, line, zone, opts.group, rep); err != nil {
				return err
			}
		}
	}

	printImportSummary(w, rep, zone, opts.dryRun)
	if rep.Failed > 0 {
		return fmt.Errorf("%d of %d task(s) could not be created; the rest were", rep.Failed, rep.Jobs)
	}
	return nil
}

// printLineRuns shows when a line would actually fire, which is the half of the
// report that catches a misreading: a phrase can look right and still mean
// something else. An unreachable daemon costs the run times, not the report —
// the conversion itself is local, and a preview that refused to print because
// the daemon was down would be useless exactly when it is most wanted.
func printLineRuns(w io.Writer, phrase, zone string, opts importOptions) {
	if opts.runs == nil || opts.count <= 0 {
		return
	}
	runs, err := opts.runs(phrase, zone, opts.count)
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

func createFromLine(w io.Writer, creator taskCreator, line cron.Line, zone, group string, rep *cron.Report) error {
	ctx, cancel := reqCtx()
	defer cancel()
	resp, err := creator.CreateTask(ctx, server.TaskCreateRequest{
		Name:     importName(line),
		Command:  line.Command,
		Args:     line.Args,
		GroupID:  group,
		Timezone: zone,
		Schedule: line.Phrase,
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
	base := filepath.Base(line.Command)
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
// paragraph is not decoration: cron carries no timezone, no catch-up, no overlap
// policy and no restart recovery, so an operator has to be told what their jobs
// just gained and what was assumed on their behalf (FR-008, FR-009).
func printImportSummary(w io.Writer, rep *cron.Report, zone string, dryRun bool) {
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
	fmt.Fprintf(w, "Cron carries no timezone, so these tasks use %s.\n", zone)
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
