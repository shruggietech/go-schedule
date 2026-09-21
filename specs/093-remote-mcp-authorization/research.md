# Research: Remote MCP Authorization

## Current protocol baseline

MCP specification revision `2026-07-28` is the current final revision. It makes HTTP MCP stateless, requires per-request protocol metadata and standard routing headers, retains Streamable HTTP, and hardens authorization. The pinned official Go SDK 1.8.0 supports `2026-07-28` plus compatibility revisions and provides maintained protected-resource metadata and bearer middleware.

## Authorization profile

The core MCP authorization specification requires an HTTP MCP resource server to publish RFC 9728 metadata, identify an authorization server, accept bearer tokens only from the authorization header, validate token audience against the exact RFC 8707 resource, and return standard challenges. The specification leaves authorization-server implementation outside the resource-server contract.

S093 uses the MCP client-credentials authorization extension because go-schedule agents are unattended machine clients and no end-user login system exists. The persistent MCP credential is the OAuth client credential. A successful standard `client_credentials` exchange returns a separate short-lived resource-bound access token. This reuses actor identity without allowing the durable client secret to cross directly into MCP requests.

## Scope model

Scopes are `mcp:observe`, `mcp:operate`, and `mcp:manage`. Tokens contain the requested scope and every lower scope so the bearer middleware can require Observe uniformly. The source actor capability is an upper bound. Enroll is intentionally not an MCP scope.

## Resource and deployment identity

RFC 8707 requires the authorization and token requests to name the exact MCP resource. RFC 9728 derives metadata from that resource. A configured `https://host[:port]/mcp` resource is therefore authoritative; the daemon bind address is only a listener detail. Direct and private deployments normally use matching values, while reverse proxies use the public proxy resource URL.

## Rejected alternatives

- Reuse the durable remote credential as the MCP bearer: rejected because it crosses API and MCP boundaries, lacks short access-token expiry, and is not explicitly resource-bound.
- Accept pairing phrases at the token endpoint: rejected because pairing is one-time enrollment proof, not an OAuth client secret or access token.
- Build an interactive authorization-code server: rejected because the product has no user account or login authority and issue #181 owns human grant administration.
- Require an external identity provider: rejected because it would make remote MCP unavailable to the existing self-hosted deployment and would create an unrelated identity dependency.
- Create a runtime MCP actor for every token: rejected because audit attribution would no longer identify the persistent remote MCP client.

## Primary references

- MCP Authorization, revision `2026-07-28`: https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization
- MCP transports, revision `2026-07-28`: https://modelcontextprotocol.io/specification/2026-07-28/basic/transports
- MCP security best practices: https://modelcontextprotocol.io/docs/2026-07-28/tutorials/security/security_best_practices
- RFC 9728, OAuth 2.0 Protected Resource Metadata: https://www.rfc-editor.org/rfc/rfc9728
- RFC 8707, Resource Indicators for OAuth 2.0: https://www.rfc-editor.org/rfc/rfc8707
- RFC 6750, OAuth 2.0 Bearer Token Usage: https://www.rfc-editor.org/rfc/rfc6750
