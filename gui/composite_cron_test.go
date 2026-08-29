package gui

import "testing"

func TestEditor_CompositeCronPreviewCreateAndPrefill(t *testing.T) {
	const expression = "*/10 9-17 * * MON,WED,FRI"
	e, fb := newTestEditor(t, nil)
	e.name.SetText("business-hours")
	e.command.SetText("cmd")
	e.schedule.SetText(expression)

	if e.save.Disabled() {
		t.Fatal("Save should be enabled for composite cron")
	}
	if _, req := fb.lastPreviewCall(); req.Schedule != expression || req.ScheduleSyntax != "cron" {
		t.Fatalf("preview request=%+v, want exact composite cron source", req)
	}
	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	if _, req := fb.lastCreateCall(); req.Schedule != expression || req.ScheduleSyntax != "cron" {
		t.Fatalf("create request=%+v, want exact composite cron source", req)
	}

	detail := recurringDetail(expression)
	detail.Schedule.SourceSyntax = "human" // Stored expression is authoritative.
	edit, editBackend := newTestEditorDetail(t, detail)
	if edit.schedule.Text != expression {
		t.Fatalf("prefill=%q, want %q", edit.schedule.Text, expression)
	}
	if _, req := editBackend.lastPreviewCall(); req.Schedule != expression || req.ScheduleSyntax != "cron" {
		t.Fatalf("prefill preview=%+v, want cron classified from current text", req)
	}
}
