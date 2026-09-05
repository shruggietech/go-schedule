# Cron Fidelity Requirements Checklist

**Purpose**: Validate that S023's ordinal-weekday requirements are complete, precise, and safe against silent timing changes
**Created**: 2026-08-28
**Feature**: [spec.md](../spec.md)

## Syntax Boundary Completeness

- [x] CHK001 - Is the supported `weekday#ordinal` shape limited to exactly one day-of-week term? [Completeness, Spec §FR-001]
- [x] CHK002 - Are numeric and named weekday inputs, including both Sunday numbers, explicitly covered? [Coverage, Spec §FR-002]
- [x] CHK003 - Is the accepted ordinal range stated exactly rather than described vaguely? [Clarity, Spec §FR-003]
- [x] CHK004 - Are day-of-month and month restrictions defined for the supported subset? [Completeness, Spec §FR-004]
- [x] CHK005 - Are lists, ranges, steps, multiple terms, and non-day-of-week placement explicitly addressed? [Edge Case, Spec §FR-012]

## Fidelity and Round-Trip Semantics

- [x] CHK006 - Does the spec require exact preservation of time, weekday, ordinal, and monthly cadence? [Consistency, Spec §FR-005–FR-008]
- [x] CHK007 - Is Sunday canonicalization defined without changing its schedule meaning? [Clarity, Spec §FR-008]
- [x] CHK008 - Are missing fifth occurrences and cron-compatible skip behavior defined? [Edge Case, Spec §FR-009]
- [x] CHK009 - Does the spec distinguish fifth-weekday policy effects from first-through-fourth cases where every month has an occurrence? [Consistency, Spec §FR-009–FR-010]
- [x] CHK010 - Are DST, month-boundary, and missing-occurrence validation windows measurable? [Acceptance Criteria, Spec §FR-011 and §SC-002–SC-003]

## Workflow and Failure Coverage

- [x] CHK011 - Are explain, convert, crontab preview/import, task authoring, and export workflows all covered? [Coverage, Spec §FR-005–FR-007 and §FR-015]
- [x] CHK012 - Is original cron source retention specified independently from generated display prose? [Clarity, Spec §FR-006]
- [x] CHK013 - Are malformed syntax and recognizable-but-unrepresentable combinations distinguished as errors or named refusals? [Clarity, Spec §FR-012]
- [x] CHK014 - Is non-mutation required for every failed or refused task boundary? [Safety, Spec §FR-013]
- [x] CHK015 - Are existing supported inputs, refusal categories, and CLI stream conventions protected? [Compatibility, Spec §FR-014–FR-015]

## Scope and Documentation

- [x] CHK016 - Does the spec identify `#` as an explicit subset extension without implying universal POSIX portability? [Assumption, Spec Assumptions]
- [x] CHK017 - Are `L`, `W`, last-weekday, Quartz, arbitrary combinations, and new recurrence modeling explicitly outside S023? [Scope, Spec Out of Scope]
- [x] CHK018 - Does the delivery requirement reference issue #22 without incorrectly closing the epic? [Traceability, Spec §FR-017]

## Notes

- All 18 requirements-quality checks pass. No critical clarification remains before planning.
