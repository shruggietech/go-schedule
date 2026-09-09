# Feature Specification: Actor Permissions and Management Audit

**Feature Branch**: `codex/076-actor-permissions-audit`

**Created**: 2026-09-09

**Status**: Implemented

**Delivery**: Actor persistence, fail-closed authorization, intent-first audit, local API, typed client, CLI, documentation, focused race tests, and canonical eight-gate verification completed 2026-09-09 on review branch `codex/076-actor-permissions-audit` for #167.

**Input**: GitHub Issue #167, "[v1.4] Implement actor permissions and management audit"

## Clarifications

- The local operating-system client is represented by one built-in actor with Enroll capability. Existing local IPC access controls remain the authentication boundary, so this slice adds no local login or bearer credential.
- Audit writes use a durable intent-first model. An audited operation is blocked if the intent cannot be persisted, and an interrupted operation remains `uncertain` rather than being reported as successful.
- Audit retention is bounded to 10,000 events and 90 days. Pruning occurs transactionally when a new intent is recorded.
- Actor records are independent of credentials. They prepare the authorization model for future remote enrollment without implementing remote access, pairing, or credential issuance.
- Revocation is evaluated on every request by reloading the actor. A revoked actor cannot authorize another request on an already-open connection.

## User Scenarios & Testing

### User Story 1 - Apply one authority model everywhere (Priority: P1)

As an operator, I want every operation classified against the same Observe, Operate, Manage, and Enroll capability hierarchy so present and future clients cannot interpret authority differently.

**Why this priority**: A shared, fail-closed authorization vocabulary is the prerequisite for safe remote and MCP expansion.

**Independent Test**: Enumerate every current API operation, verify its required capability, exercise allowed and denied requests, and confirm unknown operations are rejected.

**Acceptance Scenarios**:

1. **Given** an actor with a capability equal to or stronger than an operation requirement, **When** the actor invokes the operation, **Then** authorization succeeds.
2. **Given** an actor with a weaker, expired, revoked, malformed, or unknown identity, **When** the actor invokes an operation, **Then** authorization fails closed.
3. **Given** an operation absent from the catalog, **When** authorization is requested, **Then** the operation is denied and cannot silently inherit authority.

### User Story 2 - Preserve frictionless local operation (Priority: P1)

As a local desktop or CLI user, I want existing IPC workflows to continue without a new login while actions remain attributable to the local operating-system boundary.

**Why this priority**: The feature must improve accountability without breaking established single-user local behavior.

**Independent Test**: Upgrade an existing database, invoke representative CLI and desktop-backed API operations without credentials, and verify successful authorization under the built-in local actor.

**Acceptance Scenarios**:

1. **Given** an existing installation, **When** the database migrates, **Then** exactly one protected built-in local actor exists with Enroll capability.
2. **Given** an authenticated local IPC connection, **When** a current supported operation is invoked without an actor header or token, **Then** it runs as the built-in local actor.
3. **Given** an attempt to revoke, expire, delete, or reduce the built-in local actor, **When** the change is submitted, **Then** it is rejected.

### User Story 3 - Inspect durable management evidence (Priority: P1)

As an administrator, I want meaningful mutations, denied attempts, and privileged reads recorded with bounded retention so I can determine who attempted what and whether it completed.

**Why this priority**: Management audit is a release-blocking safety requirement for later remote administration.

**Independent Test**: Perform successful, failed, denied, and interrupted audited operations, then query and export their redacted records and verify retention pruning.

**Acceptance Scenarios**:

1. **Given** an audited operation, **When** execution begins, **Then** a durable intent records actor, daemon, operation, target, correlation identifier, and time before the side effect.
2. **Given** an audited operation completes or fails, **When** the handler returns, **Then** the same event records `succeeded` or `failed`; interruption leaves `uncertain`.
3. **Given** an authorization denial, **When** the request is rejected, **Then** a `denied` event is recorded without persisting request bodies, response bodies, headers, raw errors, commands, environment, standard input, paths, keys, or secrets.
4. **Given** audit history exceeds either retention boundary, **When** the next intent is recorded, **Then** events older than 90 days and excess events beyond the newest 10,000 are pruned.

### User Story 4 - Stage future client identities (Priority: P2)

As an administrator, I want to create, inspect, modify, revoke, filter, and export actor and audit records so future remote credentials can attach to an already-governed identity model.

**Why this priority**: Separating actor administration from credential enrollment keeps this slice useful and compatible with follow-on issues #168 and #169.

**Independent Test**: Use the local API, typed client, and CLI to manage non-built-in actors and filter or export audit events, then verify identical persisted results.

**Acceptance Scenarios**:

1. **Given** Enroll authority, **When** a valid actor is created or updated, **Then** the record persists independently of any credential.
2. **Given** Enroll authority, **When** an actor is revoked, **Then** the next request made as that actor is denied even on a reused connection.
3. **Given** audit filters, **When** history is listed or exported, **Then** matching events are returned in stable chronological form with documented fields and no secrets.

### Edge Cases

