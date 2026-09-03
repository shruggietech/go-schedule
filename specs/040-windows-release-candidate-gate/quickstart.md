# Quickstart: Windows Release Candidate Gate

This runbook is for the first real Windows candidate after S040 merges. It intentionally does not make the S040 development pull request itself a release candidate.

## 1. Stage the candidate

1. Prepare the reviewed release commit, release notes, and static README version lines.
2. Create and push the authorized `vX.Y.Z` tag.
3. Wait for the Release workflow to finish. Confirm the GitHub release is still a draft.
4. Download the Windows MSI and its candidate manifest from that draft. Do not rebuild or rename the MSI.
5. Record the workflow run ID and attempt, draft release URL, tag target commit, MSI ProductCode, byte size, and SHA-256.
6. Independently validate the manifest before attended work:

```powershell
go run ./scripts/windows-release-gate verify-candidate `
    --candidate-manifest C:\Evidence\windows-candidate-manifest.json `
    --artifact C:\Evidence\go-schedule_vX.Y.Z_windows_amd64.msi `
    --repository shruggietech/go-schedule `
    --tag vX.Y.Z `
    --commit 0123456789abcdef0123456789abcdef01234567
```

Tag creation and release staging require separate operator authorization. They are outside the S040 implementation pull request.

## 2. Prepare clean Windows 11 evidence environments

1. Restore a documented clean Windows 11 client snapshot.
2. Create the intended standard user, an unrelated local user, and an administrator used only for setup and controlled fault preparation.
3. Sign in to each supported profile once so Windows registers genuine profile roots.
4. Configure and record one standard-DPI display and one high-DPI or mixed-DPI display. Clear `FYNE_SCALE` and record the Fyne global scale setting for the clean cases.
5. Keep the exact candidate and evidence workspace on a local fixed volume.

## 3. Initialize evidence

From a non-elevated PowerShell session in the intended user's desktop:

```powershell
pwsh -NoProfile -File test/windows/Invoke-ReleaseCandidateAttended.ps1 `
    -Action Initialize `
    -MsiPath C:\Evidence\go-schedule_vX.Y.Z_windows_amd64.msi `
    -WorkspacePath C:\Evidence\go-schedule-vX.Y.Z `
    -Repository shruggietech/go-schedule `
    -Tag vX.Y.Z `
    -Commit 0123456789abcdef0123456789abcdef01234567 `
    -RunId 123456789 `
    -RunAttempt 1
```

The command refuses an existing workspace. Resume by adding observations to that workspace, never by silently recreating it. It also writes fail-closed setup/removal templates under `fragments`; every placeholder must be replaced with observed values and remain non-passing until reviewed.

## 4. Run the attended matrix

Follow `test/windows/README.md` for the detailed sequence. The matrix includes:

1. Fresh install, shortcut defaults and selection matrix, completion matrix, normal-token finish launch, cancel, maintenance, upgrade, invalid-input rejection, and rollback.
2. Normal-user CLI and GUI access, LocalSystem service identity, unrelated-user denial, and fresh-process PATH resolution.
3. Clean standard-DPI, clean high or mixed DPI, retained v0.9.1 profile, native first-launch measurements, window controls, and subsequent launch.
4. Daemon unavailable, access denied, timeout, stream disconnect, repeated refresh or reconnect, manual Retry, and recovery. Keep each repetition-sensitive failure observable for at least two minutes.
5. Public-interface manual and scheduled successful tasks, deliberate nonzero exit, and deliberate process-start failure. Retain all four canonical records in `attachments/tasks/task-runs.json` as described in `test/windows/README.md`.
6. Preserve, wipe, cancel, locked cleanup, genuine multiple profiles, reinstall recovery, clean reinstall, unaffected sentinels, and separately preserved security state.

Capture screenshots after sensitive values are excluded. Every window and error observation must reference visual evidence; window observations also reference the native measurement JSON. Record observed visible in-frame, modal, and top-level error counts; HWND enumeration alone cannot see Fyne canvas overlays.

## 5. Validate and package evidence

```powershell
pwsh -NoProfile -File test/windows/Invoke-ReleaseCandidateAttended.ps1 `
    -Action Finalize `
    -MsiPath C:\Evidence\go-schedule_vX.Y.Z_windows_amd64.msi `
    -WorkspacePath C:\Evidence\go-schedule-vX.Y.Z
```

Finalize fails on any missing, non-pass, inconsistent, or corrupted result. It writes the canonical evidence ZIP only after the shared Go validator accepts the bundle.

For an independent local check:

```powershell
go run ./scripts/windows-release-gate verify-bundle `
    --bundle C:\Evidence\go-schedule_vX.Y.Z_windows-attended-evidence.zip `
    --artifact C:\Evidence\go-schedule_vX.Y.Z_windows_amd64.msi `
    --candidate-manifest C:\Evidence\windows-candidate-manifest.json `
    --repository shruggietech/go-schedule `
    --tag vX.Y.Z `
    --commit 0123456789abcdef0123456789abcdef01234567
```

## 6. Upload and promote

1. Upload the canonical evidence ZIP to the existing draft release without replacing the MSI or candidate manifest.
2. Review the evidence and issue acceptance checklists. Close no issue whose required observation is absent.
3. Manually dispatch the Promote Release workflow with the exact tag.
4. Confirm it validates the draft state and exact assets, creates final all-asset checksums, and makes the release public only after validation.
5. Audit the published MSI and evidence hashes against the attended record, then update #94, #98, and #96 according to their individual acceptance criteria.

## Failure Recovery

- A failed or unavailable observation stays explicit. Correct the environment or product, reset the snapshot, and repeat the affected scenario against the same still-draft candidate when valid.
- If product bytes change, the previous evidence is permanently inapplicable. Stage a new candidate identity and start a new evidence bundle.
- If the draft release or exact MSI is lost, do not rebuild under the old evidence. Stage a new candidate.
- If promotion fails, inspect the complete gate diagnostics. The release must remain draft until all findings are resolved.
