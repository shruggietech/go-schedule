# Verification: Guided Windows Uninstall Entry

**Date**: 2026-09-02

**Branch**: `codex/041-guided-windows-uninstall`

## Result

S041's implementation and local repository verification are complete. The MSI
now suppresses Windows Installer's reduced-interface direct Remove entry while
retaining its full maintenance entry. The existing maintenance Remove control
continues to route through the package-owned inventory and preserve-or-wipe
dialogs.

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

The retained old-candidate rejection report is an ignored local artifact at
`dist/s041-old-candidate-rejection.md`; it is not release evidence.

## Registration evidence layers

- Source validation requires `ARPNOREMOVE=1`, forbids `ARPNOMODIFY`, and keeps
  the existing maintenance Remove route to `GoScheduleUninstallDlg`.
- Compiled-MSI inspection requires the corresponding Property-table values.
- Hosted lifecycle verification checks fresh install, repair, and upgrade for
  `NoRemove=1`, absent `NoModify` and `UninstallString`, a current native
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
