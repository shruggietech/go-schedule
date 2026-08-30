# Feature Specification: Monthly Calendar Cron Parity

**Feature Branch**: `codex/025-monthly-calendar-parity`

**Created**: 2026-08-28

**Status**: Implemented

**Delivery**: [PR #60](https://github.com/shruggietech/go-schedule/pull/60)

**Input**: Deliver one substantial monthly-calendar work slice covering cron day-of-month `L`, `W`, and `LW` semantics end to end, rather than splitting each syntax token into a separate review cycle.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Import monthly calendar selectors (Priority: P1)

As an operator with existing cron schedules, I can explain, preview, import, and use last-day, nearest-weekday, and last-weekday-of-month expressions without rewriting them by hand.

**Why this priority**: These related calendar selectors are common monthly scheduling concepts and belong in one coherent interoperability outcome.

**Independent Test**: Explain and import representative `L`, `nW`, and `LW` expressions, then compare their generated runs with hand-checked calendar dates across weekday, weekend, short-month, leap-year, and DST boundaries.

**Acceptance Scenarios**:

1. **Given** `0 9 L * *`, **When** it is explained or converted, **Then** the result is `last day of every month at 09:00`.
2. **Given** `0 9 15W * *`, **When** it is explained or converted, **Then** the result is `nearest weekday to the 15th of every month at 09:00`.
3. **Given** `0 9 LW * *`, **When** it is explained or converted, **Then** the result is `last weekday of every month at 09:00`.
4. **Given** a crontab job using any supported monthly calendar selector, **When** it is previewed or imported, **Then** it remains a job, shows its readable phrase, and retains its original cron source and syntax.

---

### User Story 2 - Author and export monthly calendar schedules (Priority: P1)

As an operator authoring schedules in plain language, I can create last-day, nearest-weekday, and last-weekday schedules and export them as canonical five-field cron without changing when they run.

**Why this priority**: Bidirectional support is the product promise, and these selectors must work consistently through the CLI, API, desktop editor, and export workflow.

**Independent Test**: Author each native phrase, export it, re-import it, and compare every generated run over at least one year containing short months, a leap-day boundary, weekends, and DST transitions.

**Acceptance Scenarios**:

1. **Given** `last day of every month at 09:00`, **When** it is converted or exported, **Then** the result is `0 9 L * *`.
2. **Given** `nearest weekday to the 15th of every month at 09:00`, **When** it is converted or exported, **Then** the result is `0 9 15W * *`.
3. **Given** `last weekday of every month at 09:00`, **When** it is converted or exported, **Then** the result is `0 9 LW * *`.
4. **Given** a supported phrase in CLI, API, or desktop schedule input, **When** it is previewed, created, edited, or exported, **Then** all interfaces preserve the same source identity and run times.

---

### User Story 3 - Preserve exact calendar behavior (Priority: P1)

As an operator, I can trust weekend adjustment, month-boundary behavior, missing-date policy, timezone handling, and daylight-saving resolution to remain explicit and faithful.

**Why this priority**: Calendar selectors are useful only if boundary dates do not silently move or duplicate runs.

**Independent Test**: Evaluate every selector against calendar matrices that include months ending on each weekday, target dates falling on every weekday, February in leap and non-leap years, and both daylight-saving transitions.

**Acceptance Scenarios**:

1. **Given** an `nW` target that falls Monday through Friday, **When** the schedule runs, **Then** it uses that date unchanged.
2. **Given** an `nW` target that falls Saturday, **When** the schedule runs, **Then** it uses the preceding Friday unless that would leave the target month, in which case it uses the following Monday.
3. **Given** an `nW` target that falls Sunday, **When** the schedule runs, **Then** it uses the following Monday unless that would leave the target month, in which case it uses the preceding Friday.
4. **Given** `LW`, **When** a month ends on Saturday or Sunday, **Then** it uses the final Friday; otherwise it uses the final calendar day.
5. **Given** an `nW` target date that does not exist in a month, **When** the task uses `skip`, `last_valid`, or `next_valid`, **Then** the existing missing-date policy first resolves or skips the intended date and the nearest-weekday adjustment is applied to the resulting valid date without producing duplicate runs.

---

### User Story 4 - Receive precise boundary feedback (Priority: P2)

As an operator, I receive an actionable error or named refusal for malformed modifiers or combinations that this focused contract cannot represent faithfully.

**Why this priority**: Adding richer syntax must preserve the no-silent-approximation guarantee.

**Independent Test**: Exercise invalid dates, bare modifiers, offsets, lists, ranges, steps, mixed modifiers, restricted months or weekdays, and selector-rich native rules, then verify clear refusal and zero task mutation.

**Acceptance Scenarios**:

1. **Given** `0W`, `32W`, bare `W`, malformed `LW`, or unsupported `L-n`, **When** conversion is requested, **Then** it returns a field-specific error or named refusal.
2. **Given** a list, range, step, multiple modifier, or mixed day-of-month selector, **When** conversion is requested, **Then** it is refused rather than simplified.
3. **Given** a supported selector combined with a restricted month or day-of-week, **When** conversion is requested, **Then** the unsupported combination is refused without approximation.
4. **Given** a failed or refused selector during task creation or update, **When** the request completes, **Then** no task or existing schedule is mutated.

### Edge Cases

- `1W` on a Saturday resolves to Monday the 3rd rather than crossing into the previous month.
- A target on the final Sunday resolves to the preceding Friday rather than crossing into the next month.
- `LW` uses Friday when the final calendar day is Saturday or Sunday.
- February 29 is present in leap years and absent in other years.
- Targets 29 through 31 can be absent, so missing-date policy can affect native nearest-weekday schedules even though imported cron defaults to skip.
- A `next_valid` resolution can enter the following month; the scheduler must retain its existing collision and duplicate-suppression behavior.
- Day-of-week `weekdayL` and `weekday#ordinal` support from S023 and S024 must remain unchanged.
- `L-n`, `L-nW`, arbitrary lists/ranges/steps containing modifiers, Quartz `?`, and six-field expressions remain outside this slice.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product MUST accept day-of-month `L` in a five-field cron expression as the final calendar day of every month.
- **FR-002**: The product MUST accept exactly one day-of-month `nW` selector where `n` is 1 through 31.
- **FR-003**: The product MUST accept day-of-month `LW` as the final Monday-through-Friday date of every month.
- **FR-004**: Supported `L`, `nW`, and `LW` import expressions MUST keep month and day-of-week unrestricted.
- **FR-005**: Explain and cron-to-human conversion MUST produce one canonical readable phrase for each supported selector while preserving time.
- **FR-006**: The shared schedule grammar MUST accept the three canonical phrase families with optional explicit time using existing time and timezone conventions.
- **FR-007**: Crontab preview and import MUST classify supported selector lines as jobs, show the canonical phrase, and retain original cron source and syntax.
- **FR-008**: Task preview, creation, and update through CLI, API, and desktop interfaces MUST accept the same canonical phrases and supported cron expressions through the shared input boundary.
- **FR-009**: Native last-day, nearest-weekday, and last-weekday schedules MUST export to canonical `L`, `nW`, and `LW` five-field expressions when their complete behavior is representable.
- **FR-010**: Nearest-weekday adjustment MUST preserve a valid weekday target, move Saturday to Friday and Sunday to Monday, and reverse direction at a month boundary so the adjusted date remains in the intended month.
- **FR-011**: Last-weekday adjustment MUST select the final calendar day when it is Monday through Friday and the preceding Friday otherwise.
- **FR-012**: For an absent `nW` target, the existing missing-date policy MUST resolve or skip the intended calendar date before nearest-weekday adjustment, while preserving existing duplicate suppression.
- **FR-013**: Export MUST allow any missing-date policy for `L`, `LW`, and `nW` targets 1 through 28 because those selectors always resolve within every month.
- **FR-014**: Export of `29W`, `30W`, or `31W` MUST require effective skip behavior; a non-skip missing-date policy MUST receive a named refusal rather than being discarded.
- **FR-015**: Round trips MUST preserve run times across at least twelve consecutive months including a leap-year boundary, both daylight-saving transitions, weekend targets, and short months.
- **FR-016**: Bare or malformed modifiers, invalid target dates, offsets, lists, ranges, steps, multiple or mixed modifiers, restricted months, restricted day-of-week, and selector-rich native schedules MUST produce field-specific errors or named refusals without approximation.
- **FR-017**: Failed or refused task creation and update MUST NOT mutate stored tasks or schedules.
- **FR-018**: Existing supported cron expressions, `weekday#ordinal`, day-of-week `weekdayL`, named refusals, source-retention behavior, and CLI text/JSON stream conventions MUST remain compatible.
- **FR-019**: The cron fidelity guide and CLI documentation MUST distinguish supported day-of-month `L`, `nW`, and `LW` from still-declined offsets, mixed forms, Quartz syntax, and arbitrary combinations.
- **FR-020**: The issue #22 parity inventory and changelog MUST record the completed `#`, day-of-week `L`, day-of-month `L`, `W`, and `LW` subsets without closing the larger epic.

### Key Entities

- **Monthly calendar selector**: A recurring monthly intent selecting the last calendar day, the weekday nearest a numbered date, or the final weekday.
- **Intended calendar date**: The numbered monthly date before nearest-weekday adjustment and before or after missing-date resolution as defined by policy.
- **Calendar adjustment**: The deterministic last-day or weekday rule applied to a monthly recurrence.
- **Canonical schedule phrase**: The stable human representation used across explain, task input, preview, editing, and conversion.
- **Named refusal**: A non-mutating result explaining why recognizable syntax or a native recurrence cannot be represented faithfully.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All three selector families convert successfully in both directions through CLI, API, desktop, crontab import, and task export coverage.
- **SC-002**: Automated calendar matrices cover all seven weekdays for `nW`, all seven month-ending weekdays for `LW`, targets 1, 15, 28, 29, 30, and 31, February in leap and non-leap years, and every missing-date policy.
- **SC-003**: Supported schedules produce identical run times before and after round trip for at least twelve consecutive months containing both daylight-saving transitions and a leap-day boundary.
- **SC-004**: Every defined malformed or unsupported boundary returns an actionable field-specific error or refusal and creates or updates zero tasks.
- **SC-005**: Existing cron, scheduling, task-boundary, CLI, API, desktop, documentation, and canonical verification suites remain green.
- **SC-006**: All eight canonical verification gates, whitespace checks, and UTF-8-without-BOM and mojibake audits pass.

## Clarifications

### Session 2026-08-28

- S025 is intentionally one outcome-sized slice covering day-of-month `L`, `nW`, and `LW`; these related selectors will not be split into separate pull requests.
- Nearest-weekday adjustment stays inside the intended month, including the special first-day and final-day weekend cases.
- Missing-date policy resolves an absent numbered target before weekday adjustment; imported cron uses the existing default skip behavior.
- Native `29W` through `31W` export only with effective skip behavior, while `L`, `LW`, and `1W` through `28W` are policy-inert.
- Offset and composite variants such as `L-3`, `L-3W`, lists, ranges, steps, and mixed modifiers remain explicit refusals.

## Assumptions

- The documented selector semantics follow the common Quartz meanings, adapted deliberately to go-schedule's five-field layout and existing Sunday numbering rather than claiming full Quartz compatibility.
- Existing task timezone, DST, catch-up, overlap, and missing-date policies remain authoritative.
- A durable representation change is acceptable within this substantial slice if research proves the current recurrence string cannot carry nearest-weekday semantics honestly.
- No new permissions, external services, notifications, or daemon lifecycle triggers are required.

## Out of Scope

- `L-n`, `L-nW`, multiple or mixed modifier terms, modifier lists/ranges/steps, month-restricted selectors, or combined day-of-month and day-of-week restrictions.
- Quartz `?`, seconds/year fields, `@reboot`, system-crontab user semantics, shell/environment emulation, or arbitrary cron field composition.
- Closing epic issue #22; S025 completes a substantial calendar-selector family but not the full audit.
