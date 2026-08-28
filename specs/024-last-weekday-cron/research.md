# Research: Last-Weekday Cron Parity

## Decision: Treat `weekdayL` as an explicit go-schedule subset extension

- **Decision**: Support one `weekdayL` term inside the existing five-field syntax and Sunday numbering (`0`, `7`, or `SUN` on input, canonical `0L` on output).
- **Rationale**: Cronie defines the five-field layout and 0/7 Sunday numbering but not `L`. Quartz defines suffix-`L` as the last named weekday under a different field layout and numbering. The product must state its exact contract.
- **Alternatives considered**: Claim Quartz compatibility (rejected because its numbering and field layout differ); refuse the extension (rejected because the native model already represents it losslessly).
- **Primary sources**: [Cronie crontab(5)](https://github.com/cronie-crond/cronie/blob/master/man/crontab.5), [Quartz CronTrigger tutorial](https://www.quartz-scheduler.org/documentation/quartz-2.5.x/tutorials/crontrigger.html).

## Decision: Reuse `cron.Field.Ordinal` with `-1`

- **Decision**: Use `Ordinal == -1` for last weekday, zero for ordinary semantics, and retain 1 through 5 for `#`.
- **Rationale**: The existing schedule parser and `rrule-go` already represent `last friday` as one weekday with `N() == -1`. Reusing that value avoids a parallel concept.
- **Alternatives considered**: Add a boolean flag (duplicates ordinal state); add a new AST node or recurrence entity (disproportionate).

## Decision: Route only one terminal day-of-week `L`

- **Decision**: Accept one numeric or named weekday prefix followed by one terminal `L`. Treat malformed atoms as day-of-week errors and recognizable combinations as named refusals.
- **Rationale**: The generic extension detector must continue naming day-of-month `L`, `W`, and unsupported combinations without accidentally accepting named forms such as `FRIL` elsewhere.
- **Alternatives considered**: Relax generic extension scanning (too broad); parse all Quartz `L` forms (outside scope and lossy).

## Decision: Reuse the existing human phrase pipeline

- **Decision**: Render supported cron as `last <weekday> of the month at HH:MM` before invoking the shared schedule parser.
- **Rationale**: That parser already produces `FREQ=MONTHLY;BYDAY=-1XX`, while shared task boundaries retain the original cron source.
- **Alternatives considered**: Construct RRULE directly in the cron package (duplicates domain parsing); modify every consumer (unnecessary).

## Decision: Preserve all exporter selector guards

- **Decision**: Broaden the existing monthly weekday helper only to permit occurrence `-1`, then format it as `weekdayL`. Keep all competing-selector, count/until, multi-time, and seconds guards.
- **Rationale**: These guards prevent a richer native recurrence from being silently reduced to a broader cron schedule.
- **Alternatives considered**: Export any negative ordinal (exceeds S024); use set-position forms (different semantics and unsupported syntax).

## Decision: Keep missing-date policy inert

- **Decision**: Permit last-weekday export under skip, last-valid, and next-valid policies.
- **Rationale**: Every month has a last occurrence of every weekday. Existing schedule and cron policy detection already excludes negative ordinal weekdays from date-bearing fallback behavior.
- **Alternatives considered**: Require skip (unnecessary refusal); change missing-date resolution (incorrectly expands scope).

## Integration Surface Findings

- `cron convert` and `cron explain` already route through shared conversion functions; no CLI production change is required.
- Crontab scanning uses `Explain`, and task creation retains the original expression with explicit cron syntax.
- `scheduleinput.Parse` is the shared API/task boundary and restores original cron source after phrase parsing.
- API create and update paths already avoid mutation when parsing fails; focused regression coverage is sufficient.
- Production changes are confined to `internal/cron/cron.go`, `internal/cron/phrase.go`, and `internal/cron/export.go`.
