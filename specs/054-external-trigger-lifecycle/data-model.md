# Data Model: External Trigger Lifecycle

## Trigger

| Field | Type | Rules |
|---|---|---|
| `id` | string | Stable generated identifier, primary key |
| `name` | string | Required after trimming, duplicates allowed |
| `key` | string | Unique opaque secret, never serialized by the ordinary domain JSON shape |
| `target_task_id` | string | Required task reference with cascading deletion |
| `enabled` | boolean | Defaults to true |
| `created_at` | timestamp | Set at creation |
| `updated_at` | timestamp | Updated on every mutation and rotation |

## Trigger readiness

Readiness is derived and is not persisted. A trigger can report `ready`, `disabled`, `target_missing`, `command_incomplete`, `task_inactive`, `task_disabled`, or `group_blocked`. The API includes a human-readable reason alongside the stable state.

## Run provenance

Runs gain nullable `source_trigger_id` provenance and the new `external_trigger` run source. This identifier is intentionally not a foreign key so historical runs remain attributable after the trigger is deleted. Raw keys never enter run records.

## Automatic-source invariants

An enabled schedule, an enabled incoming completion chain, or an enabled external trigger is an automatic source. A task without any automatic source cannot remain active. Removing, disabling, or retargeting the last enabled source and deactivating the affected task occur in the same transaction.

## Schema migration

Migration v12 creates the `triggers` table and its task and key indexes, then adds `source_trigger_id` to `runs` with an empty-string default for existing records. Existing databases retain all task, group, chain, and run data.
