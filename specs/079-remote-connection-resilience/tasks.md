# Tasks: Resilient Remote Connections

**Input**: Design documents from `specs/079-remote-connection-resilience/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/recovery.md`, `quickstart.md`

**Tests**: Required by the constitution and issue #172. Each behavioral task begins with a failing focused test.

## Phase 1: Setup and Contracts

**Purpose**: Establish the S079 lifecycle, active planning context, and stable recovery vocabulary.

- [x] T001 Register S079 and its Draft delivery state in `.specify/feature.json` and `specs/README.md`
- [x] T002 [P] Validate requirements quality in `specs/079-remote-connection-resilience/checklists/requirements.md` and `specs/079-remote-connection-resilience/checklists/resilience.md`
- [x] T003 [P] Finalize the retry, failure, freshness, and mutation contracts in `specs/079-remote-connection-resilience/contracts/recovery.md`

---

## Phase 2: Foundational Transport and State Model

**Purpose**: Add the typed client and bridge primitives required by every recovery story.

- [x] T004 Add failing transport tests for trust classification and one-attempt uncertain remote mutations in `internal/api/client/remote_test.go`
- [x] T005 Implement safe remote trust and `MutationUncertainError` classification without retries in `internal/api/client/errors.go`, `internal/api/client/client.go`, and `internal/api/client/methods.go`
- [x] T006 Add failing snapshot contract tests for freshness, retry metadata, and recovery modes in `desktop/connection/manager_test.go`
- [x] T007 Extend the connection snapshot and retry-policy seams in `desktop/connection/model.go` and `desktop/connection/manager.go`
- [x] T008 Add failing profile last-contact persistence tests in `internal/clientprofile/store_test.go` and `desktop/connections/service_test.go`
- [x] T009 Implement atomic last-success updates and generation-safe connection callbacks in `internal/clientprofile/store.go`, `desktop/connections/service.go`, and `desktop/connection/model.go`

**Checkpoint**: Typed uncertainty, terminal classification, retry metadata, and successful-contact persistence exist before recovery begins.

---

## Phase 3: User Story 1 - Recover Observation Automatically (Priority: P1)

**Goal**: Recover the selected remote profile automatically after ordinary transport loss while retaining one cancelable connection generation.

**Independent Test**: Interrupt remote health and event transport, restore it, and observe the same profile reconnect through bounded deterministic retry without user action.

- [x] T010 [US1] Add failing deterministic backoff, jitter, stream-loss, target-switch, and shutdown lifecycle tests in `desktop/connection/manager_test.go`
- [x] T011 [US1] Enable remote transient recovery, bounded retry scheduling, reset behavior, and stale-generation rejection in `desktop/connection/manager.go` and `desktop/connection/remote.go`
- [x] T012 [US1] Add reconnect integration coverage across interruptible manager backends and TLS client fixtures in `desktop/connection/manager_test.go`, `desktop/connection/remote_test.go`, and `internal/api/client/remote_test.go`
- [x] T013 [US1] Wire recovered contact persistence and immutable target reuse through `desktop/connections/service.go` and `desktop/app.go`

**Checkpoint**: Network loss, daemon restart, and equivalent suspend or resume failures recover automatically without overlapping work.

---

## Phase 4: User Story 2 - Explain Stale and Terminal States (Priority: P1)

**Goal**: Keep target-scoped data visible with explicit freshness and distinguish wait, retry, repair, update, and identity-stop conditions.

**Independent Test**: Exercise every failure category and verify safe state, retry mode, stale status, last contact, next retry, and accessible guidance.

- [x] T014 [US2] Add failing remote failure taxonomy tests for unreachable, timeout, unauthorized, revoked, forbidden, incompatible, certificate, and identity changes in `desktop/connection/remote_test.go`
- [x] T015 [US2] Implement terminal-versus-transient failure mapping and safe recovery guidance in `desktop/connection/remote.go` and `internal/api/client/errors.go`
- [x] T016 [US2] Add failing frontend tests for stale retained workspaces, retry metadata, and assistive announcements in `desktop/frontend/src/connection/store.test.ts`, `desktop/frontend/src/App.test.tsx`, and `desktop/frontend/src/settings/ConnectionsPage.test.tsx`
- [x] T017 [US2] Extend frontend recovery models and global stale presentation in `desktop/frontend/src/connection/model.ts`, `desktop/frontend/src/components/Shell.tsx`, `desktop/frontend/src/App.tsx`, and `desktop/frontend/src/settings/ConnectionsPage.tsx`
- [x] T018 [US2] Add and satisfy Chromium recovery accessibility and responsive coverage in `desktop/frontend/e2e/connection-recovery.spec.ts`

