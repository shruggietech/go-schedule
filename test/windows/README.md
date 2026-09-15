# Windows installer verification

These tools separate compiled-MSI evidence from native lifecycle evidence. They are maintainer procedures and never count a missing prerequisite as a pass.

## S088 Sandbox testing waiver

On 2026-09-15, the maintainer directed reuse of Windows Sandbox, prohibited restarting the development host, and explicitly authorized release despite incomplete testing (see [#226](https://github.com/shruggietech/go-schedule/issues/226#issuecomment-5685176399)). Use Sandbox only; do not enable host features, provision another virtualization platform, change host security settings, or install the product on the active development host.

The waiver changes release acceptance for unavailable checks, not their factual result. Run supported scenarios after reviewed-main staging and exact candidate verification. Record each unavailable scenario and capability as **not tested, release authorized by maintainer**, preserving failed and timed-out results separately. Historical candidate observations cannot be reassigned to new bytes. Do not manufacture passing fragments, omit required scenarios from the collector, or claim full qualification from partial evidence.

The existing collector and promotion validator still require the full evidence contract for fully qualified publication. A waiver is not a valid passing archive. Any waiver-aware promotion implementation must be explicitly reviewed and preserve source, artifact, checksum, and observation identity checks. Source-owned release metadata changes must merge before staging the final reviewed boundary.

Use `-ArtifactClass local-demo` for a pre-publication exploratory build. This keeps its report distinct from a workflow-staged `candidate` and a release-downloaded `published` artifact. Local-demo inspection proves compiled authoring only and cannot produce a candidate manifest or satisfy the attended gate.

## Inspect an MSI without installing it

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\artifact.md `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The inspector reads Summary Information PID 3 plus canonical product identity, icon, shortcut, PATH, and `Wix4Group` rows. PID 3 must equal `go-schedule: cross-platform task scheduler`; evidence records that value and the artifact SHA-256. It also queries `System.Subject` through the native Shell property handler that Explorer uses. The `GoScheduleAdminGroup` row must name `goschedadmin` with an empty domain. S039 inspection additionally proves the compiled shortcut feature/component relationships, both shortcut identities, the package-owned uninstall and completion dialogs, maintenance-only Windows application-management registration, independent completion controls and Finish events, secure preserve-by-default removal property, invalid-value guard, embedded windowless cleanup action, exact removal condition, and daemon/GUI close rows. This proves compiled authoring and the Explorer property-system value, not installer lifecycle execution.

## Hosted silent installer contract

`Invoke-InstallerContractCI.ps1` is destructive and refuses to run outside an elevated GitHub-hosted disposable Windows runner. The `windows-msi-contract` CI job builds two test-version MSIs from the same reviewed source, inspects the candidate database, and invokes the script with hidden `msiexec` processes:

```powershell
pwsh -NoProfile -File test/windows/Invoke-InstallerContractCI.ps1 `
  -BaselineMsiPath C:\verify\s039-baseline.msi `
  -MsiPath C:\verify\s039-candidate.msi `
  -EvidencePath C:\verify\silent-lifecycle.json `
  -ArtifactOrigin 'CI build from commit <full-commit-id>' `
  -Confirm:$false
```

The probe covers default, both, neither, and desktop-only shortcut states; maintenance transitions; same-authoring major-upgrade feature migration; invalid `GOSCHEDULE_REMOVE_DATA` rejection; repair and upgrade non-wipe controls; preserve and explicit-wipe removal; out-of-scope sentinels; retained `goschedadmin`; a locked-file partial cleanup; cleanup-result evidence; and the absence of GUI or documentation completion actions from every silent MSI log. After fresh install, repair, and upgrade, it also proves one visible product entry with `NoRemove=1`, no `NoModify`, no direct `UninstallString`, and a current MSI-owned `ModifyPath`. This is the supported Settings entry into the full maintenance wizard; direct silent `/x` remains available to administrators.

`GOSCHEDULE_REMOVE_DATA=0` (or an absent property) preserves application data. Only exact `GOSCHEDULE_REMOVE_DATA=1` requests a committed wipe. The helper's retained failure ledger is `%ProgramData%\ShruggieTech\go-schedule-uninstall\b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34\cleanup-result.json`; the matching HKLM summary records state, remaining count, and report path. Complete cleanup removes stale result evidence. MSI success proves software removal, so the lifecycle probe separately verifies the cleanup result.

The hosted runner is Windows Server and has no attended desktop session. Its evidence is explicitly labelled `hosted Windows Server silent installer contract`. It does not prove visible dialog defaults, confirmation and cancel interaction, the Windows 11 Settings wording/navigation, Explorer launches, browser handling, interactive-user integrity, or native window behavior. Those release-candidate observations remain the clean Windows 11 gate in #94, and #97/#98 remain open until that gate passes.

## Run the fresh lifecycle

Use a clean disposable Windows 11 snapshot with no product, service, install directory, PATH entry, or `goschedadmin`. Open elevated PowerShell 7:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario fresh -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\fresh.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>' `
  -Confirm:$false
```

The script installs, repairs, reinstalls, and uninstalls. It records exit codes, verbose-log paths/hashes, diagnostics, group SID/members, service state, product state, install directory, and PATH. Uninstall must preserve group membership.

## Run the v1.1.1 upgrade lifecycle

Revert to a separate clean snapshot. Install the public v1.1.1 MSI and create representative tasks, history, and appearance state before applying the exact v1.4.0 candidate:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario upgrade -MsiPath C:\verify\candidate.msi `
  -PriorMsiPath C:\verify\go-schedule_v1.1.1_windows_amd64.msi `
  -EvidencePath C:\verify\upgrade.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>' `
  -PriorArtifactOrigin 'https://github.com/shruggietech/go-schedule/releases/download/v1.1.1/go-schedule_v1.1.1_windows_amd64.msi' `
  -Confirm:$false
```

Do not reuse the fresh host without reverting it. Preserved group state makes that baseline intentionally unclean.

## Prove ordinary non-elevated access

On another disposable installation, install the candidate as the intended account, then run from normal PowerShell 7 without elevating the CLI:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario access-probe -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\access.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

The probe records whether the current token contains the `goschedadmin` SID but does not require it. It requires direct group membership, records the expected restricted descriptor, and requires `gosched health` to succeed. This covers a newly enrolled or UAC-filtered standard token through its stable user SID. An elevated probe is rejected because Administrators have independent daemon access. Uninstall afterward from an elevated session and revert the host.

## Prove the installed core path

Use the same non-elevated installed state with a new evidence path:

```powershell
pwsh .\Invoke-InstallerLifecycle.ps1 `
  -Scenario installed-core-probe -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\installed-core.md -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

This scenario first proves ordinary IPC access, then creates tasks through the installed CLI. The real LocalSystem daemon executes an absolute inbox Windows PowerShell command manually and on a five-second schedule. Evidence contains both run records, exit code 0, output, and a two-line marker file. It also records a controlled exit-code 7 run and a missing-executable process-start failure. Probe tasks are removed afterward; marker evidence is retained beside the Markdown report. Environment values are never recorded.

## CI service-boundary probe

`Invoke-ServiceCoreCI.ps1` is reserved for a clean, elevated, disposable Windows runner. CI builds the daemon and CLI, creates `goschedadmin`, registers `goschedd` as LocalSystem, and runs the same success and failure controls without building or installing an MSI. Its JSON artifact is automated service-boundary evidence, not ordinary-token IPC proof or release-equivalent installer evidence. The script removes its tasks, service, and group in a `finally` block.

## Attended release-candidate gate

### Prepared disposable sessions (S086)

`scripts/windows-qualification-session` prepares fresh and upgrade Windows Sandbox launch packages without installing software on the development host. Supply an absolute JSON manifest and an absolute new output directory whose parent already exists. Inputs must be local regular files with explicit source identity, byte length, and lowercase SHA-256. Candidate source identifies the exact hosted staging run; baseline source identifies the public v1.1.1 MSI; PowerShell and WebView2 use upstream portable/offline distributions. The tool rejects occupied destinations, links, overlap, missing inputs, and changed bytes. Do not point it at the stale v1.4.0 draft and claim the UI fixes qualified.

```json
{
  "schema_version": 1,
  "repository": "shruggietech/go-schedule",
  "tag": "v1.4.0",
  "commit": "REPLACE_WITH_EXACT_40_CHARACTER_STAGED_COMMIT",
  "run_id": 123456789,
  "run_attempt": 1,
  "inputs": [
    { "role": "candidate", "path": "C:\\candidate\\go-schedule_v1.4.0_windows_amd64.msi", "bytes": 1, "sha256": "REPLACE_WITH_EXACT_SHA256", "source": "https://github.com/shruggietech/go-schedule/actions/runs/123456789" },
    { "role": "baseline", "path": "C:\\baseline\\go-schedule_v1.1.1_windows_amd64.msi", "bytes": 1, "sha256": "REPLACE_WITH_EXACT_SHA256", "source": "https://github.com/shruggietech/go-schedule/releases/download/v1.1.1/go-schedule_v1.1.1_windows_amd64.msi" },
    { "role": "powershell", "path": "C:\\tools\\powershell.zip", "bytes": 1, "sha256": "REPLACE_WITH_EXACT_SHA256", "source": "https://github.com/PowerShell/PowerShell/releases/download/v7.5.0/PowerShell-7.5.0-win-x64.zip" },
    { "role": "webview2", "path": "C:\\tools\\webview2.exe", "bytes": 1, "sha256": "REPLACE_WITH_EXACT_SHA256", "source": "https://developer.microsoft.com/en-us/microsoft-edge/webview2/" }
  ]
}
```

The sample is intentionally invalid until actual sizes, hashes, commit, and upstream runtime version/source are recorded. A recorded URL is provenance supplied by the maintainer, not a signature or online authenticity check. Authenticate downloaded upstream packages before preparation. Generate with `go run ./scripts/windows-qualification-session --manifest C:\candidate\session.json --output C:\qualification\new-session`. The new package contains `fresh.wsb`, `upgrade.wsb`, read-only inputs, and independent initially empty export folders. Launching a `.wsb` is an intentional attended installation activity, not part of nondestructive preparation. Networking is disabled; this configuration does not qualify browser or remote-network scenarios.

The guest-only bootstrap copies verified bytes to guest-local disk, prepares offline WebView2, verifies PowerShell 7, initializes the bundled collector's unavailable templates, and opens attended MSI UI with `/L*vx!` logs. It exports ten-second progress and phase records directly to the dedicated host folder. A failure or diagnostic deadline stops further operations without killing the Windows Installer service. MSI code 3010 is recorded as completed-reboot-required, not installation failure, but stops subsequent work until the required reboot can occur in an appropriate environment. Preserve exports and reset or replace the guest before retrying with a new package. No log or phase record is an attended pass.

The collector is available at `C:\qualification\Invoke-ReleaseCandidateAttended.ps1`, its PowerShell runtime at `C:\qualification\pwsh\pwsh.exe`, and its initialized workspace at `C:\qualification\attended-workspace`. After the operator completes or explicitly leaves unavailable the walkthrough and clicks OK on the final instruction dialog, the workspace is copied create-only to the dedicated host exports. Do not close Sandbox before that export. Finalization uses the original repository collector with its Go context on the host and the exact packaged MSI; do not run the copied collector's repository-relative finalization from its guest location. Merge only genuine compatible observations from the separately required environments, preserving identities and attachment hashes. The final gate remains authoritative.

Upgrade uses a separate reset session, installs public v1.1.1, and pauses for genuine representative tasks, run history, and appearance preparation before candidate installation. Baseline and candidate installation are serialized. Do not reuse fresh-install residue as upgrade evidence. Sandbox-account/virtual-machine guards prevent accidental development-host installation; they are not a security boundary against deliberately modified packages or impersonation.

The remaining native walkthrough includes #229 (System/Light/Dark icon visibility), #230 (compact controls, hover/focus/pressed/disabled states, semantic actions, card spacing, persistent navigation/Exit), #231 (task modal sizing, examples, selectors, deletion spacing), #232 (Notifications first, collapsed advanced configuration and help), and #233 (temporary feedback, dismissible errors, aligned Agent Access/Connections, compact monospaced paths, no Copy flicker). These regressions supplement, not replace, the complete required matrix below. Normal-user token, multi-profile, Explorer/browser, and high/mixed-DPI observations require suitable separate clean Windows 11 environments. Do not ask the operator to attest an unavailable environment or manufacture passing fragments.

Historical S081 candidate commit `57555ffa413df641cb21784847ad598c113bf199` predates S083-S085. S087 preserved its bytes and diagnostics, then refreshed the draft to reviewed S086 commit `951d864d9683ec3bdcb1e37535ecddd67738bf13`. That replacement still predates the S087 timing-selector repair. Preserve both candidates' provenance; neither supplies passing evidence for a future replacement. S088 source reconciliation must merge before authorized reviewed-main restaging. Its Sandbox testing waiver permits explicitly untested criteria but does not make partial evidence satisfy the unchanged full-qualification collector or authorize implicit issue closure.

S040 adds `Invoke-ReleaseCandidateAttended.ps1` as the resumable collector for the clean Windows 11 work that cannot run credibly on a hosted server. The tag workflow first stages all platform assets in a draft GitHub release. Use the Windows MSI and `windows-candidate-manifest.json` from that exact draft. Never rebuild, rename, or substitute the MSI after evidence collection starts.

Initialize a new workspace from a normal, non-elevated PowerShell 7 session:

```powershell
$commit = '0123456789abcdef0123456789abcdef01234567'
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action Initialize `
  -MsiPath C:\verify\go-schedule_v1.4.0_windows_amd64.msi `
  -WorkspacePath C:\verify\v1.4.0-attended `
  -Tag v1.4.0 -Commit $commit -RunId 123456789 -RunAttempt 1
```

Initialization reads ProductVersion and ProductCode from the compiled MSI and records repository, tag, commit, staging run and attempt, filename, byte size, and SHA-256. It creates all 47 required observations as explicit `unavailable` placeholders plus fail-closed setup, removal, and desktop templates under `fragments`. Those templates enumerate required process/session, option, target, inventory, fingerprint, unaffected-control, security-state, reinstall, and current Wails workspace fields. The collector refuses an existing workspace so an interrupted run cannot silently erase evidence.

Each operator-reviewed fragment contains one environment and one observation. Use genuine registered local profiles and identify accounts by role plus SID, not by personal name. Record the token integrity RID (8192 for medium, 12288 for high, or 16384 for system). `RecordObservation` replaces one unused placeholder and refuses later overwrite. Every status is explicit: `pass`, `fail`, `unavailable`, `skipped`, `timed-out`, or `partial`. Only `pass` can satisfy promotion.

### Native window capture

Launch the exact installed Wails GUI, capture its process ID and a native raster screenshot, then measure the exact visible HWND. The collector derives logical content dimensions and display scale from the Win32 client rectangle and effective DPI, so no toolkit-specific instrumentation file is required:

```powershell
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action CaptureWindow `
  -WorkspacePath C:\verify\v1.4.0-attended `
  -ProcessId 1234 `
  -ObservationId window.clean-standard `
  -EnvironmentPath C:\verify\standard-environment.json `
  -ScreenshotPath C:\verify\v1.4.0-attended\attachments\windows\screen.png
```

Capture requires a screenshot and records the executable path and hash, process session, user SID and token-integrity RID, single visible top-level HWND, outer and client rectangles, monitor and work-area rectangles, effective DPI, measured restored/maximized/minimized/fullscreen state, and generic desktop content metrics. It rejects a supplied environment whose account SID or integrity RID does not match the live process token. Review the generated fragment before changing its status to `pass`. In particular, record monitor identity and confirm visible margins, title bar, resize borders, and taskbar. The same exact MSI needs separate clean standard-DPI, clean high-DPI or mixed-DPI, retained v1.1.1 profile, state-transition, and subsequent-launch observations. Use `retained-v1.1.1` as the profile state for the retained observation.

### Required attended matrix

- Prove normal-user service access and GUI task listing, LocalSystem service identity, unrelated-user pipe denial, new-process PATH resolution, and PATH absence after uninstall.
- Exercise daemon unavailable, access denied, timeout, stream disconnect, repeated refresh or reconnect, manual Retry, and recovery. Retain at least 120 seconds of timestamped samples for each repetition-sensitive condition.
- Record one in-frame incident, zero modal overlays, and zero additional top-level error windows. HWND enumeration cannot prove WebView content state, so screenshots and attended visible-surface counts are both required.
- Use distinct production run identities for manual success, scheduled success, exit-code 7, and process-start failure, while binding every result to the same candidate identity. Retain one `attachments/tasks/task-runs.json` document with schema version `1`, kind `task-run-evidence-v1`, and exactly one record per task observation. Each record preserves the task definition, captured output, completion marker, history result, production-run flag, expected/actual result, and diagnostic category. Reference that attachment from all four task observations; the gate independently hashes the retained values and compares them with the observation metrics.
- Exercise shortcut defaults and all selections, four independent completion combinations, medium-integrity Finish launch, cancel, maintenance, upgrade, invalid-input rejection, transactional rollback, preserve, wipe, locked partial cleanup, at least two genuine profiles, and reinstall after both removal modes. The upgrade observation must record the public v1.1.1 MSI version, filename, download URL, and SHA-256 (`f7ac8f56f28330b016eb6e505e424b19e9bfbe435591cfbc54a723c91ac8e567`), then independently prove tasks, run history, appearance intent, daemon identity, service operation, and local access survive while notifications, localhost MCP, and remote HTTPS remain disabled.

### Current Wails desktop release-qualification matrix

Use the exact installed release candidate. Every row is a required `pass` observation from the intended user at medium integrity with the installed LocalSystem service and at least one native raster screenshot whose bytes are validated independently of its declared media type or extension. The collector generates the metric fields and expected set values, so do not collapse or rename them.

| Scenario | Issues | Native outcome |
| --- | --- | --- |
| `desktop.appearance-standard` | #226 | At 96 DPI, exercise System, Light, and Dark. System is the clean-profile and restored default; bundled brand typography and body text remain sharp, centered, and unclipped after resize, minimize/restore, and reopen. |
| `desktop.appearance-scaled` | #226 | Repeat the complete current-theme and typography observation above 96 DPI, recording the environment's exact effective DPI. |
| `desktop.interaction-states` | #226 | In both palettes, exercise navigation, selector, ordinary, primary, danger, dialog, and table-row controls at rest, hover, focus, pressed, selected, and disabled. Text contrast is at least 4.5:1, non-text contrast is at least 3:1, and state meaning is not color-only. |
| `desktop.interaction-states-scaled` | #226 | Repeat the complete interaction-state observation above 96 DPI. |
| `desktop.navigation-options` | #226 | At 1280x800 and 800x600, prove the Tasks, Automation Sources, Schedule, Activity, Notifications, Agent Access, Connections, and Settings order, persistent target context, separate Exit, compact storage rows, exact Copy behavior, muted unavailable rows, current-option omission, and no horizontal scrollbar. |
| `desktop.navigation-options-scaled` | #226 | Repeat the complete current navigation and Settings observation above 96 DPI. |
| `desktop.scroll-input` | #226 | Exercise Tasks, Automation Sources, Schedule, Activity, Notifications, Agent Access, Connections, Settings, and the task editor. Prove responsive browser-native wheel input, no nested multiplication, preserved keyboard and focus behavior, and either physical touchpad precision or a specific unavailability reason. |
| `desktop.tasks-table` | #226 | With at least 100 rows in both palettes and sizes, prove Task, Group, State, and Schedule headers, discoverable state reasons, stable whole-row states and refresh identity, complete-value disclosure, safe removal, working toolbar and double-click actions, and no horizontal scrollbar. |
| `desktop.tasks-table-scaled` | #226 | Repeat the complete current Tasks table observation above 96 DPI. |
| `desktop.schedule-activity-tables` | #226 | With at least 100 rows per view in both palettes and sizes, prove current Schedule and Activity headers, upcoming through unavailable schedule states, run, daemon-log, and alert record types, info through error severity, stable row identity, accurate detail, range, calendar, filter, clear, and acknowledge behavior, and no horizontal scrollbar. |
| `desktop.schedule-activity-tables-scaled` | #226 | Repeat the complete current Schedule and Activity observation above 96 DPI. |

Native evidence is required because headless layout, contrast, mapping, and scroll tests cannot prove Windows text rasterization, physical input, display scaling, or interaction-state readability. Keep release issue #226 open until its acceptance criteria and formal exact-candidate evidence have been reviewed.

Finalize after all fragments are reviewed and recorded:

```powershell
pwsh -NoProfile -File .\Invoke-ReleaseCandidateAttended.ps1 `
  -Action Finalize `
  -MsiPath C:\verify\go-schedule_v1.4.0_windows_amd64.msi `
  -WorkspacePath C:\verify\v1.4.0-attended
```

Finalize hashes every referenced attachment and invokes the shared Go gate. It produces the canonical ZIP only if all identity, environment, scenario, timing, measurement, and attachment rules pass. Upload that archive to the same draft release. The manual Promote Release workflow revalidates the draft, staging run, last-observed remote tag commit, exact allowlisted asset set, manifest, archive, and exact MSI; creates the final all-asset checksum file; and only then makes the release public.

### Historical v1.0.0 issue disposition packet

The following retained procedure documents how the public v1.0.0 evidence was reconciled against its historical issues after independent verification. It is not part of the v1.4.0 release path:

```powershell
go run ./scripts/windows-release-gate render-dispositions `
  --bundle C:\verify\go-schedule_v1.0.0_windows-attended-evidence.zip `
  --candidate-manifest C:\verify\windows-candidate-manifest.json `
  --artifact C:\verify\go-schedule_v1.0.0_windows_amd64.msi `
  --repository shruggietech/go-schedule `
  --tag v1.0.0 `
  --commit 0123456789abcdef0123456789abcdef01234567 `
  --output-dir C:\verify\v1.0.0-dispositions
```

The output directory must not already exist. A successful run creates `packet.json` plus `issue-096.md`, `issue-098.md`, `issue-101.md`, `issue-104.md`, `issue-105.md`, `issue-106.md`, `issue-109.md`, `issue-111.md`, `issue-112.md`, and `issue-113.md` in one atomic directory commit. Review every record and post only the matching record to its issue. The packet does not update or close GitHub issues and does not authorize promotion; actual issue state and all individual acceptance criteria remain authoritative.

The checked-in `test/fixtures/windows-release-gate/passing` data is plain text and explicitly non-native. It proves validator behavior only. It cannot close #94 or #98 and cannot authorize promotion.

## Evidence boundaries

- Source contracts prove tracked authoring.
- MSI inspection proves compiled table and Summary Information contents.
- The Shell `System.Subject` assertion proves the native value consumed by Explorer tooltip and Properties presentation.
- Fresh and upgrade runs prove native execution on their named host.
- The access probe proves ordinary direct-member client access, including a standard token that does not yet contain the local alias SID.
- The installed-core probe proves manual and scheduled production execution in the LocalSystem service context plus diagnostic failure controls.
- The CI service-core probe continuously proves the binary-level LocalSystem boundary; it does not replace the candidate-MSI walkthrough.
- The CI MSI-contract probe proves compiled database and silent native lifecycle behavior on its disposable Windows Server runner. It does not impersonate #94's attended clean Windows 11 desktop evidence.
- `unavailable` cannot close an issue that requires runtime evidence.
