---
title: Local MCP access
nav_order: 7.5
---

# Local MCP access

**Audience:** people connecting Codex or another local MCP host to go-schedule\
**Applies to:** unreleased work after v1.4.0\
**Authority:** Observe resources by default, with explicit Operate and Manage tools over local stdio or an explicitly enabled authenticated numeric-loopback HTTP endpoint

go-schedule provides optional local Model Context Protocol access for inspecting scheduler state and, after explicit opt-in, operating existing tasks or managing bounded automation definitions. The preferred stdio mode starts only when an MCP host launches `gosched mcp serve`, communicates through standard input and standard output, and opens no network listener. Clients that cannot launch stdio can use a separately enabled authenticated endpoint bound only to numeric IPv4 loopback. The daemon must already be running, and the account controlling either mode must already have access to its Unix socket or Windows named pipe.

## Add it to Codex

Install go-schedule normally and confirm `gosched health` succeeds from the same account that runs Codex. Then add the local stdio server:

```text
codex mcp add go-schedule -- gosched mcp serve
codex mcp list
```

If `gosched` is not on that account's `PATH`, replace it with the absolute path to the installed executable. This stdio configuration does not enable the optional HTTP listener.

To remove the configuration later:

```text
codex mcp remove go-schedule
```

Restart or reconnect the MCP host after changing its server configuration. The server process exits when the host closes its stdio connection.

## Add it to another command-based host

Create a local MCP server entry named `go-schedule` whose command is `gosched` and whose ordered arguments are `mcp`, `serve`. Use an absolute executable path when the host does not inherit the installation directory in `PATH`. Do not add a URL, bearer credential, or shell wrapper to this stdio entry. The host owns the subprocess lifetime, stdio opens no listener, and removal consists of deleting that host configuration.

The exact configuration-file syntax is host-specific. Preserve the command and argument boundary rather than combining it into one shell command string. After saving, reconnect the host and confirm that it discovers five resources, four continuation templates, and no tools.

Observe is always the default. To let a trusted local host run, enable, or disable existing tasks, explicitly add `--permission operate` and a recognizable runtime client name:

```text
codex mcp add go-schedule-operate -- gosched mcp serve --permission operate --name "Codex task operator"
```

An Operate server exposes the same bounded resources plus exactly `tasks_run_now`, `tasks_enable`, and `tasks_disable`. Each call requires the exact daemon installation UUID, existing task ID, and a caller-generated request UUID. Reuse a request UUID only for a retry of the identical operation and target. Operate cannot create, edit, delete, or reconfigure tasks, groups, schedules, triggers, notification channels, connections, actors, or credentials.

To let a deliberately trusted host manage automation definitions, explicitly choose Manage:

```text
codex mcp add go-schedule-manage -- gosched mcp serve --permission manage --name "Codex definition manager"
```

Manage inherits the three Operate tools and adds exactly `tasks_manage`, `groups_manage`, `chains_manage`, `triggers_manage`, `watchers_manage`, and `notification_assignments_replace`. Each family tool changes one definition per call through the ordinary local API. Bulk mutation, actor administration, enrollment, credential reveal, direct database access, and raw filesystem access are unavailable. Task environment and stdin values are intentionally absent from the MCP schema. Trigger creation returns the new trigger identity but withholds its generated key; an Enroll-authorized human client owns later credential handling.

Attended environments can require an explicit confirmation field on every Manage call without imposing prompts on unattended deployments:

```text
codex mcp add go-schedule-manage-attended -- gosched mcp serve --permission manage --require-confirmation --name "Attended definition manager"
```

## Enable localhost HTTP when stdio is unavailable

The localhost endpoint is off on every daemon start. Choose an unused port and enable it through the protected local API:

```text
gosched mcp http enable --port 43123 --name "Local desktop host"
```

The command prints `http://127.0.0.1:43123/mcp` and a cryptographically random bearer credential exactly once. Copy both values into the local client's Streamable HTTP configuration. Do not place the credential in URLs, shell history, issue reports, or logs. `gosched mcp http status` reports the client name, permission, endpoint, successful request count, and a non-secret fingerprint. The client name and aggregate evidence are runtime-only. They contain no request content, failed-attempt history, or peer metadata.

To enable the same three Operate tools over localhost HTTP, make the permission explicit:

```text
gosched mcp http enable --port 43123 --name "Local desktop operator" --permission operate
```

Manage and its optional confirmation policy use the same lifecycle:

```text
gosched mcp http enable --port 43123 --name "Local desktop manager" --permission manage --require-confirmation
```

