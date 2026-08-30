# Feature Specification: Task-Completion Chains

**Feature Branch**: `codex/033-task-completion-chains`

**Created**: 2026-08-30

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: review branch `codex/033-task-completion-chains`; local verification completed 2026-08-30

**Input**: User description: "Create the GitHub issues needed to define task-completion chaining, then spec and deliver the complete slice end to end under autopilot."

**Traceability**: Coordinator [#72](https://github.com/shruggietech/go-schedule/issues/72), with child issues [#73](https://github.com/shruggietech/go-schedule/issues/73), [#74](https://github.com/shruggietech/go-schedule/issues/74), [#75](https://github.com/shruggietech/go-schedule/issues/75), [#76](https://github.com/shruggietech/go-schedule/issues/76), and [#77](https://github.com/shruggietech/go-schedule/issues/77). This delivers the task-completion portion of master-spec FR-007, FR-014, and SC-006.

## Clarifications

### Session 2026-08-30

- Q: Does a completion chain replace the target task's clock schedule? → A: No. It is an additional trigger; the target keeps its existing schedule and policies.
- Q: Which source outcomes can a chain select? → A: `success`, `failure`, or `any`, where `any` means either terminal executor outcome and excludes queued or skipped bookkeeping records.
- Q: What reliability guarantee applies across daemon crashes? → A: Durable at-least-once delivery. A completed delivery never repeats, but a crash after command launch and before completion recording may replay it because the external side effect cannot be observed atomically.
- Q: What happens when a target is ineligible when delivery is claimed? → A: The delivery is resolved without execution and recorded diagnostically; disabling a target does not accumulate an unbounded deferred queue.
- Q: May completion chains form cycles? → A: No. Self-links and direct or indirect cycles are refused before mutation.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run One Task After Another (Priority: P1)

As a local operator, I can connect a source task to a target task so the target runs after the source finishes with the selected outcome, without replacing either task's normal schedule.

**Why this priority**: This is the core missing v1 outcome. Without an actual completion-triggered run, the remaining management and presentation work has no product value.

**Independent Test**: Create two enabled tasks and one success chain, complete the source successfully, and observe one target run whose history names the completion origin and source.

**Acceptance Scenarios**:

1. **Given** an enabled success chain, **When** the source succeeds, **Then** the target is dispatched once through its normal overlap policy.
2. **Given** failure and any-outcome chains, **When** the source fails, **Then** both matching targets run and a success-only target does not.
3. **Given** a completion-triggered target that also has a timed schedule, **When** either trigger becomes due, **Then** both paths remain valid and are distinguishable in history.
4. **Given** a chain from A to B and another from B to C, **When** A succeeds and each target succeeds, **Then** the finite chain progresses A to B to C.

---

### User Story 2 - Recover Delivery Safely (Priority: P1)

As an operator, I can trust a completion event to survive daemon interruption and to avoid duplicate scheduling decisions, with the unavoidable external-command replay window explained honestly.

**Why this priority**: Reliability through interruption is part of the scheduler's core promise and the master specification's at-least-once requirement.

**Independent Test**: Persist a pending and a claimed delivery, reopen the scheduler, and prove pending work is dispatched, completed work is not, and ambiguous claimed work follows the documented replay policy.

**Acceptance Scenarios**:

1. **Given** a source run and matching chains, **When** the source completion is recorded, **Then** its run record and unique delivery records become durable together.
2. **Given** a pending delivery, **When** the daemon restarts, **Then** the delivery is claimed and dispatched.
3. **Given** a completed delivery, **When** the daemon restarts, **Then** it is never dispatched again.
4. **Given** a delivery claimed before an unclean stop, **When** the daemon restarts, **Then** it becomes eligible for at-least-once replay and the recovery is observable.

---

### User Story 3 - Manage Chains from CLI and API (Priority: P2)

As a local operator or client, I can list, create, change, inspect, and delete completion chains with consistent validation and stable machine-readable output.

**Why this priority**: A reliable engine still needs a complete, scriptable authoring boundary.

**Independent Test**: Exercise the chain lifecycle through the CLI and local API, including JSON output and invalid mutations, without direct database access.

**Acceptance Scenarios**:

1. **Given** valid source and target task identifiers, **When** a chain is created, **Then** both identifiers, names, selected outcome, and stable chain identity are returned.
2. **Given** an existing chain, **When** its source, target, or outcome is changed validly, **Then** later completions follow only the updated relationship.
3. **Given** an invalid identifier, duplicate, self-link, outcome, or cycle, **When** mutation is attempted, **Then** it fails without partial state and explains the correction.
4. **Given** a caller requesting JSON, **When** any read operation succeeds, **Then** the output is stable structured data rather than formatted prose.

---

### User Story 4 - Manage Chains Visually (Priority: P2)

As a desktop user, I can understand and manage task-completion chains using task names, without knowing internal identifiers or using the command line.

**Why this priority**: The desktop is the approachable interface and must not lag a core scheduling capability.

**Independent Test**: In a headless desktop test, create, update, and delete a chain through the chain view and observe live refresh plus actionable invalid/error states.

**Acceptance Scenarios**:

1. **Given** tasks and no chains, **When** the chain view opens, **Then** it explains the feature and offers creation.
2. **Given** valid distinct source and target selections, **When** the user saves, **Then** the relationship appears using task names and its selected outcome.
3. **Given** a relationship or task change elsewhere, **When** the live event arrives, **Then** the view refreshes without restarting.
4. **Given** a backend error or stale task, **When** a mutation fails, **Then** existing state remains and an actionable message appears.

---

### User Story 5 - Diagnose and Understand Chained Runs (Priority: P2)

As an operator, I can distinguish a completion-triggered run from scheduled, startup, catch-up, and manual runs and trace it back to the source task and run.

**Why this priority**: Automation chains become dangerous when history cannot explain why a command ran.

**Independent Test**: Complete a source, inspect the target through run history and activity surfaces, and follow its source task/run correlation using documented fields.

**Acceptance Scenarios**:

1. **Given** a chained run, **When** history is listed as text or JSON, **Then** it identifies `completion` origin plus source task and source run.
2. **Given** a non-completion run, **When** history is inspected, **Then** existing trigger labels and output remain compatible and no empty completion metadata is advertised.
3. **Given** public documentation, **When** an operator follows the CLI or desktop walkthrough, **Then** they can create and validate a chain without reading implementation documentation.

### Edge Cases

- A source completion matching several relationships creates one independent delivery for each relationship.
- Two relationships with the same source, target, and outcome are duplicates; relationships with different outcomes are distinct.
- `any` overlaps `success` or `failure` intentionally, so both matching relationships may dispatch the same target and its overlap policy decides concurrency.
- Queued and skipped bookkeeping records never become completion sources; only executor success or failure does.
- Self-links and longer cycles are refused for create and update, including a cycle introduced after several valid links already exist.
- Deleting a source or target removes its relationships and unresolved deliveries without affecting unrelated tasks, runs, or chains.
- A target disabled directly or through an ancestor group when delivery is claimed does not run; the delivery is resolved with diagnostic context.
- A target deleted between delivery creation and claim is resolved safely without retry loops.
- An ambiguous claimed delivery at unclean shutdown may replay. Operators must make externally side-effecting commands idempotent when crash-safe exactly-once effects matter.
- Completion cascades remain finite because the relationship graph is acyclic; fan-out and converging paths remain supported.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST let users define a completion chain with a stable identity, one existing source task, one distinct existing target task, and one outcome selector: success, failure, or any.
- **FR-002**: A completion chain MUST be an additional trigger and MUST NOT replace or mutate either task's stored schedule, timezone, catch-up policy, overlap policy, or group membership.
- **FR-003**: Only terminal executor success and failure records MUST match completion chains; queued, skipped, catch-up bookkeeping, and non-run events MUST NOT masquerade as source completions.
- **FR-004**: Recording a terminal source run MUST durably create at most one delivery for each matching relationship and source-run identity, without exposing a committed run whose matching delivery was silently omitted.
- **FR-005**: The system MUST provide durable at-least-once processing for pending and interrupted deliveries; completed deliveries MUST NOT replay, and the ambiguous post-launch crash window MUST be documented.
- **FR-006**: The target MUST dispatch through existing worker bounds and overlap handling, and its run MUST carry completion origin, source task identity, and source run identity.
- **FR-007**: The target's terminal completion MUST be eligible to create downstream deliveries, allowing finite multi-step chains.
- **FR-008**: Self-links and any create or update that would introduce a direct or indirect cycle MUST fail atomically with an actionable validation error.
- **FR-009**: A target that is missing, inactive, directly disabled, or disabled through an ancestor group when claimed MUST not execute, MUST not retry indefinitely, and MUST leave diagnostic evidence.
- **FR-010**: Deleting a task MUST remove dependent chains and unresolved deliveries without deleting unrelated chains or historical run records; deleting a chain MUST remove its unresolved deliveries.
- **FR-011**: The local API MUST support listing, creating, retrieving, updating, and deleting completion chains with stable success and validation response shapes.
- **FR-012**: The CLI MUST support the same lifecycle, follow existing command structure and exit conventions, and provide both readable text and stable JSON read output.
- **FR-013**: Chain mutations MUST publish live events so connected clients can update without polling or restart.
- **FR-014**: The desktop MUST provide a dedicated chain-management view using task names for create, edit, delete, empty, stale-data, and backend-error flows.
- **FR-015**: Run-history text, JSON, API, desktop activity, and structured logs MUST use consistent completion-origin and source-correlation terminology.
- **FR-016**: Existing time-based, one-off, startup, manual, catch-up, overlap, grouping, alerts, run history, migrations, local access control, and event streaming MUST remain compatible.
- **FR-017**: The schema change MUST be forward-only and non-destructive and MUST work for clean databases and databases that previously created and dropped the legacy trigger tables.
- **FR-018**: Restart, duplicate, fan-out, cascade, overlap, disablement, deletion, validation, authorization-boundary, GUI, race, coverage, and dispatch-latency behavior MUST have deterministic automated verification.
- **FR-019**: Public CLI, API, GUI, architecture, and scheduling documentation MUST give consistent workflows and MUST distinguish delivered task completion from deferred external named events and file watching.
- **FR-020**: The completed pull request MUST use closing keywords for coordinator #72 and child issues #73 through #77 because the slice is not complete while any child acceptance criterion remains unmet.

### Key Entities

- **Completion Chain**: A stable relationship from one source task to one distinct target task, with a success, failure, or any outcome selector and creation/update timestamps.
- **Completion Delivery**: Durable processing state for one chain matched to one immutable source run. It records pending, claimed, completed, or resolved-without-execution lifecycle plus target/source correlation and recovery evidence.
- **Run**: Existing execution history extended for completion-origin runs with optional source task and source run identities. Non-completion records leave those fields absent.
- **Task**: Existing scheduled command. It may participate in multiple incoming and outgoing chains while retaining its independent schedule and policies.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can create a two-task success chain and observe the correlated target run in under two minutes using either the CLI or desktop workflow.
- **SC-002**: One terminal source run produces exactly one durable delivery per matching relationship across at least 100 repeated duplicate-insertion attempts.
- **SC-003**: Pending and interrupted deliveries recover after restart, while 100 completed-delivery restart checks produce zero repeated dispatches.
- **SC-004**: Every attempted self-link or cycle across chains of at least 100 tasks is rejected before mutation.
- **SC-005**: A fan-out source with 100 matching chains produces 100 bounded delivery decisions without exceeding the existing 100 ms p99 dispatch-latency budget under the nominal benchmark setup.
- **SC-006**: All completion-triggered history records identify their source task and source run; all non-completion history records remain free of misleading source metadata.
- **SC-007**: CLI, API, and desktop lifecycle tests cover create, read/list, update, delete, invalid, empty, and backend-error paths.
- **SC-008**: Prior-schema databases reopen with all existing tasks, schedules, runs, groups, and alerts intact and usable after migration.
- **SC-009**: Full format, vet, lint, race, GUI, coverage, documentation, and automation gates pass with every core package at or above 80 percent coverage.
- **SC-010**: No new unbounded queue, polling loop, goroutine leak, network exposure, or external runtime dependency is introduced.

## Assumptions

- This is a local single-machine capability governed by the existing local IPC authorization boundary.
- Task commands can have external side effects; exactly-once external effects are impossible across the documented ambiguous crash window, so at-least-once delivery is authoritative.
- Existing tasks retain one normal schedule; completion chains are separate additional relationships rather than another schedule kind.
- A chain has no independent enabled flag. Users remove it to stop future delivery or edit it to change the condition.
- The source outcome selector applies to executor success/failure only. “Any” is not a wildcard for queued, skipped, or cancelled bookkeeping records.
- External named events, payloads, file/folder watching, remote execution, and notification channels remain outside S033 and tracked separately where applicable.
