# Implementation Plan: Ordinal-Weekday Cron Parity

**Branch**: `codex/023-ordinal-weekday-cron` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Add faithful import and export for one five-field cron day-of-week term in the form `weekday#ordinal`, where ordinal is 1 through 5. Represent the occurrence number on the existing parsed cron field, translate through the existing monthly ordinal-weekday phrase and RRULE model, and preserve the established task, CLI, API, and crontab boundaries without adding a recurrence type.

## Technical Context

**Language/Version**: Go 1.25.0, Markdown **Primary Dependencies**: Go standard library and existing `rrule-go` recurrence model **Storage**: Existing task records only; no schema or migration change **Testing**: Go unit/integration tests plus all eight canonical verification gates **Target Platform**: Windows and Linux local CLI, daemon API, and desktop application **Project Type**: Local scheduler with shared cron conversion library, CLI, API, and GUI **Performance Goals**: Constant-time parsing/export with no additional I/O or allocations material to scheduling **Constraints**: Five fields only, one ordinal weekday, no approximation, UTF-8 without BOM **Scale/Scope**: Three production files in `internal/cron`, shared-boundary tests, documentation, and changelog

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One optional ordinal value extends the existing parsed field and keeps parsing and rendering responsibilities local. |
| II. Testing Standards | PASS | Parser, conversion, exporter, task-boundary, CLI, and API behavior is established red before implementation, including DST and missing-fifth cases. |
| III. UX Consistency | PASS | Existing phrases, source retention, stream conventions, and named refusals remain authoritative across interfaces. |
| IV. Performance | PASS | The change performs bounded token parsing and reuses the existing recurrence engine; no benchmark-sensitive scheduler path changes. |
| V. Autonomous Execution | PASS | S023 follows Spec-Kit, the review branch, analysis, canonical gates, local commit, and the mandatory pre-publication halt. |

### Post-design re-check

All gates remain satisfied. No dependency, persistence model, public API schema, daemon lifecycle, scheduler dispatch, or governance mechanism is added.

## Architecture and Decision Log

### Model the ordinal on the existing cron field

`cron.Field` gains `Ordinal`, where zero means ordinary field semantics and 1 through 5 applies only to a single day-of-week value. This keeps Sunday normalization and `Single()` intact and avoids a second recurrence entity.

### Parse one deliberately narrow extension

The general extension refusal remains active for every non-day-of-week `#` and for `L` and `W`. Only a day-of-week token containing `#` reaches a dedicated parser. One atom and one ordinal are accepted; malformed atoms are field errors, while recognizable lists, ranges, steps, or multiple terms receive the named single-term refusal. Day-of-month and month must remain unrestricted.

### Use the existing phrase and recurrence grammar

Supported input renders as the existing phrase, such as `3rd friday monthly at 09:00`. The shared schedule parser already compiles that phrase to the correct monthly RRULE, so CLI, crontab import, task creation, updates, and preview need no production changes.

### Export only lossless native rules

A monthly RRULE exports only when it has exactly one numbered weekday with an ordinal from 1 through 5 and no competing date selector. Weekdays use the project's current Vixie-style numbering and Sunday canonicalizes to zero. The existing missing-date policy gate permits fifth occurrences only with effective skip behavior; first through fourth are always present and remain policy-inert.

### Document a project subset, not universal cron portability

Cronie documents ordinary five-field weekday numbering but not `#`; Quartz documents `#` under a different field layout and weekday numbering. Therefore go-schedule documents this as its explicit five-field subset extension rather than claiming POSIX or Quartz compatibility. This deliberately corrects the potentially misleading assumption that all cron implementations share this dialect.

## Project Structure

```text
specs/023-ordinal-weekday-cron/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/ordinal-weekday-cron.md
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

**Structure Decision**: Extend the existing cron parser, phrase renderer, and exporter, then prove propagation through existing shared boundaries. No new package or dependency is warranted.

## Implementation Phases

1. Add failing parser, phrase, conversion, and export tests for the supported matrix and refusal boundary.
2. Add failing shared task-input, CLI, crontab, and API non-mutation tests.
3. Implement the narrow parsed-field representation, parser, phrase, and exporter.
4. Prove round-trip run equivalence across DST, month boundaries, and an absent fifth occurrence.
5. Align fidelity and CLI documentation plus the chronological changelog.
6. Analyze, run focused and canonical verification, audit encoding, commit locally, and halt before publication.

## Complexity Tracking

No constitutional violations or architecture exceptions are required.
