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

The installer creates or reuses the local **`goschedadmin`** group and adds the
interactive account that launched setup. Sign out and back in once after the
first install so Windows puts that new group SID into your login token. The
installer then:

- installs `gosched-gui.exe`, `goschedd.exe`, and `gosched.exe` to
  `C:\Program Files\go-schedule\`;
- registers **`goschedd`** as a Windows service set to **start automatically**,
  so your tasks run in the background and survive reboots even with no one
  logged in;
- adds `C:\Program Files\go-schedule\` to the machine `PATH`, so `gosched` works
  as a command;
- grants daemon IPC access to SYSTEM, built-in Administrators, and members of
  `goschedadmin`;
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

The `goschedadmin` group and its memberships are also preserved. Other tools or
administrators may rely on them, so uninstall does not erase that shared OS
state. Remove them manually only if you are certain they are no longer wanted.

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

**The installer rolls back or reports error 1603.** Retry from an elevated
PowerShell window with verbose logging so the failed Windows Installer action is
preserved:

```powershell
$log = Join-Path $env:TEMP 'go-schedule-install.log'
$packages = @(Get-ChildItem -LiteralPath . `
  -Filter 'go-schedule_*_windows_amd64.msi' -File)
if ($packages.Count -ne 1) {
  throw "Expected exactly one go-schedule MSI; found $($packages.Count)."
}
& msiexec.exe /i $packages[0].FullName /L*v $log
Write-Output "Verbose installer log: $log"
```

Search that log for `Return value 3`, `0x80070005`, error `26421`,
`CreateGroup`, and `StartServices`. Confirm `Get-LocalGroup goschedadmin` works,
your account appears in `Get-LocalGroupMember goschedadmin`, and
`Get-Service goschedd` reports the service state. Error 26421 with
`0x80070005` means group provisioning was denied; on a current installer, check
local security policy or account-management restrictions and retain the log for
support. The v0.9.0 installer had a known authoring defect that produced this
combination even for local administrators; use a newer candidate or release.

**SmartScreen or antivirus warning.** The installer is currently unsigned.
Verify the SHA-256 hash against `SHA256SUMS.txt` and choose *More info → Run
anyway* if it matches.

**The GUI opens but says "daemon unreachable".** Check the service with
`gosched service status`. If it reports `stopped`, start it from an
**elevated** shell with `gosched service start`. A current daemon authorizes a
direct `goschedadmin` user by the stable user SID as well as by the group SID,
so a first launch does not require permanent elevation or a sign-out solely to
refresh group claims.

Current versions keep the application frame visible when the daemon connection
fails. One connection panel replaces repeated error dialogs and always offers
**Retry** and **Exit**. The panel distinguishes a missing or stopped daemon, a
timeout, and an IPC authorization failure. **Access denied** means Windows
rejected the current process token at the existing named pipe; it does not by
itself mean the service is stopped.

For a first-install access denial, collect these four observations in a normal,
non-elevated PowerShell window:

```powershell
Get-Service goschedd | Select-Object Name, Status, StartType
Get-LocalGroup goschedadmin | Select-Object Name, SID
Get-LocalGroupMember goschedadmin | Select-Object Name, ObjectClass, PrincipalSource
whoami /groups | Select-String -Pattern 'goschedadmin|S-1-5-32|Group Name'
```

If the local group exists and your account appears in its membership but
`whoami /groups` does not show the group's SID, the daemon's restricted pipe
still includes the stable SID of every direct user member discovered when the
service starts. Restart the service after a membership change, then choose
**Retry**. A fresh sign-in remains useful for normal Windows group-token
refresh, including nested-group membership, but is not required for the direct
installing-user path. If the account is not listed in
`Get-LocalGroupMember`, retain the installer log described above and ask an
administrator to verify installation and local account policy. Do not weaken
the pipe ACL or run the desktop application permanently elevated as a workaround.

**The service reports an `admin_group` lookup error.** Confirm the group exists
with `Get-LocalGroup goschedadmin` and that your account appears in
`Get-LocalGroupMember goschedadmin`. A foreground daemon can deliberately use
the former broad local policy by passing `--config <path>` with
`{"admin_group":""}` after stopping the service. That compatibility mode admits
Authenticated Users and emits a startup warning; normal MSI service installs
use the secure group default.

**Where are the logs?** `C:\ProgramData\goschedule\logs\goschedule.log` and its
rotated siblings, or the **Activity** view in the app, or `gosched logs`.

**A task fails on Windows.** The daemon is a noninteractive LocalSystem service,
so do not assume the interactive user's profile, mapped drives, PATH additions,
or current directory are available. Put only the executable in **Command**, put
one argument per line in **Arguments**, prefer an absolute executable path, and
choose a working directory and output path accessible to the service. A child
that starts and exits nonzero retains its exit code and output. A child that
cannot start has no exit code and reports `process start failed for
"<executable>"` with the Windows error. Arguments, stdin, and environment
values are omitted from that diagnostic because they may contain secrets.
