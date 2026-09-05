# Quickstart: Validate Dual-Syntax Task Input

**Date**: 2026-08-28

## Prerequisites

- Build or run the daemon and CLI from this branch.
- Use a temporary data directory or disposable test server.
- Choose a timezone with DST for parity checks, such as `America/New_York`.

## 1. Preview cron and human equivalents

Submit these to `POST /v1/schedules/preview` with the same timezone and policy:

```json
{"schedule":"0 9 * * 1-5","schedule_syntax":"cron","timezone":"America/New_York"}
```

```json
{"schedule":"weekdays at 09:00","schedule_syntax":"human","timezone":"America/New_York"}
```

Expected:

- source identities differ;
- RRULE and upcoming runs agree;
- the cron response retains no generated English as source.

## 2. Create and fetch a cron task

```text
gosched task add report --command /bin/true --tz America/New_York --schedule "0 9 * * 1-5"
```

Fetch the task using `gosched task show` or `GET /v1/tasks/{id}`.

Expected nested schedule evidence:

```json
{
  "expression": "0 9 * * 1-5",
  "source_syntax": "cron"
}
```

RRULE and anchor remain populated exactly as for a human-authored task.

## 3. Replace syntax and preserve unrelated policy

Create a task with a non-default missing-date policy, replace its schedule from human to cron and back, then edit only the command.

Expected:

- schedule replacements update expression and source identity together;
- command-only update changes neither;
- missing-date policy remains unchanged.

## 4. Exercise named refusals

Preview each with automatic or explicit cron selection:

```text
61 9 * * *
@reboot
0 9 1 * 1
```

Expected: status 400, field `schedule`, named cron/fidelity reason, no fallback, and no task mutation.

Also submit an invalid `schedule_syntax` and verify the error names that field.

## 5. Import and retrieve source

Import:

```text
0 9 * * 1-5 /usr/local/bin/report
```

Expected:

- report still displays the explanatory weekday phrase;
- preview runs and creation use the timing expression with cron identity;
- fetched task retains `0 9 * * 1-5` and reports `source_syntax: cron`;
- dry-run, refusal, and partial-success output retain prior behavior.

## 6. Run focused verification

```text
go test ./internal/cron ./internal/scheduleinput ./internal/api/server ./internal/cli ./gui -count=1
```

Then run the canonical foreground gate:

```text
sh scripts/verify.sh all
```

Record red/green and all eight gate results in `verification.md`.
