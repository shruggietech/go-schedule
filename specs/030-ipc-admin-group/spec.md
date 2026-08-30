# Feature Specification: Dedicated IPC Administrative Group

**Feature Branch**: `codex/030-ipc-admin-group`

**Created**: 2026-08-30

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: review branch `codex/030-ipc-admin-group`; local verification completed 2026-08-30

**Input**: Issue #13: narrow local IPC access control to a dedicated administrative group on Windows, Linux, and macOS, including installer and operator guidance.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Restrict privileged scheduling access (Priority: P1)

As an operator, I can limit daemon management to a named local administrative group so an unrelated local account cannot submit commands that run with the daemon's privileges.

**Why this priority**: Reaching the local endpoint grants control over arbitrary command execution. The access boundary must therefore be narrow by default and consistent across supported platforms.

**Independent Test**: Start the daemon with an existing administrative group, inspect the endpoint permissions, and verify that the daemon identity, platform administrators, and group members retain access while an unrelated local identity does not.

**Acceptance Scenarios**:

1. **Given** the default administrative group exists, **When** the daemon creates its local endpoint, **Then** access is limited to privileged system identities and members of that group.
2. **Given** a custom existing administrative group, **When** the daemon starts, **Then** the same restricted policy names the custom group rather than the default group.
3. **Given** a configured group that cannot be resolved, **When** the daemon starts, **Then** startup fails before serving requests and the error names `admin_group` and the missing group.

---

### User Story 2 - Choose compatibility mode deliberately (Priority: P1)

As a single-user operator who cannot provision a local group, I can explicitly select the historical broad local-access mode and can see that weaker choice in daemon logs.

**Why this priority**: Existing foreground and portable use cases need an escape hatch, but a misspelled or missing security group must never activate it accidentally.

**Independent Test**: Configure an empty administrative-group value, start the daemon, and verify that it serves with the compatibility policy and emits one structured warning describing the broader local access.

**Acceptance Scenarios**:

1. **Given** `admin_group` is explicitly empty, **When** the daemon starts, **Then** it retains the historical platform access policy and logs a warning that compatibility mode is active.
2. **Given** `admin_group` is non-empty but invalid, **When** resolution fails, **Then** the daemon does not fall back to compatibility mode.
3. **Given** restricted mode is active, **When** startup completes, **Then** the selected group and restricted access mode are visible in structured startup logging.

---

### User Story 3 - Install into a working secure default (Priority: P2)

As an installer or service operator, I receive clear platform-specific steps that provision the default group, enroll the intended administrator, and preserve or remove membership intentionally over upgrades and uninstalls.

**Why this priority**: A secure default that prevents the installed GUI and CLI from connecting is not operationally complete.

**Independent Test**: Validate the Windows installer contract and follow the documented Linux and macOS setup commands to confirm the default group exists before the service starts.

**Acceptance Scenarios**:

1. **Given** a new Windows installation, **When** installation completes, **Then** the default local group exists and the installing user is enrolled before the service starts.
2. **Given** an upgrade, **When** the default group already exists, **Then** installation preserves the group and membership without duplication or failure.
3. **Given** an uninstall, **When** the application is removed, **Then** the security group and its membership are preserved to avoid silently changing administrator access for retained scheduler data.
4. **Given** Linux or macOS service setup, **When** the operator follows the installation guide, **Then** group creation and membership occur before service startup and the need to refresh the login session is explained.

### Edge Cases

