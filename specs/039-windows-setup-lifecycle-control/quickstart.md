# Quickstart: Validate Windows Setup Lifecycle Control

## Prerequisites

- Windows 11 x64 for MSI build and native lifecycle checks.
- Go and WiX Toolset 6.0.2 with UI and Util extensions.
- An administrator-capable disposable host for installation and removal.
- A clean interactive snapshot for the downstream #94 attended walkthrough.

Do not run destructive lifecycle scenarios on a retained development workstation.

## Focused test-first checks

```powershell
go test ./internal/winuninstall ./test/integration -run 'WinUninstall|WindowsInstaller' -count=1
pwsh -NoProfile -File build/windows/verify_wxs.ps1
```

Parse every changed PowerShell file before running it. All child console processes in repository automation must remain hidden and noninteractive.

## Build and inspect an MSI

Build the three product executables and the windowless cleanup helper into a temporary stage, run the source verifier, compile `build/windows/goschedule.wxs` with the pinned UI and Util extensions, and then run:

```powershell
pwsh -NoProfile -File test/windows/inspect-installer.ps1 `
  -MsiPath <candidate.msi> `
  -ArtifactClass candidate `
  -ArtifactOrigin <commit-or-working-tree-description>
```

The inspector must prove both shortcut features, both shortcut rows, completion controls and ordered events, the secure preserve-by-default property, commit cleanup action and condition, and close-application coverage.

## Silent lifecycle matrix

On a disposable elevated Windows host, run the hosted installer-contract probe for:

1. default shortcuts and preserve removal;
2. both shortcuts and preserve removal;
3. neither shortcut and explicit wipe removal;
4. desktop-only plus wipe removal;
5. invalid wipe property, repair, reinstall, and upgrade controls.

Seed the machine root, current-profile preference root, and outside-root sentinels before removal. Preserve must retain exact hashes. Wipe must remove declared roots, preserve sentinels and `goschedadmin`, and record `complete` with no residual report on the success path.

## Canonical repository verification

Run in the foreground and watch all eight gates:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

## Attended release proof

After S039 produces a release candidate, #94 runs the interactive clean-desktop matrix. Source inspection, a silent CI runner, or screenshots alone do not substitute for that gate. Until it passes, #97 and #98 remain open.
