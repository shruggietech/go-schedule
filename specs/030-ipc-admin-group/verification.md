# Verification Record: S030

## Baseline (2026-08-30)

- Local `main` and `origin/main` pointed to `783e3eb`; the worktree was clean before `codex/030-ipc-admin-group` was created.
- `config.AdminGroup` defaulted to `goschedadmin` but no runtime code read it.
- Windows used a fixed descriptor granting Authenticated Users read/write. Unix created the endpoint and parent with process-default permissions.
- The Windows MSI did not create an administrative group or enroll the installing user. Linux and macOS guides did not provision one.
- WiX tool, UI extension, and Util extension were pinned to 5.0.2; WiX group creation first became available in v6.

## Test-First Evidence

- Config tests cover the secure default, explicit-empty compatibility value, and rejection of leading, trailing, or whitespace-only group names.
- Platform tests exercise the failure contracts as first-class cases: missing and non-group Windows accounts never open a listener; missing Unix groups never create an endpoint; unsafe existing custom parents remain unchanged; and post-listen permission failures remove the socket.
- Mutation fixtures reject incomplete installer group creation, destructive uninstall behavior, accidental user creation, missing membership, feature omission, and stale WiX pins.
- Restricted and compatibility tests assert exact descriptors or modes plus the structured startup evidence and warning order, so a permissive fallback cannot pass silently.

## Analysis Gate

- Initial cross-artifact analysis covered 18 functional requirements, 6 measurable outcomes, and 26 tasks with no constitution conflict or uncovered requirement.
- It found two planning inconsistencies: the analysis task appeared after implementation, and SC-001/SC-004 implied manual multi-account or MSI lifecycle observation despite the automated no-VM scope.
- The analysis task moved into the foundational phase. SC-001 now measures platform policy contracts, and SC-004 measures installer source plus pinned-tool build contracts.
- Implementation planning then caught a custom-path hazard: rewriting every endpoint parent could take over a shared directory such as `/tmp`. The contract now mutates only the managed default parent or a newly created custom parent; an existing custom parent is verified without mutation and rejected when unsafe.

## Final Verification

- Focused Windows config, IPC, daemon, and installer tests passed. Linux IPC tests passed under WSL. Linux and macOS IPC/daemon test binaries cross-compiled successfully.
- `build/windows/verify_wxs.ps1` passed against the final source.
- WiX 6.0.2 plus UI and Util 6.0.2 compiled a candidate MSI from the final authoring. Direct MSI table inspection recorded `Wix6ConfigureGroups_X64` at sequence 3998, `Wix4ConfigureUsers_X64` at 3999, `InstallFiles` at 4000, and `StartServices` at 5900.
- The first canonical run reached the docs gate and correctly rejected an unsupported `json` fence. The fence was changed to the repository's supported `text` category; focused docs and automation gates then passed.
- Final foreground `scripts/verify.sh all`: PASS for format, vet, lint, race, GUI, coverage, docs, and automation.
- Coverage gate results: engine 88.1%, schedule 89.0%, timezone 91.3%, store 86.9%, catchup 88.9%, logbus 91.1%.
- Final audit found no trailing whitespace, UTF-8 BOM, or mojibake. The dated Unreleased decision covers `.github/workflows/release.yml`, `build/windows/**`, and the platform installation guides. Release-facing language contains `Closes #13`.

## Review Follow-Through (2026-08-30)

- Codex correctly identified that `[LogonUser]` alone does not qualify domain or Entra identities. The installer now supplies `Domain="[%USERDOMAIN]"`, and both portable and PowerShell contracts reject an unqualified installing user.
- The initial macOS race job exposed test-only Unix-socket paths longer than macOS accepts. Unix IPC tests now use a bounded `gs-ipc-*` temporary directory while retaining automatic cleanup and the same permission assertions.
- Post-review `scripts/verify.sh all`: PASS across all eight gates. A rebuilt WiX 6.0.2 MSI decompiled with the expected `InstallingUser` domain value and retained group-before-user-before-service sequencing.
