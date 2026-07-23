---
title: Linux
parent: Installation
nav_order: 2
---

# Installing go-schedule on Linux

**Audience:** Linux users and server operators\
**Applies to:** go-schedule 0.6.0 and later\
**See also:** [`gosched` command reference](cli.md) ·
[macOS](INSTALL-macos.md) · [Windows](INSTALL-windows.md)

Linux ships as a release archive holding the daemon and the CLI. Both are
cgo-free, so there is nothing to compile and no runtime dependency to satisfy.
The desktop GUI is built separately and needs a C toolchain and OpenGL; a
headless server does not want it and does not need it.

## Contents

- [Install](#install)
- [Register the service](#register-the-service)
- [First task](#first-task)
- [Where things live](#where-things-live)
- [Running without a service](#running-without-a-service)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Troubleshooting](#troubleshooting)

## Install

From the [latest release](https://github.com/shruggietech/go-schedule/releases/latest),
download `go-schedule_<ver>_linux_<arch>.tar.gz` for your architecture — `amd64`
or `arm64`.

Verify it against `SHA256SUMS.txt` before unpacking. This is the only integrity
check you get; the archives are not signed.

```sh
sha256sum -c SHA256SUMS.txt --ignore-missing
```

```sh
tar -xzf go-schedule_*_linux_amd64.tar.gz
cd go-schedule_*_linux_amd64
```

Put the two binaries somewhere on the system `PATH` so both your shell and the
service manager can find them:

```sh
sudo install -m 0755 goschedd gosched /usr/local/bin/
```

```sh
gosched --version
```

## Register the service

The daemon runs as a systemd service so the scheduler starts on boot and keeps
running with nobody logged in. Registration writes the unit file, which is why
it needs root:

```sh
sudo gosched service install
```

```sh
sudo gosched service start
```

```sh
gosched health
```

Expect `daemon ok (version …)`. If you get a connection error instead, the
daemon is not running — check `gosched service status`, then the
[troubleshooting notes](#troubleshooting).

`service status` does not require root. `install`, `uninstall`, `start`, `stop`,
and `restart` do.

## First task

```sh
gosched task add nightly-backup \
  --command /usr/local/bin/backup.sh \
  --schedule "every day at 02:30" \
  --tz Europe/London
```

The command echoes back how it understood the schedule along with the next few
run times, so a misreading shows up now rather than at 02:30 tomorrow. Prove it
end to end without waiting:

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

The data directory is created on first run. Removing the binaries does not
remove it, which is deliberate — a reinstall keeps your tasks.

## Running without a service

For a quick trial you can run the daemon in the foreground and leave it in a
terminal:

```sh
goschedd
```

It serves the same IPC endpoint, so `gosched` works against it normally. Nothing
survives a reboot this way, and a single-instance lock stops a second daemon
from starting alongside the service — which is what you want, but it does mean
you should stop the service first if one is installed.

## Upgrading

Stop the service, replace the binaries, start it again. The database migrates
forward automatically on first start, non-destructively.

```sh
sudo gosched service stop
sudo install -m 0755 goschedd gosched /usr/local/bin/
sudo gosched service start
gosched health
```

## Uninstalling

```sh
sudo gosched service stop
sudo gosched service uninstall
sudo rm /usr/local/bin/goschedd /usr/local/bin/gosched
```

Your data is left in place. To remove it as well:

```sh
sudo rm -rf /var/lib/goschedule
```

## Troubleshooting

**`gosched: command not found`.** The binaries are not on `PATH`. Either install
them to `/usr/local/bin` as above, or invoke them by path.

**`service install` fails with a permission error.** It writes a systemd unit;
run it with `sudo`.

**`gosched health` reports the daemon unreachable.** Check
`gosched service status`. If it says `stopped`, start it. If it says the service
is not installed, install it. If it says `running` but health still fails, read
`gosched logs --severity error` — a daemon that failed at startup exits non-zero
and says why.

**Tasks run but with the wrong environment.** The service runs as root with a
minimal environment, not as your login shell. Set what a task needs explicitly
with `--env` and `--cwd` rather than relying on inherited state.

**Times drift by an hour twice a year.** They should not — schedules resolve in
the task's IANA timezone, and Daylight Saving transitions are handled
explicitly. If you see this, set `--tz` on the task rather than depending on the
system zone, and please file a report.
