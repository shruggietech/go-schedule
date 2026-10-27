# Feature Specification: Windows incident capture and causal resolution

**Feature Branch**: `codex/121-windows-incident-capture`

**Created**: 2026-10-27

**Status**: Implemented

**Delivery**: review branch codex/121-windows-incident-capture; local canonical eight-gate verification completed 2026-10-27 with 2403 stable inputs and 26 unchanged original artifacts; finalized source-bound Windows/WSL CLI report and independent #385/#387 corrections recorded in verification.md.

**Input**: Owner-authorized S121 bundles #355 and #356, including the post-S120 incident-capture recommendations. Implementation, source push, official PR, CI monitoring and at most two Codex rounds are authorized. Stop before merge, tag or promotion.

## User Scenarios & Testing

### User Story 1 - Explain an owned payload failure (Priority: P1)

A maintainer opts into incident collection for a disposable writer fixture and receives attributable evidence when an interpreter crashes, a prerequisite is refused or output cannot be collected.

**Why this priority**: #355's original runtime crash has no same-process crash artifact. Successful later repetitions do not explain it.

**Independent Test**: An owned failing child retains its original native outcome and yields an incident or an explicit collection-unavailable result.

**Acceptance Scenarios**:

1. **Given** an opted-in owned child, **When** it fails, **Then** a bounded incident identifies its exact executable, process, completed phase and independent collection outcome.
2. **Given** an unavailable or failed collector, **When** the child fails, **Then** the original failure remains primary and collection unavailability is actionable.
3. **Given** six simultaneous writers on a fresh or existing store, **When** they succeed, **Then** all six distinct literal identities and exactly one record per invocation remain intact.

### User Story 2 - Explain an owned inspection stall (Priority: P1)

A maintainer can distinguish time spent launching a command, serving an inspection, waiting for storage and cleaning up, using evidence captured before the owned runtime disappears.

**Why this priority**: #356's original command successes and independent inspection failures must remain distinguishable.

**Independent Test**: A deliberately stalled owned daemon produces correlated request, handler and storage observations, while the failed inspection remains incomplete.

**Acceptance Scenarios**:

1. **Given** an owned stalled inspection, **When** its declared limit is reached, **Then** bounded incident evidence describes available daemon progress before cleanup without treating incomplete state as complete.
2. **Given** cancellation or collection failure, **When** results are reported, **Then** command, inspection, collection and cleanup outcomes remain independently visible.
3. **Given** the original service/MCP/group scenarios, **When** they run under the supported limits, **Then** their real outcomes and exact source/binary identities remain attributable.

### User Story 3 - Act on an honest causal disposition (Priority: P2)

A maintainer can use one fixed incident description alongside the existing human/JSON QA report to reproduce a demonstrated defect or communicate a supported upstream disposition.

**Why this priority**: Repeated passing suites must not be presented as fixes for unrelated historical failures.

**Independent Test**: A report separates demonstrated corrections, unsupported hypotheses, unavailable observations and unreproduced historical failures.

**Acceptance Scenarios**:

1. **Given** a newly demonstrated defect, **When** corrected, **Then** its issue and affected component are named separately from #355/#356 historical uncertainty.
2. **Given** an unreproduced original failure, **When** follow-up results pass, **Then** the original failed artifact and unresolved causal outcome remain unchanged.

### Edge Cases

- Child exits before attachment; collector missing, denied, timed out or producing excessive output; a descendant retains streams after root exit.
- Two incidents occur concurrently; target identity changes; destination already exists or escapes the owned directory; a caller cancels during collection.
- Request is queued for storage, blocked in handler, canceled at transport or succeeds before inspection fails; daemon exits before collection.
- Incident output includes internal stacks or paths; arbitrary secrets, environment values and request payloads must not be published.
- Historical interpreter version matches but original same-process runtime patch or crash artifact is unavailable.

## Requirements

### Functional Requirements

- **FR-001**: Collection MUST be explicitly enabled and limited to owned disposable fixtures; ordinary successful output and installed runtime/service state remain unchanged.
- **FR-002**: Every incident MUST identify versioned format, source/target identity, incident category, triggering phase, original result, collection result and bounded artifact inventory.
- **FR-003**: Collection MUST retain exact child executable/process identity and available runtime identity without substituting a later process's inventory for the failed process.
- **FR-004**: Native crash evidence MUST be collected with an available process-scoped mechanism or explicitly marked unavailable; no global debugger registration or implicit license acceptance is permitted.
- **FR-005**: Owned-daemon stall evidence MUST distinguish request/handler progress, storage waiting and execution state before cleanup, without a new publicly reachable diagnostics service.
- **FR-006**: Collection MUST have finite time, output, artifact count and byte limits, preserve ownership checks and terminate its own collectors.
- **FR-007**: Collection failures MUST NOT change the original command exit, inspection completeness or cleanup result, retry mutations or convert a failure to success.
- **FR-008**: Existing six-writer identities/row counts, access control, durable audit outcomes, partial snapshots, resource limits and supported completion allowances MUST remain active.
- **FR-009**: Original #355 runtime/prerequisite and #356 inspection failures MUST retain independent causal dispositions and immutable failed evidence. Newly demonstrated defects receive separate issue records.
- **FR-010**: Maintainers MUST receive runnable incident collection and interpretation instructions plus fixed machine-readable outcomes that complement the existing QA envelope.
- **FR-011**: Source-bound Windows/WSL follow-ups MUST identify changed components and diagnostic overrides; passing repetitions MUST NOT imply historical causal resolution or issue completion.

### Key Entities

- **Incident**: A uniquely owned, versioned description binding trigger, target, original outcome and collection outcome.
- **Artifact**: A bounded local file with role, length and integrity identity, independently attributable to its collector.
- **Causal disposition**: Demonstrated project correction, evidence-matched upstream resolution, unavailable evidence or unreproduced historical cause, with explicit references.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every opted-in failure has one attributable incident or an explicit reason collection could not complete.
- **SC-002**: Collection terminates within its declared finite budget and never replaces the original outcome.
- **SC-003**: All six successful concurrent writers retain distinct literal identities and one record each on both starting-store conditions.
- **SC-004**: Maintainers can distinguish command failure, inspection failure, collection failure and cleanup failure from the recorded fields without inferring absent observations.
- **SC-005**: No original artifact, installed service, account session or public release identity is altered by the diagnostic workflow.

## Assumptions

- Accepted source is merged S120 main 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7. Reuse its contained capture and SQL cancellation corrections instead of reimplementing them.
- Existing runtime/prerequisite identities and private captures are retained. Reproduction is not guaranteed; unsupported upstream attribution is unacceptable.
- Existing owned fixtures and report schemas are reusable. Planning selects available collection mechanisms and exact budgets.
- GitHub Issues remain authoritative. Engineering verification is required for delivery but is not an issue closure criterion.

## Clarifications

### Session 2026-10-27

- Q: Does authorizing S121 authorize closure without causal evidence? A: No. Deliver incident capture and supported corrections; retain unresolved originals with Refs links.
- Q: Should incident collection change normal output or the installed runtime? A: No. Explicitly selected owned-fixture diagnostics only, with independent collection outcomes.
- Q: Must missing collector prerequisites block ordinary QA? A: No. Report collection unavailable and retain the original test/command outcome; do not silently accept licenses.
- Q: May source publication proceed without a second halt? A: Yes, the current owner kickoff explicitly authorizes push and official PR with at most two Codex rounds.

Coverage scan: functional, identity, lifecycle, UX/error, reliability, observability, access, dependencies, negative/recovery paths and completion are clear. Collector technology and exact limits proceed to planning.
