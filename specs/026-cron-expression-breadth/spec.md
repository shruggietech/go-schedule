# Feature Specification: General Five-Field Cron Breadth

**Feature Branch**: `codex/026-cron-expression-breadth`

**Created**: 2026-08-28

**Status**: Draft

**Input**: Expand faithful five-field cron support to lists, ranges, field-local steps, arbitrary weekday sets, and common cross-field combinations across conversion, task input, import, execution, editing, and export while retaining explicit refusals for cron day-of-month/day-of-week OR semantics and unsupported extensions.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run real-world composite cron schedules (Priority: P1)

As an operator with an existing cron expression, I can create and run one task for common schedules that select multiple minutes, hours, dates, months, or weekdays instead of splitting the expression into several tasks.

**Why this priority**: Lists, ranges, and combinations are the largest remaining gap between ordinary five-field cron and the product's interoperability claim.

**Independent Test**: Create tasks from a representative matrix of composite expressions and compare generated occurrences minute-for-minute with the parsed cron field sets across ordinary dates, month boundaries, leap day, and daylight-saving transitions.

**Acceptance Scenarios**:

1. **Given** `0 9,17 * * *`, **When** it is created as a task, **Then** one task runs every day at 09:00 and 17:00 in its selected timezone.
2. **Given** `30 8-17 * * 1-5`, **When** it is created, **Then** it runs at minute 30 of every hour from 08 through 17 on Monday through Friday.
3. **Given** `*/10 9-17 * * *`, **When** it is created, **Then** it runs every ten minutes during hours 09 through 17 and restarts the minute sequence at each hour boundary.
4. **Given** `0 0 1,15 JAN,MAR *`, **When** it is created, **Then** it runs at midnight on the first and fifteenth of January and March.

---

### User Story 2 - Inspect and import composite schedules safely (Priority: P1)

As an operator reviewing a crontab, I can explain, preview, dry-run, and import composite expressions with an exact readable description and upcoming runs before anything is created.

**Why this priority**: Broad execution support is safe only when preview and explanation describe the same recurrence the daemon will execute.

**Independent Test**: Explain and dry-run a crontab containing list, range, step, named-month, named-weekday, and cross-field examples, then compare descriptions, source identity, and previewed runs with the created tasks.

**Acceptance Scenarios**:

1. **Given** a supported composite expression, **When** `cron explain` or cron-to-human conversion is requested, **Then** the result names every restricted field without omitting or approximating values.
2. **Given** a supported composite crontab line, **When** import is previewed and then performed, **Then** both paths show the same description and recurrence and the stored task retains the normalized cron source.
3. **Given** invalid or unsupported syntax, **When** explain, preview, import, create, or edit is requested, **Then** it returns a field-specific error or named refusal and mutates no task.

---

### User Story 3 - Edit, restart, and export without semantic drift (Priority: P1)

As an operator, I can edit a composite cron task, restart the daemon, and export it later without changing its schedule or losing its source identity.

**Why this priority**: A schedule that only works at initial creation is not durable scheduler behavior.

**Independent Test**: Create and edit composite tasks through the shared CLI, API, and desktop boundaries, restart storage, compare upcoming runs, and export canonical five-field expressions that reproduce the same occurrences.

**Acceptance Scenarios**:

1. **Given** a composite cron task, **When** it is fetched or loaded into the desktop editor, **Then** the original normalized cron expression remains the editable source.
2. **Given** a stored composite task after daemon restart, **When** upcoming runs and catch-up are evaluated, **Then** they match the pre-restart recurrence and do not run before the task's schedule anchor.
3. **Given** a supported composite recurrence, **When** it is exported, **Then** one canonical five-field expression is produced with the same minute, hour, date, month, and weekday sets.
4. **Given** a non-default missing-date policy that would change a date-list schedule, **When** export is requested, **Then** export refuses rather than silently discarding the policy.

---

### User Story 4 - Preserve honest boundaries (Priority: P2)

As an operator, I receive a precise refusal for recognizable cron behavior this scheduler still cannot represent faithfully.

