# Feature Specification: Explicit DST Scheduling Intent

**Feature Branch**: `codex/027-dst-intent`

**Created**: 2026-08-29

**Status**: Draft

**Input**: Complete the remaining scope of issue #8 by making recurrence anchoring and daylight-saving transition behavior explicit per task across execution, persistence, API, CLI, GUI, previews, calendar, catch-up, and documentation.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Choose wall-clock or elapsed recurrence intent (Priority: P1)

As an operator, I can state whether a recurring task follows local clock readings, fixed elapsed intervals, or UTC clock readings so its behavior across offset changes matches my intent.

**Why this priority**: The scheduler currently exposes a timezone but cannot distinguish a report that must stay at 09:00 local from a backup that must run every six real hours.

**Independent Test**: Create equivalent recurrences under all three bases in `America/New_York`, enumerate runs across both 2026 DST transitions, and compare local readings and absolute gaps with the selected contract.

**Acceptance Scenarios**:

1. **Given** a daily 09:00 task using `wall_clock`, **When** the timezone offset changes, **Then** each run remains at 09:00 local and the transition gap is 23 or 25 elapsed hours.
2. **Given** an every-six-hours task using `elapsed`, **When** the timezone offset changes, **Then** every consecutive run remains exactly six elapsed hours apart and the displayed local reading shifts by one hour.
3. **Given** a daily 09:00 task using `utc`, **When** the task timezone offset changes, **Then** recurrence selection remains at 09:00 UTC and local display reflects the current offset.
4. **Given** a calendar-selecting recurrence using `elapsed`, **When** it is previewed, created, or updated, **Then** the request is refused because a variable calendar period is not a fixed elapsed duration.

---

### User Story 2 - Choose spring-gap and fall-overlap behavior (Priority: P1)

As an operator using wall-clock recurrence, I can decide whether a nonexistent local time advances or is skipped and whether a repeated local time runs at its first, both, or last occurrence.

**Why this priority**: Local clock anchoring is incomplete unless both ambiguous transition cases have explicit and testable outcomes.

**Independent Test**: Schedule tasks at 02:30 on the 2026 spring transition and 01:30 on the 2026 fall transition in `America/New_York`, then enumerate exact UTC instants for every policy.

**Acceptance Scenarios**:

1. **Given** a wall-clock task targeting nonexistent 02:30 with `next_valid`, **When** spring-forward occurs, **Then** it runs once at the first valid instant after the gap.
2. **Given** the same task with `skip`, **When** spring-forward occurs, **Then** that intended occurrence is omitted and the next ordinary recurrence remains eligible.
3. **Given** a wall-clock task targeting repeated 01:30 with `first`, `both`, or `last`, **When** fall-back occurs, **Then** it runs respectively at only the earlier instant, at both distinct instants, or only the later instant.
4. **Given** elapsed or UTC anchoring, **When** transition policies are stored, **Then** they remain available for a later basis change but do not alter the current recurrence because no local ambiguity is being resolved.

---

### User Story 3 - Configure and inspect the same policy everywhere (Priority: P1)

As an operator, I can create, edit, preview, inspect, and calendar-view DST-aware tasks through the API, CLI, and desktop editor without policy loss or contradictory next runs.

**Why this priority**: A scheduling policy is not usable if only one client can set it or if preview differs from execution.

**Independent Test**: Create and edit each policy combination through every client boundary, restart storage, and compare the task detail, preview, calendar, and engine sequences.

**Acceptance Scenarios**:

1. **Given** no DST options, **When** a task is created, **Then** it receives compatibility defaults of `wall_clock`, `next_valid`, and `first`.
2. **Given** explicit options, **When** a task is previewed or saved through CLI, API, or GUI, **Then** validation, summaries, and next runs use the same values.
3. **Given** an existing task with explicit options, **When** its schedule or an unrelated field is edited and the daemon restarts, **Then** all DST options remain unchanged.
4. **Given** a task spanning a transition, **When** detail, calendar, catch-up, or live dispatch computes occurrences, **Then** each surface returns the same ordered UTC instants without duplicates or omissions.

