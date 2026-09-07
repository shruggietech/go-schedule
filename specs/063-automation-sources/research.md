# Research: Connected Automation Sources

## Decision 1: Compose existing APIs behind a desktop service

**Decision**: Add `desktop/automation` with a narrow backend interface and a service that builds one secret-free workspace.

**Rationale**: The daemon already owns complete, tested semantics for all four source types. A desktop adapter avoids duplicating persistence and keeps Wails contracts independent from HTTP response shapes.

**Alternatives considered**: Calling generated Wails methods directly from many React components was rejected because it disperses error handling, stale-write rules, and secret handling. Adding a new daemon aggregate endpoint was rejected because it expands the public API without delivery value.

## Decision 2: Use one complete snapshot

**Decision**: Load task choices and all four source collections into one `AutomationWorkspace`; publish it only when every read succeeds.

**Rationale**: A complete timestamped snapshot lets the UI explain cross-entity relationships consistently and prevents partial results from appearing current.

**Alternatives considered**: Independent section loading was rejected because one failed section creates ambiguous freshness. Per-row target lookups were rejected because they create avoidable latency and request amplification.

## Decision 3: Make secret responses exceptional and ephemeral

**Decision**: Exclude key fields from every workspace entity and ordinary operation result. Only create, reveal, and rotate return a distinct `SecretResult`. Fire accepts a trigger identifier and resolves the key inside the Go service.

**Rationale**: This preserves explicit user intent while keeping raw keys out of routine Wails payloads, React state stores, connection events, and diagnostics.

**Alternatives considered**: Returning masked keys in lists was rejected because masking invites accidental retention of the raw value. Passing a revealed key back from React to fire was rejected because it widens the secret lifetime and browser-visible surface.

## Decision 4: Apply optimistic concurrency in the desktop adapter

**Decision**: Update drafts carry `originalUpdatedAt`. Before update, the service fetches the current entity and rejects mismatches unless `overwriteStale` is explicit.

**Rationale**: Existing APIs do not expose conditional requests, while every entity exposes an update timestamp. The established S062 pattern provides predictable reload-or-overwrite behavior without persistence changes.

**Alternatives considered**: Last-write-wins was rejected because an open editor can silently erase a newer change. Extending every daemon API with revisions was rejected as disproportionate to this migration slice.

## Decision 5: Retain daemon-owned health and readiness

**Decision**: Map trigger readiness, watcher readiness, and watcher health directly from daemon responses. For chains, derive only relationship availability from returned source and target identifiers and names, without claiming execution health.

**Rationale**: The daemon observes runtime watcher state and task readiness. React cannot infer those facts accurately from configuration.

**Alternatives considered**: Treating enabled as ready was rejected because missing targets, incomplete commands, disabled tasks, and degraded paths are distinct conditions.

## Decision 6: Share navigation and collection language, keep editors specific

**Decision**: Add one `automation` route labeled Automation Sources with four semantic sections and shared search and attention filters. Each type retains a focused editor or action dialog.

**Rationale**: Users need one mental map of what initiates task execution, but forcing unrelated fields into one generic form would make validation and accessibility worse.

**Alternatives considered**: Four navigation destinations were rejected because they preserve the fragmented legacy model. One polymorphic mega-form was rejected because it hides source-specific constraints.

## Decision 7: Preserve last complete snapshot during disconnection

**Decision**: The React store retains the most recent successful workspace, marks it read-only when unavailable, and discards stale async responses using request sequence identifiers.

**Rationale**: Operators keep useful context without being invited to make changes that cannot succeed.

**Alternatives considered**: Clearing data on disconnect was rejected because it removes diagnostic context. Queuing offline mutations was rejected because the local daemon has no conflict-safe queue contract.
