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

The inspector reads Summary Information PID 3 plus canonical product identity,
icon, shortcut, PATH, and `Wix4Group` rows. PID 3 must equal
`go-schedule: cross-platform task scheduler`; evidence records that value and
the artifact SHA-256. It also queries `System.Subject` through the native Shell
property handler that Explorer uses. The `GoScheduleAdminGroup` row must name
`goschedadmin` with an empty domain. S039 inspection additionally proves the
compiled shortcut feature/component relationships, both shortcut identities,
the package-owned uninstall and completion dialogs, independent completion
controls and Finish events, secure preserve-by-default removal property,
invalid-value guard, embedded windowless cleanup action, exact removal
condition, and daemon/GUI close rows. This proves compiled authoring and the
Explorer property-system value, not installer lifecycle execution.

## Hosted silent installer contract

`Invoke-InstallerContractCI.ps1` is destructive and refuses to run outside an
elevated GitHub-hosted disposable Windows runner. The `windows-msi-contract` CI
job builds two test-version MSIs from the same reviewed source, inspects the
candidate database, and invokes the script with hidden `msiexec` processes:

```powershell
pwsh -NoProfile -File test/windows/Invoke-InstallerContractCI.ps1 `
  -BaselineMsiPath C:\verify\s039-baseline.msi `
  -MsiPath C:\verify\s039-candidate.msi `
  -EvidencePath C:\verify\silent-lifecycle.json `
  -ArtifactOrigin 'CI build from commit <full-commit-id>' `
  -Confirm:$false
```

The probe covers default, both, neither, and desktop-only shortcut states;
maintenance transitions; same-authoring major-upgrade feature migration;
invalid `GOSCHEDULE_REMOVE_DATA` rejection; repair and upgrade non-wipe
controls; preserve and explicit-wipe removal; out-of-scope sentinels; retained
`goschedadmin`; a locked-file partial cleanup; cleanup-result evidence; and the
absence of GUI or documentation completion actions from every silent MSI log.

`GOSCHEDULE_REMOVE_DATA=0` (or an absent property) preserves application data.
Only exact `GOSCHEDULE_REMOVE_DATA=1` requests a committed wipe. The helper's
retained failure ledger is
`%ProgramData%\ShruggieTech\go-schedule-uninstall\b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34\cleanup-result.json`;
the matching HKLM summary records state, remaining count, and report path.
Complete cleanup removes stale result evidence. MSI success proves software
removal, so the lifecycle probe separately verifies the cleanup result.

The hosted runner is Windows Server and has no attended desktop session. Its
evidence is explicitly labelled `hosted Windows Server silent installer
contract`. It does not prove visible dialog defaults, confirmation and cancel
interaction, Explorer launches, browser handling, interactive-user integrity,
or native window behavior. Those release-candidate observations remain the
clean Windows 11 gate in #94, and #97/#98 remain open until that gate passes.

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

## Prove ordinary non-elevated access

On another disposable installation, install the candidate as the intended
account, then run from normal PowerShell 7 without elevating the CLI:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario access-probe -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\access.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The probe records whether the current token contains the `goschedadmin` SID but
does not require it. It requires direct group membership, records the expected
restricted descriptor, and requires `gosched health` to succeed. This covers a
newly enrolled or UAC-filtered standard token through its stable user SID. An
elevated probe is rejected because Administrators have independent daemon
access. Uninstall afterward from an elevated session and revert the host.

## Prove the installed core path

Use the same non-elevated installed state with a new evidence path:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario installed-core-probe -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\installed-core.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

This scenario first proves ordinary IPC access, then creates tasks through the
installed CLI. The real LocalSystem daemon executes an absolute
inbox Windows PowerShell command manually and on a five-second schedule.
Evidence contains both run records, exit code 0, output, and a two-line marker
file. It also records a controlled exit-code 7 run and a missing-executable
process-start failure. Probe tasks are removed afterward; marker evidence is
retained beside the Markdown report. Environment values are never recorded.

## CI service-boundary probe

`Invoke-ServiceCoreCI.ps1` is reserved for a clean, elevated, disposable
Windows runner. CI builds the daemon and CLI, creates `goschedadmin`, registers
`goschedd` as LocalSystem, and runs the same success and failure controls
without building or installing an MSI. Its JSON artifact is automated
service-boundary evidence, not ordinary-token IPC proof or release-equivalent
installer evidence. The script removes its tasks, service, and group in a
`finally` block.

## Evidence boundaries

- Source contracts prove tracked authoring.
- MSI inspection proves compiled table and Summary Information contents.
- The Shell `System.Subject` assertion proves the native value consumed by Explorer tooltip and Properties presentation.
- Fresh and upgrade runs prove native execution on their named host.
- The access probe proves ordinary direct-member client access, including a
  standard token that does not yet contain the local alias SID.
- The installed-core probe proves manual and scheduled production execution in
  the LocalSystem service context plus diagnostic failure controls.
- The CI service-core probe continuously proves the binary-level LocalSystem
  boundary; it does not replace the candidate-MSI walkthrough.
- The CI MSI-contract probe proves compiled database and silent native
  lifecycle behavior on its disposable Windows Server runner. It does not
  impersonate #94's attended clean Windows 11 desktop evidence.
- `unavailable` cannot close an issue that requires runtime evidence.
