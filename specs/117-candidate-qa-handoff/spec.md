# Feature Specification: S117 candidate QA and release handoff

**Feature Branch**: `codex/117-candidate-qa-handoff`

**Created**: 2026-10-25

**Status**: In Progress

**Delivery**: Candidate evaluation and demonstrated source corrections implemented; canonical eight gates pass. Publication and external review follow. Owner authorizes push, one official PR and at most two Codex review rounds; stop before merge, tag, staging or promotion.

**Input**: Kick off S117 under autopilot with #354-#356, use Spec-Kit, run candidate QA, implement demonstrated fixes and publish the source PR. Accepted main baseline is 11485a8a6ce97fe8e6966e1cc3cdb6a55b9d5651, with all 23 merged-main checks successful. Pending version remains v1.6.0; public release remains v1.5.1.

## Clarifications

### Session 2026-10-25

- Q: Does this kickoff authorize the release operations proposed for S117? A: No. The explicit latest boundary is authoritative: deliver candidate QA, corrections and source handoff through the reviewed PR; merge, tagging, draft staging and promotion remain later owner actions.
- Q: Must successful candidate repetitions close #355/#356? A: No. Preserve original failures and distinguish delivered diagnostics/budget behavior from the remaining original causes; close an issue only for its actual corrective outcome.
- Q: Which evaluation budget applies? A: Use the S116 supported defaults (120 seconds per case and two hours total), retaining explicit/resume budget identity, resource limits, namespace ownership and intentional expiry assertions.
- Q: How should discoveries affect scope? A: Correct demonstrated in-scope defects under regression discipline and track every distinct actionable defect against the exact source. Preserve each original report and record correction acceptance separately.
- Q: Can candidate observations stand in for final release identity? A: No. Record actual source commit, dirty flag, content and executable digests. A working-tree capture or earlier committed candidate remains that identity after later changes; do not relabel it as the final PR, merge, tag or package.

## User Scenarios & Testing

### User Story 1 - Receive complete candidate evaluation (Priority: P1)

A maintainer delegates one unattended candidate evaluation and receives a fixed-shape report of actual commands, negative inputs, diagnostic language, observations and cleanup, bound to the exact candidate inputs.

**Why this priority**: It supplies current release evidence after the supported completion defaults changed.

**Independent Test**: Use the existing owned evaluation and review pathway; every applicable manifest case and captured message has its required result/review association, and final reporting preserves the true outcome.

**Acceptance Scenarios**:

1. **Given** identified candidate inputs and the supported default budget, **When** evaluation finishes, **Then** every applicable case has its actual outcome, explicit exclusions remain distinct and all owned resources have attributable cleanup results.
2. **Given** a failed command or follow-up inspection, **When** the report is rendered, **Then** command, inspection, phase/request, cleanup and human diagnostic outcomes remain distinguishable.
3. **Given** captured communication, **When** contextual agent review finishes, **Then** every required exact-bound message has the fixed rubric and a human-readable disposition; unseen messages cannot be marked reviewed.

### User Story 2 - Correct demonstrated candidate defects (Priority: P1)

A maintainer receives fixes for observed project defects and accurate remaining dispositions for the independent Windows runtime/prerequisite and owned IPC findings.

**Why this priority**: Reliable candidate behavior is the practical outcome; passing repetitions alone do not explain historical failures.

**Independent Test**: A correction has a failing-before/passing-after regression for the demonstrated contract, while original captured results and unrelated unknown causes remain intact.

**Acceptance Scenarios**:

1. **Given** a reproduced component defect, **When** corrected, **Then** the intended operation and actionable diagnostic behavior work without weakening identity, authority, concurrency, data or resource contracts.
2. **Given** an original #355/#356 cause remains unconfirmed, **When** source results are reported, **Then** it remains visibly open with its remaining outcome and candidate impact separated from delivered diagnostics.
3. **Given** a distinct discovery, **When** durable work is identified, **Then** its GitHub record names the exact affected revision, outcome, dependencies, functional acceptance and parent release.

### User Story 3 - Receive an exact source release handoff (Priority: P2)

A maintainer can review the source increment and determine which later tag, package and publication actions still require authorization.

**Why this priority**: It advances #354 without conflating source readiness with public release.

**Independent Test**: Candidate notes, upgrade guidance and comparison decisions identify the pending version; the public selection remains unchanged and the handoff names actual source, evaluation, CI and review identities.

**Acceptance Scenarios**:

1. **Given** candidate corrections or findings, **When** release guidance is reconciled, **Then** it includes delivered scope and unresolved findings, preserving the distinction between candidate and published behavior.
2. **Given** owner-authorized publication of the source PR, **When** checks/reviews finish, **Then** every warranted comment is addressed within at most two Codex rounds and the owner receives an unmerged handoff.

