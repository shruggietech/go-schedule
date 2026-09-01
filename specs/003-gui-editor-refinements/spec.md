# Feature Specification: GUI Editor Refinements

**Feature Branch**: `003-gui-editor-refinements`

**Created**: 2026-06-20

**Status**: Implemented

**Delivery**: Commit `a21c22e`

**Input**: User description: "GUI task editor refinements (follow-up to 002) — full-size window, two-pane wider modal with Preview on the right, in-modal Help pane, remove Examples, code-block command preview, disclosure-arrow direction for Advanced Settings, right-aligned Save/Cancel, Cancel confirmation when dirty, and app-wide pointer cursor on clickable controls."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Roomier, two-pane task editor (Priority: P1)

When a user opens **New Task** / **Edit Task**, the modal is about twice as wide as before and
split into two halves: the form fields fill the left half, and the **Preview** (schedule summary +
next runs and the command preview) sits in the right half. The user can see their inputs and the
live preview side by side without the preview being squeezed into a narrow form row.

**Why this priority**: The single-column layout crowds the preview and pushes content below the
fold; the two-pane layout is the structural change the other editor items build on.

**Independent Test**: Open the editor and confirm it is roughly double the previous width, with
form fields on the left and the Preview content on the right half.

**Acceptance Scenarios**:

1. **Given** the editor opens, **When** the user views it, **Then** the dialog is about twice the
   prior width and form fields occupy the left half.
2. **Given** a valid schedule and a command with arguments, **When** the user types, **Then** the
   schedule summary, next runs, and command preview all render in the right half.
3. **Given** the editor in either Mode, **When** the user switches Recurring/One-off, **Then** the
   left-half fields update while the right-half pane remains the preview area.

---

### User Story 2 - In-modal Help (Priority: P1)

A **Help** button in the editor reveals guidance that explains each field, with examples and how
the pieces work together (schedule phrasing, the anchor, timezones, one-off times, overlap/
catch-up). The help appears in the right-half pane — the same area the Preview uses — so the user
can read guidance without leaving the dialog. The separate **Examples** button next to the
Schedule field is removed, since its content is now part of Help.

**Why this priority**: First-time users need discoverable, in-context guidance; consolidating it
into one Help pane (and dropping the thin Examples popup) is a clear usability win.

**Independent Test**: Open the editor, click Help, and confirm field-by-field guidance with
examples appears in the right pane; confirm there is no longer an Examples button by the Schedule
field.

**Acceptance Scenarios**:

1. **Given** the editor is open, **When** the user clicks Help, **Then** the right pane shows
   guidance covering every field with examples.
2. **Given** Help is showing, **When** the user dismisses Help (or toggles back), **Then** the
   right pane returns to showing the Preview.
3. **Given** the editor is open, **When** the user looks at the Schedule field, **Then** there is
   no Examples button beside it.

---

### User Story 3 - Cleaner command preview (Priority: P2)

The resolved command preview is shown as a monospace code block, without the "Will run:" text
prefix — so it reads like an actual command line.

**Why this priority**: Small but improves legibility and matches user expectation for showing a
command.

**Independent Test**: Enter a command with arguments and confirm the preview shows the command in
a monospace/code-block style with no "Will run:" prefix.

**Acceptance Scenarios**:

1. **Given** a command and arguments, **When** the user views the preview, **Then** the command
   line appears in a monospace code block.
2. **Given** the preview, **When** the user reads it, **Then** there is no "Will run:" prefix text.

---

### User Story 4 - Predictable controls: disclosure arrow, button placement, cancel safety (Priority: P2)

The **Advanced Settings** disclosure indicator points right when collapsed and down when expanded.
The **Save** and **Cancel** buttons are aligned to the right of the footer. If the user has entered
anything into the form, clicking **Cancel** asks for confirmation before discarding; if the form is
untouched, Cancel closes immediately.

**Why this priority**: These are standard desktop conventions; the cancel-confirm prevents
accidental loss of entered data.

**Independent Test**: Toggle Advanced Settings and confirm the arrow direction; confirm Save/Cancel
are right-aligned; type into a field and click Cancel to get a confirmation prompt, and on an empty
form confirm Cancel closes with no prompt.

**Acceptance Scenarios**:

1. **Given** Advanced Settings is collapsed, **When** the user views it, **Then** the disclosure
   arrow points right; **When** expanded, **Then** it points down.
2. **Given** the editor footer, **When** the user views it, **Then** Save and Cancel are aligned to
   the right.
3. **Given** the user has typed into any field, **When** they click Cancel, **Then** a confirmation
   prompt appears before the dialog closes and discards input.
4. **Given** an untouched form, **When** the user clicks Cancel, **Then** the dialog closes
   immediately with no prompt.

---

### User Story 5 - Bounded windowed launch (Priority: P3)

When the application launches, its main window opens at a useful restored size with visible
desktop margins. Specification 037 intentionally supersedes the original full-work-area choice
after first-run observation showed that it looked maximized and consumed the user's workspace.

**Why this priority**: Quality-of-life; independent of the editor changes.

**Independent Test**: Launch the app and confirm it requests 1280x800 logical units when space
permits, capped independently to 90 percent of the selected monitor's logical work area.

**Acceptance Scenarios**:

1. **Given** the app is launched, **When** the main window appears, **Then** it is restored,
   centered, and bounded within 90 percent of the selected monitor's work area.

---

### User Story 6 - Consistent pointer cursor on clickable controls (Priority: P3)

Every clickable control across the application — toolbar buttons on all tabs (Tasks, Schedule,
Groups, Triggers, Alerts), list/row actions, dialog buttons — shows a pointer/hand cursor on
hover, matching the editor's Save/Cancel buttons.

