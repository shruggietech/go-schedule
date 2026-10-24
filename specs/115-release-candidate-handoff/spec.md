# Feature Specification: S115 exact candidate source and native follow-up

**Feature Branch**: `codex/115-release-candidate-handoff`

**Created**: 2026-10-24

**Status**: In Progress

**Delivery**: In progress on review branch codex/115-release-candidate-handoff. Source corrections and native follow-up implemented; canonical acceptance and full candidate evaluation remain pending. Owner-authorized push and official PR, maximum two Codex review rounds; stop before merge, tag, staging or promotion.

**Input**: Owner kickoff S115 with #354, #355, #356 and reconciled documentation commit 18e6fd805e941184c1550a51c6de2ba922d73856. Baseline main is 7bce610396bf83210e09376a2343d2e6ad21d180.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Use one current release handoff (Priority: P1)

A release maintainer receives current source preparation that agrees on the pending version, delivered fixes, remaining defects, upgrade effects and the exact next authorized publication steps.

**Why this priority**: The release outcome is still pending after several source slices. A maintainer needs a concrete next action without mistaking a source pass for a published package.

**Independent Test**: Review the prepared candidate records as one unit; every claimed identity and delivered outcome has an attributable source and every absent publication identity is explicitly absent.

**Acceptance Scenarios**:

1. **Given** reconciled S114 documentation, **When** incorporated, **Then** it records the actual one-round owner merge and successful main deployment instead of prospective review claims.
2. **Given** pending v1.6.0 source, **When** a maintainer follows the handoff, **Then** release copy, upgrade/rollback guidance and comparison receipts agree and public latest remains v1.5.1.
3. **Given** no tag or staged package exists, **When** readiness is reported, **Then** no package hash, release identity or published outcome is invented and later authorization boundaries remain explicit.

### User Story 2 - Attribute current native outcomes (Priority: P1)

A maintainer can evaluate one identified current candidate silently on owned Windows and native Linux fixtures, retaining command, prerequisite, observation and cleanup outcomes separately.

**Why this priority**: #355 and #356 retain independent historical failures whose causes are unconfirmed. Another unrelated passing run cannot resolve those causes.

**Independent Test**: Retain a complete current-candidate report and targeted original-path observations with exact source/binary identities and unchanged case, concurrency and resource contracts.

**Acceptance Scenarios**:

1. **Given** a fresh candidate evaluation, **When** an operation or subsequent observation fails, **Then** the original exit/output, failed phase and partial state remain attributable and incomplete observations cannot satisfy full-state assertions.
2. **Given** a reproduced project defect, **When** corrected, **Then** the affected behavior has a demonstrated correction and its independent native issue record identifies the affected source.
3. **Given** no causal reproduction, **When** results are reported, **Then** #355/#356 remain open with their original evidence and current disposition, without an indefinite repetition loop or a claim of correction.

### User Story 3 - Continue after interruption or changed source (Priority: P2)

A maintainer can distinguish a reusable current observation from stale or incomplete material and resume the documented release path using the same identified candidate.

**Why this priority**: Preparation, evaluation, staging and publication occur at different times and must not silently exchange source or package identities.

**Independent Test**: Exercise applicable existing refusal and recovery boundaries with changed or incomplete inputs, retaining exact identities and the original failure outcome.

**Acceptance Scenarios**:

1. **Given** source or candidate inputs change, **When** prior material is considered, **Then** stale identities remain historical and cannot be relabelled as current acceptance.
2. **Given** a later publication attempt is interrupted, **When** the maintainer follows the handoff, **Then** the same immutable source, candidate and release identity remain authoritative and missing or conflicting inputs stop the operation.

### Edge Cases

