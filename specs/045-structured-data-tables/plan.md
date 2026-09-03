# Implementation Plan: Structured Desktop Data Tables

**Branch**: `codex/045-structured-data-tables` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/045-structured-data-tables/spec.md`

## Summary

Replace the concatenated text rows in Tasks, Schedule List, and Activity with one shared structured-list presentation: a fixed header outside a virtualized vertical row list, responsive columns that always consume the available width, ellipsized cells with full-value disclosure, stable identity-based selection, subtle alternating row surfaces, and normalized semantic labels/glyphs. Reuse Fyne's mature `widget.List` row selection, hover, focus, virtualization, and keyboard behavior instead of `widget.Table`, whose interaction model highlights individual cells and permits horizontal scrolling. No daemon, API, persistence, dependency, or scheduling contract changes are required.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Fyne v2.8.1 and the existing go-schedule GUI/theme packages; no new dependency

**Storage**: N/A; presentation-only mapping over existing in-memory task, occurrence, log, and alert snapshots

**Testing**: Go `testing`, Fyne headless test driver, focused `go test`, focused `go test -race`, full `go test ./...`, and `sh scripts/verify.sh all`

**Target Platform**: Cross-platform desktop GUI, with required native Windows dark/light qualification evidence retained for issue closure

**Project Type**: Go desktop application backed by a local daemon/API

**Performance Goals**: Preserve virtualized row creation; mapping and refresh remain O(rows × visible columns), with 100-row headless build/refresh completing within 100 ms on the local verifier

**Constraints**: No horizontal scrolling; fixed headers; 4.5:1 normal-text and 3:1 essential-indicator contrast; stable identity across refresh; no new API data; no shell/process behavior changes; existing selection and activation semantics preserved

**Scale/Scope**: Three desktop views, five task columns, four Schedule columns, four Activity columns, at least 100 rows per view, dark/light/follow-system appearances, and all five supported interface font families

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

| Principle / constraint | Evidence before research | Post-design result |
| --- | --- | --- |
| I. Code Quality | Shared row models and layout avoid three divergent string renderers; no ignored errors or new concurrency | PASS: bounded pure mapping helpers and one reusable presentation component; no new goroutines or recoverable error paths |
| II. Testing Standards | Every behavior is assigned red-first mapping, layout, selection, activation, and contrast tests | PASS: the design includes focused headless tests, race coverage where supported, full suite, and all canonical gates |
| III. UX Consistency | Explicit headers, consistent semantic roles, plain-language fallbacks, and stable interaction are core requirements | PASS: one contract governs labels, fallbacks, glyphs, column sizing, and full-value disclosure across all views |
| IV. Performance | Existing virtualized lists are retained; row work remains bounded by visible cells | PASS: no eager widget tree per full dataset and a 100-row mapping/refresh budget is part of verification |
| V. Autonomous Build-Phase Execution | S045 is issue-backed (#112, #113), uses a review branch, runs every spec-kit phase, and will publish only under explicit user authorization | PASS: user explicitly authorized push/PR and up to two review rounds; final merge remains with the maintainer |
| Go/toolchain and dependency constraint | Existing Go/Fyne stack is sufficient | PASS: no module or toolchain change |
| Supported platforms | Pure-Go GUI widget construction remains headless-testable and platform-neutral | PASS: no platform-specific implementation; Windows qualification steps documented |
| Backward compatibility | No persisted or external schema changes | PASS: all source data and backend contracts are unchanged |
| Security | Presentation only; no secrets, execution, transport, or authorization changes | PASS: full values remain local and are exposed only through existing view interactions |

### Workflow Tooling Deviation

The installed checklist prerequisite script requires `plan.md`, contradicting the required `specify → clarify → checklist → plan` order. S045 preserved the governing order by using the feature directory already resolved by the successful clarification prerequisite command and generated `checklists/ux.md` before running `setup-plan.ps1`. The product's workflow tooling is not modified because that would be unrelated scope.

## Project Structure

### Documentation (this feature)

```text
specs/045-structured-data-tables/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── structured-table.md
├── checklists/
│   ├── requirements.md
│   └── ux.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
gui/
├── structured_table.go       # shared row model, responsive columns, row widget, and table shell
├── structured_table_test.go  # shared layout, semantics, contrast, truncation, and scale tests
├── tasks.go                  # Task row mapping and retained task operations
├── app_test.go               # Task identity, refresh, and activation regressions
├── schedule.go               # Schedule occurrence mapping and list/detail integration
├── calendar_test.go          # Schedule row normalization and control regressions
├── logs.go                   # Activity row mapping and detail integration
├── logs_test.go              # Activity normalization, filtering, and activation regressions
├── theme.go                  # subtle alternating-row theme role if required
└── theme_test.go             # row and semantic contrast checks

