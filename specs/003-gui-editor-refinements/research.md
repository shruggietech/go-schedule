# Phase 0 Research: GUI Editor Refinements

Decisions resolving the Technical Context, each with rationale and alternatives. Grounded in the
Fyne v2.7.4 source (`~/go/pkg/mod/fyne.io/fyne/v2@v2.7.4`).

## 1. Open the main window maximized (FR-001)

- **Decision**: Fyne's `Window` interface exposes **no `Maximize()`** and **no screen-size API**
  (only `Resize`, `SetFullScreen`, `CenterOnScreen`). Implement a small platform helper that returns
  the screen **work area** (excludes the taskbar), convert it to Fyne units via the window canvas
  `Scale()`, then `Resize(workArea)` + `CenterOnScreen()` at startup.
  - `screen_windows.go` (`//go:build windows`): call `user32.dll!SystemParametersInfoW(SPI_GETWORKAREA)`
    via the standard-library `syscall` (`LazyDLL`/`LazyProc`) to get the work-area rect in pixels.
  - `screen_other.go` (`//go:build !windows`): return `(0,0)` → caller falls back to a generous
    fixed size (e.g. 1280×800).
  - Pure helper `windowSizeFor(workW, workH int, scale float32) fyne.Size` does the
    pixels→units division (and a small margin for window chrome) and is unit-tested.
- **Rationale**: True OS-maximize isn't reachable through Fyne's public API, and `SetFullScreen`
  is borderless (covers the taskbar and hides the title bar — which would also hide the icon just
  polished in 002). Sizing to the work area + centering looks maximized while respecting the
  taskbar. `syscall` keeps it dependency-free (constitution).
- **Alternatives considered**:
  - *`SetFullScreen(true)`* — borderless fullscreen; rejected (no title bar/taskbar, wrong feel).
  - *`golang.org/x/sys/windows`* — cleaner constants but adds a dependency; `syscall` suffices.
  - *Fixed large size* — doesn't adapt to the monitor; kept only as the non-Windows fallback.
- **HiDPI note**: convert work-area pixels to units with the canvas scale; if the scale isn't
  settled before first show, applying the resize on the first frame (or accepting a small margin)
  is sufficient — exact pixel-perfect maximize isn't required (SC-001 is "fills the work area").

## 2. Two-pane, wider editor (FR-002/FR-003)

- **Decision**: Increase the dialog size to roughly double width (≈1180×720) and lay the body out
  as a left/right split: left = the existing scrollable form (What to run / When / Advanced); right
  = a pane container that holds **either** the Preview **or** Help. Use `container.NewBorder` (or a
  fixed-ratio `container.NewGridWithColumns(2)`) so the two halves are stable; the Preview `FormItem`
  row is removed from the When form and its content moves to the right pane.
- **Rationale**: Side-by-side keeps inputs and preview visible together (SC-002). A grid/border
  split is simpler and more predictable than a draggable `container.NewHSplit` for a modal, though
  HSplit is a viable alternative if a movable divider is wanted.
- **Alternatives**: keep single column and just widen — rejected; the user specifically wants the
  preview on the right half.

## 3. In-modal Help pane (FR-004/FR-005/FR-006)

- **Decision**: Add a **Help** toggle (a cursor button in the right pane's header, or top-right of
  the modal) that swaps the right pane between the Preview and a Help view. The right pane is a
  `container.NewStack` (or swap of `Objects`) holding `previewView` and `helpView`; Help is a
  scrollable `widget.RichText` with a section per field (Name, Command, Arguments, Timezone, Mode,
  Schedule + anchor, One-off time, Overlap, Catch-up), each with a one-line explanation and an
  example. Toggling never rebuilds the form, so inputs and the computed preview persist (FR-005).
  Remove `examplesButton()` and its use beside the Schedule field (FR-006); its content folds into
  the Schedule section of Help.
- **Rationale**: One consolidated, in-context Help is more discoverable than a thin popup; reusing
  the right pane avoids a second window. RichText gives headings + monospace examples cheaply.
- **Alternatives**: a separate Help dialog/window — rejected (the spec wants it in the right pane).

## 4. Command preview as a code block, no prefix (FR-007/FR-008)

- **Decision**: Render the resolved command line with `widget.NewRichText` containing a
  `&widget.TextSegment{Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Monospace: true}}}`
  (or a code-block style), and drop the "Will run: " prefix string. `commandLinePreview` (from 002)
  still builds the text; only the presentation changes. Empty/guidance state stays as plain muted
  text.
