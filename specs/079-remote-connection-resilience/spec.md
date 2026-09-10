# Feature Specification: Resilient Remote Connections

**Feature Branch**: `codex/079-remote-connection-resilience`

**Created**: 2026-09-10

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Automatic remote recovery, explicit stale and terminal states, persisted last contact, and single-attempt mutation uncertainty completed for [#172](https://github.com/shruggietech/go-schedule/issues/172), with canonical eight-gate verification passed 2026-09-10 on review branch `codex/079-remote-connection-resilience`.

**Input**: GitHub issue [#172](https://github.com/shruggietech/go-schedule/issues/172), building on the remote profiles and target-safe clients delivered by S078.

## Clarifications

### Session 2026-09-10

- Safe reads and live subscriptions may retry automatically with bounded exponential backoff and jitter; state-changing requests are attempted once and are never replayed automatically.
- A successful connection refreshes the authoritative snapshot before live events resume. The event stream is an invalidation signal rather than a durable event log, so recovery does not require server-side event replay or an offline queue.
- Certificate or daemon identity changes are terminal for the existing profile. Recovery requires an explicit repair flow that preserves the profile only when the pinned daemon identity is unchanged; a different identity requires a separately confirmed profile.
- Remote stale data remains visible, clearly marked with the time of last successful contact, while controls that could contact the daemon are disabled until reconnection succeeds.
- Automatic retries stop for authentication, revocation, authorization, incompatibility, and identity or trust failures, but continue for ordinary transport loss, timeouts, suspend and resume, address recovery, and daemon restart.

## User Scenarios & Testing

### User Story 1 - Recover observation after ordinary connection loss (Priority: P1)

As a remote desktop user, I want the selected scheduler to reconnect after temporary network loss, sleep, or daemon restart so that I can continue observing it without repeatedly pressing Try again.

**Why this priority**: Automatic recovery is the core value of this slice and makes remote profiles practical for normal laptop and network behavior.

**Independent Test**: Connect to a remote profile, interrupt the event stream and health endpoint, restore them, and confirm the same target recovers through bounded retries without user action or duplicated operations.

**Acceptance Scenarios**:

1. **Given** a connected remote profile, **When** its live stream ends because the network or daemon is temporarily unavailable, **Then** the desktop marks existing data stale, reports retry progress, and reconnects automatically within the documented retry bounds.
2. **Given** a remote profile selected before system sleep, **When** the client resumes and the daemon is reachable, **Then** the same profile is revalidated and fresh authoritative data replaces the stale snapshot.
3. **Given** repeated transient failures, **When** automatic recovery continues, **Then** retries use bounded delays, remain cancelable, and do not create overlapping connection generations.
4. **Given** a user switches targets during recovery, **When** a prior retry or response completes, **Then** it cannot change the selected target or publish stale events.

### User Story 2 - Understand stale and terminal states (Priority: P1)

As an operator, I want connection failures classified accurately so that I know whether to wait, retry, repair credentials, update software, or reject an unexpected identity.

**Why this priority**: Automatic retry is unsafe and confusing unless transient and terminal conditions are unmistakable.

**Independent Test**: Exercise unreachable, timed-out, unauthorized, revoked, forbidden, incompatible, certificate-changed, and daemon-identity-changed cases and verify distinct state, guidance, retry policy, and stale-data presentation.

**Acceptance Scenarios**:

1. **Given** a previously successful remote view, **When** a transient failure occurs, **Then** its last complete data remains visible with a stale marker and the last successful contact time.
2. **Given** rejected or revoked credentials, **When** the next request is made, **Then** automatic retry stops and the user is directed to repair the connection or obtain an updated grant.
3. **Given** a version incompatibility, **When** negotiation occurs, **Then** automatic retry stops and the user is directed to update the client or daemon.
4. **Given** a changed certificate or daemon identity, **When** reconnect validation occurs, **Then** no feature request proceeds, automatic retry stops, and the user receives identity-specific recovery guidance.
5. **Given** a remote target is stale or recovering, **When** assistive technology reads the shell, **Then** the target, stale state, last contact, retry status, and recovery action are announced without relying on color alone.

### User Story 3 - Keep mutations single-attempt and recover deliberately (Priority: P1)

As an operator changing remote scheduler state, I want uncertain mutation outcomes reported without automatic replay so that a connection interruption cannot duplicate or reverse my intent.

**Why this priority**: Preventing duplicate state changes is a safety boundary for every remote operation.

**Independent Test**: Drop the response after a remote mutation reaches the daemon and confirm the client reports an uncertain result, performs no automatic replay, preserves the editor draft, and refreshes authoritative state after connectivity returns.

**Acceptance Scenarios**:

1. **Given** a state-changing request whose response is lost, **When** its deadline expires, **Then** the client reports that the outcome is uncertain and never sends the request again automatically.
2. **Given** an uncertain mutation, **When** connectivity returns, **Then** the client refreshes authoritative state before offering another deliberate submission.
3. **Given** a user has unsaved editor input when connectivity is lost, **When** the desktop becomes stale or reconnects, **Then** the draft remains intact unless authoritative conflict handling requires the user to reload it.
4. **Given** controls are shown while a remote target is stale, recovering, or terminally failed, **When** the user navigates the desktop, **Then** daemon-backed mutations remain disabled while local settings and connection repair remain available.

### Edge Cases

- A transport failure can occur before a request is written, while it is being written, after the daemon commits it, or while the response is being read; all state-changing cases are treated as single-attempt and potentially uncertain unless an authoritative API response proves rejection.
- Retry timers and in-flight work are canceled on target switch, application shutdown, manual retry, and successful stream activity.
- Backoff reaches a bounded maximum and uses jitter without depending on wall-clock sleeps in deterministic tests.
- A successful health response does not make the target usable until its manifest identity, compatibility, and current credential authority are validated.
- Event recovery may miss invalidation messages, so each recovered feature workspace is refreshed from authoritative reads before subsequent live invalidations are accepted.
- A certificate validation failure is distinct from an unreachable endpoint and never enables insecure trust or silent certificate replacement.
- Address changes recover only when the selected profile endpoint becomes reachable again; this slice adds no discovery, DNS override, or automatic profile rewrite.
- Data from another profile or connection generation is never retained as the stale snapshot for the current target.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST automatically retry remote health negotiation and live subscriptions after transient transport loss, timeout, daemon restart, network transition, and client resume.
- **FR-002**: Automatic retry MUST use cancelable bounded backoff with jitter, beginning within one second and never waiting more than thirty seconds between attempts.
- **FR-003**: Only safe reads, identity checks, health checks, and live subscriptions MAY be retried automatically; state-changing requests MUST be attempted exactly once per deliberate user action.
- **FR-004**: A state-changing request that loses its response MUST return a distinct uncertain-outcome result that states the operation may have completed and directs the user to refresh authoritative state before deciding whether to try again.
- **FR-005**: Recovery MUST revalidate certificate trust, pinned daemon identity, compatible API version, capability manifest, and current credential authority before feature reads or subscriptions resume.
- **FR-006**: Authentication failure, credential revocation, insufficient authority, incompatible version, certificate failure, and daemon identity mismatch MUST stop automatic retry and expose distinct actionable recovery guidance.
- **FR-007**: The connection snapshot MUST expose whether displayed daemon data is stale, when contact last succeeded, the next automatic retry time when scheduled, the current retry attempt, and whether recovery is automatic or requires user action.
- **FR-008**: The desktop MUST retain the last complete snapshot for the selected remote target during recovery or terminal failure and MUST mark it stale everywhere it remains visible.
- **FR-009**: Daemon-backed mutation controls MUST be unavailable whenever the selected remote target is not fully connected and current; desktop-local settings, target selection, credential repair, and profile removal MUST remain available.
- **FR-010**: After reconnection, every mounted remote feature workspace MUST refresh from an authoritative read before accepting new live invalidation events.
- **FR-011**: Connection generations MUST remain single-owner, race-safe, and cancelable; results, retry timers, and events from an older generation MUST never affect the current target.
- **FR-012**: Manual retry MUST remain available for transient and manually recoverable states, MUST cancel the current wait or attempt, and MUST reset the retry delay without creating an overlapping generation.
- **FR-013**: Successful remote contact MUST update the profile's last-success timestamp without persisting credentials, phrases, private keys, task inputs, or error internals.
- **FR-014**: Recovery status and stale state MUST be conveyed through text and assistive-technology announcements without relying on color, animation, or pointer interaction.
- **FR-015**: Existing local IPC behavior MUST remain unchanged and MUST NOT inherit remote backoff, stale-data, credential, certificate, or uncertain-mutation semantics.
- **FR-016**: Logs and user-visible diagnostics MUST use safe target context and failure categories without disclosing bearer values, pairing phrases, task inputs, certificate private keys, or raw transport errors containing secrets.
- **FR-017**: This slice MUST NOT add an offline mutation queue, conflict-resolution subsystem, automatic certificate acceptance, service discovery, endpoint rewriting, durable server-side event replay, or release qualification.

### Key Entities

- **Connection Snapshot**: Generation-stamped, secret-free state containing target identity, connectivity classification, freshness, last contact, retry timing, attempt count, and recovery guidance.
- **Retry Policy**: Classification and bounded schedule governing which connection operations recover automatically and which require intervention.
- **Authoritative Workspace Snapshot**: The last complete target-scoped read retained while stale and replaced after validated reconnection.
- **Mutation Outcome**: Accepted, rejected, unavailable, stale-conflict, or uncertain result for one deliberate state-changing request.
- **Connection Generation**: Single-owner lifecycle for one immutable selected target, including negotiation, refresh, event subscription, retry waits, and cancellation.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A remote target recovering from ordinary connection loss begins its first automatic retry within one second and never waits more than thirty seconds between attempts.
- **SC-002**: Sleep, resume, network loss, daemon restart, and restored-endpoint scenarios recover the same selected profile without manual reselection in 100 percent of deterministic lifecycle tests.
- **SC-003**: Lost-response mutation scenarios produce exactly one outbound mutation attempt and a clear uncertain-outcome result in 100 percent of tests.
- **SC-004**: Every retained remote workspace is visibly and programmatically marked stale during disconnection and is replaced by an authoritative refresh before live updates resume.
- **SC-005**: Unreachable, timed-out, unauthorized, revoked, forbidden, incompatible, certificate-changed, and identity-changed scenarios produce the documented distinct state and recovery guidance in 100 percent of classification tests.
- **SC-006**: One hundred repeated connect, disconnect, recover, switch, and shutdown cycles complete with no leaked owner goroutine, overlapping generation, stale event, or race finding.
- **SC-007**: Secret-canary checks across recovery logs, errors, snapshots, UI announcements, and persisted profiles find zero credentials, phrases, private keys, or task inputs.
- **SC-008**: Local IPC compatibility tests remain unchanged in behavior, and the full eight-gate verification suite passes with every core package at or above 80 percent coverage.

## Assumptions

- S078's immutable target clients, persisted profiles, identity pinning, and native credential storage are the stable starting point.
- Existing feature workspaces already retain their last complete successful data when a refresh fails; S079 adds an explicit freshness contract rather than an offline data store.
- The event stream is used to trigger authoritative refreshes. Exact event replay is unnecessary because daemon state can be reconstructed through bounded reads.
- The operating system and network stack report resume and address recovery through ordinary request success or failure; no platform-specific power-event dependency is required.
- A current successful protected read is sufficient evidence that the stored credential remains active for its granted authority.

## Dependencies

- Parent: [#18](https://github.com/shruggietech/go-schedule/issues/18).
- Completes: [#172](https://github.com/shruggietech/go-schedule/issues/172).
- Depends on: #168 and #170, both complete through S077 and S078.
- Blocks: the v1.4.0 documentation and qualification gate in #173.

## Scope Boundaries

**In scope**: Remote transient-failure classification, automatic connection and subscription recovery, bounded retry timing, stale-data presentation, last-contact persistence, single-attempt mutation safety, authoritative refresh after recovery, target-generation cancellation, accessibility, tests, and behavior documentation.

**Out of scope**: Offline mutation queues, automatic mutation replay, durable event logs, automatic certificate trust changes, daemon discovery, profile endpoint rewriting, clustered execution, new remote feature routes, and the cross-platform v1.4.0 release qualification owned by #173.
