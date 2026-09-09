# Research: Stable Daemon Identity and Capability Manifest

## Persistent Identity

**Decision**: Store a random UUID in exactly one SQLite singleton row and initialize it during `store.Open` after schema migration.

**Rationale**: The database is already the daemon's durable logical state. A database-owned random identifier survives restart, upgrade, and restore without disclosing machine attributes. The existing UUID dependency can generate it without adding supply-chain surface.

**Alternatives considered**: Hostname, machine ID, address, or data-path hashes leak identifying data and can change. A configuration-file identity can drift away from restored scheduler state. An in-memory fallback would make persistence failure look successful.

## Restore and Clone Semantics

**Decision**: A database restore preserves identity. A copied database remains the same logical daemon until an operator deliberately resets one copy before concurrent independent use.

**Rationale**: Backup and restore should recover the logical daemon as a unit. Automatically detecting a clone reliably would require machine fingerprinting or external coordination, both outside the privacy and local-only boundary.

**Alternatives considered**: Resetting on every restore breaks stable references. Resetting on host change requires prohibited host identity and misclassifies legitimate migration. Silent duplicate operation creates ambiguous targets, so the lifecycle must be explicit and documented.

## Display Name

**Decision**: Default to `go-schedule daemon`. Trim surrounding whitespace, require 1 through 80 Unicode code points, and reject Unicode control characters.

**Rationale**: A generic product name discloses no hostname or account. Code-point length matches what operators perceive more closely than byte length, while rejecting controls prevents confusing terminal, log, and UI rendering.

**Alternatives considered**: Host-derived defaults violate the privacy boundary. ASCII-only names unnecessarily exclude users. Byte-length limits treat equivalent visible names inconsistently.

## Capability Manifest

**Decision**: Expose a bounded local manifest with identity, name, product version, local API versions, empty remote API versions, operating mode, sorted feature capabilities, and OS/architecture only. Use the stable capability vocabulary `activity`, `agent-access`, `chains`, `groups`, `notifications`, `schedule`, `tasks`, `triggers`, and `watchers`.

**Rationale**: These facts are sufficient for current client presentation and feature gating. Static sorted values are deterministic and avoid endpoint probing. OS and architecture support compatibility decisions without identifying a machine.

**Alternatives considered**: Reflecting registered routes would expose implementation details and make names unstable. Hostname, address, paths, account, environment, and configuration values exceed the minimum disclosure needed. Advertising a remote protocol before #168 would make an unshipped capability appear available.

## Local API Compatibility

**Decision**: Preserve `GET /v1/health` unchanged and add `GET /v1/manifest`, `PATCH /v1/manifest`, and `POST /v1/manifest/reset` on the existing protected local transport.

**Rationale**: Additive routes let old clients keep using health. A separate manifest makes its privacy contract reviewable. Rename and reset remain locally authorized in S075; shared actor enforcement belongs to #167.

**Alternatives considered**: Expanding health couples liveness to discovery and changes an established response. Replacing health breaks compatibility. Adding remote routes or authorization now would improperly absorb later issues.

## Reset Safety

**Decision**: Reset accepts the exact current installation identifier, compares and replaces it in one database transaction, generates a new random UUID, and preserves all other data.

**Rationale**: Exact acknowledgement makes the destructive identity change deliberate and naturally rejects stale retries. Transactional comparison closes the race between confirmation and mutation.

**Alternatives considered**: A boolean confirmation provides no target assurance. An unconditional reset makes accidental replay destructive. Deleting and recreating the row creates an observable gap and complicates failure recovery.

## Desktop Integration

**Decision**: Load daemon-owned identity, name, platform, version, and capabilities after local connectivity succeeds. Keep current local permissions client-side until #167 introduces actor-scoped authorization.

**Rationale**: This removes identity and capability guesses while respecting the issue boundary. Permissions are not daemon capabilities and cannot become authoritative before the authorization model exists.

**Alternatives considered**: Keeping hardcoded identity defeats #166. Moving permissions into the manifest would imply the actor model from #167. Falling back silently when manifest retrieval fails would mask an incomplete upgraded daemon.
