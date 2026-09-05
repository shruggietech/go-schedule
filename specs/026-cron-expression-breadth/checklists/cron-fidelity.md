# Cron Fidelity Requirements Checklist: General Five-Field Cron Breadth

**Purpose**: Validate that S026 defines exact, reviewable cron semantics before implementation
**Created**: 2026-08-28
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are supported list, range, wildcard-step, range-step, and name forms defined for all five fields? [Completeness, Spec §FR-001]
- [x] CHK002 Are cross-field combination boundaries explicitly defined? [Completeness, Spec §FR-004-FR-005]
- [x] CHK003 Are explain, preview, import, authoring, editing, persistence, restart, catch-up, and export requirements all present? [Completeness, Spec §FR-007-FR-015]
- [x] CHK004 Are the remaining dialect, operational, and trigger exclusions enumerated rather than implied? [Completeness, Spec §Out of Scope]

## Requirement Clarity

- [x] CHK005 Is field-local step behavior distinguished clearly from elapsed-interval behavior? [Clarity, Spec §FR-002]
- [x] CHK006 Is the restricted day-of-month and day-of-week refusal tied explicitly to cron OR semantics? [Clarity, Spec §FR-005]
- [x] CHK007 Is the distinction between an authorable human phrase and a readable cron description explicit? [Clarity, Spec §FR-008-FR-009]
- [x] CHK008 Is canonical export defined as ordered, deduplicated, numeric, and recurrence-equivalent? [Clarity, Spec §FR-014]

## Requirement Consistency

- [x] CHK009 Do the broad-expression requirements preserve the no-silent-approximation rule across every interface? [Consistency, Spec §US2-US4]
- [x] CHK010 Are source identity, durable execution, and display-description requirements mutually consistent? [Consistency, Spec §FR-006-FR-010]
- [x] CHK011 Do missing-date, timezone, DST, anchor, restart, and catch-up requirements align across creation and export? [Consistency, Spec §FR-007, FR-015]
- [x] CHK012 Are existing simple phrases, calendar selectors, ordinal weekdays, and refusal contracts protected from regression? [Consistency, Spec §FR-009, FR-016]

## Acceptance Criteria Quality

- [x] CHK013 Can field-form breadth be measured against a concrete expression matrix? [Measurability, Spec §SC-001]
- [x] CHK014 Can end-to-end semantic parity be measured through every named interface and lifecycle stage? [Measurability, Spec §SC-002]
- [x] CHK015 Are calendar, DST, Sunday, overlap, and anchor boundaries included in a bounded parity period? [Measurability, Spec §SC-003]
- [x] CHK016 Are performance expectations quantified against both the p99 budget and a regression threshold? [Measurability, Spec §SC-006]

## Scenario and Edge-Case Coverage

- [x] CHK017 Are duplicates, overlapping ranges, mixed names and numbers, and Sunday aliases addressed? [Coverage, Spec §Edge Cases]
- [x] CHK018 Are absent dates and every missing-date policy addressed without changing imported cron defaults? [Coverage, Spec §Edge Cases]
- [x] CHK019 Are spring-forward gaps and fall-back overlaps addressed for multi-time schedules? [Coverage, Spec §Edge Cases]
- [x] CHK020 Are strictly-after-anchor and restart catch-up boundaries defined for matching creation minutes? [Coverage, Spec §Edge Cases]
- [x] CHK021 Are malformed syntax, recognizable unsupported syntax, and mutation rollback requirements all defined? [Coverage, Spec §US2-US4]

## Dependencies and Assumptions

- [x] CHK022 Is compatibility with the existing durable schedule representation stated and bounded? [Assumption, Spec §Assumptions]
- [x] CHK023 Are existing task timezone and imported policy defaults declared authoritative? [Assumption, Spec §Assumptions]
- [x] CHK024 Is the absence of new services, permissions, and dependencies explicit? [Dependency, Spec §FR-019]
