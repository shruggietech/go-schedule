# Research: Local Observe-Only MCP

## Official SDK and protocol baseline

**Decision**: Use `github.com/modelcontextprotocol/go-sdk/mcp` v1.7.0 and test negotiated protocol revisions `2026-07-28` and `2025-11-25`.

**Rationale**: The official Go SDK provides server lifecycle, stdio framing, resource registration, capability negotiation, cancellation, and protocol-version compatibility. Version 1.7.0 is the current stable release and supports both selected revisions. Using it avoids a bespoke protocol implementation and directly satisfies issue #161.

**Sources**: [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk), [server guide](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/server.md), [protocol compatibility](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/protocol.md).

**Alternatives considered**: A hand-written JSON-RPC server was rejected because it would duplicate negotiation, framing, lifecycle, and compatibility logic. An unofficial SDK was rejected because the issue explicitly requires the official SDK.

## Transport and entry point

**Decision**: Add `gosched mcp serve` using the SDK's stdio transport. Do not add a daemon listener or a separate installed executable.

**Rationale**: The existing `gosched` binary is already present on every supported installation. A subcommand is an explicit opt-in entry point, preserves stdout for protocol messages, inherits the launching identity, and avoids release-packaging changes. The MCP process reaches scheduler state only through the same protected local IPC client used by the CLI and GUI.

**Alternatives considered**: A TCP or streamable HTTP listener was rejected because it expands the attack surface and violates issue #162. A new `gosched-mcp` binary was rejected because it adds packaging and installation complexity without improving isolation for this read-only surface.

## Resource surface

**Decision**: Register first-page resources for health, active tasks, upcoming schedules, recent alerts, and recent runs. Register continuation templates for the four collections. Advertise zero tools.

**Rationale**: The five resources match the approved Observe scope and can be explained on one page. Separate task and schedule views prevent a client from inferring that task execution configuration is available. Resource templates permit bounded continuation without introducing arbitrary query tools.

**Alternatives considered**: One omnibus resource was rejected because partial failures and trust classifications become unclear. General query tools were rejected because they widen authority and make bounds harder to audit. Raw logs were deferred because they add path and structured-attribute leakage risk without being required by #162.

## Redaction and untrusted content

**Decision**: Map daemon responses into dedicated safe structs, cap display text at 2 KiB and output at 8 KiB, normalize invalid UTF-8, and include an envelope-level trust notice plus field-level untrusted metadata.

**Rationale**: Allowlisting prevents newly added domain fields from crossing the boundary accidentally. JSON encoding keeps content inert, while explicit trust labels tell agent hosts that names, descriptions, messages, and output are data rather than instructions. Bounds prevent oversized or adversarial content from dominating an agent context.

**Alternatives considered**: Reflection-based denylisting was rejected because new sensitive fields fail open. Removing output entirely was rejected because bounded run evidence is central to useful Observe behavior.

## Pagination

**Decision**: Sort each collection by stable keys, request at most 101 records for one page, return at most 100 items, and encode a versioned offset as an opaque URL-safe cursor in a continuation URI. Run output and alert messages are truncated in SQLite before entering daemon memory or IPC.

**Rationale**: Page-sized daemon reads and store-side text projection provide predictable memory, IPC, and agent-context sizes even when persisted output is extremely large. An opaque versioned cursor lets the representation evolve and rejects malformed input cleanly.

**Alternatives considered**: Unbounded arrays were rejected for context and denial-of-service risk. Database-backed snapshot cursors were rejected as disproportionate for local observation and would put MCP concerns into persistence.

## Error model and shutdown

**Decision**: Give each resource read a derived timeout no longer than 10 seconds, preserve cancellation, map client error categories into a bounded public vocabulary, and let SDK stdio EOF terminate the session.

**Rationale**: This keeps host behavior deterministic and avoids revealing configured endpoints, filesystem paths, or internal daemon messages. The adapter does not retry denied access or seek another transport.

## Future permissions

**Decision**: Document Observe, Operate, and Manage, but activate only Observe. Operate requires authenticated identity, per-action authorization, attributable audit records, and explicit arguments. Manage additionally requires deliberate confirmation for destructive or externally visible effects.

**Rationale**: A named future boundary prevents read-only convenience from silently becoming mutation authority. Authentication alone is insufficient for high-impact actions, so authorization, audit attribution, and confirmation remain independent gates.
