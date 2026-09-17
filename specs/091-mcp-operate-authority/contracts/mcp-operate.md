# Contract: MCP Operate Tools

## Discovery

Observe sessions advertise zero tools. Operate sessions advertise exactly:

- `tasks_run_now`
- `tasks_enable`
- `tasks_disable`

No tool accepts free-form commands, configuration bodies, secret references, or implicit targets.

## Input

```json
{
  "daemon_id": "installation UUID",
  "task_id": "task UUID",
  "request_id": "caller-generated UUID"
}
```

All properties are required. Unknown properties are rejected by the SDK schema. Identifiers are bounded and validated before dispatch.

## Output

```json
{
  "schema_version": "1",
  "permission": "operate",
  "operation": "tasks.run_now",
  "daemon_id": "installation UUID",
  "task_id": "task UUID",
  "request_id": "caller-generated UUID",
  "outcome": "accepted",
  "message": "Run request accepted."
}
```

`accepted` confirms daemon acceptance, not task success. `rejected` means no mutation occurred because the target or task state was invalid. `denied` means the current actor lacked authority. `uncertain` means the transport could not prove whether the daemon accepted the request, and callers must inspect current state or activity instead of retrying with a new request identity.

## Authority matrix

| Action | Observe | Operate | Manage | Enroll |
| --- | --- | --- | --- | --- |
| Discover task operation tools | No | Yes | Yes | Yes |
| Run existing task | No | Yes | Yes | Yes |
| Enable or disable existing task | No | Yes | Yes | Yes |
| Create, edit, or delete task | No | No | Yes | Yes |
| Change actors or sessions | No | No | No | Yes |
