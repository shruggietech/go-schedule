# Verification: Cron Parity Closure

## Baseline (2026-08-29)

- Focused cron, CLI, store, executor, and API server suites: PASS.
- `BenchmarkNextRun` five-run median: 32,745 ns/op, 59,376 B/op, 53 allocs/op.
- Composite broad five-run median: 1,545,428 ns/op, 54,440 B/op, 58 allocs/op.
- Composite sparse five-run median: 107,088 ns/op, 49,280 B/op, 53 allocs/op.
- No dependency, service, permission, pinned artifact, or ignore-file change is required.

## Specification analysis

- 32 functional requirements and 8 measurable outcomes map to 39 executable tasks.
- Coverage: 100 percent; unmapped requirements: 0; unmapped tasks: 0.
- Ambiguities: 0; duplications: 0; constitution conflicts: 0; critical findings: 0.
- Requirements checklist: 15/15 pass. Cron fidelity checklist: 18/18 pass.
- Deliberate deviations are explicit: imported commands replace whitespace splitting with shell execution; `TZ` remains process environment rather than schedule context; operationally lossy export becomes a refusal.

## Focused verification

- `go test ./internal/cron ./internal/cli ./internal/executor ./internal/api/server ./internal/store ./internal/schedule ./gui`: PASS.
- Explicit store-reopen and catch-up tests retain non-zero and multiple seconds: PASS.
- Quickstart conversion of `*/30 * * * * *`: PASS with a seconds-aware description.
- User-crontab dry run with `CRON_TZ`, `TZ`, `SHELL`, quoted assignment, percent stdin, and `--run-as`: PASS.
- Quartz dry run of `0 0 12 ? * MON echo noon`: PASS.
- Post-implementation benchmark medians: representative `NextRun` 32,068 ns/op (2.1 percent faster than baseline), broad composite 1,620,658 ns/op (4.9 percent slower), sparse composite 115,848 ns/op (8.2 percent slower), and seconds composite 3,470,540 ns/op. Every retained case stays within the 10 percent regression boundary and the seconds case remains far below the 100 ms budget.
- Explicit deviations: imported commands replace whitespace splitting with shell execution; `TZ` remains process environment while `CRON_TZ` owns schedule context; `LOGNAME` cannot be overridden; user crontabs require explicit `--run-as`; and operationally lossy export becomes a refusal.

## Canonical gates

- `scripts/verify.sh all`: PASS from Git Bash in one foreground run.
- Format: PASS.
- Vet: PASS.
- Lint: PASS (0 issues).
- Race: PASS across selected packages and POSIX integration tests.
- GUI: PASS.
- Coverage: PASS (engine 88.1%, schedule 88.8%, timezone 91.3%, store 86.9%, catch-up 88.9%, logbus 91.1%).
- Docs: PASS (11 pages, links, front matter, fences, theme, and product policy clean).
- Automation: PASS (approved actions, CodeQL contract, and eight-gate manifest).

## Acceptance and integrity audit

- All 39 tasks are complete; all 32 functional requirements and 8 measurable outcomes remain traceable.
- Requirements checklist: 15/15 checked. Cron fidelity checklist: 18/18 checked.
- Issue #22 matrix: all A1-A12 and B1-B9 rows have explicit dispositions.
- `git diff --check`: PASS.
- Strict UTF-8 without BOM: PASS across all 54 changed files.
- Mojibake scan: PASS across all 54 changed files.
- `.specify/feature.json` parses successfully; no unresolved placeholders remain.
- Issue closure intent: the eventual pull request must use `Closes #22`.
- Delivery boundary: commit locally on `codex/028-cron-parity-closure`, then halt before push or pull-request publication.

## Third-party review remediation (2026-08-30)

- Composite seconds evaluation now sorts and deduplicates stored field sets, starts at the relevant wall-clock cursor, and limits exhaustive comparison to the repeated fall-overlap interval. The representative seconds benchmark improved from roughly 3.3 ms and 5,789 allocations to 40-47 µs and 36 allocations per operation. Exact every-second, unordered-set, and both-fold regressions pass.
- Unix run-as environment normalization now removes every stale `USER` and `LOGNAME` entry. It also replaces inherited daemon `HOME` with the selected account's home unless the task explicitly supplied `HOME`.
- Crontab import now carries its default `SHELL=/bin/sh` into each task environment as well as using it for command execution.
- Post-remediation `scripts/verify.sh all`: PASS. Coverage: engine 88.1%, schedule 89.0%, timezone 91.3%, store 86.9%, catch-up 88.9%, logbus 91.1%.
- Unix-only executor tests compile successfully for `GOOS=linux` with `CGO_ENABLED=0`.
