# Feature Specification: Structured Desktop Data Tables

**Feature Branch**: `codex/045-structured-data-tables`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Headless GUI, focused race, full repository, and canonical eight-gate verification passed 2026-09-03 on `codex/045-structured-data-tables`; exact-candidate native Windows evidence remains assigned to release qualification for #112/#113

**Input**: User description: "S045 bundles GitHub issues #112 and #113: replace ambiguous Tasks rows with a labeled accessible table, and tabulate and color-code Schedule and Activity event rows."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand and operate on tasks (Priority: P1)

As a scheduler operator, I can scan a labeled Tasks table and distinguish whether a task is enabled from where it is in its lifecycle, without opening every task or interpreting bracketed shorthand.

**Why this priority**: Tasks are the primary management surface. The current unlabeled text can show `active` beside `disabled`, which appears contradictory and undermines confidence in routine actions.

**Independent Test**: Populate the Tasks view with enabled, disabled, completed, grouped, ungrouped, and long-named tasks; confirm every visible value is labeled, the two status dimensions are unambiguous, and selection, toolbar actions, keyboard navigation, refresh, and double-click editing still target the correct task.

**Acceptance Scenarios**:

1. **Given** tasks with different enabled and lifecycle states, **When** the operator scans the table, **Then** the headers and plain-language cell values make both state dimensions independently understandable.
2. **Given** an enabled task whose lifecycle is active and a disabled task whose stored lifecycle differs, **When** both rows are visible, **Then** neither row uses unexplained brackets or appears self-contradictory.
3. **Given** a selected task, **When** live data refreshes or the operator uses a toolbar action, keyboard selection, or double-click, **Then** the same task identity remains the target whenever that task still exists.
4. **Given** enough tasks to scroll, **When** the operator scrolls the rows, **Then** the column headers remain visible.

---

### User Story 2 - Read the schedule as structured events (Priority: P2)

As a scheduler operator, I can scan the Schedule list under stable column headers and understand whether each occurrence is upcoming or completed and, for completed occurrences, its outcome.

**Why this priority**: The current marker, timestamp, and task name are concatenated into an unlabeled line, so meaning depends on glyph recognition and spacing.

**Independent Test**: Populate the Schedule list with future occurrences and past success, failure, skipped, caught-up, queued, and unknown outcomes; confirm stable labels, semantic text and glyphs, chronological ordering, live refresh, calendar switching, and window-range selection.

**Acceptance Scenarios**:

1. **Given** upcoming and past occurrences, **When** the operator views Schedule in List mode, **Then** time, task, event type, and outcome information appear under explicit headers.
2. **Given** a past occurrence with an outcome, **When** its row is shown, **Then** both a text label and a glyph identify the meaning and share a restrained semantic color.
3. **Given** a future occurrence, **When** its row is shown, **Then** it is identified as scheduled without implying that it has already succeeded.
4. **Given** the operator switches between List and Calendar or changes the time window, **When** refreshed results arrive, **Then** the existing controls and chronological ordering remain intact.

---

### User Story 3 - Triage activity consistently (Priority: P2)

As a scheduler operator, I can scan Activity rows under stable column headers, recognize normalized event severity from both text and glyph, and open the full detail for the intended event.

**Why this priority**: Activity is the diagnostic surface. Inconsistent casing and unaligned fields slow triage, especially when messages and sources vary in length.

**Independent Test**: Populate Activity with informational, warning, error, and unknown-severity log records and alerts; confirm normalized labels, semantic text and glyph treatment, newest-first ordering, filtering, clearing, live refresh, and detail opening.

**Acceptance Scenarios**:

1. **Given** informational, warning, and error events, **When** the operator scans Activity, **Then** time, severity, source, and summary are aligned under explicit headers.
2. **Given** any supported severity, **When** its row is shown, **Then** the severity label uses the documented uppercase form and the label plus glyph communicate meaning without relying on color alone.
3. **Given** a selected Activity row, **When** the operator activates it, **Then** the existing detail view opens for exactly that event and exposes the full values.
4. **Given** an active severity filter or cleared-view cutoff, **When** live data refreshes, **Then** those behaviors remain unchanged while the visible rows retain the structured presentation.

