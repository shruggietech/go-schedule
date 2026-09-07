# Feature Specification: Task and Group Authoring

**Feature Branch**: `codex/062-task-group-authoring`

**Created**: 2026-09-07

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Local API, Go race, frontend unit/build, Chromium accessibility and scale, native Windows Wails build, real guided-command execution, and canonical eight-gate verification passed 2026-09-07 on review branch `codex/062-task-group-authoring` for [#153](https://github.com/shruggietech/go-schedule/issues/153) and [#192](https://github.com/shruggietech/go-schedule/issues/192); hosted Windows, macOS, and Linux guided execution remains pull-request evidence

**Input**: Work slice S062 bundles [#153](https://github.com/shruggietech/go-schedule/issues/153) and its child [#192](https://github.com/shruggietech/go-schedule/issues/192) to rebuild Tasks and Groups as one approachable Wails workflow, including a safe platform-native first-task example.

## User Scenarios & Testing

### User Story 1 - Understand Tasks and Groups at a Glance (Priority: P1)

An operator opens Tasks and immediately understands what exists, how tasks are organized, whether each task is ready, and whether declared or inherited state prevents execution.

**Why this priority**: Safe operation starts with an accurate overview. Editing and actions are risky when group membership, readiness, or effective state is hidden.

**Independent Test**: Load zero, one hundred, and mixed-state task collections with deeply nested groups, then confirm that search, selection, hierarchy, readiness, schedule, group path, and declared-versus-effective state remain understandable and keyboard accessible.

**Acceptance Scenarios**:

1. **Given** a connected daemon with no tasks or groups, **When** the operator opens Tasks, **Then** the application explains the empty state and offers clear task and group creation paths.
2. **Given** tasks in enabled and disabled nested groups, **When** the operator browses or filters the collection, **Then** every row exposes its group path, readiness, next occurrence or manual-only state, declared state, and effective state without relying on color alone.
3. **Given** at least one hundred tasks and deeply nested groups, **When** the operator searches, selects, and moves through the collection, **Then** the visible interaction remains responsive, selection remains stable, and long or internationalized values remain available without obscuring primary actions.

---

### User Story 2 - Create a Known-Good Task (Priority: P1)

A new operator creates an inactive task through a guided editor, understands the exact program and ordered arguments, previews its schedule, saves it, runs it intentionally, and knows where its captured output will appear.

**Why this priority**: A successful first task is the shortest path from installation to trust, and the inactive default prevents accidental execution.

**Independent Test**: On Windows, macOS, and Linux, create a task from an empty editor using the platform-native suggestion, verify exact parsing and preview, save it inactive, invoke Run now against This computer, and observe recognizable captured output in Activity.

**Acceptance Scenarios**:

1. **Given** a fresh new-task editor, **When** it opens, **Then** the command input, empty preview guidance, and help show the same safe suggestion for the active execution host: `cmd.exe /d /c ver` on Windows, `/usr/bin/sw_vers` on macOS, or `uname -a` on Linux.
2. **Given** an empty focused command input for a new task, **When** the operator presses unmodified Tab once, **Then** the suggestion is inserted, focus remains in the input, the caret moves to the end, validation refreshes, and the exact program and ordered arguments appear.
3. **Given** a nonempty command input, an edit-task flow, or Shift+Tab, **When** the operator uses keyboard traversal, **Then** no suggestion replaces stored or typed content and focus can move backward and forward without a keyboard trap.
4. **Given** valid basic task details, **When** the operator saves without opening advanced settings, **Then** the task is created inactive using safe documented defaults and the application explains how to use Run now and inspect Activity.

---

### User Story 3 - Edit Complete Task Intent (Priority: P1)

An operator can inspect and change every durable task field, including group assignment, direct command details, environment, input, working directory, run identity, recurrence or one-off timing, timezone, overlap, catch-up, missing-date, time-basis, and daylight-saving policies.

**Why this priority**: The replacement cannot become authoritative until it preserves the complete task model and its advanced scheduling semantics.

**Independent Test**: Round-trip representative recurring, one-off, manual-only, draft, legacy-empty-command, and advanced-policy tasks through the editor, including validation failures and exact schedule previews, without changing untouched fields.

**Acceptance Scenarios**:

1. **Given** an existing task, **When** the operator opens it, **Then** all stored values and the exact launch preview are represented without silent default substitution.
2. **Given** invalid command, schedule, timezone, environment, or policy input, **When** the operator attempts to save or requests a preview, **Then** the relevant field receives useful guidance, focus can reach it, and no partial mutation occurs.
3. **Given** valid recurring or one-off input, **When** the operator changes it, **Then** the next five occurrences and policy summary update before save and clearly distinguish schedule prediction from recorded execution.
4. **Given** an existing task with an empty command, **When** it is edited, **Then** the platform suggestion is never inserted automatically and the stored value remains unchanged until the operator edits it.

---

### User Story 4 - Organize Work with Nested Groups (Priority: P2)

An operator can create draft groups, rename and reparent them, assign or unassign tasks, enable or disable a subtree, and delete a group while understanding the consequences for descendants and tasks.

**Why this priority**: Groups are the organizing and eligibility boundary for larger installations, and their cascading effects must be explicit before mutation.

**Independent Test**: Build a hierarchy with duplicate names at different paths, move groups and tasks, exercise cascade state, attempt a cycle, and delete a subtree while verifying full-path identity, authoritative outcomes, and preserved ungrouped tasks.

**Acceptance Scenarios**:

1. **Given** groups with repeated leaf names, **When** the operator selects a parent or task destination, **Then** full paths uniquely distinguish every valid choice.
2. **Given** a group subtree, **When** the operator disables or enables it, **Then** the confirmation names This computer, describes the subtree effect, and the refreshed view distinguishes each descendant's declared state from effective eligibility.
3. **Given** a reparenting operation that would create a cycle, **When** the operator attempts to save, **Then** the mutation is rejected with guidance attached to the parent choice and the hierarchy remains unchanged.
4. **Given** a group containing children and tasks, **When** the operator confirms deletion, **Then** the confirmation names This computer, identifies the subtree consequence, deletes descendant groups, and leaves affected tasks ungrouped according to the existing product contract.

---

### User Story 5 - Perform Safe Task Actions (Priority: P2)

An operator can run, enable, disable, or delete a selected task with clear target identity, bounded feedback, and authoritative refresh.

**Why this priority**: Operational actions are essential, but they depend on the overview and editor accurately establishing identity and readiness first.

**Independent Test**: Exercise Run now, enable, disable, and delete across ready, manual-only, inherited-disabled, invalid, disconnected, and concurrently changed tasks, confirming safe gating, target-aware confirmations, announcements, and recovery.

**Acceptance Scenarios**:

1. **Given** a runnable task, **When** the operator confirms Run now, **Then** the action names This computer, preserves current selection, reports acceptance or a safe failure, and directs the operator to Activity.
2. **Given** a task that is not runnable, **When** actions are presented, **Then** unavailable actions explain the blocking readiness or group condition instead of failing silently.
3. **Given** a destructive task action, **When** the operator initiates it, **Then** an explicit confirmation names both the task and This computer, cancellation changes nothing, and success refreshes authoritative state.
4. **Given** connection loss or a concurrent daemon-side change, **When** an operation fails or a live event arrives, **Then** the application preserves unsaved edits, avoids optimistic success claims, and offers refresh or retry without duplicating the mutation.

### Edge Cases

- A selected task or group can be deleted by another client while its details or editor are open.
- A live refresh can arrive while the operator has unsaved edits; the draft remains intact and is marked as based on older data rather than overwritten.
- Group depth, duplicate names, missing parents, and malformed legacy references must not create ambiguous destinations or unbounded traversal.
- Long names, commands, arguments, paths, environment values, Unicode, bidirectional text, and empty optional values must remain inspectable and must not disrupt actions.
- Manual-only tasks have no next occurrence but can still be intentionally run when otherwise ready.
- A disabled task in an enabled group, an enabled task in a disabled ancestor, and an invalid task must expose different declared and effective explanations.
- Schedule preview can fail or become stale independently of the saved task list.
- Repeated action activation, slow responses, and late responses after selection changes must not duplicate writes or apply results to the wrong entity.
- The active execution host can differ from the desktop client in future remote operation; suggestions always follow the execution host identity.
- Clearing a newly inserted suggestion permits one later empty-input Tab insertion, while normal Tab traversal resumes after insertion and for all nonempty content.

## Requirements

### Functional Requirements

- **FR-001**: The Tasks area MUST provide one coherent task and group workflow with clear creation, search, filtering, selection, detail, editing, action, and empty-state paths.
- **FR-002**: Task collections MUST expose name, full group path, schedule or manual-only status, readiness, declared enabled state, effective eligibility, and the next occurrence when one exists.
- **FR-003**: Group collections MUST expose hierarchy, full path, declared enabled state, effective subtree state, child and task counts, and degraded legacy relationships without ambiguous identity.
- **FR-004**: Search and filters MUST operate over at least one hundred tasks while preserving stable selection and keyboard access through refreshes.
- **FR-005**: Operators MUST be able to create, inspect, edit, enable, disable, run, and delete tasks supported by the current product contract.
- **FR-006**: The task editor MUST preserve name, group assignment, command, ordered arguments, working directory, environment, standard input, run identity, timezone, recurring or one-off schedule, overlap, catch-up, missing-date, time-basis, daylight-saving gap, daylight-saving overlap, and enabled state.
- **FR-007**: Basic creation MUST keep advanced fields behind intentional disclosure while using documented defaults that create an inactive task without requiring advanced interaction.
- **FR-008**: Command entry MUST preserve the direct-execution grammar and show an exact, non-executable program and ordered-argument preview.
- **FR-009**: Recurring and one-off input MUST provide validation, a plain-language policy summary, and up to five next occurrences before save; manual-only work MUST be labelled explicitly.
- **FR-010**: Validation failures MUST attach to the relevant control, remain available to assistive technology, move focus to the first invalid field on submission, and prevent partial mutation.
- **FR-011**: A fresh new-task editor MUST derive one safe command suggestion from the active execution host and use the exact canonical Windows, macOS, or Linux mapping defined in User Story 2.
- **FR-012**: The suggestion MUST be identical in the command placeholder, empty preview guidance, help, and user documentation, with an explanation of captured output in Activity.
- **FR-013**: One unmodified Tab on an empty focused new-task command input MUST insert the suggestion, retain focus, place the caret at the end, refresh validation and preview, and allow subsequent forward traversal.
- **FR-014**: Suggestion insertion MUST NOT occur for nonempty input, Shift+Tab, edit flows, or unsupported execution platforms; clearing an inserted suggestion MAY make the one-shot insertion available again while the field remains a new-task draft.
- **FR-015**: Suggested commands MUST terminate promptly, produce recognizable captured output, require no optional runtime, open no listener or network request, write no file, request no elevation, and wait for no input.
- **FR-016**: Operators MUST be able to create, rename, reparent, enable, disable, and delete groups, including nested and initially disabled draft groups.
- **FR-017**: Group and task destination choices MUST use complete hierarchy paths, include an ungrouped choice for tasks and a root choice for groups, and exclude invalid descendant destinations.
- **FR-018**: Group mutation failures, including parent cycles and missing parents, MUST preserve the previous hierarchy and identify the affected field or operation.
- **FR-019**: Cascading group state MUST preserve declared child and task settings while separately explaining effective eligibility inherited from ancestors.
- **FR-020**: Deleting a group MUST follow the existing cascade contract for descendant groups while preserving affected tasks as ungrouped.
- **FR-021**: Run-now, deletion, and cascading group actions MUST require explicit confirmations that name the entity, the active daemon display name, and the relevant consequence.
- **FR-022**: Mutations MUST use daemon-authoritative results, suppress duplicate submission while pending, announce success or safe failure, and refresh affected collections without claiming optimistic success.
- **FR-023**: Live task and group events MUST refresh authoritative data without stealing focus, losing stable selection, or overwriting unsaved drafts.
- **FR-024**: If the edited entity changes or disappears remotely, the application MUST preserve the local draft, disclose that it is stale, and require an explicit reload or save decision.
- **FR-025**: Connection loss MUST retain safe read-only context when available, disable unavailable mutations with a reason, and recover through the S061 connection workflow.
- **FR-026**: No command input, environment value, standard input, run identity, path, or backend error detail MAY appear in ordinary event announcements, diagnostic messages, or confirmation copy unless the operator explicitly opens that task's details.
- **FR-027**: The workflow MUST remain understandable without color, support full keyboard operation and visible focus, preserve semantic labels and relationships, honor reduced motion, and remain usable at the supported narrow window and zoom levels.
- **FR-028**: User documentation MUST replace the previous Python listener and Windows-only file-writing examples with the canonical platform mapping and the Save, Run now, Activity walkthrough.
- **FR-029**: The existing Fyne shipping surface MUST use the same canonical suggestion mapping while it remains supported, without creating a divergent source of truth.
- **FR-030**: All task and group behavior MUST remain local-only for v1.2.0 and MUST NOT add remote transport, feature-screen migration outside Tasks and Groups, packaging cutover, or Fyne removal.

### Key Entities

- **Task**: Durable scheduled or manual-only work with direct execution, scheduling policies, group membership, declared state, readiness, and effective eligibility.
- **Task Draft**: Unsaved editor state with validation, exact command preview, schedule preview, dirty and stale indicators, and new-versus-edit identity.
- **Group**: Durable hierarchical organization with parent identity, declared state, effective state, descendants, and task membership.
- **Execution Host**: The active daemon identity and platform that determine confirmations, availability, and the safe command suggestion.
- **Command Suggestion**: The canonical safe example for an execution platform, including display command, exact program, ordered arguments, explanation, and recognizable output expectation.
- **Operation Result**: A bounded accepted or rejected outcome for a task or group mutation, safe message, affected identity, and authoritative refresh requirement.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A new operator can create an inactive known-good task from the empty state, using only the basic editor and keyboard, in under two minutes.
- **SC-002**: On each supported execution platform, the guided example exits within five seconds under normal conditions and produces recognizable captured output without network or filesystem mutation.
- **SC-003**: All current task fields and group lifecycle operations round-trip through the replacement workflow with no unexplained omissions or silent value changes.
- **SC-004**: Search, filtering, selection, and primary actions remain responsive within 200 milliseconds for one hundred tasks and a group hierarchy twenty levels deep, excluding daemon round-trip time.
- **SC-005**: One hundred consecutive create, update, action, refresh, and cancellation scenarios produce no duplicate mutations, stale-entity application, race findings, or leaked owned work.
- **SC-006**: Every destructive or cascading confirmation identifies the affected entity and active daemon, and every unavailable action exposes a non-color-only reason.
- **SC-007**: Automated accessibility checks report no serious or critical violations across empty, populated, editor, validation, confirmation, narrow-window, high-zoom, and reduced-motion states.
- **SC-008**: Keyboard-only scenarios cover creation, suggestion insertion, subsequent traversal, editing, group hierarchy management, Run now, cancellation, and error recovery without a trap or lost focus.
- **SC-009**: Windows, macOS, and Linux packaged validation follows the guided example through Run now to a successful Activity record with platform-recognizable output.
- **SC-010**: User-facing guidance contains none of the retired Python HTTP listener or Windows-only file-writing first-run examples.

## Clarifications

### Session 2026-09-07

- Q: Which source controls task and group state after writes or live updates? → A: The daemon remains authoritative; pending actions do not mutate local records optimistically, and accepted operations trigger a targeted refresh.
- Q: What happens when live data changes during an unsaved edit? → A: Preserve the local draft, mark it stale, and require the operator to reload or make an explicit save decision.
- Q: How is the command suggestion shared while both desktop implementations exist? → A: One canonical platform mapping feeds both supported desktop experiences, their acceptance evidence, and user documentation.
- Q: How much hierarchy and list scale must this slice prove? → A: One hundred tasks and twenty group levels, with 200-millisecond local interaction targets.
- Q: Should S062 migrate Activity to prove the guided run? → A: No; Run now links to the existing honest Activity placeholder, while packaged contract validation proves backend capture and #155 owns the Activity screen migration.

## Assumptions

- S061's protected local connection, execution-host identity, event channel, and shared shell primitives remain the base contract.
- The daemon and existing product contract remain authoritative; the slice does not change persisted task or group semantics.
- New tasks remain inactive by default, consistent with the established draft-task safety contract.
- Existing group deletion semantics cascade descendant groups and ungroup affected tasks.
- Existing schedule parsing, timezone, readiness, command-line parsing, and policy logic remain authoritative and are reused rather than duplicated.
- Activity screen migration, other automation-source screens, packaging cutover, remote access, and Fyne removal remain assigned to later v1.2.0 slices.
