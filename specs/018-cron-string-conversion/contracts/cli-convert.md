# CLI Contract: `gosched cron convert`

**Feature**: `018-cron-string-conversion` | **Date**: 2026-08-28

## Invocation

```text
gosched cron convert [--to cron|human] <schedule-string>
```

The command accepts exactly one positional string. Shell users quote schedules
containing spaces. The global `--json` flag applies before or after the command
group according to the existing CLI convention.

## Automatic conversion

```text
$ gosched cron convert "0 9 * * 1-5"
weekdays at 09:00

$ gosched cron convert "weekdays at 09:00"
0 9 * * 1-5
```

After trimming surrounding whitespace:

1. a value beginning with `@` is cron input;
2. a value containing five whitespace-delimited fields whose first field has
   cron's numeric/wildcard minute shape is cron input;
3. every other value is human input.

Classification occurs once. Invalid cron does not fall through to human input.

## Explicit destination

`--to` names the output:

- `--to cron` forces the supplied value to be parsed as human;
- `--to human` forces the supplied value to be parsed as cron.

Any other value is a usage error. The override supports validation while
existing five-field human phrases remain automatically recognizable.

## Text output

Success writes exactly the canonical converted string plus one newline to
stdout and nothing to stderr. Exit code: 0.

Malformed or unfaithful input writes no stdout. The stderr diagnostic names the
invalid syntax or fidelity reason. Exit code: 2.

The command never includes source labels, arrows, previews, times, daemon
status, or task data.

## Structured output

Success writes one JSON object to stdout and exits 0:

```json
{
  "input_syntax": "cron",
  "output_syntax": "human",
  "input": "0 9 * * 1-5",
  "output": "weekdays at 09:00",
  "refusal_reason": ""
}
```

Malformed or unfaithful input writes the same stable shape to stderr, leaves
stdout empty, and exits 2:

```json
{
  "input_syntax": "cron",
  "output_syntax": "human",
  "input": "61 9 * * *",
  "output": "",
  "refusal_reason": "cron: minute field: value 61 is outside 0-59"
}
```

All five keys are always present. No second plain-text diagnostic follows a
structured refusal.

## Fidelity rules

- Cron-to-human uses the existing five-field/shorthand parser and refusal
  reasons used by explain/import.
- Human-to-cron emits the existing canonical expression for a schedule only
  when the five-field form preserves run times.
- Creation-aligned schedules without an explicit phase/time are refused.
- One-offs, lossy elapsed intervals, ordinal weekdays, non-default
  missing-date behavior, and other existing unsupported shapes are refused by
  name.
- Day-of-month plus day-of-week OR semantics and field-local step behavior are
  never approximated.

## Purity boundary

The command makes no daemon, IPC, API, network, filesystem, configuration, or
task-state call. `cron explain` remains the richer preview command; import and
export retain their task/crontab roles. Cron remains invalid in task authoring
until #50 is delivered.
