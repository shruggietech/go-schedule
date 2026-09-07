package commandexample

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/commandline"
)

func TestSuggestionExactPlatformMappings(t *testing.T) {
	tests := []struct {
		platform, display, program string
		args                       []string
	}{
		{"windows", "cmd.exe /d /c ver", "cmd.exe", []string{"/d", "/c", "ver"}},
		{"macos", "/usr/bin/sw_vers", "/usr/bin/sw_vers", nil},
		{"darwin", "/usr/bin/sw_vers", "/usr/bin/sw_vers", nil},
		{"linux", "uname -a", "uname", []string{"-a"}},
	}
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			got, ok := ForPlatform(tt.platform)
			if !ok || got.Display != tt.display || got.Program != tt.program || strings.Join(got.Args, "\x00") != strings.Join(tt.args, "\x00") {
				t.Fatalf("ForPlatform(%q)=%+v,%v", tt.platform, got, ok)
			}
			parsed, err := commandline.Parse(got.Display)
			if err != nil || parsed.Program != got.Program || strings.Join(parsed.Args, "\x00") != strings.Join(got.Args, "\x00") {
				t.Fatalf("parsed=%+v err=%v suggestion=%+v", parsed, err, got)
			}
			if !got.Safe() || got.Explanation == "" {
				t.Fatalf("unsafe or unexplained mapping: %+v", got)
			}
		})
	}
	if _, ok := ForPlatform("freebsd"); ok {
		t.Fatal("unsupported platform received an executable suggestion")
	}
}

func TestCurrentPlatformSuggestionExecutesPromptlyAndProducesRecognizableOutput(t *testing.T) {
	suggestion, ok := ForPlatform(runtime.GOOS)
	if !ok {
		t.Skip("current platform is outside the packaged support matrix")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := runSuggestion(ctx, suggestion)
	if err != nil {
		t.Fatalf("safe suggestion: %v; output=%q", err, output)
	}
	if ctx.Err() != nil {
		t.Fatalf("safe suggestion exceeded five seconds: %v", ctx.Err())
	}
	if !suggestion.Recognizes(output) {
		t.Fatalf("output was not recognizable for %s: %q", runtime.GOOS, output)
	}
}

func TestDocumentationContainsExactMappingsAndNoRetiredPrimaryExamples(t *testing.T) {
	contents, err := os.ReadFile("../../docs/gui-fields.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(contents)
	for _, platform := range []string{"windows", "macos", "linux"} {
		suggestion, _ := ForPlatform(platform)
		if !strings.Contains(text, "`"+suggestion.Display+"`") {
			t.Errorf("documentation is missing %s mapping %q", platform, suggestion.Display)
		}
	}
	if strings.Contains(text, "python -m http.server") {
		t.Error("documentation still contains the retired listener example")
	}
}
