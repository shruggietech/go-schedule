# Feature Specification: Optional native desktop notifications and release preparation

**Feature Branch**: `codex/103-native-desktop-notifications`
**Created**: 2026-09-23
**Status**: In Progress
**Delivery**: Implementation and review in progress
**Input**: S103 kickoff, GitHub #177 and #185. SMTP #176 remains separate and open.

## User Scenarios & Testing

### User Story 1 - See a relevant desktop alert (Priority: P1)

As a desktop user, I can opt into native popups for outcomes that matter while the application is running, with the originating daemon and task unmistakable.

**Why this priority**: This is the core of #177 without pretending the application is an unattended monitor.
**Independent Test**: Select a daemon and condition, produce a new matching event, and observe one correctly attributed popup; old events are not replayed at startup.

**Acceptance Scenarios**:

1. **Given** popups are enabled, **when** a matching new event arrives, **then** one native popup identifies its daemon and event.
2. **Given** two daemons have similarly named tasks, **when** each produces a matching event, **then** source identity remains distinct and neither event suppresses the other.
3. **Given** the application is closed, **when** a task completes, **then** no desktop popup is promised and daemon-owned webhook delivery remains independent.

### User Story 2 - Silence and refine popups (Priority: P1)

As a user, I can turn desktop popups off or on and choose daemon, severity, and condition categories without changing scheduler outcomes or other delivery channels.

**Why this priority**: Desktop interruptions require persistent, reversible user choice.
**Independent Test**: Mute popups, create a matching event, unmute, and create a new one while inspecting unchanged webhook policy and task results.

**Acceptance Scenarios**:

1. **Given** popups are muted, **when** matching work finishes, **then** no native popup appears and task and webhook outcomes are unchanged.
2. **Given** a filter excludes a daemon or category, **when** an excluded event arrives, **then** no popup appears; included new events still appear.
3. **Given** saved preferences, **when** the application restarts, **then** mute and filters remain in effect.

### User Story 3 - Open the exact source or degrade honestly (Priority: P2)

As a user, I can activate a popup to inspect the intended daemon and activity record when supported, and see clear guidance if notifications are unavailable or denied.

**Why this priority**: A popup should lead to its evidence without a misleading platform claim.
**Independent Test**: Activate a popup after switching daemons; also exercise denied permission and unavailable facilities.

**Acceptance Scenarios**:

1. **Given** a popup from another registered daemon, **when** I activate it, **then** the application verifies the source and opens the intended activity record or explains why it no longer can.
2. **Given** permission is denied or native notifications are unsupported, **when** a matching event arrives, **then** scheduling continues without repeated error dialogs.

### User Story 4 - Prepare an honest next release (Priority: P2)

As a maintainer, I can review a coherent v1.5.0 release boundary that describes only delivered capabilities and leaves SMTP explicitly planned.

**Why this priority**: Many capabilities have accumulated since v1.4.0.
**Independent Test**: Compare release copy and changelog to merged functionality and open issues; this PR creates no tag or public release.

**Acceptance Scenarios**:

1. **Given** release documentation, **when** I review it, **then** notifications and delivered connected-control-center behavior are described without SMTP or clustered-execution claims.
2. **Given** SMTP is deferred, **when** I inspect planning, **then** #176 and #19 remain open with an explicit later disposition.

### Edge Cases

- A source disappears, is re-paired, or changes identity between popup creation and activation.
- Reconnection exposes a backlog; old events must not become a popup storm.
- A popup arrives while another daemon is selected or multiple tasks share a name.
- The user changes mute or filters during event processing.
- Native facilities, permission, or activation callbacks are missing; daemon state is never mutated.
- The application is backgrounded versus fully closed; only running-process delivery is promised.

## Requirements

### Functional Requirements

- **FR-001**: Desktop popups MUST be disabled by default until explicitly enabled, with a persistent reversible mute control.
- **FR-002**: Users MUST be able to select daemon, severity, and condition categories for local popups. Filters MUST NOT alter daemon notification assignments, webhook delivery, or task outcomes.
- **FR-003**: A popup MUST carry source daemon identity, task or health context, condition, severity, and activity identifier when available. Duplicate suppression MUST include daemon and event identity.
- **FR-004**: The application MUST emit only newly observed events while running, with bounded catch-up and no startup replay storm.
- **FR-005**: Activation MUST open the corresponding activity record after source identity validation where the platform reports activation. Unsupported activation MUST be stated rather than claimed.
- **FR-006**: Windows, macOS, and supported Linux graphical environments MUST use native operating-system notification facilities; denied or unsupported environments MUST degrade quietly.
- **FR-007**: UI and documentation MUST state that a closed application does not guarantee desktop popups. No second resident notification service may be added solely for popups.
- **FR-008**: Release preparation MUST reconcile the unreleased changelog, release notes, issue dependencies, and product wording for a forthcoming v1.5.0 without tagging or publishing in this PR.
- **FR-009**: SMTP #176 MUST remain open and explicitly deferred. Parent #19 and roadmap #146 MUST not be closed solely by this slice.

### Key Entities

- **Desktop notification preference**: Enabled/muted state, daemon selection, severity and condition selection, and local persistence version.
- **Notification candidate**: Source daemon identity, event identity, condition, severity, task context, activity reference, and observation time.
- **Activation intent**: Source identity and activity reference validated before navigation.
- **Release boundary**: Version, delivered scope, deferred scope, and source revision; publication is separate.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A newly observed matching event produces no more than one popup per daemon and event after reconnect or repeated polling.
- **SC-002**: Muted or excluded events produce zero popups and zero changes to scheduler or webhook outcomes.
- **SC-003**: Every supported popup visibly names its source daemon; activation opens the exact record or gives a clear unavailable result.
- **SC-004**: Restored preferences preserve all popup selections across application restart.
- **SC-005**: Native facility denial or absence causes zero scheduler interruption and no repeated error dialog.
- **SC-006**: Release copy names no undelivered SMTP or clustered-execution capability and leaves their issues open.

## Assumptions

- Popups are a desktop-local convenience, not a durable daemon-owned delivery channel. Existing webhooks remain the unattended monitoring path; #19 stays open for the complete cross-channel model.
- Supported platforms are Windows and macOS desktop sessions and Linux graphical sessions with a notification service; headless Linux is unsupported.
- The reviewed PR prepares release content; a public release tag and publication require separate post-merge authorization.
- The maintainer's earlier priority decision defers SMTP until after native popups; no completion claim is made for #176.
