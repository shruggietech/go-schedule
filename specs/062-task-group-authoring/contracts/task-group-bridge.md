# Task and Group Bridge Contract

All methods are local desktop bindings. They return safe result objects and do not expose backend errors or transport objects.

## `Workspace()`

Returns the authoritative detailed task collection and group summaries. One call must be sufficient for the populated overview. A failure returns `unavailable` with safe recovery guidance.

## `Task(taskID)`

Returns complete editable detail for one task. Missing or concurrently deleted identity returns `rejected` without changing selection or drafts.

## `PreviewTask(draft)`

Parses the direct command line, validates environment input and timing mode, and returns exact program/argument plus schedule previews without mutation. Validation failures identify a known field.

## `SaveTask(draft)`

Creates an inactive task or differentially updates an existing task. An edit compares `original_updated_at` against current authority. A mismatch returns `stale` unless `overwrite_stale` is true. Accepted results include fresh detail and workspace data.

## `RunTask(taskID)`, `SetTaskEnabled(taskID, enabled)`, `DeleteTask(taskID)`

Each method accepts one stable identity, performs at most one daemon mutation, and returns refreshed workspace state on acceptance. The UI owns target-aware confirmation before invoking Run or Delete. Duplicate activation is suppressed while a matching operation is pending.

## `SaveGroup(draft)`

Creates or differentially updates name, parent, and declared enabled state. Parent choice excludes the group and descendants. Timestamp conflict behavior matches task edits.

## `SetGroupEnabled(groupID, enabled)`, `DeleteGroup(groupID)`

Each performs one authoritative mutation and refresh. Delete follows the existing descendant cascade and task-ungrouping contract. The UI owns target-aware cascade confirmation.

## Safe event integration

The existing `desktop:event` channel carries only kind, entity identity, generation, and safe message. `task.*` and `group.*` events request a debounced workspace refresh. They never carry task command, environment, input, run identity, path, or backend detail. Dirty editors retain their draft and compare the refreshed entity timestamp to mark staleness.

## Field identifiers

Safe validation fields are limited to `name`, `group_id`, `command`, `environment`, `working_dir`, `stdin`, `run_as`, `timezone`, `mode`, `schedule`, `schedule_syntax`, `at`, `overlap_policy`, `catchup_policy`, `missing_date_policy`, `time_basis`, `dst_gap_policy`, `dst_overlap_policy`, `enabled`, and `parent_id`. Unknown backend fields map to an operation-level message.
