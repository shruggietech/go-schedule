# Research: Task Execution Safety and Diagnostics

## Decision 1: Use an optional initial-enabled create intent

**Decision**: Add an optional boolean creation intent. The server resolves an omitted value to enabled, while the desktop sends the explicit state of its creation-only checkbox.

**Rationale**: Create-then-disable introduces a real dispatch race after the server reloads the scheduler. Changing the universal default would silently break CLI and existing API clients. An optional value preserves omission as the old behavior and permits one atomic insert in the GUI-selected state.

**Alternatives considered**:

- Create then call disable. Rejected because a short runnable interval violates
  #118 and can execute startup or near-term schedules.
- Change all creation defaults to inactive. Rejected because #118 explicitly scopes the behavior to desktop creation and requires compatibility elsewhere.
- Add a second desktop-only endpoint. Rejected because an additive request field expresses the same domain intent with less surface area.

## Decision 2: Persist alert correlation and truncation in migration v10

**Decision**: Add nullable `alerts.run_id` and non-null `runs.output_truncated` with a false default. Keep the alert run ID as durable text without a foreign key.

**Rationale**: A failure alert currently has only a task ID and cannot identify one run among close failures. Output truncation is currently silently discarded inside the bounded buffer. Both facts must survive daemon restarts. A plain-text run ID remains available even if task deletion cascades away the run history, letting the interface state that enrichment is unavailable without erasing the identity users saw.

**Alternatives considered**:

- Encode a run ID in alert message text. Rejected as unstructured, fragile, and unsuitable for exact lookup.
- Add a foreign key with `ON DELETE SET NULL`. Rejected because it erases useful historical identity when the referenced run disappears.
- Append a truncation marker to output. Rejected because it either exceeds the configured cap or replaces retained process bytes.

## Decision 3: Retrieve a run by exact primary identity

**Decision**: Add an exact run lookup to store and the versioned local API.

**Rationale**: Listing recent runs and matching task or timestamp is ambiguous, especially for concurrent, overlap, or completion-triggered failures. A primary- key lookup is deterministic, cheap, and returns a normal not-found result after history deletion.

**Alternatives considered**:

- Load all recent runs into the desktop model. Rejected because the match still fails outside the bounded window and creates unnecessary refresh cost.
- Match the nearest run by timestamp. Rejected because it violates exact correlation and can disclose the wrong output.

## Decision 4: Retain combined process output

**Decision**: Present the existing capture as **Combined stdout/stderr**, with an explicit empty state and separate truncation flag.

**Rationale**: Both streams already write concurrently into one bounded buffer, so historical order cannot be split reliably after the fact. Honest combined labeling satisfies the issue without a larger storage/executor redesign.

**Alternatives considered**:

- Introduce separate stdout and stderr columns. Deferred because it changes capture semantics, ordering, storage volume, migration, and every run consumer.
- Infer stderr from failure status. Rejected because successful processes can write stderr and failed processes can write stdout.

## Decision 5: Extend shared group-chain reasoning

**Decision**: Add nearest-disabled-group discovery beside `ChainEnabled`, and have the GUI use both shared helpers after applying task/lifecycle precedence.

**Rationale**: Scheduler eligibility already treats any disabled ancestor and a cycle as ineligible. The GUI needs the nearest named cause but must not maintain a second policy implementation. A pure helper is deterministic, independently tested, and operates on the same group snapshot as the table.

**Alternatives considered**:

- Duplicate traversal in `gui/tasks.go`. Rejected because presentation and scheduler rules could drift.
- Add effective state to persisted tasks. Rejected because it is derived from mutable ancestors and would need redundant synchronization.
- Return only a boolean. Rejected because #120 requires the responsible group to be discoverable.

## Decision 6: Enrich Activity detail asynchronously

**Decision**: Open detail immediately from the alert's durable fields, then use exact background lookups to add run and optional task-name information. Format all states through a pure function that is headlessly testable.

**Rationale**: IPC must never block the desktop event thread. Durable IDs can be shown immediately, and exact lookup failures can update the same selectable detail with honest unavailable text.

**Alternatives considered**:

- Block while fetching before opening. Rejected because a slow/unavailable daemon would freeze the click path.
- Depend only on live run events. Rejected because alert delivery can precede the corresponding run event and application restarts clear that cache.
