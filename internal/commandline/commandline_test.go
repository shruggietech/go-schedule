package commandline

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestParsePortableGrammar(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Invocation
	}{
		{"simple", `python -m http.server`, Invocation{"python", []string{"-m", "http.server"}}},
		{"unicode whitespace", "program\talpha\u2003βeta", Invocation{"program", []string{"alpha", "βeta"}}},
		{"windows path", `"C:\Program Files\Tool\tool.exe" --name "Ada Lovelace"`, Invocation{`C:\Program Files\Tool\tool.exe`, []string{"--name", "Ada Lovelace"}}},
		{"posix path and backslash", `/usr/bin/printf '%s\n' 'hello world'`, Invocation{"/usr/bin/printf", []string{`%s\n`, "hello world"}}},
		{"empty and repeated", `program --tag one --tag two ''`, Invocation{"program", []string{"--tag", "one", "--tag", "two", ""}}},
		{"adjacent segments", `program pre"two words"post`, Invocation{"program", []string{"pretwo wordspost"}}},
		{"escaped spaces and quotes", `program a\ b say\"hi\" it\'s`, Invocation{"program", []string{"a b", `say"hi"`, "it's"}}},
		{"ordinary backslashes", `C:\Tools\app.exe \\server\share C:\trail\`, Invocation{`C:\Tools\app.exe`, []string{`\\server\share`, `C:\trail\`}}},
		{"double quote escape", `program "say \"hello\""`, Invocation{"program", []string{`say "hello"`}}},
		{"literal multiline", "program \"first\r\nsecond\"", Invocation{"program", []string{"first\r\nsecond"}}},
		{"escaped newline", "program first\\\nsecond", Invocation{"program", []string{"first\nsecond"}}},
		{"shell punctuation literal", `program '$HOME' %PATH% '|' '>' '*.txt' ';' '&' '#note'`, Invocation{"program", []string{"$HOME", "%PATH%", "|", ">", "*.txt", ";", "&", "#note"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.text)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Parse(%q) = %#v, want %#v", tt.text, got, tt.want)
			}
		})
	}
}

func TestParseRejectsInvalidInputWithLocation(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"blank", " \t\n", "enter a program"},
		{"empty program", `'' arg`, "program cannot be empty"},
		{"single quote", "program\n'unclosed", "single quote opened at line 2, column 1"},
		{"double quote", `program "unclosed`, "double quote opened at line 1, column 9"},
		{"nul", "program one\x00two", "NUL is not allowed at line 1, column 12"},
		{"invalid UTF-8", string([]byte{'p', 'r', 'o', 'g', 'r', 'a', 'm', ' ', 0x82}), "not valid UTF-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.text)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Parse(%q) error = %v, want containing %q", tt.text, err, tt.want)
			}
		})
	}
}

func TestFormatRoundTripsExactInvocation(t *testing.T) {
	tests := []Invocation{
		{Program: "program"},
		{Program: `C:\Program Files\Tool\tool.exe`, Args: []string{"--name", "Ada Lovelace", ""}},
		{Program: "/opt/工具/run", Args: []string{"héllo 世界", "--tag", "one", "--tag", "two"}},
		{Program: "program", Args: []string{"tabs\there", "lines\r\nhere", `say "hello"`, "it's", `both ' and " quotes`, `C:\trail\`, "after-trailing-slash"}},
		{Program: "program", Args: []string{"$HOME", "%PATH%", "|", ">", "*.txt", ";", "&"}},
	}
	for _, want := range tests {
		formatted, err := Format(want.Program, want.Args)
		if err != nil {
			t.Fatalf("Format(%#v): %v", want, err)
		}
		got, err := Parse(formatted)
		if err != nil {
			t.Fatalf("Parse(Format(%#v)=%q): %v", want, formatted, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip = %#v, want %#v (formatted %q)", got, want, formatted)
		}
		again, err := Format(got.Program, got.Args)
		if err != nil || again != formatted {
			t.Fatalf("canonical format = %q, %v; want %q", again, err, formatted)
		}
	}
}

func TestFormatRejectsProcessInvalidInvocation(t *testing.T) {
	for _, tt := range []struct {
		program string
		args    []string
		want    string
	}{
		{"", nil, "program cannot be empty"},
		{"bad\x00program", nil, "program contains NUL"},
		{"program", []string{"bad\x00argument"}, "argument 1 contains NUL"},
		{string([]byte{0x82}), nil, "program is not valid UTF-8"},
		{"program", []string{string([]byte{0x82})}, "argument 1 is not valid UTF-8"},
	} {
		if _, err := Format(tt.program, tt.args); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Fatalf("Format(%q, %#v) error = %v, want containing %q", tt.program, tt.args, err, tt.want)
		}
	}
}

func TestQuoteDisplayMakesBoundariesVisible(t *testing.T) {
	want := `"empty=\"\" tab=\"a\tb\" line=\"a\r\nb\" slash=\"C:\\x\" spaces=\u00a0\u2003"`
	got := QuoteDisplay("empty=\"\" tab=\"a\tb\" line=\"a\r\nb\" slash=\"C:\\x\" spaces=\u00a0\u2003")
	if got != want {
		t.Fatalf("QuoteDisplay = %s, want %s", got, want)
	}
}

func FuzzFormatRoundTrip(f *testing.F) {
	seeds := []string{"", "simple", "hello world", "héllo 世界", "tabs\there", "line\r\nbreak", `C:\Program Files\Tool`, `say "hi"`, "it's", `both ' and "`}
	for _, program := range []string{"program", `C:\Tools\app.exe`, "/opt/工具/run"} {
		for _, seed := range seeds {
			f.Add(program, seed, seed+"-2")
		}
	}
	f.Fuzz(func(t *testing.T, program, first, second string) {
		if program == "" || strings.ContainsRune(program, '\x00') || strings.ContainsRune(first, '\x00') || strings.ContainsRune(second, '\x00') {
			t.Skip()
		}
		if !utf8.ValidString(program) || !utf8.ValidString(first) || !utf8.ValidString(second) {
			if _, err := Format(program, []string{first, second}); err == nil {
				t.Fatal("Format accepted invalid UTF-8")
			}
			return
		}
		want := Invocation{Program: program, Args: []string{first, second}}
		text, err := Format(want.Program, want.Args)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Parse(text)
		if err != nil {
			t.Fatalf("Parse(%q): %v", text, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip = %#v, want %#v", got, want)
		}
	})
}
