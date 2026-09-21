# Implementation Plan: Remote MCP Authorization

**Branch**: `codex/093-remote-mcp-authorization` | **Date**: 2026-09-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/093-remote-mcp-authorization/spec.md`

## Summary

Mount a separately enabled Streamable HTTP MCP resource on the existing remote TLS server. Publish SDK-backed protected-resource metadata, provide a narrow OAuth client-credentials token endpoint, issue short-lived resource-bound opaque grants from existing MCP credentials, and create actor-bound in-process API clients so Observe, Operate, Manage, authorization, and audit behavior stay shared.

## Technical Context

**Language/Version**: Go 1.26.0

**Primary Dependencies**: official MCP Go SDK 1.8.0 `mcp`, `auth`, and `oauthex` packages; standard library HTTP, TLS, and crypto packages

**Storage**: Existing SQLite actors and client credentials; remote MCP access grants remain memory-only

**Testing**: Go unit, integration, race, SDK interoperability, protocol compatibility, deployment matrix, and canonical repository gates

**Target Platform**: Supported Windows, macOS, and Linux daemon packages using direct, private-network, or reverse-proxied HTTPS

**Performance Goals**: Constant-time token digest lookup, bounded grant count, existing 64-request global concurrency bound, and existing source and actor rate limits

**Constraints**: Default off; TLS 1.3; no pairing phrase bearer; no desktop, CLI, or JSON credential reuse; no second actor identity; no token persistence; no browser-origin access

**Scale/Scope**: One remote MCP resource per daemon, three monotonic scopes, and a bounded in-memory access-grant registry

## Constitution Check

- Existing actor, capability, API validation, audit, retry, and secret-redaction paths remain authoritative.
- OAuth metadata and bearer enforcement use maintained official MCP SDK packages.
- The token endpoint is intentionally limited to the machine-to-machine client-credentials grant and exact resource indicators.
- Configuration requires independent remote and remote-MCP opt-in and an explicit public resource identity.
- Access tokens are not accepted by the ordinary remote JSON API, and ordinary remote credentials are not accepted by MCP.
- Tests cover client and resource binding, source credential lifecycle, scope isolation, protocol compatibility, network failure, and deployment shapes.

All gates pass before implementation. Re-check after design: an embedded client-credentials authorization profile is smaller and safer than adding an interactive user account or authorization-code UI, and it matches unattended scheduler agents.

## Project Structure

### Documentation

```text
specs/093-remote-mcp-authorization/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── checklists/
└── contracts/
```

### Source Code

```text
internal/
├── config/              # explicit remote MCP resource configuration
├── enrollment/          # source credential lifecycle validation
├── api/client/          # actor-bound in-process API transport
├── remote/              # shared HTTPS limits and MCP route mounting
└── remotemcp/           # OAuth metadata, token grants, MCP dispatch

cmd/goschedd/            # compose remote API and MCP on one TLS listener
docs/                    # setup, security, and deployment guidance
```

**Structure Decision**: Keep OAuth and MCP resource behavior in `remotemcp`, let `remote` retain the shared network boundary, and compose both in the daemon without importing transport code into domain services.

## Complexity Tracking

| Addition | Why Needed | Simpler Alternative Rejected Because |
| --- | --- | --- |
| Short-lived access-grant registry | Binds each bearer to resource, daemon, client, actor, scopes, and expiry | Reusing the durable remote credential would let one bearer cross the JSON and MCP boundaries and would not provide resource-bound expiry |
| Actor-bound in-process API client | Preserves existing API behavior and original actor audit identity | Creating runtime MCP actors would introduce a second identity; direct store calls would bypass validation and audit |
| Explicit resource URL configuration | Supports reverse proxies and RFC 8707 identity correctly | Deriving from the bind address would advertise private or wildcard listener addresses as the public resource |
