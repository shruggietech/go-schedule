# Research: Agent Access Controls and MCP Release Gates

## Single-client identity boundary

**Decision**: Add one optional bounded client name to S070 enablement and status, defaulting omitted API and CLI input to `Local MCP client` for backward compatibility.

**Rationale**: The active listener has one digest and rotation replaces it globally. Naming that one credential is honest and immediately useful; multiple names would imply independent grants that do not exist.

**Alternatives considered**: A persistent client table belongs to #181 and would require durable secrets, expiry, per-client revocation, and audit semantics. An unlabelled listener does not meet #164's operator-identification requirement.

## Access evidence

**Decision**: Store only last successful authenticated access time and a saturating `uint64` request count in manager memory. Clear both on enable, rotation, disable, shutdown, and restart.

**Rationale**: Aggregate evidence answers whether the copied credential has been used without storing content or inventing a durable audit subsystem. Saturation prevents wraparound from making long-lived activity look new or absent.

**Alternatives considered**: A request ring buffer adds privacy, bounds, and synchronization complexity without acceptance value. Logging every request is not an operator workspace contract and risks user-controlled metadata. Failed-attempt counters create an ambiguous credential oracle and are excluded.

## Credential handoff

**Decision**: Copy credentials in the trusted desktop service and return only status to React. Disable the endpoint if clipboard copy fails after enable or rotation.

**Rationale**: React does not need plaintext secret material. A failed copy otherwise leaves an active credential the user cannot recover, so fail-closed rollback is the only safe and actionable state.

**Alternatives considered**: Displaying the secret in React broadens retention and screenshot exposure. Leaving an inaccessible credential active forces a separate cleanup action and weakens the stated outcome.

## Desktop domain placement

**Decision**: Add a peer `desktop/agentaccess` service with an injected local API backend and native clipboard interface, then bind small facade methods through `App`.

**Rationale**: Existing task, automation, operations, notification, and settings packages establish this boundary. Agent access is daemon-owned runtime control, not a desktop preference or connection-diagnostic detail.

**Alternatives considered**: Adding the feature to Settings conflates desktop-local preferences with a daemon listener. Direct React-to-daemon calls bypass the tested service result and rollback boundary.

## Conformance and package smoke

**Decision**: Consolidate schema and authority invariants in `internal/mcpobserve` tests and add an integration test that builds and launches the actual CLI with the official SDK command transport from a package-shaped path.

**Rationale**: Existing unit tests prove most individual behaviors, but #164 asks for one release gate over the complete surface and a real built command. The current three-platform race matrix can execute the integration without workflow expansion.

**Alternatives considered**: Requiring signed release candidates would duplicate #190. Launching a test helper rather than the built CLI repeats S069 evidence and does not prove packaged command wiring.

## Documentation boundary

**Decision**: Keep one local MCP guide with separate Codex stdio, generic stdio, and generic HTTP sections, plus authority, privacy, removal, and troubleshooting language. Link it from the CLI and desktop workspace.

**Rationale**: One guide makes credential and listener differences comparable while avoiding contradictory setup fragments.

**Alternatives considered**: Per-client pages would duplicate security language. Embedding complete setup prose in the desktop is harder to maintain and unsuitable for copyable configuration examples.
