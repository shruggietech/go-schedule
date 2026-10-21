package clievaluate

import (
	"slices"
	"strings"
)

func generatedNegativeSteps(commands []Command) []step {
	result := []step{}
	for _, command := range commands {
		if !command.Runnable || command.Hidden || strings.HasPrefix(command.Path, "completion") || command.Path == "help" {
			continue
		}
		args := append(strings.Fields(command.Path), "--qa-unknown-flag")
		arityArgs := append([]string{"--ipc-path", "${ipc}"}, strings.Fields(command.Path)...)
		arityArgs = append(arityArgs, "qa-unexpected-1", "qa-unexpected-2", "qa-unexpected-3")
		result = append(result, step{Spec: CaseSpec{ID: "negative.arity." + commandCaseID(command.Path), CommandPath: command.Path, ScenarioClass: "arity", ExecutionKind: "command_harness", Platforms: []string{"windows", "linux"}, Argv: arityArgs, ExpectedExit: 2, Checks: []CheckSpec{{Name: "unexpected positional input rejected", Expected: "true"}, {Name: "invalid invocation preserves state", Expected: "true"}}}})
		result = append(result, step{Spec: CaseSpec{ID: "negative.unknown-flag." + commandCaseID(command.Path), CommandPath: command.Path, ScenarioClass: "unknown-flag", ExecutionKind: "command_harness", Platforms: []string{"windows", "linux"}, Argv: args, ExpectedExit: 2, Checks: []CheckSpec{{Name: "actionable unknown flag diagnostic", Expected: "true"}, {Name: "rejected input preserves state", Expected: "true"}}}, OutputContains: "unknown flag"})
		for _, flag := range command.Flags {
			if flag.Scope == "inherited" || flag.Name == "help" || flag.Name == "version" {
				continue
			}
			value := ""
			switch flag.Type {
			case "bool":
				value = "not-a-bool"
			case "int", "int64", "uint", "uint64", "duration":
				value = "not-a-number"
			default:
				continue
			}
			args := append(strings.Fields(command.Path), "--"+flag.Name+"="+value)
			result = append(result, step{Spec: CaseSpec{ID: "negative.flag-type." + commandCaseID(command.Path) + "." + flag.Name, CommandPath: command.Path, ScenarioClass: "invalid-flag-value", ExecutionKind: "command_harness", Platforms: []string{"windows", "linux"}, Argv: args, ExpectedExit: 2, Checks: []CheckSpec{{Name: "actionable invalid flag diagnostic", Expected: "true"}, {Name: "rejected input preserves state", Expected: "true"}}, Coverage: []string{command.Path + " --" + flag.Name}}, OutputContains: "invalid"})
		}
	}
	return result
}

// BuildManifest binds versioned behavior to the current actual command tree.
func BuildManifest() (Manifest, []step, error) {
	commands, err := ActualCommands()
	if err != nil {
		return Manifest{}, nil, err
	}
	steps := append(metadataSteps(commands), behavioralSteps()...)
	steps = append(steps, generatedNegativeSteps(commands)...)
	steps = append(steps, additionalSteps()...)
	deferred := []step{}
	active := []step{}
	for _, s := range steps {
		if s.Hook == "identity-reset" {
			deferred = append(deferred, s)
		} else {
			active = append(active, s)
		}
	}
	steps = append(active, deferred...)
	human, humanDeferred := []step{}, []step{}
	for _, definition := range behavioralSteps() {
		definition.Human = true
		definition.Spec.ID += ".human"
		definition.Spec.Argv = append([]string{}, definition.Spec.Argv...)
		definition.Spec.Argv = slices.DeleteFunc(definition.Spec.Argv, func(arg string) bool { return arg == "--json" })
		definition.JSON = false
		definition.Spec.Coverage = coverageFor(definition.Spec.CommandPath, definition.Spec.Argv)
		if definition.Hook == "identity-reset" {
			humanDeferred = append(humanDeferred, definition)
		} else {
			human = append(human, definition)
		}
	}
	steps = append(steps, human...)
	steps = append(steps, humanDeferred...)
	manifest := Manifest{SchemaVersion: Version, Kind: "go-schedule.cli-qa-manifest", Commands: commands, Cases: []CaseSpec{}}
	for _, s := range steps {
		manifest.Cases = append(manifest.Cases, s.Spec)
	}
	return manifest, steps, nil
}

// MissingFlagCoverage reports explicit uncovered work rather than inferring it from help.
func MissingFlagCoverage(m Manifest) []string {
	covered := map[string]bool{}
	for _, c := range m.Cases {
		for _, item := range c.Coverage {
			covered[item] = true
		}
	}
	missing := []string{}
	for _, c := range m.Commands {
		for _, flag := range c.Flags {
			key := c.Path + " --" + flag.Name
			if flag.Scope == "inherited" {
				key = " --" + flag.Name
			}
			if !covered[key] {
				missing = append(missing, key)
			}
		}
	}
	slices.Sort(missing)
	return slices.Compact(missing)
}

func commandCaseID(path string) string {
	if path == "" {
		return "root"
	}
	return strings.ReplaceAll(path, " ", ".")
}
