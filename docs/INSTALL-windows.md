---
title: Windows
parent: Installation
nav_order: 1
---

# Installing go-schedule on Windows

**Audience:** Windows users installing go-schedule\
**Applies to:** MSI-based releases from 0.6.0 onward; S039 lifecycle controls are unreleased\
**See also:** [`gosched` command reference](cli.md) ·
[Linux](INSTALL-linux.md) · [macOS](INSTALL-macos.md)

> **Release status:** User-controlled shortcuts, completion actions, and the
> preserve-or-wipe uninstall flow described below are implemented for the next
> release after v0.9.1. They are not present in v0.9.1 or earlier installers.

go-schedule installs as a formal Windows application via an `.msi` package. It
puts the program in *Program Files*, runs the scheduler as an auto-starting
**Windows service**, adds the install directory to `PATH`, and lets you choose
Start Menu and desktop shortcuts. There is no "extract a zip and run an exe
from Downloads" step.

## Contents

- [Install](#install)
- [Installer choices](#installer-choices)
- [Using the CLI](#using-the-cli)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Managed and silent setup](#managed-and-silent-setup)
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
interactive account that launched setup. The daemon authorizes that direct user
by its stable user SID, so the installing user can launch immediately without
signing out or running the desktop application elevated. A fresh sign-in is
only needed when normal Windows group-token refresh matters, such as for
nested-group membership. The installer then:

- installs `gosched-gui.exe`, `goschedd.exe`, and `gosched.exe` to
  `C:\Program Files\go-schedule\`;
- registers **`goschedd`** as a Windows service set to **start automatically**,
  so your tasks run in the background and survive reboots even with no one
  logged in;
- adds `C:\Program Files\go-schedule\` to the machine `PATH`, so `gosched` works
  as a command;
- grants daemon IPC access to SYSTEM, built-in Administrators, and members of
  `goschedadmin`;
- creates the shortcuts selected in the wizard.

## Installer choices

Fresh attended setup selects a **Start Menu shortcut** by default and leaves
the **desktop shortcut** unselected. You can choose either, both, or neither.
Windows Installer maintenance mode can add or remove either shortcut later.
Every installer-created shortcut uses the go-schedule icon and opens the
windowless desktop application.

After a successful fresh attended install, the completion page offers two
independent choices:

- **Launch go-schedule** is selected by default.
- **Open the documentation website** is unselected by default and uses the
  system's default HTTPS browser.

Either, both, or neither may be selected. Setup never launches the application
or browser after a silent install, repair, modification, upgrade, failed or
canceled operation, rollback, or uninstall.

Launch **go-schedule** from either installed shortcut to open the desktop app.
It connects to the already-running service rather than starting a second copy,
and it never shows a console window.

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
preserved. Upgrades between installers that expose the shortcut choices retain
the installed feature state unless an administrator explicitly changes it.
The first upgrade from an older installer uses the new defaults because the
older package did not contain matching shortcut-feature identities.

## Uninstalling

Use **Settings -> Apps -> Installed apps -> go-schedule -> Modify**. Windows
Installer's direct **Uninstall** action can run with a reduced interface that
does not show package-authored choices, so go-schedule deliberately disables
that Settings action. **Modify** opens the full maintenance wizard. Select
**Remove** there to reach the removal inventory and these choices before any
software or data is removed:

- **Remove software and preserve application data** (default). This stops and
  removes the service, binaries, product registration, machine `PATH` entry,
  installer registry markers, and every selected installer-created shortcut.
  It preserves tasks, history, configuration, logs, runtime files, and desktop
  preferences for a later reinstall.
- **Remove software and erase application data**. This requires a separate,
  explicit confirmation. Cleanup begins only after software removal commits
  successfully.

Wipe covers the application-owned machine root
`C:\ProgramData\goschedule\` and the application preference leaf
`AppData\Roaming\fyne\tech.shruggie.goschedule` for every safely registered,
accessible local Windows profile on a fixed local volume. It does not follow
reparse points or accept path overrides. It does not delete exports,
administrator-configured paths outside the default roots, adjacent user files,
disconnected roaming copies, detached profile containers, or unregistered
profile directories.

If any candidate root is redirected or otherwise unsafe, wipe refuses before
deleting data. If a safe root cannot be completely removed, software removal
still completes and a protected report is retained at:

```text
C:\ProgramData\ShruggieTech\go-schedule-uninstall\
  b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34\cleanup-result.json
```

The report records only declared owned roots, outcomes, and Windows errors.
`HKLM\Software\ShruggieTech\go-schedule-uninstall` summarizes `State`,
`RemainingCount`, and `ReportPath`. A later successful wipe removes stale
cleanup evidence.

The `goschedadmin` group and its memberships are also preserved. Other tools or
administrators may rely on them, so uninstall does not erase that shared OS
state. Remove them manually only if you are certain they are no longer wanted.

The missing direct **Uninstall** action affects only the Windows application
list. Administrators and deployment tools can still use the silent commands
below. Running the `.msi` directly also opens the same maintenance wizard.

## Managed and silent setup

Windows Installer feature properties provide deterministic shortcut control.
The default installs the Start Menu shortcut only. These examples show the
other combinations:

```powershell
# Both shortcuts
msiexec.exe /i .\go-schedule_<ver>_windows_amd64.msi `
  ADDLOCAL=ALL /qn /norestart

# Neither shortcut
msiexec.exe /i .\go-schedule_<ver>_windows_amd64.msi `
  REMOVE=StartMenuShortcut,DesktopShortcut /qn /norestart

# Desktop only
msiexec.exe /i .\go-schedule_<ver>_windows_amd64.msi `
  ADDLOCAL=DesktopShortcut REMOVE=StartMenuShortcut /qn /norestart
```

The same `ADDLOCAL` and `REMOVE` feature names can modify an installed package.
Silent and reduced-UI setup never runs completion actions.

Silent uninstall preserves data unless the exact wipe opt-in is supplied:

```powershell
# Software only; application data is preserved
msiexec.exe /x .\go-schedule_<ver>_windows_amd64.msi `
  /qn /norestart /L*v .\go-schedule-uninstall.log

# Software and declared local application data
msiexec.exe /x .\go-schedule_<ver>_windows_amd64.msi `
  GOSCHEDULE_REMOVE_DATA=1 /qn /norestart /L*v .\go-schedule-wipe.log
```

Only `0` and `1` are valid values for `GOSCHEDULE_REMOVE_DATA`. The choice is
not persisted. Supplying the wipe property during install, repair, reinstall,
or major upgrade never schedules cleanup.

## Release-candidate safety

Starting with the next release after S040, a version tag stages a draft release
instead of immediately publishing it. The Windows installer, its candidate
manifest, the other platform assets, and the attended evidence stay non-public
until a maintainer completes the Windows 11 gate.

The gate binds the installed MSI to its repository, tag commit, staging
workflow run and attempt, ProductVersion, ProductCode, filename, byte size, and
SHA-256. It requires normal-user access, native window and DPI measurements,
two-minute connection-error observations, real manual and scheduled task runs,
and the attended setup and uninstall matrix. Missing, failed, unavailable,
skipped, timed-out, partial, stale, or altered evidence leaves the release in
draft state.

Promotion downloads and revalidates the same MSI bytes. It does not rebuild a
nominally equivalent installer after testing. Final checksums are created only
after the attended evidence archive joins the complete asset set. See
`test/windows/README.md` and
`specs/040-windows-release-candidate-gate/quickstart.md` for the maintainer
procedure. Creating a tag, promoting a draft, or publishing a release remains
a separately authorized maintainer action.

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

**Uninstall completed but data remains after a requested wipe.** Check the
protected cleanup report and summary registry key described under
[Uninstalling](#uninstalling). A `refused` result means validation found an
unsafe root and deleted nothing. A `partial` result names the declared owned
root that Windows would not remove, commonly because a process still has a file
open. Close the named application, retain the MSI log and cleanup report, and
retry removal or clean only the reported application-owned path after checking
that it is not a reparse point. Do not recursively remove a parent profile or
shared Fyne directory.

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
or current directory are available. In **Command line**, enter the executable
followed by its arguments, quote an executable path that contains spaces, prefer
an absolute path, and choose a working directory and output path accessible to
the service. The editor shows the exact Program and Arguments in order and does
not infer a shell; name `cmd`, PowerShell, or another shell explicitly when its
features are required. A child
that starts and exits nonzero retains its exit code and output. A child that
cannot start has no exit code and reports `process start failed for
"<executable>"` with the Windows error. Arguments, stdin, and environment
values are omitted from that diagnostic because they may contain secrets.
