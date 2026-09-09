---
title: Remote access architecture
nav_order: 8.6
---

# Remote access architecture

**Current status:** Architecture contract only. No remote listener is implemented.

**Primary remote transport:** HTTPS with versioned HTTP/JSON.

This page is the reviewed implementation boundary for the v1.4 remote-access milestone. It does not advertise a shipped network feature. The current daemon still serves its [local API](api.md) through a protected Unix socket or Windows named pipe, opens no remote API port, and remains fully usable offline.

## Design posture

The daemon continues to own scheduling, persistence, execution, authorization decisions, and audit creation. The local IPC mux and remote operation allowlist are separate adapters over shared daemon operations. A route being safe and authenticated over local OS-controlled IPC does not make it safe for remote exposure.

The remote adapter is deny-by-default. Each admitted operation must have a stable OpenAPI operation identifier, minimum capability, input bound, audit class, retry class, secret exclusions, local-equivalence mapping, and acceptance tests. Local-only routes are absent from the allowlist, including runtime filesystem paths, localhost MCP lifecycle and credential issuance, external-trigger secret reveal and rotation, and every future route without a completed remote record.

The design adds no custom encryption, signed token, wire protocol, user directory, or general policy engine. It uses operating-system trust and credential stores, Go's standard cryptographic and HTTP packages, and established narrowly scoped libraries.

## Ordered request boundary

Every remote connection or request crosses these controls in order:

1. The operator has explicitly enabled remote access and supplied an exact bind address.
2. Go's TLS server completes a TLS 1.3 handshake with an operator-trusted certificate.
3. A coarse source limiter bounds unauthenticated connection and credential-guessing work.
4. The server accepts only a supported `/api/v1` path, bounded headers, and an expected method.
5. The request supplies a durable bearer credential, except for the isolated enrollment exchange.
6. Authentication loads current actor state and rejects unknown, expired, or revoked authority.
7. A second limiter bounds the authenticated actor.
8. The remote allowlist resolves the operation and checks its required capability.
9. Content type, OpenAPI schema, body size, collection size, and operation-specific constraints are validated.
10. The request receives a correlation identity and audit classification.
11. The adapter invokes the shared daemon application operation.
12. The server emits a bounded response and completes any required secret-free audit event.

Authentication and capability authorization run before daemon application behavior. A failure at any step prevents every later step. Remote authentication responses do not reveal whether a credential or enrollment phrase was unknown, malformed, expired, consumed, or revoked.

Browser origins are not a v1.4 client surface. Requests carrying an `Origin` header are rejected by default, and no cross-origin resource sharing response headers are emitted. A later browser-client proposal would require an explicit origin allowlist and its own review.

## Remote operation record

Every remotely exposed operation must declare all of these properties in or beside the OpenAPI source:

| Property | Required contract |
| --- | --- |
| Operation identity | Stable and unique within one API major |
| Method and path | Under `/api/v1`; enrollment is isolated within that major |
| Local equivalent | Shared daemon operation, or a reasoned discovery or enrollment exception |
| Capability | Exactly Observe, Operate, Manage, or Enroll |
| Input bound | Header and body bytes plus collection, pagination, and duration limits |
| Audit class | None, protected read, mutation, or enrollment |
| Retry class | Safe read, resumable stream, or no automatic replay |
| Secret contract | Protected request fields and excluded response, log, event, and audit fields |
| Verification | Authentication, authorization, schema, bounds, equivalence, and failure tests |

