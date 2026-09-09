# Feature Specification: Local Observe-Only MCP

**Feature Branch**: `codex/069-local-mcp-observe`

**Created**: 2026-09-08

**Status**: Implemented

**Delivery**: Observe-only MCP contract, official SDK stdio adapter, five bounded resources, structural secret exclusion, two-revision protocol coverage, cross-platform subprocess lifecycle tests, Codex setup documentation, focused race verification, and canonical eight-gate verification passed 2026-09-08 on review branch `codex/069-local-mcp-observe` for [#161](https://github.com/shruggietech/go-schedule/issues/161) and [#162](https://github.com/shruggietech/go-schedule/issues/162).

**Input**: GitHub issues [#161](https://github.com/shruggietech/go-schedule/issues/161) and [#162](https://github.com/shruggietech/go-schedule/issues/162), define the MCP security boundary and deliver its first local observe-only transport.

## User Scenarios & Testing

### User Story 1 - Connect a local MCP host safely (Priority: P1)

A local automation host can launch `gosched mcp serve`, negotiate a supported MCP revision through the official Go SDK, and discover the scheduler's observe-only resource surface over stdio without opening a network listener.

**Why this priority**: A clean, local, interoperable protocol boundary is the prerequisite for every useful observation and establishes the security posture for future MCP work.

**Independent Test**: Launch the command as a subprocess, complete MCP initialization, list resources and templates, read health, close stdin, and verify clean termination with protocol traffic only on stdout.

**Acceptance Scenarios**:

1. **Given** a supported MCP host, **When** it launches `gosched mcp serve`, **Then** protocol initialization succeeds over stdio and advertises resources without advertising mutation tools.
2. **Given** an unsupported protocol revision, **When** initialization is attempted, **Then** the host receives a standards-compliant bounded error without a panic or sensitive diagnostic.
3. **Given** an active stdio session, **When** the host disconnects or cancels, **Then** the process terminates promptly and closes daemon IPC activity.
4. **Given** any MCP session, **When** the process is inspected, **Then** it has opened no TCP listener and writes diagnostics only to stderr.

---

### User Story 2 - Inspect scheduler state without secret disclosure (Priority: P1)

A local MCP client can inspect daemon health, active tasks, upcoming schedules, recent alerts, and recent run history through bounded JSON resources whose potentially user-controlled fields are explicitly marked untrusted.

**Why this priority**: These observations provide immediate operational value while avoiding the authority and confirmation complexity of scheduler mutation.

**Independent Test**: Seed tasks, schedules, runs, output, and alerts containing secrets and hostile text, read every approved resource, and prove that safe summaries remain useful while protected fields never cross the MCP boundary.

**Acceptance Scenarios**:

1. **Given** a healthy daemon, **When** health is read, **Then** the response identifies status and version without configuration or environment disclosure.
2. **Given** active tasks and schedules, **When** their resources are read, **Then** stable identity, operational state, safe schedule summaries, readiness, and upcoming times are returned without commands, arguments, environment, stdin, working directories, run-as identities, or raw schedule definitions.
3. **Given** alerts, task names, descriptions, or command output containing instructions or markup, **When** a resource is read, **Then** those fields remain inert JSON strings and are identified as untrusted content.
4. **Given** more records than one response permits, **When** a collection resource is read, **Then** it returns a deterministic bounded page and an opaque continuation URI until the collection is exhausted.
5. **Given** long command output or alert text, **When** it is returned, **Then** it is safely truncated with explicit truncation metadata and the complete payload is never exposed.

---

### User Story 3 - Preserve the daemon's local access decision (Priority: P2)

A local MCP client receives safe actionable failures when the daemon is unavailable or denies IPC access, and the adapter never weakens, bypasses, or replaces the existing operating-system access boundary.

**Why this priority**: MCP must remain a thin adapter over the daemon, not a second authority or business-logic implementation.

**Independent Test**: Exercise absent, denied, timed-out, malformed, and canceled IPC outcomes and verify bounded MCP errors, no fallback transport, and no leaked endpoint paths or credentials.

**Acceptance Scenarios**:

1. **Given** daemon IPC access is denied, **When** any resource is read, **Then** the MCP request is denied with safe recovery guidance and no alternate connection is attempted.
2. **Given** the daemon is absent or times out, **When** a resource is read, **Then** the request ends within the configured deadline with a bounded availability error.
3. **Given** an unexpected daemon response, **When** it is adapted, **Then** no internal path, endpoint, secret, or raw server error is returned to the MCP host.
4. **Given** the MCP adapter source, **When** its dependencies and responsibilities are reviewed, **Then** scheduler truth and mutation remain in the daemon and core packages contain no MCP-specific business logic.

### Edge Cases

- No daemon is running, the configured Unix socket or named pipe is missing, or operating-system access is denied.
- The host requests a resource before initialization, after cancellation, or with a malformed or expired continuation cursor.
- A collection is empty, exactly at the page limit, changes between pages, or contains duplicate timestamps.
- Task names, group names, descriptions, alert messages, and output contain control characters, ANSI escapes, Markdown, XML-like instructions, invalid UTF-8 replacement characters, or extremely long content.
- A run has no output, is active, has no end time or exit code, or was already truncated by the daemon.
- A task contains every protected execution field and notification credentials exist elsewhere in daemon state.
- The daemon version differs from the CLI version or returns a response the adapter cannot decode.
- The stdio host closes stdin while a daemon request is in flight.

## Requirements

### Functional Requirements

- **FR-001**: The repository MUST define a one-page initial MCP surface and a permission model with Observe, future Operate, and future Manage classes.
- **FR-002**: S069 MUST enable only Observe resources; it MUST advertise no mutation tools and MUST perform no scheduler mutation.
- **FR-003**: The adapter MUST use the official Model Context Protocol Go SDK and support protocol revisions `2026-07-28` and `2025-11-25` through SDK negotiation.
- **FR-004**: The explicit entry point MUST be `gosched mcp serve` and MUST communicate through stdin and stdout only.
- **FR-005**: The entry point MUST open no TCP listener, MUST remain disabled unless a host launches it, and MUST exit promptly when its host disconnects or cancels.
- **FR-006**: Protocol messages MUST be the only stdout content after the MCP server starts; diagnostics MUST use stderr and MUST not contain protected values.
- **FR-007**: The approved resource surface MUST include daemon health, active tasks, upcoming schedules, recent alerts, and bounded recent run history.
- **FR-008**: Collection responses MUST be deterministic, MUST contain at most 100 records, and MUST expose an opaque continuation URI when more records are available.
- **FR-009**: One returned command-output excerpt MUST contain at most 8 KiB after valid UTF-8 normalization and MUST report source or adapter truncation.
- **FR-010**: One returned user-controlled name, description, or alert message MUST contain at most 2 KiB after valid UTF-8 normalization and MUST report adapter truncation when applicable.
- **FR-011**: Resource envelopes MUST identify user-controlled fields as untrusted content and MUST instruct consuming hosts to treat them as data rather than instructions.
- **FR-012**: Task resources MUST exclude commands, arguments, environment values, stdin, working directories, run-as identities, trigger credentials, notification credentials, and raw schedule definitions.
- **FR-013**: Alert and run resources MUST exclude notification destinations, authorization values, trigger keys, environment values, stdin, and internal filesystem or IPC paths.
- **FR-014**: The implementation MUST map daemon response types into dedicated MCP-safe response types and MUST never serialize daemon domain objects directly.
- **FR-015**: Every daemon request MUST use the existing protected local IPC client, inherit its operating-system access decision, use a deadline no longer than 10 seconds, and stop on cancellation.
- **FR-016**: IPC denial MUST remain denial; the adapter MUST NOT retry through TCP, direct database access, elevated execution, or another identity.
- **FR-017**: Host-visible failures MUST use a bounded safe vocabulary for invalid resource, invalid cursor, unavailable daemon, denied access, timeout, cancellation, incompatible response, and internal failure.
- **FR-018**: Future Operate and Manage capabilities MUST remain disabled and documented as requiring explicit authentication, per-action authorization, audit attribution, and confirmation for destructive or externally visible actions.
- **FR-019**: MCP protocol and presentation logic MUST remain in an adapter package and CLI entry point; scheduler policy, persistence, execution, and authority MUST remain in existing daemon and core packages.
- **FR-020**: Codex setup documentation MUST provide a clean-install command, explain the local trust boundary and untrusted content, and state that the daemon must already be running and accessible to the launching identity.
- **FR-021**: Automated tests MUST cover SDK negotiation, discovery, every approved resource, pagination, bounds, redaction, untrusted-content labeling, protocol stream cleanliness, cancellation, host disconnect, IPC denial, no mutation surface, and cross-platform subprocess behavior.
- **FR-022**: Issues [#161](https://github.com/shruggietech/go-schedule/issues/161) and [#162](https://github.com/shruggietech/go-schedule/issues/162) MUST remain traceable through the specification, tasks, changelog, pull request, and verification record.

### Key Entities

- **Observe Resource Descriptor**: Stable URI, name, description, media type, trust classification, and bounded handler.
- **Observe Envelope**: Schema version, generation time, permission class, trust notice, pagination metadata, and safe payload.
- **Safe Task Summary**: Stable task and group identity, bounded untrusted display names, enabled and lifecycle state, timezone, safe schedule summary, readiness, policy summary, next occurrence, and update time.
- **Safe Schedule Summary**: Task identity, bounded untrusted display name, enabled and readiness state, safe human-facing schedule kind or summary, timezone, and bounded upcoming occurrence times.
- **Safe Run Summary**: Stable run and task identity, timing, outcome, exit code, trigger kind, bounded untrusted output excerpt, and truncation state.
- **Safe Alert Summary**: Stable alert and correlation identity, severity, kind, bounded untrusted message, creation time, and acknowledgement state.
- **Continuation Cursor**: Opaque versioned token that identifies a deterministic offset and rejects malformed input without exposing implementation details.
- **Permission Class**: Observe in S069, with disabled future Operate and Manage classes governed by stronger authorization and confirmation requirements.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A supported MCP client can launch `gosched mcp serve`, initialize, discover, and read all five approved Observe collections on Windows, macOS, and Linux without additional network configuration.
- **SC-002**: Automated secret canaries placed in commands, arguments, environment, stdin, working directories, run-as identities, trigger keys, notification credentials, and raw schedules produce zero matches in MCP stdout and decoded resource content.
- **SC-003**: Every collection response contains no more than 100 records, every output excerpt is at most 8 KiB, every bounded text field is at most 2 KiB, and continuation reaches every eligible record without duplication for an unchanged fixture.
- **SC-004**: Discovery reports zero tools and zero mutation capabilities, and tests observe zero scheduler state changes after reading every resource and template.
- **SC-005**: Host disconnect, cancellation, unavailable daemon, denied IPC, and request timeout all terminate or respond within 10 seconds with no goroutine leak, panic, raw endpoint, or protected value.
- **SC-006**: Focused unit, protocol, subprocess, race, and canonical `scripts/verify.sh all` verification passes before publication.

## Clarifications

### Session 2026-09-08

- Q: Which protocol revisions form the compatibility promise? A: The adapter uses the official Go SDK and explicitly supports the latest stable revision `2026-07-28` plus the broadly deployed `2025-11-25` revision; older revisions are outside the tested promise even if the SDK can negotiate them.
- Q: Are any tools included in the first surface? A: No. S069 advertises only Observe resources and resource templates. Operate and Manage are documented future classes and remain disabled.
- Q: How are large collections and untrusted output bounded? A: Collections use deterministic 100-record pages with opaque continuation URIs, output excerpts are capped at 8 KiB, and other user-controlled text is capped at 2 KiB with explicit truncation metadata.
- Q: Where does authorization occur? A: The MCP subprocess inherits the launching identity and reaches the daemon only through the existing protected Unix socket or Windows named pipe. A daemon denial is final.
- Q: May raw daemon models be serialized after omitting known fields? A: No. Dedicated allowlisted MCP response types are required so new daemon fields cannot cross the boundary accidentally.

## Assumptions

- The existing daemon IPC API, operating-system access controls, health, task-detail, alert, and run reads remain authoritative.
- Codex and other MCP hosts can launch a local command and communicate over standard input and standard output.
- A deterministic offset cursor is sufficient for the first local operational surface; pages describe a best-effort current snapshot and do not promise transactional consistency across daemon changes.
- Schedule observation uses existing safe human-facing summaries and next-run projections rather than exposing executable commands or raw recurrence definitions.
- Audit records are unnecessary for read-only Observe calls in S069; future Operate and Manage calls require attributable audit events before activation.
- Remote MCP, HTTP transports, OAuth, scheduler mutation, log resources, notification management, and arbitrary query tools are outside S069.
