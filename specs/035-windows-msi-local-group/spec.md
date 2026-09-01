# Feature Specification: Windows MSI Local-Group Recovery

**Feature Branch**: `codex/083-fix-msi-local-group`

**Created**: 2026-08-31

**Status**: In Progress

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Not delivered

**Input**: Issue #83: repair the released Windows installer regression that
routes local administrative-group provisioning through an unelevated domain
operation and aborts with access denied.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Complete a secure fresh Windows install (Priority: P1)

As a Windows operator, I can approve an elevated installation and receive a
working service, local administrative group, and intended account membership
without the installer rolling back.

**Why this priority**: The published installer currently fails on the primary
Windows installation path, leaving the product unavailable.

**Independent Test**: On a clean supported Windows 11 host, install a candidate
package with verbose logging and verify successful completion, local group and
membership creation, service startup, and client access after session refresh.

**Acceptance Scenarios**:

1. **Given** a clean supported Windows 11 host and an interactive administrator,
   **When** the candidate installer runs after UAC approval, **Then** installation
   completes without access-denied or group-creation errors.
2. **Given** no pre-existing `goschedadmin` group, **When** installation reaches
   service startup, **Then** the local group exists and contains the interactive
   installing account before `goschedd` starts.
3. **Given** the installing account refreshes its login session, **When** the GUI
   or CLI connects, **Then** the restricted daemon is reachable.

---

### User Story 2 - Preserve administrative state across MSI lifecycles (Priority: P1)

As a Windows operator, I can repair, reinstall, upgrade, and uninstall without
duplicate groups, lost memberships, or failures caused by existing security
state.

**Why this priority**: A one-time fresh-install fix is incomplete if normal MSI
lifecycle operations corrupt or reject the durable access boundary.

**Independent Test**: Exercise repair or same-version reinstall, upgrade from
the prior release, and uninstall on a disposable host while recording group,
membership, service, exit-code, and verbose-log evidence after each phase.

**Acceptance Scenarios**:

1. **Given** the group and membership already exist, **When** repair or a
   same-version reinstall runs, **Then** it completes idempotently without
   replacing or duplicating them.
2. **Given** the prior Windows release is installed, **When** the candidate
   package upgrades it, **Then** the existing group and membership remain and
   the service returns to a working state.
3. **Given** the candidate package is installed, **When** it is uninstalled,
   **Then** the service and application files are removed while `goschedadmin`
   and its membership remain.

---

### User Story 3 - Diagnose genuine provisioning restrictions (Priority: P2)

As an operator whose machine policy blocks group provisioning, I can find the
verbose installer log and distinguish a genuine account or policy restriction
from the released authoring defect.

**Why this priority**: Provisioning can still fail legitimately on managed
hosts, and operators need actionable evidence instead of a generic rollback.

**Independent Test**: Follow the Windows troubleshooting guide to locate or
produce a verbose installer log and identify the failed operation and exit code.

**Acceptance Scenarios**:

1. **Given** installation fails on a managed host, **When** the operator opens
   the troubleshooting guidance, **Then** it provides an exact verbose-log
   command and points to the provisioning failure evidence.
2. **Given** a lifecycle verification run cannot meet a prerequisite, **When**
   evidence is written, **Then** the prerequisite is reported as unavailable or
   failed rather than as a successful lifecycle.

### Edge Cases

- `goschedadmin` or the intended membership already exists before installation.
- The installing account is local, domain-backed, or represented by an account
  name whose domain differs from the computer name.
- Group creation succeeds but membership enrollment or service startup fails.
- Repair or same-version reinstall runs after a partial earlier attempt.
- Upgrade starts from v0.9.0 with preserved program data and administrative
  state.
- Uninstall runs after a failed or interrupted install.
- Host policy genuinely denies local group or membership changes.
- A lifecycle run starts from a contaminated machine or lacks elevation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Windows installer MUST provision `goschedadmin` as a
  machine-local group through an elevated local-group operation.
- **FR-002**: The local-group declaration MUST NOT provide a non-empty domain
  value that selects domain-group behavior.
- **FR-003**: The installer MUST create or reuse `goschedadmin` before enrolling
  the interactive installing account and before starting `goschedd`.
- **FR-004**: The installer MUST retain the existing interactive-account
  identity behavior while correcting only the group operation classification.
