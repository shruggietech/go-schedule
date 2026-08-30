# Quickstart: Validate Cron Parity Closure

## Seconds conversion and preview

```sh
gosched cron explain "*/30 * * * * *" --timezone UTC --count 4
gosched task preview --schedule "0 0 12 ? * MON" --schedule-syntax cron --timezone UTC
```

Expected: first output alternates second 0 and 30; second output selects Monday noon. A five-field expression still selects second zero and exports with five fields.

## User crontab fidelity

Create a temporary file:

```text
CRON_TZ=America/New_York
SHELL=/bin/sh
PATH=/usr/local/bin:/usr/bin
TZ=UTC
30 9 * * 1-5 printf '%s\n' "$TZ" && date%payload one%payload two
CRON_TZ=Europe/London
0 17 * * * echo done > /tmp/done
```

Run:

```sh
gosched cron import --file CRONTAB --dry-run
```

Expected: the first job schedules in New York, runs with `TZ=UTC`, retains shell operators, and reports stdin. The second schedules in London. Both inherit the environment snapshot.

## System and Quartz layouts

```sh
gosched cron import --file SYSTEM_CRONTAB --system --dry-run
gosched cron import --file QUARTZ_CRONTAB --dialect quartz --dry-run
```

Expected: system usernames appear as run-as identities, while Quartz files consume six timing fields. Removing either option changes only that selected layout and never triggers guessing.

## Compatibility and refusal checks

- Confirm ordinary five-field import and export fixtures are unchanged.
- Confirm `?` outside a Quartz day field and seven-field input fail without mutation.
- Confirm tasks with environment, run-as, or stdin produce explicit export refusals.
- Confirm `@reboot`, restricted DOM plus DOW, unsupported modifier composites, run-parts, anacron, and mail delivery have explicit matrix decisions.
- Migrate a schema-v7 fixture, restart, and verify stdin plus six-field source persistence.

## Verification

```sh
sh scripts/verify.sh all
```

Also run focused cron import, store migration, executor stdin, API lifecycle, and before/after recurrence benchmarks. Record exact results in `verification.md`.
