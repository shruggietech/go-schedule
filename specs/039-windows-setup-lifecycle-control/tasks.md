# Tasks: Windows Setup Lifecycle Control

**Input**: Design documents from `/specs/039-windows-setup-lifecycle-control/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/windows-setup-lifecycle.md`, `quickstart.md`

**Tests**: Required by #97, #98, and constitution principle II. Behavioral tests are written and observed failing before production authoring.

## Phase 1: Specification and Baseline

- [x] T001 Complete the specification, clarification record, requirements checklists, plan, research, data model, lifecycle contract, and quickstart under `specs/039-windows-setup-lifecycle-control/`.
- [x] T002 Run focused green baseline checks for `build/windows/verify_wxs.ps1`, `test/integration/windows_installer_contract_test.go`, and changed PowerShell parser surfaces before adding S039 behavior.
- [x] T003 Run Spec Kit's read-only consistency analysis across `spec.md`, `plan.md`, and `tasks.md`; resolve every critical or high-severity finding before implementation.

---

## Phase 2: Foundational Cleanup Safety

**Purpose**: Establish the safe installer-private cleanup contract before any destructive MSI authoring.

- [x] T004 [US4] Add failing platform-neutral planner, preflight, traversal, and result-ledger tests in `internal/winuninstall/cleanup_test.go`.
- [x] T005 [US4] Add failing Windows-specific trusted-root, fixed-volume, canonicalization, reparse, and removal tests in `internal/winuninstall/path_windows_test.go` and `internal/winuninstall/remove_windows_test.go`.
- [x] T006 [US4] Implement platform-neutral cleanup planning, two-pass refusal/deletion, reparse-aware traversal interfaces, and atomic result records in `internal/winuninstall/cleanup.go`.
- [x] T007 [US4] Implement Windows known-folder/profile discovery and path-safety adapters in `internal/winuninstall/platform_windows.go`, with a non-Windows adapter in `internal/winuninstall/platform_other.go`.
- [x] T008 [US4] Add the installer-private `wipe` entry point and Windows GUI-subsystem build contract in `cmd/gosched-cleanup/` and `test/integration/windows_installer_contract_test.go`.
- [x] T009 [US4] Run focused cleanup tests with the race detector and cross-build the windowless Windows helper.

**Checkpoint**: Unsafe roots are refused before any deletion, safe roots can be removed without following reparse entries, and complete/partial/refused evidence is deterministic.

---

## Phase 3: User Story 1 - Choose Windows Shortcuts (Priority: P1)

**Goal**: Make Start Menu and desktop shortcuts independent native installer features with correct defaults, maintenance transitions, upgrade migration, and cleanup.

**Independent Test**: Source and built-MSI contracts prove feature/component/default relationships; the hosted silent matrix exercises all four selections, maintenance transitions, upgrade migration, and uninstall isolation.

- [x] T010 [US1] Add failing semantic and mutation tests for shortcut features, defaults, directories, component identities, and attributes in `test/integration/windows_installer_contract_test.go`.
- [x] T011 [US1] Author stable `StartMenuShortcut` and `DesktopShortcut` child features and shortcut components in `build/windows/goschedule.wxs`.
- [x] T012 [US1] Extend structural source validation for both shortcut features and components in `build/windows/verify_wxs.ps1`.
- [x] T013 [US1] Extend compiled-MSI feature, component, directory, and shortcut inspection in `test/windows/inspect-installer.ps1`.

**Checkpoint**: The four shortcut states are deterministic in source and built-package metadata without component conditions.

---

## Phase 4: User Story 2 - Choose Successful-Install Actions (Priority: P1)

**Goal**: Offer independent application and documentation completion actions only after successful fresh attended installation.

**Independent Test**: Source and compiled-MSI contracts prove two independent checkboxes, defaults, ordered unelevated shell events, exact targets, and absence from every execute-sequence or non-fresh flow.

- [x] T014 [US2] Add failing semantic and mutation tests for the owned UI sequence, independent completion controls, defaults, actions, targets, event order, and guards in `test/integration/windows_installer_contract_test.go`.
- [x] T015 [US2] Replace the stock minimal UI with a package-owned FeatureTree-derived flow and custom two-checkbox success dialog in `build/windows/goschedule.wxs`.
- [x] T016 [US2] Add two guarded unelevated shell actions and ordered Finish events in `build/windows/goschedule.wxs`.
- [x] T017 [US2] Extend source and compiled-MSI verification for completion controls, events, conditions, and execute-sequence exclusion in `build/windows/verify_wxs.ps1` and `test/windows/inspect-installer.ps1`.

**Checkpoint**: The four completion selections are independent in installer metadata and no silent, maintenance, upgrade, removal, rollback, cancel, or failure path can invoke them.

---

## Phase 5: User Story 3 - Preserve Application Data During Removal (Priority: P1)

