# Feature Specification: Responsive Task and Administration Workflows

**Feature Branch**: `codex/084-responsive-workflows`

**Created**: 2026-09-11

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Implemented on review branch `codex/084-responsive-workflows`; focused frontend evidence and canonical verification are recorded in `verification.md`.

**Input**: Replace the oversized inline task editor with a focused task workflow and recompose Agent Access, Connections, and Settings into readable responsive administration surfaces, completing issues #231 and #233.

## Clarifications

### Session 2026-09-11

- Q: Which open desktop issues belong in S084? -> A: Complete #231 and #233 together because both consume the same responsive form and information-display primitives; defer Notifications information architecture in #232.
- Q: How should task creation and editing be presented? -> A: Use the same bounded modal workflow for create and edit, with common fields first, advanced settings collapsed, contained scrolling, and focus returned to the exact invoker.
- Q: Which administration controls should be progressively disclosed? -> A: Keep diagnosis and common actions visible, while inactive localhost configuration and remote pairing forms remain collapsed until explicitly opened.
- Q: How should Copy path progress and completion be shown? -> A: Only the activated path action changes state, identifies pending work, and reports completion without disabling or visually changing unrelated controls.

## User Scenarios & Testing

### User Story 1 - Create or edit a task without losing context (Priority: P1)

A desktop user opens one focused task workflow, completes the common fields, optionally expands advanced settings, previews or saves the task, and returns to the same Tasks context.

**Why this priority**: Task authoring is the product's primary workflow, and the current inline editor makes the page excessively tall while displacing the task list and navigation context.

**Independent Test**: Open task creation and editing with pointer and keyboard at 1280 by 800 and 800 by 600, complete the common path, expand advanced settings, preview, cancel, save, and confirm focus and page geometry remain stable.

**Acceptance Scenarios**:

1. **Given** the Tasks page, **When** the user invokes Create task, **Then** a bounded modal opens without inserting content into or changing the geometry of the underlying page.
2. **Given** a task is selected, **When** the user invokes Edit, **Then** the same focused workflow opens with stored values and returns focus to Edit when closed.
3. **Given** a new task workflow, **When** the user reads the safe example, **Then** concise help and an unmistakable Insert example action are available without submitting the task.
4. **Given** an open task workflow, **When** the user presses Escape or chooses Cancel, **Then** no draft is saved and focus returns to the exact invoking action.
5. **Given** a task deletion confirmation, **When** the user navigates its actions by keyboard, **Then** destructive and safe actions remain separated, clearly named, and focus returns after dismissal.

### User Story 2 - Understand and configure agent access (Priority: P1)

A desktop user first sees the current local agent-access boundary and available authority, then deliberately expands optional network configuration only when needed.

**Why this priority**: The existing form collides at constrained widths and exposes unfamiliar configuration before establishing the safer default state.

**Independent Test**: Inspect disabled and active localhost states, expand configuration, traverse every field in visual order, and verify clean stacking at 800 by 600 and 200 percent zoom.

**Acceptance Scenarios**:

1. **Given** localhost access is off, **When** Agent Access opens, **Then** stdio and authority status appear before a collapsed optional-network configuration.
2. **Given** the user expands localhost configuration, **When** the viewport narrows, **Then** every label remains above its control and fields stack before collision.
3. **Given** localhost access is active, **When** the page opens, **Then** endpoint and credential evidence are readable and lifecycle actions remain grouped with the active status.

### User Story 3 - Diagnose and pair connections without cramped forms (Priority: P1)

A desktop user reads the current connection diagnosis and saved profiles before opening the specialist remote-pairing form, with every value and action aligned at supported sizes.

**Why this priority**: Connection recovery is operationally sensitive, and the current irregular grids obscure field relationships and produce collisions.

**Independent Test**: Exercise connected, recovery, saved-profile, and pairing states at supported viewport and zoom settings while checking reading order, focus visibility, and value wrapping.

**Acceptance Scenarios**:

1. **Given** any connection state, **When** Connections opens, **Then** the current diagnosis, status, and next action appear before saved targets and specialist pairing controls.
2. **Given** saved connection profiles, **When** values are long, **Then** identifiers, endpoints, and fingerprints wrap safely without colliding with labels or actions.
3. **Given** no active pairing task, **When** Connections opens, **Then** remote pairing remains collapsed but clearly discoverable.
4. **Given** pairing is expanded, **When** the user tabs through it, **Then** focus order follows the visible top-to-bottom field order and remains visible.

