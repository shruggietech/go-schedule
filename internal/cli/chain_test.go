package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestChainCreateRequiresCompleteValidRelationship(t *testing.T) {
	for name, args := range map[string][]string{
		"missing source": {"--target", "b", "--on", "success"},
		"missing target": {"--source", "a", "--on", "success"},
		"same task":      {"--source", "a", "--target", "a", "--on", "success"},
		"bad outcome":    {"--source", "a", "--target", "b", "--on", "maybe"},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := chainCreate()
			cmd.SetArgs(args)
			err := cmd.Execute()
			if !errors.Is(err, errUsage) {
				t.Fatalf("error=%v, want usage error", err)
			}
		})
	}
}

func TestChainUpdateRequiresAtLeastOneChangedField(t *testing.T) {
	cmd := chainUpdate()
	cmd.SetArgs([]string{"chain-id"})
	if err := cmd.Execute(); !errors.Is(err, errUsage) {
		t.Fatalf("error=%v, want usage error", err)
	}
}

func TestChainCommandsExposeDocumentedLifecycleAndFlags(t *testing.T) {
	root := newChainCmd()
	want := map[string]bool{"create": false, "list": false, "show": false, "update": false, "rm": false}
	for _, command := range root.Commands() {
		if _, exists := want[command.Name()]; exists {
			want[command.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing chain subcommand %s", name)
		}
	}
	create := chainCreate()
	for _, flag := range []string{"source", "target", "on"} {
		if create.Flags().Lookup(flag) == nil {
			t.Errorf("create missing --%s", flag)
		}
	}
	if usage := create.Flags().Lookup("on").Usage; !strings.Contains(usage, "success") || !strings.Contains(usage, "failure") || !strings.Contains(usage, "any") {
		t.Fatalf("--on usage does not name supported outcomes: %q", usage)
	}
}
