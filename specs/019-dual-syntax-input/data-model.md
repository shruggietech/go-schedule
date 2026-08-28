# Data Model: Dual-Syntax Task Input Foundation

**Date**: 2026-08-28

## Syntax

Identifies the grammar used for one recurring source expression.

| Value | Meaning |
| --- | --- |
| empty | Automatic request selection, or no recurring source identity in a response |
| `human` | Existing human-readable schedule language |
| `cron` | Existing supported five-field cron dialect or supported shorthand |

Rules:

- Request hints normalize surrounding whitespace and ASCII case.
- Only empty, `human`, and `cron` are valid.
- Explicit values select exactly one parser.
- Automatic classification uses S018's shared structural detector once.

## Schedule Input Result

| Field | Rule |
| --- | --- |
| `expression` | Submitted value with surrounding whitespace removed |
| `source_syntax` | Detected or explicitly selected grammar |
| `schedule` | Existing recurring schedule with RRULE, anchor, summary, and retained expression |

Invariants:

- Success always has non-empty expression and source syntax.
- Failure creates no schedule and never retries another parser.
- The returned schedule's `Expression` equals `expression`.
- RRULE/anchor are derived through the existing human schedule parser.
- Raw cron is never evaluated by the engine.

## Stored Schedule

The existing storage shape is unchanged.

| Existing field | S019 contract |
| --- | --- |
| `rrule` | Authoritative recurring execution rule |
| `anchor` | Authoritative recurrence phase/start context |
| `human_summary` | Display description, not reparsed |
| `expression` | Inert editable source, human phrase or cron expression |
| `run_at` | One-off timing; no recurring source syntax |

`source_syntax` is response metadata, not a database column. Server response
construction derives it only when kind is recurring and expression is non-empty.

## API Request Additions

Preview, create, and update add optional `schedule_syntax`.

State rules:

- With non-empty `schedule`, empty hint means automatic selection.
- With non-empty `schedule`, valid hint selects the parser.
- A non-empty hint without a recurring schedule is invalid.
- Update without `schedule` preserves the existing stored schedule and identity.

## API Response Additions

- Recurring schedule with retained expression: `source_syntax` is `human` or
  `cron`.
- One-off schedule: omitted/empty.
- Legacy recurring schedule with empty expression: omitted/empty.
- Preview: `source_syntax` is always non-empty on success.

## Imported Job Transition

```text
crontab line
  -> existing scanner produces Line.Expr + Line.Phrase + command/args
  -> report displays Expr and Phrase
  -> preview Line.Expr with hint cron
  -> create Line.Expr with hint cron
  -> store RRULE/anchor plus Expression=Line.Expr
  -> task response derives source_syntax=cron
```

Declined/invalid lines never enter the schedule-input state transition.

## Compatibility

- No schema version changes.
- Existing human expressions infer `human` and keep prior bytes after trimming
  rules already applied by `schedule.Parse`.
- Existing expressionless schedules remain readable.
- Unknown additive JSON response fields are safe for existing clients.
- Existing request senders that omit the hint continue in automatic mode.
