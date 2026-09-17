# Feature Specification: MCP Operate Authority

**Feature Branch**: `codex/091-mcp-operate-authority`

**Created**: 2026-09-17

**Status**: In Progress

**Delivery**: Issue #178

**Input**: Work slice S091, add narrowly scoped MCP Operate authority for existing tasks.

## User Scenarios & Testing

### User Story 1 - Deliberately grant task operation (Priority: P1)

A local administrator can start an MCP connection with Observe authority by default or explicitly choose Operate authority for that connection.

**Why this priority**: Mutation tools must never appear through an implicit upgrade from the existing Observe-only contract.

**Independent Test**: Connect once with the default permission and once with Operate, then confirm that only the Operate connection discovers the three task-operation tools.

**Acceptance Scenarios**:

1. **Given** a default stdio or localhost MCP connection, **When** the client lists tools, **Then** no mutation tool is exposed.
2. **Given** an explicitly requested Operate connection, **When** the client lists tools, **Then** run-now, enable, and disable are exposed while create, edit, delete, secret, and permission tools remain absent.

### User Story 2 - Operate one explicit task on one explicit daemon (Priority: P1)

An authorized MCP client can enable, disable, or request an immediate run for an existing task only when it supplies the expected daemon identity and task identity.

**Why this priority**: Target ambiguity is unacceptable for an agent mutation.

**Independent Test**: Invoke each tool against a matching daemon and task, then repeat with the wrong daemon, missing task, disabled group, and non-runnable task.

**Acceptance Scenarios**:

1. **Given** a matching daemon and runnable task, **When** Operate requests run-now, **Then** the daemon accepts exactly one request and returns explicit target and outcome data.
2. **Given** a stale or wrong daemon identity, **When** any tool is called, **Then** the operation is rejected before mutation.
3. **Given** an Observe-only or revoked session, **When** an Operate tool is attempted, **Then** the daemon denies it and records the denial.

### User Story 3 - Diagnose every attempt safely (Priority: P2)

An administrator can distinguish accepted, rejected, denied, and uncertain MCP operations and correlate each attempt with its MCP client, daemon, task, and request identity.

**Why this priority**: Agent automation must not turn transport uncertainty or retries into invisible duplicate work.

**Independent Test**: Exercise success, validation failure, authorization denial, timeout, retry, and revocation cases and inspect both tool output and audit history.

**Acceptance Scenarios**:

1. **Given** a repeated request identifier, **When** the same tool call is retried, **Then** the original result is returned without another mutation.
2. **Given** a connection failure after dispatch cannot be distinguished from completion, **When** the tool returns, **Then** its outcome is uncertain rather than failed or accepted.
3. **Given** a completed or denied operation, **When** audit history is inspected, **Then** it names the MCP actor, daemon, operation, task, and final result.

### Edge Cases

- A task identifier is empty, oversized, or names a deleted task.
- The daemon identity changes between discovery and invocation.
- The task is disabled by its group, terminal, missing a command, or otherwise not ready.
- Two calls reuse a request identifier with different tool or target arguments.
- An HTTP credential rotates or the listener is disabled while a client is connected.
- A stdio host disconnects while an operation is in flight.
- Hostile task content resembles tool instructions or attempts to alter authority.

## Requirements

### Functional Requirements

- **FR-001**: Existing MCP connections MUST remain Observe-only unless Operate is explicitly requested.
- **FR-002**: Operate MUST expose only task run-now, enable, and disable tools.
- **FR-003**: Every tool input MUST include the expected daemon identity, task identity, and caller-generated request identity.
- **FR-004**: The daemon MUST validate its current identity before dispatching a mutation.
- **FR-005**: Operate calls MUST use the existing task API, readiness rules, authorization catalog, and intent-first audit path.
- **FR-006**: Task enable and disable MUST require Operate authority while task definition changes and deletion continue to require Manage authority.
- **FR-007**: Every Operate connection MUST have a runtime MCP actor identity that is inactive after session revocation, listener disable, or clean stdio shutdown.
- **FR-008**: Revocation MUST affect subsequent requests on existing connections without requiring daemon restart.
- **FR-009**: Tool results MUST identify permission, operation, daemon, task, request, outcome, and a safe explanation.
- **FR-010**: Identical retries MUST not repeat a mutation, while request-identity reuse with different arguments MUST be rejected.
- **FR-011**: Transport ambiguity MUST be reported as uncertain and MUST NOT be silently retried.
- **FR-012**: Observe clients MUST remain unable to discover or invoke Operate tools.
- **FR-013**: No tool or response may reveal task commands, environment variables, stdin, credentials, protected paths, or other execution inputs.
- **FR-014**: Runtime MCP sessions and their bearer values MUST be memory-only and MUST NOT survive daemon restart.

### Key Entities

- **MCP session**: Runtime identity binding a client name, capability, actor, secret digest, creation time, and revocation state.
- **Operate request**: One tool, expected daemon, task, and client-generated request identifier.
- **Operate result**: Bounded structured outcome for one attempted mutation.
- **Deduplication record**: Short-lived mapping from request identity and exact operation fingerprint to its first result.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Observe conformance still discovers five resources, four templates, and zero tools.
- **SC-002**: Operate conformance discovers exactly three mutation tools and no Manage capability.
- **SC-003**: All successful and denied Operate attempts produce attributable audit evidence.
- **SC-004**: Repeating an identical request identifier 100 times produces one underlying mutation.
- **SC-005**: Revocation prevents the next request on an already established HTTP or stdio connection.
- **SC-006**: Root race, coverage, documentation, automation, and MCP conformance gates pass.

## Assumptions

- S091 supports local stdio and authenticated numeric-loopback HTTP. Remote MCP remains issue #180.
- Runtime sessions are connection-scoped authority, not the durable grant administration planned by issue #181.
- Per-call human confirmation is not required after an administrator deliberately starts an Operate session.
- Manage authority, definition mutation, deletion, secret access, and enrollment remain out of scope.
