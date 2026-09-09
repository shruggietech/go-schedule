# Contract: Agent Access Workspace and MCP Qualification

## Additive protected local API fields

`POST /v1/mcp/http/enable` accepts the existing port and allowed-origin fields plus optional `client_name`. Existing callers that omit the field receive the compatibility name `Local MCP client`.

`GET /v1/mcp/http` and every lifecycle result add non-secret `client_name`, `request_count`, and optional `last_accessed_at`. Disabled status omits the name and timestamp and reports count zero. Enable and rotation responses still contain `credential` only for the direct protected API caller; subsequent status and disable responses never contain it.

Client-name validation failures use the existing `400 validation_failed` response with field `client_name`. No route or API version changes.

## Desktop bridge

- `AgentAccessWorkspace()` returns stdio guidance, localhost status, permission cards, and no credential.
- `EnableAgentHTTP(draft)` accepts a desktop draft with client name, port, and optional newline-separated origins. The desktop service parses and validates the draft, invokes the protected API, copies the returned credential natively, and returns safe refreshed status.
- `RotateAgentHTTP()` rotates the active credential, copies it natively, and returns safe refreshed status.
- `RevokeAgentHTTP()` disables localhost HTTP and returns Off status.
- `OpenAgentAccessGuide()` opens the fixed product documentation destination and never accepts an arbitrary URL.

Every action is local, deadline-bounded, serialized, and duplicate-safe. Enable and rotation clipboard failure trigger disablement before returning an unavailable outcome.

## Desktop presentation

The primary navigation label is `Agent Access`. The workspace has three vertically ordered sections:

1. Access overview: network listener Off or Active, with plain-language stdio availability.
2. Localhost client: enable form while Off; exact endpoint, name, fingerprint, origins, activation, last access, request count, rotate, and `Revoke and turn off` while Active.
3. Authority and setup: Observe marked Available; Operate and Manage marked Future and unavailable; fixed link to setup, privacy, and troubleshooting guidance.

New credentials are never rendered. Accepted enable and rotate messages state that the one-time credential was copied and must be stored in the intended client's protected configuration.

## Conformance gate

The official SDK gate covers protocol revisions `2026-07-28` and `2025-11-25`, exactly five static resources, four continuation templates, `application/json`, schema version `1`, Observe permission, trust notice and untrusted field metadata, safe error codes, bounded pages and strings, no protected fields, and zero tool discovery or invocation.

The package-shaped smoke builds the real `gosched` executable under a temporary path containing spaces and non-ASCII text, launches `mcp serve` with hidden Windows process flags, initializes via the official SDK command transport, discovers the exact Observe surface, and closes stdin with clean process termination. It does not require a signed release candidate or running daemon because resource reads are covered separately and formal artifact qualification belongs to #190.
