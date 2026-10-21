package clievaluate

import (
	"slices"
	"strings"
)

type assertion struct{ Path, Want string }
type step struct {
	Human          bool
	Spec           CaseSpec
	Input          string
	Save           []NamedValue
	Assert         []assertion
	Hook           string
	OutputContains string
	JSON           bool
}

func behavioralSteps() []step {
	steps := []step{}
	add := func(id, path, class string, code int, args ...string) *step {
		s := step{Spec: CaseSpec{ID: id, CommandPath: path, ScenarioClass: class, ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: append([]string{"--json", "--ipc-path", "${ipc}"}, args...), ExpectedExit: code, Checks: []CheckSpec{{Name: "required output semantics", Expected: "true"}, {Name: "independent persisted-state contract", Expected: "true"}, {Name: "history audit and log evidence", Expected: "true"}}}, JSON: true}
		steps = append(steps, s)
		return &steps[len(steps)-1]
	}
	create := func(id, path, variable, jsonPath string, args ...string) {
		s := add(id, path, "create", 0, args...)
		s.Save = []NamedValue{{Name: variable, Value: jsonPath}}
	}
	read := func(id, path string, args ...string) { add(id, path, "read", 0, args...) }
	edit := func(id, path, property, want string, args ...string) {
		s := add(id, path, "update", 0, args...)
		s.Assert = []assertion{{property, want}}
		if path == "task enable" || path == "task disable" || path == "group enable" || path == "group disable" {
			s.JSON = false
		}
	}
	negative := func(id, path string, args ...string) { add(id, path, "validation", 2, args...) }
	create("group.create", "group add", "group", "id", "group", "add", "QA 根", "--disabled")
	create("group.child", "group add", "child_group", "id", "group", "add", "QA child", "--parent", "${group}")
	read("group.tree", "group list", "group", "list", "--tree")
	edit("group.rename", "group edit", "groups.${child_group}.name", "QA renamed", "group", "edit", "${child_group}", "--name", "QA renamed", "--parent=")
	edit("group.enable", "group enable", "groups.${group}.enabled", "true", "group", "enable", "${group}")
	edit("group.disable", "group disable", "groups.${group}.enabled", "false", "group", "disable", "${group}")
	negative("group.cycle-atomic", "group edit", "group", "edit", "${group}", "--name", "must not persist", "--parent", "${group}")
	create("task.create", "task add", "task", "task.id", "task", "add", "QA 世界", "--command", "${worker}", "--arg=fixture-child", "--arg=--mode=success", "--arg=must-be-cleared", "--cwd", "${runtime}", "--env=QA_VALUE=hello 世界", "--env=REMOVE_ME=retired", "--stdin=fixture input", "--group", "${group}", "--tz=UTC", "--schedule=every 1 hour", "--syntax=human", "--disabled", "--overlap=queue_one", "--catchup=one", "--missing-date=skip", "--time-basis=wall_clock", "--dst-gap=next_valid", "--dst-overlap=first")
	read("task.filtered-list", "task list", "task", "list", "--group", "${group}", "--state=disabled")
	read("task.show", "task show", "task", "show", "${task}")
	edit("task.edit-inputs", "task edit", "tasks.${task}.name", "QA edited", "task", "edit", "${task}", "--name=QA edited", "--command", "${worker}", "--clear-args", "--clear-env", "--clear-stdin", "--clear-cwd", "--clear-run-as", "--clear-schedule", "--group=", "--tz=UTC")
	edit("task.restore-execution", "task edit", "tasks.${task}.command", "${worker}", "task", "edit", "${task}", "--arg=fixture-child", "--arg=--mode=success", "--arg=space quote \" 世界", "--env=QA_VALUE=owned", "--stdin-file", "${stdin_file}", "--cwd", "${runtime}", "--run-as=", "--schedule=0 0 1 1 *", "--syntax=cron", "--overlap=skip", "--catchup=none", "--missing-date=last_valid", "--time-basis=utc", "--dst-gap=skip", "--dst-overlap=both")
	edit("task.clear-automatic", "task edit", "tasks.${task}.schedule_id", "", "task", "edit", "${task}", "--clear-schedule", "--catchup=one", "--missing-date=next_valid", "--time-basis=elapsed", "--dst-overlap=last")
	edit("task.restore-activation", "task edit", "tasks.${task}.schedule_id", "<nonempty>", "task", "edit", "${task}", "--schedule=every 24 hours", "--time-basis=wall_clock")
	edit("task.enable", "task enable", "tasks.${task}.enabled", "true", "task", "enable", "${task}")
	editedRun := add("task.execute-edited-inputs", "task run-now", "execution", 0, "task", "run-now", "${task}")
	editedRun.JSON = false
	editedRun.Hook = "run-edited"
	edit("task.disable", "task disable", "tasks.${task}.enabled", "false", "task", "disable", "${task}")
	s := add("task.execute-success", "task run-now", "execution", 0, "task", "run-now", "${base_task}")
	s.JSON = false
	s.Hook = "run-success"
	s = add("task.execute-failure", "task run-now", "execution", 0, "task", "run-now", "${failure_task}")
	s.JSON = false
	s.Hook = "run-failure"
	s = add("task.execute-missing", "task run-now", "execution", 0, "task", "run-now", "${missing_task}")
	s.JSON = false
	s.Hook = "run-missing"
	s = add("task.execute-output-cap", "task run-now", "execution", 0, "task", "run-now", "${output_task}")
	s.JSON = false
	s.Hook = "run-cap"
	create("task.one-off-disabled", "task add", "oneoff", "task.id", "task", "add", "QA oneoff", "--command", "${worker}", "--arg=fixture-child", "--at", "${future}", "--disabled")
	create("task.startup-disabled", "task add", "startup", "task.id", "task", "add", "QA startup", "--command", "${worker}", "--arg=fixture-child", "--schedule=@reboot", "--disabled")
	negative("task.bad-timezone", "task edit", "task", "edit", "${task}", "--tz=No/Such_Zone")
	negative("task.conflicting-clear", "task edit", "task", "edit", "${task}", "--clear-env", "--env=X=Y")
	negative("task.bad-input-file", "task edit", "task", "edit", "${task}", "--stdin-file", "${runtime}/missing")
	steps[len(steps)-1].Spec.ExpectedExit = 1
	negative("task.bad-environment", "task edit", "task", "edit", "${task}", "--env=missing-equals")
	negative("task.bad-schedule", "task edit", "task", "edit", "${task}", "--schedule=not a schedule")
	negative("task.bad-at", "task edit", "task", "edit", "${task}", "--at=yesterday")
	read("schedule.preview", "schedule preview", "schedule", "preview", "every 1 hour", "--tz=UTC")
	read("schedule.dst-gap", "schedule preview", "schedule", "preview", "30 2 * * *", "--syntax=cron", "--tz=America/New_York", "--time-basis=wall_clock", "--dst-gap=skip", "--dst-overlap=both", "--missing-date=skip")
	read("schedule.month-end", "schedule preview", "schedule", "preview", "0 9 31 * *", "--syntax=cron", "--tz=UTC", "--missing-date=last_valid")
	negative("schedule.bad-policy", "schedule preview", "schedule", "preview", "every hour", "--dst-gap=invalid")
	read("calendar.window", "calendar", "calendar", "--from=2026-10-21T00:00:00Z", "--to=2026-10-22T00:00:00Z")
	read("runtime.paths", "runtime", "runtime")
	read("health.read", "health", "health")
	read("runs.history", "runs", "runs", "--task", "${base_task}", "--limit=5", "--details")
	read("runs.zero-limit", "runs", "runs", "--task", "${base_task}", "--limit=0")
	read("runs.exact", "runs get", "runs", "get", "${run}")
	read("runs.active", "runs active", "runs", "active")
	read("alerts.list", "alerts", "alerts", "--unacked")
	s = add("alerts.ack", "alerts ack", "update", 0, "alerts", "ack", "${alert}")
	s.JSON = false
	read("logs.read", "logs", "logs", "--severity=info", "--limit=8")
	read("logs.metadata", "logs", "logs", "--metadata")
	s = add("systems.summary", "systems summary", "read", 0, "systems", "summary", "--all")
	s.Spec.ExecutionKind = "command_harness"
	s = add("systems.search", "systems search", "read", 0, "systems", "search", "QA", "--all", "--kind=task", "--limit=5")
	s.Spec.ExecutionKind = "command_harness"
	read("search.selected", "search", "search", "QA", "--kind=task", "--limit=5")
	negative("search.invalid-kind", "search", "search", "QA", "--kind=invalid")
	create("chain.create", "chain create", "chain", "id", "chain", "create", "--source", "${base_task}", "--target", "${second_task}", "--on=success")
	read("chain.list", "chain list", "chain", "list")
	read("chain.show", "chain show", "chain", "show", "${chain}")
	edit("chain.update", "chain update", "chains.${chain}.on_outcome", "any", "chain", "update", "${chain}", "--source", "${base_task}", "--target", "${second_task}", "--on=any")
	negative("chain.cycle", "chain create", "chain", "create", "--source", "${second_task}", "--target", "${base_task}", "--on=success")
	s = add("chain.propagation", "task run-now", "execution", 0, "task", "run-now", "${base_task}")
	s.JSON = false
	s.Hook = "chain-propagation"
	add("chain.remove", "chain rm", "delete", 0, "chain", "rm", "${chain}").JSON = false
	create("trigger.create", "trigger create", "trigger", "trigger.id", "trigger", "create", "--name=QA trigger", "--task", "${base_task}")
	s = add("trigger.reveal", "trigger show", "read", 0, "trigger", "show", "${trigger}", "--reveal-key")
	s.Save = []NamedValue{{Name: "trigger_key", Value: "key"}}
	s.Hook = "trigger-secret"
	read("trigger.list", "trigger list", "trigger", "list")
	read("trigger.show", "trigger show", "trigger", "show", "${trigger}")
	add("trigger.fire-first", "trigger fire", "execution", 0, "trigger", "fire", "${trigger_key}").JSON = false
	add("trigger.fire-reusable", "trigger fire", "execution", 0, "trigger", "fire", "${trigger_key}").JSON = false
	edit("trigger.rename", "trigger update", "triggers.${trigger}.name", "QA trigger renamed", "trigger", "update", "${trigger}", "--name=QA trigger renamed", "--task", "${second_task}")
	add("trigger.disable", "trigger disable", "update", 0, "trigger", "disable", "${trigger}")
	add("trigger.enable", "trigger enable", "update", 0, "trigger", "enable", "${trigger}")
	s = add("trigger.rotate", "trigger rotate", "update", 0, "trigger", "rotate", "${trigger}")
	s.Save = []NamedValue{{Name: "rotated_key", Value: "key"}}
	s.Hook = "trigger-secret"
	add("trigger.old-key-rejected", "trigger fire", "validation", 1, "trigger", "fire", "${trigger_key}")
	add("trigger.rotated-key", "trigger fire", "execution", 0, "trigger", "fire", "${rotated_key}").JSON = false
	add("trigger.remove", "trigger rm", "delete", 0, "trigger", "rm", "${trigger}").JSON = false
	add("trigger.deleted-key", "trigger fire", "validation", 1, "trigger", "fire", "${rotated_key}")
	create("trigger-set.create", "trigger set create", "trigger_set", "trigger_set.id", "trigger", "set", "create", "--name=QA set", "--task", "${base_task}", "--count=2", "--disabled")
	read("trigger-set.list", "trigger set list", "trigger", "set", "list")
	read("trigger-set.show", "trigger set show", "trigger", "set", "show", "${trigger_set}")
	read("trigger-set.reveal", "trigger set reveal", "trigger", "set", "reveal", "${trigger_set}")
	add("trigger-set.rotate", "trigger set rotate", "update", 0, "trigger", "set", "rotate", "${trigger_set}")
	add("trigger-set.retarget", "trigger set retarget", "update", 0, "trigger", "set", "retarget", "${trigger_set}", "--task", "${second_task}")
	add("trigger-set.enable", "trigger set enable", "update", 0, "trigger", "set", "enable", "${trigger_set}")
	add("trigger-set.disable", "trigger set disable", "update", 0, "trigger", "set", "disable", "${trigger_set}")
	negative("trigger-set.count-zero", "trigger set create", "trigger", "set", "create", "--name=invalid", "--task", "${base_task}", "--count=0")
	negative("trigger-set.count-hundred", "trigger set create", "trigger", "set", "create", "--name=invalid", "--task", "${base_task}", "--count=100")
	create("trigger-set.count-upper", "trigger set create", "large_set", "trigger_set.id", "trigger", "set", "create", "--name=QA upper", "--task", "${base_task}", "--count=99", "--disabled")
	add("trigger-set.remove-upper", "trigger set rm", "delete", 0, "trigger", "set", "rm", "${large_set}")
	add("trigger-set.remove", "trigger set rm", "delete", 0, "trigger", "set", "rm", "${trigger_set}")
	create("watcher.create", "watcher create", "watcher", "id", "watcher", "create", "--name=QA watcher", "--task", "${base_task}", "--path", "${watch_root}", "--kind=directory", "--pattern=*.txt", "--recursive", "--debounce=25ms", "--stability=25ms")
	read("watcher.list", "watcher list", "watcher", "list")
	read("watcher.show", "watcher show", "watcher", "show", "${watcher}")
	add("watcher.update", "watcher update", "update", 0, "watcher", "update", "${watcher}", "--name=QA watcher edited", "--path", "${watch_root}", "--kind=directory", "--task", "${base_task}", "--pattern=*.txt", "--recursive=false", "--debounce=50ms", "--stability=50ms")
	s = add("watcher.event", "watcher show", "execution", 0, "watcher", "show", "${watcher}")
	s.Hook = "watcher-event"
	add("watcher.disable", "watcher disable", "update", 0, "watcher", "disable", "${watcher}")
	add("watcher.enable", "watcher enable", "update", 0, "watcher", "enable", "${watcher}")
	add("watcher.remove", "watcher rm", "delete", 0, "watcher", "rm", "${watcher}").JSON = false
	read("cron.explain", "cron explain", "cron", "explain", "0 9 * * 1-5", "--timezone=UTC", "--count=3")
	read("cron.to-human", "cron convert", "cron", "convert", "--to=human", "0 9 * * 1-5")
	read("cron.to-cron", "cron convert", "cron", "convert", "--to=cron", "every day at 09:00")
	s = add("cron.import-dry-run", "cron import", "read", 0, "cron", "import", "--file", "${cron_file}", "--dialect=unix", "--dry-run", "--timezone=UTC", "--count=2")
	s.JSON = false
	s = add("cron.import-stdin", "cron import", "create", 0, "cron", "import", "--file=-", "--dialect=unix", "--timezone=UTC", "--count=2", "--group", "${group}", "--run-as=")
	s.Input = "0 0 1 1 * echo fixture\n"
	s.JSON = false
	read("cron.export", "cron export", "cron", "export", "--task", "${task}")
	negative("cron.invalid-expression", "cron explain", "cron", "explain", "bad cron", "--timezone=UTC")
	create("notification.channel-create", "notification channel add", "channel", "id", "notification", "channel", "add", "QA webhook", "--endpoint", "${webhook}", "--authorization=${webhook_secret}", "--health-interval=0s", "--disabled")
	read("notification.channel-list", "notification channel list", "notification", "channel", "list")
	read("notification.channel-get", "notification channel get", "notification", "channel", "get", "${channel}")
	add("notification.channel-update", "notification channel update", "update", 0, "notification", "channel", "update", "${channel}", "--name=QA updated", "--endpoint", "${webhook}", "--health-interval=0s")
	add("notification.channel-rotate", "notification channel rotate", "update", 0, "notification", "channel", "rotate", "${channel}", "--authorization=${webhook_rotated}")
	add("notification.channel-enable", "notification channel enable", "update", 0, "notification", "channel", "enable", "${channel}")
	s = add("notification.channel-test", "notification channel test", "execution", 0, "notification", "channel", "test", "${channel}")
	s.Hook = "webhook-delivery"
	add("notification.group-set", "notification group set", "update", 0, "notification", "group", "set", "${group}", "--channel", "${channel}", "--on=success,failure", "--failure-threshold=2", "--duration-threshold=1s", "--quiet-period=1m", "--failure-to-start", "--recovery", "--reminder=1h")
	read("notification.group-show", "notification group show", "notification", "group", "show", "${group}", "--input")
	add("notification.task-set", "notification task set", "update", 0, "notification", "task", "set", "${base_task}", "--file", "${policy_file}")
	read("notification.task-show", "notification task show", "notification", "task", "show", "${base_task}", "--input")
	read("notification.task-effective", "notification task effective", "notification", "task", "effective", "${base_task}")
	read("notification.deliveries", "notification deliveries", "notification", "deliveries", "--channel", "${channel}", "--task", "${base_task}", "--limit=10")
	negative("notification.malformed-policy", "notification task set", "notification", "task", "set", "${base_task}", "--file", "${bad_json}")
	add("notification.channel-disable", "notification channel disable", "update", 0, "notification", "channel", "disable", "${channel}")
	add("notification.channel-remove", "notification channel rm", "delete", 0, "notification", "channel", "rm", "${channel}").JSON = false
	read("actor.list", "actor list", "actor", "list")
	create("actor.create", "actor create", "actor", "id", "actor", "create", "QA actor", "--kind=cli", "--capability=observe", "--expires=${future}")
	add("actor.update", "actor update", "update", 0, "actor", "update", "${actor}", "--name=QA actor updated", "--capability=operate", "--expires=${future}", "--state=active")
	add("actor.revoke", "actor revoke", "update", 0, "actor", "revoke", "${actor}")
	create("pairing.create", "pairing create", "pairing", "id", "pairing", "create", "QA enrollment", "--kind=cli", "--capability=manage", "--grant-duration=1h")
	s = add("profile.pair", "profile pair", "create", 0, "profile", "pair", "QA remote", "--address", "${remote}", "--expected-daemon-id", "${daemon_id}", "--pairing-id", "${pairing}", "--trusted-certificate", "${certificate}", "--client-name=QA enrollment", "--capability=manage")
	s.Spec.ExecutionKind = "command_harness"
	s.Input = "${phrase}\n"
	s.Save = []NamedValue{{Name: "profile", Value: "id"}}
	s.Hook = "profile-pair"
	read("pairing.list", "pairing list", "pairing", "list")
	create("pairing.cancel-create", "pairing create", "cancel_pairing", "id", "pairing", "create", "QA cancel", "--kind=cli", "--capability=observe")
	add("pairing.cancel", "pairing cancel", "update", 0, "pairing", "cancel", "${cancel_pairing}")
	read("credential.list", "credential list", "credential", "list")
	s = add("credential.rotate", "credential rotate", "update", 0, "credential", "rotate", "${credential}")
	s.Hook = "credential-rotate"
	add("credential.revoke", "credential revoke", "update", 0, "credential", "revoke", "${credential}")
	read("audit.list", "audit list", "audit", "list", "--operation=tasks.create", "--result=succeeded", "--limit=10", "--since=2026-01-01T00:00:00Z", "--until=2030-01-01T00:00:00Z")
	s = add("audit.export", "audit export", "read", 0, "audit", "export", "--output", "${audit_file}", "--limit=10", "--actor", "${actor}", "--operation=tasks.create", "--result=succeeded", "--since=2026-01-01T00:00:00Z", "--until=2030-01-01T00:00:00Z")
	s.JSON = false
	s.Hook = "audit-file"
	read("daemon.manifest", "daemon manifest", "daemon", "manifest")
	add("daemon.rename", "daemon rename", "update", 0, "daemon", "rename", "QA daemon")
	negative("daemon.reset-confirmation", "daemon reset-identity", "daemon", "reset-identity", "--confirm=wrong")
	steps[len(steps)-1].Spec.ExpectedExit = 1
	s = add("daemon.reset-identity", "daemon reset-identity", "update", 0, "daemon", "reset-identity", "--confirm", "${daemon_id}")
	s.Hook = "identity-reset"
	s = add("bundle.export", "bundle export", "read", 0, "bundle", "export", "--output", "${bundle_file}")
	s.JSON = false
	s.Hook = "bundle-file"
	read("bundle.validate", "bundle validate", "bundle", "validate", "${bundle_file}")
	read("bundle.compare", "bundle compare", "bundle", "compare", "${bundle_file}")
	s = add("bundle.preview", "bundle preview", "read", 0, "bundle", "preview", "${bundle_file}", "--output", "${plan_file}")
	s.JSON = false
	s.Hook = "plan-file"
	add("bundle.apply", "bundle apply", "update", 0, "bundle", "apply", "${plan_file}")
	negative("bundle.replayed-plan", "bundle apply", "bundle", "apply", "${plan_file}")
	steps[len(steps)-1].Spec.ExpectedExit = 1
	negative("bundle.bad-json", "bundle validate", "bundle", "validate", "${bad_json}")
	for _, action := range []string{"install", "uninstall", "start", "stop", "restart", "status"} {
		s = add("service."+action, "service "+action, "host-adapter", 0, "service", action)
		s.Spec.ExecutionKind = "command_harness"
		s.JSON = false
		s.Hook = "service-intent"
		if action == "install" {
			s.Spec.Argv = append(s.Spec.Argv, "--config", "${config}")
		}
	}
	s = add("gui.launch", "gui", "host-adapter", 0, "gui")
	s.Spec.ExecutionKind = "command_harness"
	s.Hook = "gui-intent"
	s.JSON = false
	for _, path := range []string{"preferences show", "preferences set", "preferences reset", "profile list", "profile show", "profile rename", "profile default set", "profile default show", "profile default reset", "profile remove"} {
		args := strings.Fields(path)
		switch path {
		case "preferences set":
			args = append(args, "--appearance=dark", "--popup-enabled", "--popup-severities=error", "--popup-conditions=failure", "--popup-daemons=fixture")
		case "preferences reset":
			args = append(args, "--scope=all")
		case "profile show", "profile default set":
			args = append(args, "--", "${profile}")
		case "profile rename":
			args = append(args, "--", "${profile}", "QA renamed profile")
		case "profile remove":
			args = append(args, "--confirm=${profile}", "--", "${profile}")
		}
		s = add(strings.ReplaceAll(path, " ", "."), path, "local-persistence", 0, args...)
		s.Spec.ExecutionKind = "command_harness"
		s.Hook = "local-state"
	}
	for _, action := range []string{"status", "enable", "rotate", "disable"} {
		args := []string{"mcp", "http", action}
		if action == "enable" {
			args = append(args, "--port", "${mcp_port}", "--permission=manage", "--name=QA MCP", "--require-confirmation", "--origin=http://127.0.0.1:43123")
		}
		s = add("mcp.http."+action, "mcp http "+action, "execution", 0, args...)
		s.Hook = "mcp-http"
	}
	s = add("mcp.stdio", "mcp serve", "execution", 0, "mcp", "serve", "--permission=manage", "--require-confirmation", "--name=QA stdio")
	s.Hook = "mcp-stdio"
	s.JSON = false
	negative("mcp.bad-port", "mcp http enable", "mcp", "http", "enable", "--port=-1")
	add("task.remove", "task rm", "delete", 0, "task", "rm", "${task}").JSON = false
	add("group.remove-child", "group rm", "delete", 0, "group", "rm", "${child_group}").JSON = false
	add("group.remove", "group rm", "delete", 0, "group", "rm", "${group}").JSON = false
	for i := range steps {
		steps[i].Spec.Coverage = coverageFor(steps[i].Spec.CommandPath, steps[i].Spec.Argv)
	}
	return steps
}

