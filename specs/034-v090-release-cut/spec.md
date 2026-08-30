# Feature Specification: v0.9.0 Release Cut

**Feature Branch**: `codex/034-v090-release-cut`

**Created**: 2026-08-30

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: review branch `codex/034-v090-release-cut`; local verification
completed 2026-08-30

**Input**: User request to cut a release after the accumulated S018 through S033
work, with GitHub release notes limited to highlights and a link to the complete
changelog.

**Tracking**: GitHub issue
[#79](https://github.com/shruggietech/go-schedule/issues/79)

## User Scenarios & Testing

### User Story 1 - Understand the release quickly (Priority: P1)

As someone evaluating the new release, I can scan a short set of meaningful
highlights and follow one link when I want the complete change history.

**Why this priority**: The maintainer explicitly wants the release page to be a
useful summary rather than a generated inventory of every merged change.

**Independent Test**: Inspect the prepared v0.9.0 release notes and establish
that they contain only selected highlights plus one link to the complete
versioned changelog.

**Acceptance Scenarios**:

1. **Given** the v0.9.0 release page, **When** a reader scans the notes, **Then**
   the page presents a concise Highlights section rather than an exhaustive
   commit or pull-request list.
2. **Given** a reader wants full detail, **When** they follow the changelog link,
   **Then** they reach the v0.9.0 section containing the complete release record.
3. **Given** a future tagged release, **When** its workflow starts, **Then** it
   selects notes belonging to that exact tag rather than silently reusing the
   v0.9.0 copy.

### User Story 2 - Preserve an honest version boundary (Priority: P1)

As a maintainer, I can identify exactly which accumulated changes belong to
v0.9.0 while retaining an empty Unreleased section for work merged afterward.

**Why this priority**: A tag without a matching changelog boundary makes both
the current release and later releases difficult to audit.

**Independent Test**: Compare the prepared changelog with v0.8.0 and verify that
all existing Unreleased detail moves under v0.9.0, a new empty Unreleased
section remains, and both comparison links identify the correct boundaries.

**Acceptance Scenarios**:

1. **Given** the detailed changes accumulated after v0.8.0, **When** the release
   is prepared, **Then** those entries appear under one dated v0.9.0 heading.
2. **Given** work lands after the release preparation, **When** the changelog is
   updated, **Then** authors have a separate Unreleased section above v0.9.0.
3. **Given** the release has not yet been tagged, **When** readers inspect the
   preparation pull request, **Then** its scope and issue state make the pending
   publication step explicit rather than claiming that assets already exist.

### User Story 3 - Download a complete verified build (Priority: P1)

As a user ready to try the accumulated work, I can download the appropriate
v0.9.0 package for my platform and verify it with the published checksum file.

**Why this priority**: The release is valuable only when the tagged workflow
finishes and all promised artifacts are present.

**Independent Test**: After separately authorized publication, audit the tagged
workflow and release page for the expected packages, checksum coverage, concise
notes, and synchronized README version strings.

**Acceptance Scenarios**:

1. **Given** the reviewed preparation merge and explicit tag authorization,
   **When** v0.9.0 is tagged, **Then** the existing release workflow builds and
   publishes the supported packages.
2. **Given** all packages are attached, **When** checksum publication completes,
   **Then** every non-checksum release asset has exactly one matching checksum.
3. **Given** any release job fails or an expected artifact is absent, **When**
   the release is audited, **Then** issue #79 remains open and the release is not
   reported as complete.

### Edge Cases

- The v0.9.0 tag or release already exists when publication is authorized.
- A tag-specific release-note file is missing or named for the wrong version.
- GitHub attempts to append automatically generated notes to the curated copy.
- The binary release exists while one or more desktop jobs are still running or
  failed.
- A previous checksum file is downloaded during a workflow re-run and could be
  included in its own replacement.
- README version synchronization completes after the release assets and must be
  audited independently.
- New changes land on `main` after the preparation merge but before the tag is
  created, which would make the reviewed release boundary stale.

## Requirements

### Functional Requirements

- **FR-001**: S034 MUST prepare release version v0.9.0 from the complete reviewed
  change set since v0.8.0.
- **FR-002**: `CHANGELOG.md` MUST preserve the full detailed history under a
  dated v0.9.0 section and retain a new Unreleased section for later changes.
- **FR-003**: Changelog comparison links MUST compare Unreleased against v0.9.0
  and v0.9.0 against v0.8.0.
- **FR-004**: The GitHub release body MUST contain a short Highlights section
  with four to six user-meaningful bullets and one link to the full v0.9.0
  changelog.
- **FR-005**: The release body MUST NOT include an automatically generated
  commit list, pull-request inventory, installation guide, or duplicate of the
  detailed changelog.
- **FR-006**: Release notes MUST be selected by the exact tag so a future release
  cannot silently publish v0.9.0 highlights.
- **FR-007**: Offline automation MUST accept the approved tag-specific notes
  contract and reject generated notes, a fixed v0.9.0 body path, missing
  changelog linkage, and overly broad release copy.
- **FR-008**: Existing release outputs MUST remain intact: one Windows amd64 MSI,
  Linux amd64 and macOS arm64 desktop bundles, Linux and macOS headless archives
  for supported architectures, and one complete SHA256SUMS.txt.
- **FR-009**: Release publication MUST use the reviewed S034 merge commit and
  MUST stop if `main` advances before tag creation until the boundary is
  revalidated.
- **FR-010**: Creating the v0.9.0 tag and running the release workflow MUST wait
  for separate explicit authorization after the preparation pull request merges.
- **FR-011**: Issue #79 MUST remain open through preparation review and MUST
  close only after the published notes, workflow, assets, checksums, and README
  synchronization are verified.
- **FR-012**: Release verification MUST require no manual installation,
  uninstallation, virtual machine, or interactive observation session.
- **FR-013**: The S034 preparation pull request MUST use `Refs #79`, not a closing
  keyword, because merging preparation alone does not publish the release.
- **FR-014**: The release preparation MUST pass all canonical local gates before
  publication and all hosted pull-request gates before merge.
- **FR-015**: The constitution's release-time compliance review MUST be recorded
  without adding new governance controls or unrelated product work.

### Key Entities

- **Release boundary**: The reviewed merge commit, version tag, date, and
  comparison range that define the contents of v0.9.0.
- **Release highlights**: The small curated set of user-facing improvements
  displayed on the GitHub release page.
- **Detailed changelog**: The authoritative complete record for the release.
- **Release artifact set**: Every supported installer, desktop bundle, headless
  archive, and the checksum file produced for v0.9.0.
- **Publication evidence**: The tag, workflow run, release URL, asset inventory,
  checksum audit, and README synchronization proving the release completed.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The published release notes contain four to six highlight bullets,
  one full-changelog link, and zero generated commit or pull-request entries.
- **SC-002**: One hundred percent of detailed post-v0.8.0 changelog entries are
  retained under v0.9.0, with no detail duplicated into the release highlights.
- **SC-003**: One hundred percent of expected release assets are present and
  every non-checksum asset has exactly one matching SHA-256 entry.
- **SC-004**: All local and hosted verification gates complete successfully
  before the release is declared complete.
- **SC-005**: The tag, release page, packages, checksum file, README version
  strings, and issue #79 final state all identify v0.9.0 consistently.
- **SC-006**: A fixture replacing curated notes with generated notes fails the
  offline release automation contract.

## Clarifications

### Session 2026-08-30

- Q: How detailed should the GitHub release body be? A: Highlights only, with a
  link to the complete changelog.
- Q: What version follows the latest v0.8.0 release? A: v0.9.0, matching the
  existing release milestone and accumulated minor-version capabilities.
- Q: Does S034 require a manual installation walkthrough? A: No. Workflow,
  artifact, and checksum verification are the release gate.

## Assumptions

- PR #78 and all earlier included slices remain merged on `main` before the
  preparation branch is published.
- The current tagged release workflow remains the supported packaging path.
- GitHub Actions and release storage are available when tag publication is
  separately authorized.
- The release date is 2026-08-30 unless review crosses a calendar boundary, in
  which case the changelog date is refreshed before merge.

## Out of Scope

- New product capabilities, including the proposed GUI run-history slice.
- Manual platform installation or uninstallation observation.
- Fixing deferred Windows icon issue #33.
- Changing supported platforms, package formats, or dependency versions.
- Publishing the tag before the reviewed preparation change merges.
