# Monthly Calendar Fidelity Requirements Checklist

**Purpose**: Validate that S025's monthly selector requirements are complete, precise, and safe against silent date movement **Created**: 2026-08-28 **Feature**: [spec.md](../spec.md)

## Selector Contract Completeness

- [x] CHK001 - Are day-of-month `L`, numbered `nW`, and `LW` each defined independently? [Completeness, Spec §FR-001–FR-003]
- [x] CHK002 - Is the valid numbered target range for `nW` explicit? [Clarity, Spec §FR-002]
- [x] CHK003 - Are required wildcard month and day-of-week fields stated? [Completeness, Spec §FR-004]
- [x] CHK004 - Are canonical human phrases defined for import and native authoring? [Consistency, Spec §FR-005–FR-006]
- [x] CHK005 - Are unsupported offsets, composite forms, lists, ranges, steps, and mixed modifiers bounded explicitly? [Coverage, Spec §FR-016]

## Calendar Semantics

- [x] CHK006 - Is nearest-weekday behavior specified for weekday, Saturday, and Sunday targets? [Completeness, Spec §FR-010]
- [x] CHK007 - Are first-day and final-day month-boundary reversals defined? [Edge Case, Spec User Story 3]
- [x] CHK008 - Is last-weekday behavior specified for every possible month-ending weekday? [Completeness, Spec §FR-011]
- [x] CHK009 - Is the ordering between missing-date resolution and weekday adjustment unambiguous? [Clarity, Spec §FR-012]
- [x] CHK010 - Is duplicate suppression required when next-valid resolution enters a later period? [Safety, Spec §FR-012]
- [x] CHK011 - Are leap years, short months, weekends, and daylight-saving boundaries measurable in acceptance coverage? [Acceptance Criteria, Spec §SC-002–SC-003]

## Fidelity and Workflow Coverage

- [x] CHK012 - Are explain, convert, crontab import, task input, desktop editing, and export all included? [Coverage, Spec §FR-005–FR-009]
- [x] CHK013 - Is original cron source retention distinct from generated display prose? [Clarity, Spec §FR-007–FR-008]
- [x] CHK014 - Are missing-date policies separated into policy-inert and export-restricted selector classes? [Consistency, Spec §FR-013–FR-014]
- [x] CHK015 - Is non-mutation required for both create and update failures? [Safety, Spec §FR-017]
- [x] CHK016 - Are existing `#`, day-of-week `L`, streams, and source-identity contracts protected? [Compatibility, Spec §FR-018]

## Documentation and Scope

- [x] CHK017 - Does the spec avoid claiming universal Quartz or POSIX compatibility? [Assumption, Spec Assumptions]
- [x] CHK018 - Are remaining issue #22 gaps distinguished from the selector family completed here? [Traceability, Spec §FR-020]
- [x] CHK019 - Is the substantial combined slice boundary explicit so these selectors are not fragmented into separate review cycles? [Scope, Spec Clarifications]

## Notes

- All 19 requirements-quality checks pass. No critical clarification remains before planning.
