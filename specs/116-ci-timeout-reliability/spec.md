# Feature Specification: S116 CI timeout reliability and release handoff

**Feature Branch**: `codex/116-ci-timeout-reliability`

**Created**: 2026-10-25

**Status**: Implemented

**Delivery**: Implemented on review branch codex/116-ci-timeout-reliability; canonical eight gates passed with native exit 0. Owner authorized push, one official PR and at most two Codex review rounds, tracked outside implementation checkboxes. Stop before merge, tag, staging or promotion.

**Input**: Owner explicitly directs increased CI completion allowances throughout this project and kicks off S116 with #372 and #354-#356, plus derived cross-CI policy #373. Baseline main is 9e744e434f07a8d22041e0d4d6318867410900b6. This instruction supersedes earlier issue guidance against increasing verification deadlines.

## Clarifications

### Session 2026-10-25

- Q: Should earlier no-timeout-increase guidance prevent this work? A: No. The owner explicitly authorizes broader CI timeout allowances; record the new policy and preserve the original failed receipts.
- Q: Does increased verification time permit skipped assertions or retries into false success? A: No. Increase the time allowed for work to complete and retain actual failures, exits and cleanup results.
- Q: Does this kickoff authorize release publication? A: It authorizes source implementation, push, PR and two review rounds, stopping before merge, tag, staging or promotion.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Let verification finish (Priority: P1)

Maintainers can run all project CI systems with increased, explicit completion allowances, including nested test runners and slower platform preparation, without losing failure diagnostics.

**Why this priority**: Short nested deadlines can terminate valid work despite a generous outer CI job allowance. The owner explicitly requires a project-wide correction.

**Independent Test**: Compare the complete CI timeout inventory with the delivered policy; controlled work lasting beyond the former allowance can finish under the new allowance, while intentionally overdue work still reports a real failure.

**Acceptance Scenarios**:

1. **Given** all project verification and publishing workflows, **When** a job performs normal bounded work, **Then** it has the declared increased completion allowance, subject to the hosting platform maximum.
2. **Given** nested package, desktop or browser verification, **When** work outlasts the old incidental deadline but remains inside the new policy, **Then** assertions can finish instead of being prematurely canceled.
3. **Given** an intentional timeout or cancellation scenario, **When** it exercises timeout behavior, **Then** its intended short deadline and failure semantics remain meaningful.

### User Story 2 - Keep search on the exact scheduler (Priority: P1)

Opening a search result that reconnects to another daemon leaves Search visible and explains the identity mismatch. Matching identities continue to open the intended context.

