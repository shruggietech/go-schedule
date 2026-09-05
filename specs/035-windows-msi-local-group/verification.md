# Verification: Windows MSI Local-Group Recovery

## Baseline

- Issue #83 records v0.9.0 exit 1603, error 26421, and `0x80070005` while creating `goschedadmin`.
- Source and both validators expected `Domain="[ComputerName]"`; the MSI inspector did not query `Wix4Group`.
- The lifecycle harness omitted group/member/service/order, repair, upgrade, exit-code, and log diagnostics.
- This task process is not elevated and has no candidate MSI, so qualifying native lifecycle evidence is unavailable at planning time.

## Implementation Evidence

- Red regression baseline: focused `TestWindowsInstallerAdminGroupRejectsBrokenLifecycle` failed because the existing validator returned no diagnostic for `Domain="[ComputerName]"`.
- Green portable contracts: focused administrative-group and evidence-tooling Go tests pass; `build/windows/verify_wxs.ps1` passes; the rewritten lifecycle utility passes the vendored PowerShell compliance checker.
- Blocking cross-artifact analysis: 14 functional requirements, 8 success criteria, and 25 tasks; no critical or high consistency, ambiguity, constitution, or coverage findings remain.

## Canonical Verification

- `format`: PASS, no files reported.
- `vet`: PASS.
- `lint`: PASS, zero issues.
- `race`: PASS, including integration tests.
- `gui`: PASS.
- `coverage`: PASS (`engine` 86.4%, `schedule` 89.2%, `timezone` 91.3%, `store` 84.1%, `catchup` 88.9%, `logbus` 91.1%).
- `docs`: PASS.
- `automation`: PASS.

## Native Windows 11 Lifecycle Evidence

Unavailable in this task process. The host is Windows 11 build 26200 and has no existing product, `goschedd` service, or `goschedadmin` group, but the process runs as non-elevated `NOVX999\\h8rt3rmin8r` and the worktree contains no candidate MSI. Fresh, v0.9.0-upgrade, and post-refresh access evidence therefore did not run. Static validation is not runtime proof, T024 remains open, S035 remains In Progress, and the pull request must use `Refs #83`.

## Delivery Audit

- Three independent agents reviewed general correctness, Windows installer/security behavior, and spec/test evidence. All findings were fixed, and all three final re-reviews reported no actionable issue at confidence 80% or higher.
- Pull-request review found that reinstall properties preceded the MSI path and that declined `ShouldProcess` actions fell through to lifecycle assertions. The argument builder now keeps the package immediately after the operation switch, and every skipped destructive phase stops without false evidence or cleanup side effects.
- UTF-8 without BOM, replacement-character, trailing-whitespace, task/checklist, issue-reference, and worktree-scope audits completed locally.
