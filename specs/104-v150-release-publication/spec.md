# Feature Specification: Publish v1.5.0

**Feature Branch**: `codex/104-v150-release`
**Created**: 2026-09-23
**Status**: In Progress
**Delivery**: Reviewed source metadata correction and all eight local verification gates are complete on `codex/104-v150-release`; public v1.5.0 publication remains pending.
**Input**: S104 kickoff and [release issue #185](https://github.com/shruggietech/go-schedule/issues/185).

## User Scenarios & Testing

### User Story 1 - Download the reviewed release (Priority: P1)

As a user, I can obtain v1.5.0 for my supported platform from the latest public release and verify that the downloads belong to one reviewed source revision.

**Why this priority**: A release is only complete when the intended binaries and integrity information are publicly available together.
**Independent Test**: Inspect the public release, its source tag, platform downloads, and checksum inventory.

**Acceptance Scenarios**:

1. **Given** reviewed main source and a completed release build, **when** v1.5.0 is published, **then** the latest-release page resolves to v1.5.0 with all intended platform artifacts.
2. **Given** a downloaded asset, **when** its digest is compared with the published inventory, **then** its content matches the exact staged asset.

### User Story 2 - Understand the supported boundary (Priority: P1)

As a user, I can read release notes and documentation that describe the connected control center honestly, including what is deferred and what has not been established by native testing.

**Why this priority**: Users must not infer clustering, SMTP, or closed-app popup delivery from a feature release.
**Independent Test**: Compare public notes and tagged documentation with delivered capabilities and the recorded qualification disposition.

**Acceptance Scenarios**:

1. **Given** the public v1.5.0 release, **when** a user reads the release notes, **then** the notes identify supported operations and do not claim clustered execution, SMTP delivery, or closed-app desktop popups.
2. **Given** any native checks that remain untested under a maintainer-authorized waiver, **when** the release is published, **then** those checks are disclosed as untested rather than passed.

### Edge Cases

- A draft may have some but not all assets while platform builds are still running. It must not be promoted during this state.
- A tag or candidate manifest that names a different source revision must stop publication.
- A failed automated check is not converted into a waiver or hidden by the release notes.
- The maintainer's standing waiver permits publication with exact-candidate attended Windows checks explicitly untested; it does not excuse any failed automated check.

## Requirements

### Functional Requirements

- **FR-001**: Publication MUST use one reviewed main commit with successful exact-commit CI and an immutable v1.5.0 tag identity.
- **FR-002**: The release MUST include all supported daemon, CLI, desktop, Windows installer, candidate-manifest, and checksum assets from that tag without rebuilding an asset after staging.
- **FR-003**: The public notes and tagged documentation MUST identify only delivered capabilities and explicitly defer SMTP and clustered execution.
- **FR-004**: The standing native-testing waiver MUST be attributable to the maintainer and disclosed for this release without representing an untested check as passed.
- **FR-005**: The final public release MUST be the repository's latest release, and issue #185 and the v1.5.0 milestone MUST be reconciled only after functional publication is confirmed.

### Key Entities

- **Release source**: The exact reviewed main commit and its v1.5.0 tag.
- **Staged asset set**: Platform packages and candidate manifest associated with the tag.
- **Publication record**: Public release notes, asset inventory, checksums, and qualification disposition.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The latest-release page provides v1.5.0 and every intended platform download from one tagged revision.
- **SC-002**: Every publicly downloadable binary or manifest has one matching checksum entry, and a fresh download validates against it.
- **SC-003**: No public v1.5.0 claim promises clustered execution, SMTP delivery, or a popup after the GUI is closed.
- **SC-004**: GitHub issue #185 closes only after the public release exists and its functional publication criteria are met.

## Assumptions

- Issue #185 is the authoritative release outcome; SMTP #176 and coordinator #19 remain planned in a future milestone.
- The maintainer clarified during S104 that the native-testing waiver is universal across releases, not a v1.5.0-specific grant. Automated CI and release artifact-integrity checks remain required by the existing pipeline.
- The existing staged-asset workflow and exact-candidate validator remain authoritative for source and artifact identity.
