# Implementation Plan: Guided Windows Uninstall Entry

**Branch**: `codex/041-guided-windows-uninstall` | **Date**: 2026-09-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/041-guided-windows-uninstall/spec.md`

## Summary

Repair the Windows Settings entry that currently starts a reduced-interface direct MSI removal and bypasses S039's preserve-or-wipe pages. Suppress direct Remove, author an MSI-owned maintenance command for the current ProductCode, and use a package-owned maintenance page whose internal Remove control routes through the inventory and confirmation pages. Extend source, compiled-MSI, and hosted lifecycle checks to prove the registration values while retaining direct silent preserve/wipe commands.

## Technical Context

**Language/Version**: WiX Toolset 6.0.2 XML; PowerShell 7; Go 1.25.0 integration tests

**Primary Dependencies**: Windows Installer 5; WiX UI and Util extensions already pinned by CI; Windows Installer COM for compiled-package inspection; Windows registry on disposable hosted runners

**Storage**: Windows Installer product registration under the machine Uninstall registry; existing `C:\ProgramData\goschedule` and per-user Fyne roots are unchanged

**Testing**: Failing Go source-contract mutation tests, PowerShell source checks, compiled-MSI inspection, hosted Windows install/upgrade/repair/preserve/wipe lifecycle probe, and canonical eight-gate verification

**Target Platform**: 64-bit per-machine MSI on Windows 11; hosted Windows Server supplies non-attended package and lifecycle evidence

**Project Type**: Desktop application packaging and lifecycle repair

**Performance Goals**: No application runtime impact; static and compiled registration checks complete in seconds; hosted lifecycle duration does not materially exceed the established S039 matrix

**Constraints**: One visible product entry; no custom ARP shadow entry; no bootstrapper; no custom-action popup; direct silent `/x` must remain functional; preserve remains default; hidden/noninteractive console execution; UTF-8 without BOM; no release, tag, or merge

**Scale/Scope**: One MSI property, one MSI-owned registry component, one package-owned maintenance page, three existing verification layers, Windows install documentation, and S041 evidence artifacts

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | The repair is declarative, uses the operating system's supported MSI registration semantics, and adds structural diagnostics at each existing verification layer. |
| II. Testing Standards | PASS | Regression assertions are written failing first in source and integration contracts; compiled and hosted evidence cover generated and installed state. Existing safety tests remain required. |
| III. UX Consistency | PASS | Windows users enter the full product maintenance wizard before removal and see the existing plain-language preserve-or-wipe inventory. Documentation names the exact Settings journey. |
| IV. Performance | PASS | No scheduler or GUI runtime path changes. One declarative registration property adds no material installation cost. |
| V. Autonomous Execution | PASS WITH RECORDED OPERATOR OVERRIDE | Full Spec Kit and analyze gates remain mandatory. The operator explicitly authorized automatic S041 push and PR plus at most one additional Codex review round; merge, tag, and release remain outside authority. |

### Post-design re-check

All principles remain satisfied. The MSI source and verification scripts are pinned release artifacts whose modification is required by #98; the decision is recorded in `CHANGELOG.md`. No additional executable, dependency, privilege boundary, or data deletion path is introduced.

## Architecture and Decision Log

### Route attended system management through native MSI maintenance

Set the Windows Installer `ARPNOREMOVE` product property to `1` and deliberately leave `ARPNOMODIFY` unset. Windows Installer suppresses its generated direct `UninstallString`/Remove action. Hosted evidence then proved that it also omits `ModifyPath`, contrary to the original one-property plan. An MSI-owned registry component therefore writes the expandable `MsiExec.exe /I[ProductCode]` maintenance command into the same native product key.

WiX 6.0.2's stock `MaintenanceTypeDlg` also disables its Remove control when `ARPNOREMOVE` is set. S041 owns the maintenance page so that external direct removal remains suppressed while internal guided removal remains available. Change and Repair retain their standard modes; Remove resets the existing choice to preserve and routes to `GoScheduleUninstallDlg`. Direct Windows Installer command-line and API removal remain supported.

### Keep unattended removal explicit and stable

The declarative ARP property does not block `msiexec /x`. Existing commands remain authoritative:

- absent or `GOSCHEDULE_REMOVE_DATA=0`: remove software and preserve application data;
- exact `GOSCHEDULE_REMOVE_DATA=1`: remove software and request the post-commit safe wipe;
- any other value: fail before product mutation.

No condition on the cleanup action, completion actions, shortcut features, or service lifecycle changes in S041.

### Prove source, package, and installed registration separately

`build/windows/verify_wxs.ps1` rejects source that omits `ARPNOREMOVE=1`, suppresses maintenance through `ARPNOMODIFY`, omits the owned `/I` registry value, or lets `ARPNOREMOVE` disable the package-owned Remove control. `test/integration/windows_installer_contract_test.go` carries the same regression and mutation cases. `test/windows/inspect-installer.ps1` proves Property, Registry, Dialog, Control, ControlCondition, and ControlEvent rows. `test/windows/Invoke-InstallerContractCI.ps1` asserts that installed state has `NoRemove=1`, no `NoModify`, no `UninstallString`, the owned current-ProductCode `ModifyPath`, and one visible go-schedule identity.

Compiled rows prove what Windows Installer will register; the installed registry probe proves standard actions actually produced that state. Neither is misrepresented as the final Windows 11 Settings observation, which remains in #94.

### Reject custom registration and execute-sequence UI

A hidden native MSI plus a hand-authored shadow ARP key could force any command, but duplicates Windows Installer's ownership and complicates repair, upgrade, rollback, product identity, and removal. A custom executable/bootstrapper adds a new signed artifact and release architecture. An execute-sequence popup would contaminate silent administration and run across an elevation boundary. Those alternatives are disproportionate and less reliable than the documented maintenance route.

## External Research

- [Microsoft: Configuring Add/Remove Programs with Windows Installer](https://learn.microsoft.com/en-us/windows/win32/msi/configuring-add-remove-programs-with-windows-installer) documents that `ARPNOREMOVE` hides direct Remove while Change can still remove products whose package UI offers removal.
- [Microsoft: Windows Installer properties for the Uninstall registry key](https://learn.microsoft.com/en-us/windows/win32/msi/uninstall-registry-key) identifies `ModifyPath` and `UninstallString` as Windows Installer-generated values.
- [Microsoft: Registry table](https://learn.microsoft.com/en-us/windows/win32/msi/registry-table) defines MSI-owned, formatted, expandable registry values and their component lifecycle.
- [WiX 6.0.2: MaintenanceTypeDlg source](https://github.com/wixtoolset/wix/blob/v6.0.2/src/ext/UI/wixlib/MaintenanceTypeDlg.wxs) shows the stock Remove control is disabled by `ARPNOREMOVE`.
- [Microsoft: msiexec](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/msiexec) confirms `/x` remains the direct administrative uninstall mechanism and UI levels are independently selectable.

## Project Structure

```text
build/windows/
├── goschedule.wxs
└── verify_wxs.ps1
test/integration/windows_installer_contract_test.go
test/windows/
├── inspect-installer.ps1
├── Invoke-InstallerContractCI.ps1
└── README.md
docs/INSTALL-windows.md
specs/041-guided-windows-uninstall/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/windows-application-management.md
├── checklists/
├── tasks.md
└── verification.md
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend the existing MSI authoring and its three established verification layers. The repair changes product registration, not application runtime, so no new Go production package or executable is warranted.

## Complexity Tracking

No constitution violation or exceptional complexity is introduced.
