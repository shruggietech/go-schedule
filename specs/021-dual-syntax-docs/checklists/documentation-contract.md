# Documentation Contract Requirements Checklist

**Purpose**: Review whether S021's requirements fully and consistently define the dual-syntax documentation posture before planning
**Created**: 2026-08-28
**Feature**: [spec.md](../spec.md)
**Depth**: Standard PR-review gate
**Audience**: Author and third-party reviewers

## Requirement Completeness

- [x] CHK001 Are all current discovery surfaces (README, root CLI help, task help, GUI guidance, API contract, and cron guide) explicitly included? [Completeness, Spec §FR-001–FR-008]
- [x] CHK002 Are equivalent copy/pasteable human and cron creation examples required rather than merely suggested? [Completeness, Spec §FR-003]
- [x] CHK003 Are automatic detection, explicit syntax identity, no fallback, retention, and response identity all defined for API guidance? [Completeness, Spec §FR-007]
- [x] CHK004 Are `convert`, `explain`, `import`, and `export` each required to have distinct documented purposes? [Completeness, Spec §FR-005]
- [x] CHK005 Are cron expression and crontab line/file explicitly distinguished? [Completeness, Spec §FR-012]
- [x] CHK006 Are master-spec, historical-spec, source-comment, policy-test, and changelog follow-through all included? [Completeness, Spec §FR-013–FR-018]

## Requirement Clarity

- [x] CHK007 Is "dual syntax" defined as human phrases plus the supported five-field cron subset, rather than unrestricted cron? [Clarity, Spec §Context, FR-001]
- [x] CHK008 Is the teaching priority clear without making human syntax technically privileged? [Clarity, Spec §US1, FR-001]
- [x] CHK009 Is the exact dialect-detail inventory named, including fields, names, collections, steps, descriptors, and extensions? [Clarity, Spec §FR-008]
- [x] CHK010 Is DOM/DOW behavior stated as traditional OR plus a named product refusal, avoiding an intersection implication? [Clarity, Spec §FR-009]
- [x] CHK011 Is field-local step meaning distinguished from an interval starting at the field minimum? [Clarity, Spec §FR-010]
- [x] CHK012 Is the meaning of "historical" separated from current authoritative product guidance? [Clarity, Spec §Edge Cases, Assumptions]

## Requirement Consistency

- [x] CHK013 Do README, CLI, GUI, and API requirements use the same acceptance boundary? [Consistency, Spec §FR-001–FR-007]
- [x] CHK014 Do timezone and DST requirements preserve one compiled recurrence model across both source syntaxes? [Consistency, Spec §FR-011]
- [x] CHK015 Does preserving S008 chronology align with updating the authoritative master specification? [Consistency, Spec §FR-013–FR-014]
- [x] CHK016 Do issue-closing requirements align with leaving fidelity breadth in #22? [Consistency, Spec §FR-020]
- [x] CHK017 Does the no-runtime-change boundary remain consistent with CLI help and source-comment updates? [Consistency, Spec §FR-004, FR-015, FR-019]

## Acceptance Criteria Quality

- [x] CHK018 Can current-surface consistency be measured through a bounded inventory and policy search? [Measurability, Spec §SC-001, SC-004]
- [x] CHK019 Can equivalent examples be objectively compared for the same recurrence? [Measurability, Spec §SC-002]
- [x] CHK020 Is complete dialect-topic coverage measured against the explicit issue inventory? [Measurability, Spec §SC-003]
- [x] CHK021 Can historical supersession coverage be counted without rewriting original content? [Measurability, Spec §SC-005]
- [x] CHK022 Are policy fixtures required to prove both allowed and rejected documentation states? [Measurability, Spec §SC-006]

## Scenario and Edge-Case Coverage

- [x] CHK023 Are discovery, detailed migration, and maintainer-policy journeys independently specified? [Coverage, Spec §US1–US3]
- [x] CHK024 Are invalid, unsupported, or lossy cron forms required to retain named refusal guidance? [Coverage, Spec §US1.3, US2.1–US2.3]
- [x] CHK025 Are historical changelog entries excluded from categorical-current-language policy scans? [Edge Case, Spec §Edge Cases]
- [x] CHK026 Are accurate human-parser comments protected from over-broad replacement? [Edge Case, Spec §Edge Cases, FR-015]
- [x] CHK027 Are cross-platform examples, link behavior, and UTF-8 punctuation addressed? [Edge Case, Spec §Edge Cases]

## Dependencies and Boundaries

- [x] CHK028 Is shipped S018-S020 behavior identified as the source of truth and not reopened? [Dependency, Spec §Assumptions]
- [x] CHK029 Is the fidelity table identified as canonical and future breadth assigned solely to #22? [Dependency, Spec §Assumptions]
- [x] CHK030 Is the calendar-step refusal correction bounded separately from excluded syntax breadth, API shape, persistence, engine, GUI, security, packaging, and release changes? [Boundary, Spec §FR-019, FR-021]
- [x] CHK031 Is the silent-approximation defect specified with a measurable named-refusal outcome and no-mutation requirement? [Correctness, Spec §FR-019, SC-008]

## Notes

- All 31 requirements-quality checks pass. The checklist emphasizes contract completeness, dialect accuracy, historical supersession, and reviewable closure evidence. It does not test implementation behavior.
