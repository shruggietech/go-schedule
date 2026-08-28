# Windows installer verification

These tools separate compiled MSI evidence from observations that only a clean
Windows desktop can provide. They are maintainer procedures, not release
automation, and they never count an unavailable prerequisite as a pass.

## Inspect an MSI without installing it

On Windows with PowerShell 7:

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath C:\path\to\go-schedule_vX.Y.Z_windows_amd64.msi `
  -EvidencePath C:\path\to\artifact-evidence.md `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The command reads the MSI's Icon, Property, Shortcut, and Environment tables.
It does not install, repair, uninstall, start a service, or change PATH. For a
downloaded release asset, use `-ArtifactClass published` and supply its absolute
HTTPS release-asset URL as `-ArtifactOrigin`. The generated record keeps these
evidence classes distinct.

## Run the clean lifecycle

Use a fresh Windows VM, Windows Sandbox, or a physical machine that has never
had go-schedule installed. Take a disposable snapshot first. Copy the MSI and
this directory into that environment, open an elevated PowerShell 7 session,
and run:

```powershell
pwsh .\install-lifecycle.ps1 `
  -MsiPath C:\verify\go-schedule_vX.Y.Z_windows_amd64.msi `
  -EvidencePath C:\verify\lifecycle-evidence.md `
  -ArtifactClass published `
  -ArtifactOrigin 'https://github.com/shruggietech/go-schedule/releases/download/<tag>/<asset>' `
  -PauseForNativeObservation
```

The script refuses to proceed if `gosched` already resolves, the install
directory exists, an installed-product entry exists, or machine PATH already
contains the install directory. It installs silently, probes commands with a
freshly composed machine/user PATH, reinstalls, and uninstalls. It always tries
to remove the package in `finally` after installation begins. Failed runs write
cleanup results and the final product, directory, and PATH state into evidence.

## Record native icon observations

While the package is installed during a paused clean lifecycle run, inspect
these four surfaces. For each prompt, enter `proven`, `failed`, or `unavailable`,
then enter an optional note or screenshot reference:

1. Start Menu shortcut.
2. Settings > Apps > Installed apps.
3. GUI window title area after launching go-schedule.
4. Taskbar entry while the GUI is running.

The script retains those responses even when a later reinstall or uninstall
step fails. If Windows shows an older cached icon, refresh the Start Menu/taskbar
shell view or sign out and in, then record both the initial and refreshed result.
Do not repair or reinstall the package merely to make a cached image change;
that would blur package behavior with shell-cache behavior.

## Evidence boundaries

- A locally built candidate can prove future MSI table contents.
- Issue #16 requires a downloaded published MSI and a clean environment.
- Issue #33 requires native window and taskbar observation.
- The `.msi` file's own Explorer icon is a Windows Installer surface and is not
  part of this verification.
