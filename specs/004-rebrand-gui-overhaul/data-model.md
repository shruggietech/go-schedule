# Data Model: Rebrand + GUI & Installer Overhaul

Covers new/changed entities, events, and the store migration. Times are UTC (RFC 3339), per the
existing convention.

## New entity: LogRecord

A unified log entry surfaced in the Logs view and persisted to the on-disk log file. Produced by
the daemon's `slog` handler; not stored in SQLite (see research §3).

| Field        | Type                | Notes                                                        |
|--------------|---------------------|-------------------------------------------------------------|
| `id`         | string              | Stable id (monotonic seq or ULID) for dedupe in the GUI     |
| `time`       | time.Time (UTC)     | When the record was emitted                                 |
| `severity`   | enum                | `info` \| `warning` \| `error` (mapped from slog level)     |
| `source`     | string              | Component / slog group (e.g. `engine`, `executor`, `api`)   |
| `message`    | string              | Short human message                                         |
| `task_id`    | string (optional)   | Correlation id when the record relates to a task/run        |
| `run_id`     | string (optional)   | Correlation id when related to a specific run               |
| `attrs`      | map[string]any      | Structured slog attributes = the "cause/context" detail     |

**Severity mapping**: slog `Debug`,`Info` → `info`; `Warn` → `warning`; `Error` → `error`.

**Validation/handling**: records are append-only and immutable; the GUI never edits them. The ring
holds the most recent N (default 1000); the file holds the rotation window (default 5 × 10 MiB).

**JSONL line shape** (on disk): one JSON object per line with the fields above (`attrs` inlined as
top-level keys is acceptable as long as the reserved keys above are not shadowed). See
[contracts/log-file.md](contracts/log-file.md).

## Changed entity: Alert (unchanged shape, unified presentation)

`domain.Alert` keeps its existing fields (`id`, `task_id`, `severity`, `kind`, `message`,
`created_at`, `acknowledged`) and table. In the GUI it is presented in the same Logs list as
`LogRecord`, tagged `source = "alert"`, severity reused directly. "Dismiss All" acknowledges the
shown alerts and clears the in-memory log view; it does not delete the alerts table rows or the
log file.

## Events (broker) — additions

`events.Event` gains kinds and payloads (broker stays non-blocking / drop-on-full):

| Kind          | Payload                          | Published when                                  |
|---------------|----------------------------------|-------------------------------------------------|
| `run`         | `*domain.Run` (existing)         | run state changes (existing)                    |
| `alert`       | `*domain.Alert` (existing)       | alert raised (existing)                         |
| `log` (NEW)   | `*domain.LogRecord`              | any daemon log record emitted                   |
| `task` (NEW)  | `{verb, *domain.Task | id}`      | task created / updated / deleted / enabled etc. |
| `group` (NEW) | `{verb, *domain.Group | id}`     | group created / updated / deleted / enabled etc.|

`verb` ∈ `created` | `updated` | `deleted`. For `deleted`, only the id need be populated.

## Removed entities (Triggers)

Deleted from `domain` and the store:

- `Trigger` (and `TriggerOutcome` enum: `OnSuccess`/`OnFailure`/`OnAny`)
- `DedupLedger`

De-emphasized (kept for history compatibility, no longer produced/dispatched):

- `ScheduleKind` value `event` (`ScheduleEvent`) and `Schedule.TriggerID` — columns remain in the
  `schedules` table but are unused; new schedules are never of kind `event`.
- `RunTrigger` value `event` (`TriggerEvent`) — retained as a historical enum value in
  `runs.trigger`; no new `event` runs are produced.

## Store migration v3

Appended to `internal/store/store.go` `migrations`:

```sql
-- v3: remove the Triggers feature
DROP TABLE IF EXISTS dedup_ledger;   -- child first (FK to triggers)
DROP TABLE IF EXISTS triggers;
```

**Behavior**:
- On a DB that contains triggers → tables dropped; daemon starts clean (FR-020).
- On a DB that never had triggers → `DROP ... IF EXISTS` is a no-op (FR-020).
- The migration framework records version 3; no data loss outside trigger tables.
- Optional: before dropping, count rows in `triggers`; if > 0, the daemon logs a one-time
  `warning` LogRecord noting how many trigger definitions were removed (troubleshooting aid).

**Idempotency / safety**: runs inside the existing transactional migration loop; failure rolls
back and aborts startup with a wrapped error naming migration 3.

## Configuration additions (`config.Config`)

| Field                 | Type   | Default            | Notes                                          |
|-----------------------|--------|--------------------|------------------------------------------------|
| `LogFilePath`         | string | `<DataDir>/logs/goschedule.log` | Empty ⇒ derive from DataDir       |
| `LogMaxSizeBytes`     | int    | `10 << 20` (10 MiB)| Rotation threshold per file                    |
| `LogMaxFiles`         | int    | `5`                | Rotated files retained (bounds disk, FR-017)   |
| `LogRingSize`         | int    | `1000`             | Recent records served by `GET /v1/logs`        |

Validation: positive integers (fail-fast, naming the field), consistent with existing `Validate()`
style. `DataDir`/`DBPath` rename handling (`goschedule` dir, `goschedule.db`) per research §1.

## Entity relationships (after this feature)

```
Group 1──* Task *──1 Schedule        (Schedule.kind ∈ {one_off, recurring};   event retired)
Task  1──* Run                       (Run.trigger ∈ {schedule, catchup, manual}; event retired)
Task  0..1──* Alert                  (existing; shown in Logs view)
(LogRecord)                          (not persisted in SQLite; file + ring + stream)
— Trigger, DedupLedger: REMOVED —
```
