# Feature Specification: Desktop UX and Options

**Feature Branch**: `codex/042-desktop-ux-options`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/042-desktop-ux-options`; headless GUI, focused race, and canonical eight-gate verification passed on 2026-09-03; native Windows visual acceptance remains part of the exact-candidate gate in #94

**Input**: Bundle GitHub issues #101, #103, #104, #105, and #106 into one coherent desktop-interface slice covering appearance preferences, storage-path visibility, navigation layout, orderly exit, task-row editing, and sharp Info-page text.

## Context

The current desktop interface is dark-only, exposes no preference or storage
view, uses a narrow leading tab rail that cannot anchor a command at its bottom,
requires a separate Edit-button action after task selection, and renders two
centered Info-page body lines more softly than adjacent text on Windows. These
reports share the same interface shell, typography, and headless verification
surface. Delivering them together avoids repeatedly restructuring navigation
and appearance state.

## Clarifications

### Session 2026-09-03

- Q: Which appearance modes and initial default should S042 define? -> A: Dark, light, and follow-system modes, with dark retained when no valid saved preference exists.
- Q: Which bounded font choices should the first Options view expose? -> A: Brand, system, and bundled monospace, with Brand retained as the default.
- Q: Can automated evidence close the native visual criteria in the bundled issues? -> A: No. S042 implements and tests the contracts, while #94 retains exact-candidate Windows 11 visual acceptance at standard and scaled DPI before those issues close.
- Q: How should unknown or unavailable paths be represented? -> A: Show only meaningful application-owned categories, label unavailable paths honestly, and never invent or scan unrelated profile locations.
- Q: How should Exit relate to navigation? -> A: Present it in the navigation rail but execute the same one-shot orderly shutdown path as the title-bar close control.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Choose and retain a readable appearance (Priority: P1)

As a desktop user, I can choose dark, light, or system-following appearance and
a supported interface font, immediately see that choice across the application,
restore the defaults, and have valid choices return on the next launch.

**Why this priority**: Appearance is the foundation for the new Options view and
the Info-page text correction. Theme and font state affect every other view and
dialog, so their lifecycle must be settled before the navigation shell is
considered complete.

**Independent Test**: Start with empty, valid, and malformed preference values;
change each appearance choice, rebuild the interface from the same preferences,
and verify the selected values, complete-interface update, safe fallback, and
default restoration without using a physical display.

**Acceptance Scenarios**:

1. **Given** no saved appearance preferences, **When** the application starts, **Then** it uses the documented dark and Brand defaults.
2. **Given** a user selects light mode and System font, **When** the choice is applied and the application is restarted, **Then** both choices remain active throughout every view and dialog.
3. **Given** follow-system mode, **When** the operating-system preference reports either light or dark, **Then** the complete interface uses the corresponding readable palette.
4. **Given** corrupt or unsupported saved values, **When** the application starts, **Then** it falls back to documented defaults without preventing startup.
5. **Given** non-default choices, **When** the user restores defaults, **Then** dark mode and Brand font apply immediately and remain after restart.
6. **Given** the Info view at any supported appearance and font, **When** version and attribution text is laid out, **Then** it remains centered, unwrapped, unclipped, and uses the same sharp text path as comparable nearby labels.

---

### User Story 2 - Understand local application storage (Priority: P1)

As a user diagnosing or preparing to uninstall the application, I can inspect
and copy the resolved locations of application-owned machine data, database,
logs, runtime state, desktop preferences, installed files, documentation, and
maintenance evidence, together with ownership, existence, and uninstall
retention information.

**Why this priority**: Users currently have no in-application explanation of
where machine-wide and per-user state lives. This information is necessary to
make preserve-or-wipe behavior understandable and supportable.

**Independent Test**: Supply controlled platform, configuration, executable,
and filesystem inputs; verify every category, resolved value, scope, existence,
and lifecycle label, then copy each available value through the view.

**Acceptance Scenarios**:

1. **Given** default storage locations, **When** Options opens, **Then** the view distinguishes machine-wide data from per-user desktop preferences.
2. **Given** an application-owned path exists or is absent, **When** it is displayed, **Then** its existence state is accurate without creating or deleting it.
3. **Given** an available resolved path, **When** the user selects its text or activates Copy, **Then** the exact displayed value is available for pasting.
4. **Given** a category that is not meaningful or discoverable on the current platform, **When** Options opens, **Then** the view omits it or marks it unavailable without presenting a fabricated path.
5. **Given** an administrator-configured location outside application-owned defaults, **When** paths are presented, **Then** the view does not claim that unrelated data will be erased by uninstall.

---

### User Story 3 - Navigate comfortably and exit predictably (Priority: P2)

As a desktop user, I can use a spacious leading navigation rail with Options
immediately above Info and a visually separate Exit command anchored at the
bottom-right of that rail.

**Why this priority**: Options and Exit cannot be placed correctly in the
existing fixed tab stack. One shared navigation-shell correction should provide
adequate width, selection state, keyboard access, and stable bottom placement.

**Independent Test**: Construct the complete shell at the default and minimum
supported content sizes, select every view, activate Exit by pointer and
keyboard paths, and verify geometry, focus order, selection semantics, and
one-shot shutdown.

**Acceptance Scenarios**:

1. **Given** the complete view list, **When** the application opens, **Then** every label has balanced horizontal breathing room and the longest supported label is unclipped.
2. **Given** ordinary view navigation, **When** Options or Info is selected, **Then** the selected view is visibly identified and Options is directly above Info.
3. **Given** any window height, **When** the shell is laid out, **Then** Exit remains separated from ordinary views at the bottom-right of the rail.
4. **Given** Exit or the title-bar close control, **When** either is activated repeatedly, **Then** cancellation and window close occur exactly once through the same orderly shutdown behavior.
5. **Given** the supported small-window clamp, **When** the shell is resized, **Then** navigation remains readable without consuming unreasonable content width.

---

### User Story 4 - Edit the intended task directly (Priority: P2)

As a task author, I can double-click a task row to open exactly one fully
populated Edit dialog for that task while single-click, keyboard, and toolbar
editing continue to behave as before.

**Why this priority**: This is a small, independently valuable desktop
interaction, but it must preserve identity through live list refreshes and share
the existing guarded editor lifecycle rather than opening duplicate dialogs.

**Independent Test**: Drive row selection and double activation with controlled
task refreshes and detail-fetch outcomes, and assert editor identity, single
dialog ownership, fallback behavior, and unchanged selection behavior.

**Acceptance Scenarios**:

1. **Given** a visible task row, **When** it is clicked once, **Then** it is selected and no editor opens.
2. **Given** a visible task row, **When** it is double-clicked, **Then** exactly one Edit dialog opens using the existing current-detail lookup path.
3. **Given** the list order changes after a row is rendered, **When** that row is double-clicked, **Then** the application edits the row's task identity rather than whichever task now occupies its former index.
4. **Given** a detail lookup failure, **When** a task row is double-clicked, **Then** the existing degraded edit fallback remains available.
5. **Given** rapid repeated double-clicks or an already-open task editor, **When** another edit activation occurs, **Then** no stacked editor opens.
6. **Given** empty space or a stale task identity, **When** it receives an activation, **Then** no editor opens and no unrelated task is selected.

### Edge Cases

- Saved appearance strings may be empty, malformed, or from a future version.
- The operating-system theme preference may change while the application is open.
- A path may be unavailable, relative, nonexistent, permission-restricted, or
  point outside the application-owned default roots.
- The running executable may be a development binary rather than an installed
  package, so installed documentation might not be discoverable.
- The cleanup-evidence path is Windows-specific and normally absent after a
  successful cleanup.
- A live refresh may remove or reorder a task between pointer-down and double
  activation.
- Shutdown may be requested from the title bar and Exit control in quick
  succession.
- Long translated or badge-bearing navigation labels must not collapse the rail
  or starve the content area.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The navigation rail MUST contain Tasks, Groups, Chains, Schedule, Activity, Options, and Info in that order, with Options immediately above Info.
- **FR-002**: Exit MUST be visually separated from ordinary views, anchored at the bottom-right of the navigation rail, keyboard reachable, and never represented as selected content.
- **FR-003**: Exit and title-bar close MUST share one idempotent orderly shutdown path that cancels background activity and closes the window exactly once.
- **FR-004**: The navigation rail MUST derive a stable minimum width from its supported labels plus balanced theme spacing and MUST remain usable at the supported small-window clamp.
- **FR-005**: The Options view MUST offer Dark, Light, and Follow system appearance modes; Dark MUST remain the default when no valid preference exists.
- **FR-006**: The Options view MUST offer Brand, System, and Monospace font choices; Brand MUST remain the default when no valid preference exists.
- **FR-007**: Appearance changes MUST apply immediately to all current views and dialogs and MUST apply to subsequently created interface elements.
- **FR-008**: Valid appearance choices MUST persist per user and restore on the next launch; invalid stored values MUST fall back safely without blocking startup.
- **FR-009**: Users MUST be able to restore both appearance choices to their documented defaults in one action.
- **FR-010**: Every offered appearance and font combination MUST retain readable contrast, visible focus, unclipped controls, and consistent dialog presentation.
- **FR-011**: The Info-page version and attribution lines MUST remain centered and unwrapped, preserve their exact semantic text and dynamic version source, and use the same body-text rendering contract as sharp nearby interface text.
- **FR-012**: Options MUST present resolved application-owned categories for the machine data root, task database, logs, runtime state, per-user desktop preferences, executable directory, installed documentation when discoverable, and platform-specific maintenance evidence when meaningful.
- **FR-013**: Every storage row MUST state its category, ownership scope, current existence, and expected behavior under software-only removal and explicit data wipe; development paths MUST NOT be represented as installer-owned.
- **FR-014**: Every available displayed path MUST be selectable and have an explicit Copy action that copies the exact resolved value.
- **FR-015**: Storage discovery MUST be read-only, MUST NOT scan unrelated profiles, MUST NOT create missing paths, and MUST NOT imply that user-created exports or configured external locations are application-owned.
- **FR-016**: Double-clicking a valid task row MUST open the same Edit workflow and current-detail lookup used by the toolbar Edit command.
- **FR-017**: Task-row activation MUST retain stable task identity across live refreshes and MUST ignore stale or absent identities.
- **FR-018**: Single-click selection, toolbar editing, and keyboard interaction MUST remain available and MUST NOT open an editor unexpectedly.
- **FR-019**: The application MUST prevent simultaneous stacked task editors and MUST release the guard whenever an editor closes through Save, Cancel, or dismissal.
- **FR-020**: Headless tests MUST cover preference defaults, persistence, invalid-value fallback, every theme/font choice, storage classification and copying, navigation order and geometry, shared one-shot shutdown, sharp Info-label invariants, and task-row single/double activation with refresh and lookup failures.
- **FR-021**: Native visual acceptance MUST cover standard DPI and at least one scaled-DPI Windows configuration, light and dark appearance, navigation placement, and Info text sharpness against the exact staged candidate before issues requiring visual proof close.

### Key Entities

- **Appearance preferences**: Per-user mode and font identifiers, each constrained to a documented bounded set with safe defaults.
- **Storage location**: A category label, resolved value or unavailable state, ownership scope, existence state, and preserve/wipe lifecycle description.
- **Navigation destination**: An ordinary selectable view with a stable order, label, content, and selected state.
- **Navigation command**: A focusable action placed within the rail but excluded from destination selection.
- **Task-row identity**: The stable task identifier bound to a rendered row and resolved against current model state before editing.
- **Shutdown state**: A one-way application lifecycle transition that guarantees cancellation and close happen at most once.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 3 appearance modes and all 3 font choices apply immediately and restore correctly after a simulated restart for 100 percent of valid combinations.
- **SC-002**: Every invalid or missing appearance value falls back to Dark and Brand without a startup error.
- **SC-003**: Every displayed available storage location is copyable byte-for-byte and 100 percent of rows identify scope, existence, and both uninstall outcomes.
- **SC-004**: At the 1280 by 800 default content size and supported 800 by 600 launch viewport, all navigation labels remain unclipped, Options precedes Info, Exit remains bottom anchored, and the content view remains usable.
- **SC-005**: Single-click opens zero editors and each valid double-click opens exactly one correct editor across normal, reordered, removed, and detail-failure scenarios.
- **SC-006**: Repeated or concurrent close requests perform application cancellation and window closure exactly once without panic or leaked background work.
- **SC-007**: Headless GUI tests, race tests for non-widget state, and all eight canonical repository gates pass without reducing existing coverage.
- **SC-008**: Exact-candidate Windows evidence at 100 percent and one scaled-DPI setting shows the two Info lines as sharp as adjacent body text, with readable light/dark palettes and correctly placed navigation controls.

## Assumptions

- Existing per-user preferences are stored by the application framework under
  the established application identifier and are an appropriate persistence
  mechanism for non-sensitive appearance choices.
- Dark and Brand remain defaults to preserve the current visual identity for
  users who make no choice.
- System font means the framework's platform-appropriate default family;
  Monospace means the existing bundled monospace face.
- The machine data directory represents the owned umbrella for task data,
  default configuration/runtime state, database, and logs. An optional daemon
  configuration overlay outside that root cannot be discovered from the current
  local client contract and must not be claimed as wipe-owned.
- Installed-file and documentation locations are derived from the running
  executable and may be labeled unavailable in development or non-installed
  contexts.
- This slice changes no scheduler, daemon protocol, installer deletion boundary,
  or release-promotion rule.

## Scope Boundaries

### In scope

- GitHub issues #101, #103, #104, #105, and #106.
- One coherent navigation and appearance architecture for the desktop GUI.
- Per-user appearance persistence and read-only path presentation.
- Automated cross-platform and headless verification plus exact-candidate
  attended validation instructions.

### Out of scope

- Closing #94, #98, or coordinator #96.
- Automating or substituting for the attended Windows 11 release-candidate gate.
- Arbitrary user-supplied fonts, downloaded fonts, font-file browsing, or a full
  theme editor.
- Editing daemon configuration or changing storage locations from Options.
- Opening, deleting, migrating, or repairing displayed paths.
- Changing task execution diagnostics tracked by #102.
- Packaging, tagging, publishing, merging, or cutting a release.
