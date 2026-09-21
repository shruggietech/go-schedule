# Feature Specification: Actionable Notification Conditions

**Feature Branch**: `codex/095-notification-conditions`

**Created**: 2026-09-21

**Status**: Implemented

**Delivery**: Focused Go, desktop, frontend, Playwright accessibility, native Windows Wails, canonical race, integration, coverage, documentation, and automation gates passed on review branch `codex/095-notification-conditions`; hosted review is tracked by the S095 pull request.

**Input**: Work slice S095, complete issue #175 by adding durable, explainable notification conditions beyond per-run success and failure.

## Clarifications

### Session 2026-09-21

- Q: How can a stopped daemon report unhealthy status? -> A: An enabled destination may opt into bounded health heartbeats. Receivers treat a missing heartbeat as unhealthy and a resumed heartbeat as recovery; go-schedule does not claim it can send while stopped.
- Q: How do existing failure assignments migrate? -> A: Every existing `on_failure` assignment becomes a failure condition with threshold one, preserving current behavior.
- Q: How do multiple simultaneous problems behave? -> A: One deterministic primary condition is selected in the order failure-to-start, consecutive failure, then duration exceeded. A changed primary condition is a new transition; repeated identical conditions follow reminder and quiet-period rules.
- Q: What does quiet period suppress? -> A: It suppresses repeated reminders and routine success messages. First problem transitions and recovery transitions are never hidden by the quiet period.

## User Scenarios & Testing

### User Story 1 - Alert on actionable task problems (Priority: P1)

As an operator, I can notify only after a selected number of consecutive failures, immediately when a process cannot start, or when a run exceeds an expected duration so routine outcomes do not create noise.

**Why this priority**: Useful conditions are the core outcome and must remain independent from notification transport behavior.

**Independent Test**: Configure each condition on task and group policies, record controlled runs, and inspect deliveries to prove the correct transition, reason, streak, threshold, and source policy.

**Acceptance Scenarios**:

1. **Given** a failure threshold of three, **When** two failures occur, **Then** no notification is created; the third consecutive failure creates one condition delivery.
2. **Given** a failure-to-start condition, **When** process creation fails, **Then** one immediate delivery identifies failure to start without exposing the command or error details.
3. **Given** a duration threshold, **When** a terminal run meets or exceeds it, **Then** one delivery identifies the observed and configured durations.
4. **Given** a group policy and no direct task policy, **When** a descendant task meets a condition, **Then** the nearest group policy evaluates with the same deterministic rules as a direct task policy.

---

### User Story 2 - Suppress repeats and report recovery (Priority: P1)

As an operator, I can receive a first problem notification, bounded reminders, and one recovery notification without a stream of duplicates.

**Why this priority**: Thresholds without durable suppression and recovery would merely move notification noise rather than solve it.

**Independent Test**: Drive a task through failure, repeated failure, daemon restart, reminder eligibility, policy replacement, and recovery using an injected clock, then inspect durable condition state and deliveries.

**Acceptance Scenarios**:

1. **Given** an active problem and reminders disabled, **When** equivalent failing runs continue, **Then** they create no duplicate delivery.
2. **Given** a reminder interval, **When** the same problem continues before and after the interval, **Then** only the eligible later run creates a reminder.
3. **Given** an active problem and recovery enabled, **When** the next run no longer matches any configured problem, **Then** one recovery delivery is created and later healthy runs do not repeat it.
4. **Given** a daemon restart, **When** evaluation resumes, **Then** streak, active problem, notification time, and duplicate suppression continue from durable state.
5. **Given** a changed effective policy, **When** the next run is evaluated, **Then** stale state is reset before the new policy is applied.

---

### User Story 3 - Monitor daemon presence honestly (Priority: P2)

As an operator, I can opt a destination into periodic daemon-health heartbeats so an external receiver can identify a missing daemon and recognize when heartbeats resume.

**Why this priority**: A daemon cannot emit while stopped, so an explicit heartbeat contract is the smallest truthful daemon-health condition.

