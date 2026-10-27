# Feature Specification: S119 v1.6.0 release execution

**Feature Branch**: `codex/119-release-execution`

**Created**: 2026-10-26

**Status**: Implemented

**Delivery**: Review branch codex/119-release-execution implements source discovery reconciliation and #379 after authorized release v1.6.0 (ID 423908120) at approved commit 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b; staging 40278046423 and promotion 40450478494 succeeded. All eight foreground gates passed in 703.326s with 2,360 frozen inputs unchanged. Source PR push and at most two Codex rounds are authorized; merge remains owner-controlled. See delivery.md and verification.md.

**Input**: Execute S119 for #354 with independent #355/#356 dispositions, using the S118 handoff. Deliver the actual release pathway and a reviewed source PR rather than another unchanged candidate replay.

## Clarifications

### Session 2026-10-26

- Q: Which exact source/tag may stage? A: The owner explicitly approved v1.6.0 at 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b; this approval does not transfer to a changed source.
- Q: Does draft staging approval include public promotion? A: No. Present the actual release/run/attempt/asset identities for a separate decision after staging succeeds.
- Q: Must #355/#356 close to publish? A: No. Retain independent known-issue disclosures and original failed evidence; no causal fix is claimed without a demonstrated outcome.
- Q: Should S119 repeat candidate preparation or create publication checkboxes? A: No. Reuse the completed source/evidence, track release operations as external delivery receipts and keep source implementation tasks distinct.

### Session 2026-10-27

- Q: Is public promotion authorized for the actual staged candidate? A: Yes. The owner approved release 423908120, source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b, staging run 40278046423 attempt 1 and all eight presented original assets. Guarded promotion 40450478494 subsequently published that exact candidate.

## User Scenarios & Testing

### User Story 1 - Receive one identifiable staged candidate (Priority: P1)

The owner receives a draft v1.6.0 release built from the exact approved source, with complete downloadable platform assets and unambiguous candidate provenance.

**Why this priority**: A concrete candidate enables an informed publication decision and advances the remaining release outcome.

**Independent Test**: Resolve the approved tag, successful staging attempt and draft; every expected asset identifies the same source and original candidate bytes.

**Acceptance Scenarios**:

1. **Given** exact-tag approval for accepted source, **When** the existing release pathway stages v1.6.0, **Then** one immutable tag and draft identify that source and the complete intended asset set.
2. **Given** an existing conflicting tag, draft, failed stage or incomplete asset set, **When** execution encounters it, **Then** the original identities and failure remain visible and no replacement or unsupported success is claimed.

### User Story 2 - Publish the approved candidate and discover it accurately (Priority: P1)

After reviewing the actual staged candidate, the owner can approve its public release. Readers find the intended downloads, version, comparison and installation guidance, with truthful known-issue disclosures.

**Why this priority**: #354 is a public delivery outcome, not a source-preparation record.

**Independent Test**: Actual public release, latest pointer, complete checksums and discovery references identify the approved candidate; pending-only claims are reconciled through review.

**Acceptance Scenarios**:

1. **Given** approval of the concrete staged identities, **When** promotion completes, **Then** the original assets become the public v1.6.0 release with matching checksums and current latest discovery.
2. **Given** promotion is not yet approved or fails, **When** source records are updated, **Then** they retain the actual pending state and do not advertise a completed public release.
3. **Given** a completed public release, **When** source discovery changes are prepared, **Then** the reviewed PR accurately selects that release and preserves historical snapshots.

### User Story 3 - Preserve independent defects and bounded review (Priority: P2)

Maintainers retain the original #355 runtime/prerequisite and #356 IPC diagnoses while reviewing the actual release increment.

**Why this priority**: Publication must not erase unresolved causes or reassign older failed evidence to new bytes.

**Independent Test**: Each issue and delivery record names its own delivered and remaining outcome; source checks and reviews are bound to the actual revision.

**Acceptance Scenarios**:

1. **Given** the original causes are unconfirmed, **When** release delivery proceeds, **Then** both defects remain disclosed and open independently of passing repeats or publication.
2. **Given** source publication authority, **When** the PR is reviewed, **Then** warranted feedback is addressed within at most two Codex rounds and the owner receives an unmerged handoff with terminal successful checks.

