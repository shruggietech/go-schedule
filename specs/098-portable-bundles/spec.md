# Feature Specification: Portable Automation Bundles and Drift Reporting

**Feature Branch**: `codex/098-portable-bundles`

**Created**: 2026-09-21

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Implementation complete, pull request and review pending

**Input**: Work slice S098 and GitHub issue #184.

## Clarifications

### Session 2026-09-21

- Q: Does a bundle identify a machine or carry secret material? -> A: No. Bundles carry portable automation intent only. They exclude daemon identity, credentials, trigger keys, notification channel secrets, run history, and machine-local paths or users.
- Q: Does omission remove target state? -> A: No. Apply is additive and update-only by default. Removal is unavailable unless a future explicitly scoped format revision introduces an explicit removal list.
- Q: Does a bundle synchronize multiple daemons? -> A: No. Every export, comparison, preview, and apply names exactly one daemon. There is no background reconciliation, desired state, cross-daemon transaction, or retry of an uncertain remote mutation.
- Q: What is portable in S098? -> A: Groups, task scheduling intent, and completion chains are portable in v1. Watchers, external triggers, trigger sets, notification policy references, and every secret-bearing or machine-specific field are explicitly reported as exclusions rather than silently approximated.

## User Scenarios & Testing

### User Story 1 - Export a reviewable portable definition (Priority: P1)

An operator can export a selected daemon's automation definitions into a deterministic, versioned bundle that can be reviewed or stored without exposing secrets or machine identity.

**Why this priority**: A safe export format is the foundation for every later comparison and transfer operation.

**Independent Test**: Create equivalent automation state in different insertion orders, export twice, and confirm byte-for-byte equivalent canonical JSON with no credentials, raw trigger keys, history, or daemon identifier.

**Acceptance Scenarios**:

1. **Given** a daemon with tasks, groups, sources, and policies, **When** an operator exports a bundle, **Then** the bundle has a declared version and deterministic ordering.
2. **Given** secret-bearing triggers, notification channels, and machine-local fields, **When** the bundle is exported, **Then** secret values and host identity are absent and nonportable fields are reported as exclusions.
3. **Given** semantically equivalent state created in a different order, **When** it is exported, **Then** the resulting bundle is identical.

---

### User Story 2 - Inspect compatibility, conflict, and drift without mutation (Priority: P1)

An operator can validate a bundle against one explicitly selected daemon and see an exhaustive, stable plan of create, update, unchanged, skipped, conflict, and incompatibility outcomes before anything changes.

**Why this priority**: Operators need a trustworthy review surface before applying portable configuration to a separate system.

**Independent Test**: Compare a valid bundle against a target containing matching, changed, and additional definitions; assert that the comparison reports each difference and leaves all target records unchanged.

**Acceptance Scenarios**:

1. **Given** a valid bundle and an explicit target, **When** the operator requests a preview, **Then** every planned change identifies that target and its action without modifying it.
2. **Given** a bundle that needs unsupported platform capabilities or has invalid references, **When** it is validated, **Then** the operator receives actionable findings before an apply option is offered.
3. **Given** target-only definitions, **When** the operator compares the bundle and target, **Then** they are shown as read-only drift rather than inferred deletions.

---

### User Story 3 - Apply a reviewed plan to one daemon (Priority: P1)

An operator can explicitly apply the exact reviewed plan to one selected daemon, with stable matching, deterministic order, per-item outcomes, and clear handling of partial or uncertain remote outcomes.

**Why this priority**: Transfer has value only when review and execution have the same meaning.

**Independent Test**: Preview a bundle, apply the returned plan to a selected daemon, and confirm that only planned create or update operations occur in plan order while omitted target state remains untouched.

**Acceptance Scenarios**:

1. **Given** a valid preview and a confirmed target, **When** the operator applies it, **Then** only the previewed non-conflicting operations execute against that exact target.
2. **Given** an item changes after preview or an operation fails, **When** apply runs, **Then** the affected item reports rejected, failed, or uncertain without replaying or rolling back unrelated items.
3. **Given** a target that is not selected or no longer matches the preview identity, **When** apply is requested, **Then** no mutation occurs.

---

### User Story 4 - Work from the desktop and CLI (Priority: P2)

An operator can export, inspect, compare, and apply a bundle through the desktop control center or headless CLI with the same safety information and target context.

**Why this priority**: Portable configuration must remain useful to both desktop and headless deployments.

**Independent Test**: Run the CLI JSON workflow and the desktop bundle workflow against a local daemon and confirm equivalent plans, target labels, conflict information, and outcome summaries.

