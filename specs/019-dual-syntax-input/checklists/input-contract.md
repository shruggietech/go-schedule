# Requirements Quality Checklist: Dual-Syntax Input Contract

**Purpose**: Test whether S019's input, fidelity, compatibility, persistence, and import requirements are complete enough for implementation and PR review
**Created**: 2026-08-28
**Feature**: [spec.md](../spec.md)
**Audience**: Author and third-party PR reviewers
**Depth**: Standard correctness gate

## Requirement Completeness

- [x] CHK001 Are all non-GUI authoring surfaces that must share the central input boundary enumerated? [Completeness, Spec §FR-001]
- [x] CHK002 Are automatic and explicitly hinted classification paths both specified? [Completeness, Spec §FR-002-004]
- [x] CHK003 Are source retention and authoritative execution data distinguished for every recurring input? [Completeness, Spec §FR-005-006]
- [x] CHK004 Are preview, create, update, read, CLI, and import outcomes covered individually? [Coverage, Spec §FR-008-013]
- [x] CHK005 Are one-off and legacy expressionless schedules included in the response-identity contract? [Completeness, Spec §FR-008]
- [x] CHK006 Are migration and backfill expectations explicitly bounded? [Completeness, Spec §FR-014]

## Requirement Clarity

- [x] CHK007 Is the cron-shaped structural detector stated precisely enough to distinguish existing five-word human phrases? [Clarity, Spec §FR-002]
- [x] CHK008 Are the only valid hint values, omission behavior, and invalid-hint result unambiguous? [Clarity, Spec §FR-003-004]
- [x] CHK009 Is “retain the submitted expression” bounded by a stated whitespace-normalization rule? [Clarity, Spec §FR-006, Edge Cases]
- [x] CHK010 Is “specific fidelity reason” tied to the existing converter rather than an undefined new error taxonomy? [Clarity, Spec §FR-007]
- [x] CHK011 Is source syntax identity defined for recurring, one-off, legacy, and unchanged-update responses? [Clarity, Spec §FR-008-009]

## Requirement Consistency

- [x] CHK012 Does explicit hint precedence agree with the no-fallback rule? [Consistency, Spec §FR-003-004]
- [x] CHK013 Do source retention requirements consistently keep raw cron outside the execution path? [Consistency, Spec §FR-005-006]
- [x] CHK014 Do CLI thin-client requirements agree with central API validation? [Consistency, Spec §FR-001, FR-011]
- [x] CHK015 Do import retention requirements preserve existing per-line preview and partial-success behavior? [Consistency, Spec §FR-012]
- [x] CHK016 Do compatibility requirements avoid contradicting the intentionally superseded human-only policy? [Consistency, Spec §FR-010, FR-015]

## Acceptance Criteria Quality

- [x] CHK017 Can preview/create equivalence be objectively evaluated from recurrence and upcoming-run outputs? [Measurability, Spec §US1, SC-001-002]
- [x] CHK018 Can no-fallback failure behavior be verified through output, mutation, and reason criteria? [Measurability, Spec §US1, SC-003]
- [x] CHK019 Can human compatibility be compared against pre-feature expression, RRULE, anchor, summary, and timing evidence? [Measurability, Spec §FR-010, SC-004]
- [x] CHK020 Can import retention be verified from the created task's expression and syntax identity? [Measurability, Spec §US3, SC-005]

## Scenario and Edge-Case Coverage

- [x] CHK021 Are supported shorthand, standard expression, invalid cron-shaped, lossy, and five-word-human examples addressed? [Coverage, Spec §Edge Cases]
- [x] CHK022 Are field-local steps and combined DOM/DOW behavior deliberately bounded to the existing faithful dialect? [Coverage, Spec §Edge Cases]
- [x] CHK023 Are DST, month-boundary, timezone, and missing-date parity requirements present without creating a second timing model? [Coverage, Spec §FR-013]
- [x] CHK024 Is the behavior of a schedule update that omits the schedule field specified? [Coverage, Spec §FR-009]
- [x] CHK025 Are invalid hints, empty expressions, and expressionless legacy rows covered by deterministic outcomes? [Coverage, Spec §FR-003, FR-008, Edge Cases]

## Dependencies, Security, and Scope

- [x] CHK026 Is S018 explicitly identified as the detector and fidelity dependency? [Dependency, Spec §Assumptions]
- [x] CHK027 Are GUI adoption, broad posture documentation, dialect expansion, and raw-cron execution excluded? [Scope, Spec §Out of Scope]
- [x] CHK028 Is the absence of an IPC, authorization, secret, daemon, or execution-boundary change explicit? [Security, Spec §FR-016]
- [x] CHK029 Is issue closure behavior explicit so this partial slice cannot accidentally close epic #50? [Traceability, Spec §Traceability]
- [x] CHK030 Is the GUI's temporary human-only request behavior specified without claiming GUI cron adoption? [Scope, Spec §FR-017]

## Notes

- Validated against the clarified specification on 2026-08-28: 30/30 items pass. No missing or ambiguous requirement remains for the selected scope.
