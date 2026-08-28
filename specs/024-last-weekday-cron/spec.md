# Feature Specification: Last-Weekday Cron Parity

**Feature Branch**: `codex/024-last-weekday-cron`
**Created**: 2026-08-28
**Status**: Draft
**Input**: Add faithful two-way conversion for the cron day-of-week `L` extension as the next focused slice of issue #22.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand and import a last-weekday expression (Priority: P1)

As an operator with an existing cron expression such as `0 9 * * 5L`, I can explain, convert, preview, and import it as the last Friday of every month without rewriting the schedule by hand.

**Why this priority**: The scheduler already understands last-weekday schedules in human form, so refusing the equivalent cron form is an unnecessary interoperability gap.

**Independent Test**: Explain or import numeric and named last-weekday expressions and compare their generated run times with the requested weekday across several month boundaries.

**Acceptance Scenarios**:

1. **Given** `0 9 * * 5L`, **When** it is explained or converted, **Then** the result is `last friday of the month at 09:00` with no refusal.
2. **Given** `30 14 * * WEDL`, **When** it is explained or converted, **Then** weekday names are accepted case-insensitively and the result is the last Wednesday of every month at 14:30.
3. **Given** a crontab job using one supported last weekday, **When** it is previewed or imported, **Then** it remains a job, shows the readable phrase, and retains the original cron source at task creation.
4. **Given** any calendar month, **When** future runs are calculated, **Then** exactly its final occurrence of the selected weekday is used.

---

### User Story 2 - Export a last-weekday schedule (Priority: P1)

As an operator with a native last-weekday schedule, I can convert or export it to a canonical five-field cron expression using the day-of-week `L` extension so that the schedule can round-trip without changing its run times.

**Why this priority**: Two-way support makes the conversion promise coherent and covers the remaining monthly weekday form already represented by the product.

**Independent Test**: Export all seven last-weekday schedules, re-import the expressions, and compare generated runs over windows containing DST transitions and multiple month boundaries.

**Acceptance Scenarios**:

1. **Given** `last wednesday of the month at 14:00`, **When** it is converted or exported, **Then** the result is the canonical expression `0 14 * * 3L`.
2. **Given** each weekday from Sunday through Saturday, **When** its last monthly occurrence is exported, **Then** the weekday is preserved exactly.
3. **Given** any missing-date behavior, **When** a last-weekday schedule is exported, **Then** export remains faithful because every month has a last occurrence of every weekday.
4. **Given** a successful round trip, **When** schedules are evaluated through a DST transition and at least three month boundaries, **Then** their run times remain identical.

---

### User Story 3 - Receive precise boundary feedback (Priority: P2)

As an operator, I receive an actionable error or named refusal when `L` is malformed, appears in day-of-month, or is combined with recurrence restrictions outside this focused subset.

**Why this priority**: The product must expand compatibility without silently approximating unsupported cron semantics.

**Independent Test**: Exercise bare `L`, invalid weekdays, multiple terms, mixed extensions, restricted date/month combinations, and day-of-month `L`, verifying rejection without task mutation.

**Acceptance Scenarios**:

1. **Given** bare `L` or an invalid weekday before `L`, **When** the expression is parsed, **Then** it fails with an error naming the day-of-week field and valid weekday forms.
2. **Given** a list, range, step, multiple term, or mixed `L` and `#` expression, **When** conversion is requested, **Then** it is refused or rejected with a reason that says only one last-weekday term is supported.
3. **Given** `L` in day-of-month, **When** conversion is requested, **Then** it remains a named refusal and is not interpreted as last weekday.
4. **Given** a valid last weekday combined with a restricted day-of-month or month, **When** conversion is requested, **Then** the unsupported combination is refused without approximation.

### Edge Cases

