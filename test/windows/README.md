# Windows installer verification

These tools separate compiled-MSI evidence from native lifecycle evidence. They
are maintainer procedures and never count a missing prerequisite as a pass.

Use `-ArtifactClass local-demo` for a pre-publication exploratory build. This
keeps its report distinct from a workflow-staged `candidate` and a release-
downloaded `published` artifact. Local-demo inspection proves compiled authoring
only and cannot produce a candidate manifest or satisfy the attended gate.

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
the package-owned uninstall and completion dialogs, maintenance-only Windows
application-management registration, independent completion
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
After fresh install, repair, and upgrade, it also proves one visible product
entry with `NoRemove=1`, no `NoModify`, no direct `UninstallString`, and a
current MSI-owned `ModifyPath`. This is the supported Settings entry into the full
maintenance wizard; direct silent `/x` remains available to administrators.

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
interaction, the Windows 11 Settings wording/navigation, Explorer launches,
browser handling, interactive-user integrity, or native window behavior. Those
release-candidate observations remain the clean Windows 11 gate in #94, and
#97/#98 remain open until that gate passes.

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

## Attended release-candidate gate

S040 adds `Invoke-ReleaseCandidateAttended.ps1` as the resumable collector for
the clean Windows 11 work that cannot run credibly on a hosted server. The tag
workflow first stages all platform assets in a draft GitHub release. Use the
Windows MSI and `windows-candidate-manifest.json` from that exact draft. Never
rebuild, rename, or substitute the MSI after evidence collection starts.

Initialize a new workspace from a normal, non-elevated PowerShell 7 session:

```powershell
$commit = '0123456789abcdef0123456789abcdef01234567'
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action Initialize `
  -MsiPath C:\verify\go-schedule_v1.0.0_windows_amd64.msi `
  -WorkspacePath C:\verify\v1.0.0-attended `
  -Tag v1.0.0 -Commit $commit -RunId 123456789 -RunAttempt 1
```

Initialization reads ProductVersion and ProductCode from the compiled MSI and
records repository, tag, commit, staging run and attempt, filename, byte size,
and SHA-256. It creates all 36 required observations as explicit
`unavailable` placeholders plus fail-closed setup/removal templates under
`fragments`. Those templates enumerate required process/session, option,
target, inventory, fingerprint, unaffected-control, security-state, and
reinstall fields. The collector refuses an existing workspace so an
interrupted run cannot silently erase evidence.

Each operator-reviewed fragment contains one environment and one observation.
Use genuine registered local profiles and identify accounts by role plus SID,
not by personal name. Record the token integrity RID (8192 for medium, 12288
for high, or 16384 for system). `RecordObservation` replaces one unused
placeholder and refuses later overwrite. Every status is explicit: `pass`, `fail`, `unavailable`,
`skipped`, `timed-out`, or `partial`. Only `pass` can satisfy promotion.

### Native window capture

Launch the exact installed GUI with
`GOSCHEDULE_WINDOW_EVIDENCE_PATH` set to a new JSON file beneath the evidence
workspace's `attachments` directory. The opt-in file records literal Fyne
canvas width, height, and scale. Then capture the exact process and visible
HWND:

```powershell
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action CaptureWindow `
  -WorkspacePath C:\verify\v1.0.0-attended `
  -ProcessId 1234 `
  -ObservationId window.clean-standard `
  -EnvironmentPath C:\verify\standard-environment.json `
  -FyneEvidencePath C:\verify\v1.0.0-attended\attachments\windows\fyne.json `
  -ScreenshotPath C:\verify\v1.0.0-attended\attachments\windows\screen.png
