# Feature Specification: Remote MCP Authorization

**Feature Branch**: `codex/093-remote-mcp-authorization`

**Created**: 2026-09-20

**Status**: Implemented

**Delivery**: implemented and locally verified on review branch `codex/093-remote-mcp-authorization`; publication and hosted review proceed through the S093 pull request for issue #180.

**Input**: Work slice S093, add explicitly enabled remote MCP over HTTPS with standard authorization.

## User Scenarios & Testing

### User Story 1 - Discover and authorize a remote MCP resource (Priority: P1)

An MCP client can discover the protected resource and authorization server metadata for one explicitly enabled remote go-schedule daemon, exchange a previously issued MCP client credential through the standard client-credentials grant, and connect through Streamable HTTP.

**Independent Test**: Start remote HTTPS with remote MCP enabled, fetch both metadata documents, request a token with the exact resource indicator, and complete MCP discovery using the issued bearer token.

**Acceptance Scenarios**:

1. **Given** remote access is disabled or remote MCP is disabled, **When** a client requests MCP or OAuth endpoints, **Then** the endpoints are absent.
2. **Given** both controls are enabled, **When** a client requests protected resource metadata, **Then** the response identifies the exact configured MCP resource, authorization server, bearer-header method, and supported scopes.
3. **Given** an active MCP credential, matching client ID, allowed scopes, and exact resource, **When** the token endpoint receives a client-credentials request, **Then** it returns a short-lived bearer token.
4. **Given** a pairing phrase, desktop, CLI, JSON credential, wrong client ID, wrong resource, or excessive scope, **When** token exchange is attempted, **Then** authorization fails without identifying which secret was wrong.

### User Story 2 - Enforce daemon and authority boundaries (Priority: P1)

Remote MCP exposes Observe, Operate, and Manage surfaces according to the issued grant while retaining the existing actor, authorization, validation, retry, redaction, and audit paths.

**Independent Test**: Exchange tokens for Observe, Operate, and Manage MCP actors, compare discovery, exercise allowed and denied calls, and inspect actor-attributed audit evidence.

**Acceptance Scenarios**:

1. **Given** an Observe token, **When** tools are listed, **Then** no mutation tool is exposed.
2. **Given** an Operate token, **When** tools are listed, **Then** exactly the three existing-task operations are added.
3. **Given** a Manage token, **When** tools are listed, **Then** Operate and the six bounded definition tools are available.
4. **Given** a token minted for another resource, client, or daemon, **When** it is presented, **Then** the request receives a standard authorization failure before MCP dispatch.
5. **Given** an accepted mutation, **When** audit history is inspected, **Then** it identifies the original persistent MCP actor rather than a transport-only surrogate.

### User Story 3 - Revoke, expire, and deploy safely (Priority: P1)

Operators can revoke or rotate the source MCP credential, rely on bounded access-token expiry, and deploy direct, private-network, or reverse-proxied HTTPS without weakening host, origin, TLS, or rate-limit controls.

**Independent Test**: Connect through each deployment shape, rotate or revoke the source credential, expire a token, interrupt transport, and retry with a newly issued token.

**Acceptance Scenarios**:

1. **Given** an active access token, **When** its source credential is rotated, revoked, or expires, **Then** the access token fails on its next request.
2. **Given** an expired access token, **When** it is presented, **Then** the server returns `401` with protected-resource discovery information.
3. **Given** browser-origin traffic, malformed host data, excessive request rate, query-string tokens, or cookies, **When** the remote endpoint receives it, **Then** the request fails before MCP dispatch.
4. **Given** supported and unsupported MCP protocol revisions, **When** version negotiation occurs, **Then** the official SDK accepts supported revisions and rejects unsupported revisions predictably.
5. **Given** network loss during a mutation, **When** the caller cannot prove the result, **Then** existing uncertain-outcome and request-ID reconciliation rules remain authoritative.

## Requirements

### Functional Requirements

