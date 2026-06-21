# Contract: Task Editor Dialog UI (003 refinements)

Updates the 002 editor contract. Items map to FR-001…FR-014. Behaviors not listed here are
unchanged from 002 ([002 contract](../../002-gui-task-editor-ux/contracts/editor-ui.md)).

## Window

- On launch the main window fills the screen work area (maximized appearance), respecting the
  taskbar (FR-001). Non-Windows falls back to a generous default size.

## Modal layout (two panes)

- The modal is ≈2× the prior width and split left/right (FR-002).
- **Left pane**: scrollable form — "What to run" (Name*, Command*, Arguments+caption), "When"
  (Timezone, Mode, Schedule | One-off date+time, Start at when sub-daily), and the Advanced
  collapsible. No "Preview" row in the form anymore.
- **Right pane**: shows the **Preview** by default, or **Help** when toggled (FR-003/FR-004/FR-005).

## Preview (right pane, default)

- Schedule summary + next runs (recurring), guidance when empty, "⚠ …" when invalid — as 002.
- Command line rendered in a **monospace code block**, with **no "Will run:" prefix** (FR-007/008).

## Help (right pane, on toggle)

- A Help control toggles the right pane to guidance covering every field with an example each
  (Name, Command, Arguments, Timezone, Mode, Schedule incl. anchor, One-off time, Overlap,
  Catch-up) (FR-004).
- Toggling back returns to the live Preview without losing input (FR-005).
- There is **no Examples button** beside the Schedule field (FR-006).

## Advanced Settings (custom collapsible)

- Starts collapsed; header arrow ▶ (`NavigateNext`) when collapsed, ▼ (`MenuDropDown`) when
  expanded (FR-009). Contains overlap + catch-up with the 002 human-readable labels.

## Footer

- Save and Cancel are right-aligned (FR-010), both pointer-cursor buttons.
- **Cancel when dirty** → confirmation prompt before discarding/closing (FR-011). **Cancel when
  untouched** → closes immediately (FR-012). Window-chrome/Escape dismissal follows the same
  dirty-guard where feasible.

## Cursor (app-wide)

- Every enabled clickable button across all tabs (Tasks/Schedule/Groups/Triggers/Alerts toolbars,
  list-action buttons) and dialogs shows the pointer/hand cursor on hover (FR-013). Table/list row
  selection is not a button and is out of scope.

## Compatibility

- All 002 editor behavior (mode visibility, validation + Save gating, combined preview content,
  interval anchor, timezone combo, one-off picker, advanced labels→wire values) remains intact
  (FR-014). No daemon/CLI/API/store changes.
