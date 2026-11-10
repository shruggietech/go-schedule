# Tasks: Correlated parallel task joins

## Phase 1: Setup

- [x] T001 Establish spec, traceability, baseline failures and preserved evidence in specs/144-parallel-task-joins and .test-temp/s144 (FR-001 through FR-019).

## Phase 2: Shared foundation

- [x] T002 Add settled join, invocation, filter and delivery/run correlation domain types in internal/domain/joins.go and domain.go (FR-001, FR-002, FR-003, FR-006).

## Phase 3: US1 complete parallel execution

Independent criterion: Prepare starts two owned branches and Publish is admitted once after their exact successes, with overlapping/same-task invocations isolated.

- [x] T003 [P] [US1] Write failing identity, duplicate-task, interleaving, atomic-last-branch and graph tests in internal/store/joins_test.go (FR-001, FR-002, FR-003, FR-004, FR-010).
- [x] T004 [US1] Implement schema25 definitions/snapshots/role-specific delivery identity and atomic record/claim/resolve decisions in internal/store/current_schema.sql, joins.go, chains.go and store.go (FR-001, FR-002, FR-003, FR-004, FR-006, FR-007, FR-010).
- [x] T005 [US1] Add typed activation sources and current/snapshotted join readiness in internal/task/readiness.go and internal/store CRUD/search/observation helpers (FR-005, FR-009, FR-014).
- [x] T006 [P] [US1] Write failing engine parallel, queued-marker, duplicate/collapse and current-admission regressions in internal/engine/joins_test.go and queued eligibility tests (FR-002, FR-003, FR-004, FR-005, FR-008, FR-017).
- [x] T007 [US1] Integrate exact join delivery origins, active metadata and #471 shared fresh queue/worker eligibility in internal/engine/engine.go and overlap.go (FR-002, FR-003, FR-004, FR-005, FR-008, FR-014, FR-017).

## Phase 4: US2 equivalent complete interfaces

Independent criterion: local/remote authorized authoring, bounded exact observation and portable round-trip through all supported control surfaces.

- [x] T008 [P] [US2] Write failing API/auth/remote/schema/refusal/metadata tests in internal/api/server and internal/remote join tests (FR-006, FR-010, FR-011, FR-019).
- [x] T009 [US2] Add join CRUD/instance paging/live merge, revision conflicts, join events and remote chain/join routes in internal/api/server, internal/authorization, internal/remote and internal/events (FR-001, FR-006, FR-010, FR-011, FR-019).
- [x] T010 [P] [US2] Write failing portable remapping/runtime-exclusion/drift/mixed-cycle tests in internal/bundle and internal/api/server bundle tests (FR-010, FR-012, FR-014).
- [x] T011 [US2] Add portable v3 join configuration, capability checks and atomic final-graph application in internal/bundle, internal/api/server/bundles.go and store portable helpers (FR-010, FR-012, FR-014).
- [x] T012 [P] [US2] Write failing client/CLI human/JSON/exact-instance/file/stdin tests in internal/api/client and internal/cli join tests (FR-001, FR-006, FR-011).
- [x] T013 [US2] Implement join client methods and CLI family in internal/api/client, internal/cli/join.go and command registration (FR-001, FR-006, FR-011).
- [x] T014 [P] [US2] Write failing MCP authority/dedup/paging/redaction/exact-instance tests in internal/mcpmanage and internal/mcpobserve (FR-006, FR-011, FR-015).
- [x] T015 [US2] Add joins_manage and bounded definition/instance/exact/branch resources, updating affected MCP transport/discovery assertions in internal/mcpmanage, internal/mcpobserve, internal/mcphttp, internal/remotemcp and test/integration/mcp_packaged_test.go (FR-001, FR-006, FR-011, FR-015).
- [x] T016 [P] [US2] Write failing desktop join/revision/selected-target/uncertainty/remote-workspace tests in desktop/automation and desktop App tests (FR-001, FR-006, FR-011, FR-013, FR-018, FR-019).
- [x] T017 [US2] Add complete Wails/service/operations models and shared #472 uncertain mutation handling with current expected-daemon/Manage guards in desktop/automation, desktop/operations and desktop/app.go (FR-001, FR-006, FR-011, FR-013, FR-018, FR-019).
- [x] T018 [P] [US2] Write failing keyboard branch editor, same-task identity, exact inspection, duplicate-submit, target-switch and failed-refresh uncertainty tests in desktop/frontend automation/Activity and e2e tests (FR-001, FR-003, FR-006, FR-011, FR-013, FR-018, FR-019).
- [x] T019 [US2] Implement Parallel joins authoring/observation, full remote Automation, exact Activity attribution and target-scoped mutation/read lifecycle in desktop/frontend (FR-001, FR-003, FR-006, FR-011, FR-013, FR-018, FR-019).

