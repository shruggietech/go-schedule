# Feature Specification: Desktop Visual and Shell Foundations

**Feature Branch**: `codex/083-desktop-foundations`

**Created**: 2026-09-11

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Shared component, Windows-hosted interaction, native Wails build, and canonical eight-gate verification passed on 2026-09-11 on review branch `codex/083-desktop-foundations` for issues [#229](https://github.com/shruggietech/go-schedule/issues/229) and [#230](https://github.com/shruggietech/go-schedule/issues/230); pull-request review remains publication evidence.

**Input**: Restore the shared desktop visual hierarchy, control feedback, theme consistency, fixed navigation shell, and transient status feedback exposed during attended v1.4.0 qualification.

## Clarifications

### Session 2026-09-11

- Q: Which issue boundary belongs in S083? -> A: Complete the shared foundations in #229 and #230; defer page-specific task authoring, Notifications information architecture, and dense administration forms to #231 through #233.
- Q: How compact may controls become? -> A: Reduce visual height and padding by about 25 percent while preserving a minimum 24 by 24 CSS-pixel target or equivalent spacing and localization headroom.
- Q: Which feedback remains persistent? -> A: Routine success feedback is timed and transient; correctable or ongoing problems remain dismissible message bars or contextual field errors.

## User Scenarios & Testing

### User Story 1 - Read and operate a coherent desktop (Priority: P1)

A desktop user can distinguish action priority, hover and keyboard focus, disabled state, and semantic intent while cards and dialogs retain consistent breathing room in light, dark, and follow-system appearances.

**Why this priority**: The current candidate presents unclear and oversized controls across every workflow, so shared primitives must be corrected before page-specific redesign begins.

**Independent Test**: Exercise the shared control catalog and representative task actions under light, dark, and follow-system appearances using pointer, keyboard, automated contrast checks, and focused native Windows observation.

**Acceptance Scenarios**:

1. **Given** any supported appearance, **When** a pointer hovers or presses a button, **Then** the button changes visibly while its label and essential boundary remain readable.
2. **Given** keyboard navigation, **When** focus reaches any action, **Then** a clear focus indicator identifies the action without relying on color alone.
3. **Given** task actions and a confirmation dialog, **When** Enable, Delete, Confirm, and Close are shown, **Then** their semantic roles, compact dimensions, and spacing are immediately distinguishable.
4. **Given** follow-system appearance, **When** the operating-system theme resolves to light or dark, **Then** the application surface and visible brand imagery use a consistent palette.

### User Story 2 - Navigate long pages without losing the shell (Priority: P1)

A desktop user can scroll any active page while application identity, navigation destinations, target context, appearance control, and Exit remain available.

**Why this priority**: Losing navigation and Exit on long pages creates a basic operability failure and compounds every page-level density problem.

**Independent Test**: At 1440 by 900, 900 by 650, and 800 by 600, open every route, force long content, scroll the page region, and confirm the shell remains fixed with no document-level or horizontal overflow.

**Acceptance Scenarios**:

1. **Given** a page taller than the viewport, **When** the user scrolls from top to bottom, **Then** only the active page content moves and all shell controls remain available.
2. **Given** an 800 by 600 viewport or 200 percent zoom, **When** the interface reflows, **Then** navigation, Exit, target context, and Appearance remain reachable without horizontal document scrolling.

### User Story 3 - Receive useful feedback without layout disruption (Priority: P1)

A desktop user receives one clear action result at a time without permanent status text, obscured controls, or page geometry changes.

**Why this priority**: Current success and error strings accumulate at the viewport edge or become oversized page content, producing both visual and accessibility failures.

**Independent Test**: Trigger repeated success, warning, and error outcomes; verify transient messages replace one another, pause while hovered or focused, dismiss automatically or explicitly as appropriate, and preserve page geometry.

**Acceptance Scenarios**:

1. **Given** a routine successful action, **When** its confirmation appears, **Then** one bounded toast is announced, overlays rather than shifts content, and clears automatically.
2. **Given** repeated routine outcomes, **When** a new outcome arrives, **Then** it replaces the prior toast rather than concatenating into a permanent line.
3. **Given** a correctable or ongoing problem, **When** it is presented, **Then** it appears as a contextual field error or dismissible message bar with enough information to recover.
4. **Given** a transient message under hover, keyboard focus, or reduced-motion preference, **When** dismissal timing applies, **Then** the user has sufficient time to read and interact with it without motion-dependent behavior.

### Edge Cases

- The saved appearance is `system` and the operating system changes palette while the application is open.
- A message changes while the previous toast is hovered or focused.
- An empty message arrives after a visible toast and must clear it without leaving an empty live-region container visible.
- A narrow or zoomed viewport cannot accommodate the full horizontal navigation treatment.
- A dialog contains multiple actions, long localized labels, or a scrollable body.
- A disabled semantic action must remain identifiable while clearly unavailable.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST use one documented token scale for surfaces, text, spacing, control dimensions, radii, borders, elevation, focus indicators, and semantic action colors.
- **FR-002**: Shared buttons MUST expose primary, secondary, subtle, affirmative, and destructive variants with distinct rest, hover, focus-visible, pressed, selected where applicable, disabled, and pending states.
- **FR-003**: Button labels MUST reach a 4.5:1 contrast ratio and essential boundaries, icons, and focus indicators MUST reach 3:1 in every supported appearance and state.
- **FR-004**: Ordinary desktop buttons MUST be materially more compact than the v1.4.0 candidate while preserving a minimum 24 by 24 CSS-pixel target or equivalent target spacing and room for localized labels.
- **FR-005**: Cards and dialogs MUST use the shared spacing scale for deliberate top, bottom, and side padding.
- **FR-006**: Dialog actions MUST have consistent dimensions, at least 8 CSS pixels of separation, a clear primary or destructive action hierarchy, Escape dismissal, focus containment, and focus return to the exact invoker.
- **FR-007**: Task Enable and Delete actions MUST use affirmative and destructive variants respectively and remain distinguishable through text, boundary, and state treatment in light and dark palettes.
- **FR-008**: Follow-system appearance MUST resolve the same palette for application surfaces and visible brand imagery and MUST respond when the operating-system palette changes.
- **FR-009**: The application shell MUST own the viewport and MUST constrain vertical scrolling to one active-page content region.
- **FR-010**: Application identity, navigation destinations, target context, Appearance, and Exit MUST remain reachable at every supported viewport and content length.
- **FR-011**: The active page MUST reflow at 800 by 600 and 200 percent zoom without document-level horizontal scrolling or loss of shell functionality.
- **FR-012**: Routine successful action feedback MUST appear in a bounded overlay live region, replace prior routine feedback, avoid page geometry changes, and dismiss automatically after 5 seconds.
- **FR-013**: Automatic toast dismissal MUST pause while the toast is hovered or contains keyboard focus and MUST resume with sufficient remaining reading time afterward.
- **FR-014**: The feedback layer MUST maintain safe viewport margins and MUST NOT overlap the Appearance control, navigation, window boundary, or currently focused control.
- **FR-015**: Correctable or ongoing errors MUST use contextual field errors or dismissible persistent message bars; they MUST NOT become undismissable oversized page sections.
- **FR-016**: Status feedback MUST provide appropriate polite or assertive live-region semantics and explicit dismissal controls without depending on animation.
- **FR-017**: Component, integration, accessibility, reflow, and native Windows checks MUST cover the shared states and shell behavior introduced by this slice.
- **FR-018**: S083 MUST NOT redesign task creation, Notifications information architecture, Agent Access, Connections, or Settings page composition beyond consuming the shared foundations.

### Key Entities

- **Appearance Resolution**: The saved light, dark, or follow-system preference and the currently resolved visual palette.
- **Control Variant**: A shared action role with semantic priority, interaction states, dimensions, and accessibility rules.
- **Application Shell**: The persistent application identity, navigation, target context, page viewport, Appearance control, and Exit action.
- **Feedback Item**: One transient or persistent user-facing outcome with tone, announcement priority, dismissal behavior, and replacement identity.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of shared button variants have readable rest, hover, focus-visible, pressed, disabled, and pending treatments in light, dark, and follow-system appearances.
- **SC-002**: Representative button label contrast is at least 4.5:1 and essential visual contrast is at least 3:1 in every tested state and palette.
- **SC-003**: Ordinary button height and horizontal padding are at least 20 percent smaller than the v1.4.0 candidate while every target remains at least 24 by 24 CSS pixels or has equivalent target spacing.
- **SC-004**: Every dialog action group preserves at least 8 CSS pixels between adjacent actions.
- **SC-005**: All eight primary routes retain navigation and Exit from top to bottom at 1440 by 900, 900 by 650, and 800 by 600, with no document-level horizontal overflow.
- **SC-006**: Routine feedback changes page geometry by 0 CSS pixels, shows no more than one message at a time, and clears 5 seconds after the latest message unless paused by hover or focus.
- **SC-007**: Automated accessibility scans report zero serious or critical findings for the representative light, dark, and follow-system states.
- **SC-008**: Focused native Windows evidence demonstrates pointer hover, keyboard focus, theme resolution, fixed shell scrolling, and transient feedback against one exact S083 build.
- **SC-009**: All canonical verification gates pass with zero S083 exclusions.

## Assumptions

- The desktop remains a keyboard-and-pointer Windows application with an 800 by 600 minimum supported viewport.
- The existing locally bundled fonts and brand mark remain available; theme consistency may use a palette-aware treatment of the current neutral mark rather than introduce a new brand asset.
- Five seconds is sufficient for routine status text because persistent or actionable errors use a different feedback class and because hover or focus pauses dismissal.
- The shared primitives may correct page surfaces that consume them, but page-specific information architecture and form composition remain assigned to #231 through #233.

## Dependencies

- Parent: [#228](https://github.com/shruggietech/go-schedule/issues/228).
- Completes [#229](https://github.com/shruggietech/go-schedule/issues/229) and [#230](https://github.com/shruggietech/go-schedule/issues/230).
- Enables [#231](https://github.com/shruggietech/go-schedule/issues/231), [#232](https://github.com/shruggietech/go-schedule/issues/232), and [#233](https://github.com/shruggietech/go-schedule/issues/233) to consume stable shared foundations.
- Blocks public v1.4.0 promotion tracked by [#226](https://github.com/shruggietech/go-schedule/issues/226).

## Scope Boundaries

**In scope**: Shared design tokens and controls, cards and dialogs, theme-aware imagery, viewport shell ownership, active-page scrolling, responsive shell reflow, transient toasts, dismissible message bars, accessibility coverage, and focused native Windows evidence.

**Out of scope**: Task editor modal conversion, task-form density, Notifications navigation or disclosure redesign, Agent Access layout, Connections pairing layout, Settings path presentation, release tagging, release drafting, release promotion, and unrelated daemon or CLI behavior.
