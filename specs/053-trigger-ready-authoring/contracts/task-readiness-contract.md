# Contract: Task Readiness and Draft Authoring

## Create task

`POST /v1/tasks` accepts omitted `name`, `command`, `schedule`, and `at` fields. Existing validation remains for every nonempty supplied value.

- A request with no automatic source creates no schedule row and stores no schedule reference.
- A request missing command readiness or automatic activation readiness is stored with `enabled: false`, even if `enabled: true` was requested.
- A fully configured request that omits `enabled` preserves the historical API default of enabled.
- The desktop client continues to send its explicit create checkbox choice.

## Update task

`PATCH /v1/tasks/{id}` preserves omission as no change and adds explicit clear intent:

```json
{
  "clear_name": true,
  "clear_command": true,
  "clear_schedule": true
}
```

- A clear flag conflicts with a nonempty replacement for the same field and returns field-specific validation.
- Clearing command or the final automatic source atomically sets `enabled` to false.
- Failed validation leaves task and schedule references unchanged.

## Task detail

Configured task:

```json
{
  "task": {},
  "schedule": {},
  "readiness": {
    "command_ready": true,
    "activation_ready": true,
    "automatic_sources": ["schedule"],
    "status": "ready",
    "reason": ""
  }
}
```

Unscheduled task:

```json
{
  "task": {},
  "schedule": null,
  "readiness": {
    "command_ready": true,
    "activation_ready": false,
    "automatic_sources": [],
    "status": "manual_only",
    "reason": "No automatic activation source is configured."
  }
}
```

List responses keep the existing task array. Clients derive display status from task fields plus the already loaded completion-chain and group snapshots, using the same domain vocabulary as detail responses.

## Enable task

`POST /v1/tasks/{id}/enable` succeeds only when the task is command-ready, activation-ready, active, and backed by valid referenced state at the mutation boundary.

Readiness failures return `400 validation_failed` with field `enabled` and a specific reason. Missing task identity remains `404`.

## Run now

`POST /v1/tasks/{id}/run-now` requires command readiness and otherwise preserves explicit manual execution semantics regardless of local enabled or ancestor-group state.

Missing command returns `400 validation_failed` with field `command` before scheduler dispatch. Missing task identity remains `404`.

## CLI

- `gosched task add [name]` accepts zero or one name and no longer requires command or timing flags.
- Omitted values create a disabled draft.
- `gosched task edit <id> --name ""`, `--command ""`, or `--clear-schedule` expresses explicit clear intent.
- Human output prints `unnamed`, `schedule: none`, and the readiness status without manufacturing values.
- JSON output mirrors the API exactly.

## Compatibility

Existing configured-task create, edit, list, detail, enable, disable, Run now, schedule, startup, and completion-chain calls keep their behavior and JSON value shapes. The only new response shape occurs for a task that could not exist before S053.
