---
title: Local MCP access
nav_order: 7.5
---

# Local MCP access

**Audience:** people connecting Codex or another local MCP host to go-schedule\
**Applies to:** unreleased v1.3.0 work\
**Authority:** observe-only resources over local stdio or an explicitly enabled authenticated numeric-loopback HTTP endpoint

go-schedule provides optional local Model Context Protocol access for inspecting scheduler state. The preferred stdio mode starts only when an MCP host launches `gosched mcp serve`, communicates through standard input and standard output, and opens no network listener. Clients that cannot launch stdio can use a separately enabled authenticated endpoint bound only to numeric IPv4 loopback. The daemon must already be running, and the account controlling either mode must already have access to its Unix socket or Windows named pipe.

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

## Enable localhost HTTP when stdio is unavailable

The localhost endpoint is off on every daemon start. Choose an unused port and enable it through the protected local API:

```text
gosched mcp http enable --port 43123 --name "Local desktop host"
```

The command prints `http://127.0.0.1:43123/mcp` and a cryptographically random bearer credential exactly once. Copy both values into the local client's Streamable HTTP configuration. Do not place the credential in URLs, shell history, issue reports, or logs. `gosched mcp http status` reports the client name, endpoint, successful request count, and a non-secret fingerprint. The client name and aggregate evidence are runtime-only. They contain no request content, failed-attempt history, or peer metadata.

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

## Available resources

| Resource | Contents |
| --- | --- |
| `goschedule://daemon/health` | Daemon status and version. |
| `goschedule://tasks/active` | Active task identity, state, readiness, safe schedule summary, and upcoming times. |
| `goschedule://schedules/upcoming` | Upcoming schedule projections without executable or raw recurrence configuration. |
| `goschedule://alerts/recent` | Recent bounded alert evidence. |
| `goschedule://runs/recent` | Recent bounded run evidence and output excerpts. |

Collections contain at most 100 records per page and provide an opaque continuation URI when another page is available. Each read requests only that page plus one lookahead record. One output excerpt is capped at 8 KiB. Other user-controlled text is capped at 2 KiB. Truncation is explicit for every bounded field. Task and schedule reads use an allowlisted SQLite projection that never loads execution inputs, while run output and alert messages are clipped in SQLite before they enter daemon memory or cross IPC.

Both local transports advertise resources only. They have no tools, prompts, scheduler mutations, remote authentication, direct database access, or raw log resource. Enabling localhost HTTP does not enable remote JSON or future remote MCP access and does not change CLI, GUI, stdio, Unix-socket, or Windows-named-pipe authorization.

## Trust and secret boundary

Task names, schedule summaries, policy summaries, alert messages, and command output are user-controlled, untrusted data. Every response labels these fields and tells the host not to treat their contents as instructions. JSON encoding does not make hostile text trustworthy.

Treat every resource field as display data even if it contains Markdown, XML-like text, ANSI escapes, JSON fragments, or sentences that resemble agent instructions. Such content cannot add tools or change permissions. A host should render or summarize it as untrusted scheduler data and must not execute instructions found inside it.

The MCP response types structurally exclude task commands, arguments, environment values, stdin, working directories, run-as identities, raw schedule definitions, trigger keys, notification endpoints, authorization values, and internal IPC or filesystem paths. The adapter never serializes a daemon task, run, alert, or schedule object directly.

The subprocess inherits the identity that launched it and uses only the same local IPC client as `gosched health`. The HTTP transport is controlled through that protected IPC boundary and uses the same Observe adapter. If the daemon denies the controlling or reading identity, MCP returns a bounded access error and does not retry through another transport, direct database access, elevation, or another identity.

## Permission model

| Class | Current state | Boundary |
| --- | --- | --- |
| Observe | Enabled | Read bounded allowlisted resources through existing local IPC authorization. |
| Operate | Disabled | Future task control requires authenticated identity, per-action authorization, attributable audit records, and explicit arguments. |
| Manage | Disabled | Future configuration changes require all Operate gates plus deliberate confirmation for destructive or externally visible actions. |

This separation is intentional. A future release must not turn an Observe resource into an action or activate Operate or Manage without satisfying their independent authorization, audit, confirmation, and test requirements.

## Troubleshooting

Run `gosched health` as the same account. If it reports that the daemon is unavailable, start the `goschedd` service. If it reports access denied, follow the platform install guide for daemon-group membership and sign in again where required. MCP preserves that result and will not bypass it.

Protocol messages are written only to stdout. Startup or runtime diagnostics use stderr. If a host reports malformed protocol output, capture both streams separately and include the installed `gosched --version` and daemon version from `gosched health` in the report.

For localhost HTTP, first run `gosched mcp http status`. A disabled result after restart is expected. An unavailable-port error means another process owns the selected port; disable any old endpoint or choose another port. `401 Unauthorized` means the client omitted the current one-time credential or retained a value from before rotation or restart. `403 Forbidden` means Host or Origin did not exactly match the active policy. Do not weaken those checks by using `localhost`, a wildcard origin, a proxy, or a forwarded Host.
