# Tasks: CLI parity and isolated operation

**Input**: Design artifacts in specs/107-cli-parity-isolation/. Tests are required by constitution principles I/II and the autopilot skill. All six issue functional outcomes are in scope.

## Phase 1: Setup

- [x] T001 Synchronize clean main, verify the six issue scopes and create codex/107-cli-parity-isolation (spec.md).
- [x] T002 Run installed specify/clarify/checklist/plan instructions and create design artifacts (specs/107-cli-parity-isolation/).
- [x] T003 Run installed analyze after complete tasks and resolve blocking findings (plan.md, tasks.md).

## Phase 2: Foundation

- [x] T004 Add failing independent-root, explicit-selector and cancellation tests in internal/cli/invocation_test.go (FR-001, FR-002).
- [x] T005 Replace globals/direct I/O with scoped invocation dependencies and captured-client binding in internal/cli/cli.go, target.go, mcp.go and command handlers (FR-001, FR-002, FR-013).
- [x] T006 Add explicit IPC validation, local-command annotation and dependency interfaces in internal/cli/invocation.go (FR-001, FR-002).

## Phase 3: US1, owned operation and adapters

**Goal**: #288 complete. **Independent test**: selected packaged fixture and recording host-operation handlers affect only owned resources.

- [x] T007 [US1] Add authoritative IPC/PID runtime metadata and tests in cmd/goschedd and internal/api/server (FR-003).
- [x] T008 [US1] Add failing containment/identity/process-cleanup tests in internal/clifixture (FR-003, FR-014).
- [x] T009 [US1] Implement owned workspace/target guard and Windows job/Unix group supervisor in internal/clifixture (FR-003, FR-014).
- [x] T010 [US1] Add recording service/secret/GUI adapters, focused real-handler harness and evidence kinds in internal/clifixture and internal/cli (FR-004).
- [x] T011 [US1] Exercise real packaged CLI and benign task against owned daemon in test/integration/cli_isolation_test.go (FR-003, FR-004, FR-014).

## Phase 4: US2, direct task/group authoring

**Goal**: #283 and #287 complete. **Independent test**: set/clear/omit, disabled create, hierarchy and preview against fixture state.

- [x] T012 [P] [US2] Add failing shared-client environment-clear regression in internal/api/server and desktop/taskgroup tests (FR-005).
- [x] T013 [US2] Implement ClearEnv transport/validation and desktop mapping in internal/api/server/update.go and desktop/taskgroup/service.go (FR-005).
- [x] T014 [US2] Add command request regressions and complete task set/clear/disabled/syntax flags in internal/cli/task.go and task_test.go (FR-005, FR-006, FR-007).
- [x] T015 [US2] Add group edit/disabled create and atomic hierarchy request tests in internal/cli/group.go, internal/api/server/groups.go and internal/store/group.go, including discovered #304 (FR-006).
- [x] T016 [US2] Add complete non-mutating schedule preview and policy/syntax tests in internal/cli/schedule.go (FR-007).

## Phase 5: US3, notification policy

**Goal**: #284 complete. **Independent test**: heterogeneous file/stdin export/reapply, failed replacement and inherited reset.

- [x] T017 [US3] Add input/export/round-trip/conflict tests in internal/cli/notification_test.go (FR-008).
- [x] T018 [US3] Implement strict bounded structured input/canonical export while retaining shorthand in internal/cli/notification.go (FR-008, FR-013).

## Phase 6: US4, operational views

**Goal**: #285 complete. **Independent test**: bounded calendar/exact runs, preserved legacy output and successful partial observation.

- [x] T019 [P] [US4] Add observation identity/access/partial-failure/cancellation tests in internal/clientobserve (FR-010).
- [x] T020 [US4] Implement bounded shared observation helper in internal/clientobserve (FR-010).
- [x] T021 [US4] Add calendar/runtime/active/exact/detail/log metadata commands and compatibility tests in internal/cli (FR-009).
- [x] T022 [US4] Add systems summary/search with selected/all-target output and independent errors in internal/cli (FR-010, FR-013).

## Phase 7: US5, durable preferences

**Goal**: #286 complete. **Independent test**: shared GUI/CLI migration/scoped updates/default selection without daemon/native secrets.

- [x] T023 [P] [US5] Add shared migration/invalid-input/concurrent-update regressions in internal/userprefs (FR-011).
- [x] T024 [US5] Extract preference model/store and cross-process atomic scoped updates; retain desktop wrappers in internal/userprefs and desktop/settings (FR-011).
- [x] T025 [US5] Add headless preference read/set/reset and isolated-path tests in internal/cli/preferences.go (FR-011, FR-013).
- [x] T026 [US5] Add saved profile-default read/set/reset and no-network/no-secret tests in internal/cli/profile.go (FR-012).

## Phase 8: Integration and publication boundary

- [x] T027 Document every new command/input/output and isolation contract in docs/cli.md and relevant maintainer docs (FR-013, SC-001).
- [x] T028 Record dated architecture decisions and unreleased feature/fix notes in CHANGELOG.md; update specs/README.md (SC-001, SC-005).
- [x] T029 Run focused regressions and canonical sh scripts/verify.sh all through completion; record exact evidence in verification.md (SC-001 through SC-005).
- [x] T030 Review six-issue coverage, encoding/GitHub format and diff; commit locally with co-author and stop before push (verification.md).

## Dependencies and execution order

T001-T003 precede implementation. T004-T006 establish command boundary before command-surface changes. T012-T013 may proceed independently of the CLI foundation. T013 precedes T014 environment clears. T007-T010 precede packaged mutation T011. T019/T020 and T023/T024 are independent shared modules; command integration T021/T022/T025/T026 waits for T005/T006. T027-T030 require complete issue outcomes.

## Parallel opportunities and team strategy

Use disjoint file ownership for independent shared modules after foundation design: invocation/isolation, userprefs/desktop settings, and clientobserve. Never edit the same command file concurrently. Research agents already inspected those areas; implementation tasks may be assigned to team members only with these boundaries. Parent integrates task/group/notifications, root registrations, documentation and final verification. This follows the installed task template's parallel team strategy.

US1 tests of ownership and recording adapters form an independently useful foundation; US2 transport tests can run beside it. US3 tests remain in notification files. US4 shared observation and US5 shared persistence have independent files and tests. Complete all five stories before the single review-boundary halt; no partial issue closure is implied.
