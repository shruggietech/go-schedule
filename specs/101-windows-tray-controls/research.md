# S101 research

## Per-session architecture

**Decision**: A separate windowless companion owns the notification-area icon. The Windows service remains headless and Wails remains optional.
**Rationale**: Windows services do not safely display interactive UI in the signed-in user's session, and closing Wails must not remove the icon.
**Alternatives considered**: Service-owned UI and Wails-owned tray both fail the required lifecycle.
**Sources**: [Microsoft interactive services](https://learn.microsoft.com/en-us/windows/win32/services/interactive-services), [notification area](https://learn.microsoft.com/en-us/windows/win32/shell/notification-area).

## Status and recovery

**Decision**: Query the SCM with status-only rights, preserve pending states, and probe local IPC with a short timeout. Treat SCM Running plus failed/stale health as Unknown/unreachable. Restore the icon on TaskbarCreated.
**Rationale**: Existing service status collapses pending states and has no health freshness. Explorer recreation discards notification-area registrations.
**Alternatives considered**: Existing kardianos Status cannot preserve pending states; daemon health alone cannot prove installed service state.
**Sources**: [Shell_NotifyIcon](https://learn.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shell_notifyiconw), [Windows taskbar](https://learn.microsoft.com/en-us/windows/win32/shell/taskbar).

## Privileged actions

**Decision**: Use an asInvoker tray and GUI with a narrow elevated service-action mode invoked by runas. Confirmation happens before Stop/Restart. The caller re-reads SCM and IPC to establish final result.
**Rationale**: Standard users can inspect status but not mutate the service; a broad elevated UI would unnecessarily increase risk.
**Alternatives considered**: Elevating the whole companion or permanently expanding the service ACL.
**Source**: [Microsoft least-privilege elevation guidance](https://learn.microsoft.com/en-us/windows/win32/secbp/running-with-administrator-privileges).

## Installation and launch

**Decision**: Install the companion with the MSI and register logon launch for signed-in users; retire running instances and autorun on upgrade/uninstall. GUI activation is session-scoped.
**Rationale**: An icon must exist without first launching Wails, and uninstall must not leave a stale icon or registry entry.
**Alternatives considered**: GUI-only startup and daemon startup both fail one of the lifecycle requirements.
**Source**: [Microsoft Run keys](https://learn.microsoft.com/en-us/windows/win32/setupapi/run-and-runonce-registry-keys).

## Small-size artwork

**Decision**: Export exact-geometry white and black versions of the canonical reduced mark into 16, 20, 24, and 32 px ICO resources, and add these named derivatives to the brand guide and integrity manifest.
**Rationale**: The existing approved monochrome marks use full-detail geometry that loses legibility at tray size. The reduced geometry is explicitly required for small surfaces. A geometry-regression test prevents accidental redraw.
**Alternatives considered**: Reusing the full-detail monochrome mark violates the small-size rule; using the full-color reduced mark alone loses contrast against some dark taskbars.

## Scope

**Decision**: Defer Linux tray (#256), native desktop notification delivery and optional silencing (#177), and SMTP (#176).
**Rationale**: These have distinct platform and UX contracts. S101 can deliver a complete Windows service-control capability independently.
