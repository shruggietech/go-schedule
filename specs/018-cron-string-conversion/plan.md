# Implementation Plan: Pure Schedule String Conversion

**Branch**: `codex/018-cron-string-conversion` | **Date**: 2026-08-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/018-cron-string-conversion/spec.md`

## Summary

Add `gosched cron convert` as a symmetric, daemon-free one-string conversion surface. A pure conversion service in the existing `internal/cron` boundary will classify or explicitly select the input syntax, reuse the current cron-to-phrase path, parse human schedules through `internal/schedule`, and extract a schedule-only cron renderer from the existing task export path. Conversions that depend on an implicit creation anchor or otherwise lose timing meaning are refused. Text output stays one-line and pipeline-friendly; structured failures are written to stderr without a second plain-text error.

## Technical Context

**Language/Version**: Go 1.25.0 (`go.mod`, `GOTOOLCHAIN=auto`)

**Primary Dependencies**: Existing `github.com/spf13/cobra`, `github.com/teambition/rrule-go`, and standard library only. No dependency is added.

**Storage**: N/A. Conversion is an in-memory pure operation and changes no stored task, configuration, or file.

**Testing**: Table-driven package tests in `internal/cron`; CLI rendering, stream, JSON, and exit-classification tests in `internal/cli`; canonical foreground verification through `scripts/verify.sh all`.

**Target Platform**: Linux, macOS, and Windows CLI; normal POSIX and PowerShell argument quoting.

**Project Type**: Single Go module containing daemon, CLI, and desktop GUI.

**Performance Goals**: One conversion completes through bounded parsing and formatting proportional to the input and recurrence option size. It is not a scheduling hot path, so the existing dispatch benchmark is unaffected and no new performance benchmark is warranted.

**Constraints**: Fully local; exact one-string output; exit 2 for malformed and unfaithful conversion; no task-authoring, API, GUI, store, configuration, network, IPC, or cron-dialect expansion.

**Scale/Scope**: One new conversion module file plus tests in `internal/cron`, one CLI subcommand and tests, two focused documentation sections, changelog, and Spec-Kit artifacts.

## Constitution Check

*GATE: passed before Phase 0 and re-checked after Phase 1 design.*

| Principle | Assessment |
| --- | --- |
| **I. Code Quality** | The conversion stays in the existing cron boundary, exposes explicit typed syntax/result values, uses named refusals instead of panics, and extracts rather than duplicates the schedule-to-cron renderer. Exported internal-package symbols receive contract comments. |
| **II. Testing (NON-NEGOTIABLE)** | Tests precede implementation and cover success, malformed input, named refusals, exact streams, structured output, auto/forced classification, and recurrence parity across DST and month boundaries. No clock, migration, recovery, concurrency, or IPC safety surface is weakened. |
| **III. UX Consistency** | The command remains verb-noun (`cron convert`), honors global `--json`, places results on stdout and diagnostics on stderr, uses exit 2 for validation/fidelity refusal, and gives copyable examples. |
| **IV. Performance** | Work is bounded local parsing/formatting outside the dispatch path. No unbounded search, goroutine, network wait, or allocation-heavy recurrence enumeration is introduced into production conversion. |
| **V. Autonomous execution** | S018 follows specify through analyze, implementation, all eight foreground gates, local commit, and the mandatory pre-publication halt on a review branch. |

No constitutional violation or complexity exception is present.

## Project Structure

### Documentation (this feature)

```text
specs/018-cron-string-conversion/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-convert.md
├── checklists/
│   ├── requirements.md
│   └── conversion.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/
├── cron/
│   ├── convert.go          # syntax classification and pure bidirectional conversion
│   ├── convert_test.go     # exact output, refusal, and calendar-parity corpus
│   ├── export.go           # extract schedule-only rendering from task export
│   └── export_test.go      # preserve task-export behavior
└── cli/
    ├── cron.go             # add cron convert command and rendering
    ├── cron_test.go        # conversion command, stream, JSON, and no-daemon tests
    ├── cli.go              # suppress duplicate text after reported JSON diagnostics
    └── cli_test.go         # exit classification where needed

docs/
├── cli.md                  # convert contract and copyable examples
└── cron.md                 # distinguish convert, explain, import, and export

CHANGELOG.md                # Unreleased Added entry; no dated architecture decision
```

**Structure Decision**: Keep conversion in `internal/cron`, which already owns both cron parsing/phrasing and schedule export. Creating a new central task-input package before #50 defines source retention would overreach this slice. The extracted schedule-only renderer removes task-state policy from pure string conversion without duplicating RRULE mapping.

## Implementation Approach

### Phase A - Pure conversion contract

- Add `Syntax` values for cron and human plus a `Conversion` result with stable identity, input, output, and refusal fields.
- Auto-classify trimmed `@` input and five fields with a cron-shaped minute field as cron; classify all other input as human. An explicit destination forces the opposite input.
- Cron input delegates to the existing `Explain` path. Errors and unsupported outcomes become named conversion refusals, never human fallbacks.
- Correct the already-supported minute-step, every-minute, and hourly phrases to state `starting at 00:00`; omitting that phase would make imported or round-tripped schedules depend on task-creation time and change cron timing.
- Human input delegates to `schedule.Parse` using a fixed UTC planning anchor, then to a schedule-only cron renderer. The anchor exists only to satisfy the parser; implicit-anchor schedules are identified and refused rather than having the chosen instant leak into output.

### Phase B - Faithful human-to-cron rendering

- Extract `ExportSchedule(schedule, missingDatePolicy)` from `Export(task, schedule)`. The task wrapper retains enabled/state checks and delegates the recurrence mapping, preserving every existing export result.
- Require an explicit time/phase wherever the human phrase would otherwise inherit conversion time. Sub-daily intervals without `starting at`/`from` and calendar rules without a complete cron calendar position are refused.
- Keep the existing five-field dialect and named loss boundaries unchanged.

### Phase C - CLI and diagnostics

- Register `convert` beside explain/import/export with one positional string; implement and validate optional `--to cron|human` alongside the text-mode user stories rather than deferring it to structured automation.
- Default success writes exactly the output plus newline. Default refusal returns the established usage/validation class for exit 2.
- Structured success writes the five-field object to stdout. Structured refusal writes the same object to stderr and returns a reported usage error that the root executor maps to exit 2 without printing a duplicate plain diagnostic.
- Keep command execution injectable through Cobra output/error writers so exact streams are tested without replacing process globals.

### Phase D - Documentation and traceability

- Document both directions, forced destination, quoting, JSON, exit behavior, and refusal examples in `docs/cli.md` and `docs/cron.md`.
- Preserve the boundary-only task-authoring posture: conversion is allowed, but `task add/edit --schedule` remains human-only until #50.
- Add an Unreleased changelog entry and close only #51. The future PR references
  #50 rather than closing it.

## Verification

1. Capture the existing cron/CLI focused baseline.
2. Record red tests before each conversion and CLI implementation phase.
3. Run focused package tests for both directions, parity, streams, JSON, and exit classification.
4. Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits.
5. Run `sh scripts/verify.sh all` in the foreground through format, vet, lint, race, GUI, coverage, docs, and automation.

## Complexity Tracking

No entries. The design adds no package, dependency, process, storage model, network path, or alternative scheduler.
