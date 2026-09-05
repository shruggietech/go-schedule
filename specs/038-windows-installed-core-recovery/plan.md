# Implementation Plan: Windows Installed Core Recovery

**Branch**: `codex/038-windows-installed-core-recovery` | **Date**: 2026-09-01 | **Spec**: [spec.md](spec.md)

## Summary

Restore the installed Windows core path by expanding the restricted named-pipe descriptor with validated direct user member SIDs from the configured local group, then extend the native installer lifecycle probe through real task creation, manual execution, scheduled execution, and controlled failures. Improve executor start-failure diagnostics without logging arguments or environment values. Preserve LocalSystem service isolation, windowless children, the configured-group ACE, and the S036 in-frame recovery UI.

## Technical Context

**Language/Version**: Go 1.25, PowerShell 7, Markdown

**Primary Dependencies**: Go standard library, `github.com/Microsoft/go-winio`, `golang.org/x/sys/windows`, Netapi32 Windows API, existing Cobra CLI and HTTP/JSON API

**Storage**: Existing SQLite task/run store; native verification writes Markdown evidence and deterministic marker files

**Testing**: Go unit tests under Windows build tags, executor tests, existing scheduler integration tests, a disposable Windows CI LocalSystem probe, native installer lifecycle probe, canonical eight-gate suite

**Target Platform**: Windows 11 desktop client with LocalSystem Windows service; cross-platform regression gates remain mandatory

**Project Type**: Go desktop application, daemon/service, CLI, and Windows MSI

**Performance Goals**: Resolve one bounded local-group membership snapshot at daemon startup; no per-connection group enumeration and no scheduling hot-path regression

**Constraints**: Fail closed; no Authenticated Users/Everyone ACE, permanent elevation, UAC disablement, interactive service, visible console window, secret-bearing diagnostics, fake-runner acceptance, or new external dependency

**Scale/Scope**: One configured local alias, its direct member set, one native access probe, and four execution outcomes (manual success, scheduled success, nonzero exit, process-start failure)

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Membership enumeration and descriptor construction are separated behind injected operations; executor diagnostics remain bounded and secret-safe. |
| II. Testing Standards | PASS | Tests are written first for descriptor policy and execution diagnostics; native proof uses real installed binaries and service wiring. |
| III. UX Consistency | PASS | Existing in-frame recovery remains; failed history differentiates a child exit from failure to start without exposing command arguments. |
| IV. Performance | PASS | Netapi32 is called once at listener startup; member SIDs are sorted and de-duplicated once. No dispatch-loop changes are planned. |
| V. Autonomous Execution | PASS | Full Spec Kit sequence, read-only analyze gate, eight canonical gates, review branch, and authorized PR/review loop are required. |

### Post-design re-check

All gates remain PASS. `docs/INSTALL-windows.md`, `test/windows/Invoke-InstallerLifecycle.ps1`, and `CLAUDE.md` are pinned surfaces and change because the issues require a corrected installed authorization contract, a real service execution walkthrough, and Spec Kit's enabled context hook. The dated decision is recorded in `CHANGELOG.md`. The plan intentionally deviates from S036's ACL exclusion because issue #90 proves that presentation-only recovery did not restore the product outcome.

## Architecture and Decision Log

### Expand verified direct user members at listener startup

Use `NetLocalGroupGetMembers` information level 1 to retrieve each direct member's SID and account type. Add only `SidTypeUser` SIDs to the DACL, retaining the local-group ACE for fresh tokens and group-valued members. This handles the verified state where the installing user's stable user SID is present but a newly added local alias SID is absent from the standard token. Granting Authenticated Users, requiring elevation, and silently falling back to compatibility mode were rejected because they weaken or bypass the established authorization boundary.

### Build the descriptor from typed validated principals

Resolve the configured group, enumerate member records, validate every SID with Windows parsing, de-duplicate case-insensitively, sort member SID strings, and render one ACE per principal. The listener does not log member names or SIDs. Enumeration failure aborts startup; a record with no usable SID is skipped only when Windows identifies it as deleted or unknown. This keeps SDDL construction deterministic and fail-closed.

### Prove execution through the installed service boundary

Extend `Invoke-InstallerLifecycle.ps1` with an execution-probe scenario that requires an elevated installed lifecycle context only for setup, but calls the installed CLI and daemon for all product operations. It creates deterministic tasks using absolute inbox Windows PowerShell, direct .NET marker-file I/O, and explicit noninteractive arguments. It captures manual and scheduled run JSON, marker hashes/content, service account, and correlated logs. The existing fake integration runner remains useful for scheduler timing but is not accepted as Windows execution proof.

Add a narrower disposable Windows CI probe that builds the daemon and CLI, creates the authorization alias, registers the daemon as LocalSystem, and executes the same success and failure controls. This automates the service boundary without claiming ordinary-token IPC, candidate-MSI, or release-gate coverage.

### Distinguish process-start failure without leaking inputs

When `cmd.Run` fails before an exit code exists and captured output is empty, retain a stable `process start failed for <quoted executable>: <wrapped OS error>` message. Do not include arguments, stdin, working environment values, or shell-expanded text. A started child that exits nonzero retains its exit code and captured output unchanged.

### Keep the broader release gate separate

S038 supplies reusable native access and task-execution probes. S039/#94 will decide how candidate artifacts are provisioned and gated in release automation across the full installed workflow. S038 does not claim that release gate is complete.

## Project Structure

```text
internal/ipc/ipc_windows.go
internal/ipc/ipc_windows_test.go
internal/executor/executor.go
internal/executor/executor_test.go
test/windows/Invoke-InstallerLifecycle.ps1
test/windows/Invoke-ServiceCoreCI.ps1
test/windows/README.md
.github/workflows/ci.yml
docs/INSTALL-windows.md
specs/036-ipc-access-denied-recovery/spec.md
specs/038-windows-installed-core-recovery/
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend existing IPC, executor, and Windows lifecycle boundaries. No new runtime package, database migration, public API route, or dependency is warranted.

## Complexity Tracking

No constitution violation requires justification.
