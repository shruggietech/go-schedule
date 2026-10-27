# Feature Specification: Windows incident replay and causal triage

**Feature Branch**: `codex/122-windows-incident-replay`

**Created**: 2026-10-27

**Status**: Implemented

**Delivery**: Implemented on review branch codex/122-windows-incident-replay in PR #391, including hosted CI artifact correction #392; verification.md preserves exact application QA and separate corrected-script acceptance; original #355/#356 causes remain open.

**Input**: Owner-authorized S122 bundles #355 and #356, uses Spec-Kit, performs bounded investigation, implements demonstrated fixes, publishes one source PR, monitors CI and handles at most two Codex rounds. Stop before merge, tagging or promotion.

## User Scenarios & Testing

### User Story 1 - Reproduce and explain payload failures (Priority: P1)

A maintainer replays the original interpreter crash and separate prerequisite refusal on the retained and current runtime, preserving each invocation's original result and incident identity.

**Why this priority**: Runtime failure and missing prerequisites require different corrections; successful later invocations do not explain an earlier failure.

**Independent Test**: One bounded runtime campaign produces attributable outcomes for concurrent fresh/existing-state writers and prerequisite-selection contrasts without changing the installed interpreter.

**Acceptance Scenarios**:

1. **Given** retained and current runtimes, **When** a bounded replay runs, **Then** each outcome names the selected runtime, command, original exit, candidate/phase observations and available incident artifacts.
2. **Given** six simultaneous writer invocations, **When** an invocation succeeds or fails, **Then** exact literal identities and record counts remain asserted and failures remain visible without retrying into success.
3. **Given** a reproduced failure, **When** its cause is evaluated, **Then** a demonstrated project defect receives a separately tracked correction or an evidence-matched external disposition; an unreproduced original remains explicitly unresolved.

### User Story 2 - Explain owned command and inspection stalls (Priority: P1)

A maintainer replays the five originally failing service, MCP and group scenarios while retaining independent command, protocol, inspection, collection and cleanup results.

**Why this priority**: A successful mutation followed by a failed inspection must not be misreported as either a complete pass or a failed mutation.

**Independent Test**: A bounded owned replay distinguishes a controlled queued-store request from a handler-only stall and records the original five scenario outcomes with their source and process identities.

**Acceptance Scenarios**:

1. **Given** fresh owned memory and durable stores, **When** original scenarios run, **Then** their outcomes retain live incident metadata when available and explicit collection dispositions otherwise.
2. **Given** cancellation or exhausted collection allowance, **When** inspection and cleanup complete, **Then** original results, partial observations and collection failures remain independent.
3. **Given** an actual recurrent stall, **When** its private handler, pool and stack observations are analyzed, **Then** any causal conclusion cites the matching incident rather than unrelated passing tests.

### User Story 3 - Repeat a bounded investigation and act on its result (Priority: P2)

A maintainer follows a small reusable replay recipe and receives a fixed hypothesis/outcome matrix that another agent can use to implement or file precise defects.

**Independent Test**: The guide names exact commands, prerequisites, bounds and outcome fields; every investigation result maps to one original symptom and source identity.

**Acceptance Scenarios**:

1. **Given** a repeatable recipe, **When** a campaign completes, **Then** the report lists planned and observed repetitions, source/binary/runtime identities, original outcomes, incident references, cleanup and causal disposition for every selected condition.
2. **Given** an interrupted campaign, **When** its results are retained, **Then** completed observations remain available and omitted conditions are explicitly incomplete.

### Edge Cases

Unavailable retained runtime or collector; invalid or linked output directory; changed runtime identity; absent, oversized or malformed artifact; exhausted per-daemon snapshot allowance; interrupted replay; partial inspection; native runtime crash; deliberate bad prerequisite selection; successful expected refusal; unavailable native analysis tooling.

## Requirements

### Functional Requirements

