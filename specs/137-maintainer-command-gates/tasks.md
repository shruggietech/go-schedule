# Tasks: Maintainer command gates

## Phase 1: Setup

- [x] T001 Confirm accepted main, preserve prior evidence and complete S137 specification/design in specs/137-maintainer-command-gates/.

## Phase 2: Foundation

- [x] T002 Extract the source-only fixture writer into test/scripts/lib/automation-fixture.sh without changing existing policy cases.
- [x] T003 Add failing-first integration regressions in internal/maintainercheck/command_gates_test.go.

## Phase 3: US1 - Integrated local contracts

Independent criterion: normal automation propagates all focused contract failures; an isolated broken-stop driver remains finite.

- [x] T004 [US1] Implement isolated status/diagnostic/order and broken-stop guards in test/scripts/verify-driver_test.sh; delegate the existing failure-propagation mode.
- [x] T005 [US1] Compose lifecycle, public entrypoint and driver guards in test/scripts/maintainer-commands_test.sh and wire scripts/verify.sh automation.

## Phase 4: US2 - Native platform contract

Independent criterion: both real wrappers inspect every selected literal fixture root; the native test matrix cannot silently drop the shared corpus.

- [x] T006 [US2] Implement the fixed root-selection matrix in test/scripts/public-entrypoints_test.sh using the shared fixture writer.
- [x] T007 [US2] Add scoped integration enforcement in internal/maintainercheck/command_gates.go and wire .github/workflows/ci.yml jobs.test.

## Phase 5: Polish and verification

- [x] T008 Update CONTRIBUTING.md, CHANGELOG.md and specs/README.md with accurate integration and process decisions.
- [x] T009 Run focused regressions and all eight foreground canonical gates, audit preserved inputs and record specs/137-maintainer-command-gates/verification.md.

## Dependencies and implementation strategy

T001 precedes T002/T003. The shared fixture is foundational for T006. T003 precedes T007. T004 and T006 feed T005, then CI integration and guidance complete the end-to-end slice. The fixture extraction and policy regression authoring affect different files but are performed sequentially by the primary agent. US1 local enforcement is the first usable increment; US2 adds the existing native delivery boundary. Publication, reviews and owner merge are delivery operations rather than unchecked implementation tasks.
