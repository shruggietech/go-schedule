# Feature Specification: Desktop Notification Management

**Feature Branch**: `codex/068-desktop-notifications`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Desktop channel, policy, and redacted delivery-history management plus focused Go race, React, Chromium accessibility and zoom, native Windows Wails build, and canonical eight-gate verification passed 2026-09-08 on review branch `codex/068-desktop-notifications` for [#160](https://github.com/shruggietech/go-schedule/issues/160).

**Input**: GitHub issue [#160](https://github.com/shruggietech/go-schedule/issues/160), add clear notification setup and delivery history to the desktop control center.

## User Scenarios & Testing

### User Story 1 - Configure and repair a webhook channel (Priority: P1)

A desktop user can create a named webhook destination, provide an optional authorization value, send an unmistakable test, disable or re-enable delivery, replace protected values when repairing the destination, and remove the channel without using the command line.

**Why this priority**: Every policy and delivery-history workflow depends on a trustworthy destination that users can administer without exposing stored credentials.

**Independent Test**: Create a channel with and without authorization, test it, disable and re-enable it, replace its endpoint and authorization, then remove it while verifying protected values are never redisplayed.

**Acceptance Scenarios**:

1. **Given** no channels, **When** a user enters a name, HTTPS endpoint, optional authorization, and saves, **Then** the channel appears with a redacted destination summary and authorization-presence indicator while both protected inputs are cleared.
2. **Given** an existing channel, **When** a user edits it, **Then** stored endpoint and authorization values remain hidden and are changed only through explicit replacement intent.
3. **Given** an enabled channel, **When** a user sends a test, **Then** the result is announced and its history record is labeled as a test rather than a task outcome.
4. **Given** a disabled channel, **When** a user tries to test it, **Then** testing is unavailable and the interface explains that the channel must be enabled first.
5. **Given** an existing channel, **When** a user disables, re-enables, repairs, or removes it, **Then** the resulting state is refreshed and announced without exposing arbitrary daemon error content.

---

### User Story 2 - Assign understandable task and group policies (Priority: P1)

A desktop user can choose a task or group, assign channels for failure or success outcomes, and understand whether a task uses its own policy, an inherited group policy, or no notification policy.

**Why this priority**: A working channel has no production value until users can bind it to outcomes, and inheritance must be explained to prevent accidental silence or noise.

**Independent Test**: Configure a group failure policy, inspect a child task's inherited policy, add a direct task override, enable success notification deliberately, and clear the direct override to restore inheritance.

**Acceptance Scenarios**:

1. **Given** a group and at least one channel, **When** a user assigns failure delivery, **Then** the group policy is saved as a complete replacement and displayed as directly configured.
2. **Given** a task without a direct policy whose group has assignments, **When** the task is selected, **Then** the interface names the inherited group source and lists the effective channel outcomes separately from direct assignments.
3. **Given** a task with inherited policy, **When** a user saves any direct assignments, **Then** the interface explains that the task policy now overrides all group inheritance.
4. **Given** a direct task policy, **When** the user clears all direct assignments, **Then** the effective policy refreshes to the nearest configured ancestor or none.
5. **Given** success delivery is selected, **When** the user reviews or saves the policy, **Then** a non-blocking warning explains its higher notification volume.

---

### User Story 3 - Understand delivery progress and failures (Priority: P2)

A desktop user can inspect bounded recent notification delivery history, distinguish tests from production task outcomes, filter by channel and state, and see enough redacted evidence to understand queued, retrying, successful, and failed deliveries.

**Why this priority**: Delivery evidence closes the operational loop after setup and policy assignment, especially when a receiver is unavailable or returns an error.

**Independent Test**: Load mixed test and production records across pending, claimed, succeeded, and failed states, filter them, select each record, and verify destination, attempts, timing, status, and safe failure details remain understandable.

**Acceptance Scenarios**:

1. **Given** mixed history, **When** it loads, **Then** every row identifies test or task outcome, channel, destination summary, state, attempt count, and relevant time without exposing endpoint or authorization secrets.
2. **Given** a pending delivery after at least one attempt, **When** it is displayed, **Then** it is identified as retry scheduled and shows the next attempt time rather than appearing identical to a first queued attempt.
3. **Given** a claimed delivery, **When** it is displayed, **Then** it is identified as sending in progress and is not presented as a completed result.
4. **Given** a failed delivery, **When** it is selected, **Then** the interface shows the final safe error, response status when available, and correlation identifiers needed for investigation.
5. **Given** no matching history or a refresh failure, **When** filters or connection state change, **Then** the interface preserves the last complete snapshot and offers accessible empty or recovery guidance.

### Edge Cases

- A create or repair request has a blank name, a blank replacement endpoint, a malformed endpoint, or a non-HTTPS endpoint.
- An authorization replacement is intentionally blank to remove a stored value.
- A channel changes or disappears between loading and mutation.
- A channel is removed while assignments or unfinished deliveries reference it; the desktop reflects the authoritative server result and retained terminal evidence.
- A task or group changes or disappears while its policy editor is open.
- A task belongs to nested groups and the effective source is an ancestor rather than its immediate group.
- A policy contains a disabled or subsequently deleted channel.
- No tasks, groups, channels, direct assignments, inherited assignments, or delivery records exist.
- A delivery has no task or run because it is a channel test.
- A delivery is pending before its first attempt, pending after a failed attempt, claimed, succeeded, or terminally failed.
- Refresh responses arrive out of order or a domain event arrives during a pending mutation.
- The daemon disconnects while a complete notification snapshot is visible.
- The interface is used entirely by keyboard, with a screen reader, or at 80, 100, 150, and 200 percent browser zoom.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST provide a complete Notifications route for channel setup, policy assignment, and recent delivery history.
- **FR-002**: Users MUST be able to create, edit, disable, enable, test, and remove webhook channels without using the CLI.
- **FR-003**: Channel creation MUST require a non-empty name and valid HTTPS endpoint and MUST allow an optional authorization value.
- **FR-004**: Stored endpoint and authorization values MUST remain write-only; the desktop MUST display only a redacted destination summary and whether authorization exists.
- **FR-005**: Editing an existing channel MUST keep protected fields blank and MUST require explicit replacement intent before changing the endpoint or authorization.
- **FR-006**: Authorization replacement MUST support both setting a new value and deliberately clearing the stored value without redisplaying the previous value.
- **FR-007**: A successful channel mutation MUST clear protected draft input, refresh authoritative channel metadata, preserve stable selection when possible, and announce the result.
- **FR-008**: Testing MUST be available only for enabled channels, MUST use the production transport path, and MUST create history that is visibly identified as a test.
- **FR-009**: Destructive channel removal MUST require explicit confirmation and MUST explain that terminal redacted history can remain while active configuration is removed.
- **FR-010**: The policy editor MUST list current tasks and groups with stable identity and enough hierarchy context to distinguish similarly named scopes.
- **FR-011**: Users MUST be able to replace the complete direct notification policy for one selected task or group with zero or more channel assignments.
- **FR-012**: Each assignment MUST independently support failure and success outcomes and MUST require at least one selected outcome.
- **FR-013**: Enabling success notifications MUST display a persistent non-blocking volume warning before save without preventing deliberate use.
- **FR-014**: Task policy detail MUST distinguish direct assignments from the effective policy and MUST identify whether the effective source is the task, a named group ancestor, or no scope.
- **FR-015**: Clearing all direct task assignments MUST restore and display inherited group policy when one exists; direct assignments MUST override rather than merge with inherited assignments.
- **FR-016**: Group policy detail MUST identify its assignments as direct configuration and explain that descendant tasks inherit them only when no nearer policy overrides them.
- **FR-017**: The desktop MUST show up to 200 recent redacted deliveries ordered newest first and allow filtering by channel and lifecycle state.
- **FR-018**: Delivery rows and detail MUST visibly distinguish tests from production task outcomes and MUST not imply that a test represents task completion.
- **FR-019**: Delivery presentation MUST distinguish first-attempt queued, retry scheduled, sending, successful, and terminal failure states using text in addition to color.
- **FR-020**: Delivery detail MUST show available channel, redacted destination, task, group, run, attempt, next-attempt, completion, HTTP status, and safe failure evidence without showing payload secrets or protected channel values.
- **FR-021**: Asynchronous reads MUST discard stale responses, relevant notification events MUST refresh through a bounded debounce, and failed refreshes MUST preserve the last complete successful snapshot.
- **FR-022**: Mutation controls MUST be disabled while disconnected or while the relevant mutation is pending, and duplicate activation MUST not create duplicate mutations.
- **FR-023**: Validation, conflict, missing-resource, and daemon-unavailable outcomes MUST use a bounded user-safe vocabulary with actionable recovery guidance.
- **FR-024**: Empty, disabled, queued, retrying, sending, failed, and successful states MUST be keyboard accessible, screen-reader meaningful, understandable without color, and readable from 80 through 200 percent zoom.
- **FR-025**: Automated Go, React, accessibility, stale-response, secret-boundary, policy-inheritance, history-state, Chromium interaction, native-build, and canonical repository verification MUST pass.
- **FR-026**: Issue [#160](https://github.com/shruggietech/go-schedule/issues/160) MUST remain traceable through the specification, tasks, changelog, pull request, and verification record.

### Key Entities

- **Notification Channel Summary**: Stable identity, name, transport kind, redacted destination summary, authorization-presence indicator, enabled state, and update time.
- **Channel Draft**: User-entered name plus write-only endpoint and authorization replacement values and explicit replacement intent.
- **Policy Scope**: One task or group identified by stable ID, human name, hierarchy context, and scope type.
- **Notification Assignment Draft**: One channel selection with independent success and failure outcome flags within a complete scope replacement.
- **Effective Policy**: The source scope selected for a task and the assignments inherited from it, distinct from directly configured task assignments.
- **Delivery Record**: Redacted delivery evidence including test or task event kind, channel and destination summaries, correlation identifiers, lifecycle state, attempts, timing, status, and safe final error.
- **Notification Workspace**: One complete snapshot of channels, available policy scopes, and bounded recent delivery history.
- **Notification Operation Result**: Accepted, rejected, conflict, stale, or unavailable outcome with safe message, optional field, refreshed workspace, or selected policy.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A keyboard-only user can create, test, disable, repair, re-enable, and remove a channel and can configure both a task policy and group policy without using the CLI.
- **SC-002**: Automated secret-boundary tests find zero stored endpoint or authorization values in any desktop response, frontend state initialized from server data, visible text, accessibility name, or history detail.
- **SC-003**: Every supported policy fixture identifies direct assignments and the effective source correctly for task override, immediate-group inheritance, ancestor-group inheritance, and no-policy cases.
- **SC-004**: Every delivery lifecycle fixture is labeled with its user-facing state, and all test fixtures are distinguishable from production task outcomes in both list and detail views.
- **SC-005**: A 200-delivery snapshot remains filterable and keyboard navigable at 80, 100, 150, and 200 percent zoom without clipped primary actions or horizontal page overflow at the supported desktop viewport.
- **SC-006**: Focused tests, native Wails build, and `sh scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-07

- Q: How are stored endpoint and authorization values edited without disclosure? A: Existing editors start blank and expose separate replace controls. A blank authorization replacement deliberately clears it; a blank endpoint replacement is rejected.
- Q: Does a task policy merge with inherited group assignments? A: No. Any direct task assignment replaces inheritance as the effective policy, while clearing all direct assignments restores the nearest configured ancestor.
- Q: What history is in scope? A: The 200 most recent redacted deliveries, newest first, with channel and state filters. Broader pagination remains outside this slice.
- Q: How are pending and claimed deliveries explained? A: A never-attempted pending record is queued, an attempted pending record is retry scheduled, and a claimed record is sending.
- Q: Is success notification discouraged or blocked? A: It remains available but carries a persistent non-blocking volume warning whenever selected.

## Assumptions

- The notification API, protected local transport, encrypted or permission-bounded secret persistence, dispatcher, and terminal delivery evidence delivered by issues #158 and #159 are authoritative and unchanged by this slice.
- The existing Wails shell, desktop connection manager, task and group reads, event stream, component primitives, and validation gates provide the implementation foundation.
- The desktop manages the current local daemon only; remote targets, multiple transports, arbitrary payload editing, delivery replay, and bulk channel import are outside this slice.
- The server's complete-replacement policy behavior and nearest-ancestor precedence are the source of truth; the desktop explains them but does not calculate an alternate policy.
- History payload bodies are intentionally excluded from desktop presentation even though the server retains a redacted structured event.
- This slice fully resolves issue #160. Parent issue #19 remains open for the remaining v1.3.0 notification and MCP work.
