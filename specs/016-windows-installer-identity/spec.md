# Feature Specification: Windows Installer Identity and Verification

**Feature Branch**: `codex/016-windows-installer-identity`

**Created**: 2026-08-27

**Status**: Implemented

**Delivery**: [PR #48](https://github.com/shruggietech/go-schedule/pull/48)

**Input**: User description: "Turn the remaining Windows installer identity and clean-install P1 issues into one focused Spec-Kit slice and run it end-to-end under autopilot."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Recognize the installed application (Priority: P1)

As a Windows operator, I can recognize go-schedule by its official brand mark in the Start Menu and in the installed-apps list instead of seeing a generic placeholder.

**Why this priority**: These are persistent operating-system surfaces created by the installer. A generic icon makes the installed product look incomplete and makes it harder to distinguish from unrelated applications.

**Independent Test**: Install a candidate package in a disposable Windows environment and observe that both installed surfaces use the same official mark as the application.

**Acceptance Scenarios**:

1. **Given** a clean supported Windows environment, **When** go-schedule is installed, **Then** its Start Menu shortcut shows the official brand mark.
2. **Given** go-schedule is installed, **When** the operator views installed applications, **Then** the go-schedule entry shows the official brand mark.
3. **Given** go-schedule is upgraded, **When** either installed surface is viewed after normal shell refresh, **Then** it still shows the official mark rather than a stale or generic image.

---

### User Story 2 - Recognize the running application (Priority: P1)

As a Windows operator, I can recognize the running go-schedule GUI by its official mark in the window title area and taskbar.

**Why this priority**: The running application must present the same identity as the installer-created surfaces, especially when several applications are open.

**Independent Test**: Launch the installed GUI in a disposable Windows environment and observe both native running-application surfaces.

**Acceptance Scenarios**:

1. **Given** go-schedule was installed from a candidate package, **When** the GUI is launched, **Then** the window and taskbar use the official brand mark.
2. **Given** the shell icon cache initially contains an older go-schedule icon, **When** the documented normal refresh procedure is applied, **Then** the current official mark appears without requiring a product reinstall.

---

### User Story 3 - Trust command availability across the install lifecycle (Priority: P1)

As a Windows operator using a published installer, I can invoke `gosched` from a newly opened shell after installation, and uninstalling or upgrading does not leave a missing, duplicate, or malformed machine search-path entry.

**Why this priority**: Every documented command depends on this behavior, and a development machine can mask a broken installer by already having the install directory on its search path.

**Independent Test**: In a clean Windows environment with no prior go-schedule installation or search-path entry, exercise install, new-shell command use, same-version reinstall or upgrade, and uninstall while recording the machine search path at each stage.

**Acceptance Scenarios**:

1. **Given** a clean Windows environment where `gosched` cannot be resolved, **When** the latest published installer is installed and a new shell opens, **Then** `gosched --version` and `gosched task list` resolve by name.
2. **Given** go-schedule is installed, **When** it is reinstalled or upgraded, **Then** exactly one install-directory entry exists in the machine search path.
3. **Given** go-schedule is installed, **When** it is uninstalled, **Then** its search-path entry is removed without changing unrelated entries or leaving an empty segment.

---

### User Story 4 - Preserve the installer contract (Priority: P2)

As a maintainer, I receive a clear local failure before publication if the installer loses its brand-icon or search-path declarations, and I have a repeatable evidence procedure for observations that require a clean Windows desktop.

**Why this priority**: Static package declarations can be guarded on every change, while native shell rendering and lifecycle behavior require a separate environment. Keeping those evidence classes explicit prevents both regressions and false claims of verification.

**Independent Test**: Mutate each protected declaration in an isolated fixture and confirm the maintainer check rejects it with a specific message; review the clean-environment procedure and evidence record for all native observations.

**Acceptance Scenarios**:

1. **Given** a valid installer definition, **When** the maintainer check runs, **Then** all required icon and search-path relationships pass.
2. **Given** any one protected declaration is missing or points at the wrong brand asset, **When** the check runs, **Then** it fails and names the broken contract.
3. **Given** no clean Windows environment is available during a work slice, **When** verification is reported, **Then** native observations are marked unavailable with the exact prerequisite and their issues remain open.

### Edge Cases

- A four-space-indented command or unrelated executable resource must not be mistaken for installer evidence.
- Windows may retain a stale shell icon after an upgrade; the verification record must distinguish a cache refresh from a package repair.
- The installer file itself keeps Windows Installer's generic file icon; this slice covers installed-product surfaces, not the `.msi` file in Explorer.
- A search path that already contains the install directory is not a clean baseline and cannot prove installation behavior.
- An already-open shell retains its inherited environment and cannot be used to disprove a successful machine search-path update.
- A local candidate package can prove future package behavior but cannot by itself satisfy an issue that explicitly requires a published artifact.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The installed Start Menu shortcut MUST explicitly use the official go-schedule brand mark.
- **FR-002**: The installed-apps entry MUST explicitly use the official go-schedule brand mark.
- **FR-003**: The shortcut and installed-apps entry MUST derive from one canonical icon asset rather than independently maintained copies.
- **FR-004**: The Windows GUI artifact MUST carry the official brand mark for native title-area and taskbar use.
- **FR-005**: The installer MUST add its install directory to the machine search path on install, preserve a single entry across reinstall or upgrade, and remove only that entry on uninstall.
- **FR-006**: The maintainer check MUST reject a missing icon declaration, a missing installed-apps icon relationship, a missing shortcut icon relationship, or any relationship that points at a noncanonical asset.
- **FR-007**: The maintainer check MUST continue to reject a missing, per-user, permanent, replacing, or misdirected search-path declaration.
- **FR-008**: Every protected installer-contract failure MUST identify the specific missing or invalid relationship.
- **FR-009**: Regression coverage MUST prove that each newly protected icon relationship fails independently before the valid definition passes.
- **FR-010**: A clean-environment procedure MUST cover negative baseline, install, new-shell command use, installed icon surfaces, running icon surfaces, reinstall or upgrade, uninstall, and evidence capture.
- **FR-011**: Verification reporting MUST separate definition checks, built artifact checks, published-artifact checks, and native desktop observations.
- **FR-012**: An issue requiring native or published-artifact evidence MUST remain open when that evidence was not obtained, even if its declarations and regression checks pass.
- **FR-013**: The work slice MUST NOT publish a release, alter signing policy, introduce a bootstrap executable, redesign the brand, or change application behavior outside Windows identity and install lifecycle evidence.
- **FR-014**: The existing service registration, installation scope, install directory, application data preservation, and supported Windows architecture MUST remain unchanged.

### Key Entities

- **Canonical brand mark**: The single multi-resolution Windows icon asset used by installer-created surfaces and the GUI artifact.
- **Installer identity contract**: The declared relationships among the brand mark, Start Menu shortcut, installed-apps entry, GUI artifact, and package.
- **Clean-environment evidence**: A dated record naming Windows build, artifact version and origin, baseline state, lifecycle steps, observations, and any unavailable prerequisite.
- **Verification status**: One of proven, failed, or unavailable for each evidence class; unavailable is never treated as proven.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Both installer-created Windows identity surfaces specify the same official brand mark, with zero independently maintained icon copies.
- **SC-002**: Removing or misdirecting any one of the three installer icon relationships produces one or more specific maintainer-check failures.
- **SC-003**: The valid installer definition passes all existing and new identity and search-path contract checks.
- **SC-004**: The clean-environment procedure covers 100 percent of the eight lifecycle and observation stages named in FR-010.
- **SC-005**: Every P1 issue in scope has an evidence record that states proven, failed, or unavailable and cites the exact artifact and environment used.
- **SC-006**: No scoped issue is marked complete without all evidence its own acceptance criteria require.
- **SC-007**: All eight repository verification gates pass with no product behavior, performance budget, or coverage regression.

## Assumptions

- Windows continues to use the package's primary icon relationship for the installed-apps entry and the shortcut's explicit icon relationship for the Start Menu surface.
- The existing official multi-resolution Windows icon is the canonical asset; brand redesign and asset regeneration belong to issue #10.
- A disposable virtual machine, Windows Sandbox, or a genuinely clean physical machine is acceptable for native observations. The current development host is not clean because go-schedule is already installed and on the machine search path.
- Enabling virtualization features, rebooting the host, uninstalling the operator's current installation, and cutting a release are outside this slice's authorization.
- Issue #16 requires a published installer. Candidate artifacts may validate future behavior but do not close that issue.
