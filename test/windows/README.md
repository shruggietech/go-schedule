# Windows installer verification

These tools separate compiled-MSI evidence from native lifecycle evidence. They
are maintainer procedures and never count a missing prerequisite as a pass.

## Inspect an MSI without installing it

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\artifact.md `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The inspector reads icon, shortcut, PATH, and `Wix4Group` rows. The
`GoScheduleAdminGroup` row must name `goschedadmin` with an empty domain. This
proves compiled authoring, not runtime execution.

## Run the fresh lifecycle

Use a clean disposable Windows 11 snapshot with no product, service, install
directory, PATH entry, or `goschedadmin`. Open elevated PowerShell 7:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario fresh -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\fresh.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>' `
  -Confirm:$false
```

The script installs, repairs, reinstalls, and uninstalls. It records exit codes,
verbose-log paths/hashes, diagnostics, group SID/members, service state, product
state, install directory, and PATH. Uninstall must preserve group membership.

## Run the v0.9.0 upgrade lifecycle

Revert to a separate clean snapshot. The script preprovisions `goschedadmin`
and the current account so v0.9.0 can establish an upgrade baseline:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario upgrade -MsiPath C:\verify\candidate.msi `
  -PriorMsiPath C:\verify\v0.9.0.msi `
  -EvidencePath C:\verify\upgrade.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>' `
  -PriorArtifactOrigin 'https://github.com/shruggietech/go-schedule/releases/download/v0.9.0/go-schedule_v0.9.0_windows_amd64.msi' `
  -Confirm:$false
```

Do not reuse the fresh host without reverting it. Preserved group state makes
that baseline intentionally unclean.

## Prove refreshed non-elevated access

On another disposable installation, install the candidate, sign out and back in
as the enrolled account, then run from normal PowerShell 7:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario access-probe -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\access.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The probe requires the current token to contain the `goschedadmin` SID and
requires `gosched health` to succeed. An elevated probe is rejected because
Administrators have independent daemon access. Uninstall afterward from an
elevated session and revert the host.

## Evidence boundaries

- Source contracts prove tracked authoring.
- MSI inspection proves compiled table contents.
- Fresh and upgrade runs prove native execution on their named host.
- The access probe proves post-refresh membership-based client access.
- `unavailable` cannot close an issue that requires runtime evidence.
