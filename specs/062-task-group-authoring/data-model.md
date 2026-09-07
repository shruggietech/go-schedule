# Data Model: Task and Group Authoring

## Workspace Snapshot

| Field | Type | Rules |
| --- | --- | --- |
| tasks | Task Summary array | Authoritative detailed list, stable task identity |
| groups | Group Summary array | Flat transport form with derived path and depth |
| loaded_at | timestamp | UTC observation time, never used as entity authority |

## Task Summary

| Field | Type | Rules |
| --- | --- | --- |
| id | string | Stable daemon identity |
| name | string | Empty legacy value displays as `unnamed` |
| group_id | string | Empty means ungrouped |
| group_path | string | Full hierarchy path or `Not assigned` |
| command_configured | boolean | Does not expose command text in the list |
| declared_enabled | boolean | Stored task flag |
| effective_state | enum | `runnable`, `manual_only`, `task_disabled`, `group_disabled`, `not_runnable`, `terminal`, `invalid_group` |
| effective_reason | string | Safe user-facing reason, names a blocking group path when applicable |
| lifecycle | enum | Existing task lifecycle value |
| timezone | string | Existing zone identity |
| schedule_summary | string | Recurrence/one-off summary or `Manual only` |
| policy_summary | string | Empty when not meaningful |
| next_runs | timestamp array | At most five UTC instants |
| updated_at | timestamp | Stale-draft comparison token |

## Task Detail and Draft

Task Detail contains every editable task value: identity, name, group, formatted direct command line, working directory, ordered environment rows, standard input, run identity, enabled state, timezone, timing mode, schedule expression/syntax or one-off instant, overlap, catch-up, missing-date, time-basis, daylight-saving gap, daylight-saving overlap, schedule and policy summaries, next runs, readiness, and `updated_at`.

Task Draft adds `is_new`, `original_updated_at`, `overwrite_stale`, and client-only dirty/pending state. New drafts use defaults: inactive, Local timezone, recurring mode with blank schedule, queue-one overlap, one-run catch-up, skip missing dates, wall-clock basis, next-valid gap, first overlap, and no group.

### Task Draft transitions

```text
clean -> dirty -> validating -> saving -> clean
                  |             |
                  v             v
                invalid       failed
dirty + newer authority -> stale -> reload -> clean
                              |-> explicit overwrite -> saving
```

## Group Summary and Draft

| Field | Type | Rules |
| --- | --- | --- |
| id | string | Stable daemon identity |
| name | string | Required, trimmed for creation and rename |
| parent_id | string | Empty means root |
| path | string | Full hierarchy path |
| depth | integer | Zero-based, bounded traversal |
| declared_enabled | boolean | Stored group flag |
| effective_enabled | boolean | False when self or an ancestor is disabled or invalid |
| effective_reason | string | Safe inherited or invalid-hierarchy reason |
| child_count | integer | Direct child groups |
| task_count | integer | Direct assigned tasks |
| descendant_count | integer | All descendant groups |
| updated_at | timestamp | Stale-draft comparison token |

Group Draft contains optional identity, name, parent identity, declared enabled state, original timestamp, and overwrite intent. Parent choices exclude self and descendants. Draft group creation can start disabled.

## Command Suggestion

| Platform | Display command | Program | Ordered arguments | Recognizable output |
| --- | --- | --- | --- | --- |
| windows | `cmd.exe /d /c ver` | `cmd.exe` | `/d`, `/c`, `ver` | Contains `Windows` |
| macos | `/usr/bin/sw_vers` | `/usr/bin/sw_vers` | none | Contains `ProductName` and `ProductVersion` |
| linux | `uname -a` | `uname` | `-a` | Nonempty kernel/system line |

Unsupported platforms expose no insertable suggestion. The mapping contains no free-form executable input and is exact-allowlisted.

## Operation Result

| Field | Type | Rules |
| --- | --- | --- |
| action | enum | `load`, `preview`, `save_task`, `run_task`, `toggle_task`, `delete_task`, `save_group`, `toggle_group`, `delete_group` |
| outcome | enum | `accepted`, `rejected`, `stale`, `unavailable` |
| message | string | Safe user-facing result |
| field | string | Optional known form field for validation focus |
| entity_id | string | Optional affected identity |
| workspace | Workspace Snapshot | Present after accepted collection mutation |
| task | Task Detail | Present after load, preview, or accepted task save as appropriate |

No backend error, endpoint, command, environment, standard input, run identity, or path is permitted in ordinary result messages.
