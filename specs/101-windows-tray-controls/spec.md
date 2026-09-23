# Feature Specification: Windows local service tray and controls

**Feature Branch**: `codex/101-windows-tray-controls`
**Created**: 2026-09-22
**Status**: Implemented
**Delivery**: S101 review branch pending merge
**Input**: GitHub issue #255 and the approved S101 outline. Linux tray (#256), native desktop notifications (#177), and SMTP (#176) remain separate work.

## User Scenarios & Testing

### User Story 1 - See truthful local service status (Priority: P1)

As a Windows user, I can identify whether my local scheduler is running without opening the desktop GUI. The notification-area companion stays available for my signed-in session after I close the GUI.

**Why this priority**: Service presence is the foundation for every other control.
**Independent Test**: Install and launch the product, then observe the companion across service stop/start, GUI close, Explorer restart, and session restart.

**Acceptance Scenarios**:

1. **Given** an installed service with fresh local health, **when** I view the tray icon and tooltip, **then** I see Running.
2. **Given** a stopped or absent service, **when** I view the tray, **then** I see Stopped or Not installed, never a false Running.
3. **Given** a pending SCM transition or unreachable daemon, **when** I view the tray, **then** I see Starting, Stopping, or a clear Unknown/unreachable status.
4. **Given** the GUI is closed or Explorer restarts, **when** I return to the notification area, **then** one companion icon remains or recovers without a stale duplicate.

### User Story 2 - Control the local service (Priority: P1)

As a Windows user, I can start, stop, or restart the installed local service from the tray or desktop GUI without a flashing console. Protected actions request elevation only when needed and report their actual outcome.

**Why this priority**: A status-only icon does not deliver the approved management capability.
**Independent Test**: Exercise each action as a standard user and an administrator, including cancellation, timeout, and service failure.

**Acceptance Scenarios**:

1. **Given** the service is stopped, **when** I request Start and approve elevation if required, **then** the control waits for a healthy Running state or reports failure.
2. **Given** the service is running, **when** I request Stop and confirm, **then** the service stops and the GUI does not silently restart it.
3. **Given** an action is cancelled or denied, **when** the operation ends, **then** the actual state remains visible and no success is claimed.
4. **Given** a remote connection is selected in the GUI, **when** I inspect service controls, **then** the local service state and actions remain visible and clearly separate from the remote connection.

### User Story 3 - Open the GUI from the tray (Priority: P2)

As a Windows user, I can open the desktop GUI from the tray. If it is already open or hidden, that window is focused instead of a second GUI being created.

**Why this priority**: A predictable entry point avoids duplicate windows and confusing controls.
**Independent Test**: Repeatedly choose Open from the tray while the GUI is absent, visible, minimized, and hidden.

**Acceptance Scenarios**:

1. **Given** no GUI window exists, **when** I choose Open, **then** one GUI appears.
2. **Given** the GUI exists, **when** I choose Open, **then** that same window becomes visible and focused.
3. **Given** an installed service is deliberately stopped, **when** the GUI opens, **then** opening the GUI does not start the service.

### Edge Cases

- Explorer restarts, user session ends, or the companion is started twice in the same session.
- The service is not installed, is pending, is running but local IPC is stale, or has insufficient query rights.
- Elevation is refused, the service changes state while the menu is open, or an action exceeds its wait limit.
- The GUI is already active, minimized, or hidden; local and remote connections have different status.
- Installer upgrade or removal encounters a running companion.

## Requirements

### Functional Requirements

- **FR-001**: The Windows installation MUST provide one windowless notification-area companion per signed-in user session, independent of GUI lifetime, and remove its icon on exit.
- **FR-002**: The companion MUST restore its icon after Explorer restarts and avoid duplicate companion icons in one session.
- **FR-003**: Tray icon, tooltip, and menu MUST communicate Running, Stopped, Not installed, Starting, Stopping, and Unknown/unreachable, and MUST use approved reduced go-schedule brand artwork with suitable light, dark, and high-contrast variants and small-size legibility.
- **FR-004**: Running MUST require both the Windows service controller's running state and a fresh successful local daemon health response.
- **FR-005**: Tray and GUI MUST offer Open, local service status/details, Start, Stop, Restart, and Quit companion where applicable. Disabled or unavailable actions MUST be explained.
- **FR-006**: Stop and Restart MUST request confirmation. Privileged changes MUST use a narrow elevation path and never elevate the whole GUI or companion.
- **FR-007**: A control MUST show pending progress and report only the observed final state or an explicit cancellation/failure. It MUST NOT optimistically claim success.
- **FR-008**: All S101-launched Windows child processes MUST be noninteractive and hidden, except the intended UAC consent prompt and GUI window.
- **FR-009**: Open MUST focus an existing GUI instance, including a minimized or hidden one, rather than create another instance.
- **FR-010**: The GUI MUST NOT automatically restart a deliberately stopped installed service. Standalone GUI behavior without an installed service may continue.
- **FR-011**: GUI local service status and controls MUST remain visible independently of the currently selected local or remote daemon connection.
- **FR-012**: Installation, upgrade, and removal MUST provision, start for future user logons, and retire the companion without leaving stale autorun registration or binaries. User-facing documentation MUST describe controls and elevation.

### Key Entities

- **Local service snapshot**: SCM state, local health freshness, combined user-visible state, observed time, and error.
- **Service action**: Requested operation, confirmation/elevation outcome, pending state, and observed final result.
- **Session companion**: One notification-area process and icon associated with a signed-in Windows session, separate from the daemon and GUI.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Each signed-in session exposes at most one go-schedule tray icon, including after repeated launch and Explorer restart.
- **SC-002**: All defined service conditions map to a truthful status in both tray and GUI; a stale health response never displays Running.
- **SC-003**: Start, Stop, and Restart show an observed terminal result or explicit failure/cancellation within a bounded wait.
- **SC-004**: Repeated Open commands result in one foreground GUI, and opening it never reverses a deliberate installed-service stop.
- **SC-005**: Standard-user use creates no unexpected visible console windows, and upgrade/removal leave no stale companion registration.

## Assumptions

- Windows installed service and local IPC are the existing source of service and health truth. Remote daemon status does not substitute for local service status.
- A separate per-session companion is required because the service runs noninteractively and the GUI may be closed.
- Windows UAC is the authority boundary for service mutations. A user may decline elevation without changing the service.
- S101 does not deliver Linux tray behavior, native desktop notification delivery, or SMTP.
