# Feature Specification: Early-crash incident evidence

**Feature Branch**: `codex/138-early-crash-diagnostics`

**Created**: 2026-11-08

**Status**: Implemented

**Delivery**: Review branch codex/138-early-crash-diagnostics; [delivery](delivery.md), [verification](verification.md); owner final review and merge pending

**Input**: Owner-approved S138 scope: improve early-crash diagnostics and record recovered historical evidence; #355/#356 explicitly remain open. Push, official PR, CI monitoring and at most two Codex reviews including automatic review are authorized; stop before owner merge.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Retain evidence when a child fails before its script starts (Priority: P1)

A maintainer using existing opt-in owned-script incident capture receives bounded Windows runtime-event evidence even if a crashing interpreter never reaches its script's identity prelude or produces a dump.

**Why this priority**: The historical crash stack is in pre-entry metadata resolution, outside the script prelude's observation window.

**Independent Test**: An owned Windows child fails before authored identity code runs; its incident retains the original native result and independently reports available matching runtime events.

**Acceptance Scenarios**:

1. **Given** an opt-in owned child with a known PID and execution interval, **When** matching Windows runtime events are available, **Then** a bounded private evidence artifact retains only those events and identifies their weaker event-based provenance.
2. **Given** absent, inaccessible, malformed or over-budget event evidence, **When** collection runs, **Then** it reports the precise collection disposition while preserving the native exit and ordinary diagnostics.
3. **Given** another process or a reused PID outside the owned interval, **When** events are evaluated, **Then** unrelated events are excluded.
4. **Given** incident capture is disabled or the interpreter never launches, **When** the fixture completes, **Then** no event query runs.

### User Story 2 - Use recovered history without inventing a cause (Priority: P2)

A maintainer sees the recovered historical runtime event separately from the original writer failure and from controlled collector demonstrations.

**Why this priority**: A temporal match narrows runtime evidence but does not identify a specific writer or establish an upstream correction.

**Independent Test**: Read the investigation record and distinguish verified event fields, missing original provenance and unresolved causal outcomes.

**Acceptance Scenarios**:

1. **Given** event 229042 near the original failure timestamp, **When** its evidence is documented, **Then** .NET 10.0.12/CoreCLR 10.0.1226.41902 is reported as a temporal association, with original writer PID attribution unproven.
2. **Given** no original dump or supported matching upstream disposition, **When** the slice completes, **Then** #355 and #356 remain open and no shared cause or runtime fix is claimed.

### Edge Cases

- Event record timestamps, provider, application identity and PID disagree with the selected child.
- No event has been flushed by the time the bounded query completes.
- The query process fails, output is truncated, or XML is invalid.
- The owned child exits ordinarily, is canceled, or never starts.
- Script identity/dump collection and event collection have different outcomes.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Existing opt-in script incident capture MUST support bounded Windows runtime-event readback for an actually launched owned child without relying on the script prelude.
- **FR-002**: Events MUST match the exact owned PID, execution interval, approved runtime/error providers and interpreter identity; unrelated or ambiguous data MUST NOT become attributed evidence.
- **FR-003**: Query duration, event count and retained bytes MUST be bounded. Missing, unavailable, failed and truncated evidence MUST remain distinguishable.
- **FR-004**: Native exit, cancellation, output/stream and cleanup results MUST remain primary and unchanged by event collection.
- **FR-005**: Collection MUST reuse the existing private incident artifact/envelope and opt-in, with hidden noninteractive Windows children, no new dependency, no elevation, and no global event/dump configuration change.
- **FR-006**: Non-Windows hosts, disabled capture and launch failures MUST retain their existing behavior without a Windows event query.
- **FR-007**: Documentation MUST distinguish newly recovered historical event evidence, controlled demonstrations, absent original evidence and unresolved #355/#356 causes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A controlled pre-entry Windows failure produces an attributable event artifact or an explicit unavailable/absent disposition without changing its native exit.
- **SC-002**: Every wrong-process, wrong-time, wrong-provider and wrong-application case is excluded from attributed event evidence.
- **SC-003**: Collection stops within its declared bounds and retains no more than the existing metadata allowance.
- **SC-004**: The original causal issues remain visibly open and historical inputs remain byte-identical.

## Assumptions

- The owner approved this diagnostic scope after preflight; causal resolution is not represented as delivered.
- Reuse existing incident capture and source verification. No new QA framework, broad unchanged replay, installed-service/runtime modification, release action or timeout-policy change is in scope.
- New event evidence can be unavailable under host policy or log retention; that condition is literal, not a passing original operation or proof that no crash occurred.
