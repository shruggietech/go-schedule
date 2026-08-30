# Research: Task-Completion Chains

## R1 - Relationship Model

**Decision**: Model a completion chain as a separate relationship between existing tasks. It supplements the target's schedule instead of becoming a new schedule kind.

**Rationale**: Current scheduling, preview, catch-up, and editing assume one authoritative schedule. A separate relationship preserves those invariants and permits fan-in and fan-out.

**Alternatives considered**: Restoring the removed trigger subsystem would recreate discarded design and make completion mutually exclusive with timing. Embedding one source on the target cannot represent multiple edges cleanly.

## R2 - Delivery Guarantee

**Decision**: Durable at-least-once delivery. One `(chain, source run)` record suppresses duplicate event insertion. Pending and interrupted claims recover; completed deliveries never replay. A crash after command launch but before completion storage may replay.

**Rationale**: SQLite state cannot commit atomically with arbitrary external command side effects. Retrying ambiguous work avoids silent loss and matches master-spec FR-014.

**Alternatives considered**: Completing before launch can lose work. Never retrying claims can also lose work. Claiming exactly-once external effects would be technically false.

## R3 - Atomic Recording

**Decision**: Insert a terminal run, complete its incoming delivery when present, and create matching outgoing deliveries in one transaction.

**Rationale**: A committed source success can never silently omit downstream work, and every target can cascade through the same path.

**Alternatives considered**: An in-memory post-run hook loses events on process failure. Periodic history scans add polling, duplicate complexity, and unbounded work.

## R4 - Claim Ownership

**Decision**: The engine owns recovery, claiming, eligibility, and dispatch; the store owns transactions and state transitions. Process claims at start and after terminal completions, while commands stay bounded by the worker pool.

**Rationale**: No extra long-lived goroutine or timer is needed. The engine already owns overlap, shutdown, and task eligibility.

**Alternatives considered**: A delivery service duplicates lifecycle logic. API-driven processing fails with no connected client.

## R5 - Dispatch Origin

**Decision**: Use an internal origin value carrying trigger, source task, source run, and delivery IDs. Keep the external Runner contract unchanged and enrich its returned run in the engine.

**Rationale**: Executors run commands; scheduler bookkeeping stays in the engine. This minimizes platform execution churn.

**Alternatives considered**: Expanding every Runner leaks persistence correlation into execution. Log-only correlation cannot support recovery or structured history.

## R6 - Cycle Safety

**Decision**: Validate the prospective directed graph during create and update. Reject self-links and any edge whose target can already reach its source.

**Rationale**: Acyclic graphs guarantee finite cascades while allowing fan-in, fan-out, and long pipelines.

**Alternatives considered**: Runtime depth limits permit invalid persistent graphs. Per-run dedup does not stop cycles because every downstream run has a new identity.

## R7 - Ineligible Targets

**Decision**: If the target is missing, inactive, disabled, or under a disabled group when claimed, resolve without execution and write diagnostic evidence.

**Rationale**: Disablement means do not run. Retaining events creates a hidden unbounded queue and surprising delayed commands.

**Alternatives considered**: Waiting for re-enable is unbounded. Silent deletion is undiagnosable.

## R8 - User Surfaces

**Decision**: Add `gosched chain` lifecycle commands, local `/v1/chains` routes, and a dedicated desktop Chains tab. Mutations use stable IDs and presentation includes names.

**Rationale**: This follows current task/group patterns and keeps relationships discoverable without crowding task schedule editing.

**Alternatives considered**: Advanced-editor controls overload an already complex form. CLI-only delivery makes the desktop incomplete.