- **FR-001**: Preserve the original #355 and #356 failed evidence and keep interpreter, prerequisite and IPC diagnoses separate.
- **FR-002**: Provide a repeatable bounded runtime replay using retained/current baselines, six simultaneous writers, fresh/existing schemas and metadata/writer plus implicit/explicit prerequisite-selection contrasts.
- **FR-003**: Replay all five original owned command/inspection scenario shapes against fresh memory and durable stores, retaining command, protocol, inspection, collection and cleanup independently.
- **FR-004**: Reuse the delivered process-bound and private daemon collectors; retain exact runtime/process/source identities, original exits and artifact availability without changing ordinary outcomes.
- **FR-005**: Bound each selected campaign condition to at most ten repetitions and the combined campaign execution to thirty minutes, with sequential condition groups and existing resource controls. Verification and a subsequent full evaluation are separate activities. Stop early on a causal recurrence to analyze it before further experiments.
- **FR-006**: Respect the existing per-daemon collection count and time/byte limits by using fresh owned sessions; collection failures or missing analysis tools must not rewrite original results.
- **FR-007**: Deliver a fixed hypothesis/outcome matrix with planned/observed conditions, identities, original outcomes, incident references, cleanup, causal conclusion and separately tracked correction or unresolved disposition.
- **FR-008**: Correct demonstrated in-scope project defects with discriminating regressions and actionable user-facing diagnostics; preserve access controls, six-writer assertions, partial-state semantics and the established QA report contract.
- **FR-009**: Keep all work silent and confined to owned fixtures. Do not change installed runtimes/services, account sessions, public endpoints or release identities.
- **FR-010**: Reconcile S121's owner merge and exact-main delivery records. Keep #355/#356 open when their individual causal outcomes remain incomplete; successful verification is not an issue-closure gate.

### Key Entities

- **Campaign**: Bounded set of identified conditions, source baseline, start/stop disposition and cleanup inventory.
- **Condition observation**: Original symptom/track, repetition, runtime/store/prerequisite selection, command and independent outcomes.
- **Causal disposition**: Matching incident references, supported conclusion, demonstrated defect/fix or explicit remaining uncertainty.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All selected conditions have an attributable completed, failed, unavailable or interrupted disposition; none disappear from the investigation matrix.
- **SC-002**: Every causal conclusion links the exact observed failure and identities; zero successful repetitions are presented as explanations of historical failures.
- **SC-003**: All original successful concurrent invocations retain their exact identities and expected records, and all original failures remain visible.
- **SC-004**: Campaigns stay within the declared repetition/time/resource limits and remove every disposable owned runtime or explicitly identify remaining cleanup obligations.
- **SC-005**: Every newly demonstrated defect has separate GitHub traceability and any implemented correction preserves the original contract.

## Assumptions

Merged S120/S121 source at 6f1de579bed866eb2bb328a25a0ad165caaa278a is the accepted baseline. The retained portable runtime and current Windows/WSL tooling are available for initial owned investigation. No original crash dump exists retrospectively. Native artifact analysis is attempted only when a matching actual recurrence needs it. One full source-bound Windows/WSL catalog follows actual application changes; unchanged full-catalog repetition is not an investigation prerequisite. Source push and official PR are explicitly authorized; no additional publication halt is needed.

## Clarifications

### Session 2026-10-27

- Q: How does the campaign stop? A: At ten repetitions per selected condition, thirty minutes combined, interruption, or a matching causal recurrence requiring analysis. Independent verification is outside that investigation bound.
- Q: What if the original fault does not recur? A: Retain the bounded observations, report uncertainty and leave that original issue open.
- Q: What changes are authorized? A: Owned replay tooling/recipes, demonstrated in-scope fixes and delivery records, followed by one PR and at most two Codex rounds. No merge, tag or promotion.

Clarification coverage scan: functional scope, entities/identity, maintainer flow, limits, reliability, privacy, dependencies, interruption, terminology and completion are clear. No critical ambiguity remains; decisions follow existing issue scope and owner autonomy.
