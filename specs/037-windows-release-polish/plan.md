# Implementation Plan: Windows Release Polish

**Branch**: `codex/037-windows-release-polish` | **Date**: 2026-08-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/037-windows-release-polish/spec.md`

## Summary

Replace the full-work-area GUI startup request with a pure, table-tested logical sizing function capped at 90 percent of the selected Windows monitor's work area. Discover that monitor through pointer location plus `MonitorFromPoint` and `GetMonitorInfoW`, preserving fallback and normal window state controls. Revise the WiX Summary Subject, extend source and compiled-MSI contracts, correct specification 003, and record evidence without conflating native Explorer observation with static checks.

## Technical Context

**Language/Version**: Go 1.25; PowerShell 7; WiX Toolset 6

**Primary Dependencies**: Fyne v2; Windows user32 through the standard-library syscall bridge; Windows Installer COM

**Storage**: N/A

**Testing**: Go table tests and headless GUI tests; PowerShell source validation; compiled candidate-MSI inspection; canonical eight-gate aggregate

**Target Platform**: Cross-platform desktop GUI with Windows-specific monitor discovery and MSI packaging

**Project Type**: Desktop application and release packaging

**Performance Goals**: Constant-time startup calculation with no new dependency or persistent state

**Constraints**: 1280x800 preferred logical size, independent 90 percent caps, hard small-display bound, exact MSI PID 3 value, no release workflow execution

**Scale/Scope**: One startup window, six top-level tabs, one Windows amd64 MSI metadata field, two GitHub issues

## Constitution Check

| Principle | Pre-design | Post-design | Evidence |
| --- | --- | --- | --- |
| I. Code Quality | PASS | PASS | Pure sizing logic, narrow Windows adapter, explicit API failure fallback, no new dependency. |
| II. Testing Standards | PASS | PASS | Regression tests are written before implementation; eight gates remain mandatory. |
| III. User Experience Consistency | PASS | PASS | Restored startup is bounded while user-controlled window states remain unchanged. |
| IV. Performance Requirements | PASS | PASS | No scheduler hot path changes; monitor lookup and sizing are constant time. |
| V. Autonomous Build-Phase Execution | PASS | PASS | S037 uses the full Spec Kit sequence, review branch, pre-publication breakdown, and one PR. |

No constitution deviation is required.

## Decisions

### 2026-08-31: Select the launch monitor from the pointer

- **Decision**: Query the monitor nearest the current pointer before sizing.
- **Rationale**: No stable native window monitor exists before first show. The pointer identifies the user's active monitor and permits a deterministic constrained-secondary-monitor contract.
- **Alternatives considered**: Keep primary-monitor `SystemParametersInfoW` (fails issue #89), inspect the Fyne driver for private native handles (couples to internals), or persist prior placement (out of scope).

### 2026-08-31: Apply the 90 percent cap in logical units with no minimum enlargement

- **Decision**: Convert each physical work-area dimension by the effective scale, multiply by 0.9, and take the smaller of that result and the preferred dimension.
- **Rationale**: This directly implements the pinned issue formula and makes the work area a hard upper bound.
- **Alternatives considered**: Subtract fixed chrome margins (not DPI-stable), retain 640x480 clamps (can exceed the display), or use OS maximize state (repeats the reported problem).

### 2026-08-31: Verify MSI Subject at source and compiled-artifact layers

- **Decision**: Parse the WiX Subject in the source verifier and read PID 3 from the compiled MSI Summary Information stream in the artifact inspector.
- **Rationale**: Source validation catches review-time drift; compiled inspection proves the shipped database metadata. Explorer observation remains separate native evidence.
- **Alternatives considered**: Source regex only (cannot prove the artifact) or changing the release workflow to publish new evidence automatically (unnecessary release-pipeline scope).

## Project Structure

```text
gui/{app.go,screen.go,screen_test.go,screen_windows.go,screen_windows_test.go}
build/windows/{goschedule.wxs,verify_wxs.ps1}
test/windows/{inspect-installer.ps1,README.md}
specs/003-gui-editor-refinements/{spec.md,verification.md}
specs/037-windows-release-polish/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
└── tasks.md
CHANGELOG.md
specs/README.md
CLAUDE.md
```

**Structure Decision**: Extend the existing GUI startup and Windows installer verification seams. Do not add a new package, dependency, placement store, or release workflow.

## Complexity Tracking

No constitution violations or additional complexity require justification.