### Edge Cases

- Empty Tasks, Schedule, and Activity views retain visible headers and an understandable empty body without placeholder rows being mistaken for data.
- Missing task groups, time zones, event sources, outcomes, and severities use one documented fallback value rather than leaving ambiguous blank columns.
- Unknown future lifecycle, outcome, or severity values remain readable as normalized text and use a neutral visual treatment rather than being misclassified.
- Long task names, group names, sources, and messages do not force horizontal scrolling; the complete value remains discoverable through selection, detail, or an equivalent accessible disclosure.
- Narrow window sizes preserve the highest-priority columns and meaning; lower-priority content may be compacted or disclosed elsewhere but must not overlap or become unreadable.
- A selected row that disappears during live refresh clears safely; a selected row that remains is still associated with its stable identity even if its visual index changes.
- Rows remain distinguishable in both supported palettes and when the operating system or application font changes.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Tasks, Schedule List, and Activity views MUST present row data under explicit, persistent column headers.
- **FR-002**: Headers MUST remain visible while the corresponding row body scrolls vertically.
- **FR-003**: The structured views MUST fit within the application's supported content width without requiring horizontal scrolling.
- **FR-004**: Every view MUST define a deterministic responsive priority for columns so narrow layouts preserve identity and meaning before secondary detail.
- **FR-005**: Truncated, compacted, or otherwise hidden cell content MUST remain discoverable without editing underlying data.
- **FR-006**: Empty views MUST retain their headers and MUST NOT display a template row as actual data.
- **FR-007**: Missing and unknown values MUST use documented, plain-language fallbacks and a neutral semantic treatment.
- **FR-008**: Alternating row backgrounds MUST be subtle in dark, light, and follow-system modes.
- **FR-009**: Hover, keyboard focus, and selection MUST each be distinguishable, MUST preserve readable text, and MUST preserve selected identity when states overlap.
- **FR-010**: Row state MUST never be conveyed by color alone; text, glyph shape, borders, or another non-color cue MUST remain available.
- **FR-011**: Live refresh MUST preserve selection by stable record identity when the record remains and MUST clear selection safely when it disappears.
- **FR-012**: The Tasks table MUST label and display task name, enabled status, lifecycle state, timezone, and group membership.
- **FR-013**: Tasks MUST express enabled status and lifecycle state as separate, plain-language concepts whose labels and values cannot reasonably be read as one contradictory status.
- **FR-014**: Tasks MUST NOT use unexplained square-bracket decoration around lifecycle or enabled values.
- **FR-015**: Tasks MUST preserve single-click and keyboard selection, all existing toolbar actions, and double-click-to-edit for the correct stable task identity.
- **FR-016**: Task ordering MUST remain deterministic across refreshes, and a retained task selection MUST follow identity rather than row index.
- **FR-017**: Schedule List rows MUST label and display occurrence time, task name, event type, and outcome where an outcome exists.
- **FR-018**: Schedule MUST distinguish upcoming occurrences from completed occurrences and MUST NOT represent an upcoming occurrence as a completed success.
- **FR-019**: Schedule outcome labels and glyphs MUST cover success, failure, skipped, caught up, queued, missing outcome, and unknown values using a documented normalization scheme.
- **FR-020**: Schedule MUST preserve chronological ordering, live refresh, range selection, and switching between List and Calendar views.
- **FR-021**: Activity rows MUST label and display event time, severity, source, and summary.
- **FR-022**: Activity severity labels MUST use the documented uppercase forms `INFO`, `WARNING`, and `ERROR`; unsupported values MUST be normalized to a readable uppercase label without being assigned a false severity.
- **FR-023**: Activity severity text and its leading glyph MUST share the intended restrained semantic color while retaining independent textual or shape meaning.
- **FR-024**: Activity MUST preserve newest-first ordering, severity filtering, clearing and alert acknowledgement behavior, live refresh, and opening the full detail for the correct event.
- **FR-025**: Equivalent row concepts across the three views MUST use consistent casing, alignment, glyph conventions, semantic color roles, and fallback language.
- **FR-026**: Normal-sized text MUST maintain at least 4.5:1 contrast in rest, hover, focus, and selected states; essential non-text indicators MUST maintain at least 3:1 against adjacent colors.
- **FR-027**: Structured row presentation MUST remain usable with every supported interface font and appearance mode.
- **FR-028**: Automated verification MUST cover row mapping, normalization, fallbacks, responsive column priority, selection stability, activation behavior, and semantic contrast; native Windows evidence MUST cover populated dark and light views.

