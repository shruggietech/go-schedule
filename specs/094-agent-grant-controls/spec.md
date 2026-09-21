# Feature Specification: Agent Grant Controls

**Feature Branch**: `codex/094-agent-grant-controls`

**Created**: 2026-09-20

**Status**: Implemented

**Delivery**: Focused Go, desktop, frontend, Playwright accessibility, native Windows Wails, canonical race, integration, coverage, documentation, and automation gates passed on review branch `codex/094-agent-grant-controls`; hosted review is tracked by the S094 pull request.

**Input**: Work slice S094, complete issue #181 and the remaining acceptance sweep for MCP epic #148.

## Clarifications

### Session 2026-09-20

- Q: Which durations are practical defaults? A: 1 hour, 24 hours, 7 days, 30 days, and an explicit Non-expiring choice.
- Q: Can an existing grant be widened or extended? A: No. The desktop may only narrow authority, impose an expiry, shorten an existing expiry, or revoke. Higher or longer authority requires a new grant.
- Q: How is secret-free enrollment handed off? A: Copy a one-time enrollment bundle through the native clipboard boundary, never render it in the web interface, and cancel the pairing if the copy fails.
- Q: Does this slice make remote listener configuration part of grant management? A: No. Listener lifecycle and permission grants remain separate controls so enabling one transport never changes another.
- Q: What is the bounded recent-audit view? A: The newest 25 shared audit events associated with the selected actor.

## User Scenarios & Testing

### User Story 1 - Understand current agent access (Priority: P1)

An administrator opens Agent Access and can immediately determine whether MCP is off or active, which named agents have access to this daemon, their Observe, Operate, or Manage authority, the transport they use, when access was created and last used, when it expires, and whether it remains active.

**Why this priority**: Safe administration starts with an accurate, secret-free inventory that answers who can do what and where.

**Independent Test**: Seed active, expired, and revoked MCP actors across stdio, localhost HTTP, and remote HTTPS, then verify the workspace presents each non-secret grant fact and the independent state of every transport.

**Acceptance Scenarios**:

1. **Given** no optional listener and no persistent or runtime MCP grant, **When** Agent Access loads, **Then** it clearly reports MCP Off while explaining that stdio remains available only when explicitly launched.
2. **Given** active grants on this daemon, **When** Agent Access loads, **Then** each row names the client, daemon, authority, transport, creation, last use, expiry, and current state without displaying a bearer, pairing phrase, digest, or other protected secret.
3. **Given** Observe, Operate, and Manage grants, **When** they are compared, **Then** Manage has distinct high-impact language and visual treatment while all three remain understandable without protocol terminology.
4. **Given** localhost HTTP, remote HTTPS, or stdio is enabled or used, **When** another transport is inspected, **Then** its state remains independent and is not presented as implicitly enabled.

---

### User Story 2 - Grant and constrain remote MCP access (Priority: P1)

An administrator creates a named remote MCP grant with the minimum required authority and a practical duration, or deliberately chooses non-expiring access, then transfers one short-lived enrollment bundle without exposing a durable credential in the interface.

**Why this priority**: The existing command-line enrollment workflow is secure but does not provide the simple desktop control promised by the MCP epic.

**Independent Test**: Create grants for each authority and duration choice, exchange one enrollment, and verify the resulting actor retains the selected expiry while the desktop never renders or retains the pairing phrase.

**Acceptance Scenarios**:

1. **Given** the grant form, **When** an administrator creates access, **Then** client name, target daemon, remote HTTPS transport, authority, and duration are explicit before submission.
2. **Given** a finite duration, **When** the pairing is exchanged, **Then** the resulting actor expires at the originally selected grant deadline rather than at the enrollment deadline.
3. **Given** non-expiring access, **When** it is selected, **Then** the choice requires a deliberate control and the resulting grant is labeled Non-expiring.
4. **Given** successful grant creation, **When** the one-time enrollment bundle is produced, **Then** it is copied through the native clipboard boundary and the interface displays only safe metadata.

---

### User Story 3 - Narrow, expire, revoke, and audit grants (Priority: P1)

An administrator can reduce an active agent's authority, impose or shorten an expiry, revoke it, and inspect recent actions from the shared audit history. Existing connections lose superseded authority on their next request.

**Why this priority**: A grant control is incomplete unless policy changes take effect for already connected clients and leave understandable evidence.

**Independent Test**: Keep stdio, localhost, remote JSON, and remote MCP connections open while narrowing, expiring, and revoking their actors; assert the next operation is denied as appropriate and the recent-action view identifies actor, daemon, action, target, result, and time.

**Acceptance Scenarios**:

1. **Given** a Manage or Operate grant, **When** it is narrowed, **Then** the next request on an existing connection is evaluated against the lower authority.
2. **Given** active non-expiring access, **When** a future expiry is imposed, **Then** access ends at that deadline and the state becomes Expired without requiring a restart.
3. **Given** an active grant, **When** it is revoked, **Then** the next request on every related live connection is denied and the grant cannot be reactivated.
4. **Given** an agent with audit history, **When** its recent actions are opened, **Then** the view presents bounded, secret-free shared audit records associated with that actor.
5. **Given** keyboard-only or assistive-technology use, **When** grant dialogs, authority controls, destructive actions, and audit disclosure are operated, **Then** labels, focus order, focus return, state announcements, and confirmation semantics remain usable.

### Edge Cases

