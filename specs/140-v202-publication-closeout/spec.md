# Feature Specification: v2.0.2 publication closeout

**Feature Branch**: `codex/140-v202-publication-closeout`

**Created**: 2026-11-08

**Status**: Implemented

**Delivery**: Pre-tag source review branch codex/140-v202-publication-closeout; [verification](verification.md), [delivery](delivery.md); owner final review/merge and later actual publication remain pending

**Input**: Owner-approved S140 under [#444](https://github.com/shruggietech/go-schedule/issues/444): finish the v2.0.2 publication continuation, first delivering the missing comparison records through the pre-tag source PR. Push/open the official PR automatically, monitor CI and resolve findings, use at most two total Codex reviews including automatic review, and stop for owner final review/merge. Later exact-source tagging/draft staging, inspected-asset promotion and post-publication discovery remain the same S140 continuation with their existing explicit decisions. Preserve the prohibition on operator-desktop focus and non-admin operator interaction.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review the maintenance candidate accurately (Priority: P1)

A maintainer can inspect a complete v2.0.2 comparison whose product context and evidence describe the prepared maintenance source, with every capability assessed against the selected go-schedule behavior and the seven named vendors.

**Why this priority**: The existing release path requires a reviewed candidate comparison before it can stage the maintenance packages.

**Independent Test**: The candidate comparison contains all eight product contexts, ten capability rows and 80 evidence cells, and each row has an explicit review decision.

**Acceptance Scenarios**:

1. **Given** the accepted S139 maintenance implementation, **When** a maintainer reviews v2.0.2, **Then** its comparison includes those delivered source behaviors and retains the actual limits on native and cross-OS evidence.
2. **Given** current official vendor evidence, **When** a cell is reviewed, **Then** its status, edition/deployment context, explanation, source and review date agree with that evidence.
3. **Given** unavailable or insufficient evidence, **When** the comparison is prepared, **Then** uncertainty remains explicit and unsupported claims are not inferred from missing documentation.

### User Story 2 - Bind release review to the reviewed content (Priority: P1)

A maintainer can distinguish the exact candidate comparison, maintenance changelog and release notes that were assessed, while earlier released comparisons remain unchanged.

**Why this priority**: A changed candidate must not reuse approval for different content or rewrite the history of an existing release.

**Independent Test**: All ten decisions identify changed or no-change outcomes, and the current review corresponds to the exact selected candidate content.

**Acceptance Scenarios**:

1. **Given** a complete candidate dataset, changelog and notes, **When** the review is completed, **Then** it identifies v2.0.2 and binds all three reviewed contents.
2. **Given** an existing v2.0.1 or earlier snapshot/review, **When** S140 prepares the new candidate, **Then** its earlier bytes remain unchanged.
3. **Given** content changed after review, **When** a maintainer checks readiness, **Then** the existing guard rejects the stale review instead of silently accepting it.

### User Story 3 - Continue publication from an honest source checkpoint (Priority: P2)

The owner can merge the pre-tag source work and then follow the remaining release stages without confusing preparation, staging, public availability or installed/native observations.

**Why this priority**: The source PR must not claim an uncreated tag, uninspected package or unpublished release.

**Independent Test**: The continuation names the order, source identity selection, existing approval boundaries and actual completion conditions for the remaining release work.

**Acceptance Scenarios**:

1. **Given** public v2.0.1 packages, **When** this source PR is delivered, **Then** current public selection remains v2.0.1 and v2.0.2 is described as prepared source.
2. **Given** the prerequisite source is owner-merged, **When** publication continues, **Then** the exact accepted main is selected before the existing tag/draft decision, actual asset inspection and separate promotion decision.
3. **Given** actual promotion later completes, **When** discovery is reconciled, **Then** public download readback and current availability guidance refer to the promoted release, while #444 remains open until that outcome exists.

### Edge Cases

- A vendor source moves, is inaccessible, or describes another edition: use an accessible official equivalent where supported; retain explicit uncertainty and the last verified evidence date where carried forward. Do not fabricate a current retrieval or unsupported feature absence.
- Candidate notes or changelog change during preparation: reassess affected capability decisions and regenerate the content identities only after that review.
- No capability status changes: record explicit no-change decisions with a reason; do not invent new functionality to fill a maintenance release.
- A source branch is still unmerged: keep the tag/staging/publication continuation pending and do not select that branch as accepted main.
- No isolated native desktop or actual cross-OS connectivity exists: preserve S139's bounded observations and disclosures; do not request a manual walkthrough, focus a shared desktop, or reconfigure the host.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The prepared v2.0.2 comparison MUST cover the existing eight product identities, ten capability identities and 80 cells with explicit context, status, explanation, official source and evidence date.
- **FR-002**: Every go-schedule cell MUST reflect accepted S139 source behavior, distinguish prepared source from published packages, and preserve unresolved native, cross-OS, historical-incident and clustered-execution limits.
- **FR-003**: All seven named competitors MUST be assessed against official evidence, preserving edition/add-on/deployment distinctions and explicit unknowns when evidence is insufficient.
- **FR-004**: The release review MUST identify v2.0.2, contain exactly one changed or no-change decision for every capability, and bind the exact comparison, changelog section and tag-specific notes that were reviewed.
- **FR-005**: Historical snapshots/reviews, published v2.0.1 package identities and protected S139 evidence MUST remain unchanged.
- **FR-006**: Current public-release selection and availability MUST remain v2.0.1 during this pre-tag source checkpoint; candidate review MUST NOT publish v2.0.2 as available.
- **FR-007**: The release continuation MUST state the pre-tag source merge, exact accepted-main selection, existing tag/draft and inspected-asset promotion decisions, independent public readback and post-publication discovery merge, with #444 remaining open until actual delivery.
- **FR-008**: All work MUST avoid operator-desktop focus, non-admin operator requests, shared native GUI/keyring actions and global host reconfiguration; public evidence MUST omit PII and secrets.
- **FR-009**: Every actionable defect discovered during this slice MUST be deduplicated and tracked with functional criteria under the existing release coordinator; planned comparison preparation MUST remain distinguishable from product defects.

### Key Entities

- **Candidate comparison**: Release identity, review date, product contexts, capability evidence cells and opportunity links for prepared v2.0.2 source.
- **Content-bound review**: Candidate identity, the three exact content identities and one decision for each capability.
- **Release continuation**: Ordered source, staging, inspection, promotion and discovery states, with actual evidence and explicit authority at each transition.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A maintainer can inspect all eight product contexts and 80 evidence cells, including every one of the seven competitor assessments, without an undocumented capability or edition assumption.
- **SC-002**: All ten capabilities have explicit decisions, and the review identifies the exact three reviewed contents without changing a historical release record.
- **SC-003**: Every remaining release stage has an explicit prerequisite and completion outcome; the current source handoff makes zero claims that v2.0.2 packages already exist or are public.
- **SC-004**: Protected evidence remains unchanged and the work causes zero operator-desktop focus acquisition or non-admin operator interaction.

## Assumptions

- The owner-approved S140 continuation starts at accepted main d5bb2632d24bbeb0615b041533acae4866761d28, whose source tree equals the reviewed S139 candidate.
- This run ends at the first official pre-tag source PR ready for final owner review/merge. Tagging, staging, promotion and the later discovery PR stay in the same release continuation and are not unfinished source implementation tasks.
- The repository's existing comparison, release, standing-waiver, CI and review machinery remains the implementation path. No new framework, dependency, scheduling feature or verification policy is needed.
- Original #355/#356 await matching evidence, and the owner-deferred clustered lane remains unchanged. The current maintenance comparison does not authorize those features.

## Clarifications

### Session 2026-11-08

No critical ambiguities detected. The owner accepted the ordered two-source-checkpoint S140 proposal and explicitly authorized this pre-tag PR's publication/CI continuation. Functional scope, data, interaction, quality, dependencies, failure handling, constraints, terminology and completion taxonomy are clear; no question is required.

## Bounded implementation discovery

Deduplicated #453 removes the cron guide's obsolete unreleased-after-0.8.0 notice for functionality already shipped in v2.0.1. This documentation correction follows FR-009 and does not change cron behavior, supported dialects or fidelity. The pre-tag PR individually delivers #452, #453 and the pending-source highlight-contract correction #454 while #444 remains open.