- **Rationale**: Monospace reads as a real command line; removing the prefix declutters. RichText
  monospace is the lightest way to get code styling in Fyne.
- **Alternatives**: a bordered `Entry` set read-only with monospace — heavier and looks editable;
  rejected.

## 5. Advanced Settings disclosure arrow ▶/▼ (FR-009)

- **Decision**: Replace Fyne's `widget.Accordion` (its renderer hardcodes
  `IconNameArrowDropDown` when closed and `IconNameArrowDropUp` when open — the opposite of the
  requested convention) with a small custom **collapsible** widget: a header `cursorButton` whose
  icon is `theme.NavigateNextIcon()` (▶) when collapsed and `theme.MenuDropDownIcon()` (▼) when
  expanded, toggling the visibility of a content container. Used for Advanced Settings (overlap +
  catch-up), starting collapsed.
- **Rationale**: The icon direction can't be changed through the Accordion API, so a custom
  collapsible is required; it's tiny, testable (icon/state assertions), and as a bonus the header
  uses the pointer cursor (FR-013).
- **Alternatives**: post-process the Accordion's internal button icon — not reachable (unexported
  renderer); rejected.

## 6. Right-aligned footer buttons (FR-010)

- **Decision**: Build the footer as `container.NewBorder(nil, nil, nil, container.NewHBox(cancel,
  save))` (buttons in the trailing/right slot), or an HBox with a leading `layout.NewSpacer()`
  (which expands) before the buttons. The current code uses a non-expanding `widget.NewLabel("")`
  spacer — replace it with a real `layout.NewSpacer()` or the Border approach.
- **Rationale**: `layout.NewSpacer()` is the idiomatic expander; Border's trailing slot is equally
  clean. Either right-aligns reliably.

## 7. Cancel confirmation when dirty (FR-011/FR-012)

- **Decision**: Track a baseline of field values captured at open (empty for New, prefilled for
  Edit). `isDirty()` compares current values to the baseline (ignoring unchanged default timezone
  and unchanged overlap/catch-up defaults). The Cancel handler and the dialog close-intercept call
  `confirmCancel`: if `isDirty()`, show `dialog.NewConfirm("Discard changes?", …)` and only close on
  confirm; otherwise close immediately. Wire the same path to the dialog's dismiss/Escape via
  `CustomDialog`'s dismiss handling where feasible.
- **Rationale**: Prevents accidental data loss (SC-005) while staying silent for untouched forms.
  A captured baseline makes Edit-mode dirtiness correct (only changes from prefilled values count).
- **Alternatives**: a global "touched" flag set on any `OnChanged` — simpler but would treat the
  programmatic prefill as dirty unless carefully gated; the baseline comparison is more robust.

## 7a. Dialog dismiss interception

- **Decision**: The editor uses `dialog.NewCustomWithoutButtons` with our own footer buttons. To
  also guard window-chrome/Escape dismissal, prefer routing close through our `confirmCancel`. If
  the dialog type can't intercept Escape directly, document the limitation; the primary Cancel
  button is fully guarded.
- **Rationale**: Covers the common path; matches the spec's "where feasible".

## 8. App-wide pointer cursor (FR-013)

- **Decision**: Add cursor-aware constructors in `widgets.go` mirroring the stock ones —
  `newToolbarButton(label, icon, tapped)` and a no-icon variant — returning the existing
  `cursorButton` (which implements `desktop.Cursorable`). Swap the ~14 `widget.NewButton` /
  `widget.NewButtonWithIcon` call sites in `tasks.go`, `groups.go`, `triggers.go`, `schedule.go`,
  `alerts.go` to these. The custom collapsible header is already a `cursorButton`.
- **Rationale**: Reuses the proven 002 `cursorButton`; a single swap per call site gives full
  coverage (SC-006) with no behavior change.
- **Scope note**: `widget.Button` is the clickable control everywhere here. Table/list **row**
  selection is not a button and Fyne doesn't expose a per-row cursor; rows are out of scope for the
  pointer treatment (the spec's "row/list actions" are realized as buttons, which are covered).
  Log this boundary so coverage isn't overclaimed.

## Cross-cutting confirmations

- **No 002 regressions (FR-014)**: mode visibility, validation/Save gating, anchor, timezone combo,
  one-off picker, advanced labels are preserved; only their container/placement changes. Existing
  `gui` tests must keep passing, with `whenLabels`-style helpers updated for the new layout.
- **Testing**: headless tests for dirty-detection, Preview↔Help toggle, absence of the Examples
  button, code-block preview text, collapsible icon per state, and cursor on the new buttons; a
  pure unit test for `windowSizeFor`.
