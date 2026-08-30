# Feature Specification: v0.9 Closure and Maintenance Automation

**Feature Branch**: `codex/029-v09-closure`

**Created**: 2026-08-30

**Status**: Implemented

**Delivery**: Review branch `codex/029-v09-closure`

**Input**: Complete the remaining actionable v0.9 maintenance work in one substantial slice: reconcile and enforce Spec-Kit lifecycle metadata (#42), add low-noise dependency update proposals (#40), enable the remaining useful supported secret-scanning controls (#39), and move the evidence-blocked Windows icon report (#33) out of the active milestone.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Trust the specification history (Priority: P1)

As the maintainer, I can scan the specification inventory and immediately distinguish delivered history from active or deferred work.

**Why this priority**: Every completed slice currently looks unfinished in at least one metadata surface, so the planning record cannot guide the next decision without a manual archaeology pass.

**Independent Test**: Audit every specification against its task list and delivery evidence, then run the lifecycle check over both the real repository and deliberately inconsistent fixtures.

**Acceptance Scenarios**:

1. **Given** a delivered specification, **When** its header and inventory entry are read, **Then** both identify it as implemented and cite verifiable delivery evidence.
2. **Given** a draft specification with no remaining task or an implemented specification with an unresolved required task, **When** the lifecycle check runs, **Then** it fails with the exact specification and contradiction.
3. **Given** an older task file whose notation or unchecked state does not match what shipped, **When** the audit is complete, **Then** its reconciliation is explicit without changing historical design decisions.

---

### User Story 2 - Receive dependency work as reviewable proposals (Priority: P1)

As the maintainer, I receive bounded, consistently labeled update pull requests for Go modules and GitHub Actions without performing routine discovery by hand.

**Why this priority**: The repository currently has no durable update-discovery channel, while the project intentionally uses pull requests for third-party review.

**Independent Test**: Validate the dependency automation contract offline, including ecosystem coverage, cadence, grouping, labels, and pull-request limits.

**Acceptance Scenarios**:

1. **Given** Go module updates, **When** the weekly dependency scan runs, **Then** compatible minor and patch updates are grouped into a bounded review proposal labeled `dependencies`.
2. **Given** GitHub Actions updates, **When** the monthly scan runs, **Then** they arrive as a separate bounded proposal with the same label.
3. **Given** a dependency proposal, **When** it opens against `main`, **Then** the existing pull-request verification workflows apply without a bypass or new approval rule.

---

### User Story 3 - Finish useful repository security controls (Priority: P1)

As the maintainer, I have the remaining useful secret-scanning protections enabled when GitHub supports them, with the actual hosted state recorded honestly.

**Why this priority**: Detection already exists, but supported prevention and validation controls remain disabled.

**Independent Test**: Read the repository security settings before and after activation and confirm push protection, non-provider patterns, and validity checks report enabled or have an explicit provider constraint.

**Acceptance Scenarios**:

1. **Given** a supported requested control, **When** S029 activates it, **Then** GitHub reports it enabled.
2. **Given** a requested control unavailable to this repository, **When** activation is attempted, **Then** the constraint is recorded instead of imitating the feature with workflow machinery.
3. **Given** unrelated optional controls, **When** S029 completes, **Then** they retain their prior state unless separately justified.

---

### User Story 4 - Keep v0.9 focused on actionable work (Priority: P2)

As the sole maintainer, I can finish the active milestone without being interrupted for a Windows install observation that I explicitly deferred.

**Why this priority**: Issue #33 has no current reproduction evidence and its prerequisite manual workflow conflicts with the maintainer's stated priorities.

**Independent Test**: Inspect issue #33 after triage and confirm it remains open, retains its verification need, and is assigned to the future milestone rather than v0.9.

**Acceptance Scenarios**:

1. **Given** issue #33 remains unverified, **When** S029 triages it, **Then** it remains open and moves to `Post-v1` with a non-blocking priority.
2. **Given** the v0.9 milestone, **When** S029 is ready for review, **Then** only the three issues delivered by this slice remain as merge-closing work.

### Edge Cases

- A specification can be implemented before its pull request is published; lifecycle state describes implementation readiness, while delivery evidence records publication separately.
- Historical tasks may include publication bookkeeping or manual observations that were superseded, waived, or completed outside the original checklist. Reconciliation must say which, not silently invent execution.
- A task file using a legacy non-checkbox format must still have auditable delivery evidence and may not make the validator silently skip the specification.
- Security API support can vary by repository plan and feature rollout. Partial success must be verified control by control.
- Dependency automation must not bundle major upgrades or actual version changes into S029.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST define a small, canonical set of specification lifecycle states and objective transition evidence.
- **FR-002**: Every specification from 001 through 029 MUST declare one allowed lifecycle state.
- **FR-003**: Every historical implemented specification MUST cite a delivery pull request, release, or commit.
- **FR-004**: Historical task metadata MUST be reconciled with delivered evidence; any non-executed item MUST retain an explicit disposition.
- **FR-005**: An offline lifecycle check MUST inspect every specification and reject missing or invalid states, missing implemented-delivery evidence, draft specifications with no actionable task, and implemented specifications with unresolved required tasks.
- **FR-006**: Lifecycle validation MUST be part of the existing automation gate and MUST have positive and negative fixture coverage.
- **FR-007**: The active Spec-Kit templates and autopilot guidance MUST keep lifecycle state current and MUST separate implementation tasks from post-implementation publication bookkeeping.
- **FR-008**: Dependency automation MUST cover Go modules at repository root on a weekly cadence.
- **FR-009**: Dependency automation MUST cover GitHub Actions at repository root on a monthly cadence.
- **FR-010**: Routine compatible dependency updates MUST be grouped to reduce review noise, while major updates and security updates remain independently reviewable.
- **FR-011**: Dependency proposals MUST carry the existing `dependencies` label and MUST use bounded open-pull-request limits.
- **FR-012**: The existing pull-request workflows MUST remain the verification path for automated dependency proposals; S029 MUST add no branch protection or approval policy.
- **FR-013**: Dependabot security updates MUST be confirmed enabled.
- **FR-014**: Secret-scanning push protection, non-provider patterns, and validity checks MUST be enabled when the repository supports them, with each final state verified through GitHub.
- **FR-015**: S029 MUST NOT enable unrelated secret-scanning controls or add substitute workflow machinery.
- **FR-016**: Issue #33 MUST remain open, retain its verification requirement, move from v0.9 to `Post-v1`, and no longer carry release-blocking priority.
- **FR-017**: The change log MUST record the lifecycle contract, dependency automation, hosted security activation, and the proportional decision to defer #33.
- **FR-018**: The review pull request for S029 MUST use `Closes #42`, `Closes #39`, and `Refs #40`; it MUST NOT close #33. Issue #40 remains open until the first generated dependency proposal is reviewed as its acceptance criteria require.
- **FR-019**: No application behavior, dependency version, branch-protection setting, release artifact, or Windows installer behavior may change in this slice.

### Key Entities

- **Specification lifecycle**: The implementation-readiness state of one feature specification plus the evidence supporting it.
- **Delivery evidence**: A pull request, release, or commit proving that an implemented specification reached the repository.
- **Resolved task disposition**: Completion, supersession, waiver, or external completion recorded against historical task metadata.
- **Dependency proposal channel**: The cadence, grouping, label, and limit governing automated update pull requests.
- **Hosted security control**: A repository setting whose requested and observed states are independently recorded.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 29 specifications have valid states and all implemented specifications have at least one verifiable delivery reference.
- **SC-002**: The lifecycle validator accepts the audited repository and rejects 100 percent of fixtures covering each prohibited contradiction.
- **SC-003**: Exactly two dependency ecosystems are configured, with no more than five open routine proposals per ecosystem and the `dependencies` label on every proposal.
- **SC-004**: All three requested secret-scanning controls report enabled, or every unavailable control has a provider response recorded.
- **SC-005**: Issue #33 remains open in `Post-v1`, retains `needs: verification`, and carries no P1 priority.
- **SC-006**: The full eight-gate verification, shell checks, UTF-8-without-BOM audit, whitespace audit, and mojibake scan pass.

## Clarifications

### Session 2026-08-30

- Lifecycle state describes implementation maturity, not whether the review branch has already merged. Publication evidence is a separate field.
- `Implemented` requires all required work resolved. A historical task may be explicitly waived or superseded, but it may not remain an unexplained open checkbox.
- Routine minor and patch updates are grouped. Major and security updates remain separate for clear risk review.
- S029 activates only the three controls named by #39. AI detection and delegated alert dismissal remain unchanged.
- Issue #33 is deferred, not closed. No Windows VM, install walkthrough, or speculative icon change belongs in this slice.

## Assumptions

- GitHub-hosted repository settings are reachable with the maintainer's existing authorization.
- Existing pull-request workflows run for Dependabot proposals without additional branch policy.
- Historical merge commits and release tags are durable enough to serve as delivery evidence.

## Out of Scope

- Dependency upgrades, automated merges, branch protection, required reviewers, or release work.
- Application, installer, scheduler, CLI, or GUI behavior changes.
- Reproducing or fixing issue #33, including any Windows VM or manual install session.
- Enabling secret-scanning AI detection or delegated alert dismissal.
