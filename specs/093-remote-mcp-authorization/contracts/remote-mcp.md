# Contract: Remote MCP Authorization

## Public paths when enabled

| Method | Path | Authentication | Purpose |
| --- | --- | --- | --- |
| GET | RFC 9728 path derived from the resource | None | Protected resource metadata |
| GET | `/.well-known/oauth-authorization-server` | None | Authorization server metadata |
| POST | `/oauth/token` | HTTP Basic MCP client credential | Client-credentials token exchange |
| POST | `/mcp` | Resource-bound access token | Stateless Streamable HTTP MCP request |

All paths are absent when remote MCP is disabled.

## Token request

Content type is `application/x-www-form-urlencoded`. Required fields are `grant_type=client_credentials` and the exact configured `resource`. `scope` is optional and defaults to `mcp:observe`. Unknown or duplicated fields fail closed. Client authentication uses one HTTP Basic header; body credentials are not accepted.

## Token response

```json
{
  "access_token": "opaque-short-lived-value",
  "token_type": "Bearer",
  "expires_in": 600,
  "scope": "mcp:observe"
}
```

The response and every authorization endpoint use `Cache-Control: no-store`.

## Error response

OAuth token errors use the standard JSON `error` field with a bounded optional `error_description`. Authentication failures use `invalid_client`; invalid resource, grant, or scope requests use the corresponding standard error without revealing credential existence.

Protected MCP failures use `401` or `403` plus a `WWW-Authenticate: Bearer` challenge containing the RFC 9728 metadata URL and applicable scope.
