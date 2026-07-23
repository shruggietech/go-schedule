---
title: macOS
parent: Installation
nav_order: 3
---

# Installing go-schedule on macOS

**Audience:** macOS users, both desktop and headless\
**Applies to:** go-schedule 0.6.0 and later\
**See also:** [`gosched` command reference](cli.md) ·
[Linux](INSTALL-linux.md) · [Windows](INSTALL-windows.md)

macOS has two downloads, and which one you want depends on whether you want a
desktop app or a background scheduler on a machine you mostly reach over SSH.
The difference that matters is not the GUI — it is **what starts on boot**.

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
| `go-schedule-desktop_<ver>_darwin_<arch>` | GUI, daemon, and CLI in one `.app` | Not until you register the service |
| `go-schedule_<ver>_darwin_<arch>.tar.gz` | Daemon and CLI | Once you register the service |

Both are available for `amd64` (Intel) and `arm64` (Apple silicon). Verify
either against `SHA256SUMS.txt` before opening it — the builds are not signed or
notarized.

```sh
shasum -a 256 -c SHA256SUMS.txt --ignore-missing
```

## Desktop bundle

`go-schedule-desktop_<ver>_darwin_<arch>` contains `gosched-gui.app`, with the
daemon and CLI inside it at
`gosched-gui.app/Contents/MacOS/`. Open it:

```sh
open gosched-gui.app
```

The first launch starts the background daemon itself, so there is nothing to
configure. That daemon keeps running after you close the window, so tasks
continue to fire.

> **It does not survive a reboot.** An auto-started daemon is a plain background
> process, not a registered service. If you want the scheduler running after the
> machine restarts — which is the usual reason to want a scheduler — register
> the service as well. That is the single most common surprise on macOS, and it
> is worth doing at install time rather than discovering it later.

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

The daemon runs under `launchd` so it starts at boot. Registration writes a
launch daemon plist, which is why it needs root:

```sh
sudo gosched service install
```

```sh
sudo gosched service start
```

```sh
gosched health
```

Expect `daemon ok (version …)`.

`service status` does not require root. `install`, `uninstall`, `start`, `stop`,
and `restart` do.

If you installed the desktop bundle and had already been using its auto-started
daemon, register the service and let it take over — a single-instance lock stops
two daemons running at once, so the second one to start simply refuses rather
than corrupting anything.

## First task

```sh
gosched task add nightly-backup \
  --command /usr/local/bin/backup.sh \
  --schedule "every day at 02:30" \
  --tz America/Los_Angeles
```

The command echoes back how it understood the schedule and the next few run
times. Prove it works without waiting for one:

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
| Database | `/Library/Application Support/goschedule/goschedule.db` |
| Logs | `/Library/Application Support/goschedule/logs/goschedule.log` |
| IPC socket | under `/Library/Application Support/goschedule/` |

Created on first run, and left alone when you remove the binaries.

## Upgrading

For the desktop bundle, replace the `.app` with the new one. For the archive,
stop the service, replace the binaries, start it again:

```sh
sudo gosched service stop
sudo install -m 0755 goschedd gosched /usr/local/bin/
sudo gosched service start
gosched health
```

The database migrates forward automatically and non-destructively on first
start.

## Uninstalling

```sh
sudo gosched service stop
sudo gosched service uninstall
sudo rm -f /usr/local/bin/goschedd /usr/local/bin/gosched
```

Then delete `gosched-gui.app` if you installed the bundle. Data is left in
place; to remove it too:

```sh
sudo rm -rf "/Library/Application Support/goschedule"
```

## Troubleshooting

**"go-schedule cannot be opened because the developer cannot be verified."** The
builds are unsigned. Verify the SHA-256 hash against `SHA256SUMS.txt`, then
right-click the app and choose *Open*, or clear the quarantine attribute:

```sh
xattr -d com.apple.quarantine gosched-gui.app
```

**Tasks stop firing after a reboot.** The GUI's auto-started daemon is not a
service. Register it — see [above](#register-the-service).

**`gosched health` reports the daemon unreachable.** Check
`gosched service status` first. If it says `running` but health still fails,
read `gosched logs --severity error`; a daemon that failed at startup exits
non-zero and says why.

**Tasks run but cannot find a tool they need.** A launchd service starts with a
minimal environment, not your login shell's. Give the task what it needs
explicitly with `--env` and `--cwd` rather than relying on inherited state.