**Why this priority**: Breadth must not weaken the project's no-silent-approximation rule.

**Independent Test**: Exercise restricted day-of-month plus restricted day-of-week, unsupported modifier composites, Quartz fields, boot triggers, and malformed values, then verify stable refusals and zero mutation.

**Acceptance Scenarios**:

1. **Given** `0 0 13 * 5`, **When** conversion or task creation is requested, **Then** it is refused because cron uses day-of-month OR day-of-week while this scheduler's recurrence model intersects them.
2. **Given** `@reboot`, a six-field expression, Quartz `?`, or unsupported `L`, `W`, or `#` combinations, **When** conversion is requested, **Then** the existing named refusal remains.
3. **Given** a malformed range, step, name, or out-of-range value, **When** it is parsed, **Then** the error identifies the field and invalid token.

### Edge Cases

- Wildcard and range steps restart at each field boundary, including uneven minute steps such as `*/7` and stepped subranges such as `10-20/2`.
- Lists and ranges are sets: duplicates, overlapping ranges, Sunday `0`/`7`, and mixed names/numbers normalize to one canonical ordered set.
- A day-of-month value absent from a month follows the task's existing missing-date policy; imported cron defaults to `skip` for exact cron behavior.
- Month and weekday names are case-insensitive and normalize to canonical numeric export.
- A schedule created during a matching minute must only run strictly after its anchor; restart catch-up must not manufacture an occurrence before that anchor.
- A nonexistent local wall time advances to the next valid instant and an ambiguous wall time runs at its first occurrence, consistently for every selected hour and minute.
- Calendar modifiers delivered by S023 through S025 remain supported only in their documented focused shapes; mixing them with general lists, ranges, or restricted fields remains outside this slice.
- Very broad expressions such as `* * * * *` must remain bounded in memory and must not regress scheduling decision performance by more than ten percent.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product MUST accept ordinary five-field cron lists, inclusive ascending ranges, wildcard steps, range steps, and names in every field where those forms are valid.
- **FR-002**: Accepted minute and hour sets MUST support multiple daily times, range combinations, and field-local step sequences, including uneven steps, without approximating them as elapsed intervals.
- **FR-003**: Accepted day-of-month, month, and day-of-week sets MUST support multiple values and ranges, with Sunday `0` and `7` treated as the same value.
- **FR-004**: The product MUST accept cross-field combinations when at most one of day-of-month and day-of-week is restricted.
- **FR-005**: A restricted day-of-month combined with a restricted day-of-week MUST remain a named refusal because cron's OR behavior is not representable by the current single recurrence contract.
- **FR-006**: Accepted composite expressions MUST execute with their exact cron field semantics and MUST NOT depend on whether their readable description is accepted as human input.
- **FR-007**: The compiled recurrence MUST remain strictly after its schedule anchor and MUST preserve task timezone, daylight-saving, missing-date, restart, catch-up, and overlap semantics.
- **FR-008**: Cron explanation and cron-to-human conversion MUST return a stable readable description that accounts for every restricted field and does not claim to be an authorable human phrase when no equivalent phrase exists.
- **FR-009**: Existing simple expressions MUST retain their established canonical phrases and conversion output.
- **FR-010**: Crontab preview and import MUST classify every newly supported composite expression as a job, show its exact description and upcoming runs, and retain normalized cron source identity.
- **FR-011**: CLI, API, and desktop task preview, creation, and update MUST accept the same newly supported expressions through the shared schedule-input boundary.
- **FR-012**: Failed or refused preview, creation, or update MUST mutate no task.
- **FR-013**: Persisted composite schedules MUST survive storage and daemon restart without a schema change or semantic drift.
- **FR-014**: Export MUST recognize supported composite recurrence shapes and produce one canonical numeric five-field expression with ordered, deduplicated values.
- **FR-015**: Export MUST refuse when task state, source-independent recurrence behavior, calendar adjustment, phase, or missing-date policy cannot be represented faithfully in one five-field expression.
- **FR-016**: Existing human-authored schedule conversion and all previously supported cron macros, calendar selectors, ordinal weekdays, refusals, source retention, and stream/exit conventions MUST remain compatible.
- **FR-017**: Parser and compiler diagnostics MUST identify the affected field and invalid token or the exact unsupported semantic boundary.
- **FR-018**: The cron fidelity guide, CLI documentation, README claim, issue #22 inventory, and changelog MUST distinguish the newly supported standard combinations from the remaining deliberate refusals.
- **FR-019**: The implementation MUST add no external service, permission, or third-party dependency.
- **FR-020**: The scheduling hot path for representative broad expressions MUST remain within the existing p99 dispatch budget and MUST NOT regress the relevant benchmark by more than ten percent without recorded justification.

