package clievaluate

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestGeneratedCaseIDsRemainStableWhenCommandTreeGrows(t *testing.T) {
	health := Command{Path: "health", Runnable: true, Aliases: []string{"check"}, Flags: []Flag{{Name: "limit", Type: "int", Scope: "local"}}}
	ids := func(commands []Command) []string {
		steps := append(metadataSteps(commands), generatedNegativeSteps(commands)...)
		result := []string{}
		for _, step := range steps {
			if step.Spec.CommandPath == "health" {
				result = append(result, step.Spec.ID)
			}
		}
		return result
	}
	if before, after := ids([]Command{health}), ids([]Command{{Path: "added", Runnable: true}, health}); !reflect.DeepEqual(before, after) {
		t.Fatalf("existing case IDs changed: %v -> %v", before, after)
	}
}

func TestCompleteManifestAndDriftAreSeparateFromHelp(t *testing.T) {
	m, _, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	if err := InventoryCoverage(m); err != nil {
		t.Fatal(err)
	}
	if len(m.Cases) < 300 {
		t.Fatal("behavior/metadata/negative matrix unexpectedly missing")
	}
	changed := append([]Command{}, m.Commands...)
	changed[0].Short += " changed"
	if err := CheckInventory(changed); err == nil {
		t.Fatal("metadata drift silently accepted")
	}
	copy := m
	copy.Cases = nil
	for _, c := range m.Cases {
		if c.CommandPath != "task edit" || c.ScenarioClass == "help" || c.ScenarioClass == "unknown-flag" || c.ScenarioClass == "arity" {
			copy.Cases = append(copy.Cases, c)
		}
	}
	if err := ValidateManifest(copy); err == nil {
		t.Fatal("help/error-only command falsely counted as behavior")
	}
}

func TestEveryCatalogFlagExistsInActualTree(t *testing.T) {
	m, steps, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	commands := map[string]Command{}
	for _, c := range m.Commands {
		commands[c.Path] = c
	}
	for _, step := range steps {
		if strings.HasPrefix(step.Spec.CommandPath, "@") || step.Spec.ScenarioClass == "unknown-flag" || strings.HasPrefix(step.Spec.CommandPath, "__complete") {
			continue
		}
		flags := map[string]bool{}
		for _, f := range commands[step.Spec.CommandPath].Flags {
			flags[f.Name] = true
		}
		for _, arg := range step.Spec.Argv {
			if arg == "--" {
				break
			}
			if strings.HasPrefix(arg, "--") {
				name := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)[0]
				if !flags[name] {
					t.Errorf("%s references nonexistent --%s", step.Spec.ID, name)
				}
			}
		}
	}
	data, err := Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSchema("manifest", data); err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
}

func TestManifestValidationRejectsNewFlagWithoutBehaviorCoverage(t *testing.T) {
	m, _, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	m.Commands[0].Flags = append(m.Commands[0].Flags, Flag{Name: "new-fixture-option", Scope: "local", Type: "string"})
	if err := ValidateManifest(m); err == nil {
		t.Fatal("new unexercised flag silently accepted")
	}
}
