# Feature Specification: S120 Windows causal diagnostics

**Feature Branch**: `codex/120-windows-causal-diagnostics`

**Created**: 2026-10-27

**Status**: Implemented

**Delivery**: [PR #383](https://github.com/shruggietech/go-schedule/pull/383), review branch codex/120-windows-causal-diagnostics; reviewed implementation cb48f4813d87fde27c245b53e807604e609d0562. Owner-authorized source publication; stop before merge, tagging or promotion.

**Input**: Implement the full #355/#356 bundle under autopilot using targeted causal investigation, existing isolated fixtures and truthful individual dispositions.

## Clarifications

### Session 2026-10-27

- Q: Must an unconfirmed original cause be closed to complete a useful S120 increment? A: No. Investigate every original category and deliver supported corrections or explicit remaining scope; keep any unmet issue outcome open.
- Q: May diagnostics replace operation results or change successful output? A: No. Collection is separately attributed, bounded and confined to owned fixtures; original exits, successful output and assertions stay authoritative.
- Q: Which existing work must be reused? A: The delivered six-writer/probe diagnostics, partial snapshots, request timings, cancellation recovery and supported 120-second/two-hour QA defaults.
- Q: Is source publication already approved? A: Yes. One PR and in-scope review fixes are authorized, with at most two Codex rounds and no merge, tag or promotion.

## User Scenarios & Testing

### User Story 1 - Explain runtime and prerequisite failures (Priority: P1)

A maintainer investigating a failed script invocation can distinguish interpreter failure, prerequisite discovery refusal and application writer failure, retain the original outcome and determine the supported corrective action.

**Why this priority**: #355 contains two independent historical failures whose improved diagnostics and successful repetitions have not explained their causes.

**Independent Test**: Each targeted experiment identifies its source, exact inputs, runtime and original result; demonstrated corrections preserve distinct concurrent session records.

**Acceptance Scenarios**:

1. **Given** the original interpreter crash or prerequisite refusal evidence, **When** the maintainer evaluates a specific causal hypothesis, **Then** original evidence, runtime identity, experiment conditions and the supported conclusion remain distinguishable.
2. **Given** a demonstrated project defect, **When** the correction is applied, **Then** successful concurrent invocations retain one distinct record each and failed invocations retain their primary failure.
3. **Given** an upstream-only defect or an unreproduced symptom, **When** the investigation reports its disposition, **Then** it provides evidence-matched supported guidance or explicitly retains the unconfirmed cause without claiming a fix.

### User Story 2 - Locate an owned IPC or observation stall (Priority: P1)

A maintainer can attribute an intermittent operation stall to the affected request, process and execution phase rather than confusing a successful command with failed inspection or cleanup.

**Why this priority**: #356 remains open after partial observations, phase timings, cancellation recovery and supported completion budgets were delivered.

**Independent Test**: Targeted original service/MCP/group request shapes retain complete ownership and separate command, inspection and cleanup results; any added causal diagnostic is bounded and tied to the owned process.

**Acceptance Scenarios**:

1. **Given** a stalled owned operation, **When** existing timings cannot discriminate its cause, **Then** narrowly scoped additional evidence identifies the active component without changing the original result.
2. **Given** cancellation, incomplete inspection or process cleanup failure, **When** evidence is retained, **Then** preceding partial data survives and incomplete state never satisfies complete-state assertions.
3. **Given** a demonstrated affected component, **When** corrected, **Then** operations and independent observations remain responsive under their declared limits without blind mutation retries or relaxed assertions.

### User Story 3 - Review attributable outcomes (Priority: P2)

Maintainers receive a source-bound investigation and implementation record for each issue while the original failures and public release remain unchanged.

**Why this priority**: Separate causes and release identities must remain reviewable across a bundled slice.

**Independent Test**: The source record maps each issue to experiments, delivered corrections and remaining scope, and reconciles S119's actual merge without reassigning earlier evidence.

**Acceptance Scenarios**:

1. **Given** mixed passing and failing experiments, **When** the report is produced, **Then** every native outcome retains its identity and no successful repetition substitutes for causal evidence.
2. **Given** source publication authorization, **When** the implementation is reviewed, **Then** warranted findings are addressed within two Codex rounds and the owner receives an unmerged PR with successful applicable checks.

### Edge Cases

- Historical logs exist without a crash dump or a current reproduction.
- A collector is unavailable, exceeds its bounded allowance or perturbs the observation.
- An interpreter exits before collection, or a failed prerequisite probe leaves inherited output handles open.
- A command succeeds but its after-state inspection or cleanup fails.
- Cancellation is observed by the client before server/store work finishes.
- A runtime, dependency overlay or diagnostic build differs from the published package.
- Original failure shapes pass on Windows or Linux without identifying a cause.
- A diagnostic reveals a new independent defect requiring its own issue.

## Requirements

### Functional Requirements

- **FR-001**: Preserve original #355/#356 failed evidence and exact source, binary and runtime attribution; successful repeats MUST NOT be represented as causal fixes.
- **FR-002**: Investigate the interpreter crash and prerequisite refusal as separate hypotheses, with bounded targeted experiments and individually supported dispositions.
- **FR-003**: Correct demonstrated in-scope project defects, or document evidence-matched supported upstream guidance; unreproduced original causes MUST remain explicitly unresolved.
- **FR-004**: Preserve six concurrent writers, fresh/existing schemas, distinct literal identities, exact row counts and original process outcomes.
- **FR-005**: Target original owned service/MCP/group and metadata operation shapes and distinguish client, transport, server/storage, observation and cleanup progress.
- **FR-006**: Reuse delivered partial snapshots, phase/request diagnostics, cancellation recovery and supported budgets rather than duplicating them.
- **FR-007**: Any added diagnostic MUST be bounded, opt-in or failure-triggered, scoped to owned fixtures, and retain collection status separately from the original operation result.
- **FR-008**: Diagnostics MUST exclude real credentials and arbitrary private payloads, create no public listener, and preserve ownership and process/resource limits.
- **FR-009**: Preserve command exit/output, incomplete observations and original cleanup failures separately; incomplete state MUST NOT satisfy complete-state assertions.
- **FR-010**: Compare environments or dependencies only to discriminate an explicit hypothesis, with changed conditions and diagnostic-only identity recorded.
- **FR-011**: File or update a separate attributable GitHub bug for each newly demonstrated defect; retain issue-level outcome mapping and existing milestones/labels.
- **FR-012**: Reconcile S119's completed PR #380 merge and exact-main acceptance as historical delivery metadata, preserving immutable public v1.6.0 identities.
- **FR-013**: Follow the installed Spec-Kit sequence and blocking analysis; complete all eight foreground gates before any source commit.
- **FR-014**: Publish the authorized source PR, address every warranted review finding within at most two Codex rounds and stop before merge, tagging or promotion.

### Key Entities

- **Original failure**: Immutable source/runtime/operation identity, native result and retained evidence.
- **Causal experiment**: Named hypothesis, exact controlled inputs, environment difference, owned target, bounded diagnostics and observed result.
- **Diagnostic attachment**: Attributed collection result, bounded content and original operation association.
- **Issue disposition**: Demonstrated project correction, evidence-supported upstream resolution or explicitly unconfirmed remaining cause.
- **Source delivery**: S120 implementation identity, original evidence references and separate S119/public-release history.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Both issues receive distinct evidence-supported dispositions covering every original failure category, with no unexplained source or result substitution.
- **SC-002**: All diagnostic outcomes remain attributable to one owned target and retain the original operation outcome, including collection failure or cancellation.
- **SC-003**: Demonstrated corrections preserve every existing writer, state-completeness and ownership guarantee.
- **SC-004**: The implemented source has successful applicable checks and no unresolved warranted review findings within two Codex rounds, remaining unmerged for the owner.

## Assumptions

- Existing owned Windows/WSL fixtures and retained evidence are reused; an unchanged full-catalog replay is not a substitute for hypothesis-directed work.
- No installed-service, login/logout, system-wide diagnostic configuration or owner-runtime replacement is required by default.
- A missing native crash artifact cannot be invented; upstream attribution requires a matching signature and supported evidence.
- Public v1.6.0 source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b and S119 merge dac2e97895353435ad30515402b9f99d5e5d4ebc remain separate from S120.
