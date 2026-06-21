# Implementation Plan: GUI Editor Refinements

**Branch**: `003-gui-editor-refinements` | **Date**: 2026-06-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/003-gui-editor-refinements/spec.md`

## Summary

A focused set of Fyne GUI polish changes building on 002: open the main window at the screen work
area; widen the task editor and move Preview into a right-hand pane; add an in-modal **Help** view
in that pane; drop the Schedule **Examples** button; show the command preview as a monospace code
block without the "Will run:" prefix; replace the Advanced Settings accordion with a custom
collapsible whose arrow points ▶ (collapsed) / ▼ (expanded); right-align the footer buttons; confirm
on Cancel when the form is dirty; and give every clickable control an app-wide pointer cursor.

Two items need more than a trivial edit because of Fyne API gaps: **maximize** (no Fyne API → a
platform work-area helper) and the **disclosure arrow** (Fyne's `Accordion` hardcodes the icons →
a small custom collapsible). Everything else is layout/wiring in `gui/`.

## Technical Context

**Language/Version**: Go (project toolchain); GUI via Fyne v2.7.4.

**Primary Dependencies**: Fyne v2 (`widget`, `container`, `layout`, `dialog`, `theme`,
`driver/desktop`, `canvas`); `fyne.io/x/fyne` (already added in 002, unaffected). The Windows
work-area query uses the standard library `syscall` against `user32.dll` — **no new module
dependency**.

**Storage**: none (GUI-only; no daemon/CLI/API/schema changes).

**Testing**: `go test ./gui/...` headless (Fyne test driver), plus the CI `-race` set for any
non-GUI helper. New unit tests for dirty-detection, pane toggle (Preview↔Help), the custom
collapsible state/icon, the cursor-aware buttons, and the work-area→Fyne-units conversion (pure
function).

**Target Platform**: Linux, macOS, Windows desktop (`gosched-gui`). Maximize is implemented for
Windows (primary target) with a generous fallback elsewhere.

**Project Type**: Desktop app (thin GUI client).

**Performance Goals**: No scheduling/hot-path impact; UI interactions remain instant.

**Constraints**: Keep all 002 behavior intact (FR-014). Avoid new dependencies (constitution).

**Scale/Scope**: ~1 dialog rebuild (`editor.go`), one window-sizing helper (+ platform files),
one custom collapsible widget, and cursor-aware button constructors swapped into ~14 call sites.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality** — PASS. New helpers (work-area query, collapsible, cursor buttons) are small,
  single-purpose, documented. The platform helper isolates `syscall` behind a build-tagged file
  with a clean fallback; no `panic`, errors handled.
- **II. Testing Standards** — PASS (planned). Headless GUI tests cover dirty-detection, pane
  toggling, collapsible icon/state, button cursor, and footer alignment; the work-area→units
  conversion is a pure unit-tested function. No wall-clock dependence.
- **III. UX Consistency** — PASS, directly advanced: consistent pointer affordance app-wide,
  standard disclosure-arrow direction, right-aligned actions, and safe-cancel confirmation.
- **IV. Performance** — PASS. No engine/hot-path change; no new allocations on any dispatch path.

No violations → Complexity Tracking empty.

## Project Structure

### Documentation (this feature)

```text
specs/003-gui-editor-refinements/
├── plan.md
├── spec.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── editor-ui.md        # updated dialog UI contract (panes, help, footer, cancel)
└── checklists/requirements.md
```

### Source Code (repository root)

```text
gui/
├── app.go                 # open window at work area (maximize-on-launch)
├── screen_windows.go      # NEW: work-area size via user32 (//go:build windows)
├── screen_other.go        # NEW: fallback work-area (//go:build !windows)
├── screen_test.go         # NEW: unit test for work-area→Fyne-units conversion
├── widgets.go             # cursor-aware button ctors + NEW custom collapsible widget
├── widgets_test.go        # collapsible icon/state + cursor button tests
├── editor.go              # two-pane layout, Help pane, code-block preview, right footer, cancel-confirm
├── editor_test.go         # dirty-detection, pane toggle, no-Examples, preview style
├── tasks.go               # cursor buttons
├── groups.go              # cursor buttons
├── triggers.go            # cursor buttons
├── schedule.go            # cursor buttons
└── alerts.go              # cursor buttons
```

**Structure Decision**: Single Go module; all work in `gui/`. No daemon/CLI/API/store changes.
The two non-trivial items get dedicated, isolated pieces: a build-tagged window-size helper and a
custom collapsible widget in `widgets.go`. Cursor coverage is achieved by swapping button
constructors at their ~14 call sites.

## Complexity Tracking

> No constitution violations — section intentionally empty.
