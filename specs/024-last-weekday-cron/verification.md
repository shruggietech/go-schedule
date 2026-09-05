# Verification: Last-Weekday Cron Parity

## Baseline (2026-08-28)

- PASS: `go test ./internal/cron ./internal/cli ./internal/scheduleinput ./internal/api/server -count=1`
- PASS: `sh scripts/docs-check.sh`
- PASS: requirements checklist 16/16 and last-weekday fidelity checklist 17/17

## Test-first evidence

- RED: `go test ./internal/cron ./internal/cli ./internal/scheduleinput ./internal/api/server -count=1` failed on the new supported import, seven-day export matrix, policy, round-trip, crontab, shared-input, and API expectations while production still returned the prior `L` and last-weekday refusals.
- RED: Existing malformed-input behavior also exposed why named forms such as `WEDL` need a dedicated terminal-suffix parser instead of generic name truncation.

## Focused verification

- PASS: `go test ./internal/cron ./internal/cli ./internal/scheduleinput ./internal/api/server -count=1`
- PASS: all seven last weekdays, numeric and named Sunday aliases, canonical `0L` output, policy-inert export, selector guards, task non-mutation, and a January-through-May last-Friday round trip across the March DST transition.
- PASS: `sh scripts/docs-check.sh` (11 pages plus current-surface policy fixtures).

## Canonical verification

- PASS: all eight `scripts/verify.sh` gates completed in the foreground.
  1. format (Windows; no files reported)
  2. vet (Windows)
  3. lint (Windows; 0 issues)
  4. race (Ubuntu WSL with GCC; all selected packages passed)
  5. GUI (Windows)
  6. coverage (core packages 87.0% through 91.5%)
  7. docs (11 pages and policy fixtures)
  8. automation (approved actions, CodeQL contract, and exact eight-gate manifest)
- NOTE: The initial Windows `all` invocation reached race after format, vet, and lint passed, then stopped because this host has no Windows C compiler. Race was rerun successfully under the available Ubuntu GCC toolchain; the remaining gates were rerun individually on Windows. No gate was skipped or weakened.

## Encoding and repository audits

- PASS: `git diff --check`
- PASS: strict UTF-8 decoding without BOM across all 25 changed files
- PASS: mojibake signature audit across all 25 changed files
