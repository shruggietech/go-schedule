# Implementation Plan: Cron Parity Closure

**Branch**: `codex/028-cron-parity-closure` | **Date**: 2026-08-29 | **Spec**: [spec.md](spec.md)

## Summary

Close issue #22 with one substantial fidelity slice. Make crontab import preserve ordered timezone, environment, shell, user, and stdin behavior; map a bounded six-field Quartz subset into the existing RRULE scheduler; add durable task stdin; explicitly ratify incompatible audit rows; and scope public claims to the tested contract.

## Technical Context

**Language/Version**: Go 1.25.0 and Markdown
**Primary Dependencies**: Existing Cobra, `rrule-go`, SQLite, and Fyne dependencies; no additions
**Storage**: SQLite schema v8 with one additive task stdin column
**Testing**: Go unit/integration tests, migration/restart/executor fixtures, CLI and GUI headless tests, recurrence benchmarks, and eight canonical gates
**Target Platform**: Windows, Linux, and macOS-supported Go paths; imported cron execution follows Unix shell semantics
**Project Type**: Local daemon, CLI, API, and desktop application
**Performance Goals**: Existing p99 dispatch target below 100 ms; representative next-run benchmarks within ten percent of baseline
**Constraints**: No silent approximation, one authoritative recurrence evaluator, deterministic import layout, additive migration, UTF-8 without BOM
**Scale/Scope**: Cron parser/compiler/export, schedule missing-date seconds, crontab scanner, task model/store/API/executor, CLI import/reporting, lifecycle surfaces, docs, and issue #22

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Stateful import context and a dialect-aware field model replace lossy token splitting without adding a second scheduler. |
| II. Testing Standards | PASS | Parser, importer, executor, migration, lifecycle, compatibility, and performance tests begin red before implementation. |
| III. UX Consistency | PASS | Explicit layout flags prevent guessing; every unsupported row and lossy export has a visible reason. |
| IV. Performance | PASS | Seconds compile into existing RRULE sets; baseline and post-change benchmarks enforce the ten-percent threshold. |
| V. Autonomous Execution | PASS | S028 follows full Spec-Kit, blocking analysis, test-first implementation, all gates, local commit, and one pre-publication halt. |

### Post-design re-check

All gates remain satisfied. Schema v8 is forward-only and additive. No dependency, service, permission, pinned artifact, workflow, or release change is required.

## Architecture and Decision Log

### Replace lossy command splitting

This slice deliberately changes the prior documented importer behavior. Imported cron commands now run through the effective shell with the original command string because whitespace splitting corrupts valid cron. The change is limited to crontab import; ordinary task argv authoring remains unchanged.

### Keep schedule and process timezone separate

`CRON_TZ` becomes per-line schedule context. `TZ` remains an environment variable visible to the child. An explicit `--timezone` remains an operator override. This follows cron semantics and avoids silently moving runs.

### Persist stdin as execution data

Task gains one optional string and schema v8 supplies an empty compatibility default. PATCH uses pointer semantics to support clearing. Executor attaches an exact reader. No temporary file or shell rewrite is introduced.

### Require explicit file dialect and system layout

Import accepts `--dialect unix|quartz` and `--system`. These choices determine exact field counts. Auto-detection is rejected because a valid command can begin with a numeric or cron-shaped token.

### Normalize two weekday dialects into one recurrence

Five-field input keeps Unix weekday numbering and synthetic second zero. Six-field input uses Quartz weekday numbering and supports `?` only in day fields. Both normalize into the existing internal weekday set and compile to one RRULE with `BYSECOND`.

### Preserve canonical output and reject operational loss

Second-zero recurrences keep five-field output; seconds-bearing recurrences use six fields. Export refuses tasks whose environment, run-as identity, or stdin cannot fit faithfully into a standalone line. This is an intentional tightening of an existing lossy boundary.

### Close the audit through explicit disposition

The documentation matrix records one decision for A1-A12 and B1-B9. Event, notification, alternate-file, composite-OR, and unsupported modifier work stays outside this slice with rationale. README copy describes a faithful documented subset instead of universal cron parity.

## Project Structure

```text
specs/028-cron-parity-closure/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/cron-parity.md
├── checklists/
└── tasks.md

internal/domain/
internal/cron/
internal/schedule/
internal/store/
internal/api/server/
internal/executor/
internal/cli/
internal/scheduleinput/
internal/catchup/
internal/engine/
gui/
docs/
CHANGELOG.md
README.md
```

**Structure Decision**: Extend the existing single Go module and established package boundaries. No new service, dependency, or parallel evaluator is introduced.

## Complexity Tracking

No constitution violations require justification.