**Goal**: Explain removal scope and preserve application data by default in full UI and unattended removal.

**Independent Test**: UI/source contracts prove the inventory and default; the hosted native preserve scenario proves exact machine and preference bytes survive software removal and reinstall.

- [x] T018 [US3] Add failing semantic and mutation tests for removal inventory, preserve default, direct-full-UI interception, and managed default behavior in `test/integration/windows_installer_contract_test.go`.
- [x] T019 [US3] Add the removal inventory dialog, preserve/wipe selection, maintenance/direct-remove routing, and secure default property in `build/windows/goschedule.wxs`.
- [x] T020 [US3] Extend source and compiled-MSI verification for removal UI routing, property security, defaults, and absence of destructive actions in preserve paths in `build/windows/verify_wxs.ps1` and `test/windows/inspect-installer.ps1`.

**Checkpoint**: Every flow preserves unless an exact explicit wipe is confirmed or supplied.

---

## Phase 6: User Story 4 - Explicitly Wipe Application Data (Priority: P1)

**Goal**: Wipe all declared local application-owned roots only after explicit confirmation and successful software-removal commit, with safe refusal and accurate residual evidence.

**Independent Test**: Helper tests prove safety and result states; MSI contracts prove commit-only gating; hosted native wipe, repair, upgrade, invalid-value, sentinel, and locked-item scenarios consume the result ledger.

- [x] T021 [US4] Add failing installer tests for exact wipe property validation, destructive confirmation, GUI close behavior, embedded helper, commit-only execution, complete-removal condition, and ignored-return safety in `test/integration/windows_installer_contract_test.go`.
- [x] T022 [US4] Add the wipe confirmation dialog, invalid-value guard, GUI close action, embedded helper, commit wipe action, and truthful retained-result guidance in `build/windows/goschedule.wxs`.
- [x] T023 [US4] Extend source and compiled-MSI verification for Binary, Property, CustomAction, sequence, dialog, condition, and close-application relationships in `build/windows/verify_wxs.ps1` and `test/windows/inspect-installer.ps1`.
- [x] T024 [US4] Add a disposable hosted Windows installer contract matrix with hidden noninteractive child processes in `test/windows/Invoke-InstallerContractCI.ps1`.
- [x] T025 [US4] Add the Windows MSI build, inspection, silent lifecycle, evidence upload, and helper window-subsystem gate to `.github/workflows/ci.yml` and synchronize `.github/workflows/release.yml`.

**Checkpoint**: Source, built package, helper tests, and hosted silent Windows execution prove preserve/wipe gating without claiming #94's attended desktop evidence.

---

## Phase 7: Documentation, Evidence, and Lifecycle

- [x] T026 Update `docs/INSTALL-windows.md` and `test/windows/README.md` with shortcut, completion, maintenance, preserve/wipe, property, scope, result-report, and #94 attended-proof guidance.
- [x] T027 Update `CHANGELOG.md` with S039 behavior and dated architectural decisions for pinned installer, workflow, and documentation changes.
- [x] T028 Run all focused Go, PowerShell parser, WiX source, helper cross-build, and locally available compiled-MSI checks and record outcomes in `specs/039-windows-setup-lifecycle-control/verification.md`.
- [x] T029 Run all eight canonical gates with `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` and record exact outcomes in `specs/039-windows-setup-lifecycle-control/verification.md`.
- [x] T030 Check changed text for UTF-8 without BOM and mojibake, resolve every required task, and update `spec.md` plus `specs/README.md` to `Implemented` with objective delivery evidence and an explicit #94/#97/#98 disposition.
- [x] T031 Review the final diff for destructive-scope safety, hidden-process guarantees, issue traceability, workflow integrity, and accidental unrelated changes.

## Dependencies & Execution Order

- Phase 1 establishes the approved design and blocking analyze gate.
- Cleanup safety in Phase 2 completes before any MSI destructive action is authored.
- Shortcut and completion work share the owned UI but remain independently testable.
- Preserve routing completes before wipe routing.
- Hosted lifecycle and documentation depend on all four user stories.
- Pull-request publication, external review, merge, and branch cleanup remain delivery evidence rather than feature tasks.

## Parallel Opportunities

- After Phase 2, shortcut contract work and cleanup-helper refinements can proceed independently when they touch different files.
- Source verification and compiled-MSI inspection may be extended in parallel only when file ownership is separated; changes to the same verifier remain sequential.
- Documentation drafts may proceed after the public property and feature identities stabilize.

## Implementation Strategy

1. Establish failing cleanup and installer contracts.
2. Implement the cleanup helper and prove its fail-closed boundaries.
3. Deliver shortcut features as the first independently testable installer increment.
4. Replace the UI once, adding completion and removal routing coherently.
5. Add commit-only wipe authoring and hosted silent lifecycle proof.
6. Run canonical verification, record unavailable attended evidence honestly, and leave #97/#98 open pending #94.
