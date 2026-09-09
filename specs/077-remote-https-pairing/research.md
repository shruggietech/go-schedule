# Research: Authenticated Remote Access and Pairing

## Decision 1: Combine transport and enrollment

**Decision**: Implement #168 and #169 together so the HTTPS boundary ships with its durable authentication path and the enrollment path is exercised against its real transport.

**Rationale**: Shipping HTTPS with a temporary shared secret would create throwaway security behavior. Shipping pairing without the remote endpoint would leave the critical exchange untestable.

**Alternatives considered**: A configured static token was rejected because it bypasses named actors, independent revocation, and the approved credential lifecycle. Pairing-only storage was rejected because it has no end-to-end client value.

## Decision 2: Reuse application handlers through an explicit adapter

**Decision**: Build a separate remote route table. Each record names a remote path, approved local equivalent, operation ID, capability, target class, audit class, retry class, body limit, and secret exclusions. The adapter rewrites only matched requests into a separately configured local API server instance.

**Rationale**: This preserves identical domain behavior without mounting the local mux wholesale. The explicit table keeps runtime paths, MCP controls, and secret endpoints unreachable by omission.

**Alternatives considered**: Duplicating handlers was rejected because semantics would drift. Prefix-rewriting every local route was rejected because future local additions would become remote accidentally.

## Decision 3: Keep bearer authority server-owned

**Decision**: Generate 32 random bytes, return base64url once, persist SHA-256 only, and load credential plus actor state per request.

**Rationale**: Indexed digest lookup is sufficient for high-entropy values and immediate revocation. Token claims and signing add no value for one daemon that is always online for authorization.

**Alternatives considered**: JWTs, reversible tokens, shared passwords, and credential-derived authority were rejected by S074.

## Decision 4: Use ten random words and Argon2id

**Decision**: Select ten words independently from a 64-entry unambiguous list with uniform six-bit indexing, yielding exactly 60 bits. Store a unique 16-byte salt and a 32-byte Argon2id verifier using one iteration, 64 MiB, and four lanes. Each session expires after ten minutes and permits five failed attempts.

**Rationale**: The online attempt budget makes 60 bits ample while ten short words remain transferable. The RFC 9106 memory-constrained profile avoids weak fast digests and remains bounded for one enrollment request.

**Alternatives considered**: A shorter phrase was rejected below the explicit entropy requirement. Durable-strength Argon2 parameters were rejected because phrases are short-lived and already online-limited. Custom PAKE was rejected as unnecessary protocol design.

## Decision 5: Bound both unauthenticated and authenticated work

**Decision**: Use `golang.org/x/time/rate` v0.15.0 with source buckets before authentication and actor buckets afterward. Each map retains at most 4,096 buckets for the listener lifetime, and active requests are capped at 64. The newer v0.16.0 requires Go 1.26, so the latest release compatible with the repository's Go 1.25 baseline is selected.

**Rationale**: Two stages prevent guessing and isolate noisy valid actors without an unbounded identity map.

**Alternatives considered**: Proxy-only and global-only limits were rejected because direct HTTPS is supported and one client must not starve all others.

## Decision 6: Pin OpenAPI as reviewed source

**Decision**: Commit an OpenAPI 3.1 contract, generate the typed client boundary with pinned `oapi-codegen` v2.8.0, and verify its operation IDs against the runtime table while keeping security middleware explicit.

**Rationale**: S077 needs mechanical reachability and contract drift evidence. Authentication, limits, and operation policy remain inappropriate for generated code.

**Alternatives considered**: Prose-only documentation and annotation-derived contracts were rejected because neither makes the network allowlist independently reviewable.

## Decision 7: Isolate native credential storage

**Decision**: Add a minimal `clientsecret` adapter over `github.com/zalando/go-keyring` v0.2.8 and use it from a bounded desktop pairing flow. It stores only opaque credentials under daemon installation identity and client identifier, exposes no file fallback, and supports save, load, and delete for the follow-on connection-profile work. Argon2id uses `golang.org/x/crypto` v0.55.0, the latest release compatible with Go 1.25.

**Rationale**: This fulfills #169's durable client-side security and desktop enrollment boundary without prematurely defining persistent connection switching or target navigation.

**Alternatives considered**: Application files and three custom native implementations were rejected by S074. Deferring the adapter would leave #169 incomplete.

The new libraries retain permissive licensing: `x/crypto` and `x/time` use BSD-3-Clause, while `go-keyring`, `oapi-codegen`, and `oapi-codegen/runtime` use MIT. Their notices remain available through the Go module cache and upstream source distributions.

## Decision 8: Keep stream resilience bounded to authority

**Decision**: S077 revalidates an event stream every 15 seconds and closes it within 30 seconds of credential or actor invalidation. Resume identities, reconnect UX, stale state, and network-transition behavior remain in #172.

**Rationale**: Immediate authority loss belongs to authentication. Full connection resilience is a separate reviewed issue and should not expand this security slice.

**Alternatives considered**: Caching stream authority until disconnect was rejected as unsafe. Implementing the full #172 state machine was rejected as scope expansion.
