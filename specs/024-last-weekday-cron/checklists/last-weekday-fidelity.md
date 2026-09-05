# Last-Weekday Fidelity Requirements Checklist

**Purpose**: Validate that S024's last-weekday requirements are complete, precise, and safe against silent timing changes **Created**: 2026-08-28 **Feature**: [spec.md](../spec.md)

## Syntax Boundary Completeness

- [x] CHK001 - Is the supported `weekdayL` shape limited to exactly one day-of-week atom? [Completeness, Spec §FR-001]
- [x] CHK002 - Are numeric and named weekdays, including both Sunday numbers, explicitly covered? [Coverage, Spec §FR-002]
- [x] CHK003 - Are wildcard day-of-month and month requirements stated exactly? [Clarity, Spec §FR-003]
- [x] CHK004 - Are bare, malformed, multiple, list, range, step, and mixed extension forms addressed? [Edge Case, Spec §FR-010]
- [x] CHK005 - Is day-of-month `L` explicitly distinguished and declined? [Scope, Spec §FR-011]

## Fidelity and Round-Trip Semantics

- [x] CHK006 - Does the spec require exact preservation of time, weekday, and monthly cadence? [Consistency, Spec §FR-004–FR-009]
- [x] CHK007 - Is Sunday canonicalization defined without changing schedule meaning? [Clarity, Spec §FR-007]
- [x] CHK008 - Is missing-date policy correctly defined as inert for last weekdays? [Consistency, Spec §FR-008]
- [x] CHK009 - Are DST and month-boundary validation windows measurable? [Acceptance Criteria, Spec §FR-009 and §SC-002]

## Workflow and Failure Coverage

- [x] CHK010 - Are explain, convert, crontab preview/import, task authoring, and export workflows covered? [Coverage, Spec §FR-004–FR-006 and §FR-014]
- [x] CHK011 - Is original cron source retention specified independently from generated prose? [Clarity, Spec §FR-005]
- [x] CHK012 - Are malformed syntax and recognizable unsupported combinations distinguished as errors or named refusals? [Clarity, Spec §FR-010]
- [x] CHK013 - Is non-mutation required for failed or refused task boundaries? [Safety, Spec §FR-012]
- [x] CHK014 - Are existing accepted inputs, refusal categories, and stream conventions protected? [Compatibility, Spec §FR-013–FR-014]

## Scope and Documentation

- [x] CHK015 - Does the spec identify `weekdayL` as an explicit subset extension without implying universal portability? [Assumption, Spec Assumptions]
- [x] CHK016 - Are day-of-month `L`, `W`, Quartz forms, and broader combinations outside S024? [Scope, Spec Out of Scope]
- [x] CHK017 - Does delivery reference issue #22 without incorrectly closing the epic? [Traceability, Spec §FR-016]

## Notes

- All 17 requirements-quality checks pass. No critical clarification remains before implementation planning.
