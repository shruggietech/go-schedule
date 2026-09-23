# Feature Specification: Linux desktop daemon presence and controls

**Feature Branch**: `codex/102-linux-daemon-indicator`
**Created**: 2026-09-23
**Status**: Implemented
**Delivery**: S102 review branch `codex/102-linux-daemon-indicator` delivers the per-session Linux indicator, local GUI service controls, packaging, and documentation for #256. Linux unit and race suites, Wails build, frontend tests, and canonical verification passed 2026-09-23; pull request review and merge remain pending.
**Input**: GitHub issue #256 and the approved S102 recommendation. Native run-outcome notifications (#177) and SMTP (#176) remain separate.

## User Scenarios & Testing

### User Story 1 - See truthful local daemon presence (Priority: P1)

As a Linux desktop user with a compatible status area, I can tell whether this computer's installed scheduler is running without keeping its main window open.

**Why this priority**: Presence must be truthful before any control is useful.
**Independent Test**: Observe the item after login, service stop/start, GUI close, and service failure.

**Acceptance Scenarios**:

1. **Given** an installed and healthy local service, **when** I inspect the indicator, **then** it says Running.
2. **Given** a stopped, missing, transitioning, or unreachable local service, **when** I inspect the indicator, **then** it shows the corresponding state instead of falsely saying Running.
3. **Given** no compatible status host, **when** I open the GUI, **then** I can still see local status and the scheduler itself continues normally.

### User Story 2 - Control the local service safely (Priority: P1)

As a Linux desktop user, I can start, stop, or restart the installed service from the indicator or GUI, with confirmation for work-interrupting actions and truthful authorization feedback.

**Why this priority**: A passive status item alone does not meet the local-management outcome.
**Independent Test**: Exercise each action with an authorized account, denied authorization, cancellation, and service failure.

**Acceptance Scenarios**:

1. **Given** a stopped installed service, **when** I choose Start, **then** Running appears only after fresh local health succeeds.
2. **Given** a running service, **when** I choose Stop or Restart, **then** I first see the impact on local work and can cancel without changing the service.
3. **Given** a denied, failed, or timed-out action, **when** it ends, **then** I see the current observed state and an actionable error.
4. **Given** a remote daemon selected in the GUI, **when** I use the local controls, **then** they still target This computer only.
5. **Given** an installed local service with a custom `--config` path, **when** I open the GUI or indicator, **then** both use that service's configured IPC endpoint rather than the default.

### User Story 3 - Open one GUI and retain a clean session lifecycle (Priority: P2)

As a user, I can open or focus the GUI from the indicator, and logging out or removing the desktop integration leaves no stale item.

**Why this priority**: A long-lived session item must not create duplicate windows or linger after use.
**Independent Test**: Repeat Open, close the GUI, restart the item, log out, and remove the desktop integration.

**Acceptance Scenarios**:

1. **Given** the GUI already exists, **when** I choose Open, **then** that window is focused rather than duplicated.
2. **Given** the service was deliberately stopped, **when** I open the GUI, **then** the GUI does not silently restart it.
3. **Given** an active graphical session with one compatible host, **when** startup is repeated, **then** at most one go-schedule indicator is present.

### Edge Cases

- The graphical session lacks a session bus or compatible status host, or the host appears after the item starts.
- The service is missing, transitional, or reported running without fresh local health.
- The installed unit retains a custom configuration path and IPC endpoint.
- Authorization is denied, cancelled, or unavailable; the service changes state during an action.
- Multiple GUI or indicator launch requests occur in one session.
- A status host changes panel theme or restarts while the item is active.

## Requirements

### Functional Requirements

- **FR-001**: A supported Linux graphical session with a compatible status host MUST show at most one local-daemon indicator independent of GUI lifetime and daemon process lifetime.
- **FR-002**: The indicator and GUI MUST share clear local state names for Running, Stopped, Not installed, Starting, Stopping, and Unknown/unreachable. Running MUST require both service-manager state and fresh local health.
- **FR-003**: Indicator activation MUST open or focus the GUI. The indicator menu MUST expose local status, Open, available Start/Stop/Restart actions, and Quit indicator.
- **FR-004**: Start/Stop/Restart MUST target the installed local system service only. Stop/Restart MUST explain impact and require confirmation. Authorization denial, cancellation, failure, and timeout MUST never be presented as success.
- **FR-005**: The GUI MUST provide local service status and actions even without a status host and while a remote daemon is selected. Opening the GUI MUST NOT auto-start a deliberately stopped installed service.
- **FR-006**: A missing session bus, absent host, or unsupported desktop MUST leave daemon scheduling and GUI operation intact and offer an alternate way to inspect status.
- **FR-007**: Session startup, shutdown, and removal MUST not leave stale registrations or duplicate items. The item MUST use a compact, legible approved mark and plain-text status.
- **FR-008**: Documentation MUST name the supported desktop contract, explain authorization and status-host limitations, and keep run-outcome notifications separate.

### Key Entities

- **Local service observation**: Service-manager state, fresh local health, user-visible state, detail, observation time.
- **Local service action**: Requested operation, confirmation and authorization outcome, observed final state.
- **Session indicator**: One user-session status item and menu, independent of the GUI and daemon lifetimes.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A supported signed-in session has no more than one go-schedule indicator after repeated launches and host restarts.
- **SC-002**: Every defined local service condition maps to a truthful state in both indicator and GUI; stale health never appears as Running.
- **SC-003**: Every action reaches an observed terminal state or reports a bounded failure/cancellation without claiming success prematurely.
- **SC-004**: Repeated Open requests result in one GUI window and do not reverse a deliberate service stop.
- **SC-005**: A missing status host or session bus does not prevent daemon scheduling or GUI status/control access.

## Assumptions

- The supported indicator contract is Linux graphical sessions with a D-Bus session bus and a compatible StatusNotifier host. Other desktops receive GUI and CLI fallback, not a promise of a visible tray icon.
- Service mutations use the system's existing authorization mechanism; the GUI and indicator remain unprivileged.
- The installed service is system-wide and continues running without a logged-in desktop session.
- Native run-outcome desktop notifications (#177) and SMTP (#176) are out of scope.
