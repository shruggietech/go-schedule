# Feature Specification: Windows Setup Lifecycle Control

**Feature Branch**: `codex/039-windows-setup-lifecycle-control`

**Created**: 2026-09-02

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/039-windows-setup-lifecycle-control`; canonical eight-gate verification and local WiX 6.0.2 compiled-MSI inspection recorded in `verification.md`; hosted silent evidence runs in pull-request CI; #94 retains attended Windows 11 acceptance, so #97/#98 remain open

**Input**: Bundle GitHub issues #97 and #98 into one Windows setup lifecycle slice: user-controlled Start Menu and desktop shortcuts, independent post-install GUI and documentation actions, and a transparent preserve-or-wipe uninstall choice with safe cleanup and complete verification.

## Context and Scope

The v0.9.1 Windows setup always creates a Start Menu shortcut, cannot create a desktop shortcut, offers no useful completion actions, and removes the software without giving users an attended choice to preserve or erase application data. S039 makes those lifecycle decisions visible and deterministic while retaining safe defaults.

This slice covers initial install, maintenance modification, repair, upgrade, uninstall, and managed unattended operation. It completes the implementation work in #97 and #98 and supplies reusable verification for the downstream exact-candidate release gate in #94. It does not claim that release gate, publish a release, change application runtime behavior, or remove uncertain/shared operating-system security state.

## Clarifications

### Session 2026-09-02

- Q: Should the application and documentation completion actions be mutually exclusive? -> A: No. They are independent checkboxes because either, both, or neither can be useful.
- Q: What does “all application-related user data” include? -> A: The declared machine application root and application-specific preference roots belonging to safely registered local Windows profiles, never arbitrary exports or adjacent user files.
- Q: When may irreversible data cleanup begin? -> A: Only after the software-removal transaction commits successfully, so a canceled, failed, or rolled-back uninstall cannot erase application data.
- Q: How should `goschedadmin` and uncertain shared security state be handled? -> A: Preserve it and describe it separately unless exclusive installer ownership can be proven.
- Q: Does the operator's up-front publication direction satisfy the usual autopilot pre-push authorization? -> A: Yes for this S039 review branch and pull request only; merge, tag, and release remain unauthorized.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Choose Windows Shortcuts (Priority: P1)

As a person installing go-schedule, I choose whether setup creates Start Menu and desktop shortcuts instead of receiving shell integration I did not request.

**Why this priority**: Shortcut choice is the primary setup interaction and establishes the selectable installation state used by maintenance, upgrade, and uninstall.

**Independent Test**: Exercise all four shortcut combinations through built-package feature selection, then verify only the selected shortcuts exist and each row carries the windowless desktop application's canonical identity. The attended launch observation remains part of #94.

**Acceptance Scenarios**:

1. **Given** a fresh attended install, **When** the shortcut options first appear, **Then** Start Menu is selected and desktop is unselected.
2. **Given** any of the four shortcut selections, **When** installation succeeds, **Then** exactly the selected shortcuts are present in their standard Windows locations.
3. **Given** an installed product, **When** maintenance modification changes either shortcut selection, **Then** the corresponding shortcut is added or removed without reinstalling unrelated product components.
4. **Given** an upgrade over an existing installation, **When** the user does not explicitly change shortcut choices, **Then** the existing choices are preserved.
5. **Given** installer-created shortcuts, **When** the product is removed, **Then** those shortcuts and an empty product Start Menu folder are removed while unrelated user-created shortcuts remain.

---

### User Story 2 - Choose Successful-Install Actions (Priority: P1)

As a person completing an attended initial installation, I independently choose whether to open the desktop application and whether to open the project documentation.

**Why this priority**: The completion page is the user's transition from setup into successful first use and was explicitly missing in v0.9.1.

**Independent Test**: Prove in the built package that all four independent completion combinations map to ordered fresh-full-UI shell actions and that no action is available from other flows. #94 observes the targets in a real interactive desktop session.

**Acceptance Scenarios**:

1. **Given** a successful fresh attended install, **When** the completion page appears, **Then** application launch is selected and documentation launch is unselected.
2. **Given** the completion page, **When** either checkbox is changed, **Then** the other selection remains unchanged.
3. **Given** any of the four completion-action combinations, **When** Finish is selected, **Then** each chosen target opens once through the user's normal Windows application or HTTPS handler and each unchosen target remains closed.
4. **Given** cancel, failure, rollback, repair, upgrade, uninstall, or unattended execution, **When** setup ends, **Then** neither completion action runs.

---

### User Story 3 - Preserve Application Data During Removal (Priority: P1)

As a person uninstalling go-schedule, I can remove the software while preserving tasks, configuration, logs, runtime state, and desktop preferences for a later reinstall.

**Why this priority**: Preservation is the safe default and protects users who remove software temporarily or expect reinstall continuity.

**Independent Test**: Populate machine application data and per-user desktop preferences, remove the product using the default mode, then compare preserved content byte-for-byte and prove a reinstall can access the prior tasks.

**Acceptance Scenarios**:

1. **Given** an attended uninstall, **When** removal choices appear, **Then** the page identifies software always removed, data preserved by default, data eligible for wipe, and separately handled security state.
2. **Given** no explicit wipe selection, **When** attended or unattended removal succeeds, **Then** product software and installer-owned integration are removed while application data remains byte-for-byte unchanged.
3. **Given** the removal choice page, **When** the user cancels, **Then** no software or application data is removed.
4. **Given** a preserve-mode uninstall followed by reinstall, **When** the application starts, **Then** the prior tasks and supported preferences remain available.

---

### User Story 4 - Explicitly Wipe Application Data (Priority: P1)

As a person intentionally retiring go-schedule, I can explicitly confirm removal of all declared application-owned data without risking files outside those roots.

**Why this priority**: Data erasure is required for a complete uninstall but is destructive, so its scope, confirmation, failure reporting, and safety boundaries are release-critical.

**Independent Test**: Populate every declared machine and supported-user data root alongside out-of-scope controls, explicitly select and confirm wipe, then prove declared data is absent, controls remain, and reinstall starts clean.

**Acceptance Scenarios**:

1. **Given** attended removal, **When** wipe is selected, **Then** a separate confirmation names the destructive effect before removal can continue.
2. **Given** an explicitly confirmed wipe, **When** removal succeeds, **Then** product software plus the declared machine and all-supported-user application data are removed.
3. **Given** empty, relative, redirected, reparse-point, or out-of-scope candidate paths, **When** wipe evaluates them, **Then** deletion is refused and the incomplete cleanup is reported accurately.
4. **Given** a locked or inaccessible owned item, **When** wipe cannot remove it, **Then** removal does not claim a complete wipe and actionable failure evidence identifies the remaining owned path without exposing unrelated user information.
5. **Given** a wipe-mode uninstall followed by reinstall, **When** the application starts, **Then** no prior task, configuration, log, runtime, or supported desktop-preference state is present.

### Edge Cases

- Upgrade, repair, rollback, or a failed uninstall must never trigger the destructive wipe path.
- Unattended removal defaults to preservation unless the documented wipe property carries its one exact opt-in value.
- A user who selects wipe but does not pass the separate confirmation cannot proceed with destructive removal.
- Shortcut choice survives upgrade even when an older release had only the mandatory Start Menu shortcut.
- Removal deletes an installer-created Start Menu directory only when it is empty and leaves unrelated or user-created content untouched.
- Completion actions never inherit installer elevation and never run when no interactive desktop user is available.
- Wipe scope includes supported local Windows profiles that are not currently signed in, but excludes orphaned paths that cannot be tied safely to a registered local profile.
- Shared or uncertain `goschedadmin` group and membership state is preserved and described separately rather than silently classified as application data.
- Partial data cleanup cannot be rolled back into restored user data, so cleanup begins only after the software-removal transaction has committed successfully.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Attended initial setup MUST present independently selectable Start Menu and desktop shortcut choices before files are installed.
- **FR-002**: The Start Menu shortcut MUST be selected by default and the desktop shortcut MUST be unselected by default on a fresh install.
- **FR-003**: Each installed shortcut MUST use the canonical go-schedule name, icon, description, desktop-application target, and working location and MUST launch without a console window.
- **FR-004**: Maintenance modification MUST allow either shortcut to be added or removed independently without changing unrelated product components.
- **FR-005**: Upgrade MUST preserve the installed shortcut feature state unless the user explicitly changes it.
- **FR-006**: Uninstall MUST remove every installer-created shortcut and an empty installer-created Start Menu folder without removing unrelated user-created content.
- **FR-007**: Successful initial attended setup MUST offer independent application-launch and documentation-launch checkboxes on its completion page.
- **FR-008**: The application action MUST default selected, the documentation action MUST default unselected, and changing either MUST NOT change the other.
- **FR-009**: Selected completion targets MUST open exactly once as the interactive unelevated user; documentation MUST use the system-defined default HTTPS handler for `https://shruggietech.github.io/go-schedule/`.
- **FR-010**: Completion actions MUST NOT run after cancel, failure, rollback, repair, maintenance modification, upgrade, uninstall, or unattended execution.
- **FR-011**: Managed installation MUST expose deterministic documented values for shortcut selection and MUST never launch completion actions unattended.
- **FR-012**: Attended uninstall MUST present a concrete removal inventory and a preserve-or-wipe choice, with preserve selected by default.
- **FR-013**: The inventory MUST distinguish software and installer integration always removed, application data preserved by default, application data erased by wipe, and shared or uncertain security state handled separately.
- **FR-014**: Default attended and unattended uninstall MUST preserve the machine data root and every supported per-user desktop-preference root byte-for-byte.
- **FR-015**: Wipe MUST require both an explicit destructive selection and a separate attended confirmation; unattended wipe MUST require one exact documented opt-in property value.
- **FR-016**: A confirmed wipe MUST remove the task database, configuration, logs, runtime files, and application-specific desktop-preference data from the declared machine and supported local-user roots.
- **FR-017**: Wipe MUST start only after the software-removal transaction commits successfully and MUST NOT run during install, upgrade, repair, rollback, cancellation, or failed removal.
- **FR-018**: Wipe MUST derive product-owned roots from trusted operating-system locations, resolve canonical absolute paths, and refuse empty, relative, redirected, reparse-point, unregistered-profile, or out-of-scope deletion targets.
- **FR-019**: Wipe MUST leave user files outside declared product-owned roots untouched, including exports and unrelated content placed beside installer shortcuts.
- **FR-020**: Wipe MUST report partial cleanup or access failures accurately, identify remaining owned scope in actionable installation evidence, and MUST NOT claim success when declared data remains.
- **FR-021**: Removal MUST stop the desktop application and service before software removal or data cleanup and MUST preserve the established installer rollback behavior for software state.
- **FR-022**: Removal MUST preserve shared or uncertain `goschedadmin` group and membership state; any state known to be exclusively installer-owned MUST still be described and verified separately before it may be removed.
- **FR-023**: Reinstall after preserve MUST recover prior task state, while reinstall after a successful wipe MUST start without prior application-owned state.
- **FR-024**: S039 automated verification MUST cover requirements quality, source authoring, built-package tables, helper safety, all silent shortcut states, preserve/wipe, upgrade, repair, invalid values, current-profile cleanup, locked cleanup, and unattended non-launch behavior.
- **FR-025**: S039 native silent evidence MUST record package identity and hash, selected modes, execution context, shortcut state, populated data inventory, cleanup results, and unaffected control files. Downstream #94 MUST add attended action, confirmation, cancellation, multiple-profile, and interactive-user observations before #97 or #98 closes.
- **FR-026**: The complete S039 implementation MUST pass the repository's canonical format, vet, lint, race, GUI, coverage, documentation, and automation gates without weakening existing safety-critical coverage.