### Key Entities

- **Cron field set**: The normalized ordered values selected by one cron field, together with whether the source used wildcard or step semantics.
- **Composite cron schedule**: One five-field expression whose minute, hour, date, month, or weekday selections contain multiple values or interact across fields.
- **Compiled recurrence**: The source-independent authoritative schedule used for execution, persistence, restart, preview, and catch-up.
- **Readable cron description**: A stable explanation of every selected field; it is display output and is not treated as execution input unless it belongs to the established human grammar.
- **Canonical cron expression**: A normalized numeric five-field expression that reproduces the compiled recurrence exactly.
- **Named refusal**: A non-mutating result that identifies a recognizable but unrepresentable cron semantic.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: At least twenty representative composite expressions spanning all five fields, lists, ranges, wildcard steps, range steps, names, and cross-field restrictions produce exact expected occurrences.
- **SC-002**: Supported composite tasks produce identical upcoming runs through explain, preview, import, create, edit, restart, catch-up, and export/re-import scenarios.
- **SC-003**: A twelve-month parity matrix includes leap day, short months, both daylight-saving transitions, Sunday normalization, overlapping sets, and strictly-after-anchor behavior with zero occurrence drift.
- **SC-004**: Every defined unsupported or malformed boundary returns an actionable refusal or field error and creates or updates zero tasks.
- **SC-005**: Existing simple human and cron conversion fixtures retain their exact outputs unless a documented normalization improvement is required by this feature.
- **SC-006**: Broad-expression recurrence evaluation remains within the documented p99 budget and benchmark performance stays within ten percent of the recorded baseline.
- **SC-007**: All eight canonical verification gates, whitespace checks, and UTF-8-without-BOM and mojibake audits pass.

## Clarifications

### Session 2026-08-28

- Composite cron expressions use the same durable execution contract as human schedules; a generated English description is display metadata, not a second execution language.
- S026 covers standard five-field lists, ranges, wildcard steps, range steps, names, arbitrary weekday sets, and cross-field conjunctions, including uneven field-local steps.
- Cron day-of-month/day-of-week OR semantics, modifier composites, Quartz fields, `@reboot`, system-crontab users, and crontab environment or shell emulation remain outside this slice.
- Existing simple cron expressions keep their concise authorable phrases; broader expressions receive exact readable field descriptions without pretending the human grammar can express every combination.
- Issue #22 remains open after S026 because operational crontab fidelity and deliberately excluded dialect features remain unresolved.

## Assumptions

- The existing durable schedule representation can express the newly accepted conjunctions and persist their value sets without a migration when day-of-month and day-of-week are not both restricted.
- Imported tasks continue to default to the existing cron-faithful missing-date, catch-up, and overlap policies.
- Five-field cron retains one-minute resolution and no timezone field; the task's selected IANA timezone remains authoritative.

## Out of Scope

- Cron's restricted day-of-month OR restricted day-of-week union semantics.
- Quartz `?`, seconds or year fields, and arbitrary mixtures of `L`, `W`, or `#` modifiers.
- `@reboot`, boot events, anacron, run-parts directories, or other non-clock triggers.
- Crontab environment assignments, `CRON_TZ`, `MAILTO`, system-crontab user columns, run-as-user behavior, percent-to-stdin processing, or shell emulation.
- A free-form natural-language grammar capable of authoring every composite cron set.
- Closing epic issue #22.
