# Research: Cron Parity Closure

## Decision 1: Preserve cron's execution model at import

**Decision**: Scan crontabs in order with an explicit context containing schedule timezone, environment, and shell. Snapshot that context onto each job. Execute imported commands as `SHELL -c <original command>` and parse cron percent input before the shell sees the command.

**Rationale**: Cron applies assignments to following lines, runs the entire command through a shell, and transforms unescaped percent signs before shell evaluation. The current whitespace split changes pipelines, quoting, redirection, and standard input silently.

**Alternatives rejected**:

- Continue requiring explicit shell wrappers: keeps documented behavior but does not preserve imported cron.
- Parse shell syntax into argv: shell syntax is not safely reducible to argv and would require a second shell implementation.
- Treat `TZ` as schedule timezone: contradicts Cronie's separate `CRON_TZ` schedule directive and process `TZ` environment.

**Primary reference**: [crontab(5)](https://man7.org/linux/man-pages/man5/crontab.5.html).

## Decision 2: Add persisted task standard input

**Decision**: Add optional task `stdin`, store it in additive schema v8, expose it on create and pointer-based update, and attach an exact reader to the child process.

**Rationale**: Cron percent semantics cannot be represented by command arguments or environment. A pointer in PATCH distinguishes omitted input from an explicit clear.

**Alternatives rejected**:

- Temporary files: adds lifecycle, permission, and cleanup risks.
- Shell here-doc rewriting: changes quoting and can collide with payload content.
- Import-time refusal: leaves a core audited behavior unresolved despite a small native model addition.

## Decision 3: Make file layout explicit

**Decision**: `cron import` defaults to Unix five-field user-crontab layout and accepts explicit `--dialect quartz` and `--system` options. The scanner consumes exactly the selected number of timing fields and optional user column.

**Rationale**: A Unix command may begin with a number or cron-shaped token. Automatic six-field or username detection can therefore reinterpret valid commands.

**Alternatives rejected**:

- Detect six fields from tokens: ambiguous for valid command text.
- Treat every six-token prefix as Quartz: breaks existing user crontabs.
- Separate commands for every layout: multiplies surface area without improving fidelity.

## Decision 4: Support a bounded Quartz six-field subset

**Decision**: Standalone six-field input is the Quartz-style subset `second minute hour day-of-month month day-of-week`. It supports field sets, value-start steps, exactly placed `?`, and Quartz numeric weekday values. Five-field input retains Unix weekday values and implies second zero.

**Rationale**: RRULE already carries `BYSECOND`, so seconds need no new scheduler. Dialect-specific weekday mapping prevents one-day drift. Quartz `x/n` steps select from `x` through the field maximum, which the existing single-value step parser does not yet model.

**Alternatives rejected**:

- Apply Unix weekday numbers to six fields: contradicts the named Quartz contract in issue #22.
- Add Quartz's optional year and all extensions: expands beyond the issue's seconds and `?` gap and cannot map every extension faithfully.
- Add a cron runtime evaluator: violates the single authoritative scheduler invariant.

**Primary reference**: [Quartz CronTrigger tutorial](https://www.quartz-scheduler.org/documentation/quartz-2.5.x/tutorials/crontrigger.html).

## Decision 5: Preserve five-field output and refuse lossy execution export

**Decision**: Export a five-field timing expression when seconds are exactly zero and a six-field Quartz expression when seconds are required. Refuse tasks carrying environment, run-as, or standard input until a full operational export format can serialize them.

**Rationale**: Existing output remains stable. A plain cron line cannot carry all operational context without surrounding directives or system-layout selection, so emitting one would claim a round trip that does not exist.

**Alternatives rejected**:

- Ignore execution context: silently changes behavior.
- Always emit six fields: breaks conventional crontab consumers and existing snapshots.
- Build a multi-file export bundle: disproportionate to the import-focused issue.

## Decision 6: Keep incompatible rows explicit

**Decision**: Defer `@reboot` to event triggers and mail delivery to notifications. Ratify DOM/DOW OR semantics and unsupported modifier composites outside the single-recurrence model. Treat run-parts and anacron as separate import formats, not crontab expressions.

**Rationale**: Issue #22 accepts explicit disposition with rationale. Pretending these are line-conversion features would either approximate behavior or broaden this slice into multiple subsystems.

## Baseline

- Focused cron, CLI, store, executor, and API suites pass before S028.
- `BenchmarkNextRun` median: 32,745 ns/op, 59,376 B/op, 53 allocs/op.
- Composite broad median: 1,545,428 ns/op; composite sparse median: 107,088 ns/op.
- No new dependency, permission, service, pinned artifact, or ignore rule is required.
