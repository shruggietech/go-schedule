# Quickstart: Validate Filesystem Watchers

## Preconditions

- Build or run `goschedd`, `gosched`, and `gosched-gui` from the S056 branch.
- Create an enabled timeless task whose command appends one line to a temporary output file.
- Use a separate temporary watch directory and keep the output file outside that directory.

## Directory watcher lifecycle

1. Create an enabled recursive directory watcher with pattern `*.ready`, debounce `250ms`, stability `500ms`, and the task ID.
2. List and show it in human and JSON modes; verify the cleaned absolute root, target, timing, readiness, and active health.
3. Write one matching file in several quick chunks, wait for settlement, and verify exactly one run with trigger `filesystem_watcher` and the watcher ID.
4. Create a temporary file with a nonmatching name and rename it to a matching final name; verify exactly one run after stability.
5. Repeat in a nested real directory and verify recursion; repeat through a linked directory and verify no run.

## File watcher and recovery

1. Create a file-kind watcher for a path that does not exist and verify degraded health with a missing-root reason.
2. Create the parent directory if necessary, wait for bounded recovery, and verify active health without a run for the file's prior existence.
3. Write or atomically replace the exact file and verify one settled run.
4. Remove and recreate the observed parent, verify degraded then active health, and verify only post-recovery changes run.

## Live mutation and task eligibility

1. Change pattern, recursion, path, debounce, stability, and target without restarting the daemon; prove old-generation events cannot dispatch afterward.
2. Disable the watcher during a pending stability period and prove the candidate is cancelled.
3. Exercise task disabled, inactive, incomplete-command, blocked-group, queue-one, skip, and concurrent overlap behavior and confirm the existing dispatcher remains authoritative.
4. Delete the watcher and prove subsequent events do not run the task.

## Desktop workflow

1. Open Triggers and confirm External Triggers and Filesystem Watchers are separate sections.
2. Create, inspect, edit, disable, enable, and delete a watcher using keyboard and pointer input.
3. Confirm path and health details remain readable in light and dark modes and destructive deletion requires confirmation.

## Automated verification

```bash
go test ./internal/watcher ./internal/store ./internal/engine ./internal/api/server ./internal/api/client ./internal/cli ./gui/viewmodel
go test -race ./internal/watcher ./internal/store ./internal/engine ./internal/api/server ./internal/api/client ./internal/cli ./gui/viewmodel
go test ./gui
sh scripts/verify.sh all
```

Expected result: all focused tests and all eight canonical gates pass; core coverage remains at least 80 percent; no race, leak, replay, duplicate-dispatch, hard-wrap, or em-dash defect remains.
