# Implementation Plan: Dual-Syntax Product Documentation

**Branch**: `codex/021-dual-syntax-docs` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Align current product documentation, help, and the authoritative S001 contract
with the human-or-cron scheduling boundary delivered in S019 and S020. Define
the exact supported five-field cron subset, preserve historical specifications
with supersession notices, and add a fixture-backed documentation policy check.
Also refuse the three calendar-field wildcard-step forms that currently lose
timing information by being silently reduced to daily schedules.

## Technical Context

**Language/Version**: Go 1.25.0, POSIX shell, Markdown

**Dependencies**: Existing standard library and repository scripts only

**Storage**: No schema or persistence change

**Testing**: Focused Go tests, shell policy fixtures, and all eight canonical
verification gates

**Scope**: Current documentation/help surfaces, S001 authority, S008
supersession notes, documentation policy scripts, and one parser safety refusal

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Existing parser remains authoritative; the fix rejects unrepresentable input at its conversion boundary. |
| II. Testing Standards | PASS | Regression tests and policy fixtures precede implementation; canonical gates remain mandatory. |
| III. UX Consistency | PASS | CLI, GUI, README, guides, and contracts describe the same dual-syntax behavior and named refusals. |
| IV. Performance | PASS | The parser check is bounded to three already parsed fields; documentation checks scan a fixed inventory. |
| V. Autonomous Execution | PASS | Full Spec-Kit sequence, review branch, local commit, and mandatory pre-push halt are retained. |

### Explicit scope deviation

The original slice was documentation-only. Research reproduced a correctness
defect where `*/n` in day-of-month, month, or day-of-week is accepted and
silently approximated as daily. The proportional fix is a named refusal for
those forms. No additional cron syntax is introduced; broader fidelity remains
tracked by issue #22.

## Architecture and Decision Log

### One human-first product posture

Current surfaces say that operators may author recurring schedules in readable
phrases or in go-schedule's supported five-field cron subset. Equivalent
weekday-at-09:00 examples make the two paths comparable without implying that
cron knowledge is required.

### Name the exact subset

Documentation separates tokens the parser can read from recurrence shapes the
product can faithfully represent. It covers supported macros, weekday rules,
field-local steps, named refusals, timezone ownership, DST behavior, and the
difference between an expression and a crontab file.

### Preserve chronology and supersede it visibly

S001 is the authoritative product contract and is updated. S008 remains a
historical record, with prominent notices on the pages whose old human-only
boundary is now false.

### Bound documentation policy

A small POSIX helper checks an explicit inventory of current-facing files for
required dual-syntax concepts and targeted obsolete claims. A fixture harness
proves both acceptance and rejection. Historical specs and the changelog are
excluded from the stale-copy scan.

### Refuse lossy calendar steps

After parsing, wildcard steps greater than one in day-of-month, month, and
day-of-week return a fidelity refusal. `*/1` remains equivalent to unrestricted.
Creation boundaries must surface the refusal and must not mutate stored tasks.

## Source References

- [POSIX crontab](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html)
- [Linux crontab(5)](https://man7.org/linux/man-pages/man5/crontab.5.html)
- [robfig/cron expression format](https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format)

## Project Structure

```text
specs/021-dual-syntax-docs/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/documentation-policy.md
├── checklists/
└── tasks.md

README.md                 docs/*.md
internal/cli/             internal/cron/
internal/scheduleinput/   specs/001-task-scheduler/
specs/008-cron-import/    specs/010-docs-site-pages/
scripts/docs-check.sh     scripts/docs-policy-check.sh
test/scripts/docs-policy-check_test.sh
CHANGELOG.md
```

## Implementation Phases

1. Record baselines and add failing parser/boundary regression tests.
2. Implement the narrow calendar-step fidelity refusal.
3. Align current product surfaces, help, and the authoritative S001 contract.
4. Add S008 supersession notices and the fixture-backed documentation policy.
5. Update the changelog, analyze coverage, run canonical verification, audit
   encoding, commit locally, and halt before publication.

## Complexity Tracking

The only departure from the documentation-only boundary is the reproduced
silent timing error described above. It is addressed with a local refusal and
tests instead of expanding the recurrence model or cron dialect.
