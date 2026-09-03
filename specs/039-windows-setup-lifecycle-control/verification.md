# S039 Verification: Windows Setup Lifecycle Control

**Date**: 2026-09-02

**Branch**: `codex/039-windows-setup-lifecycle-control`

**Disposition**: Local implementation verification complete. Pull-request CI supplies the hosted Windows Server silent lifecycle evidence. Attended clean Windows 11 acceptance remains #94, so #97 and #98 remain open.

## Test-first evidence

- Cleanup planner tests were observed failing before `Target`, `Result`, and `Run` existed.
- Cleanup command tests were observed failing before the fixed `wipe` verb and exit-state mapping existed.
- Installer lifecycle authoring initially produced 36 missing-contract diagnostics before the shortcut, UI, and removal authoring was added.
- Mutation tests were added for feature defaults and ownership, completion guards and sequence uniqueness, destructive confirmation, helper embedding, commit timing, upgrade exclusion, and invalid wipe values.
- Internal review later reproduced `TestRunRewritesTerminalStateAfterInitialResultFailure` failing with one write instead of the required running-plus-terminal retry before the implementation was corrected.
- Hosted run `33700662461` exposed a clean-runner empty GUI-process snapshot that the lifecycle harness rejected during parameter binding. `TestWindowsInstallerCIScriptAcceptsCleanGUIProcessSnapshot` was added and observed failing before the parameter accepted an empty collection, then passed after the correction.

## Focused local verification

| Check | Result |
| --- | --- |
| `go test -race ./internal/winuninstall ./cmd/gosched-cleanup` | PASS |
| `go vet ./internal/winuninstall ./cmd/gosched-cleanup` | PASS |
| `go test -race ./internal/winuninstall ./cmd/gosched-cleanup ./test/integration` | PASS |
| `pwsh -NoProfile -File build/windows/verify_wxs.ps1` | PASS |
| PowerShell parser for `verify_wxs.ps1`, `inspect-installer.ps1`, and `Invoke-InstallerContractCI.ps1` | PASS |
| UTF-8 without BOM and mojibake scan across 32 changed paths | PASS |
| `git diff --check` | PASS |

The local Windows build staged `goschedd.exe`, `gosched-gui.exe`, `gosched.exe`, and the windowless `gosched-cleanup.exe`, then linked a WiX 6.0.2 MSI with UI and Util 6.0.2. Compiled inspection passed for package identity, shortcut feature/component mappings, dialogs and controls, ordered completion events, secure removal property, cleanup Binary and custom action, execute/UI sequence boundaries, and daemon/GUI close rows.

- Local working-tree MSI SHA-256: `f3dcd159907ee6cb56f5f84048f4d6e40cec32d6409f7e5918ef537e5cc2b235`
- Artifact class: local candidate built from the final pre-commit working tree
- Artifact origin: local S039 verification only, not a published release asset

WiX ICE validation passes when suppressing the known `ICE38`, `ICE43`, and `ICE57` static findings for the context-dependent `ProgramMenuFolder` and `DesktopFolder` shortcut components, plus the existing intentional `ICE61` from `AllowSameVersionUpgrades`. The package is fixed per-machine (`ALLUSERS=1`), so Windows Installer resolves those standard directories to their all-users locations. No other schema, link, UI sequence, or ICE finding remains.

## Canonical repository verification

`C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` passed from the final post-review working tree:

1. format: PASS
2. vet: PASS
3. lint: PASS (`0 issues`)
4. race: PASS
5. GUI: PASS
6. coverage: PASS (`engine 86.4%`, `schedule 89.2%`, `timezone 91.3%`, `store 84.1%`, `catchup 88.9%`, `logbus 91.1%`)
7. documentation: PASS (14 pages plus product-policy fixtures)
8. automation: PASS

## Internal review closure

Five independent review perspectives covered project directives, shallow correctness, history and blame, prior pull-request decisions, and changed-code comments. Confirmed findings were corrected before the final gate run:

- an initial evidence-write failure now retries a terminal `internal-error` record and has a regression test;
- cleanup API contracts are documented and close/write/temporary-cleanup errors are no longer silently discarded;
- the registered-profile parent key remains open through relative SID-key enumeration;
- the destructive CI harness now requires GitHub's `RUNNER_ENVIRONMENT=github-hosted` discriminator and exits cleanly when `ShouldProcess` is declined;
- the Windows guide distinguishes unreleased S039 behavior from older MSI releases;
- uninstall completion copy distinguishes cleared successful evidence from retained refused/partial evidence;
- the Spec Kit feature pointer, managed both-shortcut example, and cleanup-result fallback count are corrected.

## Hosted and attended evidence disposition

The pull request's `windows-msi-contract` job is the authoritative hosted Windows Server silent run. It builds baseline and candidate MSIs, uploads compiled and runtime evidence, and exercises all four shortcut states, maintenance, same-authoring upgrade migration, invalid wipe values, repair/upgrade non-wipe controls, preserve removal, complete wipe, locked-file partial cleanup, out-of-scope sentinels, group preservation, and unattended non-launch behavior.

The retained development workstation was not used for machine-wide install/uninstall tests. The CI lifecycle script refuses any environment other than an elevated GitHub-hosted disposable Windows runner and launches `msiexec.exe` hidden and noninteractive.

S039 does not claim attended interaction or exact release-candidate proof. Issue #94 must still verify visible defaults, independent completion actions, default handlers, cancel/back behavior, multiple genuine profiles, interactive-user integrity, startup window bounds, and recurring-error behavior on a clean Windows 11 candidate. Issues #97 and #98 therefore remain open and are referenced, not closed, by the S039 pull request.
