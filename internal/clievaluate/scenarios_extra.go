package clievaluate

import "strings"

func additionalSteps() []step {
	steps := []step{}
	add := func(id, path, class, kind string, exit int, hook string, args ...string) *step {
		s := step{Spec: CaseSpec{ID: id, CommandPath: path, ScenarioClass: class, ExecutionKind: kind, Platforms: []string{"windows", "linux"}, Argv: args, ExpectedExit: exit, Checks: []CheckSpec{{Name: "scenario invariant", Expected: "true"}, {Name: "independent evidence contract", Expected: "true"}}, Coverage: coverageFor(path, args)}, Hook: hook}
		steps = append(steps, s)
		return &steps[len(steps)-1]
	}
	add("negative.unknown-command", "", "unknown-command", "command_harness", 2, "unchanged", "qa-command-does-not-exist")
	add("help.remote-selection-independent", "help", "documentation", "packaged_cli", 0, "read", "--profile=qa-missing-profile", "help", "task", "show")
	add("help.unknown-topic", "help", "validation", "command_harness", 2, "unchanged", "help", "task", "qa-missing-topic")
	add("completion.reject-operands", "completion bash", "arity", "command_harness", 2, "unchanged", "completion", "bash", "qa-unexpected")
	add("bundle.occupied-output", "bundle export", "validation", "packaged_cli", 1, "unchanged", "--ipc-path", "${ipc}", "bundle", "export", "--output", "${runtime}")
	add("bundle.readonly-output", "bundle export", "validation", "packaged_cli", 1, "readonly-output", "--ipc-path", "${ipc}", "bundle", "export", "--output", "${readonly_file}")
	add("daemon.invalid-config", "@isolation", "configuration", "command_harness", 1, "invalid-daemon-config")
	add("daemon.occupied-config", "@isolation", "configuration", "command_harness", 1, "occupied-daemon-config")
	for _, entry := range []struct {
		path string
		args []string
	}{{"task show", []string{"task", "show"}}, {"group edit", []string{"group", "edit"}}, {"chain show", []string{"chain", "show"}}, {"watcher show", []string{"watcher", "show"}}, {"trigger show", []string{"trigger", "show"}}, {"runs get", []string{"runs", "get"}}} {
		args := append([]string{"--ipc-path", "${ipc}"}, entry.args...)
		args = append(args, "missing-fixture-id")
		if entry.path == "group edit" {
			args = append(args, "--name=must not persist")
		}
		add("negative.missing."+strings.ReplaceAll(entry.path, " ", "."), entry.path, "validation", "packaged_cli", 1, "unchanged", args...)
	}
	add("root.version", "", "version", "packaged_cli", 0, "version", "--version")
	add("completion.callback", "__complete", "completion", "packaged_cli", 0, "completion", "__complete", "--help").Spec.Coverage = []string{"__complete --help"}
	s := add("completion.no-desc", "__complete", "completion", "packaged_cli", 0, "completion", "__completeNoDesc", "")
	s.Spec.Coverage = append(s.Spec.Coverage, "__complete alias=__completeNoDesc")
	add("selector.conflicts", "health", "validation", "command_harness", 2, "unchanged", "health", "--profile=fixture", "--endpoint=https://127.0.0.1:1", "--daemon-id=fixture", "--credential-id=fixture", "--certificate-file", "${certificate}", "--ipc-path", "${ipc}")
	add("selector.profile", "health", "remote-target", "command_harness", 0, "remote-authority", "health", "--profile", "${remote_profile}", "--json")
	add("selector.explicit-tls", "health", "remote-target", "command_harness", 0, "remote-authority", "health", "--endpoint", "${remote}", "--daemon-id", "${daemon_id}", "--credential-id", "${remote_credential}", "--certificate-file", "${certificate}", "--json")
	add("selector.daemon-endpoint", "health", "remote-target", "command_harness", 0, "remote-authority", "health", "--daemon-endpoint", "${remote}", "--daemon-id", "${daemon_id}", "--credential-id", "${remote_credential}", "--certificate-file", "${certificate}", "--json").Spec.Coverage = []string{" --daemon-endpoint"}
	add("selector.wrong-certificate", "health", "validation", "command_harness", 2, "wrong-certificate", "health", "--endpoint", "${remote}", "--daemon-id", "${daemon_id}", "--credential-id", "${remote_credential}", "--certificate-file", "${bad_certificate}")
	add("selector.wrong-identity", "health", "validation", "command_harness", 1, "unchanged", "health", "--endpoint", "${remote}", "--daemon-id=wrong-daemon", "--credential-id", "${remote_credential}", "--certificate-file", "${certificate}")
	add("selector.no-fallback", "health", "validation", "packaged_cli", 1, "no-fallback", "health", "--ipc-path", "${missing_ipc}")
	add("selector.empty-ipc", "health", "validation", "command_harness", 2, "unchanged", "health", "--ipc-path=")
	add("selector.empty-remote", "health", "validation", "command_harness", 2, "unchanged", "health", "--daemon-endpoint=")
	add("selector.canceled-secret", "systems summary", "cancellation", "command_harness", 1, "canceled-secret", "systems", "summary", "--all")
	add("selector.partial-systems", "systems summary", "partial-result", "command_harness", 1, "partial-systems", "systems", "summary", "--all", "--json")
	add("transport.authority-audit", "@transport", "authority", "command_harness", 0, "authority-audit")
	add("transport.expired-pairing", "@transport", "authority", "deterministic_engine", 0, "expired-pairing")
	add("transport.event-cancel-reconnect", "@transport", "stream", "command_harness", 0, "event-stream")
	add("transport.mcp-permission-confirmation", "@transport", "authority", "packaged_cli", 0, "mcp-authority")
	add("isolation.private-secret-entry", "@isolation", "credential-boundary", "packaged_cli", 1, "private-secret", "--internal-native-credential-read")
	add("isolation.descendant-drain", "@isolation", "owned-cleanup", "command_harness", 0, "descendant-drain")
	s = add("isolation.unix-socket-bound", "@isolation", "path-boundary", "command_harness", 0, "socket-bound")
	s.Spec.Platforms = []string{"linux"}
	s.Spec.ApplicabilityReason = "unix_socket_only"
	for _, name := range []string{"dst-gap-fold", "month-boundaries", "catchup-startup", "one-off-completion", "overlap-policies", "watcher-recovery", "notification-retry", "notification-transport", "smtp-transport", "popup-filtering", "bundle-reference-integrity"} {
		add("engine."+name, "@engine", "deterministic-policy", "deterministic_engine", 0, name)
	}
	add("task.stdin-file-create", "task add", "create", "packaged_cli", 0, "input-create", "--json", "--ipc-path", "${ipc}", "task", "add", "QA file input", "--command", "${worker}", "--arg=fixture-child", "--stdin-file", "${stdin_file}", "--disabled")
	s = add("task.run-as-unix", "task add", "create", "packaged_cli", 0, "input-create", "--json", "--ipc-path", "${ipc}", "task", "add", "QA owned user", "--command", "${worker}", "--run-as=${current_user}", "--disabled")
	s.Spec.Platforms = []string{"linux"}
	s.Spec.ApplicabilityReason = "unix_run_as_supported"
	s = add("task.run-as-windows", "task add", "validation", "packaged_cli", 2, "unchanged", "--json", "--ipc-path", "${ipc}", "task", "add", "QA unsupported run as", "--command", "${worker}", "--run-as=${current_user}", "--disabled")
	s.Spec.Platforms = []string{"windows"}
	s.Spec.ApplicabilityReason = "windows_run_as_refusal"
	add("task.edit-stdin-at", "task edit", "update", "packaged_cli", 0, "input-edit", "--json", "--ipc-path", "${ipc}", "task", "edit", "${base_task}", "--stdin=inline 世界", "--at", "${future}")
	add("actor.clear-expiration", "actor update", "update", "packaged_cli", 0, "read", "--json", "--ipc-path", "${ipc}", "actor", "update", "${actor}", "--clear-expiration")
	add("audit.actor-filter", "audit list", "validation", "packaged_cli", 2, "unchanged", "--json", "--ipc-path", "${ipc}", "audit", "list", "--actor", "${actor}", "--limit=0")
	add("notification.task-shorthand", "notification task set", "update", "packaged_cli", 0, "task-policy", "--json", "--ipc-path", "${ipc}", "notification", "task", "set", "${base_task}", "--channel", "${second_channel}", "--on=failure", "--failure-threshold=1", "--duration-threshold=1s", "--quiet-period=0s", "--failure-to-start", "--recovery", "--reminder=1h")
	s = add("notification.group-stdin", "notification group set", "update", "packaged_cli", 0, "group-policy", "--json", "--ipc-path", "${ipc}", "notification", "group", "set", "${base_group}", "--file=-")
	s.Input = "${policy_json}"
	s = add("notification.task-resume-inheritance", "notification task set", "update", "packaged_cli", 0, "task-inherit", "--json", "--ipc-path", "${ipc}", "notification", "task", "set", "${base_task}", "--file=-")
	s.Input = `{"assignments":[]}`
	add("notification.group-effective", "notification task effective", "read", "packaged_cli", 0, "inherited-policy", "--json", "--ipc-path", "${ipc}", "notification", "task", "effective", "${base_task}")
	add("notification.delivery-filters", "notification deliveries", "read", "packaged_cli", 0, "read", "--json", "--ipc-path", "${ipc}", "notification", "deliveries", "--run", "${run}", "--state=succeeded", "--limit=0")
	add("bundle.watcher-bindings", "bundle preview", "validation", "packaged_cli", 2, "unchanged", "--json", "--ipc-path", "${ipc}", "bundle", "preview", "${bundle_file}", "--watcher-path=missing=${watch_root}")
	for _, scope := range []string{"appearance", "popups"} {
		add("preferences.reset-"+scope, "preferences reset", "local-persistence", "command_harness", 0, "local-state", "--json", "preferences", "reset", "--scope="+scope)
	}
	add("preferences.invalid-store-repair", "preferences reset", "local-persistence", "command_harness", 0, "repair-preferences", "--json", "preferences", "reset", "--scope=all")
	add("service.recorded-failure", "service start", "host-adapter", "command_harness", 1, "service-failure", "service", "start")
	add("gui.recorded-failure", "gui", "host-adapter", "command_harness", 1, "gui-failure", "gui")
	add("cron.import-file", "cron import", "create", "packaged_cli", 0, "cron-file", "--ipc-path", "${ipc}", "cron", "import", "--file", "${cron_file}", "--timezone=UTC", "--dialect=unix", "--count=1")
	add("cron.system-dry-run", "cron import", "read", "packaged_cli", 0, "cron-system", "--ipc-path", "${ipc}", "cron", "import", "--file=-", "--dry-run", "--system", "--dialect=unix", "--timezone=UTC").Input = "0 0 1 1 * fixture-user echo fixture\n"
	add("cron.quartz-dry-run", "cron import", "read", "packaged_cli", 0, "cron-quartz", "--ipc-path", "${ipc}", "cron", "import", "--file=-", "--dry-run", "--dialect=quartz", "--timezone=UTC").Input = "0 0 9 ? * MON-FRI echo fixture\n"
	add("cron.partial-invalid", "cron import", "validation", "packaged_cli", 2, "cron-partial", "--ipc-path", "${ipc}", "cron", "import", "--file=-", "--dry-run", "--timezone=UTC").Input = "0 0 1 1 * echo fixture\ninvalid line\n"
	for i := range steps {
		steps[i].Spec.Coverage = append(steps[i].Spec.Coverage, coverageFor(steps[i].Spec.CommandPath, steps[i].Spec.Argv)...)
		if strings.HasPrefix(steps[i].Spec.ID, "selector.") {
			for j, arg := range steps[i].Spec.Coverage {
				if strings.HasSuffix(arg, " --endpoint") {
					steps[i].Spec.Coverage[j] = " --endpoint"
				}
			}
		}
	}
	return steps
}
