# Performance & Measurement Requirements Checklist: Run engine benchmarks in CI and enforce the p99 dispatch-latency budget

**Purpose**: Validate that the requirements for the p99 latency gate and CI
benchmark wiring are complete, unambiguous, measurable, and consistent — before
planning and implementation. These are "unit tests for the requirements," not
tests of the code.
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)

## Measurement Correctness

- [x] CHK001 Is the measured quantity defined precisely as scheduled-time → execution-start, with command execution explicitly excluded? [Clarity, Spec §FR-001, Key Entities]
- [x] CHK002 Is "nominal load" for the budget defined so the measurement conditions are unambiguous (scheduling overhead in isolation)? [Ambiguity, Spec §Assumptions]
- [x] CHK003 Is the statistic to assert specified as the p99 (not mean/median), matching the constitution's budget phrasing? [Consistency, Spec §FR-001, §FR-003]
- [x] CHK004 Is the number of samples over which the p99 is computed specified as a fixed, documented value? [Completeness, Spec §Assumptions, Edge Cases]
- [x] CHK005 Are requirements clear that the measurement excludes unrelated work (no command execution, no I/O) so the number reflects dispatch overhead only? [Clarity, Spec §FR-001, §FR-011]

## Budget Definition & Traceability

- [x] CHK006 Is the budget defined as a single named value co-located with the engine code it governs, so budget and code cannot drift? [Completeness, Spec §FR-002]
- [x] CHK007 Do the requirements state that the assertion checks against that same named budget value (not a separately-hardcoded number)? [Consistency, Spec §FR-002]
- [x] CHK008 Is the budget value (p99 < 100 ms) traceable to the constitution's Performance principle? [Traceability, Spec §Overview, §FR-002]

## Failure Behavior & Acceptance Criteria

- [x] CHK009 Is the failure condition specified (p99 exceeds budget) with a requirement that the failure message names observed p99 vs budget? [Clarity, Spec §FR-003]
- [x] CHK010 Are the acceptance scenarios for pass (within budget) and fail (over budget) both defined and objectively verifiable? [Acceptance Criteria, Spec §US1]
- [x] CHK011 Is there a measurable success criterion that a regression pushing p99 over budget is demonstrably caught? [Measurability, Spec §SC-003]

## Non-Flakiness & Determinism

- [x] CHK012 Is a non-flakiness requirement stated for shared/loaded CI hardware, with the headroom rationale that makes it stable? [Completeness, Spec §FR-004, Edge Cases]
- [x] CHK013 Do the requirements forbid dependence on real sleeps for the correctness of the assertion, consistent with the constitution's testing standard? [Consistency, Spec §FR-004]
- [x] CHK014 Is coarse timer granularity on the target platforms addressed as an edge case that must not cause false failures? [Edge Case, Spec §Edge Cases]
- [x] CHK015 Is the requirement that the check runs in the standard test suite (locally and in CI) without a C toolchain stated? [Completeness, Spec §FR-005]

## CI Wiring Correctness

- [x] CHK016 Is it required that CI actually executes the engine benchmarks on each run (not merely compiles them)? [Completeness, Spec §FR-006, §SC-004]
- [x] CHK017 Is publishing the benchmark output as a retrievable artifact specified as a distinct, verifiable requirement? [Completeness, Spec §FR-006, §US2]
- [x] CHK018 Is it required that the CI benchmark run executes benchmarks only (not the functional suite again) and needs no C toolchain? [Clarity, Spec §FR-007]
- [x] CHK019 Is the goroutine-leak test's existing race-detector coverage stated as a confirmation item with no new code, and is non-exclusion from the race selection required? [Coverage, Spec §FR-008, §US3]

## Decision Record & Pinned Artifacts

- [x] CHK020 Is the absolute-budget-over-relative-delta gate captured as an explicit, recorded decision with rationale (not left implicit)? [Completeness, Spec §FR-009, §SC-006]
- [x] CHK021 Is the requirement to record any pinned-artifact (CI workflow) change as a dated changelog decision present? [Consistency, Spec §FR-010]
- [x] CHK022 Do the requirements consistently satisfy the constitution's "benchmark regression checks" obligation via the absolute gate, with the deviation from the relative-10 % clause justified? [Conflict, Spec §Assumptions, §FR-009]

## Scope Boundaries

- [x] CHK023 Is it explicitly required that engine runtime dispatch behavior is not altered and no latency logging/metrics are added to production paths? [Completeness, Spec §FR-011]
- [x] CHK024 Is the exclusion of benchstat baseline machinery and signing/release concerns clearly bounded out of scope? [Coverage, Spec §Overview scope-out, §Assumptions]

## Notes

- All items interrogate the **requirements**, not the eventual code. An unchecked
  item means the spec needs tightening before implementation, not that a test
  failed.
- Highest-risk items for this feature: CHK004 (sample count), CHK012–CHK014
  (non-flakiness), and CHK022 (constitutional consistency of the absolute gate).
