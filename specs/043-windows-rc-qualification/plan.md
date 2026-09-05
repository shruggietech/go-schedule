# Implementation Plan: Windows Demo Qualification

**Branch**: `codex/043-windows-rc-qualification` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/043-windows-rc-qualification/spec.md`

## Summary

Create an evidence-first pre-PR qualification package for the merged Windows release-readiness work. The implementation adds no speculative product behavior. It defines a reproducible local demo build, executes existing automated gates, records the compiled artifact identity, and presents only the remaining attended steps to the operator. Formal #94 qualification remains bound to a later exact Release-workflow artifact.

## Technical Context

**Language/Version**: Go 1.25; PowerShell 7; WiX Toolset 6.0.2

**Primary Dependencies**: Fyne 2.8.1, modernc SQLite, Windows Installer, existing release-gate and installer inspection tools

**Storage**: Local `dist/` artifact and Markdown verification record; no new runtime storage

**Testing**: `scripts/verify.sh all`, focused Go/race tests, PowerShell parser and WiX source checks, compiled-MSI inspection

**Target Platform**: Windows 11 x64 for attended testing; repository gates remain cross-platform

**Project Type**: Go daemon, CLI, and Fyne desktop application with WiX MSI

**Performance Goals**: No runtime performance change; retain existing benchmark and dispatch-latency gates

**Constraints**: No destructive automated install on a non-disposable host; no network publication before attended demo completion; strict fail-closed evidence

**Scale/Scope**: One demo MSI, one condensed operator checklist, seven linked release-readiness issues, existing verification infrastructure

## Constitution Check

### Pre-research gate

- **I. Code Quality**: PASS. Existing build and inspection paths are reused; any correction requires idiomatic code and contextual errors.
- **II. Testing Standards**: PASS. Automated gates run before handoff, missing attended evidence remains non-pass, and any behavior change begins red.
- **III. UX Consistency**: PASS. The attended matrix covers installer copy, navigation, appearance, storage disclosure, and recovery behavior.
- **IV. Performance**: PASS. No hot-path change is planned; the canonical benchmark gate remains mandatory.
- **V. Autonomous Execution**: PASS. S043 is issue-traceable, uses a review branch, runs the complete Spec Kit sequence, and halts before publication. The user's stricter instruction additionally withholds the PR until attended testing ends.

### Post-design gate

PASS. The design adds no dependency, persistence format, privilege expansion, or release bypass. The two-tier demo/formal-candidate model prevents exploratory evidence from weakening #94.

## Research decisions

See [research.md](research.md). The key decision is to build a distinctly named local demo from the eventual S043 commit, not to create a provisional GitHub tag or draft release. A tag would violate both the user's no-publication instruction and the formal release authorization boundary.

## Project Structure

```text
specs/043-windows-rc-qualification/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── demo-handoff.md
├── checklists/
│   ├── requirements.md
│   └── attended-demo.md
├── tasks.md
└── verification.md

build/windows/
├── goschedule.wxs
└── verify_wxs.ps1

test/windows/
├── inspect-installer.ps1
├── Invoke-InstallerLifecycle.ps1
└── Invoke-ReleaseCandidateAttended.ps1

dist/
├── go-schedule_s043-demo_windows_amd64.msi
└── s043-demo-installer-inspection.md
```

**Structure Decision**: Retain the existing single Go repository and Windows/ release tooling. Compiled execution exposed that the inspector could only label a local demo as a formal candidate, so S043 adds one bounded `local-demo` provenance value without relaxing candidate-manifest or published-origin checks.

## Implementation phases

1. Produce and analyze the S043 evidence contract.
2. Establish a clean automated baseline on the merged S042 source.
3. Install the pinned WiX tool only if absent, then build all four Windows executables with the S043 demo version and compile one MSI.
4. Inspect the compiled MSI and run focused and canonical verification.
5. Record identity, outcomes, and unavailable native boundaries.
6. Commit the source/evidence package locally, rebuild from that exact commit if the commit changes artifact inputs, and hand off the final demo MSI.
7. After operator testing, remediate confirmed defects under proof-first rules or, if clean, prepare the still-unpushed branch for the separately authorized PR.

## Complexity Tracking

No constitutional violation or new architectural complexity is introduced.
