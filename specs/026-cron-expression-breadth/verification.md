# Verification: General Five-Field Cron Breadth

**Date**: 2026-08-28
**Branch**: `codex/026-cron-expression-breadth`
**Issue**: #22 (partial; remains open)

## Scope Matrix

The acceptance matrix covers:

- minute/hour lists, ranges, wildcard steps, range steps, overlaps, and uneven steps;
- day-of-month and month lists/ranges/names;
- arbitrary weekday sets, names, and Sunday `0`/`7` normalization;
- safe cross-field conjunctions with at most one restricted day field;
- creation, explanation, conversion, crontab preview/import, editing, restart, catch-up, and export;
- `skip`, `last_valid`, and `next_valid` date-set behavior;
- strict schedule anchors, leap day, short months, and both DST transitions;
- retained DOM/DOW OR, Quartz, boot, and modifier-composite refusals.

## Pre-change Baseline

Windows amd64, Intel Core Ultra 9 285K, three samples:

| Benchmark | Time | Allocations |
| --- | ---: | ---: |
| Dispatch | 33.95-35.91 us/op | 1,309 B/op, 27 allocs/op |
| NextRun daily | 39.88-41.51 us/op | 59,304 B/op, 51 allocs/op |
| NextRun nearest weekday | 43.01-46.00 us/op | 47,248 B/op, 31 allocs/op |

Command:

```text
go test ./internal/engine -run '^$' -bench . -benchmem -count=3
```

## Test-First Evidence

- Parser fixtures initially exposed the standing calendar-step refusals and the phrase-only compilation bottleneck. Direct compilation tests then drove the constrained daily RRULE shape and retained-source contract.
- Next-run matrices drove calendar-set resolution, DST wall-time normalization, strict anchor handling, and missing-date collision suppression.
- The twenty-expression occurrence matrix found a two-year search bound that exhausted a leap-day-only schedule. The final resolver covers the maximum Gregorian leap-day gap, including a non-leap century.
- API, CLI, storage, catch-up, and desktop tests were added before their shared paths were accepted as complete. Existing paths required no parallel feature implementation once the shared compiler contract was in place.

## Focused Verification

All focused suites passed:

```text
go test ./internal/cron ./internal/scheduleinput ./internal/api/server ./internal/cli ./internal/store ./internal/catchup ./internal/engine ./gui -count=1
go test ./internal/schedule -run CompositeCron -count=1
Git Bash: go test ./... -count=1
```

The full Git Bash suite includes the integration package and passed in 57.680s.

Post-change benchmark samples on the same Windows host:

| Benchmark | Time | Allocations |
| --- | ---: | ---: |
| Dispatch | 32.79-33.28 us/op | 1,308-1,309 B/op, 27 allocs/op |
| Existing NextRun | 32.34-34.14 us/op | 59,304 B/op, 51 allocs/op |
| Composite NextRun, broad | 60.14-63.00 us/op | 47,696 B/op, 32 allocs/op |
| Composite NextRun, sparse leap day | 59.78-61.22 us/op | 47,376 B/op, 35 allocs/op |
| Cron compile, simple | 30.18-30.80 us/op | 49,268-49,269 B/op, 70 allocs/op |
| Cron compile, broad | 37.16-37.68 us/op | 52,659 B/op, 114 allocs/op |
| Cron compile, sparse | 34.19-35.66 us/op | 49,391 B/op, 93 allocs/op |

The existing dispatch and next-run controls improved from baseline; neither regressed. Composite evaluation remains below the 100ms p99 dispatch budget by more than three orders of magnitude.

## Canonical Gates

| Gate | Result | Evidence |
| --- | --- | --- |
| format | PASS | Canonical Git Bash driver |
| vet | PASS | Canonical Git Bash driver |
| lint | PASS | golangci-lint v2.1.6, 0 issues |
| race | PASS | WSL Ubuntu-24.04 host split; all selected packages and integration passed |
| gui | PASS | `go test ./gui/...` |
| coverage | PASS | engine 88.1%, schedule 91.1%, timezone 88.9%, store 87.0%, catchup 87.5%, logbus 91.1% |
| docs | PASS | 11 pages, links, front matter, fences, theme, and product policy clean |
| automation | PASS | approved actions, CodeQL contract, and eight-gate manifest |

The canonical Windows run passed format, vet, and lint, then reached the known host limitation at race because `gcc` is absent. Race passed through the repository's established WSL Ubuntu-24.04 split; the remaining gates passed in Git Bash. This is an environment split, not a skipped gate.

## Final Audits

- Task completion: 31/31 complete
- Acceptance criteria: SC-001 through SC-007 satisfied
- `git diff --check`: pass
- UTF-8 without BOM: pass
- Mojibake scan: pass

## Review Remediation

Codex review on PR #61 identified two presentation gaps, both corrected before merge:

- composite date sets now include an effective `skip`, `last_valid`, or `next_valid` clause in preview and task summaries when a selected date can be absent;
- uneven wildcard minute steps render their exact field-local minute set rather than implying a constant elapsed interval across hour boundaries.