### Edge Cases

- A case can exit successfully while independent inspection or cleanup fails; retain both outcomes.
- A known input/resource/ownership refusal must retain its original result and must not become a false pass through retries or a new allowance.
- A missing interpreter or platform-only operation must retain its defined exclusion/prerequisite result, not an invented execution result.
- A candidate capture may fail and later corrections pass targeted acceptance; preserve the captured failure and separate source/binary identities.
- Incomplete communication review, interrupted capture, stale source inputs or concurrent report updates must remain visible and use the established validation/recovery pathway.
- Hosted CI can fail on a platform not available locally; preserve the failure, investigate and correct only demonstrated in-scope causes.
- Public release and tag may not exist; source preflight cannot create or claim those identities.

## Requirements

### Functional Requirements

- **FR-001**: Candidate inputs MUST retain exact source, platform, executable and budget identities before execution and across report finalization.
- **FR-002**: The existing unattended routine MUST cover every applicable manifest case, declared negative scenario and output mode, preserving explicit platform exclusions.
- **FR-003**: Evaluation MUST preserve supported case/total budgets, CPU/memory/process/output bounds, namespace ownership, hidden/noninteractive execution and resource cleanup; it MUST not operate the owner's installed services, sessions or unrelated state.
- **FR-004**: Every captured message requiring review MUST receive actual contextual review using the established fixed rubric and exact case/message/digest associations, with validated canonical human and machine reports.
- **FR-005**: Original failures, command exits, inspection/request/phase diagnostics, cleanup outcomes and later corrective acceptance MUST remain separately attributable; no stale report may be relabelled as a new source result.
- **FR-006**: Demonstrated in-scope project defects MUST receive proportionate corrections with meaningful failing-before/passing-after regression evidence; security, authority, data, identity and concurrency assertions MUST remain intact.
- **FR-007**: #355 runtime/prerequisite and #356 IPC causes MUST remain independent; unconfirmed causes MUST remain open with precise remaining scope. Successful repeats and verification activity MUST NOT be issue-closure gates.
- **FR-008**: Every actionable distinct discovery MUST have a nonduplicated GitHub record under the existing taxonomy and milestone, with parent, dependencies and functional acceptance criteria.
- **FR-009**: Pending v1.6.0 source notes, changelog, upgrade guidance and comparison content bindings MUST reflect delivered changes and known issues, while public selection remains v1.5.1.
- **FR-010**: The handoff MUST name actual candidate and final source identities and distinguish source validation, exact-head hosted CI, later merged-main acceptance, tags, staged assets and public release.
- **FR-011**: S117 MUST use the complete Spec-Kit sequence, resolve blocking analysis findings before implementation and complete the existing canonical eight verification gates before source delivery.
- **FR-012**: Owner-authorized push and one official PR MUST include exact-head CI monitoring and individual feedback resolution within at most two Codex rounds, stopping before merge, tagging, staging or promotion.

### Key Entities

- **Candidate identity**: Source commit, dirty state, content digest, executable digests, platform/runtime identity and selected budget.
- **Evaluation record**: Manifest/case execution, observation completeness, communication, cleanup and immutable original result.
- **Message review**: Exact case/message binding, fixed rubric, reviewed context and disposition.
- **Defect disposition**: Reproduced behavior, affected identity, corrective outcome or explicitly unresolved cause, and linked GitHub issue.
- **Source handoff**: Pending version/content bindings, actual validation/CI/review identities and later owner-controlled release actions.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One fresh candidate evaluation accounts for 100% of applicable manifest cases and required captured-message reviews, with each explicit exclusion and incomplete outcome distinguishable.
- **SC-002**: Every delivered correction has an independently attributable before/after contract result; 100% of original failure records remain intact.
- **SC-003**: All discovered actionable defects have linked, nonduplicated issue dispositions, with no unconfirmed original cause closed because a repetition passed.
- **SC-004**: Candidate guidance, source handoff and public selection contain no conflicting version or publication claims.
- **SC-005**: The published source PR has terminal successful applicable checks, all warranted review feedback addressed and no more than two Codex rounds; it remains unmerged for the owner.

## Assumptions

- Existing GitHub access, native Windows tools and the existing nonroot Ubuntu-24.04 WSL environment remain available; ask for authentication only on actual failure.
- The established evaluator, schemas, review prompts and release preflight are reused; no new QA framework, production timeout policy or dependency migration is presumed.
- #354 remains open for its later public release/discovery outcome. Original #355/#356 closure depends on their own delivered corrective outcomes, not slice membership.
- A focused diagnostic or correction can use owned fixtures; no owner workstation installation, login/logout, public DNS mutation or release operation is part of this source slice.
