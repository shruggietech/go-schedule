---
title: macOS
parent: Installation
nav_order: 3
---

# Installing go-schedule on macOS

**Audience:** macOS users, both desktop and headless\
**Applies to:** go-schedule 0.6.0 and later\
**See also:** [`gosched` command reference](cli.md) · [Linux](INSTALL-linux.md) · [Windows](INSTALL-windows.md)

macOS has two downloads, and which one you want depends on whether you want a desktop app or a background scheduler on a machine you mostly reach over SSH. The difference that matters is not the GUI, it is **what starts on boot**.

Beginning with v1.4.0, the desktop archive contains the Wails bundle described below. The [latest-release page](https://github.com/shruggietech/go-schedule/releases/latest) is authoritative while a candidate is staged and qualified.

## Contents

- [Which download](#which-download)
- [Desktop bundle](#desktop-bundle)
- [Daemon and CLI only](#daemon-and-cli-only)
- [Register the service](#register-the-service)
- [First task](#first-task)
- [Where things live](#where-things-live)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Troubleshooting](#troubleshooting)

## Which download

| Download | Contents | Starts on boot |
| --- | --- | --- |
| `go-schedule-desktop_<ver>_darwin_arm64.tar.gz` | Wails desktop, daemon, and CLI in one `.app` | Not until you register the service |
| `go-schedule_<ver>_darwin_<arch>.tar.gz` | Daemon and CLI | Once you register the service |

The desktop bundle targets `arm64` (Apple silicon). The daemon-and-CLI archive remains available for `amd64` (Intel) and `arm64`. Verify either against `SHA256SUMS.txt` before opening it, the builds are not signed or notarized.

```sh
shasum -a 256 -c SHA256SUMS.txt --ignore-missing
```

## Desktop bundle

`go-schedule-desktop_<ver>_darwin_arm64.tar.gz` contains the Wails application `gosched-gui.app`, with the daemon and CLI inside it at `gosched-gui.app/Contents/MacOS/`. Open it:

```sh
open gosched-gui.app
```

The first launch starts the background daemon itself, so there is nothing to configure beyond the one-time administrative-group setup below. That daemon keeps running after you close the window, so tasks continue to fire.

> **It does not survive a reboot.** An auto-started daemon is a plain background process, not a registered service. If you want the scheduler running after the machine restarts, which is the usual reason to want a scheduler, register the service as well. That is the single most common surprise on macOS, and it is worth doing at install time rather than discovering it later.

To use the bundled CLI from a terminal, either call it by path or link it:

```sh
sudo ln -s "$PWD/gosched-gui.app/Contents/MacOS/gosched" /usr/local/bin/gosched
```

## Daemon and CLI only

For a headless Mac, take the archive instead:

```sh
tar -xzf go-schedule_*_darwin_arm64.tar.gz
cd go-schedule_*_darwin_arm64
sudo install -m 0755 goschedd gosched /usr/local/bin/
gosched --version
```

## Register the service

The daemon runs under `launchd` so it starts at boot. Registration writes a launch daemon plist, which is why it needs root:

Before opening the desktop app or starting the service for the first time, create the dedicated group and add your account:

```sh
dseditgroup -o read goschedadmin >/dev/null 2>&1 || sudo dseditgroup -o create goschedadmin
sudo dseditgroup -o edit -a "$USER" -t user goschedadmin
```

Sign out and back in after changing membership. The daemon fails closed if its configured non-empty group does not exist. Its default data directory uses group `goschedadmin` and mode `0770`; the socket uses mode `0660`.

```sh
sudo gosched service install
```

The service loads `/Library/Application Support/goschedule/config.json` automatically when it exists. For a different location, install it with `sudo gosched service install --config "/Library/Application Support/goschedule/daemon.json"`. The file and any referenced TLS private key must be readable by the service identity. Follow [Remote access](remote-access.md#operator-runbook) to enable the optional HTTPS listener; installation alone never opens one.

```sh
sudo gosched service start
```

```sh
gosched health
```

Expect `daemon ok (version …)`.

`service status` does not require root. `install`, `uninstall`, `start`, `stop`, and `restart` do.

If you installed the desktop bundle and had already been using its auto-started daemon, register the service and let it take over, a single-instance lock stops two daemons running at once, so the second one to start simply refuses rather than corrupting anything.

## First task

Follow [Your first useful task](first-task.md) to create a supplied local file-snapshot script, run it immediately, inspect its exact output and repair an intentional missing-input failure. The POSIX instructions work on Linux and macOS with `/bin/sh`; no pre-existing backup script is required. Use paths accessible to the daemon identity. The complete CLI add creates an active schedule, so the guide disables automatic runs during the exercise.

## Where things live

| What | Path |
| --- | --- |
| Database | `/Library/Application Support/goschedule/goschedule.db` |
| Logs | `/Library/Application Support/goschedule/logs/goschedule.log` |
| IPC socket | under `/Library/Application Support/goschedule/` |

Created on first run, and left alone when you remove the binaries.

## Upgrading

Before replacing the `.app` or archive binaries, stop the daemon that owns the effective configuration and data directory. Closing the desktop window does not stop a GUI-started daemon. Stop the registered service with the command below, or stop the exact manually or GUI-started daemon process if no service owns it. Back up the complete stopped data directory (default `/Library/Application Support/goschedule`), including the database, any SQLite sidecar files, configuration and protected credential material, and retain the matching older bundle or binaries. A custom `--config` can select a different data directory.

For the desktop bundle, replace the `.app` only after stopping and backing up, then reopen it or start its registered service. For the daemon-and-CLI archive:

```sh
sudo gosched service stop
sudo install -m 0755 goschedd gosched /usr/local/bin/
sudo gosched service start
gosched health
```

Published v2.0.0 initializes schema 24 directly and reopens valid current stores without changing their data or identity. It provides no historical schema conversion or old-version rollback guarantee. Unsupported, partial, or foreign stores are refused without deletion. See [Current state and explicit reset](install.md#current-state-and-explicit-reset) for the deliberate development reset procedure. SMTP credentials remain target-local and are not exported by portable task bundles.

The service definition retains a custom `--config` argument, and the default configuration remains beside the preserved database. An upgrade does not enable remote access or replace certificate and credential state.

## Uninstalling

```sh
sudo gosched service stop
sudo gosched service uninstall
sudo rm -f /usr/local/bin/goschedd /usr/local/bin/gosched
```

Then delete `gosched-gui.app` if you installed the bundle. Data is left in place; to remove it too:

```sh
sudo rm -rf "/Library/Application Support/goschedule"
```

## Troubleshooting

**"go-schedule cannot be opened because the developer cannot be verified."** The builds are unsigned. Verify the SHA-256 hash against `SHA256SUMS.txt`, then right-click the app and choose *Open*, or clear the quarantine attribute:

```sh
xattr -d com.apple.quarantine gosched-gui.app
```

**Tasks stop firing after a reboot.** The GUI's auto-started daemon is not a service. Register it, see [above](#register-the-service).

**`gosched health` reports the daemon unreachable.** Check `gosched service status` first. If it says `running` but health still fails, read `gosched logs --severity error`; a daemon that failed at startup exits non-zero and says why.

**The daemon reports an `admin_group` lookup or permission error.** Verify the group with `dseditgroup -o read goschedadmin` and sign in again after membership changes. For a custom `ipc_path`, its existing parent must already belong to the configured group with exact mode `0770`; the daemon will not modify it. A foreground daemon can deliberately use the former broad local policy by passing `--config <path>` with `{"admin_group":""}`. This compatibility mode creates a `0666` socket and emits a startup warning; registered services remain on the secure default unless explicitly configured otherwise.

**Tasks run but cannot find a tool they need.** A launchd service starts with a minimal environment, not your login shell's. Give the task what it needs explicitly with `--env` and `--cwd` rather than relying on inherited state.
