# Quickstart: Verify the Windows MSI Local-Group Fix

## Portable contracts

```powershell
go test ./test/integration -run `
  'TestWindowsInstaller(AdminGroupContract|AdminGroupRejectsBrokenLifecycle|EvidenceToolingContract)$' `
  -count=1
pwsh -NoLogo -NoProfile -NonInteractive -File build/windows/verify_wxs.ps1
```

## Inspect a candidate

```powershell
pwsh -NoLogo -NoProfile -NonInteractive `
  -File test/windows/inspect-installer.ps1 `
  -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\candidate-msi.md `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

Evidence must report `GoScheduleAdminGroup`, `goschedadmin`, and an empty domain.

## Fresh lifecycle

On a clean elevated disposable Windows 11 snapshot:

```powershell
pwsh -NoLogo -NoProfile -NonInteractive `
  -File .\test\windows\Invoke-InstallerLifecycle.ps1 `
  -Scenario fresh -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\fresh.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

## v0.9.0 upgrade lifecycle

On a separate reset snapshot:

```powershell
pwsh -NoLogo -NoProfile -NonInteractive `
  -File .\test\windows\Invoke-InstallerLifecycle.ps1 `
  -Scenario upgrade -MsiPath C:\verify\candidate.msi `
  -PriorMsiPath C:\verify\v0.9.0.msi `
  -EvidencePath C:\verify\upgrade.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>' `
  -PriorArtifactOrigin 'https://github.com/shruggietech/go-schedule/releases/download/v0.9.0/go-schedule_v0.9.0_windows_amd64.msi'
```

Both evidence files must be proven. Then run `sh scripts/verify.sh all` in the foreground through all eight gates.
