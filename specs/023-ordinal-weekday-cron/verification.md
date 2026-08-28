# Verification: Ordinal-Weekday Cron Parity

## Baseline (2026-08-28)

- PASS: `go test ./internal/cron ./internal/cli ./internal/scheduleinput ./internal/api/server -count=1`
- PASS: `sh scripts/docs-check.sh`
- PASS: requirements checklist 16/16 and cron-fidelity checklist 18/18

## Test-first evidence

- RED: `go test ./internal/cron -count=1` failed on the new supported import,
  35-case export matrix, policy, and round-trip expectations while production
  still returned the prior `#` and ordinal-weekday refusals.
- RED: `go test ./internal/scheduleinput ./internal/cli ./internal/api/server
  -count=1` failed at the shared input, crontab import, conversion, preview, and
  create boundaries while `#` remained refused.

## Focused verification

- PASS: `go test ./internal/cron ./internal/scheduleinput ./internal/cli
  ./internal/api/server -count=1`
- PASS: all 35 weekday/ordinal exports, Sunday aliases, malformed/refused
  shapes, effective policy matrix, task non-mutation, and January-through-May
  fifth-Friday round trip across the March DST transition.
- PASS: `sh scripts/docs-check.sh` (11 pages plus current-surface policy fixtures).
- PASS: final Spec-Kit analysis found 100% task coverage for all 17 functional
  requirements and 6 success criteria, zero unmapped tasks, zero ambiguity or
  duplication findings, and zero constitutional conflicts.

## Canonical verification

- PASS: `sh scripts/verify.sh all`
  1. format (no files reported)
  2. vet
  3. lint (0 issues)
  4. race
  5. GUI
  6. coverage (core packages 87.0% through 91.5%)
  7. docs (11 pages and policy fixtures)
  8. automation (approved actions and exact eight-gate manifest)

## Encoding and repository audits

- PASS: `git diff --cached --check`
- PASS: strict UTF-8 decoding without BOM across all 25 changed files
- PASS: mojibake signature audit across all 25 changed files
