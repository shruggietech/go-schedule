# Feature Specification: S125 Reliability Delivery and Scheduling Decision

**Feature Branch**: `codex/125-reliability-delivery-decision`
**Created**: 2026-10-28
**Status**: Implemented
**Delivery**: Source outcomes implemented; controlled browser regression and all 35 browser cases passed, all eight canonical gates passed in 933.338 seconds with zero source/evidence drift; see [delivery.md](delivery.md) and [verification.md](verification.md). Release publication and conditional issue outcomes remain explicitly incomplete.
**Input**: Owner-authorized S125 with #397, #394, #186, #189 and #20; #187/#188 only as their gates permit. Push and open one PR, monitor CI and at most two Codex review rounds. Stop before merging, tagging or promotion.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Reliable keyboard search and diagnosable failure (Priority: P1)

A desktop user moves to Search and enters a query without an asynchronous route-focus handoff interrupting the interaction. Maintainers retain useful diagnostics when the production browser contract fails.

**Why this priority**: The merged-main browser contract failed and blocks release-source selection (#397).
**Independent Test**: Exercise a deliberately delayed route-focus handoff, enter archive with the keyboard, and retain the existing compact-viewport, zoom, filter and overflow expectations.
**Acceptance Scenarios**:

1. **Given** route navigation has a pending focus handoff, **When** the search interaction begins, **Then** route focus settles before query entry, and the exact query remains intact.
2. **Given** a failed browser exercise, **When** its runner exits, **Then** failure context and trace remain retrievable and the original failure status remains failed.
3. **Given** normal navigation, **When** Search opens, **Then** accessible main-content focus and exact daemon identity behavior remain unchanged.

### User Story 2 - One release delivery path (Priority: P1)

A maintainer can take the prepared v1.7.0 reliability source through the existing release workflow without reconstructing release requirements or confusing source preparation with published assets.

**Why this priority**: #394 remains open after S124 source preparation.
**Independent Test**: Follow the release handoff from this source proposal through the exact-source approval boundaries using documented existing commands; all unreleased/public identities remain distinguishable.
**Acceptance Scenarios**:

1. **Given** the owner-merged S124 source and failed main browser check, **When** S125 reconciles release records, **Then** the failure is tracked independently and that red source is not selected for tagging.
2. **Given** this source proposal remains unmerged, **When** its delivery status is reported, **Then** the next exact main source and staged assets are explicitly unresolved and public discovery remains v1.6.0.
3. **Given** future owner approval, **When** delivery proceeds after merge, **Then** the same existing tag/draft, asset approval, promotion and discovery path can finish #394 without another preparation-only feature.

### User Story 3 - Findings-based distributed scheduling direction (Priority: P2)

The owner receives a concrete assessment of unmet clustered-execution need and a bounded decision record covering #186, #189 and #20, with #187/#188 activated only by their documented gates.

**Why this priority**: Five related conditional issues should be handled as one decision chain rather than fragmented future slices.
**Independent Test**: Trace every claimed scenario, user, scale and tradeoff to available evidence; verify that each conditional child has an explicit proceed/hold/stop disposition.
**Acceptance Scenarios**:

1. **Given** supported concrete unmet scenarios, **When** the demand recommendation authorizes exploration, **Then** failure semantics precede approved technology/prototype work.
2. **Given** insufficient real-use evidence, **When** the assessment is written, **Then** unknowns remain explicit, a compact intake describes exactly what is missing, and architecture/prototype gates stay closed.
3. **Given** a decision requires owner judgment, **When** the source PR is handed off, **Then** a bounded recommendation and alternatives are reviewable without silently closing or cancelling the product proposal.
4. **Given** a documented stop before prototype work, **When** #189 records the disposition, **Then** skipped architecture/prototype work is never described as completed implementation.

### Edge Cases

- Delayed animation-frame focus occurs after the initial navigation click completes.
- Browser failure creates no trace because startup failed; artifact collection warns without hiding the original exit.
- The owner merges source but has not authorized an immutable tag or actual assets; no release mutation follows automatically.
- Search current/closed issues reveal aspiration but no actual workload evidence; ambition is preserved without fabricated user counts.
- Historical #355/#356 repeat successfully; original causal outcomes remain unresolved.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Search interaction MUST synchronize the authored route-focus handoff before keyboard query entry while preserving accessible navigation, exact query, viewport/zoom, filter and overflow behavior (#397).
- **FR-002**: Failed browser exercises MUST leave retrievable bounded failure context and traces without replacing their failure status (#397).
- **FR-003**: Release records MUST reconcile S124 owner merge, its independent red main receipt and the S125 correction; prepared v1.7.0 source MUST remain distinct from public v1.6.0 and future tagged assets (#394).
- **FR-004**: The release handoff MUST describe the complete existing delivery sequence, exact identity selection and owner approval boundaries; this slice MUST stop before merge/tag/promotion and MUST NOT imply #394 is complete.
- **FR-005**: The demand assessment MUST identify actual supporting sources, concrete scenarios, unmet alternatives, bounded scale and failure tradeoffs when available, and explicitly distinguish unknowns from findings (#186).
- **FR-006**: The scheduling decision MUST cover continue, hold and stop alternatives, value/correctness/operability/support consequences and next actions without inventing demand or unapproved product restrictions (#189/#20).
- **FR-007**: #187 MUST proceed only after a positive #186 recommendation; #188 MUST proceed only after approved semantics. A failed gate MUST leave an explicit reason and issue-level incomplete outcome.
- **FR-008**: GitHub and source records MUST preserve issue-level traceability, original #355/#356 causal uncertainty, and shipped-versus-planned claims. Only individually delivered outcomes may use closing references.
- **FR-009**: Existing CI completion allowances, identity/security assertions and release guards MUST remain effective. New owned Windows launches MUST remain hidden, noninteractive and bounded.

### Key Entities

- **Release source**: reviewed branch head, owner merge, exact main acceptance and public/tag identities, each recorded separately.
- **Demand finding**: source, actual actor/workload, scale, unmet alternative, failure tradeoff, confidence and missing inputs.
- **Decision gate**: prerequisite issue, decision, authority, allowed next action and explicit hold/stop reason.
- **Browser failure receipt**: original source, failure context, trace and status, separate from corrected-source evidence.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The delayed-focus search scenario retains the complete archive query and all existing keyboard/zoom expectations.
- **SC-002**: Every failed browser exercise with generated diagnostics has a retrievable artifact without changing its failure outcome.
- **SC-003**: One handoff covers every remaining release-delivery step and identifies all three owner boundaries (merge, exact tag/draft and concrete-asset promotion).
- **SC-004**: All seven scoped issues have an explicit delivered, partial or gated disposition with zero unsupported closure or shipped-capability claims.
- **SC-005**: Every distributed-demand assertion names its evidence or an explicit missing input; no unsupported user, scale or availability commitment is introduced.

## Assumptions

- The owner has explicitly authorized push/PR and up to two Codex rounds, superseding the skill's default pre-push pause.
- This request authorizes assessment and routine decisions, not an unsupported change to the existing distributed-scheduling gates.
- The source fix changes browser fixture/diagnostic behavior; reuse existing assembled application QA unless a genuine application change justifies fresh application QA.
- Public release remains v1.6.0 until separately approved publication; original Windows incidents remain independently tracked.

## Clarifications

### Session 2026-10-28

- Q: Does publication authorization permit tagging or merging? A: No. Push/PR and in-scope review fixes are authorized; merge, tag and promotion remain excluded.
- Q: Does missing demand justify rejecting clustered execution? A: No. Assess available evidence, record a hold when insufficient and preserve owner intent and open outcomes.
- Q: Must unchanged application QA be repeated for a browser fixture correction? A: No. Preserve prior assembled QA identities and run meaningful focused browser regressions plus canonical source gates.

Coverage scan: functional scope, entities, interaction, quality attributes, dependencies, edge cases, constraints, terminology and completion are clear. No critical ambiguity requires an operator answer to perform the authorized source work.
