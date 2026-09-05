# Data Model: Filesystem Watchers

## FilesystemWatcher

Durable configuration owned by the scheduler store.

| Field | Type | Rules |
| --- | --- | --- |
| `id` | string | Stable generated identifier, primary key |
| `name` | string | Trimmed, required, not globally unique |
| `kind` | enum | `file` or `directory` |
| `path` | string | Required, cleaned absolute host path |
| `pattern` | string | Directory only, defaults to `*`, valid base-name glob without separators |
| `recursive` | bool | Directory only; false for file watchers |
| `debounce` | duration | 25 milliseconds through 1 hour, default 250 milliseconds |
| `stability` | duration | 25 milliseconds through 1 hour, default 500 milliseconds |
| `target_task_id` | string | Required task foreign key with cascade delete |
| `enabled` | bool | Whether the runtime attempts observation and dispatch |
| `created_at` | timestamp | UTC, immutable |
| `updated_at` | timestamp | UTC, advanced on mutation |

### Validation

- File watchers store an empty pattern and false recursion.
- Directory watcher patterns are validated without requiring the directory to exist.
- Paths are normalized before equality and persistence.
- The target task must exist in the same transaction as each mutation.

### Lifecycle

```text
created enabled -> runtime registration attempted
created disabled -> runtime health disabled
updated -> old generation invalidated, new registration attempted
enabled -> runtime registration attempted
disabled -> pending candidates cancelled, health disabled
deleted -> pending candidates cancelled, runtime entry removed
target task deleted -> watcher cascades and runtime reloads through the task event path
```

## WatcherHealth

Ephemeral runtime state returned with API responses.

| Field | Type | Rules |
| --- | --- | --- |
| `state` | enum | `active`, `disabled`, or `degraded` |
| `reason` | string | Empty for active, stable explanation for other states, bounded to 512 characters |
| `changed_at` | timestamp | Injected-clock UTC time of the most recent state or reason transition |

No health field is persisted.

## PendingCandidate

Ephemeral state keyed by watcher ID and cleaned absolute file path.

| Field | Type | Rules |
| --- | --- | --- |
| `generation` | integer | Must match the current runtime generation before dispatch |
| `debounce_until` | timestamp | Moved forward by every matching candidate signal |
| `snapshot` | optional pair | Regular-file size and modification time captured after debounce |
| `stable_since` | timestamp | First time the current snapshot was observed |
| `due_at` | timestamp | Next evaluation deadline |

Candidates are deleted after dispatch, disappearance, degradation, disable, deletion, reload, or shutdown.

## Run provenance extension

| Field | Type | Rules |
| --- | --- | --- |
| `trigger` | enum extension | Adds `filesystem_watcher` |
| `source_watcher_id` | nullable string | Stable watcher ID for filesystem-originated runs; empty for all other origins |

The matched path is not stored in Run.

## Relationships

```text
Task 1 <- 0..N FilesystemWatcher
FilesystemWatcher 1 <- 0..N Run provenance references without a foreign key
FilesystemWatcher 1 <- 0..N ephemeral PendingCandidate
FilesystemWatcher 1 <- 1 ephemeral WatcherHealth while the daemon is running
```

## Migration v14

- Create `filesystem_watchers` and indexes on target task and enabled state.
- Add nullable `source_watcher_id` to `runs`.
- Preserve every version-13 row unchanged.
- Keep run provenance intentionally independent of watcher deletion so historical source identity survives lifecycle cleanup.
