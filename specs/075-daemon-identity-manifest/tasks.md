# Tasks: Stable Daemon Identity and Capability Manifest

**Input**: Design documents from `specs/075-daemon-identity-manifest/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/local-manifest.md`, and both checklists

**Tests**: Required for every persistent, API, client, CLI, and desktop behavior in issue #166.

## Phase 1: Specification and Design

- [x] T001 [US1] Specify and clarify issue #166 in `specs/075-daemon-identity-manifest/spec.md`.
- [x] T002 [US1] Validate specification quality in `specs/075-daemon-identity-manifest/checklists/requirements.md`.
- [x] T003 [US3] Validate identity lifecycle, reset safety, privacy, and compatibility requirements in `specs/075-daemon-identity-manifest/checklists/identity-lifecycle.md`.
- [x] T004 [US1] Record identity, restore, clone, manifest, compatibility, reset, and desktop decisions in `specs/075-daemon-identity-manifest/research.md`.
- [x] T005 [US1] Define singleton identity and manifest projection boundaries in `specs/075-daemon-identity-manifest/data-model.md`.
- [x] T006 [US2] Define local manifest, rename, reset, privacy, error, and compatibility contracts in `specs/075-daemon-identity-manifest/contracts/local-manifest.md`.
- [x] T007 [US3] Document focused review and operator validation in `specs/075-daemon-identity-manifest/quickstart.md`.

## Phase 2: Analysis Gate

- [x] T008 Audit specification, plan, research, data model, contract, quickstart, checklists, and tasks through the spec-kit analyze gate; resolve every critical or high finding before implementation.

## Phase 3: Persistent Identity Foundation

- [x] T009 [US1] Write failing `internal/store/migration_v16_test.go` coverage for prior-schema scheduler preservation and exactly one initialized identity.
- [x] T010 [US1] Write failing `internal/store/identity_test.go` coverage for 100 fresh unique identities, restart and restore stability, singleton concurrency, and contextual persistence failures.
- [x] T011 [US3] Extend the store tests with failing rename validation and atomic exact-confirmation reset cases that preserve display name and scheduler data.
- [x] T012 [US1] Add validated `domain.DaemonIdentity`, schema migration v16, fail-closed initialization, and store lifecycle methods in `internal/domain/identity.go`, `internal/store/store.go`, and `internal/store/identity.go`.
- [x] T013 [US1] Run the focused domain and store tests and resolve every failure.

## Phase 4: Manifest API and Shared Client

- [x] T014 [US2] Write failing `internal/api/server/manifest_test.go` coverage for bounded deterministic fields, privacy exclusions, validation errors, reset conflict, and unchanged health behavior.
- [x] T015 [US2] Write failing `internal/api/client/client_test.go` coverage for typed manifest, rename, reset, and structured errors.
- [x] T016 [US2] Implement the manifest projection and local routes in `internal/api/server/manifest.go` and `internal/api/server/server.go`.
- [x] T017 [US2] Add typed shared-client manifest operations in `internal/api/client/client.go` and supporting client files.
- [x] T018 [US2] Run the focused API server and client tests and resolve every failure.

## Phase 5: Operator and Desktop Integration

- [x] T019 [US3] Write failing CLI tests for command registration, human and JSON manifest output, rename validation, required reset confirmation, and error propagation.
- [x] T020 [US2] Write failing desktop connection tests proving connected target identity, name, platform, version, and capabilities come from the daemon manifest.
- [x] T021 [US3] Implement `gosched daemon manifest`, `rename`, and `reset-identity` in `internal/cli/daemon.go` and register them in `internal/cli/cli.go`.
- [x] T022 [US2] Replace desktop identity and capability placeholders with daemon manifest values while retaining local permissions until #167.
- [x] T023 [US2] Run focused CLI and desktop connection tests and resolve every failure.

## Phase 6: Documentation, Lifecycle, and Focused Verification

- [x] T024 [US3] Author `docs/daemon-identity.md` with install, upgrade, restore, clone, rename, reset, privacy, and remote-stage behavior; link it from API and architecture documentation.
- [x] T025 [US2] Document the manifest wire contract and CLI commands in `docs/api.md` and current user documentation.
- [x] T026 [US1] Record the dated S075 identity decision and user-visible additions in `CHANGELOG.md`.
- [x] T027 Run focused migration, store, server, client, CLI, desktop, lifecycle, publication-format, and diff checks; resolve every failure.
- [x] T028 Advance the specification to `In Progress`, update the project item to S075 and In progress, and record focused evidence in `specs/075-daemon-identity-manifest/verification.md`.

## Phase 7: Canonical Verification and Delivery

- [x] T029 Run `sh scripts/verify.sh all` in the foreground with the repository baseline and resolve every failure.
- [x] T030 Mark every required task complete, advance the specification and inventory to `Implemented`, and record exact local delivery evidence.
- [x] T031 Run UTF-8, BOM, mojibake, Unicode em dash, Markdown wrapping, diff, and repository-status audits.
- [x] T032 Prepare the complete S075 slice for commit as `feat(075): add daemon identity manifest` with the required co-author trailer.

## Authorized Publication Runbook

After the verified commit, push the authorized branch, open a structured pull request that closes #166, move the project item to PR review, and process every CI and review result. At most one manual second `@Codex` review round may be requested. Do not request a third round. Once final-head CI is green and all review conversations are resolved, ask the maintainer to perform the final review and merge ritual.

## Dependencies and Execution Order

- Phase 1 completes specify, clarify, checklist, and plan work before implementation.
- T008 is the blocking analysis gate.
- T009 through T011 establish red evidence before T012 makes persistence green.
- T014 and T015 establish red evidence before T016 and T017 implement the API boundary.
- T019 and T020 establish red evidence before T021 and T022 integrate operator and desktop experiences.
- Documentation and lifecycle work depend on the settled implementation contract.
- Canonical verification and publication depend on every focused test passing.

## Scope Guard

- Close #166 only after its acceptance criteria and verification are complete.
- Do not implement actor persistence, permission enforcement, audit records, credentials, pairing, remote profiles, remote HTTPS, certificates, or network listeners from #167 onward.
- Do not disclose hostname, address, data path, account, environment, command, credential, or secret fields.
- Do not change the existing health response contract or remove existing local routes.
