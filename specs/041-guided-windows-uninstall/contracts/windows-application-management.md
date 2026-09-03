# Contract: Windows Application Management and Removal Entry

## Attended Settings contract

1. Windows Settings lists one go-schedule product entry with its canonical product identity.
2. The entry exposes maintenance/Modify and does not expose direct Remove as the supported attended path.
3. Maintenance opens the full go-schedule package wizard.
4. Selecting Remove in that wizard opens `GoScheduleUninstallDlg` before execution.
5. Preserve is initially selected.
6. Selecting wipe requires `GoScheduleWipeConfirmDlg` before execution.
7. Cancel from maintenance, choice, or confirmation removes nothing.

Windows releases may label or arrange the maintenance action differently. Acceptance depends on the command and resulting wizard journey, not one exact label position.

## Native registration contract

For the installed current ProductCode under the machine application-registration root:

- `DisplayName` equals `go-schedule`;
- `NoRemove` equals integer `1`;
- `NoModify` is absent;
- `ModifyPath` is non-empty and identifies the current ProductCode through Windows Installer maintenance;
- `UninstallString` is absent;
- no second visible entry identifies the same product.

## Direct unattended contract

```powershell
msiexec.exe /x <candidate.msi> /qn /norestart
msiexec.exe /x <candidate.msi> GOSCHEDULE_REMOVE_DATA=0 /qn /norestart
msiexec.exe /x <candidate.msi> GOSCHEDULE_REMOVE_DATA=1 /qn /norestart
```

- The first two forms preserve application data.
- The third form requests the existing safe post-commit wipe.
- Values other than `0` and `1` fail before mutation.
- No form launches the GUI, documentation, or an attended dialog.

## Evidence boundary

Source verification and compiled-MSI inspection must fail when the registration contract drifts. Hosted lifecycle evidence must record installed registry values after fresh install, repair, and upgrade. Exact-candidate Windows 11 evidence in #94 must still observe the Settings-to-maintenance-to-removal-choice journey and cancellation safety.
