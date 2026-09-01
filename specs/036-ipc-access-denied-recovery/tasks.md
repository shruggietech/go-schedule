# Tasks: IPC Access-Denied Recovery

**Input**: Design documents from `specs/036-ipc-access-denied-recovery/`

**Tests**: Required by constitution Principle II and issue #90.

## Phase 1: Setup and Red Baseline

- [x] T001 Record issue #90, v0.9.1 amplification, current error copy, native
  prerequisites, and security-boundary baseline in `specs/036-ipc-access-denied-recovery/verification.md`.
- [x] T002 Update lifecycle inventory in `specs/README.md` and active Spec Kit
  context in `CLAUDE.md`.
- [x] T003 Add failing connection-classification regressions in
  `internal/api/client/errors_test.go`.
- [x] T004 Add a failing simultaneous model/calendar/stream incident regression
  in `gui/connection_test.go` and record the focused red result in
  `specs/036-ipc-access-denied-recovery/verification.md`.

## Phase 2: Foundational Classification and Incident State

- [x] T005 Implement typed transport failure categories and wrapping in
  `internal/api/client/errors.go`, `internal/api/client/client.go`,
  `internal/api/client/methods.go`, and `internal/api/client/events.go`.
- [x] T006 Implement the mutex-owned incident model, copy selection, deduplication,
  and bounded backoff helpers in `gui/connection.go`.
- [x] T007 Add read-only Windows diagnosis and portable fallback behind build
  tags in `gui/access_diagnosis_windows.go` and `gui/access_diagnosis_other.go`.

## Phase 3: User Story 1 - Recover Unauthorized Session (Priority: P1)

**Goal**: Keep one reachable access-denied recovery state and restore after the
Windows login token is refreshed.

**Independent Test**: Concurrent startup denial yields one panel; verified
stale-token evidence selects sign-out guidance; a successful retry clears it.

- [x] T008 [US1] Add verified stale-token, absent-membership, and unknown-cause
  diagnosis tests in `gui/connection_test.go` and Windows build tests.
- [x] T009 [US1] Render the persistent Retry/Exit connection panel above normal
  tabs and wire incident updates in `gui/app.go` and `gui/connection.go`.
- [x] T010 [US1] Route model and schedule startup transport failures through the
  incident coordinator in `gui/app.go` and `gui/schedule.go`.
- [x] T011 [US1] Implement successful Retry clearing and application Exit
  cancellation in `gui/app.go` and `gui/connection.go`.

## Phase 4: User Story 2 - Accurate Failure Guidance (Priority: P1)

**Goal**: Distinguish unavailable, access denied, timeout, other transport, and
API response errors.

**Independent Test**: Each injected category produces its own actionable copy;
API errors remain operation errors; denial never says only to check the daemon.

- [x] T012 [US2] Complete table-driven wrapped-error and copy tests in
  `internal/api/client/errors_test.go` and `gui/connection_test.go`.
- [x] T013 [US2] Keep unrelated API/operation errors on the existing error path
  while suppressing only duplicate connectivity errors in `gui/util.go` and
  `gui/connection.go`.

## Phase 5: User Story 3 - Stable Background Recovery (Priority: P2)

**Goal**: Maintain one cancelable bounded reconnect loop while disconnected.

**Independent Test**: Injected waits observe 2, 4, 8, 16, 30, 30 seconds,
immediate Retry interruption, one stream loop, and cancellation on exit.

- [x] T014 [US3] Add deterministic backoff, immediate-retry, and cancellation
  tests in `gui/connection_test.go`.
- [x] T015 [US3] Replace recursive fixed-delay stream refresh with the single
  reconnect coordinator in `gui/app.go` and `gui/connection.go`.

## Phase 6: Documentation, Analysis, and Verification

- [x] T016 Update fresh-install diagnosis and recovery walkthrough in the pinned
  artifact `docs/INSTALL-windows.md`.
- [x] T017 Record the architecture decision, pinned-document change, and user
  fix in `CHANGELOG.md`.
- [x] T018 Run blocking cross-artifact analysis and resolve every critical or
  high finding across S036 `spec.md`, `plan.md`, and `tasks.md`.
- [x] T019 Run focused client and headless GUI tests, including Windows compile
  coverage, and record red/green evidence in `verification.md`.
- [x] T020 Run native Windows diagnosis plus deterministic authorized recovery
  and record service, group, membership, token, pipe, guidance, Retry, and
  usable-GUI evidence in `verification.md`.
- [x] T021 Run `sh scripts/verify.sh all` and record all eight gates in
  `verification.md`.
- [x] T022 Mark S036 Implemented with objective delivery evidence in `spec.md`,
  `tasks.md`, and `specs/README.md`; run lifecycle and UTF-8/mojibake checks.

## Dependencies & Execution Order

- Phase 1 establishes the test-first baseline.
- Phase 2 blocks all user stories.
- US1 and US2 share the incident surface and proceed in that order; US3 depends
  on their coordinated retry semantics.
- Documentation and full verification follow all implementation stories.

## Parallel Opportunities

- T003 and T004 affect separate packages and can be authored independently.
- T005 and T007 affect separate packages after their tests exist.
- T016 can be drafted while focused tests run, but its evidence cannot be
  finalized before the native walkthrough.

## Implementation Strategy

Deliver the typed failure boundary first, then the one in-frame incident as the
P1 recovery increment, then bounded background recovery. Preserve traceability
to every issue #90 criterion and do not treat publication or PR review as an
unchecked implementation task.
