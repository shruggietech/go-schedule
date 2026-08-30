# Data Model: Task-Completion Chains

## CompletionChain

| Field | Meaning | Validation |
| --- | --- | --- |
| `id` | Stable chain identity | Generated once, non-empty |
| `source_task_id` | Observed task | Existing task, different from target |
| `target_task_id` | Dispatched task | Existing task, different from source |
| `on_outcome` | `success`, `failure`, or `any` | Exact enum |
| `created_at` | UTC creation | RFC 3339 |
| `updated_at` | UTC mutation | RFC 3339, not before creation |

Unique identity is `(source_task_id, target_task_id, on_outcome)`. The directed graph must remain acyclic. Differently conditioned edges between the same tasks are allowed and may both match when one selector is `any`.

## CompletionDelivery

| Field | Meaning | Validation |
| --- | --- | --- |
| `id` | Stable delivery identity | Generated once |
| `chain_id` | Matched relationship | Existing while unresolved |
| `source_task_id` | Immutable source task | Copied from chain |
| `target_task_id` | Immutable target task | Copied from chain |
| `source_run_id` | Causing terminal run | Immutable run identity |
| `state` | `pending`, `claimed`, `completed`, or `resolved` | Exact transition |
| `attempts` | Claim count | Non-negative, increments on claim |
| `created_at` | Match time | UTC RFC 3339 |
| `claimed_at` | Latest claim | Optional UTC RFC 3339 |
| `completed_at` | Terminal time | Optional UTC RFC 3339 |
| `target_run_id` | Resulting run | Present for completed execution |
| `resolution` | Non-execution or recovery detail | Bounded text |

Unique identity is `(chain_id, source_run_id)`. State transitions:

```text
pending -> claimed -> completed
pending -> claimed -> resolved
claimed -> pending
```

Completed and resolved states are terminal.

## Run Extension

Existing `Run` gains optional `source_task_id` and `source_run_id`. Both are present together only when `trigger == completion`. Existing origins omit them. `RunTrigger` gains `completion`.

## Transaction Boundary

For executor success/failure, one transaction inserts the run, completes its incoming delivery, inserts unique matching outgoing deliveries, then commits. Failure rolls back every effect. Queued/skipped bookkeeping never creates outgoing deliveries; a skipped completion delivery resolves explicitly.

## Migration

Schema v9 creates tables and indexes named separately from removed legacy `triggers` and `dedup_ledger`, then adds nullable correlation columns to `runs`. Existing rows are not rewritten and decode with empty correlation.
