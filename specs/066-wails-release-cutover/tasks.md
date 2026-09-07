# Tasks: Wails Release Cutover

**Input**: Design documents from `specs/066-wails-release-cutover/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/release-candidate.md, quickstart.md

**Tests**: Required by the specification and constitution. Contract tests are authored before their workflow and removal changes.

## Phase 1: Setup and Inventory

**Purpose**: Freeze the candidate boundary and identify every maintained legacy dependency before deletion.

- [x] T001 Record the complete legacy-to-Wails capability inventory in `specs/066-wails-release-cutover/parity.md`
- [x] T002 [P] Inventory production package, CI, dependency, documentation, and issue-template Fyne references across `.github/`, `build/`, `cmd/`, `docs/`, `gui/`, `scripts/`, `test/`, `README.md`, `go.mod`, and `go.sum`
- [x] T003 [P] Define supported package payloads and exact-candidate gates in `specs/066-wails-release-cutover/contracts/release-candidate.md`
- [x] T004 Add S066 to the lifecycle index in `specs/README.md`

---

## Phase 2: Foundational Release Contracts

**Purpose**: Add failure-closed tests before changing the production desktop and release paths.

- [x] T005 Add failing Wails cutover source, dependency, naming, workflow, and documentation contract tests in `test/integration/desktop_cutover_contract_test.go`
- [x] T006 Update failing Windows release workflow expectations for native Wails packaging in `test/integration/windows_installer_contract_test.go`
- [x] T007 Update failing automation-policy fixtures and assertions for production-only Wails gates in `test/scripts/automation-check_test.sh`
- [x] T008 Extend the automation checker with production desktop, release payload, and retired-implementation residue requirements in `scripts/automation-check.sh`

**Checkpoint**: Cutover contracts fail against the Fyne shipping path and describe the replacement precisely.

---

## Phase 3: User Story 1 - Install the complete Wails desktop (Priority: P1)

**Goal**: Produce supported desktop distributions containing the native Wails app and current companion payloads.

**Independent Test**: Native workflow builds and static package-contract tests prove each supported platform payload and stable application identity.

- [x] T009 [US1] Configure the production Wails output and application metadata for the stable launch contract in `desktop/wails.json`
- [x] T010 [US1] Replace the release GUI build with native Wails dependency, build, and staging steps in `.github/workflows/release.yml`
- [x] T011 [US1] Preserve daemon, CLI, branding, macOS bundle, Linux integration, and Windows MSI payload assembly in `.github/workflows/release.yml`
- [x] T012 [US1] Update the Windows MSI process, file, and shortcut comments while preserving stable executable targets in `build/windows/goschedule.wxs`
- [x] T013 [US1] Replace root Fyne CI with production Wails verification and package-contract coverage in `.github/workflows/ci.yml`
- [x] T014 [US1] Validate release and installer contract tests in `test/integration/desktop_cutover_contract_test.go` and `test/integration/windows_installer_contract_test.go`

**Checkpoint**: Every supported production package sources its desktop entry from `desktop/` and rejects incomplete payloads.

---

## Phase 4: User Story 2 - Upgrade without losing operation or data (Priority: P1)

**Goal**: Preserve upgrade, repair, uninstall, task data, and bounded preference migration behavior.

**Independent Test**: Existing MSI source, compiled-package, silent-lifecycle, cleanup, IPC, and settings migration suites remain green against Wails staging.

- [x] T015 [US2] Extend Windows source-contract tests for Wails replacement at the stable installed path in `test/integration/windows_installer_contract_test.go`
- [x] T016 [US2] Update Windows CI MSI staging to build the Wails desktop and root companions in `.github/workflows/ci.yml`
- [x] T017 [US2] Retain and run preference migration, daemon storage, cleanup, service, shortcut, PATH, and uninstall lifecycle tests across `desktop/settings/`, `test/integration/`, and `test/windows/`
- [x] T018 [US2] Document the Fyne-to-Wails upgrade and preference boundary in `docs/install.md`, `docs/INSTALL-windows.md`, `docs/INSTALL-macos.md`, and `docs/INSTALL-linux.md`

**Checkpoint**: Existing installations upgrade at the same external identity without changing daemon-owned data or cleanup policy.

---

## Phase 5: User Story 3 - Qualify one auditable release candidate (Priority: P1)

**Goal**: Remove retired implementations only after the production replacement and gates are defined.

**Independent Test**: Canonical verification, desktop verification, hosted native matrix, and residue contracts pass for one commit.

- [x] T019 [US3] Replace the canonical Fyne GUI gate with production desktop Go, frontend, and build verification in `scripts/verify.sh`
- [x] T020 [US3] Remove the retired Fyne source and entry point from `gui/` and `cmd/gosched-gui/`
- [x] T021 [US3] Remove Fyne dependencies and tidy the root module in `go.mod` and `go.sum`
- [x] T022 [US3] Remove the superseded S060 executable proof from `experiments/wails-foundation/` and its CI policy requirements
- [x] T023 [US3] Update Dependabot grouping for the root and desktop module dependency boundaries in `.github/dependabot.yml`
- [x] T024 [US3] Complete the parity and intentional-deviation disposition record in `specs/066-wails-release-cutover/parity.md`

**Checkpoint**: One production desktop remains, every parity entry is explained, and no current build depends on Fyne.

---

## Phase 6: User Story 4 - Follow current desktop guidance (Priority: P2)

**Goal**: Make every maintained user and contributor surface describe the Wails desktop and current verification path.

**Independent Test**: Documentation and residue checks pass with remaining Fyne mentions limited to history or migration.

- [x] T025 [P] [US4] Update desktop architecture and build guidance in `README.md`, `desktop/README.md`, `CONTRIBUTING.md`, and `docs/build-autopilot.md`
- [x] T026 [P] [US4] Update user desktop, CLI, field, installation, and architecture guidance across `docs/`
- [x] T027 [P] [US4] Update desktop choices in `.github/ISSUE_TEMPLATE/bug_report.yml`, `.github/ISSUE_TEMPLATE/feature_request.yml`, and `.github/ISSUE_TEMPLATE/config.yml`
- [x] T028 [US4] Classify retained historical and migration Fyne references and enforce current-product wording in `test/integration/desktop_cutover_contract_test.go`

**Checkpoint**: Current guidance has one production desktop story and historical references remain accurate.

---

## Phase 7: Polish and Cross-Cutting Concerns

- [x] T029 Add the S066 feature and dated architecture decision to `CHANGELOG.md`
- [x] T030 Reconcile the completed S065 branch-cleanup note and S066 lifecycle row in `specs/README.md`
- [x] T031 Run focused root, desktop, frontend, browser, installer-source, documentation, automation, residue, and lifecycle tests
- [x] T032 Run native Windows Wails build and inspect the stable output name
- [x] T033 Run `sh scripts/verify.sh all` in the foreground and record every gate in `specs/066-wails-release-cutover/verification.md`
- [x] T034 Run `go run ./scripts/github-format` and UTF-8/mojibake scans over changed publication content
- [x] T035 Mark all tasks complete, set the specification to Implemented, and re-run read-only cross-artifact analysis

---

## Dependencies and Execution Order

- Setup establishes the inventory before removal.
- Foundational contract tests must fail against the old shipping boundary before workflow changes.
- User Story 1 establishes package construction before User Story 2 validates lifecycle preservation.
- User Story 3 may remove Fyne and the proof only after User Stories 1 and 2 have executable gates.
- User Story 4 follows the final file and command names from User Story 3.
- Final verification depends on every story and blocks publication on any red or unavailable required local gate.

## Parallel Opportunities

- T002 and T003 inspect independent surfaces.
- T025, T026, and T027 update distinct documentation and issue-template files after implementation names settle.
- Root static contract tests, desktop frontend tests, and desktop Go tests can be run independently before the canonical sequential verifier.

## Implementation Strategy

The minimum safe increment is User Stories 1 through 3 together because production packaging, lifecycle preservation, and old-implementation removal form one atomic cutover. User Story 4 completes the public contract in the same slice. Do not ship a state where release automation names Wails but current documentation or root dependencies still maintain Fyne.
