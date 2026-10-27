# Feature Specification: Historical Windows incident differential diagnosis

**Feature Branch**: `codex/123-historical-windows-diagnosis`

**Created**: 2026-10-27

**Status**: Implemented

**Delivery**: Owner-merged PR #393 as a2e4653021b6007135550ed8f84f1421f9072db0 on 2026-10-27. All 23 exact-main checks and Pages deployment succeeded. delivery.md and verification.md retain the separate local/campaign/review identities, reconstruction limits and unresolved #355/#356 causes.

**Input**: Owner-authorized S123 bundles #355/#356, uses Spec-Kit, performs bounded historical comparison, corrects demonstrated defects, publishes a PR and handles at most two Codex rounds. Stop before merge, tagging or promotion.

## User Scenarios & Testing

### User Story 1 - Attribute original runtime failures (Priority: P1)

A maintainer compares the original and current writer behavior without confusing a historical runtime crash, prerequisite refusal or separately corrected observer defect.

**Why this priority**: #355 has no demonstrated causal resolution despite successful follow-ups.

**Independent Test**: Match preserved input identities and retain each six-writer outcome on fresh/existing stores.

**Acceptance Scenarios**:

1. **Given** retained original inputs, **When** comparison is prepared, **Then** source/binary identities and instrumentation availability are explicit and immutable originals remain unchanged.
2. **Given** concurrent invocations, **When** an interpreter, prerequisite or writer fails, **Then** its original result remains visible without retrying into success or reducing record expectations.

### User Story 2 - Distinguish historical IPC behavior (Priority: P1)

A maintainer compares the five original operation shapes on historical/current candidates while distinguishing execution, protocol, inspection and cleanup.

**Why this priority**: #356 included successful commands followed by failed inspection.

**Independent Test**: Exercise fresh owned candidates with selected diagnostic allowances and independent outcomes.

**Acceptance Scenarios**:

1. **Given** historical short and current supported allowances, **When** compared, **Then** overrides remain attributable without changing defaults.
2. **Given** incompatible runners or unavailable collectors, **When** preparation checks them, **Then** limitations and mixed-component identities are explicit without modifying historical binaries.
3. **Given** a demonstrated component defect, **When** corrected, **Then** the outcome is restored with authorization, record and lifecycle contracts intact.

### User Story 3 - Receive a bounded causal disposition (Priority: P2)

A maintainer receives separate incident conclusions, linked corrections where justified and supported guidance.

**Why this priority**: Successful repetitions are follow-up observations, not historical cause.

**Independent Test**: Review a symptom/hypothesis/outcome matrix with precise identities and remaining uncertainty.

**Acceptance Scenarios**:

1. **Given** no matching recurrence before the shared deadline, **When** the campaign finishes, **Then** the affected original issue stays open and speculative replay expansion stops.
2. **Given** merged S122, **When** records are reconciled, **Then** merged-main acceptance, application-head reviews and later CI-only correction remain distinct.

### Edge Cases

- Unrecoverable dirty source requires explicit reconstruction differences, never a clean-commit claim.
- Identity mismatch or incompatibility stops the affected condition before mutation.
- Unexpected failure stops its experiment for analysis, preserving the original deadline and first failure.
- Partial inspection cannot satisfy complete-state assertions.
- Historical absence of a collector cannot be presented as a successful collected observation.

## Requirements

### Functional Requirements

- **FR-001**: Preserve original failures, finalized artifacts and binary inputs unchanged; use isolated copies and read-only historical references.
- **FR-002**: Identify each source, runtime, payload, runner, component mix and diagnostic override before execution.
- **FR-003**: Compare recovered historical/current payload behavior across retained/current runtimes with six concurrent writers, fresh/existing schemas, literal identities and exact record assertions.
- **FR-004**: Attribute runtime crashes, launch failures, prerequisite refusals and writer failures separately, preserving native outcomes.
- **FR-005**: Compare the five original service/MCP/group shapes with prerequisites and separate execution, protocol, inspection and cleanup outcomes.
- **FR-006**: Distinguish historical short-allowance diagnostics from supported defaults, actual memory/WAL selection and available instrumentation.
- **FR-007**: Declare the matrix, modest resource bounds, one shared execution deadline and stop conditions before execution. Report preparation and later correction acceptance separately.
- **FR-008**: Correct demonstrated defects under separate linked bugs with first failures and discriminating regressions; reuse existing collectors and fixed QA reporting.
- **FR-009**: Leave original issues open if their causal outcomes remain unconfirmed; require matching evidence for project/upstream attribution.
- **FR-010**: Deliver a concise causal matrix, maintainer guidance and reconciled S122 merge records.
- **FR-011**: Run silently on owned fixtures with access controls intact, no installed-service/account-session operations or system-wide stress, and stop before merge or release operations.

### Key Entities

- **Candidate**: source/binary identity, dirty state, historical/current role and compatibility.
- **Condition**: payload/runtime or case/store/allowance selection, repetition, shared deadline and independent outcomes.
- **Disposition**: symptom, hypothesis, evidence, demonstrated correction or uncertainty and associated issue.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every attempted condition has complete identity and an original result or explicit unavailable/blocked disposition.
- **SC-002**: Zero retained original/finalized artifacts change and zero installed services or account sessions are operated.
- **SC-003**: Each original issue has a separate evidence-supported disposition; neither closes solely because checks pass.
- **SC-004**: Investigation stops within its declared allowance, with preparation, correction and acceptance timing reported separately.

## Assumptions

- Accepted baseline is merged S122 e8cefc1c8059364a0ee84fe22ac740842885d53d. Existing containment, cancellation and collectors are reused.
- Original S112 is dirty source ff465894a3ab9e7a90720f5b7388a5f44d7a40b25bf3bddaf2c8637f61f991af based on 2b827e0a128cff5698b84ff84a83628e31063423, not a clean historical commit.
- Existing Windows/WSL owned fixtures suffice unless an observation demonstrates otherwise.
- Conditional clustered scheduling and release operations remain separate.

## Clarifications

### Session 2026-10-27

- Q: Does no recurrence justify closing the original issues? A: No; retain needs: reproduction and stop speculative expansion.
- Q: May current-controller observations inherit historical-controller identity? A: No; label each binary/controller/payload and unavailable historical instrumentation separately.
- Q: How does publication proceed? A: The owner explicitly authorized push, PR and up to two Codex rounds in this kickoff; stop before merge, tagging or promotion.