**Why this priority**: The exact merged-main macOS job failed this scenario at its former five-second whole-test limit (#372).

**Independent Test**: The mismatch and matching-identity paths retain their identity assertions with sufficient time for their awaited interactions to complete.

**Acceptance Scenarios**:

1. **Given** a result for one daemon and a selection resolving to another, **When** Open completes, **Then** Search remains visible with the existing actionable mismatch message and no context opens for the wrong daemon.
2. **Given** the expected daemon, **When** Open completes, **Then** the intended context opens, including the existing popup activity path.

### User Story 3 - Attribute native verification outcomes (Priority: P1)

Maintainers can distinguish Windows interpreter failures, prerequisite refusal, preparation time, command execution, follow-up inspection and cleanup using exact source/runtime receipts and supported verification allowances.

**Why this priority**: #355 and #356 retain independent original failures; S115 added three metadata deadline observations. The owner now authorizes increased verification allowances as a corrective outcome.

**Independent Test**: Targeted owned Windows/Linux paths report original native outcomes and phase timings under the explicitly selected new budget, with all concurrent identities and cleanup assertions retained.

**Acceptance Scenarios**:

1. **Given** original native failures, **When** a later timeout policy is exercised, **Then** each result identifies its own source, runtime and budget and never overwrites an earlier receipt.
2. **Given** six concurrent successful writers, **When** they finish within the supported allowance, **Then** six distinct literal records remain; failed writers retain their original exit and diagnostics.
3. **Given** a failed or incomplete inspection, **When** reported, **Then** command outcome, completed partial data, failed phase and cleanup remain distinguishable.

### User Story 4 - Hand off the real release source (Priority: P2)

A release maintainer receives current source notes, known-issue dispositions and a concrete later staging/promotion sequence for the reviewed v1.6.0 candidate.

**Why this priority**: #354 remains a publication outcome; source preparation must not be mistaken for a released package.

**Independent Test**: Release metadata and issue dispositions name the actual reviewed source and clearly identify absent tag/package/publication identities.

**Acceptance Scenarios**:

1. **Given** S115 is merged, **When** records are reconciled, **Then** they retain the actual failed main CI and successful website/CodeQL outcomes.
2. **Given** current source corrections, **When** handed off, **Then** known issues, upgrade/rollback and comparison bindings agree, public selection remains v1.5.1 and later authorization steps are explicit.

### Edge Cases

- An inner timeout expires before a larger outer job limit.
- A job already inherits the hosting platform maximum.
- A test deliberately checks a short timeout, cancellation or scheduler latency.
- A declared or resumed QA budget differs from the new default.
- A whole-test deadline races an equally long element lookup.
- The original CLR failure does not recur and its cause remains unconfirmed.
- Source verification succeeds while no tagged package exists.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Inventory every repository-owned CI job, nested test runner and verification helper deadline; classify incidental completion limits separately from intentional timeout/cancellation assertions and production behavior.
- **FR-002**: Increase finite incidental completion allowances across all CI systems consistently, declare outer job allowances explicitly up to the platform maximum, and document exact previous/new values and any genuine hosting ceiling.
- **FR-003**: Apply the same supported timeout policy to canonical local verification and hosted execution so local success does not rely on a different incidental cap.
- **FR-004**: Preserve all functional, identity, concurrency, resource ownership, security, coverage and failure assertions; intentional timeout behavior continues to fail when its declared boundary is exceeded.
- **FR-005**: Preserve #372 exact-target navigation and actionable mismatch behavior while allowing the full asynchronous interaction to finish under the new policy.
- **FR-006**: Provide a supported increased native evaluation completion allowance, retain explicit caller/resume budget identity, and keep historical versus current outcomes separately attributable.
- **FR-007**: Retain #355 interpreter/prerequisite distinctions and #356 command/partial-inspection/cleanup distinctions; correct independently demonstrated component defects without inventing a shared cause or rewriting failed evidence.
- **FR-008**: Reconcile #354 source handoff, known issues, upgrade/rollback guidance and comparison identity without claiming a tag, staged package or public release.
- **FR-009**: Maintain native GitHub issue-level traceability and record the owner-directed timeout-policy change, including superseded guidance and per-issue functional dispositions.
- **FR-010**: Publish one reviewed source PR, monitor exact-head CI and resolve all warranted review findings in at most two Codex rounds; perform no merge, tagging, release staging or promotion.

### Key Entities

- **Timeout policy**: Explicit completion limits, their scope, former value, selected value, platform ceiling and classification.
- **Native observation**: Source/runtime/binary identity, selected budget, original result, phase timings and cleanup.
- **Release handoff**: Reviewed source, pending version, known defects, absent/present package identities and later owner decisions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every CI system and nested verification runner appears in the timeout inventory, with no unexplained incidental deadline left at its former restrictive value.
- **SC-002**: All search identity scenarios preserve their exact-target outcomes under the increased completion allowance.
- **SC-003**: Current native outcomes retain six distinct writer identities where successful and exact failure/cleanup results where unsuccessful; no historical failure is relabelled.
- **SC-004**: All required verification and exact-head hosted checks complete successfully, and no review finding remains unresolved after at most two authorized rounds.
- **SC-005**: Every included issue has an explicit functional disposition; no source-only result advertises a published v1.6.0 release.

## Assumptions

- The latest owner instruction is authoritative over prior issue advice about keeping incidental verification deadlines unchanged.
- Longer completion allowances do not alter intentional timeout semantics or scheduler performance requirements.
- A supported caller-supplied or resumed evaluation budget remains authoritative for that run; existing evidence remains immutable.
- Original #355/#356 causes may remain unconfirmed despite delivered timeout-policy and diagnostic outcomes; their individual closure must follow actual functional delivery.
- No owner workstation login/logout, installed-service operation, runtime replacement, DNS change or production credential is needed.
