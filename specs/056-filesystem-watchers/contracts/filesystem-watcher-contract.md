# Filesystem Watcher Contract

## Local API

All routes use the existing protected local IPC transport and standard JSON error envelope.

| Method | Route | Outcome |
| --- | --- | --- |
| `GET` | `/v1/watchers` | List durable definitions joined to current runtime health |
| `POST` | `/v1/watchers` | Create one watcher and reload runtime observation |
| `GET` | `/v1/watchers/{id}` | Return one watcher and current health |
| `PATCH` | `/v1/watchers/{id}` | Replace supplied mutable fields atomically and reload runtime observation |
| `DELETE` | `/v1/watchers/{id}` | Delete one watcher, cancel pending candidates, and return no content |
| `POST` | `/v1/watchers/{id}/enable` | Enable one watcher and attempt registration |
| `POST` | `/v1/watchers/{id}/disable` | Disable one watcher and cancel pending candidates |

Create accepts `name`, `kind`, `path`, `pattern`, `recursive`, `debounce`, `stability`, `target_task_id`, and optional `enabled`. Update accepts pointers to the same mutable fields. Durations use Go duration strings at the API boundary and canonical nanoseconds in domain and persistence code.

Responses include all durable fields plus `target_task_name`, `readiness`, `readiness_reason`, and `health`. Health contains `state`, `reason`, and `changed_at`.

Validation errors identify `name`, `kind`, `path`, `pattern`, `debounce`, `stability`, `target_task_id`, or the complete body. Missing watcher IDs use `not_found`; configuration conflicts use `conflict`; runtime degradation is represented in a successful response rather than turning a valid persisted definition into an API failure.

## CLI

```text
gosched watcher create --name <name> --kind <file|directory> --path <path> --task <task-id> [--pattern <glob>] [--recursive] [--debounce 250ms] [--stability 500ms] [--disabled] [--json]
gosched watcher list [--json]
gosched watcher show <id> [--json]
gosched watcher update <id> [--name <name>] [--kind <file|directory>] [--path <path>] [--task <task-id>] [--pattern <glob>] [--recursive=<bool>] [--debounce <duration>] [--stability <duration>] [--json]
gosched watcher enable <id> [--json]
gosched watcher disable <id> [--json]
gosched watcher delete <id> [--yes] [--json]
```

Human list output uses stable columns for ID, name, kind, path, pattern or exact-file marker, recursive state, target, enabled state, readiness, and health. JSON output emits the API response shape. Mutations return the current object except delete, which reports the deleted ID. Errors go to stderr and produce a nonzero exit status.

## Desktop

The Triggers destination contains two vertically arranged sections: External Triggers and Filesystem Watchers. The watcher table uses Name, Kind, Path, Pattern, Recursive, Target, Enabled, Readiness, Health, and ID columns with the existing persisted adjustable-width behavior. Its toolbar provides New, Edit, Enable or disable, and Delete actions. Forms expose kind-sensitive pattern and recursion controls, duration entries, target selection, and enabled intent. Health reason is available through full-cell detail and degraded rows remain selectable.

## Event and log contract

Watcher configuration and health transitions use a distinct `filesystem_watcher` event kind with stable watcher ID, name, verb, health state, and bounded reason. Paths are omitted. Create, update, enable, disable, and delete each publish one lifecycle event. Runtime health publishes only when state or reason changes. Structured logs use watcher ID, name, state, and reason without path.

## Filesystem semantics

- File kind observes its parent and filters exact cleaned path.
- Directory kind observes the root and, when recursive, every real subdirectory.
- The pattern matches a candidate's base name only.
- Create and Write are candidate signals. A final Create resulting from rename into place is sufficient.
- Remove, Rename of the old path, Chmod, directories, and non-regular files never dispatch directly.
- Two equal size and modification-time snapshots separated by stability duration are required after debounce.
- Observer error invalidates pending work and initiates full bounded recovery.
- Startup and recovery do not scan or replay existing files.