```

Capture requires a screenshot and records the executable path and hash,
process session, user SID and token-integrity RID, single visible
top-level HWND, outer and client rectangles, monitor and work-area rectangles,
effective DPI, measured restored/maximized/minimized/fullscreen state, and Fyne metrics. It
rejects a supplied environment whose account SID or integrity RID does not
match the live process token. Review the generated
fragment before changing its status to `pass`. In particular, record monitor
identity and confirm visible margins, title bar, resize borders, and taskbar.
The same exact MSI needs separate clean standard-DPI, clean high-DPI or
mixed-DPI, retained v0.9.1 profile, state-transition, and subsequent-launch
observations.

### Required attended matrix

- Prove normal-user service access and GUI task listing, LocalSystem service
  identity, unrelated-user pipe denial, new-process PATH resolution, and PATH
  absence after uninstall.
- Exercise daemon unavailable, access denied, timeout, stream disconnect,
  repeated refresh or reconnect, manual Retry, and recovery. Retain at least
  120 seconds of timestamped samples for each repetition-sensitive condition.
- Record one in-frame incident, zero modal overlays, and zero additional
  top-level error windows. HWND enumeration cannot see Fyne canvas overlays,
  so screenshots and attended visible-surface counts are both required.
- Use distinct production run identities for manual success, scheduled
  success, exit-code 7, and process-start failure, while binding every result
  to the same candidate identity. Retain one
  `attachments/tasks/task-runs.json` document with schema version `1`, kind
  `task-run-evidence-v1`, and exactly one record per task observation. Each
  record preserves the task definition, captured output, completion marker,
  history result, production-run flag, expected/actual result, and diagnostic
  category. Reference that attachment from all four task observations; the
  gate independently hashes the retained values and compares them with the
  observation metrics.
- Exercise shortcut defaults and all selections, four independent completion
  combinations, medium-integrity Finish launch, cancel, maintenance, upgrade,
  invalid-input rejection, transactional rollback, preserve, wipe, locked
  partial cleanup, at least two genuine profiles, and reinstall after both
  removal modes.

### S047 desktop release-qualification matrix

Use the exact installed release candidate. Every row is a required `pass`
observation from the intended user at medium integrity with the installed
LocalSystem service and at least one native screenshot. The collector generates
the metric fields and expected set values, so do not collapse or rename them.

| Scenario | Issues | Native outcome |
| --- | --- | --- |
| `desktop.appearance-standard` | #101, #106 | At 96 DPI, exercise Dark and Light plus System, Geist, Inter, Ubuntu, and Monospace. System is the clean-profile and restored default; text remains sharp, centered, and unclipped after resize, minimize/restore, and reopen. |
| `desktop.appearance-scaled` | #101, #106 | Repeat the complete appearance observation above 96 DPI, recording the environment's exact effective DPI. |
| `desktop.interaction-states` | #109 | In both palettes, exercise navigation, selector, ordinary, primary, danger, dialog, and table-row controls at rest, hover, focus, pressed, selected, and disabled. Text contrast is at least 4.5:1, non-text contrast is at least 3:1, and state meaning is not color-only. |
| `desktop.navigation-options` | #104, #105, #106 | At 1280x800 and 800x600, prove destination order, spacing, full-height boundary, separate bottom-right Exit, compact storage rows, exact Copy behavior, muted unavailable rows, current-option omission, and no horizontal scrollbar. |
| `desktop.scroll-input` | #111 | At 1x, 2x, and 4x, exercise Options, Info, command, schedule, and Help surfaces. Prove immediate persisted wheel scaling, no nested multiplication, preserved keyboard input, and either physical touchpad precision or a specific unavailability reason. |
| `desktop.tasks-table` | #112 | With at least 100 rows in both palettes and sizes, prove fixed labeled headers, distinct Enabled/Lifecycle fields, stable whole-row states and refresh identity, complete-value disclosure, safe removal, working toolbar/double-click actions, and no horizontal scrollbar. |
| `desktop.schedule-activity-tables` | #113 | With at least 100 rows per view in both palettes and sizes, prove fixed labeled headers, exact schedule-state and uppercase severity sets, matching text/glyph semantics, stable row identity, accurate detail/range/calendar/filter/clear/acknowledge behavior, and no horizontal scrollbar. |

Native evidence is required because headless layout, contrast, mapping, and
scroll tests cannot prove Windows text rasterization, physical input, display
scaling, or interaction-state readability. Keep each issue open until its own
acceptance criteria and formal exact-candidate evidence have been reviewed.

Finalize after all fragments are reviewed and recorded:

```powershell
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action Finalize `
  -MsiPath C:\verify\go-schedule_v1.0.0_windows_amd64.msi `
  -WorkspacePath C:\verify\v1.0.0-attended
```

Finalize hashes every referenced attachment and invokes the shared Go gate.
It produces the canonical ZIP only if all identity, environment, scenario,
timing, measurement, and attachment rules pass. Upload that archive to the
same draft release. The manual Promote Release workflow revalidates the draft,
staging run, last-observed remote tag commit, exact allowlisted asset set,
manifest, archive, and exact MSI; creates the final
all-asset checksum file; and only then makes the release public.

The checked-in `test/fixtures/windows-release-gate/passing` data is plain text
and explicitly non-native. It proves validator behavior only. It cannot close
#94 or #98 and cannot authorize promotion.

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