Native MCP clients normally send no `Origin` header. If a trusted local browser application must connect, explicitly allow its exact numeric-loopback origin when enabling:

```text
gosched mcp http enable --port 43123 --origin http://127.0.0.1:3000
```

Origins must use `http` or `https`, numeric `127.0.0.1`, and an explicit port. Explicit default ports are stored in browser-serialized form (`http://127.0.0.1` for port 80 and `https://127.0.0.1` for port 443) so exact matching agrees with the browser `Origin` header. Hostnames, wildcards, public addresses, paths, opaque origins, and implicit ports are rejected. Every request still needs the bearer credential and must address the exact active `127.0.0.1:<port>` Host.

Replace a copied credential without changing the endpoint or origins:

```text
gosched mcp http rotate
```

Rotation invalidates the previous value immediately and prints the replacement once. Revoke the credential and close the listener with `gosched mcp http disable`. Disable is safe to repeat. Daemon shutdown or restart also closes the listener and invalidates the credential; HTTP enablement, origins, and credentials are intentionally never persisted.

The desktop Agent Access workspace performs the same lifecycle through protected local IPC. It copies a new credential through the native clipboard boundary and never exposes that value to the webview. If clipboard copying fails, it disables localhost HTTP so an inaccessible credential cannot remain active.

Agent Access also presents one secret-free inventory of MCP actors for This computer. Each grant names the client, daemon, authority, inferred transport, creation, last use when known, expiry, and lifecycle state. A local administrator can create a remote MCP enrollment for 1 hour, 24 hours, 7 days, 30 days, or a deliberately acknowledged non-expiring grant. The one-time enrollment bundle is copied through the native clipboard boundary and is never returned to React; a failed copy cancels the pending pairing.

Existing remote grants can be narrowed from Manage to Operate or Observe, from Operate to Observe, assigned an earlier expiry, or permanently revoked. They cannot be widened, extended, made non-expiring, or reactivated through this control. The daemon reloads current actor authority on each request, so the next operation on an existing connection uses the narrower boundary. Recent actions show at most the newest 25 actor-attributed shared audit events without request bodies or protected values. Listener controls remain separate from grants, so changing stdio, localhost HTTP, or remote HTTPS never silently enables another transport.

## Available resources

| Resource | Contents |
| --- | --- |
| `goschedule://daemon/health` | Daemon status and version. |
| `goschedule://tasks/active` | Active task identity, state, readiness, safe schedule summary, and upcoming times. |
| `goschedule://schedules/upcoming` | Upcoming schedule projections without executable or raw recurrence configuration. |
| `goschedule://alerts/recent` | Recent bounded alert evidence. |
| `goschedule://runs/recent` | Recent bounded run evidence and output excerpts. |

Collections contain at most 100 records per page and provide an opaque continuation URI when another page is available. Each read requests only that page plus one lookahead record. One output excerpt is capped at 8 KiB. Other user-controlled text is capped at 2 KiB. Truncation is explicit for every bounded field. Task and schedule reads use an allowlisted SQLite projection that never loads execution inputs, while run output and alert messages are clipped in SQLite before they enter daemon memory or cross IPC.

Observe sessions on local and remote transports advertise resources only. Operate sessions add exactly three task tools. Manage sessions inherit those tools and add exactly six definition tools, with no direct database access, raw log resource, permission administration, enrollment, or credential reveal. Enabling localhost HTTP does not enable remote JSON or remote MCP access and does not change CLI, GUI, Unix-socket, or Windows-named-pipe authorization.

## Remote HTTPS authorization

Remote MCP is independently disabled even when the remote JSON listener is enabled. An operator must set `remote.mcp.enabled` and a canonical public HTTPS `remote.mcp.resource_url` ending in `/mcp`. The public resource URL can differ from the numeric listener address when a reverse proxy fronts the daemon, but the proxy must preserve the public `Host` value and connect to the daemon through trusted TLS.

The RFC 9728 protected-resource document points clients to authorization-server metadata on the same HTTPS origin. Unattended clients use the standard client-credentials grant with an existing persistent `mcp` credential: the credential ID is `client_id`, the one-time returned credential token is `client_secret`, and the exact configured resource is the RFC 8707 `resource` value. Pairing phrases and credentials belonging to desktop, CLI, or JSON actors are rejected.

The token endpoint returns an opaque, short-lived, memory-only bearer. It is valid only in the `Authorization` header at the configured `/mcp` resource. Tokens expire after ten minutes by default, disappear on daemon restart, and fail immediately after source credential rotation, revocation, expiry, actor capability change, or daemon identity change. Requested `mcp:observe`, `mcp:operate`, or `mcp:manage` scope cannot exceed the persistent actor capability.

