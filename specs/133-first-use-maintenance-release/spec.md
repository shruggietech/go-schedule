# Feature Specification: S133 First-use Maintenance Release

**Feature Branch**: `codex/133-first-use-maintenance-release`

**Created**: 2026-11-06

**Status**: Implemented

**Delivery**: Review branch codex/133-first-use-maintenance-release; #424 source correction and #419 maintenance preparation verified by the full eight-gate aggregate, 35 desktop/11 site browser cases and 1763 current QA passes. See [verification.md](verification.md) and [delivery.md](delivery.md); owner merge and actual public delivery remain separate.

**Input**: Owner-authorized S133 autopilot for #419 and #424, consolidating dependency PRs #420-#423, automatic push/official PR, monitored CI and at most two total Codex review requests including automatic review. Stop ready for final owner review and merge.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build one consistent maintenance candidate (Priority: P1)

A maintainer can build the same selected desktop dependency baseline through local, CI and release entry points, including the upstream Linux child-frame bridge correction, without an older build tool silently remaining selected.

**Why this priority**: #424 identifies an incomplete integration in pending PR #420; a maintenance candidate needs a coherent build-tool baseline.

**Independent Test**: Inspect the runtime and every active desktop build entry point, then exercise existing build and automation contracts.

**Acceptance Scenarios**:

1. **Given** the proposed desktop runtime update, **When** a supported build entry point is used, **Then** runtime and build tool select the same patch release.
2. **Given** the four dependency proposals, **When** combined, **Then** manifests and lockfiles describe the complete selected baseline without unrelated dependency changes.

### User Story 2 - Obtain accurately described first-use packages (Priority: P1)

An operator obtains a maintenance release whose Run details expose the own Run ID separately from Source run, with setup, failure and same-task recovery guidance matching the available package version.

**Why this priority**: The field and guide are merged, but the field is absent from public v2.0.0 packages. #419 tracks actual public delivery.

**Independent Test**: Check candidate version/notes/comparison and Run details correlation, then follow existing immutable publication after later owner approvals.

**Acceptance Scenarios**:

1. **Given** public v2.0.0 and pending source, **When** current guidance is read, **Then** source-only behavior remains distinguished from downloadable behavior.
2. **Given** an approved future publication, **When** packages and documentation are selected, **Then** version, package identity, Run ID guidance and discovery agree.
3. **Given** an older release, **When** S133 source is prepared, **Then** its assets, notes, comparison and tag remain unchanged.

### Edge Cases

- Cancelled dependency-PR checks do not establish acceptance of the combined candidate.
- Upstream's Linux correction does not establish application exploitability; same-origin frames retain their origin authority.
- A prepared PR is distinct from public delivery; #419 remains open at this handoff.
- Historical Windows #355/#356 and owner-deferred clustering remain independent.
- Current valid state remains durable; no migration, old-version upgrade campaign, service installation or active-workstation GUI operation is introduced.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Align the desktop runtime and every active local, CI and release build-tool version, including existing automation expectations (#424).
- **FR-002**: Integrate the exact bounded changes of PRs #420-#423, with complete root/desktop module metadata and matched frontend/site manifests and lockfiles.
- **FR-003**: Classify and prepare one maintenance version with matching source metadata, unique calendar-dated changelog heading, tag-specific notes and source-bound comparison (#419).
- **FR-004**: Describe own Run ID, first-task setup/failure/recovery and dependency changes accurately; preserve source-versus-public attribution until publication.
- **FR-005**: Preserve historical release records/assets, diagnostic evidence, current-state correctness, local IPC authority and existing resource/assertion contracts.
- **FR-006**: Reuse existing staging and standing-waiver promotion machinery and record same-slice continuation after owner source integration, immutable tag/draft approval and separate concrete-asset promotion approval.
- **FR-007**: Stop this execution at the official PR ready for owner review/merge. No merge, tag, staging or promotion. Automatic review counts as request one; at most one additional request is allowed.
- **FR-008**: Keep issue completion distinct: #424's source correction may close through merge; #419 stays open until public package/discovery delivery.

### Key Entities

- **Maintenance source candidate**: Reviewed revision, selected versions, content-bound comparison and delivery state.
- **Published release**: Immutable tag/source and platform asset/checksum inventory; v2.0.0 stays unchanged during this execution.
- **Review ledger**: Exact source/check identities and no more than two Codex requests including automatic review.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All active desktop build entry points agree with the selected runtime patch.
- **SC-002**: All four proposals have individually attributable representation without unrelated version expansion.
- **SC-003**: Candidate version, notes, comparison and first-use guidance identify one maintenance outcome while current public availability stays truthful.
- **SC-004**: Historical bytes/evidence remain unchanged and the report separates source readiness, public availability and untested attended behavior.

## Assumptions

- v2.0.1 is the intended patch: presentation/documentation and dependency repairs do not change public contracts or formats.
- Automatic source publication and in-scope feedback fixes are authorized; final merge belongs to the owner.
- Release delivery follows the explicitly requested source-PR stop; release authorization is not inferred.
- The attended-native-testing waiver applies; automated verification and artifact identity remain required.

## Clarifications

### Session 2026-11-06

- Q: Does kickoff authorize release publication? A: No. Stop at source review/merge; tag/draft and inspected-asset approvals remain later.
- Q: Add a dependency-management framework? A: No. Align existing pins/automation and integrate the four bounded proposals.
- Q: What completes each issue? A: #424 is the source correction; #419 is public delivery and stays open here.
