# Feature Specification: S135 v2.0.1 publication closeout

**Feature Branch**: `codex/135-v201-publication-closeout`

**Created**: 2026-11-06

**Status**: Implemented

**Delivery**: Public release v2.0.1 and review branch `codex/135-v201-publication-closeout`; all 13 tasks and foreground source qualification complete. See [delivery](delivery.md) and [verification](verification.md). Owner integration and current-guide deployment remain separate.

**Input**: Kick off S135 under autopilot, push the official PR automatically, monitor CI and resolve findings, limit Codex reviews to two including automatic review, and stop for final owner review and merge. Complete the remaining #419 publication and release-discovery outcome after S134 PR #431.

## User Scenarios & Testing

### User Story 1 - Obtain the accepted maintenance release (Priority: P1)

A user can download the maintenance release containing the first-use improvements and integrated maintenance corrections from the public release selected by the project.

**Why this priority**: #419 remains functionally incomplete while those changes exist only in accepted source.

**Independent Test**: The public release exposes the complete supported package set under v2.0.1 and identifies its accepted source.

**Acceptance Scenarios**:

1. **Given** accepted main at `25f42d6eb35ff0cec2753f0f2934396ba50df954`, **When** the owner authorizes exact-source staging and later approves the inspected draft assets for promotion, **Then** the public v2.0.1 release contains those same selected bytes and source identity.
2. **Given** native testing remains waived, **When** a user reads the release, **Then** the published disclosure states that native installation and attended GUI testing are untested.

### User Story 2 - Discover the release consistently (Priority: P1)

A prospective user sees consistent release availability across the README, website fallback, comparison catalog and first-task guidance.

**Why this priority**: Publishing packages without updating the supported discovery surfaces leaves the parent outcome incomplete.

**Independent Test**: Current discovery surfaces select v2.0.1 and describe the first-use behavior now included in the public packages without changing historical release records.

**Acceptance Scenarios**:

1. **Given** v2.0.1 is public, **When** current discovery records are integrated and deployed, **Then** their version and links identify that public release.
2. **Given** publication has not completed, **When** a discovery surface is served, **Then** it does not advertise v2.0.1 as already public.

### User Story 3 - Review a bounded closeout (Priority: P2)

The owner receives one official source PR containing the publication closeout and discovery changes, with issue traceability and explicit evidence boundaries.

**Why this priority**: The owner retains final review and merge authority while autopilot completes the authorized engineering work.

**Independent Test**: The PR names #419, distinguishes completed publication from changes awaiting merge/deployment, and records at most two Codex review requests including automatic review.

**Acceptance Scenarios**:

1. **Given** a finding affects the candidate, **When** it is resolved, **Then** the PR and checks identify the revised candidate and preserve prior failures diagnostically.
2. **Given** the PR is ready, **When** autopilot hands off, **Then** it remains open for owner review and merge.

### Edge Cases

- An existing tag or public release with a different source or asset set stops the affected release operation without replacing it.
- Failed or partial staging leaves the release unpublished; fresh diagnostic records identify the failed run and attempt.
- A changed draft asset set invalidates the inspected-asset selection and requires a new concrete promotion decision.
- Public metadata, checksums or downloads that disagree remain an explicit defect rather than a successful publication claim.
- A failed discovery deployment leaves public packages and current website availability separately reported.

## Requirements

### Functional Requirements

- **FR-001**: Release v2.0.1 MUST identify accepted main `25f42d6eb35ff0cec2753f0f2934396ba50df954` and its immutable tag.
- **FR-002**: Public delivery MUST include four daemon/CLI archives, two desktop archives, the Windows MSI, its candidate manifest and the complete checksum inventory.
- **FR-003**: Promotion MUST preserve the selected draft package bytes and expose their exact identifiers, sizes and digests in the closeout record.
- **FR-004**: Release notes MUST retain the standing native-testing disclosure and distinguish source, hosted builds, package readback and attended native observations.
- **FR-005**: Current README release discovery, website fallback, comparison selection, first-task guidance and related installation/desktop/website availability notices MUST identify the actual public release consistently after publication.
- **FR-006**: Historical v2.0.0 comparison records and existing protected qualification evidence MUST remain unchanged.
- **FR-007**: The official PR MUST preserve issue-level traceability to #419 and distinguish public publication from source integration and Pages deployment.
- **FR-008**: Existing release/promotion mechanisms MUST be reused unless a demonstrated defect requires an in-scope correction with a linked GitHub bug record.
- **FR-009**: The closeout MUST preserve failed observations and report incomplete functional outcomes literally; verification activities MUST NOT become issue closure criteria.
- **FR-010**: Autopilot MUST stop before final PR merge, count automatic review toward the two-review maximum, and obtain the existing exact-tag/draft and separate actual-asset promotion authorizations.

### Key Entities

- **Release candidate**: Accepted source identity, immutable version tag and staging run/attempt.
- **Selected asset**: Release ID, asset ID, name, size, digest and associated candidate identity.
- **Discovery record**: Current release identity and public links, separate from immutable historical comparison data.
- **Delivery record**: Publication status, pending integration/deployment, linked issue and evidence limitations.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The public v2.0.1 release contains exactly nine expected assets, with eight content entries covered by its checksum file.
- **SC-002**: All selected package identities agree between inspected draft and public readback, with zero unexplained replacements.
- **SC-003**: All current discovery surfaces and availability guides identify the actual public version after integration/deployment, with zero historical-record changes.
- **SC-004**: One official PR is ready for final owner review with zero unresolved applicable findings and no more than two Codex review requests.

## Assumptions

- S134 PR #431 is already owner-merged. Its accepted main candidate has successful exact-main CI, CodeQL and Pages deployment.
- The owner approved exact-source tagging and draft staging on 2026-11-06; the owner subsequently approved that inspected asset selection and promotion succeeded. These are operation authorizations, not new product eligibility rules.
- The existing native-testing waiver continues. Original Windows investigations #355/#356 and owner-deferred clustered execution remain outside this publication increment.
- This slice reuses #419 and its existing milestone. It does not start another dependency preparation campaign or introduce a new qualification framework.
