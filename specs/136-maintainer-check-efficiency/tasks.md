# Tasks: Maintainer check efficiency

**Input**: S136 specification, plan, research, model and command contract. Tests are required by FR-004 and the constitution.

## Phase 1: Setup

- [x] T001 Preserve unchanged source, separate native measurements and contract research in .test-temp/s136 and specs/136-maintainer-check-efficiency/research.md.
- [x] T002 Record bundled #433/#434/#435/#436 dependencies and frozen input inventory in specs/136-maintainer-check-efficiency/verification.md.

## Phase 2: Foundation

- [x] T003 Add evaluation-local read/YAML handling, explicit errors and regression scaffolding in internal/maintainercheck/check.go and check_test.go.

## Phase 3: US1 - Efficient automation policy

Independent acceptance: intended existing automation fixture decisions remain, action-like comments/script strings are ignored and executable unknown actions still fail.

- [x] T004 [US1] Add failing action-field, malformed/duplicate YAML and publication-order regressions in internal/maintainercheck/check_test.go.
- [x] T005 [US1] Implement structured policy checks and retained release/notes/operations contracts in internal/maintainercheck/automation.go.
- [x] T006 [US1] Retain the full corpus with valid synthetic YAML/lifecycle data and a reused checker in test/scripts/automation-check_test.sh.

## Phase 4: US2 - Efficient lifecycle

Independent acceptance: existing states, task/delivery rules, inventory and disposition failure categories remain.

- [x] T007 [US2] Add lifecycle boundary and mutation/cache regressions in internal/maintainercheck/lifecycle_test.go before implementing internal/maintainercheck/lifecycle.go.
- [x] T008 [US2] Reuse one built checker in test/scripts/spec-lifecycle-check_test.sh without removing its existing cases.

## Phase 5: US3 - Compatible execution

Independent acceptance: existing shell commands resolve caller roots, preserve failure propagation and launch native children hidden on Windows.

- [x] T009 [US3] Implement scripts/maintainer-check/main.go and compatible scripts/automation-check.sh and scripts/spec-lifecycle-check.sh adapters with explicit child failure handling.
- [x] T010 [US3] Add command/root/child regressions and relevant benchmarks in scripts/maintainer-check/main_test.go and internal/maintainercheck/*_test.go.

## Phase 6: Polish and integrated verification

- [x] T011 Update maintainer guidance, dated decision, changelog and specs/README.md for the unchanged commands and execution path.
- [x] T012 Run focused race/corpus/benchmark and native timing comparisons; run all eight foreground canonical gates and preserve source/input identity in specs/136-maintainer-check-efficiency/verification.md.
- [x] T013 Resolve implementation findings, complete task/lifecycle records, formatting/encoding and immutable-evidence checks in specs/136-maintainer-check-efficiency/ and CHANGELOG.md.

## Dependencies and execution

T001 -> T002/T003 -> US1 and US2 -> US3 -> T011-T013. Tests precede each policy implementation. Shared cache/files require sequential edits; independent documentation can be drafted alongside test execution. US1 is the first useful increment; all three stories are required for #433. GitHub push, PR review and final owner merge are delivery evidence, not unchecked implementation tasks.
