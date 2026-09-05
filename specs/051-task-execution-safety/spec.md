# Feature Specification: Task Execution Safety and Diagnostics

**Feature Branch**: `codex/051-task-execution-safety`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/051-task-execution-safety`; focused, race, full-suite, and canonical eight-gate verification passed 2026-09-05 for issues #102, #118, and #120

**Input**: Bundle GitHub issues [#102](https://github.com/shruggietech/go-schedule/issues/102), [#118](https://github.com/shruggietech/go-schedule/issues/118), and [#120](https://github.com/shruggietech/go-schedule/issues/120) into one post-v1 task-execution safety slice.

## Problem Statement

The desktop currently creates a task in an active state before the operator has explicitly chosen to run it, presents an individually enabled task as runnable even when a disabled group suppresses it, and reduces failed-run activity to a generic message that omits the run identity and retained process diagnostics. Together these gaps make the task lifecycle harder to control and explain.

S051 must make the safe state the creation default, distinguish configured from effective eligibility, and turn a failed-run activity record into a bounded, selectable diagnostic tied to the exact task and run. The three outcomes form a single operator journey: create safely, understand eligibility, and diagnose an execution failure.

## Scope

### In scope

- An explicit creation-time choice to activate a newly saved desktop task, cleared by default and retained through validation errors.
- Atomic creation in the chosen enabled state so an inactive task cannot become transiently eligible between separate operations.
- Existing non-desktop creation behavior remains compatible when callers do not express an enabled-state preference.
- A Tasks-table effective-state value that remains distinct from the task's own enabled flag and lifecycle state.
- Identification of the disabled direct group or nearest disabled ancestor that suppresses an otherwise eligible task.
- Exact correlation from each failed-run alert to the persisted run that caused it, including durable forward migration for existing installations.
- Activity detail that identifies task, run, trigger, outcome, exit status or launch failure, and bounded combined process output with truncation status.
- Headless, persistence, API, engine, client, GUI, race, and integration coverage for the bundled behavior.

### Out of scope

- Changing CLI task-creation defaults or requiring existing clients to send a new field.
- Separating stdout and stderr into distinct persisted streams. The existing combined capture is labeled honestly.
- Persisting adjustable table-column widths, tracked separately by #119.
- Changing group eligibility, task lifecycle, overlap, catch-up, scheduling, or manual-run policy.
- Recording arguments, standard input, environment values, or other new secret- bearing material in activity diagnostics.
- Reworking unrelated Activity entries or adding retention controls.

## Clarifications

### Session 2026-09-05

- Q: Should inactive GUI creation be implemented as create-then-disable or one atomic operation? → A: One atomic operation; a cleared choice must never create a transiently runnable task.
- Q: How should group suppression interact with the task's own state? → A: Preserve separate configured, lifecycle, and effective values; name the nearest disabled group only for an otherwise eligible task.
- Q: How should process output be presented when streams are not stored separately? → A: Label it as bounded combined stdout/stderr and disclose truncation without exposing new inputs or environment data.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Diagnose an Exact Failed Run (Priority: P1)

An operator opens a failed-run entry in Activity and sees enough retained, correlated information to distinguish a process that exited unsuccessfully from one that could not start, without searching logs or guessing which run failed.

**Why this priority**: Generic failure messages block recovery from ordinary configuration and command errors. The execution layer already retains most of the required evidence, so losing it at the presentation boundary is a direct operator-support defect.

**Independent Test**: Run controlled tasks that succeed, exit nonzero, fail to start, emit no output, and exceed the output cap. Each failed alert opens the matching task and run detail with honest status and bounded output.

**Acceptance Scenarios**:

1. **Given** two task runs fail close together, **when** either alert is opened, **then** the detail names only that alert's exact task and run.
2. **Given** a process starts and exits nonzero, **when** its detail opens, **then** the numeric exit code and retained combined output are visible.
3. **Given** a process cannot start, **when** its detail opens, **then** the absence of an exit code is identified as a launch failure and the retained diagnostic is visible.
4. **Given** captured output reached its configured cap, **when** detail opens, **then** the output remains bounded and the truncation is stated explicitly.
5. **Given** the task was deleted or detail retrieval fails, **when** the alert opens, **then** durable task/run identifiers remain visible and unavailable optional fields are described without substituting another run.

---

### User Story 2 - Create a Task Without Accidental Execution (Priority: P2)

An operator creates a desktop task in an inactive state by default and may deliberately opt in to activation as part of the same save.

**Why this priority**: A newly authored command should not become eligible to run until its author explicitly chooses that outcome.

**Independent Test**: Save otherwise identical new tasks with the activation choice cleared and selected, then verify their initial eligibility and ensure editing an existing task never changes its enabled state implicitly.

**Acceptance Scenarios**:

1. **Given** a fresh new-task dialog, **when** it opens, **then** its activation choice is visible and cleared.
2. **Given** the choice remains cleared, **when** the task is saved, **then** it is created inactive and is never transiently eligible for dispatch.
3. **Given** the operator selects activation, **when** the task is saved, **then** it is created active and participates in normal scheduling.
4. **Given** validation rejects a draft, **when** the operator corrects it, **then** the previously selected activation choice remains unchanged.
5. **Given** an existing task is edited, **when** changes are saved, **then** its enabled state is preserved unless changed through the existing explicit enable/disable action.

---

### User Story 3 - Understand Effective Task Eligibility (Priority: P2)

An operator viewing Tasks can distinguish the task's own enabled setting from its effective scheduling eligibility and can identify a disabled group that is blocking an otherwise active task.

**Why this priority**: Group cascading already prevents execution correctly, but the current row can simultaneously say Enabled and appear runnable, leaving the operator without an explanation.

**Independent Test**: Place active tasks beneath enabled, directly disabled, and ancestor-disabled groups, then mutate group membership and enablement. The effective-state value updates immediately while configured state remains true.

**Acceptance Scenarios**:

1. **Given** an enabled active task with no disabled group ancestor, **when** it is listed, **then** its effective state says it is runnable.
2. **Given** an enabled active task beneath a disabled group, **when** it is listed, **then** its configured value remains Enabled while effective state names the nearest disabled group.
3. **Given** a task is disabled itself or is not in the active lifecycle state, **when** it is listed, **then** effective state explains that condition without blaming a group.
4. **Given** a group is enabled, disabled, reparented, deleted, or the task is reassigned, **when** the live update arrives, **then** the same task row refreshes to the new effective state without an application restart.
5. **Given** a narrow window or either appearance mode, **when** a blocking value is shortened in the row, **then** its full labeled value remains available through the table's disclosure mechanism and not by color alone.

### Edge Cases

- An alert created before run correlation existed remains readable as a legacy entry and says that exact run diagnostics are unavailable.
- A correlated run that was removed with its task does not fall back to the most recent run for another task or to a timestamp guess.
- Empty process output is labeled as empty rather than omitted or confused with a retrieval failure.
- A failed process with no exit code is presented as a launch/setup failure; a process with any recorded exit code is presented as an exited process.
- Unicode and multiline output preserve line breaks and remain selectable.
- Nested disabled groups use the nearest disabled group on the path from task to root, providing the most actionable cause while preserving cycle safety.
- Missing or cyclic group references produce a conservative non-runnable explanation without hanging the interface.
- A caller that omits the optional create-state preference receives the same active-by-default behavior as before S051.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A fresh desktop new-task dialog MUST present an **Activate task after saving** choice that is cleared by default.
- **FR-002**: Saving with that choice cleared MUST create the task disabled in one authoritative operation, with no interval in which it is eligible to run.
- **FR-003**: Saving with that choice selected MUST create the task enabled and otherwise preserve existing scheduling behavior.
- **FR-004**: Validation failure MUST preserve the current activation choice, and opening a later fresh creation dialog MUST restore the cleared default.
- **FR-005**: Editing an existing task MUST NOT expose the creation-only choice or change the task's enabled value implicitly.
- **FR-006**: Existing creation callers that omit an enabled-state preference MUST retain active-by-default behavior.
- **FR-007**: The Tasks view MUST display configured enabled state, lifecycle state, and effective scheduling state as distinct labeled values.
- **FR-008**: An enabled active task beneath a disabled direct or ancestor group MUST be identified as group-blocked and MUST name the nearest disabled group.
- **FR-009**: A task that is disabled itself or is not lifecycle-active MUST receive an effective-state explanation that takes precedence over group suppression.
- **FR-010**: Effective state MUST use the existing group-chain policy rather than a second, divergent scheduling rule.
- **FR-011**: Effective state MUST refresh after live task/group creation, update, reparenting, deletion, and reassignment while preserving row identity.
- **FR-012**: Effective state and its full explanation MUST be available without relying on color, including in both appearance modes and narrow windows.
- **FR-013**: Every newly recorded failed-run alert MUST durably identify the exact persisted run that caused it.
- **FR-014**: Existing databases MUST migrate forward without losing alerts, runs, tasks, schedules, groups, or prior schema semantics.
- **FR-015**: Failed-run Activity detail MUST show task identifier, available task name, run identifier, trigger, outcome, and either numeric exit code or an explicit launch/setup failure state.
- **FR-016**: Failed-run Activity detail MUST show retained output in a selectable multiline region labeled **Combined stdout/stderr**.
- **FR-017**: Output capture MUST retain its configured byte bound and MUST persist and display whether content was truncated.
- **FR-018**: Empty output, legacy uncorrelated alerts, deleted tasks/runs, and retrieval failures MUST have explicit, non-misleading fallback text.
- **FR-019**: Run resolution MUST use the durable run identifier and MUST NOT infer identity from timestamps, task recency, or list ordering.
- **FR-020**: Diagnostics MUST NOT newly expose command arguments, standard input, environment values, or secrets beyond the already retained output.
- **FR-021**: Manual, scheduled, catch-up, startup, and completion-triggered failures MUST all retain their actual trigger and resolve correctly.
- **FR-022**: All behavioral changes MUST have deterministic regression, persistence, interface, and headless coverage and MUST pass the repository's race, coverage, documentation, and automation gates.

### Key Entities

- **Task activation intent**: The optional creation-time choice between an initially disabled and initially enabled task. Omission by legacy callers retains the existing enabled default.
- **Effective task state**: A presentation derived from the task's configured enabled value, lifecycle state, and established group-chain eligibility, including the responsible group when blocked.
- **Failed-run alert correlation**: The durable link from a scheduler failure alert to exactly one run and its task.
- **Run diagnostic**: Trigger, outcome, optional exit code, bounded combined output, and an explicit truncation indicator.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In 100 percent of fresh desktop creation flows, the activation choice begins cleared and an unchanged save produces an inactive task.
- **SC-002**: Across repeated creation and immediate scheduler reload tests, zero opt-out tasks become transiently eligible for execution.
- **SC-003**: Every tested task state (runnable, task-disabled, lifecycle- inactive, direct-group-blocked, and ancestor-group-blocked) produces one unambiguous full effective-state explanation.
- **SC-004**: One hundred percent of newly generated failed-run alerts in the supported trigger modes resolve to their exact persisted run even when failures occur consecutively.
- **SC-005**: Nonzero exit, launch failure, empty output, multiline output, and truncated output scenarios expose all required diagnostic fields without exceeding the configured output cap.
- **SC-006**: Existing stored data upgrades without loss, and legacy callers plus legacy alerts retain documented backward-compatible behavior.
- **SC-007**: The complete repository verification suite passes with every core package remaining at or above 80 percent coverage and no race or lint finding.

## Assumptions

- The existing combined stdout/stderr capture is the authoritative output for this slice; preserving stream order as two independent channels is future work.
- The existing task/group event stream is sufficient to refresh effective state.
- A missing exit code on a failed run means execution did not produce a process exit status and is presented as a launch/setup failure, including run-as setup.
- Task deletion may remove run history, so durable identifiers and honest unavailable-state text are required even when enrichment cannot be retrieved.
- Issues #102, #118, and #120 close only when their individual acceptance criteria and the bundled verification evidence pass.