**Independent Test**: Configure a heartbeat interval, advance an injected clock across due and not-due boundaries, restart the dispatcher, and verify stable bounded scheduling without duplicate due work.

**Acceptance Scenarios**:

1. **Given** health heartbeats are disabled, **When** the daemon runs, **Then** no health delivery is created.
2. **Given** a supported heartbeat interval, **When** it becomes due, **Then** one `daemon.health` delivery reports current healthy presence and the expected next heartbeat.
3. **Given** a restart before the next due time, **When** scheduling resumes, **Then** no duplicate heartbeat is created.
4. **Given** a receiver observes no heartbeat for the documented missing interval, **When** it alerts, **Then** documentation makes clear that this is receiver-side absence detection rather than a daemon-sent outage claim.

---

### User Story 4 - Configure and understand conditions (Priority: P1)

As a desktop or CLI user, I can configure advanced conditions in plain language and see why each delivery fired and when another can occur.

**Why this priority**: Condition behavior is unsafe if it is not explainable at configuration and evidence surfaces.

**Independent Test**: Configure task and group policies through API, CLI, and desktop, then inspect effective policy and delivery evidence at 800 by 600, 200 percent zoom, and keyboard-only operation.

**Acceptance Scenarios**:

1. **Given** an assignment editor, **When** advanced conditions are opened, **Then** thresholds, failure-to-start, duration, recovery, reminders, and quiet period have concise definitions and validated inputs.
2. **Given** a saved policy, **When** direct or effective assignments are inspected, **Then** the complete condition configuration and inheritance source are visible.
3. **Given** a condition delivery, **When** recent results or diagnostics are inspected, **Then** the condition type and safe reason explain why it fired.

### Edge Cases

- A failure threshold is changed while a streak is active; the next evaluation starts fresh under the new policy fingerprint.
- A run both fails to start and satisfies a consecutive-failure threshold; failure-to-start is the primary condition and only one delivery is created per channel.
- A long successful run exceeds its duration threshold; it opens a duration problem rather than also sending a routine success notification.
- A task is deleted while condition state exists; task-scoped state is removed without changing terminal delivery history.
- A channel is disabled while condition state is active; evaluation updates durable state but creates no outbound work until an enabled assignment becomes eligible under a later transition or reminder.
- Clock movement cannot make a reminder or heartbeat eligible earlier than its persisted UTC due time.
- A heartbeat delivery is still pending when another interval passes; no second non-terminal heartbeat is created for the same channel.

## Requirements

### Functional Requirements

