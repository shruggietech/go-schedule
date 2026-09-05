# Data Model: Activity Diagnostics Clarity

## Recent activity response

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `logs` | array of log records | yes | Bounded recent records, newest first. |
| `log_path` | string | yes | Exact configured daemon log path; empty means unavailable. |

The path is metadata about the complete rotating log and is independent of whether `logs` contains records.

## GUI state

| Field | Type | Meaning |
| --- | --- | --- |
| `Logs` | log-record slice | Current bounded records. |
| `LogPath` | string | Last exact path learned from a successful refresh. |

A successful refresh replaces both values under one lock. Live log events may change `Logs` but retain `LogPath`.

## Activity diagnostics text

- Always states that Activity is a limited recent view and that older daemon records are in the full log.
- Non-empty path: appends the exact value after `Full daemon log:`.
- Empty path: appends `Full daemon log: unavailable until daemon responds.`

## Startup completion record

| Property | Value |
| --- | --- |
| Message | `daemon startup complete` |
| `endpoint` | Resolved local IPC endpoint. |
| `db` | Resolved database path. |
| `log_path` | Resolved configured rotating-log path. |

The record is emitted once at the existing ready-to-serve boundary.
