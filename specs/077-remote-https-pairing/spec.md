# Feature Specification: Authenticated Remote Access and Pairing

**Feature Branch**: `codex/077-remote-https-pairing`

**Created**: 2026-09-09

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Implemented on review branch `codex/077-remote-https-pairing`; canonical verification is recorded in `verification.md` and pull-request review remains pending.

**Input**: GitHub issues #168 and #169, implementing the reviewed v1.4 remote-access architecture from S074 on the identity and authorization foundations delivered by S075 and S076.

## Clarifications

### Session 2026-09-09

- Remote access is one opt-in HTTPS listener. Local IPC remains enabled, unchanged, credential-free, and the only listener in default configuration.
- S077 completes the server boundary and enrollment lifecycle through local administration plus a minimal desktop pairing journey that stores the issued credential in the native operating-system credential store. Persistent profile switching, repair, and target-wide navigation remain in #170; named CLI profiles remain in #171.
- A pairing phrase is valid for 10 minutes, permits at most five failed attempts, is single-use, and contains at least 60 bits of randomly selected entropy.
- Durable credentials are independent 256-bit opaque bearer values, returned once, retained only as digests, and attached to server-owned actors whose current authority is loaded on every request.
- The first remote contract exposes a bounded useful subset of existing operations across Observe, Operate, Manage, and Enroll. Local-only runtime paths, MCP listener administration, and secret reveal or rotation endpoints are excluded.

## User Scenarios & Testing

### User Story 1 - Enable a safe remote boundary (Priority: P1)

As an administrator, I want remote access to exist only after I provide an exact address and trusted certificate so installation, upgrade, and ordinary local use never expose the daemon accidentally.

**Why this priority**: An explicit encrypted listener is the non-negotiable boundary for every remote operation and enrollment exchange.

**Independent Test**: Start the daemon with default, invalid, loopback, private-address, wildcard, and public-listener configurations; inspect open sockets and negotiate supported and rejected transport versions.

**Acceptance Scenarios**:

1. **Given** default or explicitly disabled remote settings, **When** the daemon starts, **Then** no TCP listener exists and local IPC remains available.
2. **Given** enabled settings with an exact bind address and valid certificate pair, **When** the daemon starts, **Then** it serves the remote contract through TLS 1.3 while preserving local IPC.
3. **Given** incomplete, plaintext, invalid-certificate, wildcard, or publicly exposed settings without explicit acknowledgement, **When** startup is attempted, **Then** startup fails before binding.
4. **Given** remote access is shutting down, **When** requests or streams are active, **Then** new acceptance stops and bounded graceful shutdown does not interrupt local persistence integrity.

### User Story 2 - Pair a named client once (Priority: P1)

As an administrator using an SSH terminal, I want to create a short-lived word phrase that a remote client can exchange once for a named durable relationship after verifying the intended daemon identity.

**Why this priority**: Pairing supplies the production authentication material needed by the remote listener without creating accounts, reusable passwords, or custom cryptography.

**Independent Test**: Create a phrase locally, complete the desktop pairing form over trusted HTTPS with the expected daemon identity, verify protected credential storage, use the returned credential, then repeat with replay, mismatch, expiry, cancellation, exhaustion, unsupported credential storage, and concurrent exchange attempts.

**Acceptance Scenarios**:

1. **Given** Enroll authority on local IPC, **When** an administrator creates a phrase for a client name, kind, and capability, **Then** the phrase is displayed exactly once with its daemon identity and expiry.
2. **Given** an unused phrase and matching daemon identity, **When** one remote client exchanges it, **Then** the phrase is atomically consumed and one independent actor and credential are created.
3. **Given** a consumed, cancelled, expired, exhausted, malformed, or wrong-daemon phrase, **When** exchange is attempted, **Then** the same generic failure is returned and no credential or actor is created.
4. **Given** concurrent exchanges of one phrase, **When** they race, **Then** exactly one can succeed.
5. **Given** a supported native credential store, **When** desktop pairing succeeds, **Then** the credential is stored under daemon and client identities without appearing in desktop state, settings files, logs, or diagnostics.
6. **Given** no supported native credential store, **When** desktop pairing is attempted, **Then** the flow fails closed without writing a plaintext fallback or consuming the phrase.

### User Story 3 - Use and revoke durable authority (Priority: P1)

As a paired client, I want my opaque credential to authorize only operations allowed by my current actor capability, and as an administrator I want rotation or revocation to take effect immediately.

**Why this priority**: Durable authentication is useful only if it remains bounded by the shared S076 authorization and audit model.

**Independent Test**: Exercise every remote allowlist operation with missing, malformed, unknown, valid, underprivileged, expired, revoked, rotated, and rate-limited credentials while checking side effects and audit evidence.

**Acceptance Scenarios**:

1. **Given** a valid credential, **When** its actor invokes an allowed operation, **Then** the result matches the local operation and carries safe daemon and correlation context.
2. **Given** missing, malformed, unknown, expired, or revoked authentication, **When** a protected operation is requested, **Then** one generic `401` response is returned without revealing credential state.
3. **Given** valid authentication without sufficient authority, **When** a stronger operation is requested, **Then** the shared authorization model returns `403`, records denial, and performs no side effect.
4. **Given** a rotated or revoked credential, **When** the old value is used for the next request, **Then** authentication fails immediately.

### User Story 4 - Consume a stable, bounded contract (Priority: P2)

As a client developer, I want a machine-readable versioned contract with predictable errors, request limits, and retry classifications so integrations can avoid unsafe guesses.

**Why this priority**: A reviewable contract prevents the network surface from drifting into an accidental mirror of privileged local routes.

**Independent Test**: Compare the machine-readable contract, remote allowlist, and registered handlers; exercise malformed media types, oversized inputs, unsupported paths, hostile browser origins, rate limits, event reconnection, and local-equivalence cases.

**Acceptance Scenarios**:

1. **Given** a clean checkout, **When** contract verification runs, **Then** every registered remote route has one operation record and no undocumented route is reachable.
2. **Given** an unsupported path, method, media type, body, API major, or browser origin, **When** a request is made, **Then** it fails with a stable bounded response before application behavior.
3. **Given** a safe read or event stream, **When** ordinary interruption occurs, **Then** its retry classification is explicit; no state-changing request is automatically replayed.

### Edge Cases

- Header bytes, body bytes, collection limits, timeouts, and concurrent connections are bounded before expensive processing.
- Forwarding headers are ignored in S077 because trusted-proxy attribution is not required to authenticate or authorize a bearer actor.
- Phrase verification consumes an attempt for every syntactically valid but incorrect phrase lookup without revealing which component failed.
- Credential, phrase, certificate-key, authorization-header, and verifier values never appear in logs, errors, audit records, exports, events, diagnostics, or command history guidance.
- A database upgrade creates no phrase or credential and never enables remote access.
- A certificate or key file becoming unreadable prevents a new listener from starting; it does not weaken or disable local IPC.
- An event stream periodically reloads credential and actor state and terminates within 30 seconds after revocation or expiry.
- Mutations whose response is lost remain non-replayable; clients must refresh authoritative state before deliberate retry.

## Requirements

### Functional Requirements

- **FR-001**: Remote access MUST default to disabled and MUST require an explicit enabled flag, exact bind address, certificate path, and private-key path.
- **FR-002**: Wildcard or non-private exposure MUST require a separate explicit acknowledgement, and invalid remote configuration MUST fail before any TCP bind.
- **FR-003**: The remote listener MUST accept TLS 1.3 or newer only, refuse plaintext, bound headers and connection timeouts, and shut down gracefully within five seconds.
- **FR-004**: Local IPC MUST remain enabled, offline-capable, and governed by its existing operating-system access controls regardless of remote configuration.
- **FR-005**: Remote routes MUST live under `/api/v1`, be admitted through a separate explicit allowlist, and MUST NOT mount or derive access from the complete local mux.
- **FR-006**: Every allowed remote operation MUST declare a stable operation ID, local equivalent, capability, target class, audit class, retry class, body limit, and secret exclusions.
- **FR-007**: Runtime paths, localhost MCP administration, external-trigger secret reveal or rotation, notification-secret rotation, and undeclared future routes MUST remain unreachable remotely.
- **FR-008**: Protected remote requests MUST accept credentials only through RFC 6750 `Authorization: Bearer` syntax and MUST reject credential-bearing query parameters, cookies, and alternate headers.
- **FR-009**: Durable credentials MUST contain 256 random bits, use opaque base64url encoding, be returned only at issuance or rotation, and be retained only as a digest plus safe fingerprint.
- **FR-010**: Authentication MUST load the current credential and actor state on every request and MUST use constant-time verifier comparison.
- **FR-011**: Missing, malformed, unknown, expired, or revoked credentials MUST produce the same stable `401` response with `WWW-Authenticate: Bearer` and no identifying detail.
- **FR-012**: Authenticated requests MUST pass the shared S076 capability and intent-first audit boundary before application behavior.
- **FR-013**: The system MUST apply bounded source and actor token-bucket rate limits with stable `429` responses and bounded `Retry-After` values.
- **FR-014**: An Enroll-authorized local administrator MUST be able to create and cancel pairing phrases and list active phrase metadata without retrieving phrase values.
- **FR-015**: Pairing phrases MUST contain at least 60 bits of random entropy, expire after 10 minutes, permit five failed attempts, be single-use, and be stored only as uniquely salted Argon2id verifiers.
- **FR-016**: Enrollment exchange MUST require the expected daemon installation identity, requested client display name, actor kind, and phrase through one isolated HTTPS endpoint with generic failure behavior.
- **FR-017**: Successful phrase exchange MUST atomically consume the phrase and create exactly one actor and one independent durable credential.
- **FR-018**: Concurrent exchange, cancellation, expiry, and attempt exhaustion MUST never create more than one client relationship from one phrase.
- **FR-019**: The desktop MUST provide a bounded pairing form for HTTPS address, expected daemon identity, trusted-certificate input, phrase, client display name, and requested capability; it MUST store successful credentials only through supported native operating-system credential storage.
- **FR-020**: Desktop pairing MUST fail before exchange when native credential storage is unavailable, retain no bearer value in application state after storage, and provide no plaintext or application-file fallback.
- **FR-021**: An Enroll-authorized local administrator MUST be able to list safe credential metadata, rotate a credential, and revoke a credential together with its actor relationship.
- **FR-022**: Rotation MUST atomically replace the digest and invalidate the previous credential; revocation MUST reject the next request and close existing event streams within 30 seconds.
- **FR-023**: A machine-readable OpenAPI 3.1 contract MUST describe every remote route, authentication rule, bounded schema, stable error, and response shape, with automated allowlist and handler completeness checks.
- **FR-024**: Requests carrying an `Origin` header MUST be rejected and remote responses MUST NOT emit permissive cross-origin headers.
- **FR-025**: Remote responses, logs, audit records, events, diagnostics, exports, and desktop state MUST structurally exclude bearer values, phrases, digests, verifiers, certificate keys, and existing task or notification secrets.
- **FR-026**: Documentation MUST cover enablement, disablement, certificate trust, private-network, SSH-tunnel, reverse-proxy, and direct-HTTPS use, exposure acknowledgement, phrase enrollment, credential rotation and revocation, error behavior, and unsupported deployment claims.
- **FR-027**: This slice MUST NOT add automatic certificate issuance, plaintext listener modes, JWTs, mutual-TLS client identity, OAuth, browser CORS access, remote MCP mutation, service discovery, NAT traversal, offline mutation queues, or general policy rules.