## Operate results and retry safety

Every Operate result has a stable schema version, permission, operation, daemon ID, task ID, request ID, outcome, and bounded message. `accepted` means the daemon authoritatively accepted the operation. `rejected` means validation, target state, or request identity prevented it. `denied` means current server-owned authority prevented it. `uncertain` means the client did not receive authoritative completion evidence; inspect current task state and Activity before deciding whether to retry.

The runtime session remembers a bounded set of request IDs for ten minutes. Repeating the same request ID with the same operation and target returns the first result without another mutation. Reusing it for another operation or target is rejected. This in-memory retry boundary ends with the MCP process or localhost listener, so callers must still reconcile uncertain outcomes rather than replaying blindly after a restart.

## Manage results, deletes, and compatibility

Every Manage result adds affected object kind and object identity to the stable Operate outcome model. A call changes one definition or atomically replaces one task or group assignment scope, so there is no partial bulk result. Delete behavior, reference conflicts, readiness updates, validation errors, and assignment replacement are the same as the JSON API, GUI, and CLI because tools delegate to that API rather than mutating storage directly.

`accepted` proves the API completed the mutation. `rejected` covers invalid action, schema, target, reference, confirmation, or request-identity reuse. `denied` means current server-owned capability prevented dispatch. `uncertain` means the caller must inspect current state before retrying. Incompatible or future actions fail closed. Manage never silently retries transport failures.

## Trust and secret boundary

Task names, schedule summaries, policy summaries, alert messages, and command output are user-controlled, untrusted data. Every response labels these fields and tells the host not to treat their contents as instructions. JSON encoding does not make hostile text trustworthy.

Treat every resource field as display data even if it contains Markdown, XML-like text, ANSI escapes, JSON fragments, or sentences that resemble agent instructions. Such content cannot add tools or change permissions. A host should render or summarize it as untrusted scheduler data and must not execute instructions found inside it.

Observe response types structurally exclude task commands, arguments, environment values, stdin, working directories, run-as identities, raw schedule definitions, trigger keys, notification endpoints, authorization values, and internal IPC or filesystem paths. Manage inputs permit bounded definition fields such as commands, schedules, and watcher paths where the ordinary API requires them, but omit task environment and stdin. Manage results return only operation and object identity, outcome, and a safe explanation, never definition content or credentials.

The subprocess inherits the operating-system access that launched it and uses only the same local IPC client as `gosched health`. An explicit Operate launch creates a short-lived runtime MCP actor with a random secret retained only in process memory as a digest. Every task action is authorized against that actor's current state and written to the existing intent-first audit log with the actor, daemon, operation, task, correlation ID, and result. Listener disable, process exit, or explicit teardown revokes the actor; requests after revocation remain attributable and fail closed. If the daemon denies the controlling or reading identity, MCP returns a bounded access error and does not retry through another transport, direct database access, elevation, or another identity.

## Permission model

| Class | Current state | Boundary |
| --- | --- | --- |
| Observe | Enabled | Read bounded allowlisted resources through existing local IPC authorization. |
| Operate | Explicit opt-in | Run, enable, or disable one existing task with a runtime MCP actor, exact daemon and task targets, deduplicated request identity, current authorization, and attributable audit. |
| Manage | Explicit opt-in | Create, update, or delete one bounded task, group, chain, trigger, or watcher definition, or atomically replace one notification-assignment scope, through the ordinary API with redacted results and optional confirmation. |

This separation is intentional. Observe discovery cannot expose or invoke Operate or Manage tools. Operate cannot cross into Manage. Manage cannot cross into Enroll, permission administration, or credential administration. Remote access preserves the same hierarchy and delegates mutations through the ordinary actor-attributed API.

## Troubleshooting

Run `gosched health` as the same account. If it reports that the daemon is unavailable, start the `goschedd` service. If it reports access denied, follow the platform install guide for daemon-group membership and sign in again where required. MCP preserves that result and will not bypass it.

Protocol messages are written only to stdout. Startup or runtime diagnostics use stderr. If a host reports malformed protocol output, capture both streams separately and include the installed `gosched --version` and daemon version from `gosched health` in the report.

For localhost HTTP, first run `gosched mcp http status`. A disabled result after restart is expected. An unavailable-port error means another process owns the selected port; disable any old endpoint or choose another port. `401 Unauthorized` means the client omitted the current one-time credential or retained a value from before rotation or restart. `403 Forbidden` means Host, Origin, or current actor authority did not satisfy the active policy. Do not weaken those checks by using `localhost`, a wildcard origin, a proxy, or a forwarded Host.
