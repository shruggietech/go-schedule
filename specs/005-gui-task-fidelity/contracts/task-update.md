# Contract delta: task update and task detail

**Feature**: [spec.md](../spec.md) · **Date**: 2026-07-22

Two contract changes. Both are additive at the wire level; one is a breaking
change to the Go type, internal to this repository.

## 1. `PATCH /v1/tasks/{id}` — tri-state group membership

### Change

`TaskUpdateRequest.GroupID` becomes `*string`.

```go
GroupID *string `json:"group_id,omitempty"`
```

### Semantics

| Wire form | Go value | Meaning |
|---|---|---|
| field omitted | `nil` | leave group membership unchanged |
| `"group_id": ""` | `ptr("")` | remove the task from all groups |
| `"group_id": "<id>"` | `ptr("<id>")` | assign the task to that group |

Before this change the second row was inexpressible: the handler tested
`if req.GroupID != ""`, so an empty value was indistinguishable from an omitted
one and no client could ungroup a task (spec §FR-014).

### Validation

A non-empty ID is resolved before assignment. Unknown IDs return:

```http
HTTP/1.1 400 Bad Request
{"code":"validation","field":"group_id","message":"group not found"}
```

not a foreign-key 500 (§FR-016, R6). An empty string skips the lookup — clearing
is always valid.

### Compatibility

- **Wire**: unchanged for every existing client. Omission still means unchanged;
  no existing caller sends `"group_id": ""` because it previously did nothing.
- **Go**: breaking for in-repo callers of the struct literal. Call sites are the
  CLI (`internal/cli/task.go`), the GUI (`gui/editor.go`), and server tests.
- **Store**: no schema change; `tasks.group_id` is already nullable with
  `ON DELETE SET NULL`.

### CLI mapping

`gosched task edit` decides which of the three intents to send using
`cmd.Flags().Changed("group")`:

| Invocation | Sent | Meaning |
|---|---|---|
| no `--group` | `nil` | unchanged — **existing behavior preserved** (§FR-015) |
| `--group ""` | `ptr("")` | ungroup |
| `--group <id>` | `ptr("<id>")` | assign |

## 2. `GET /v1/tasks/{id}` — schedule carries the operator's phrase

### Change

The `schedule` object in `TaskResponse` gains one field:

```json
{
  "task": { "...": "..." },
  "schedule": {
    "id": "...",
    "kind": "recurring",
    "rrule": "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
    "human_summary": "Every weekday at 09:00",
    "expression": "weekdays at 09:00"
  },
  "next_runs": ["..."]
}
```

### Semantics

- `expression` is the phrase the schedule was created from — re-submittable as
  the `schedule` field of a create or update.
- Omitted (empty) when the schedule is one-off, of kind `event`, or was stored
  before migration v4. Nothing reconstructs it in that last case — an empty
  `schedule` on update means "leave the schedule unchanged".
- **Never an execution input.** Clients round-trip it; the engine ignores it
  (§FR-011a).

### Compatibility

Purely additive. Existing clients ignore the new field. Applies to every endpoint
returning a task detail (`POST /v1/tasks`, `GET /v1/tasks/{id}`,
`PATCH /v1/tasks/{id}`), since all three render through the same `taskDetail`
helper.
