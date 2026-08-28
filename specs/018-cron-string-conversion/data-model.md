# Data Model: Pure Schedule String Conversion

**Date**: 2026-08-28

This feature persists no data. Its model consists of immutable values for one
local conversion.

## Syntax

Identifies one side of a conversion.

| Value | Meaning |
| --- | --- |
| `cron` | Supported five-field cron expression or shorthand |
| `human` | Existing human-readable schedule language |

Validation rules:

- A requested destination must be empty (automatic), `cron`, or `human`.
- Input and output syntax always differ.
- Automatic input identity is derived once and never changed after a parse
  failure.
- Five-field auto-detection also requires the first field to have a cron minute
  shape, preserving current five-field human phrases.

## Conversion

| Field | Rule |
| --- | --- |
| `input_syntax` | Detected or forced syntax of the supplied value |
| `output_syntax` | Opposite syntax and requested destination |
| `input` | Supplied value with surrounding whitespace removed |
| `output` | Canonical converted string on success; empty on refusal |
| `refusal_reason` | Empty on success; specific validation/fidelity reason on refusal |

Invariants:

- Exactly one of `output` and `refusal_reason` is non-empty.
- A refusal never carries an approximated expression.
- The value contains no daemon response, task identifier, timestamp, preview,
  or persisted state.

## State Transition

```text
raw request
  -> normalize surrounding whitespace
  -> validate/derive syntax identities
  -> parse selected input syntax exactly once
  -> success(output) OR refusal(reason)
  -> render text or structured stream
```

There is no retry into the other syntax and no storage lifecycle.

## Existing Domain Relationship

- Cron-to-human uses the current parsed cron specification and named
  `Unsupported` outcome.
- Human-to-cron uses a transient recurring `domain.Schedule` and its existing
  missing-date fidelity rules.
- Normal task export additionally owns enabled/state policy. The extracted
  schedule renderer does not weaken those task-level checks.
