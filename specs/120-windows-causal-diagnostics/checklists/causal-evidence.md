# Causal Evidence Requirements Checklist

**Purpose**: Review requirement completeness, attribution and ownership before planning/implementation.

**Created**: 2026-10-27

**Feature**: [spec.md](../spec.md)

- [x] CHK001 Are original interpreter and prerequisite failure requirements distinct? [Completeness, FR-002]
- [x] CHK002 Is unresolved cause explicitly separated from delivered correction? [Clarity, FR-001/FR-003]
- [x] CHK003 Are all six writer identity and row-count guarantees explicit? [Completeness, FR-004]
- [x] CHK004 Are command, inspection and cleanup outcomes separately required? [Consistency, FR-009]
- [x] CHK005 Does the spec require reuse of existing diagnostics and budgets? [Completeness, FR-006]
- [x] CHK006 Are diagnostic ownership and size/time bounds required? [Security, FR-007/FR-008]
- [x] CHK007 Are collection failure and perturbation edge cases covered? [Coverage, Edge Cases]
- [x] CHK008 Is original evidence immutable across experiments? [Integrity, FR-001]
- [x] CHK009 Are environment/dependency changes tied to a specific hypothesis? [Clarity, FR-010]
- [x] CHK010 Are missing crash artifacts handled without invented conclusions? [Coverage, Assumptions]
- [x] CHK011 Are new independent defects assigned separate planning records? [Traceability, FR-011]
- [x] CHK012 Are original and new release/source identities distinguished? [Consistency, FR-012]
- [x] CHK013 Are state completeness assertions preserved under partial collection? [Consistency, FR-009]
- [x] CHK014 Are user-facing success criteria measurable without prescribing a tool? [Measurability, SC-001/SC-002]
- [x] CHK015 Are publication and review bounds explicit? [Authority, FR-014]
- [x] CHK016 Are all issue categories covered without claiming repeats prove causes? [Coverage, FR-001/FR-002/FR-005]

All 16 requirements-quality items pass. Standard depth, author/PR-reviewer audience, original-cause integrity and bounded ownership focus. No mandatory checklist hooks are registered. The installed prerequisite helper requires a plan scaffold; setup-plan copied its real template before this checklist, with substantive planning following clarification/checklist completion.
