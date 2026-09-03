# Quickstart: Validate Structured Desktop Data Tables

## Prerequisites

- Go toolchain selected from `go.mod`
- C compiler available for the race detector where required
- Repository checkout on `codex/045-structured-data-tables`
- For attended evidence, Windows with the current exact-candidate installer

## 1. Focused headless verification

```bash
go test ./gui -run 'Test(Structured|TaskRow|TaskList|ScheduleRow|ScheduleList|ActivityRow|ActivityTable|Semantic)' -count=1
go test -race ./gui -run 'Test(Structured|TaskRow|TaskList|ScheduleRow|ScheduleList|ActivityRow|ActivityTable|Semantic)' -count=1
```

Expected: mapping, normalization, row identity, responsive widths, interaction, and contrast tests pass without a display server.

## 2. Full repository verification

```bash
go test ./... -count=1
sh scripts/verify.sh all
```

Expected: the full suite passes and all eight canonical gates report success in order: format, vet, lint, race, gui, coverage, docs, automation.

## 3. Tasks walkthrough

Populate at least 100 tasks including:

- enabled active, disabled, and completed lifecycle combinations;
- grouped, ungrouped, missing-group, and long Unicode names;
- distinct time zones and one missing/unknown display value in test data.

Confirm the contract in [contracts/structured-table.md](contracts/structured-table.md): fixed headers, plain-language state columns, quiet alternating rows, full-value disclosure, no horizontal scrollbar, keyboard and pointer selection, toolbar targeting, double-click editing, selection after reorder, and safe clearing after removal.

## 4. Schedule walkthrough

In List mode, provide upcoming and past occurrences covering success, failure, skipped, caught-up, queued, missing, and unknown outcomes. Confirm fixed headers, chronological ordering, normalized event/outcome text, paired glyph/color semantics, full-value disclosure, range changes, live refresh, and Calendar round trips.

## 5. Activity walkthrough

Provide log records and alerts covering info, warning, error, empty, and unknown severity plus long sources/messages. Confirm fixed headers, newest-first ordering, `INFO`/`WARNING`/`ERROR` casing, paired glyph/color semantics, severity filters, Clear View acknowledgement, live refresh, and correct complete detail activation.

## 6. Appearance and sizing matrix

Repeat populated view checks in:

- Dark mode with System font
- Light mode with System font
- Follow system in both available operating-system variants
- Inter, Ubuntu, Geist, and Monospace fonts
- Default launch size and the smallest size the application permits

Expected: no horizontal scrolling, overlap, clipped primary meaning, unreadable interaction state, or header/body misalignment.

## 7. Native Windows evidence

Follow the S045 section in `test/windows/README.md`. Record exact installer identity and populated dark/light screenshots before closing #112 or #113. Headless results alone do not satisfy that attended acceptance.
