# Verification Evidence: Windows Installer Identity and Verification

**Date**: 2026-08-27
**Branch**: `codex/016-windows-installer-identity`

## Baseline

| Evidence | Status | Observation |
|---|---|---|
| Latest published artifact | proven | v0.8.0 published 2026-07-23; Windows asset `go-schedule_v0.8.0_windows_amd64.msi` |
| Published MSI identity tables | failed | No Icon table, no `ARPPRODUCTICON`, and `GuiShortcut.Icon_` is empty |
| Published MSI PATH table | proven | `PathEnv`, machine append/remove encoding, `[INSTALLFOLDER]`, component `Gosched` |
| Installed-app registry identity | failed | Installed v0.8.0 `DisplayIcon` is empty |
| Installed Start Menu identity | failed | Shortcut targets `gosched-gui.exe`; `IconLocation` is `,0` |
| Current machine PATH cardinality | proven | Exactly one `C:\Program Files\go-schedule\` entry |
| Clean-host prerequisite | unavailable | v0.8.0 is installed and resolves from PATH on this host |
| Windows Sandbox prerequisite | unavailable | `WindowsSandbox.exe` is absent; optional-feature inspection requires elevation |

The published MSI was downloaded from the v0.8.0 GitHub release and inspected
read-only through Windows Installer COM. No install, repair, uninstall, service
change, PATH change, virtualization change, or reboot occurred.

## Test-first evidence

The focused test failed before the installer edit with exactly three expected
contract failures:

- canonical Icon declaration is missing;
- `ARPPRODUCTICON` property is missing;
- `GuiShortcut` does not reference the canonical Icon.

The independent mutation cases already pass because they exercise an isolated
valid fixture and remove or redirect one relationship at a time.

After adding one canonical package Icon plus its ARP and shortcut consumers,
`go test ./test/integration -run WindowsInstaller -count=1` passes. The release
workflow resource contract also passes without a workflow edit.

## Candidate artifact evidence

**Status: proven.** A disposable candidate MSI was compiled with WiX 5.0.2 from
the changed WXS and placeholder staged executables. This proves package tables,
not executable behavior. SHA-256:
`f21e77ad70529a7353e866797bd9f69e9a1b5e4ecb59e7707406b2c9ae093959`.

| Table relationship | Observation |
|---|---|
| Icon | `GoSchedule.ico` |
| Property | `ARPPRODUCTICON=GoSchedule.ico` |
| Shortcut | `GuiShortcut.Icon_=GoSchedule.ico` |
| Environment | `=-*PATH`, `[~];[INSTALLFOLDER]`, component `Gosched` |

The same inspector rejects published v0.8.0 with three specific failures for
its missing ARP property, Icon row, and shortcut Icon reference.

## Clean lifecycle evidence

**Status: unavailable.** The guarded lifecycle tool refused before mutation
because PowerShell was not elevated, v0.8.0 is installed, `gosched` resolves,
the install directory and product registry entry exist, and machine PATH already
contains the directory. Issue #16 remains open unless a published MSI runs on a
clean disposable or physical Windows environment.

## Native desktop evidence

**Status: unavailable.** The portable test proves the existing release workflow
generates a 64-bit executable resource from the canonical icon. This host cannot
provide a clean native observation. Issue #33 remains open unless the installed
GUI title area and taskbar are observed in a suitable Windows desktop session.

## Repository verification

- `sh scripts/verify.sh all`: passed all eight canonical gates.
- Format: passed.
- Vet: passed.
- Lint: passed with zero issues.
- Race: passed, including integration tests.
- GUI: passed for GUI and viewmodel packages.
- Coverage: passed all enforced package thresholds.
- Docs: passed with 11 pages clean.
- Automation: passed.
- `Invoke-ScriptAnalyzer`, excluding only `PSUseBOMForUnicodeEncodedFile` to
  preserve the repository's required UTF-8-without-BOM encoding: passed.

The focused Go contract test, WXS verifier, candidate MSI inspector, and guarded
lifecycle refusal were rerun after the final script corrections and passed their
respective expected outcomes.

## Review remediation

The three AI review findings on PR #48 were accepted and corrected:

- paused native observations are captured as validated statuses plus optional
  notes and rendered into final evidence, including later failure paths;
- failed lifecycle evidence is written only after cleanup and records the final
  installed-product, install-directory, and machine-PATH state;
- artifact inspection and lifecycle evidence require an explicit `candidate` or
  `published` class and origin, with published origins restricted to this
  repository's HTTPS release-asset URLs.

Focused regression tests, PowerShell parsing and analysis, candidate and
published artifact checks, the non-mutating lifecycle refusal, encoding audits,
and all eight canonical verification gates passed after these corrections.

## Issue disposition

- #24: eligible to close after merge; explicit shortcut icon proven in source
  and candidate MSI.
- #25: eligible to close after merge; explicit `ARPPRODUCTICON` proven in source
  and candidate MSI.
- #16: remains open; published clean-install lifecycle evidence unavailable.
- #33: remains open; native running-GUI observation unavailable.
