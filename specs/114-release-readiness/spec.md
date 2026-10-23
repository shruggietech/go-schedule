# Feature Specification: v1.6.0 release readiness and Windows failure diagnosis

**Feature Branch**: `codex/114-release-readiness`

**Created**: 2026-10-23

**Status**: In Progress

**Delivery**: Source commit 15e2578554a5ebc1b4b9b0ea7eec2f2113620575 is published in PR #365. Canonical local application verification passed; hosted CI/review handling remains in progress, maximum two Codex rounds. Original #355/#356 causes remain unconfirmed and open; #354 awaits later authorized publication. Stop before merge, tag or public promotion.

**Input**: Owner kickoff S114 with #354, #355 and #356, following the recommendation recorded in #354.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Attribute payload failures (Priority: P1)

A maintainer investigating the retained Windows payload failures can distinguish an interpreter crash, rejected prerequisite and writer failure, and can follow a supported resolution when a cause is established.

**Why this priority**: Passing repetitions have not explained the original failures; incorrect attribution would hide real defects.

**Independent Test**: Compare attributable runs of the existing six simultaneous writers and injected candidate failures with retained native exits and exact records.

**Acceptance Scenarios**:

1. **Given** an interpreter or candidate process fails, **When** its result is recorded, **Then** the original exit/output and executable identity remain visible without retrying it into success.
2. **Given** a demonstrated project defect, **When** corrected, **Then** successful invocations retain one distinct record each and failed invocations remain visible.
3. **Given** no causal reproduction, **When** the investigation is reported, **Then** #355 stays open with exact scope and known uncertainty; a passing run is never labelled its fix.

### User Story 2 - Explain stalled observations (Priority: P1)

A maintainer can identify which owned operation or inspection consumed its deadline and correct the demonstrated component while retaining partial observations.

**Why this priority**: A successful command and failed subsequent inspection have different operational meanings.

**Independent Test**: Exercise sustained owned Windows and Linux requests and incomplete inspection paths using unchanged declared case and resource limits.

**Acceptance Scenarios**:

1. **Given** a completed command and failed follow-up observation, **When** reported, **Then** command outcome, incomplete observation and cleanup results are distinct.
2. **Given** a reproduced component defect, **When** corrected, **Then** operation and independent inspection satisfy the existing limits and complete-state assertions.
3. **Given** no causal reproduction, **When** reporting results, **Then** #356 remains open and original failures stay attributable.

### User Story 3 - Prepare an exact release handoff (Priority: P2)

A release maintainer receives one coherent v1.6.0 preparation that states the reviewed source, upgrade/rollback effects, remaining bugs and exact next publication steps.

**Why this priority**: Source preparation must not advertise an unpublished package as available.

**Independent Test**: Inspect the release metadata and run existing release/comparison checks against reviewed source without creating a tag or publishing assets.

**Acceptance Scenarios**:

1. **Given** pending v1.6.0 source, **When** preparation changes, **Then** release notes, upgrade guidance and comparison receipts agree and public latest remains v1.5.1.
2. **Given** the PR is ready, **When** handing it to the owner, **Then** issue-level functional dispositions, reviewed source, CI and remaining tag/promotion boundaries are explicit.
3. **Given** publication remains unauthorized, **When** source work ends, **Then** no tag, staging run or public promotion is performed and #354 remains open.

### Edge Cases

