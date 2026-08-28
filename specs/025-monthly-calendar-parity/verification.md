# Verification: Monthly Calendar Cron Parity

## Baseline (2026-08-28)

- PASS: `go test ./internal/schedule ./internal/cron ./internal/scheduleinput ./internal/cli ./internal/api/server ./gui/... -count=1`
- PASS: `sh scripts/docs-check.sh`
- PASS: requirements checklist 16/16 and calendar-fidelity checklist 19/19

## Test-first evidence

- RED: `go test ./internal/schedule ./internal/cron -count=1` failed before
  production implementation. The schedule package did not compile because the
  typed adjustment did not exist; cron still returned the prior named `L`/`W`
  refusals; native selector phrases did not parse.
- RED coverage established the durable marker, canonical grammar, exact
  weekend/month-boundary dates, all three missing-date policies, invalid stored
  shapes, accepted five-field forms, surrounding refusals, export eligibility,
  and twelve-month round trips.

## Focused verification

- PASS: `go test ./internal/domain ./internal/store ./internal/schedule
  ./internal/cron ./internal/scheduleinput ./internal/cli
  ./internal/api/server ./gui/... -count=1`.
- PASS: `go test ./... -count=1` across all packages.
- PASS: schedule matrices cover every target weekday, every month-ending
  weekday, day 1 boundary reversal, days 1/15/28/29/30/31, common/leap
  February, skip/last-valid/next-valid, duplicate suppression, and DST.
- PASS: schema-v6 legacy default, marked-schedule create/read/reopen, API
  preview/create/reload, atomic invalid update, shared-input identity, CLI
  text/import, catch-up, and desktop preview/save behavior.
- PASS: twelve-month cron round trips from September 2027 through August 2028
  include fall and spring DST transitions plus leap day.
- PASS: `go test ./internal/engine -run '^$' -bench 'BenchmarkNextRun'
  -benchmem -count=3`; adjusted next-run remained approximately 40-46 us/op
  on the local Windows host, comparable to the existing 31-45 us/op monthly
  benchmark and far below the 100 ms dispatch budget.
- PASS: `bash scripts/docs-check.sh` (11 pages plus policy fixtures).
- PASS: issue #22 remains open and its A9, A10, and A11 inventory rows now
  record the implemented focused subsets and remaining composite refusals.

## Canonical verification

- PASS: all eight `scripts/verify.sh` gates completed in the foreground.
  1. format (Windows; no files reported)
  2. vet (Windows)
  3. lint (Windows; 0 issues)
  4. race (Ubuntu 24.04 WSL with GCC; all selected packages passed)
  5. GUI (Windows)
  6. coverage (core packages 87.0% through 91.8%)
  7. docs (11 pages and policy fixtures)
  8. automation (approved actions, CodeQL contract, and exact eight-gate manifest)
- NOTE: The initial Windows `all` invocation passed format, vet, and lint, then
  reached the known host limitation: no Windows GCC for `runtime/cgo`. Race
  passed under the available Ubuntu GCC toolchain; GUI, coverage, docs, and
  automation then passed on Windows. No gate was skipped or weakened.

## Encoding and repository audits

- PASS: `git diff --check`.
- PASS: strict UTF-8 decoding without BOM across all 31 changed paths.
- PASS: mojibake-signature audit across all 31 changed paths.
