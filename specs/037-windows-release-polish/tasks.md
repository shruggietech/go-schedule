# Tasks: Windows Release Polish

**Input**: Design documents from `specs/037-windows-release-polish/`

**Tests**: Required and written before their corresponding implementation changes.

## Phase 1: Specification and Baseline

- [x] T001 Read complete issue #88 and #89 bodies and acceptance criteria and record issue-level traceability in `specs/037-windows-release-polish/spec.md`.
- [x] T002 Generate and validate requirements and Windows-release checklists in `specs/037-windows-release-polish/checklists/`.
- [x] T003 Record bounded-window, monitor-selection, installer-subject, and evidence-boundary decisions in `specs/037-windows-release-polish/plan.md` and supporting artifacts.
- [x] T004 Advance S037 to In Progress in `specs/037-windows-release-polish/spec.md` and `specs/README.md` before implementation begins.

## Phase 2: User Story 1 - Reachable Windowed Launch (Priority: P1)

**Goal**: Request a restored 1280x800 window capped independently at 90 percent of the selected monitor's logical work area.

**Independent Test**: Required table cases and the 800x600 root reachability contract fail against the old full-work-area behavior and pass after implementation.

- [x] T005 [US1] Add the required table-driven sizing regression cases to `gui/screen_test.go` and record the expected red result in `specs/037-windows-release-polish/verification.md`.
- [x] T006 [US1] Implement the pure 90 percent logical sizing contract in `gui/screen.go` and update startup intent in `gui/app.go`.
- [x] T007 [US1] Add a Windows monitor-selection contract test in `gui/screen_windows_test.go` before changing the adapter.
- [x] T008 [US1] Replace primary-work-area discovery with nearest-pointer monitor work-area discovery in `gui/screen_windows.go`.
- [x] T009 [US1] Add and satisfy a headless root reachability regression at effective 800x600 in `gui/app_test.go`.
- [x] T010 [US1] Correct specification 003 FR-001, its user story, success criterion, assumptions, and add reversal evidence in `specs/003-gui-editor-refinements/spec.md` and `specs/003-gui-editor-refinements/verification.md`.

## Phase 3: User Story 2 - Concise Explorer Metadata (Priority: P1)

**Goal**: Compile and inspect the exact approved Windows MSI Subject while preserving all unrelated installer identity.

**Independent Test**: Static verification rejects the old Subject, and compiled inspection fails any PID 3 value other than the approved copy while recording the hash and observed value.

- [x] T011 [US2] Extend the source verification contract in `build/windows/verify_wxs.ps1` to reject the old or any non-approved Subject and record the expected red result.
- [x] T012 [US2] Change only `SummaryInformation Description` in `build/windows/goschedule.wxs` to the approved value.
- [x] T013 [US2] Extend `test/windows/inspect-installer.ps1` to read and assert Summary Information PID 3 and include it in evidence output.
- [x] T014 [US2] Update `test/windows/README.md` with the compiled Subject and Explorer evidence boundaries.
- [x] T015 [US2] Run the source verifier, compiled candidate inspection, native Shell Subject assertion, and PowerShell parse checks; record actual results in `specs/037-windows-release-polish/verification.md`.

## Phase 4: Polish and Lifecycle Completion

- [x] T016 Record dated pinned-artifact decisions and issue-specific change entries in `CHANGELOG.md`.
- [x] T017 Refresh the managed plan reference in `CLAUDE.md` through the agent-context hook.
- [x] T018 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/037-windows-release-polish/verification.md`.
- [x] T019 Audit UTF-8 without BOM, mojibake, placeholders, unchecked required tasks, diff errors, installer identity preservation, and scope boundaries.
- [x] T020 Advance S037 to Implemented with delivery evidence in `specs/037-windows-release-polish/spec.md`, `specs/README.md`, and `specs/037-windows-release-polish/verification.md`.
- [x] T021 Commit S037 locally with the repository's conventional feature message and co-author trailer.

## Dependencies and Execution Order

- Phase 1 blocks implementation.
- US1 and US2 are independently testable after Phase 1 but are completed sequentially because this run has one implementer.
- Each story's regression test or verifier change precedes its implementation change.
- Canonical verification and lifecycle completion follow both stories.

## Issue Traceability

- **Issue #89**: T005 through T010, plus T018 through T020.
- **Issue #88**: T011 through T015, plus T018 through T020.

## Implementation Strategy

Deliver both P1 stories as one coherent Windows release-polish increment. Do not split the bundled slice, persist window placement, resize task-editor dialogs, change unrelated MSI identity, merge the PR, tag, or run a release workflow.
