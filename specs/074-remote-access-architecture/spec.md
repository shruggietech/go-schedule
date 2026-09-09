# Feature Specification: Remote Access Architecture

**Feature Branch**: `codex/074-remote-access-architecture`

**Created**: 2026-09-09

**Status**: Implemented

**Delivery**: Review branch `codex/074-remote-access-architecture`; complete architecture, fixture-backed contract, and canonical eight-gate verification passed 2026-09-09.

**Input**: User request to run S074 under autopilot, implementing issue [#165](https://github.com/shruggietech/go-schedule/issues/165).

## User Scenarios & Testing

### User Story 1 - Review a bounded trust boundary (Priority: P1)

As a maintainer, I can review one complete architecture that states what is exposed, who is trusted, and which component owns every security control before any remote listener is implemented.

**Why this priority**: Issues #166 through #173 depend on a stable boundary. Allowing implementation to begin from scattered assumptions would make later identity, authorization, transport, pairing, and client work conflict.

**Independent Test**: A reviewer can start from the architecture page, follow its trust-boundary and deployment-mode tables, and determine the owner and required verification for every network-facing control without consulting chat history.

**Acceptance Scenarios**:

1. **Given** the default installation, **When** the architecture is evaluated, **Then** local IPC remains the only daemon API listener and no TCP port is opened.
2. **Given** remote access is explicitly enabled, **When** a request crosses the network boundary, **Then** HTTPS, authentication, authorization, rate limiting, request bounds, and attributable audit requirements all precede daemon operations.
3. **Given** local and remote callers target the same operation, **When** their authorities are equivalent, **Then** both reach the same daemon-owned domain behavior rather than separate implementations.

---

### User Story 2 - Choose a supported deployment without hidden ownership (Priority: P1)

As an operator, I can distinguish local IPC, private-network HTTPS, SSH-tunneled HTTPS, reverse-proxied HTTPS, and direct HTTPS deployments, including what go-schedule owns and what infrastructure I must provide.

**Why this priority**: Security claims are only meaningful when certificate, routing, name resolution, firewall, proxy, and tunnel responsibilities are explicit.

**Independent Test**: Each documented mode identifies the listener location, encryption boundary, certificate owner, network owner, client trust requirement, and support posture.

**Acceptance Scenarios**:

1. **Given** an operator does not configure remote access, **When** the daemon starts, **Then** existing offline and local-client behavior is unchanged.
2. **Given** an operator uses a VPN or SSH tunnel, **When** reviewing the deployment, **Then** the document recommends the simpler private path while still requiring authenticated HTTPS at the go-schedule boundary.
3. **Given** an operator selects direct public HTTPS, **When** reviewing prerequisites, **Then** operator-supplied certificate lifecycle, DNS, firewall, and exposure monitoring are explicit and no zero-configuration claim is made.

---

### User Story 3 - Build later slices against stable contracts (Priority: P2)

As an implementer, I can derive identity, authorization, HTTPS, pairing, desktop, CLI, and resilience work from a versioning, dependency, credential, failure, and verification policy that does not invent product behavior prematurely.

**Why this priority**: The architecture must unblock future implementation while preserving issue-level sequencing and reviewability.

**Independent Test**: Each blocked issue can cite a named boundary, role, credential lifecycle stage, API compatibility rule, or acceptance-test class and can identify which decisions remain in that issue.

**Acceptance Scenarios**:

1. **Given** a future endpoint is proposed for remote exposure, **When** it is compared with the contract, **Then** it requires an OpenAPI operation, capability mapping, request bounds, audit classification, and local-equivalence evidence.
2. **Given** a pairing phrase is exchanged successfully, **When** durable access begins, **Then** the phrase is no longer accepted and the issued opaque credential has independent identity, rotation, expiration, and revocation.
3. **Given** an API change is not backward compatible within `/v1`, **When** it is reviewed, **Then** it requires a new path major and an overlap and deprecation plan rather than silent breakage.

### Edge Cases

- A reverse proxy sends forwarding headers from an untrusted address. go-schedule ignores them unless explicit trusted-proxy configuration matches the immediate peer.
- A remote listener is configured without a certificate and key, with a wildcard address but no explicit exposure acknowledgement, or with an insecure TLS minimum. Startup fails before binding.
- A valid credential lacks the capability for the requested operation. The request is denied before any handler mutation and an audit record identifies the actor and outcome without recording the credential.
- A credential is revoked while a live event stream exists. The stream is terminated within the revocation bound and cannot authorize a new request.
- The request outcome becomes uncertain after connection loss. Clients do not automatically replay a state-changing request.
- A client presents a pairing phrase as an everyday bearer credential. Authentication fails without revealing whether the phrase was once valid.
- The OpenAPI document and generated or implemented routes diverge. Contract verification fails before publication.
- A dependency becomes unmaintained or publishes a security fix. The owning implementation slice must upgrade, replace, or record a bounded exception before release qualification.

## Requirements

### Functional Requirements

- **FR-001**: The architecture MUST name HTTPS carrying versioned HTTP/JSON as the sole primary remote transport and MUST reject custom encryption, signed-token formats, and wire protocols.
- **FR-002**: Remote access MUST remain opt-in, fail closed on invalid configuration, and preserve the current Unix-socket or Windows-named-pipe path and offline behavior when disabled.
- **FR-003**: The architecture MUST define one remote adapter that invokes an explicit allowlist of existing daemon operations after boundary controls, without exposing the local handler wholesale.
- **FR-004**: The architecture MUST define Observe, Operate, Manage, and Enroll capabilities with a minimum mapping policy, while reserving concrete actor persistence and per-operation mapping for #167.
- **FR-005**: Every protected remote request MUST require a named actor credential and MUST pass authentication, capability authorization, bounded input, rate limiting, correlation, and audit classification before reaching a daemon operation.
- **FR-006**: The architecture MUST separate short-lived single-use enrollment phrases from opaque durable bearer credentials and MUST define independent issuance, digest-only server retention, expiration, rotation, revocation, and recovery boundaries.
- **FR-007**: The architecture MUST define local IPC, private-network HTTPS, SSH-tunneled HTTPS, reverse-proxied HTTPS, and direct HTTPS modes with explicit product and operator ownership.
- **FR-008**: Network TLS MUST use Go's standard library, TLS 1.3 minimum, operator-provided certificates or an operator-owned TLS reverse proxy, bounded server timeouts, and no product-owned ACME lifecycle in v1.4.
- **FR-009**: Rate limiting MUST use `golang.org/x/time/rate` with separate coarse source and authenticated-actor bounds; rate-limit state MUST remain ephemeral and denial MUST be predictable.
- **FR-010**: The remote contract MUST use OpenAPI 3.1 as its source of truth and a pinned `oapi-codegen` tool to generate strict standard-library HTTP server and client boundaries when #168 implements endpoints.
- **FR-011**: API compatibility MUST use URI major versioning, preserve additive compatibility within a major, reject incompatible clients during capability discovery, and require documented overlap before removing an old major.
- **FR-012**: Live updates MUST use authenticated Server-Sent Events over HTTPS with resumable event identity; state-changing requests MUST never be automatically replayed after an uncertain result.
- **FR-013**: The architecture MUST map concrete threats to bounded controls and acceptance tests covering exposure defaults, transport security, credential guessing and replay, privilege escalation, proxy spoofing, request exhaustion, secret leakage, revocation, wrong-target actions, and version skew.
- **FR-014**: The dependency policy MUST identify ownership, pinning, update review, license and vulnerability checks, and replacement criteria for every non-standard component.
- **FR-015**: The architecture MUST state non-goals including user accounts, SSO, federation, general policy language, custom certificate authority, automatic public exposure, offline mutation queues, and remotely exposed MCP mutation.
- **FR-016**: A repository check and adversarial fixtures MUST fail when required architecture decisions, deployment modes, ownership, or issue sequencing disappear from the maintained architecture page.
- **FR-017**: The architecture MUST be recorded as a dated decision in `CHANGELOG.md` and linked from the existing architecture and API documentation.
- **FR-018**: The feature specification and project item MUST preserve traceability to S074 and issue #165, and #165 MUST remain open until the pull request is merged.

### Key Entities

- **Daemon installation**: Stable remote target identity, display name, supported API majors, capabilities, and operating mode. Its persistence semantics belong to #166.
- **Actor**: A local OS caller or named remote client with assigned capabilities and lifecycle state. Persistence and audit storage belong to #167.
- **Enrollment phrase**: Short-lived, single-use bootstrap secret accepted only by the enrollment endpoint. Its protocol belongs to #169.
- **Durable credential**: High-entropy opaque bearer secret issued uniquely to one actor; only a digest and safe fingerprint are retained by the daemon.
- **Remote operation**: One allowlisted daemon action with OpenAPI operation identity, required capability, request bounds, audit class, and replay behavior.
- **Deployment mode**: A supported arrangement of daemon listener, TLS endpoint, network boundary, and operator infrastructure.
- **Audit event**: Secret-free evidence of actor, daemon, operation, target, result, correlation identity, and time. Its durable model belongs to #167.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One architecture page names exactly one primary remote transport and all five supported deployment modes.
- **SC-002**: Every deployment-mode row assigns listener, TLS, certificate, routing, and monitoring ownership without an unowned security control.
- **SC-003**: Every threat in the bounded threat table maps to at least one control and one acceptance-test class.
- **SC-004**: All four capabilities and both credential types have explicit lifecycle and non-goal boundaries.
- **SC-005**: API versioning, deprecation, live updates, failure behavior, and dependency maintenance can be checked without relying on an implementation that does not yet exist.
- **SC-006**: Repository fixtures demonstrate that removal of a required transport, mode, control, or sequencing gate causes verification to fail.
- **SC-007**: Canonical local verification passes without adding a runtime dependency or opening a network listener.
- **SC-008**: Issues #166 through #173 remain open and their dependency order is unchanged at publication time.

## Assumptions

- The existing HTTP/JSON daemon operations remain authoritative and continue to be served locally over protected IPC.
- Issue #145 is resolved, satisfying the remote release's Unix credential prerequisite.
- The v1.4 milestone targets individually administered installations, not a multi-tenant hosted service.
- Operators choosing network access can provide a routable address and either trusted TLS termination or a certificate and key.
- Concrete endpoint enumeration, actor persistence, pairing protocol details, client UX, and release qualification remain in #166 through #173.

## Dependencies

- Parent: [#18](https://github.com/shruggietech/go-schedule/issues/18).
- Release prerequisite [#145](https://github.com/shruggietech/go-schedule/issues/145) is resolved.
- This feature blocks [#166](https://github.com/shruggietech/go-schedule/issues/166) through [#173](https://github.com/shruggietech/go-schedule/issues/173).

## Clarifications

### Session 2026-09-09

- Existing HTTP/JSON domain behavior is reused, but its local mux is not exported directly to the network. A remote allowlist and boundary middleware are mandatory.
- HTTPS remains required even when a VPN, SSH tunnel, or reverse proxy provides an additional private or encrypted path. This keeps one testable application boundary and avoids mode-specific plaintext exceptions.
- TLS certificates remain operator-owned in v1.4. Direct TLS loads supplied certificate and key files; reverse-proxy mode binds an explicitly selected private address and authenticates forwarded requests itself.
- Durable credentials use standard OAuth 2.0 Bearer header semantics with opaque 256-bit random values. They are not JWTs and carry no embedded authority.
- Server-Sent Events are selected for one-way activity updates; remote mutation remains ordinary request/response HTTP and is never silently replayed.
- #166 is not bundled into S074 because #165 explicitly requires design review before network-related implementation begins.
