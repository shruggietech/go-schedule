# Verification: Filesystem Watchers

## Spec Kit analysis

The blocking cross-artifact analysis covered 25 functional requirements, 9 measurable outcomes, 3 user stories, and 47 implementation tasks. Initial coverage was 100 percent with no requirement or story left unmapped.

One CRITICAL ordering defect was found: the generated analysis task followed implementation instead of gating it. The task was moved into Phase 1 as T005 and every later task was renumbered in dependency order before product code began. No unresolved critical, high, ambiguity, duplication, constitution, or coverage finding remains.

Both requirements-quality checklists pass: 16 of 16 general items and 33 of 33 filesystem watcher contract items.

## Implementation evidence

Schema v14 migrates a real schema-v13 fixture, preserves existing data, adds durable watcher definitions, and adds nullable watcher run provenance. Store lifecycle tests cover normalization, validation, target cascade, timeless-task activation readiness, run provenance, and a 100-definition lifecycle budget below one second.

The daemon now owns one buffered fsnotify observer and one generation-scoped event loop. Deterministic tests cover exact and glob selection, recursive depth, 100-event write coalescing, debounce, stability, obsolete-generation cancellation, missing-root recovery without replay, overflow degradation, transition deduplication, and handle closure. Native Windows tests complete 100 atomic-replacement trials with exactly one dispatch per trial, detect a direct write in an existing nested directory, and reject linked roots when the operating system permits link creation.

The shared engine dispatcher enforces task lifecycle, command readiness, group eligibility, overlap policy, and worker bounds for watcher requests. Completed runs carry `filesystem_watcher` and `source_watcher_id`; matched file paths are absent from durable runs, events, and logs.

The local API, typed client, `gosched watcher` command family, and the desktop Triggers view support create, list, show, update, enable, disable, and delete. API tests prove string-duration validation, health joins, one runtime reload, and one path-free lifecycle event. Headless GUI tests prove the separate structured watcher table and its selection, health, and readiness fields.

Focused verification passed with `go test ./...`, `go test -race ./internal/watcher -count=5`, and race-enabled watcher, engine, store, API, and event suites. Cgo-free daemon, CLI, and internal package builds passed for `linux/amd64`, `darwin/amd64`, and `windows/amd64`; the native Windows headless GUI suite passed separately because Fyne desktop cross-compilation requires platform toolchains.

The canonical `scripts/verify.sh all` run passed all eight gates through the installed WSL shell with the Windows Go binaries explicitly selected: format, vet, lint with zero findings, race, GUI in 223.752 seconds, coverage, docs, and automation. Coverage results were engine 81.9 percent, schedule 89.2 percent, timezone 91.3 percent, store 80.1 percent, catchup 88.9 percent, and logbus 91.1 percent. The publication validator reported no em dashes or hard-wrapped Markdown prose.

## Review evidence

Pending pull-request review.
