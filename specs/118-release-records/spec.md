# Feature Specification: S118 release record reconciliation

**Feature Branch**: `codex/118-release-records`

**Created**: 2026-10-25

**Status**: Implemented

**Delivery**: Implemented on review branch codex/118-release-records with all eight local gates passed in 700.437s and 2347 frozen inputs unchanged; see [verification.md](verification.md), [delivery.md](delivery.md) and [handoff.md](handoff.md). Owner-authorized PR publication and bounded external review follow; stop before merge, tag, staging or promotion.

**Input**: Reconcile release records for #354-#356 using Spec-Kit and deliver a reviewed source PR from accepted main 9ad13821c825d01efcc038b6b1634b187a854758. Source corrections #375/#376 merged through PR #377; public release remains v1.5.1 and v1.6.0 publication is pending.

## Clarifications

### Session 2026-10-25

- Q: Does S118 authorize a tag, draft staging, merge or promotion? A: No. The latest request authorizes only the reviewed source PR and its bounded CI/review delivery.
- Q: Must #355/#356 be closed as part of this slice? A: No. Reconcile their actual individual remaining outcomes, preserving the original unconfirmed causes and failed evidence.
- Q: Are more full candidate sweeps required? A: No. This records-focused increment reuses the immutable S117 capture and separately identifies corrected-source and merged-main acceptance. Run the mandated canonical source checks, not an unchanged catalog replay.
- Q: Should publication tasks be left unchecked in the new implementation ledger? A: No. Keep publication/review/merge receipts in Delivery and the PR; implementation tasks include only source and engineering verification work.

## User Scenarios & Testing

### User Story 1 - Understand completed delivery (Priority: P1)

A maintainer reads the S117 feature records and sees its completed source outcome, actual merge and delivery evidence without publication steps masquerading as unfinished implementation.

**Why this priority**: Stale lifecycle records obscure what has already been delivered and cause unnecessary repeat work.

**Independent Test**: Compare the specification, task ledger, inventory and actual merged PR; all describe the same implementation maturity and preserve the original candidate outcome.

**Acceptance Scenarios**:

1. **Given** the S117 source is merged, **When** its lifecycle is reconciled, **Then** required implementation work is resolved, the record identifies the actual reviewed and merged source, and publication receipts appear as delivery evidence.
2. **Given** an older candidate report failed and later corrections passed, **When** completed delivery is recorded, **Then** the original report remains failed and separately attributable.

### User Story 2 - Know exactly what remains for release (Priority: P1)

A maintainer can distinguish the pending v1.6.0 source, known defects, future tagged candidate and public v1.5.1 package and can follow one concrete handoff for later release operations.

**Why this priority**: Correct identities and issue dispositions prevent misleading release claims and accidental reuse of an earlier candidate.

**Independent Test**: Every current issue and release record names its delivered outcome, remaining scope and applicable source or publication boundary.

**Acceptance Scenarios**:

1. **Given** #355/#356 original causes remain unconfirmed, **When** their release disposition is reconciled, **Then** both remain visibly open and independently described without claiming a causal fix from passing repetitions.
2. **Given** pending v1.6.0 guidance, **When** source readiness is recorded, **Then** the actual reviewed notes, upgrade/rollback guidance and comparison bindings agree while the published version remains v1.5.1.
3. **Given** later tag and promotion decisions, **When** the owner reviews the handoff, **Then** it names the required source selection, actual staging identities and separate owner decisions without fabricating nonexistent assets.

### Edge Cases

- A squash merge has a different commit identity from the reviewed head; preserve both and their actual tree relationship.
- A valid failed candidate capture cannot be relabelled as final-source success after corrections.
- A completed publication task is evidence of delivery, not a remaining implementation task.
- A source record may be current while a release tag or draft does not exist; preserve the pending state.
- A later source correction or merge changes the candidate source identity and requires a fresh exact-source release decision.
- Newly discovered defects remain separate issue records; membership in this slice does not close #354/#355/#356.

## Requirements

### Functional Requirements

- **FR-001**: S117 specification, required task disposition and inventory MUST consistently record its completed source implementation with actual delivery evidence.
- **FR-002**: Publication, review, merge and branch-cleanup steps MUST appear as workflow evidence rather than unresolved implementation tasks.
- **FR-003**: Reviewed head, merge source, candidate inputs, reports and public package MUST retain distinct truthful identities and outcomes.
- **FR-004**: #355 and #356 MUST retain separate original diagnoses, actual delivered diagnostic improvements and explicitly unconfirmed remaining causes.
- **FR-005**: #354 MUST remain open for its public release/discovery outcome, with completed source children and remaining dependencies reconciled.
- **FR-006**: Pending v1.6.0 release guidance and content bindings MUST reflect delivered source without advertising an unpublished package or changing public v1.5.1 selection.
- **FR-007**: A concrete maintainer handoff MUST distinguish reviewed source preparation, later exact-source tag/staging, separately approved promotion and public discovery reconciliation.
- **FR-008**: Original retained evidence MUST remain unchanged; no runtime, installed-service, account-session, DNS, tag, staging or promotion operation is part of this source slice.
- **FR-009**: Source delivery MUST follow the complete Spec-Kit sequence and existing canonical verification, with current slice identity aligned throughout.
- **FR-010**: Owner-authorized PR publication MUST address all warranted individual feedback within at most two Codex rounds and stop unmerged.

### Key Entities

- **Implementation record**: Feature maturity, required task dispositions and objective source delivery references.
- **Release disposition**: Issue-level delivered outcome, unresolved cause, source reference and publication state.
- **Evidence identity**: Reviewed head, merged source, frozen candidate inputs and unchanged historical report outcome.
- **Release handoff**: Later exact-source, staged-candidate, owner-approval and public-discovery transitions.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All S117 required tasks and lifecycle records have consistent completed implementation dispositions, with zero unresolved publication tasks represented as implementation.
- **SC-002**: Every included issue (#354, #355 and #356) has a current individual functional disposition and no unsupported closure or fix claim.
- **SC-003**: All pending/public release records agree on their actual state and existing evidence, with zero invented tag, asset or publication identities.
- **SC-004**: Every retained candidate/report hash remains unchanged, and the reviewed source PR finishes with successful applicable checks and addressed feedback inside the two-round limit.

## Assumptions

- Existing source and hosted verification remain available; ask for authentication only on an actual access failure.
- Source reconciliation is the authorized S118 increment. Tagging, staging, promotion and merge require later owner action.
- Existing release and evaluation tooling are reused. No application behavior, timeout policy, dependency migration, full candidate replay or replacement QA framework is presumed.
- Any concrete product defect found during this work is tracked separately against its affected source and corrected only within demonstrable scope.
