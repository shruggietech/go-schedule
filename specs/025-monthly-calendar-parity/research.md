# Research: Monthly Calendar Cron Parity

## Decision: Support an explicit five-field subset

- **Decision**: Accept day-of-month `L`, `1W` through `31W`, and `LW` only with wildcard month and day-of-week.
- **Rationale**: Cronie defines the five-field baseline and day-field OR behavior but not these modifiers. Quartz defines their meaning in a different layout. The product must name its bounded contract.
- **Alternatives considered**: Claim Quartz compatibility (incorrect field contract); accept composites (not faithfully modeled); split into more slices (contrary to the requested outcome-sized delivery).
- **Primary sources**: [Cronie crontab(5)](https://github.com/cronie-crond/cronie/blob/master/man/crontab.5), [Quartz CronTrigger tutorial](https://www.quartz-scheduler.org/documentation/quartz-2.5.x/tutorials/crontrigger.html).

## Decision: Represent `L` and `LW` with standard RRULE

- **Decision**: Compile `L` as `FREQ=MONTHLY;BYMONTHDAY=-1` and `LW` as monthly weekdays with `BYSETPOS=-1`.
- **Rationale**: These shapes select exactly one intended date every month and are executable by `rrule-go`.
- **Alternatives considered**: Custom metadata for all selectors (duplication); source-text computation (violates separation).

## Decision: Add a typed durable adjustment for `nW`

- **Decision**: Add `Schedule.CalendarAdjustment` with one non-empty value, `nearest_weekday`, persisted by schema v6. The carrier RRULE holds the numbered target and time.
- **Rationale**: `nW` conditionally selects `n`, `n-1`, `n+1`, `n+2`, or `n-2`; no single RFC 5545 RRULE represents that across all months. Execution metadata must survive restart independently of editable source.
- **Alternatives considered**: Multiple RRULEs (model and collision expansion); fake X-properties (parser rejects them); infer from `Expression` (inert and possibly absent); decline `nW` (fails scope).

## Decision: Extend the existing monthly period walker

- **Decision**: Resolve marked rules by month, apply existing missing-date policy to the numbered target, then a bounded nearest-weekday adjustment and existing wall-time normalization.
- **Rationale**: This preserves established policy and DST meaning. Unadjusted schedules retain the current path.
- **Alternatives considered**: Post-process rrule-go occurrences (cannot recover skipped intent); scan days indefinitely (weaker bound).

## Decision: Reject invalid stored combinations

- **Decision**: Unknown adjustment values, non-monthly carriers, multiple/non-positive month days, and competing selectors return actionable errors or refusals.
- **Rationale**: Silently ignoring persisted execution metadata can run a job on the wrong date.
- **Alternatives considered**: Treat malformed metadata as empty (unsafe); normalize on read (hides corruption).

## Decision: Preserve policy fidelity in export

- **Decision**: `L`, `LW`, and `1W` through `28W` export under any policy. `29W` through `31W` export only under effective skip.
- **Rationale**: Targets through 28 always exist; a non-skip policy changes larger targets and cron cannot carry it.
- **Alternatives considered**: Require skip for all `nW` (unnecessary refusal); always export larger targets (lossy).

## Decision: Reuse shared workflows and fix GUI preview

- **Decision**: Keep `scheduleinput.Parse` as the one task-entry boundary. Pass selected missing-date policy in editor previews and refresh on policy changes.
- **Rationale**: Existing interfaces already share parsing and retained-source behavior. The preview omission is an existing consistency bug exposed by policy-sensitive `nW`.

## Integration Findings

- Store CRUD enumerates schedule columns, so migration and all queries change together.
- JSON uses `omitempty` for the adjustment, preserving ordinary response shape.
- Failed update tests include another requested field mutation to prove whole-request atomicity.
- Crontab import and editing retain exact normalized cron and `cron` syntax.
- DOM `LW` (last weekday) remains distinct from DOW `5L` (last Friday).
- Issue #22 remains a broader epic; S025 uses `Refs #22` and does not close it.
