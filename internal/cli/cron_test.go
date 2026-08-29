package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/cron"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestCronConvert_DefaultTextIsOneLocalLine(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  string
	}{
		{input: "0 9 * * 1-5", want: "weekdays at 09:00\n"},
		{input: "weekdays at 09:00", want: "0 9 * * 1-5\n"},
		{input: "0 9 * * 5#3", want: "3rd friday monthly at 09:00\n"},
		{input: "3rd friday monthly at 09:00", want: "0 9 * * 5#3\n"},
		{input: "0 9 * * 5L", want: "last friday of the month at 09:00\n"},
		{input: "last friday of the month at 09:00", want: "0 9 * * 5L\n"},
		{input: "0 9 L * *", want: "last day of every month at 09:00\n"},
		{input: "last day of every month at 09:00", want: "0 9 L * *\n"},
		{input: "0 9 15W * *", want: "nearest weekday to the 15th of every month at 09:00\n"},
		{input: "nearest weekday to the 15th of every month at 09:00", want: "0 9 15W * *\n"},
		{input: "0 9 LW * *", want: "last weekday of every month at 09:00\n"},
		{input: "last weekday of every month at 09:00", want: "0 9 LW * *\n"},
		{input: "*/10 9-17 * * MON,WED,FRI", want: "every 10 minutes during hours 9 through 17 on Monday, Wednesday, and Friday\n"},
		{input: "*/7 * * * *", want: "at minutes 0, 7, 14, 21, 28, 35, 42, 49, and 56 of each hour every day\n"},
	} {
		cmd := cronConvert()
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		cmd.SetArgs([]string{tt.input})

		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
		if got := stdout.String(); got != tt.want {
			t.Errorf("stdout = %q, want %q", got, tt.want)
		}
		if stderr.Len() != 0 {
			t.Errorf("stderr = %q, want empty", stderr.String())
		}
	}
}

func TestImport_MonthlyCalendarSelectorsRemainCronSource(t *testing.T) {
	for _, tt := range []struct{ expr, phrase string }{
		{"0 9 L * *", "last day of every month at 09:00"},
		{"0 9 15W * *", "nearest weekday to the 15th of every month at 09:00"},
		{"0 9 LW * *", "last weekday of every month at 09:00"},
	} {
		rep := scan(t, tt.expr+" /usr/local/bin/report\n")
		if rep.Jobs != 1 || rep.Lines[0].Phrase != tt.phrase {
			t.Fatalf("classification %q = %+v", tt.expr, rep)
		}
		fc := &fakeCreator{}
		var out bytes.Buffer
		if err := runImport(&out, rep, importOptions{timezone: "UTC"}, fc); err != nil {
			t.Fatal(err)
		}
		if len(fc.reqs) != 1 || fc.reqs[0].Schedule != tt.expr || fc.reqs[0].ScheduleSyntax != "cron" {
			t.Fatalf("import %q = %+v", tt.expr, fc.reqs)
		}
	}
}

func TestImport_LastWeekdayRemainsCronSource(t *testing.T) {
	rep := scan(t, "0 9 * * 5L /usr/local/bin/report\n")
	if rep.Jobs != 1 || len(rep.Lines) != 1 || rep.Lines[0].Phrase != "last friday of the month at 09:00" {
		t.Fatalf("last-weekday line classification = %+v", rep)
	}
	fc := &fakeCreator{}
	var out bytes.Buffer
	if err := runImport(&out, rep, importOptions{timezone: "UTC"}, fc); err != nil {
		t.Fatal(err)
	}
	if len(fc.reqs) != 1 || fc.reqs[0].Schedule != "0 9 * * 5L" || fc.reqs[0].ScheduleSyntax != "cron" {
		t.Fatalf("last-weekday import request = %+v", fc.reqs)
	}
}

func TestImport_OrdinalWeekdayRemainsCronSource(t *testing.T) {
	rep := scan(t, "0 9 * * 5#3 /usr/local/bin/report\n")
	if rep.Jobs != 1 || len(rep.Lines) != 1 || rep.Lines[0].Phrase != "3rd friday monthly at 09:00" {
		t.Fatalf("ordinal line classification = %+v", rep)
	}
	fc := &fakeCreator{}
	var out bytes.Buffer
	if err := runImport(&out, rep, importOptions{timezone: "UTC"}, fc); err != nil {
		t.Fatal(err)
	}
	if len(fc.reqs) != 1 || fc.reqs[0].Schedule != "0 9 * * 5#3" || fc.reqs[0].ScheduleSyntax != "cron" {
		t.Fatalf("ordinal import request = %+v", fc.reqs)
	}
}

