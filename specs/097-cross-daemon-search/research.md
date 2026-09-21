# Research: Cross-Daemon Search and Target-Safe Actions

## Decision: Use one daemon-owned bounded search contract

**Rationale**: Existing observation endpoints are individually safe but would require several requests and transfer unrelated records for every daemon. A single endpoint lets SQLite apply matching and limits at the authoritative source, returns one observation timestamp, and preserves partial target isolation.

**Alternatives considered**:

- Fan out existing tasks, groups, runs, calendar, and alert endpoints, then filter in the desktop. Rejected because it multiplies requests and transfers unbounded or loosely bounded data.
- Maintain a desktop fleet index. Rejected because it introduces durable stale state, synchronization, migrations, and background networking outside issue #183.

## Decision: Search only safe allowlisted fields

**Rationale**: Names, identifiers, group paths, state, readiness, schedule summaries, occurrence times, failure times, alert kind, and severity are sufficient for discovery. Commands, arguments, environments, output, alert messages, credentials, and endpoints can contain sensitive data and are unnecessary.

**Alternatives considered**:

- Full-text search over run output and alert messages. Rejected because it expands the secret boundary and response volume.
- Exact identifier lookup only. Rejected because it does not satisfy operator discovery by human-readable names.

## Decision: Reuse existing Operate-authorized mutation endpoints

**Rationale**: Task enable, disable, run-now, and alert acknowledgement already carry authorization, audit, retry, and error semantics. A new cross-daemon batch API would duplicate those controls and falsely suggest a transaction boundary.

**Alternatives considered**:

- Add a daemon batch mutation endpoint. Rejected because batching has no value across independent daemons and complicates per-object uncertainty.
- Switch the global desktop connection before every action. Rejected because it changes unrelated workspace state and increases wrong-target race risk.

## Decision: Revalidate target and object before mutation

**Rationale**: Search observations are stale immediately after creation. Reconstructing the immutable registration, confirming pinned daemon identity and Operate authority, and reloading the exact object closes the time-of-check to time-of-use gap as far as the independent daemon contract permits.

**Alternatives considered**:

- Trust the result's observation timestamp for a freshness window. Rejected because time alone cannot detect profile edits, identity replacement, revocation, deletion, or state changes.
- Require the operator to rerun search manually. Rejected because it moves a security invariant into fallible user procedure.

## Decision: Permit deliberate multi-target actions without rollback

**Rationale**: Issue #183 explicitly requires per-target partial failure reporting. One action can span targets only when every selected object supports it, confirmation groups every target, and execution returns an outcome per object. Success on one daemon is never reversed because another daemon fails.

**Alternatives considered**:

- Restrict all actions to one target. Rejected because it leaves the multi-daemon partial-outcome requirement incomplete.
- Provide compensating rollback. Rejected because run-now and alert acknowledgement are not generally reversible and independent daemons do not share a transaction coordinator.

## Decision: Bound fan-out at eight targets with three-second target deadlines

**Rationale**: Search is explicit and interactive. Eight workers improve progressive coverage for up to 100 registrations while remaining small enough to bound sockets, goroutines, and local resource pressure. A three-second deadline matches existing desktop operation timeouts and keeps failures actionable.

**Alternatives considered**:

- Reuse the All Systems four-worker, five-second bounds. Rejected because worst-case interactive search completion is unnecessarily slow.
- Launch one goroutine per profile. Rejected because it creates a large burst and weakens lifecycle control.

## Decision: Add Search as a primary desktop destination

**Rationale**: Search is a fleet-level workflow like All Systems, not a filter inside one selected daemon page. A dedicated route can preserve query, progressive results, selection, confirmation, and outcomes while other pages remain target-specific.

**Alternatives considered**:

- Add search controls to All Systems. Rejected because it would mix operational summary sorting with object discovery and actions.
- Add global search to the title bar. Rejected because the confirmation and partial-outcome workflow needs substantial dedicated space.
