# Filesystem Watcher Requirements Checklist

**Purpose**: Validate the completeness, clarity, consistency, and measurability of the S056 watcher contract before implementation
**Created**: 2026-09-05
**Feature**: [spec.md](../spec.md)

## Selection and Dispatch Semantics

- [x] CHK001 Is the distinction between file and directory watcher paths explicit even when the path does not exist? [Clarity, Spec FR-001 and FR-002]
- [x] CHK002 Are exact-file, directory-depth, and base-name pattern selection rules independently specified? [Completeness, Spec FR-003]
- [x] CHK003 Are candidate event classes and deliberately ignored event classes enumerated? [Completeness, Spec FR-004]
- [x] CHK004 Is rename-into-place behavior defined by the resulting stable file rather than a platform-specific low-level event name? [Portability, Spec FR-004]
- [x] CHK005 Are debounce and stability separate, ordered, and objectively measurable concepts? [Clarity, Spec FR-005 and FR-006]
- [x] CHK006 Are temporary disappearance and continued mutation during settling defined? [Edge Case, Spec FR-006]
- [x] CHK007 Is watcher dispatch explicitly subordinate to every existing task eligibility and concurrency rule? [Consistency, Spec FR-007]
- [x] CHK008 Is run provenance defined without retaining the potentially sensitive matched path? [Privacy, Spec FR-008]

## Lifecycle and Recovery Semantics

- [x] CHK009 Is startup behavior prospective and explicit about the absence of durable replay? [Clarity, Spec FR-009]
- [x] CHK010 Are live create, update, state, and delete effects defined without a restart? [Completeness, Spec FR-010]
- [x] CHK011 Does the requirements set define cancellation of pending candidates across configuration generations? [Race Safety, Spec FR-010]
- [x] CHK012 Are every runtime health state, its meaning, its reason, and its transition timestamp specified? [Completeness, Spec FR-011]
- [x] CHK013 Are missing, replaced, inaccessible, unsupported-link, overflow, and observer-failure recovery paths covered? [Recovery, Spec FR-012]
- [x] CHK014 Is repeated-failure deduplication defined for both structured events and logs? [Observability, Spec FR-013]
- [x] CHK015 Are existing and newly created subdirectory rules specified for recursive watchers? [Coverage, Spec FR-014]
- [x] CHK016 Are symbolic-link and Windows-junction traversal boundaries explicit? [Security, Spec FR-014]
- [x] CHK017 Are network and removable path guarantees honestly bounded to host support? [Portability, Spec FR-015]

## Administration and Data Integrity

- [x] CHK018 Is the complete local API lifecycle and standard failure format required? [Completeness, Spec FR-016]
- [x] CHK019 Are human and machine-readable CLI output, duration parsing, errors, and exit behavior specified? [Consistency, Spec FR-017]
- [x] CHK020 Is the desktop distinction between opaque-key triggers and filesystem watchers required? [UX Clarity, Spec FR-018]
- [x] CHK021 Are desktop health, selection, lifecycle, accessibility, and destructive confirmation requirements present? [UX Coverage, Spec FR-018]
- [x] CHK022 Are enabled watchers consistently included in the trigger-ready task activation contract? [Consistency, Spec FR-019]
- [x] CHK023 Are retarget, disable, delete, and task-cascade readiness effects defined? [Data Integrity, Spec FR-019]
- [x] CHK024 Are path normalization and portable pattern validation rules unambiguous? [Validation, Spec FR-020]
- [x] CHK025 Are timing defaults and minimum and maximum bounds quantified? [Measurability, Spec FR-021]
- [x] CHK026 Does the migration requirement preserve existing data and require a real prior-schema fixture? [Migration Safety, Spec FR-024]

## Non-Functional and Completion Criteria

- [x] CHK027 Is nominal configuration scale explicit for lifecycle performance? [Performance, Spec FR-022 and SC-005]
- [x] CHK028 Are goroutine, timer, observer-handle, and recovery-loop termination requirements complete? [Concurrency, Spec FR-023]
- [x] CHK029 Do cross-platform requirements avoid depending on identical native event sequences? [Portability, Spec FR-025]
- [x] CHK030 Are write-storm and rename trial counts sufficient to expose intermittent coalescing defects? [Measurability, Spec SC-001 and SC-002]
- [x] CHK031 Are degraded and recovered health time budgets quantified? [Measurability, Spec SC-003]
- [x] CHK032 Are no-replay, no-flood, race, leak, coverage, and full-gate outcomes measurable? [Completion, Spec SC-004, SC-006, SC-008, and SC-009]
- [x] CHK033 Are authorization, payload, remote delivery, Chain targeting, and replay exclusions explicit? [Scope, Spec Dependencies and Scope]

## Notes

- All 33 requirement-quality questions passed during the S056 clarification review. The installed checklist prerequisite incorrectly requires `plan.md` even though the project mandates checklist before plan, so S056 used the feature path already returned by the successful clarification prerequisite command, matching the established project workaround.
