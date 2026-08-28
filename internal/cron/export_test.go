package cron

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

var exportAnchor = time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

func taskFor(phrase string) (domain.Task, domain.Schedule, error) {
	sch, err := schedule.Parse(phrase, "UTC", exportAnchor)
	if err != nil {
		return domain.Task{}, domain.Schedule{}, err
	}
	return domain.Task{Name: "t", Command: "/bin/true", Enabled: true, State: domain.TaskActive}, sch, nil
}

// TestExport_Expressible covers FR-012: a schedule cron can carry comes out as a
// crontab line.
func TestExport_Expressible(t *testing.T) {
	for _, c := range []struct{ phrase, want string }{
		{"every 15 minutes", "*/15 * * * *"},
		{"every hour", "0 * * * *"},
		{"every hour starting at 00:30", "30 * * * *"},
		{"every 2 hours starting at 08:30", "30 */2 * * *"},
		{"every day at 09:00", "0 9 * * *"},
		{"weekdays at 09:00", "0 9 * * 1-5"},
		{"every wednesday at 14:00", "0 14 * * 3"},
		{"on the 1st of every month at 09:00", "0 9 1 * *"},
		{"every year on february 29 at 00:00", "0 0 29 2 *"},
	} {
		t.Run(c.phrase, func(t *testing.T) {
			task, sch, err := taskFor(c.phrase)
			if err != nil {
				t.Fatal(err)
			}
			got, bad, ok := Export(task, sch)
			if !ok {
				t.Fatalf("Export refused %q: %s", c.phrase, bad.Reason)
			}
			if got != c.want {
				t.Fatalf("Export(%q) = %q, want %q", c.phrase, got, c.want)
			}
		})
	}
}

func TestExportSchedule_PreservesTaskExportMapping(t *testing.T) {
	task, sch, err := taskFor("weekdays at 09:00")
	if err != nil {
		t.Fatal(err)
	}
	want, wantBad, wantOK := Export(task, sch)
	got, gotBad, gotOK := ExportSchedule(sch, task.MissingDatePolicy)
	if got != want || gotBad != wantBad || gotOK != wantOK {
		t.Fatalf("schedule export = (%q, %+v, %v), task export = (%q, %+v, %v)",
			got, gotBad, gotOK, want, wantBad, wantOK)
	}

	task.Enabled = false
	if _, bad, ok := Export(task, sch); ok || !strings.Contains(bad.Reason, "disabled") {
		t.Fatalf("task wrapper lost disabled-state refusal: ok=%v bad=%+v", ok, bad)
	}
}

func TestExport_UsesImplicitDailyAnchorWithoutInventingPrecision(t *testing.T) {
	t.Run("minute-aligned", func(t *testing.T) {
		anchor := time.Date(2026, 6, 1, 12, 34, 0, 0, time.UTC)
		sch, err := schedule.Parse("every day", "UTC", anchor)
		if err != nil {
			t.Fatal(err)
		}
		task := domain.Task{Enabled: true, State: domain.TaskActive}
		got, bad, ok := Export(task, sch)
		if !ok {
			t.Fatalf("Export refused minute-aligned daily schedule: %s", bad.Reason)
		}
		if got != "34 12 * * *" {
			t.Fatalf("Export = %q, want anchor-preserving %q", got, "34 12 * * *")
		}
	})

	t.Run("sub-minute", func(t *testing.T) {
		anchor := time.Date(2026, 6, 1, 12, 34, 56, 0, time.UTC)
		sch, err := schedule.Parse("every day", "UTC", anchor)
		if err != nil {
			t.Fatal(err)
		}
		task := domain.Task{Enabled: true, State: domain.TaskActive}
		if _, bad, ok := Export(task, sch); ok || !strings.Contains(bad.Reason, "one-minute resolution") {
			t.Fatalf("sub-minute daily export = ok:%v bad:%q", ok, bad.Reason)
		}
	})
}

