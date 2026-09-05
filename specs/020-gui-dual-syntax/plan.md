# Implementation Plan: GUI Dual-Syntax Scheduling

**Branch**: `codex/020-gui-dual-syntax` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/020-gui-dual-syntax/spec.md`

## Summary

Adopt S019's shared `scheduleinput` boundary in the existing desktop Schedule field. The editor will classify and validate the current text once, then use the returned normalized expression and source syntax for preview and submission. Cron tasks retain their original expression on edit. Human schedules, one-offs, legacy expressionless tasks, and the scheduling engine remain unchanged.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Existing Fyne v2.7.4 GUI and internal `scheduleinput`, `schedule`, API, domain, and timezone packages; no new module dependency

**Storage**: Existing SQLite schedule expression and compiled recurrence data; no schema or persistence change

**Testing**: Deterministic Go GUI/unit tests with recording backend doubles, `go test -race`, and repository coverage gates

**Target Platform**: Existing Windows/Linux/macOS Fyne desktop application

**Project Type**: Local daemon with JSON-over-IPC API, CLI, and desktop GUI

**Performance Goals**: One local classification/parse per validation, preview, or form build; no scheduler hot-path or dispatch-latency change

**Constraints**: One field and no syntax state; no GUI timing evaluator; exact cron round trip; human anchor behavior preserved; all eight verification gates must pass

**Scale/Scope**: One editor implementation file, focused GUI tests, GUI help, one field guide, changelog, and Spec-Kit evidence

## Constitution Check

*GATE: Passed before research and re-checked after design.*

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One editor helper delegates classification and parsing to `scheduleinput`; no detector, evaluator, dependency, or stored UI syntax is added. |
| II. Testing Standards | PASS | GUI request, prefill, switching, invalid/refused cron, legacy, one-off, and human regression tests precede implementation; race and coverage gates remain mandatory. |
| III. UX Consistency | PASS | The single field uses the same accepted syntax and named refusals as API/CLI; retained cron is not translated behind the operator's back. |
| IV. Performance | PASS | Work occurs only during editor interaction, outside dispatch hot paths, and is bounded to one existing parser call. |
| V. Autonomous Execution | PASS | The complete Spec-Kit sequence, analysis gate, review branch, local commit, and mandatory pre-push halt are preserved. |

Engineering constraints pass: no dependency, migration, security-boundary, daemon, execution, or scheduling-engine change. Linux and Windows remain covered by the canonical repository verification.

### Post-design re-check

The design consumes the additive S019 request/response contract without changing it, derives syntax from current field text rather than stale response metadata, and contains documentation to GUI-local surfaces. No constitution exception or complexity deviation is required.

## Architecture and Decision Log

### One field, one shared classifier

Add a focused editor helper that calls `scheduleinput.Parse` with an empty hint, the effective Schedule text, the selected timezone, and the current anchor. Its `Expression`, `Syntax`, and compiled schedule become the single source for local validity, preview requests, and recurring submissions.

A syntax selector or cached syntax property was rejected because either can disagree with edited text. A GUI-specific cron detector was rejected because it would duplicate the no-fallback contract delivered in S019.

### Current text controls current intent

`domain.Schedule.Expression` remains the edit prefill. `SourceSyntax` is useful response evidence but is not editor state. Each validation, preview, and save reclassifies the current text through the shared boundary, so cron-to-human and human-to-cron edits cannot inherit stale metadata.

### Submission parity

Extend the private `taskForm` with the selected recurring syntax. `buildForm` stores the normalized expression and syntax returned by the same shared parser. `submitTask` sends both for create and schedule-replacing update. One-offs and blank expressionless legacy updates leave both recurring fields empty.

### Human-only Start-at affordance

Keep `schedule.IsSubDailyInterval` solely for the existing human `Start at` field and anchor composition. Cron text does not expose or receive this GUI-only affordance. Timing and previews still come from the daemon's central boundary.

### Error and documentation containment

Strip internal parser/conversion prefixes from displayed validation errors while retaining the named fidelity refusal. Update in-editor help, the Schedule field hint, `docs/gui-fields.md`, the one newly false GUI sentence in `docs/cron.md`, and the Unreleased changelog. Issue #52 retains the broad product documentation rewrite.

## Project Structure

### Documentation (this feature)

```text
specs/020-gui-dual-syntax/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── gui-schedule-editor.md
├── checklists/
│   ├── requirements.md
│   └── gui-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
gui/
├── editor.go
├── editor_test.go
├── editor_prefill_test.go
└── app_test.go
docs/
├── gui-fields.md
└── cron.md
CHANGELOG.md
```

**Structure Decision**: Modify the existing editor boundary and recording GUI tests. Reuse `internal/scheduleinput` without changing central parsing, API, domain, persistence, or engine packages.

## Implementation Phases

### Phase 1 - Baseline and editor contract tests

Record the current focused suite and add failing tests for cron classification, preview/submission identity, refusals, exact prefill, syntax switching, and compatibility paths.

### Phase 2 - Dual-syntax editor boundary

Replace human-only local parsing and hard-coded request hints with one shared input helper. Carry its normalized expression and syntax through preview and the private submission form.

### Phase 3 - Edit compatibility and guidance

Protect retained cron editing, human/cron switching, one-off isolation, legacy preservation, and human anchors. Update GUI-local hints, help, documentation, and the chronological changelog.

### Phase 4 - Analysis, verification, and local commit

Re-run Spec-Kit analysis, audit encoding and whitespace, execute the eight canonical gates in the foreground, commit locally, and halt before publication.

## Complexity Tracking

No constitution violation or justified excess complexity.
