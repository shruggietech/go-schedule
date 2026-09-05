# Implementation Plan: General Five-Field Cron Breadth

**Branch**: `codex/026-cron-expression-breadth` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Remove the human-phrase bottleneck from standard five-field cron input. Compile validated field sets directly into the existing durable recurrence, retain source solely for editing, preserve concise descriptions for simple shapes, generate exact descriptions for broad shapes, extend policy-aware date resolution and canonical export, and prove fidelity through every task interface and lifecycle stage.

## Technical Context

**Language/Version**: Go 1.25.0 and Markdown **Primary Dependencies**: Existing `rrule-go`, Cobra, SQLite, and Fyne dependencies; no additions **Storage**: Existing schedules schema v6; no migration **Testing**: Go unit/integration tests, deterministic timezone matrices, restart/catch-up coverage, benchmarks, and eight canonical gates **Target Platform**: Windows, Linux, and macOS-supported Go paths **Project Type**: Local daemon, CLI, API, and desktop application **Performance Goals**: Existing p99 dispatch target below 100 ms; broad next-run benchmark within ten percent of baseline **Constraints**: Exact five-field semantics, bounded field sets, inert source text, no DOM/DOW OR approximation, UTF-8 without BOM **Scale/Scope**: Cron parser/compiler/description/export, shared schedule input/evaluation, CLI/API/GUI boundaries, docs, and issue inventory

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One compiler and one durable evaluator replace phrase-mediated execution without adding a parallel engine. |
| II. Testing Standards | PASS | Parser, compilation, lifecycle, policy, timezone, refusal, and regression tests begin red and stay deterministic. |
| III. UX Consistency | PASS | Every interface shares parsing, descriptions, source identity, preview, and refusal behavior. |
| IV. Performance | PASS | Fixed cron bounds, targeted benchmarks, and the existing p99 budget constrain broad recurrence work. |
| V. Autonomous Execution | PASS | S026 follows full Spec-Kit, analysis, implementation, verification, local commit, and one pre-publication halt. |

### Post-design re-check

All gates remain satisfied. No schema, dependency, permission, service, or pinned workflow change is required. The architecture change is proportional: it removes an accidental coupling between two already distinct source grammars while preserving one authoritative recurrence evaluator.

## Architecture and Decision Log

### Compile cron independently of human grammar

`internal/cron` gains a compiler that consumes its existing validated `Spec`, resolves task timezone, and returns a recurring `domain.Schedule`. `internal/scheduleinput` selects cron and invokes this compiler directly. `Expression` remains inert, while RRULE plus adjustment remain authoritative.

This deliberately deviates from the prior rule that cron must become a parseable phrase first. That rule prevented standard field combinations the schedule model can represent and duplicated the description's role as execution input. The new rule is stronger: both source grammars compile into one durable model, and neither display summary is reparsed.

### Represent standard cron as constrained daily field sets

The compiler emits an unbounded daily recurrence with explicit `BYHOUR`, `BYMINUTE`, and `BYSECOND=0`, plus optional `BYMONTH`, positive `BYMONTHDAY`, or ordinary `BYDAY` filters. Parser field values are already ordered and bounded. Field-local steps expand before compilation, so uneven sequences retain boundary resets.

Focused `L`, `W`, and `#` schedules continue through their existing paths because they carry calendar behavior beyond ordinary field sets.

### Preserve concise output, add deterministic broad descriptions

`Phrase` continues to own established concise phrases. A new fallback description renders otherwise supported standard field sets without claiming that the result belongs to the authoring grammar. Explain, convert, crontab preview, API summary, and desktop preview use the same description.

### Resolve date sets under task policy

The missing-date resolver expands from one intended day to an ordered set for recognized daily date-set recurrences. It evaluates each intended date and selected wall time, applies `skip`, `last_valid`, or `next_valid`, deduplicates collisions, normalizes DST only after choosing the calendar date, and requires candidates on or after the anchor and strictly after the cursor.

### Export source-independent canonical fields

Before the existing frequency switch, export recognizes the constrained daily field-set shape. It reconstructs minute, hour, date, month, and weekday fields from recurrence options and emits deterministic numeric wildcards, full-range steps, ranges, or lists. Existing human-generated schedules continue through their established exporters. A non-skip date policy refuses whenever it changes cron's skip behavior.

## Project Structure

```text
specs/026-cron-expression-breadth/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/composite-cron.md
├── checklists/
└── tasks.md

internal/cron/
internal/scheduleinput/
internal/schedule/
internal/catchup/
internal/api/server/
internal/cli/
gui/
docs/
CHANGELOG.md
```

**Structure Decision**: Extend the existing packages and schedule entity in place. No second evaluator, persistence model, interface-specific parser, or new dependency is warranted.

## Implementation Phases

1. Establish red parser, compiler, description, occurrence, policy, export, interface, restart, and benchmark coverage.
2. Add direct standard-field compilation and shared schedule-input routing while preserving focused modifier paths.
3. Add exact broad descriptions and propagate them through conversion and crontab preview.
4. Generalize missing-date date-set evaluation and verify anchor and DST boundaries.
5. Add source-independent canonical composite export and round-trip parity.
6. Prove CLI, API, desktop, persistence, restart, catch-up, and refusal non-mutation behavior.
7. Update product claims, fidelity tables, issue inventory, changelog, and Spec-Kit verification evidence.
8. Run analysis, all canonical gates, encoding audits, commit locally, and halt before publication.

## Complexity Tracking

| Added complexity | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| Direct cron-to-recurrence compiler | Standard cron composition exceeds the human phrase grammar | Expanding natural language or splitting tasks changes semantics and product behavior. |
| Policy-aware multi-date resolver | Task policy must remain meaningful for date lists | Ignoring or forbidding policy would create silent or surprising behavior. |
| Composite export recognizer | Export must prove recurrence fidelity without trusting source | Echoing retained input would violate the inert-source invariant. |
