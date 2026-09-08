# Tasks: Dependable Webhook Notifications

**Input**: Design documents from `specs/067-webhook-notifications/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/local-api.md, contracts/webhook-v1.schema.json, quickstart.md

**Tests**: Required by the specification and constitution. Behavioral and contract tests are written and observed failing before implementation.

## Phase 1: Setup and Contracts

**Purpose**: Register the slice and freeze receiver-facing and security boundaries.

- [x] T001 Add S067 lifecycle inventory and current-plan context in `specs/README.md` and `CLAUDE.md`
- [x] T002 [P] Validate JSON Schema syntax and webhook fixture expectations in `specs/067-webhook-notifications/contracts/webhook-v1.schema.json`
- [x] T003 [P] Add migration v15 failure-first coverage in `internal/store/migration_v15_test.go`

---

## Phase 2: Foundational Notification Model

**Purpose**: Establish plain entities, forward-only storage, and redaction before delivery behavior.

- [x] T004 Add failing domain serialization and validation tests in `internal/domain/notification_test.go`
- [x] T005 Define notification channel, assignment, effective-policy, delivery, state, event, filter, and payload entities in `internal/domain/notification.go`
- [x] T006 Implement forward-only notification schema migration v15 in `internal/store/store.go`
- [x] T007 Add failing channel lifecycle, assignment replacement, precedence, removal, filtering, retention, and secret-erasure store tests in `internal/store/notifications_test.go`
- [x] T008 Implement channel, assignment, effective-policy, delivery claim/recovery/transition/filter/prune persistence in `internal/store/notifications.go`
- [x] T009 Extend run transaction tests for matching atomic notification creation and task-outcome isolation in `internal/store/chains_test.go` and `internal/store/notifications_test.go`
- [x] T010 Extend `RecordRunAndCreateDeliveries` to create immutable matching notification work atomically in `internal/store/chains.go`

**Checkpoint**: Notification intent is durable, redacted, and transactionally coupled only to recorded run identity.

---

## Phase 3: User Story 1 - Send a Safe Run-Outcome Webhook (Priority: P1)

**Goal**: Deliver the documented payload for selected task outcomes without exposing task or channel secrets.

**Independent Test**: A controlled receiver validates payload, headers, filtering, and redaction for successful and failed task runs.

- [x] T011 [P] [US1] Add failing endpoint validation, summary, payload, header, redirect, and diagnostic-redaction tests in `internal/notification/webhook_test.go`
- [x] T012 [US1] Implement endpoint validation, safe summaries, v1 payload construction, and one-attempt HTTP transport in `internal/notification/webhook.go`
- [x] T013 [US1] Add failing engine integration tests for post-commit notification wakeup and source-run invariance in `internal/engine/completion_test.go`
- [x] T014 [US1] Add a post-commit notification wake callback without network work in the task worker in `internal/engine/engine.go`

**Checkpoint**: One matching run produces a safe standards-based request and a distinct durable delivery.

---

## Phase 4: User Story 2 - Apply Task and Group Policy Predictably (Priority: P1)

**Goal**: Configure and explain nearest-scope replacement for task and nested-group assignments.

**Independent Test**: Store and API tests demonstrate task, child-group, parent-group, and empty-scope resolution with multiple channels and no duplicate assignments.

- [x] T015 [P] [US2] Add failing assignment and effective-policy API contract tests in `internal/api/server/notifications_test.go`
- [x] T016 [US2] Implement task and group assignment replacement plus effective-policy handlers in `internal/api/server/notifications.go` and register them in `internal/api/server/server.go`
- [x] T017 [US2] Add shared client assignment and effective-policy methods with request tests in `internal/api/client/notifications.go` and `internal/api/client/notifications_test.go`

**Checkpoint**: Every task's selected policy is inspectable and comes from one deterministic scope.

---

## Phase 5: User Story 3 - Survive Delivery Failure and Restart (Priority: P1)

**Goal**: Process webhook work asynchronously with bounded retries, recovery, and retention.

**Independent Test**: Deterministic runtime tests cover success, rejection, transport failure, cancellation, stable duplicate IDs, preserved attempt counts, exhaustion, pruning, and goroutine shutdown.

- [x] T018 [US3] Add failing dispatcher lifecycle, retry, recovery, exhaustion, pruning, worker-isolation, and race-safe shutdown tests in `internal/notification/dispatcher_test.go`
- [x] T019 [US3] Implement the bounded dispatcher, injected delay/clock boundaries, recovery, wake coalescing, and terminal transitions in `internal/notification/dispatcher.go`
- [x] T020 [US3] Wire dispatcher ownership, engine wakeup, startup recovery, and shutdown into `cmd/goschedd/main.go`, covered through dispatcher and daemon-owned integration lifecycle tests
- [x] T021 [US3] Publish structured delivery failure and recovery logs containing only safe correlation fields in `internal/notification/dispatcher.go`

**Checkpoint**: Receiver failure and daemon restart remain bounded and cannot consume scheduler task workers or mutate run outcomes.

---

## Phase 6: User Story 4 - Manage and Inspect Webhook Channels (Priority: P2)

**Goal**: Provide complete local API and CLI channel lifecycle plus redacted delivery evidence.

**Independent Test**: API, client, CLI, and installed-daemon tests create, test, update, disable, rotate, filter, and remove a channel while preserving terminal history.

- [x] T022 [P] [US4] Add failing channel lifecycle, test-delivery, deletion, delivery-filter, and secret-redaction API tests in `internal/api/server/notifications_test.go`
- [x] T023 [US4] Implement channel and delivery handlers and dispatcher wake interface in `internal/api/server/notifications.go` and `internal/api/server/server.go`
- [x] T024 [US4] Add shared client channel and delivery methods with request tests in `internal/api/client/notifications.go` and `internal/api/client/notifications_test.go`
- [x] T025 [US4] Add failing CLI management-surface registration and outcome-validation tests in `internal/cli/notification_test.go`, with request and response behavior covered at the shared client and API boundaries
- [x] T026 [US4] Implement `gosched notification` channel, assignment, policy, test, and delivery commands in `internal/cli/notification.go` and register them in `internal/cli/cli.go`
- [x] T027 [US4] Add daemon-owned webhook delivery lifecycle and task-outcome isolation scenarios in `test/integration/notifications_test.go`

**Checkpoint**: Operators can manage and diagnose the full webhook lifecycle without direct database access.

---

## Phase 7: Documentation and Verification

- [x] T028 Document receiver setup, payload, commands, precedence, security boundary, retries, duplicates, retention, and common generic destinations in `docs/notifications.md`, `docs/api.md`, `docs/cli.md`, and `docs/README.md`
- [x] T029 Add S067 features and dated architecture decisions to `CHANGELOG.md`
- [x] T030 Run focused domain, store, notification, engine, daemon, API, client, CLI, and integration tests
- [x] T031 Run `go run ./scripts/github-format` and UTF-8/mojibake scans over changed publication content
- [x] T032 Run `sh scripts/verify.sh all` in the foreground and record every gate in `specs/067-webhook-notifications/verification.md`
- [x] T033 Mark all tasks complete, set the specification to Implemented, update `specs/README.md` delivery evidence, and re-run read-only cross-artifact analysis

## Dependencies and Execution Order

- Setup freezes contracts before storage implementation.
- Foundational model and migration work blocks every user story.
- User Story 1 establishes safe payload and transport behavior.
- User Story 2 can proceed after foundational persistence and does not depend on outbound runtime completion.
- User Story 3 depends on the one-attempt transport and durable claim transitions.
- User Story 4 composes channel persistence, policies, and dispatcher wakeup into operator interfaces.
- Documentation and canonical verification depend on every included story.

## Parallel Opportunities

- T002 and T003 touch independent contract and migration-test files.
- T011 and T015 can establish transport and policy API contracts independently after storage foundations exist.
- Focused package tests can run independently before the sequential canonical verifier.

## Implementation Strategy

The smallest releasable increment is User Stories 1 through 3 together because durable policy, safe delivery, and failure isolation form one correctness boundary. User Story 4 completes issue #159's required lifecycle and evidence surface in the same slice. Do not ship delivery without its management, restart, redaction, and removal contracts.
