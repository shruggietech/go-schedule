# Implementation Plan: Dual-Syntax Task Input Foundation

**Branch**: `codex/019-dual-syntax-input` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/019-dual-syntax-input/spec.md`

## Summary

Add a single recurring schedule-input boundary that classifies an expression as human or cron, honors an optional explicit hint, compiles both forms into the existing RRULE/anchor model, and retains the submitted expression for editing. Adopt it in API preview/create/update, CLI-backed task requests, and cron import; derive non-persisted source identity in responses. No schema or engine change is needed. The GUI continues to force human input until its later adoption slice.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library; existing rrule-go, Cobra, SQLite, and Fyne dependencies only; no new module dependency

**Storage**: Existing SQLite `schedules.expression`, RRULE, and anchor columns; no migration or backfill

**Testing**: Go unit and HTTP integration-style tests, deterministic recurrence windows, `go test -race`, repository coverage gate

**Target Platform**: Linux, macOS, and Windows daemon/CLI; unchanged Fyne GUI

**Project Type**: Local daemon with JSON-over-IPC API, thin CLI, and desktop GUI

**Performance Goals**: Parsing remains bounded to one classification plus one selected syntax parse; no scheduler hot-path or dispatch-latency change

**Constraints**: RRULE/anchor remain authoritative; no parse fallback; no raw cron in the engine; additive wire fields; backward-compatible persisted state; all eight verification gates must pass

**Scale/Scope**: One input package, three server request paths, response enrichment, CLI task help, import preview/create, narrow current-policy docs, and compatibility tests

## Constitution Check

*GATE: Passed before research and re-checked after design.*

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One orchestration package reuses existing parsers and detector; typed syntax and contextual errors prevent duplicated branch logic. |
| II. Testing Standards | PASS | Central, server, import, persistence, GUI-preservation, DST/month, and failure tests precede production changes; canonical race/coverage gates remain mandatory. |
| III. UX Consistency | PASS | API request/response fields are shared by CLI and future GUI work; errors name `schedule` or `schedule_syntax`; CLI stays a thin client. |
| IV. Performance | PASS | The change is outside dispatch hot paths and performs one selected parse per request; no benchmark-sensitive loop changes. |
| V. Autonomous Execution | PASS | Spec-Kit sequence, analysis gate, review branch, local commit, and pre-push halt are preserved. |

Engineering constraints pass: no dependency, no destructive migration, no security-boundary change, no pinned process artifact, and Linux/Windows support remains under the repository verification suite.

### Post-design re-check

The API contract is additive, the storage model is unchanged, current GUI behavior is explicitly pinned to `human`, and all named failure/security paths have test tasks. No constitution exception or complexity deviation is required.

## Architecture and Decision Log

### Central package boundary

Create `internal/scheduleinput`, which imports `internal/cron` and `internal/schedule`. Export `cron.DetectSyntax` so S018 structural detection has one implementation. The new package owns authoring concerns: hint normalization, single-pass selection, cron explanation into the human parser, expression retention, and source-identity derivation.

This deliberately differs from S018's decision not to create a conversion wrapper. S019 introduces a durable cross-syntax authoring boundary consumed by multiple product surfaces. Placing all human input under `internal/cron` would make the cron package own a general task-authoring contract; copying detection into server code would allow drift. One focused orchestration package is the smallest clean dependency graph:

```text
api/server -> scheduleinput -> cron -> schedule -> domain
gui (later) -> scheduleinput
engine -> schedule -> domain
```

The engine never imports `scheduleinput` or reads `Schedule.Expression`.

### Input and hint semantics

- Normalize surrounding input whitespace once.
- Normalize hint case and surrounding whitespace.
- Valid hints are empty, `human`, and `cron`.
- Empty uses `cron.DetectSyntax`; explicit values select exactly one parser.
- Cron uses `cron.Explain`, then `schedule.Parse` in the real task timezone and request anchor. A named refusal becomes a schedule validation error.
- Human uses `schedule.Parse` directly.
- Cron success overwrites only `Schedule.Expression` with the retained cron input. RRULE, anchor, and summary come from the existing schedule parser.

### Source identity without schema state

Add an inert `SourceSyntax` JSON field to `domain.Schedule`, but do not add it to SQLite INSERT/SELECT statements. Server response construction derives it from a non-empty recurring `Expression` through the shared detector. Preview returns the same field. Empty-expression legacy rows and one-offs omit it.

Alternatives rejected:

- A persisted syntax column can disagree with the expression and requires a migration without changing execution capability.
- A top-level task-only field separates identity from the expression it describes and complicates future GUI consumption.
- Reclassification in each client duplicates product logic.

### API validation and versioning

Add optional `schedule_syntax` to preview/create/update request objects. Invalid values produce validation field `schedule_syntax`; selected-parser failures use field `schedule`. A non-empty hint without a recurring schedule is rejected rather than silently ignored. Additive JSON fields remain compatible with existing clients and stored rows.

### Import parity

Keep the scanner's expression and explanatory phrase. Preview upcoming runs and create tasks from `Line.Expr` with explicit cron identity; continue printing `Line.Phrase` for readable audit output. Preserve dry-run, unreachable-daemon, per-line refusal, command parsing, partial-success, and summary behavior.

### GUI scope containment

GUI validation and editing remain human-only. Its preview/create/update requests set `schedule_syntax: human` explicitly so the new automatic API behavior cannot expand GUI acceptance accidentally. Cron prefill/editing is a known capability for the next GUI slice and is not claimed here.

### Current documentation truth

Update only live statements made false by this slice: task CLI help, current package/domain comments, narrow README/CLI/cron/API guidance, and the Unreleased changelog decision. Do not rewrite historical Spec 008 or perform issue #52's broad product-positioning and dialect tutorial.

## Project Structure

### Documentation (this feature)

```text
specs/019-dual-syntax-input/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── schedule-input-api.md
├── checklists/
│   ├── requirements.md
│   └── input-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── cron/
│   ├── convert.go
│   └── convert_test.go
├── scheduleinput/
│   ├── input.go
│   └── input_test.go
├── domain/domain.go
├── api/server/
│   ├── tasks.go
│   ├── update.go
│   ├── tasks_test.go
│   ├── expression_test.go
│   └── update_test.go
└── cli/
    ├── task.go
    ├── task_test.go
    ├── cron.go
    └── cron_test.go
gui/
└── editor.go
docs/
├── api.md
├── cli.md
└── cron.md
README.md
CHANGELOG.md
```

**Structure Decision**: Add one internal orchestration package and modify the existing API, client-facing CLI/import, GUI request construction, and live documentation surfaces. Store SQL and scheduling-engine packages remain unchanged.

## Implementation Phases

### Phase 1 - Baseline and central input

Record the current human-only and S018 behavior, add failing detector/input contracts, export the shared detector, and implement the central parser under test-first discipline.

### Phase 2 - API preview/create/read

Add failing additive-wire, hint, no-fallback, persistence, compatibility, and run-parity tests. Adopt the input parser in preview/create and derive response identity.

### Phase 3 - Update and CLI task authoring

Protect unchanged/replacement/no-mutation behavior, adopt the parser in update, change task help, and pin GUI requests to explicit human syntax.

### Phase 4 - Import retention

Replace the obsolete human-only import assertion with source-retention and preview/create parity tests, then submit the scanner expression plus cron hint.

### Phase 5 - Documentation and verification

Synchronize live comments/docs and changelog, run Spec-Kit analysis, audit UTF-8/no-BOM/mojibake and whitespace, execute all eight gates in the foreground, and commit locally.

## Complexity Tracking

No constitution violation or justified excess complexity.
