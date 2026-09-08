# Research: Dependable Webhook Notifications

## Outbound HTTP stack

**Decision**: Use Go's standard `net/http`, `net/url`, and `encoding/json` packages with one purpose-built `http.Client` that refuses redirects.

**Rationale**: These maintained production packages already provide request contexts, timeouts, URL parsing, headers, JSON encoding, and portable transport. A vendor SDK or retry library would add dependency and policy surface for a fixed three-attempt schedule.

**Alternatives considered**: Resty and retryablehttp offer broader convenience but duplicate a small, fixed contract and enlarge supply-chain scope. Vendor-specific clients would undermine the generic webhook promise.

## Secret protection

**Decision**: Store full endpoint and optional authorization values as write-only SQLite columns protected with Windows DPAPI under the daemon service identity on Windows. On Linux and macOS, restrict the daemon data directory to mode 0700 and database to mode 0600. Expose only sanitized scheme-and-host summaries and boolean credential presence. Clear secret snapshots from terminal deliveries and erase them on channel removal.

**Rationale**: Windows ProgramData inheritance can vary, so DPAPI supplies service-identity protection through the already-pinned system binding without custom cryptography or a new dependency. Unix daemon-only permissions fit unattended system-service operation. The design states the different platform guarantees explicitly and does not claim universal application-level encryption at rest.

**Alternatives considered**: Cross-platform desktop keyrings are frequently unavailable to headless system services and add platform-specific unlock/session failure modes. Custom encryption merely moves key custody and is prohibited. File permissions alone were rejected on Windows because inherited ProgramData ACLs are not a sufficiently precise credential boundary. Task environment fields would leak channel credentials into the wrong lifecycle and are explicitly rejected.

## Policy inheritance

**Decision**: The nearest scope containing any assignments replaces the complete inherited policy. Task scope wins, then nearest ancestor group, then outward ancestors.

**Rationale**: Replacement is explainable from one source scope and prevents accidental duplicate delivery when administrators add a task exception. It also permits multiple channels at the selected scope.

**Alternatives considered**: Union inheritance makes suppression difficult and can multiply notifications through deep group trees. Per-channel overlay is expressive but creates a more complex deletion and explanation model than the first release needs.

## Durable delivery boundary

**Decision**: Extend the existing run transaction to create matching delivery rows before commit, then wake a dedicated notification runtime after the engine's post-commit callback.

**Rationale**: A successfully recorded run cannot lose its notification intent between commits, while no outbound work occurs inside the transaction or task worker. The existing completion-chain transaction demonstrates the repository's durable fan-out pattern.

**Alternatives considered**: An in-memory callback loses work on crash. Calling HTTP inline couples receiver latency to scheduler workers. A message broker adds unnecessary infrastructure for local single-daemon scale.

## Retry and duplicate semantics

**Decision**: Provide at-least-once delivery with three total attempts, a 5-second request timeout, 1 and 2-second backoff, no redirects, and a stable delivery ID on every attempt.

**Rationale**: The limits cap shutdown and outage cost while handling transient failures. A stable idempotency key gives receivers a standard deduplication option when a successful response is lost.

**Alternatives considered**: Exactly-once delivery is impossible across an ordinary HTTP boundary without receiver participation. Unbounded exponential retries create permanent work and disk growth. A single attempt is too fragile for routine transient failures.

## History and deletion

**Decision**: Retain the newest 1,000 terminal deliveries, never prune pending or claimed work, and keep immutable safe snapshots after channel or task deletion. Channel deletion atomically removes assignments and non-terminal work and clears the channel secret.

**Rationale**: Operators retain useful evidence while storage remains bounded and deleted credentials cannot continue sending. Safe snapshots preserve historical meaning without foreign-key lifetime coupling.

**Alternatives considered**: Cascade-deleting all history impairs incident diagnosis. Retaining endpoint and authorization snapshots after terminal completion creates needless secret copies. Time-based retention depends on a policy clock and is less predictable than a fixed first-release cap.
