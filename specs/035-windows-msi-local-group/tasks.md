# Tasks: Windows MSI Local-Group Recovery

**Input**: Design documents from `specs/035-windows-msi-local-group/`

**Tests**: Required by constitution Principle II and issue #83.

## Phase 1: Setup and Red Baseline

- [x] T001 Record issue #83, v0.9.0, host, elevation, artifact, and false-green validator baselines in `specs/035-windows-msi-local-group/verification.md`.
- [x] T002 Update active lifecycle inventory in `specs/README.md` and Spec-Kit context in `CLAUDE.md`.
- [x] T003 Add a failing non-empty group-domain mutation to `test/integration/windows_installer_contract_test.go`.
- [x] T004 Run the focused Go contract against current authoring and record the expected red result in `specs/035-windows-msi-local-group/verification.md`.

## Phase 2: User Story 1 - Secure Fresh Install (Priority: P1)

**Goal**: Route `goschedadmin` creation through the elevated local-group path.

**Independent Test**: Source and compiled contracts require an empty group domain, and a clean Windows 11 install records successful provisioning before service startup.

- [x] T005 [US1] Remove the non-empty group domain while preserving all other group and user attributes in `build/windows/goschedule.wxs`.
- [x] T006 [P] [US1] Require an empty effective group domain in `build/windows/verify_wxs.ps1`.
- [x] T007 [US1] Update positive and negative Go group contracts in `test/integration/windows_installer_contract_test.go`.
- [x] T008 [US1] Query and validate the compiled `Wix4Group` row in `test/windows/inspect-installer.ps1`.
- [x] T009 [US1] Extend the evidence-tooling contract for compiled group-row assertions in `test/integration/windows_installer_contract_test.go`.
- [x] T010 [US1] Run focused Go and PowerShell source contracts and record green evidence in `specs/035-windows-msi-local-group/verification.md`.

## Phase 3: User Story 2 - Preserve MSI Lifecycle State (Priority: P1)

**Goal**: Capture honest fresh, repair, reinstall, v0.9.0 upgrade, and uninstall evidence for group, membership, service, logs, and final state.

**Independent Test**: Separately reset Windows 11 runs produce complete proven fresh and upgrade evidence, or explicitly report failed/unavailable fields.

- [x] T011 [US2] Add scenario, prior-artifact, operation exit-code, log-path, log-hash, and diagnostic evidence to `test/windows/Invoke-InstallerLifecycle.ps1`.
- [x] T012 [US2] Add group SID, normalized membership, intended-account, and service observations to `test/windows/Invoke-InstallerLifecycle.ps1`.
- [x] T013 [US2] Add explicit repair and same-version reinstall phases with hidden noninteractive installer launches to `test/windows/Invoke-InstallerLifecycle.ps1`.
- [x] T014 [US2] Add the separately reset v0.9.0 upgrade scenario and preserved group/member invariance checks to `test/windows/Invoke-InstallerLifecycle.ps1`.
- [x] T015 [US2] Strengthen lifecycle-script contract assertions in `test/integration/windows_installer_contract_test.go`.
- [x] T016 [US2] Document fresh and upgrade invocation, isolation, session-token, and evidence boundaries in `test/windows/README.md`.
- [x] T017 [US2] Add a post-refresh, non-elevated access-probe mode and its evidence contract to `test/windows/Invoke-InstallerLifecycle.ps1` and `test/integration/windows_installer_contract_test.go`.

## Phase 4: User Story 3 - Diagnose Provisioning Restrictions (Priority: P2)

**Goal**: Give Windows operators an actionable verbose MSI log workflow.

**Independent Test**: The guide supplies a runnable command and explains the relevant exit code, group, membership, service, and access-denied evidence.

- [x] T018 [US3] Add verbose logging and provisioning diagnosis guidance to `docs/INSTALL-windows.md`.
- [x] T019 [US3] Update S016's historical evidence boundary to point to S035's stronger lifecycle contract in `specs/016-windows-installer-identity/contracts/verification-evidence.md`.

## Phase 5: Documentation, Analysis, and Verification

- [x] T020 Add the pinned-installer deviation and verification decision to `CHANGELOG.md`.
- [x] T021 Run blocking cross-artifact analysis and resolve every critical or high finding across S035 `spec.md`, `plan.md`, and `tasks.md`.
- [x] T022 Run focused Go tests and PowerShell compliance/source checks.
- [x] T023 Run `sh scripts/verify.sh all` in the foreground through format, vet, lint, race, GUI, coverage, docs, and automation.
- [ ] T024 Run qualifying fresh, access-probe, and v0.9.0-upgrade Windows 11 lifecycles when prerequisites exist; otherwise record the exact unavailable prerequisites in `specs/035-windows-msi-local-group/verification.md` without claiming success.
- [x] T025 Audit UTF-8 without BOM, mojibake, whitespace, checklist completion, issue traceability, and worktree state; update S035 lifecycle fields and `specs/README.md` with the supported evidence state.

## Dependencies and Execution Order

- Phase 1 establishes a failing regression contract before authoring changes.
- US1 provides corrected authoring and static/compiled contracts.
- US2 depends on US1's named evidence contract but can be implemented without a candidate MSI; its native runs remain environment-gated.
- US3 can proceed after evidence terminology stabilizes in US2.
- Phase 5 requires all implementation tasks and checklists to be complete.

## Parallel Opportunities

- T006 and the T007 fixture updates touch independent validators after T005.
- T008 and T011 begin on separate inspector and lifecycle scripts.
- T017 and T018 are independent documentation surfaces after US2 terminology is stable.

## Implementation Strategy

Deliver one reviewable regression-fix slice with test-first source protection, compiled-artifact inspection, and a complete native evidence harness. Publish the authorized review branch and PR after canonical local gates pass. Use `Closes #83` only if qualifying Windows 11 evidence is proven; otherwise use `Refs #83` and leave the runtime acceptance gap visible for reviewers.