### User Story 4 - Read and copy local paths without interface-wide feedback (Priority: P1)

A desktop user can scan storage ownership, read complete filesystem paths in an appropriate code style, and copy one exact path without unrelated buttons changing state.

**Why this priority**: Oversized proportional path text and interface-wide pending feedback make Settings difficult to scan and visually unstable.

**Independent Test**: Load records containing long Windows and Unix paths, copy each copyable record, and verify typography, complete-value access, isolated pending or success feedback, keyboard order, and theme behavior.

**Acceptance Scenarios**:

1. **Given** a long filesystem path, **When** Settings opens, **Then** the complete value uses compact monospace text, deliberate spacing, and safe wrapping within its card.
2. **Given** multiple actions on Settings, **When** one Copy path action runs, **Then** only that action shows pending and completion state while all unrelated actions remain visually stable.
3. **Given** light, dark, or follow-system appearance, **When** Settings reflows, **Then** paths, ownership descriptions, badges, and actions remain readable without horizontal scrolling.

### Edge Cases

- Task validation moves focus to a field inside a scrollable dialog that is currently outside the visible dialog body.
- Advanced task settings contain enough fields to exceed the available window height.
- A task save completes while its modal is open, and the original invoker remains mounted or disappears because the task list refreshes.
- An agent-access origin, connection endpoint, fingerprint, or filesystem path is longer than one line and contains no convenient whitespace.
- A copy operation succeeds or fails after the user moves to another record or starts another independent settings action.
- A connection profile list is empty, the daemon is unavailable, or a pairing repair workflow targets an existing profile.
- At 200 percent zoom, a nominal 800 by 600 window has only a narrow single-column content region.

## Requirements

### Functional Requirements

- **FR-001**: Create task and Edit MUST open the same bounded modal workflow without adding an editor card to the Tasks page document flow.
- **FR-002**: The task modal MUST contain its own vertical scrolling when necessary, avoid horizontal scrolling, retain its heading and actions, and fit within the supported viewport.
- **FR-003**: The task modal MUST contain keyboard focus, close through Escape and Cancel, and return focus to the exact Create task or Edit invoker whenever that invoker remains available.
- **FR-004**: Common task fields MUST precede an Advanced settings disclosure that is collapsed when the workflow opens.
- **FR-005**: The safe command example MUST use concise help and a visually explicit secondary Insert example action that changes only the command field and never submits the task.
- **FR-006**: Task fields and selectors MUST share the same control dimensions, label placement, responsive grid, focus treatment, and error presentation.
- **FR-007**: Task preview and validation results MUST remain associated with the modal, bring invalid fields into view, and preserve entered values.
- **FR-008**: Task deletion confirmation MUST use the shared dialog contract with separated destructive and safe actions, Escape dismissal, pending protection, and focus restoration.
- **FR-009**: The desktop MUST provide reusable responsive form-grid, description-list, path-display, status-label, and card-section presentation primitives for the administration pages in this slice.
- **FR-010**: Agent Access MUST present stdio and authority status before optional localhost HTTP configuration and MUST keep inactive network configuration collapsed by default.
- **FR-011**: Active localhost HTTP evidence and lifecycle actions MUST remain visible, aligned, and readable without exposing secret credential values.
- **FR-012**: Connections MUST present current diagnosis and recovery actions before saved profiles and remote-pairing configuration.
- **FR-013**: Remote pairing MUST remain collapsed by default except when repairing a selected profile and MUST use the shared responsive form composition when expanded.
- **FR-014**: Administration labels MUST appear above controls, related fields MUST align consistently, and multi-column regions MUST stack before content collides.
- **FR-015**: Long endpoints, identifiers, fingerprints, origins, and paths MUST wrap safely within their region without horizontal page scrolling or loss of complete-value access.
- **FR-016**: Every filesystem path MUST use smaller monospace text with deliberate upper spacing and retain the complete literal value in selectable text.
- **FR-017**: A Copy path operation MUST update only its exact action's pending and completion presentation and MUST NOT disable, flicker, or relabel unrelated buttons.
- **FR-018**: Copy path actions MUST preserve exact record identity when operations overlap or finish out of order.
- **FR-019**: Visual reading order and keyboard focus order MUST match across task, agent-access, connection, pairing, and settings forms.
- **FR-020**: The completed workflows MUST remain usable at 1280 by 800, 800 by 600, 200 percent zoom, supported scaled DPI, and light, dark, and follow-system appearances.
- **FR-021**: Automated component, interaction, accessibility, and responsive-layout tests plus a native Windows build and focused walkthrough MUST cover the behavior introduced by S084.
- **FR-022**: S084 MUST NOT redesign Notifications, change daemon or network contracts, expose secrets, mutate persisted schemas, or alter release tags and hosted release state.

