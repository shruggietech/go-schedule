# Implementation Plan: Last-Weekday Cron Parity

**Branch**: `codex/024-last-weekday-cron` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Add faithful import and export for one five-field cron day-of-week term in the
form `weekdayL`. Reuse the existing ordinal-weekday representation with `-1`
meaning the last occurrence, translate through the established last-weekday
phrase and RRULE model, and preserve existing task, CLI, API, and crontab
boundaries without adding a recurrence type.

## Technical Context

**Language/Version**: Go 1.25.0, Markdown
**Primary Dependencies**: Go standard library and existing `rrule-go` recurrence model
**Storage**: Existing task records only; no schema or migration change
**Testing**: Go unit/integration tests plus all eight canonical verification gates
**Target Platform**: Windows and Linux local CLI, daemon API, and desktop application
**Project Type**: Local scheduler with shared cron conversion library, CLI, API, and GUI
**Performance Goals**: Constant-time parsing/export with no material scheduling-path overhead
**Constraints**: Five fields only, one last-weekday atom, no approximation, UTF-8 without BOM
**Scale/Scope**: Three production files in `internal/cron`, shared-boundary tests, documentation, and changelog

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One existing ordinal sentinel keeps parsing and rendering responsibilities local without a new abstraction. |
| II. Testing Standards | PASS | Parser, conversion, exporter, task-boundary, CLI, and API behavior is established red before implementation, including DST and refusal cases. |
| III. UX Consistency | PASS | Existing phrases, source retention, stream conventions, and named refusals remain authoritative across interfaces. |
| IV. Performance | PASS | The change performs bounded token parsing and reuses the recurrence engine; no benchmark-sensitive scheduler path changes. |
| V. Autonomous Execution | PASS | S024 follows Spec-Kit, the review branch, analysis, canonical gates, local commit, and the mandatory pre-publication halt. |

### Post-design re-check

All gates remain satisfied. No dependency, persistence model, public API schema,
daemon lifecycle, scheduler dispatch, or governance mechanism is added.

## Architecture and Decision Log

### Reuse the existing ordinal representation

`cron.Field.Ordinal` uses `-1` for a last-weekday term, matching the existing
RRULE representation. Zero remains ordinary semantics and 1 through 5 remain
numbered occurrences. This is a deliberate extension of the S023 field contract
rather than a new cron AST or domain entity.

### Parse one deliberately narrow extension

Only a day-of-week atom ending in `L` reaches the dedicated parser. Numeric and
named weekdays are accepted, Sunday normalizes to zero, and bare `L`, lists,
ranges, steps, multiple terms, mixed `L`/`#`, restricted dates, and restricted
months remain errors or named refusals. Day-of-month `L` remains declined.

### Use the existing last-weekday grammar

Supported input renders as `last <weekday> of the month at HH:MM`. The shared
schedule parser already compiles this phrase to `FREQ=MONTHLY;BYDAY=-1XX`, so
CLI, crontab import, task creation, updates, and preview need no production
changes.

### Export only lossless native rules

A monthly RRULE exports as `weekdayL` only when it has exactly one `-1`
numbered weekday and no competing selector, bound, multi-time, or sub-minute
shape. Missing-date policy is intentionally inert because every month has a
last occurrence of every weekday. Existing selector guards remain unchanged.

### Document a project subset, not universal portability

Cronie documents the chosen five-field layout and Sunday numbering but not
`L`. Quartz documents suffix-`L` semantics under a different field layout and
weekday numbering. Go-schedule therefore documents `weekdayL` as its explicit
five-field extension rather than claiming compatibility with either dialect.

## Project Structure

```text
specs/024-last-weekday-cron/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/last-weekday-cron.md
├── checklists/
└── tasks.md

internal/
├── cron/
├── scheduleinput/
├── cli/
└── api/server/

docs/cron.md
docs/cli.md
CHANGELOG.md
```

**Structure Decision**: Extend the existing cron parser, phrase renderer, and
exporter, then prove propagation through existing shared boundaries. No new
package or dependency is warranted.

## Implementation Phases

1. Add failing parser, phrase, conversion, and export tests for the supported matrix and refusal boundary.
2. Add failing shared task-input, CLI, crontab, and API non-mutation tests.
3. Implement the narrow parsed-field sentinel, parser, phrase, and exporter.
4. Prove round-trip run equivalence across DST and at least three month boundaries.
5. Align fidelity and CLI documentation plus the chronological changelog.
6. Analyze, run focused and canonical verification, audit encoding, commit locally, and halt before publication.

## Complexity Tracking

No constitutional violations or architecture exceptions are required.
