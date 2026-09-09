# Research: Authenticated Localhost MCP

## Decision 1: Runtime-only daemon ownership

**Decision**: The running daemon owns one optional listener and exposes status, enable, rotate, and disable through its existing protected IPC API. Listener state, allowed origins, and credential material are not persisted.

**Rationale**: This directly satisfies default-off and restart-stale behavior, gives GUI work in #164 a safe control surface, and keeps authorization for changing network exposure behind the existing operating-system IPC boundary.

**Alternatives considered**: A standalone `gosched mcp http serve` process would make lifecycle and future desktop controls indirect; persistent configuration would violate the issue's restart fail-closed criterion; command-line credentials would leak through process inspection.

## Decision 2: Fixed numeric IPv4 loopback

**Decision**: Accept only a port and bind `tcp4` to `127.0.0.1:<port>`.

**Rationale**: A fixed numeric address removes ambiguous host/interface configuration, avoids dual-stack wildcard behavior, gives Host validation one canonical value, and remains broadly compatible with local MCP clients.

**Alternatives considered**: `localhost` depends on resolver and address-family behavior; `::1` plus IPv4 requires two coordinated listeners and broader Host normalization; caller-supplied hosts create unnecessary exposure risk.

## Decision 3: Stateless official Streamable HTTP

**Decision**: Use `mcp.NewStreamableHTTPHandler` with `Stateless: true`, `PropagateRequestCancellation: true`, the SDK localhost protection enabled, and `MaxRequestBodyBytes: 1 << 20`.

**Rationale**: The Observe surface has no server-to-client requests or durable subscriptions, and the SDK documents stateless mode as the sessionless direction with support for current protocol revisions. Stateless operation also avoids abandoned-session retention.

**Alternatives considered**: Stateful sessions add lifecycle and resource-retention risk without Observe value; legacy SSE is not the current transport; a custom transport would duplicate protocol logic already supplied by the official SDK.

## Decision 4: Digest-only bearer verification

**Decision**: Generate 32 random bytes, encode them with unpadded base64url, retain only SHA-256 of the credential, compare request digests with `subtle.ConstantTimeCompare`, and expose a short digest fingerprint in status.

**Rationale**: The credential has at least 256 bits of entropy, is copyable by clients, never needs recovery, and plaintext does not remain in daemon state after enable or rotation returns.

**Alternatives considered**: Persisted secrets conflict with restart staleness; plaintext retention increases memory disclosure impact; reusable IPC or trigger credentials would couple independent authority surfaces.

## Decision 5: Exact request context before protocol parsing

**Decision**: Require exact active Host, allow Origin absence for native clients, require exact membership for any present Origin, require exactly one canonical Bearer header, set `Cache-Control: no-store`, and then call the SDK handler.

**Rationale**: Local TCP peers are not proof of intent. Ordered middleware blocks DNS rebinding, hostile browser origins, and ambient unauthenticated access before request bodies or MCP messages are processed.

**Alternatives considered**: Permissive localhost aliases broaden rebinding edge cases; wildcard CORS grants browser authority; query-string credentials leak through URLs; SDK localhost protection alone does not provide client authentication or an explicit application origin allowlist.

## Decision 6: Existing Observe client as the only data path

**Decision**: Each HTTP request receives the same `mcpobserve.NewServer` construction used by stdio, backed by an existing local IPC client.

**Rationale**: One registration and mapping path guarantees no transport-specific tools, redaction, pagination, or scheduler authority. The daemon remains the sole source of scheduler truth.

**Alternatives considered**: Direct store reads bypass API policy; duplicated HTTP resource handlers invite schema drift; proxying raw MCP frames from another subprocess adds needless processes and failure modes.