The initial allowlist is defined only when [issue #168](https://github.com/shruggietech/go-schedule/issues/168) implements the machine-readable contract. No prose statement on this page exposes an endpoint by itself.

## Capability model

Observe, Operate, Manage, and Enroll are the complete capability vocabulary.

| Capability | Permitted intent | Explicit exclusions |
| --- | --- | --- |
| Observe | Read safe scheduler state, history, alerts, and live activity | Executable configuration, secret reveal, task execution, mutation |
| Operate | Observe plus deliberate task execution and non-configuration operational acknowledgement | Task, group, automation, credential, or notification configuration |
| Manage | Operate plus scheduler object and notification-policy administration | Actor authority, enrollment, credential issuance, and server exposure configuration |
| Enroll | Manage plus pairing phrase creation, actor capability assignment, credential rotation, and revocation | Federation, arbitrary policy, certificate-authority operations |

The classes are monotonic for product comprehension, but every operation still declares one minimum class. The daemon loads current server-owned authority for every request. Capability data is never trusted from a bearer value, a forwarded header, a display name, or client-supplied claims.

Local IPC authorization remains independent. A local OS caller does not receive an implicit bearer credential, and enabling HTTPS cannot widen Unix-socket ownership or Windows named-pipe access.

## Credential lifecycle

An enrollment phrase is short-lived and single-use. It exists only to bootstrap one named client relationship, is visibly tied to the intended daemon identity, has a strict expiry and attempt budget, is retained only through a salted Argon2id verifier, and cannot be used as an ordinary API credential. Successful exchange atomically consumes the phrase and issues unrelated durable material. Cancellation, expiry, exhaustion, or success permanently ends that phrase.

A durable credential is an opaque 256-bit bearer value retained only as a digest. Go's `crypto/rand` creates the raw value, base64url is its opaque transport encoding, RFC 6750 `Authorization: Bearer` syntax carries it, SHA-256 stores a verifier, and constant-time comparison checks it. The token contains no claims, authority, identity, expiry, or signature. Server-owned actor and credential records provide those facts so revocation takes effect immediately.

Each client installation receives an independent credential and safe fingerprint. Issuance returns the raw value once. The desktop and named CLI profiles store it through supported operating-system credential storage, never in application JSON, process arguments, shell-history examples, logs, exports, diagnostics, or screenshots. Rotation replaces the verifier atomically and invalidates the prior value. A lost credential is revoked and replaced, never recovered or re-derived from an enrollment phrase.

The daemon revalidates actor and credential state periodically during a live stream. Revocation, expiry, or actor disablement terminates authority for new requests immediately and closes existing streams within the documented revalidation bound defined by #167 and #168.

## API compatibility and live updates

OpenAPI 3.1 is the source of truth for the remote contract. Issue #168 will pin `github.com/oapi-codegen/oapi-codegen/v2` as a Go tool and generate strict standard-library HTTP server and client boundaries. Authentication and authorization stay explicit middleware because generated routing does not implement product security policy. Generated files are committed, and a clean regeneration check prevents source, server, and client drift.

Local IPC retains `/v1`. Remote paths use `/api/v1` so network compatibility can evolve without silently changing local transport assumptions. Additive optional response fields, new operations, and new manifest capabilities are compatible within a major. Removing or renaming fields, changing their meaning, tightening previously accepted values, or changing success semantics requires a new path major.

The capability manifest reports stable daemon identity, product version, supported API majors, features, and safe compatibility facts before clients present actions. A new major overlaps the prior supported major for a documented migration interval defined by its delivery issue. Deprecation appears in OpenAPI, client diagnostics, documentation, and release notes before removal. This architecture invents no calendar deadline.

Server-Sent Events provide authenticated one-way live activity over HTTPS. Each event has stable resumable identity, and clients reconnect with `Last-Event-ID`. Duplicate event delivery is tolerated because events do not mutate daemon state. Safe reads and streams may retry within documented bounds. Mutations have no automatic replay when a disconnect leaves their outcome uncertain; the client refreshes authoritative state and requires a deliberate resubmission.

## Failure behavior

| Condition | Remote result | Required evidence |
| --- | --- | --- |
| Missing, malformed, unknown, expired, or revoked bearer credential | `401`, one generic stable error code, and `WWW-Authenticate: Bearer` | Bounded secret-free security evidence, aggregated where necessary |
| Authenticated actor lacks capability | `403` and stable denial code | Actor-attributed audit denial |
| Unsupported API path or negotiated incompatibility | `404` for an unknown path or `409` for known incompatibility | Safe compatibility diagnostic |
| Invalid content type or schema | `400` or `415` with bounded field detail | Actor-attributed failed request without protected input |
| Header or body exceeds its limit | `413` or connection-level header rejection | Bounded security diagnostic |
| Source or actor rate is exceeded | `429` with bounded `Retry-After` | Aggregate limiter evidence without credential disclosure |
| Daemon conflict or domain validation fails | Existing stable JSON error envelope | Actor, operation, result, and target identity |
| Client loses a mutation response | No automatic replay | Authoritative refresh before explicit retry |
| Event stream disconnects | Resume after the last accepted event identity | Duplicate-safe event delivery without mutation |

HTTP server configuration includes finite read-header, read, idle, and operation-aware write behavior, bounded maximum headers, bounded request bodies, graceful shutdown, and connection-lifecycle ownership. Streaming endpoints use heartbeat and maximum-idle rules instead of an ordinary fixed response timeout. Shutdown stops accepting connections, allows bounded in-flight operations to finish, terminates streams, and leaves local IPC lifecycle independently available until daemon shutdown.

## Deployment modes

All network modes require explicit enablement, an exact bind address, a certificate and key readable only by the daemon identity, and explicit exposure acknowledgement for wildcard or public addresses. Invalid configuration fails before bind. HTTPS is retained inside VPN, tunnel, and proxy deployments so the application boundary has one security contract.

| Mode | Listener and TLS | Product ownership | Operator ownership | Support posture |
| --- | --- | --- | --- | --- |
| Local IPC | Unix socket or Windows named pipe; no TCP | IPC permissions, local API, daemon lifecycle | OS account and group administration | Default and unchanged |
| Private-network HTTPS | Explicit private address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit | VPN or private routing, firewall, DNS, trusted certificate | Recommended remote mode |
| SSH-tunneled HTTPS | Explicit loopback or private address with direct go-schedule TLS inside SSH | HTTPS boundary and all application controls | SSH identity, tunnel lifecycle, port forwarding, trusted certificate | Recommended for headless administration |
| Reverse-proxied HTTPS | Explicit private address with direct go-schedule TLS; proxy adds public TLS | Application TLS configuration validation, authentication, authorization, trusted-proxy allowlist, limits, audit | Proxy lifecycle; public and backend certificate issuance and renewal; backend hostname trust and key permissions; routing, forwarding policy, firewall | Supported with explicit immediate-peer trust |
| Direct HTTPS | Explicit network address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit, fail-closed startup | Certificate issuance and renewal, key permissions, DNS, firewall, exposure monitoring | Supported for experienced operators, never automatic |

Forwarded headers are ignored unless trusted-proxy configuration explicitly matches the immediate peer. Trusting a proxy changes only attribution of client address and scheme metadata; it never delegates bearer authentication or capability authorization. A proxy cannot inject actor identity or bypass application limits.

Disabling remote access stops acceptance, drains bounded in-flight operations, closes streams, clears ephemeral limiter state, and preserves local IPC. Installation, upgrade, restore, and ordinary daemon restart never enable a listener implicitly.

## Threats and verification

The threat boundary covers individually administered daemons and their approved clients. It does not claim defense against an administrator who controls the daemon process, database, executable, or host operating system.

| ID | Risk | Required control | Acceptance-test class |
| --- | --- | --- | --- |
| T01 | Installation or upgrade accidentally exposes a port | Default disabled, explicit address, pre-bind validation | Clean install, upgrade, restart, and socket-inventory tests |
| T02 | Passive capture or transport downgrade reveals data or bearer values | HTTPS everywhere, TLS 1.3 minimum, trusted certificate | Protocol-version, certificate, hostname, and plaintext-refusal tests |
| T03 | An unauthenticated source guesses credentials or phrases | Coarse source limiter, phrase attempt budget, generic failure | Burst, refill, expiry, exhaustion, and enumeration tests |
| T04 | A captured phrase or old credential is replayed | Single-use phrase, independent credential, atomic rotation | Concurrent exchange, replay, rotation, and cancellation tests |
| T05 | A valid client invokes excessive authority | Server-owned capability mapping before handler dispatch | Per-operation allow, deny, stale-authority, and no-side-effect tests |
| T06 | A client spoofs proxy metadata or actor identity | Immediate-peer proxy allowlist, ignored forwarding headers, bearer-owned actor | Trusted and untrusted proxy integration tests |
| T07 | Oversized or slow requests exhaust daemon resources | Header, body, collection, rate, timeout, and connection bounds | Slow-header, large-body, concurrency, cancellation, and leak tests |
| T08 | Secrets appear in responses, logs, events, audit, or diagnostics | Explicit secret inventory, allowlisted response types, structural exclusion | Canary credential, phrase, trigger, environment, and webhook leakage tests |
| T09 | Revoked or expired authority survives in caches or streams | Current-state lookup and periodic stream revalidation | Revocation race, expiry boundary, open-stream, and restart tests |
| T10 | Same-named daemons cause a wrong-target action | Stable daemon identity, manifest negotiation, target-visible mutations | Clone, rename, profile mismatch, confirmation, and stale-cache tests |
| T11 | Version skew changes request meaning or hides capability | URI major, OpenAPI generation, manifest negotiation, overlap policy | Older-client, older-daemon, unknown-major, and deprecation tests |
| T12 | A hostile web origin drives a credentialed browser request | Reject `Origin` by default and emit no permissive CORS policy | Cross-origin preflight, simple-request, and trusted-native-client tests |

Security tests are bounded to these risks and supported modes. A review may add a threat only when it identifies an exposed asset, attacker capability, required control, and executable acceptance class.

## Dependency ownership

Go standard library owns HTTP, TLS, randomness, digests, and constant-time comparison. `golang.org/x/time/rate` owns rate limiting. OpenAPI 3.1 and pinned `oapi-codegen/v2` own API description and generated boundaries. `golang.org/x/crypto/argon2` owns enrollment phrase verification. Pinned `zalando/go-keyring` owns native client credential storage.

| Component | Introduction owner | Maintenance rule | Replacement trigger |
| --- | --- | --- | --- |
| `net/http` and `crypto/*` | #168 and #169 through the supported Go toolchain | Follow Go security releases and retain focused boundary tests | A proven requirement cannot be met by the supported standard library |
| `golang.org/x/time/rate` | #168 | Pin directly, include Dependabot coverage, review releases, and test with explicit timestamps | Maintenance or security posture fails, or required bounds cannot be expressed |
| OpenAPI 3.1 and `oapi-codegen/v2` | #168 | Pin the tool, commit source and output, verify clean regeneration, review license and advisories | Contract drift, maintenance failure, or inability to represent a required stable operation |
| `golang.org/x/crypto/argon2` | #169 | Pin directly, use Argon2id parameters benchmarked and recorded on supported platforms, review Go security releases | Supported guidance changes or resource bounds cannot be met |
| `zalando/go-keyring` | #169 | Permit native Keychain, Credential Manager, and Secret Service only; run native tests; review releases, license, and advisories | Maintenance failure, interactive or plaintext fallback, or supported-platform contract failure |

S074 adds none of these as a new direct dependency. The owning issue selects and pins a then-current reviewed version, proves clean restoration and supported-platform behavior, and records any deviation. Unsupported native credential storage fails closed; there is no application-file fallback.

## Non-goals

No JWT, user-account system, SSO, custom certificate authority, automatic public exposure, general policy language, offline mutation queue, or remote MCP mutation authority is included.

The boundary also excludes federation, teams, multi-tenancy, OAuth authorization-server behavior, automatic ACME, NAT traversal, service discovery, per-resource ACLs, deny expressions, background mutation replay, multi-master synchronization, conflict resolution, credential export, and credential recovery. Lost credentials are replaced through independently authorized administration.

## Required implementation order

1. #166 establishes daemon identity and discovery.
2. #167 establishes actors, capabilities, and audit.
3. #168 implements HTTPS and OpenAPI.
4. #169 implements enrollment and credential storage.
5. #170 and #171 implement desktop and CLI clients.
6. #172 implements resilience.
7. #173 qualifies the release.

No network implementation begins until S074 and [issue #165](https://github.com/shruggietech/go-schedule/issues/165) are reviewed and merged. Each downstream issue remains open until its own acceptance criteria and verification are complete.
