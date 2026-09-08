# Release Readiness Requirements Checklist: Wails Release Cutover

**Purpose**: Validate that the requirements define a complete, reviewable, and rollback-aware production desktop cutover
**Created**: 2026-09-07
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are required desktop platforms, architectures, package formats, and payloads specified? [Completeness, Spec FR-001, FR-002]
- [x] CHK002 Are stable executable, shortcut, bundle, and installer identity requirements documented? [Completeness, Spec FR-003, FR-004, FR-005]
- [x] CHK003 Are upgrade, repair, preservation, explicit wipe, and migration outcomes all defined? [Coverage, Spec User Story 2, FR-004, FR-018]
- [x] CHK004 Is the legacy capability inventory required to cover supported workflows and intentional retirements? [Completeness, Spec FR-009, FR-010]
- [x] CHK005 Are removal requirements defined for source, dependencies, build paths, tests, proof applications, and maintained guidance? [Completeness, Spec FR-011, FR-013, FR-017]

## Requirement Clarity

- [x] CHK006 Is the meaning of the stable gosched-gui launch contract unambiguous across platforms? [Clarity, Spec FR-003]
- [x] CHK007 Is the boundary between current guidance and legitimate historical or migration references clear? [Clarity, Spec FR-012]
- [x] CHK008 Is candidate identity defined strongly enough to reject evidence from another revision? [Clarity, Spec FR-016]
- [x] CHK009 Is the no-release boundary distinct from package and release-workflow preparation? [Clarity, Spec FR-019]

## Requirement Consistency

- [x] CHK010 Do Fyne removal requirements remain consistent with one-time preference migration? [Consistency, Spec FR-011, FR-018]
- [x] CHK011 Do native desktop package requirements preserve the independently distributed cgo-free daemon and CLI matrix? [Consistency, Spec FR-001, FR-014]
- [x] CHK012 Do issue completion requirements align with the explicit post-merge milestone ritual? [Consistency, Spec FR-019, FR-020]

## Acceptance Criteria Quality

- [x] CHK013 Can package completeness be measured per target without subjective judgment? [Measurability, Spec SC-001]
- [x] CHK014 Can parity be accepted only with zero unexplained inventory entries? [Measurability, Spec SC-002]
- [x] CHK015 Can verification distinguish passed, failed, and skipped gates for the exact candidate? [Measurability, Spec SC-003]
- [x] CHK016 Can current-product Fyne residue be measured independently from immutable history? [Measurability, Spec SC-005]

## Scenario and Edge-Case Coverage

- [x] CHK017 Are partial artifact, naming mismatch, platform failure, and stale-evidence scenarios addressed? [Coverage, Spec Edge Cases]
- [x] CHK018 Are running-process, same-version, downgrade, preserve, repair, and wipe paths addressed? [Coverage, Spec Edge Cases]
- [x] CHK019 Are malformed, missing, unreadable, and valid legacy preference states addressed? [Coverage, Spec Edge Cases]
- [x] CHK020 Are accessibility, zoom, offline behavior, and native platform evidence included in the release boundary? [Non-Functional, Spec FR-007, SC-006]

## Dependencies and Assumptions

- [x] CHK021 Are parent and child issue dependencies explicitly recorded? [Dependency, Spec Assumptions]
- [x] CHK022 Are supported desktop and server-only architecture assumptions explicit? [Assumption, Spec Assumptions]
- [x] CHK023 Is continued daemon and CLI contract authority explicit and consistent with the cutover? [Assumption, Spec Assumptions]
- [x] CHK024 Is the limitation of pull-request evidence versus post-merge public release activity explicit? [Assumption, Spec Assumptions]

## Notes

- Standard depth for release reviewers, with emphasis on packaging completeness, rollback-safe upgrades, exact-candidate evidence, and removal boundaries.
- All items passed after the specification review; implementation verification is tracked separately in quickstart.md and verification.md.
