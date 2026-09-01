# Feature Specification: IPC Access-Denied Recovery

**Feature Branch**: `codex/036-ipc-access-denied-recovery`

**Created**: 2026-08-31

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/036-ipc-access-denied-recovery`; native Windows diagnosis, deterministic recovery, and canonical verification passed 2026-09-01

**Input**: Issue #90: make a fresh Windows GUI launch recoverable when the
daemon named pipe rejects the current logon token, without weakening IPC
authorization.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Recover from an unauthorized login token (Priority: P1)

As a newly installed Windows user whose current login token does not yet contain
`goschedadmin`, I can see one accurate connection state, follow session-refresh
guidance, retry after signing in again, and reach the normal interface without
reinstalling.

**Why this priority**: The released first-run path can trap the user behind an
unbounded sequence of misleading dialogs and deny all useful application access.

**Independent Test**: On native Windows, prove that the installed restricted
pipe denies a verified stale token and selects the matching guidance. Then
deterministically transition the same GUI incident to an authorized backend and
use Retry to restore the interface.

**Acceptance Scenarios**:

1. **Given** the daemon pipe exists but denies the GUI token, **When** model,
   calendar, and event-stream requests fail together, **Then** the application
   frame stays visible and exactly one access-denied recovery state appears.
2. **Given** the account belongs to `goschedadmin` but the process token lacks
   its SID, **When** access is denied, **Then** the recovery state specifically
   tells the user to sign out and back in, and offers Retry and Exit.
3. **Given** the active access-denied incident and a backend that now authorizes
   the request, **When** the user selects Retry, **Then** normal data loads and
   the existing incident is cleared without reinstalling.

---

### User Story 2 - Understand other connection failures (Priority: P1)

As a user unable to connect for another reason, I receive a single actionable
state that distinguishes daemon absence, timeout, and other transport failures
from an API response error.

**Why this priority**: Accurate classification is necessary for a recovery
action to be safe and useful, and prevents access denial from being described
as daemon absence.

**Independent Test**: Inject each failure category and assert its classification,
guidance, and single-state presentation independently.

**Acceptance Scenarios**:

1. **Given** no daemon endpoint exists, **When** a request starts, **Then** the
   state identifies daemon unavailability and suggests checking service status.
2. **Given** a transport timeout, **When** the request deadline expires, **Then**
   the state identifies a timeout and offers Retry and Exit.
3. **Given** the daemon returns an API error response, **When** the GUI receives
   it, **Then** it remains an API error rather than being classified as a
   transport incident.

---

### User Story 3 - Remain stable while disconnected (Priority: P2)

As a disconnected user, I can leave the application open without recurring
dialogs while background reconnect attempts remain bounded and update the one
existing state.

**Why this priority**: Recovery controls are not usable if independent refresh
and stream loops continuously create new interruptions.

**Independent Test**: Keep all startup operations failing across multiple
reconnect intervals and assert one visible incident, bounded backoff, no modal
creation, and clean shutdown.

**Acceptance Scenarios**:

1. **Given** a connection incident is active, **When** tab refreshers and the
   event stream encounter the same failure, **Then** they update or coalesce
   into that incident and create no additional visible errors.
2. **Given** the daemon remains unreachable, **When** background reconnection
   continues, **Then** retries use bounded increasing delays and remain
   interruptible by application exit.

### Edge Cases

- The service is stopped, the pipe does not exist, or the daemon exits during a retry.
- The account is not a member of `goschedadmin`, versus membership exists but
  the current process token is stale.
- Group or token diagnostics cannot be queried; guidance remains accurate and
  does not assert an unverified cause.
- Several request paths report differently wrapped forms of the same operating
  system error at nearly the same time.
- An API response error occurs while a transport incident is active.
- The user selects Retry repeatedly or exits while a retry is in flight.
- Connectivity returns during a scheduled background retry.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The client layer MUST classify daemon unavailable, IPC access
  denied, timeout, other transport failures, and API response errors as distinct
  categories while retaining the original cause.
- **FR-002**: Access denial MUST NOT include or imply the daemon-absence message
  “is the daemon running?”.
- **FR-003**: One underlying connectivity incident MUST produce at most one
  visible application state across model refresh, calendar refresh, event
  streaming, and later refreshers.
- **FR-004**: Startup transport failures MUST be presented inside the reachable
  application frame and MUST NOT use an OK-only blocking modal.
- **FR-005**: Every active connection incident MUST offer Retry and Exit.
- **FR-006**: On Windows access denial, diagnostics MUST distinguish verified
  stale-token membership, verified absent group membership, and an unverified
  cause without weakening or bypassing authorization.
- **FR-007**: Verified stale-token membership guidance MUST tell the user to
  sign out and back in before retrying.
- **FR-008**: Retry MUST perform one coordinated refresh/reconnect attempt,
  update the existing incident on failure, and clear it on successful recovery.
- **FR-009**: Background reconnection MUST use bounded exponential backoff,
  MUST update the existing incident rather than create dialogs, and MUST have a
  defined cancellation path when the application exits.
- **FR-010**: An active connectivity incident MUST suppress duplicate visible
  errors from tab refreshers without hiding unrelated non-transport operation
  errors.
- **FR-011**: Tests MUST simulate concurrent tasks, calendar, and event-stream
  access denial and prove one deduplicated incident.
- **FR-012**: Tests MUST distinguish daemon absence, access denial, timeout,
  other transport failure, and API response error classifications.
- **FR-013**: A native Windows diagnostic run MUST record service status,
  `goschedadmin` existence, installing-account membership, process-token group
  SID membership, live restricted-pipe denial, and selection of stale-session
  guidance when the account is enrolled but the token lacks the SID.
- **FR-014**: The implementation MUST preserve the restricted named-pipe
  authorization boundary and MUST NOT grant access to additional identities.
- **FR-015**: The slice MUST pass all eight canonical repository verification
  gates without weakening IPC security tests.

### Key Entities

- **Connection failure**: A classified error with category, actionable message,
  retained cause, and transport/API distinction.
- **Connection incident**: The single visible state representing current
  connectivity, including category, guidance, retry status, and lifecycle.
- **Windows access diagnosis**: Evidence about service, local group, account
  membership, and current process-token SID membership, where each observation
  can also be unknown.
- **Reconnect policy**: An interruptible bounded sequence of attempt delays that
  resets after successful connectivity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Three simultaneous startup access-denied failures produce exactly
  one visible connection incident and zero error dialogs.
- **SC-002**: Acknowledging, retrying, or waiting with an active incident creates
  zero additional visible incidents for the same connectivity condition.
- **SC-003**: Each required failure class is identified correctly in automated
  tests, with access denial never described as daemon absence.
- **SC-004**: Every connection incident keeps the application frame, Retry, and
  Exit reachable throughout recovery.
- **SC-005**: Background retry delays increase to a finite maximum, are
  cancellation-aware, and never create more than one concurrent reconnect loop.
- **SC-006**: A deterministic authorization transition proves one Retry restores
  the normal interface without reinstalling.
- **SC-007**: Combined native Windows diagnostics and deterministic recovery
  evidence prove that real token state selects the guidance and successful
  authorization restores usable tabs without forcing an active session to end.
- **SC-008**: All eight canonical gates pass and the IPC authorization policy is
  unchanged.

## Clarifications

### Session 2026-08-31

- Q: Where is the incident presented? A: As a persistent in-frame recovery
  panel above the normal tab content, never as a blocking modal.
- Q: How is stale membership asserted? A: Only when account membership is
  verified and the current process token lacks the group SID; otherwise the
  guidance states what is known and how to inspect it.
- Q: What retry policy applies? A: One reconnect loop uses exponential backoff
  from 2 seconds up to 30 seconds and resets after success.

## Assumptions

- Windows local-group and token inspection can fail under policy restrictions;
  unknown evidence is a first-class diagnostic result.
- Signing out and back in remains the user-facing remedy for a verified stale
  token, but terminating an active desktop session is not a development
  verification prerequisite; recovery is proven by a deterministic
  access-denied-to-authorized transition.
- API response errors are operation errors and do not activate the global
  transport incident unless their underlying cause is transport-related.

## Out of Scope

- Changing named-pipe ACLs, `goschedadmin`, or the daemon authorization model.
- Automatically modifying group membership, elevating the GUI, restarting the
  service, signing the user out, or bypassing access checks.
- Installer authoring changes beyond documentation or verification evidence
  needed for this recovery path.
- Publishing a release, merging the pull request, or tagging a version.