### Key Entities

- **Shortcut selection**: The independent installed state of the Start Menu and desktop shortcuts across initial setup, maintenance, upgrade, and removal.
- **Completion selection**: Two independent, session-only choices that can launch the desktop application and project documentation after a successful fresh attended install.
- **Removal mode**: Preserve or wipe, defaulting to preserve and becoming destructive only after explicit confirmation or the exact managed opt-in.
- **Removal inventory**: The user-visible classification of software state, machine application data, supported-user preference data, and separately governed security state.
- **Owned data root**: A canonical application-specific path derived from a trusted Windows location and eligible for deletion only after scope and traversal-safety validation.
- **Cleanup result**: Complete, partial, refused, or not requested, with evidence describing any declared owned state that remains.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All four fresh-install shortcut combinations result in exactly the two selected/absent states, and maintenance can transition between all combinations without reinstalling unrelated software.
- **SC-002**: The built installer represents all four completion-action combinations independently and limits them to fresh full-UI success; 100 percent of unattended and non-fresh automated flows execute neither target. #94 supplies native interactive execution proof.
- **SC-003**: Preserve-mode removal retains 100 percent of populated declared application-data bytes and a reinstall reads the prior task set successfully.
- **SC-004**: Successful hosted wipe-mode removal leaves zero files or directories in every populated declared application-owned root discovered on the disposable test machine; #94 extends this to multiple genuine profiles on the clean candidate host.
- **SC-005**: Wipe-mode tests leave 100 percent of out-of-scope control files unchanged and refuse every unsafe candidate-path class in the requirements.
- **SC-006**: Cancel, upgrade, repair, rollback, and failed-removal tests cause zero application-data deletions.
- **SC-007**: Every partial or refused cleanup produces a non-success cleanup result and names the remaining owned scope in retained diagnostic evidence.
- **SC-008**: Setup and removal documentation lets an administrator predict every installed, removed, preserved, and optionally erased item, including exact unattended choices, without consulting source code.
- **SC-009**: All eight canonical repository gates pass, and the implementation leaves #94 with a reusable lifecycle matrix for exact-candidate release proof.

