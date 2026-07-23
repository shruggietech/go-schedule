# Phase 1 Quickstart: Validating Feature 006

**Feature**: 006-maintainer-test-scripts | **Date**: 2026-07-23

How to prove this feature works, end to end. Each scenario maps to a spec acceptance
criterion. Reference material lives in [contracts/cli.md](contracts/cli.md) and
[data-model.md](data-model.md) rather than being repeated here.

## Prerequisites

```bash
sqlite3 --version
```

Must report 3.33.0 or later. If absent, either install it from the platform package manager
or let the scripts do it:

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -InstallSqlite
```

Also needed: `pwsh` 7+ for the PowerShell twins, `bash` for the POSIX twins, and a running
`goschedd` for anything involving the scheduler.

---

## Scenario 1 — One beat, no scheduler (SC-004, US1 scenario 3)

The cheapest possible confidence check.

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -Label smoke
```

Expect: exit 0, one line on stdout naming the resolved database path.

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query summary
```

Expect: exactly 1 record, 1 session. Then run the POSIX twin and confirm the count becomes 2
with both `script_flavor` values represented — that is the twin-parity check (SC-004).

---

## Scenario 2 — On-time firing (SC-001, SC-002, US1)

```bash
gosched task add "hb-verify" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -IntervalSeconds --arg 60 --schedule "every 1 minute"
```

Wait ten minutes, then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query drift -IntervalSeconds 60
```

**Expect**: about 10 beats; `expected_source` reported as `boundary` for all of them; a p99
drift figure with its source labelled; and nothing flagged unreliable. Then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query gaps -IntervalSeconds 60
```

**Expect**: no gaps.

Read the p99 against the project's documented dispatch budget. This measures and compares; it
does not certify (SC-002).

---

## Scenario 3 — Restart survival (US2 scenario 1)

With scenario 2's task still registered:

```bash
gosched service restart
```

Wait two more minutes, then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query restarts
```

**Expect**: beats on both sides of the restart. Note that each *invocation* is its own
session, so under a scheduler the session boundary marks every beat — the meaningful evidence
is the uninterrupted `started_ms` sequence across the restart, which is what the query
reports.

---

## Scenario 4 — Downtime catch-up (US2 scenario 2)

```bash
gosched service stop
```

Wait past at least one scheduled firing, then:

```bash
gosched service start
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query gaps -IntervalSeconds 60
```

**Expect**: one gap covering the downtime, and — under a make-up-once catch-up policy —
exactly one beat shortly after restart before the cadence resumes. More than one make-up beat,
or none, is a genuine finding about the scheduler.

---

## Scenario 5 — Overlap policies (US2 scenario 3)

For each of `queue_one`, `skip`, `allow_concurrent`, with a run deliberately longer than its
interval:

```bash
gosched task add "hb-overlap" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -SleepSeconds --arg 90 --schedule "every 1 minute" --overlap queue_one
```

After several minutes:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query overlaps
```

**Expect**, per policy: `queue_one` → no overlapping ranges, runs serialized; `skip` → no
overlaps and visibly fewer beats than elapsed intervals; `allow_concurrent` → overlapping
ranges present. Delete and recreate the task between policies.

This is also the concurrent-writer check (SC-007): under `allow_concurrent` every beat must be
present despite simultaneous writers.

---

## Scenario 6 — Failure reporting (US1 scenario 4)

```bash
gosched task add "hb-fail" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -FailWith --arg 1 --schedule "every 1 minute"
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query failures
```

**Expect**: beats recorded with `outcome = failed` and `exit_code = 1`, *and* corresponding
alerts from `gosched alerts`. The beat existing proves the record survived the failure; the
alert proves the daemon saw it.

---

## Scenario 7 — Host snapshot (US3)

```bash
pwsh -File test/scripts/Test-GetSystemInfo.ps1 -InvocationSource quickstart
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database system -Query listeners
```

**Expect**: one snapshot with addresses and listening ports attached. `process_name` is
commonly `NULL` without elevation — that is normal (see data-model). Run the POSIX twin too
and confirm equivalent results.

---

## Scenario 8 — Missing prerequisite (SC-005)

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -SqliteExe /nonexistent/sqlite3
```

**Expect**: exit **2**, with a message naming both `-InstallSqlite` and the platform package
manager command. Exit 1 here would be a contract violation.

---

## Scenario 9 — Repository hygiene (SC-006)

```bash
git status --porcelain .claude
```

**Expect**: skills only. No settings or credential file. Then clone to a temporary directory
and confirm `.claude/skills/` arrives populated with no post-clone step.

---

## Automated coverage

```bash
go test ./test/integration/ -run TestScripts -v
```

Covers scenarios 1, 6, 8, the bounded-loop guarantee, and concurrent writers. It **skips with
a stated reason** when `sqlite3`, `pwsh`, or `bash` is unavailable — a skip is not a pass, and
the reason is printed so nobody reads one as the other.
