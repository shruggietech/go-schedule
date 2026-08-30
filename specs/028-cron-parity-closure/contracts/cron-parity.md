# Contract: Cron Parity Closure

## Import layout

```text
gosched cron import --file PATH [--dialect unix|quartz] [--system | --run-as USER] [--dry-run]
```

- `unix` is the default and consumes five timing fields.
- `quartz` consumes six timing fields beginning with seconds.
- `--system` consumes one user field after timing.
- `--run-as` supplies the owner identity absent from a user crontab file.
- Layout flags are orthogonal and never inferred from command tokens.

## Assignment contract

- `CRON_TZ=value` controls the schedule timezone of following jobs.
- `NAME=value` controls the child environment of following jobs.
- Matching outer single or double quotes preserve boundary whitespace and are removed.
- `SHELL=value` also selects the program invoked with `-c`.
- `MAILTO` and `MAILFROM` remain visible deferred warnings.
- An explicit CLI timezone overrides file schedule timezones; otherwise each line uses its effective `CRON_TZ` or `Local`.

## Command and percent contract

The command portion before the first unescaped percent is passed intact as one argument to the effective shell. Escaped percent becomes literal percent. The first unescaped percent starts standard input; later unescaped percent signs become newlines. Quote characters do not suppress cron percent processing.

## Six-field timing contract

Accepted shape:

```text
second minute hour day-of-month month day-of-week
```

- Seconds accept 0 through 59, lists, ranges, wildcard steps, range steps, and value-start steps.
- `?` is accepted only as an entire day-of-month or day-of-week field.
- Quartz numeric weekdays are 1=Sunday through 7=Saturday.
- Seven-field year expressions are refused.
- Existing DOM/DOW and modifier fidelity refusals remain.

## API task input contract

Create adds optional `stdin`. Update adds optional nullable `stdin`; omitted preserves, string replaces, and empty string clears. Task detail returns the stored field through the existing local API.

## Export contract

- Second zero produces conventional five-field timing.
- Required seconds produce six-field Quartz timing.
- Task environment, run-as, or standard input causes a named operational-context refusal rather than lossy output.
- Every task still produces either a line or a commented refusal.

## Error contract

Malformed fields, invalid timezones, missing system users or commands, misplaced `?`, unsupported Quartz year, and incompatible day combinations identify the affected field or feature. Import continues reporting independent lines; task creation failures do not roll back tasks already created.
