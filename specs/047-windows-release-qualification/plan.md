# Implementation Plan: Windows Release Qualification

**Branch**: `codex/047-windows-release-qualification` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/047-windows-release-qualification/spec.md`

## Summary

Extend the existing cross-platform Windows release validator with eleven fixed, issue-traceable desktop observations and scenario-specific metric validation. Update the resumable PowerShell collector to generate matching templates, expand the synthetic fixture and mutation tests, and document one consolidated native matrix. After implementation and canonical verification, commit the source and build one `local-demo` MSI with an embedded S047 identity for the pre-push attended walkthrough. Formal release evidence remains reserved for the exact workflow-staged MSI after review, merge, and authorized tagging.

## Technical Context

**Language/Version**: Go 1.25.0 with toolchain 1.25.1; PowerShell 7.6-compatible scripts; WiX Toolset 6.0.2 for local packaging

**Primary Dependencies**: Go standard library; existing Fyne desktop runtime; Windows Installer COM inspection; existing WiX extensions

**Storage**: Versioned JSON evidence workspace and ZIP archive; repository fixture files; local untracked `dist/` artifacts

**Testing**: Table-driven and mutation-based Go tests, PowerShell parser checks, Windows integration contract tests, compiled MSI inspection, `scripts/verify.sh all`

**Target Platform**: Validator on Linux and Windows; collector and native acceptance on Windows 11 client; Windows amd64 MSI

**Project Type**: Desktop application plus daemon/CLI, installer, release tooling, and evidence validator

**Performance Goals**: Evidence validation remains bounded by the existing 16 MiB JSON, 512-entry archive, 128 MiB per-entry, and 512 MiB total limits; no scheduler hot path changes

**Constraints**: UTF-8 without BOM; no new dependency; no visible or interactive child process from automation; no destructive lifecycle automation on the non-disposable development host; no formal native claim from fixtures/headless evidence; no push, PR, tag, workflow dispatch, or release

**Scale/Scope**: Preserve 36 existing scenarios and add 11 desktop scenarios (47 total); eight linked desktop issues plus #98 and coordinator #96; one local demo MSI and inspection report

## Constitution Check

*GATE: Passed before research and rechecked after design.*

| Principle | Requirement | Design disposition |
| --- | --- | --- |
| I. Code Quality | Small, explicit, idiomatic validation with contextual errors | PASS. One `desktop.` dispatcher and focused validators reuse typed metric helpers; no dependency or concurrency is added. |
| II. Testing Standards | Behavioral changes have regression tests; race and coverage gates remain | PASS. Tests are written before validator/collector changes, mutate every metric family, preserve fixture-class separation, and run under the canonical aggregate. |
| III. UX Consistency | Evidence must prove actionable, consistent desktop behavior | PASS. Native requirements name themes, font behavior, controls, navigation, scrolling, tables, and readable semantic states with measurable floors. |
| IV. Performance | No unexplained hot-path regression | PASS. Release validation is offline and bounded by existing caps; scheduler/runtime paths are unchanged, so no new benchmark is warranted. |
| V. Autonomous Execution | Full Spec Kit order, analyze gate, review branch, one pre-push halt | PASS. S047 runs under explicit autopilot authority and stops with the local demo plus remaining attended work before publication. |

Pinned `test/windows/**`, release-contract documentation, fixture, and `CHANGELOG.md` surfaces change because native desktop evidence must become a promotion invariant. The dated decision records that expansion. No constitutional deviation is required.

## Technical Design

### Fixed observations

Append these identities to the existing canonical sequence:

1. `desktop.appearance-standard`
2. `desktop.appearance-scaled`
3. `desktop.interaction-states`
4. `desktop.interaction-states-scaled`
5. `desktop.navigation-options`
6. `desktop.navigation-options-scaled`
7. `desktop.scroll-input`
8. `desktop.tasks-table`
9. `desktop.tasks-table-scaled`
10. `desktop.schedule-activity-tables`
11. `desktop.schedule-activity-tables-scaled`

Every `desktop.` observation requires one intended-user, medium-integrity Windows 11 environment with the installed LocalSystem service and at least one raster-image attachment validated from its bytes. Every visual family has a 96-DPI observation and a distinct observation above 96 DPI; scroll input remains the standard-DPI physical-input scenario. The four table observations require at least 100 rows per applicable view. Palette and window-size coverage use exact string sets to prevent a single screenshot from being described as both configurations.

### Metric shape

Use flat, typed metrics rather than a second nested schema. Existing helper functions already emit precise missing/type/value diagnostics and strict JSON decoding rejects unknown record fields. Booleans represent reviewed outcomes; exact strings/sets represent palette, size, state, font, and event coverage; integers/floats represent row counts and contrast floors.

Comma-separated set metrics are normalized by trimming, rejecting empty or duplicate elements, sorting, and comparing with the required set. This keeps operator fragments readable while preventing order-sensitive false failures.

### Collector templates

Generalize the lifecycle-only template helper into a scenario template helper. It emits the established lifecycle fields unchanged and adds explicit desktop metric keys plus screenshot paths. `Initialize` creates templates for setup, remove, and desktop families, while window capture remains restricted to its five native geometry observations.

### Local demo

The source must be committed before packaging so all binaries can embed `1.0.0-s047-demo.<short-commit>`. The MSI keeps numeric ProductVersion `1.0.0`, uses a filename containing `s047-demo`, and is inspected with artifact class `local-demo`. `dist/` remains untracked. Rebuilding after any source change invalidates prior hash-bound attended results.

## Project Structure

### Documentation (this feature)

```text
specs/047-windows-release-qualification/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── checklists/
│   ├── requirements.md
│   ├── release-evidence.md
│   └── attended-demo.md
├── contracts/
│   └── desktop-release-evidence.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/releasegate/
├── model.go
├── validate.go
└── validate_test.go

test/windows/
├── Invoke-ReleaseCandidateAttended.ps1
└── README.md

test/fixtures/windows-release-gate/passing/
└── evidence.json

scripts/windows-release-gate/
├── main.go
└── main_test.go

docs/
└── INSTALL-windows.md

specs/README.md
CLAUDE.md
CHANGELOG.md
```

**Structure Decision**: Extend the established S040 release-gate package, collector, fixture, and Windows documentation in place. A separate validator or evidence schema would duplicate promotion semantics and create a bypass risk.

## Implementation Sequence

1. Pin all eleven scenario identities and failing validator/fixture tests.
2. Add metric-set parsing and the desktop scenario validators.
3. Add failing collector contract assertions, then generate desktop templates and update the canonical fixture.
4. Update the Windows runbook, issue mapping, and attended demo checklist.
5. Run focused validation, parser, integration, and full canonical gates.
6. Commit the verified source, then build and inspect the hash-bound local demo.
7. Halt for the attended checklist before any push.

## Post-Design Constitution Recheck

All five principles remain PASS. The plan preserves existing security, IPC, installer, scheduler, race, coverage, and promotion controls. It introduces no new process launcher, privilege, dependency, persistence format used by the product, or runtime performance surface. Native observations remain human- attended and machine-validated instead of being simulated.

## Complexity Tracking

No constitution violation requires justification.
