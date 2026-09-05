# Tasks: Windows Installer Identity and Verification

**Input**: Design documents from `specs/016-windows-installer-identity/`
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Required by the feature specification and constitution. Bug-fix regressions precede WXS implementation.

## Phase 1: Setup and evidence baseline

**Purpose**: Establish the exact published and local baseline without changing the installed application.

- [X] T001 Record the latest published MSI metadata and current-host prerequisite status in `specs/016-windows-installer-identity/verification.md`
- [X] T002 Record v0.8.0 MSI Icon, Property, Shortcut, and Environment table observations in `specs/016-windows-installer-identity/verification.md`
- [X] T003 Record current installed-app registry, Start Menu shortcut, and machine PATH observations in `specs/016-windows-installer-identity/verification.md`

---

## Phase 2: Foundational source-contract regression tests

**Purpose**: Write portable tests that fail on the current missing installer icon relationships and protect the existing executable resource path.

**Critical**: These tests must fail before the WXS implementation task.

- [X] T004 Create structural WXS contract parsing helpers in `test/integration/windows_installer_contract_test.go`
- [X] T005 Add failing canonical-source assertions for the Icon, `ARPPRODUCTICON`, and `GuiShortcut` relationships in `test/integration/windows_installer_contract_test.go`
- [X] T006 Add independent mutation regressions for missing and misdirected icon relationships in `test/integration/windows_installer_contract_test.go`
- [X] T007 Add passing assertions for the existing 64-bit GUI executable resource pipeline in `.github/workflows/release.yml` to `test/integration/windows_installer_contract_test.go`
- [X] T008 Run the focused Windows installer contract test and record the expected pre-fix failure in `specs/016-windows-installer-identity/verification.md`

**Checkpoint**: The source regression suite fails only because the WXS lacks the three new icon relationships.

---

## Phase 3: User Story 1 - Recognize the installed application (Priority: P1)

**Goal**: Give the Start Menu shortcut and installed-apps entry one explicit official brand mark.

**Independent Test**: The structural regression suite and compiled candidate MSI both show one canonical Icon row consumed by the ARP property and shortcut.

- [X] T009 [US1] Add the canonical package Icon declaration in `build/windows/goschedule.wxs`
- [X] T010 [US1] Add the `ARPPRODUCTICON` relationship in `build/windows/goschedule.wxs`
- [X] T011 [US1] Add the explicit `GuiShortcut` icon relationship in `build/windows/goschedule.wxs`
- [X] T012 [US1] Run the focused Go contract test in `test/integration/windows_installer_contract_test.go` and confirm the WXS regression is green

**Checkpoint**: User Story 1 source contract is independently complete.

---

## Phase 4: User Story 2 - Recognize the running application (Priority: P1)

**Goal**: Protect the existing GUI executable icon resource and define honest native observation.

**Independent Test**: Portable tests prove the resource-generation contract; the disposable-host guide names window and taskbar observations separately.

- [X] T013 [US2] Document candidate, published, and native icon observations plus cache refresh boundaries in `test/windows/README.md`
- [X] T014 [US2] Document how to append Start Menu, installed-apps, title-area, and taskbar observations to lifecycle evidence in `test/windows/README.md`
- [X] T015 [US2] Record native observation as unavailable on this host and keep issue #33 open in `specs/016-windows-installer-identity/verification.md`

**Checkpoint**: Executable resource contract is proven; native observation is explicitly unavailable rather than misreported.

---

## Phase 5: User Story 3 - Trust command availability across install lifecycle (Priority: P1)

**Goal**: Supply a guarded clean-host procedure that records PATH lifecycle without risking a developer's existing installation.

**Independent Test**: On the current contaminated host the tool refuses before mutation and names every failed clean-baseline prerequisite; on a disposable host it covers all lifecycle stages.

- [X] T016 [US3] Implement clean-host, elevation, artifact, and evidence-path preconditions in `test/windows/install-lifecycle.ps1`
- [X] T017 [US3] Implement silent install, new-process CLI probes, and machine PATH cardinality capture in `test/windows/install-lifecycle.ps1`
- [X] T018 [US3] Implement reinstall, uninstall, unrelated-PATH preservation, and cleanup evidence in `test/windows/install-lifecycle.ps1`
- [X] T019 [US3] Run the lifecycle tool on the current host and record its expected non-mutating refusal in `specs/016-windows-installer-identity/verification.md`

