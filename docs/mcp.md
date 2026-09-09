---
title: Local MCP access
nav_order: 7.5
---

# Local MCP access

**Audience:** people connecting Codex or another local MCP host to go-schedule\
**Applies to:** unreleased v1.3.0 work\
**Authority:** observe-only resources over the installed `gosched` CLI and the existing protected daemon IPC endpoint

go-schedule provides an optional local Model Context Protocol server for inspecting scheduler state. It starts only when an MCP host launches `gosched mcp serve`, communicates through standard input and standard output, and opens no network listener. The daemon must already be running, and the launching account must already have access to its Unix socket or Windows named pipe.

## Add it to Codex

Install go-schedule normally and confirm `gosched health` succeeds from the same account that runs Codex. Then add the local stdio server:

```text
codex mcp add go-schedule -- gosched mcp serve
codex mcp list
```

If `gosched` is not on that account's `PATH`, replace it with the absolute path to the installed executable. Do not add `--url`: this server has no HTTP or TCP transport.

To remove the configuration later:

```text
codex mcp remove go-schedule
```

Restart or reconnect the MCP host after changing its server configuration. The server process exits when the host closes its stdio connection.

## Available resources

| Resource | Contents |
| --- | --- |
| `goschedule://daemon/health` | Daemon status and version. |
| `goschedule://tasks/active` | Active task identity, state, readiness, safe schedule summary, and upcoming times. |
| `goschedule://schedules/upcoming` | Upcoming schedule projections without executable or raw recurrence configuration. |
| `goschedule://alerts/recent` | Recent bounded alert evidence. |
| `goschedule://runs/recent` | Recent bounded run evidence and output excerpts. |

Collections contain at most 100 records per page and provide an opaque continuation URI when another page is available. Each read requests only that page plus one lookahead record. One output excerpt is capped at 8 KiB. Other user-controlled text is capped at 2 KiB. Truncation is explicit for every bounded field. Task and schedule reads use an allowlisted SQLite projection that never loads execution inputs, while run output and alert messages are clipped in SQLite before they enter daemon memory or cross IPC.

The initial server advertises resources only. It has no tools, prompts, scheduler mutations, TCP listener, remote authentication, direct database access, or raw log resource.

## Trust and secret boundary

Task names, schedule summaries, policy summaries, alert messages, and command output are user-controlled, untrusted data. Every response labels these fields and tells the host not to treat their contents as instructions. JSON encoding does not make hostile text trustworthy.

The MCP response types structurally exclude task commands, arguments, environment values, stdin, working directories, run-as identities, raw schedule definitions, trigger keys, notification endpoints, authorization values, and internal IPC or filesystem paths. The adapter never serializes a daemon task, run, alert, or schedule object directly.

The subprocess inherits the identity that launched it and uses only the same local IPC client as `gosched health`. If the daemon denies that identity, MCP returns a bounded access error and does not retry through another transport, direct database access, elevation, or another identity.

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
