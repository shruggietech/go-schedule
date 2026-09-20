# Feature Specification: MCP Manage Authority

**Feature Branch**: `codex/092-mcp-manage-authority`

**Created**: 2026-09-20

**Status**: Implemented

**Delivery**: review branch `codex/092-mcp-manage-authority`; canonical eight-gate verification passed 2026-09-20 for issue #179.

**Input**: Work slice S092, add explicit MCP Manage authority for bounded automation-definition changes.

## User Scenarios & Testing

### User Story 1 - Deliberately grant Manage authority (Priority: P1)

A local administrator can start an MCP connection with Observe authority by default, Operate authority for existing-task controls, or Manage authority for bounded definition changes.

**Why this priority**: Definition mutation must never appear through an implicit upgrade or an Operate grant.

**Independent Test**: Connect with each permission and compare discovery. Only Manage discovers definition tools, while Manage also retains the lower-capability Observe and Operate surfaces.

**Acceptance Scenarios**:

1. **Given** a default or Observe connection, **When** tools are listed, **Then** no Operate or Manage tools are exposed.
2. **Given** an Operate connection, **When** tools are listed, **Then** only the three existing task-operation tools are exposed.
3. **Given** an explicitly requested Manage connection, **When** tools are listed, **Then** bounded task, group, chain, trigger, watcher, and notification-assignment management tools are exposed without permission, enrollment, credential-reveal, database, or raw-filesystem tools.

### User Story 2 - Manage one definition through shared lifecycle rules (Priority: P1)

An authorized MCP client can create, update, or delete one supported definition on one exact daemon through the same API validation and audit paths used by human clients.

**Why this priority**: A second mutation path would drift from scheduler validation, readiness, and audit behavior.

**Independent Test**: Exercise each supported definition family with valid and invalid inputs, wrong daemon identity, missing references, protected data, and incompatible actions, then compare API behavior and audit evidence.

**Acceptance Scenarios**:

1. **Given** a matching daemon and valid request, **When** a Manage tool changes one definition, **Then** the ordinary API lifecycle completes and the result returns only bounded identity and outcome data.
2. **Given** a stale daemon identity, invalid definition, protected secret request, or insufficient session, **When** a Manage tool is invoked, **Then** the mutation is rejected before unauthorized state is exposed or changed.
3. **Given** a trigger is created, **When** the API generates a trigger secret, **Then** the MCP result withholds the secret and returns only the trigger identity and a safe explanation.

### User Story 3 - Apply optional confirmation and diagnose every attempt (Priority: P2)

An administrator can require caller confirmation for Manage mutations on a connection, and can diagnose accepted, rejected, denied, and uncertain attempts without receiving sensitive definition content.

**Why this priority**: Destructive automation changes need a deployable safety policy and attributable outcomes without forcing interactive confirmation in unattended environments.

**Independent Test**: Run the same mutation with confirmation policy disabled and enabled, retry request identifiers, simulate transport ambiguity, and inspect audit records.

**Acceptance Scenarios**:

1. **Given** confirmation policy is enabled, **When** a Manage call omits confirmation, **Then** it is rejected without dispatch.
2. **Given** confirmation policy is disabled, **When** a valid Manage call is made, **Then** no extra interactive prompt is required.
3. **Given** a repeated request identifier for the exact same mutation, **When** it is retried, **Then** the first result is returned without another mutation; reuse for different input is rejected.
4. **Given** a response cannot prove the final state, **When** the call completes, **Then** it reports an uncertain outcome and instructs the caller to inspect current state before retrying.

### Edge Cases

- A create or update embeds task-controlled text that resembles MCP instructions or requests more authority.
- A request identifier is reused after a session is revoked or with different definition content.
- A delete targets a missing object, a group with children, or an object referenced by another definition.
- A task update attempts to carry environment variables or stdin content through the MCP contract.
- A trigger create succeeds but its generated secret cannot be returned to the MCP caller.
- An assignment replacement names a missing channel or mixes task and group scope.
- A client requests unsupported bulk mutation or a future definition kind.
- The daemon changes identity between discovery and mutation.

