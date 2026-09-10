# Research: Resilient Remote Connections

## Decision 1: Retry only connection negotiation and subscriptions

**Decision**: The desktop connection manager automatically retries remote health, manifest validation, and event stream establishment. Shared API feature methods remain single-attempt.

**Rationale**: HTTP method alone is not a sufficient replay-safety signal, and several POST routes trigger actions. Central connection retries recover observation without creating accidental repeated mutations.

**Alternatives considered**: A retrying HTTP transport was rejected because it could replay writes after ambiguous failures. Service-specific retries were rejected because they would duplicate policy across every desktop feature.

## Decision 2: Refresh authoritative state before reopening live invalidations

**Decision**: A recovered connection publishes a new connected generation, causing mounted workspaces to reload from their authoritative read endpoints before event-driven reloads resume.

**Rationale**: Current SSE events notify clients that an entity changed but are not the source of truth. Full workspace reads already exist, are bounded, and close any gap created while disconnected.

**Alternatives considered**: Durable event IDs and a replay buffer were rejected as new daemon persistence and protocol scope. Blind stream restart was rejected because changes could be missed during the outage.

## Decision 3: Bounded exponential backoff with jitter

**Decision**: Retry delays grow from a sub-second base to a 30-second cap and add bounded jitter supplied through an injectable policy seam.

**Rationale**: Fast first recovery feels responsive, while the cap prevents a disconnected target from generating excessive load. Injection keeps tests deterministic and avoids real sleeps.

**Alternatives considered**: Fixed cadence was rejected because many clients could synchronize after daemon restart. Unbounded exponential backoff was rejected because ordinary recovery would become too slow.

## Decision 4: Terminal failure taxonomy

**Decision**: Credential rejection or revocation, forbidden authority, API incompatibility, certificate validation, and daemon identity mismatch stop automatic recovery. Timeouts, refused connections, EOF, temporary DNS or network errors, and stream loss remain transient.

**Rationale**: Retrying a trust or authorization decision cannot repair it and can hide a security-relevant change. Transport recovery is expected after normal mobility and daemon lifecycle events.

**Alternatives considered**: Retrying every failure was rejected as noisy and unsafe. Making every failure manual was rejected because it preserves S078's core usability gap.

## Decision 5: Typed uncertain mutation result

**Decision**: When a remote state-changing request fails at the transport layer, the client returns a typed uncertainty that carries only the operation name and safe guidance. It never retries.

**Rationale**: The daemon may have committed before the response vanished. An explicit type lets each bridge service distinguish uncertainty from validation rejection and local unavailability without exposing raw transport details.

**Alternatives considered**: Treating every failure as unavailable was rejected because it encourages blind retries. Idempotency keys were deferred because retrofitting durable deduplication across all mutation routes is larger than this slice.

## Decision 6: Persist only last successful contact

**Decision**: Successful remote negotiation updates the existing profile `last_successful_at` field through the profile store, using the selected profile ID and generation-safe callback.

**Rationale**: This supplies durable operator context without adding cached daemon data or secrets to disk.

**Alternatives considered**: Session-only timestamps were rejected because restart would erase useful diagnostics. Persisting workspace data was rejected as an offline cache outside scope.
