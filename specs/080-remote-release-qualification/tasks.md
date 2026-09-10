# Tasks: Remote Access Release Qualification

**Input**: Design documents from `specs/080-remote-release-qualification/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/qualification.md`, `quickstart.md`

**Tests**: Required by the constitution and issue #173. Behavioral changes begin with failing focused tests.

## Phase 1: Specification and Contracts

- [x] T001 Register S080 and its Draft delivery state in `.specify/feature.json` and `specs/README.md`
- [x] T002 [P] Validate specification quality in `specs/080-remote-release-qualification/checklists/requirements.md` and `specs/080-remote-release-qualification/checklists/release.md`
- [x] T003 [P] Finalize research, data, qualification, and verification contracts in `specs/080-remote-release-qualification/research.md`, `specs/080-remote-release-qualification/data-model.md`, `specs/080-remote-release-qualification/contracts/qualification.md`, and `specs/080-remote-release-qualification/quickstart.md`

## Phase 2: Service Configuration Foundation

- [x] T004 Add failing service-install configuration tests in `internal/cli/service_test.go`
- [x] T005 Implement validated absolute `service install --config` argument binding in `internal/cli/service.go`

## Phase 3: User Story 1 - Operate a Remote Headless Daemon (Priority: P1)

**Goal**: Prove the chronological package-shaped lifecycle on every supported operating system.

**Independent Test**: Run one real daemon through default-off, enable, pairing, use, denial, audit, revocation, disable, and replacement-binary restarts.

- [x] T006 [US1] Add a failing package-shaped v1.4 lifecycle test in `test/integration/v14_remote_release_qualification_test.go`
- [x] T007 [US1] Complete the daemon lifecycle harness and exact secret-free assertions in `test/integration/v14_remote_release_qualification_test.go`
- [x] T008 [US1] Add the named three-platform qualification matrix to `.github/workflows/ci.yml`

## Phase 4: User Story 2 - Choose a Supported Deployment Safely (Priority: P1)

**Goal**: Publish one clear operational path from installation to recovery without unsafe zero-configuration claims.

**Independent Test**: Run the documentation contract and manually map every lifecycle and deployment responsibility to a visible section.

- [x] T009 [P] [US2] Add service configuration syntax to `docs/cli.md`
- [x] T010 [US2] Expand enablement, deployment, pairing, revocation, disablement, upgrade, backup, and recovery guidance in `docs/remote-access.md`
- [x] T011 [P] [US2] Link platform-specific service configuration and remote guidance from `docs/INSTALL-linux.md`, `docs/INSTALL-macos.md`, and `docs/INSTALL-windows.md`

## Phase 5: User Story 3 - Review Release Evidence (Priority: P1)

**Goal**: Map every release criterion to deterministic and hosted evidence without publishing artifacts.

**Independent Test**: Run focused and canonical gates, inspect the cross-platform workflow, and complete the source-bound evidence matrix.

- [x] T012 [US3] Add the S080 feature and service-configuration decision to `CHANGELOG.md`
- [x] T013 [US3] Run read-only spec-kit analysis and resolve all critical or high findings across `specs/080-remote-release-qualification/spec.md`, `specs/080-remote-release-qualification/plan.md`, and `specs/080-remote-release-qualification/tasks.md`
- [x] T014 [US3] Run focused race, lifecycle, documentation, formatting, secret-canary, and diff checks from `specs/080-remote-release-qualification/quickstart.md`
- [x] T015 [US3] Run `sh scripts/verify.sh all` and record all eight gates in `specs/080-remote-release-qualification/verification.md`
- [x] T016 [US3] Advance S080 to Implemented, update `specs/README.md`, complete all tasks, and record delivery evidence in `specs/080-remote-release-qualification/spec.md`
- [x] T017 [US3] Run lifecycle, mojibake, publication-format, and repository-cleanliness checks before commit

## Dependencies and Execution Order

- Phase 1 establishes the executable contract.
- Phase 2 blocks service-based headless enablement and documentation.
- Phase 3 composes the shipped remote journey and blocks evidence closure.
- Phase 4 may proceed after Phase 2 fixes the exact command syntax.
- Phase 5 follows all implementation and documentation work.

## Parallel Opportunities

- T002 and T003 affect independent specification artifacts.
- T009 and T011 affect independent documentation files after T005 fixes command syntax.
- Focused package groups in T014 may be invoked separately, while the canonical aggregate remains one foreground sequential run.

## Implementation Strategy

1. Make registered daemon configuration explicit and deterministic.
2. Compose existing remote capabilities through a package-shaped lifecycle rather than duplicating their unit tests.
3. Add one named three-platform CI surface for the release gate.
4. Publish the operational contract and close with exact evidence.

## Done Gate

S080 is implemented only when the service retains a validated configuration path, the complete remote lifecycle is reproducible, every issue #173 criterion maps to passing evidence or published guidance, all eight canonical gates pass, and no tag or release is created.
