# Feature Specification: S131 v2.0.0 Publication

**Feature Branch**: `codex/131-v2-publication`

**Created**: 2026-10-31

**Status**: Implemented

**Delivery**: Public/latest v2.0.0, release 430039348, approved immutable source 75eae1ae11ad40db38341dc7d63f32876cba6483; verified current-source reconciliation and exact-byte receipts in verification.md and delivery.md. Official PR/owner merge remain workflow evidence.

**Input**: Owner-authorized S131 continues the recommended actual v2.0.0 publication outcome, automatically pushes one official reconciliation PR, handles CI and all warranted review findings with at most two total Codex requests, and stops before owner merge.

## Context and authority

[#412](https://github.com/shruggietech/go-schedule/issues/412) and milestone 19 are the existing public-delivery outcome. Accepted release source is owner-merged S130 main 75eae1ae11ad40db38341dc7d63f32876cba6483, following S128's unused prelaunch simplification and S129's prepared v2.0.0 identity. Root and desktop select SQLite 1.61.1/libc 1.79.0. Public latest is v1.7.0 at selection. This slice executes the existing [S129 release continuation](../129-prelaunch-release/handoff.md), rather than creating a preparation-only release framework.

Source push and PR publication are explicitly authorized. The owner separately approved exact immutable tag/draft staging at accepted source and promotion of inspected release 430039348, whose eight original assets remain unchanged. Neither release decision was inferred from the source autopilot request. Final reconciliation integrates by owner merge.

## Clarifications

### Session 2026-10-31

- Q: Which source should be staged? A: The already accepted S130 main 75eae1ae11ad40db38341dc7d63f32876cba6483, subject to fresh exact-source checks and explicit tag/draft approval. S131 records and later fallback reconciliation do not require rebuilding that source.
- Q: Does the kickoff authorize release mutation? A: It authorizes execution of S131 and source publication; preserve the separately stated exact tag/draft and concrete-asset promotion decisions.
- Q: Are migrations or historical upgrade qualification needed? A: No. There are no users; retain current initialization, valid persistence, explicit obsolete-state refusal and access security.
- Q: Does source PR completion close the release outcome? A: Only after actual public delivery/discovery and integration of remaining current-source fallback wording are complete. CI and attended testing are engineering evidence, not issue functional closure gates.

## User Scenarios & Testing

### User Story 1 - Stage an attributable release candidate (Priority: P1)

The owner can approve one exact source and inspect one complete draft candidate with consistent product and platform identity before publication.

**Why this priority**: A source-only change does not provide downloadable v2.0.0 packages.

**Independent Test**: Inspect the actual draft inventory, package identities and source attribution against the approved source without publishing it.

**Acceptance Scenarios**:

1. **Given** accepted source and explicit staging approval, **When** the existing release process completes, **Then** one immutable v2.0.0 tag identifies that source and one draft contains all eight expected original assets.
2. **Given** incomplete, failed or conflicting candidate state, **When** the candidate is inspected, **Then** it remains unpublished and its actual state and failure are recorded without silently replacing accepted bytes or moving a tag.
3. **Given** native observations unavailable under the standing waiver, **When** evidence is reported, **Then** the untested attended behavior is explicit and source/local tests are not presented as acceptance of staged bytes.

### User Story 2 - Publish the approved bytes (Priority: P1)

After reviewing the actual inventory, the owner can authorize publishing those same bytes as public/latest v2.0.0.

**Why this priority**: Publication is the missing functional outcome in #412.

**Independent Test**: Resolve the public release and compare every original asset identity and digest with the approved draft receipt.

**Acceptance Scenarios**:

1. **Given** separate promotion approval naming the concrete candidate, **When** the existing promotion completes, **Then** the same release is public/latest with eight original assets plus one checksum file containing eight entries.
2. **Given** interrupted promotion, **When** work resumes, **Then** actual identities and outcomes are inspected before any bounded retry, preserving original failures and accepted asset bytes.

### User Story 3 - Discover the published baseline consistently (Priority: P2)

A visitor sees the actual public release in comparison and installation guidance; current repository fallback and wording agree with it after normal PR integration.

**Why this priority**: Package availability and current instructions must agree, including offline website builds.

**Independent Test**: Read public/latest, HTTPS home, comparison, installation and representative downloads, then inspect the current-source reconciliation diff.

**Acceptance Scenarios**:

1. **Given** actual v2.0.0 publication, **When** public discovery builds, **Then** current pages select v2.0.0 and do not advertise unshipped behavior.
2. **Given** publication receipts, **When** the source reconciliation is reviewed, **Then** fallback/current wording and S131 delivery records reflect those receipts while tagged notes, comparison bindings and prior release assets remain unchanged.
3. **Given** satisfied source CI and reviews, **When** handoff occurs, **Then** the official PR remains unmerged for the owner and accurately identifies any remaining functional delivery outcome.

### Edge Cases

- An existing tag or nonmatching draft requires inspection; never move the immutable tag, overwrite an accepted asset, or invent an absent source/run receipt.
- Changed release ID, asset ID, name, size, checksum, manifest source or successful run identity blocks promotion until its cause is understood.
- Unavailable authentication, failed staging or pending CI remains an explicit incomplete state.
- A Pages delay or fallback mismatch is attributed separately from successful release publication.
- A genuine new defect is tracked separately against the affected source/release; unrelated #355/#356 causes remain unconfirmed and clustering remains owner-deferred.
- The final PR's automatic review consumes request one. At most one manual second request is authorized, including after corrections; an unreviewed corrected head must be reported honestly.

## Requirements

### Functional Requirements

- **FR-001**: Deliver v2.0.0 from the exact owner-approved accepted source through the existing immutable tag and draft process, with eight complete platform/manifest assets.
- **FR-002**: Attribute each actual candidate asset to its release, successful staging run/attempt and source, retaining its identity, size and digest plus embedded product/version and manifest agreement.
- **FR-003**: Publish only after separate approval of the inspected inventory and preserve all eight accepted original assets unchanged; add one matching checksum file and select that release as latest.
- **FR-004**: Make public release discovery, comparison, current installation/README wording and repository fallback agree with actual publication, with source changes delivered through one official review PR.
- **FR-005**: Preserve previous release assets, notes, immutable comparison datasets/reviews and all 221 retained inputs; preserve exact pending, failed and successful observations separately.
- **FR-006**: Keep no-user current-state behavior, access security and truthful standing native-testing disclosure; introduce no migrations, old-version qualification, cluster work or release framework.
- **FR-007**: Reconcile #412 and milestone 19 according to actual functional delivery; keep unrelated original incidents and deferred outcomes open.
- **FR-008**: Automatically push the verified source branch and official PR, address every warranted CI/review finding within scope, consume no more than two Codex requests, and stop before owner merge.

### Key Entities

- **Accepted source**: Immutable commit, matching reviewed source tree and exact-source acceptance receipts.
- **Candidate**: Release ID, tag/source and staging run/attempt plus the complete eight-asset inventory.
- **Publication**: The same candidate's public/latest state, checksum identity and discovery receipts.
- **Reconciliation PR**: Current fallback/wording and S131 execution records with independent source/review identity.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One public/latest v2.0.0 release provides nine downloadable files with eight checksum entries and one coherent source/product version.
- **SC-002**: All eight original accepted asset identities and bytes are retained unchanged through publication.
- **SC-003**: Current public home, comparison and installation routes resolve over HTTPS, and applicable version selectors/current guidance agree with v2.0.0 after reconciliation integration.
- **SC-004**: All 221 retained inputs and previous immutable release records remain unchanged.
- **SC-005**: One official reconciliation PR is ready for owner merge with truthful terminal CI/review evidence and at most two review requests; no merge is performed by the agent.

## Assumptions

- The existing release and standing-waiver workflows remain the delivery authority; no new dependencies, service or framework are needed.
- The owner retains separate exact-source staging and actual-asset promotion decisions, and final PR merge.
- Attended native behavior remains untested under the standing waiver; no active-workstation installation, focus capture or service restart is needed.
- No tag/draft approval for earlier v1.7.0 carries forward to v2.0.0.
