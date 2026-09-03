# Verification: Structured Desktop Data Tables

## Baseline (2026-09-03)

- Repository: `github.com/shruggietech/go-schedule`
- Branch: `codex/045-structured-data-tables`
- Base: `66cf77582bf6f9654971fa6a6464c84ac8c246c1` (merged S044)
- Issues: [#112](https://github.com/shruggietech/go-schedule/issues/112), [#113](https://github.com/shruggietech/go-schedule/issues/113)
- Worktree began clean and synchronized with `origin/main`.
- Existing `go test ./gui -count=1`: PASS (`ok`, 52.065s).
- Existing presentation evidence: `gui/tasks.go`, `gui/schedule.go`, and `gui/logs.go` each concatenate multiple unlabeled fields into a single `widget.Label`; Tasks uses square brackets around lifecycle, Schedule uses unlabeled marker/time/name text, and Activity emits lowercase `info` beside uppercase `ERROR`.
- `.gitignore` already covers Go binaries/tests/coverage, generated Windows resources, local runtime data, module/workspace files, editors, operating-system cruft, and local distribution output. No ignore change is necessary.

## Spec-kit Analysis Gate

- Functional requirements: 28
- Measurable outcomes: 8
- Tasks: 44
- Requirement/outcome task coverage: 36/36 (100%)
- Critical findings: 0
- High findings: 0
- Unmapped tasks: 0
- Checklists: requirements 16/16, UX 27/27

## Red/Green Evidence

- Shared foundation red: `go test ./gui -run 'Test(Responsive|Structured)' -count=1`
  failed to compile because the new structured column, row, and list contracts
  did not yet exist.
- Shared foundation green: responsive/layout/row/list tests passed in 0.091s;
  dark/light interaction and semantic contrast tests passed in 0.053s.
- Tasks checkpoint: `go test ./gui -run 'Test(Task|TasksTab|Structured)' -count=1`
  passed in 0.245s.
- Schedule checkpoint: the focused column, row mapping, table, ordering,
  occurrence, and construction selection passed in 0.171s.
- Activity red: the focused selection failed to compile on missing
  `activityColumns`, source identity, and `activityRowModels` before the view
  integration existed.
- Activity checkpoint: `go test ./gui -run
  'Test(Activity|MergeLogEntries|UI_Activity)' -count=1` passed in 0.197s.
- Combined S045 focused selection passed in 0.352s.

## Verification Gates

- `gofmt` on all changed Go files: PASS.
- `go test ./gui -count=1`: PASS (`ok`, 72.043s before the final Unicode
  normalization test; final canonical GUI rerun passed in 79.986s afterward).
- Deterministic S045 race selection covering layout, mapping, normalization,
  identity, virtualization, summary, and contrast: PASS (`ok`, 1.338s).
- A broader headless GUI race diagnostic reproduced Fyne 2.8.1's font-cache
  race when the test driver rendered while the Schedule loader used
  `fyne.Do`. The trace is confined to `fyne.io/fyne/v2/internal/cache`; the
  normal full GUI suite and the deterministic S045 race selection pass. GUI is
  intentionally excluded from the repository's canonical race package set for
  this upstream driver limitation.
- `go test ./... -count=1`: PASS, including GUI in 73.065s and integration in
  47.367s.
- `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all`: PASS through all
  eight gates: format, vet, lint (0 issues), race, GUI, coverage, docs, and
  automation. Coverage thresholds passed (engine 86.4%, schedule 89.2%,
  timezone 91.3%, store 84.1%, catchup 88.9%, logbus 91.1%).
- `git diff --check`: PASS.
- UTF-8 strict decode, no-BOM validation, and mojibake scan: PASS for all 31
  changed or added files.
- Final `scripts/spec-lifecycle-check.sh .`: PASS (45 lifecycle-consistent
  specifications).
- Final `scripts/docs-check.sh`: PASS (15 pages, links, front matter, fences,
  theme, and product policy clean).

## Requirements Audit

| Requirement / outcome | Delivery evidence |
| --- | --- |
| FR-001..FR-006, SC-002..SC-003 | `structuredList` owns a fixed header, virtualized vertical body, proportional no-overflow allocator, ellipsis, complete-value disclosure, and empty body; focused tests cover exact width conservation and 100 rows. |
| FR-007..FR-010, FR-025..FR-027, SC-004..SC-005 | Pure view mappers document missing/unknown values, pair semantic glyphs with text, and reuse S044 importance colors plus a translucent alternate-row surface. Dark/light rest, alternate, hover, focus, selection, and semantic contrast math passes. Native rasterization remains below. |
| FR-011, SC-006 | `structuredList.setRows` reconciles selection by stable identity and clears removed records; tests cover reorder, retention, removal, duplicates, stale activation, and Unicode. |
| FR-012..FR-016, SC-001 | Tasks exposes Task, Enabled, Lifecycle, Time zone, and Group as separate cells. Pure and integration tests cover lifecycle values, missing values, toolbar selection, refresh, and double-click edit by current task ID. |
| FR-017..FR-020 | Schedule exposes When, Task, Event, and Outcome; every known outcome plus missing and unknown values has a normalized mapping. Existing range/Calendar behavior remains intact and chronological sorting passes. |
| FR-021..FR-024 | Activity exposes When, Severity, Source, and Summary; INFO, WARNING, ERROR, empty, and unknown severity mappings pass. Existing newest-first filtering, Clear View acknowledgement, live refresh, and exact identity-to-detail resolution remain covered. |
| FR-028, SC-007..SC-008 | Focused, full GUI, race, repository, and canonical gates pass. The exact native Windows dark/light matrix is documented in `test/windows/README.md` and remains a release-candidate gate. |

All 28 functional requirements and all 8 measurable outcomes have repository
delivery evidence. No daemon, API, persistence, command-entry (#110), or richer
failure-diagnostic (#102) scope was introduced.

## Native Windows Qualification Boundary

S045 supplies automated headless evidence and attended test instructions. Populated dark/light screenshots from the exact release-candidate MSI remain required before issues #112 or #113 close.