func coverageFor(path string, args []string) []string {
	result := []string{}
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if strings.HasPrefix(arg, "--") {
			name := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)[0]
			key := path + " --" + name
			if slices.Contains([]string{"json", "profile", "daemon-id", "credential-id", "certificate-file", "ipc-path"}, name) {
				key = " --" + name
			}
			if !slices.Contains(result, key) {
				result = append(result, key)
			}
		}
	}
	return result
}

func metadataSteps(commands []Command) []step {
	result := []step{}
	for _, command := range commands {
		if command.Hidden {
			continue
		}
		args := append(strings.Fields(command.Path), "--help")
		s := step{Spec: CaseSpec{ID: "help." + commandCaseID(command.Path), CommandPath: command.Path, ScenarioClass: "help", ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: args, ExpectedExit: 0, Checks: []CheckSpec{{Name: "help contains usage", Expected: "true"}}, Coverage: []string{command.Path + " --help"}}, OutputContains: "Usage:"}
		result = append(result, s)
		for _, alias := range command.Aliases {
			parts := strings.Fields(command.Path)
			if len(parts) == 0 {
				continue
			}
			parts[len(parts)-1] = alias
			result = append(result, step{Spec: CaseSpec{ID: "alias." + commandCaseID(command.Path) + "." + alias, CommandPath: command.Path, ScenarioClass: "alias", ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: append(parts, "--help"), ExpectedExit: 0, Checks: s.Spec.Checks, Coverage: []string{command.Path + " alias=" + alias}}, OutputContains: "Usage:"})
		}
	}
	for _, shell := range []string{"bash", "fish", "powershell", "zsh"} {
		result = append(result, step{Spec: CaseSpec{ID: "completion." + shell, CommandPath: "completion " + shell, ScenarioClass: "completion", ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: []string{"completion", shell, "--no-descriptions"}, ExpectedExit: 0, Checks: []CheckSpec{{Name: "completion script describes gosched", Expected: "true"}}, Coverage: []string{"completion " + shell + " --no-descriptions"}}, OutputContains: "gosched"})
	}
	result = append(result, step{Spec: CaseSpec{ID: "help.command", CommandPath: "help", ScenarioClass: "documentation", ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: []string{"help", "task"}, ExpectedExit: 0, Checks: []CheckSpec{{Name: "command help", Expected: "true"}}}, OutputContains: "Usage:"})
	return result
}
