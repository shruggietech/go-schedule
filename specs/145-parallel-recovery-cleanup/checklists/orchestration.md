# Orchestration Requirements Checklist: Correlated recovery and cleanup

**Purpose**: Author/reviewer requirements clarity before implementation.

**Created**: 2026-11-10

**Feature**: [spec.md](../spec.md)

- [x] CHK001 Are branch settlement and route selection distinguished from early blocked status? [Completeness, Spec FR-002, FR-003]
- [x] CHK002 Are cancellation request failure and actual native termination distinct? [Completeness, Spec Edge Cases, FR-002, FR-005]
- [x] CHK003 Are recovery cause defaults and explicit selections unambiguous? [Completeness, Spec FR-004]
- [x] CHK004 Is cleanup ordering defined when the follow-up is unselected, fails or cannot start? [Completeness, Spec FR-005]
- [x] CHK005 Are per-instance and per-stage duplicate/overlap identities specified? [Completeness, Spec FR-006, SC-002]
- [x] CHK006 Are restart uncertainty and external command effects described truthfully? [Completeness, Spec FR-007]
- [x] CHK007 Are admission, deletion and future-only edits consistent across every new route? [Completeness, Spec FR-008, FR-009]
- [x] CHK008 Are mixed graph replacement and portable task references complete? [Completeness, Spec FR-010, FR-012]
- [x] CHK009 Are full authoring text and bounded observations explicitly separated? [Completeness, Spec FR-014]
- [x] CHK010 Are all surfaces, target fencing, revisions, uncertainty and keyboard use covered? [Completeness, Spec FR-011, FR-013]
- [x] CHK011 Are unchanged joins and manual/ordinary chain behavior preserved? [Completeness, Spec FR-015]
- [x] CHK012 Are current-only state, no-focus verification and source/release boundaries explicit? [Completeness, Spec FR-016, FR-017, Assumptions]

All 12 requirements-quality items pass. No additional owner question is needed under the approved autopilot decision policy.
