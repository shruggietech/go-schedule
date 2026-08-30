package cron

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Dialect identifies the timing-field layout used by a crontab file.
type Dialect string

const (
	// DialectUnix consumes five timing fields in conventional crontab order.
	DialectUnix Dialect = "unix"
	// DialectQuartz consumes six timing fields beginning with seconds.
	DialectQuartz Dialect = "quartz"
)

// ScanOptions selects file layout that cannot be inferred safely from tokens.
type ScanOptions struct {
	Dialect Dialect
	System  bool
}

// LineKind classifies one scanned crontab line.
type LineKind int

const (
	LineSkipped LineKind = iota
	LineJob
	LineDeclined
	LineError
)

// Line contains one crontab line's conversion outcome and effective context.
type Line struct {
	Number      int
	Raw         string
	Kind        LineKind
	Expr        string
	Phrase      string
	Command     string
	Args        []string
	CommandText string
	Stdin       string
	HasStdin    bool
	Env         map[string]string
	Timezone    string
	RunAs       string
	Reason      string
	Warnings    []string
}

// Report accounts for every line and task-creation outcome in one import.
type Report struct {
	Lines    []Line
	Read     int
	Jobs     int
	Skipped  int
	Declined int
	Errors   int
	Created  int
	Failed   int
	Warnings []string
}

type scanContext struct {
	timezone string
	env      map[string]string
	shell    string
}

var reAssignment = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*=(.*)$`)

// ScanCrontab reads a conventional five-field user crontab.
func ScanCrontab(r io.Reader) (Report, error) {
	return ScanCrontabWithOptions(r, ScanOptions{})
}

// ScanCrontabWithOptions reads a crontab using an explicit timing and user layout.
func ScanCrontabWithOptions(r io.Reader, opts ScanOptions) (Report, error) {
	if opts.Dialect == "" {
		opts.Dialect = DialectUnix
	}
	if opts.Dialect != DialectUnix && opts.Dialect != DialectQuartz {
		return Report{}, fmt.Errorf("cron: import dialect must be unix or quartz, got %q", opts.Dialect)
	}
	ctx := scanContext{env: map[string]string{"SHELL": "/bin/sh"}, shell: "/bin/sh"}
	var rep Report
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		rep.Read++
		line := convertLine(rep.Read, sc.Text(), opts, &ctx, &rep)
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

func convertLine(num int, raw string, opts ScanOptions, ctx *scanContext, rep *Report) Line {
	line := Line{Number: num, Raw: raw}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		line.Kind = LineSkipped
		return line
	}
	if m := reAssignment.FindStringSubmatch(trimmed); m != nil {
		return applyAssignment(line, m[1], m[2], ctx, rep)
	}

	expr, runAs, command, ok := splitTiming(trimmed, opts)
	if !ok {
		line.Kind = LineError
		line.Reason = "no command follows the schedule"
		return line
	}
	line.Expr, line.RunAs = expr, runAs
	line.Timezone = ctx.timezone
	line.Env = cloneEnv(ctx.env)
	line.CommandText, line.Stdin, line.HasStdin = splitPercent(command)
	if strings.TrimSpace(line.CommandText) == "" {
		line.Kind = LineError
		line.Reason = "no command follows the schedule"
		return line
	}
	if ctx.shell == "" {
		line.Kind = LineError
		line.Reason = "SHELL is empty"
		return line
	}
	line.Command = ctx.shell
	line.Args = []string{"-c", line.CommandText}

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
	}
	return line
}

func applyAssignment(line Line, name, rawValue string, ctx *scanContext, rep *Report) Line {
	line.Kind = LineSkipped
	value, err := assignmentValue(rawValue)
	if err != nil {
		line.Kind = LineError
		line.Reason = err.Error()
		return line
	}
	switch name {
	case "CRON_TZ":
		ctx.timezone = value
	case "MAILTO", "MAILFROM":
		note := fmt.Sprintf("line %d: %s=%s is not carried across because cron mail delivery is deferred to notification support", line.Number, name, value)
		line.Warnings = append(line.Warnings, note)
		rep.Warnings = append(rep.Warnings, note)
	case "LOGNAME":
		note := fmt.Sprintf("line %d: LOGNAME cannot be overridden by a crontab and was ignored", line.Number)
		line.Warnings = append(line.Warnings, note)
		rep.Warnings = append(rep.Warnings, note)
	default:
		ctx.env[name] = value
		if name == "SHELL" {
			ctx.shell = value
		}
	}
	return line
}

func assignmentValue(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	if value[0] != '\'' && value[0] != '"' {
		return value, nil
	}
	if len(value) < 2 || value[len(value)-1] != value[0] {
		return "", fmt.Errorf("cron: environment value has an unmatched quote")
	}
	return value[1 : len(value)-1], nil
}

func cloneEnv(env map[string]string) map[string]string {
	if len(env) == 0 {
		return nil
	}
	out := make(map[string]string, len(env))
	for key, value := range env {
		out[key] = value
	}
	return out
}

func splitTiming(s string, opts ScanOptions) (expr, runAs, command string, ok bool) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return "", "", "", false
	}
	n := 5
	if opts.Dialect == DialectQuartz {
		n = 6
	}
	if strings.HasPrefix(fields[0], "@") {
		n = 1
	}
	prefixFields := n
	if opts.System {
		prefixFields++
	}
	if len(fields) <= prefixFields {
		return strings.Join(fields, " "), "", "", false
	}
	expr = strings.Join(fields[:n], " ")
	if opts.System {
		runAs = fields[n]
	}
	idx := 0
	for i := 0; i < prefixFields; i++ {
		rel := strings.Index(s[idx:], fields[i])
		if rel < 0 {
			return expr, runAs, "", false
		}
		idx += rel + len(fields[i])
	}
	return expr, runAs, strings.TrimSpace(s[idx:]), true
}

func splitPercent(command string) (commandText, stdin string, present bool) {
	var before, after strings.Builder
	target := &before
	for i := 0; i < len(command); i++ {
		if command[i] == '\\' && i+1 < len(command) && command[i+1] == '%' {
			target.WriteByte('%')
			i++
			continue
		}
		if command[i] == '%' {
			if !present {
				present = true
				target = &after
			} else {
				after.WriteByte('\n')
			}
			continue
		}
		target.WriteByte(command[i])
	}
	return before.String(), after.String(), present
}
