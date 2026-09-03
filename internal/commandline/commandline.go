// Package commandline translates the desktop editor's portable command text
// into the direct program-plus-arguments model used by tasks. It is deliberately
// lexical only: no shell expansion, substitution, piping, or redirection occurs.
package commandline

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Invocation is the exact direct process boundary represented by editor text.
type Invocation struct {
	Program string
	Args    []string
}

// SyntaxError identifies an invalid character or unmatched quote in editor
// text. Line and Column are one-based Unicode code-point positions.
type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("command line: %s at line %d, column %d", e.Message, e.Line, e.Column)
}

type quoteState rune

const (
	quoteNone   quoteState = 0
	quoteSingle quoteState = '\''
	quoteDouble quoteState = '"'
)

// Parse applies go-schedule's platform-independent direct-command grammar.
func Parse(text string) (Invocation, error) {
	if !utf8.ValidString(text) {
		return Invocation{}, fmt.Errorf("command line: text is not valid UTF-8")
	}
	runes := []rune(text)
	lines, columns := positions(runes)

	var (
		words        []string
		word         strings.Builder
		started      bool
		quote        quoteState
		quoteStarted int
	)
	flush := func() {
		if !started {
			return
		}
		words = append(words, word.String())
		word.Reset()
		started = false
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == 0 {
			return Invocation{}, syntaxErrorAt(lines, columns, i, "NUL is not allowed")
		}

		switch quote {
		case quoteSingle:
			if r == '\'' {
				quote = quoteNone
			} else {
				word.WriteRune(r)
			}
			continue
		case quoteDouble:
			switch {
			case r == '"':
				quote = quoteNone
			case r == '\\' && i+1 < len(runes) && runes[i+1] == '"':
				word.WriteRune('"')
				i++
			default:
				word.WriteRune(r)
			}
			continue
		}

		switch {
		case unicode.IsSpace(r):
			flush()
		case r == '\'' || r == '"':
			started = true
			quoteStarted = i
			quote = quoteState(r)
		case r == '\\' && i+1 < len(runes) && (unicode.IsSpace(runes[i+1]) || runes[i+1] == '\'' || runes[i+1] == '"'):
			started = true
			word.WriteRune(runes[i+1])
			i++
		default:
			started = true
			word.WriteRune(r)
		}
	}

	if quote != quoteNone {
		name := "single quote"
		if quote == quoteDouble {
			name = "double quote"
		}
		return Invocation{}, syntaxErrorAt(lines, columns, quoteStarted, name+" opened")
	}
	flush()
	if len(words) == 0 {
		return Invocation{}, &SyntaxError{Line: 1, Column: 1, Message: "enter a program"}
	}
	if words[0] == "" {
		return Invocation{}, &SyntaxError{Line: 1, Column: 1, Message: "program cannot be empty"}
	}

	invocation := Invocation{Program: words[0]}
	if len(words) > 1 {
		invocation.Args = words[1:]
	}
	return invocation, nil
}

// Format returns a stable, lossless editor representation of an invocation.
func Format(program string, args []string) (string, error) {
	if program == "" {
		return "", fmt.Errorf("command line: program cannot be empty")
	}
	if !utf8.ValidString(program) {
		return "", fmt.Errorf("command line: program is not valid UTF-8")
	}
	if strings.ContainsRune(program, 0) {
		return "", fmt.Errorf("command line: program contains NUL")
	}
	for i, arg := range args {
		if !utf8.ValidString(arg) {
			return "", fmt.Errorf("command line: argument %d is not valid UTF-8", i+1)
		}
		if strings.ContainsRune(arg, 0) {
			return "", fmt.Errorf("command line: argument %d contains NUL", i+1)
		}
	}

	parts := make([]string, 0, len(args)+1)
	parts = append(parts, formatWord(program))
	for _, arg := range args {
		parts = append(parts, formatWord(arg))
	}
	return strings.Join(parts, " "), nil
}

// QuoteDisplay renders an exact value for the launch preview. It is display
// only and is never fed back to Parse or a process launcher.
func QuoteDisplay(value string) string { return strconv.QuoteToGraphic(value) }

func formatWord(value string) string {
	if value == "" {
		return "''"
	}
	if isBareWord(value) {
		return value
	}
	// A trailing backslash inside double quotes would escape the closing quote
	// under this grammar, so use another reversible form for that one case.
	if !strings.ContainsRune(value, '"') && !strings.HasSuffix(value, `\`) {
		return `"` + value + `"`
	}
	if !strings.ContainsRune(value, '\'') {
		return "'" + value + "'"
	}
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func isBareWord(value string) bool {
	// Format joins words with a space. A bare trailing backslash would escape
	// that separator and merge the following argument.
	if strings.HasSuffix(value, `\`) {
		return false
	}
	for _, r := range value {
		if unicode.IsSpace(r) || r == '\'' || r == '"' {
			return false
		}
	}
	return true
}

func positions(runes []rune) ([]int, []int) {
	lines := make([]int, len(runes))
	columns := make([]int, len(runes))
	line, column := 1, 1
	for i, r := range runes {
		lines[i], columns[i] = line, column
		if r == '\n' {
			line, column = line+1, 1
		} else {
			column++
		}
	}
	return lines, columns
}

func syntaxErrorAt(lines, columns []int, index int, message string) error {
	if index < 0 || index >= len(lines) {
		return &SyntaxError{Line: 1, Column: 1, Message: message}
	}
	return &SyntaxError{Line: lines[index], Column: columns[index], Message: message}
}