- **FR-001**: Notification assignments MUST retain success and failure compatibility and add failure threshold, failure-to-start, duration threshold, recovery, reminder interval, and quiet-period settings.
- **FR-002**: Existing failure assignments MUST migrate to a failure threshold of one without creating, deleting, or sending any notification during migration.
- **FR-003**: Failure thresholds MUST accept integers from 1 through 100 when failure notification is enabled and MUST count only consecutive terminal failures for the same task and channel under one unchanged effective policy.
- **FR-004**: A run MUST record whether process creation failed as structured state rather than by parsing captured output.
- **FR-005**: Failure-to-start MUST be independently selectable and MUST take precedence over other conditions for the same run and channel.
- **FR-006**: Duration thresholds MUST accept zero for disabled or 1 second through 30 days and MUST compare against non-negative terminal run duration.
- **FR-007**: Problem precedence MUST be failure-to-start, consecutive failure, then duration exceeded, with at most one task-condition delivery per channel and run.
- **FR-008**: Condition evaluation and source-run recording MUST occur in one database transaction before outbound work begins.
- **FR-009**: Durable state MUST retain the policy fingerprint, consecutive failure count, active primary condition, active-since time, last notification time, last evaluated run, and update time for each task and channel.
- **FR-010**: Equivalent repeated problems MUST be suppressed unless a positive reminder interval has elapsed; reminder intervals MUST accept zero for disabled or 1 minute through 30 days.
- **FR-011**: Recovery MUST clear active problem state and create at most one recovery delivery when enabled and the current run matches no configured problem.
- **FR-012**: Quiet periods MUST accept zero for disabled or 1 minute through 30 days and suppress only routine success messages and reminders; first problem and recovery transitions MUST remain observable.
- **FR-013**: Replacing or inheriting a materially different effective assignment MUST change its deterministic fingerprint and reset stale evaluation state before the next run.
- **FR-014**: Disabled or removed channels MUST create no new delivery; deleting a task or channel MUST remove its condition state through foreign-key lifecycle rules.
- **FR-015**: A notification channel MAY opt into daemon-health heartbeats with zero for disabled or an interval from 1 minute through 24 hours.
- **FR-016**: Heartbeat scheduling MUST be durable, create no duplicate non-terminal work for a channel, and continue from its persisted next-due time after restart.
- **FR-017**: Heartbeat payloads MUST identify `daemon.health`, current healthy presence, creation time, and expected next time without containing task configuration, credentials, host secrets, or unsupported outage claims.
- **FR-018**: Webhook payloads and delivery evidence MUST include a bounded condition kind and safe explanation for condition-generated deliveries while remaining backward compatible with the existing webhook schema.
- **FR-019**: API, client, CLI, MCP Manage assignment input, desktop backend, and desktop frontend contracts MUST preserve and expose the complete condition configuration.
- **FR-020**: Boundary validation MUST reject condition-only recovery without a problem condition, out-of-range values, duplicate channels, and assignments with no selected condition.
- **FR-021**: Notification evaluation, heartbeat creation, delivery retry, and receiver failure MUST remain isolated from scheduler outcome and worker capacity.
- **FR-022**: The desktop MUST explain thresholds, failure-to-start, duration, recovery, reminders, quiet periods, heartbeat absence, and effective inheritance in plain language and without color-only meaning.
- **FR-023**: User documentation MUST publish state transitions, precedence, suppression, restart behavior, heartbeat absence semantics, safe defaults, API and CLI examples, and receiver expectations.
- **FR-024**: S095 MUST complete the functional acceptance criteria of issue #175 and update its parent coordinators without claiming SMTP or native desktop notification delivery.

### Key Entities

- **Notification condition configuration**: Advanced per-channel settings stored on a task or group assignment.
- **Notification condition state**: Durable per-task, per-channel evaluation state tied to a deterministic effective-policy fingerprint.
- **Condition delivery**: Existing notification delivery extended with a safe condition kind and explanation.
- **Daemon-health heartbeat state**: Durable per-channel next-due schedule for opt-in healthy-presence events.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every configured task condition produces exactly zero or one matching delivery per channel and source run according to the documented transition table.
- **SC-002**: A threshold of N produces its first failure delivery on failure N for every N from 1 through 100 and resets after a non-failing run or policy change.
- **SC-003**: Repeated identical problems create no delivery before the reminder interval and exactly one eligible reminder after it, including across daemon restart.
- **SC-004**: Recovery creates exactly one delivery per active problem transition and leaves subsequent healthy runs eligible only for explicitly selected routine success notifications.
- **SC-005**: Heartbeat scheduling produces no more than one non-terminal heartbeat per channel and preserves its due schedule across restart.
- **SC-006**: Existing success and failure assignments retain behavior after migration with no outbound work created by migration.
- **SC-007**: All condition configuration and evidence remains usable by keyboard with no serious accessibility findings or document-level horizontal overflow at 800 by 600 and 200 percent zoom.

## Assumptions

- Existing webhook transport, protected secrets, delivery retry, inheritance, and retention remain authoritative.
- A missing heartbeat is evaluated by the receiver because a stopped daemon cannot send an outage event. SMTP and native desktop transports remain outside S095.
- UTC persisted timestamps and injected clocks provide deterministic reminders and heartbeats; monotonic process time is not durable across restart.
- Health heartbeat delivery failure remains ordinary delivery evidence and cannot recursively create another health notification.

## Dependencies and Traceability

- Parent: #19.
- Completes: #175 when its functional acceptance criteria are delivered.
- Depends on completed #158, #159, and #160.
- Leaves #176 and #177 open.
