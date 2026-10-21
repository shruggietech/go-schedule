package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCommandInventoryUsesRealInitializedTree(t *testing.T) {
	data, err := CommandInventoryJSON()
	if err != nil {
		t.Fatal(err)
	}
	var records []struct {
		Path    string   `json:"path"`
		Aliases []string `json:"aliases"`
		Flags   []struct {
			Name  string `json:"name"`
			Scope string `json:"scope"`
		} `json:"flags"`
	}
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatal(err)
	}
	paths := map[string]bool{}
	alias := false
	shadow := false
	for _, r := range records {
		if paths[r.Path] {
			t.Fatal("duplicate command path")
		}
		paths[r.Path] = true
		if r.Path == "watcher" {
			alias = strings.Join(r.Aliases, ",") == "watch"
		}
		if r.Path == "notification channel add" {
			for _, f := range r.Flags {
				if f.Name == "endpoint" {
					shadow = f.Scope == "local"
				}
			}
		}
	}
	for _, path := range []string{"", "runs", "alerts", "systems search", "completion bash", "help"} {
		if !paths[path] {
			t.Fatalf("missing real/generated command %q", path)
		}
	}
	if !alias || !shadow {
		t.Fatal("alias or effective flag shadow missing")
	}
	other, err := CommandInventoryJSON()
	if err != nil || string(other) != string(data) {
		t.Fatal("inventory is not deterministic")
	}
}
