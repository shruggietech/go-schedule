# Quickstart: Cron interoperability and calendar-anomaly policy

**Feature**: `008-cron-interop` | **Date**: 2026-07-23

End-to-end walkthrough for verifying the feature on a built daemon. Assumes
`goschedd` is running and `gosched` is on `PATH`.

## 1. Translate one expression

```bash
gosched cron explain "0 9 * * 1-5"
```

Expect the phrase `weekdays at 09:00` and three upcoming run times. Nothing is
created.

```bash
gosched cron explain "@reboot"
```

Expect a named refusal, exit code 0 — a refusal is an answer, not a failure.

```bash
gosched cron explain "0 0 13 * 5"
```

Expect a refusal citing cron's day-of-month/day-of-week OR semantics.

```bash
gosched cron explain "*/7 * * * *"
```

Expect a refusal citing the non-dividing step.

## 2. Preview an import

Create a sample crontab:

```bash
printf 'MAILTO=ops@example.com\n# nightly\n0 2 * * * /usr/local/bin/backup\n*/15 * * * * /usr/local/bin/probe\n@reboot /usr/local/bin/warm\n' > /tmp/sample.crontab
```

```bash
gosched cron import --file /tmp/sample.crontab --dry-run
```

Expect: the `MAILTO` line reported as a warning, the comment skipped, two lines
translated with phrases and run times, the `@reboot` line declined, a summary
reading 5 lines / 0 created / 1 skipped / 1 declined, and the fidelity paragraph
naming the timezone applied and the policy defaults. Confirm with
`gosched task list` that nothing was created.

## 3. Import for real

```bash
gosched cron import --file /tmp/sample.crontab
```

Expect the identical report with `2 created`. Then confirm the phrase the
preview promised is the phrase the task reports:

```bash
gosched task list
gosched task show <id>
```

## 4. Export back out

```bash
gosched cron export
```

Expect a crontab line per expressible task and a `# declined:` comment naming
every task that cannot be represented — including any one-off and any disabled
task. Every task must appear exactly once.

## 5. Calendar-date scheduling and the missing-date policy

```bash
gosched task add month-end --command /usr/local/bin/close \
  --schedule "on the 31st of every month at 09:00" --tz America/New_York
gosched task show <id>
```

With the default policy, expect the next runs to skip 30-day months and February,
and the schedule description to say so rather than claiming every month.

```bash
gosched task edit <id> --missing-date last_valid
gosched task show <id>
```

Expect a run in every month — the 31st, the 30th, or the last day of February —
and a description naming the fallback.

```bash
gosched task edit <id> --missing-date next_valid
gosched task show <id>
```

Expect a 30-day month to produce a run on the 1st of the following month, with
that month's own 31st still present.

Then confirm the policy and phrase are independent:

```bash
gosched task edit <id> --schedule "on the 15th of every month at 09:00"
gosched task show <id>
```

The policy must still read `next_valid`.

## 6. Yearly and leap day

```bash
gosched task add leapday --command /usr/local/bin/audit \
  --schedule "every year on february 29 at 09:00" --missing-date last_valid
gosched task show <id>
```

Expect a run on 29 February in leap years and 28 February otherwise.

## 7. GUI

Open the GUI (`gosched gui`), edit a task, expand **Advanced Settings**, and
confirm a **Missing dates** selector sits alongside Overlap and Catch-up,
prefilled from the task and round-tripping on save. Confirm no cron input exists
anywhere in the interface.

## 8. Migration safety

Open a database written before this feature (or a copy of one) and confirm every
task reads `missing_date_policy: skip` and that `gosched task show` reports the
same upcoming runs it did before the upgrade.
