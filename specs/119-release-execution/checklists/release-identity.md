# Release identity requirements checklist: S119

**Purpose**: Evaluate completeness of release-state and authority requirements.

**Created**: 2026-10-26

**Feature**: [spec.md](../spec.md)

## Source and candidate identity

- [x] CHK001 Exact approved tag/source is named. [FR-001]
- [x] CHK002 Main acceptance and later branch identity are distinguished. [FR-001, FR-009]
- [x] CHK003 Tag/draft conflicts have an explicit non-overwrite outcome. [FR-003]
- [x] CHK004 Actual successful staging run/attempt is required. [FR-002]
- [x] CHK005 Complete platform asset identities are required. [FR-002]
- [x] CHK006 Original Windows manifest and candidate bytes remain bound. [FR-002]

## Publication and discovery

- [x] CHK007 Draft and public release are distinct states. [FR-004, FR-005]
- [x] CHK008 Promotion approval requires concrete staged identities. [FR-004]
- [x] CHK009 Existing waiver and checksum contracts are retained. [FR-004]
- [x] CHK010 Interrupted or already-published recovery cannot replace content. [FR-003, FR-005]
- [x] CHK011 Public latest/download/checksum identity is complete. [FR-005]
- [x] CHK012 Source discovery changes follow actual publication. [FR-006]

## Defects and delivery boundaries

- [x] CHK013 Original #355/#356 causes and evidence remain independent. [FR-007]
- [x] CHK014 New defects have attributable separate records. [FR-008]
- [x] CHK015 Spec-Kit/blocking analysis and canonical acceptance remain required. [FR-009]
- [x] CHK016 Source review authority stops at unmerged PR and at most two Codex rounds. [FR-010]

## Notes

16/16 requirements-quality items pass. No new approval policy, native qualification prerequisite or speculative runtime change is introduced.
