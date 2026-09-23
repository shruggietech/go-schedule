---
title: Linux
parent: Installation
nav_order: 2
---

# Installing go-schedule on Linux

**Audience:** Linux users and server operators\
**Applies to:** go-schedule 0.6.0 and later\
**See also:** [`gosched` command reference](cli.md) · [macOS](INSTALL-macos.md) · [Windows](INSTALL-windows.md)

Beginning with v1.4.0, releases include a Linux desktop archive beside the existing server archive. The server archive contains the cgo-free daemon and CLI. The desktop archive adds the native Wails application and freedesktop integration assets; it uses the system WebKitGTK runtime, while a headless server does not need desktop or WebView dependencies. The [latest-release page](https://github.com/shruggietech/go-schedule/releases/latest) is authoritative while a candidate is staged and qualified.

## Contents

- [Install](#install)
- [Desktop bundle](#desktop-bundle)
- [Register the service](#register-the-service)
- [First task](#first-task)
- [Where things live](#where-things-live)
- [Running without a service](#running-without-a-service)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Troubleshooting](#troubleshooting)

## Install

From the [latest release](https://github.com/shruggietech/go-schedule/releases/latest), download `go-schedule_<ver>_linux_<arch>.tar.gz` for your architecture, `amd64` or `arm64`.

Verify it against `SHA256SUMS.txt` before unpacking. This is the only integrity check you get; the archives are not signed.

```sh
sha256sum -c SHA256SUMS.txt --ignore-missing
```

```sh
tar -xzf go-schedule_*_linux_amd64.tar.gz
cd go-schedule_*_linux_amd64
```

Put the two binaries somewhere on the system `PATH` so both your shell and the service manager can find them:

```sh
sudo install -m 0755 goschedd gosched /usr/local/bin/
```

```sh
gosched --version
```

## Desktop bundle

For an amd64 workstation, download `go-schedule-desktop_<ver>_linux_amd64.tar.gz`. It contains `gosched-gui`, `gosched-indicator`, `goschedd`, `gosched`, and a `share/` tree with desktop entries and icons. The Wails application uses the system WebKitGTK 4.1 runtime. Keep the four binaries together or install them on `PATH`, then install the integration assets if wanted:

```sh
sudo install -m 0755 gosched-gui gosched-indicator goschedd gosched /usr/local/bin/
sudo cp -R share/applications share/icons /usr/local/share/
gosched gui
```

The `arm64` release remains daemon-and-CLI only.

### Local daemon indicator

For a Linux graphical session with a D-Bus session bus and compatible StatusNotifier host, install the supplied per-user autostart entry and sign out and back in:

```sh
mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}/autostart"
cp share/autostart/go-schedule-indicator.desktop "${XDG_CONFIG_HOME:-$HOME/.config}/autostart/"
```

The item runs in your session, independent of the system daemon and GUI. Its menu shows the local service state, Open go-schedule, available Start/Stop/Restart actions, and Quit indicator. Running requires both an active system service and fresh local daemon health. Stop and Restart first ask for confirmation using `zenity` or `kdialog`; without either dialog helper, use the GUI's Connections page instead. An unprivileged change uses `pkexec` and the session's graphical Polkit authentication agent. Denial, cancellation, or a missing agent leaves the actual service state visible. Quit indicator never stops the daemon.

KDE Plasma, Xfce, Cinnamon, MATE, and GNOME with a StatusNotifier/AppIndicator extension are potential hosts, but availability depends on the installed panel and session configuration. A stock GNOME Shell without an extension and a session without a D-Bus bus do not promise a visible item. The GUI Connections page and `gosched service status` remain the fallback, and the daemon continues scheduling without any desktop session. This indicator is long-lived service presence, not native task-outcome notifications.

To stop launching the item at login, remove only your installed autostart copy and use Quit indicator for the current session. Removing the desktop bundle should also remove its binary and any installed autostart copy; the system daemon and scheduler data are separate.

## Register the service

The daemon runs as a systemd service so the scheduler starts on boot and keeps running with nobody logged in. Registration writes the unit file, which is why it needs root:

Create the dedicated administrative group and add the account that will use the CLI or desktop app before the daemon's first start:

```sh
getent group goschedadmin >/dev/null || sudo groupadd --system goschedadmin
sudo usermod -aG goschedadmin "$USER"
```

Sign out and back in after changing membership. The daemon fails closed if its configured non-empty group does not exist. It creates its default data directory with group `goschedadmin` and mode `0770`, and its socket with mode `0660`.

```sh
sudo gosched service install
```

The service loads `/var/lib/goschedule/config.json` automatically when it exists. For a different location, install it with `sudo gosched service install --config /etc/goschedule/daemon.json`. The file must be readable by both the service identity and any signed-in user who uses the desktop GUI or indicator, since those clients read the service unit's effective `--config` argument to find its local IPC endpoint. A group-readable file is sufficient; the referenced TLS private key need only be readable by the service identity. Follow [Remote access](remote-access.md#operator-runbook) to enable the optional HTTPS listener; installation alone never opens one.

```sh
sudo gosched service start
```

```sh
gosched health
```

Expect `daemon ok (version …)`. If you get a connection error instead, the daemon is not running, check `gosched service status`, then the [troubleshooting notes](#troubleshooting).

`service status` does not require root. `install`, `uninstall`, `start`, `stop`, and `restart` do.

## First task

```sh
gosched task add nightly-backup \
  --command /usr/local/bin/backup.sh \
  --schedule "every day at 02:30" \
  --tz Europe/London
```

The command echoes back how it understood the schedule along with the next few run times, so a misreading shows up now rather than at 02:30 tomorrow. Prove it end to end without waiting:

```sh
gosched task run-now <id>
```

```sh
gosched runs --task <id>
```

The full command set is in the [reference](cli.md).

## Where things live

| What | Path |
| --- | --- |
| Database | `/var/lib/goschedule/goschedule.db` |
| Logs | `/var/lib/goschedule/logs/goschedule.log` (plus rotated siblings) |
| IPC socket | under `/var/lib/goschedule/` |

The data directory is created on first run. Removing the binaries does not remove it, which is deliberate, a reinstall keeps your tasks.

## Running without a service

For a quick trial you can run the daemon in the foreground and leave it in a terminal:

```sh
goschedd
```

It serves the same IPC endpoint, so `gosched` works against it normally. Nothing survives a reboot this way, and a single-instance lock stops a second daemon from starting alongside the service, which is what you want, but it does mean you should stop the service first if one is installed.

## Upgrading

Stop the service, replace the binaries, start it again. The database migrates forward automatically on first start, non-destructively.

```sh
sudo gosched service stop
sudo install -m 0755 goschedd gosched /usr/local/bin/
sudo gosched service start
gosched health
```

The service definition retains a custom `--config` argument, and the default `/var/lib/goschedule/config.json` remains in the preserved data directory. An upgrade does not enable remote access or replace certificate and credential state.

## Uninstalling

```sh
sudo gosched service stop
sudo gosched service uninstall
sudo rm /usr/local/bin/goschedd /usr/local/bin/gosched
```

If you installed the desktop bundle, first choose Quit indicator, then remove `gosched-gui` and `gosched-indicator` from the directory where you installed them. Remove your own `${XDG_CONFIG_HOME:-$HOME/.config}/autostart/go-schedule-indicator.desktop` copy if you enabled session startup. These steps do not remove the saved scheduler database.

Your data is left in place. To remove it as well:

```sh
sudo rm -rf /var/lib/goschedule
```

## Troubleshooting

**`gosched: command not found`.** The binaries are not on `PATH`. Either install them to `/usr/local/bin` as above, or invoke them by path.

**`service install` fails with a permission error.** It writes a systemd unit; run it with `sudo`.

**`gosched health` reports the daemon unreachable.** Check `gosched service status`. If it says `stopped`, start it. If it says the service is not installed, install it. If it says `running` but health still fails, read `gosched logs --severity error`, a daemon that failed at startup exits non-zero and says why.

**The Linux indicator is missing.** Confirm a D-Bus graphical session and a compatible StatusNotifier host, plus the installed `gosched-indicator` binary, compact icon tree, and autostart entry. Start `gosched-indicator` from a terminal to see a specific startup error. The Connections page and `gosched service status` remain available even without a host.

**The daemon reports an `admin_group` lookup or permission error.** Confirm the group and your fresh login membership with `getent group goschedadmin` and `id -nG`. For a custom `ipc_path`, its existing parent must already belong to the configured group with exact mode `0770`; the daemon does not rewrite a custom directory. As a temporary single-user compatibility choice, an operator can launch `goschedd --config <path>` with this explicit overlay:

```text
{"admin_group":""}
```

That mode deliberately restores a `0666` socket and emits a startup warning. Registered services use the secure default unless their service definition is explicitly given the config argument.

**Tasks run but with the wrong environment.** The service runs as root with a minimal environment, not as your login shell. Set what a task needs explicitly with `--env` and `--cwd` rather than relying on inherited state.

**Times drift by an hour twice a year.** They should not, schedules resolve in the task's IANA timezone, and Daylight Saving transitions are handled explicitly. If you see this, set `--tz` on the task rather than depending on the system zone, and please file a report.