func TestCronConvert_ForcedTextAndInvalidDestination(t *testing.T) {
	cmd := cronConvert()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--to", "human", "0 9 * * 1-5"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := stdout.String(); got != "weekdays at 09:00\n" {
		t.Errorf("stdout = %q", got)
	}

	cmd = cronConvert()
	cmd.SetArgs([]string{"--to", "yaml", "0 9 * * 1-5"})
	if err := cmd.Execute(); !errors.Is(err, errUsage) {
		t.Fatalf("invalid destination error = %v, want usage", err)
	}
}

func TestCronConvert_JSONUsesStableSuccessAndRefusalStreams(t *testing.T) {
	jsonOut = true
	t.Cleanup(func() { jsonOut = false })

	t.Run("success on stdout", func(t *testing.T) {
		cmd := cronConvert()
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		cmd.SetArgs([]string{"3rd friday monthly at 09:00"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
		if stderr.Len() != 0 {
			t.Fatalf("stderr = %q, want empty", stderr.String())
		}
		var got cron.Conversion
		if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
			t.Fatalf("stdout is not conversion JSON: %v\n%s", err, stdout.String())
		}
		if got.InputSyntax != cron.SyntaxHuman || got.OutputSyntax != cron.SyntaxCron ||
			got.Input != "3rd friday monthly at 09:00" || got.Output != "0 9 * * 5#3" || got.RefusalReason != "" {
			t.Errorf("success object = %+v", got)
		}
	})

	t.Run("refusal on stderr", func(t *testing.T) {
		cmd := cronConvert()
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		cmd.SetArgs([]string{"every 7 minutes"})
		err := cmd.Execute()
		if !errors.Is(err, errUsage) || !isReported(err) {
			t.Fatalf("error = %v, want reported usage", err)
		}
		if stdout.Len() != 0 {
			t.Fatalf("stdout = %q, want empty", stdout.String())
		}
		var got cron.Conversion
		if err := json.Unmarshal(stderr.Bytes(), &got); err != nil {
			t.Fatalf("stderr is not conversion JSON: %v\n%s", err, stderr.String())
		}
		if got.InputSyntax != cron.SyntaxHuman || got.OutputSyntax != cron.SyntaxCron ||
			got.Input != "every 7 minutes" || got.Output != "" || got.RefusalReason == "" {
			t.Errorf("refusal object = %+v", got)
		}
	})
}

// fakeCreator records what the import asked for, standing in for the daemon so
// the reporting behavior can be exercised without one.
type fakeCreator struct {
	reqs []server.TaskCreateRequest
	err  error
}

type scriptedCreator struct {
	calls int
	reqs  []server.TaskCreateRequest
}

func (f *scriptedCreator) CreateTask(_ context.Context, req server.TaskCreateRequest) (server.TaskResponse, error) {
	f.calls++
	if f.calls == 1 {
		return server.TaskResponse{}, context.DeadlineExceeded
	}
	f.reqs = append(f.reqs, req)
	return server.TaskResponse{Task: domain.Task{ID: "created", Name: req.Name}}, nil
}

func (f *fakeCreator) CreateTask(_ context.Context, req server.TaskCreateRequest) (server.TaskResponse, error) {
	if f.err != nil {
		return server.TaskResponse{}, f.err
	}
	f.reqs = append(f.reqs, req)
	return server.TaskResponse{Task: domain.Task{ID: "id-" + req.Name, Name: req.Name}}, nil
}

const sampleCrontab = `MAILTO=ops@example.com
# nightly backup
0 2 * * * /usr/local/bin/backup --full

*/15 * * * * /usr/local/bin/probe
@reboot /usr/local/bin/warm
*/7 * * * * /usr/local/bin/odd
0 9 * * *
`

func scan(t *testing.T, text string) *cron.Report {
	t.Helper()
	rep, err := cron.ScanCrontab(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	return &rep
}

// TestImport_DryRunCreatesNothing is FR-005: the preview produces the whole
// report and creates nothing.
func TestImport_DryRunCreatesNothing(t *testing.T) {
	rep := scan(t, sampleCrontab)
	var buf bytes.Buffer
	if err := runImport(&buf, rep, importOptions{dryRun: true, timezone: "UTC"}, nil); err != nil {
		t.Fatalf("dry run returned an error: %v", err)
	}
	out := buf.String()

	if rep.Created != 0 {
		t.Errorf("dry run created %d task(s), want 0", rep.Created)
	}
	for _, want := range []string{
		"MAILTO",                     // the assignment is warned about, not dropped
		"every 15 minutes",           // a translated line
		"every day at 02:00",         // another
		"boot",                       // @reboot declined by name
		"at minutes 0, 7, 14",        // uneven field step is described as a field set
		"no command follows",         // the schedule-only line is an error
		"This was a preview",         // the run says it changed nothing
		"Cron carries no timezone",   // fidelity statement
		"no catch-up, overlap, or r", // fidelity statement
	} {
		if !strings.Contains(out, want) {
			t.Errorf("dry-run output missing %q\n---\n%s", want, out)
		}
	}
}

// TestImport_CountsEveryLine is FR-010 and SC-002: every line is accounted for,
// so a silently dropped line would show up as a counting error.
func TestImport_CountsEveryLine(t *testing.T) {
	rep := scan(t, sampleCrontab)
	if got, want := rep.Read, 8; got != want {
		t.Fatalf("read %d lines, want %d", got, want)
	}
	if got, want := rep.Jobs, 3; got != want {
		t.Errorf("jobs = %d, want %d", got, want)
	}
	if got, want := rep.Declined, 1; got != want {
		t.Errorf("declined = %d, want %d", got, want)
	}
	if got, want := rep.Errors, 1; got != want {
		t.Errorf("errors = %d, want %d", got, want)
	}
	if got, want := rep.Skipped, 3; got != want { // MAILTO, comment, blank
		t.Errorf("skipped = %d, want %d", got, want)
	}
	if total := rep.Jobs + rep.Declined + rep.Errors + rep.Skipped; total != rep.Read {
		t.Errorf("line accounting does not balance: %d classified vs %d read", total, rep.Read)
	}
}

// TestImport_CreatesSupportedLines covers FR-005a: declined lines do not stop
// the supported ones from being created, and the payload comes across.
func TestImport_CreatesSupportedLines(t *testing.T) {
	rep := scan(t, sampleCrontab)
	fc := &fakeCreator{}
	var buf bytes.Buffer
	if err := runImport(&buf, rep, importOptions{timezone: "UTC", group: "g1"}, fc); err != nil {
		t.Fatalf("import returned an error: %v", err)
	}
	if len(fc.reqs) != 3 {
		t.Fatalf("created %d task(s), want 3", len(fc.reqs))
	}
	first := fc.reqs[0]
	if first.Command != "/usr/local/bin/backup" {
		t.Errorf("command = %q, want the crontab's program", first.Command)
	}
	if len(first.Args) != 1 || first.Args[0] != "--full" {
		t.Errorf("args = %v, want [--full]", first.Args)
	}
	if first.Schedule != "0 2 * * *" || first.ScheduleSyntax != "cron" {
		t.Errorf("schedule input = %q (%q), want retained cron source", first.Schedule, first.ScheduleSyntax)
	}
	if first.Timezone != "UTC" || first.GroupID != "g1" {
		t.Errorf("timezone/group not applied: %q / %q", first.Timezone, first.GroupID)
	}
	if rep.Created != 3 {
		t.Errorf("report Created = %d, want 3", rep.Created)
	}
}

// TestImport_PartialFailureKeepsWhatWasCreated is the rest of FR-005a: a
// creation failure is reported alongside the created count rather than rolled
// back or hidden.
func TestImport_PartialFailureKeepsWhatWasCreated(t *testing.T) {
	rep := scan(t, "0 2 * * * /usr/local/bin/backup\n*/15 * * * * /usr/local/bin/probe\n")
	fc := &fakeCreator{err: context.DeadlineExceeded}
	var buf bytes.Buffer
	err := runImport(&buf, rep, importOptions{timezone: "UTC"}, fc)
	if err == nil {
		t.Fatal("a run where every creation failed should report a failure")
	}
	if rep.Failed != 2 {
		t.Errorf("failed = %d, want 2", rep.Failed)
	}
	if !strings.Contains(buf.String(), "not created") {
		t.Errorf("the failure was not reported per line:\n%s", buf.String())
	}
}

func TestImport_MixedPartialSuccessContinuesWithCronSource(t *testing.T) {
	rep := scan(t, "0 2 * * * /usr/local/bin/backup\n*/15 * * * * /usr/local/bin/probe\n")
	creator := &scriptedCreator{}
	var buf bytes.Buffer
	err := runImport(&buf, rep, importOptions{timezone: "UTC"}, creator)
	if err == nil {
		t.Fatal("partial failure should be reported after all lines are attempted")
	}
	if creator.calls != 2 || rep.Failed != 1 || rep.Created != 1 {
		t.Fatalf("calls=%d failed=%d created=%d", creator.calls, rep.Failed, rep.Created)
	}
	if len(creator.reqs) != 1 || creator.reqs[0].Schedule != "*/15 * * * *" || creator.reqs[0].ScheduleSyntax != "cron" {
		t.Fatalf("successful request = %+v, want retained second cron line", creator.reqs)
	}
	if out := buf.String(); !strings.Contains(out, "not created") || !strings.Contains(out, "created:") {
		t.Fatalf("mixed result not fully reported:\n%s", out)
	}
}

// TestImport_PreviewInputMatchesCreatedTask keeps the daemon preview and create
// request on the same retained cron source while the report remains readable.
func TestImport_PreviewInputMatchesCreatedTask(t *testing.T) {
	const text = "0 2 * * * /usr/local/bin/backup\n0 9 * * 1-5 /usr/local/bin/report\n0 9 1 * * /usr/local/bin/invoice\n"

	preview := scan(t, text)
	type previewCall struct{ input, syntax, timezone string }
	var previews []previewCall
	var previewOut bytes.Buffer
	if err := runImport(&previewOut, preview, importOptions{
		dryRun: true, timezone: "UTC", count: 1,
		runs: func(input, syntax, timezone string, _ int) ([]time.Time, error) {
			previews = append(previews, previewCall{input, syntax, timezone})
			return nil, nil
		},
	}, nil); err != nil {
		t.Fatal(err)
	}

	actual := scan(t, text)
	fc := &fakeCreator{}
	var actualOut bytes.Buffer
	if err := runImport(&actualOut, actual, importOptions{timezone: "UTC"}, fc); err != nil {
		t.Fatal(err)
	}

	if len(fc.reqs) != 3 {
		t.Fatalf("created %d task(s), want 3", len(fc.reqs))
	}
	if len(previews) != 3 {
		t.Fatalf("previewed %d task(s), want 3", len(previews))
	}
	for i, line := range jobLines(preview) {
		if previews[i].input != line.Expr || previews[i].syntax != "cron" || previews[i].timezone != "UTC" {
			t.Errorf("line %d preview input = %+v, want cron %q in UTC", line.Number, previews[i], line.Expr)
		}
		if fc.reqs[i].Schedule != line.Expr || fc.reqs[i].ScheduleSyntax != "cron" {
			t.Errorf("line %d: task input = %q (%q), want cron source %q",
				line.Number, fc.reqs[i].Schedule, fc.reqs[i].ScheduleSyntax, line.Expr)
		}
		if !strings.Contains(previewOut.String(), line.Phrase) {
			t.Errorf("line %d: phrase %q was not shown in the preview", line.Number, line.Phrase)
		}
	}
}

func jobLines(rep *cron.Report) []cron.Line {
	var out []cron.Line
	for _, l := range rep.Lines {
		if l.Kind == cron.LineJob {
			out = append(out, l)
		}
	}
	return out
}

// TestCronImportUsesTheTaskInputBoundary pins the architecture: raw cron still
// cannot enter the human parser, but import retains it through the API's typed
// authoring boundary rather than replacing it with generated prose.
func TestCronImportUsesTheTaskInputBoundary(t *testing.T) {
	anchor := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	for _, expr := range []string{
		"0 9 * * 1-5", "*/15 * * * *", "@daily", "0 0 1 * *",
	} {
		if _, err := schedule.Parse(expr, "UTC", anchor); err == nil {
			t.Errorf("schedule.Parse accepted cron %q; cron must enter through the shared schedule-input boundary", expr)
		}
	}

	rep := scan(t, "0 9 * * 1-5 /usr/local/bin/report\n")
	fc := &fakeCreator{}
	if err := runImport(&bytes.Buffer{}, rep, importOptions{timezone: "UTC"}, fc); err != nil {
		t.Fatal(err)
	}
	if len(fc.reqs) != 1 {
		t.Fatalf("created %d task(s), want 1", len(fc.reqs))
	}
	if fc.reqs[0].Schedule != "0 9 * * 1-5" || fc.reqs[0].ScheduleSyntax != "cron" {
		t.Errorf("import request = %+v, want retained cron with explicit hint", fc.reqs[0])
	}
}

// TestExportLines_EveryTaskAppearsOnce is FR-012: a task is either a line or a
// named refusal, never absent.
func TestExportLines_EveryTaskAppearsOnce(t *testing.T) {
	anchor := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	recurring, err := schedule.Parse("every day at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	details := []server.TaskResponse{
		{Task: domain.Task{ID: "1", Name: "report", Command: "/bin/report", Enabled: true, State: domain.TaskActive}, Schedule: recurring},
		{Task: domain.Task{ID: "2", Name: "once", Enabled: true, State: domain.TaskActive}, Schedule: schedule.NewOneOff(anchor)},
		{Task: domain.Task{ID: "3", Name: "off", Command: "/bin/x", Enabled: false, State: domain.TaskDisabled}, Schedule: recurring},
	}

	lines := exportLines(details)
	if len(lines) != len(details) {
		t.Fatalf("exported %d entries for %d tasks", len(lines), len(details))
	}
	if lines[0].Line != "0 9 * * * /bin/report" {
		t.Errorf("line = %q, want the crontab form", lines[0].Line)
	}
	for _, i := range []int{1, 2} {
		if lines[i].Declined == "" {
			t.Errorf("task %q was exported as a live line, want a refusal", lines[i].Name)
		}
	}

	var buf bytes.Buffer
	printExport(&buf, details)
	out := buf.String()
	if strings.Count(out, "# declined:") != 2 {
		t.Errorf("expected two commented refusals:\n%s", out)
	}
	if !strings.Contains(out, "3 task(s)") {
		t.Errorf("header did not account for every task:\n%s", out)
	}
}

// TestExport_EmptySetSucceeds is FR-011a: an empty export is recognizably an
// empty export rather than no output at all.
func TestExport_EmptySetSucceeds(t *testing.T) {
	var buf bytes.Buffer
	printExport(&buf, nil)
	if !strings.Contains(buf.String(), "0 task(s)") {
		t.Errorf("empty export produced %q", buf.String())
	}
}

// TestPrintExplain_RefusalIsAnAnswer pins the rendering of a refusal: it is
// reported, not raised.
func TestPrintExplain_RefusalIsAnAnswer(t *testing.T) {
	var buf bytes.Buffer
	printExplain(&buf, explainResult{Expression: "@reboot", Unsupported: "fires at boot", Timezone: "UTC"})
	out := buf.String()
	if !strings.Contains(out, "unsupported: fires at boot") {
		t.Errorf("refusal not rendered:\n%s", out)
	}
	if strings.Contains(out, "phrase:") {
		t.Errorf("a refusal must not print a phrase:\n%s", out)
	}
}

// TestImport_ShowsNextRunTimes is FR-004: the report shows when each line would
// actually fire. A phrase can read correctly and still mean something else, so
// the run times are the half of the preview that catches a misreading.
func TestImport_ShowsNextRunTimes(t *testing.T) {
	rep := scan(t, "0 2 * * * /usr/local/bin/backup\n")
	fixed := time.Date(2026, 7, 24, 2, 0, 0, 0, time.UTC)
	opts := importOptions{
		dryRun: true, timezone: "UTC", count: 2,
		runs: func(_, _, _ string, count int) ([]time.Time, error) {
			out := make([]time.Time, 0, count)
			for i := 0; i < count; i++ {
				out = append(out, fixed.AddDate(0, 0, i))
			}
			return out, nil
		},
	}
	var buf bytes.Buffer
	if err := runImport(&buf, rep, opts, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"2026-07-24T02:00:00Z", "2026-07-25T02:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing run time %q\n---\n%s", want, out)
		}
	}
}

// TestImport_SurvivesAnUnreachableDaemon: a preview is most wanted exactly when
// the daemon is not running, so losing the run times must not lose the report.
func TestImport_SurvivesAnUnreachableDaemon(t *testing.T) {
	rep := scan(t, "0 2 * * * /usr/local/bin/backup\n")
	opts := importOptions{
		dryRun: true, timezone: "UTC", count: 3,
		runs: func(string, string, string, int) ([]time.Time, error) { return nil, context.DeadlineExceeded },
	}
	var buf bytes.Buffer
	if err := runImport(&buf, rep, opts, nil); err != nil {
		t.Fatalf("an unreachable daemon must not fail the preview: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "every day at 02:00") {
		t.Errorf("the conversion itself is local and must still be reported:\n%s", out)
	}
	if !strings.Contains(out, "unavailable") {
		t.Errorf("the missing run times should be stated, not silently absent:\n%s", out)
	}
}
