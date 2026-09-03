# Quickstart: Validate Guided Windows Uninstall Entry

## Safety boundary

Use only a disposable Windows test machine or hosted CI runner for install/remove scenarios. Do not launch an attended installer through this repository's headless development session.

## 1. Source and cross-platform contract

```powershell
& ./build/windows/verify_wxs.ps1
go test ./test/integration -run WindowsInstallerContract -count=1
```

Expected: source and mutation contracts require `ARPNOREMOVE=1`, the owned `/I[ProductCode]` registration, the package-owned maintenance page with guided Remove enabled, and every existing removal guard.

## 2. Build and inspect the MSI

Use the pinned Windows CI build path from `.github/workflows/ci.yml`, then run:

```powershell
& ./test/windows/inspect-installer.ps1 `
  -MsiPath <candidate.msi> `
  -ArtifactClass candidate `
  -ArtifactOrigin <immutable-origin>
```

Expected: compiled evidence records direct Remove disabled and maintenance enabled without weakening preserve/wipe sequencing.

## 3. Run hosted lifecycle evidence

Run `test/windows/Invoke-InstallerContractCI.ps1` through the Windows installer CI job.

Expected after fresh install, repair, and upgrade:

- one visible go-schedule application entry;
- `NoRemove=1`;
- `NoModify` absent;
- MSI-owned `ModifyPath` present for the current ProductCode;
- `UninstallString` absent;
- every established silent preserve/wipe and safety scenario passes.

## 4. Perform exact-candidate attended observation

On the clean Windows 11 snapshot required by #94:

1. Install the exact candidate.
2. Open Settings > Apps > Installed apps > go-schedule.
3. Select the available maintenance/Modify action.
4. Observe the full go-schedule maintenance wizard.
5. Select Remove.
6. Observe the removal inventory and preserve-or-wipe choice with preserve selected.
7. Cancel and prove software and data inventories are unchanged.
8. Repeat preserve and confirmed-wipe journeys according to the #94 runbook.

This step is mandatory before #98 closes. Static, compiled, or hosted evidence cannot substitute for it.

## 5. Run canonical verification

```powershell
C:\Program Files\Git\bin\bash.exe scripts/verify.sh all
```

Expected: all eight repository gates pass.
