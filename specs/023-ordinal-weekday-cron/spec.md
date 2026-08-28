# Feature Specification: Ordinal-Weekday Cron Parity

**Feature Branch**: `codex/023-ordinal-weekday-cron`
**Created**: 2026-08-28
**Status**: Draft
**Input**: Add faithful two-way conversion for the cron `#` ordinal-weekday extension as the next focused slice of issue #22.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand and import an ordinal-weekday expression (Priority: P1)

As an operator with an existing cron expression such as `0 9 * * 5#3`, I can explain, convert, preview, and import it as the third Friday of every month without rewriting the schedule by hand.

**Why this priority**: The scheduler already understands ordinal-weekday schedules in human form, so refusing the equivalent cron form is an unnecessary interoperability gap.

**Independent Test**: Explain or import numeric and named ordinal-weekday expressions and compare their generated run times with the requested weekday and ordinal across month boundaries.

**Acceptance Scenarios**:

1. **Given** `0 9 * * 5#3`, **When** it is explained or converted, **Then** the result is `3rd friday monthly at 09:00` with no refusal.
2. **Given** `30 14 * * WED#2`, **When** it is explained or converted, **Then** weekday names are accepted case-insensitively and the result is the second Wednesday monthly at 14:30.
3. **Given** a crontab job using one supported ordinal weekday, **When** it is previewed or imported, **Then** it remains a job, shows the readable phrase, and retains the original cron source at task creation.
4. **Given** a fifth-weekday expression in a month without that fifth weekday, **When** future runs are calculated with cron-compatible skip behavior, **Then** that month produces no run and later matching months remain correct.

---

### User Story 2 - Export an ordinal-weekday schedule (Priority: P1)

As an operator with a native ordinal-weekday schedule, I can convert or export it to a canonical five-field cron expression using the `#` extension so that the schedule can round-trip without changing its run times.

**Why this priority**: Two-way support makes the conversion promise coherent and removes a current refusal for a recurrence the product already models.

**Independent Test**: Export first-through-fifth weekday schedules, re-import the expressions, and compare every generated run over windows containing DST transitions, month boundaries, and a missing fifth occurrence.

**Acceptance Scenarios**:

1. **Given** `3rd wednesday monthly at 14:00`, **When** it is converted or exported, **Then** the result is the canonical expression `0 14 * * 3#3`.
2. **Given** first-through-fifth ordinal weekdays, **When** each is exported, **Then** the ordinal and weekday are preserved exactly.
3. **Given** a fifth-weekday schedule with skip behavior, **When** it round-trips through cron, **Then** months without that occurrence remain skipped.
4. **Given** a fifth-weekday schedule with an effective missing-date behavior cron cannot represent, **When** export is requested, **Then** it is refused by name rather than approximated.

---

### User Story 3 - Receive precise boundary feedback (Priority: P2)

As an operator, I receive an actionable error or named refusal when a `#` expression is malformed or combines the extension with recurrence restrictions outside this focused subset.

**Why this priority**: Expanding syntax must retain the project's no-silent-approximation guarantee.

**Independent Test**: Exercise invalid ordinals, multiple ordinal weekdays, non-weekday placement, restricted month/date combinations, and surrounding unsupported syntax and verify each is rejected without task mutation.

**Acceptance Scenarios**:

1. **Given** ordinal zero, an ordinal above five, or a nonnumeric ordinal, **When** the expression is parsed, **Then** it fails with an error naming the day-of-week field and valid ordinal range.
2. **Given** a list, range, or step combined with `#`, **When** conversion is requested, **Then** it is refused or rejected with a reason that says only one weekday and one ordinal are supported.
3. **Given** `#` in any field other than day-of-week, **When** parsing occurs, **Then** it is rejected by name rather than interpreted as another value.
4. **Given** a valid ordinal weekday combined with a restricted day-of-month or month, **When** conversion is requested, **Then** the unsupported combination is refused without approximation.

### Edge Cases

