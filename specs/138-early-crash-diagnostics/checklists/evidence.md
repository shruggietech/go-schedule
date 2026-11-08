# Evidence Requirements Checklist: Early-crash incident evidence

**Purpose**: Assess requirement quality for incident attribution and bounded collection.

**Created**: 2026-11-08

**Feature**: [spec.md](../spec.md)

## Attribution and bounds

- [x] CHK001 Are ownership, PID, time, provider and application requirements explicit? [Completeness, FR-002]
- [x] CHK002 Are duration, count and byte bounds required? [Measurability, FR-003]
- [x] CHK003 Are evidence absence and collection failure distinct from the native outcome? [Clarity, FR-003/FR-004]
- [x] CHK004 Is opt-in behavior and launch-failure handling defined? [Coverage, FR-001/FR-006]
- [x] CHK005 Are hidden child execution and configuration boundaries preserved? [Consistency, FR-005]
- [x] CHK006 Is weak historical attribution distinguished from exact owned evidence? [Clarity, FR-007]
- [x] CHK007 Are the original causal issues explicitly left open? [Scope, FR-007]
- [x] CHK008 Are dependency and framework limits stated? [Consistency, Assumptions]
