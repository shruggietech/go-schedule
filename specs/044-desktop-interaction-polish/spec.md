# Feature Specification: Desktop Interaction and Appearance Polish

**Feature Branch**: `codex/044-desktop-interaction-polish`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Focused, race, full-suite, and canonical eight-gate verification passed 2026-09-03 on `codex/044-desktop-interaction-polish`; exact-candidate native Windows evidence remains assigned to release qualification.

**Input**: User description: "S044 bundles GitHub issues #101, #104, #105, #106, #109, and #111 into one desktop interaction and appearance hardening slice."

**Issue Traceability**: [#101](https://github.com/shruggietech/go-schedule/issues/101), [#104](https://github.com/shruggietech/go-schedule/issues/104), [#105](https://github.com/shruggietech/go-schedule/issues/105), [#106](https://github.com/shruggietech/go-schedule/issues/106), [#109](https://github.com/shruggietech/go-schedule/issues/109), [#111](https://github.com/shruggietech/go-schedule/issues/111)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Read Every Interactive State (Priority: P1)

As a desktop user, I can read labels and recognize controls while they are at rest, hovered, focused, selected, pressed, or unavailable in every supported color mode.

**Why this priority**: The current active-navigation and dialog-button hover states can make their text effectively disappear. This blocks confident use of ordinary controls and is the highest-priority defect in the slice.

**Independent Test**: Exercise representative navigation, dialog, selection, button, link, and unavailable states in dark, light, and follow-system modes and verify that foregrounds, backgrounds, outlines, and non-color cues remain distinct and readable.

**Acceptance Scenarios**:

1. **Given** a selected navigation destination, **When** the pointer hovers it, **Then** its label remains readable and the destination remains visibly selected.
2. **Given** a primary dialog action, **When** it is hovered, focused, pressed, or unavailable, **Then** its label and state remain unambiguous.
3. **Given** any supported color mode, **When** representative shared controls enter every supported interaction state, **Then** normal text reaches at least 4.5:1 contrast and essential non-text indicators reach at least 3:1.
4. **Given** a selected, focused, disabled, success, warning, or error state, **When** color perception is unavailable, **Then** another visible cue still communicates the state.

---

### User Story 2 - Start With Crisp, Comfortable Typography (Priority: P1)

As a Windows desktop user, I receive the crisp framework-selected System font by default and can choose among a small, clearly named set of readable alternatives without seeing stale or partially updated text.

**Why this priority**: Attended S043 testing showed that ordinary brand-font body text looked fuzzy on a QHD ultrawide display while selecting System fixed the problem immediately.

**Independent Test**: Start with clean, legacy, invalid, and explicitly saved preferences, verify the effective font choice, change among every offered font, and compare the Info version and attribution text with nearby interface text across reopen, resize, minimize/restore, and both color palettes.

**Acceptance Scenarios**:

1. **Given** a new user or an invalid font preference, **When** the application starts, **Then** System is selected and used throughout the interface.
2. **Given** a user who explicitly selected a valid font before this upgrade, **When** the application starts, **Then** that explicit selection is preserved.
3. **Given** the font selector, **When** it is opened, **Then** it offers System, Monospace, and a small curated set of familiar sans-serif alternatives under family-level names.
4. **Given** any offered font, **When** it is selected or restored after restart, **Then** existing and subsequently opened views and dialogs use it consistently.
5. **Given** the Info view in light or dark mode, **When** it is resized, minimized/restored, or reopened, **Then** its version and attribution remain centered, unclipped, correct, and as visually sharp as nearby interface text.

---

### User Story 3 - Scan and Copy Storage Information Efficiently (Priority: P2)

As a user inspecting application storage, I can scan compact aligned rows, understand unavailable entries at a glance, copy available paths, and use each appearance selector without being offered its current value as though it were a new action.

**Why this priority**: The current card-per-path layout consumes excessive space, unavailable rows look partly actionable, and selectors misleadingly repeat the active choice.

**Independent Test**: Render mixed available and unavailable storage entries at default and supported-small window sizes, inspect row alignment and muted states, copy every available path, and open each selector before and after a selection.

**Acceptance Scenarios**:

1. **Given** application storage entries, **When** Options is opened, **Then** the entries appear as compact aligned rows with category and path visible and Copy at the far right.
2. **Given** supporting ownership, scope, existence, or removal detail, **When** the row cannot show it without crowding, **Then** the detail remains available through a concise hover/focus explanation without horizontal scrolling.
3. **Given** an unavailable path, **When** its row is rendered, **Then** the entire row is visibly muted, the path cannot be selected or copied, and its reason remains available.
4. **Given** a valid available path, **When** Copy is activated, **Then** the exact path enters the clipboard.
5. **Given** a selector with a current value, **When** its choices are shown, **Then** the current value is identified as current and cannot be mistaken for an unapplied alternative.

---

### User Story 4 - Navigate and Scroll Comfortably (Priority: P2)

As a desktop user, I see a clearly bounded navigation rail with balanced spacing and a recognizable Exit command, and I can traverse long vertical views with a responsive wheel setting that I can tune and restore.

**Why this priority**: Navigation placement now works, but its boundary and Exit glyph need visual finishing. Windows wheel movement remains slow enough to make the longer Options view frustrating.

**Independent Test**: Exercise the navigation rail at default and supported small sizes, activate Exit through pointer and keyboard paths, and drive a long view using normalized wheel input at every supported sensitivity while checking keyboard, touchpad, drag, and scrollbar behavior.

**Acceptance Scenarios**:

1. **Given** the desktop shell, **When** any destination is shown, **Then** a visible separator defines the rail boundary and every label has balanced horizontal space without clipping.
2. **Given** the Exit command, **When** it is displayed, **Then** it remains bottom-right anchored, visually separated from destinations, and has a semantically colored glyph without appearing selected.
3. **Given** one conventional mouse-wheel detent on a long application-owned vertical view, **When** the default sensitivity is active, **Then** content moves immediately by a useful, bounded amount.
4. **Given** the Scroll sensitivity option, **When** a user changes it, **Then** the bounded value applies immediately to every application-owned vertical view, persists across restart, and can be restored to the documented default.
5. **Given** precision touchpad, keyboard, drag, or scrollbar input, **When** sensitivity changes, **Then** those input methods remain controllable and are not multiplied as mouse-wheel detents.

### Edge Cases

- Clean preferences, S042-era preferences, unknown values, and incomplete preference writes all resolve to safe documented defaults without blocking startup.
- A color-mode change while a control is hovered or focused refreshes the full state without briefly producing an unreadable foreground/background pair.
- A font change while a dialog is open updates both the dialog and underlying view without stale fonts.
- Very long storage paths wrap or elide vertically and expose their exact value; they never introduce a horizontal scrollbar.
- Missing storage paths, missing clipboard access, and runtime-metadata refresh failure do not make unavailable data look actionable or erase the last known truthful inventory.
- The shortest supported window height keeps Exit separate from the destination stack and retains access to long Options content.
- Repeated close requests still execute orderly shutdown exactly once.
- Nested content does not apply wheel sensitivity more than once to the same mouse-wheel event.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The application MUST define readable foreground, background, border, and glyph treatment for rest, hover, selected, selected-plus-hover, focus, pressed, and disabled states in dark, light, and follow-system modes.
- **FR-002**: Normal interactive text MUST reach at least 4.5:1 contrast in each state, and essential non-text indicators MUST reach at least 3:1.
- **FR-003**: Selected, focused, disabled, success, warning, and error states MUST retain a non-color cue.
- **FR-004**: Shared state treatment MUST cover navigation controls, ordinary and dialog buttons, selectors, links, and representative list/row controls.
- **FR-005**: System MUST be the default interface font for clean, reset, missing, unknown, or corrupt font preferences.
- **FR-006**: An explicitly saved supported font choice MUST survive upgrade and restart.
- **FR-007**: The interface font selector MUST offer System, Monospace, and a small curated set of familiar sans-serif family choices; each family MUST be presented once under a clear family-level name.
- **FR-008**: Every offered font MUST apply immediately and consistently to existing and subsequently opened views, dialogs, normal text, bold text, monospace-semantic text, and symbols as appropriate.
- **FR-009**: The Info version MUST continue to derive from the runtime build version, and the version and attribution MUST remain centered, unclipped, and readable in both palettes at supported sizes.
- **FR-010**: Options MUST present known storage locations as compact aligned rows rather than one card per location.
- **FR-011**: Each storage row MUST keep category, resolved path or unavailability, and Copy action legible without horizontal scrolling; supporting details MUST remain discoverable by pointer and keyboard.
- **FR-012**: An unavailable storage entry MUST mute the entire row and MUST NOT expose path selection or Copy as active behavior.
- **FR-013**: An available storage entry MUST copy its exact resolved path and MUST preserve the existing ownership and removal claims.
- **FR-014**: Appearance selectors MUST distinguish the current value from unapplied alternatives and MUST NOT present the current value as a second actionable choice.
- **FR-015**: A visible separator MUST define the leading navigation rail's content boundary while preserving balanced insets, unclipped labels, and the supported small-window layout.
- **FR-016**: Exit MUST remain an application command anchored at the rail's bottom-right, use a semantic colored glyph, remain distinct from selected destinations, and share the one-shot orderly shutdown path.
- **FR-017**: Options MUST expose a bounded Scroll sensitivity setting with a documented responsive default, immediate application, persistence, safe normalization, and restore-default behavior.
- **FR-018**: Scroll sensitivity MUST affect application-owned vertical mouse- wheel handling consistently without adding horizontal scrolling, multiplying nested events, or altering precision touchpad, keyboard, drag, or scrollbar semantics.
- **FR-019**: Automated checks MUST cover palette contrast, representative interaction states, font defaults and migration, selector semantics, storage row availability, navigation geometry, scroll normalization/delta behavior, and one-shot shutdown.
- **FR-020**: Native Windows release evidence MUST cover dark and light modes, 100% and one scaled-DPI setting, conventional wheel input, and precision touchpad behavior when such hardware is available. Hardware explicitly unavailable in the attended environment MUST be recorded as unavailable, not represented as passed.

### Key Entities

- **Appearance preference**: Per-user color mode and interface-font choice, including its safe default and normalization rules.
- **Interaction state style**: The coordinated foreground, background, border, glyph, and non-color cue for a control state under a color mode.
- **Storage row**: A known storage category, resolved path or unavailable reason, ownership/removal details, availability state, and Copy capability.
- **Scroll preference**: The per-user bounded sensitivity level, its documented default, and the normalized mouse-wheel movement derived from it.
- **Navigation command**: A bottom-anchored action that does not represent or retain a content-view selection.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every tested normal-text interaction state in all three color modes reaches at least 4.5:1 contrast; every tested essential non-text indicator reaches at least 3:1.
- **SC-002**: In a clean profile, System is the active font on first launch in 100% of tests; all previously supported explicit valid choices survive upgrade and restart.
- **SC-003**: A user can identify a storage category, read or obtain its exact path, understand its availability, and activate Copy without horizontal scrolling at both 1280 by 800 and the supported small-window clamp.
- **SC-004**: One conventional mouse-wheel detent moves every long application- owned vertical view immediately at the default setting, and every supported sensitivity level remains within its documented bound.
- **SC-005**: Navigation labels and Exit remain unclipped and non-overlapping at default size and the supported small-window clamp, with Exit still anchored at the bottom-right.
- **SC-006**: All appearance, storage, navigation, and scroll preference regression checks pass, along with the complete canonical verification suite.
- **SC-007**: Native Windows inspection finds no unreadable hover/focus/pressed/ disabled state and no visibly fuzzy Info version or attribution at 100% and at least one scaled-DPI setting.

## Assumptions

- S042's working Options, storage inventory, navigation placement, double-click editing, and orderly shutdown are the baseline; this slice refines them rather than replacing their behavior.
- System is the new default because the attended S043 comparison found it materially sharper than Brand on the reporter's QHD ultrawide display.
- Existing valid explicit font preferences are user choices and therefore survive the default change; resetting or an invalid/missing preference chooses System.
- Font families are curated for readability and reliable distribution, not for an exhaustive inventory of locally installed fonts.
- Mouse-wheel sensitivity modifies discrete wheel movement only. Precision touchpad and non-wheel input retain their platform behavior.
- Existing storage ownership and wipe semantics remain authoritative; this slice changes presentation only.
- Tasks, Schedule, and Activity table redesigns (#112 and #113), and task command entry redesign (#110), remain outside S044.
- Exact native hardware evidence can be completed during release-candidate qualification when a specific installer is available; automated and local evidence MUST describe that boundary honestly.