**Checkpoint**: The procedure exists and refuses false evidence locally; issue #16 remains open until it runs against a published MSI on a clean host.

---

## Phase 6: User Story 4 - Preserve the installer contract (Priority: P2)

**Goal**: Give maintainers specific local failures and compiled-artifact evidence without installing the MSI.

**Independent Test**: Source verification rejects each broken relationship; candidate MSI inspection reports the correct table rows.

- [X] T020 [US4] Add icon source, ARP, and shortcut assertions with specific failures to `build/windows/verify_wxs.ps1`
- [X] T021 [P] [US4] Implement Windows Installer COM table inspection and evidence output in `test/windows/inspect-installer.ps1`
- [X] T022 [US4] Install WiX 5.0.2 as a local verification prerequisite outside tracked source and build a candidate MSI from `build/windows/goschedule.wxs`
- [X] T023 [US4] Inspect the candidate MSI with `test/windows/inspect-installer.ps1` and record the table evidence in `specs/016-windows-installer-identity/verification.md`

**Checkpoint**: Source and candidate-artifact evidence are `proven`.

---

## Phase 7: Polish and cross-cutting verification

**Purpose**: Synchronize policy, documentation, issue traceability, and all quality gates.

- [X] T024 Update `CHANGELOG.md` Unreleased Fixed and dated Decisions entries for #24/#25 and pinned `build/windows` changes
- [X] T025 Validate all Spec-Kit requirements and checklists against implementation and record issue closure eligibility in `specs/016-windows-installer-identity/verification.md`
- [X] T026 Run PowerShell syntax analysis and the focused source/artifact checks for `build/windows/verify_wxs.ps1` and `test/windows/*.ps1`
- [X] T027 Run `git diff --check`, UTF-8-without-BOM, and mojibake audits across all changed files
- [X] T028 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/016-windows-installer-identity/verification.md`
- [X] T029 Mark every completed task `[X]` in `specs/016-windows-installer-identity/tasks.md` and rerun `/speckit-analyze`
- [X] T030 Commit locally as `feat(016): harden Windows installer identity` with the required co-author trailer

---

## Dependencies and execution order

### Phase dependencies

- Phase 1 has no dependencies.
- Phase 2 depends on Phase 1 and is the blocking TDD gate.
- Phase 3 depends on the expected failure from Phase 2.
- Phases 4 and 5 depend on the source contracts from Phase 2 and may proceed in order after Phase 3.
- Phase 6 depends on the valid WXS from Phase 3.
- Phase 7 depends on all user-story phases.

### User story dependencies

- **US1** is the MVP and has no other story dependency after foundational tests.
- **US2** reuses the foundational executable-resource contract and is otherwise independent of PATH lifecycle work.
- **US3** depends only on the shared evidence vocabulary and clean-host rules.
- **US4** depends on US1 because compiled artifact inspection needs a valid WXS.

### Parallel opportunities

- Documentation tasks T013 and T014 are sequential because they share one file.
- T021 can be implemented independently after US1 because it affects a separate artifact-inspection script; T020 must complete before the release verifier is green.
- No tasks that mutate the WXS or evidence record run in parallel.

## Parallel example: US4

```text
Task T020: extend build/windows/verify_wxs.ps1 source assertions
Task T021: implement test/windows/inspect-installer.ps1 table inspection
Then T022: build candidate MSI
Then T023: inspect and record evidence
```

## Implementation strategy

### MVP first

1. Establish baseline and red source regression.
2. Implement US1's one-icon/two-consumer WXS contract.
3. Prove the contract green before adding evidence tooling.

### Incremental delivery

1. Source definition becomes correct and regression-tested.
2. Candidate MSI tables prove compiled authoring.
3. Guarded lifecycle tooling makes clean published verification repeatable.
4. Native observations stay visible until a suitable environment exists.
