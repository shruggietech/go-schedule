# Feature Specification: v1.1.0 Release

**Feature Branch**: `codex/057-v110-release`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Release preparation and all eight canonical verification gates completed on `codex/057-v110-release`; public release operations remain tracked by issue [#140](https://github.com/shruggietech/go-schedule/issues/140).

**Input**: The maintainer authorized a new release with minimal interaction and requested very brief release notes ending with a link to the broader changelog.

## Clarifications

### Session 2026-09-05

- Q: Which semantic version fits the accumulated work? -> A: v1.1.0, because the changes add backward-compatible capabilities after v1.0.0.
- Q: How much release copy should be published? -> A: Four concise highlights followed by one tagged full-changelog link.
- Q: May release operations proceed autonomously? -> A: Yes, through every safe automated step, with user input requested only for an unavoidable attended candidate check.
- Q: Does the release include the remaining roadmap epics? -> A: No, issues #18, #19, and #20 remain deferred and outside the release boundary.

## User Scenarios & Testing

### User Story 1 - Understand the release quickly (Priority: P1)

A person opening the v1.1.0 release can identify its main outcomes immediately and follow one link for the complete record.

**Why this priority**: Release notes are the public entry point and the maintainer explicitly requested brevity.

**Independent Test**: Validate that the tag-specific note contains exactly four one-line highlights and ends with one tagged changelog link.

**Acceptance Scenarios**:

1. **Given** the v1.1.0 release page, **When** a reader scans the notes, **Then** the four main outcomes are visible without implementation detail or duplicated changelog prose.
2. **Given** a reader wants complete detail, **When** the reader follows the final link, **Then** the tagged v1.1.0 changelog section opens.

### User Story 2 - Install one trustworthy artifact set (Priority: P1)

A user downloads v1.1.0 and receives version-consistent platform artifacts with a complete checksum inventory.

**Why this priority**: Public release identity and artifact integrity are inseparable from the release itself.

**Independent Test**: Stage from one immutable tag, validate the exact Windows candidate, promote without rebuilding, download every public asset, and verify every checksum.

**Acceptance Scenarios**:

1. **Given** reviewed release-preparation source, **When** v1.1.0 is staged, **Then** every artifact and candidate manifest identifies the same tag and commit.
2. **Given** a qualified draft, **When** it is promoted, **Then** the workflow publishes the existing assets and a checksum file covering every payload.

### User Story 3 - Preserve an auditable release boundary (Priority: P2)

A maintainer can reconstruct why v1.1.0 was cut, what it contains, and which automated or attended evidence authorized publication.

**Why this priority**: Immutable tags and public artifacts require durable provenance.

**Independent Test**: Compare the tag, main commit, workflow run, manifest, release metadata, changelog, README example, and downloaded binary versions.

**Acceptance Scenarios**:

1. **Given** a completed release, **When** its identity is audited, **Then** every source and artifact surface reports v1.1.0 and the release is marked latest.

### Edge Cases

- An existing local tag, remote tag, draft, or public v1.1.0 release stops tag creation until its identity is reconciled.
- A failed staging job leaves the release draft and prohibits promotion.
- A changed candidate byte invalidates prior qualification and cannot be relabeled as the same release.
- Missing, extra, empty, or checksum-mismatched assets prohibit publication.
- A release-note link that targets HEAD or another tag fails the publication-format gate.

## Requirements

### Functional Requirements

- **FR-001**: The repository MUST contain an empty Unreleased section followed by a dated v1.1.0 changelog section containing the complete current Unreleased record without content loss.
- **FR-002**: The README health example MUST identify version 1.1.0 before tag creation.
- **FR-003**: The tag-specific release note MUST contain exactly four concise one-line highlights and MUST end with one link to the tagged v1.1.0 changelog section.
- **FR-004**: The release preparation MUST preserve the existing draft-only staging, immutable candidate, checksum, and no-rebuild promotion controls.
- **FR-005**: The annotated v1.1.0 tag MUST be created only from the reviewed release-preparation merge commit after confirming clean synchronized main and absent conflicting tag or release state.
- **FR-006**: The Release workflow MUST complete successfully and leave one draft containing the expected platform packages and Windows candidate manifest.
- **FR-007**: The exact staged Windows MSI MUST pass the applicable candidate and attended qualification gates before promotion.
- **FR-008**: Promotion MUST verify the staged artifact set, qualification evidence, tag identity, and final checksum inventory without rebuilding artifacts.
- **FR-009**: The final audit MUST verify public, latest, tag, commit, asset, checksum, note, changelog, README, and binary version consistency.
- **FR-010**: Deferred issues #18, #19, and #20 MUST remain outside the v1.1.0 release scope.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Release notes contain four highlights and one final full-changelog link, with no additional prose.
- **SC-002**: All eight local verification gates and every required hosted check pass on the reviewed release-preparation commit.
- **SC-003**: Every staged and public artifact is non-empty and represented exactly once in the final checksum inventory.
- **SC-004**: The public release, tag, main boundary, README, changelog, manifest, and sampled binaries all identify v1.1.0.
- **SC-005**: No deferred roadmap issue is closed or advertised as shipped by this release.

## Assumptions

- The accumulated S051 through S056 changes are backward-compatible and warrant a minor release.
- The existing Release and Promote Release workflows remain the authoritative artifact pipeline.
- One concise attended Windows confirmation may still be necessary because native installation and appearance cannot be proven from headless CI alone.
