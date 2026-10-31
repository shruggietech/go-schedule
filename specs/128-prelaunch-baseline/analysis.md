# S128 Specification Analysis

**Date**: 2026-10-31
**Method**: Installed /speckit-analyze, read-only comparison after successful prerequisite resolution.
**Result**: PASS before implementation.

## Findings and resolution

| ID | Severity | Finding | Resolution |
|---|---|---|---|
| C1 | CRITICAL, resolved | Constitution3.1.0 mandates historical migration compatibility, conflicting with the owner's latest explicit no-users/no-migrations requirement. | The owner approved S128's policy alignment. T012 formally amends the constitution to 4.0.0 and synchronizes current guidance. No security/current-state correctness or verification requirement is waived. |
| S1 | HIGH, resolved | Metrics adapters and full qualification tools also enforce historical conversions/upgrade matrices; limiting work to application schema and Fyne imports would leave FR-005/FR-006 incomplete. | Research mapped these active surfaces; T007/T013/T014 include direct current metrics4, fresh qualification and current evidence2. Historical packets remain immutable. |
| S2 | MEDIUM, resolved | Historical migration-named files contain current behavior tests that must not be deleted wholesale. | T006 explicitly ports/renames current identity, query-index, notification, epoch and portable-data assertions. |
| S3 | MEDIUM, resolved | Arbitrary metrics database paths currently permit a shared file; strict separate-version validation would add an unrequested restriction. | The complete current metrics4 layout serves either writer and shared paths, without conversion of old formats. |

No unresolved critical/high/medium/low finding remains. No scope amendment or further owner decision is required beyond the explicit approved proposal.

## Requirement coverage

| Requirement | Task coverage |
|---|---|
| FR-001 | T004, T005, T007 |
| FR-002 | T004, T005, T006 |
| FR-003 | T004, T005, T007 |
| FR-004 | T008, T009, T010 |
| FR-005 | T012, T013, T014, T015 |
| FR-006 | T006, T007, T010, T011, T013, T014, T016 |
| FR-007 | T005, T007, T015 |
| FR-008 | T002, T015, T016, T017 |
| FR-009 | T012, T015, T018 |
| FR-010 | T001, T002, T003, T012, T018 |
| SC-001 | T004, T005, T007, T016 |
| SC-002 | T004, T006, T016 |
| SC-003 | T008, T009, T010, T016 |
| SC-004 | T011, T013, T014, T016, T017 |
| SC-005 | T012, T015, T017, T018 |

## Metrics and next action

15 requirements/outcomes;18 tasks;100 percent coverage; zero unmapped tasks, unresolved placeholders or open findings. All22 requirements checklist items pass. Source identity is S128 / specs/128-prelaunch-baseline / codex/128-prelaunch-baseline / native #408 / milestone 18.

Proceed with /speckit-implement using the disjoint file groups and mandatory red/green tests, then integrated foreground verification. This report was persisted after the read-only analysis completed.

## Implementation discovery review (2026-10-31)

The complete aggregate exposed #409, an existing date-dependent desktop assertion. The bounded injected-clock correction preserves runtime defaults and DST semantics and satisfies the existing deterministic-verification principle. It is traced in plan decision 11 and T019, not a new migration or scheduler outcome. No new constitutional conflict or unresolved scope finding.

The coverage-discovered #410 is within US1's explicit concurrent-helper/current-refusal requirement. Plan decision 12 and T020 retain the first failure and require bounded retry rather than weakened validation. Source-state refusal remains mandatory; no new scope or constitutional exception.
