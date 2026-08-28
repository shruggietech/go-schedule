# Windows installer verification

These tools separate compiled MSI evidence from observations that only a clean
Windows desktop can provide. They are maintainer procedures, not release
automation, and they never count an unavailable prerequisite as a pass.

## Inspect an MSI without installing it

On Windows with PowerShell 7:

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath C:\path\to\go-schedule_vX.Y.Z_windows_amd64.msi `
  -EvidencePath C:\path\to\artifact-evidence.md
```

The command reads the MSI's Icon, Property, Shortcut, and Environment tables.
It does not install, repair, uninstall, start a service, or change PATH.

## Run the clean lifecycle

Use a fresh Windows VM, Windows Sandbox, or a physical machine that has never
had go-schedule installed. Take a disposable snapshot first. Copy the MSI and
this directory into that environment, open an elevated PowerShell 7 session,
and run:

```powershell
pwsh .\install-lifecycle.ps1 `
  -MsiPath C:\verify\go-schedule_vX.Y.Z_windows_amd64.msi `
  -EvidencePath C:\verify\lifecycle-evidence.md `
  -PauseForNativeObservation
```

The script refuses to proceed if `gosched` already resolves, the install
directory exists, an installed-product entry exists, or machine PATH already
contains the install directory. It installs silently, probes commands with a
freshly composed machine/user PATH, reinstalls, and uninstalls. It always tries
to remove the package in `finally` after installation begins.

## Record native icon observations

While the package is installed during a clean lifecycle run, inspect these four
surfaces and append the result plus screenshots to the evidence record:

1. Start Menu shortcut.
2. Settings > Apps > Installed apps.
3. GUI window title area after launching go-schedule.
4. Taskbar entry while the GUI is running.

Use `proven`, `failed`, or `unavailable` for each surface. If Windows shows an
older cached icon, refresh the Start Menu/taskbar shell view or sign out and in,
then record both the initial and refreshed result. Do not repair or reinstall
the package merely to make a cached image change; that would blur package
behavior with shell-cache behavior.

## Evidence boundaries

- A locally built candidate can prove future MSI table contents.
- Issue #16 requires a downloaded published MSI and a clean environment.
- Issue #33 requires native window and taskbar observation.
- The `.msi` file's own Explorer icon is a Windows Installer surface and is not
  part of this verification.
