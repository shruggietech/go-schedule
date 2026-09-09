# Contract: Profiles, Target Selection, and Desktop Lifecycle

## CLI Global Selection

```text
gosched [--profile <id-or-label>] [--endpoint <https-origin> --daemon-id <id> --credential-id <id> --certificate-file <path>] <command>
```

No remote selection means This computer. Named and explicit forms conflict. Explicit selection requires every field. Bearer values and pairing phrases are never global flags.

For a remote human-mode invocation, stderr begins with one line:

```text
target: <profile-label> (<https-origin>, daemon <short-id>)
```

JSON stdout retains the command's existing schema and contains no banner.

## CLI Profile Lifecycle

```text
gosched profile pair <label> --endpoint <https-origin> --daemon-id <id> --pairing-id <id> --certificate-file <path> --capability <value> --phrase-stdin
gosched profile list
gosched profile show <id-or-label>
gosched profile rename <id-or-label> <new-label>
gosched profile rm <id-or-label>
```

`profile pair` reads exactly one phrase line from standard input, exchanges it once, saves the bearer in native credential storage, and writes secret-free metadata. `profile rm` deletes the native credential before metadata. Every command supports `--json`.

## Desktop Bridge

```text
ConnectionProfiles() -> { action, outcome, message, workspace? }
SelectConnection(profileID) -> { action, outcome, message, workspace? }
RenameConnection(profileID, label) -> { action, outcome, message, workspace? }
RemoveConnection(profileID) -> { action, outcome, message, workspace? }
PairRemote(draft with optional repairProfileID) -> { action, outcome, message, credentialID?, profileID? }
```

The workspace contains one synthetic local entry plus safe remote profiles. Results never contain certificate PEM or bearer values. Selection publishes the requested target immediately, cancels the prior generation, and negotiates before feature operations become available.

## Remote Transport

- Map local client path `/v1/...` to selected origin path `/api/v1/...`.
- Require TLS 1.3 and the selected PEM roots.
- Reject all redirects.
- Load one bearer from native storage before construction.
- Call public manifest discovery and require exact daemon installation ID before exposing an authenticated feature client.
- Never retry a mutation automatically.

## Capability Gating

| Surface | Required advertised capability | Minimum authority |
| --- | --- | --- |
| Tasks read | `tasks` | Observe |
| Task run, enable, disable | `tasks` | Operate |
| Task create, update, delete | `tasks` | Manage |
| Groups read | `groups` | Observe |
| Schedule | `schedule` | Observe |
| Activity | `activity` | Observe |
| Actors and audit | `actor-authorization`, `management-audit` | Manage |
| Automation, Notifications, Agent Access | Remote route support not advertised in S078 | Disabled for remote targets |

Desktop-local Settings remain available regardless of selected daemon.
