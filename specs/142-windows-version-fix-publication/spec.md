# Feature Specification: Ship Windows version identity fix

**Feature Branch**: `codex/142-windows-version-fix-publication`

**Created**: 2026-11-09

**Status**: Implemented

**Delivery**: Public/latest v2.0.3 release 439902909 and all nine anonymous downloads are verified. Current-discovery source is implemented and locally verified on continuation branch codex/142-v203-publication-closeout; its final integration/deployment continues automatically under existing full owner authorization.

**Input**: Owner directs completion of the release without further confirmation and closure of issues already corrected in code. Track actual publication in #462.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Download the corrected maintenance release (Priority: P1)

Users obtain public v2.0.3 packages carrying the existing Windows version-identity correction.

**Independent Test**: Read public release metadata, original package bytes and Windows version fields without launching the desktop.

**Acceptance Scenarios**:

1. Given the integrated #458 correction, when the maintenance release is published, then its selected version identifies the desktop resources, application and installer consistently.
2. Given the completed publication, when opening current release guidance and comparison, then they identify v2.0.3 and its actual downloads.

### User Story 2 - See completed work closed (Priority: P2)

Maintainers see implemented outcomes closed without waiting for speculative historical causal attribution.

**Independent Test**: Compare each issue's functional outcome to integrated source and record the resulting disposition.

**Acceptance Scenarios**:

1. Given delivered functional corrections, when reconciling their records, then they close with concrete implementation references and honest historical limits.
2. Given an unimplemented deferred feature, when reconciling this release, then it is not represented as completed.

### Edge Cases

Incomplete or changed assets, moved tags and unsuccessful required checks refuse publication. Existing historical artifacts remain attributable. The latest explicit owner instruction replaces earlier per-stage approval pauses; it does not waive source correctness or asset integrity.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Public v2.0.3 MUST contain the integrated Windows selected-version correction.
- **FR-002**: Tag, packages, original asset bytes and public checksum inventory MUST identify one accepted source.
- **FR-003**: Current public documentation and comparison discovery MUST follow actual publication.
- **FR-004**: Delivered issue outcomes MUST be closed without treating more verification or unresolved historical attribution as independent closure gates.
- **FR-005**: Historical tags, assets, release notes and comparison records MUST remain unchanged.
- **FR-006**: Execution MUST remain unattended, avoid operator desktop focus and sanitize public evidence.

### Key Entities

Accepted source, immutable tag, draft/public release, inspected assets, published discovery and issue disposition.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Public/latest v2.0.3 includes the corrected version identities.
- **SC-002**: All public assets match the inspected originals and checksum inventory.
- **SC-003**: Public guidance selects the actual published release.
- **SC-004**: Every delivered issue outcome in this slice has an explicit completed disposition.

## Assumptions

Existing release machinery, standing native-testing waiver and PR-based source integration remain. The owner authorizes automatic merge, tagging and promotion for this release. Earlier ceremony/approval-only pauses no longer apply.

## Clarifications

### Session 2026-11-09

- Q: Are separate tag/promotion approvals needed? A: No, the owner explicitly authorizes completion without further confirmation.
- Q: Do causal investigations keep already corrected functional issues open? A: No; record implemented outcomes and preserve historical uncertainty honestly.
