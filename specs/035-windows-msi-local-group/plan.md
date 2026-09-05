# Implementation Plan: Windows MSI Local-Group Recovery

**Branch**: `codex/083-fix-msi-local-group` | **Date**: 2026-08-31 |
**Spec**: [spec.md](spec.md)

## Summary

Remove the group domain that makes WiX select its unelevated domain-group path, then harden every evidence layer that previously blessed the defect. Portable tests and source validation reject any non-empty group domain, compiled-MSI inspection records the group row, and the disposable-host lifecycle verifier records fresh install, repair, reinstall, v0.9.0 upgrade, uninstall, security state, service state, exit codes, and verbose-log evidence.

## Technical Context

**Language/Version**: Go 1.25, WiX 6.0.2 XML, PowerShell 7, Markdown

**Primary Dependencies**: Go standard library, Windows Installer automation, PowerShell local-account/service cmdlets, WiX Util extension

**Storage**: Versioned installer and verifier sources; operator-selected evidence and verbose-log outputs

**Testing**: Go XML contracts, PowerShell source validation, compiled-MSI inspection, disposable Windows 11 lifecycle runs, eight canonical gates

**Target Platform**: Windows 11 x64 installation; portable contracts in CI

**Project Type**: Installer regression fix and verification tooling

**Performance Goals**: No scheduler runtime change or installed overhead

**Constraints**: Preserve installing-user identity, preserve group/membership on uninstall, hidden noninteractive child processes, UTF-8 without BOM, never count unavailable native evidence as passing

**Scale/Scope**: One group, one account, one service, separate fresh and upgrade disposable baselines, one candidate MSI, one v0.9.0 baseline MSI

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One semantic authoring removal, explicit diagnostics, and hidden child-process launches. |
| II. Testing Standards | PASS | Regression mutations, compiled-artifact inspection, and native lifecycle evidence cover distinct failure classes. |
| III. UX Consistency | PASS | Troubleshooting gains an exact verbose-log workflow and actionable checks. |
| IV. Performance | PASS | No runtime path changes; verification runs only during packaging review. |
| V. Autonomous Execution | PASS | Full spec-kit sequence, analyze gate, canonical verification, commit, authorized push, and PR. |

### Post-design re-check

Pinned `build/windows/**` and `docs/INSTALL-windows.md` change, so a dated changelog decision is mandatory. Native Windows 11 evidence remains required; if elevation or artifacts are unavailable, the PR uses `Refs #83` and leaves the issue open.

The checklist helper requires `plan.md`, conflicting with the mandated checklist-before-plan order. As in S034, the checklist was produced directly from the completed spec. This preserves the gate without modifying Spec Kit.

## Architecture and Decision Log

### Omit the group domain, preserve the user domain

Remove `Domain="[ComputerName]"` only from `GoScheduleAdminGroup`. WiX 6 treats every non-empty group domain as a domain operation; omission selects its elevated local-group operation. `InstallingUser` retains `Domain="[%USERDOMAIN]"` because it identifies the existing account. An empty attribute was rejected in favor of canonical omission, and a custom action was rejected as unsupported duplication.

### Test semantics and inspect the compiled row

Both source validators require an empty effective group domain and retain all other lifecycle attributes. A `[ComputerName]` mutation must fail with a routing-specific diagnostic. The MSI inspector queries `Wix4Group` because custom-action presence cannot reveal which runtime path the row selects.

### Split fresh and upgrade evidence

The fresh scenario starts without product or group, then installs, repairs, reinstalls, and uninstalls the candidate. A separate reset host preprovisions the durable group/member, installs the otherwise-broken v0.9.0 baseline, upgrades to the candidate, and uninstalls. One host cannot supply both clean baselines because uninstall deliberately preserves the group.

### Record runtime chronology honestly

Each operation records exit code, log path/hash, group SID and members, service state, product/PATH state, access-denied and rollback diagnostics, and local group action ordering before service startup. A separate post-refresh, non-elevated access-probe phase records the current token's group SID and client success before cleanup. Missing observations are failures or unavailable, never passes. GitHub's Windows Server runner is supplementary, not a substitute for the required supported Windows 11 run.

## Project Structure

```text
build/windows/{goschedule.wxs,verify_wxs.ps1}
test/integration/windows_installer_contract_test.go
test/windows/{README.md,inspect-installer.ps1,Invoke-InstallerLifecycle.ps1}
docs/INSTALL-windows.md
specs/035-windows-msi-local-group/{spec,plan,research,data-model,quickstart,verification}.md
specs/035-windows-msi-local-group/contracts/
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend existing installer, validation, lifecycle, and documentation surfaces. Add no workflow, runtime package, or ninth gate.

## Complexity Tracking

No constitution violation requires justification. Two lifecycle scenarios are necessary because the required preserved group prevents baseline reuse.
