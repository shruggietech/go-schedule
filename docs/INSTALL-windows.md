---
title: Windows
parent: Installation
nav_order: 1
---

# Installing go-schedule on Windows

**Audience:** Windows users installing go-schedule\
**Applies to:** go-schedule 0.6.0 and later\
**See also:** [`gosched` command reference](cli.md) ·
[Linux](INSTALL-linux.md) · [macOS](INSTALL-macos.md)

go-schedule installs as a formal Windows application via an `.msi` package. It
puts the program in *Program Files*, runs the scheduler as an auto-starting
**Windows service**, adds the install directory to `PATH`, and adds a Start-Menu
shortcut. There is no "extract a zip and run an exe from Downloads" step.

## Contents

- [Install](#install)
- [Using the CLI](#using-the-cli)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Troubleshooting](#troubleshooting)

## Install

1. From the [latest release](https://github.com/shruggietech/go-schedule/releases/latest),
   download **`go-schedule_<ver>_windows_amd64.msi`**.

2. (Recommended) Verify the download against `SHA256SUMS.txt`:

   ```powershell
   Get-FileHash .\go-schedule_*_windows_amd64.msi -Algorithm SHA256
   ```

3. **Double-click the `.msi`** and complete the wizard. Windows prompts for
   administrator approval (UAC) — this is required because the installer
   registers a system service. Approve it.

That is all. The installer:

- installs `gosched-gui.exe`, `goschedd.exe`, and `gosched.exe` to
  `C:\Program Files\go-schedule\`;
- registers **`goschedd`** as a Windows service set to **start automatically**,
  so your tasks run in the background and survive reboots even with no one
  logged in;
- adds `C:\Program Files\go-schedule\` to the machine `PATH`, so `gosched` works
  as a command;
- adds a **go-schedule** shortcut to the Start Menu.

Launch **go-schedule** from the Start Menu to open the desktop app. It connects
to the already-running service rather than starting a second copy, and it never
shows a console window.

> **Data location:** tasks and logs live under `C:\ProgramData\goschedule\` — the
> database `goschedule.db` and the `logs\` folder. Both are created
> automatically on first run.

## Using the CLI

The CLI is installed alongside the app and is on your `PATH`. Open a **new**
PowerShell window — an already-open shell inherited its environment before the
install and will not see the new `PATH` entry until it is restarted.

```powershell
gosched health
```

```powershell
gosched service status
```

```powershell
gosched task add backup `
  --command "C:\Windows\System32\cmd.exe" --arg "/c" --arg "echo backup" `
  --schedule "every weekday at 09:00"
```

```powershell
gosched task list
```

```powershell
gosched logs --severity error
```

`service status` works from an ordinary, non-elevated shell. The subcommands
that change the service — `install`, `uninstall`, `start`, `stop`, `restart` —
require an elevated one. Full detail in the
[command reference](cli.md#service).

If you would rather not open a new window, the full path works in the shell you
already have:

```powershell
& "C:\Program Files\go-schedule\gosched.exe" health
```

## Upgrading

Download the newer `.msi` and run it. It performs an in-place major upgrade: the
old version is removed and the new one installed, your `PATH` entry is replaced
rather than duplicated, and your data under `C:\ProgramData\goschedule\` is
preserved.

## Uninstalling

Use **Settings → Apps → Installed apps → go-schedule → Uninstall**. This stops
and removes the service, deletes the program files, removes the `PATH` entry,
and removes the Start-Menu shortcut.

Your data under `C:\ProgramData\goschedule\` is **left in place**, so a later
reinstall keeps your tasks. To remove it completely, delete that folder after
uninstalling:

```powershell
Remove-Item -Recurse -Force "C:\ProgramData\goschedule"
```

## Troubleshooting

**`gosched` is not recognized as a command.** The `PATH` entry is added at
install time but is not broadcast into shells that were already open. Close the
window and open a new one. If a brand-new shell still cannot find it, check that
the entry exists:

```powershell
([Environment]::GetEnvironmentVariable('Path','Machine') -split ';') |
  Where-Object { $_ -like '*go-schedule*' }
```

Versions before 0.6.0 did not add the entry at all; upgrading fixes it.

**UAC prompt on install.** Expected — registering a system service needs
elevation. Declining cancels the install cleanly, leaving nothing behind.

**SmartScreen or antivirus warning.** The installer is currently unsigned.
Verify the SHA-256 hash against `SHA256SUMS.txt` and choose *More info → Run
anyway* if it matches.

**The GUI opens but says "daemon unreachable".** Check the service with
`gosched service status`. If it reports `stopped`, start it from an
**elevated** shell with `gosched service start`.

**Where are the logs?** `C:\ProgramData\goschedule\logs\goschedule.log` and its
rotated siblings, or the **Logs** view in the app, or `gosched logs`.
