# Verification: Guided Windows Uninstall Entry

**Date**: 2026-09-02

**Branch**: `codex/041-guided-windows-uninstall`

## Result

S041's first implementation passed local repository verification, but hosted
installed-state evidence correctly rejected its one-property assumption. The
corrected MSI suppresses Windows Installer's reduced-interface direct Remove
entry, owns the current ProductCode's `/I` maintenance registration, and uses a
package-owned maintenance page whose Remove control routes through the inventory
and preserve-or-wipe dialogs.

This result does not claim exact-candidate attended acceptance. GitHub issue
#98 remains open until the Windows 11 journey and cancellation behavior are
observed against the byte-identical staged candidate under #94.

## Red-to-green evidence

Before `ARPNOREMOVE` was added to the WiX source:

| Check | Expected red result |
| --- | --- |
| Focused Go installer lifecycle and evidence-tooling contracts | Failed with `ARPNOREMOVE must disable the reduced-interface direct removal entry` |
| `build/windows/verify_wxs.ps1` | Failed because `ARPNOREMOVE=1` was absent |
| Updated compiled-MSI inspector against the retained S039 candidate | Failed because `Property.ARPNOREMOVE` was empty rather than `1` |

After the declarative WiX change:

| Check | Result |
| --- | --- |
| `go test ./test/integration -run 'WindowsInstaller|Cleanup' -count=1` | PASS |
| `go test ./cmd/gosched-cleanup/... -count=1` | PASS |
| `build/windows/verify_wxs.ps1` | PASS |
| PowerShell parser checks for all three modified verification scripts | PASS |
| `scripts/verify.sh automation` after lifecycle bookkeeping correction | PASS |

## Hosted correction evidence

PR #107's first Windows MSI run built and inspected the package successfully,
then failed the fresh installed-state assertion:

```text
default install ModifyPath '' does not open maintenance for
{C71909CC-85FE-4764-8713-E8BB82CC65E7}
```

Inspection of that compiled MSI also showed WiX 6.0.2's stock
`MaintenanceTypeDlg.RemoveButton` had a `Disable` condition containing
`ARPNOREMOVE`. The original plan would therefore have hidden the unsafe direct
action without leaving a usable guided action.

New regression contracts were made red before the correction. They required:

- an MSI-owned expandable `MsiExec.exe /I[ProductCode]` Registry-table row in
  the native current-product uninstall key;
- installation of that row through the required Main feature;
- a package-owned maintenance page whose Remove control is not disabled by
  `ARPNOREMOVE`;
- an exact control-event route from that Remove control to
  `GoScheduleUninstallDlg`.

The source contracts and PowerShell structural checks turned green after the
owned registration component and maintenance page were added. Positive compiled
evidence then passed locally with pinned WiX 6.0.2:

- candidate SHA-256: `1196eca252e9c275c205ba8e1fcb4aa579032103d01901b8ec9c77c22035fef4`;
- Registry value: `#%MsiExec.exe /I[ProductCode]`;
- ordered Remove events: set maintenance mode, reset wipe property, reset the
  visible choice to preserve, then open `GoScheduleUninstallDlg`;
- no `ARPNOREMOVE` ControlCondition on the package-owned Remove control.

Installed-state evidence remains pending the corrected hosted run because the
local workstation is not the disposable environment authorized for MSI
lifecycle mutation.

The retained old-candidate rejection report is an ignored local artifact at
`dist/s041-old-candidate-rejection.md`; it is not release evidence.

## Registration evidence layers

- Source validation requires `ARPNOREMOVE=1`, forbids `ARPNOMODIFY`, requires
  the owned current-ProductCode `/I` registration, and keeps the package-owned
  maintenance Remove route to `GoScheduleUninstallDlg` enabled.
- Compiled-MSI inspection requires the corresponding Property, Registry,
  FeatureComponents, Dialog, Control, ControlCondition, and ControlEvent rows.
- Hosted lifecycle verification checks fresh install, repair, and upgrade for
  `NoRemove=1`, absent `NoModify` and `UninstallString`, the current MSI-owned
  `ModifyPath`, and exactly one visible go-schedule entry.
- Exact Windows 11 Settings wording, interaction, defaults, and cancellation
  remain attended evidence owned by #94.

## Canonical verification

Command:

```text
C:\Program Files\Git\bin\bash.exe scripts/verify.sh all
```

The initial run passed the first seven gates and correctly rejected stale Spec
Kit lifecycle bookkeeping at `automation`. After aligning the feature state and
completed tasks, a clean uninterrupted rerun passed all eight gates:

1. `format`
2. `vet`
3. `lint` (`0 issues.`)
4. `race`
5. `gui`
6. `coverage` (engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, logbus 91.1%)
7. `docs`
8. `automation` (including release staging/promotion and Spec Kit lifecycle)

## Local review

The final review covers requirements traceability, the prior S039 removal
route, WiX and Windows Installer registration semantics, failure behavior in
the hosted registry probe, direct silent removal preservation, documentation,
UTF-8 without BOM, mojibake, and diff integrity.

No finding at or above the 80-percent confidence threshold remained after the
review. All 21 changed or new text files decoded as strict UTF-8, none had a
UTF-8 BOM or mojibake signature, and `git diff --check` passed.

After the hosted correction, the focused installer and cleanup suites, pinned
WiX build and inspection, PowerShell parsing, and all eight canonical gates
passed again before the corrective commit. The 15 corrected text files also
passed strict UTF-8, BOM, mojibake, and diff-integrity checks.