**Checkpoint**: Users can identify current versus stale data and every terminal recovery path without relying on color.

---

## Phase 5: User Story 3 - Keep Mutations Single-Attempt (Priority: P1)

**Goal**: Report ambiguous remote write outcomes without automatic replay and preserve deliberate user recovery state.

**Independent Test**: Drop each representative mutation response after request receipt and prove one outbound attempt, an `uncertain` result, retained draft or workspace, and a later authoritative refresh.

- [x] T019 [US3] Add failing service-result tests for remote uncertain task and alert mutations in `desktop/taskgroup/service_test.go` and `desktop/operations/service_test.go`
- [x] T020 [US3] Map typed mutation uncertainty to stable bridge results in `desktop/taskgroup/service.go` and `desktop/operations/service.go`
- [x] T021 [US3] Add failing task-editor and operation-store tests for draft preservation and refresh-after-recovery in `desktop/frontend/src/tasks/TaskEditor.test.tsx`, `desktop/frontend/src/tasks/store.test.ts`, and `desktop/frontend/src/operations/store.test.ts`
- [x] T022 [US3] Preserve editor state, retain last complete workspaces, and refresh them on the recovered generation in `desktop/frontend/src/tasks/TasksPage.tsx`, `desktop/frontend/src/tasks/store.ts`, and `desktop/frontend/src/operations/store.ts`

**Checkpoint**: Remote mutations remain one attempt and uncertain outcomes cannot be mistaken for safe failures to retry blindly.

---

## Phase 6: Documentation and Closure

**Purpose**: Synchronize supported behavior, analyze cross-artifact coverage, and gather canonical evidence.

- [x] T023 [P] Document retry, stale-state, trust recovery, and uncertain-mutation boundaries in `docs/remote-access.md`, `docs/api.md`, and `docs/gui-fields.md`
- [x] T024 [P] Add the S079 feature and architecture decisions to `CHANGELOG.md`
- [x] T025 Update S079 to In Progress and run the read-only spec-kit analysis across `specs/079-remote-connection-resilience/spec.md`, `specs/079-remote-connection-resilience/plan.md`, and `specs/079-remote-connection-resilience/tasks.md`
- [x] T026 Run focused Go race, frontend, Chromium, native Wails, secret-canary, encoding, and `git diff --check` validation described by `specs/079-remote-connection-resilience/quickstart.md`
- [x] T027 Run `sh scripts/verify.sh all` and record all eight gates in `specs/079-remote-connection-resilience/verification.md`
- [x] T028 Mark all tasks complete, advance S079 to Implemented, update `specs/README.md`, and record delivery evidence in `specs/079-remote-connection-resilience/spec.md`
- [x] T029 Run `go run ./scripts/github-format`, `sh scripts/spec-lifecycle-check.sh .`, mojibake scans, and final repository cleanliness checks before commit

---

## Dependencies and Execution Order

- Phase 1 establishes lifecycle and contract traceability.
- Phase 2 blocks all three user stories.
- User Story 1 establishes recovery generations consumed by User Stories 2 and 3.
- User Story 2 adds the shared freshness and failure vocabulary consumed by retained feature views.
- User Story 3 completes the mutation boundary and frontend preservation behavior.
- Phase 6 follows all user stories and is the final done gate.

## Parallel Opportunities

- T002 and T003 affect independent specification artifacts.
- T023 and T024 affect independent documentation files once behavior is stable.
- Focused Go and frontend tests in T026 may be invoked as separate foreground commands, but the canonical aggregate remains one foreground sequential run.

## Implementation Strategy

1. Define typed transport outcomes and the expanded connection snapshot.
2. Implement automatic recovery entirely inside the existing connection owner.
3. Project freshness and retry state globally before adapting feature stores.
4. Complete representative task and alert mutation uncertainty paths used by the current remote allowlist.
5. Close with documentation, analysis, native binding validation, canonical verification, and lifecycle evidence.

## Done Gate

S079 is implemented only when all tasks are checked, every #172 acceptance criterion is represented by passing deterministic tests, remote state-changing requests are proven single-attempt under lost responses, all retained remote data is marked stale until authoritative refresh, local IPC compatibility remains intact, and all eight canonical gates pass.
