# Research: S098 Portable Automation Bundles

## Decision: Use canonical JSON with a declared `go-schedule.bundle/v1` schema

**Rationale**: Stable field order and collection sort order make review, diffing, digesting, fixtures, and plan binding reproducible without a new external dependency.

**Alternatives considered**:

- YAML: Human-friendly but has several equivalent encodings and weak canonicalization.
- Reusing daemon database records: Leaks internal IDs and platform execution inputs, and cannot travel safely.

## Decision: Keep only portable declarative intent in v1

**Rationale**: Commands, working directories, environments, run-as identities, watcher paths, notification destinations, and trigger keys are either machine-specific or secret-bearing. A safe first transfer format must reject rather than guess at them.

**Alternatives considered**:

- Copy every field: Would export credentials and turn a portable bundle into a host clone.
- Silently omit unsupported fields: Would create a misleading successful transfer.

## Decision: Bind apply to an explicit target-bound preview

**Rationale**: A bundle digest alone proves document equality, not that the operator reviewed the current selected daemon. The target ID and fingerprint prevent a stale plan from quietly mutating a different or changed daemon.

**Alternatives considered**:

- Apply raw bundle directly: Bypasses review and allows stale conflicts.
- Persist server-side plans: Adds cleanup and state management without improving the safety guarantee.

## Decision: Omission means drift, never removal

**Rationale**: A partial export must not delete unrelated target automation. Target-only state is visible as drift and remains unchanged.

**Alternatives considered**:

- Desired-state reconciliation: Contradicts the issue boundary and introduces hidden synchronization.

## Decision: Execute dependent operations in one stable order

**Rationale**: Groups and tasks must exist before chains, sources, and policy references. Each operation has a visible outcome, so a failure does not obscure later independent outcomes.

**Order**: groups, tasks, notification-policy references, chains, trigger sets, external triggers, watchers.