- Crash exit codes, missing/version-refused prerequisites and application failures remain separate; no credentials or arbitrary environment dump enters diagnostics.
- Interrupted or partial requests cannot satisfy full-state assertions; mutation retries cannot obscure uncertain outcomes.
- Evidence from historical source/binaries remains immutable; diagnostic overlays and changed dependency baselines receive separate attribution.
- No historical cause is inferred solely from CPU utilization or successful repetitions.
- No installed service, account session, system DNS or system runtime replacement is needed for owned investigation.
- Release assets do not yet exist for an uncreated tag; preparation must not invent candidate hashes, release IDs, dates or publication evidence.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Investigate #355's interpreter crash and prerequisite refusal independently, retaining original native outcomes and distinguishing demonstrated project causes from supported upstream dispositions and unconfirmed observations.
- **FR-002**: Preserve six concurrent writers, literal identity and exact record-count assertions; do not serialize, retry failures into success or increase contention budgets.
- **FR-003**: Investigate #356 with attributable phase/request evidence and fresh owned Windows/Linux comparison, including prior/current dependency baselines when relevant to a causal hypothesis.
- **FR-004**: Preserve case/total deadlines, resource limits, complete-state assertions, command outcomes, partial observation data and cleanup results.
- **FR-005**: Correct demonstrated in-scope causes with regression coverage; retain unresolved original causes as open issue outcomes with reproducible context rather than asserting completion.
- **FR-006**: Reconcile #354's source notes, schema-v24 upgrade/rollback guidance, pending status and comparison review receipts after demonstrated corrections.
- **FR-007**: Preserve the existing immutable-source/tag, exact release-ID, candidate-manifest, asset allowlist, checksum and standing-waiver contracts; describe their later sequence without staging or publishing.
- **FR-008**: Maintain original evidence identities and bound new observations to exact source, inputs, environment and outcome.
- **FR-009**: Deliver one official PR and handle every actual CI/review finding; consume at most two Codex rounds, then stop for owner review before merge, tagging or public promotion.
- **FR-010**: Keep published v1.5.1 discovery unchanged until a later authorized promotion; keep #354 open and close #355/#356 only if their own functional outcomes are delivered.

### Key Entities

- **Investigation observation**: Source/input identities, affected runtime, owned scope, original process or request outcome, elapsed phase evidence and causal disposition.
- **Release handoff**: Pending version, reviewed source, upgrade/rollback contract, remaining known bugs, comparison receipts and later publication instructions.
- **Issue disposition**: Delivered correction, authoritative supported upstream resolution or unconfirmed cause, with independent issue state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every original failure category has a distinct recorded disposition and retains original outcomes without claiming unsupported fixes.
- **SC-002**: All six successful simultaneous writer invocations remain individually attributable; incomplete observations remain unsuccessful.
- **SC-003**: Every release-owned preparation record agrees on pending v1.6.0 and the schema-v24 upgrade; published discovery remains v1.5.1.
- **SC-004**: The final owner handoff identifies each included issue's delivered and incomplete outcomes, exact reviewed revision, terminal CI and no more than two review rounds.
- **SC-005**: The slice performs zero merges, tags, release staging runs, public promotions or installed-system mutations.

## Assumptions

- Current baseline is merged main 7f7c5ec2c5252bb4b116b8d3b365438f30196779.
- S113 diagnostic and partial-observation functionality is already delivered; this slice does not relabel that work as new.
- Original evidence is retained privately under the S112/S113 archives, and external/public records must not include credentials.
- Lack of reproduction leaves the underlying bugs open. Investigation completion and bug correction are distinct.
- Existing release workflows remain the publication path. The owner expressly excludes tags and promotion from this run.

## Clarifications

### Session 2026-10-23

- Q: Does this kickoff authorize release staging or publication? A: No. It authorizes push and one PR with at most two Codex rounds; stop before merge, tag or public promotion.
- Q: Can repeated passes close an unknown historical cause? A: No. Keep #355/#356 open unless their individual corrective or supported upstream outcome is established.

Coverage: functional scope, data identity, lifecycle, failure handling, constraints, dependencies, completion signals and security are clear. Runtime comparisons and exact diagnostic experiment selection are deferred to planning, with no unresolved product decision.

### Explicit launcher-output exception

The owner approved corrected Bash launcher failure output during S114: retain the original launcher status and explanation instead of a secondary missing-capture-file error. Successful output and writer assertions remain unchanged. Expanded writer phases are opt-in through GOSCHEDULE_PROBE_DIAGNOSTICS=1.
