# Research: Trigger Sets

## Membership Model

**Decision**: Store one Trigger Set entity and optional set identity plus permanent position on each ordinary external trigger.

**Rationale**: A trigger belongs to zero or one set, firing behavior is identical for both member and standalone triggers, and the database can enforce unique positions per set without duplicating credentials or dispatch logic.

**Alternatives considered**: A separate membership table was rejected because many-to-many membership is out of scope and would create an additional deletion and consistency surface. A JSON member list was rejected because it weakens foreign keys, indexed lookup, and transactional member operations.

## Transaction Boundary

**Decision**: Create, retarget, enable, disable, rotate, and delete complete sets inside one database transaction, including automatic-source readiness effects.

**Rationale**: The issue permits either atomic operations or precise partial results. Atomicity is simpler for callers, prevents mixed-key generations, and matches the single-writer SQLite architecture.

**Alternatives considered**: Per-member API loops were rejected because caller interruption and intermediate errors create ambiguous partial state. Compensating rollback was rejected because SQLite already supplies the required transaction boundary.

## Ordered Bulk Output

**Decision**: Assign positions 1 through the creation count once, retain gaps after deletion, and emit secrets in ascending position. Human output is exactly one complete command per line with one final newline; JSON includes set ID, position, trigger ID, key, and command.

**Rationale**: Permanent positions keep external caller assignments stable and make repeated bulk copy byte-deterministic.

**Alternatives considered**: Ordering by mutable name or generated trigger ID was rejected because neither communicates creation order. Dense renumbering after deletion was rejected because it changes unaffected members.

## Interface Shape

**Decision**: Add versioned `/v1/trigger-sets` resources and nested `gosched trigger set` commands while extending ordinary trigger responses with optional set metadata.

**Rationale**: Both interfaces preserve the existing trigger vocabulary and let the desktop use the same contracts as scripted administration.

**Alternatives considered**: A top-level CLI noun was rejected as inconsistent with the established command hierarchy. Client-side batching of individual endpoints was rejected because it cannot guarantee atomicity.

## Event Propagation

**Decision**: Publish one redacted Trigger Set lifecycle event for each broad mutation and let connected desktop clients refresh their cached triggers and sets once.

**Rationale**: Publishing a member event for every row would turn a maximum-size mutation into a refresh storm and expose implementation sequencing.

**Alternatives considered**: Suppressing events was rejected because the GUI must converge live. Publishing 99 individual events was rejected for avoidable load and confusing intermediate views.

## Performance Budget

**Decision**: Benchmark maximum-size create, reveal, retarget, state change, rotate, and delete operations against the one-second nominal budget.

**Rationale**: The feature is bounded at 99 members, so a direct transaction with indexed queries should remain far below the user-visible threshold without speculative optimization.

**Alternatives considered**: A per-member latency budget was rejected because users experience the complete action. A tighter unmeasured target was rejected because the constitution requires evidence-based performance claims.
