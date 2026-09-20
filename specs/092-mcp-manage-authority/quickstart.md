# Quickstart: MCP Manage Authority

Start a local stdio server with explicit Manage permission:

```text
gosched mcp serve --permission manage --name "Trusted automation client"
```

For attended deployments, require confirmation on every Manage call:

```text
gosched mcp serve --permission manage --require-confirmation --name "Attended automation client"
```

The localhost HTTP equivalent adds `--permission manage` and optional `--require-confirmation` to `gosched mcp http enable`. Observe remains the default.

Callers first obtain the exact daemon identity through Observe resources. A Manage call then supplies that `daemon_id`, a fresh `request_id`, one action, and the family-specific definition. Reuse a request ID only to retry the exact same call. An uncertain result requires state inspection before any new request.