- The configured group name is empty, whitespace-only, unknown, duplicated across local and domain scopes, or resolves to a non-group identity.
- The default endpoint parent already exists with broader permissions from an older release, or a custom endpoint names an existing shared parent that the daemon must not rewrite.
- A stale endpoint exists from a previous daemon run.
- Endpoint permission application fails after the listener is created; the daemon must close and remove the endpoint rather than leave a partially configured listener.
- The daemon runs in a foreground development session whose user cannot change group ownership.
- A custom endpoint lives outside the default data directory.
- A Windows group name contains characters that cannot safely appear in an access-control descriptor.
- An installer upgrade or repair runs when the group or membership already exists.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The built-in configuration MUST select `goschedadmin` as the default administrative group.
- **FR-002**: A non-empty `admin_group` MUST select restricted mode and MUST be resolved to an operating-system group before the daemon begins serving requests.
- **FR-003**: Failure to resolve a non-empty configured group MUST fail daemon startup with an actionable error naming `admin_group` and the requested value.
- **FR-004**: Restricted mode MUST grant endpoint access only to the daemon identity, platform administrative identities required for recovery, and members of the configured group.
- **FR-005**: Restricted mode MUST prevent unrelated local identities from opening the endpoint through its parent directory; it MUST secure the managed default parent, create a missing custom parent securely, and reject rather than rewrite an unsafe existing custom parent.
- **FR-006**: The daemon MUST apply and verify restricted endpoint permissions before announcing readiness or accepting requests.
- **FR-007**: If restricted permission setup fails after endpoint creation, the daemon MUST close the listener and remove the newly created endpoint where the platform exposes it as a filesystem entry.
- **FR-008**: An explicitly empty `admin_group` MUST activate the historical broad local-access compatibility policy.
- **FR-009**: Compatibility mode MUST emit a structured warning at startup and MUST NOT be activated by group lookup or permission errors.
- **FR-010**: Restricted startup MUST emit structured information identifying the selected access mode and configured group without exposing unrelated account data.
- **FR-011**: The Windows installer MUST create or reuse the default group and enroll the installing interactive user before starting the daemon service.
- **FR-012**: Windows installer upgrade, repair, and uninstall behavior MUST be idempotent and MUST preserve the group and membership on uninstall.
- **FR-013**: Linux and macOS installation guidance MUST provision the default group, enroll the intended administrator, explain session refresh, and do so before service startup.
- **FR-014**: Operator documentation MUST explain restricted mode, explicit compatibility mode, custom group configuration, startup failures, and the privilege implications of group membership.
- **FR-015**: Automated tests MUST cover restricted and compatibility policy construction, missing-group fail-closed behavior, Unix permission application and cleanup, Windows descriptor construction, installer provisioning, and structured startup logging.
- **FR-016**: The change MUST retain existing endpoint names and client connection behavior for authorized users.
- **FR-017**: The implementation MUST remain compatible with daemon and CLI cross-compilation for all currently tested targets.
- **FR-018**: The S030 pull request MUST use `Closes #13` because this slice completes the issue across runtime, installer, tests, and documentation.

### Key Entities

- **Administrative group**: The configured local operating-system group whose members may manage the daemon.
- **IPC access mode**: Either restricted group access or explicit broad compatibility access.
- **Endpoint policy**: The platform-specific permissions applied before the daemon serves local requests.
- **Provisioned membership**: The durable relationship between an intended operator account and the default administrative group.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On every supported platform, committed access-policy contracts exclude unrelated local identities while preserving the daemon identity, recovery administrators, and configured group.
- **SC-002**: One hundred percent of missing or invalid non-empty group cases fail before daemon readiness, with an error that names the affected setting and value.
- **SC-003**: One hundred percent of explicit compatibility-mode startups emit a warning, while restricted startups identify their selected group without a warning.
- **SC-004**: The Windows installer source and pinned-tool build satisfy committed fresh-install, upgrade, repair, and uninstall group-provisioning contracts without requiring manual machine observation.
- **SC-005**: An operator can locate and complete the group setup for any supported platform from its installation guide before the first service start.
- **SC-006**: All eight canonical repository verification gates pass, and the existing core-package coverage floor remains satisfied.

## Assumptions

- The daemon continues to run with system-level privileges when installed as a service; group membership therefore grants equivalent scheduler-management power.
- The default Windows installer has access to the elevated installing context and can determine the interactive user who initiated installation.
- Linux and macOS remain archive-based distributions, so their group provisioning is documented operator work rather than a new package-manager integration.
- Existing custom endpoint locations are supported, but an existing custom parent must already be dedicated to the daemon with the required group ownership and mode; the daemon will verify it but will not take ownership of an unrelated directory.
- Empty `admin_group` is the sole compatibility opt-in. Whitespace-only values are invalid configuration rather than an alias for empty.

## Out of Scope

- Remote authentication, network listeners, per-task authorization, or multiple privilege tiers.
- Changing the privilege identity under which scheduled commands execute.
- Creating Linux distribution packages or a macOS graphical installer.
- Deleting an existing administrative group or removing its members during uninstall.
- Signing or notarizing release artifacts.
