# Phase 1 Data Model: Maintainer Test Scripts

**Feature**: 006-maintainer-test-scripts | **Date**: 2026-07-23

Two independent SQLite databases, each created on demand. They share no schema and are never
joined; keeping them separate means deleting one to reset a heartbeat experiment does not
discard accumulated host history.

**Conventions used by both**

- Timestamps are stored as **integer milliseconds since the Unix epoch**, in UTC, alongside a
  human-readable RFC 3339 string carrying the local offset. The integer is what queries sort
  and subtract; the string is what a maintainer reads. Storing both is redundant on purpose —
  the alternative is every ad-hoc query needing a conversion function.
- `NULL` means *could not be determined*. It is never used for a legitimate zero.
- Both databases run in WAL mode with a 5000 ms busy timeout (research §R5).
- Schema version is recorded in `meta` so a future change can migrate rather than guess.

---

## `heartbeat.db`

### `meta`

| Column | Type | Notes |
|---|---|---|
| `key` | TEXT PRIMARY KEY | |
| `value` | TEXT | |

Seeded with `schema_version = 1` and `created_utc`.

### `beat`

One row per completed heartbeat run. Written once, at the end of the run (FR-021c).

| Column | Type | Null? | Notes |
|---|---|---|---|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | | |
| `session_id` | TEXT NOT NULL | | Random per process invocation. Groups the beats of one continuous run; a change of `session_id` is a restart boundary. |
| `sequence` | INTEGER NOT NULL | | 1-based position **within the session**. In default single-shot mode this is always 1. |
| `label` | TEXT | yes | Caller-supplied tag, so several scheduled tasks can share one database and stay distinguishable. |
| `hostname` | TEXT NOT NULL | | |
| `username` | TEXT | yes | |
| `pid` | INTEGER NOT NULL | | |
| `started_ms` | INTEGER NOT NULL | | Run start, captured in memory before any work. |
| `started_iso` | TEXT NOT NULL | | RFC 3339 with local offset. |
| `finished_ms` | INTEGER NOT NULL | | Run end, immediately before the write. |
| `duration_ms` | INTEGER NOT NULL | | `finished_ms - started_ms`. Stored rather than computed so the overlap query stays a plain join. |
| `expected_ms` | INTEGER | yes | Expected firing moment. `NULL` when no source was available. |
| `expected_source` | TEXT NOT NULL | | One of `env`, `boundary`, `none`. Never absent — a drift figure without its provenance is the thing FR-003 forbids. |
| `drift_ms` | INTEGER | yes | `started_ms - expected_ms`. `NULL` exactly when `expected_ms` is. Positive means late. |
| `interval_seconds` | INTEGER | yes | The interval the caller declared; the basis of `boundary` snapping and of the reliability check. |
| `exit_code` | INTEGER NOT NULL | | What the process returned to the scheduler. `0` on the normal path. |
| `outcome` | TEXT NOT NULL | | `ok` or `failed`. Redundant with `exit_code` by design — FR-017 requires the two agree, and storing both makes disagreement detectable rather than assumed away. |
| `sched_env` | TEXT | yes | JSON object of any `GOSCHED_*` variables found. Empty today (research §R1); captured so a future release's context is recorded without a schema change. |

**Indexes**: `started_ms`; `(session_id, sequence)`.

**Validation rules**

- `finished_ms >= started_ms`.
- `drift_ms IS NULL` if and only if `expected_ms IS NULL`, which holds exactly when
  `expected_source = 'none'`.
- `expected_source = 'boundary'` requires `interval_seconds` to be present.
- `outcome = 'ok'` if and only if `exit_code = 0`.
- `sequence` is unique within `session_id`.

**Derived, not stored** — computed by reader queries:

- *Inter-beat interval*: difference of consecutive `started_ms` within a session.
- *Gap*: an inter-beat interval exceeding twice `interval_seconds`.
- *Overlap*: two beats whose `[started_ms, finished_ms]` ranges intersect. Decidable only
  because both endpoints are stored — the reason FR-002 requires the finish moment.
- *Unreliable drift*: `expected_source = 'boundary'` and
  `abs(drift_ms) > interval_seconds * 250`. Reported as unreliable rather than as fact,
  because past that magnitude a late firing and an early next one are indistinguishable.

---

## `system.db`

### `meta`

Same shape as above.

### `snapshot`

One row per invocation of the host-inspection script.

| Column | Type | Null? | Notes |
|---|---|---|---|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | | |
| `unixtime_ms` | INTEGER NOT NULL | | |
| `iso_local` | TEXT NOT NULL | | RFC 3339 with offset. |
| `iso_utc` | TEXT NOT NULL | | |
| `tz_offset_minutes` | INTEGER NOT NULL | | Stored explicitly so DST transitions are visible in the data rather than hidden inside a string. |
| `hostname` | TEXT NOT NULL | | |
| `username` | TEXT | yes | |
| `process_count` | INTEGER | yes | `NULL` = probe unavailable. `0` would mean genuinely none. |
| `uptime_seconds` | INTEGER | yes | |
| `os_platform` | TEXT NOT NULL | | `windows`, `linux`, `darwin`. |
| `os_release` | TEXT | yes | |
| `script_pid` | INTEGER NOT NULL | | |
| `script_flavor` | TEXT NOT NULL | | `powershell` or `posix` — which twin wrote the row. Makes twin-parity differences findable in the data. |
| `invocation_source` | TEXT NOT NULL | | Caller-supplied; defaults to `manual`. |

**Index**: `unixtime_ms`.

### `snapshot_address`

| Column | Type | Null? | Notes |
|---|---|---|---|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | | |
| `snapshot_id` | INTEGER NOT NULL | | → `snapshot(id)` ON DELETE CASCADE |
| `family` | TEXT NOT NULL | | `ipv4` or `ipv6`. |
| `address` | TEXT NOT NULL | | |
| `interface` | TEXT | yes | |
| `scope` | TEXT | yes | `loopback`, `link`, `global` where determinable. |

**Index**: `snapshot_id`.

### `snapshot_port`

| Column | Type | Null? | Notes |
|---|---|---|---|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | | |
| `snapshot_id` | INTEGER NOT NULL | | → `snapshot(id)` ON DELETE CASCADE |
| `protocol` | TEXT NOT NULL | | `tcp` or `udp`. |
| `family` | TEXT | yes | |
| `address` | TEXT | yes | Bound address; `NULL` where the probe does not report it. |
| `port` | INTEGER NOT NULL | | |
| `process_name` | TEXT | yes | Frequently `NULL` — most platforms require elevation to attribute a socket to a process. Absence here is normal, not a defect. |

**Index**: `(snapshot_id, protocol, port)`.

**Validation rules**

- A snapshot with zero address rows is valid: the probe may have been unavailable. The
  distinction between *no addresses found* and *could not look* lives in the warning on
  stderr and is not reconstructible from the table alone — an accepted limitation, recorded
  here so nobody later mistakes an empty child table for a determined result.
- `port` is between 1 and 65535.
- Child rows are only ever written for a snapshot that was successfully inserted, so a
  partial failure never orphans them.

---

## Growth

No pruning, rotation, or size cap (FR-022a). A beat row is roughly 200 bytes; a one-minute
schedule produces about 525,000 rows and on the order of 100 MB per year. A snapshot with its
children is roughly 2 KB, so an hourly schedule is about 18 MB per year. Both are disposable:
deleting the file is the documented reset, and it is the only one.
