// Package cron converts between crontab expressions and this scheduler's
// human-readable schedule phrases. It exists at the boundary only: cron is an
// interchange format for import and export, never an authoring syntax. Nothing
// in this package produces a recurrence directly — an expression becomes the
// phrase a user would have typed, and that phrase goes through
// internal/schedule like any other, so a preview cannot disagree with what is
// stored.
//
// The parser is written here rather than taken from a dependency because the
// work is deciding what *cannot* be represented. A scheduling library normalizes
// exactly the distinctions this package must detect: it will happily accept
// "*/7" and hand back a schedule, when the honest answer is a refusal.
package cron

import (
	"fmt"
	"strconv"
	"strings"
)

// Field is one parsed crontab field: the set of values it matches, and whether
// it was a bare "*" (which means "every", not "all values enumerated" — the
// distinction matters when deciding whether a phrase exists).
type Field struct {
	Values   []int
	Wildcard bool
	// Step is the divisor from a "*/n" form, or 0 when there was none.
	Step int
	// min and max bound the field's range, for validation and step checks.
	min, max int
}

// Spec is a parsed cron timing expression.
type Spec struct {
	Minute, Hour, DOM, Month, DOW Field
	// Shorthand records the "@daily"-style macro this came from, if any, purely
	// so messages can quote what the operator wrote.
	Shorthand string
}

// Unsupported is a named refusal. It is a value rather than an error because a
// refusal is an outcome the caller reports, not a failure of the run: a crontab
// of three @reboot lines converts successfully to three refusals.
type Unsupported struct {
	Input  string
	Reason string
}

func (u Unsupported) String() string { return u.Reason }

// Result is the outcome of converting one expression: exactly one of Spec or
// Unsupported is meaningful, discriminated by OK.
type Result struct {
	Spec Spec
	Bad  Unsupported
	OK   bool
}

// field bounds, in crontab order.
var bounds = [5][2]int{
	{0, 59}, // minute
	{0, 23}, // hour
	{1, 31}, // day of month
	{1, 12}, // month
	{0, 7},  // day of week (both 0 and 7 are Sunday)
}

var monthNames = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

var dowNames = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

// shorthands are the macros with an exact five-field equivalent. @reboot is
// deliberately absent: it has no schedule at all, and is refused by name.
var shorthands = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

// Parse turns a cron timing expression into a Spec, or into a named refusal.
// A malformed expression is an error; a well-formed one this package will not
// represent is a refusal. The distinction is the caller's exit code (a typo is
// the operator's mistake; an unsupported extension is not).
func Parse(expr string) (Result, error) {
	raw := strings.TrimSpace(expr)
	if raw == "" {
		return Result{}, fmt.Errorf("cron: empty expression")
	}

	shorthand := ""
	if strings.HasPrefix(raw, "@") {
		key := strings.ToLower(raw)
		if key == "@reboot" {
			return refuse(raw, "@reboot fires at boot rather than on a schedule, which has no equivalent here"), nil
		}
		expanded, ok := shorthands[key]
		if !ok {
			return Result{}, fmt.Errorf("cron: unknown shorthand %q", raw)
		}
		shorthand, raw = key, expanded
	}

	parts := strings.Fields(raw)
	switch {
	case len(parts) == 6:
		return refuse(expr, "six-field (Quartz-style, seconds-precision) expressions are not supported"), nil
	case len(parts) != 5:
		return Result{}, fmt.Errorf("cron: expected 5 fields, got %d in %q", len(parts), expr)
	}

	// The non-standard day-of-month extensions have no equivalent and must be
	// named rather than silently misread as a literal.
	for i, p := range parts {
		if bad, ok := extensionIn(p); ok {
			return refuse(expr, fmt.Sprintf("the %q extension in the %s field is not supported", bad, fieldName(i))), nil
		}
	}

	spec := Spec{Shorthand: shorthand}
	targets := [5]*Field{&spec.Minute, &spec.Hour, &spec.DOM, &spec.Month, &spec.DOW}
	for i, p := range parts {
		f, err := parseField(p, bounds[i][0], bounds[i][1], i)
		if err != nil {
			return Result{}, fmt.Errorf("cron: %s field: %w", fieldName(i), err)
		}
		*targets[i] = f
	}

	// Cron ORs a restricted day-of-month with a restricted day-of-week; the
	// recurrence model intersects them. Rather than silently changing a weekly
	// job into a handful of runs a year, refuse and say why.
	if !spec.DOM.Wildcard && !spec.DOW.Wildcard {
		return refuse(expr, "restricting both day-of-month and day-of-week means \"either\" in cron, which has no equivalent here"), nil
	}

	return Result{Spec: spec, OK: true}, nil
}

func refuse(input, reason string) Result {
	return Result{Bad: Unsupported{Input: input, Reason: reason}}
}

