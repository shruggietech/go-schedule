# Requirements Checklist: DST Policy Contract

**Purpose**: Unit-test the policy requirements for completeness, clarity, consistency, and measurable coverage **Created**: 2026-08-29

## Requirement Completeness

- [x] CHK001 Are all basis values, defaults, applicable schedule shapes, and refusal behavior specified? [Completeness, Spec FR-001 to FR-004]
- [x] CHK002 Are both transition classes, all policy values, defaults, and inert-basis behavior specified? [Completeness, Spec FR-005 to FR-008]
- [x] CHK003 Are persistence, migration, restart, and schedule-replacement requirements documented? [Completeness, Spec FR-016 to FR-019]

## Requirement Clarity and Consistency

- [x] CHK004 Is the distinction among wall-clock, elapsed, and UTC objectively defined without inferring intent from frequency? [Clarity, Spec FR-002 to FR-004]
- [x] CHK005 Are transition and process-overlap policies explicitly distinguished? [Consistency, Spec Out of Scope]
- [x] CHK006 Is ordering between missing-date, calendar adjustment, and DST resolution unambiguous? [Consistency, Spec FR-010]
- [x] CHK007 Is the corrected classification of five/seven-hour wall-clock gaps consistent across requirements and documentation outcomes? [Conflict, Spec Clarifications]

## Scenario and Edge Coverage

- [x] CHK008 Are first, both, and last outcomes covered with exact ordered-instant expectations, including a cursor between folds? [Coverage, Spec US2 and Edge Cases]
- [x] CHK009 Are gap collision, non-one-hour transitions, duplicate suppression, and inert one-off/event behavior addressed? [Coverage, Spec Edge Cases]
- [x] CHK010 Are preview, create, edit, detail, calendar, catch-up, restart, and dispatch required to agree? [Coverage, Spec FR-011 to FR-018]
- [x] CHK011 Are invalid enum and incompatible elapsed failures required to be field-specific and non-mutating? [Exception Flow, Spec US4]

## Acceptance Criteria Quality

- [x] CHK012 Can the transition matrix, lifecycle parity, migration result, and benchmark threshold be measured objectively? [Measurability, Spec SC-001 to SC-006]
- [x] CHK013 Are compatibility expectations measurable for both stored tasks and omitted in-memory values? [Acceptance Criteria, Spec FR-016 and FR-019]

## Dependencies and Assumptions

- [x] CHK014 Are task timezone presentation semantics under elapsed and UTC explicitly documented? [Assumption, Spec Assumptions]
- [x] CHK015 Are dependency, permission, service, and scope boundaries explicit? [Dependency, Spec FR-021 and Out of Scope]
