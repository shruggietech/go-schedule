# Phase 1 Data Model: GUI Editor Refinements

GUI-only; no persistent data. The entities below are transient UI state.

## 1. Editor right-pane state

| State | Meaning |
|-------|---------|
| `panePreview` | right pane shows the live Preview (schedule summary + next runs + command code block) — the default |
| `paneHelp` | right pane shows the Help guidance |

Transitions: Help toggle flips `panePreview ⇄ paneHelp`. Toggling never rebuilds the form, so
inputs and the computed preview persist (FR-005).

## 2. Form baseline & dirtiness (drives Cancel confirm)

Baseline captured when the editor opens:

| Mode | Baseline |
|------|----------|
| New | empty fields; defaults Timezone=`Local`, Overlap/Catch-up = their default labels |
| Edit | the task's prefilled values |

`isDirty()` = any of {name, command, args, timezone, schedule, startAt, oneOffDate, oneOffTime,
overlap label, catch-up label} differs from its baseline value. Unchanged defaults are **not**
dirty. Used by `confirmCancel` (FR-011/FR-012).

## 3. Window sizing

| Item | Meaning |
|------|---------|
| work-area (px) | screen area excluding the taskbar (Windows: `SPI_GETWORKAREA`; else 0,0) |
| scale | window canvas scale factor (Fyne units = px / scale) |
| `windowSizeFor(w,h,scale)` | pure conversion → `fyne.Size` in units, minus a small chrome margin; falls back to a generous default when work-area is unknown |

## 4. Custom collapsible (Advanced Settings)

| Field | Meaning |
|-------|---------|
| `open` | whether the content is expanded |
| header icon | `NavigateNextIcon` (▶) when `!open`, `MenuDropDownIcon` (▼) when `open` (FR-009) |
| content | the overlap + catch-up form, shown only when `open`; starts collapsed |

## 5. Cursor-aware buttons (presentation only)

All app toolbar/dialog buttons are constructed as `cursorButton` (implements `desktop.Cursorable`
→ `PointerCursor`). No behavioral change vs `widget.Button`; only the hover cursor differs
(FR-013). Disabled buttons need not show the pointer.
