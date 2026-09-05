# Requirements Quality Checklist: GUI Dual-Syntax Contract

**Purpose**: Validate that S020's GUI input, preview, editing, error, compatibility, and discoverability requirements are complete enough for implementation and PR review **Created**: 2026-08-28 **Feature**: [spec.md](../spec.md) **Audience**: Author and third-party PR reviewers **Depth**: Standard correctness and UX gate

## Requirement Completeness

- [x] CHK001 Are create, preview, prefill, unchanged save, replacement, and help journeys all enumerated? [Completeness, Spec §User Stories]
- [x] CHK002 Are both accepted syntax forms covered at every recurring editor boundary? [Completeness, Spec §FR-001-008]
- [x] CHK003 Are new-task, existing-task, imported-task, one-off, and expressionless-legacy states distinguished? [Completeness, Spec §US1-US2, Edge Cases]
- [x] CHK004 Are local validation, daemon preview, and final submission responsibilities separated? [Completeness, Spec §FR-002, FR-005-006, Assumptions]
- [x] CHK005 Are GUI-local guidance changes and product-wide documentation work explicitly separated? [Completeness, Spec §FR-012, Out of Scope]

## Requirement Clarity

- [x] CHK006 Is the single-field decision explicit enough to rule out a selector, builder, or second cron field? [Clarity, Spec §Clarifications, Out of Scope]
- [x] CHK007 Is source identity defined as following current field contents rather than stale prefill metadata? [Clarity, Spec §FR-008]
- [x] CHK008 Is "supported cron" tied to the S019 faithful boundary instead of an undefined dialect? [Clarity, Spec §FR-001-002, Assumptions]
- [x] CHK009 Is the normalized-expression retention rule clear for leading and trailing whitespace? [Clarity, Spec §Edge Cases]
- [x] CHK010 Is the degraded legacy blank-field behavior distinguished from an invalid empty new-task field? [Clarity, Spec §US2, Edge Cases]

## Requirement Consistency

- [x] CHK011 Do preview and create/update require the same expression and source identity? [Consistency, Spec §FR-005-006]
- [x] CHK012 Do invalid-cron requirements agree with the central no-fallback rule? [Consistency, Spec §FR-002, FR-004]
- [x] CHK013 Do prefill and syntax-switch requirements agree on when response metadata stops governing? [Consistency, Spec §FR-007-008]
- [x] CHK014 Do human compatibility requirements preserve existing start-at, timezone, and policy behavior? [Consistency, Spec §FR-009]
- [x] CHK015 Do scope and traceability requirements consistently leave #50/#52 open? [Consistency, Spec §FR-016, Traceability]

## Acceptance Criteria Quality

- [x] CHK016 Can successful cron creation be objectively evaluated from preview and request identity? [Measurability, Spec §US1, SC-001]
- [x] CHK017 Can round-trip editing be evaluated from the exact retained expression before and after save? [Measurability, Spec §US2, SC-004]
- [x] CHK018 Can invalid/refused input be evaluated through save eligibility, named reason, and absence of mutation? [Measurability, Spec §US1, SC-003]
- [x] CHK019 Can compatibility be evaluated across the named existing fixture classes? [Measurability, Spec §SC-005]
- [x] CHK020 Is repository completion tied to the exact eight-gate verification contract? [Measurability, Spec §SC-006]

## Scenario and Edge-Case Coverage

- [x] CHK021 Are automatic cron, explicit human, invalid cron-shaped, unsupported cron, and five-word-human cases addressed? [Coverage, Spec §Edge Cases]
- [x] CHK022 Is switching syntax more than once in the same editing session addressed? [Coverage, Spec §US2, Edge Cases]
- [x] CHK023 Are preview transport failure and local validation outcomes independently specified? [Coverage, Spec §Edge Cases]
- [x] CHK024 Are one-off requests explicitly prevented from carrying recurring syntax identity? [Coverage, Spec §FR-010]
- [x] CHK025 Are timezone, DST, month-boundary, and missing-date parity tied to one central timing model? [Coverage, Spec §Edge Cases, SC-002]

## UX, Accessibility, Security, and Scope

- [x] CHK026 Are actionable visible-error requirements present without requiring an editor redesign? [UX, Spec §FR-011, Assumptions]
- [x] CHK027 Does the discoverability contract keep human schedules prominent while adding copy/pasteable cron guidance? [UX, Spec §US3, FR-012]
- [x] CHK028 Are existing focus, keyboard, and form interaction patterns preserved by the no-redesign scope? [Accessibility, Assumption, Out of Scope]
- [x] CHK029 Is the absence of API, storage, engine, IPC, authorization, daemon, and execution changes explicit? [Security, Scope, Spec §FR-014]
- [x] CHK030 Are dialect expansion, broad documentation, and issue closure deliberately deferred with traceability? [Scope, Spec §Out of Scope, Traceability]

## Notes

- Validated against the clarified specification on 2026-08-28: 30/30 items pass. No unresolved requirement-quality gap remains for the selected scope.
