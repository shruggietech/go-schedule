# Data Model: GUI Dual-Syntax Scheduling

## Editor recurring input

| Field | Type | Source | Rule |
| --- | --- | --- | --- |
| Raw text | string | Schedule entry | Operator-editable; surrounding whitespace is insignificant. |
| Expression | string | `scheduleinput.Parse` | Normalized text sent to preview and save. |
| Syntax | `scheduleinput.Syntax` | `scheduleinput.Parse` | `human` or `cron`, derived from current text. |
| Compiled schedule | `domain.Schedule` | `scheduleinput.Parse` | Used only to prove local validity; daemon supplies visible preview timing. |
| Validation error | error | `scheduleinput.Parse` | Named parser or fidelity refusal; no fallback. |

## Editor form

The private `taskForm` gains `scheduleSyntax string` beside `schedule`.

- Recurring create: both values are required after successful parsing.
- Recurring update with replacement: both values are sent.
- Expressionless legacy update: both remain empty, preserving stored timing.
- One-off create/update: both remain empty and `at` is sent instead.

## Task detail prefill

`domain.Task.Schedule.Expression` is the sole Schedule field prefill. The
response's `SourceSyntax` is descriptive metadata, not mutable editor state.
The current text is reclassified after every edit.

## State transitions

```text
blank new recurring -> invalid
supported human -> valid human preview/save
supported cron -> valid cron preview/save
cron-shaped invalid/refused -> invalid, no human fallback
valid syntax A -> edited syntax B -> preview/save use B
blank existing legacy recurrence -> valid preserve operation
one-off mode -> recurring expression and syntax omitted
```

No persisted entity or API schema changes in this slice.
