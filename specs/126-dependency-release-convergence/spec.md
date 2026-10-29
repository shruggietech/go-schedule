# Feature Specification: S126 Dependency and Release Convergence

**Feature Branch**: `codex/126-dependency-release-convergence`

**Created**: 2026-10-29

**Status**: Implemented

**Delivery**: Selected dependency convergence, source-record reconciliation and release continuation implemented; fresh Windows/WSL QA, both module graphs, 45 browser cases and all eight canonical gates passed. See verification.md and delivery.md; actual #394 publication remains beyond owner merge/asset approvals.

**Input**: Owner-authorized S126, following #394's recommendation: converge #401-#403, consider #400/#404 only with a suitable published driver, qualify the assembled candidate, and publish one reviewed source PR. Stop for owner review before merging, tagging or promotion.

## User Scenarios & Testing

### User Story 1 - Consume compatible dependency updates (Priority: P1)

Maintainers can build the scheduler, desktop and documentation with the selected compatible dependency updates while retaining current scheduling, configuration and durable-state behavior.

**Why this priority**: Ready maintenance proposals should deliver together without waiting for an unsuitable storage target.

**Independent Test**: Restore both application modules and frontend/site lockfiles from the assembled candidate and exercise the existing application and site contracts.

**Acceptance Scenarios**:

1. **Given** proposals #401-#403, **When** their selected updates are integrated, **Then** all affected manifests and lockfiles describe one coherent supported candidate.
2. **Given** #400 targets a driver with a documented regression, **When** no suitable corrected release is published, **Then** the existing driver is retained and #400/#404 remain deferred without blocking the ready updates.
3. **Given** a suitable corrected published driver, **When** adopted, **Then** both application modules consume its complete intended graph without rewriting metadata, resolving #404.

### User Story 2 - Inspect truthful release readiness (Priority: P1)

Maintainers can see the actual selected dependencies, application results and source provenance for the prepared reliability release, while public users continue to see the release that is actually published.

**Why this priority**: #394 needs delivery continuity without premature public capability or package claims.

**Independent Test**: Inspect release notes, content-bound comparison review and source records against the selected dependency identities and actual publication state.

**Acceptance Scenarios**:

1. **Given** S125 is merged with successful main checks, **When** source records are reconciled, **Then** they identify that actual merge and distinguish its acceptance from S126 candidate acceptance.
2. **Given** S126 changes runtime dependencies, **When** candidate QA runs, **Then** its report identifies the new candidate and preserves prior failed evidence and unresolved #355/#356 causes.
3. **Given** no v1.7.0 assets have been approved or published, **When** source preparation finishes, **Then** public discovery remains v1.6.0 and #394 remains open for actual delivery.

### User Story 3 - Complete source review within explicit boundaries (Priority: P1)

The owner receives one source PR with addressed CI and reviewer feedback and an executable continuation to real release delivery after owner merge.

**Why this priority**: One coherent delivery path avoids another preparation-only planning cycle.

**Independent Test**: Audit the exact PR head, terminal checks, every review comment and review receipts against the two-round cap and publication boundary.

**Acceptance Scenarios**:

1. **Given** source publication is explicitly authorized, **When** local verification succeeds, **Then** the branch is pushed and a PR opened without another permission halt.
2. **Given** reviewers leave comments, **When** feedback is resolved, **Then** warranted corrections and explanations are recorded, with at most two Codex rounds.
3. **Given** the final PR is ready, **When** the owner is notified, **Then** the PR remains unmerged and tags, draft staging and promotion await their separate approvals.

### Edge Cases

- A changelog documents a patch that the module registry does not publish: retain the supported current driver, record the observed identities and defer the update.
- Root-only changes leave the desktop graph stale: regenerate both modules rather than permit build-time rewrites.
- Individual proposal checks pass but the combined tree fails: correct the combined failure and preserve first evidence.
- New review feedback after the second round: resolve it within scope without requesting a third round.
- The owner merges before monitoring finishes: stop PR review monitoring and do not infer release authorization.

## Requirements

### Functional Requirements

- **FR-001**: Integrate the compatible selected changes from #401-#403 as one coherent dependency set with complete lockfiles.
- **FR-002**: Adopt #400 only if a suitable corrected driver is published; otherwise preserve the current driver and explicitly defer #400/#404.
- **FR-003**: If storage adoption proceeds, align both module graphs and preserve stored state, configuration, cancellation and audit contracts (#404).
- **FR-004**: Preserve all existing product interfaces and scheduler capability claims; no new clustered scheduling or storage architecture is introduced.
- **FR-005**: Reconcile actual S125 merge/source acceptance and accurately disclose S126 dependency selection in pending reliability-release records.
- **FR-006**: Keep public discovery at the actually published release until separately approved promotion; retain original #355/#356 causal uncertainty.
- **FR-007**: Publish one source PR and address every review comment within a maximum of two Codex rounds; stop before merge/tag/release mutations.
- **FR-008**: Maintain one concrete release continuation through actual merged-source selection, exact tag/draft approval, inspection of real assets, separate promotion approval and discovery reconciliation.

### Key Entities

- **Dependency selection**: proposal, observed version, selected or deferred outcome, complete manifests and rationale.
- **Candidate acceptance**: assembled source identity, QA and verification receipts distinct from earlier source and public artifacts.
- **Release continuation**: accepted main, approved immutable source, real release/run/asset identities and separately authorized transitions.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All three ready proposals are represented in one coherent source delivery; each conditional proposal has an explicit disposition.
- **SC-002**: Existing task lifecycle, interface, configuration and persistence outcomes remain available on the assembled candidate.
- **SC-003**: Release records identify selected changes and actual publication state without inventing an asset or source identity.
- **SC-004**: The owner receives one ready source PR after no more than two Codex rounds and retains the merge decision.

## Assumptions

- Existing release and unattended QA infrastructure is reused. No new release framework or unchanged historical repeat campaign is needed.
- #394's actual publication outcome remains incomplete at the source-review boundary; source PR integration cannot close it.
- Tag/draft and concrete-asset promotion approvals from earlier candidates do not authorize a new candidate.
- No new dependency proposal is silently absorbed after source selection; material scope expansion is recorded and assessed.

## Clarifications

### Session 2026-10-29

- Q: How far does this kickoff authorize delivery? A: Implement, push, open the source PR and resolve CI/reviews; stop for the owner's merge ritual. Subsequent exact tag/draft and asset promotion approvals remain distinct.
- Q: Should an unavailable SQLite fix delay the ready set? A: No. Preserve the current driver and keep #400/#404 visibly open and conditional.
