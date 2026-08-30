# Data Model: Cron Parity Closure

## Task standard input

Task gains one optional execution field:

- `Stdin`: exact UTF-8 string supplied to the child process on each run; empty means no supplied payload.

Storage schema v8 adds `tasks.stdin TEXT NOT NULL DEFAULT ''`. The total empty default preserves every existing task. Create accepts a string. Update uses a nullable request field so omission preserves the value and an explicit empty string clears it.

## Crontab scan context

An in-memory context evolves from top to bottom:

- `ScheduleTimezone`: effective `CRON_TZ`, empty until assigned.
- `Environment`: effective ordinary variables, copied for each job.
- `Shell`: effective `SHELL`, default `/bin/sh`.

`CRON_TZ` changes only schedule context. `MAILTO` and `MAILFROM` produce deferred warnings and are not represented as notification delivery. `SHELL` remains in environment and also selects the executable.

## Imported job

Each accepted line carries:

- source line number and raw text;
- normalized timing expression and dialect;
- readable schedule description;
- schedule timezone snapshot;
- shell executable and `-c` command string;
- standard-input payload and whether a percent separator appeared;
- environment snapshot;
- optional run-as user from system layout;
- warnings and refusal reason.

Maps and strings are copied so later assignments cannot mutate earlier jobs.

## Cron timing specification

The parsed timing value contains six field sets:

- second;
- minute;
- hour;
- day of month;
- month;
- day of week.

Five-field Unix input supplies a synthetic singleton second zero and uses Unix weekday numbers. Six-field Quartz input parses all six source fields and uses Quartz weekday numbers. Both normalize into the recurrence model's Sunday-through-Saturday representation.

`?` is a full-field no-specific-value marker for one Quartz day field. It becomes an unrestricted internal field but remains distinguishable during validation.

## Durable recurrence mapping

Supported standard fields map to the existing recurring Schedule:

- `RRULE` contains selected seconds, minutes, hours, and optional calendar fields;
- `Anchor` remains the creation instant;
- `HumanSummary` is display-only and seconds-aware;
- `Expression` retains normalized source and remains inert;
- existing calendar adjustment metadata remains authoritative for focused modifiers.

No cron runtime evaluator or new schedule table field is added.

## Compatibility invariants

- Existing five-field expressions still execute at second zero.
- Existing five-field export remains five fields.
- Existing tasks migrate with empty standard input.
- GUI updates that do not expose stdin preserve imported stdin through pointer-based API omission.
- Platform run-as validation remains authoritative at task creation.
