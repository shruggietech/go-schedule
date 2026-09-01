# Feature Specification: Windows Installed Core Recovery

**Feature Branch**: `codex/038-windows-installed-core-recovery`

**Created**: 2026-09-01

**Status**: In Progress

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: implementation and verification in progress on `codex/038-windows-installed-core-recovery`

**Input**: GitHub issues #90 and #93: restore ordinary-user IPC access and real service-hosted task execution on installed Windows systems.

## User Scenarios & Testing

### User Story 1 - Ordinary Desktop Access (Priority: P1)

A user enrolled in the local `goschedadmin` group launches the installed GUI or CLI from a normal, non-elevated desktop session and reaches the LocalSystem-hosted daemon without signing out, elevating the application, or broadening access to every authenticated user.

**Why this priority**: The installed application is unusable when its intended operator cannot cross the local IPC boundary. Elevation succeeding only proves that the Built-in Administrators ACE works; it is not an acceptable routine access path.

**Independent Test**: Create a restricted named pipe using a local alias whose direct user member SID is absent from the current token, then prove that the direct member connects while an unrelated user SID is not granted. Confirm the same descriptor retains SYSTEM, Built-in Administrators, and the configured alias.

**Acceptance Scenarios**:

1. **Given** a standard desktop token whose user is a direct member of `goschedadmin` but whose token omits the alias SID, **When** the daemon creates its restricted pipe, **Then** that user's SID receives read/write access and the client connects without elevation.
2. **Given** the same user launches the GUI from the Start Menu or the CLI from a normal shell, **When** the installed service is healthy, **Then** the complete UI loads and `gosched health` succeeds without a permanent elevation requirement.
3. **Given** an unrelated local user who is not a direct group member, **When** that user attempts to connect, **Then** the restricted pipe denies access.
4. **Given** a configured group with zero direct members, **When** the daemon starts, **Then** the descriptor still authorizes SYSTEM, Built-in Administrators, and the configured group without granting Authenticated Users.
5. **Given** group membership cannot be enumerated, **When** the daemon starts, **Then** startup fails with an actionable authorization error rather than silently broadening access.

---

### User Story 2 - Service-Hosted Task Execution (Priority: P1)

An ordinary authorized user creates a Windows task through the installed product, triggers it manually, and observes the same task execute successfully on schedule through the real LocalSystem-hosted executor.

**Why this priority**: Scheduling a command is the product's core promise. Existing in-process tests substitute an always-successful runner and therefore cannot detect a broken service process boundary.

**Independent Test**: Against a service-hosted daemon, create a task that invokes an absolute Windows system executable with deterministic arguments and a writable marker path. Prove manual and scheduled runs produce exit code 0, expected captured output, and distinct marker side effects.

**Acceptance Scenarios**:

1. **Given** the installed daemon runs as LocalSystem, **When** a task invokes an absolute system command with explicit arguments and a writable marker path, **Then** a manual Run now request produces a successful run record, exit code 0, expected output, and the marker side effect.
2. **Given** the same task has a near-term recurring schedule, **When** its due time arrives, **Then** the scheduler dispatches the production executor and records an independently observable successful scheduled run.
3. **Given** a command starts successfully and exits nonzero, **When** the run completes, **Then** history retains the exit code and captured output and classifies it separately from a process-start failure.
4. **Given** the executable cannot be started, **When** execution fails, **Then** history identifies the command boundary and retains the Windows error without exposing environment values or secrets.
5. **Given** any service-hosted task executes, **When** the child process starts, **Then** no console window or interactive prompt appears and the daemon remains isolated in session 0.

### Edge Cases

- Direct local-group entries can be users, nested groups, deleted accounts, or duplicate SIDs. Only resolvable direct user members are expanded; the configured group ACE remains authoritative for fresh tokens and nested membership.
- Member SID ordering from Windows APIs is not stable. Generated descriptors must be deterministic and de-duplicate principals.
- SDDL construction must reject empty, malformed, or delimiter-bearing SID text rather than interpolating unchecked data.
- A task executable may be missing, inaccessible to LocalSystem, or represented by a relative/bare name whose service PATH differs from the interactive user's PATH.
- A working directory may be absent or inaccessible from session 0.
- Captured stdout/stderr remains bounded by the existing output cap.
- Cancellation and timeout errors must remain distinguishable from ordinary exit-code failure when no output is available.

## Requirements

### Functional Requirements

