package cli

import "testing"

func TestWatcherCommandRegistersCompleteLifecycle(t *testing.T) {
	command := newWatcherCmd()
	want := map[string]bool{"create": true, "list": true, "show": true, "update": true, "enable": true, "disable": true, "rm": true}
	for _, child := range command.Commands() {
		delete(want, child.Name())
	}
	if len(want) != 0 {
		t.Fatalf("missing commands: %v", want)
	}
}

func TestWatcherKindRejectsUnknownValue(t *testing.T) {
	var value watcherKindValue
	if err := value.Set("socket"); err == nil {
		t.Fatal("expected invalid kind error")
	}
}