func fieldName(i int) string {
	return [...]string{"minute", "hour", "day-of-month", "month", "day-of-week"}[i]
}

// extensionIn reports a non-standard extension token in a field.
func extensionIn(p string) (string, bool) {
	upper := strings.ToUpper(p)
	if strings.Contains(p, "#") {
		return "#", true
	}
	// L and W are only extensions when they stand as day specifiers, not when
	// they appear inside a month or weekday name (JUL contains no L... but WED
	// contains no W either; check the standalone forms).
	for _, ext := range []string{"L", "W"} {
		for _, tok := range strings.Split(upper, ",") {
			tok = strings.TrimSpace(tok)
			if tok == ext || strings.HasSuffix(tok, ext) && isNumericPrefix(tok[:len(tok)-1]) {
				return ext, true
			}
		}
	}
	return "", false
}

func isNumericPrefix(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// parseField parses one field: "*", "*/n", a list, ranges, steps, and names.
func parseField(p string, min, max, idx int) (Field, error) {
	f := Field{min: min, max: max}

	if p == "*" {
		f.Wildcard = true
		f.Values = rangeOf(min, max)
		return f, nil
	}

	// "*/n" — the only step form that also stays a wildcard in meaning.
	if step, ok := strings.CutPrefix(p, "*/"); ok {
		n, err := strconv.Atoi(step)
		if err != nil || n < 1 {
			return Field{}, fmt.Errorf("invalid step %q", step)
		}
		f.Wildcard = true
		f.Step = n
		for v := min; v <= max; v += n {
			f.Values = append(f.Values, v)
		}
		return f, nil
	}

	seen := map[int]bool{}
	for _, part := range strings.Split(p, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return Field{}, fmt.Errorf("empty element in %q", p)
		}
		step := 1
		if base, s, ok := strings.Cut(part, "/"); ok {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 {
				return Field{}, fmt.Errorf("invalid step %q", s)
			}
			step, part = n, base
		}
		lo, hi, err := parseRange(part, min, max, idx)
		if err != nil {
			return Field{}, err
		}
		for v := lo; v <= hi; v += step {
			seen[normalize(v, idx)] = true
		}
	}
	for v := min; v <= max; v++ {
		if seen[normalize(v, idx)] {
			f.Values = append(f.Values, normalize(v, idx))
		}
	}
	f.Values = dedupe(f.Values)
	if len(f.Values) == 0 {
		return Field{}, fmt.Errorf("no values matched %q", p)
	}
	return f, nil
}

// normalize folds day-of-week 7 onto 0, since cron accepts both for Sunday.
func normalize(v, idx int) int {
	if idx == 4 && v == 7 {
		return 0
	}
	return v
}

func parseRange(part string, min, max, idx int) (int, int, error) {
	if lo, hi, ok := strings.Cut(part, "-"); ok {
		l, err := parseValue(lo, idx)
		if err != nil {
			return 0, 0, err
		}
		h, err := parseValue(hi, idx)
		if err != nil {
			return 0, 0, err
		}
		if l < min || h > max || l > h {
			return 0, 0, fmt.Errorf("range %q is outside %d-%d", part, min, max)
		}
		return l, h, nil
	}
	v, err := parseValue(part, idx)
	if err != nil {
		return 0, 0, err
	}
	if v < min || v > max {
		return 0, 0, fmt.Errorf("value %d is outside %d-%d", v, min, max)
	}
	return v, v, nil
}

// parseValue resolves a number or a three-letter month/day name.
func parseValue(s string, idx int) (int, error) {
	s = strings.TrimSpace(s)
	if v, err := strconv.Atoi(s); err == nil {
		return v, nil
	}
	key := strings.ToLower(s)
	if len(key) > 3 {
		key = key[:3]
	}
	switch idx {
	case 3:
		if v, ok := monthNames[key]; ok {
			return v, nil
		}
	case 4:
		if v, ok := dowNames[key]; ok {
			return v, nil
		}
	}
	return 0, fmt.Errorf("%q is not a valid value", s)
}

func rangeOf(min, max int) []int {
	out := make([]int, 0, max-min+1)
	for v := min; v <= max; v++ {
		out = append(out, v)
	}
	return out
}

func dedupe(in []int) []int {
	out := in[:0]
	var prev = -1
	for _, v := range in {
		if v != prev {
			out = append(out, v)
			prev = v
		}
	}
	return out
}

// Single reports the field's only value when it has exactly one, which is how
// the phrase builder recognizes a fixed time-of-day or date.
func (f Field) Single() (int, bool) {
	if !f.Wildcard && len(f.Values) == 1 {
		return f.Values[0], true
	}
	return 0, false
}

// EveryValue reports whether the field matches every value in its range with no
// step — a bare "*", or a step of 1.
func (f Field) EveryValue() bool {
	return f.Wildcard && (f.Step == 0 || f.Step == 1)
}
