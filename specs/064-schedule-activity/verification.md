# Verification: Operational Schedule and Activity

**Date**: 2026-09-07

**Branch**: `codex/064-schedule-activity`

**Issue**: [#155](https://github.com/shruggietech/go-schedule/issues/155)

## Spec-kit analysis

The specification, clarifications, checklists, plan, data model, bridge contract, and tasks were compared before implementation. Every functional requirement maps to acceptance scenarios and executable tasks, terminology is consistent, no clarification markers remain, no constitution conflicts were found, and the analysis gate completed with zero findings.

## Focused evidence

| Command | Result |
| --- | --- |
| `go test ./...` from `desktop/` | PASS, operations facade and service mapping included |
| `go test -race ./operations ./...` from `desktop/` | PASS, all desktop packages race-clean |
| `npm test -- --run` from `desktop/frontend/` | PASS, 15 files and 45 tests |
| `npm run build` from `desktop/frontend/` | PASS, TypeScript and production Vite build |
| `npx playwright test e2e/operations.spec.ts` from `desktop/frontend/` | PASS, Chromium keyboard, accessibility, zoom, and 105-row Schedule plus 107-row Activity interaction |
| `go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean` from `desktop/` | PASS, native Windows amd64 `go-schedule.exe` built with generated bindings and packaged assets |
| `go run ./scripts/github-format` | PASS, no em dashes or hard-wrapped Markdown prose |

## Canonical verification

`bash -lc 'GO="/mnt/c/Program Files/Go/bin/go.exe" GOFMT="/mnt/c/Program Files/Go/bin/gofmt.exe" sh scripts/verify.sh all'` completed successfully in the foreground.

| Gate | Result |
| --- | --- |
| format | PASS |
| vet | PASS |
| lint | PASS, 0 issues |
| race | PASS |
| gui | PASS |
| coverage | PASS, engine 82.6%, schedule 89.2%, timezone 91.3%, store 80.2%, catchup 88.9%, logbus 91.1% |
| docs | PASS, 15 pages plus policy fixtures |
| automation | PASS, workflow and fixture checks |

## Acceptance mapping

- Schedule retains agenda and accessible calendar projections, 1-day, 7-day, and 30-day windows, stable identity selection, and explicit Prediction versus Recorded run labeling.
- Activity snapshots active executions before persisted history for a race-free completion handoff, then retains recent daemon logs, bounded alerts, exact log-path metadata, text and typed filters, detailed output and provenance, individual acknowledgement, and non-destructive Clear View behavior.
- Running, success, failure, skipped, caught-up, queued, upcoming, acknowledged, unacknowledged, and unavailable states use text labels and styled shapes rather than color alone.
- Request sequence tests reject stale responses, view state remains component-owned across refresh, and failed loads retain the last complete snapshot.
- Automated React and Chromium fixtures exceed the required 100 rows for both Schedule and Activity while preserving keyboard focus and responsive reflow.
- Additive read-only active-run and bounded-alert API behavior was added without changing persistence schema, scheduling policy, retention, or packaging inputs.

## Encoding sanity

The changed specification, Go, TypeScript, CSS, test, and GitHub-facing files pass the repository formatter and targeted mojibake scan as UTF-8 without BOM.
