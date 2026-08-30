---
title: Architecture
nav_order: 8.5
---

# Architecture

The daemon owns scheduling, persistence, execution, and completion delivery.
The CLI and desktop remain thin local clients.

## Completion delivery

When an executor returns a terminal success or failure, one SQLite transaction
records the run and inserts a unique pending delivery for each matching chain.
The engine claims deliveries in bounded batches, rechecks the target's current
task and group eligibility, and sends eligible work through the existing worker
and overlap path. The target run completes the incoming delivery and can create
the next finite cascade in the same transaction.

Each delivery is unique for `(chain_id, source_run_id)`. Completed and resolved
deliveries never replay. On daemon start, an interrupted `claimed` delivery
returns to `pending` and is attempted again. This provides durable at-least-once
delivery without a polling loop. A crash after the external command launches
but before its completion transaction can repeat that command because SQLite
cannot atomically observe an external process's side effects.

Completion chains form a directed acyclic graph. Validation rejects self-links
and any update or insertion that can reach its own source. The dedicated v9
schema intentionally does not restore the obsolete generic `triggers` or
`dedup_ledger` tables.