docs/
└── gui-fields.md             # visible column meanings, labels, fallbacks, and semantic conventions

test/windows/
└── README.md                 # attended S045 populated dark/light qualification steps
```

**Structure Decision**: Keep the feature inside the existing `gui` package. Add one shared structured-table presentation primitive, then retain source-specific mapping and behavior in `tasks.go`, `schedule.go`, and `logs.go`. This keeps table mechanics consistent without moving domain or API responsibilities into UI code.

## Phase 0: Research Decisions

Research findings and rejected alternatives are recorded in [research.md](research.md). The decisive choices are:

1. Use a fixed header plus virtualized `widget.List` rows, not `widget.Table`, to preserve full-row hover/selection and eliminate horizontal scrolling.
2. Allocate every row and header through the same bounded weighted layout, with ellipsis and a full-value disclosure strip rather than horizontal overflow.
3. Normalize domain values into pure row models before binding widgets so unknowns, casing, glyphs, semantic importance, and stable identity are testable without rendering.
4. Reuse S044 theme roles and Fyne label importance for semantic colors, with a single quiet alternate-row role and explicit contrast tests.
5. Preserve existing per-view order and controls; do not add sorting, column customization, or backend fields.

## Phase 1: Design

### Shared presentation

- `structuredColumn` defines header, minimum share, expansion weight, alignment, and whether the cell is primary.
- `structuredCell` carries display text, complete text, and semantic importance.
- `structuredRowModel` carries stable identity, ordered cells, and a complete accessible summary.
- `structuredRow` is a reusable tappable/double-tappable widget containing an alternating background and truncated labels.
- `structuredList` composes a fixed header, a `widget.List`, and an optional selectable full-value disclosure label. It owns index-to-identity reconciliation while view code owns source data.
- `responsiveColumnWidths` distributes the exact available width with bounded minimum shares and weights; it never reports content wider than its container.

### View mappings

- Tasks: `Task | Enabled | Lifecycle | Time zone | Group`. Values are `Enabled`/`Disabled`, `Active`/`Completed`/`Disabled`/normalized unknown, and `Not assigned`/`Unknown` fallbacks.
- Schedule List: `When | Task | Event | Outcome`. Upcoming rows use `SCHEDULED`; past rows use `COMPLETED` plus normalized outcome. Outcome glyph and semantic role encode success, failure, skipped, caught up, queued, unavailable, or neutral unknown.
- Activity: `When | Severity | Source | Summary`. Severity maps to `INFO`, `WARNING`, `ERROR`, or normalized neutral `UNKNOWN`; the leading glyph and label share one semantic importance.

### Interaction and disclosure

- Task row selection continues to drive existing toolbar actions; double activation edits by stable task ID.
- Schedule row selection displays a complete, read-only occurrence summary below the list.
- Activity row selection opens the existing complete detail dialog and then clears visual selection, preserving current behavior.
- Retained Task selection is restored by stable ID after refresh. Removed IDs clear. Activity and Schedule use model identity only for activation/disclosure, not backend mutation.

### Verification strategy

- Red-first pure mapping tests cover every known and unknown lifecycle/outcome/severity value.
- Shared tests cover header/cell count, width conservation, narrow widths, ellipsis, semantic importance, alternating rows, stable identity, empty and 100-row datasets.
- Integration tests cover task selection after reorder/removal, task double activation, Schedule controls/order, Activity filters/clear/detail, and theme/font refresh.
- Contrast tests blend the alternate, hover, selection, and focus surfaces against both palette backgrounds and enforce the spec thresholds.
- The native Windows checklist records populated Tasks, Schedule List, and Activity in dark and light modes without claiming unattended visual acceptance.

## Complexity Tracking

No constitution violations require justification. The shared presentation primitive is the smallest coherent implementation because the two approved issues explicitly require coordinated table behavior across three views; three independent renderers would duplicate responsive, semantic, and selection logic.