### Key Entities

- **Structured Row**: A stable record identity plus ordered display cells, semantic presentation metadata, and a route to complete values when compacted.
- **Column Definition**: A header label, content role, alignment, responsive priority, sizing rule, and full-value disclosure rule shared by a structured view.
- **Task Row**: A structured representation of task identity, name, enabled status, lifecycle state, timezone, and group.
- **Schedule Row**: A structured representation of occurrence identity, time, task, event type, outcome, glyph, and semantic role.
- **Activity Row**: A structured representation of event identity, time, severity, source, summary, glyph, semantic role, and full detail.
- **Semantic Role**: A normalized meaning such as informational, scheduled, success, warning, error, disabled, or neutral unknown, expressed with both non-color cues and restrained palette values.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In a representative mixed-state Tasks view, every displayed value can be assigned to the correct meaning from visible headers alone, with zero unexplained bracketed status tokens.
- **SC-002**: All three structured views retain visible headers while scrolling at least 100 populated rows.
- **SC-003**: At the default launch size and the smallest application-supported window size, all three views remain operable without horizontal scrolling, overlapping text, or loss of record identity and primary meaning.
- **SC-004**: Across dark, light, and follow-system appearances, every tested text/state combination meets the 4.5:1 contrast threshold and every essential non-text indicator meets 3:1.
- **SC-005**: One hundred percent of supported Schedule outcomes and Activity severities map to documented labels, glyphs, and semantic roles; unknown values map to a neutral readable fallback.
- **SC-006**: In automated refresh scenarios that reorder, update, or remove rows, retained selections target the same stable identity and removed selections produce no unintended action.
- **SC-007**: Existing task actions, task double-click editing, Schedule controls, Activity filters, clearing, live refresh, and Activity detail opening pass without behavioral regression.
- **SC-008**: Headless tests exercise empty, populated, long-value, unknown-value, responsive, selection, and activation cases, and a native Windows walkthrough records populated dark and light evidence before issue closure.

## Assumptions

- This slice changes desktop presentation and interaction only; daemon, API, persistence, scheduling, and execution contracts remain unchanged.
- Current backend ordering remains authoritative: Schedule is chronological and Activity is newest-first. S045 does not add user-selectable sorting.
- Schedule Calendar mode retains its existing calendar presentation; the structured table applies to Schedule List mode.
- Existing Activity detail is the authoritative full-value disclosure for Activity rows. Tasks and Schedule may use an equivalent non-mutating disclosure suited to their interaction model.
- Semantic colors reuse the application theme roles established in S044 rather than introducing high-contrast bespoke palettes.
- Native Windows visual evidence is a release-qualification responsibility. This implementation supplies headless evidence and does not claim attended acceptance on behalf of the operator.
- GitHub issues #112 and #113 remain individually traceable and close only when their respective acceptance criteria, including required native evidence, are complete.

## Scope Boundaries

### In Scope

- GitHub issues #112 and #113.
- Shared structured-row and responsive-column behavior needed by Tasks, Schedule List, and Activity.
- Presentation-specific tests, documentation, and qualification guidance.

### Out of Scope

- Task command-line entry redesign (#110).
- Richer task-run failure capture and diagnostics (#102).
- New daemon/API fields, persistence migrations, scheduling semantics, or execution behavior.
- New user-selectable sorting, column customization, exporting, or global search.
- Redesign of Schedule Calendar mode, Groups, or Chains.