## Requirements

### Functional Requirements

- **FR-001**: Existing MCP connections MUST remain Observe-only unless a higher permission is explicitly requested.
- **FR-002**: Operate MUST remain limited to run-now, enable, and disable for existing tasks and MUST NOT discover or invoke Manage tools.
- **FR-003**: Manage MUST expose bounded tools for tasks, groups, completion chains, external triggers, filesystem watchers, and task or group notification assignments.
- **FR-004**: Each Manage call MUST change at most one definition or atomically replace one assignment scope; bulk mutation is unsupported.
- **FR-005**: Every call MUST include exact daemon identity and a caller-generated request UUID, and update or delete calls MUST include exact object identity.
- **FR-006**: Manage handlers MUST call the ordinary versioned API and MUST NOT mutate the store, scheduler, permissions, actors, sessions, or files directly.
- **FR-007**: Tool inputs MUST omit task environment and stdin values, trigger keys, actor capabilities, enrollment credentials, and raw database or filesystem operations.
- **FR-008**: Tool results MUST include schema version, permission, operation, daemon, affected object kind and identity, request identity, outcome, and a safe explanation only.
- **FR-009**: Tool results MUST NOT return commands, arguments, paths, environment values, stdin, notification authorization, trigger keys, or other stored sensitive definition content.
- **FR-010**: Trigger creation MUST withhold generated keys and commands from the MCP result; later credential reveal or rotation remains outside Manage authority.
- **FR-011**: A connection MAY require an explicit confirmation boolean on every Manage mutation, while the default policy MUST permit deliberate unattended Manage sessions without per-call prompts.
- **FR-012**: Identical retries MUST not repeat a mutation, while request-identity reuse with different arguments MUST be rejected.
- **FR-013**: Transport or server ambiguity MUST report an uncertain outcome and MUST NOT trigger automatic retry.
- **FR-014**: Validation, reference integrity, readiness effects, delete semantics, and assignment replacement MUST match the existing GUI, CLI, and JSON API behavior.
- **FR-015**: Every dispatched request MUST use the MCP actor so shared intent-first audit records identify the proposed operation, final result, daemon, and affected object.
- **FR-016**: Runtime Manage sessions MUST remain memory-only and subsequent requests MUST fail after revocation, listener disable, or clean stdio shutdown.
- **FR-017**: Unsupported actions, bulk input, incompatible schema versions, and attempts to raise authority through managed content MUST fail closed.

### Key Entities

- **Manage session**: Runtime MCP identity with explicit Manage capability, optional confirmation policy, expiry, and revocation state.
- **Manage request**: One tool operation bound to exact daemon, request, object kind, optional object identity, and bounded definition input.
- **Manage result**: Redacted outcome projection containing identity, status, and safe explanation without definition contents.
- **Deduplication record**: Bounded runtime record binding request identity to the exact mutation fingerprint and first result.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Observe discovery exposes zero tools, Operate exposes exactly three tools, and Manage adds exactly six bounded definition tools.
- **SC-002**: Conformance covers create, update, delete, assignment replacement, confirmation, denial, revocation, retry, incompatible input, and uncertain outcome behavior.
- **SC-003**: Repeating an identical request identifier 100 times produces one underlying mutation, while any differing retry is rejected.
- **SC-004**: Every dispatched Manage attempt has attributable audit evidence and no response contains a seeded secret or protected execution input.
- **SC-005**: All root race, coverage, GUI, documentation, automation, and MCP conformance gates pass.

## Assumptions

- S092 remains local to stdio and authenticated numeric-loopback HTTP; remote MCP transport remains issue #180.
- The existing versioned JSON API is the source of truth for validation and lifecycle semantics.
- One tool per definition family keeps discovery bounded while an action discriminator selects create, update, or delete.
- Task environment and stdin management remain available to human clients but are intentionally excluded from the MCP contract until protected-reference support exists.
- Trigger creation is allowed, but its generated credential is withheld; an Enroll-capable human client must perform any later reveal or rotation.
- Persistent client grants and permission administration remain issue #181 and are outside this slice.
