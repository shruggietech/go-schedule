# Tasks: Authenticated Remote Access and Pairing

**Input**: Design documents from `specs/077-remote-https-pairing/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Security, migration, contract, lifecycle, concurrency, and cross-platform behavior require test-first coverage under the constitution.

## Phase 1: Setup

- [x] T001 Record direct dependency pins and licenses in `go.mod`, `go.sum`, and `specs/077-remote-https-pairing/research.md`
- [x] T002 [P] Add the OpenAPI 3.1 source and generation contract in `api/openapi/remote-v1.yaml` and `api/openapi/remote-v1.cfg.yaml`
- [x] T003 [P] Add S077 inventory and unreleased decision records in `specs/README.md` and `CHANGELOG.md`

## Phase 2: Foundational Storage and Configuration

- [x] T004 [P] Write remote configuration boundary tests in `internal/config/config_test.go`
- [x] T005 Implement default-disabled listener configuration and pre-bind validation in `internal/config/config.go`
- [x] T006 [P] Write pairing and credential domain validation tests in `internal/domain/remote_access_test.go`
- [x] T007 Implement pairing and credential domain types plus secret-safe response types in `internal/domain/remote_access.go`
- [x] T008 Write schema-v18 migration, CRUD, and concurrency tests in `internal/store/migration_v18_test.go` and `internal/store/remote_access_test.go`
- [x] T009 Implement digest-only credential and Argon2id pairing persistence in `internal/store/store.go` and `internal/store/remote_access.go`

## Phase 3: User Story 1 - Safe Opt-in HTTPS Boundary

**Goal**: Serve no TCP port by default and one bounded TLS 1.3 remote adapter when explicitly configured.

**Independent Test**: Exercise disabled and invalid configuration, TLS negotiation, plaintext refusal, exact route admission, local IPC coexistence, and graceful shutdown.

- [x] T010 [US1] Write remote listener and route-admission tests in `internal/remote/server_test.go`
- [x] T011 [US1] Implement the explicit remote operation table and `/api/v1` path adapter in `internal/remote/operations.go`
- [x] T012 [US1] Implement TLS lifecycle, request bounds, concurrency bounds, origin rejection, and shutdown in `internal/remote/server.go`
- [x] T013 [US1] Wire the optional listener beside local IPC in `cmd/goschedd/main.go` and cover lifecycle behavior in `cmd/goschedd/main_test.go`

## Phase 4: User Story 2 - One-time Named Client Pairing

**Goal**: Create a local phrase and atomically exchange it once over the intended HTTPS daemon.

**Independent Test**: Prove success plus wrong-daemon, malformed, replay, expiry, cancellation, exhaustion, and concurrent exchange behavior.

- [x] T014 [US2] Write entropy, verifier, lifecycle, and concurrent-exchange tests in `internal/enrollment/service_test.go`
- [x] T015 [US2] Implement cryptographic phrase generation, verification, credential issuance, and atomic exchange in `internal/enrollment/service.go`
- [x] T016 [US2] Add audited local pairing administration handlers and catalog records in `internal/api/server/access.go`, `internal/api/server/server.go`, and `internal/authorization/catalog.go`
- [x] T017 [US2] Add the isolated generic-failure remote enrollment handler in `internal/remote/server.go`
- [x] T018 [US2] Add typed local-client and CLI pairing administration in `internal/api/client/access.go`, `internal/cli/access.go`, and their tests
- [x] T019 [US2] Add the bounded desktop pairing bridge, native-storage handoff, form, and tests in `desktop/app.go`, `desktop/frontend/src`, and `internal/clientsecret`

## Phase 5: User Story 3 - Durable Authority Lifecycle

**Goal**: Authenticate every protected request through current credential and actor state and support safe rotation and revocation.

**Independent Test**: Exercise the credential-state and four-capability matrices, rotate and revoke values, and inspect secret-free audit results.

- [x] T020 [US3] Write bearer parsing, generic failure, current-state reload, capability, rotation, revocation, and secret-canary tests in `internal/remote/auth_test.go`
- [x] T021 [US3] Implement strict bearer authentication and request-scoped actor resolution in `internal/remote/auth.go`
- [x] T022 [US3] Implement bounded source and actor rate limiting in `internal/remote/limiter.go` with deterministic tests in `internal/remote/limiter_test.go`
- [x] T023 [US3] Add audited local credential list, rotation, and revocation operations across `internal/api/server`, `internal/api/client`, `internal/cli`, and `internal/authorization`
- [x] T024 [P] [US3] Add the native credential-store adapter and injected tests in `internal/clientsecret/keyring.go` and `internal/clientsecret/keyring_test.go`
- [x] T025 [US3] Revalidate actor and credential authority during remote event streams in `internal/remote/server.go`

## Phase 6: User Story 4 - Stable Machine Contract

**Goal**: Keep OpenAPI, runtime allowlisting, bounds, errors, and local equivalence mechanically aligned.

**Independent Test**: Compare all route and operation identities, reject excluded local paths, and exercise every stable error class.

- [x] T026 [US4] Complete OpenAPI schemas, operations, security declarations, and stable errors in `api/openapi/remote-v1.yaml`
- [x] T027 [US4] Add OpenAPI and runtime completeness tests in `internal/remote/contract_test.go`
- [x] T028 [US4] Add HTTPS end-to-end integration tests in `test/integration/remote_access_test.go`

## Phase 7: Documentation and Completion

- [x] T029 [P] Document operator enablement, trust, pairing, rotation, revocation, and deployment examples in `docs/remote-access.md`, `docs/api.md`, and `docs/cli.md`
- [x] T030 [P] Update architecture and security boundaries in `docs/architecture.md`, `docs/access-control.md`, and `CHANGELOG.md`
- [x] T031 Run focused tests, clean contract generation, dependency/license review, and secret/mojibake scans; record results in `specs/077-remote-https-pairing/verification.md`
- [x] T032 Run `sh scripts/verify.sh all`, resolve failures, and record all eight gates in `specs/077-remote-https-pairing/verification.md`
- [x] T033 Mark the spec Implemented, complete this task list, run `/speckit-analyze` closure checks, and commit the review branch

## Dependencies and Execution Order

- Setup precedes storage and configuration.
- T004 to T009 establish the configuration, domain, and transaction boundaries required by every story.
- User Story 1 provides the listener and explicit route adapter.
- User Story 2 depends on the listener and store transaction.
- User Story 3 depends on issued credentials and completes normal protected requests.
- User Story 4 validates the combined reachable surface.
- Documentation and canonical verification follow all stories.

## Parallel Opportunities

- T002 and T003 can proceed independently after T001.
- Configuration tests and domain tests affect separate packages.
- Native credential-store work is independent of server authentication once the credential response is stable.
- Operator, API, and architecture documentation can be authored independently after contracts stabilize.

## Implementation Strategy

Implement test-first in dependency order: configuration and schema, TLS adapter, phrase exchange, bearer lifecycle, rate limits, contract completeness, integration evidence, then documentation and full verification. The independently useful MVP is User Story 1 plus User Story 2 because it establishes one authenticated remote relationship without a temporary credential mechanism.
