# Contract: Localhost MCP HTTP Control and Transport

## Protected local IPC API

### `GET /v1/mcp/http`

Returns non-secret endpoint status. Disabled status is `{"enabled":false,"allowed_origins":[]}`.

### `POST /v1/mcp/http/enable`

Request: `{"port":43123,"allowed_origins":["http://127.0.0.1:3000"]}`.

Success returns status plus `credential`, which is disclosed only in this response. Invalid port or origin returns `400 validation_failed`; already enabled or occupied port returns `409 conflict` with no partial listener.

### `POST /v1/mcp/http/rotate`

Success preserves endpoint and origins and returns status plus a replacement `credential`. Disabled state returns `409 conflict`.

### `POST /v1/mcp/http/disable`

Revokes the credential and stops the listener. Repetition succeeds and returns disabled status.

## Public Streamable HTTP endpoint

The only route is `POST http://127.0.0.1:<port>/mcp`. The handler requires exact Host, an optional exact allowed Origin, exactly one `Authorization: Bearer <credential>` value, MCP `Content-Type` and `Accept` headers, and a body no larger than 1 MiB. Host or Origin rejection uses `403`; authorization rejection uses `401` and `WWW-Authenticate: Bearer`; protocol and media errors remain official SDK responses. Responses use `Cache-Control: no-store` and do not grant wildcard CORS.

The SDK handler is stateless and creates the existing Observe server for each request. The transport advertises resources and templates only, and its resource payloads follow [S069's Observe contract](../../069-local-mcp-observe/contracts/observe-resources.md).

## CLI

- `gosched mcp http status` prints disabled or the non-secret active status.
- `gosched mcp http enable --port <port> [--origin <origin>...]` prints endpoint, one-time credential, and allowed origins.
- `gosched mcp http rotate` prints the replacement one-time credential and fingerprint.
- `gosched mcp http disable` prints confirmation and is idempotent.
- Global `--json` returns the API response shape. Credentials are never accepted as arguments or emitted by status and disable.
