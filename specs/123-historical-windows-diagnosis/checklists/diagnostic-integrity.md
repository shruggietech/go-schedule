# Diagnostic integrity requirements checklist

**Purpose**: Review the completeness and consistency of S123 requirements.

**Created**: 2026-10-27

**Feature**: [spec.md](../spec.md)

## Completeness and attribution

- [x] CHK001 Are original source, wrapper, runtime and binary identities distinguished? [Completeness, FR-001/FR-002]
- [x] CHK002 Are reconstructed and mixed-component observations explicitly separated? [Clarity, FR-002]
- [x] CHK003 Are all six simultaneous writer outcomes and both schema conditions required? [Coverage, FR-003]
- [x] CHK004 Are CLR, launch, prerequisite and writer outcomes separate? [Consistency, FR-004]
- [x] CHK005 Are all five IPC targets and their prerequisite setup in scope? [Completeness, FR-005]
- [x] CHK006 Are command, protocol, observation and cleanup outcomes distinct? [Clarity, FR-005]
- [x] CHK007 Are historical and supported allowances independently attributed? [Consistency, FR-006]
- [x] CHK008 Is missing historical instrumentation an explicit unavailable state? [Exception coverage, FR-006]

## Boundaries and completion

- [x] CHK009 Is the single deadline defined before experiment execution? [Measurability, FR-007]
- [x] CHK010 Are preparation and later acceptance times separate? [Clarity, FR-007]
- [x] CHK011 Are first failures and separately demonstrated defects preserved? [Recovery coverage, FR-008]
- [x] CHK012 Is no recurrence explicitly inconclusive? [Acceptance quality, FR-009]
- [x] CHK013 Are upstream claims bound to matching evidence? [Consistency, FR-009]
- [x] CHK014 Does scope preserve original immutable artifacts and owner boundaries? [Security, FR-001/FR-011]
- [x] CHK015 Does the deliverable reconcile S122 without rewriting earlier history? [Traceability, FR-010]
- [x] CHK016 Is speculative replay expansion excluded after the bounded disposition? [Scope, US3]
