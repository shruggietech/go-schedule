# Quickstart: Validate General Five-Field Cron Breadth

## Local conversion and explanation

```sh
gosched cron convert "0 9,17 * * *"
gosched cron explain "*/10 9-17 * * 1-5" --timezone America/New_York --count 5
```

Expected: both expressions are accepted, every restriction appears in the description, and upcoming runs match the field sets.

## Task preview and authoring

```sh
gosched task preview --schedule "30 8-17 * * 1-5" --schedule-syntax cron --timezone America/New_York
gosched task add --name business-hours --command /usr/bin/report --schedule "30 8-17 * * 1-5" --schedule-syntax cron --timezone America/New_York
```

Expected: preview and created task have the same normalized cron source, readable summary, and upcoming runs.

## Crontab round trip

Create a temporary crontab containing:

```text
0 9,17 * * * /usr/bin/report
0 0 1,15 JAN,MAR * /usr/bin/archive
```

Run dry-run import, real import, and export. Expected: two jobs in each preview, one task per line, and canonical exported expressions that re-import to identical run times.

## Boundary checks

- Compare `*/7 * * * *` across `:56` to the next hour's `:00`.
- Compare a date list across February in leap and common years under every missing-date policy.
- Compare selected `02:xx` times across spring-forward and fall-back in `America/New_York`.
- Confirm `0 0 13 * 5`, `@reboot`, Quartz fields, and modifier composites remain non-mutating refusals.

## Verification

```sh
sh scripts/verify.sh all
```

Also run the broad-recurrence benchmarks and compare against the baseline recorded in `verification.md`.