- **FR-001**: Remote MCP MUST be absent unless `remote.enabled` and `remote.mcp.enabled` are both true.
- **FR-002**: Enabling remote MCP MUST require an explicit canonical HTTPS resource URL ending in `/mcp` with no query, fragment, or user information.
- **FR-003**: The server MUST publish RFC 9728 protected resource metadata using the official MCP Go SDK handler and authorization server metadata for the embedded client-credentials flow.
- **FR-004**: The authorization server MUST accept only `client_credentials`, HTTP Basic client authentication, an exact RFC 8707 resource indicator, and supported MCP scopes.
- **FR-005**: The client ID MUST equal the stored credential ID and the client secret MUST authenticate that exact credential.
- **FR-006**: Only credentials whose actor kind is `mcp` MAY obtain remote MCP access tokens; desktop, CLI, JSON, local, pairing, and runtime-session credentials MUST be rejected.
- **FR-007**: A pairing phrase MUST never be accepted as a client secret or MCP access token.
- **FR-008**: Issued access tokens MUST be cryptographically random, memory-only, bound to the exact resource, client credential, actor, daemon, and granted scopes, and expire after a bounded lifetime.
- **FR-009**: Every access-token verification MUST confirm the source credential and actor remain active so rotation, revocation, and expiry take effect without waiting for token expiry.
- **FR-010**: Tokens MUST be accepted only from the `Authorization` header and MUST be rejected from URI queries, cookies, duplicated headers, or malformed bearer schemes.
- **FR-011**: Observe, Operate, and Manage scopes MUST be monotonic and MUST NOT exceed the persistent MCP actor capability.
- **FR-012**: Remote MCP requests MUST reuse existing source and actor rate limits, concurrency bounds, host and origin rejection, TLS 1.3, request bounds, and no-store behavior.
- **FR-013**: MCP resources and tools MUST reuse the existing Observe, Operate, and Manage implementations without widening their schemas or redaction boundaries.
- **FR-014**: Mutation calls MUST execute through the ordinary versioned API with the persistent MCP actor so shared authorization and audit middleware remain authoritative.
- **FR-015**: The official MCP Go SDK MUST negotiate its supported current and compatibility protocol revisions; no custom wire protocol is permitted.
- **FR-016**: Authorization failures MUST use standard `401` or `403` responses and include a `WWW-Authenticate` protected-resource challenge when applicable.
- **FR-017**: Direct HTTPS, private-network HTTPS, and reverse-proxied HTTPS guidance MUST distinguish the configured canonical public resource from the listener bind address.

### Key Entities

- **Remote MCP configuration**: Explicit enablement and canonical HTTPS resource identity nested under remote listener configuration.
- **MCP client credential**: Existing persistent credential owned by an actor whose kind is `mcp` and whose capability bounds token scopes.
- **Access grant**: Memory-only record binding an opaque access-token digest to resource, daemon, actor, source credential, scopes, issue time, and expiry.
- **OAuth metadata**: Public RFC 9728 protected-resource document and authorization-server metadata derived from the configured resource.
- **Actor-bound API client**: In-process versioned API client whose request context resolves to the persistent MCP actor.

## Success Criteria

### Measurable Outcomes

- **SC-001**: When either enablement control is false, all remote MCP and OAuth paths return not found.
- **SC-002**: A standards client can discover metadata, acquire a resource-bound token, negotiate MCP, list resources, and use its allowed tools without a custom transport or authentication extension.
- **SC-003**: Observe discovers zero mutation tools, Operate discovers exactly three, and Manage discovers those three plus exactly six definition tools.
- **SC-004**: Wrong-daemon, wrong-resource, wrong-client, excessive-scope, expired, rotated, revoked, non-MCP, pairing-phrase, query-token, and cookie-token matrices all fail before MCP dispatch.
- **SC-005**: Accepted mutations retain persistent actor attribution and all existing redaction, retry, uncertain-outcome, and audit guarantees.
- **SC-006**: Focused race tests, official SDK interoperability, supported protocol versions, remote deployment matrices, and all canonical repository gates pass.

## Assumptions

- Machine-to-machine client credentials are the correct authorization profile for unattended remote schedulers and are supported by the MCP client-credentials authorization extension.
- An administrator still creates the persistent MCP actor and credential through the protected existing enrollment workflow; issue #181 will improve grant administration without changing this protocol boundary.
- Access tokens remain opaque and memory-only. Daemon restart intentionally invalidates them and clients obtain replacements with their still-active source credential.
- The resource URL is operator-configured because reverse proxies can present a public origin different from the daemon bind address.
