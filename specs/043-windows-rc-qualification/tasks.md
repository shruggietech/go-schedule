# Tasks: Windows Demo Qualification

**Input**: Design documents from `specs/043-windows-rc-qualification/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, and `contracts/demo-handoff.md`

## Phase 1: Spec Kit gate

- [x] T001 Create S043 on `codex/043-windows-rc-qualification` and register `.specify/feature.json`.
- [x] T002 Author `spec.md` with explicit demo/formal evidence boundaries and issue traceability.
- [x] T003 Resolve clarification decisions in `spec.md` without an inter-step halt.
- [x] T004 Create requirements and attended-demo checklists under `checklists/`.
- [x] T005 Create `research.md`, `data-model.md`, `contracts/demo-handoff.md`, `quickstart.md`, and `plan.md`.
- [x] T006 Create this task plan and run the blocking cross-artifact analyze gate.

## Phase 2: Automated baseline

- [x] T007 [US2] Run focused release-gate, installer, cleanup, GUI, API, daemon, and executor tests from merged S042.
- [x] T008 [US2] Parse relevant PowerShell scripts and run WiX source validation.
- [x] T009 [US2] Run the complete eight-gate canonical suite in the foreground.
- [x] T010 [US2] Record destructive installed automation as unavailable unless the host satisfies its elevation and disposability checks.

## Phase 3: Demo candidate

- [x] T011 [US1] Ensure WiX Toolset 6.0.2 and matching UI/Util extensions are available without changing repository dependencies.
- [x] T012 [US1] Add a failing contract then extend `test/windows/inspect-installer.ps1` with truthful `local-demo` provenance without relaxing candidate or published validation.
- [x] T013 [US1] Build GUI, daemon, CLI, and cleanup helper with the S043 demo version into a clean local staging directory.
- [x] T014 [US1] Validate the stage with `build/windows/verify_wxs.ps1` and compile one numeric-version 1.0.0 MSI.
- [x] T015 [US1] Inspect the compiled MSI with `test/windows/inspect-installer.ps1` and retain its report.
- [x] T016 [US1] Independently verify embedded versions, ProductVersion, ProductCode, byte size, SHA-256, and source commit.

## Phase 4: Handoff and attended testing

- [x] T017 [US3] Record all automated results and boundaries in `verification.md`.
- [x] T018 [US3] Update `specs/README.md` and `CHANGELOG.md` for the in-progress local demo qualification.
- [x] T019 [US3] Run final UTF-8 without BOM, mojibake, diff-integrity, and local review checks.
- [x] T020 [US3] Commit the initial reviewed S043 source and evidence locally without pushing.
- [x] T021 [US1] Rebuild after the local-demo correction commit, then link the exact final MSI with its complete identity.
- [x] T022 [US3] Hand off `checklists/attended-demo.md`; wait for the operator to complete demo testing before any push or PR.

## Phase 5: Post-demo disposition

- [x] T023 [US3] Record operator observations against the handed-off demo hash, including explicit incomplete/failed checks and the post-wipe audit.
- [x] T024 [US3] Disposition confirmed non-blocking UI findings to #101, #104, #105, #106, #109, #110, #111, #112, and #113 without expanding S043 product scope; no product correction or rebuild was required.
- [x] T025 [US3] After demo completion, prepare the branch for the authorized pull-request CI and review boundary.
- [x] T026 [US3] **Deferred by scope:** after review/merge and separate tag authorization, the release ritual must run the formal exact-candidate matrix before closing #94/#98/#96. S043 does not claim or waive that future gate.

## Dependencies & execution order

- Phase 1 blocks all implementation and qualification work.
- Automated baseline must pass before the demo is presented.
- The demo artifact must be bound to the final local pre-handoff commit.
- Operator testing blocks branch publication.
- Formal candidate proof depends on reviewed merge and separately authorized tag staging.

## Parallel opportunities

No concurrent agent work is used. Read-only focused checks may be batched, but the canonical gate and candidate build run in the foreground, and artifact identity steps remain sequential.