## Assumptions

- The supported Windows installer remains per-machine and therefore uses all-users Start Menu and desktop shortcut locations.
- The normal desktop application and documentation actions are useful together, so their checkboxes are independent rather than mutually exclusive.
- “All application-related user data” means the application-owned machine root plus the application-specific preference roots of safely registered local Windows profiles; it does not mean arbitrary exports or user-authored files outside those roots.
- Preserve is the only safe implicit uninstall mode. Destructive unattended cleanup is opt-in and intentionally unsuitable for generic package-manager removal unless the exact wipe property is supplied.
- Product data removal is irreversible and therefore occurs only after successful software-removal commit, rather than pretending deleted user data can participate in installer rollback.
- Exact release-candidate MSI proof, attended installer behavior, multi-profile interactive validation, native window sizing, connection-error observation, and release blocking remain the responsibility of #94 after S039 implementation lands.

## Dependencies and Traceability

- Implements the code and automated proof for #97 and #98; both issues remain open until #94 completes their attended/native acceptance evidence.
- Child of coordinator #96.
- Supplies installer lifecycle inputs to #94 but does not close or partially claim that issue.
- Retained preference state remains relevant to the completed startup-window correction in #89.

## Out of Scope

- Publishing or tagging a release, merging the pull request, or completing #94's exact-candidate release gate.
- Changing GUI startup sizing, connection-recovery behavior, scheduling, IPC authorization, task execution, or persistence schemas.
- Deleting user exports, arbitrary files outside declared application roots, unknown/orphaned profiles, or shared/uncertain operating-system group state.
- Adding a general-purpose filesystem cleanup facility or exposing destructive cleanup through the normal application UI or public CLI.
