# Data Model: Remote Access Architecture

S074 defines conceptual boundaries only. It adds no database schema. The owning downstream issues must preserve these relationships and lifecycle invariants.

## Daemon installation

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `installation_id` | Stable opaque target identity | Unique across clones after deliberate reset handling | #166 |
| `display_name` | User-editable target label | Not an authorization or uniqueness key | #166 |
| `product_version` | Running daemon version | Safe discovery value | #166 |
| `api_majors` | Supported remote API path majors | Non-empty only when remote support exists | #166 and #168 |
| `capabilities` | Supported product features | Contains no secret or unnecessary host identity | #166 |
| `remote_mode` | Disabled or configured listener mode | Disabled by default | #168 |

## Actor

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `actor_id` | Stable server-owned identity | Never derived from display name or credential | #167 |
| `kind` | Local OS, desktop, CLI, JSON, or MCP client | Closed documented vocabulary | #167 |
| `display_name` | Human-recognizable client name | Duplicate names do not merge identities | #167 |
| `capability` | Observe, Operate, Manage, or Enroll | No general policy expression | #167 |
| `state` | Active, expired, or revoked | Revoked actors authorize nothing | #167 |
| `created_at` | UTC creation time | RFC 3339 in API output | #167 |
| `expires_at` | Optional UTC expiry | Checked before every protected operation | #167 |

## Durable credential

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `credential_id` | Stable revocation identity | Independent for every client installation | #169 |
| `actor_id` | Owning actor | One actor may rotate credentials without changing identity | #169 |
| `digest` | SHA-256 of 256-bit opaque bearer value | Raw value is never persisted server-side | #169 |
| `fingerprint` | Short safe display identifier | Cannot authenticate or reconstruct the bearer value | #169 |
| `issued_at` | UTC issuance time | Recorded without raw secret | #169 |
| `expires_at` | Required or policy-defined expiry | Expired values fail as unauthorized | #169 |
| `revoked_at` | Optional UTC revocation time | Revocation terminates ongoing authority within the documented bound | #169 |

The raw value exists only in the successful issuance response and the client's native credential store. Rotation creates a new raw value and digest, invalidates the old digest atomically, and does not reuse enrollment material.

## Enrollment phrase

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `enrollment_id` | Attempt-lifecycle identity | Not a durable actor identity | #169 |
| `phrase_digest` | One-way verifier for the displayed phrase | Phrase cannot be recovered from server state | #169 |
| `daemon_id` | Intended target | Client verifies it before storing a credential | #166 and #169 |
| `expires_at` | Short UTC deadline | Expiration is mandatory | #169 |
| `attempts_remaining` | Bounded guessing budget | Failure response does not reveal existence | #169 |
| `state` | Active, consumed, expired, or cancelled | Only Active may be exchanged once | #169 |

## Remote operation

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `operation_id` | Stable OpenAPI operation identity | Unique within one API major | #168 |
| `method` and `path` | HTTP boundary | Under `/api/v1`; the local `/v1` path remains unchanged | #168 |
| `required_capability` | Minimum actor authority | Exactly one of the four capability classes | #167 and #168 |
| `request_limit` | Maximum headers and body | Explicit for every operation class | #168 |
| `audit_class` | None, protected read, mutation, or enrollment | Secrets and protected inputs are never recorded | #167 |
| `retry_class` | Safe read, resumable stream, or no automatic replay | Mutations are never automatically replayed | #168 and #172 |
| `local_equivalent` | Existing daemon application operation | Required except remote-only discovery and enrollment | #168 |

## Audit event

| Field | Meaning | Constraint | Owner |
| --- | --- | --- | --- |
| `event_id` | Stable record identity | Supports bounded pagination and export | #167 |
| `actor_id` | Requesting actor | Present for every authenticated remote event | #167 |
| `daemon_id` | Target installation | Prevents same-name target ambiguity | #166 and #167 |
| `operation_id` | Requested remote operation | Stable across display-copy changes | #167 |
| `target_kind` and `target_id` | Affected domain object | Omitted only when no object exists | #167 |
| `result` | Allowed success, allowed failure, or denied | Authorization denial is visible | #167 |
| `correlation_id` | Safe request identity | Never a credential or pairing secret | #167 |
| `occurred_at` | UTC event time | RFC 3339 in API output | #167 |

## Lifecycle relationships

1. A daemon installation exposes no remote operation while remote mode is disabled.
2. An administrator creates one enrollment phrase through an already authorized local or remote Enroll operation.
3. A client verifies daemon identity and exchanges the active phrase once.
4. The daemon creates a named actor and independent durable credential, stores only the digest, and marks the phrase consumed atomically.
5. A protected request authenticates the credential, loads current actor state, checks capability, applies bounds, runs one allowlisted operation, and emits the required audit event.
6. Rotation replaces only credential material. Revocation or expiry removes authority from every later request and ongoing stream within the documented bound.
7. Recovery never recreates a lost raw credential. The administrator revokes or replaces the client relationship through an independently authorized path.
