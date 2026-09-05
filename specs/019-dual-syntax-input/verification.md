# Verification Evidence: Dual-Syntax Task Input Foundation

**Date**: 2026-08-28

**Branch**: `codex/019-dual-syntax-input`

## Baseline

- `go test ./internal/cron ./internal/schedule ./internal/api/server ./internal/cli ./gui -count=1` passed before test or production changes.
- S018 owns the only structural cron detector and pure cron-to-human refusal contract; the detector is private to `internal/cron`.
- API preview/create/update route recurring input directly through the human-only `schedule.Parse` function.
- `Schedule.Expression` is already stored by SQLite migration v4 and CRUD, and is inert with respect to engine execution.
- Cron import reports both `Line.Expr` and `Line.Phrase` but previews/creates from the phrase, losing cron source identity.
- CLI task add/edit pass `--schedule` to the API without a local parser, while current help describes the value as human-only.
- GUI validation is locally human-only and its preview/create/update requests currently omit a syntax hint.
- No schema migration, new dependency, IPC/security change, or engine timing path is required by the planned design.

## Test-first evidence

### Foundational red phase

- Added central automatic/forced selection, hint validation, retained source, named refusal, human compatibility, source identity, and DST/month parity regressions before production code.
- Added an exported-detector parity regression to keep S018 conversion and task input on one classification rule.
- `go test ./internal/cron ./internal/scheduleinput -count=1` failed as expected because `DetectSyntax`, the `scheduleinput` types, and parser did not exist.

### Create/preview/read red phase

- Added cron preview/create parity, explicit and automatic syntax, invalid-hint routing, forced-parser no-fallback, retained storage source, one-off identity, and expressionless legacy regressions before changing the API contract.
- The focused server test build failed as expected because `schedule_syntax` and transient `source_syntax` did not yet exist.

### Update/CLI/GUI red phase

- Added replacement, preservation, invalid-input no-mutation, policy independence, dual-syntax help, and explicit GUI human-hint regressions before changing those callers.
- The focused build/tests failed as expected because update requests lacked `schedule_syntax`, CLI help remained human-only, and GUI requests omitted the containment hint.

### Import red phase

- Replaced the obsolete phrase-substitution policy test with retained cron source/hint, preview/create parity, readable phrase output, and mixed partial-success regressions.
- The focused import test build failed as expected because the injected preview boundary still accepted only a generated phrase, timezone, and count.

## Focused verification

- `go test ./internal/cron ./internal/scheduleinput -count=1` passed after the shared detector and central schedule-input boundary were implemented.
- The boundary retains trimmed cron source while compiling only the resulting human phrase through `schedule.Parse`; raw cron does not reach the engine.
- `go test ./internal/api/server ./internal/cron ./internal/scheduleinput -count=1` passed with automatic and explicit cron preview/create, matching human RRULE and run times, persisted source reload, and derived response identity. One-offs and expressionless legacy rows expose no source identity.
- `go test ./internal/api/server ./internal/cli ./gui ./internal/cron ./internal/scheduleinput -count=1` passed with human-to-cron and cron-to-human replacement, omitted-schedule preservation, invalid-update no-mutation, independent missing-date policy changes, dual-syntax CLI help, and explicit human hints on the deliberately human-only GUI path.
- `go test ./internal/cli ./internal/cron ./internal/api/server ./internal/scheduleinput -count=1` passed with import preview and create using the same normalized `Line.Expr` and explicit cron identity. Phrase reporting, dry-run behavior, unreachable-daemon reporting, refusal continuation, command parsing, and mixed partial-success accounting remain covered.

## Repository verification

- `sh scripts/verify.sh all` passed in the foreground:
  1. format
  2. vet
  3. lint (`0 issues`)
  4. race, including integration tests
  5. GUI and view-model tests
  6. coverage (`engine 88.1%`, `schedule 91.5%`, `timezone 88.9%`, `store 87.0%`, `catchup 87.5%`, `logbus 91.1%`)
  7. documentation (`11 pages`, links, front matter, fences, and theme clean)
  8. automation (approved actions, CodeQL contract, and eight-gate manifest)

## Spec-Kit readiness and quickstart

- Specification quality checklist: 16/16 complete.
- Dual-syntax input contract checklist: 30/30 complete.
- Focused quickstart behavior is covered by the green boundary/API/CLI/import suites: automatic and explicit preview, create/read retention, update in both directions, no-fallback validation, and import source parity all pass.
- No database schema or backfill was added; the existing v4 expression column stores the editable source and `source_syntax` is response-only.
- No authorization, IPC transport, secret handling, command execution, daemon lifecycle, dependency, or engine evaluation boundary changed. RRULE and anchor remain the only recurring execution definition.
- The planned `docs/api.md` target was corrected to the focused Spec-Kit API contract because this repository has no such live documentation file. A new broad API manual would exceed this slice and issue #52 remains open.
- `git diff --check` passed. All 37 changed or untracked files decoded as strict UTF-8 without BOM and passed the mojibake marker scan.

## Issue disposition

- #50: partial delivery only; pull request uses `Refs #50` and leaves the epic open for GUI adoption and remaining end-to-end acceptance.
- #52: remains open for the complete dual-syntax documentation posture.
- #22: remains open for additional cron dialect breadth.

## Final Spec-Kit analysis

- 23/23 buildable requirements and measurable outcomes map to tasks (100%).
- 35/35 tasks use valid, unique IDs and map to a requirement, user story, or constitution-mandated delivery gate.
- No unresolved ambiguity, duplication, constitution conflict, dependency inversion, or unmapped work remains.
- The nonexistent `docs/api.md` planning target is explicitly resolved by T028 and this evidence log through the focused API contract, without expanding into issue #52.
