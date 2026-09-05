# Contract: Dual-Syntax Schedule Input API

**Feature**: `019-dual-syntax-input` | **Date**: 2026-08-28

## Shared request field

The following request bodies accept:

```json
{
  "schedule": "0 9 * * 1-5",
  "schedule_syntax": "cron"
}
```

`schedule_syntax` is optional. Accepted non-empty values are `human` and `cron` after trimming/case normalization. Omission selects automatic classification.

Applies to:

- `POST /v1/schedules/preview`
- `POST /v1/tasks`
- `PATCH /v1/tasks/{id}` when replacing a recurring schedule

## Automatic selection

After trimming the schedule:

1. `@` prefix is cron;
2. exactly five fields whose minute field has cron numeric/wildcard shape is cron;
3. every other value is human.

Selection happens once. A selected parser failure never falls back.

## Preview success

Request:

```json
{
  "schedule": "0 9 * * 1-5",
  "timezone": "America/New_York"
}
```

Response status `200`:

```json
{
  "rrule": "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=9;BYMINUTE=0",
  "human_summary": "Every weekday at 09:00",
  "source_syntax": "cron",
  "next_runs": []
}
```

The exact RRULE ordering remains the existing parser contract. `next_runs` contains the normal preview count.

## Task create/update/read success

The existing task response retains its shape. The nested schedule adds:

```json
{
  "schedule": {
    "expression": "0 9 * * 1-5",
    "source_syntax": "cron"
  }
}
```

- `expression` is durable editable source.
- `source_syntax` is derived response metadata.
- RRULE/anchor remain the execution definition.
- One-off and expressionless legacy schedules omit both source identity and any fabricated recurring expression.

## Validation errors

Invalid hint:

```json
{
  "error": {
    "code": "validation_failed",
    "field": "schedule_syntax",
    "message": "must be human or cron"
  }
}
```

Selected-parser failure:

```json
{
  "error": {
    "code": "validation_failed",
    "field": "schedule",
    "message": "cron: minute field: value 61 is outside 0-59"
  }
}
```

Both return status `400`. No task/schedule mutation occurs. A non-empty syntax hint without a recurring schedule also names `schedule_syntax`.

## Explicit-hint examples

- `schedule_syntax: cron` with `weekdays at 09:00` fails as cron; no human fallback.
- `schedule_syntax: human` with `0 9 * * 1-5` fails as human; no cron fallback.
- An omitted hint accepts each through structural selection.

## CLI contract

`gosched task add/edit --schedule` accepts a human-readable or supported cron recurrence and sends it unchanged to the daemon. There is no separate CLI-only detector or hint flag in this slice.

Cron import sends `Line.Expr` plus explicit cron identity for both upcoming-run preview and task creation. It continues to display `Line.Phrase` as explanation.

## GUI compatibility contract

Existing GUI preview/create/update requests explicitly send `human`. This preserves its current validator and does not claim GUI cron support.
