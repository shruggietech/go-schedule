# CLI Contract: `gosched cron` and the missing-date flag

**Feature**: `008-cron-interop` | **Date**: 2026-07-23

Conventions inherited from `internal/cli`: results on stdout, diagnostics and
errors on stderr, exit 0 success / 1 runtime error / 2 usage or validation, and
`--json` honored on every subcommand.

## `gosched cron explain <expression>`

Translates one expression. Creates nothing, contacts the daemon only to resolve
the default timezone.

```text
$ gosched cron explain "0 9 * * 1-5"
0 9 * * 1-5
  phrase: weekdays at 09:00
  next:   2026-07-24 09:00:00 -0400 EDT
          2026-07-27 09:00:00 -0400 EDT
          2026-07-28 09:00:00 -0400 EDT
```

| Flag | Meaning |
| --- | --- |
| `--timezone` | IANA zone the run times are shown in; defaults to the daemon's task default |
| `--count` | how many upcoming runs to show (default 3) |

Exit codes: 0 on a successful translation **and** on a well-formed expression
that is declined (the refusal is the answer); 2 on a malformed expression, with
the offending field named.

```text
$ gosched cron explain "@reboot"
@reboot
  unsupported: @reboot has no equivalent — it fires at boot, not on a schedule
```

## `gosched cron import`

```text
$ gosched cron import --file crontab --dry-run
```

| Flag | Meaning |
| --- | --- |
| `--file` | crontab path, or `-` for standard input (required) |
| `--dry-run` | produce the identical report, create nothing |
| `--timezone` | IANA zone for the created tasks; defaults to the task default |
| `--group` | group ID to place the imported tasks in |
| `--count` | upcoming runs shown per line (default 3) |

Per line the report shows the source expression, the phrase, the resolved
command, and the next runs — or the refusal and its reason. Warnings
(`MAILTO=`, variable assignments) are attached to the line that produced them.

The summary is mandatory and states, in this order: the counts (read, created,
skipped, declined, failed); that cron carries no timezone and which one was
applied; and that the imported tasks received catch-up, overlap, and
missing-date defaults that cron has no notion of, naming them.

Exit codes: 0 whenever the input was read, whatever the mix of outcomes; 1 if a
task creation failed (already-created tasks remain, and the summary says so); 2
for an unreadable file or an invalid timezone.

## `gosched cron export`

```text
$ gosched cron export
# gosched cron export — 2026-07-23T15:04:05Z — 4 tasks
0 9 * * 1-5 /usr/bin/report --daily
# declined: "nightly backup" — one-off schedules have no cron equivalent
# declined: "health probe" — task is disabled; cron has no disabled state
*/15 * * * * /usr/local/bin/probe
```

| Flag | Meaning |
| --- | --- |
| `--task` | export a single task by ID |

Every task appears exactly once, as a line or as a `# declined:` comment naming
the task and the reason. An empty task set produces the header comment and
nothing else. Exit code 0 unless the daemon call fails.

## `gosched task add` / `gosched task edit` — new flag

```text
--missing-date skip|last_valid|next_valid
```

Default `skip` on create; omitted on edit means unchanged. An unrecognized value
is exit 2 naming the flag. The value is shown by `gosched task show` alongside
the overlap and catch-up policies, and included in `--json` output as
`missing_date_policy`.

The values are spelled with underscores because the existing policy flags are:
`--overlap queue_one|skip|allow_concurrent` and `--catchup one|none`. A hyphenated
`last-valid` would read better in isolation and would make this the only policy
flag in the CLI whose values are punctuated differently from every other —
constitution principle III (consistent interfaces) settles it against the
prettier spelling.