---

### User Story 4 - Preserve compatibility and bounded execution (Priority: P2)

As an existing operator, I can upgrade without rewriting tasks, while invalid combinations fail clearly and scheduling remains bounded.

**Why this priority**: This slice changes a safety-critical recurrence boundary and must not create silent migration or performance regressions.

**Independent Test**: Migrate a pre-S027 database, run the existing scheduling suite, exercise invalid values and combinations, and compare representative next-run benchmarks.

**Acceptance Scenarios**:

1. **Given** a database created before S027, **When** it opens under the new version, **Then** every task receives compatibility defaults and its prior wall-clock outcomes remain unchanged.
2. **Given** an unknown basis or transition policy, **When** preview, creation, or update is attempted, **Then** the request identifies the invalid field and mutates no task.
3. **Given** a valid fixed-interval or wall-clock recurrence, **When** upcoming runs are requested, **Then** evaluation terminates within established limits and never emits the same instant twice.

### Edge Cases

- Spring gaps and fall overlaps whose offset change is not exactly one hour must follow the same policy contract.
- A `next_valid` spring occurrence can collide with another intended occurrence; one UTC instant is emitted once.
- `both` must return two ordered UTC instants for one repeated wall reading, including when the cursor lies between them.
- A transition policy must compose with missing-date resolution and monthly calendar adjustments, with missing date selected before local-time resolution.
- One-off and event schedules persist the options for task consistency but do not apply recurrence anchoring.
- UTC and fixed elapsed bases do not consult the host or task timezone when selecting recurrence instants.
- Empty stored values from in-memory callers are normalized to compatibility defaults at the boundary.
- Catch-up-one still emits at most one catch-up execution even if multiple transition occurrences were missed.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Every task MUST carry a recurrence basis with values `wall_clock`, `elapsed`, or `utc`; omitted values MUST default to `wall_clock`.
- **FR-002**: `wall_clock` MUST preserve intended local calendar fields across offset changes, allowing elapsed gaps to vary.
- **FR-003**: `elapsed` MUST preserve fixed real-time spacing for supported interval recurrences and MUST reject recurrence shapes without a fixed duration.
- **FR-004**: `utc` MUST evaluate recurrence fields against UTC and MUST use the task timezone only when presenting an instant locally.
- **FR-005**: Every task MUST carry a spring-gap policy with values `next_valid` or `skip`; omitted values MUST default to `next_valid`.
- **FR-006**: Every task MUST carry a fall-overlap policy with values `first`, `both`, or `last`; omitted values MUST default to `first`.
- **FR-007**: Spring-gap and fall-overlap policies MUST affect only wall-clock recurrence resolution and MUST remain persisted when another basis makes them inert.
- **FR-008**: The resolver MUST support real IANA transitions without assuming a one-hour offset change.
- **FR-009**: When multiple intended occurrences resolve to one instant, output MUST contain that instant once; `both` MUST otherwise retain both distinct overlap instants in order.
- **FR-010**: Missing-date and calendar adjustment decisions MUST be applied before DST transition resolution.
- **FR-011**: Preview, next-run detail, calendar, catch-up, restart, and dispatch MUST share one authoritative policy-aware recurrence contract.
- **FR-012**: API task create and update requests MUST accept all three fields, validate values and basis/schedule compatibility, and return field-specific non-mutating failures.
- **FR-013**: Schedule preview MUST accept the same policy fields and return the same upcoming instants that a task created from that request would use.
- **FR-014**: CLI task add and edit MUST expose `--time-basis`, `--dst-gap`, and `--dst-overlap`; task show MUST name their effective values.
- **FR-015**: The desktop task editor MUST expose friendly controls for all three values inside Advanced Settings, prefill stored values, include them in live preview, and submit them on create and edit.
- **FR-016**: Storage MUST migrate existing tasks forward with `wall_clock`, `next_valid`, and `first`, preserve explicit choices across restart, and require no operator intervention.
- **FR-017**: One-off and event schedules MUST retain task policy values without applying recurrence-basis or transition logic.
- **FR-018**: Catch-up-one MUST remain capped at one execution even when `both` makes two missed instants eligible.
- **FR-019**: Existing tasks and callers that omit the new fields MUST retain established day-or-coarser `next_valid`/`first` wall-clock outcomes.
- **FR-020**: Documentation and human summaries MUST explain the effective basis and transition policies without claiming a five/seven-hour wall-clock gap is an elapsed-time defect.
- **FR-021**: The implementation MUST add no external service, permission, or third-party dependency.
- **FR-022**: Representative policy-aware next-run evaluation MUST remain inside the existing p99 dispatch budget and MUST NOT regress the relevant benchmark by more than ten percent without recorded justification.

