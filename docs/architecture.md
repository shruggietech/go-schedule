---
title: Architecture
nav_order: 8.5
---

# Architecture

The daemon owns scheduling, persistence, execution, and completion delivery. The CLI and desktop remain thin local clients.

## Completion delivery

When an executor returns a terminal success or failure, one SQLite transaction records the run and inserts a unique pending delivery for each matching chain. The engine claims deliveries in bounded batches, rechecks the target's current task and group eligibility, and sends eligible work through the existing worker and overlap path. The target run completes the incoming delivery and can create the next finite cascade in the same transaction.

Each delivery is unique for `(chain_id, source_run_id)`. Completed and resolved deliveries never replay. On daemon start, an interrupted `claimed` delivery returns to `pending` and is attempted again. This provides durable at-least-once delivery without a polling loop. A crash after the external command launches but before its completion transaction can repeat that command because SQLite cannot atomically observe an external process's side effects.

Completion chains form a directed acyclic graph. Validation rejects self-links and any update or insertion that can reach its own source. The dedicated v9 schema did not restore the obsolete v2 generic `triggers` or `dedup_ledger` tables.

## External trigger dispatch

Schema v12 introduces a distinct `external_triggers` table. A generated 256-bit key enters through the existing authenticated local IPC API, so the feature adds no TCP listener or second resident service. After key and target eligibility checks, the engine submits the request through the existing overlap-aware dispatcher and records `external_trigger` plus the stable trigger ID in run history.

Schema v13 adds `external_trigger_sets` plus optional set identity and permanent position columns on ordinary triggers. A Trigger Set is an administration boundary only: all members use the existing local IPC fire path, and the feature adds no listener, remote surface, payload, watcher, task Group relationship, or Chain target. Set creation and broad lifecycle mutations use database transactions, preserve automatic-source eligibility, and publish one redacted `trigger_set` event only after commit. Individual members retain ordinary trigger identity and lifecycle operations, but target changes are set-scoped to preserve the one-target invariant.

## Filesystem watcher dispatch

Schema v14 adds durable `filesystem_watchers` definitions and `runs.source_watcher_id`. The daemon owns one buffered fsnotify observer behind a portable adapter and one timer-driven event loop behind the project clock interface. Configuration mutation replaces the full observer generation, closes obsolete handles, and discards obsolete settling candidates before they can dispatch. Exact-file watches observe the parent directory so atomic replacement remains visible. Recursive directory watches register existing and newly created real directories without following symbolic links or junctions.

Create and write notifications are hints. A candidate first waits for its watcher debounce interval, then requires two equal regular-file size and modification-time snapshots separated by its stability interval. A stable candidate enters the existing eligibility and overlap-aware dispatcher with `filesystem_watcher` and the stable watcher ID. The matched path remains ephemeral and is excluded from runs, events, and logs. Missing roots and observer failure produce deduplicated degraded health transitions and a bounded two-second registration retry. Startup and recovery are prospective only and never scan or replay unobserved changes.

The user-scoped database retains the recoverable key because the desktop can explicitly reveal and copy it after restart. Ordinary list and detail responses, live events, run history, logs, and errors omit the key. Create, rotate, and explicit reveal are the only disclosure operations.
