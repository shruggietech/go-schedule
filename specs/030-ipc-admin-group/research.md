# Research: Dedicated IPC Administrative Group

## R1. Missing-group behavior

**Decision**: Fail daemon startup whenever a non-empty configured group cannot be resolved or applied. Reserve compatibility mode for an explicitly empty `admin_group`.

**Rationale**: Automatic fallback would turn a spelling or provisioning error into privilege escalation. The daemon runs arbitrary commands with service privileges, so a security configuration failure must be visible and fail closed.

**Alternatives considered**:

- Warn and admit all local users, as proposed in issue #13. Rejected because it makes the secure setting advisory.
- Remove compatibility mode entirely. Rejected because foreground and portable use without group-provisioning rights remains a supported practical case.

## R2. IPC package contract

**Decision**: `Listen` accepts the endpoint plus group name and returns both the listener and immutable access-mode evidence.

**Rationale**: Platform code has the facts needed to select and enforce policy. Returning evidence lets the daemon log one consistent structured record without passing logging dependencies into the transport.

**Alternatives considered**:

- Pass the whole application configuration. Rejected because IPC needs only two values.
- Pass a logger into the IPC package. Rejected because it couples transport construction to logging and complicates focused tests.

## R3. Unix ownership and modes

**Decision**: In restricted mode, set the managed default parent or a newly created custom parent to the configured group with mode `0770`, set the socket to the same group with mode `0660`, then read both back. Verify but never rewrite an existing custom parent. In compatibility mode, manage the default parent as `0755`, leave existing custom parents unchanged, and set the socket to `0666`.

**Rationale**: Explicit endpoint modes make behavior independent of the process umask. Tightening the known application-owned directory removes traversal by unrelated users, while verification catches unsafe custom parents without risking permission changes to shared directories.

**Alternatives considered**:

- Change every socket parent. Rejected because a custom endpoint could name `/tmp` or another shared directory the daemon must never take over.
- Change only the socket. Rejected because an existing broad or inaccessible parent can defeat or break the intended boundary.
- Create a second endpoint-only directory. Rejected because it changes established endpoint locations and complicates client discovery.

## R4. Unix testability

**Decision**: Put group lookup, ownership, mode, listen, stat, and cleanup operations behind a small unexported controller passed to the internal implementation.

**Rationale**: Tests can prove ordering and cleanup without root privileges or mutating real host groups. The production path still calls standard operating-system primitives directly.

**Alternatives considered**:

- Package-global function variables. Rejected because they are easier to leak across tests and race under parallel execution.
- Root-only integration tests. Rejected because they are unsuitable for routine local and hosted CI.

## R5. Windows group identity

**Decision**: Resolve the configured name to a Windows SID, accept alias, group, and well-known-group account types, and render only the canonical SID string into the descriptor.

**Rationale**: SID-based descriptors are stable across localization and safe from configuration-string descriptor injection. Rejecting user-account SID types prevents a misleading `admin_group` value.

**Alternatives considered**:

- Put the configured name directly into SDDL. Rejected because names are not valid SDDL trustees and would require unsafe escaping.
- Keep Authenticated Users alongside the group. Rejected because it nullifies the restriction.

## R6. Windows installer capability

**Decision**: Upgrade the release workflow's WiX tool, UI extension, and Util extension together from 5.0.2 to 6.0.2. Use WiX Util `Group`, `User`, and `GroupRef` declarations rather than a custom executable action.

**Rationale**: WiX v6 added local/domain group creation while retaining high source compatibility with v5. The declarative extension owns install sequencing and repair behavior. Official references: [WiX release notes](https://docs.firegiant.com/wix/whatsnew/releasenotes/), [Group](https://docs.firegiant.com/wix/schema/util/group/), [User](https://docs.firegiant.com/wix/schema/util/user/), and [GroupRef](https://docs.firegiant.com/wix/schema/util/groupref/).

**Alternatives considered**:

- Stay on WiX 5 and run `net localgroup` from a custom action. Rejected because quoting, rollback, account lookup, and error handling become repository-owned security code.
- Upgrade to WiX 7. Rejected because v6.0.2 supplies the required feature with less unrelated toolchain movement.

## R7. Installer account selection and lifecycle

**Decision**: Create/reuse a machine-local `goschedadmin` group, enroll Windows Installer's `[LogonUser]`, and preserve both the group and membership during uninstall.

**Rationale**: Windows Installer defines `LogonUser` as the currently logged-on user, which represents the person who initiated the elevated install rather than the LocalSystem service identity. Preserving membership matches retained scheduler data and avoids surprising access loss on upgrade or reinstall. Official reference: [Windows Installer LogonUser property](https://learn.microsoft.com/en-us/windows/win32/msi/logonuser).

**Alternatives considered**:

- Enroll all built-in Administrators. Rejected because the runtime descriptor already grants that identity and it does not help a filtered non-elevated token.
- Delete the group on uninstall. Rejected because recreated groups receive new SIDs and retained data may still reference the old administrative intent.

## R8. Platform packaging boundary

**Decision**: Automate Windows provisioning in the existing MSI; document elevated group setup for Linux and macOS archive/service installs.

**Rationale**: Only Windows has an installer artifact in this repository. Adding Linux packages or a macOS installer is a separate distribution initiative, while runtime enforcement and clear pre-service instructions fully close the authorization gap.
