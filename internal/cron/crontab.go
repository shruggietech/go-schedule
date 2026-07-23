package cron

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// LineKind classifies what a crontab line turned out to be.
type LineKind int

const (
	LineSkipped  LineKind = iota // a comment or a blank line
	LineJob                      // a schedule and a command
	LineDeclined                 // well-formed, but not representable
	LineError                    // malformed
)

// Line is one crontab line's conversion result. Preview and a real import
// produce identical slices of these; whether tasks were created is the only
// difference between the two runs.
type Line struct {
	Number  int
	Raw     string
	Kind    LineKind
	Expr    string // the timing portion, for a job or a declined line
	Phrase  string // the human phrase, for a job
	Command string
	Args    []string
	Reason  string // why it was declined, or what went wrong
	// Warnings record what was read but not carried across — a MAILTO, a shell
	// variable assignment. They are attached to the line that produced them
	// rather than dropped, because a dropped MAILTO silently changes where a
	// job's output goes.
	Warnings []string
}

// Report is the account of one conversion run.
type Report struct {
	Lines    []Line
	Read     int
	Jobs     int
	Skipped  int
	Declined int
	Errors   int
	Created  int
	// Failed counts jobs that converted but could not be created. Non-zero means
	// a partial import: the tasks that were created remain.
	Failed int
	// Warnings are file-level notes (variable assignments seen before any job).
	Warnings []string
}

// reAssignment matches a crontab environment assignment (NAME=value), which
// crontab applies to every subsequent line. We report these rather than apply
// them: applying one would change the meaning of every line after it in a way
// the preview could not show.
var reAssignment = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*=(.*)$`)

// ScanCrontab reads a crontab and converts every line. It never returns a
// partial result: a line that cannot be converted becomes a declined or errored
// Line, so the caller can account for every line of the input.
func ScanCrontab(r io.Reader) (Report, error) {
	var rep Report
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	num := 0
	for sc.Scan() {
		num++
		raw := sc.Text()
		rep.Read++
		line := convertLine(num, raw, &rep)
		rep.Lines = append(rep.Lines, line)
		switch line.Kind {
		case LineSkipped:
			rep.Skipped++
		case LineJob:
			rep.Jobs++
		case LineDeclined:
			rep.Declined++
		case LineError:
			rep.Errors++
		}
	}
	if err := sc.Err(); err != nil {
		return rep, fmt.Errorf("cron: read crontab: %w", err)
	}
	return rep, nil
}

func convertLine(num int, raw string, rep *Report) Line {
	line := Line{Number: num, Raw: raw}
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		line.Kind = LineSkipped
		return line
	}

	if m := reAssignment.FindStringSubmatch(trimmed); m != nil {
		line.Kind = LineSkipped
		note := fmt.Sprintf("line %d: %s is not carried across — crontab variables have no equivalent here", num, m[1])
		if strings.EqualFold(m[1], "MAILTO") {
			note = fmt.Sprintf("line %d: MAILTO=%s is not carried across — run output is recorded in the run history instead of mailed",
				num, strings.TrimSpace(m[2]))
		}
		line.Warnings = append(line.Warnings, note)
		rep.Warnings = append(rep.Warnings, note)
		return line
	}

	expr, command, ok := splitTiming(trimmed)
	if !ok {
		line.Kind = LineError
		line.Reason = "no command follows the schedule"
		return line
	}
	line.Expr, line.Command = expr, command

	phrase, bad, err := Explain(expr)
	switch {
	case err != nil:
		line.Kind = LineError
		line.Reason = err.Error()
	case bad.Reason != "":
		line.Kind = LineDeclined
		line.Reason = bad.Reason
	default:
		line.Kind = LineJob
		line.Phrase = phrase
		line.Command, line.Args = splitCommand(command)
	}
	return line
}

// splitTiming separates a line's timing fields from its command. A shorthand
// takes one field; a standard expression takes five.
func splitTiming(s string) (expr, command string, ok bool) {
	fields := strings.Fields(s)
	n := 5
	if strings.HasPrefix(fields[0], "@") {
		n = 1
	}
	if len(fields) <= n {
		return strings.Join(fields, " "), "", false
	}
	expr = strings.Join(fields[:n], " ")
	// Recover the command with its original internal spacing rather than the
	// tokenized form, so a quoted argument survives.
	idx := 0
	for i := 0; i < n; i++ {
		idx = strings.Index(s[idx:], fields[i]) + idx + len(fields[i])
	}
	return expr, strings.TrimSpace(s[idx:]), true
}

// splitCommand splits a command line into its program and arguments. It is
// deliberately whitespace-only: crontab hands the whole string to a shell, and
// re-implementing shell quoting here would introduce a second, subtly different
// parser. A command needing shell semantics keeps them by being imported as a
// shell invocation, which the operator can see in the preview.
func splitCommand(s string) (string, []string) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}
