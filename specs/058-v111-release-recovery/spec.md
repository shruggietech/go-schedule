# Feature Specification: v1.1.1 Release Recovery

**Feature Branch**: `codex/058-v111-release-recovery`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: v1.1.1 release recovery preparation and all eight canonical verification gates completed on `codex/058-v111-release-recovery`; publication operations remain tracked by issue [#140](https://github.com/shruggietech/go-schedule/issues/140).

**Input**: The maintainer requested S058 after the v1.1.0 release attempt exposed a production watcher defect and an exact-commit CI gap. PR #143 corrected both defects, but the unpublished v1.1.0 tag and draft must not be mistaken for the qualified public release.

## Clarifications

### Session 2026-09-05

- Q: Should the existing v1.1.0 tag be moved to the corrected commit? -> A: No. Preserve the immutable tag as unpublished history and use v1.1.1 for the corrected public release.
- Q: What happens to the unpublished v1.1.0 draft? -> A: Retire only the draft during post-merge release operations so no unqualified v1.1.0 artifacts remain available for publication.
- Q: Should corrective release tracking duplicate issue #140? -> A: No. Rename and update the existing release coordinator and milestone to v1.1.1 while retaining links to the interrupted v1.1.0 attempt.
- Q: How much release copy should be published? -> A: Four concise highlights followed by one tagged full-changelog link.

## User Scenarios & Testing

### User Story 1 - Recover without rewriting history (Priority: P1)

A maintainer can recover from the interrupted v1.1.0 attempt without moving its tag or publishing its unqualified draft.

**Why this priority**: Immutable release identity and the absence of unsafe public artifacts are the foundation for every later operation.

**Independent Test**: Verify the old tag still resolves to its original commit, the old release remains unpublished until its draft is retired, and all new source identity targets v1.1.1.

**Acceptance Scenarios**:

1. **Given** the existing v1.1.0 tag, **When** recovery begins, **Then** the tag remains unchanged and no public v1.1.0 release is created.
2. **Given** the existing unpublished v1.1.0 draft, **When** v1.1.1 release operations begin after merge, **Then** only that draft is retired and the historical tag remains intact.

### User Story 2 - Understand the corrected release quickly (Priority: P1)

A person opening the v1.1.1 release can identify the principal feature outcomes and watcher reliability correction immediately, then follow one link for the complete record.

**Why this priority**: v1.1.1 becomes the first public release of the accumulated v1.1 capability set and must explain that outcome without lengthy copy.

**Independent Test**: Validate that the tag-specific note contains exactly four one-line highlights and ends with one tagged changelog link.

**Acceptance Scenarios**:

1. **Given** the v1.1.1 release page, **When** a reader scans the notes, **Then** four concise outcomes cover the feature set and corrected filesystem watching.
2. **Given** a reader wants complete detail, **When** the reader follows the final link, **Then** the tagged v1.1.1 changelog section opens.

### User Story 3 - Publish one trustworthy corrected artifact set (Priority: P1)

A user downloads v1.1.1 and receives artifacts built from the reviewed corrective merge only after CI for that exact commit succeeds.

**Why this priority**: The previous attempt proved that branch-level success is insufficient when the release source commit differs.

**Independent Test**: Stage from one immutable v1.1.1 tag, confirm the Release workflow found successful main-branch CI for the exact tag commit, qualify the exact Windows candidate, promote without rebuilding, and verify every public checksum.

**Acceptance Scenarios**:

1. **Given** reviewed release-recovery source, **When** v1.1.1 is staged, **Then** the tag, CI run, packages, and candidate manifest identify the same commit.
2. **Given** a qualified draft, **When** it is promoted, **Then** the workflow publishes the existing artifacts without rebuilding them.

### Edge Cases

- A moved or recreated v1.1.0 tag stops recovery because immutable history no longer matches the recorded attempt.
- A public v1.1.0 release stops automatic draft retirement and requires explicit maintainer reconciliation.
- An existing local tag, remote tag, draft, or public v1.1.1 release stops tag creation until identity is reconciled.
- A missing, pending, failed, cancelled, or timed-out CI run for the exact v1.1.1 tag commit prohibits staging.
- Missing, extra, empty, changed, or checksum-mismatched assets prohibit publication.

## Requirements

### Functional Requirements

- **FR-001**: The v1.1.0 tag MUST remain at its original commit and MUST NOT be moved, deleted, or recreated by recovery.
- **FR-002**: The unpublished v1.1.0 draft MUST remain non-public and MUST be retired during post-merge release operations before v1.1.1 promotion.
- **FR-003**: Issue #140 and its milestone MUST become the authoritative v1.1.1 release record without creating duplicate corrective-release tracking.
- **FR-004**: The repository MUST contain an empty Unreleased section followed by dated v1.1.1 and v1.1.0 changelog sections, with the corrective changes recorded under v1.1.1.
- **FR-005**: The README health example MUST identify version 1.1.1 before tag creation.
- **FR-006**: The tag-specific v1.1.1 release note MUST contain exactly four concise one-line highlights and MUST end with one link to the tagged v1.1.1 changelog section.
- **FR-007**: The annotated v1.1.1 tag MUST be created only at the reviewed S058 merge commit after confirming clean synchronized main and absent conflicting v1.1.1 state.
- **FR-008**: The Release workflow MUST locate and require successful main-branch CI for the exact v1.1.1 tag commit before building or uploading artifacts.
- **FR-009**: The exact staged Windows MSI MUST pass the applicable candidate and attended qualification gates before promotion.
- **FR-010**: Promotion MUST verify the staged artifact set, qualification evidence, tag identity, and final checksum inventory without rebuilding artifacts.
- **FR-011**: The final audit MUST verify public, latest, tag, commit, asset, checksum, note, changelog, README, and binary version consistency for v1.1.1.
- **FR-012**: The canonical automation gate MUST execute the approved-workflow fixture regression suite so validator contract changes cannot leave a stale fixture undetected.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The original v1.1.0 tag commit remains unchanged and no public v1.1.0 release exists.
- **SC-002**: Release notes contain four highlights and one final tagged full-changelog link, with no additional prose.
- **SC-003**: All eight local verification gates and every required hosted check pass on the reviewed S058 commit.
- **SC-004**: Every staged and public v1.1.1 artifact is non-empty and represented exactly once in the final checksum inventory.
- **SC-005**: The public release, tag, main boundary, README, changelog, manifest, and sampled binaries all identify v1.1.1.

## Assumptions

- v1.1.0 was never published and its existing GitHub release remains a draft.
- PR #143 contains the complete watcher and exact-commit CI corrections required before another release attempt.
- The existing Release and Promote Release workflows remain the authoritative artifact pipeline.
- One concise attended Windows confirmation may still be necessary because native installation behavior cannot be proven completely by headless CI.
