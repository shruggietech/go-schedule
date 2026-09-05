# Implementation Plan: Monthly Calendar Cron Parity

**Branch**: `codex/025-monthly-calendar-parity` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Deliver day-of-month `L`, `nW`, and `LW` as one complete interoperability slice. Represent `L` and `LW` with standard RRULE selectors, add one durable `nearest_weekday` calendar adjustment for the only shape RRULE cannot express, and propagate exact parsing, execution, conversion, persistence, preview, editing, export, refusal, and documentation behavior through every interface.

## Technical Context

**Language/Version**: Go 1.25.0, SQL, Markdown
**Primary Dependencies**: Existing `rrule-go`, SQLite, Cobra, and Fyne dependencies
**Storage**: Additive SQLite schema version 6 on `schedules`; backward-compatible JSON API
**Testing**: Go unit/integration tests, migration/restart coverage, and all eight canonical gates
**Target Platform**: Windows and Linux CLI, daemon API, and desktop application
**Project Type**: Local scheduler with shared recurrence, cron, API, CLI, and GUI boundaries
**Performance Goals**: Constant-time parsing; bounded monthly walking; unchanged ordinary-schedule hot path
**Constraints**: Five fields, one selector atom, no approximation, forward migration, UTF-8 without BOM
**Scale/Scope**: Domain, persistence, recurrence, cron, shared interfaces, GUI preview, docs, and issue inventory

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One enum-valued adjustment isolates irreducible `nW` behavior; incompatible stored shapes error. |
| II. Testing Standards | PASS | Parser, calendar matrices, persistence/restart, non-mutation, interfaces, and round trips start red. |
| III. UX Consistency | PASS | Canonical phrases, source identity, refusals, and policy-sensitive preview stay uniform. |
| IV. Performance | PASS | Ordinary RRULE schedules keep their path; adjusted rules use a bounded monthly walk. |
| V. Autonomous Execution | PASS | Full Spec-Kit, review branch, analysis, gates, local commit, and one pre-publication halt. |

### Post-design re-check

All gates remain satisfied. The additive migration is necessary because `nW` is not representable as one RFC 5545 RRULE and source text is deliberately non-authoritative. No package, dependency, service, widget, permission, or governance control is added.

## Architecture and Decision Log

### Use native RRULE selectors where faithful

`L` compiles to monthly `BYMONTHDAY=-1`. `LW` compiles to the last position in the Monday-through-Friday set. Both are complete recurrence definitions and need no custom execution metadata.

### Persist one typed nearest-weekday adjustment

`nW` conditionally depends on each month's weekday and boundary position. `Schedule.CalendarAdjustment` accepts only empty or `nearest_weekday`; its carrier RRULE remains monthly with one positive `BYMONTHDAY=n`. This explicitly changes the authoritative recurrence contract from RRULE alone to RRULE plus an optional typed adjustment. `Expression` remains inert display/edit source.

Schema version 6 adds `schedules.calendar_adjustment TEXT NOT NULL DEFAULT ''`. Existing rows retain identical timing behavior. CRUD, migration, restart, and recovery tests prove adjusted schedules survive storage.

### Apply missing-date policy before weekday adjustment

The adjusted monthly walker resolves or skips an absent numbered date using existing policy rules, then moves the valid result to its nearest weekday without crossing the resolved date's month. Existing wall-time and DST normalization run last. Strictly-after iteration suppresses duplicate instants.

### Keep the cron extension narrow

Day-of-month accepts only `L`, `LW`, or one `1W` through `31W` atom, with wildcard month and day-of-week. Lists, ranges, steps, offsets, mixtures, restricted fields, Quartz `?`, and extra fields are errors or named refusals. Export is canonical numeric output after full fidelity checks.

### Repair policy-sensitive desktop preview

The editor saves missing-date policy but currently omits it from PreviewRequest and does not refresh when it changes. S025 wires both so preview and save agree, with no new GUI control.

### Document the product dialect precisely

Cronie defines the five-field baseline but not these selectors. Quartz documents their calendar meaning in a six-field dialect. Go-schedule documents its deliberate Quartz-derived five-field subset rather than claiming universal portability.

## Project Structure

```text
specs/025-monthly-calendar-parity/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/monthly-calendar-cron.md
├── checklists/
└── tasks.md

internal/{domain,store,schedule,cron,scheduleinput,cli,api/server}/
gui/
docs/
CHANGELOG.md
```

**Structure Decision**: Extend the existing schedule entity, migration chain, monthly resolver, cron AST, and shared boundaries in place. No parallel recurrence engine or interface-specific implementation is warranted.

## Implementation Phases

1. Establish red domain, migration, parser, grammar, and calendar-matrix tests.
2. Add schema version 6 and the typed adjustment with restart coverage.
3. Implement exact `nW` execution and native `L`/`LW` grammar.
4. Implement narrow cron parsing, canonical export, refusals, and year-long run parity.
5. Prove CLI, API, crontab, desktop, source retention, and non-mutation; repair desktop preview policy propagation.
6. Update fidelity, CLI, GUI, issue-inventory, and changelog documentation.
7. Analyze, verify, audit encoding, commit locally, and halt before publication.

## Complexity Tracking

| Added complexity | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| One persisted adjustment enum and schema v6 column | Faithfully execute and round-trip `nW` after restart | Source inference violates the inert-expression contract; fake RRULE tokens are invalid; declining `nW` fails scope. |
