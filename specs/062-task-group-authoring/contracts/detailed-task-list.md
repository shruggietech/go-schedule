# Detailed Task List Contract

## Existing request

`GET /v1/tasks` retains its existing response envelope and task array. Existing `group` and `state` filters remain unchanged.

## Opt-in request

`GET /v1/tasks?details=true` applies the same filters and returns:

```json
{
  "tasks": [
    {
      "task": {},
      "schedule": {},
      "readiness": {},
      "policy_summary": "",
      "next_runs": []
    }
  ]
}
```

Each array member is exactly the existing single-task `TaskResponse` representation. `schedule` is null and `next_runs` is empty for manual-only tasks. The list is empty, never null, when no tasks match.

The client exposes `ListTaskDetails(ctx, group, state)` separately from `ListTasks`, so callers must opt into the richer shape. One bad referenced schedule fails the detailed request safely rather than returning a partially authoritative collection.

## Complete clear semantics

`PATCH /v1/tasks/{id}` adds `clear_working_dir` and `clear_run_as` booleans. Each conflicts with a nonempty replacement for the same field and otherwise clears that value. Existing requests are unchanged. Environment and arguments use non-null empty collections to clear, and standard input retains its pointer-based clear convention.

`POST /v1/groups` adds an optional `enabled` boolean using the task-create convention. Omitted values preserve the historical enabled default; explicit false creates an inactive draft group atomically.
