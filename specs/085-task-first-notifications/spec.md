# Feature Specification: Task-First Notifications

**Feature Branch**: `codex/085-task-first-notifications`

**Created**: 2026-09-11

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Implemented and verified in `verification.md`

**Input**: Make Notifications immediately answer whether notifications are active and what recently happened, while moving webhook and policy administration behind clear progressive disclosure, completing issue #232.

## Clarifications

### Session 2026-09-11

- Q: What belongs on the initial Notifications surface? -> A: A plain-language status overview, configured-scope summary, recent task outcomes, and the next useful action appear before specialist setup.
- Q: How should specialist administration remain reachable? -> A: Separate native disclosures for destinations, assignment rules, and delivery diagnostics remain collapsed by default and open with one explicit action.
- Q: How should delivery states guide users? -> A: Disabled, healthy, queued, retrying, sending, successful, and failed states each use plain language, readable text independent of color, and a specific next action when action is useful.
- Q: Does S085 change delivery behavior or secret storage? -> A: No. S085 reorganizes and explains the existing capability while retaining current redaction, write-only secret replacement, policy, test, retry, and refresh behavior.

## User Scenarios & Testing

### User Story 1 - Understand notification status immediately (Priority: P1)

A desktop user opens Notifications and can immediately tell whether any destination is active, whether tasks or groups can notify, what recently happened, and what to do next without understanding webhook terminology.

**Why this priority**: The destination currently leads with specialist configuration and fails its basic information promise.

**Independent Test**: Open Notifications with empty, disabled, healthy, in-progress, and failed fixtures and verify the first visible region communicates status, recent task outcomes, and the next action at 800 by 600.

**Acceptance Scenarios**:

1. **Given** no configured destination, **When** Notifications opens, **Then** the overview explains that notifications are off and offers one direct path to setup.
2. **Given** enabled destinations and configured scopes, **When** Notifications opens, **Then** the overview summarizes active destinations and the task or group coverage before any advanced form.
3. **Given** recent deliveries, **When** Notifications opens, **Then** recent task outcomes are visible before configuration and each state explains whether action is needed.
4. **Given** a failed or retrying delivery, **When** the user reads it, **Then** plain-language guidance identifies the next useful action without requiring diagnostic terminology.

### User Story 2 - Configure notifications deliberately (Priority: P1)

A user who needs setup expands the relevant destination or assignment section, receives concise definitions in context, completes the existing operation, and returns to the overview without unrelated administration dominating the page.

**Why this priority**: Setup remains essential, but it should appear only after the user deliberately selects that task.

**Independent Test**: Expand destination and assignment sections by pointer and keyboard, create or edit a channel, load a task or group policy, save it, collapse the section, and confirm secret handling and focus order remain safe.

**Acceptance Scenarios**:

1. **Given** the initial page, **When** the user has not requested setup, **Then** destination and assignment forms remain collapsed.
2. **Given** a setup need, **When** the user expands Destinations or Assignment rules, **Then** its specialist terms have concise inline explanations associated with the relevant control or section.
3. **Given** stored authorization and endpoint values, **When** a destination is edited, **Then** secrets remain write-only and replacement requires explicit intent.
4. **Given** an inherited policy, **When** assignment rules are opened, **Then** inheritance and override consequences are explained in plain language.

### User Story 3 - Investigate delivery evidence on demand (Priority: P1)

A user can scan recent outcomes first and expand detailed delivery diagnostics only when a result needs investigation.

**Why this priority**: Recent outcomes are useful to everyone, while correlation identifiers and transport details serve a narrower diagnostic task.

**Independent Test**: Filter and select test and task-outcome deliveries across all states, open detailed diagnostics, and confirm redacted evidence remains accessible without making the initial page dense.

**Acceptance Scenarios**:

1. **Given** recent deliveries, **When** Notifications opens, **Then** a bounded recent-results summary is visible and detailed filters and evidence remain collapsed.
2. **Given** a selected result, **When** diagnostics are expanded, **Then** redacted channel, destination, attempts, timestamps, failure detail, and correlation identifiers are available.
3. **Given** no recent deliveries, **When** Notifications opens, **Then** the empty state explains what notifications do, why no results appear, and the first useful setup or execution action.

### Edge Cases

- All destinations exist but are disabled.
- Configured scopes refer to a removed or disabled destination.
- The latest result is a transport test rather than a task outcome.
- A delivery transitions from queued to retrying, sending, successful, or failed while the page is open.
- A selected diagnostic record leaves the latest complete snapshot or active filters.
- There are 200 delivery records and long task, group, endpoint, error, or correlation values.
- The scheduler disconnects while the user is reading overview data or editing advanced configuration.
- The page is used at 800 by 600, 200 percent zoom, or entirely by keyboard.

## Requirements

### Functional Requirements

