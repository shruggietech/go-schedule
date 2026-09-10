# Research: Remote Connection Profiles and Target-Safe Clients

## Decision 1: Share metadata, not target defaults

**Decision**: Desktop and CLI read the same user-scoped profile collection. Only the desktop persists an active profile; CLI commands remain local unless the invocation explicitly supplies `--profile` or all explicit endpoint fields.

**Rationale**: Shared metadata avoids duplicate credentials and configuration, while invocation-scoped CLI selection prevents unattended scripts from silently changing machines when a desktop preference changes.

**Alternatives considered**: A shared global default was rejected as a wrong-target risk. Separate desktop and CLI profile stores were rejected because they duplicate trust and lifecycle state.

## Decision 2: Use immutable clients behind a synchronized provider

**Decision**: Construct one immutable local or remote API client per target. Desktop services resolve the provider's current client for each operation. A switch cancels the old connection generation, installs the new client, and then begins negotiation.

**Rationale**: An in-flight request cannot have its endpoint or credential changed underneath it. All feature services still share the existing API method implementation.

**Alternatives considered**: Mutating one client's transport was rejected because concurrent operations could cross targets. Reimplementing every service against generated remote types was rejected as duplicated semantics.

## Decision 3: Store a versioned atomic JSON document

**Decision**: Store profile metadata in `profiles.json` under the existing desktop user configuration root. Serialize writers with an adjacent lock, write a permission-restricted temporary file, flush, atomically rename, and reject unknown forward versions.

**Rationale**: Profiles belong to the interactive user and contain no bearer secrets. A single document makes desktop selection and cross-process updates transactional.

**Alternatives considered**: Daemon SQLite was rejected because profiles are client-owned and would be unavailable before target selection. One file per profile was rejected because active-selection updates would not be atomic with lifecycle changes.

## Decision 4: Pin identity and trust independently of labels

**Decision**: Treat the daemon installation ID and canonical HTTPS origin as identity. Treat the local profile label as editable presentation. Persist certificate PEM for trust configuration and expose only its SHA-256 fingerprint in diagnostics.

**Rationale**: Display names are intentionally mutable and non-unique. Identity and trust must survive renaming and distinguish clones.

**Alternatives considered**: Label-keyed profiles and trust-on-first-use were rejected as ambiguous and unsafe. System trust alone remains usable only by importing the operator certificate into the explicit profile trust input.

## Decision 5: Keep credential lifecycle fail-closed

**Decision**: Pairing first validates the daemon response, stores the new bearer natively, and then commits profile metadata. Repair commits new metadata and credential before deleting the old native credential. Removal deletes the credential before metadata.

**Rationale**: Metadata must never claim a usable relationship when no credential exists, and failed credential deletion must remain visible for cleanup.

**Alternatives considered**: Metadata-first writes were rejected because crashes create broken profiles. Deleting old credentials before a repair commit was rejected because a failed repair would destroy the working relationship.

## Decision 6: Preserve CLI stdout contracts

**Decision**: Remote target context is printed once to stderr for human-mode commands. JSON stdout remains the existing command payload. Profile lifecycle commands expose target metadata directly in their own JSON contracts.

**Rationale**: Scripts already parse stdout. The explicit `--profile` or endpoint flags identify invocation intent, while stderr gives humans unmistakable context without corrupting machine output.

**Alternatives considered**: Wrapping every JSON response was rejected as a broad breaking change. Printing a human banner to stdout was rejected because it changes command result contracts.

## Decision 7: Defer automatic resilience

**Decision**: S078 supports bounded initial connection, explicit retry, reselection, and repair. It does not add resume cursors, sleep recovery, automatic mutation retry, or certificate-change workflows.

**Rationale**: Those behaviors have independent failure semantics and acceptance matrices in #172. Combining them would weaken reviewability of the profile and wrong-target boundary.

**Alternatives considered**: A full reconnect state machine was rejected as scope expansion. Silent fallback to local IPC was rejected as a severe wrong-target defect.
