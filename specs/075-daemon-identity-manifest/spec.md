# Feature Specification: Stable Daemon Identity and Capability Manifest

**Feature Branch**: `codex/075-daemon-identity-manifest`

**Created**: 2026-09-09

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/075-daemon-identity-manifest`; persistent identity, protected local manifest and lifecycle operations, CLI and desktop integration, documentation, focused tests, and canonical eight-gate verification completed 2026-09-09 for #166

**Input**: User description: "Complete issue #166 end-to-end by giving every daemon a persistent installation identity, editable display name, safe capability manifest, explicit upgrade, reinstall, restore, clone, and reset behavior, local API compatibility, documentation, and tests. Exclude actor authorization, remote HTTPS, pairing, remote profiles, and issue #167 onward behavior."

## Clarifications

### Session 2026-09-09

- Q: Does restoring the daemon database preserve or replace the installation identity? → A: Preserve it, because a restore continues the same logical daemon unless an operator deliberately resets identity.
- Q: How must a copied or cloned installation avoid duplicate target identity? → A: The clone remains the same logical daemon until an operator performs an explicit compare-and-confirm identity reset before concurrent use.
- Q: May the default display name disclose the machine hostname? → A: No. Use a generic product-owned default and disclose no hostname.
- Q: What compatibility boundary applies while adding discovery? → A: Preserve the existing health contract and add discovery without removing or changing existing local operations.

## User Scenarios & Testing

### User Story 1 - Identify the Connected Daemon (Priority: P1)

As a local client user, I want every daemon to report an immutable installation identifier separately from its editable name so that same-named daemons cannot be mistaken for one another.

**Why this priority**: Stable target identity is the prerequisite for every future remote profile, audit record, and wrong-target safeguard.

**Independent Test**: Start a fresh daemon, read its manifest repeatedly and after restart, then compare identity, display name, and safe platform facts without relying on an address or hostname.

**Acceptance Scenarios**:

1. **Given** a fresh data store, **When** the daemon first starts, **Then** it receives one opaque installation identifier and a generic display name that remain stable across ordinary restarts.
2. **Given** two independent fresh data stores with the same display name, **When** their manifests are compared, **Then** their installation identifiers differ.
3. **Given** an existing local client that uses the health or other local endpoints, **When** it connects to the upgraded daemon, **Then** its existing contract remains compatible.

---

### User Story 2 - Discover Compatibility Before Actions (Priority: P1)

As a client author, I want one bounded manifest that reports product version, protocol versions, operating mode, features, and safe platform facts so that incompatible or unavailable actions are not presented.

**Why this priority**: Clients must negotiate from daemon-owned facts instead of duplicating capability lists or guessing from a display name.

**Independent Test**: Request the manifest through the protected local transport and verify that its stable, ordered values account for the daemon's current features while remote protocol support remains explicitly absent.

**Acceptance Scenarios**:

1. **Given** the current local-only daemon, **When** a client reads the manifest, **Then** it receives the installation identifier, display name, product version, supported local API versions, empty remote API versions, local-only operating mode, safe platform facts, and a deterministic feature list.
2. **Given** a client evaluates a feature, **When** that feature is absent from the manifest, **Then** the client can withhold the corresponding action without probing a privileged endpoint.
3. **Given** a local desktop connection succeeds, **When** it publishes target context, **Then** the target uses the daemon-owned identifier, display name, version, platform, and capabilities rather than client-owned placeholders.

---

### User Story 3 - Rename or Deliberately Reset Identity (Priority: P2)

As an operator with local daemon access, I want to rename the daemon and deliberately reset a cloned identity so that target labels remain useful without making destructive identity changes accidental.

**Why this priority**: Names are operational aids, while identity reset is a rare lifecycle tool that must require exact acknowledgement.

**Independent Test**: Rename a daemon, restart it, reject invalid names and mismatched reset confirmations, then reset with the exact current identifier and verify that only the identifier changes.

**Acceptance Scenarios**:

1. **Given** an operator supplies a valid display name, **When** the name is updated, **Then** the trimmed name persists and the installation identifier does not change.
2. **Given** an empty, control-character-bearing, or oversized name, **When** an update is attempted, **Then** it fails predictably without changing stored identity.
3. **Given** a reset request that does not exactly confirm the current identifier, **When** the daemon handles it, **Then** the request fails as a conflict and nothing changes.
4. **Given** a reset request that exactly confirms the current identifier, **When** the daemon handles it, **Then** a new identifier is persisted atomically while the display name and scheduler data remain unchanged.

### Edge Cases

- Two processes attempt first-use initialization. Exactly one singleton identity is retained and every reader receives it.
- A migration or first-use identity write fails. Daemon startup fails with contextual diagnostics rather than running with an empty or volatile identity.
- A database backup is restored over the same installation. The restored logical daemon keeps the backed-up identity and display name.
- A database is copied to create an independently operated clone. Documentation requires a deliberate identity reset before both copies operate concurrently.
- An identity reset is repeated using the previous identifier. The second request fails because its confirmation is stale.
- An unrecognized future manifest field or feature is received by an older local client. Existing fields and endpoints continue to decode and behave normally.
- Product version metadata is a development string. The manifest reports it verbatim without using it as installation identity.

## Requirements

### Functional Requirements

- **FR-001**: Every initialized daemon store MUST contain exactly one opaque installation identifier that is independent of display name, hostname, address, product version, and scheduler-object identities.
- **FR-002**: Independently initialized daemon stores MUST receive different installation identifiers, while ordinary restart and upgrade MUST preserve an existing identifier.
- **FR-003**: The daemon MUST retain a nonempty editable display name separately from installation identity and MUST default to a generic name that reveals no hostname or account information.
- **FR-004**: Display names MUST be trimmed, contain 1 through 80 Unicode characters, and reject control characters without mutating stored state.
- **FR-005**: The daemon MUST expose one bounded manifest through the protected local API containing installation identity, display name, product version, supported local API versions, supported remote API versions, operating mode, deterministic feature capabilities, and safe platform facts.
- **FR-006**: The S075 manifest MUST report local API `v1`, no remote API version, and `local_only` operating mode without opening a network listener or implying remote support is shipped.
- **FR-007**: Manifest collections MUST be deterministic, duplicate-free, and safe for direct client display and compatibility decisions.
- **FR-008**: The manifest MUST exclude hostnames, network addresses, storage paths, account names, credentials, trigger keys, environment values, commands, and other secrets or unnecessarily identifying data.
- **FR-009**: The desktop local connection MUST consume daemon-owned manifest identity and capabilities while preserving the existing protected local transport and local authorization experience.
- **FR-010**: Existing health, runtime-information, scheduling, activity, notification, agent-access, and other local API contracts MUST remain available and backward compatible.
- **FR-011**: An authorized local operator MUST be able to read the manifest and update the display name through both the shared local API client and CLI with human-readable and JSON output.
- **FR-012**: Identity reset MUST require the caller to provide the exact current installation identifier, MUST atomically create a different identifier, MUST preserve display name and scheduler data, and MUST reject stale or mismatched confirmation as a conflict.
- **FR-013**: A restored database MUST preserve its stored identity; an independently operated clone MUST be reset deliberately before concurrent use; a clean installation with no retained data MUST create a new identity.
- **FR-014**: Migration from the prior schema MUST be forward-only, preserve every existing scheduler record, initialize one valid identity, and remain covered by prior-schema verification.
- **FR-015**: Identity initialization, reads, renames, and resets MUST return contextual errors and MUST never substitute an in-memory identity after persistence failure.
- **FR-016**: S075 MUST NOT add actors, permissions enforcement, audit history, bearer credentials, pairing phrases, remote profiles, remote HTTPS, certificate handling, or a remote listener.

### Key Entities

- **Daemon Identity**: The singleton persistent installation identifier, editable display name, and lifecycle timestamps for one logical daemon.
- **Capability Manifest**: A safe immutable response projection containing identity, compatibility, operating mode, feature, and platform facts.
- **Identity Reset Request**: A deliberate compare-and-confirm operation carrying the exact current identifier and producing a replacement identifier without changing scheduler state.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred independently initialized test stores produce 100 distinct nonempty installation identifiers, and each identifier remains unchanged across repeated reads and restart.
- **SC-002**: The prior-schema migration preserves 100 percent of seeded scheduler records and creates exactly one valid daemon-identity record.
- **SC-003**: Every manifest response contains all required fields, zero prohibited host or secret fields, deterministic duplicate-free collections, and an empty remote API version list.
- **SC-004**: Valid renames survive restart; 100 percent of invalid-name and stale-confirmation cases leave identity state unchanged.
- **SC-005**: A successful reset changes exactly one installation identifier while preserving the display name and all seeded scheduler records.
- **SC-006**: Existing local health and client compatibility tests continue to pass without requiring callers to adopt the manifest.
- **SC-007**: The complete repository verification aggregate passes with no weakened scheduling, migration, concurrency, desktop, IPC, documentation, or automation gate.

## Assumptions

- The daemon database represents the logical daemon for backup and restore purposes, so restoring it intentionally restores daemon identity.
- Copying the complete database creates a logical clone; operators who intend simultaneous independent operation must reset one copy before starting both.
- The existing protected local IPC boundary is sufficient authority for rename and reset in S075. The shared actor and capability enforcement planned by #167 remains out of scope.
- Generic platform operating-system and architecture values are safe compatibility facts; hostname, account, address, and path values are not.
- Feature capability names describe supported product surfaces, not the authorization granted to a caller.

## Dependencies

- Parent: #18.
- Depends on completed #165 and the S074 architecture contract.
- Blocks #167, #168, #169, #170, #171, #172, and #173 where stable target identity or capability negotiation is required.
- Implements and closes #166 only.
