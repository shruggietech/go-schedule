# Operational Safety Requirements Checklist

**Purpose**: Validate that task activation, effective eligibility, and failure-diagnostic requirements are complete and reviewable **Created**: 2026-09-05 **Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are requirements defined for both activation choices and their default? [Completeness, Spec §FR-001..FR-003]
- [x] CHK002 Is creation atomicity specified strongly enough to exclude transient eligibility? [Completeness, Spec §FR-002]
- [x] CHK003 Are validation recovery and fresh-dialog reset requirements both defined? [Completeness, Spec §FR-004]
- [x] CHK004 Are create and edit behaviors explicitly distinguished? [Completeness, Spec §FR-005]
- [x] CHK005 Is omitted-field behavior defined for existing non-desktop callers? [Completeness, Spec §FR-006]
- [x] CHK006 Are configured, lifecycle, and effective task states separately defined? [Completeness, Spec §FR-007..FR-009]
- [x] CHK007 Are direct and ancestor group suppression requirements both covered? [Completeness, Spec §FR-008]
- [x] CHK008 Are all failed-run detail fields enumerated? [Completeness, Spec §FR-015..FR-017]
- [x] CHK009 Are legacy alert, missing task/run, empty output, and retrieval failure states specified? [Completeness, Spec §FR-018]
- [x] CHK010 Are all supported run-trigger origins included? [Completeness, Spec §FR-021]

## Requirement Clarity

- [x] CHK011 Is the precedence between task-disabled, lifecycle-inactive, and group-blocked states unambiguous? [Clarity, Spec §FR-009]
- [x] CHK012 Is the responsible group deterministically identified as the nearest disabled group? [Clarity, Spec §FR-008]
- [x] CHK013 Is run identity explicitly durable rather than inferred by time or ordering? [Clarity, Spec §FR-013, §FR-019]
- [x] CHK014 Is a missing exit code given one honest, bounded interpretation? [Clarity, Spec §Edge Cases]
- [x] CHK015 Is combined process output labeled without implying separate streams? [Clarity, Spec §FR-016]
- [x] CHK016 Is truncation disclosure distinct from output contents and cap enforcement? [Clarity, Spec §FR-017]

## Requirement Consistency

- [x] CHK017 Does the inactive GUI default remain consistent with backward-compatible omitted-field behavior? [Consistency, Spec §FR-001, §FR-006]
- [x] CHK018 Does effective-state presentation preserve rather than redefine scheduler eligibility? [Consistency, Spec §FR-010]
- [x] CHK019 Do exact-run diagnostics and deleted-history fallback avoid contradictory lookup behavior? [Consistency, Spec §FR-018..FR-019]
- [x] CHK020 Do output visibility requirements align with the prohibition on newly exposing secret-bearing inputs? [Consistency, Spec §FR-016, §FR-020]

## Acceptance Criteria Quality

- [x] CHK021 Can the activation default and zero transient dispatch outcome be objectively measured? [Measurability, Spec §SC-001..SC-002]
- [x] CHK022 Can every effective-state category be evaluated from one full explanation? [Measurability, Spec §SC-003]
- [x] CHK023 Is exact correlation measurable across consecutive failures and trigger modes? [Measurability, Spec §SC-004]
- [x] CHK024 Are failure-output variants and byte-bound behavior measurable? [Measurability, Spec §SC-005]
- [x] CHK025 Are migration, compatibility, coverage, race, and lint completion signals defined? [Measurability, Spec §SC-006..SC-007]

## Scenario and Edge-Case Coverage

- [x] CHK026 Are primary, alternate, exception, recovery, and compatibility flows represented across the three stories? [Coverage, Spec §User Scenarios]
- [x] CHK027 Are nested groups, cycles, missing references, deletion, and live refresh addressed? [Coverage, Spec §Edge Cases, §FR-011]
- [x] CHK028 Are nonzero exit, launch failure, empty, multiline, Unicode, and truncated output addressed? [Coverage, Spec §US1, §Edge Cases]
- [x] CHK029 Are narrow-window, dark/light, full-disclosure, and non-color-only requirements present? [Accessibility, Spec §FR-012]
- [x] CHK030 Are sensitive arguments, stdin, environment data, and unrelated Activity redesign explicitly excluded? [Security, Scope]

## Dependencies and Assumptions

- [x] CHK031 Is the dependency on established group-chain policy explicit? [Dependency, Spec §FR-010]
- [x] CHK032 Is the existing combined-output boundary documented as an assumption rather than separate-stream delivery? [Assumption, Spec §Assumptions]
- [x] CHK033 Are #119 column persistence and broader Activity redesign clearly outside this slice? [Boundary, Spec §Out of scope]
- [x] CHK034 Are individual issue closure conditions preserved for #102, #118, and #120? [Traceability, Spec §Assumptions]

## Notes

- Standard-depth reviewer checklist focused on execution correctness, backward compatibility, accessibility, and diagnostic privacy.
- Validated 2026-09-05: 34/34 requirements-quality checks pass.
