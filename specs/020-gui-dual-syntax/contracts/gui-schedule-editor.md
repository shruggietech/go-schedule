# Contract: GUI Schedule Editor

## Accepted input

The existing Schedule field accepts the human phrases and supported five-field cron expressions accepted by `internal/scheduleinput`. Automatic selection uses the central structural detector. Cron-shaped failures are not retried as human.

## Preview request

For non-empty recurring text that passes local parsing, the editor sends:

```json
{
  "schedule": "0 9 * * 1-5",
  "schedule_syntax": "cron",
  "timezone": "America/New_York"
}
```

The expression is trimmed and otherwise retained. Human input sends `human`. The daemon remains authoritative for summary and upcoming runs.

## Create and update requests

- A recurring create sends the normalized expression and selected syntax.
- A recurring replacement update sends both values from current field text.
- An unchanged retained cron expression stays cron.
- A syntax switch follows the new text, not response metadata.
- A one-off sends `at` and no recurring expression or syntax.
- A blank expression on a compatible expressionless legacy edit omits syntax and preserves the stored schedule.

## Validation and refusal

Save eligibility requires a successful shared parse plus existing required-field and timezone validation. Invalid or faithfully unsupported cron keeps Save disabled and displays the central reason without internal conversion wording.

## Prefill and help

Recurring edit prefill uses `Schedule.Expression` exactly. Help keeps human examples first, adds copy/pasteable cron examples, names the five-field scope, and points to `docs/cron.md` for fidelity details.