// TestExport_Declines covers FR-012 and FR-012a: what cron cannot carry is
// refused by name, never approximated, and never silently omitted.
func TestExport_Declines(t *testing.T) {
	t.Run("one-off", func(t *testing.T) {
		task := domain.Task{Name: "once", Enabled: true, State: domain.TaskActive}
		sch := schedule.NewOneOff(exportAnchor.Add(24 * time.Hour))
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export produced a line for a one-off schedule")
		}
		if !strings.Contains(bad.Reason, "once") {
			t.Errorf("reason = %q, want it to mention firing once", bad.Reason)
		}
	})

	t.Run("disabled task", func(t *testing.T) {
		task, sch, err := taskFor("every day at 09:00")
		if err != nil {
			t.Fatal(err)
		}
		task.Enabled = false
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export produced a live line for a disabled task")
		}
		if !strings.Contains(bad.Reason, "disabled") {
			t.Errorf("reason = %q, want it to mention the disabled state", bad.Reason)
		}
	})

	t.Run("sub-minute", func(t *testing.T) {
		task, sch, err := taskFor("every 30 seconds")
		if err != nil {
			t.Fatal(err)
		}
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export produced a line for a sub-minute schedule")
		}
		if !strings.Contains(bad.Reason, "sub-minute") {
			t.Errorf("reason = %q, want it to mention sub-minute resolution", bad.Reason)
		}
	})

	t.Run("ordinal weekday", func(t *testing.T) {
		task, sch, err := taskFor("3rd wednesday monthly at 14:00")
		if err != nil {
			t.Fatal(err)
		}
		if _, _, ok := Export(task, sch); ok {
			t.Fatal("Export produced a line for an ordinal-weekday rule")
		}
	})

	t.Run("non-default missing-date policy", func(t *testing.T) {
		task, sch, err := taskFor("on the 31st of every month at 09:00")
		if err != nil {
			t.Fatal(err)
		}
		task.MissingDatePolicy = domain.MissingDateLastValid
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export silently dropped a missing-date policy cron cannot carry")
		}
		if !strings.Contains(bad.Reason, "missing-date policy") {
			t.Errorf("reason = %q, want it to name the policy", bad.Reason)
		}
	})

	t.Run("skip policy on the same rule is expressible", func(t *testing.T) {
		// The default policy *is* cron's behavior, so this must still export.
		task, sch, err := taskFor("on the 31st of every month at 09:00")
		if err != nil {
			t.Fatal(err)
		}
		task.MissingDatePolicy = domain.MissingDateSkip
		got, bad, ok := Export(task, sch)
		if !ok {
			t.Fatalf("Export refused a skip-policy rule: %s", bad.Reason)
		}
		if got != "0 9 31 * *" {
			t.Fatalf("got %q, want %q", got, "0 9 31 * *")
		}
	})

	t.Run("multi-day interval", func(t *testing.T) {
		task, sch, err := taskFor("every 3 days")
		if err != nil {
			t.Fatal(err)
		}
		if _, _, ok := Export(task, sch); ok {
			t.Fatal("Export produced a line for an every-3-days rule, which cron cannot express")
		}
	})

	t.Run("misaligned sub-daily phase", func(t *testing.T) {
		task, sch, err := taskFor("every 15 minutes starting at 00:05")
		if err != nil {
			t.Fatal(err)
		}
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export approximated a :05/:20/:35/:50 interval as */15")
		}
		if !strings.Contains(bad.Reason, "phase") {
			t.Errorf("reason = %q, want it to name the phase", bad.Reason)
		}
	})

	t.Run("misaligned hourly phase", func(t *testing.T) {
		task, sch, err := taskFor("every 2 hours starting at 09:30")
		if err != nil {
			t.Fatal(err)
		}
		_, bad, ok := Export(task, sch)
		if ok {
			t.Fatal("Export approximated an odd-hour interval as */2")
		}
		if !strings.Contains(bad.Reason, "phase") {
			t.Errorf("reason = %q, want it to name the phase", bad.Reason)
		}
	})
}

// TestRoundTrip_CrossesDSTAndMonthBoundary is FR-013 / SC-003, the strengthening
// the issue asked for: cron → phrase → schedule → cron, with the resulting
// expression producing the same run times as the original over a window spanning
// the 2026-03-08 US spring-forward and a month boundary.
func TestRoundTrip_CrossesDSTAndMonthBoundary(t *testing.T) {
	const tz = "America/New_York"
	start := time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)

	for _, expr := range []string{
		"*/15 * * * *",
		"0 9 * * *",
		"0 9 * * 1-5",
		"30 2 * * *", // 02:30 — inside the hour that does not exist on 2026-03-08
		"0 9 1 * *",  // crosses the month boundary
		"0 0 4 3 *",
	} {
		t.Run(expr, func(t *testing.T) {
			phrase, bad, err := Explain(expr)
			if err != nil || bad.Reason != "" {
				t.Fatalf("Explain(%q): err=%v refusal=%q", expr, err, bad.Reason)
			}
			sch, err := schedule.Parse(phrase, tz, start)
			if err != nil {
				t.Fatalf("Parse(%q): %v", phrase, err)
			}
			task := domain.Task{Enabled: true, State: domain.TaskActive, Timezone: tz}

			back, bad, ok := Export(task, sch)
			if !ok {
				t.Fatalf("Export refused a schedule that came from cron: %s", bad.Reason)
			}

			// Re-import the exported expression and compare every run time in
			// the window. Identical phrases would make this circular, so the
			// comparison is on run times.
			phrase2, bad, err := Explain(back)
			if err != nil || bad.Reason != "" {
				t.Fatalf("Explain(%q) after export: err=%v refusal=%q", back, err, bad.Reason)
			}
			sch2, err := schedule.Parse(phrase2, tz, start)
			if err != nil {
				t.Fatalf("Parse(%q): %v", phrase2, err)
			}

			a := runsBetween(t, sch, tz, start, end)
			b := runsBetween(t, sch2, tz, start, end)
			if len(a) == 0 {
				t.Fatalf("the original produced no runs in the window; the test would prove nothing")
			}
			if len(a) != len(b) {
				t.Fatalf("round trip changed the run count: %d -> %d (%q -> %q)", len(a), len(b), expr, back)
			}
			for i := range a {
				if !a[i].Equal(b[i]) {
					t.Fatalf("round trip moved run %d: %v -> %v (%q -> %q)", i, a[i], b[i], expr, back)
				}
			}
		})
	}
}
