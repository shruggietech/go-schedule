# Tasks: S109 CLI contract and diagnostic repairs

Input: [specification](spec.md), [plan](plan.md), [research](research.md) and [contracts](contracts/cli.md). Tests are mandatory under constitution II. Full bundle: #311-#322. No publication task is part of implementation completion.

## Phase 1: Setup

- [x] T001 Synchronize baseline and create specs/109-cli-contract-diagnostics/spec.md with full issue traceability.
- [x] T002 Complete clarification, requirements checklists and design in specs/109-cli-contract-diagnostics/.

## Phase 2: Foundation

- [x] T003 Add failing parser, operand, runtime-category and no-side-effect regressions in internal/cli/invocation_contract_test.go (#311/#312).
- [x] T004 Implement typed parser/operand boundaries and default arity in internal/cli/cli.go and internal/cli/invocation.go (#311/#312, FR-001/FR-002).

## Phase 3: User Story 1, reliable scripted invocations

Independent outcome: invalid invocations have correct status and no side effects; targeting and reports retain their intended behavior.

- [x] T005 [US1] Add endpoint ownership and selector conflict regressions in internal/cli/invocation_contract_test.go (#314).
- [x] T006 [US1] Add --daemon-endpoint and ownership-aware resolution in internal/cli/cli.go and internal/cli/invocation.go (#314, FR-004).
- [x] T007 [P] [US1] Add failing unknown/duplicate binding tests in internal/cli/invocation_contract_test.go and internal/api/server/bundles_test.go (#315).
- [x] T008 [US1] Refuse duplicate CLI and raw JSON bindings plus unknown watcher identities in internal/cli/bundle.go and internal/api/server/bundles.go (#315, FR-005).
- [x] T009 [P] [US1] Add failing mixed valid/unsupported/invalid import regressions in internal/cli/cron_test.go (#316).
- [x] T010 [US1] Retain full import reports and classify invalid input in internal/cli/cron.go (#316, FR-006).
- [x] T011 [US1] Classify malformed local bundle input as usage while preserving file I/O runtime failures in internal/cli/bundle.go (#311, FR-001).

## Phase 4: User Story 2, useful human and machine output

Independent outcome: JSON stdout is clean and human help/list/warning messages explain actual operations and data.

- [x] T012 [US2] Add failing alerts-channel, help and column-heading regressions in internal/cli/communication_contract_test.go (#313/#317/#320).
- [x] T013 [US2] Route alerts notice to stderr in internal/cli/notification.go or the actual alerts handler (#313, FR-003).
- [x] T014 [US2] Describe credential/pairing/service effects and relevant inputs in internal/cli/access.go and internal/cli/service.go (#317, FR-007).
- [x] T015 [US2] Label actor/audit/pairing/profile/channel/watcher columns in internal/cli/access.go, profile.go, notification.go and watcher.go (#320, FR-010).
- [x] T016 [P] [US2] Add warning contract regression and explain effective IPC settings in cmd/goschedd/main.go and an adjacent test (#322, FR-012).

## Phase 5: User Story 3, accurate diagnostic and state context

Independent outcome: diagnostics identify the failed input/object and rotation/watcher responses show truthful association/observation context.

- [x] T017 [US3] Add failing lookup, field, schedule, audit, certificate and malformed-file regressions in internal/cli/diagnostics_test.go (#318).
- [x] T018 [US3] Add typed contextual diagnostics and field-specific correction in internal/cli, preserving codes and secret privacy (#318, FR-008).
- [x] T019 [P] [US3] Add failing rotation identity, token, last-used and failure-atomicity tests in internal/enrollment/service_test.go and internal/store/remote_access_test.go (#319).
- [x] T020 [US3] Assemble rotation identity context before commit in internal/store/remote_access.go and internal/enrollment/service.go (#319, FR-009).
- [x] T021 [P] [US3] Add failing definition-observation, coalescing and unchanged-health regressions in internal/watcher/manager_test.go and internal/api/server/filesystem_watchers_test.go (#321).
- [x] T022 [US3] Correlate actual observed definition with saved watcher state in internal/domain, internal/watcher and internal/api/server/filesystem_watchers.go (#321, FR-011).
- [x] T023 [US3] Explain pending observations in internal/cli/watcher.go and propagate additive context through desktop automation mappings/model/card with regression tests (#321, FR-011/FR-013).

## Phase 6: Integration and local delivery

- [x] T024 Refresh internal/clievaluate/inventory-v1.json and targeted scenarios without discarding historical report/case identity; cover --daemon-endpoint (#314/#317, FR-013).
- [x] T025 Update docs/cli.md, relevant response documentation and CHANGELOG.md with corrected contracts and dated decisions (FR-013).
- [x] T026 Run focused Windows and native Linux race regressions for affected packages and inspect secret/output compatibility (FR-001 through FR-013).
- [x] T027 Run the complete foreground scripts/verify.sh all aggregate and record every gate in specs/109-cli-contract-diagnostics/verification.md.
- [x] T028 Build fresh receipt-bound inputs, complete Windows/selected-WSL capture and all contextual reviews, and finalize/validate/render the canonical evaluation in .test-temp/s109-*; record each F001-F012 disposition in specs/109-cli-contract-diagnostics/verification.md.
- [x] T029 Complete all tasks, update specification/index to Implemented, write local delivery and PR draft in specs/109-cli-contract-diagnostics/, audit UTF-8/format and commit feat(109) locally before the pre-publication halt.

## Dependencies and execution strategy

Setup precedes foundation. US1 uses the classification boundary. US3 diagnostics compose that boundary; US2 output changes can proceed in independent files. Enrollment/observation internals can proceed independently after the contracts gate, while shared CLI files remain single-writer. T024-T029 integrate all three stories. The MVP is US1, but this authorized delivery includes all three stories and all 12 issues.

Parallel examples: enrollment regressions and atomic response work (T019/T020) do not conflict with the CLI foundation. Watcher runtime/server/desktop work (T021/T022 and desktop part of T023) can proceed independently while the main implementer owns CLI edits. Binding decoder/server work can be isolated from human output edits. Finish integration before aggregate verification; do not commit or change source during receipt-bound capture.

## Workflow

Installed /speckit-tasks produced 29 tasks: setup 2, foundation 2, US1 7, US2 5, US3 7, integration 6. All tasks have IDs, appropriate story labels and paths. No task extension hooks are registered.
