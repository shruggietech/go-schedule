# Feature Specification: Reliability release preparation

**Feature Branch**: `codex/124-reliability-release-preparation`

**Created**: 2026-10-27

**Status**: Implemented

**Delivery**: Implemented on review branch codex/124-reliability-release-preparation; 19 source tasks completed, full current Windows/WSL evaluation and all eight canonical gates passed. verification.md preserves exact source stages and first failures. Authorized PR publication, CI/reviews and final owner handoff are tracked separately; release delivery remains open on #394.

**Input**: Owner-authorized S124 with [#394](https://github.com/shruggietech/go-schedule/issues/394): prepare the post-v1.6 reliability release, push and open its source PR, monitor CI, and handle at most two Codex review rounds. Stop before merge, tag creation, staging or promotion.

## Clarifications

### Session 2026-10-27

- Q: Patch or minor version? A: Select v1.7.0. The proposed v1.6.1 was conditional; documented additive maintainer CLI incident and replay selections require a minor increment under the declared SemVer policy. No breaking change was identified.
- Q: Does source preparation close #394? A: No. This source slice references #394; actual delivery remains a separately authorized operation. #355 and #356 remain open because their historical causes are unconfirmed.
- Q: What version do candidate QA binaries report? A: The existing source QA builder reports 0.0.1-dev. Preserve that truth and exact assembled-source identity; do not fabricate a future tagged candidate manifest or reuse historical QA as a current result.
- Q: How should competitor evidence change? A: Reassess every existing cell, retain qualified unknowns, and add a new v1.7.0 snapshot/review without modifying v1.6.0 records or public selection.

## User Scenarios & Testing

### User Story 1 - Understand and prepare the reliability update (Priority: P1)

A maintainer can identify the selected release, understand delivered fixes and additive diagnostics, and follow upgrade and rollback instructions without mistaking preparation for a published release.

**Why this priority**: A coherent source candidate is required before any later release operation.

**Independent Test**: Read the prepared release notes, changelog, README and platform guides together; every current candidate reference agrees and published links still resolve to v1.6.0.

**Acceptance Scenarios**:

1. **Given** the accepted post-v1.6 source, **When** a maintainer reads the candidate notes, **Then** v1.7.0 identifies delivered reliability changes and additive diagnostic controls, with #355/#356 explicitly unresolved.
2. **Given** an existing installation, **When** a maintainer follows the platform guide, **Then** it explains stopping owned services, preserving a complete database/configuration/credential backup, matching binaries for rollback, and schema v24 continuity from v1.6.0.
3. **Given** an unmerged source proposal, **When** a reader follows the public release badge or website default, **Then** v1.6.0 remains the published release and no future assets are advertised as present.

### User Story 2 - Compare supported capabilities accurately (Priority: P2)

A prospective user can view a release-scoped comparison based on current official evidence, with edition and deployment qualifications and explicit unknowns.

**Why this priority**: Release preparation must refresh the maintained comparison without inventing feature improvements.

**Independent Test**: Select the prepared snapshot in the comparison tooling and trace all 80 cells to qualified product notes and sources; all ten capability review decisions are present.

**Acceptance Scenarios**:

1. **Given** the eight-product, ten-capability comparison, **When** the v1.7.0 candidate is selected explicitly, **Then** it reflects current official documentation and the delivered go-schedule interface changes.
2. **Given** mixed vendor editions or incomplete evidence, **When** a cell is presented, **Then** its note preserves the applicable edition, preview or availability limits.
3. **Given** preserved v1.6.0 records, **When** source preparation finishes, **Then** those historical records and the public default remain unchanged.

### User Story 3 - Hand off a traceable source candidate (Priority: P2)

A maintainer receives an exact-source handoff that distinguishes local Windows/WSL evaluation, source CI, bot reviews and later tagged assets, and can authorize the next release operation from concrete records.

**Why this priority**: Evidence must identify what was exercised and what still requires owner action.

**Independent Test**: Trace the handoff to its source identity, binary hashes, generated QA records, CI head and review receipts without treating any of these as a historical root-cause finding.

**Acceptance Scenarios**:

1. **Given** the assembled candidate, **When** unattended QA completes, **Then** its fixed-schema report retains first outcomes and separately records semantic review of every actual message.
2. **Given** green source checks and satisfied reviews, **When** the source PR is handed off, **Then** it remains unmerged and no tag, draft release or promotion has occurred.
3. **Given** a newly discovered defect, **When** it affects the candidate, **Then** a specific existing or new bug record tracks it independently of #394 and the unconfirmed historical incidents.

### Edge Cases

- A new candidate version has no remote tag yet: source-only guards and future tag URLs are described as preparation, never proof of tagged assets.
- A vendor source is conflicting, inaccessible or scoped to a different offering: retain unknown or narrow the note; do not infer absence or borrow another edition.
- An interrupted QA run retains its first attempt and original source/binary/runtime selections; resume under the same overall budget only through the existing coordinator contract.
- A current failure does not prove a historical cause. Preserve historical inputs and file a separate attributable defect.
- Later factual handoff edits are distinguished from the frozen source exercised by runtime QA; material implementation edits require appropriate fresh qualification.

## Requirements

### Functional Requirements

- **FR-001**: Candidate metadata, changelog, release notes and current preparation examples MUST consistently select v1.7.0 with a dated SemVer rationale.
- **FR-002**: Notes MUST identify delivered cancellation, lifecycle, script-output and diagnostic changes without claiming resolution of historical #355/#356 causes.
- **FR-003**: Upgrade/rollback guidance MUST cover complete database sidecars, configuration, protected credentials, matching older binaries and schema v24 continuity; it MUST NOT describe v24 as newly introduced by v1.7.0.
- **FR-004**: Public release selection, latest badge/link and immutable v1.6.0 notes/comparison/review records MUST remain unchanged.
- **FR-005**: A new comparison snapshot MUST assess all eight products and ten capabilities using official sources reviewed 2026-10-27, with qualified Resolve Actions Pro 8.1 MCP/HA evidence and explicit unresolved vendor claims.
- **FR-006**: The companion review MUST bind the selected version and settled content digests and record changed/no-change decisions for every capability; source-only identity MUST NOT be described as an existing remote tag.
- **FR-007**: A source handoff MUST distinguish actual dev-version QA binaries, assembled-source identity, current evaluation results, exact-head CI, review receipts and pending future release asset identities.
- **FR-008**: Preparation MUST retain the standing native-testing waiver wording and outstanding attended observations without advertising automated evidence as attended results.
- **FR-009**: S123 records MUST reconcile owner-merged PR #393 and verified accepted main without overwriting historical verification provenance.
- **FR-010**: Issue #394 and its parent roadmap MUST show the selected version and source-preparation state while retaining their historical proposal and actual-delivery boundary. #355/#356 MUST remain open.
- **FR-011**: All source changes MUST be pushed to an official PR with issue traceability; every actionable review comment MUST be addressed within at most two Codex rounds. The workflow MUST stop before merging, tagging, staging or promotion.

### Key Entities

- **Prepared release**: selected version, accepted base, source proposal, release notes, unresolved issues and separately pending publication identities.
- **Comparison snapshot**: release-scoped product contexts, capability cells, source URLs and review dates.
- **Comparison review**: version, source reference, content digests and ten semantic decisions.
- **Source qualification record**: immutable input hashes, actual binary versions, generated run/corpus/report identities, outcomes, review attribution and cleanup receipts.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A reader can identify one selected pending version and one current public version across all current release-preparation surfaces without contradictory claims.
- **SC-002**: All 80 comparison cells and all ten review decisions have current qualified evidence, while historical v1.6.0 records retain their exact bytes.
- **SC-003**: Platform guidance provides a complete recovery path using preserved state and matching binaries, without implying old binaries can read migrated state.
- **SC-004**: The owner receives one traceable, unmerged source PR and a concrete list of subsequent release operations; #394 remains open until actual delivery.

## Assumptions

- Accepted main is a2e4653021b6007135550ed8f84f1421f9072db0. Existing release workflows, package identity derivation, schemas and public capability status remain authoritative.
- No new application capability, dependency framework, schema migration or unchanged historical replay campaign is needed for this source-preparation slice.
- Engineering verification includes full new assembled-source Windows/WSL CLI QA, semantic message evaluation, release/comparison guards and the canonical foreground aggregate. These are implementation checks, not GitHub issue closure criteria.
- Existing headless launchers and resource-bounded QA are used; no user login/logout, service disruption or foreground consoles are introduced.
