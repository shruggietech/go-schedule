# Tasks: Dedicated IPC Administrative Group

**Input**: Design documents from `specs/030-ipc-admin-group/` **Tests**: Required by constitution Principle II and the safety-critical local IPC access-control contract.

## Phase 1: Setup and Baseline

- [x] T001 Record pre-change config, IPC, installer, toolchain, and documentation evidence in `specs/030-ipc-admin-group/verification.md`.
- [x] T002 Update the active feature context and lifecycle inventory in `CLAUDE.md`, `.specify/feature.json`, and `specs/README.md`.

## Phase 2: Foundational Contracts

**Purpose**: Establish shared configuration and access-policy types before platform implementations.

- [x] T003 Add failing `admin_group` default, explicit-empty, and whitespace validation tests in `internal/config/config_test.go`.
- [x] T004 [P] Add failing access-mode evidence contract tests in `internal/ipc/ipc_test.go`.
- [x] T005 Implement `admin_group` validation in `internal/config/config.go` until T003 passes.
- [x] T006 Implement the shared listener/access-info contract in `internal/ipc/ipc.go` until T004 passes.
- [x] T007 Run the blocking Spec-Kit analysis and resolve every consistency, coverage, and constitution finding across `specs/030-ipc-admin-group/`.

**Checkpoint**: Configuration distinguishes secure, invalid, and explicit compatibility values; platform implementations can return consistent evidence.

## Phase 3: User Story 1 - Restrict Privileged Scheduling Access (Priority: P1)

**Goal**: Apply and verify a fail-closed group policy before serving on Unix or Windows.

**Independent Test**: Platform unit tests prove allowed identities, resolved group use, exact Unix modes/ownership, descriptor construction, and cleanup after partial failure.

- [x] T008 [P] [US1] Add failing restricted Unix lookup, permission, readback, stale-endpoint, and cleanup tests in `internal/ipc/ipc_unix_test.go`.
- [x] T009 [P] [US1] Add failing restricted Windows SID-type, descriptor, and missing-group tests in `internal/ipc/ipc_windows_test.go`.
- [x] T010 [P] [US1] Add failing restricted startup logging tests in `cmd/goschedd/main_test.go`.
- [x] T011 [US1] Implement restricted Unix group resolution, parent/socket permission application, readback, and cleanup in `internal/ipc/ipc_unix.go`.
- [x] T012 [US1] Implement restricted Windows group SID resolution and protected descriptor construction in `internal/ipc/ipc_windows.go`.
- [x] T013 [US1] Wire configured group input and structured restricted-policy logging through `cmd/goschedd/main.go`.

**Checkpoint**: A valid group yields a verified restricted listener; any non-empty lookup or permission error fails before readiness.

## Phase 4: User Story 2 - Choose Compatibility Mode Deliberately (Priority: P1)

**Goal**: Preserve broad local access only for an explicit empty configuration value and make the choice observable.

**Independent Test**: Empty-group platform tests select the historical broad policy and daemon tests capture exactly one warning; non-empty failures never transition to compatibility.

- [x] T014 [P] [US2] Add failing Unix and Windows compatibility-policy tests in `internal/ipc/ipc_unix_test.go` and `internal/ipc/ipc_windows_test.go`.
- [x] T015 [P] [US2] Add failing compatibility warning and ordering tests in `cmd/goschedd/main_test.go`.
- [x] T016 [US2] Implement explicit compatibility modes and access evidence in `internal/ipc/ipc_unix.go` and `internal/ipc/ipc_windows.go`.
- [x] T017 [US2] Implement structured compatibility warning behavior in `cmd/goschedd/main.go`.

**Checkpoint**: Empty means observable compatibility; every non-empty failure remains fail closed.

## Phase 5: User Story 3 - Install into a Working Secure Default (Priority: P2)

**Goal**: Provision Windows automatically and make Unix service setup clear before first start.

**Independent Test**: Installer XML contracts accept the exact group/membership lifecycle and reject broken variants; documentation integrity passes with group setup ordered before service start.

- [x] T018 [US3] Add failing group creation, installing-user membership, preservation, sequencing, and WiX pin assertions in `test/integration/windows_installer_contract_test.go` and `build/windows/verify_wxs.ps1`.
- [x] T019 [US3] Upgrade the pinned tool and extensions to WiX 6.0.2 in `.github/workflows/release.yml` and implement declarative group/membership provisioning in `build/windows/goschedule.wxs`.
- [x] T020 [P] [US3] Document secure default, compatibility configuration, and troubleshooting in `docs/INSTALL-windows.md`, `docs/INSTALL-linux.md`, and `docs/INSTALL-macos.md`.
- [x] T021 [P] [US3] Replace the obsolete permissive-boundary description in `SECURITY.md` and summarize the administrative-group boundary in `README.md`.

**Checkpoint**: Windows has an idempotent installer contract and all platforms have actionable pre-service provisioning instructions.

## Phase 6: Cross-Cutting Verification and Completion

- [x] T022 Record the fail-closed, Unix-directory, and WiX 6.0.2 architecture decisions in the dated `CHANGELOG.md` Unreleased section.
- [x] T023 Run focused config, IPC, daemon, installer, cross-compilation, PowerShell, and WiX build checks; record test-first and focused results in `specs/030-ipc-admin-group/verification.md`.
- [x] T024 Run `sh scripts/verify.sh all` in the foreground through format, vet, lint, race, GUI, coverage, docs, and automation.
- [x] T025 Audit the diff, UTF-8 without BOM, mojibake, trailing whitespace, pinned-artifact decision evidence, and `Closes #13` language; finalize `specs/030-ipc-admin-group/verification.md`.
- [x] T026 Mark required tasks resolved and advance `specs/030-ipc-admin-group/spec.md` plus `specs/README.md` to Implemented with local review-branch delivery evidence.

## Dependencies and Execution Order

- Phase 2 blocks both runtime user stories.
- US1 establishes restricted platform construction before US2 adds the explicit broad alternative.
- US3 depends on US1's default group contract but is independently testable through installer and documentation contracts.
- T022 through T026 require all three user stories to be complete.

## Parallel Opportunities

- T003 and T004 touch separate packages.
- T007, T008, and T009 establish independent Unix, Windows, and daemon failures before implementation.
- T013 and T014 cover separate platform and logging surfaces.
- T019 and T020 update separate documentation groups after the runtime contract is stable.

## Implementation Strategy

Implement one substantial security slice in dependency order: shared contract, restricted Unix and Windows enforcement, explicit compatibility mode, installer provisioning, operator documentation, then complete local verification. No publication bookkeeping appears as an implementation task; autopilot halts after the local review-ready commit.