- Interpreter crash, launcher failure, prerequisite refusal, application failure and failed follow-up observation remain distinct; expanded writer phases remain opt-in and existing original-output corrections remain intact.
- Six concurrent writers retain literal identity and exactly one record per successful invocation; no serialization, retry-to-success or relaxed contention budget hides a failure.
- Historical captures and dependency overlays retain their own source and environment identities; a newer pass does not rewrite them.
- An absent runtime is an explicit unavailable observation, not a pass. Windows and Linux observations are separately identified.
- Quiet owned evaluation must not operate installed services, replace system runtimes, change account sessions, mutate public endpoints or DNS, or steal desktop focus.
- Changed, missing, malformed or interrupted evidence must remain distinguishable from a complete current outcome; secrets and arbitrary environment dumps stay out of shared material.
- Actual tags, packages and public release selection remain unchanged until a later owner-authorized operation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Incorporate commit 18e6fd8 and preserve the original 574d071 history, with actual S114 delivery, one review round, merge, deployment and independent issue dispositions accurately recorded.
- **FR-002**: Deliver one coherent pending v1.6.0 source handoff whose release copy, upgrade/rollback guidance and comparison receipts agree, including a read-only pre-tag source validator shared with Release that checks badge, health version, unique calendar-dated heading and exact notes link; distinguish source preparation from absent staged and public artifacts.
- **FR-003**: Bind fresh candidate observations and any reused release material to their exact source, executable, content and environment identities; refuse stale or incomplete inputs where the existing workflow requires complete current inputs.
- **FR-004**: Perform a bounded, silent owned Windows/native Linux follow-up of the current application and original #355/#356 paths, retaining exact outcomes, phase attribution, partial observations and cleanup separately without changing existing limits or assertions.
- **FR-005**: Correct demonstrated in-scope defects with attributable regressions and native GitHub records; retain #355/#356 as open unknown-cause investigations unless their own corrective or supported upstream outcome is established.
- **FR-006**: Preserve the exact immutable-source, candidate, release-identity, checksum, asset allowlist and standing-waiver contracts in the later publication path, including failure and interruption guidance.
- **FR-007**: Keep publication/review/merge workflow outside implementation-task checkboxes, maintain issue-level traceability and keep #354 open for its actual public-release outcome.
- **FR-008**: Publish one official source PR and consider every actual review finding within at most two Codex rounds; perform zero merges, tags, release staging runs or public promotions under this kickoff.

### Key Entities

- **Candidate source**: One identified source revision and content set, with separately attributable executable builds and pending-publication state.
- **Native observation**: Owned scope, environment, exact source/build identities, original execution and observation outcomes, elapsed phases, partial state and cleanup result.
- **Release handoff**: Pending version, content identities, known defects, upgrade/rollback contract and later immutable staging/promotion sequence.
- **Issue disposition**: Delivered correction, supported upstream resolution or unconfirmed cause, recorded independently for each issue.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All carried S114 documentation states the actual merge and single completed review; no prospective second-review or open-PR claim remains as current status.
- **SC-002**: Every fresh native result is attributable to one current source/build and preserves the original outcome; every known historical failure category receives a separate current disposition.
- **SC-003**: Every demonstrated correction has its own functional issue outcome, while unresolved #355/#356 causes remain visibly open.
- **SC-004**: All release-owned preparation records agree on pending v1.6.0; public release discovery remains v1.5.1 throughout this source slice.
- **SC-005**: The source PR reaches owner review with terminal applicable checks and all actual review findings addressed within the two-round cap, with no merge, tag, staging or promotion performed.

## Assumptions

- Current main 7bce610 already contains the S114 helper corrections and release-note disclosure; those changes are not relabelled as new S115 fixes.
- The S114 local reconciliation commit is an authorized input, not an already merged or published source change.
- The existing complete unattended catalog and targeted native diagnostics are the starting point; routine experiment selection is a planning decision, not a new owner prerequisite.
- A passing candidate observation is evidence about that candidate, not proof that an intermittent historical cause is corrected.
- Publication uses the existing reviewed path after a later owner merge and separate tag/promotion authorization; this kickoff does not authorize staging.

## Clarifications

### Session 2026-10-24

- Q: Does this kickoff include tag-only staging or publication? A: No. It authorizes source implementation, push, one PR and up to two Codex rounds; stop before merge, tag, staging or promotion.
- Q: Can passing candidate observations close an unexplained historical cause? A: No. Preserve the original receipts and leave #355/#356 open unless their individual corrective or supported upstream outcome is established.
- Q: Does carrying 18e6fd8 imply it already reached main? A: No. Incorporate it as an identified S115 source change after the blocking design analysis.

Functional scope, identities, lifecycle, failure handling, resource constraints, authorization and completion signals are clear. Experiment selection and executable implementation detail are planning decisions. No owner question is required; all 16 specification-quality items remain passing. No clarify hook is registered.