**Acceptance Scenarios**:

1. **Given** a local or enrolled remote daemon, **When** an operator opens the bundle workflow, **Then** the selected daemon is visible throughout export, preview, comparison, and apply.
2. **Given** a plan with conflicts or compatibility findings, **When** it is rendered by either interface, **Then** action-required state is distinguishable without color alone.

## Edge Cases

- An unknown future bundle version, duplicate logical identity, malformed JSON, or unsupported field is rejected without target mutation.
- A bundle may reference a group, task, chain, trigger, watcher, trigger set, or policy that is unavailable on the target. The resulting item is a named conflict or compatibility finding, never an implicit approximation.
- Circular chains, invalid schedule expressions, missing source task references, and platform-specific watcher paths are caught during preview.
- A remote connection loss after a mutation request yields an uncertain per-item outcome and is never automatically replayed.
- A second export, preview, or apply may not depend on process-local hidden state; only the supplied bundle, selected target, and returned plan identity are authoritative.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST define a versioned bundle document with a deterministic canonical representation for portable groups, task scheduling intent, and completion chains, plus explicit exclusions for unsupported record families.
- **FR-002**: Export MUST omit daemon identity, credentials, client credentials, protected secret references, raw trigger keys, notification endpoints and authorization, run history, audit records, and machine-specific filesystem paths, users, and executable locations.
- **FR-003**: Export MUST sort every collection and normalize optional values so equivalent supported state produces byte-identical bundle content.
- **FR-004**: Validation and comparison MUST name the selected target, validate schema version and references, and report compatibility, conflict, create, update, unchanged, skipped, and target-only drift findings without mutation.
- **FR-005**: Matching MUST use stable portable logical identities supplied in the bundle, never display names alone. Duplicate identities, name ambiguity, or missing references MUST be conflicts.
- **FR-006**: Apply MUST require an explicit target and a preview bound to that target's daemon identity and current state fingerprint. It MUST execute only reviewable create or update operations in deterministic dependency order.
- **FR-007**: Omitted bundle objects MUST never delete or disable target state. Removal requires explicit future product scope and is not implemented by this feature.
- **FR-008**: Apply MUST return a per-item terminal outcome: applied, unchanged, rejected, failed, or uncertain. It MUST not replay uncertain mutations or roll back previously successful independent items.
- **FR-009**: Local and remote access MUST use the same versioned daemon contract and established authorization boundary. Read-only bundle operations require Observe authority, while apply requires Manage authority.
- **FR-010**: The CLI MUST support machine-readable and human-readable export, validate, compare, preview, and apply output. The desktop MUST expose the corresponding workflow with visible target context and confirmation.
- **FR-011**: Drift comparison MUST remain read-only and explain target-only, bundle-only, and changed portable intent.
- **FR-012**: The feature MUST not introduce background synchronization, desired-state reconciliation, leader election, distributed locking, cross-daemon transactions, or clustered task execution.

### Key Entities

- **Automation bundle**: A versioned, canonical declaration of portable automation intent and logical identities.
- **Bundle item**: One portable object, its logical identity, portable data, and references to other bundle items.
- **Bundle plan**: A target-bound, deterministic preview containing compatibility findings and intended item operations.
- **Plan item outcome**: The auditable result of one apply attempt, including identity, action, terminal state, and explanatory message.
- **Drift finding**: A read-only difference between a bundle and one selected daemon.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Re-exporting equivalent supported automation state yields byte-identical bundle data in 100 repeated comparisons.
- **SC-002**: Validation identifies every malformed version, duplicate identity, unavailable dependency, and unsupported platform field before an operator can apply a plan.
- **SC-003**: A preview and apply display the same target identity and item actions, with no target mutation during preview or drift comparison.
- **SC-004**: A bundle operation against 100 portable objects completes with deterministic ordering and an item-level outcome for every planned object.
- **SC-005**: Secret scanning of exported examples and automated tests finds no credential, trigger key, notification endpoint, or daemon identity value.

## Assumptions

- Existing task, group, chain, watcher, trigger, trigger-set, notification, identity, remote-access, and authorization foundations remain authoritative.
- Logical identities are durable user-facing portable keys generated once for pre-existing records when first bundled; display names remain editable metadata.
- The v1 bundle format has no removal operations and carries no machine-specific execution inputs.
- The system may reject unsupported content rather than attempting lossy translation.
- This feature delivers a development capability through a PR, not a public release.
