# Security and UX Requirements Checklist: External Trigger Lifecycle

**Purpose**: Validate security, lifecycle, diagnostics, accessibility, and desktop requirements before implementation
**Created**: 2026-09-05
**Feature**: [spec.md](../spec.md)

## Key Security Requirements

- [x] CHK001 Are entropy, generation authority, representation, and uniqueness requirements specified for keys? [Completeness, Spec §FR-002]
- [x] CHK002 Are all ordinary surfaces from which raw keys must be excluded explicitly identified? [Coverage, Spec §FR-005, Spec §FR-014]
- [x] CHK003 Are the limited actions authorized to return a raw key unambiguous? [Clarity, Spec §FR-005]
- [x] CHK004 Is immediate old-key invalidation defined as part of the atomic rotation result? [Consistency, Spec §FR-006]
- [x] CHK005 Is the local trust boundary explicit without implying that keys are safe to log? [Clarity, Spec §Assumptions]
- [x] CHK006 Are deleted and replaced key behaviors specified to avoid disclosing former existence? [Security, Spec §Edge Cases]

## Invocation and Failure Requirements

- [x] CHK007 Are all prerequisites for an eligible invocation enumerated separately? [Completeness, Spec §FR-010]
- [x] CHK008 Are machine-readable outcomes specified for every known rejection class? [Coverage, Spec §FR-012]
- [x] CHK009 Is the accepted-request guarantee distinguished from exactly-once external effects? [Clarity, Spec §Clarifications]
- [x] CHK010 Are concurrent calls and overlap-policy ownership specified consistently? [Consistency, Spec §SC-003]
- [x] CHK011 Are daemon termination and non-replay behavior defined? [Recovery, Spec §Edge Cases]
- [x] CHK012 Is the no-additional-listener boundary explicit and objectively reviewable? [Scope, Spec §FR-009]

## Lifecycle and Readiness Requirements

- [x] CHK013 Are create, update, enable, disable, rotate, delete, and target-deletion transitions covered? [Completeness, Spec §FR-004, Spec §FR-007]
- [x] CHK014 Are final-source disable rules and multi-source preservation rules both defined? [Consistency, Spec §FR-016]
- [x] CHK015 Is it clear that trigger creation and enablement never silently enable a task? [Clarity, Spec §FR-017]
- [x] CHK016 Are readiness refresh inputs complete across trigger, task, command, lifecycle, enabled, and group states? [Coverage, Spec §FR-025]

## Desktop Interaction Requirements

- [x] CHK017 Are table identity, state, readiness, and redaction requirements specified together? [Completeness, Spec §FR-021]
- [x] CHK018 Are every desktop lifecycle and copy action named? [Coverage, Spec §FR-022]
- [x] CHK019 Are confirmation requirements defined for each destructive operation and tied to visible identity? [Clarity, Spec §FR-023]
- [x] CHK020 Are keyboard, assistive labeling, appearance, and accessible confirmation requirements measurable? [Acceptance Criteria, Spec §SC-008]

## Dependencies and Boundaries

- [x] CHK021 Are the existing daemon, dispatch, readiness, group, and IPC contracts named as dependencies? [Dependency, Spec §FR-009, Spec §FR-011, Spec §FR-015, Spec §FR-028]
- [x] CHK022 Are Trigger Sets, watchers, remote access, payloads, Chain targets, and replay explicitly excluded? [Scope, Spec §Assumptions]
- [x] CHK023 Are migration and persistence requirements tied to the preceding schema and restart behavior? [Recovery, Spec §FR-003, Spec §FR-027]
- [x] CHK024 Are performance and scale targets quantified for both dispatch latency and concurrency? [Measurability, Spec §SC-003, Spec §SC-006]

## Notes

- Standard reviewer-depth checklist focused on the two highest-risk domains, key security and desktop lifecycle usability.
- All requirements-quality items passed on 2026-09-05 after the clarification decisions were integrated.
