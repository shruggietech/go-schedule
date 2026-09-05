# Implementation Plan: Explicit DST Scheduling Intent

**Branch**: `codex/027-dst-intent` | **Date**: 2026-08-29 | **Spec**: [spec.md](spec.md)

## Summary

Complete issue #8 with one task-level scheduling policy set. Preserve the existing wall-clock default, add fixed elapsed and UTC recurrence bases, make spring gaps and fall overlaps configurable, migrate stored tasks, and route the same policy through every occurrence-producing and task-authoring surface.

## Technical Context

**Language/Version**: Go 1.25.0 and Markdown
**Primary Dependencies**: Existing `rrule-go`, Cobra, SQLite, and Fyne dependencies; no additions
**Storage**: SQLite schema v7 with three additive task columns and one nullable schedule epoch
**Testing**: Go unit/integration tests, real IANA transition matrices, migration/restart/catch-up coverage, GUI headless tests, benchmarks, and eight canonical gates
**Target Platform**: Windows, Linux, and macOS-supported Go paths
**Project Type**: Local daemon, CLI, API, and desktop application
**Performance Goals**: Existing p99 dispatch target below 100 ms; representative policy-aware next-run benchmark within ten percent of baseline
**Constraints**: Compatibility defaults, exact UTC instants, deterministic tests, one authoritative recurrence contract, UTF-8 without BOM
**Scale/Scope**: Domain policy model, timezone resolver, recurrence/catch-up/engine, migration/CRUD, API/CLI/GUI, calendar, docs, and issue #8

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One typed policy set replaces additional positional parameters and is consumed by the existing recurrence engine. |
| II. Testing Standards | PASS | Resolver, basis, migration, lifecycle, API, CLI, GUI, and real-transition tests begin red and remain deterministic. |
| III. UX Consistency | PASS | API, CLI, desktop, preview, detail, calendar, catch-up, and dispatch share the same values and evaluator. |
| IV. Performance | PASS | Fixed-duration arithmetic is constant time; wall-clock search is bounded and benchmarked against the ten-percent threshold. |
| V. Autonomous Execution | PASS | S027 follows full Spec-Kit, blocking analysis, test-first implementation, all eight gates, local commit, and one pre-publication halt. |

### Post-design re-check

All gates remain satisfied. Schema v7 is forward-only and additive. No dependency, permission, service, pinned workflow, or release-process change is required.

## Architecture and Decision Log

### Correct the defect classification without changing compatibility

The issue calls a five-hour spring or seven-hour fall gap in a six-hour local cycle a defect. That sequence is valid wall-clock anchoring because local readings remain six clock hours apart. S027 deliberately deviates from that diagnosis: `wall_clock` remains the default, while `elapsed` supplies exact-duration behavior. This avoids a silent schedule change for existing tasks.

### Pass one task-level policy value

`domain.SchedulePolicy` carries the existing missing-date policy plus time basis, gap policy, and overlap policy. `Effective()` normalizes zero values to compatibility defaults. Schedule, catch-up, engine, preview, detail, and calendar accept this value rather than growing a fragile list of positional policy arguments.

### Separate three recurrence bases

- `wall_clock` evaluates recurrence fields as floating local calendar intent and resolves each intent through the selected gap/overlap policies.
- `elapsed` is limited to fixed-duration interval shapes. It computes the next instant arithmetically from a persisted absolute epoch, preserving exact duration and avoiding timezone iteration. The epoch is bound once in the authoring timezone, so later presentation-timezone changes cannot move the phase.
- `utc` evaluates recurrence fields and missing-date choices in UTC. The task timezone remains presentation metadata only.

Calendar-selected monthly/yearly shapes are refused under `elapsed`. Daily and weekly interval shapes are accepted only when they select one occurrence per period and therefore have an exact 24-hour or 168-hour duration.

### Resolve wall intent to zero, one, or two instants

The timezone package gains a policy-aware resolver that discovers exact IANA mappings rather than assuming a one-hour transition. A missing wall reading returns no instant for `skip` or the first later valid instant for `next_valid`; an ambiguous reading returns first, both, or last. Existing `WallTime` remains a compatibility wrapper for `next_valid` plus `first`.

Wall-clock recurrence evaluation uses floating calendar values so the recurrence library cannot silently choose an offset. Missing-date and calendar adjustment paths select their intended date first, then use the same wall resolver. Duplicate concrete instants are suppressed.

### Migrate and expose independently of schedule replacement

Schema v7 adds `time_basis`, `dst_gap_policy`, and `dst_overlap_policy` with compatibility defaults, plus a nullable absolute elapsed epoch on Schedule. Like missing-date policy, the policy values live on Task and survive schedule replacement or unrelated edits. Create, update, preview, CLI flags, and GUI Advanced Settings validate the same enum values and elapsed compatibility.

### Bound fall-overlap lookup to the transition

An ambiguous wall interval is located from the nearby IANA offset transition. The evaluator jumps directly to the first recurrence intent that could still produce a second-fold instant, then compares it with the normal forward candidate. Dense rules therefore do not enumerate a 52-hour window.

## Project Structure

```text
specs/027-dst-intent/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/dst-policy.md
├── checklists/
└── tasks.md

internal/domain/
internal/timezone/
internal/schedule/
internal/catchup/
internal/engine/
internal/store/
internal/api/server/
internal/cli/
gui/
docs/
CHANGELOG.md
README.md
```

**Structure Decision**: Extend the existing single Go module and its established package boundaries. No new service or parallel scheduling engine is introduced.

## Complexity Tracking

No constitution violations require justification.
