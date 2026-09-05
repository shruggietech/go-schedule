# Feature Specification: Trigger-Ready Task Authoring

**Feature Branch**: `codex/053-trigger-ready-authoring`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/053-trigger-ready-authoring`; focused functional and race suites, the full headless GUI suite, canonical format, vet, lint, race, coverage, documentation, and automation gates passed 2026-09-05.

**Input**: GitHub issues [#129](https://github.com/shruggietech/go-schedule/issues/129), [#130](https://github.com/shruggietech/go-schedule/issues/130), and [#131](https://github.com/shruggietech/go-schedule/issues/131).

## Problem Statement

Task creation currently rejects incomplete records, even when an operator needs to preserve work while researching a command or activation schedule. The interface also conflates whether a task can execute, whether it has an automatic activation source, whether it is enabled, and whether an ancestor group blocks it. This blocks the timeless task model needed by future external triggers. At the same time, sidebar destinations lack thematic grouping and simple group creation cannot be completed with Enter.

## Scope

### In scope

- Save and reopen tasks with an omitted name, command line, or time schedule.
- Display an omitted task name as `unnamed` without requiring uniqueness.
- Preserve empty groups as valid organizational containers.
- Distinguish persisted configuration, command readiness, automatic activation readiness, local enabled state, ancestor-group eligibility, and terminal lifecycle.
- Allow a task with a valid command and no automatic activation source to run manually.
- Prevent a task without a runnable command from being enabled or run.
- Keep syntactically invalid supplied commands, schedules, dates, timezones, and policies subject to validation instead of persisting unsafe malformed values.
- Present task readiness and activation state clearly in the Tasks and Groups views.
- Visually separate Tasks, Groups, and Chains from Schedule, Activity, Options, and Info while reserving the task-definition section for the future Triggers view.
- Save a valid new group when Enter is pressed in its name field.
- Deterministic API, persistence, engine, CLI, and headless GUI coverage plus canonical repository verification.

### Out of scope

- External trigger persistence, invocation, or a Triggers view from #132 and #133.
- Trigger Sets and filesystem watchers from #134 and #135.
- Persisting malformed executable or schedule syntax as arbitrary draft text.
- Nameless groups, duplicate group names, or changes to group hierarchy semantics.
- Direct Chain execution or a new workflow abstraction.
- Sidebar destination reordering beyond the requested thematic grouping.

## Clarifications

### Session 2026-09-05

- Q: Which incomplete values can be saved? -> A: Omitted values can be saved; malformed supplied values remain validation errors so the scheduler never stores executable or timing input it cannot interpret safely.
- Q: What does activation-ready mean? -> A: A task is activation-ready only when it is runnable and has an automatic source such as a schedule, startup event, completion chain, or future external trigger; a runnable task without one remains available for manual execution.
- Q: How does enabling interact with drafts? -> A: Saving never silently enables an incomplete task, and an enable request is rejected until the task is both runnable and activation-ready.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preserve Incomplete Task Work (Priority: P1)

An operator can save a task before its name, command, or timing configuration is complete, close the editor, and resume later without inventing placeholder values.

**Why this priority**: Preserving incomplete work is the central outcome and establishes the durable task model used by every later behavior.

**Independent Test**: Save tasks with each supported omission and with all supported omissions together, reload the application and daemon state, and verify every supplied value and omission survives unchanged.

**Acceptance Scenarios**:

1. **Given** a new task with no name, **when** it is saved, **then** the record persists and every list displays `unnamed` while the stored name remains empty.
2. **Given** multiple nameless tasks, **when** they are listed or edited, **then** each displays `unnamed` and remains distinguishable by stable identity.
3. **Given** a task with no command, **when** it is saved, **then** it persists disabled and is marked as not runnable with a command-specific reason.
4. **Given** a task with a valid command but no automatic activation source, **when** it is saved, **then** it persists as manual-only and can be run with Run now.
5. **Given** malformed supplied command or timing input, **when** Save is requested, **then** the record is not mutated and the offending field receives actionable validation feedback.
6. **Given** an empty group, **when** state is reloaded, **then** the group remains available for later task assignment.

---

### User Story 2 - Understand and Control Readiness (Priority: P1)

An operator can tell whether a task is runnable, automatically activation-ready, enabled, blocked by its group hierarchy, or terminal without interpreting contradictory labels.

**Why this priority**: Draft persistence is safe only when execution eligibility remains explicit and enforced consistently.

**Independent Test**: Exercise the supported combinations of command readiness, activation sources, local enabled state, group state, and lifecycle, then compare displayed status with enable and run behavior.

**Acceptance Scenarios**:

1. **Given** a task without a valid command, **when** its status is shown, **then** the interface identifies it as not runnable and enabling or Run now is refused with the same reason.
2. **Given** a runnable task without a schedule, startup source, or incoming completion chain, **when** its status is shown, **then** the interface identifies it as manual-only and Run now remains available.
3. **Given** a runnable task with an automatic activation source, **when** it is enabled and every ancestor group is enabled, **then** the interface identifies it as automatically eligible.
4. **Given** an otherwise eligible task under a disabled ancestor group, **when** its status is shown, **then** the nearest blocking group remains identified.
5. **Given** a manual-only task later becomes the target of an incoming completion chain, **when** state refreshes, **then** it becomes activation-ready without acquiring a time schedule.
6. **Given** a completed one-off task, **when** readiness is displayed, **then** terminal lifecycle remains distinct from configuration and eligibility status.

---

### User Story 3 - Navigate by Product Concept (Priority: P2)

An operator can visually distinguish task-definition destinations from operational and application destinations in the sidebar.

**Why this priority**: The grouping makes the current application easier to scan and creates a clear insertion point for Triggers without exposing unfinished functionality.

**Independent Test**: Open both themes at supported window sizes and scaling levels, then verify the two destination sections, Exit placement, keyboard order, focus, hover, and selected states.

**Acceptance Scenarios**:

1. **Given** the sidebar, **when** it is displayed, **then** Tasks, Groups, and Chains form one section while Schedule, Activity, Options, and Info form another.
2. **Given** either light or dark mode, **when** a destination is selected, hovered, or focused, **then** the section boundary and control text remain readable without exaggerated contrast.
3. **Given** keyboard navigation, **when** focus moves through the sidebar, **then** logical order matches visual order and Exit remains bottom-anchored.
4. **Given** #133 has not shipped, **when** the sidebar is displayed, **then** no nonfunctional Triggers destination appears.

---

### User Story 4 - Create a Group from the Keyboard (Priority: P2)

An operator can type a group name and press Enter to perform the same validated creation as the Create button.

**Why this priority**: A single-field creation flow should not force a pointer transition, and the behavior is small enough to ship with the related authoring improvements.

**Independent Test**: Submit valid, blank, and input-method-composition cases through Enter and compare them with button submission.

**Acceptance Scenarios**:

1. **Given** a valid group name, **when** Enter is pressed in the name field, **then** exactly one group is created and the dialog closes.
2. **Given** a blank group name, **when** Enter is pressed, **then** no group is created, the entered value is retained, and validation remains visible.
3. **Given** active input-method composition, **when** Enter confirms a composed character, **then** the dialog does not submit prematurely.
4. **Given** the Create button, **when** it is used, **then** it follows the same validation and single-submission path as Enter.

### Edge Cases

- A draft is created with every omittable task field blank.
- A previously runnable task is edited to remove its command or final automatic activation source.
- An enable request races with a task edit, completion-chain deletion, or source-task deletion that removes activation readiness.
- A task's schedule reference is absent, malformed, or missing from an older or externally modified database.
- A task has a schedule but no command.
- A task has a command but only a disabled or deleted incoming completion chain.
- Several nameless tasks sort adjacent to one another and are edited or deleted by stable identity.
- A group has no child groups and no tasks before and after restart.
- The sidebar is shorter than the combined minimum height of both sections and Exit.
- Enter key repeat, a click immediately after Enter, or delayed asynchronous completion attempts duplicate group creation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Task creation MUST accept omitted name, command, and automatic activation source values and persist the remaining supplied configuration.
- **FR-002**: An omitted task name MUST remain empty in storage and MUST display as exactly `unnamed` in every user-facing task label.
- **FR-003**: Multiple nameless tasks MUST be supported and task actions MUST continue to use stable identity rather than display name.
- **FR-004**: Empty groups MUST remain valid persisted containers and MUST survive reload without synthetic members.
- **FR-005**: The system MUST derive command readiness from whether the stored command line is present and valid for execution.
- **FR-006**: The system MUST derive automatic activation readiness from command readiness plus at least one valid automatic source: a time schedule, startup event, incoming completion chain, or future external trigger.
- **FR-007**: A command-ready task without an automatic source MUST be identified as manual-only and MUST remain runnable through Run now regardless of its local enabled state.
- **FR-008**: A task without command readiness MUST refuse Run now with an actionable command-specific validation result before process launch.
- **FR-009**: An incomplete task MUST be saved disabled even when a caller requests activation during creation.
- **FR-010**: Enabling a task MUST be rejected unless it is command-ready, automatic activation-ready, nonterminal, and backed by valid referenced configuration at the same committed state boundary.
- **FR-011**: Removing a command or final automatic source from an enabled task MUST atomically disable it.
- **FR-012**: Syntactically invalid supplied command, schedule, date, timezone, execution identity, or policy values MUST remain validation errors and MUST NOT partially mutate a stored task.
- **FR-013**: Task list presentation MUST distinguish command readiness, automatic activation state, local enabled state, ancestor-group blockage, and terminal lifecycle without contradictory labels.
- **FR-014**: Group-tree task labels MUST use the same fallback name and status vocabulary as the Tasks view.
- **FR-015**: Incoming completion chains MUST count as automatic sources for their valid target task without changing that task's time schedule.
- **FR-016**: Deleting the final incoming completion chain, directly or by deleting its source task, MUST recompute the affected target's activation readiness and disable it atomically when no automatic source remains.
- **FR-017**: Existing scheduled, startup, completion-chain, manual-run, group-cascade, overlap, catch-up, and completed one-off behavior MUST remain unchanged for fully configured tasks.
- **FR-018**: Persistence migration MUST preserve every existing task, schedule, run, alert, group, completion chain, and completion delivery while allowing tasks without a schedule reference.
- **FR-019**: Reads of a task without a schedule MUST return a valid task detail with an explicitly absent schedule instead of an internal error.
- **FR-020**: Calendar and scheduler recomputation MUST omit tasks without time schedules without logging an error or manufacturing an occurrence.
- **FR-021**: The sidebar MUST place Tasks, Groups, and Chains in a visually bounded task-definition section and Schedule, Activity, Options, and Info in a separate section.
- **FR-022**: Exit MUST remain bottom-anchored and visually separate from both destination sections.
- **FR-023**: Sidebar grouping MUST preserve destination identity, content selection, notification badges, keyboard focus order, accessible names, and readable light and dark interaction states.
- **FR-024**: The task-definition section MUST permit a future Triggers destination without displaying it before functional trigger support ships.
- **FR-025**: Pressing Enter in the new-group name field MUST invoke the same validated submission path as the Create button.
- **FR-026**: Group creation submission MUST create at most one group per user action and MUST ignore Enter used for active input-method composition.
- **FR-027**: Failed group validation MUST keep the dialog and entered value available for correction.
- **FR-028**: API, CLI, GUI, persistence, engine, migration, and interaction tests MUST cover all new state combinations and preserve compatibility for fully configured callers.
- **FR-029**: The complete change MUST pass canonical format, vet, lint, race, GUI, coverage, documentation, and automation gates.

### Key Entities

- **Task configuration**: The persisted task fields, which may omit display name, executable command, and a time schedule while retaining stable identity and safe defaults.
- **Command readiness**: A derived result stating whether the task has a valid executable command and, when false, the specific reason manual execution is unavailable.
- **Automatic activation readiness**: A derived result stating whether a command-ready, nonterminal task has at least one valid automatic source.
- **Effective eligibility**: A derived result combining readiness, local enabled state, terminal lifecycle, and ancestor-group state for automatic dispatch.
- **Automatic activation source**: A valid schedule, startup event, incoming completion chain, or future external trigger associated with a task.
- **Navigation section**: An ordered, accessible group of related sidebar destinations that does not change destination identity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One hundred percent of supported omission combinations round-trip through save, restart, list, detail, and edit without placeholder data or value loss.
- **SC-002**: One hundred percent of incomplete-command tasks are prevented from enabling and manual execution before process launch with a field-specific reason.
- **SC-003**: One hundred percent of command-ready tasks without automatic sources remain manually runnable and produce no scheduler error during reload or recomputation.
- **SC-004**: Every tested transition that adds or removes the final automatic source updates activation readiness and enabled state in the same committed operation.
- **SC-005**: Existing fully configured task fixtures retain identical schedule, lifecycle, enabled, group, and execution behavior after migration.
- **SC-006**: All user-facing task labels use `unnamed` for an omitted name while stable-ID actions remain correct across at least ten adjacent nameless records.
- **SC-007**: Both navigation sections, all destination interaction states, and the bottom Exit control remain readable and keyboard reachable in light and dark mode at three representative widths.
- **SC-008**: Each Enter action creates zero or one group, never more than one, across valid input, blank input, key repeat, composition, and click-after-key scenarios.
- **SC-009**: The canonical eight-gate repository verification completes successfully with every core package at or above 80 percent coverage.

## Assumptions

- Blank omitted values are durable draft intent, while malformed supplied values are mistakes that remain rejected.
- A display fallback does not mutate persisted user input.
- Enabled means eligible for automatic dispatch; manual Run now remains an explicit operator action governed by command readiness rather than the enabled toggle.
- Completion-chain target membership is the only non-time automatic source available during S053; external triggers join the same derived model in #132.
- Existing groups already permit zero members, so S053 protects and documents that invariant rather than adding artificial membership state.
- No new dependency is needed for state derivation, migration, sidebar grouping, or keyboard submission.
