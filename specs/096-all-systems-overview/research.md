# Research: All Systems Operational Overview

## Decision 1: Compute summaries inside each daemon

**Decision**: Add one `GET /v1/system-summary` endpoint that returns a bounded daemon-local observation.

**Rationale**: The daemon owns its clock, schedules, task state, alerts, and delivery history. Computing locally keeps semantics consistent, prevents raw history fan-out, and makes local IPC and remote HTTPS behavior identical.

**Alternatives considered**:

- Reuse existing list endpoints from the desktop. Rejected because task and history volume would become the fan-out cost, secrets and run output would be easier to expose accidentally, and the desktop would duplicate scheduling rules.
- Add a persistent fleet database. Rejected because S096 is an observation surface for independent daemons, not clustered state or durable monitoring.

## Decision 2: Use bounded aggregate queries

**Decision**: Add store methods that count matching rows with indexed time and state predicates and fetch only one representative row per category.

**Rationale**: Counts remain authoritative while response size and Go allocations remain constant as history grows. Existing run, alert, and notification indexes support the required predicates; any missing composite index is added only if query evidence requires it.

**Alternatives considered**:

- Load existing paged endpoints and count in Go. Rejected because page limits would make counts inaccurate or total work unbounded.
- Maintain summary counters transactionally. Rejected because it would expand every mutation path and add migration risk for a read-only feature.

## Decision 3: Reuse the connection failure taxonomy

**Decision**: Temporary overview clients use the existing local and remote connection backends so identity, trust, authorization, compatibility, and timeout failures map to the same typed states used elsewhere.

**Rationale**: Operators should see one vocabulary and one recovery action for the same failure. Reimplementing TLS, daemon identity, or permission classification would create security drift.

**Alternatives considered**:

- Infer state from error strings. Rejected because string parsing is brittle and could collapse security-significant distinctions.
- Select each profile through the global router. Rejected because background refresh would repeatedly redirect the desktop's active mutation target.

## Decision 4: Keep refresh orchestration in the desktop backend

**Decision**: A new `desktop/systems` service snapshots registrations, builds independent clients, enforces four-worker concurrency and five-second per-target contexts, owns generation cancellation, and retains session-only successful summaries.

**Rationale**: Go provides deterministic concurrency, cancellation, and typed error handling without tying refresh correctness to React component lifetime. The frontend receives one coherent generation and remains responsible for display state only.

**Alternatives considered**:

- Fan out from React through many bridge calls. Rejected because cancellation and bounded concurrency would be duplicated in JavaScript and desktop credentials would need a broader bridge surface.
- Refresh sequentially. Rejected because one timeout would delay every later target and violate independent progress expectations.

## Decision 5: Drill down only after exact connection selection

**Decision**: A drill-down intent uses the local key or saved profile ID, destination route, and optional source identifiers. The existing selection service must accept that exact registration before navigation occurs.

**Rationale**: Display names are mutable and non-unique. Profile IDs preserve trust and authority boundaries, while delayed navigation prevents stale observations from opening controls against the wrong daemon.

**Alternatives considered**:

- Navigate first and connect afterward. Rejected because mutable controls could briefly render for the prior daemon.
- Deduplicate profiles by daemon identity. Rejected because separately saved registrations can intentionally carry different credentials or trust state.

## Decision 6: Present a dense responsive table with expandable detail

**Decision**: Use shared filter controls, stable sorting, a compact desktop table, and responsive stacked rows and details at narrow or high-zoom layouts. Status text and icons supplement color, and refresh announcements use an ARIA live region without moving focus.

**Rationale**: Up to 101 registrations require comparison density, while the existing desktop accessibility conventions require keyboard reachability, no document-level horizontal overflow, and explicit status language.

**Alternatives considered**:

- Render one large card per daemon at all widths. Rejected because scanning 100 systems would become excessively long and comparison would be poor.
- Use horizontal scrolling for the table. Rejected because it conflicts with the 800 by 600 and 200 percent zoom success criterion.

## Decision 7: Avoid new dependencies

**Decision**: Implement aggregation, concurrency, sorting, and responsive presentation with the Go and frontend libraries already present.

**Rationale**: The standard library and current stack fully cover the feature. A new data-grid, cache, or concurrency dependency would add review and supply-chain cost without proportional value.
