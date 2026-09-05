# UX Requirements Checklist: Persisted Adjustable Columns

**Purpose**: Test accessibility, resilience, and consistency of the written requirements **Created**: 2026-09-05 **Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are both affected views and explicit exclusions defined? [Completeness, Spec §Scope]
- [x] CHK002 Are pointer, keyboard, focus, naming, reset, persistence, and responsive requirements documented? [Completeness, Spec §FR-001–FR-011]
- [x] CHK003 Are preserved table behaviors listed? [Completeness, Spec §FR-013–FR-014]

## Requirement Clarity and Consistency

- [x] CHK004 Is adjacent-only boundary behavior unambiguous? [Clarity, Spec §FR-004]
- [x] CHK005 Are minimum and below-minimum behaviors distinguished? [Clarity, Spec §FR-005–FR-006]
- [x] CHK006 Is persistence normalized, versioned, current-user, and per-view? [Clarity, Spec §FR-007–FR-010]
- [x] CHK007 Do resize, restore, and reset preserve width and isolation? [Consistency, Spec §FR-004, §FR-008, §FR-011]
- [x] CHK008 Does accessibility align with list navigation and disclosure? [Consistency, Spec §FR-003, §FR-014]

## Acceptance and Coverage

- [x] CHK009 Can conservation, restoration, fallback, isolation, and default improvement be measured? [Measurability, Spec §SC-001–SC-005]
- [x] CHK010 Does every requirement map to a scenario or outcome? [Traceability, Spec §User Scenarios, §Success Criteria]
- [x] CHK011 Are termination, tiny widths, malformed data, schema change, rebuilds, key limits, and focused reset covered? [Coverage, Spec §Edge Cases]
- [x] CHK012 Are display scale, font, theme, and recreation scenarios specified? [Coverage, Spec §US2, §FR-013]

## Dependencies and Assumptions

- [x] CHK013 Is preference ownership stated and daemon storage excluded? [Assumption, Spec §Assumptions, §Scope]
- [x] CHK014 Is Tasks adjustment deliberately excluded? [Boundary, Spec §Scope, §Assumptions]