### Edge Cases

- The release tag exists with a different source, or a draft contains conflicting/incomplete bytes.
- A staging attempt fails or is rerun; the actual successful attempt and original failures remain attributable.
- Tag/draft approval does not authorize public promotion; exact staged identities must be presented first.
- A complete already-published release may need latest-pointer recovery without content replacement.
- A source correction changes the prospective candidate; approval for the earlier exact SHA cannot silently transfer.
- Public comparison may follow the actual latest release before source-owned pending prose is merged; retain both states explicitly.
- Missing authentication or workflow-write authority requires actual owner involvement, not a substitute release or guard bypass.

## Requirements

### Functional Requirements

- **FR-001**: Staging MUST use only v1.6.0 at owner-approved 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b, preserving immutable tag/source identity and exact-main acceptance.
- **FR-002**: The staged draft MUST retain the actual release ID, workflow run/attempt, complete platform asset identities and original Windows candidate manifest/bytes.
- **FR-003**: Conflicts, missing assets, failed attempts and authentication failures MUST retain their actual outcomes without overwrites, fabricated identities or weakened guards.
- **FR-004**: Public promotion MUST require the owner's separate approval of the concrete staged identities and preserve the existing standing native-testing waiver and artifact/checksum contracts.
- **FR-005**: After approved promotion, public release, latest pointer, downloads and checksums MUST identify the same approved candidate.
- **FR-006**: Source-owned published selection, pending prose and public discovery MUST accurately reflect actual publication state, with historical snapshots and upgrade/rollback guidance preserved.
- **FR-007**: #355/#356 MUST retain independent unresolved original causes and immutable failed evidence; passing repetitions, staging or release activity MUST NOT imply closure.
- **FR-008**: Every demonstrated new defect MUST receive a separate attributable issue and proportionate correction, without speculative runtime/IPC changes or unchanged full-catalog replay.
- **FR-009**: S119 MUST follow the complete installed Spec-Kit sequence, blocking analysis and all eight canonical foreground gates before any source commit.
- **FR-010**: Authorized source PR publication MUST address every warranted review finding within at most two Codex rounds, record exact-head successful checks and stop unmerged.

### Key Entities

- **Approved source**: Exact immutable source, tag, owner authority and successful merged-main acceptance.
- **Staged candidate**: Draft release ID, successful workflow/attempt, complete assets, manifest, original sizes and hashes.
- **Publication decision**: Separate owner approval tied to the actual candidate identities.
- **Public release**: Published ID/tag/source, latest pointer, downloads, checksums and discovery selection.
- **Defect disposition**: Individual demonstrated correction or explicitly unconfirmed original cause with unchanged evidence.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One approved immutable source produces one identifiable draft with 100% of the expected assets and no unexplained identity mismatch.
- **SC-002**: After separate approval, all public assets, checksums and latest discovery identify the intended v1.6.0 release.
- **SC-003**: Every included issue has an individual truthful outcome; no original failure or unresolved cause is overwritten or implicitly closed.
- **SC-004**: The source PR has successful applicable checks, zero unresolved warranted feedback and no more than two Codex rounds, while remaining unmerged for owner review.

## Assumptions

- Existing staging and standing-waiver promotion workflows are reused. No new release framework, dependency migration or installed-system operation is presumed.
- Approved tag source is existing accepted main; S119 planning/discovery records have a separate source identity and cannot be silently substituted into that approved tag.
- Source PR publication is authorized. Tag/staging approval was separately received; the owner separately approved the exact staged candidate on 2026-10-27 and guarded promotion completed.
- #354 remains incomplete until its actual public release/discovery outcome is delivered. #355/#356 remain open absent their individual corrective outcomes.

## Historical delivery reconciliation by S120

PR #380 was owner-merged on 2026-10-27 as dac2e97895353435ad30515402b9f99d5e5d4ebc; its remote branch was deleted. Exact-main CI, CodeQL, Website and Pages succeeded with 23 successful applicable checks. See [delivery.md](delivery.md#owner-merge-and-post-merge-acceptance-2026-10-27) for identities and public HTTPS observations. Earlier pre-merge statements describe the original handoff, not current PR state. Immutable v1.6.0 source, assets and original failure evidence remain unchanged.
