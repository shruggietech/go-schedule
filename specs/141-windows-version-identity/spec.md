# Feature Specification: Windows version identity

**Feature Branch**: `codex/141-windows-version-identity`

**Created**: 2026-11-09

**Status**: Implemented

**Delivery**: Implemented locally on review branch codex/141-windows-version-identity; [verification](verification.md) records successful canonical verification and compiled metadata. Official PR and hosted checks follow under owner authorization, then stop for final owner review/merge.

**Input**: Owner kickoff S141 under autopilot, implementing #458.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Identify the selected Windows release (Priority: P1)

A maintainer selects a release version once and obtains desktop executable metadata, application version and installer identity that describe that release.

**Why this priority**: Published v2.0.2 GUI File.Version is 1.0.0.0 while its installer ProductVersion is 2.0.2.

**Independent Test**: Build a selected Windows candidate and inspect its executable and installer without executing or installing them.

**Acceptance Scenarios**:

1. Given a selected release version, when packaging the desktop, then numeric and textual version resources identify the selected release consistently with runtime and MSI identity.
2. Given a deliberately mismatched executable, when packaging, then the mismatch is refused before release upload.

### User Story 2 - Build an identifiable development desktop (Priority: P2)

Developers can build through maintained local targets without inheriting a misleading released version.

**Independent Test**: Inspect a development build and an explicit-version local build.

**Acceptance Scenarios**:

1. Given no selected release, when building locally, then development identity is explicit and numeric resource fields remain valid.
2. Given an invalid selected version, when building, then a clear error prevents misleading output.

### Edge Cases

Missing version, malformed version, numeric component overflow, development labels, stale resources from earlier builds and mismatch between installer and executable must have explicit dispositions. Existing published bytes cannot be overwritten.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Windows desktop resources MUST identify the selected release consistently with application and MSI versions.
- **FR-002**: Maintained local and release builds MUST use one selected version input without a separately maintained version pin.
- **FR-003**: Development builds MUST be labeled explicitly and retain valid numeric resources.
- **FR-004**: Invalid versions and packaged metadata mismatches MUST fail with actionable errors before upload.
- **FR-005**: Published v2.0.2 tags, assets and historical evidence MUST remain unchanged.
- **FR-006**: All work MUST avoid operator-desktop focus. Compiled metadata candidates MUST be inspected without launching the desktop or installing the package; canonical owned CLI/daemon fixtures remain headless. Publication remains a later separate decision.

### Key Entities

- Selected version: release label and its numeric Windows representation.
- Candidate identity: executable resources, linked application version and installer version.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A selected candidate has matching executable, application and installer version identities.
- **SC-002**: Development builds identify development status and malformed inputs produce no newly accepted candidate.
- **SC-003**: A deliberately mismatched candidate is refused by maintained packaging checks.
- **SC-004**: Published asset identities and original evidence remain unchanged, with zero operator focus acquisition.

## Assumptions

Existing Wails, WiX and release machinery are retained. #355/#356 and the deferred roadmap are independent. No new release version is selected by this source slice.

## Clarifications

### Session 2026-11-09

- Q: Does S141 rebuild public v2.0.2? A: No, deliver build behavior through a source PR, preserve immutable public assets.
- Q: Is an installed GUI walkthrough required? A: No, inspect compiled metadata headlessly and retain native runtime limitations.

- Q: Does metadata-only qualification prohibit canonical CLI/daemon fixtures? A: No. The desktop/MSI candidate is inspected without launch or installation, and canonical owned CLI/daemon fixtures continue headlessly.