- A pairing can expire before exchange; it creates no actor and remains visibly separate from active grants.
- A grant deadline can pass while the Agent Access page is open; refresh and subsequent authorization both report it as expired.
- An actor may have no durable credential because it represents a runtime stdio or localhost session; transport inference remains explicit and does not invent a credential.
- An Observe-only localhost listener has no mutation actor; it is represented by its safe listener metadata rather than a fabricated actor grant.
- A credential may remain stored after its actor is revoked; the actor state remains authoritative and the credential cannot authenticate successfully.
- A wrong-daemon request remains denied before any action even when the same client name and authority exist on two daemons.
- An administrator cannot widen an existing grant from the desktop control; higher authority requires a new deliberate grant.

## Requirements

### Functional Requirements

- **FR-001**: Agent Access MUST present an MCP summary and independent stdio, localhost HTTP, and remote HTTPS transport states for the selected local daemon.
- **FR-002**: The grant inventory MUST include only MCP actors and MUST present client name, target daemon name and identifier, authority, transport, creation time, last-use time when known, expiry or Non-expiring, and Active, Expired, or Revoked state.
- **FR-003**: The inventory MUST correlate actors with safe credential, runtime listener, and audit metadata without returning or rendering credential values, pairing phrases, verifiers, digests, or protected configuration.
- **FR-004**: Observe, Operate, and Manage MUST use distinct labels and plain-language descriptions, and Manage MUST explicitly identify definition-changing authority.
- **FR-005**: Creating a remote MCP grant MUST require a client name, one of Observe, Operate, or Manage, and one of 1 hour, 24 hours, 7 days, 30 days, or a deliberate Non-expiring duration.
- **FR-006**: A finite grant deadline MUST be recorded when the pairing request is created and MUST become the persistent MCP actor expiry if that pairing is exchanged.
- **FR-007**: Successful grant creation MUST copy a one-time enrollment bundle through the native clipboard boundary, clear the phrase from service memory as soon as the copy completes, and return only secret-free metadata to the web interface.
- **FR-008**: Failed native clipboard transfer MUST cancel the unexchanged pairing and report that no usable grant was issued.
- **FR-009**: Administrators MUST be able to narrow Manage to Operate or Observe, narrow Operate to Observe, impose an expiry on non-expiring access, shorten a finite expiry, and revoke any non-built-in MCP actor.
- **FR-010**: The desktop control MUST NOT widen an existing actor capability, remove an existing finite expiry, extend a finite expiry, reactivate a revoked actor, or modify the protected built-in actor.
- **FR-011**: Actor capability, expiry, and revocation changes MUST be loaded by shared authorization on the next request so existing stdio, localhost HTTP, remote JSON, and remote MCP connections cannot retain superseded access.
- **FR-012**: Agent Access MUST show up to 25 newest shared audit events for a selected actor, including daemon, operation, target type and identifier, result, and occurrence time.
- **FR-013**: Audit presentation MUST use bounded fields and MUST never include request bodies, task inputs, environment values, credentials, pairing phrases, tokens, secrets, or raw internal error causes.
- **FR-014**: Listener lifecycle controls MUST remain separate from actor grants; enabling or disabling one transport MUST NOT enable, disable, or imply enablement of another transport.
- **FR-015**: Grant creation, narrowing, expiry, revocation, audit inspection, and confirmation dialogs MUST be operable by keyboard, expose programmatic names and states, preserve focus, and announce asynchronous outcomes.
- **FR-016**: Wrong-target, expired, and revoked requests MUST fail closed and retain actor-attributed denial evidence where the shared authorization boundary can identify the actor.
- **FR-017**: S094 MUST retain the completed MCP epic guarantees for default-off network access, Observe, Operate, and Manage separation, stdio without a listening port, authenticated HTTP transports, client-attributed audit, redaction, hostile-content isolation, and supported protocol compatibility.

### Key Entities

- **Agent grant**: Secret-free administrative projection of one MCP actor, its target daemon, authority, inferred transport, lifecycle timestamps, state, and recent activity.
- **Enrollment request**: Ten-minute one-time request containing the intended MCP client name, authority, and optional persistent grant deadline.
- **Client credential metadata**: Safe identifier, fingerprint, lifecycle timestamps, and last-use evidence associated with a persistent actor; no bearer value is included.
- **Transport state**: Independent availability and active-state projection for stdio, localhost HTTP, and remote HTTPS.
- **Recent agent action**: Bounded shared audit record correlated to an MCP actor and target daemon.

## Success Criteria

### Measurable Outcomes

- **SC-001**: For every seeded grant, a user can answer who, which daemon, which authority, which transport, and until when from one Agent Access view without opening another page.
- **SC-002**: All five duration choices produce the expected expiry semantics, and no protected secret appears in rendered text, logs, snapshots, bridge results, or audit output.
- **SC-003**: Capability narrowing, imposed or shortened expiry, and revocation deny a now-disallowed request on the next attempt across every applicable live transport.
- **SC-004**: Observe, Operate, and Manage are distinguishable by accessible name, explanatory copy, and non-color-only visual treatment at supported window sizes and themes.
- **SC-005**: Keyboard and assistive-technology checks complete grant creation, safe modification, destructive confirmation, and recent-audit inspection without a pointer.
- **SC-006**: The issue #181 matrix and the parent #148 acceptance matrix pass together with the repository's complete unit, integration, race, frontend, conformance, formatting, and canonical verification gates.

## Assumptions

- Administrative grant controls operate against This computer because enrollment authority is intentionally local; remote clients can observe safe actor and audit summaries but cannot administer another daemon's grants.
- The existing daemon identity, actor, credential, pairing, authorization, and audit models remain authoritative.
- Remote HTTPS is the only transport that needs a durable pairing grant. Stdio is launched explicitly by its host, and localhost HTTP retains its separate listener setup and one-time credential copy flow.
- Recent means the newest 25 actor-associated audit events, while the existing retention policy remains unchanged.
- A non-expiring grant is permitted only as an explicit initial selection. Existing finite grants can be shortened or revoked, not extended or converted to non-expiring.
