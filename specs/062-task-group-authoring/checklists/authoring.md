# Task and Group Authoring Requirements Checklist

**Purpose**: Review the completeness, clarity, consistency, and measurability of S062 requirements before implementation
**Created**: 2026-09-07
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are overview requirements defined for empty, populated, high-volume, and degraded task and group collections? [Completeness, Spec §User Story 1, §FR-001-004]
- [x] CHK002 Are all durable task fields and scheduling policies explicitly enumerated? [Completeness, Spec §FR-006]
- [x] CHK003 Are every supported task and group lifecycle mutation and its cancellation path specified? [Completeness, Spec §FR-005, §FR-016, §FR-021-022]
- [x] CHK004 Are basic, advanced, recurring, one-off, manual-only, and legacy-empty task scenarios documented? [Coverage, Spec §User Stories 2-3]
- [x] CHK005 Are declared state, inherited state, readiness, and effective eligibility requirements independently defined? [Completeness, Spec §FR-002-003, §FR-019]
- [x] CHK006 Are the safe command mapping, insertion behavior, guidance surfaces, execution constraints, and documentation obligations all specified? [Completeness, Spec §FR-011-015, §FR-028-029]

## Requirement Clarity

- [x] CHK007 Is daemon authority after mutations and live updates unambiguous? [Clarity, Spec §Clarifications, §FR-022-024]
- [x] CHK008 Is the difference between a full group path, group name, root choice, and ungrouped choice explicit? [Clarity, Spec §FR-003, §FR-017]
- [x] CHK009 Is exact command preview defined as a non-executable program plus ordered arguments? [Clarity, Spec §FR-008]
- [x] CHK010 Is suggestion insertion limited by editor mode, input state, modifier key, host support, and one-shot traversal behavior? [Clarity, Spec §FR-013-014]
- [x] CHK011 Are destructive confirmations required to name both the entity and active daemon plus the relevant consequence? [Clarity, Spec §FR-021]
- [x] CHK012 Are scale and responsiveness targets quantified separately from daemon round-trip time? [Clarity, Spec §SC-004]

## Requirement Consistency

- [x] CHK013 Do inactive-by-default creation requirements align with Run now guidance and current draft safety? [Consistency, Spec §User Story 2, §FR-007]
- [x] CHK014 Do suggestion requirements use the execution host consistently instead of the desktop client platform? [Consistency, Spec §FR-011, §Edge Cases]
- [x] CHK015 Do group deletion requirements align with preserving affected tasks as ungrouped? [Consistency, Spec §User Story 4, §FR-020]
- [x] CHK016 Do live-refresh requirements preserve stable selection and unsaved drafts without permitting optimistic local records? [Consistency, Spec §FR-022-024]
- [x] CHK017 Does the local-only boundary remain consistent with future execution-host and multi-daemon readiness? [Consistency, Spec §FR-030, §Assumptions]

## Acceptance Criteria Quality

- [x] CHK018 Can known-good task completion be measured by time, platform coverage, bounded command duration, and observable output? [Measurability, Spec §SC-001-002, §SC-009]
- [x] CHK019 Can task field and group operation parity be audited without accepting silent value changes? [Measurability, Spec §SC-003]
- [x] CHK020 Can concurrency and lifecycle quality be evaluated through a quantified one-hundred-cycle outcome? [Measurability, Spec §SC-005]
- [x] CHK021 Are accessibility outcomes quantified by severity, state coverage, and keyboard journey coverage? [Measurability, Spec §SC-007-008]
- [x] CHK022 Can removal of unsafe first-run guidance be objectively audited? [Measurability, Spec §SC-010]

## Scenario and Edge-Case Coverage

- [x] CHK023 Are remote deletion, concurrent changes, stale drafts, slow responses, and late responses addressed? [Coverage, Spec §Edge Cases, §FR-022-024]
- [x] CHK024 Are deep hierarchies, duplicate names, missing parents, malformed references, and cycle attempts addressed? [Coverage, Spec §Edge Cases, §FR-017-018]
- [x] CHK025 Are long, Unicode, bidirectional, and empty values covered without weakening inspectability or actions? [Coverage, Spec §Edge Cases]
- [x] CHK026 Are schedule-preview failures separated from saved-task and execution-record behavior? [Coverage, Spec §Edge Cases, §FR-009]
- [x] CHK027 Are repeat submission, pending-state suppression, cancellation, connection loss, and safe recovery requirements defined? [Coverage, Spec §User Story 5, §FR-022, §FR-025]

## Security, Privacy, and Accessibility

- [x] CHK028 Are ordinary announcements and confirmations prohibited from exposing command, input, identity, path, environment, or backend-error detail? [Security, Spec §FR-026]
- [x] CHK029 Are validation association, focus movement, keyboard operation, visible focus, non-color meaning, reduced motion, zoom, and narrow-window requirements documented? [Accessibility, Spec §FR-010, §FR-027]
- [x] CHK030 Are suggested-command safety properties explicit for network, filesystem, privilege, interactivity, runtime availability, duration, and observable output? [Security, Spec §FR-015]

## Dependencies and Boundaries

- [x] CHK031 Are S061 connection dependencies and existing scheduling authorities identified without duplicating their logic? [Dependency, Spec §Assumptions]
- [x] CHK032 Are Activity migration, other feature screens, packaging cutover, remote access, and Fyne removal explicitly excluded? [Boundary, Spec §FR-030, §Assumptions]
- [x] CHK033 Is the temporary dual-desktop suggestion obligation bounded to the period before cutover? [Dependency, Spec §FR-029]
