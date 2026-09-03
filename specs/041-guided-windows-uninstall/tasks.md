# Tasks: Guided Windows Uninstall Entry

**Input**: Design documents from `/specs/041-guided-windows-uninstall/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Regression tests are mandatory. Add them before the WiX source change and demonstrate the red-to-green transition.

**Organization**: Tasks are grouped by independently testable user story. S041 is intentionally focused on #98 because destructive uninstall entry and registration semantics require a coherent Windows-specific review.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it changes a different file after shared expectations are fixed
- **[Story]**: Maps a task to one user story from spec.md

## Phase 1: Setup and Baseline

**Purpose**: Preserve the reported failure and lock down the supported operating-system contract before implementation.

- [x] T001 Record the maintainer's failed Windows Settings removal and authoritative Windows Installer semantics in `specs/041-guided-windows-uninstall/research.md`
- [x] T002 Confirm the existing maintenance Remove route, silent removal guards, compiled inspector, hosted lifecycle probe, and exact-candidate evidence boundary across `build/windows/goschedule.wxs`, `test/windows/`, and `specs/040-windows-release-candidate-gate/`

---

## Phase 2: Foundational Regression Contracts

**Purpose**: Create failing assertions for the registration behavior that S039 did not verify.

**CRITICAL**: These tests must fail against the S039 package source before the WiX property is added.

- [x] T003 [P] Add source-contract and mutation assertions for maintenance-only ARP registration in `test/integration/windows_installer_contract_test.go`
- [x] T004 [P] Add structural source assertions for `ARPNOREMOVE=1` and absent `ARPNOMODIFY` in `build/windows/verify_wxs.ps1`
- [x] T005 [P] Add compiled Property-table evidence for maintenance-only ARP registration in `test/windows/inspect-installer.ps1`
- [x] T006 [P] Add installed registry assertions and evidence recording for fresh install, repair, and upgrade in `test/windows/Invoke-InstallerContractCI.ps1`
- [x] T007 Run the focused source regression tests before implementation and record the expected failures in `specs/041-guided-windows-uninstall/verification.md`

**Checkpoint**: The regression suite demonstrates that S039 permits the bypassing direct Settings removal entry.

---

## Phase 3: User Story 1 - Enter Guided Removal from Windows Settings (Priority: P1) MVP

**Goal**: Make maintenance the supported attended Windows application-management entry so Remove reaches the existing preserve-or-wipe flow.

**Independent Test**: Source and compiled MSI contracts show direct Remove suppressed, maintenance retained, and the maintenance Remove control routed to `GoScheduleUninstallDlg`.

### Tests for User Story 1

- [x] T008 [US1] Re-run the focused source regression tests and preserve the failing result before editing `build/windows/goschedule.wxs`

### Implementation for User Story 1

- [x] T009 [US1] Author the maintenance-only Windows application-management property and rationale in `build/windows/goschedule.wxs`
- [x] T010 [US1] Run focused source contracts and confirm all new registration assertions turn green
- [x] T011 [US1] Extend the attended Settings contract and evidence boundary in `test/windows/README.md`

**Checkpoint**: Source authoring rejects the bypass path and retains full maintenance routing.

---

## Phase 4: User Story 2 - Preserve Managed and Silent Removal Contracts (Priority: P1)

**Goal**: Keep direct silent preserve/wipe, repair, upgrade, invalid-input, shortcut, cleanup, and non-launch behavior unchanged.

**Independent Test**: The hosted Windows installer job builds the MSI, proves compiled and installed registration, and passes the complete lifecycle matrix.

### Tests for User Story 2

- [x] T012 [US2] Review every existing cleanup and execution condition for unintended S041 changes in `build/windows/goschedule.wxs` and `test/windows/Invoke-InstallerContractCI.ps1`

### Implementation for User Story 2

- [x] T013 [US2] Update guided and unattended removal instructions in `docs/INSTALL-windows.md`
- [x] T014 [US2] Run the local compiled inspection against the prior S039 candidate to prove the new inspector rejects the old registration contract
- [x] T015 [US2] Run all locally available installer, cleanup-helper, and integration contracts

**Checkpoint**: S041 changes only the attended entry contract and detects the old candidate as nonconforming.

---

## Phase 5: Polish and Cross-Cutting Verification

**Purpose**: Complete traceability, lifecycle status, encoding integrity, and repository-wide evidence.

- [x] T016 [P] Add the S041 feature and dated maintenance-routing decision to `CHANGELOG.md`
- [x] T017 [P] Add S041 to the chronological feature index in `specs/README.md`
- [x] T018 Update `specs/041-guided-windows-uninstall/spec.md` to Implemented and create `specs/041-guided-windows-uninstall/verification.md` with red/green, gate, and remaining attended-risk evidence
- [x] T019 Run `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` in the foreground and record all eight canonical gate results in `specs/041-guided-windows-uninstall/verification.md`
- [x] T020 Perform a systematic local diff, history, safety, and requirements review; fix every finding with confidence at least 80 percent
- [x] T021 Validate UTF-8 without BOM, scan changed files for mojibake, run `git diff --check`, and mark every completed task in `specs/041-guided-windows-uninstall/tasks.md`
- [x] T022 Commit the complete verified slice on `codex/041-guided-windows-uninstall` with issue #98 traceability

---

## Phase 6: Hosted Evidence Correction

**Purpose**: Correct the Windows Installer behavior discovered only after PR publication.

- [x] T023 Record the failed hosted assertion that `ARPNOREMOVE` omitted `ModifyPath` and inspect the compiled WiX 6.0.2 stock maintenance control conditions
- [x] T024 Add failing source, compiled-table, and mutation contracts for an MSI-owned `/I[ProductCode]` registration and package-owned enabled Remove control
- [x] T025 Implement the owned application-management registry component and maintenance page, then update every affected specification and decision record
- [ ] T026 Run focused and canonical verification, amend and push the correction, then require green compiled/installed hosted evidence and resolve the authorized second Codex review

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependency.
- **Regression Contracts (Phase 2)**: Depends on the baseline; blocks implementation.
- **User Story 1 (Phase 3)**: Depends on failing registration tests.
- **User Story 2 (Phase 4)**: Depends on the new registration property and retains all S039 lifecycle behavior.
- **Polish (Phase 5)**: Depends on both stories.

### User Story Dependencies

- **User Story 1** is the MVP and directly repairs the reported Settings bypass.
- **User Story 2** validates that the repair does not trade away safe unattended administration.

### Parallel Opportunities

- T003 through T006 touch independent verification layers after the registration contract is fixed.
- T016 and T017 touch independent documentation files after implementation stabilizes.
- No multi-agent delegation is used because current session policy reserves it for explicit operator requests.

## Implementation Strategy

1. Preserve the maintainer report and authoritative registration decision.
2. Make all four verification layers fail against the current source/candidate where locally possible.
3. Add the declarative ARP property, owned maintenance registration, and owned maintenance page.
4. Turn local contracts green and use hosted CI for fresh compiled/installed evidence.
5. Correct every hosted or review finding without weakening the attended-evidence boundary.
6. Leave #98 open until #94 captures the exact-candidate Windows 11 attended journey.

## Notes

- The operator has already authorized push, PR creation, and verified in-scope review-fix pushes for S041.
- At most one manual `@Codex review` request may follow the automatic first review round.
- Merge, tag, release, and material scope expansion remain unauthorized.
