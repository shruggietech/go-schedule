# Data Model: S098

## BundleDocument

| Field | Meaning |
|---|---|
| `schema` | Fixed `go-schedule.bundle/v1` identifier |
| `groups` | Canonically ordered portable groups |
| `tasks` | Portable task intent with schedule and policies only |
| `chains` | Portable source and target task identities |
| `trigger_sets` | Secret-free trigger-set intent |
| `triggers` | Secret-free external-trigger intent |
| `watchers` | Portable watcher selection metadata, or a compatibility exclusion when path semantics cannot move |
| `notification_policies` | References only, never channels or destinations |
| `exclusions` | Explicit omitted machine-local or protected material |

## PortableIdentity

`portable_id` is a stable opaque UUID associated with a persisted transferable object. It is independent from the daemon-local record ID and immutable after creation. References in the document always use it.

## BundlePlan

| Field | Meaning |
|---|---|
| `id` | Random preview identity generated per plan |
| `bundle_digest` | SHA-256 of canonical document bytes |
| `target_daemon_id` | Explicit selected daemon identity |
| `target_fingerprint` | SHA-256 of canonical target snapshot at preview time |
| `items` | Deterministically ordered create, update, unchanged, conflict, incompatible, skipped, or target-only drift results |

## BundlePlanItem

Each item has an object kind, portable identity, display name, action, reason, optional dependency identities, and terminal apply outcome. Apply executes only `create` and `update` plan items, leaving all other actions untouched.

## State Transitions

```text
Bundle document -> validate -> target-bound preview -> explicit confirmation -> apply item outcomes
                                    |                                        |
                                    +-> compare drift (read-only)             +-> no replay or rollback
```