**Why this priority**: Consistency; signals clickability everywhere, not just in the editor.

**Independent Test**: Hover the toolbar buttons on each tab and other clickable controls and
confirm the cursor becomes a pointer/hand.

**Acceptance Scenarios**:

1. **Given** any tab's toolbar buttons, **When** the user hovers a button, **Then** the cursor is a
   pointer/hand.
2. **Given** clickable row/list actions and dialog buttons, **When** the user hovers them, **Then**
   the cursor is a pointer/hand.

---

### Edge Cases

- **Dirty detection for Cancel**: "entered anything" means any user edit to a field beyond its
  initial/prefilled value. In Edit mode, the prefilled values are the baseline — only changes from
  that baseline count as dirty. The default Timezone ("Local") and default overlap/catch-up
  selections are not, by themselves, "entered" data.
- **Help vs Preview toggle**: showing Help must not lose the user's inputs or the computed preview;
  returning from Help shows the current live preview.
- **Closing via window chrome / Escape**: dismissing the dialog by means other than the Cancel
  button should follow the same dirty-confirmation behavior where feasible (so input isn't lost
  silently).
- **Very small screens**: the main window and wider two-pane modal must stay within the available
  screen work area, using existing scroll regions rather than forcing a nominal minimum beyond it.
- **Multi-monitor launch**: startup sizing uses the selected monitor's reserved work area rather
  than dimensions from a different, larger primary monitor.
- **Pointer cursor coverage**: disabled buttons need not show the pointer; only enabled, clickable
  controls do.

## Requirements *(mandatory)*

### Functional Requirements

#### Window & layout

- **FR-001**: As corrected by specification 037, the main application window MUST open restored at
  `min(1280x800, 90 percent of the selected monitor work area)` in logical units. This requirement
  intentionally reverses the original full-work-area behavior.
- **FR-002**: The task editor modal MUST be approximately twice its previous width and MUST present
  a two-pane layout: form fields on the left, a right-hand pane for Preview/Help.
- **FR-003**: The Preview content (schedule summary + next runs, and the command preview) MUST
  render in the right-hand pane.

#### Help

- **FR-004**: The editor MUST provide a Help control that displays field-by-field guidance with
  examples (covering name, command, arguments, timezone, mode, schedule incl. the anchor, one-off
  time, overlap, catch-up) in the right-hand pane.
- **FR-005**: The user MUST be able to return the right-hand pane from Help back to the live
  Preview; toggling MUST not lose entered input or the computed preview.
- **FR-006**: The standalone Examples control next to the Schedule field MUST be removed; its
  content is covered by Help.

#### Command preview

- **FR-007**: The command preview MUST render the resolved command line in a monospace/code-block
  style.
- **FR-008**: The command preview MUST NOT include a "Will run:" (or similar) text prefix.

#### Controls & affordances

- **FR-009**: The Advanced Settings disclosure indicator MUST point right when collapsed and down
  when expanded.
- **FR-010**: The Save and Cancel buttons MUST be aligned to the right of the editor footer.
- **FR-011**: When the form has unsaved user input (changed from its baseline), clicking Cancel
  MUST prompt for confirmation before discarding and closing.
- **FR-012**: When the form is untouched, Cancel MUST close the dialog immediately without a prompt.
- **FR-013**: All enabled, clickable controls across the application (toolbar buttons on every tab,
  row/list actions, and dialog buttons) MUST display a pointer/hand cursor on hover.

#### Compatibility

- **FR-014**: All behavior delivered in feature 002 (mode-driven field visibility, required-field
  validation and Save gating, combined preview content, interval anchor, timezone combo, one-off
  picker, advanced settings with human-readable labels) MUST remain intact.

### Key Entities *(include if feature involves data)*

- **Editor pane state**: which content the right-hand pane currently shows (Preview vs Help).
- **Form dirtiness**: whether the user has changed any field from its baseline (drives the Cancel
  confirmation).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On launch, the main window requests 1280x800 logical units where space permits and
  never exceeds 90 percent of either available logical work-area dimension.
- **SC-002**: In the editor, a user can view their inputs and the live preview side by side without
  scrolling the preview out of view.
- **SC-003**: A first-time user can open Help and find an explanation and example for every field
  without leaving the dialog.
- **SC-004**: The command preview reads as a plain command line (monospace, no prefix).
- **SC-005**: A user who has typed into the form cannot lose that input by an accidental single
  Cancel click — a confirmation is always shown first.
- **SC-006**: Every enabled clickable control in the app shows a pointer cursor on hover (0 controls
  that look clickable but show the default arrow).
- **SC-007**: All 002 editor behaviors continue to pass their existing checks (no regressions).

## Assumptions

- Target is the Go-native Fyne desktop GUI (`gosched-gui`); no daemon/CLI/API changes are required.
- The startup window remains restored. User-controlled maximize, restore, resize, minimize, and
  fullscreen behavior remains available after launch.
- "Twice the width" is approximate; the exact dimensions are chosen so the two panes are comfortable
  and the modal still fits common screens.
- The right-hand pane defaults to showing the Preview; Help is shown on demand and toggles back.
- Help content is static, in-app text (no external links required), maintained alongside the editor.
- "Dirty" is judged against the field baselines at open (empty for New, prefilled for Edit);
  default selections that the user did not change do not count as dirty.
- App-wide pointer cursor is achieved by routing clickable controls through the same cursor-aware
  control used for the editor buttons in 002; where a control type cannot expose a custom cursor, an
  equivalent affordance is acceptable.
