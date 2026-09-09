# Feature Specification: Agent Access Controls and MCP Release Gates

**Feature Branch**: `codex/071-agent-access-controls`

**Created**: 2026-09-09

**Status**: Implemented

**Delivery**: Agent Access controls, runtime named-client evidence, host guidance, conformance, package-shaped smoke, accessibility coverage, and canonical verification passed 2026-09-09 on review branch `codex/071-agent-access-controls` for [#164](https://github.com/shruggietech/go-schedule/issues/164).

**Input**: GitHub issue [#164](https://github.com/shruggietech/go-schedule/issues/164), complete the observe-only local MCP release with understandable desktop controls, supported-client guidance, and conformance evidence.

## User Scenarios & Testing

### User Story 1 - Understand and control local agent access (Priority: P1)

A desktop user can open one Agent Access workspace, understand that stdio is available only when an MCP host launches it, see whether the optional localhost HTTP listener is off or active, identify the active named localhost client, and revoke that client by turning the listener off.

**Why this priority**: Users need one truthful place to determine whether agent access is active and to close the only optional network listener.

**Independent Test**: Open Agent Access against disabled and enabled daemon fixtures, inspect the stdio and HTTP explanations, enable one named localhost client, observe its safe status, record authenticated use, rotate access, and revoke it without exposing a credential after its one-time copy.

**Acceptance Scenarios**:

1. **Given** a connected daemon with localhost HTTP disabled, **When** Agent Access opens, **Then** it reports no active network listener and describes stdio as an on-demand subprocess rather than an active connection.
2. **Given** valid client, port, and optional browser-origin input, **When** localhost access is enabled, **Then** one runtime-only named client becomes active and its newly issued credential is copied once without being retained in the workspace response.
3. **Given** an active named client, **When** authenticated requests succeed, **Then** the workspace shows bounded runtime-only access evidence consisting of the last successful access time and request count.
4. **Given** an active named client, **When** access is rotated, **Then** the old credential is invalidated, the replacement is copied once, and the client name and listener policy remain unchanged.
5. **Given** an active named client, **When** the user revokes it, **Then** the credential is invalidated, the listener closes, and Agent Access returns to Off.
6. **Given** credential copying fails during enable or rotation, **When** the operation completes, **Then** localhost HTTP is disabled so no inaccessible credential remains active.

---

### User Story 2 - Understand the Observe boundary and configure a supported host (Priority: P1)

A user can follow concise Codex or generic MCP-host instructions and understand the exact difference between stdio process access, HTTP bearer authorization, Observe authority, and unavailable future Operate and Manage authority.

**Why this priority**: Setup guidance is part of the supported product contract, and ambiguous authority or credential language can cause unsafe configuration.

**Independent Test**: Follow both clean-install setup paths using packaged command layouts, initialize an official client, discover the five Observe resources, and verify the guidance never claims that stdio opens a listener or that Operate, Manage, remote access, or durable grants are available.

**Acceptance Scenarios**:

1. **Given** a clean installation and an accessible daemon, **When** the Codex stdio instructions are followed, **Then** Codex can initialize the installed `gosched mcp serve` command without a network listener or credential.
2. **Given** a generic host that supports command-based MCP, **When** the generic stdio instructions are followed, **Then** the host can initialize the same Observe surface through the installed executable.
3. **Given** a host that requires Streamable HTTP, **When** localhost access is explicitly enabled, **Then** the guidance requires the endpoint plus its separate one-time bearer credential and explains rotation and revocation.
4. **Given** any Agent Access surface or documentation, **When** authority is described, **Then** Observe is the only available class and Operate and Manage are visibly unavailable future capabilities.

---

### User Story 3 - Qualify the complete local MCP contract (Priority: P2)

A maintainer can run deterministic conformance and packaging gates that prove every exposed resource remains bounded, redacted, compatible, data-only, and available from packaged command layouts.

**Why this priority**: The release should not depend on manual confidence that earlier slices remain compatible after desktop and packaging integration.

**Independent Test**: Run the repository gates on each supported platform and exercise official protocol negotiation, every resource and continuation template, schema shape, safe errors, hostile content, zero tools, built-command stdio startup, and desktop Agent Access behavior.

**Acceptance Scenarios**:

1. **Given** every approved resource and continuation template, **When** conformance runs for both supported protocol revisions, **Then** schema, pagination, bounds, trust labeling, redaction, error vocabulary, and zero-tool authority pass.
2. **Given** task-controlled strings that resemble instructions or protocol content, **When** resources are read, **Then** they remain bounded JSON data and cannot change discovery, permissions, or tool availability.
3. **Given** a freshly built CLI in a package-shaped directory, **When** an official client launches `mcp serve`, **Then** initialization and discovery succeed and process shutdown is clean on Windows, macOS, and Linux.
4. **Given** the desktop at supported viewport and zoom combinations, **When** Agent Access is operated by keyboard, **Then** serious and critical accessibility violations are absent and no horizontal overflow occurs.

### Edge Cases

- The desktop loads while the daemon is unavailable, reconnects, or becomes unavailable during a lifecycle action.
- Client names are blank, whitespace-only, duplicated across restarts, over the limit, invalid UTF-8, or contain control characters and instruction-like text.
- A one-time credential is issued but clipboard access fails, the frontend bridge rejects, or the user navigates away immediately.
- Status, successful authorization, rotation, disablement, and daemon shutdown overlap.
- Recent access evidence has no successful request yet, reaches its counter limit, or is read concurrently with authentication.
- Browser origins are empty, repeated, invalid, or use explicit default ports.
- A built executable path contains spaces or non-ASCII characters.
- A malicious resource field contains Markdown, XML-like instructions, ANSI escapes, or JSON fragments that name tools or permissions.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST provide one Agent Access workspace reachable from primary navigation.
- **FR-002**: The workspace MUST present stdio as available on demand through `gosched mcp serve`, MUST state that it opens no listener, and MUST NOT describe stdio as enabled, connected, revocable, or bearer-authenticated.
- **FR-003**: The workspace MUST report localhost HTTP as Off or Active from authoritative daemon status and MUST display the exact endpoint, active client name, allowed origins, credential fingerprint, activation time, last successful access time, and successful request count when available.
- **FR-004**: The desktop MUST let a connected local user enable localhost HTTP with a client name, valid port, and zero or more explicit browser origins; rotate its credential; and revoke the client by disabling the listener.
- **FR-005**: One localhost HTTP listener MUST represent at most one active named client and one current credential. Durable or simultaneous multi-client grants remain outside S071 and assigned to issue #181.
- **FR-006**: Client names MUST be trimmed, valid UTF-8, between 1 and 64 bytes, and free of control characters. Existing CLI and API callers that omit a name MUST receive the stable compatibility name `Local MCP client`.
- **FR-007**: Client names and access evidence MUST remain runtime-only and MUST be cleared on disablement, daemon shutdown, or restart with the endpoint, origins, and credential digest.
- **FR-008**: A successful authenticated non-preflight MCP request MUST atomically update a saturating successful-request count and last-successful-access timestamp before SDK dispatch; rejected requests and CORS preflights MUST NOT update evidence.
- **FR-009**: Status MUST expose only non-secret data. It MUST NOT expose plaintext credentials, authorization headers, request bodies, resource URIs, task content, network peer metadata, or failed credential attempts.
- **FR-010**: Desktop enable and rotation MUST copy the newly issued credential through the native clipboard boundary and MUST NOT return it to the frontend workspace. If copying fails, the desktop service MUST disable the endpoint before returning a bounded failure.
- **FR-011**: Desktop actions MUST be serialized, duplicate-safe, deadline-bounded, and return actionable states without raw daemon, clipboard, endpoint-internal, or credential details.
- **FR-012**: The workspace MUST explain Observe in plain language and MUST show Operate and Manage as unavailable future authority without controls that imply activation.
- **FR-013**: Documentation MUST provide Codex and generic command-based stdio setup, generic Streamable HTTP setup, clean removal, troubleshooting, privacy, hostile-content, rotation, and revocation guidance.
- **FR-014**: Documentation MUST distinguish process-launched stdio with inherited local IPC access from separately authorized localhost HTTP, and MUST NOT advertise remote MCP, durable grants, Operate, or Manage as shipped.
- **FR-015**: Conformance tests MUST cover both supported protocol revisions, five static resources, four continuation templates, schema version, media type, pagination, text and output bounds, safe error vocabulary, structural secret exclusion, trust metadata, and zero advertised or callable tools.
- **FR-016**: Hostile task-controlled strings MUST remain JSON data and MUST NOT alter resource discovery, permission class, server instructions, or tool availability.
- **FR-017**: A packaged-command smoke test MUST build the real CLI into a package-shaped path, launch `gosched mcp serve` with the official SDK, initialize and discover the Observe surface, verify zero tools, and shut down cleanly on Windows, macOS, and Linux with hidden Windows child-process creation.
- **FR-018**: Desktop tests MUST cover disconnected, disabled, enabled, no-activity, active-use, validation, clipboard-failure rollback, rotation, revocation, stale result, accessibility, keyboard, and supported zoom states.
- **FR-019**: Existing CLI status, enable, rotate, disable, and JSON contracts MUST remain backward compatible apart from additive non-secret status fields and the optional client-name input.
- **FR-020**: S071 MUST NOT persist agent grants, add remote listeners, add scheduler mutation tools, collect request content, or change the S069 Observe resource payload contract.
- **FR-021**: GitHub issue [#164](https://github.com/shruggietech/go-schedule/issues/164) MUST remain traceable through the specification, tasks, changelog, pull request, and verification record.

### Key Entities

- **Agent Access Workspace**: Desktop projection of stdio availability, localhost HTTP status, Observe authority, future authority, setup guidance, and safe lifecycle actions.
- **Named Localhost Client**: The one runtime-only display identity attached to the active localhost listener and current credential.
- **Access Evidence**: Runtime-only last-successful-access time and saturating successful-request count for the active credential generation.
- **One-Time Credential Handoff**: Backend-only result that copies a newly issued credential to the native clipboard and returns only non-secret status to the frontend.
- **Conformance Matrix**: Deterministic coverage of revisions, resources, templates, schemas, trust, redaction, bounds, errors, authority, subprocess lifecycle, and package-shaped execution.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A user can determine within one workspace whether any MCP network listener is active and revoke an active localhost client in no more than two actions.
- **SC-002**: One hundred percent of desktop and API status responses contain no plaintext credential, and every clipboard-failure test leaves localhost HTTP disabled.
- **SC-003**: One hundred percent of successful authenticated non-preflight requests update bounded runtime evidence, while rejected requests and preflights produce zero evidence updates.
- **SC-004**: Official clients negotiate both supported revisions, discover exactly five resources and four templates, and discover and call zero tools across stdio, localhost HTTP, and packaged-command conformance tests.
- **SC-005**: Secret and hostile-content canaries produce zero protected-value matches and zero changes to discovery or authority in every decoded MCP response.
- **SC-006**: The Agent Access desktop passes component tests plus automated WCAG 2.2 AA serious/critical checks and horizontal-overflow checks at 80, 100, 150, and 200 percent zoom in a 900 by 650 viewport.
- **SC-007**: Focused Go race tests, frontend tests and build, packaged-command smoke tests, and canonical `scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-09

- Q: Does naming introduce durable multi-client grants? A: No. S071 names the one runtime-only S070 credential; simultaneous durable grants remain assigned to #181.
- Q: What does revoke mean for the single localhost client? A: Revoke disables the listener and invalidates its credential; rotation replaces the credential while preserving the client name and endpoint policy.
- Q: What recent evidence is retained? A: Only runtime last-successful-access time and a saturating successful-request count; no request content or failed-attempt history is collected.
- Q: How does the desktop receive a new credential? A: The Go desktop service copies it through the native clipboard and returns only non-secret status; copy failure disables the endpoint.
- Q: Is stdio ever shown as active? A: No. It is described only as an on-demand command that opens no listener and inherits local IPC authorization.

## Assumptions

- The S069 Observe contracts and S070 localhost transport remain authoritative and are extended only with non-secret runtime client and access metadata.
- The existing Wails application, protected local IPC client, native clipboard wrapper, and repository browser-test harness are available.
- Users who choose HTTP can provide an available port and copy the one-time clipboard value into their client before overwriting it.
- Package-shaped command smoke proves installed command layout and process behavior; formal signed or release-candidate artifact qualification remains assigned to issue #190.
- Remote MCP, durable named grants, expiration, per-action audit, Operate, Manage, and multiple simultaneous clients remain outside S071.
