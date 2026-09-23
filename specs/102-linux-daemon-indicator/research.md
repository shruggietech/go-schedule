# S102 research and decisions

## StatusNotifier integration

**Decision**: Use `github.com/gogpu/systray` v0.3.0 for the Linux-only indicator. Its Linux implementation uses D-Bus StatusNotifierItem and dbusmenu, supports menu item updates and host restart registration, and has an MIT license. The repository already depends on godbus indirectly.

**Alternatives**: `shelepuginivan/systray` is a host implementation, not an application item; `knightpp/sni` supplies SNI/menu but has lower recent activity; hand-written SNI/dbusmenu would duplicate a protocol stack. The older CGO systray variants would add toolkit coupling.

**Source**: [StatusNotifierItem specification](https://specifications.freedesktop.org/status-notifier-item/latest-single/) and [gogpu/systray](https://github.com/gogpu/systray).

## Linux service control

**Decision**: Use a bounded systemd read for installed, active, and transition state, and the existing local daemon IPC health for final Running. Use fixed systemctl operations through graphical Polkit (`pkexec`) when an unprivileged user changes state. The GUI and indicator remain unprivileged.

**Alternatives**: Launching `sudo` would ask for a terminal and fail in the tray. Elevating the GUI or indicator would grant broad UI-process authority. Bare `systemctl` may depend on desktop-specific authorization behavior.

**Review refinement**: A service installed with `--config` can use a nondefault `data_dir` or `ipc_path`. Read the effective `ExecStart` argument vector from systemd's system D-Bus property, including drop-in overrides, and load that config for both desktop processes. Parsing `systemctl show` prose or assuming `/var/lib/goschedule/config.json` can select the wrong daemon. A custom config must be readable by the signed-in desktop user; failure to resolve an installed unit is an error, not a false healthy state.

## Session and distribution

**Decision**: Use an XDG autostart desktop entry alongside the portable Linux desktop archive; instruct users and downstream packages to install it where appropriate. A per-user file lock prevents duplicate indicator or GUI instances. A local activation socket allows focus without a D-Bus session.

**Alternatives**: A system service cannot display user-session UI. A user-systemd unit is less universal than XDG autostart and still needs a graphical session. A D-Bus-only singleton would deny the GUI fallback when no bus exists.