- Sunday appears as `0`, `7`, or `SUN`; all forms identify Sunday and export canonically as `0L`.
- Weekday names use the existing case-insensitive three-letter-name rule.
- Every month contains a last occurrence of every weekday, so missing-date behavior is inert.
- A valid day-of-week `L` expression may add a month restriction that the native monthly phrase cannot preserve.
- Bare `L`, multiple terms, lists, ranges, steps, mixed `L` and `#`, day-of-month `L`, `W`, Quartz seconds, and `@reboot` remain outside this slice.
- A failed or refused cron task input must not mutate stored tasks.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product MUST accept exactly one `weekdayL` term in the day-of-week field of a five-field cron expression.
- **FR-002**: Supported weekdays MUST include numeric `0` through `7` and existing case-insensitive weekday names, with both `0` and `7` representing Sunday.
- **FR-003**: Supported import expressions MUST keep day-of-month and month unrestricted so the resulting rule means the last selected weekday of every month.
- **FR-004**: Explain and cron-to-human conversion MUST produce the existing readable last-weekday phrase with the exact time and weekday preserved.
- **FR-005**: Crontab preview and import MUST classify a supported last-weekday line as a job, show its readable phrase, and retain its original cron expression and explicit cron syntax at the task boundary.
- **FR-006**: Native last-weekday schedules for all seven weekdays MUST export to a canonical numeric five-field expression using `weekdayL`.
- **FR-007**: Sunday MUST export canonically as `0L` regardless of whether import used `0L`, `7L`, or a named form.
- **FR-008**: Last-weekday schedules MUST export with any missing-date setting because every month contains the requested occurrence.
- **FR-009**: Cron-to-human-to-cron and human-to-cron-to-human round trips MUST preserve run times across DST transitions and at least three month boundaries.
- **FR-010**: Bare or malformed `L`, invalid weekdays, multiple terms, lists, ranges, steps, mixed `L` and `#`, and unsupported date/month combinations MUST produce a field-specific error or named refusal without approximation.
- **FR-011**: Day-of-month `L` MUST remain a named refusal and MUST NOT be included in this support expansion.
- **FR-012**: Failed or refused last-weekday task creation or update MUST NOT mutate stored tasks.
- **FR-013**: Existing supported cron expressions and named refusals for `W`, six-field input, `@reboot`, lists, steps, day-of-month `L`, and lossy field combinations MUST remain compatible.
- **FR-014**: CLI text and JSON success/refusal stream conventions MUST remain compatible across convert, explain, import, and export workflows.
- **FR-015**: The cron fidelity documentation MUST distinguish supported day-of-week `weekdayL` from declined day-of-month `L`, `W`, and broader arbitrary combinations.
- **FR-016**: The changelog MUST record this focused parity expansion and reference issue #22 without closing the epic.

### Key Entities

- **Last-weekday term**: One weekday whose final occurrence in every calendar month is selected.
- **Canonical cron expression**: A five-field expression whose day-of-week field uses numeric `weekdayL`.
- **Last-weekday schedule**: The existing monthly recurrence that retains the selected weekday, last-occurrence selector, time, timezone, and effective missing-date behavior.
- **Named refusal**: A non-mutating result explaining why otherwise recognizable syntax cannot be represented faithfully.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All seven canonical weekday combinations, plus numeric `7L` and named Sunday aliases, preserve the selected weekday in automated coverage.
- **SC-002**: Every last-weekday schedule exports and re-imports with identical run times over at least one DST transition and three month boundaries.
- **SC-003**: Every defined malformed or unsupported boundary returns an actionable field-specific error or refusal and creates or updates zero tasks.
- **SC-004**: Existing cron, task-boundary, CLI output, documentation, and canonical verification suites remain green.
- **SC-005**: All eight canonical verification gates, whitespace checks, and UTF-8-without-BOM/mojibake audits pass.

## Clarifications

### Session 2026-08-28

- Support is limited to one day-of-week `weekdayL` term; lists, ranges, steps, multiple terms, and mixed `L`/`#` remain outside S024.
- Import requires unrestricted day-of-month and month fields because the existing native phrase represents the selected last weekday of every month.
- Export uses numeric weekday values and canonicalizes Sunday to `0L`.
- Day-of-month `L` remains explicitly declined because it represents a different recurrence family.
- Missing-date behavior does not constrain export because every month has a last occurrence of every weekday.

## Assumptions

- The day-of-week `L` extension is an explicitly supported addition to go-schedule's documented five-field subset, not a claim of universal POSIX cron portability.
- The existing human schedule grammar and recurrence model remain authoritative for generated phrases and run-time calculation.
- No persistence, API schema, timezone ownership, daemon lifecycle, or scheduler dispatch change is required.

## Out of Scope

- Day-of-month `L`, `W`, Quartz `?`, seconds fields, `@reboot`, and arbitrary last-day conversion.
- Arbitrary weekday lists, multiple last-weekday terms, lists/ranges/steps combined with `L`, mixed `L`/`#`, or month-restricted last weekdays.
- New recurrence entities, database migrations, cron daemon execution, shell semantics, environment variables, notifications, or installer work.
- Closing epic issue #22; S024 records partial progress only.
