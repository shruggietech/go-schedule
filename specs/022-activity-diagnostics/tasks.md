# Tasks: Activity Diagnostics Clarity

**Input**: Design documents from `specs/022-activity-diagnostics/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required. Each behavioral seam is recorded red before implementation.

## Phase 1: Setup and Baseline

- [X] T001 Create the review branch, record focused test and documentation baselines in `verification.md`, and validate the requirements checklist
- [X] T002 Complete API, Activity UI, daemon logging, documentation, and test-surface research in `research.md`

---

## Phase 2: User Story 1 - Find the Complete Daemon Log (Priority: P1)

**Goal**: Activity truthfully describes its bounded records and displays the exact daemon-reported full-log path in all record states.

### Tests

- [X] T003 [US1] Add failing server tests for exact `log_path` metadata on populated, empty, and nil-ring responses in `internal/api/server/logs_test.go`
- [X] T004 [US1] Add a failing view-model refresh test for atomic record and exact-path propagation in `gui/viewmodel/viewmodel_test.go`
- [X] T005 [US1] Add failing pure presentation tests for exact Windows/custom paths, bounded-view wording, and unavailable metadata in `gui/logs_test.go`
- [X] T006 [US1] Add failing Activity integration tests for refreshed path and initial unavailable wording in `gui/app_test.go`
- [X] T007 [US1] Add failing CLI regression tests proving human output remains a table and JSON remains a bare record array in `internal/cli/logs_test.go`
- [X] T008 [US1] Run the focused server/view-model/GUI/CLI tests and record expected pre-implementation failures in `verification.md`

### Implementation

- [X] T009 [US1] Extend the server and logs response with exact configured path metadata, including nil/empty rings, in `internal/api/server/server.go` and `internal/api/server/logs.go`
- [X] T010 [US1] Return the typed response from `internal/api/client/logs.go` and explicitly preserve CLI record-only human and JSON output in `internal/cli/logs.go`
- [X] T011 [US1] Store records and log path atomically while preserving path metadata across live events in `gui/viewmodel/viewmodel.go`
- [X] T012 [US1] Add the passive, wrapped Activity diagnostics label without changing existing controls in `gui/logs.go`
- [X] T013 [US1] Update all affected fakes and internal server constructor call sites
- [X] T014 [US1] Run focused API, CLI, view-model, and GUI suites green and record evidence

---

## Phase 3: User Story 2 - Recognize Completed Daemon Startup (Priority: P2)

**Goal**: One startup-completion record clearly names the ready event and its endpoint, database, and log path.

### Tests

- [X] T015 [US2] Add a failing startup helper test for exact message, one record, and `endpoint`, `db`, and `log_path` attributes in `cmd/goschedd/main_test.go`
- [X] T016 [US2] Run the focused daemon test and record the expected pre-implementation failure in `verification.md`

### Implementation

- [X] T017 [US2] Extract and use the single `daemon startup complete` helper and pass the resolved log path into the server in `cmd/goschedd/main.go`
- [X] T018 [US2] Run the focused daemon suite green and record evidence

---

## Phase 4: Documentation and Compatibility

- [X] T019 [US1] Document bounded Activity semantics and the exact configured path in `docs/cli.md`
- [X] T020 [US1] Update `specs/004-rebrand-gui-overhaul/contracts/api-logs.md` with the additive response metadata and remove its stale Logs-view wording
- [X] T021 [US1] Run documentation integrity checks and the CLI output-shape regression suite
- [X] T022 [US1] [US2] Add the newest chronological Unreleased changelog entry for issues #27 and #31

---

## Phase 5: Polish and Cross-Cutting Verification

- [X] T023 Validate both checklists and run Spec-Kit analysis to zero critical findings
- [X] T024 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across changed files
- [X] T025 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results
- [X] T026 Mark every completed task `[X]`, rerun analysis, and commit locally with the required co-author trailer
- [X] T027 Halt before push and pull-request creation with a concise compliance and delivery summary that requires `Closes #27` and `Closes #31` in the PR body

## Dependencies and Execution Order

- Baseline, clarification, checklists, research, and analysis precede implementation.
- US1 response tests establish the shared contract before client, view-model, and UI changes.
- US2 is independently testable after the server constructor carries the resolved path.
- Documentation follows the stable implementation contract.
- Analysis, audits, and canonical gates follow all code and documentation work.

## Implementation Strategy

1. Prove exact metadata behavior at the daemon response boundary.
2. Preserve public CLI output while propagating the richer response to Activity.
3. Present one truthful, passive diagnostics message in every UI state.
4. Make the startup record discrete and directly testable.
5. Document, analyze, verify, commit locally, and halt before publication.