- **FR-001**: Notifications MUST lead with a plain-language overview before any destination, policy, testing, or diagnostic form.
- **FR-002**: The overview MUST report enabled destination count, configured task and group coverage, recent task-outcome count, and an overall state derived from the latest complete snapshot.
- **FR-003**: Overall state MUST distinguish setup needed, disabled, healthy, work in progress, retrying, and failed conditions with readable text independent of color.
- **FR-004**: Each overview state MUST identify the next useful action, or explicitly state that no action is needed.
- **FR-005**: A bounded recent-results list MUST appear on the initial surface and prioritize task outcomes over transport tests while retaining access to both.
- **FR-006**: Recent-result rows MUST name the task or identify a test, name the destination, expose the state and time, and provide plain-language state guidance.
- **FR-007**: Destination administration, assignment rules, and delivery diagnostics MUST be separate, clearly named disclosures that are collapsed on initial presentation.
- **FR-008**: Each advanced disclosure MUST open through one keyboard-operable action, expose its expanded state, and retain logical document and focus order.
- **FR-009**: The destination section MUST retain create, edit, enable, disable, test, and remove capabilities without redisplaying endpoints or authorization secrets.
- **FR-010**: The assignment section MUST retain task and group selection, direct assignment, inherited policy explanation, success-volume warning, override clearing, and save behavior.
- **FR-011**: The diagnostics section MUST retain channel and state filtering, selected-result detail, missing-result handling, and redacted operational evidence.
- **FR-012**: Channel, destination, assignment rule, policy, inheritance, test notification, task notification, retry, and failure concepts MUST have concise contextual definitions associated with their relevant section or control.
- **FR-013**: Empty, disconnected, disabled, healthy, queued, retrying, sending, successful, and failed presentations MUST explain current meaning and the next useful action in plain language.
- **FR-014**: Initial and advanced surfaces MUST avoid document-level horizontal scrolling and remain keyboard usable at 800 by 600 and 200 percent zoom.
- **FR-015**: Long values MUST wrap within their region and detailed history MUST remain bounded rather than forcing the entire page to grow with all 200 records.
- **FR-016**: Automated component, service, browser, responsive, keyboard, and accessibility tests plus a native production build MUST cover the S085 behavior.
- **FR-017**: S085 MUST NOT change notification dispatch semantics, persistence schemas, authorization boundaries, daemon APIs, release tags, or hosted release state.

### Key Entities

- **Notification Overview**: Snapshot-derived counts, overall state, current guidance, configured scope summaries, and the bounded recent-results set.
- **Configured Scope Summary**: A task or group with a direct or effective assignment and the enabled or unavailable destinations it can use.
- **Recent Result**: A task outcome or transport test with state, time, destination, task context, and plain-language guidance.
- **Advanced Section**: One collapsed destination, assignment, or diagnostic task with an explicit summary and local controls.
- **Delivery Evidence**: Existing redacted transport and correlation detail shown only on diagnostic request.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A first-time user can identify whether notifications are active, whether action is required, and the latest task outcome from the initial 800 by 600 viewport without opening an advanced section.
- **SC-002**: Destination, assignment, and diagnostic administration consume 0 visible form rows on initial presentation and each opens through exactly one explicit disclosure action.
- **SC-003**: One hundred percent of configured task and group coverage represented by the latest complete snapshot is named or counted in the overview.
- **SC-004**: One hundred percent of supported delivery states include plain-language meaning and next-action guidance that remains understandable without color.
- **SC-005**: One hundred percent of specialist concepts named in FR-012 have an accessible inline explanation associated with the relevant section or control.
- **SC-006**: The initial page and every expanded advanced state produce no document-level horizontal overflow at 800 by 600 or 200 percent zoom.
- **SC-007**: Automated accessibility scans report zero serious or critical findings in empty, disabled, healthy, retrying, and failed fixtures.
- **SC-008**: Existing secret-redaction, write-only replacement, policy, delivery-filter, and refresh regressions continue to pass unchanged in behavior.
- **SC-009**: All canonical verification gates pass with zero S085 exclusions.

## Assumptions

- Existing notification channels, policy resolution, delivery records, and bridge operations remain authoritative.
- The latest complete workspace can be extended with secret-free configured-scope summaries without changing daemon persistence or dispatch behavior.
- Recent results are a small initial subset of the already bounded 200-record workspace; detailed filters retain the full set.
- Native disclosure is sufficient secondary navigation because there are three closely related tasks within one destination.
- A configured scope whose only destination is disabled remains configured but contributes to the disabled or setup-needed guidance.

## Dependencies

- Parent: [#228](https://github.com/shruggietech/go-schedule/issues/228).
- Completes [#232](https://github.com/shruggietech/go-schedule/issues/232).
- Depends on completed shared foundations [#229](https://github.com/shruggietech/go-schedule/issues/229) and [#230](https://github.com/shruggietech/go-schedule/issues/230).
- Related notification capability history: [#160](https://github.com/shruggietech/go-schedule/issues/160), [#19](https://github.com/shruggietech/go-schedule/issues/19), [#175](https://github.com/shruggietech/go-schedule/issues/175), [#176](https://github.com/shruggietech/go-schedule/issues/176), and [#177](https://github.com/shruggietech/go-schedule/issues/177).
- Blocks public v1.4.0 promotion tracked by [#226](https://github.com/shruggietech/go-schedule/issues/226).

## Scope Boundaries

**In scope**: Secret-free notification coverage summaries, task-first overview, bounded recent results, progressive disclosure for existing administration, contextual explanations, state guidance, responsive and accessibility coverage, native build evidence, and canonical verification.

**Out of scope**: New notification transports, dispatch or retry changes, persistence migration, authorization changes, daemon API changes, release restaging, tag mutation, release publication, and unrelated desktop destinations.
