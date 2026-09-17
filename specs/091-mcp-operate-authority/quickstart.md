# Quickstart: MCP Operate Authority

## Observe remains the default

```text
gosched mcp serve
```

This advertises five resources and no tools.

## Start an explicit stdio Operate session

```text
gosched mcp serve --permission operate --name "Codex local operator"
```

The runtime session ends and is revoked when the stdio process exits.

## Enable an explicit localhost HTTP Operate session

```text
gosched mcp http enable --port 43123 --name "Local desktop operator" --permission operate
```

The bearer credential is shown once. `gosched mcp http disable` revokes the session and closes the listener. Rotation invalidates the old bearer without changing the MCP actor.

## Call a tool

Supply the exact installation identifier reported by the daemon, the task identifier, and one UUID that remains stable across retries of the same logical request.

```json
{
  "daemon_id": "f6932ee9-a164-4df3-a0c2-bcbeb719b172",
  "task_id": "dba80df0-12a1-46a8-a357-b391211f578f",
  "request_id": "75c89425-d83a-4298-8d56-cfa90afac5c5"
}
```

Do not retry an uncertain run-now with a new request identifier. Inspect task activity first.