### Key Entities

- **Recurrence basis**: A task-level choice describing whether recurrence fields are interpreted as local wall readings, fixed elapsed intervals, or UTC readings.
- **Spring-gap policy**: A task-level choice to advance a nonexistent wall reading to the first valid instant or omit that occurrence.
- **Fall-overlap policy**: A task-level choice to select the first, both, or last concrete instant represented by an ambiguous wall reading.
- **Scheduling policy set**: The recurrence basis, gap policy, overlap policy, and existing missing-date policy passed together through every occurrence-producing path.
- **Wall-clock intent**: Calendar components selected by a recurrence before timezone transition resolution.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A real-transition matrix covering all three bases, both spring policies, all three fall policies, and at least two IANA zones produces exactly the specified UTC instants with zero duplicates.
- **SC-002**: Fixed elapsed examples retain exact durations across both 2026 New York transitions, while equivalent wall-clock examples retain local readings and UTC examples retain UTC readings.
- **SC-003**: API, CLI, GUI, task detail, calendar, catch-up, restart, and dispatch scenarios agree on ordered next runs for every applicable policy.
- **SC-004**: A pre-S027 database migrates with 100 percent of existing tasks assigned compatibility defaults and no observed recurrence drift in regression fixtures.
- **SC-005**: Every invalid value or incompatible elapsed recurrence returns a field-specific failure and creates or updates zero tasks.
- **SC-006**: Representative next-run benchmarks remain within ten percent of the recorded baseline and within the existing p99 dispatch budget.
- **SC-007**: All eight canonical verification gates, whitespace checks, UTF-8-without-BOM checks, and mojibake audits pass.

## Clarifications

### Session 2026-08-29

- The five-hour spring and seven-hour fall gaps in an anchored six-hour local cycle are correct `wall_clock` results, not a defect. S027 preserves that default and supplies `elapsed` for exact six-hour intent.
- `elapsed` is accepted only for recurrence shapes with a fixed duration; calendar-selected monthly or yearly behavior must use `wall_clock` or `utc`.
- Transition policies are retained but inert under `elapsed` and `utc`, which avoids destructive resets when an operator changes basis.
- Missing-date resolution precedes DST resolution, and catch-up-one remains capped at one execution.
- S027 completes and closes issue #8; the already-delivered missing-date behavior remains unchanged except for composition coverage.

## Assumptions

- Existing task timezones remain presentation zones under `elapsed` and `utc` even though they do not select recurrence instants.
- Existing one-off and event behavior does not need new timing semantics.
- The current API is local and versioned in place, so additive optional fields are backward compatible.

## Out of Scope

- Changing the default away from wall-clock compatibility.
- New timezone databases, user-defined transition tables, leap-second handling, or astronomical time.
- Per-occurrence overrides, holiday calendars, or new missing-date choices.
- Reworking overlap-of-process execution policy, which is distinct from repeated local wall times.
