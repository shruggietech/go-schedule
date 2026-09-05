package cli

import (
	"bytes"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestPrintTriggerSetSecretIsOrderedOneCommandPerLine(t *testing.T) {
	old := jsonOut
	jsonOut = false
	defer func() { jsonOut = old }()
	var out bytes.Buffer
	result := server.TriggerSetSecretResponse{Members: []server.TriggerSetSecretMember{{Position: 1, Command: "gosched trigger fire first"}, {Position: 3, Command: "gosched trigger fire third"}}}
	if err := printTriggerSetSecret(&out, result); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "gosched trigger fire first\ngosched trigger fire third\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestTriggerSetCommandRegistersCompleteLifecycle(t *testing.T) {
	cmd := newTriggerSetCmd()
	want := map[string]bool{"create": true, "list": true, "show": true, "retarget": true, "enable": true, "disable": true, "reveal": true, "rotate": true, "rm": true}
	for _, child := range cmd.Commands() {
		delete(want, child.Name())
	}
	if len(want) != 0 {
		t.Fatalf("missing commands: %v", want)
	}
}
