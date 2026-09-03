# Feature Specification: Guided Windows Uninstall Entry

**Feature Branch**: `codex/041-guided-windows-uninstall`

**Created**: 2026-09-02

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/041-guided-windows-uninstall`; operator-authorized automatic push and pull request; merge, tag, and release remain unauthorized

**Input**: Repair GitHub issue #98 after the current unreleased MSI removed go-schedule directly from Windows Settings without showing its package-owned preserve-or-wipe flow.

## Context and Scope

The S039 package contains an attended removal inventory, preserve-or-wipe choice, destructive confirmation, safe cleanup helper, and deterministic silent properties. A maintainer nevertheless selected Uninstall from Windows Settings and the package was removed without those choices. Windows Settings can invoke a native MSI removal through reduced interface behavior that does not show package-authored wizard pages. S041 must make the operating-system management entry lead users into the full maintenance wizard before removal while retaining direct silent administration.

This is a focused repair of the entry path into the existing #98 workflow. It does not redesign data ownership, broaden wipe scope, add a bootstrapper, or claim the exact-candidate attended release gate in #94.

## Clarifications

### Session 2026-09-02

- Q: Should S041 try to force custom dialogs into the reduced native removal path? -> A: No. Use Windows Installer's supported maintenance entry, then let users select Remove inside the full package-owned wizard.
- Q: What should remain available to automation? -> A: Direct unattended removal remains supported and preserves data unless the exact documented wipe opt-in is supplied.
- Q: Does the operator's up-front publication direction satisfy the usual autopilot pre-push authorization? -> A: Yes for this S041 review branch and pull request only; merge, tag, and release remain unauthorized.
- Q: What if `ARPNOREMOVE` suppresses the maintenance command or disables Remove inside the stock WiX maintenance page? -> A: Own the current ProductCode's `/I` registration and maintenance page explicitly; keep the external direct Remove action suppressed without disabling guided removal inside the package.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enter Guided Removal from Windows Settings (Priority: P1)

As a Windows user managing an installed copy of go-schedule, I enter the product's maintenance wizard from the system application list and choose Remove before any software or application data is deleted.

**Why this priority**: The reported defect bypasses the only place where the user can understand and choose the data-removal policy.

**Independent Test**: Install the candidate on a disposable Windows host, inspect the registered management commands, invoke the visible maintenance action, select Remove, and observe the removal inventory with preserve selected before continuing.

**Acceptance Scenarios**:

1. **Given** go-schedule is installed, **When** the user opens its Windows application-management actions, **Then** the supported visible action opens the full go-schedule maintenance wizard rather than immediately removing the product.
2. **Given** the maintenance wizard, **When** the user selects Remove, **Then** the existing removal inventory and preserve-or-wipe choice appear with preserve selected.
3. **Given** either the maintenance page or removal-choice page, **When** the user cancels, **Then** no product software or application data is removed.
4. **Given** the obsolete direct system removal action could bypass the choice page, **When** the candidate is registered, **Then** that action is not offered as the attended user path.

---

### User Story 2 - Preserve Managed and Silent Removal Contracts (Priority: P1)

As an administrator, I can still repair, change, preserve-remove, or explicitly wipe go-schedule through documented commands without interactive prompts or completion launches.

**Why this priority**: Fixing the attended entry must not break enterprise automation or change the safe preservation default.

**Independent Test**: Run the established hosted Windows lifecycle matrix with default and explicit removal properties, and compare the resulting product, shortcut, data, cleanup-ledger, and process states.

**Acceptance Scenarios**:

1. **Given** an unattended direct removal with no wipe property, **When** removal succeeds, **Then** software is removed and all declared application data remains unchanged.
2. **Given** an unattended direct removal with the exact wipe opt-in, **When** removal succeeds, **Then** the existing safe cleanup contract runs after software removal commits.
3. **Given** repair, modification, upgrade, or invalid removal input, **When** the operation runs, **Then** no unintended wipe or completion action occurs.

### Edge Cases

- Windows management surfaces that honor only a reduced native removal interface must not become a bypass around the package-owned choice page.
- Hiding the unsafe direct attended removal action must not hide the product itself or its maintenance action.
- A direct administrator invocation remains possible and must retain documented safe defaults even though it is not the supported attended Settings path.
- Upgrade from an older package must replace the registered management contract without producing duplicate visible product entries.
- Repair must restore the intended management registration if it is missing or damaged.
- Cancellation before execution must leave product files, registration, shortcuts, service state, and data unchanged.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Windows application list MUST expose one visible go-schedule product entry and a supported action that opens the full maintenance wizard.
- **FR-002**: The registered attended management action MUST NOT preselect direct removal before the package-owned wizard can collect the preserve-or-wipe choice.
- **FR-003**: The unsafe direct attended removal action that bypassed the choice page MUST NOT be offered by the system application list.
- **FR-004**: Selecting Remove in the maintenance wizard MUST route through the existing removal inventory and preserve-or-wipe page with preserve selected by default.
- **FR-005**: Canceling from maintenance or removal choice MUST remove no software or application data.
- **FR-006**: Direct unattended removal MUST remain available, default to preserve, and require the existing exact opt-in for wipe.
- **FR-007**: Repair, feature modification, upgrade, rollback, invalid input, and silent removal MUST retain their existing no-unintended-wipe and no-completion-launch guarantees.
- **FR-008**: Upgrade and repair MUST maintain exactly one visible application-management entry with current product identity and working maintenance behavior.
- **FR-009**: Source and compiled-package checks MUST prove the management registration contract, the full-interface removal route, and the absence of competing or bypassing attended routes.
- **FR-010**: Hosted Windows lifecycle evidence MUST record the installed registration values and prove silent preserve, wipe, repair, upgrade, and invalid-input behavior.
- **FR-011**: Windows installation and test documentation MUST explain how to start guided removal from Settings and distinguish it from documented unattended commands.
- **FR-012**: Exact-candidate attended observation of the Settings-to-choice-page journey MUST remain required by #94 before #98 closes.
- **FR-013**: The implementation MUST pass every canonical repository gate without weakening the existing installer cleanup, shortcut, LocalSystem, or release-candidate evidence contracts.

### Key Entities

- **Application-management entry**: The single visible Windows record that identifies go-schedule and exposes its supported maintenance action.
- **Maintenance action**: The attended operating-system entry that opens the package-owned change, repair, and remove wizard without preselecting destructive work.
- **Direct removal action**: A reduced-interface operating-system path that can remove a native MSI before custom removal choices are collected and therefore must not be the supported attended entry.
- **Removal mode**: The existing preserve or wipe decision, safe by default and destructive only after explicit confirmation or exact unattended opt-in.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A fresh install and an upgrade each produce exactly one visible go-schedule entry whose supported attended action reaches the maintenance wizard before any removal begins.
- **SC-002**: In 100 percent of attended Settings-path observations, the removal inventory and preserve-or-wipe choice appear before software removal, with preserve selected initially.
- **SC-003**: Canceling from either pre-removal page leaves 100 percent of inventoried software and application data unchanged.
- **SC-004**: All existing hosted silent preserve, wipe, repair, upgrade, invalid-input, shortcut, cleanup-safety, and non-launch scenarios continue to pass.
- **SC-005**: Source and compiled-candidate inspection reject a package that restores the bypassing attended action, hides the maintenance action, or registers duplicate visible entries.
- **SC-006**: Documentation allows a Windows user to reach the guided removal choice from Settings without consulting source code.
- **SC-007**: All eight canonical repository gates pass and #94 retains an explicit exact-candidate attended verification step for this repaired journey.

## Assumptions

- Windows Installer permits disabling direct Remove while a package-owned maintenance command and wizard still provide guided removal; the package must explicitly preserve both pieces on Windows versions where `ARPNOREMOVE` suppresses the generated `ModifyPath` and WiX disables its stock Remove control.
- Windows Settings wording can vary by Windows release; the contract is the visible maintenance action and guided flow, not one immutable button label.
- The S039 preserve/wipe implementation and cleanup ownership boundaries remain correct unless S041 testing exposes a regression.
- Hosted Windows Server can prove registration and silent lifecycle semantics, while the final Windows 11 Settings interaction remains attended evidence owned by #94.

## Dependencies and Traceability

- Repairs the attended entry portion of #98 after the maintainer's 2026-09-03 exploratory failure report.
- Child of coordinator #96.
- Depends on the S039 setup lifecycle implementation from #97.
- Supplies a corrected candidate contract to #94; it does not close #94, #98, or #96 without exact-candidate attended proof.

## Out of Scope

- Publishing or tagging a release, merging the S041 pull request, or claiming exact-candidate Windows 11 acceptance.
- Replacing MSI packaging with a bootstrapper or custom executable uninstaller.
- Broadening or changing application-data, user-profile, reparse-point, cleanup-ledger, or `goschedadmin` ownership rules.
- Implementing unrelated GUI ergonomics issues #101, #103, #104, #105, or #106.
- Changing scheduler, task execution, IPC, persistence, startup window sizing, or error presentation.