- **FR-001**: Restricted Windows IPC MUST grant full access to SYSTEM and Built-in Administrators and read/write access to the configured `admin_group` SID.
- **FR-002**: Restricted Windows IPC MUST additionally grant read/write access to each unique, resolvable direct user member SID returned for the configured local group.
- **FR-003**: Direct-member expansion MUST NOT grant access to Authenticated Users, Everyone, an unrelated local user, or a user absent from the configured group.
- **FR-004**: The configured group ACE MUST remain present so fresh tokens and supported nested-group membership continue to work through Windows token authorization.
- **FR-005**: Group lookup or enumeration failure MUST fail closed with an actionable startup error; deleted or unresolvable individual member records MAY be skipped only when Windows supplies no usable SID for that record.
- **FR-006**: Generated SDDL MUST be deterministic, de-duplicate SIDs case-insensitively, and validate every interpolated SID using Windows SID parsing before the listener opens.
- **FR-007**: IPC access reporting MUST identify restricted mode and the configured group without logging member identities.
- **FR-008**: Production daemon wiring MUST use the real executor in service-hosted integration coverage; an injected fake runner cannot satisfy Windows execution verification.
- **FR-009**: The Windows execution probe MUST create a task through the public API or CLI, invoke an absolute inbox system executable, use explicit arguments and a deterministic service-writable marker under the daemon data directory, retain a copy with the probe evidence, and redact environment values.
- **FR-010**: The probe MUST prove both manual and scheduled execution with distinct run records, exit code 0, expected output, and marker side effects.
- **FR-011**: The probe MUST exercise a deliberate nonzero exit and distinguish it from a process-start failure with retained output or operating-system error detail.
- **FR-012**: Process-start failures MUST record the executable boundary and wrapped OS error while avoiding arguments, stdin, and environment values that may contain secrets.
- **FR-013**: Windows child processes MUST retain the existing hidden-process guarantee and disable interactive prompts.
- **FR-014**: The daemon MUST remain a LocalSystem service; this slice MUST NOT switch it to the interactive user or add privilege escalation.
- **FR-015**: Native verification MUST record Windows version, candidate hash/version, service account/state, installing identity, group SID and direct membership, current token SIDs, generated pipe descriptor, ordinary connection result, task definition with redacted environment values, run records, marker evidence, and correlated daemon logs.
- **FR-016**: S036's exclusion of named-pipe ACL and installer authorization changes is superseded only to the extent required to resolve #90; its single in-frame recovery presentation remains intact.
- **FR-017**: The release-wide installed smoke gate tracked by #94 remains outside S038 and MUST NOT be represented as completed by this slice.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every test descriptor contains exactly one ACE for each unique authorized SID and zero Authenticated Users or Everyone ACEs in restricted mode.
- **SC-002**: A standard user directly enrolled in `goschedadmin` reaches the real named pipe without elevation even when the current token omits the alias SID; an unrelated user remains denied.
- **SC-003**: A service-hosted task completes one manual and one scheduled run with exit code 0, expected output, and two independently verified marker effects.
- **SC-004**: One deliberate exit-code failure and one process-start failure produce diagnostically distinct run records without secret-bearing values.
- **SC-005**: All eight canonical repository gates pass with no regression to IPC restriction, hidden Windows child processes, scheduling semantics, or GUI recovery behavior.

## Assumptions

- The installer continues to enroll the intended account as a direct member of the local `goschedadmin` alias.
- Windows evaluates a user-SID ACE against the standard token even when UAC or session staleness omits the newly added local alias SID.
- Direct-member expansion is evaluated at daemon startup; membership changes become effective after service restart or the normal fresh-token path.
- `%SystemRoot%\\System32\\cmd.exe` is present on supported Windows hosts and can provide deterministic output and marker effects when invoked with explicit noninteractive arguments.
- A clean release-equivalent MSI walkthrough may require a disposable Windows host. Repository tests and scripts must report unavailable native evidence honestly rather than substitute a fake backend.

## Out of Scope

- Granting access to Authenticated Users or Everyone, disabling UAC, or requiring permanent application elevation.
- Recursively expanding nested groups or implementing domain-controller membership traversal.
- Running tasks as the interactive desktop user, loading an interactive profile, or supporting per-task Windows `run_as`.
- General shell parsing of a command-line string. Command and arguments remain separate structured fields.
- Publishing a release, merging the pull request, or completing the broader release-smoke gate in #94.

## Clarifications

### Session 2026-09-01

- Q: Which issues belong in S038? -> A: Bundle #90 and #93 because together they restore the installed Windows core path; defer the broader release gate #94 to S039.
- Q: May ordinary access be restored by granting Authenticated Users? -> A: No. Preserve the configured-group boundary and add only verified direct user member SIDs needed for stale or filtered desktop tokens.
- Q: What constitutes task-execution proof? -> A: Real service-hosted manual and scheduled runs with output, exit code, marker side effects, and a failing control; fake runners do not qualify.
- Q: Should tasks run as the desktop user? -> A: No. Keep LocalSystem service isolation and make service-context behavior explicit and verifiable.
