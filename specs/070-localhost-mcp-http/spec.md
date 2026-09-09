# Feature Specification: Authenticated Localhost MCP

**Feature Branch**: `codex/070-localhost-mcp-http`

**Created**: 2026-09-08

**Status**: Implemented

**Delivery**: Runtime-only numeric-loopback MCP lifecycle, protected IPC and CLI controls, one-time digest-backed credentials, exact Host and Origin enforcement, official stateless Streamable HTTP Observe parity, adversarial security tests, focused race verification, and canonical eight-gate verification passed 2026-09-08 on review branch `codex/070-localhost-mcp-http` for [#163](https://github.com/shruggietech/go-schedule/issues/163).

**Input**: GitHub issue [#163](https://github.com/shruggietech/go-schedule/issues/163), add an optional authenticated loopback-only HTTP transport for the existing Observe surface.

## User Scenarios & Testing

### User Story 1 - Enable a temporary localhost endpoint (Priority: P1)

A local operator can explicitly start an MCP Streamable HTTP endpoint through the protected daemon API, receive a newly generated client credential once, inspect its non-secret status, and turn it off without affecting ordinary CLI, GUI, or stdio access.

**Why this priority**: A fail-closed lifecycle and separately provisioned credential are prerequisites for exposing any MCP traffic over TCP.

**Independent Test**: Start a daemon with no listener, enable a chosen free port through IPC, initialize MCP using the returned credential, disable it, and prove the port closes while local IPC and stdio remain usable.

**Acceptance Scenarios**:

1. **Given** a freshly started daemon, **When** endpoint status is requested, **Then** it reports disabled and no TCP listener exists.
2. **Given** a valid loopback port and origin policy, **When** the operator enables the endpoint, **Then** the daemon binds only numeric loopback, returns the endpoint and a cryptographically random bearer credential once, and reports enabled status without revealing the credential later.
3. **Given** an enabled endpoint, **When** it is disabled, **Then** the listener shuts down promptly, the credential is revoked, and local IPC, GUI, CLI, and stdio behavior is unchanged.
4. **Given** an endpoint that was enabled before daemon shutdown, **When** the daemon restarts, **Then** the endpoint remains disabled and the old credential is stale.
5. **Given** a port already in use or an invalid port or origin, **When** enablement is requested, **Then** it fails atomically with an actionable bounded error and leaves endpoint status disabled.

---

### User Story 2 - Read the Observe surface through authenticated HTTP (Priority: P1)

A local MCP client that cannot launch stdio can connect to the temporary endpoint using its bearer credential and discover and read the same bounded Observe resources delivered by S069.

**Why this priority**: Transport parity is the outcome of the slice, and reusing the S069 server prevents HTTP from becoming a separate authority or data contract.

**Independent Test**: Connect official SDK clients over stdio and HTTP to the same daemon fixture, negotiate each supported revision, discover resources and templates, read every resource, and compare decoded results while confirming that neither transport advertises tools.

**Acceptance Scenarios**:

1. **Given** a valid credential and acceptable request context, **When** a supported MCP client initializes over Streamable HTTP, **Then** official SDK negotiation succeeds and advertises the same resources, templates, instructions, and zero tools as stdio.
2. **Given** the same daemon state, **When** each Observe resource is read through stdio and HTTP, **Then** the safe payload contracts match apart from expected generation timestamps.
3. **Given** a canceled request or daemon shutdown, **When** HTTP work is in flight, **Then** cancellation propagates, sessions terminate, and listener shutdown completes within the bounded lifecycle deadline.
4. **Given** an unsupported protocol revision, malformed body, wrong media headers, or oversized request, **When** the endpoint receives it, **Then** it returns a bounded protocol or HTTP error without reaching scheduler mutation code.

---

### User Story 3 - Reject local web and rebinding attacks (Priority: P1)

The localhost endpoint rejects requests that lack the current credential, present a non-loopback Host, or carry a browser Origin outside the explicit allowlist, even though the TCP peer is local.

**Why this priority**: Browser requests and DNS rebinding can turn a loopback listener into ambient authority unless every layer fails closed.

**Independent Test**: Send direct and browser-shaped requests across a matrix of current, missing, malformed, rotated, revoked, and stale credentials; allowed and hostile origins; numeric loopback, localhost, arbitrary, forwarded, and malformed Host values; and IPv4 and IPv6 connection attempts.

**Acceptance Scenarios**:

1. **Given** no bearer credential, a malformed authorization header, or a non-current credential, **When** any request reaches the endpoint, **Then** it is rejected before MCP parsing with no credential oracle beyond unauthorized status.
2. **Given** a request whose Host is not the exact numeric loopback address and active port, **When** it reaches the endpoint, **Then** it is rejected before authorization or MCP parsing.
3. **Given** a request with an Origin header, **When** the origin is not an exact member of the enablement allowlist, **Then** it is rejected without permissive CORS response headers.
4. **Given** a non-browser MCP client with no Origin header, **When** Host and authorization are valid, **Then** the request is accepted.
5. **Given** credential rotation, disablement, or daemon restart, **When** the prior credential is retried, **Then** it is rejected and local IPC authorization remains unchanged.

### Edge Cases

- The requested port is zero, outside the valid range, privileged, already occupied, released between validation and bind, or races with another enable request.
- Enable, rotate, disable, status, and shutdown overlap while MCP requests are active.
- The request uses `localhost`, an IPv4-mapped IPv6 address, an alternate textual loopback form, user information, trailing dots, multiple Host values, absolute-form targets, or forwarded host headers.
- Origin is absent, opaque (`null`), malformed, case-varied, contains credentials, uses a disallowed scheme, differs only by port, or has a path, query, fragment, or trailing dot.
- Authorization contains another scheme, duplicate headers, leading or trailing whitespace, invalid UTF-8 bytes, a current token plus suffix, or a token from before rotation.
- A request body is chunked or claims a misleading Content-Length and exceeds the configured maximum.
- Shutdown begins during initialization, resource discovery, resource read, or credential rotation.

## Requirements

### Functional Requirements

- **FR-001**: The daemon MUST expose runtime-only localhost MCP lifecycle operations through its existing protected IPC API: status, enable, credential rotation, and disable.
- **FR-002**: The endpoint MUST be disabled on every daemon start and MUST NOT persist enablement, port, allowed origins, or plaintext credentials to configuration, storage, logs, process arguments, or environment variables.
- **FR-003**: Enablement MUST bind only `127.0.0.1` on an operator-selected port from 1 through 65535; callers MUST NOT supply a host or network interface.
- **FR-004**: The implementation MUST reject port conflicts and invalid configuration atomically and MUST leave the prior lifecycle state unchanged on failure.
- **FR-005**: A successful enable or rotation MUST generate a new credential using a cryptographically secure random source, return it only in that operation's protected IPC response, and expose only a non-secret fingerprint in status.
- **FR-006**: Credential comparison MUST avoid content-dependent early exit, and missing, malformed, duplicated, stale, rotated, and revoked credentials MUST all fail closed before MCP parsing.
- **FR-007**: Every MCP HTTP request MUST require exactly one `Authorization: Bearer <credential>` header, and responses MUST prevent credential caching without echoing credentials.
- **FR-008**: Every request MUST have a Host equal to the active endpoint's exact numeric IPv4 loopback address and port; hostname, public, ambiguous, malformed, alternate-loopback, and forwarded values MUST NOT influence acceptance.
- **FR-009**: Requests without an Origin header MAY proceed after other checks; requests with Origin MUST match one exact normalized origin explicitly supplied during enablement.
- **FR-010**: Allowed origins MUST use `http` or `https`, contain a host and explicit valid port, and contain no credentials, path other than `/`, query, fragment, wildcard, opaque value, or non-loopback host.
- **FR-011**: The HTTP transport MUST use the official Go SDK's stateless Streamable HTTP handler, retain SDK localhost protection, propagate request cancellation, and bound request bodies to 1 MiB.
- **FR-012**: The HTTP transport MUST instantiate the same `mcpobserve.NewServer` registration used by stdio and MUST advertise only the S069 Observe resources and templates with zero tools.
- **FR-013**: HTTP and stdio MUST produce equivalent safe resource contracts for the same daemon state; neither transport may add scheduler authority, direct database access, or alternate data mapping.
- **FR-014**: Enable while enabled MUST fail with a conflict rather than silently replacing a listener; rotation MUST preserve the endpoint and allowed origins while immediately invalidating the prior credential.
- **FR-015**: Disable MUST be idempotent, revoke the active credential before or concurrently with listener shutdown, reject new requests, and wait no longer than five seconds for graceful shutdown before forcing closure.
- **FR-016**: Daemon shutdown MUST revoke the credential, stop accepting connections, cancel in-flight HTTP work, and return no lifecycle error solely because the endpoint was never enabled.
- **FR-017**: Localhost HTTP lifecycle operations MUST NOT change daemon IPC access, stdio MCP availability, GUI behavior, remote JSON settings, future remote MCP settings, scheduler data, or notification credentials.
- **FR-018**: CLI commands MUST provide `gosched mcp http status`, `enable --port <port> [--origin <origin>...]`, `rotate`, and `disable`, support machine-readable output, print newly issued credentials only for enable and rotate, and avoid placing credentials in command arguments.
- **FR-019**: Public responses, CLI diagnostics, and logs MUST use bounded actionable messages and MUST exclude credentials, authorization headers, protected daemon values, and internal implementation details.
- **FR-020**: Automated tests MUST cover default-off behavior, lifecycle races, restart staleness, port conflict, SDK compatibility, stdio parity, cancellation, graceful shutdown, DNS rebinding Host matrices, exact Origin matrices, authorization matrices, request-size bounds, and absence of authority changes.
- **FR-021**: GitHub issue [#163](https://github.com/shruggietech/go-schedule/issues/163) MUST remain traceable through the specification, tasks, changelog, pull request, and verification record.

### Key Entities

- **Local MCP HTTP Manager**: Daemon-owned concurrency-safe runtime component controlling the single optional listener, its HTTP server, active credential verifier, origin policy, and shutdown lifecycle.
- **Endpoint Status**: Non-secret state containing enabled status, exact numeric loopback endpoint, allowed origins, activation time, and credential fingerprint.
- **Client Credential**: High-entropy bearer secret issued once over protected IPC and held only in daemon memory while current.
- **Origin Policy**: Canonical exact-match set of explicitly approved loopback browser origins.
- **Authenticated MCP Handler**: Ordered request boundary enforcing Host, Origin, authorization, cache, and size policy before delegating to the official SDK handler.
- **Observe Server Factory**: Existing S069 construction path that registers the allowlisted resource surface for either stdio or Streamable HTTP.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A fresh or restarted daemon has no TCP listener attributable to MCP, and disablement closes an enabled endpoint within five seconds.
- **SC-002**: One hundred percent of missing, malformed, stale, rotated, revoked, duplicated, or incorrect authorization cases are rejected before MCP dispatch in the automated security matrix.
- **SC-003**: One hundred percent of non-exact Host and Origin cases in the DNS-rebinding and browser-request matrix are rejected, while exact configured requests and originless native-client requests succeed.
- **SC-004**: Official SDK clients initialize and read every approved Observe resource through HTTP for both supported protocol revisions, and parity tests find no schema or authority difference from stdio.
- **SC-005**: Port conflict, invalid configuration, repeated enable, rotation, disable, cancellation, and daemon shutdown tests complete without leaked listeners, goroutines, credentials, or altered local IPC access.
- **SC-006**: Focused security and lifecycle tests, race verification, and canonical `scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-08

- Q: Does enablement survive daemon restart? A: No. Listener state, origins, and credentials are runtime-only so every daemon start fails closed and every pre-restart credential becomes stale.
- Q: Which address is exposed? A: The daemon binds only `127.0.0.1`; the operator supplies a port but never a host, interface, wildcard, or IPv6 address.
- Q: How are browser origins handled? A: Originless native clients are accepted after Host and authorization checks; browser-shaped requests require an exact explicitly configured loopback HTTP or HTTPS origin with an explicit port.
- Q: Is the HTTP credential shared with daemon IPC, stdio, remote JSON, or future remote MCP? A: No. It is a separately generated in-memory bearer credential that gates only this listener.
- Q: Which HTTP mode is used? A: The official SDK's stateless Streamable HTTP handler is used with request cancellation and a 1 MiB body limit; stateful sessions and legacy SSE are outside S070.
- Q: What happens when enable is called while active? A: It returns conflict. Rotation changes only the credential; changing port or origins requires disable followed by enable.

## Assumptions

- Operators can choose an available nonzero port and copy a one-time credential into a local MCP client configuration.
- The existing protected IPC API is the authoritative control channel for listener lifecycle and secret issuance.
- Numeric IPv4 loopback is sufficient for the first HTTP delivery and is clearer to secure and document than a dual-stack or hostname listener.
- S069 resource bounds, redaction, trust labeling, protocol compatibility, and zero-tool authority remain authoritative and are reused without modification unless parity testing finds a defect.
- Persistent named clients, recent access evidence, desktop controls, packaging conformance, and broader guidance remain in issue #164.
- Remote HTTP, OAuth, TLS termination, Operate and Manage tools, durable grants, expiry, and audit history remain outside S070.