### Key Entities

- **Remote Listener Configuration**: Explicit enablement, bind address, certificate and key locations, and exposure acknowledgement.
- **Remote Operation Record**: One fail-closed declaration connecting a remote method and path to an existing local operation and its security bounds.
- **Pairing Session**: A one-time enrollment intent with verifier, salt, client metadata, capability, expiry, attempts, and terminal state.
- **Client Credential**: A durable server-side relationship containing an identifier, actor relationship, digest, safe fingerprint, lifecycle state, and timestamps.
- **Enrollment Result**: The one-time response containing daemon identity, actor metadata, credential identifier, fingerprint, and raw bearer value.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Clean install, upgrade, restart, and disabled-configuration tests observe zero TCP listeners.
- **SC-002**: One valid local configuration serves the documented remote contract through TLS 1.3 on Linux, macOS, and Windows while local IPC remains functional.
- **SC-003**: Contract tests account for 100 percent of registered remote handlers and allowlist records, with zero local-only or secret-management routes reachable.
- **SC-004**: A successful enrollment completes in one phrase exchange, returns one credential once, and creates exactly one actor and credential under concurrent load.
- **SC-005**: Replay, mismatch, cancellation, expiry, exhaustion, malformed input, and wrong-daemon cases create zero credentials and return indistinguishable failure envelopes.
- **SC-006**: Authentication and capability matrices cover every credential state and all four capability levels without unauthorized side effects.
- **SC-007**: Rotation and revocation invalidate old credentials on the next request, and active event streams lose authority within 30 seconds.
- **SC-008**: Burst tests prove source and actor limits return bounded `429` responses without credential-state disclosure or unbounded limiter growth.
- **SC-009**: Structural canary tests find zero protected values in responses, logs, audits, events, diagnostics, or exports.
- **SC-010**: Full repository verification, race detection, coverage, documentation, contract drift, and security checks pass.

## Assumptions

- S074 remains the architecture of record, and S075/S076 provide daemon identity, actor authority, operation classification, and audit persistence.
- Operators provide and renew trusted server certificates; automatic certificate lifecycle is outside v1.4.
- S077 adds only the desktop enrollment and native credential-storage handoff required by #169. Persistent connection profiles, switching, repair, removal, selected-target context, and named CLI profiles remain in #170 and #171.
- Remote resilience beyond stream revocation and bounded server shutdown remains in #172.

## Dependencies

- Parent roadmap issue: #18.
- Completed dependencies: #165, #166, and #167.
- This slice completes: #168 and #169.
- Follow-on consumers: #170, #171, #172, and #173.
