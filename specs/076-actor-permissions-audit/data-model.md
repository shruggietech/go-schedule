# Data Model: Actor Permissions and Management Audit

## Actor

| Field | Rule |
|---|---|
| `id` | Opaque UUID, immutable |
| `kind` | `local_os`, `desktop`, `cli`, `json`, or `mcp` |
| `display_name` | Trimmed 1 to 80 Unicode characters, no controls |
| `capability` | `observe`, `operate`, `manage`, or `enroll` |
| `state` | `active`, `expired`, or `revoked` |
| `builtin` | True only for the protected local actor |
| `created_at` | UTC timestamp, immutable |
| `updated_at` | UTC timestamp |
| `expires_at` | Optional UTC timestamp |

The built-in actor is unique, active, non-expiring, and fixed at Enroll. Non-built-in actors may transition from active to revoked, and active actors become effectively expired after `expires_at` even before their stored state is normalized.

## Audit Event

| Field | Rule |
|---|---|
| `id` | Opaque UUID |
| `actor_id` | Optional only when identity resolution failed |
| `daemon_id` | Current persisted daemon UUID |
| `operation` | Identifier from the shared catalog |
| `target_kind` | Non-secret classification |
| `target_id` | Optional non-secret identifier |
| `result` | `uncertain`, `succeeded`, `failed`, or `denied` |
| `correlation_id` | Opaque UUID for request correlation |
| `occurred_at` | UTC intent or denial time |
| `completed_at` | Optional UTC completion time |

## Operation Definition

An operation definition contains method, path pattern, stable identifier, required capability, target kind, and audit class (`none`, `privileged_read`, or `mutation`). It is static application metadata rather than persisted policy.

## Storage and migration

Schema v17 adds `actors` and `audit_events`, supporting indexes for built-in actor uniqueness, chronological audit traversal, actor filtering, operation filtering, and result filtering. Migration preserves all scheduler and daemon identity rows and initializes the built-in actor only after migrations succeed.
