# Tasks: Correlated parallel recovery and cleanup

## Phase 1: Setup

- [x] T001 Retain accepted baseline, authority and protected input inventories in .test-temp/s145 and the traceability records in specs/145-parallel-recovery-cleanup (FR-016, FR-017).

## Phase 2: Shared foundation

- [x] T002 Add settled condition/settlement/route constants, optional definition fields, separate instance stages and roles in internal/domain/joins.go (FR-001, FR-002, FR-004, FR-006, FR-014).

## Phase 3: US1 exact aggregate recovery

Independent criterion: one exact aggregate selects only success or configured recovery after all branch settlements.

- [x] T003 [P] [US1] Write failing omission/empty policy, exact settlement, duplicate/reopen and routing regressions in internal/store/join_recovery_test.go (FR-001, FR-002, FR-003, FR-004, FR-006, FR-007).
- [x] T004 [US1] Implement schema26 initialization/refusal, complete definition/snapshot projections and unique stage identities in internal/store/current_schema.sql, joins.go and join_definition.go (FR-001, FR-004, FR-006, FR-007, FR-014, FR-016).
- [x] T005 [US1] Preserve #478 exact legacy blocked reasons and implement transactional branch settlement and selected follow-up advancement in internal/store/join_instances.go and join_runs.go (FR-002, FR-003, FR-004, FR-006, FR-007).
- [x] T006 [US1] Extend causal graph, activation/readiness, immutable snapshot eligibility and invalidation in internal/store/chains.go, completion_admission.go, join_sources.go, join_delete.go and internal/task (FR-008, FR-009, FR-010, FR-015).
- [x] T007 [P] [US1] Write failing no-start wake, active sibling/failed cancellation, queued/capacity, exact role and #481 skip-policy same-task reuse regressions in internal/engine/recovery_test.go, recovery_overlap_test.go and recovery_no_start_overlap_test.go (FR-002, FR-003, FR-006, FR-008, FR-015).
- [x] T008 [US1] Implement coalesced event-loop wake, zero-claim drain termination, #481 ownership-before-publication and new exact roles in internal/engine/engine.go and delivery.go (FR-002, FR-003, FR-006, FR-008, FR-015).
- [x] T009 [P] [US1] Write failing authority/identity/audit/strict policy and live stage merge tests in internal/api/server and internal/remote join tests (FR-001, FR-004, FR-011, FR-014).
- [x] T010 [US1] Extend join DTOs, full authoring, bounded exact observations/state filters and capability manifest in internal/api/server and remote regressions (FR-001, FR-004, FR-011, FR-014).

## Phase 4: US2 ordered final cleanup

Independent criterion: cleanup waits for exact branch and chosen follow-up settlement, including explicit no-start, without rewriting outcomes.

- [x] T011 [P] [US2] Write failing follow-up failure/no-start/unselected, delayed native termination and cleanup idempotence tests in internal/store/join_recovery_test.go (FR-002, FR-005, FR-006, FR-007, FR-009).
- [x] T012 [US2] Complete atomic cleanup admission/settlement and invalidation-aware late evidence in internal/store join lifecycle/readiness files (FR-005, FR-006, FR-007, FR-008, FR-009).
- [x] T013 [US2] Add cleanup engine/API integration regressions for no-start-only progress, active ownership and current admission in internal/engine and internal/api/server (FR-002, FR-005, FR-008, FR-011).

## Phase 5: US3 equivalent complete interfaces

Independent criterion: every maintained interface authors, preserves and observes the same complete routing contract.

