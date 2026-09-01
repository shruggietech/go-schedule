# Tasks: Windows Installed Core Recovery

**Input**: Design documents from `/specs/038-windows-installed-core-recovery/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/windows-installed-core.md`, `quickstart.md`

**Tests**: Required by issues #90 and #93. Story tests are written and observed failing before production implementation.

## Phase 1: Specification and Baseline

- [x] T001 Complete specification, clarification, requirements checklist, plan, research, data model, verification contract, and quickstart under `specs/038-windows-installed-core-recovery/`.
- [x] T002 Run Spec Kit's read-only consistency analysis across `spec.md`, `plan.md`, and `tasks.md`; resolve every critical or high-severity finding before implementation.
- [x] T003 Record clean branch state and run focused baseline tests for `internal/ipc` and `internal/executor`.

---

## Phase 2: User Story 1 - Ordinary Desktop Access (Priority: P1)

**Goal**: Let a verified direct `goschedadmin` user reach the restricted named pipe from a standard token that lacks the alias SID, without broadening access.

**Independent Test**: Descriptor tests prove direct-user expansion, unrelated-user omission, group-ACE retention, deterministic de-duplication, invalid SID rejection, and enumeration failure before listener creation.

### Tests

- [x] T004 [US1] Add failing Windows descriptor-policy tests in `internal/ipc/ipc_windows_test.go` for direct users, group-valued members, duplicates, ordering, invalid SIDs, and enumeration failure.
- [x] T005 [US1] Add a native current-token access contract or probe assertion in `test/windows/Invoke-InstallerLifecycle.ps1` that records direct membership, token omission, descriptor, and ordinary health access.

### Implementation

- [x] T006 [US1] Implement bounded Netapi32 direct-member enumeration in `internal/ipc/ipc_windows.go`, including buffer release and typed member usage.
- [x] T007 [US1] Build deterministic validated restricted SDDL in `internal/ipc/ipc_windows.go`, preserving SYSTEM, Built-in Administrators, and configured-group ACEs while adding only direct user SIDs.
- [x] T008 [US1] Update `docs/INSTALL-windows.md` and `specs/036-ipc-access-denied-recovery/spec.md` to document the corrected routine-access contract and the narrow S036 supersession.
- [x] T009 [US1] Run focused Windows IPC tests and cross-compile/test the Windows-tagged package.

---

## Phase 3: User Story 2 - Service-Hosted Task Execution (Priority: P1)

**Goal**: Prove real LocalSystem service execution succeeds manually and on schedule, with useful controls for exit and start failures.

**Independent Test**: The installed lifecycle probe creates tasks through the installed CLI, observes production manual and scheduled success records and marker effects, then verifies one nonzero exit and one process-start failure.

### Tests

- [x] T010 [US2] Add failing executor tests in `internal/executor/executor_test.go` for process-start boundary text, executable-only disclosure, and preservation of nonzero exit output.
- [x] T011 [US2] Add execution-probe assertions and evidence fields in `test/windows/Invoke-InstallerLifecycle.ps1` before implementing helper behavior.

### Implementation

- [x] T012 [US2] Implement secret-safe process-start diagnostics in `internal/executor/executor.go` without changing successful or nonzero-exit semantics.
- [x] T013 [US2] Implement the installed service execution probe in `test/windows/Invoke-InstallerLifecycle.ps1` using absolute system commands, manual and scheduled markers, bounded polling, and failure controls.
- [x] T014 [US2] Document execution-probe prerequisites, invocation, evidence, and service-context command rules in `test/windows/README.md` and `docs/INSTALL-windows.md`.
- [x] T015 [US2] Run focused executor tests and parse/lint the PowerShell probe.

---

## Phase 4: Integration, Evidence, and Lifecycle

- [x] T016 Update `CHANGELOG.md` with the S038 correction and a dated decision covering pinned Windows surfaces and the S036 authorization-scope deviation.
- [x] T017 Run the native evidence probe when host prerequisites are available; otherwise record an explicit unavailable result without substituting fake evidence in `specs/038-windows-installed-core-recovery/verification.md`.
- [x] T018 Run all eight canonical gates with `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` and record exact outcomes in `specs/038-windows-installed-core-recovery/verification.md`.
- [x] T019 Check all text outputs for UTF-8 without BOM and mojibake, then update `spec.md`, `tasks.md`, and `specs/README.md` to `Implemented` with objective verification evidence.
- [x] T020 Review the final diff for scope, security boundary, issue traceability, and accidental unrelated changes.

## Dependencies & Execution Order

- Phase 1 completes before production changes.
- User Story 1 and User Story 2 are independently testable but execute chronologically in this autopilot run to keep evidence and file ownership clear.
- Integration and lifecycle work depends on both user stories.
- Pull-request publication and external review are delivery workflow evidence, not feature tasks.
