# Data Model: Remote Connection Profiles and Target-Safe Clients

## Profile Collection

| Field | Type | Constraint |
| --- | --- | --- |
| `version` | integer | Exactly 1; higher versions fail closed |
| `active_desktop_profile_id` | string | Empty means This computer; otherwise references one profile |
| `profiles` | ordered array | At most 100 unique profile IDs |

## Connection Profile

| Field | Type | Constraint |
| --- | --- | --- |
| `id` | string | Random opaque local profile ID; immutable and unique |
| `label` | string | 1 through 80 Unicode characters; controls rejected; duplicates allowed |
| `endpoint` | string | Canonical HTTPS origin; no user info, query, fragment, or non-root path |
| `daemon_id` | string | Expected installation ID; immutable during repair |
| `credential_id` | string | Native-store lookup key; never a bearer value |
| `certificate_pem` | string | One or more trusted certificates; no private key block |
| `certificate_fingerprint` | string | Lowercase SHA-256 hex of the normalized leaf certificate |
| `client_kind` | string | `desktop`, `cli`, or `json` |
| `capability` | string | `observe`, `operate`, `manage`, or `enroll` |
| `daemon_display_name` | string | Last successfully negotiated daemon label |
| `platform` | string | Last negotiated execution platform |
| `architecture` | string | Last negotiated execution architecture |
| `product_version` | string | Last negotiated daemon version |
| `created_at` | RFC 3339 timestamp | UTC and immutable |
| `updated_at` | RFC 3339 timestamp | UTC and monotonic per committed update |
| `last_successful_at` | RFC 3339 timestamp | Optional, safe cached status only |

## Target Selection

| Mode | Fields | Resolution |
| --- | --- | --- |
| Local | none | Existing protected IPC client |
| Named profile | profile ID or unambiguous label | Exact stored profile plus native credential |
| Explicit remote | endpoint, daemon ID, credential ID, certificate file | Invocation-only target; not persisted |

Named selection rejects an ambiguous duplicate label and recommends the immutable profile ID. Explicit selection is all-or-nothing and conflicts with named selection.

## Lifecycle

```text
pair -> active profile -> rename -> active profile
active profile -> repair -> active profile with replacement credential
active profile -> remove -> deleted metadata and deleted native credential
active desktop profile -> remove -> This computer selected, then deletion
```

No lifecycle transition reveals the bearer credential after initial enrollment handoff to native storage.
