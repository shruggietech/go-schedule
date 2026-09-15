# Implementation Plan: Final v1.4.0 boundary and qualification

**Branch**: `codex/088-release-boundary-qualification` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Status**: Source preparation implemented; verification and review in progress. Sandbox-only execution approved with a maintainer waiver for unavailable native tests.

## Summary

Bundle #231, #228, and #226. Reconcile cumulative v1.4.0 metadata, obtain a reviewed source boundary, stage that exact boundary under explicit release authorization, and qualify its exact MSI using the existing attended collector. Public promotion remains separately authorized.

## Technical Context

Reuse the Go application, Wails desktop, PowerShell qualification scripts, GitHub Actions release workflow, and existing verification gates. Target disposable Windows 11 environments with ordinary-user process evidence, independent registered profiles, fresh installation, v1.1.1 upgrade, standard and scaled display measurements, and native UI access. All 47 required observations and their attachments remain mandatory.

## Constitution Check

Preserve exact-artifact provenance, issue-level acceptance criteria, reviewed-source staging, and fail-closed qualification. Missing native evidence cannot be replaced by hosted contract tests, fixture results, or screenshots of another candidate. No gate reduction or new testing framework is proposed.

## Research Findings

- Reviewed main `c48ee096251b33187deda61200eaa25edd0e36ab` includes the S087 selector-height repair.
- The draft v1.4.0 tag identifies `951d864d9683ec3bdcb1e37535ecddd67738bf13`, which predates that repair.
- S083-S085 improvements and the S087 repair remain under Unreleased. Tag-specific release notes need cumulative reconciliation. Retain the existing dated heading unless the maintainer authorizes a date change.
- No Hyper-V module, VM management tools, configured PowerShell remoting sessions, or active Sandbox session were found. The Hyper-V management CIM namespace is unavailable.
- Sandbox availability alone does not establish the required multi-profile, native process, display, and upgrade capabilities. Environment access and provisioning authority remain unresolved.

## Execution Order

1. Reuse Windows Sandbox without host changes or restart. Inventory supported checks and list unavailable checks with the explicit maintainer waiver rather than blocking on another VM platform.
2. Complete spec-kit research, design artifacts, tasks, and consistency analysis.
3. Reconcile `CHANGELOG.md`, `.github/release-notes/v1.4.0.md`, and affected documentation. Run all eight local gates and the GitHub formatter, publish a structured preparation PR, handle at most two review rounds, and obtain green latest-head CI.
4. Obtain the maintainer's preparation merge and explicit release-operation authorization. The Release workflow requires successful exact-commit main push CI before staging reviewed-main artifacts.
5. Verify all eight staged assets and manifest identities. Complete fresh, upgrade, lifecycle, display, desktop, and task workflow evidence against the final exact MSI using existing candidate and bundle validators.
6. Report each issue's completed and remaining criteria. Keep public-promotion criteria open until separately authorized and verified.

## Project Structure

Reuse `.github/release-notes/`, `.github/workflows/`, `CHANGELOG.md`, operator documentation, `test/windows/`, and `scripts/`. Slice artifacts live here. Historical assets and evidence under ignored `dist/` remain intact.

## Pending Decisions

Do not install on the active development host, change host security settings, enable virtualization features, reboot, or provision another virtualization platform. Final candidate execution depends on the reviewed-boundary merge. Unavailable Sandbox checks may remain untested under the explicit maintainer waiver; full qualification still requires actual passing evidence.

## Maintainer Override (2026-09-15)

The maintainer rejected host restarts and requested the previously used Sandbox approach, then explicitly authorized incomplete items to be released without testing. This changes the earlier all-environments prerequisite: execute supported tests, retain unsupported observations as untested, and carry a visible release waiver. Do not manufacture passing collector fragments, call partial evidence fully qualified, disable CI, or silently treat known failures as approved. Reviewed-main staging still follows the preparation merge. Any release-gate implementation change needed to represent waived qualification must be explicit, narrowly scoped, tested, and reviewed rather than bypassed outside the normal workflow.