- Sunday appears as `0`, `7`, or `SUN`; all forms identify Sunday and export canonically as `0`.
- Weekday names use the existing case-insensitive three-letter-name rule.
- The fifth occurrence of a weekday does not exist in every month.
- A cron expression may use valid `#` syntax but add a month restriction that the native monthly phrase cannot preserve.
- Multiple `#` terms, lists, ranges, steps, `L`, `W`, Quartz seconds, and `@reboot` remain outside this slice.
- A failed or refused cron task input must not mutate stored tasks.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product MUST accept exactly one `weekday#ordinal` term in the day-of-week field of a five-field cron expression.
- **FR-002**: Supported weekdays MUST include numeric `0` through `7` and existing case-insensitive weekday names, with both `0` and `7` representing Sunday.
- **FR-003**: Supported ordinals MUST be integers from `1` through `5` inclusive.
- **FR-004**: Supported import expressions MUST keep day-of-month and month unrestricted so the resulting rule means the selected ordinal weekday of every month.
- **FR-005**: Explain and cron-to-human conversion MUST produce the existing readable ordinal-weekday phrase with the exact time, ordinal, and weekday preserved.
- **FR-006**: Crontab preview and import MUST classify a supported ordinal-weekday line as a job, show its readable phrase, and retain its original cron expression and explicit cron syntax at the task boundary.
- **FR-007**: Native first-through-fifth ordinal-weekday schedules MUST export to a canonical numeric five-field expression using `weekday#ordinal`.
- **FR-008**: Sunday MUST export canonically as numeric `0` regardless of whether import used `0`, `7`, or a name.
- **FR-009**: A fifth-weekday schedule MUST export only when its effective missing-date behavior matches cron's skip behavior; incompatible behavior MUST receive a named refusal.
- **FR-010**: First-through-fourth ordinal weekdays MAY export with any missing-date setting because those occurrences exist in every month and the setting does not alter run times.
- **FR-011**: Cron-to-human-to-cron and human-to-cron-to-human round trips MUST preserve run times across month boundaries, DST transitions, and absent fifth occurrences.
- **FR-012**: Invalid ordinals, malformed `#` tokens, multiple terms, lists, ranges, steps, non-day-of-week placement, and unsupported date/month combinations MUST produce a field-specific error or named refusal without approximation.
- **FR-013**: Failed or refused ordinal-weekday task creation or update MUST NOT mutate stored tasks.
- **FR-014**: Existing supported cron expressions and existing named refusals for `L`, `W`, six-field input, `@reboot`, lists, and lossy field combinations MUST remain compatible.
- **FR-015**: CLI text and JSON success/refusal stream conventions MUST remain compatible across convert, explain, import, and export workflows.
- **FR-016**: The cron fidelity documentation MUST distinguish the newly supported `#` subset from still-declined `L`, `W`, and broader arbitrary combinations.
- **FR-017**: The changelog MUST record this focused parity expansion and reference issue #22 without closing the epic.

### Key Entities

- **Ordinal-weekday term**: One weekday plus an occurrence number from first through fifth within each month.
- **Canonical cron expression**: A five-field expression whose day-of-week field uses numeric `weekday#ordinal`.
- **Ordinal-weekday schedule**: The existing monthly recurrence that retains the selected weekday, ordinal, time, timezone, and effective missing-date behavior.
- **Named refusal**: A non-mutating result explaining why otherwise recognizable syntax cannot be represented faithfully.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 35 canonical weekday-and-ordinal combinations (seven weekdays times five ordinals), plus numeric `7` and named Sunday aliases, preserve the selected weekday and ordinal in automated coverage.
- **SC-002**: Every first-through-fifth ordinal-weekday schedule exports and re-imports with identical run times over at least one DST transition and three month boundaries.
- **SC-003**: At least one fifth-weekday round trip proves that a month without the requested occurrence remains skipped.
- **SC-004**: Every defined malformed or unsupported boundary returns an actionable field-specific error or refusal and creates or updates zero tasks.
- **SC-005**: Existing cron, task-boundary, CLI output, documentation, and canonical verification suites remain green.
- **SC-006**: All eight canonical verification gates, whitespace checks, and UTF-8-without-BOM/mojibake audits pass.

## Clarifications

### Session 2026-08-28

- Support is limited to one day-of-week `weekday#ordinal` term with ordinal `1` through `5`; lists, ranges, steps, and multiple terms remain outside S023.
- Import requires unrestricted day-of-month and month fields because the existing native phrase represents an ordinal weekday of every month.
- Export uses numeric weekday values and canonicalizes Sunday to `0`.
- `last weekday` remains outside S023 because cron's `#` extension expresses numbered occurrences, not the separate `L` extension.
- Fifth-weekday export requires cron-compatible skip behavior; first through fourth exist every month and are unaffected by missing-date policy.

## Assumptions

- The `#` extension is an explicitly supported addition to go-schedule's documented five-field subset, not a claim of universal POSIX cron portability.
- The existing human schedule grammar and recurrence model remain authoritative for generated phrases and run-time calculation.
- No persistence, API schema, timezone ownership, daemon lifecycle, or scheduler dispatch change is required.

## Out of Scope

- `L`, `W`, Quartz `?`, seconds fields, `@reboot`, and last-weekday conversion.
- Arbitrary weekday lists, multiple ordinal terms, lists/ranges/steps combined with `#`, or month-restricted ordinal weekdays.
- New recurrence entities, database migrations, cron daemon execution, shell semantics, environment variables, notifications, or installer work.
- Closing epic issue #22; S023 records partial progress only.
