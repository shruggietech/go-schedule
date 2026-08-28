# Research: Ordinal-Weekday Cron Parity

## Decision: Treat `#` as an explicit go-schedule subset extension

- **Decision**: Support `weekday#ordinal` inside the existing five-field syntax
  and existing Sunday numbering (`0` or `7` on input, canonical `0` on output).
- **Rationale**: Cronie documents the ordinary five-field day-of-week field as
  0 through 7 or names, but does not define `#`. Quartz defines `#` in a
  different dialect and weekday numbering. The product must describe its exact
  contract instead of implying universal portability.
- **Alternatives considered**: Claim Quartz compatibility (rejected because its
  field layout and numbering differ); refuse the extension (rejected because
  the native recurrence model already represents it losslessly).
- **Primary sources**:
  [Cronie crontab(5)](https://github.com/cronie-crond/cronie/blob/master/man/crontab.5),
  [Quartz CronTrigger tutorial](https://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html).

## Decision: Add one optional value to `cron.Field`

- **Decision**: Add `Ordinal int`; zero means ordinary field semantics and 1
  through 5 modifies one normalized day-of-week value.
- **Rationale**: The ordinal is syntax-level metadata attached to a weekday.
  Keeping the weekday in `Values` preserves existing normalization and
  `Single()` behavior without contaminating the domain recurrence model.
- **Alternatives considered**: Encode the pair into `Values` (opaque and error
  prone); add a new cron AST node or schedule entity (disproportionate to one
  bounded extension).

## Decision: Separate malformed syntax from recognizable unsupported syntax

- **Decision**: Missing halves, extra separators, nonnumeric ordinals, and
  ordinals outside 1 through 5 are day-of-week errors. Lists, ranges, steps,
  multiple ordinal terms, non-day-of-week placement, and restricted date/month
  combinations are named refusals.
- **Rationale**: Operators can distinguish a typo from a valid-looking dialect
  or combination that go-schedule intentionally does not model.

## Decision: Reuse the existing monthly phrase pipeline

- **Decision**: Render supported cron as `<ordinal> <weekday> monthly at HH:MM`
  before invoking the existing human schedule parser.
- **Rationale**: That parser already produces `FREQ=MONTHLY;BYDAY=+nXX`, and all
  task boundaries already retain original cron source after successful parsing.
- **Alternatives considered**: Construct RRULE directly in the cron package
  (duplicates domain parsing and risks drift); modify every consumer (unneeded).

## Decision: Gate export by selector shape and effective missing policy

- **Decision**: Export exactly one positive numbered weekday from a monthly rule
  when no competing date selector exists. Permit first through fourth under any
  missing-date policy; require effective skip for fifth.
- **Rationale**: Every month has the first four occurrences. A fifth can be
  absent, and cron skips that month. The current exporter already detects this
  through its date-bearing policy logic.
- **Alternatives considered**: Reject non-skip policy for every ordinal
  (unnecessarily refuses equivalent behavior); approximate fifth last/next
  behavior (violates fidelity).

## Integration Surface Findings

- `cron convert` and `cron explain` already route through `Convert` and
  `Explain`; no CLI implementation change is required.
- Crontab scanning uses `Explain`, and task creation retains the original
  expression with `schedule_syntax: cron`.
- `scheduleinput.Parse` is the shared API/task boundary and restores original
  cron source after phrase parsing.
- Existing API create/update flows already avoid mutation when parsing fails;
  focused regression coverage is sufficient.
- Production changes are confined to `internal/cron/cron.go`,
  `internal/cron/phrase.go`, and `internal/cron/export.go`.