### Key Entities

- **Focused Task Workflow**: The create or edit mode, current task draft, validation and preview result, pending action, scroll position, and exact invoking control.
- **Administration Section**: A titled status, detail, configuration, or action grouping whose reading order and responsive behavior remain stable.
- **Definition Item**: A label and complete value pair for operational metadata, including values that require safe wrapping or code typography.
- **Path Presentation**: A path label, complete literal path, availability, copyability, and record-scoped copy state.
- **Disclosure State**: The explicit collapsed or expanded state of advanced task, localhost HTTP, or remote-pairing configuration.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Opening task creation or editing changes the underlying Tasks page height and scroll position by 0 CSS pixels.
- **SC-002**: One hundred percent of task fields, dialog actions, and administration controls remain reachable by keyboard at 800 by 600 and 200 percent zoom without document-level horizontal scrolling.
- **SC-003**: Task dialogs maintain at least 16 CSS pixels from each viewport edge and at least 8 CSS pixels between adjacent actions at every tested size.
- **SC-004**: Advanced task settings, inactive localhost HTTP configuration, and ordinary remote pairing are collapsed on initial presentation and require exactly one explicit disclosure action to open.
- **SC-005**: One hundred percent of long operational values remain fully available and contained within their card at widths down to the supported minimum.
- **SC-006**: During a Copy path operation, zero unrelated controls change disabled state, label, size, or visual styling.
- **SC-007**: Automated accessibility scans report zero serious or critical findings across the four workflows in light, dark, and follow-system appearances.
- **SC-008**: Focused native Windows evidence demonstrates task create and cancel, task edit, delete confirmation, Agent Access disclosure, Connections pairing disclosure, long Settings paths, and isolated copy feedback against one exact S084 build.
- **SC-009**: All canonical verification gates pass with zero S084 exclusions.

## Assumptions

- S083's shared control, dialog, theme, shell, and feedback contracts are the required foundation and will be extended rather than replaced.
- Existing task, agent-access, connection, pairing, and settings bridges remain authoritative and require no backend contract changes.
- Complete-value access means literal selectable text with safe wrapping; truncation-only presentation is insufficient for operational identifiers and paths.
- A repair request is an active pairing task and may open its pairing disclosure automatically so the requested action is not hidden.
- The desktop remains a keyboard-and-pointer application with an 800 by 600 minimum supported viewport.

## Dependencies

- Parent: [#228](https://github.com/shruggietech/go-schedule/issues/228).
- Completes [#231](https://github.com/shruggietech/go-schedule/issues/231) and [#233](https://github.com/shruggietech/go-schedule/issues/233).
- Depends on completed shared foundations [#229](https://github.com/shruggietech/go-schedule/issues/229) and [#230](https://github.com/shruggietech/go-schedule/issues/230).
- Leaves [#232](https://github.com/shruggietech/go-schedule/issues/232) open for a dedicated Notifications information-architecture slice.
- Blocks public v1.4.0 promotion tracked by [#226](https://github.com/shruggietech/go-schedule/issues/226).

## Scope Boundaries

**In scope**: Focused task create and edit modal, task form hierarchy and deletion confirmation, reusable administration presentation primitives, Agent Access composition, Connections diagnosis and pairing composition, Settings path presentation, record-scoped copy feedback, responsive and accessibility tests, and focused native Windows evidence.

**Out of scope**: Notifications redesign, new notification capabilities, daemon or API changes, network authorization changes, persisted schema changes, release restaging, tag mutation, release publication, and unrelated desktop pages.
