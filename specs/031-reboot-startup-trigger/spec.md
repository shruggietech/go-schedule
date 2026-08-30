# Feature Specification: Scheduler Startup Trigger

**Feature Branch**: `codex/031-reboot-startup-trigger`

**Created**: 2026-08-30

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/031-reboot-startup-trigger`; local verification completed 2026-08-30

**Input**: Issue #64: support `@reboot` as a first-class, once-per-daemon-start task trigger across authoring, execution, history, cron interoperability, and documentation.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run work when the scheduler starts (Priority: P1)

As an operator, I can create an active task with `@reboot` or `at scheduler startup`, and it runs exactly once each time the scheduler daemon starts.

**Why this priority**: Once-per-start execution is the core capability. Without a lifecycle boundary that differs from reload, the syntax would be misleading and unsafe.

**Independent Test**: Persist one startup task, start two independent engine lifecycles, and observe exactly one startup-origin run in each lifecycle while a reload during either lifecycle adds no run.

**Acceptance Scenarios**:

1. **Given** an enabled startup task in an enabled group, **When** the daemon starts after opening its store and initializing execution, **Then** the task is dispatched exactly once with a startup run origin.
2. **Given** the same persisted task, **When** the daemon stops and starts again, **Then** it is dispatched exactly once again.
3. **Given** a running daemon, **When** its engine reloads or a startup task is created, edited, enabled, imported, moved, or reloaded, **Then** the task waits for the next daemon start.
4. **Given** a disabled startup task or one beneath a disabled group, **When** the daemon starts, **Then** no startup run is dispatched.
5. **Given** a startup run still executing, **When** a manual run is requested, **Then** the task's configured overlap policy controls the result.

---

### User Story 2 - Author and inspect startup schedules everywhere (Priority: P1)

As an operator, I can create, preview, inspect, and edit the startup schedule through the CLI, local API, and desktop Schedule field without fabricated clock occurrences.

**Why this priority**: A first-class trigger must be reachable through every supported task-authoring surface and must describe its non-clock semantics honestly.

**Independent Test**: Preview and create each canonical syntax through the API and CLI, reopen the task in the desktop editor, and confirm the canonical summary appears with an empty next-runs list.

**Acceptance Scenarios**:

1. **Given** `@reboot`, **When** it is explained or previewed, **Then** the result says `At scheduler startup`, succeeds, and contains no upcoming clock times.
2. **Given** either `@reboot` or `at scheduler startup`, **When** it is submitted for create or edit, **Then** the stored schedule is the same startup event and retains the submitted expression for editing.
3. **Given** a stored startup task, **When** it is reopened in the desktop editor, **Then** its Schedule field is populated and can be saved or changed through the existing recurring-mode workflow.
4. **Given** startup scheduling, **When** recurrence-only calendar policies are present, **Then** compatibility defaults remain harmless and no clock occurrence is invented.

---

### User Story 3 - Import, export, and convert `@reboot` faithfully (Priority: P2)

As an operator migrating cron jobs, I can import and export faithfully representable `@reboot` entries and convert them to or from the canonical human phrase without loss.

**Why this priority**: Cron interoperability is the source of the feature request and must preserve the operational context already supported for timed entries.

**Independent Test**: Dry-run and perform a crontab import containing `@reboot`, inspect the persisted task context, export it, and round-trip both string forms.

**Acceptance Scenarios**:

1. **Given** an `@reboot command` entry, **When** import runs in dry-run mode, **Then** it reports the entry as creatable and creates nothing.
2. **Given** the same entry, **When** a real import runs, **Then** command, shell, environment, stdin, run-as, group, timezone, and duplicate behavior match timed-entry semantics.
3. **Given** a faithfully exportable startup task, **When** it is exported, **Then** the timing expression is `@reboot`; operational context that a standalone line cannot preserve remains a named refusal.
4. **Given** either canonical syntax, **When** string conversion targets the other syntax, **Then** the output is exact and converts back without loss.

### Edge Cases

- A startup event is daemon-start semantics, not evidence of a physical host reboot; dependent services may not yet be ready.
- A startup schedule has no future clock occurrence, no catch-up window, no DST behavior, and no completion transition.
- A startup task created, imported, edited, enabled, or moved after startup is not retroactively dispatched in the current process.
- Multiple reload notifications coalesce without ever becoming startup notifications.
- Event rows with an unknown trigger identifier are retained but are not dispatched as startup tasks.
- Existing databases may already contain event-shaped schedule rows; opening them must not recreate the removed trigger or deduplication tables.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST recognize `@reboot` as the canonical cron form of a scheduler-daemon startup event.
- **FR-002**: The system MUST recognize `at scheduler startup` as the canonical human form of the same event.
- **FR-003**: Cron explain, conversion, task preview, task create, and task update MUST accept the startup event without generating upcoming clock times.
- **FR-004**: A stored startup schedule MUST remain an event schedule with the stable trigger identity `scheduler_startup`, a human summary of `At scheduler startup`, and the submitted source expression.
- **FR-005**: Each engine start MUST load eligible startup tasks once and dispatch each exactly once after store initialization.
- **FR-006**: Engine reload, task or group mutation, import, client reconnect, and enable operations MUST NOT dispatch a startup event in the current engine lifecycle.
- **FR-007**: A later independent daemon or engine start MUST dispatch each then-eligible startup task again; no durable fired marker or catch-up record may suppress it.
- **FR-008**: Disabled tasks, non-active tasks, and tasks beneath any disabled group MUST NOT run at startup.
- **FR-009**: Startup dispatch MUST use the existing overlap policy and worker-pool path.
- **FR-010**: Run history MUST identify startup execution with a `startup` origin distinct from `schedule`, `catchup`, `manual`, and the legacy generic `event` value.
- **FR-011**: Catch-up evaluation and clock-based schedule advancement MUST ignore startup event schedules.
- **FR-012**: Crontab dry-run and real import MUST treat `@reboot` as a supported job while preserving the same command, shell, environment, stdin, run-as, group, timezone, and duplicate semantics as timed jobs.
- **FR-013**: Cron export MUST emit `@reboot` for startup tasks when existing task-context fidelity rules permit export and MUST retain those refusals otherwise.
- **FR-014**: The persisted startup model MUST reuse the existing schedule event shape without reviving the removed trigger or deduplication tables.
- **FR-015**: Existing databases MUST reopen without destructive migration, data loss, schedule corruption, or restoration of removed trigger tables.
- **FR-016**: CLI, GUI, API, cron fidelity, import, changelog, and test-script documentation MUST explain the supported syntax, daemon-start semantics, no-next-runs behavior, and service-readiness limitation.

### Key Entities

- **Startup schedule**: A non-clock schedule whose event identity is `scheduler_startup`, with canonical human and cron forms and retained authoring expression.
- **Engine lifecycle**: One invocation of the scheduling engine's start loop. It owns a single startup-dispatch boundary that reload cannot repeat.
- **Startup run**: An append-only execution record whose `startup` origin distinguishes it from timed, catch-up, manual, and legacy event runs.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Across two independent controlled engine starts, every eligible startup task records exactly two startup runs total, while at least one reload inside a start records zero additional startup runs.
- **SC-002**: All supported authoring surfaces accept both canonical forms and return the same event meaning with zero upcoming clock times.
- **SC-003**: Dry-run import performs zero mutations; real import and export retain 100% of the operational context covered by existing timed-entry fidelity rules.
- **SC-004**: Disabled-task and disabled-group startup scenarios produce zero executions.
- **SC-005**: An existing prior-version database reopens with all prior task and schedule data intact and with removed trigger tables still absent.
- **SC-006**: The repository's eight canonical verification gates pass without reducing coverage thresholds or excluding startup lifecycle behavior.

## Assumptions

- “Startup” means scheduler daemon or engine process start, not physical host boot detection.
- The existing recurring-mode Schedule field remains the authoring location for non-one-off startup events; adding a third editor mode would add UI complexity without new capability.
- Event schedules remain active after firing because they are eligible again on the next daemon start.
- Existing recurrence policy fields remain wire-compatible on startup tasks but do not affect event timing.
- External named events, file watching, task-completion triggers, durable event deduplication, and service-readiness dependencies remain outside this slice and under the broader event-trigger roadmap.