## Phase 5: US3 truthful durable outcomes

Independent criterion: partial completion survives current reopen; failure/cancellation/deletion blocks the correct instance without changing unrelated native work.

- [x] T020 [P] [US3] Write failing partial-reopen/recovery/idempotence/edit/delete/late-native-result and notification-preservation tests in internal/store join tests and current-schema fixtures (FR-007, FR-008, FR-009, FR-014).
- [x] T021 [US3] Finish snapshot readiness, deletion resolution, current-format fixtures and exact terminal observation paths in internal/store and internal/engine (FR-005, FR-006, FR-007, FR-008, FR-009, FR-014).
- [x] T022 [US3] Add real owned packaged CLI/daemon two-branch/overlap/cancel/portable observations in test/integration/parallel_joins_test.go (FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-011, FR-012, FR-015).

## Phase 6: Polish and whole-feature acceptance

- [x] T023 Update OpenAPI/generated remote contracts and complete CLI command/scenario catalogue in api/openapi/remote-v1.yaml, internal/remote/openapi.gen.go and internal/clievaluate (FR-006, FR-011, FR-012, FR-019).
- [x] T024 Update CLI/API/MCP/desktop/portable/readiness docs and source-only comparison boundaries in docs and CHANGELOG.md (FR-001 through FR-016, FR-018, FR-019).
- [x] T025 Measure paired engine benchmarks and new join decision performance using internal/engine benchmarks (FR-004, FR-005, FR-014).
- [x] T026 Run owned hidden Windows/Linux/WSL and headless complete UI verification, retaining original failures and coverage limits in specs/144-parallel-task-joins/verification.md (FR-011, FR-013, FR-014, FR-015).
- [x] T027 Freeze source/historical inputs and pass the complete watched foreground eight-gate aggregate in specs/144-parallel-task-joins/verification.md (all FRs, SC-001 through SC-005).
- [x] T028 Finish source status/index/delivery and UTF8/no-BOM/format/privacy/retention checks in specs/144-parallel-task-joins, specs/README.md and CHANGELOG.md (FR-015, FR-016).

## Dependencies and parallel execution

T001 -> T002 -> independent [P] store (T003-T005), engine/API/portable (T006-T011) and interface (T012-T019) file-owned work, with failing tests before each implementation. Store owns internal/store and internal/task; engine/API owns internal/engine, API/server, auth/remote operations, bundle/events and daemon wiring; interfaces own API/client, CLI/MCP and desktop. Root owns shared domain, OpenAPI/generated contract, catalogue, docs and packaged/performance/final verification. Agents coordinate method changes through the settled contract; no simultaneous writes to shared files. T020-T022 close durable integrated behavior before T023-T028 whole-source acceptance.

## Implementation strategy

Deliver the complete parallel block in one source PR, including #471/#472 and required remote chain CRUD. A backend-only partial result is not S144 completion. Push/PR/CI/review are continuing authorized publication workflow, not unchecked implementation tasks.
