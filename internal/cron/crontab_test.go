package cron

import (
	"reflect"
	"strings"
	"testing"
)

func scanFixture(t *testing.T, text string, opts ScanOptions) Report {
	t.Helper()
	rep, err := ScanCrontabWithOptions(strings.NewReader(text), opts)
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

func jobs(rep Report) []Line {
	var out []Line
	for _, line := range rep.Lines {
		if line.Kind == LineJob {
			out = append(out, line)
		}
	}
	return out
}

func TestScanCrontab_ContextSnapshotsAndShell(t *testing.T) {
	rep := scanFixture(t, `CRON_TZ=America/New_York
PATH = /usr/local/bin:/usr/bin
TZ=UTC
SHELL="/bin/bash"
0 9 * * * printf '\%s\n' "$TZ" && date
CRON_TZ=Europe/London
PATH=/bin
0 17 * * * echo done > /tmp/done
`, ScanOptions{})
	got := jobs(rep)
	if len(got) != 2 {
		t.Fatalf("jobs = %d, report=%+v", len(got), rep)
	}
	if got[0].Timezone != "America/New_York" || got[1].Timezone != "Europe/London" {
		t.Fatalf("timezones = %q, %q", got[0].Timezone, got[1].Timezone)
	}
	if got[0].Command != "/bin/bash" || !reflect.DeepEqual(got[0].Args, []string{"-c", `printf '%s\n' "$TZ" && date`}) {
		t.Fatalf("first command = %q %#v", got[0].Command, got[0].Args)
	}
	if got[0].Env["PATH"] != "/usr/local/bin:/usr/bin" || got[0].Env["TZ"] != "UTC" || got[0].Env["SHELL"] != "/bin/bash" {
		t.Fatalf("first env = %#v", got[0].Env)
	}
	if got[1].Env["PATH"] != "/bin" || got[0].Env["PATH"] != "/usr/local/bin:/usr/bin" {
		t.Fatalf("environment snapshots alias: first=%#v second=%#v", got[0].Env, got[1].Env)
	}
}

func TestScanCrontab_DefaultShellIsExecutionEnvironment(t *testing.T) {
	rep := scanFixture(t, "0 9 * * * printf '%s' \"$SHELL\"\n", ScanOptions{})
	got := jobs(rep)
	if len(got) != 1 || got[0].Command != "/bin/sh" || got[0].Env["SHELL"] != "/bin/sh" {
		t.Fatalf("default shell context = %#v", got)
	}
}

func TestScanCrontab_AssignmentQuotesAndMailWarnings(t *testing.T) {
	rep := scanFixture(t, "EMPTY=\nPAD='  kept  '\nLOGNAME=wrong\nMAILTO=ops@example.com\nMAILFROM=cron@example.com\n0 0 * * * echo ok\n", ScanOptions{})
	got := jobs(rep)
	if len(got) != 1 || got[0].Env["EMPTY"] != "" || got[0].Env["PAD"] != "  kept  " {
		t.Fatalf("job environment = %#v", got)
	}
	if _, ok := got[0].Env["MAILTO"]; ok {
		t.Fatalf("MAILTO was represented as delivery: %#v", got[0].Env)
	}
	if _, ok := got[0].Env["LOGNAME"]; ok {
		t.Fatalf("LOGNAME override was applied: %#v", got[0].Env)
	}
	if len(rep.Warnings) != 3 {
		t.Fatalf("warnings = %#v", rep.Warnings)
	}
}

func TestScanCrontab_InvalidAssignmentAndDialectAreVisible(t *testing.T) {
	rep := scanFixture(t, "PATH='unterminated\n0 0 * * * echo ok\n", ScanOptions{})
	if rep.Errors != 1 || rep.Jobs != 1 || !strings.Contains(rep.Lines[0].Reason, "unmatched quote") {
		t.Fatalf("report = %+v", rep)
	}
	if _, err := ScanCrontabWithOptions(strings.NewReader(""), ScanOptions{Dialect: "guess"}); err == nil {
		t.Fatal("invalid dialect was accepted")
	}
}

func TestScanCrontab_PercentSemantics(t *testing.T) {
	rep := scanFixture(t, "0 0 * * * printf 'quoted%still-special'\\%literal%line1%line2\\%tail\n0 1 * * * echo empty%\n", ScanOptions{})
	got := jobs(rep)
	if len(got) != 2 {
		t.Fatalf("jobs = %d", len(got))
	}
	if got[0].CommandText != "printf 'quoted" || got[0].Stdin != "still-special'%literal\nline1\nline2%tail" || !got[0].HasStdin {
		t.Fatalf("first percent result = command %q stdin %q present=%v", got[0].CommandText, got[0].Stdin, got[0].HasStdin)
	}
	if got[1].CommandText != "echo empty" || got[1].Stdin != "" || !got[1].HasStdin {
		t.Fatalf("empty stdin separator lost: %+v", got[1])
	}
}

func TestScanCrontab_ExplicitSystemLayout(t *testing.T) {
	text := "0 2 * * * root /usr/local/bin/backup --full\n"
	user := jobs(scanFixture(t, text, ScanOptions{}))[0]
	if user.RunAs != "" || user.CommandText != "root /usr/local/bin/backup --full" {
		t.Fatalf("user layout = %+v", user)
	}
	system := jobs(scanFixture(t, text, ScanOptions{System: true}))[0]
	if system.RunAs != "root" || system.CommandText != "/usr/local/bin/backup --full" {
		t.Fatalf("system layout = %+v", system)
	}
}

func TestScanCrontab_ExplicitQuartzLayout(t *testing.T) {
	text := "0 0 12 ? * MON echo noon\n"
	unix := scanFixture(t, text, ScanOptions{})
	if unix.Jobs != 0 {
		t.Fatalf("unix layout guessed Quartz: %+v", unix)
	}
	quartz := scanFixture(t, text, ScanOptions{Dialect: DialectQuartz})
	got := jobs(quartz)
	if len(got) != 1 || got[0].Expr != "0 0 12 ? * MON" || got[0].CommandText != "echo noon" {
		t.Fatalf("quartz layout = %+v", quartz)
	}
}
