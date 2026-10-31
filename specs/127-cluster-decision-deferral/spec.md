# Feature Specification: S127 Clustered-Execution Deferral

**Feature Branch**: `codex/127-cluster-decision-deferral`

**Created**: 2026-10-31

**Status**: Implemented

**Delivery**: On review branch codex/127-cluster-decision-deferral, owner-directed deferral and all six native open dispositions recorded; full eight-gate verification, 10 site browser cases and zero source/evidence drift passed. See delivery.md and verification.md; PR publication is separately authorized and merge remains excluded.

**Input**: Owner-authorized S127 under Spec-Kit/autopilot, automatic push and one official PR, CI and external review handling, at most two Codex review rounds, stop before merge. The owner selected: "Explicitly defer clustered execution."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand the actual product decision (Priority: P1)

An operator or maintainer can distinguish the owner's explicit deferral from the earlier evidence-specific hold, a permanent rejection, and a delivered clustered capability.

**Why this priority**: S125 documented missing demand inputs but could not choose the owner's product direction. The owner has now chosen deferral.

**Independent Test**: Read the current decision page and S127 decision record; both identify the owner direction, date, retained capability boundary and unfinished outcomes without requiring historical chat.

**Acceptance Scenarios**:

1. **Given** the owner selected explicit deferral, **When** the current assessment is read, **Then** it records that direction as the current disposition and preserves the earlier S125 assessment as history.
2. **Given** the deferral, **When** issue outcomes are reported, **Then** #186-#189, #20 and #146 remain open for their unfinished outcomes; none is represented as completed research, architecture, prototype or final go/no-go.
3. **Given** shipped v1.7.0, **When** capabilities are described, **Then** independent-daemon control remains distinct from shared execution ownership and host failover.

### User Story 2 - Resume deliberately when direction changes (Priority: P2)

A maintainer has one clear reactivation path without a recurring reassessment or speculative implementation.

**Why this priority**: A reversible deferral needs an attributable way to resume, rather than an invented date or an implicit permanent restriction.

**Independent Test**: Trace a future workload report or explicit owner revision through the retained issue dependencies and identify the next allowed stage.

**Acceptance Scenarios**:

1. **Given** a later concrete workload report, **When** the owner revisits the deferral, **Then** #186 assesses its ownership/recovery gap, scale, failure tradeoffs and alternatives before recommending #187.
2. **Given** continued exploration, **When** technology or prototype work is proposed, **Then** #187's maintainer semantics approval still precedes #188.
3. **Given** no new direction, **When** maintenance continues, **Then** no deadline, automatic resumption, scheduled check, implementation roadmap or new eligibility restriction is inferred.
4. **Given** an owner revision of the existing gates, **When** it is recorded, **Then** its changed scope and authority are explicit; this kickoff alone supplies no such revision.

### Edge Cases

- New aspirations or elapsed time do not become actual workload evidence.
- Deferral does not settle a final no-go or assert that alternatives meet every workload.
- An issue remains open even though its parent release #394 is delivered.
- The owner may later stop the proposal or revise its gates; record that separate decision when made.
- Public release records remain immutable; no new runtime or release behavior is claimed.
- Sensitive workload details may stay private with an attributable publishable summary.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Current documentation MUST record the owner's explicit 2026-10-31 deferral and distinguish it from S125's insufficient-evidence hold and a final stop decision (#20/#189).
- **FR-002**: Source and native planning records MUST assign each of #186, #187, #188, #189, #20 and #146 a truthful open/deferred or prerequisite-gated disposition without claiming unmet functional outcomes are delivered.
- **FR-003**: The reactivation path MUST require a later owner revisit or explicit revision; actual workload evidence is assessed under #186 and does not automatically resume deferred work.
- **FR-004**: The existing sequence MUST remain #186 recommendation, #187 semantics and maintainer approval, #188 isolated unshipped prototype, then #189 final decision; no backend, credentials, support scale or release schedule is selected.
- **FR-005**: Current guidance MUST distinguish shipped independent-daemon functionality from deferred clustered execution and describe the existing workload intake and alternatives without duplicating the intake or inventing demand.
- **FR-006**: S127 MUST preserve historical records, release identities, source/runtime behavior, existing security assertions and frozen evidence. #355/#356 and #400/#404 MUST retain their independent unresolved outcomes.
- **FR-007**: Source, branch, feature selector, issue mapping, delivery records and PR MUST consistently identify S127; publication/review bookkeeping MUST remain separate from implementation tasks and functional issue completion.

### Key Entities

- **Owner disposition**: attributable direction, date, exact choice, scope and reversible status.
- **Issue disposition**: existing issue, unmet outcome, allowed next stage, retained dependency and open state.
- **Reactivation input**: owner revisit/revision and any supplied workload evidence, with explicit unknowns.
- **Delivery receipt**: verified source and publication/review identities, separate from product or research outcomes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One current decision page and one S127 record express the same owner-directed deferral, with zero claims of final rejection or shipped cluster behavior.
- **SC-002**: All six scoped issues have explicit, consistent dispositions and remain open; none acquires an unsupported closing reference.
- **SC-003**: The reactivation instructions identify all four existing stages and the maintainer approval boundary, with no invented date or automatic resumption.
- **SC-004**: All current claims trace to owner direction, existing issue requirements or shipped source; missing workload data remains unknown.
- **SC-005**: Retained evidence and published release identities remain unchanged.

## Assumptions

- The owner's direction authorizes the deferral itself; it does not authorize an implementation roadmap, gate removal or permanent no-go.
- No runtime, storage, workflow or dependency change is needed to deliver this planning outcome.
- Existing labels/milestones can describe these records; adding a new tracking framework or duplicate issues is unnecessary.
- Push, PR and verified in-scope review-fix pushes are explicitly authorized; merge, tags and releases remain excluded.

## Clarifications

### Session 2026-10-31

- Q: Which product direction should S127 pursue? A: The owner selected "Explicitly defer clustered execution."
- Q: Does deferral complete #186 or #189? A: No. It records current direction while preserving their missing evidence and final-decision outcomes.
- Q: Does a workload report automatically restart the effort? A: No. It supplies useful input for an owner revisit; retained dependency and semantics-approval gates still apply.
