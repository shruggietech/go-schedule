# Feature Specification: S129 Prelaunch Baseline Publication

**Feature Branch**: `codex/129-prelaunch-release`
**Created**: 2026-10-31
**Status**: In Progress
**Delivery**: Verified v2.0.0 source on review branch codex/129-prelaunch-release; source PR publication/review and later actual public delivery remain separate. See verification.md, delivery.md and handoff.md.
**Native scope**: [#412](https://github.com/shruggietech/go-schedule/issues/412), milestone 19 (Prelaunch baseline publication).

**Input**: The owner authorizes Spec-Kit/autopilot, automatic source push and one official PR, in-scope CI/review corrections and at most two Codex requests. Stop at final owner review/merge. This continues the recommended release of the simplified prelaunch baseline after S128.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review one coherent release identity (Priority: P1)

A maintainer can review the next release's source, notes and current-state behavior as one coherent increment, while downloads continue to identify the actual published version until a later approved publication.

**Why this priority**: S128 is integrated at 0c1ab32129a3bd3fcb4694915aaccef4166a502c, but public v1.7.0 packages still predate the simplified baseline.
**Independent Test**: Read the candidate identity, release notes and release-source checks together and distinguish them from live public discovery.

**Acceptance Scenarios**:

1. **Given** the merged simplified baseline, **When** the next release source is prepared, **Then** a single version consistent with declared versioning policy identifies its source-owned candidate records.
2. **Given** public v1.7.0 remains available, **When** source preparation is reviewed, **Then** current download selection and published examples still identify the actual public package.
3. **Given** the owner states no users and no migration obligations, **When** release behavior is described, **Then** fresh initialization, valid current-state reopen and explicit obsolete-state refusal are described without converters or historical upgrade requirements.

### User Story 2 - Publish the accepted bytes through the existing process (Priority: P1)

After source integration and the owner's later exact tag/draft approval, a maintainer stages the current packages, inspects their actual identities and assets, then seeks separate concrete-asset promotion approval. The same accepted bytes become public and discoverable.

**Why this priority**: S129 is a release-delivery slice, not another preparation-only design slice or a new release framework.
**Independent Test**: Follow the concrete same-slice handoff from reviewed source to exact staged assets and final published readback, keeping every approval and evidence boundary explicit.

**Acceptance Scenarios**:

1. **Given** only the current source/PR authorization, **When** CI and review finish, **Then** the agent stops before owner merge, tagging, staging or promotion.
2. **Given** accepted main and later exact tag/draft approval, **When** existing staging runs, **Then** the tag, source, staging attempt, candidate manifest and complete expected asset set identify the same candidate.
3. **Given** concrete-asset promotion approval, **When** existing promotion runs, **Then** it publishes the same inspected assets with checksums and accurate current discovery; partial or failed outcomes remain explicit.

### User Story 3 - Keep release claims and deferred work truthful (Priority: P2)

A maintainer sees current capability comparison, current installation guidance, known unresolved outcomes and retained evidence without confusing source checks, attended testing or published-package acceptance.

**Why this priority**: The prelaunch simplification removes historical obligations while keeping current correctness, security and artifact identity requirements.
**Independent Test**: Trace current release claims and comparisons to their source and public authority, and identify the work that remains deferred.

**Acceptance Scenarios**:

1. **Given** automated and headless checks, **When** their results are reported, **Then** attended native observations remain explicitly untested under the standing waiver.
2. **Given** immutable old releases, notes, comparisons and failed evidence, **When** the new candidate is prepared, **Then** the old records and attribution remain unchanged.
3. **Given** unresolved SQLite, historical runtime and clustered-execution work, **When** release scope is finalized, **Then** those outcomes are not silently included, closed or advertised as fixed.

### Edge Cases

- A prepared version does not yet have a tag, draft or public release.
- Public latest remains older than the current source baseline.
- A tag already exists, staging is incomplete, assets change, a checksum is absent, or a workflow fails. Existing refusal/idempotency behavior remains authoritative.
- Comparison content has changed without a matching current review, or an older review belongs to different bytes/source.
- The final authorized review returns a code-review result rather than a separate security-review receipt; report the actual mode.
- An obsolete or malformed development store is refused; startup does not convert or silently delete it.
- The owner merges/closes the PR before monitoring completes; respect the new terminal state and do not merge or retrigger reviews.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The next candidate MUST have one documented version selected from actual changes under the declared versioning policy, with no tag or publication implied by that selection.
- **FR-002**: Source-owned release notes, current-candidate guidance and release identity checks MUST agree on that candidate and describe the merged prelaunch baseline accurately.
- **FR-003**: Public discovery and download examples MUST keep the actual published release selected until later approved promotion; preparation MUST NOT advertise unbuilt or unpublished packages as available.
- **FR-004**: Current release guidance MUST describe direct application initialization, valid current-state persistence, current-only preferences and helper stores, explicit reset/refusal and current qualification without historical migrations/upgrade matrices.
- **FR-005**: The candidate comparison MUST describe current shipped capabilities accurately, retain official-source attribution and current review binding, and preserve immutable historical datasets.
- **FR-006**: The same S129 slice MUST retain a concrete continuation using existing release workflows, exact source/tag/attempt/asset identities and separate owner source-merge, tag/draft and concrete-promotion boundaries. No new release framework is required.
- **FR-007**: Publication MUST retain accepted asset bytes and existing checksums/identity safeguards; interruption or failure MUST remain distinguishable from a completed release.
- **FR-008**: Attended native testing MUST remain explicitly unperformed under the standing waiver; automated CI, source checks and headless results MUST NOT be represented as installed/tagged-package or physical-input evidence.
- **FR-009**: Historical release assets, notes, comparison bindings and 215 retained inputs MUST remain unchanged; original failures remain attributable.
- **FR-010**: SQLite 1.60.1/libc 1.77.1, existing runtime behavior and deferred #355/#356/#400/#404/cluster outcomes MUST remain outside source expansion unless a genuine in-scope CI/review finding warrants a bounded correction.
- **FR-011**: Current authorization MUST produce one official source PR after complete foreground verification and handle every received finding, with at most two Codex requests including automatic PR-open review; stop before owner merge.
- **FR-012**: Native planning MUST distinguish completed source readiness from actual public release delivery, keeping the release outcome open until its individual functional outcome is delivered.

### Key Entities

- **Prepared release**: selected version, reviewed source, notes and current capability claims.
- **Accepted candidate**: immutable tag/source, successful staging attempt and actual asset inventory/digests.
- **Published release**: the approved accepted assets, complete checksums and public discovery.
- **Evidence receipt**: exact source/artifact identity, actual result and expressly untested scope.
- **Owner boundary**: source integration, exact tag/draft staging, then concrete-asset promotion.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every source-owned candidate identity resolves to one selected version, with zero premature public download claims.
- **SC-002**: Current release claims cover all merged baseline changes and distinguish every unresolved deferred outcome without inventing a fix or migration guarantee.
- **SC-003**: The continuation identifies every expected asset and each remaining owner boundary, with no rebuild after concrete asset acceptance.
- **SC-004**: Current and historical evidence remain independently attributable, with zero drift in retained historical inputs.
- **SC-005**: The current source PR reaches its owner review boundary with terminal successful applicable CI and all received findings reconciled, while actual release completion remains separately tracked.

## Assumptions

- Latest explicit owner authorization satisfies the normal pre-push halt for this one PR. It does not grant source merge, tagging, staging, promotion or recurring monitoring.
- Version selection is a reversible source decision; creating/pushing the selected immutable tag remains a later exact approval.
- No users and no migration obligations remain current owner policy. Preserve valid current state and explicit reset control.
- The existing release/standing-waiver workflows are reused; no new dependency, storage framework or product capability is needed.

## Clarifications

### Session 2026-10-31

- Q: Should S129 add historical upgrade or migration support? A: No, the explicit owner prelaunch policy remains authoritative.
- Q: Does this kickoff authorize release mutation now? A: It explicitly authorizes source push/PR and stops for the owner merge ritual. Exact tag/draft and concrete-asset promotion approvals remain later boundaries in this same slice.
- Q: Should public latest change during source preparation? A: No, it continues to identify the actual published release until approved promotion.

## Implementation discovery

- #413 corrects the source guard and current health example to the full v-prefixed version actually injected by release builds. It is a bounded FR-002 correction, with existing runtime and tag/publication boundaries unchanged. The native bug has its own source outcome; #412 remains actual public delivery.
