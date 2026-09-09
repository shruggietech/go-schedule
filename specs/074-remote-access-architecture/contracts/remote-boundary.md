# Contract: Remote Access Boundary

This contract governs issues #166 through #173. It describes the review gate for future implementation and does not claim that remote access exists in S074.

## Boundary order

Every remote connection and request must cross these controls in order before daemon application behavior runs:

1. Explicitly configured listener and bind address.
2. TLS 1.3 handshake using an operator-trusted server certificate.
3. Coarse source rate limit.
4. Exact remote API major path and bounded headers.
5. Bearer credential authentication or the isolated enrollment exchange.
6. Current actor-state and expiration check.
7. Authenticated actor rate limit.
8. Operation allowlist lookup and capability authorization.
9. Request content-type, schema, and body-size validation.
10. Correlation and audit classification.
11. Existing daemon application operation.
12. Bounded response, stable error envelope, and secret-free audit completion.

Failure at any step prevents every later step. Authentication failures do not distinguish unknown, expired, revoked, or malformed credentials to the remote caller.

## Remote allowlist record

Every operation admitted by #168 must have one reviewable record containing:

| Required property | Rule |
| --- | --- |
| OpenAPI operation ID | Stable and unique within the API major |
| Method and path | Under `/api/v1`; enrollment is isolated under the same major |
| Local equivalent | Names the shared daemon application operation, or explains the discovery or enrollment exception |
| Capability | Exactly Observe, Operate, Manage, or Enroll |
| Input limit | Maximum header and body bytes plus collection bounds |
| Audit class | None, protected read, mutation, or enrollment |
| Retry class | Safe read, resumable stream, or no automatic replay |
| Secrets | Explicit input and output exclusions |
| Acceptance tests | Authentication, authorization, schema, bounds, equivalence, and failure cases |

Local-only endpoints are denied by omission. The initial denylist includes runtime filesystem paths, localhost MCP lifecycle and credential issuance, external-trigger reveal and rotate operations, and any future endpoint lacking a completed allowlist record.

## Capability semantics

| Capability | Permitted intent | Explicit exclusions |
| --- | --- | --- |
| Observe | Read safe scheduler state, history, alerts, and live activity | Executable configuration, secret reveal, task execution, mutation |
| Operate | Observe plus deliberate task execution and non-configuration operational acknowledgement | Task, group, automation, credential, or notification configuration |
| Manage | Operate plus scheduler object and notification-policy administration | Actor authority, pairing, credential issuance, server exposure configuration |
| Enroll | Manage plus pairing phrase creation, actor capability assignment, credential rotation and revocation | General identity federation, arbitrary policy, certificate-authority operations |

The classes are monotonic for product comprehension, but every operation still declares one minimum class. Local OS authorization does not become a bearer credential and remains enforced by the IPC layer.

## Versioning and compatibility

- Local IPC keeps its current `/v1` contract. Remote paths use `/api/v1` so network compatibility can evolve without silently changing local transport assumptions.
- Additive optional response fields, new operations, and new declared capability-manifest entries are compatible within a major.
- Removing or renaming fields, changing meanings, tightening accepted values for existing clients, or changing an operation's success semantics requires a new remote path major.
- The capability manifest reports daemon identity, product version, supported remote API majors, feature capabilities, and minimum compatible client information before ordinary actions are shown.
- A new major must overlap the preceding supported major for a documented migration interval defined by its delivery issue. No calendar deadline is invented by this architecture.
- Deprecation appears in the OpenAPI source, client diagnostics, documentation, and release notes before removal.
- Generated server and client boundaries must reproduce cleanly from the pinned OpenAPI source. Drift is a failing contract test.

## Failure contract

| Condition | Remote result | Required evidence |
| --- | --- | --- |
| Missing, malformed, unknown, expired, or revoked bearer credential | `401` with one generic stable error code and `WWW-Authenticate: Bearer` | Secret-free failed-authentication security event, subject to aggregation |
| Authenticated actor lacks capability | `403` with stable denial code | Actor-attributed audit denial |
| Unsupported API major or incompatible capability | `404` for unknown path or `409` for negotiated incompatibility | Safe compatibility diagnostic |
| Invalid content type or schema | `400` or `415` with bounded field error | Actor-attributed failed request without protected input |
| Header or body exceeds limit | `413` or connection-level header rejection | Bounded security diagnostic |
| Source or actor rate exceeded | `429` with bounded `Retry-After` | Aggregate limiter evidence without credential disclosure |
| Daemon conflict or domain validation failure | Existing stable JSON error envelope | Actor-attributed result and target identity |
| Client loses a mutation response | No automatic retry | Client refreshes authoritative state and asks before deliberate resubmission |
| Event stream disconnects | Resume from last accepted event identity | Duplicate event delivery is tolerated; no mutation occurs |

## Deployment contract

| Mode | Listener and TLS | Product ownership | Operator ownership | Support posture |
| --- | --- | --- | --- | --- |
| Local IPC | Unix socket or Windows named pipe; no TCP | IPC permissions, local API, daemon lifecycle | Local OS account and group administration | Default and unchanged |
| Private-network HTTPS | Explicit private address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit | VPN or private routing, firewall, DNS, trusted certificate | Recommended remote mode |
| SSH-tunneled HTTPS | Explicit loopback or private address with direct go-schedule TLS inside SSH | HTTPS boundary and all application controls | SSH identity, tunnel lifecycle, port forwarding, trusted certificate | Recommended for headless administration |
| Reverse-proxied HTTPS | Explicit private address with direct go-schedule TLS; proxy adds public TLS | Application TLS configuration validation, authentication, authorization, trusted-proxy allowlist, limits, audit | Proxy lifecycle; public and backend certificate issuance and renewal; backend hostname trust and key permissions; routing, forwarded-header policy, firewall | Supported with explicit peer trust |
| Direct HTTPS | Explicit network address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit, fail-closed startup | Certificate issuance and renewal, private-key permissions, DNS, firewall, exposure monitoring | Supported for experienced operators, never automatic |

All network modes require an explicit enable flag, exact bind address, certificate and key, and exposure acknowledgement for wildcard or public addresses. Invalid configuration fails before bind. Disabling remote access stops accepting connections, drains bounded in-flight work, terminates streams, and leaves local IPC available.

## Non-goals

- User accounts, passwords, SSO, OAuth authorization-server behavior, federation, teams, or multi-tenancy.
- JWTs, custom encryption, a custom certificate authority, automatic ACME, NAT traversal, service discovery, or automatic public exposure.
- Per-resource ACLs, policy expressions, deny rules, inheritance, or an authorization language.
- Offline mutation queues, background replay of uncertain writes, multi-master synchronization, or conflict resolution.
- Remote MCP mutation authority. MCP remains independently bounded by its own Observe contract until a later reviewed issue explicitly changes it.
- Credential recovery or export. Lost credentials are replaced, not reconstructed.

## Required downstream sequence

1. #166 establishes stable daemon identity and capability discovery.
2. #167 establishes actors, the capability matrix, and durable audit records.
3. #168 implements the opt-in HTTPS boundary and OpenAPI contract.
4. #169 implements one-time phrase pairing and native client credential storage.
5. #170 and #171 add desktop and CLI or JSON client profiles.
6. #172 adds bounded reconnect and stale-state behavior.
7. #173 qualifies supported deployments and the complete release.

No network implementation begins until S074 is reviewed and merged.
