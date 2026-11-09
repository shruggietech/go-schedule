# S139 CLI surface register

The actual initialized inventory contains 152 nodes, 125 runnable nodes, two aliases and 1590 effective flag records. Every node has captured cases in both native environments. Counts below are original frozen observations, including help/negative cases, and are not a claim that every case uses a packaged executable. The retained semantic failures are resolved by the separately identified fresh checks in [verification](verification.md); they remain failures in the original capture. See [audit](audit.md) for boundaries and findings.

| Command | Runnable | Captured cases | Original semantic-linked failures | Actual boundary |
| --- | --- | --- | --- | --- |
| Root | no | 6 | 0 | Namespace/help |
| `__complete` | yes | 4 | 0 | Case-specific executable, handler or deterministic boundary |
| `actor` | no | 2 | 0 | Namespace/help |
| `actor create` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `actor list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `actor revoke` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `actor update` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `alerts` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `alerts ack` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `audit` | no | 2 | 0 | Namespace/help |
| `audit export` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `audit list` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `bundle` | no | 2 | 0 | Namespace/help |
| `bundle apply` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `bundle compare` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `bundle export` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `bundle preview` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `bundle validate` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `calendar` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `chain` | no | 2 | 0 | Namespace/help |
| `chain create` | yes | 14 | 4 | Case-specific executable, handler or deterministic boundary |
| `chain list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `chain rm` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `chain show` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `chain update` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `completion` | no | 2 | 0 | Namespace/help |
| `completion bash` | yes | 6 | 0 | Case-specific executable, handler or deterministic boundary |
| `completion fish` | yes | 4 | 0 | Case-specific executable, handler or deterministic boundary |
| `completion powershell` | yes | 4 | 0 | Case-specific executable, handler or deterministic boundary |
| `completion zsh` | yes | 4 | 0 | Case-specific executable, handler or deterministic boundary |
| `credential` | no | 2 | 0 | Namespace/help |
| `credential list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `credential revoke` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `credential rotate` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `cron` | no | 2 | 0 | Namespace/help |
| `cron convert` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `cron explain` | yes | 16 | 0 | Case-specific executable, handler or deterministic boundary |
| `cron export` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `cron import` | yes | 28 | 0 | Case-specific executable, handler or deterministic boundary |
| `daemon` | no | 2 | 0 | Namespace/help |
| `daemon manifest` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `daemon rename` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `daemon reset-identity` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `group` | no | 2 | 0 | Namespace/help |
| `group add` | yes | 16 | 0 | Case-specific executable, handler or deterministic boundary |
| `group disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `group edit` | yes | 16 | 4 | Case-specific executable, handler or deterministic boundary |
| `group enable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `group list` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `group rm` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `gui` | yes | 12 | 0 | Recorded launcher; actual missing-binary refusal separately observed |
| `health` | yes | 28 | 2 | Case-specific executable, handler or deterministic boundary |
| `help` | yes | 8 | 0 | Case-specific executable, handler or deterministic boundary |
| `logs` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `mcp` | no | 2 | 0 | Namespace/help |
| `mcp http` | no | 2 | 0 | Namespace/help |
| `mcp http disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `mcp http enable` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `mcp http rotate` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `mcp http status` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `mcp serve` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification` | no | 2 | 0 | Namespace/help |
| `notification channel` | no | 2 | 0 | Namespace/help |
| `notification channel add` | yes | 32 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel disable` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel enable` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel get` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel rm` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel rotate` | yes | 22 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel test` | yes | 22 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification channel update` | yes | 22 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification deliveries` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification group` | no | 2 | 0 | Namespace/help |
| `notification group set` | yes | 24 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification group show` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification task` | no | 2 | 0 | Namespace/help |
| `notification task effective` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification task set` | yes | 30 | 0 | Case-specific executable, handler or deterministic boundary |
| `notification task show` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `pairing` | no | 2 | 0 | Namespace/help |
| `pairing cancel` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `pairing create` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `pairing list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `preferences` | no | 2 | 0 | Scoped preference files |
| `preferences reset` | yes | 16 | 0 | Scoped preference files |
| `preferences set` | yes | 12 | 0 | Scoped preference files |
| `preferences show` | yes | 10 | 0 | Scoped preference files |
| `profile` | no | 2 | 0 | Namespace/help |
| `profile default` | no | 2 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile default reset` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile default set` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile default show` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile list` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile pair` | yes | 10 | 2 | Scoped profile store and memory credentials; native keyring blocked |
| `profile remove` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile rename` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `profile show` | yes | 10 | 0 | Scoped profile store and memory credentials; native keyring blocked |
| `runs` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `runs active` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `runs get` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `runtime` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `schedule` | no | 2 | 0 | Namespace/help |
| `schedule preview` | yes | 22 | 0 | Case-specific executable, handler or deterministic boundary |
| `search` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `service` | no | 2 | 0 | Namespace/help |
| `service install` | yes | 10 | 0 | Recorded control adapter; native read-only status separately observed |
| `service restart` | yes | 10 | 0 | Recorded control adapter; native read-only status separately observed |
| `service start` | yes | 12 | 0 | Recorded control adapter; native read-only status separately observed |
| `service status` | yes | 10 | 0 | Recorded control adapter; native read-only status separately observed |
| `service stop` | yes | 10 | 0 | Recorded control adapter; native read-only status separately observed |
| `service uninstall` | yes | 10 | 0 | Recorded control adapter; native read-only status separately observed |
| `systems` | no | 2 | 0 | Namespace/help |
| `systems search` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `systems summary` | yes | 16 | 0 | Case-specific executable, handler or deterministic boundary |
| `task` | no | 2 | 0 | Namespace/help |
| `task add` | yes | 26 | 0 | Case-specific executable, handler or deterministic boundary |
| `task disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `task edit` | yes | 60 | 4 | Case-specific executable, handler or deterministic boundary |
| `task enable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `task list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `task rm` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `task run-now` | yes | 34 | 0 | Case-specific executable, handler or deterministic boundary |
| `task show` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger` | no | 2 | 0 | Namespace/help |
| `trigger create` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger enable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger fire` | yes | 26 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger rm` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger rotate` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set` | no | 2 | 0 | Namespace/help |
| `trigger set create` | yes | 26 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set enable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set retarget` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set reveal` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set rm` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set rotate` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger set show` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger show` | yes | 18 | 0 | Case-specific executable, handler or deterministic boundary |
| `trigger update` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher` | no | 4 | 0 | Namespace/help |
| `watcher create` | yes | 14 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher disable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher enable` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher list` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher rm` | yes | 10 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher show` | yes | 16 | 0 | Case-specific executable, handler or deterministic boundary |
| `watcher update` | yes | 12 | 0 | Case-specific executable, handler or deterministic boundary |
