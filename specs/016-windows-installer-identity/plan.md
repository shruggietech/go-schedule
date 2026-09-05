# Implementation Plan: Windows Installer Identity and Verification

**Branch**: `codex/016-windows-installer-identity` | **Date**: 2026-08-27 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/016-windows-installer-identity/spec.md`

## Summary

Give the Windows MSI one explicit canonical icon relationship shared by the Start Menu shortcut and installed-apps entry, preserve the existing executable resource pipeline for the running GUI, and make both contracts regression- tested. Add Windows-only artifact inspection and clean-install lifecycle tools that distinguish what was proven from what requires a disposable desktop.

The local host can prove source contracts and candidate MSI tables. It cannot honestly prove the clean-machine or native-shell observations because go-schedule v0.8.0 is already installed and Windows Sandbox is not enabled. Those evidence classes remain `unavailable`, and issues #16 and #33 remain open.

## Technical Context

**Language/Version**: Go 1.25.0, PowerShell 7, WiX Toolset 5.0.2 XML.

**Primary Dependencies**: Go standard library XML decoder for portable contract tests; Windows Installer COM automation and `msiexec` for Windows-only artifact and lifecycle verification; existing `goversioninfo` v1.4.1 release step for the GUI executable resource. No new runtime or module dependency.

**Storage**: Evidence is Markdown written to an operator-selected path. Product data storage and migrations are untouched.

**Testing**: Table-driven Go regression tests mutate each protected installer and executable-resource relationship independently; the existing PowerShell WiX sanity check gains specific icon-contract assertions; a candidate MSI is built and inspected without installing it; canonical `sh scripts/verify.sh all` remains the merge gate.

**Target Platform**: Windows amd64 for MSI artifacts and native observations; all CI platforms for source-contract regression tests.

**Project Type**: Cross-platform daemon, CLI, and Fyne desktop application with a Windows MSI release artifact.

**Performance Goals**: None engaged. Installer metadata and offline tests do not touch scheduler dispatch paths or budgets.

**Constraints**: `build/**` and `.github/workflows/release.yml` are pinned. The WXS and verifier changes require a dated changelog decision. The release workflow is inspected but does not need modification because it already embeds the canonical icon into the GUI executable. No release, uninstall of the operator's current installation, host virtualization change, or reboot.

**Scale/Scope**: One WXS icon declaration, two references, verifier assertions, one portable Go regression file, two Windows evidence scripts, one maintainer guide, and Spec-Kit artifacts. No application Go behavior changes.

## Constitution Check

*Re-checked after Phase 1 design. All gates pass; no Complexity Tracking entries.*

- **I. Code Quality.** The package declaration remains explicit and readable. The tests use the standard XML decoder instead of treating XML as arbitrary text. PowerShell reports the exact broken relationship. No error is ignored.
- **II. Testing Standards.** This bug fix is test-first: current source fails the new icon-contract test, then the WXS change makes it pass. Every new relationship has an independent negative mutation. Existing safety-critical scheduling, migration, recovery, concurrency, and IPC tests are untouched.
- **III. User Experience Consistency.** Installed and running Windows surfaces reference the same official icon. Evidence uses `proven`, `failed`, and `unavailable` consistently and never calls an unavailable desktop check green.
- **IV. Performance Requirements.** Not engaged. No runtime or hot-path code changes.
- **V. Autonomous Build-Phase Execution.** Full Spec-Kit sequence, blocking analysis, local commit, and one halt before publication. Pinned WXS and verifier changes receive a dated changelog decision.

## Project Structure

### Documentation (this feature)

```text
specs/016-windows-installer-identity/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── contracts/
│   ├── installer-identity.md
│   └── verification-evidence.md
├── quickstart.md
├── verification.md
├── checklists/
│   ├── requirements.md
│   └── installer.md
└── tasks.md
```

### Source Code (repository root)

```text
build/windows/
├── goschedule.wxs                 # canonical Icon + ARP/shortcut refs [pinned]
└── verify_wxs.ps1                 # source sanity assertions             [pinned]

test/integration/
└── windows_installer_contract_test.go  # portable mutation regressions

test/windows/
├── README.md                      # disposable-host procedure
├── inspect-installer.ps1          # MSI table evidence, no installation
└── install-lifecycle.ps1          # clean install/reinstall/uninstall evidence

CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Keep package authoring in `build/windows`, portable source regressions in the existing integration test package, and state-changing Windows procedures in a clearly named `test/windows` surface. No production package or reusable library is introduced for build metadata.

## Design

### One icon row, two installer consumers

Add one package-level icon whose source is the existing multi-resolution `cmd/gosched-gui/icon.ico`. `ARPPRODUCTICON` and the Start Menu shortcut both reference its identifier. This is the smallest supported WiX shape and prevents two installer surfaces from drifting to separate files.

### Preserve the existing executable-resource path

The release workflow already generates a 64-bit COFF resource from the same canonical icon before building `gosched-gui.exe`. No workflow edit is needed. The portable regression test asserts that source, output, and architecture contract so a future release edit cannot silently remove the native fallback.

### Layer evidence instead of overstating it

Portable tests prove source relationships. `inspect-installer.ps1` proves the compiled MSI's Property, Icon, Shortcut, and Environment table rows without installation. `install-lifecycle.ps1` refuses a contaminated baseline, then records install, new-shell commands, reinstall, and uninstall on a disposable Windows host. Native title/taskbar and shell icon observations remain explicit human evidence. Each layer records `proven`, `failed`, or `unavailable`.

### Issue closure boundary

The deterministic WXS change and candidate MSI inspection complete #24 and #25. Issue #16 explicitly requires a published MSI on a clean host, and #33 requires native running-GUI observation; both stay open until that evidence exists. The PR will use `Closes #24`, `Closes #25`, `Refs #16`, and `Refs #33`.

## Complexity Tracking

No constitution violations. Table intentionally empty.
