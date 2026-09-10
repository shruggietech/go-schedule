# Tasks: Remote Connection Profiles and Target-Safe Clients

**Input**: Design documents from `specs/078-remote-connection-profiles/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Profile persistence, target isolation, credential lifecycle, transport security, CLI compatibility, and accessibility require test-first coverage under the constitution.

## Phase 1: Setup

- [x] T001 Add S078 inventory and architecture decisions in `specs/README.md`, `CLAUDE.md`, and `CHANGELOG.md`
- [x] T002 [P] Define profile and selection contracts in `specs/078-remote-connection-profiles/contracts/profiles.md`
- [x] T003 [P] Add verification placeholders in `specs/078-remote-connection-profiles/verification.md`

## Phase 2: Shared Profile and Transport Foundation

- [x] T004 Write validation, secret-exclusion, corruption, forward-version, atomicity, and concurrency tests in `internal/clientprofile/profile_test.go` and `internal/clientprofile/store_test.go`
- [x] T005 Implement the versioned profile model, canonical validation, fingerprinting, locking, and atomic persistence in `internal/clientprofile/profile.go` and `internal/clientprofile/store.go`
- [x] T006 Write local and remote transport tests for TLS, redirect refusal, bearer headers, path mapping, status classification, and identity pinning in `internal/api/client/client_test.go` and `internal/api/client/remote_test.go`
- [x] T007 Implement immutable local and authenticated remote constructors plus stable status errors in `internal/api/client/client.go`, `internal/api/client/methods.go`, `internal/api/client/events.go`, and `internal/api/client/access.go`
- [x] T008 Implement a synchronized target provider that prevents in-flight client retargeting in `internal/api/client/client.go` with race tests in `internal/api/client/remote_test.go`

## Phase 3: User Story 1 - Persistent Desktop Selection

**Goal**: Pair, persist, restart, and select exact remote daemon profiles while retaining This computer.

**Independent Test**: Pair two same-named daemons, restart, select each stable profile, and verify global identity plus request routing.

- [x] T009 [US1] Write remote backend health, manifest identity, authority, and safe-failure tests in `desktop/connection/remote_test.go`
- [x] T010 [US1] Implement remote health and event adaptation in `desktop/connection/remote.go`
- [x] T011 [US1] Write selection cancellation, stale-generation, and target-publication tests in `desktop/connection/manager_test.go`
- [x] T012 [US1] Extend the connection manager with safe backend switching in `desktop/connection/manager.go` and target-specific messages in `desktop/connection/model.go`
- [x] T013 [US1] Persist a profile after successful desktop pairing and preserve secret rollback semantics in `desktop/remotepairing/service.go` and `desktop/remotepairing/service_test.go`
- [x] T014 [US1] Wire shared provider, profile store, keyring, and startup selection in `desktop/main.go` and `desktop/app.go`
- [x] T015 [US1] Add frontend profile bridge and model contracts in `desktop/frontend/src/connection/bridge.ts` and `desktop/frontend/src/connection/model.ts`

## Phase 4: User Story 2 - Desktop Profile Lifecycle and Target Safety

**Goal**: List, rename, repair, and remove profiles while making target and unsupported authority unmistakable.

**Independent Test**: Exercise lifecycle failure order, same-name disambiguation, active removal, capability gating, and accessible mutation confirmation.

- [x] T016 [US2] Write desktop profile lifecycle service tests in `desktop/connections/service_test.go`
- [x] T017 [US2] Implement secret-free workspaces and select, rename, repair handoff, and remove orchestration in `desktop/connections/service.go` and `desktop/connections/model.go`
- [x] T018 [US2] Expose connection lifecycle methods and target-aware action results through `desktop/app.go`
- [x] T019 [US2] Write Connections view, shell target, disabled-route, keyboard, and confirmation tests in `desktop/frontend/src/settings/ConnectionsPage.test.tsx`, `desktop/frontend/src/App.test.tsx`, and related component tests
- [x] T020 [US2] Implement profile cards, selection, rename, removal confirmation, repair mode, and target-aware shell copy in `desktop/frontend/src/settings/ConnectionsPage.tsx`, `desktop/frontend/src/remotepairing/PairingForm.tsx`, `desktop/frontend/src/components/Shell.tsx`, and `desktop/frontend/src/App.tsx`
- [x] T021 [US2] Add capability and authority gating without local fallback in `desktop/frontend/src/App.tsx` and feature action surfaces

## Phase 5: User Story 3 - Explicit CLI Remote Profiles

**Goal**: Pair and administer profiles, then target one remote daemon per deliberate CLI invocation.

**Independent Test**: Compare no-flag local behavior, named profile behavior, explicit selection, human diagnostics, JSON stdout, and every invalid selection class.

- [x] T022 [US3] Write global target-resolution and local-compatibility tests in `internal/cli/cli_test.go` and `internal/cli/profile_test.go`
- [x] T023 [US3] Implement profile and explicit endpoint global flags with immutable client resolution in `internal/cli/cli.go` and `internal/cli/target.go`
- [x] T024 [US3] Write profile pair and lifecycle command tests including protected stdin and secret canaries in `internal/cli/profile_test.go`
- [x] T025 [US3] Implement profile pair, list, show, rename, and remove commands in `internal/cli/profile.go`
- [x] T026 [US3] Add remote target human diagnostics while preserving JSON stdout and exit codes in `internal/cli/cli.go`

## Phase 6: User Story 4 - Documented JSON Clients

**Goal**: Publish safe executable remote-client workflows without encouraging credential or mutation-replay hazards.

**Independent Test**: Run private HTTPS and SSH tunnel enrollment and authenticated-read examples from a clean disposable environment.

- [x] T027 [P] [US4] Document named and explicit CLI profiles in `docs/cli.md`
- [x] T028 [P] [US4] Document enrollment, identity pinning, protected bearer input, pagination, errors, and mutation uncertainty in `docs/api.md`
- [x] T029 [P] [US4] Update operator profile, repair, removal, private HTTPS, and SSH tunnel workflows in `docs/remote-access.md`
- [x] T030 [US4] Update architecture and access-control boundaries in `docs/architecture.md`, `docs/access-control.md`, and `CHANGELOG.md`

## Phase 7: Completion

- [x] T031 Run `/speckit-analyze`, resolve every finding, and record the coverage result in `specs/078-remote-connection-profiles/verification.md`
- [x] T032 Run focused Go race, frontend, target-isolation, secret-canary, and documentation checks; record results in `specs/078-remote-connection-profiles/verification.md`
- [x] T033 Run `sh scripts/verify.sh all`, resolve failures, and record all eight gates in `specs/078-remote-connection-profiles/verification.md`
- [x] T034 Run UTF-8, mojibake, secret, generated-binding, and GitHub-format checks; mark the spec Implemented and every task complete
- [x] T035 Commit the verified S078 review branch with issue traceability

## Dependencies and Execution Order

- Setup precedes the shared profile and immutable transport foundation.
- T004 through T008 block every user story.
- User Story 1 establishes selection and request routing.
- User Story 2 depends on User Story 1 and completes the desktop lifecycle.
- User Story 3 reuses the shared foundation and is independently testable after Phase 2.
- User Story 4 follows stable CLI and profile contracts.
- Completion follows all desired stories.

## Parallel Opportunities

- Contract and verification artifacts are independent after T001.
- Profile persistence and transport tests affect separate packages before provider integration.
- CLI profile commands can proceed after Phase 2 while desktop lifecycle work continues.
- CLI, API, and remote-access documentation affect separate files after contracts stabilize.

## Implementation Strategy

Implement test-first in dependency order: profile validation and atomic storage, immutable HTTPS clients, synchronized desktop provider, desktop selection and lifecycle, CLI selection and lifecycle, documentation, analysis, then full verification. The independently useful MVP is User Story 1, but #170 and #171 close only after every story and its verification evidence is complete.
