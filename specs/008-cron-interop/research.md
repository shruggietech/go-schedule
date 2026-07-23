# Phase 0 Research: Cron interoperability and calendar-anomaly policy

**Feature**: `008-cron-interop` | **Date**: 2026-07-23

This records what was read in the existing tree, the decisions taken under the
Build-Phase Autopilot Protocol decision policy, and the alternatives rejected.

## What the tree already provides

| Concern | Where it lives today | Bearing on this feature |
| --- | --- | --- |
| Phrase → RRULE | `internal/schedule/parse.go` — four regex-driven parse arms, each returning RRULE parts plus a `HumanSummary`, funnelled through `finish` | Two new arms follow the identical shape; helpers (`maybeTime`, `byTime`, `parseTimeOfDay`, `ordinalWord`) are reused verbatim |
| Recurrence evaluation | `internal/schedule/recur.go` — `nextRecurring` builds an rrule-go rule at the task's anchor and normalizes day-or-coarser occurrences through `timezone.WallTime` | The missing-date policy resolves *before* that normalization; the DST step is untouched |
| Cron equivalence | `internal/schedule/cronparity_test.go` — five cron patterns paired with phrases, asserting matching run times | Becomes the seed corpus for the converter's tests, extended rather than replaced |
| Schedule persistence | `internal/store/store.go` migrations v1–v4; v4 is the model for an additive column with a total default | Migration v5 copies that pattern exactly |
| Policy plumbing | `domain.OverlapPolicy` / `domain.CatchupPolicy` → `TaskCreateRequest`/`TaskUpdateRequest` validation → `internal/cli/task.go` flags → `gui/editor_data.go` generic `policyChoice[T]` | A third policy is a strict repetition of an established path, including the generic GUI choice type |
| CLI command shape | `internal/cli/cli.go` `newRoot()` registering `newTaskCmd()` etc.; `errUsage` for exit code 2; `--json` persistent flag | `newCronCmd()` slots in with no new conventions |

## Decisions

### D1 — The missing-date policy lives on the task, not the schedule

**Chosen**: add `MissingDatePolicy` to `domain.Task`; migration v5 adds the
column to `tasks`. `schedule.NextRun` and `schedule.UpcomingRuns` take the
policy as a parameter.

**Alternatives**:

- *On the schedule row.* Superficially cheaper — `NextRun(sch, tz, after)`
  already receives the schedule, so no signature changes. Rejected on reading
  `internal/api/server/update.go:122-138`: an edit that supplies a new phrase
  **creates a new schedule row** and repoints the task at it. A policy stored on
  the schedule would therefore be silently reset to the default on any phrase
  edit unless a carry-over is remembered at that one site. That is a silent
  run-time change — precisely what FR-024a and FR-026 forbid, and the class of
  defect issue #4 already produced in the task editor. Correctness that depends
  on remembering to copy a field is not correctness.
- *A parallel policy table.* Rejected as unjustified complexity for a
  single-valued enum (constitution: simplest design that satisfies the
  principles).

**Cost accepted**: six call sites change signature
(`internal/engine/engine.go:146,235`, `internal/catchup/catchup.go:31`,
`internal/api/server/calendar.go:74`, `internal/api/server/tasks.go:209,216`).
All are compile-checked and all already hold the task, so none needs new
plumbing. This reverses the working assumption in the pre-approval plan file,
which had proposed the schedule row before the update path was read.

### D2 — Cron converts by way of the phrase, never directly to a recurrence

**Chosen**: `cron.Phrase(spec) (string, error)` produces the phrase a user would
have typed; that phrase goes through the existing `schedule.Parse`. An
expression with no phrase is declined.

**Rationale**: makes the preview structurally truthful (FR-005, SC-002a) — the
preview shows the string that is actually parsed. Keeps one implementation of
phrase→RRULE. Makes "cron is not an authoring syntax" (FR-014) a property of the
architecture rather than of discipline: the converter has no privileged entry
point the user does not also have.

**Alternative rejected**: cron → RRULE directly. Faster to write and able to
express more (arbitrary by-minute lists), but it creates a second timing path
that can drift from the phrase grammar, and it would let cron express schedules
no user could author or edit afterwards — a task nobody could subsequently
change through either interface.

### D3 — The cron parser is written in-tree, with no new dependency

**Chosen**: `internal/cron` implements the five-field grammar itself.

**Rationale**: the constitution's Engineering Constraints prefer the standard
library where it suffices and require justification for every dependency. A cron
library (`robfig/cron`, `adhocore/gronx`) would parse fields we then have to
re-inspect field by field anyway, because the work is the *mapping* to phrases
and the *refusal* of what cannot map — neither of which a scheduling library
exposes. The grammar is roughly 120 lines and is fully covered by table tests.

**Alternative rejected**: `robfig/cron/v3`'s parser. It normalizes away exactly
the distinctions FR-002 and FR-003b require us to detect (it accepts `*/7`
happily and reports a schedule, not a field structure), so we would parse twice.

### D4 — Declines are outcomes, not failures

**Chosen**: an import or preview that reads its input reports success (exit 0)
even when every line is declined; the per-line report and the summary counts
carry the outcome. Failure is reserved for unreadable input, an invalid
timezone, or a malformed single expression given to `explain` (exit 2, via the
existing `errUsage` convention).

**Rationale**: constitution principle III — conventional exit codes, with errors
that state what to do. A crontab containing three `@reboot` lines is a
successful, informative conversion, not a command failure.

### D5 — Ordinal-weekday rules join the policy rather than being special-cased

**Chosen**: `5th friday monthly` is treated as a rule that may not occur in a
period, so `last_valid` resolves it to the last Friday and `next_valid` to the
first Friday of the next month.

**Rationale**: the spec's evidence (issue #8, finding 2) shows this rule already
skips two thirds of the calendar while claiming otherwise. Fixing the label
without offering the fallback would leave the operator informed but still unable
to state their intent.

### D6 — Non-dividing steps and DOM+DOW combinations are declined

Recorded in the spec's Clarifications (session 2026-07-23) and restated here
because they bound the converter's scope: `*/n` is accepted only where `n`
divides its field's range, and any expression restricting both day-of-month and
day-of-week is declined. Both are instances of FR-002's no-approximation rule.

## Open questions

None. Every clarification raised in `/speckit-clarify` was resolved under the
decision policy and integrated into the spec.