- **FR-005**: Fresh installation MUST complete after UAC approval without
  access-denied HRESULT `0x80070005`, installer error 26421, or rollback caused
  by administrative-group creation.
- **FR-006**: Repair and same-version reinstall MUST be idempotent when the
  group and membership already exist.
- **FR-007**: Upgrade from the prior Windows release MUST preserve and reuse the
  existing group and membership.
- **FR-008**: Uninstall MUST preserve `goschedadmin` and its membership while
  removing the product service and files.
- **FR-009**: Automated source contracts MUST reject any non-empty local-group
  domain value and retain the group, membership, preservation, feature-linkage,
  and service-ordering requirements.
- **FR-010**: The Windows lifecycle verifier MUST record artifact identity,
  installer exit codes, verbose-log locations, group existence, membership,
  service state and ordering evidence, reinstall or repair, upgrade, uninstall,
  and final state.
- **FR-011**: The lifecycle verifier MUST refuse contaminated or unelevated
  baselines and MUST never represent unavailable runtime evidence as passing.
- **FR-012**: Windows installation guidance MUST provide an actionable verbose
  logging command and identify group, membership, service, and access-denied
  evidence relevant to provisioning failures.
- **FR-013**: The fix MUST preserve the restricted IPC security design and MUST
  NOT broaden endpoint authorization or change non-Windows provisioning.
- **FR-014**: The implementation MUST pass all eight canonical repository
  verification gates without weakening installer or IPC contracts.

### Key Entities

- **Local administrative group**: The durable machine-local `goschedadmin`
  security principal used by restricted IPC authorization.
- **Installing account membership**: The preserved relationship enrolling the
  interactive account in `goschedadmin`.
- **Installer lifecycle run**: A fresh install, repair or reinstall, upgrade,
  and uninstall sequence executed on a disposable supported Windows host.
- **Lifecycle evidence**: Artifact identity, verbose logs, exit codes, security
  state, service state, and final cleanup observations recorded without
  overstating unavailable results.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A candidate MSI completes one clean elevated Windows 11 fresh
  install with exit code 0 or restart-required success and zero occurrences of
  `0x80070005` or error 26421 in its verbose log.
- **SC-002**: In 100 percent of recorded successful install phases,
  `goschedadmin` exists, the intended account is a member, and `goschedd` starts
  only after provisioning succeeds.
- **SC-003**: One repair, one same-version reinstall, and one upgrade from
  v0.9.0 complete without duplicating or replacing the administrative group.
- **SC-004**: One uninstall preserves the group and intended membership while
  leaving no installed product, service, install directory, or product PATH
  entry.
- **SC-005**: Every source-contract mutation that adds any non-empty group
  domain is rejected before packaging.
- **SC-006**: The lifecycle evidence contains every required observation or an
  explicit unavailable or failed result; zero missing observations are counted
  as passes.
- **SC-007**: After session refresh, one non-elevated client probe reaches the
  restricted daemon through the intended account's group membership.
- **SC-008**: All eight canonical verification gates pass and core coverage
  remains at or above its existing floor.

## Clarifications

### Session 2026-08-31

- Q: Which identity behavior changes? A: Only the group declaration loses its
  domain; the installing-user identity contract remains unchanged.
- Q: Can static source or MSI-table inspection satisfy runtime acceptance?
  A: No; native lifecycle evidence is required and unavailable evidence remains
  explicitly unresolved.
- Q: What administrative state survives uninstall? A: Both `goschedadmin` and
  its existing membership remain preserved.

## Assumptions

- Issue #83's reproduced v0.9.0 failure and root-cause evidence are the accepted
  red baseline; this slice does not need to reinstall the known-bad artifact to
  rediscover it.
- WiX's supported machine-local group behavior is selected by omitting the
  group domain rather than formatting the computer name into that field.
- A disposable supported Windows 11 environment, the v0.9.0 MSI, and a
  candidate MSI are available for the native lifecycle run.
- Success exit code 3010 remains acceptable when Windows requires a restart,
  provided all required state and logs are recorded honestly.

## Out of Scope

- Redesigning local IPC authorization or broadening its authorized identities.
- Changing Linux or macOS group provisioning.
- Removing `goschedadmin` or its members during uninstall.
- Signing, publishing, or cutting a replacement release artifact.
- Treating source, XML, or MSI-table inspection as a substitute for native
  installation evidence.