- [x] T014 [P] [US3] Write failing portable v3/v4 remap, empty policy, drift, graph and preview/apply capability regressions in internal/bundle and internal/api/server (FR-001, FR-004, FR-010, FR-012, FR-015).
- [x] T015 [US3] Implement conditional v4 export/readers, complete reference remapping/fingerprints and target compatibility in internal/bundle and internal/api/server/bundles_joins.go (FR-001, FR-004, FR-010, FR-012, FR-015).
- [x] T016 [P] [US3] Write failing pointer-preserved client packets and CLI recovery/cleanup flags, human/JSON/clear/error cases in internal/api/client and internal/cli/join_test.go (FR-001, FR-004, FR-011, FR-014).
- [x] T017 [US3] Correct #479 human output error propagation and extend client/CLI sparse writable fields and exact separate stage output in internal/api/client/joins.go and internal/cli/join.go (FR-001, FR-004, FR-011, FR-014).
- [x] T018 [P] [US3] Write failing MCP management empty/default policy and bounded exact stage projection regressions in internal/mcpmanage and internal/mcpobserve (FR-004, FR-006, FR-011, FR-014).
- [x] T019 [US3] Propagate controlled management and allowlisted bounded observations in internal/mcpmanage/joins.go and internal/mcpobserve/joins.go (FR-001, FR-004, FR-011, FR-014).
- [x] T020 [P] [US3] Write failing actual facade/full-name sparse edits, capability, target/revision and uncertainty regressions in desktop/automation and desktop app tests (FR-001, FR-004, FR-011, FR-013, FR-014).
- [x] T021 [US3] Extend capability-gated full authoring drafts, sparse CAS and field errors in desktop/app.go and desktop/automation (FR-001, FR-004, FR-011, FR-013, FR-014).
- [x] T022 [P] [US3] Write failing keyboard condition/routing preview/exact stages and Activity role tests in desktop/frontend unit/e2e files (FR-001, FR-004, FR-005, FR-011, FR-013, FR-014).
- [x] T023 [US3] Complete recovery/cleanup controls, vertical routing help and exact stage/Activity views in desktop/frontend automation and operations (FR-001, FR-004, FR-005, FR-011, FR-013, FR-014).

## Phase 6: Polish and complete acceptance

- [x] T024 Update OpenAPI/generated contracts and the truthful CLI catalogue in api/openapi/remote-v1.yaml, internal/remote/openapi.gen.go and internal/clievaluate (FR-011, FR-014, FR-017).
- [x] T025 Update current CLI/API/MCP/desktop/portable/architecture guidance and dated CHANGELOG decisions in docs and CHANGELOG.md (FR-001 through FR-017).
- [x] T026 Add actual owned Windows/Linux packaged sequence/overlap/cancel/portable regressions in test/integration/parallel_recovery_test.go (FR-002, FR-003, FR-005, FR-006, FR-007, FR-011, FR-012, FR-015; SC-001 through SC-004).
- [x] T027 Run complete headless browser and actual Linux/WSL focused runtime coverage, retaining scope limits in specs/145-parallel-recovery-cleanup/verification.md (FR-011, FR-013, FR-017; SC-001 through SC-004).
- [x] T028 Measure unchanged accepted-main/candidate dispatch/fanout and new route decisions with internal/engine benchmarks and record results in specs/145-parallel-recovery-cleanup/verification.md (FR-006, FR-008, FR-015; SC-005).
- [x] T029 Freeze source/retained inputs and pass complete watched foreground eight-gate acceptance recorded in specs/145-parallel-recovery-cleanup/verification.md (FR-001 through FR-017; SC-001 through SC-005).
- [x] T030 Finish implemented lifecycle/index/delivery and UTF8/no-BOM/format/privacy/integrity checks in specs/145-parallel-recovery-cleanup and specs/README.md (FR-016, FR-017).

## Dependencies and parallel execution

T001 -> T002 -> file-owned store/task (T003-T006, T011-T012), engine/API/portable (T007-T010, T013-T015) and interfaces (T016-T023). Tests precede implementation in each ownership track. US2 shares the stage foundation with US1; its native-ordering criterion is independently exercised. Root owns shared domain, generated OpenAPI/catalogue, docs and T024-T030 qualification. Store owns internal/store/internal/task; engine owns internal/engine except root's new route benchmark, API/server, bundle and remote operation tests; interfaces own API/client, CLI/MCP and desktop. Changes to the same files remain sequential and coordinated.

Each independent [P] task has separate owned files after T002. A full root race/build or complete aggregate starts only when all source writers are idle. Resource-bounded focused tests may run through verified hidden foreground sessions.

## Implementation strategy

Deliver the entire #477 outcome in one source PR, not a backend-only partial. Retain the original failing observations, current-only refusal and external-command uncertainty. All authoring/observation surfaces, portable fidelity and exact native ordering are required. Publication/CI/review/owner merge are workflow evidence, not unchecked implementation tasks. No tag, release or cluster work is part of S145.
