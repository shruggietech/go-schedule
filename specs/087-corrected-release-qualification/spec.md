# Feature Specification: Corrected v1.4.0 qualification

**Feature Branch**: `codex/087-corrected-release-qualification`

**Created**: 2026-09-15

**Status**: In Progress

**Delivery**: Reviewed candidate tag refreshed and fresh install/launch observed on 2026-09-15; reproduced selector repair prepared for review; [staging, repair, and qualification limitations](verification.md). Native qualification remains blocked pending reviewed repaired bytes.

**Input**: User description: "Run S087 under spec-kit autopilot, qualify corrected v1.4.0, publish a review PR, handle at most two AI review rounds, and report green CI and merge readiness."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Qualify the corrected candidate (Priority: P1)

The maintainer receives results for the reviewed corrections rather than the obsolete candidate that exposed the original UI defects.

**Why this priority**: Qualification of stale software cannot establish readiness.

**Independent Test**: Compare the tag, staging run, manifest, downloaded installer, and every observation against one exact reviewed revision and asset digest.

**Acceptance Scenarios**:

1. **Given** an unpublished obsolete candidate and explicit refresh approval, **When** staging is refreshed, **Then** every staged identity matches the reviewed revision.
2. **Given** missing refresh approval or mismatched bytes, **When** qualification begins, **Then** release mutations or passing claims are refused.

### User Story 2 - Complete fresh and upgrade observations (Priority: P1)

The maintainer obtains genuine independent fresh-install and public-v1.1.1-upgrade evidence, with S086 handling repetitive setup and diagnostic export.

**Why this priority**: Automated setup does not prove the Windows 11 user experience.

**Independent Test**: Validate the complete existing native matrix, baseline identity, required environment measurements, and candidate-bound attachments through the existing gate.

**Acceptance Scenarios**:

1. **Given** verified packages, **When** separate fresh and upgrade sessions run, **Then** setup and logs are automated while actual observations are recorded truthfully.
2. **Given** absent normal-user, multiple-profile, or high/mixed-DPI observations, **When** evidence is finalized, **Then** qualification remains incomplete.
3. **Given** a timeout or UI regression, **When** detected, **Then** logs are retained and qualification stops without a fabricated pass.

### User Story 3 - Deliver an honest review handoff (Priority: P2)

The maintainer reviews validated results without confusing merged tooling with a published release.

**Why this priority**: Public claims must follow actual release gates.

**Independent Test**: Audit PR, issue, and release records for evidence-backed outcomes and explicit remaining promotion approval.

**Acceptance Scenarios**:

1. **Given** complete validated evidence, **When** merge readiness is reported, **Then** all review findings have individual dispositions and latest-head CI is green.
2. **Given** no public-promotion approval, **When** qualification completes, **Then** the release remains a draft.

### Edge Cases

- Tag or assets change between validation and observation.
- Role inputs alias or offline prerequisites fail their recorded identity.
- Sandbox cannot represent a required environment or physical input.
- Issue #226 is closed despite unfinished publication criteria.
- Second-round findings require fixes, without a third requested review round.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Bind all artifacts and observations to one exact reviewed candidate identity.
- **FR-002**: Require explicit approval for existing-tag replacement, draft-asset refresh, hosted staging, and attended installation operations.
- **FR-003**: Preserve the entire existing fresh-install and public-v1.1.1-upgrade matrix without weakening its gate.
- **FR-004**: Complete the required native walkthrough of #229 through #233 for #228.
- **FR-005**: Reuse S086 preparation, offline prerequisites, and diagnostic export instead of adding another evidence framework.
- **FR-006**: Keep missing observations nonpassing and never treat automated logs as attended evidence.
- **FR-007**: Leave public promotion outside this slice pending separate approval.
- **FR-008**: Respond to each review finding and fix warranted defects, requesting no more than two AI review rounds.
- **FR-009**: Keep durable planning and release claims aligned with actual qualification and publication outcomes.
- **FR-010**: Repair a confirmed qualification-blocking selector-height regression from #231 with a shared intrinsic field-layout correction and browser regression coverage. Do not count observations of the existing candidate as verification of repaired bytes.

### Key Entities

- **Candidate identity**: Repository, version, reviewed revision, staging run and attempt, asset bytes and digest.
- **Qualification session**: Independent scenario, environment, observations, measurements, logs, and attachments.
- **Review handoff**: Traceable results, limitations, review disposition, and latest-head verification.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All candidate identities agree, with zero obsolete observations counted as passes.
- **SC-002**: Every required native scenario and attachment validates before qualification is declared complete.
- **SC-003**: Zero unauthorized release mutations or public-promotion actions occur.
- **SC-004**: Every review finding is dispositioned and latest-head checks are green, with no third requested round.

## Assumptions

- Scope traces to #226 and #228; no v1.5 feature work belongs in this slice.
- Reviewed source is `951d864d9683ec3bdcb1e37535ecddd67738bf13`; the obsolete v1.4.0 tag resolves to `57555ffa413df641cb21784847ad598c113bf199` and its release is still a draft.
- PR publication and the proposed existing-tag replacement, draft refresh, hosted staging, disposable qualification installs, and reopening #226 are approved by the maintainer's continuation instruction on 2026-09-15. Public promotion remains excluded.
- Closure of #226 does not prove its unfinished release criteria; it will be reopened until its release criteria are actually satisfied.

## Clarifications

### Session 2026-09-15

- Q: Which release operations are authorized? A: The maintainer directed immediate continuation after the explicit operation list; proceed with that listed scope and no further authorization pause, excluding public promotion.
- Q: Which source revision is qualified? A: The already-reviewed S086 merge commit, not an unreviewed S087 branch or previously captured obsolete observations.
- Execution finding: The fresh native walkthrough exposed inconsistent Timing mode height. A browser regression reproduced a 10.296875 px difference. S087 now includes the bounded #231 repair under the authorized defect-remediation workflow. Full release qualification remains blocked until the repair is reviewed, merged, restaged, and observed; this PR must not claim release qualification completion.
