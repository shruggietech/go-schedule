# Tasks: Active-run cancellation

## Phase 1: Setup

- [x] T001 Establish S143 scope, authority, baseline probes and #465/#466/#467 traceability in specs/143-active-run-cancellation/spec.md (FR-001 through FR-014).

## Phase 2: Shared foundation

- [x] T002 Add shared result/active observation fields and context-cause/failure-feedback contracts in internal/domain/cancellation.go, internal/domain/domain.go and internal/runcontrol/control.go (FR-003, FR-005, FR-008).

## Phase 3: US1 stop one execution

Independent criterion: cancel one exact active Run ID, preserve its sibling and reclaim capacity only after native completion.

- [x] T003 [P] [US1] Write failing cancellation/process isolation/input regressions in internal/executor/cancellation_test.go and platform process tests (FR-002, FR-010).
- [x] T004 [US1] Implement owned platform process/stream lifecycle in internal/executor/executor.go, process.go, process_unix.go and process_windows.go (FR-002, FR-003, FR-010).
- [x] T005 [P] [US1] Write failing exact-ID/race/finishing/cleanup-feedback tests in internal/engine/cancellation_test.go (FR-001, FR-004, FR-005, FR-008).
- [x] T006 [US1] Implement per-run handles, completion retirement and active failure observation in internal/engine/cancellation.go and engine.go (FR-001, FR-003, FR-005, FR-008).
- [x] T007 [US1] Retain the #467 failing skip/queue regressions and fix per-task in-flight accounting in internal/engine/engine.go, overlap.go and overlap_test.go (FR-004, FR-014).
- [x] T008 [US1] Add metadata-only context-aware completed-run fallback and durable result coverage in internal/store/runs_cancel.go and tests (FR-005, FR-011).

## Phase 4: US2 consistent interfaces

Independent criterion: equivalent authorized local/remote cancellation and observed result through API, CLI, MCP and desktop; refusals never mutate work.

- [x] T009 [P] [US2] Write failing local/remote actor, revocation, target-identity and API contract regressions in internal/api/server and internal/remote tests (FR-006).
- [x] T010 [US2] Add runs.cancel Operate/audit/RetryNever routing and local API endpoint in internal/authorization/catalog.go, internal/remote/operations.go and internal/api/server/runs.go/server.go (FR-006, FR-008).
- [x] T011 [P] [US2] Write failing CLI human/JSON and MCP idempotency/refusal tests in internal/cli/cancellation_test.go and internal/mcpoperate tests (FR-006, FR-008).
- [x] T012 [US2] Add API client, active ID discovery, runs cancel and controlled runs_cancel tooling in internal/api/client/methods.go, internal/cli/runs.go and internal/mcpoperate/operate.go (FR-006, FR-008).
- [x] T013 [P] [US2] Retain failing #466 local/remote v2.0.3 probes and actual unsupported/missing API rejection coverage in desktop/connection tests (FR-013).
- [x] T014 [US2] Replace product-major gating with advertised API negotiation in desktop/connection/local.go and remote.go and update explicit API fixtures (FR-013).
- [x] T015 [P] [US2] Write failing selected-daemon cancellation/exact active-history handoff tests in desktop/operations tests (FR-005 through FR-008).
- [x] T016 [US2] Add operations/backend and Wails cancellation bridge plus exact-record refresh (#468) in desktop/operations and desktop App sources (FR-006, FR-007, FR-008).
- [x] T017 [P] [US2] Add failing headless component/interaction tests for Cancel run, pending/confirmed/failure states, identity switches, same-render lock and exact refresh in desktop/frontend/src/operations and e2e tests (FR-005, FR-007, FR-008, FR-012).
- [x] T018 [US2] Wire Activity/Run details cancellation, bounded reconciliation, accessible controls and target-scoped lifecycle in desktop/frontend/src/operations and App.tsx (FR-006, FR-007, FR-008).

## Phase 5: US3 truthful results and workflow effects

Independent criterion: retain one canceled result/output across restart with no completion-chain dispatch or false notification streak update.

- [x] T019 [P] [US3] Write canceled/cancellation-failed persistence, chain suppression and outbound notification streak regressions in internal/store and internal/engine tests (FR-003, FR-009, FR-011).
- [x] T020 [US3] Integrate cancellation outcome handling in internal/store notification/completion code, internal/store/search.go and desktop outcome models/labels (FR-003, FR-009).
- [x] T021 [US3] Add owned daemon/CLI integration cancellation with concurrent parent/child work and authoritative final observation in test/integration (FR-001 through FR-006, FR-008, FR-012).

## Phase 6: Polish and final source acceptance

- [x] T022 Update existing OpenAPI/CLI/MCP/desktop/catalogue documentation and cancellation messages in docs, internal/api and existing contract catalogues (FR-006 through FR-010).
- [x] T023 Measure dispatch/cancellation registration benchmarks and retain the existing dispatch-latency contract in internal/engine (FR-004, FR-005).
- [x] T024 Run focused Windows/WSL owned-process and headless UI verification, preserve failure evidence and record actual coverage limits in specs/143-active-run-cancellation/verification.md (FR-002, FR-007, FR-012).
- [x] T025 Freeze source/historical inputs, pass the complete foreground eight-gate canonical aggregate and record exact-source results in specs/143-active-run-cancellation/verification.md (all FRs, SC-001 through SC-005).
- [x] T026 Complete source lifecycle/index/changelog/delivery records and final UTF-8/no-BOM/GitHub-format/privacy checks in specs/143-active-run-cancellation, specs/README.md and CHANGELOG.md (FR-011, FR-012).

## Dependencies and parallel execution

T001 -> T002 -> independent executor (T003-T004), engine/API (T005-T010) and interface work (T011-T018), with tests before their corresponding source changes. Engine and executor share only the settled runcontrol/domain contract. Interface compile validation follows the API client contract. T019-T021 integrate the shared source result, then T022-T026 close the whole capability. Tasks touching the same files remain sequential; root coordinates shared changes and final acceptance. Push/PR/review/merge are publication workflow, not unchecked source implementation tasks.

## Implementation strategy

Use the exact-run backend first, then connect every supported interface and workflow consumer in one source PR. Do not substitute a backend-only partial result for S143 completion. The confirmed prerequisite and exact-refresh corrections remain in the same coherent slice.
