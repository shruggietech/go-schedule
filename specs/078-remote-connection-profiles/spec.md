# Feature Specification: Remote Connection Profiles and Target-Safe Clients

**Feature Branch**: `codex/078-remote-connection-profiles`

**Created**: 2026-09-09

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Shared profile, trust, selection, and transport contracts for the desktop, CLI, and documented JSON clients completed for [#170](https://github.com/shruggietech/go-schedule/issues/170) and [#171](https://github.com/shruggietech/go-schedule/issues/171), with canonical eight-gate verification passed 2026-09-09 on review branch `codex/078-remote-connection-profiles`.

**Input**: GitHub issues #170 and #171 on the authenticated HTTPS and native credential-storage foundation delivered by S077.

## Clarifications

### Session 2026-09-09

- The desktop persists one active connection selection, while the CLI remains local by default and requires `--profile` or a complete explicit endpoint selection for every remote invocation.
- A profile pins the daemon installation ID, HTTPS endpoint, trusted certificate, credential ID, safe display metadata, and last successful contact. The bearer value remains only in the operating-system credential store.
- Pairing creates a profile only after identity validation, credential storage, and profile persistence all succeed. Repair replaces the credential and trust material for the same pinned daemon identity.
- Desktop-local Settings remain available for every target. Daemon feature routes are enabled only when the active target advertises the capability and the credential grants the required authority.
- Automatic reconnect, resume cursors, certificate-change recovery, and stale-data retention beyond an explicit failed attempt remain assigned to #172.

## User Scenarios & Testing

### User Story 1 - Pair and select a persistent desktop target (Priority: P1)

As a desktop user, I want a paired daemon to become a durable named connection that I can select without re-entering trust or credential material.

**Why this priority**: Persistent target selection turns S077 enrollment into a usable remote workflow.

**Independent Test**: Pair two same-named daemons, restart the desktop, select each by its stable identity, and confirm the global target context and requests follow the selected profile.

**Acceptance Scenarios**:

1. **Given** a successful pairing response, **When** protected storage and profile persistence complete, **Then** one profile is listed and may become the desktop-active target.
2. **Given** multiple profiles with the same display name, **When** the connection list is shown, **Then** endpoint and shortened stable daemon identity distinguish them.
3. **Given** a persisted remote selection, **When** the desktop restarts, **Then** it attempts only that exact profile and identifies it before any feature request.
4. **Given** This computer is selected, **When** the desktop operates normally, **Then** local IPC remains credential-free and unchanged.

### User Story 2 - Manage, repair, and remove desktop profiles safely (Priority: P1)

As a desktop user, I want to rename, repair, and remove connections while seeing honest trust, compatibility, authority, and recent-contact state.

**Why this priority**: Durable profiles need a complete lifecycle and must not strand native credentials.

**Independent Test**: Rename, reject an identity-changing repair, repair a revoked credential, remove active and inactive profiles, and verify both metadata and native credential deletion.

**Acceptance Scenarios**:

1. **Given** a profile, **When** its local label is renamed, **Then** only the local label changes and the pinned daemon identity does not.
2. **Given** a replacement phrase for the same daemon, **When** repair succeeds, **Then** the new credential is stored, the old native credential is deleted, and the profile remains the same identity.
3. **Given** a repair response for another daemon identity, **When** validation occurs, **Then** repair fails without changing the profile or native credential.
4. **Given** an inactive profile, **When** removal is confirmed, **Then** its native credential is deleted before its metadata disappears.
5. **Given** the active remote profile, **When** removal succeeds, **Then** the desktop switches to This computer and does not send another request to the removed target.

### User Story 3 - Run CLI commands against an explicit remote profile (Priority: P1)

As a CLI user, I want named remote profiles and one-command endpoint selection without credentials in arguments or an unexpected remote default.

**Why this priority**: Headless administration and scripting must gain remote access without weakening existing local behavior.

**Independent Test**: Create a profile through protected phrase input, invoke supported commands with `--profile`, use explicit endpoint selection with a native credential reference, and prove commands without selection still use local IPC.

**Acceptance Scenarios**:

1. **Given** any saved desktop or CLI profile, **When** a command names it with `--profile`, **Then** the command pins the stored daemon identity and reports safe target context in human output.
2. **Given** no target flags, **When** any existing command runs, **Then** it uses This computer exactly as before.
3. **Given** a complete explicit endpoint selection, **When** a command runs, **Then** the bearer value is loaded from native storage and never appears in process arguments.
4. **Given** incomplete, conflicting, missing, revoked, or identity-mismatched selection, **When** a command starts, **Then** it fails before the requested mutation and returns a stable nonzero exit.
5. **Given** JSON output, **When** a remote command runs, **Then** stdout remains valid command JSON and target context does not corrupt it.

### User Story 4 - Integrate through documented JSON safely (Priority: P2)

As a client developer, I want copyable examples that enroll, pin daemon identity, authenticate, classify failures, and avoid replaying uncertain mutations.

**Why this priority**: The network contract is incomplete as a product surface without safe operational guidance.

**Independent Test**: Follow the private HTTPS and SSH tunnel examples from a clean environment with a disposable phrase, then exercise wrong identity, unauthorized, timeout, and uncertain mutation cases.

**Acceptance Scenarios**:

1. **Given** a phrase and trusted certificate, **When** the documented enrollment exchange succeeds, **Then** the client verifies the returned daemon ID before retaining the credential.
2. **Given** a stored credential, **When** the documented request example runs, **Then** the credential is read from protected input rather than embedded in a URL, source file, or command history.
3. **Given** a lost mutation response, **When** the guidance is followed, **Then** the client refreshes authoritative state before any deliberate retry.

### Edge Cases

- Profile IDs and local labels are distinct; labels may repeat, but exact profile references and daemon installation IDs do not become ambiguous.
- Endpoints are canonical HTTPS origins with no user information, query, fragment, or path beyond an optional trailing slash.
- Certificate PEM is trust material rather than a bearer secret, but profile diagnostics expose only its SHA-256 fingerprint.
- A profile file never contains a bearer token, pairing phrase, environment secret, or task execution input.
- Profile writes are atomic, permission-restricted where supported, forward-version rejected, and recoverable from a missing file as an empty collection with This computer active.
- Credential deletion failure leaves profile metadata intact and reports actionable recovery instead of claiming removal.
- An unavailable keyring, profile file, certificate, endpoint, or daemon produces no fallback to plaintext storage or insecure TLS.
- The selected target is visible globally and in mutation confirmations; same-named targets include endpoint and short daemon identity.
- Unsupported remote feature areas remain visible only with a useful disabled explanation and never fall back to local IPC silently.
- Desktop selection cancels the prior connection generation before the new target can publish events.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST maintain one versioned user-scoped profile collection shared by the desktop and CLI, with atomic writes and a default empty state.
- **FR-002**: Every remote profile MUST contain a stable profile ID, local label, canonical HTTPS endpoint, pinned daemon installation ID, credential ID, trusted certificate, certificate fingerprint, client kind, granted capability, safe platform and version metadata, creation and update timestamps, and optional last-success timestamp.
- **FR-003**: Profile files MUST exclude bearer credentials, pairing phrases, certificate private keys, task inputs, notification secrets, and other authority-bearing values.
- **FR-004**: The desktop MUST persist one active profile reference or This computer; a missing or invalid active reference MUST fail safely to This computer with a visible explanation.
- **FR-005**: The CLI MUST use This computer unless each remote invocation deliberately supplies one named `--profile` or the complete explicit endpoint, daemon ID, credential ID, and certificate-file selection.
- **FR-006**: Named profile and explicit endpoint selection MUST be mutually exclusive, and incomplete selection MUST fail before contacting a daemon.
- **FR-007**: Remote clients MUST use TLS 1.3 or newer, trust only the selected profile certificate roots, reject redirects, attach one native-store bearer credential, and pin the returned manifest installation ID before feature operations.
- **FR-008**: Pairing from desktop or CLI MUST store the credential natively and persist profile metadata only after validating the expected daemon identity; partial failure MUST leave no usable half-profile.
- **FR-009**: CLI phrase and bearer input MUST never be accepted as ordinary command arguments; enrollment MUST read the phrase from protected standard input and normal commands MUST load bearer values by credential reference.
- **FR-010**: The desktop Connections view MUST list This computer plus every remote profile with label, endpoint, shortened daemon ID, trust fingerprint, capability, platform, version, connection state, and last successful contact.
- **FR-011**: The desktop MUST support selecting, locally renaming, repairing, and removing profiles without configuration-file editing.
- **FR-012**: Repair MUST require a fresh phrase and trusted certificate, MUST preserve the pinned daemon identity, MUST replace the stored credential only after successful enrollment, and MUST delete the superseded native credential.
- **FR-013**: Removal MUST delete the native credential before deleting profile metadata; removing the active remote profile MUST select This computer first.
- **FR-014**: Every desktop screen MUST display the active target globally, and every desktop mutating action MUST include the target label plus disambiguating endpoint or short daemon ID before submission.
- **FR-015**: Desktop feature areas and actions MUST be enabled only when the active daemon advertises their capability and the current credential grants the required authority; unsupported areas MUST provide an explanation and MUST NOT route to This computer.
- **FR-016**: A target switch MUST cancel the previous health and event generation, reject stale results, configure all daemon-backed desktop services for the new target, and publish the selected target before new requests begin.
- **FR-017**: CLI human output for remote invocations MUST identify the selected profile, endpoint, and short daemon ID on stderr so stdout result contracts remain unchanged; JSON stdout MUST remain machine-readable without diagnostic prefixes.
- **FR-018**: Profile list, show, rename, remove, and pair commands MUST support human and JSON output, stable exit codes, and actionable failure messages.
- **FR-019**: Documented JSON examples MUST cover enrollment, protected credential handling, manifest identity pinning, Bearer authentication, pagination, stable errors, timeouts, SSH tunnels, private HTTPS, direct HTTPS prerequisites, revocation, and non-replay of uncertain mutations.
- **FR-020**: Logs, errors, UI state, profile output, JSON output, and diagnostics MUST never disclose bearer values, pairing phrases, certificate private keys, or task execution secrets.
- **FR-021**: Profile selection and administration MUST be race-safe across desktop and CLI processes; concurrent writes MUST either serialize or fail without corrupting the last complete document.
- **FR-022**: This slice MUST NOT add automatic retry loops, offline mutation queues, service discovery, certificate rotation recovery, browser credential storage, insecure TLS modes, remote MCP mutation, or a CLI-wide default remote target.

### Key Entities

- **Connection Profile**: User-scoped metadata that pins one client relationship to one HTTPS daemon identity without storing its bearer value.
- **Profile Collection**: Versioned atomic document containing profiles and the desktop-active reference.
- **Target Selection**: One invocation-scoped CLI selection or persistent desktop selection resolving to local IPC or one exact profile.
- **Remote Client Configuration**: In-memory endpoint, trust roots, bearer value, expected daemon identity, and request path mapping.
- **Profile Lifecycle Result**: Secret-free action result for pair, select, rename, repair, or remove.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A paired desktop profile remains selectable after restart and reaches the same daemon identity without re-entering a phrase, certificate, or credential.
- **SC-002**: Tests with two identical display names always distinguish profiles by endpoint and stable identity and send 100 percent of requests to the explicitly selected daemon.
- **SC-003**: Existing CLI commands with no remote flags retain byte-compatible local stdout behavior and use no TCP connection.
- **SC-004**: One named CLI profile or complete explicit selection can run every operation supported by its remote capability while invalid selections perform zero requested mutations.
- **SC-005**: Secret-canary scans across profiles, output, diagnostics, logs, and UI state find zero bearer values or phrases.
- **SC-006**: Pair, repair, removal, corrupt-file, forward-version, concurrent-write, keyring-failure, identity-mismatch, wrong-target, and unsupported-capability tests pass without metadata corruption or insecure fallback.
- **SC-007**: Every desktop mutation presents the selected target before submission, and keyboard plus assistive-technology tests can identify and cancel the action.
- **SC-008**: JSON client guidance completes enrollment and one authenticated read through private HTTPS and an SSH tunnel while preserving identity verification and stable error handling.
- **SC-009**: The full eight-gate verification suite passes with core coverage at or above 80 percent and no race finding.

## Assumptions

- The S077 HTTPS, enrollment, keyring, authorization, audit, and remote operation contracts are the stable foundation for this slice.
- Profile metadata is user-scoped because desktop and CLI credentials belong to the interactive user, while daemon configuration remains machine-scoped.
- The remote operation allowlist determines which existing workflows are available. This slice does not broaden that allowlist merely to make every local screen active.
- #172 will add automatic recovery after ordinary connection loss. S078 performs a bounded initial connection and explicit retry or reselection only.

## Dependencies

- Parent: #18.
- Completes: #170 and #171.
- Depends on: #157, #166, #167, #168, and #169, all complete.
- Blocks: #172 and contributes to the release gate in #173.

## Scope Boundaries

**In scope**: Shared profile persistence, desktop pairing persistence and lifecycle, desktop target selection and context, CLI profile administration and explicit remote targeting, remote transport configuration, permission and capability gating, direct JSON guidance, tests, and documentation.

**Out of scope**: Automatic reconnect and resume behavior from #172, v1.4 release qualification from #173, new daemon routes, remote notification or automation management, remote MCP, discovery, NAT traversal, hosted services, and release publication.