- Actor display names are trimmed, must contain 1 to 80 Unicode characters, and reject control characters.
- Capability, actor kind, actor state, audit result, and operation identifiers reject unknown values.
- Concurrent requests by an actor being revoked are authorized from the state loaded for each request; all requests beginning after revocation is committed are denied.
- Audit completion failure leaves the persisted intent `uncertain` and returns an internal failure without claiming success.
- Empty audit filters return the retained history; invalid timestamps, limits, actors, operations, or results return validation errors.
- Export ordering is deterministic when multiple events share a timestamp.
- Existing task execution logs remain separate from management audit records.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST define a closed, monotonic capability hierarchy of Observe, Operate, Manage, and Enroll.
- **FR-002**: The system MUST define actor kinds for local operating-system access, desktop, CLI, JSON clients, and MCP clients.
- **FR-003**: The system MUST persist actors with an opaque identifier, kind, display name, capability, lifecycle state, creation time, update time, and optional expiration time.
- **FR-004**: The system MUST create exactly one protected built-in local operating-system actor with Enroll capability during database initialization and migration.
- **FR-005**: Existing local IPC clients MUST resolve to the built-in local actor without a new login, header, token, or credential.
- **FR-006**: The built-in local actor MUST NOT be deleted, revoked, expired, or assigned a capability below Enroll.
- **FR-007**: The system MUST maintain one shared operation catalog that maps every registered management API operation to a minimum capability, target classification, and audit classification.
- **FR-008**: The operation catalog MUST classify safe reads as Observe, execution and acknowledgement actions as Operate, configuration mutations as Manage, and actor, credential, enrollment, listener, and identity administration as Enroll.
- **FR-009**: Authorization MUST fail closed for unknown operations, actors, capabilities, kinds, or states and for expired or revoked actors.
- **FR-010**: Actor state MUST be loaded for each request so revocation applies to the next request on a persistent connection.
- **FR-011**: The API router and operation catalog MUST have an automated completeness check preventing an uncataloged registered management route.
- **FR-012**: The system MUST audit meaningful mutations, authorization denials, actor administration, and privileged reads of actor and audit data while leaving ordinary task execution logs separate.
- **FR-013**: An audit event MUST contain an opaque event identifier, actor identifier when known, daemon identifier, operation identifier, target kind and optional target identifier, result, correlation identifier, occurrence time, and optional completion time.
- **FR-014**: Before an audited operation executes, the system MUST durably persist an `uncertain` audit intent and MUST reject the operation if the intent cannot be persisted.
- **FR-015**: After an audited operation returns, the system MUST update the same event to `succeeded` or `failed`; an uncompleted intent MUST remain `uncertain`.
- **FR-016**: Authorization denials MUST be recorded as `denied` without executing the protected operation.
- **FR-017**: Audit records MUST NOT store request bodies, response bodies, authorization headers, raw error text, commands, environment variables, standard input, filesystem paths, cryptographic keys, credentials, or secrets.
- **FR-018**: Audit retention MUST keep no more than the newest 10,000 events and no event older than 90 days, pruning both limits transactionally when a new audited intent is inserted.
- **FR-019**: Enroll-authorized clients MUST be able to list, create, update, and revoke non-built-in actors through the local API, typed Go client, and CLI.
- **FR-020**: Enroll-authorized clients MUST be able to list filtered audit events and export them as deterministic newline-delimited JSON through the local API, typed Go client, and CLI.
- **FR-021**: Actor display names MUST contain 1 to 80 Unicode characters after trimming and MUST reject control characters.
- **FR-022**: Actor records MUST remain independent of credentials, and this slice MUST NOT add remote listeners, TLS, pairing phrases, bearer authentication, keyring storage, or MCP mutation authority.
- **FR-023**: Documentation MUST describe the capability matrix, local actor trust boundary, actor lifecycle, audit schema, redaction, retention, filtering, export, migration behavior, and follow-on integration points for issues #168 and #169 without claiming general policy support.

### Key Entities

- **Capability**: One level in the closed Observe, Operate, Manage, and Enroll authority hierarchy.
- **Actor**: A durable client identity and lifecycle record that can later be associated with a credential without containing one.
- **Operation Definition**: The canonical mapping from an API or future transport operation to required capability, target, and audit treatment.
- **Audit Event**: A durable, redacted record of authorization and execution outcome for a management action or privileged read.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Automated tests account for 100 percent of registered management API operations in the shared catalog.
- **SC-002**: Every capability level is tested for both an allowed operation and a denied stronger operation.
- **SC-003**: Existing local CLI, desktop-backed, and typed-client workflows complete without new authentication input after migration.
- **SC-004**: Successful, failed, denied, and interrupted audited operations produce the expected durable result without storing prohibited sensitive fields.
- **SC-005**: Revocation committed between two requests on one persistent connection causes the second request to be denied.
- **SC-006**: Retention tests prove that no more than 10,000 events and no event older than 90 days remain after a new intent.
- **SC-007**: API, typed client, and CLI tests demonstrate actor administration plus deterministic filtered audit listing and export.
- **SC-008**: The full repository verification suite, including race detection, formatting, documentation validation, and security checks, passes.

## Assumptions

- Existing local IPC ownership and access controls continue to authenticate the local operating-system boundary.
- The daemon identity introduced in S074 supplies the audit daemon identifier.
- Credential issuance and remote transport identity resolution remain owned by issues #168 and #169.
- The current MCP listener remains Observe-only and uses the shared capability and operation vocabulary without gaining actor administration or mutation endpoints in this slice.

## Dependencies

- Parent roadmap issue: #18.
- Completed architecture dependencies: #165 and #166.
- Follow-on consumers: #168 and #169.
