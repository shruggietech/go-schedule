# Data Model: MCP Operate Authority

## Runtime session

| Field | Rules |
| --- | --- |
| Session ID | Random UUID, public correlation value |
| Actor ID | Existing SQLite MCP actor identifier |
| Client name | Valid normalized actor display name |
| Capability | Observe or Operate only in S091 |
| Secret digest | SHA-256 digest in memory; raw secret returned once |
| Created at | UTC timestamp |
| Revoked | Current actor state; the digest mapping remains in memory only so a revoked attempt can be attributed and denied |

## Tool input

| Field | Rules |
| --- | --- |
| `daemon_id` | Required, exact current installation identity |
| `task_id` | Required, non-empty, bounded identifier |
| `request_id` | Required UUID used for retry deduplication |

## Tool result

| Field | Rules |
| --- | --- |
| `schema_version` | `1` |
| `permission` | `operate` |
| `operation` | `tasks.run_now`, `tasks.enable`, or `tasks.disable` |
| `daemon_id` | Validated daemon identity |
| `task_id` | Requested task identity |
| `request_id` | Caller identity for this logical attempt |
| `outcome` | accepted, rejected, denied, or uncertain |
| `message` | Bounded safe explanation with no execution inputs |

## Deduplication record

| Field | Rules |
| --- | --- |
| Request ID | UUID key |
| Fingerprint | Operation plus daemon and task identity |
| Result | Immutable first result |
| Expires at | Bounded retention timestamp |

## State transitions

```text
created -> active -> revoked
                 -> expired
```

An active session resolves to its actor. A revoked or expired session still resolves to its inactive actor so authorization can record an attributed denial. Unknown and malformed secrets resolve to no actor and fail closed.
