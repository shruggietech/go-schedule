# Implementation Plan: Authenticated Localhost MCP

**Branch**: `codex/070-localhost-mcp-http` | **Date**: 2026-09-08 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/070-localhost-mcp-http/spec.md`

## Summary

Add a concurrency-safe `internal/mcphttp` runtime manager owned by `goschedd`, controlled only through new protected IPC API methods and CLI commands. The manager binds `127.0.0.1` on an operator-selected port, keeps its credential and policy only in memory, validates exact Host, optional exact Origin, and bearer authorization before delegating to the official SDK's stateless Streamable HTTP handler, and constructs the same `mcpobserve.NewServer` used by stdio. Disablement, rotation, and daemon shutdown revoke access without changing any other local or future remote access surface.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, existing local daemon API client/server, existing official `github.com/modelcontextprotocol/go-sdk/mcp` v1.7.0

**Storage**: Runtime memory only for endpoint state and credential digest; existing daemon store remains read-only through the S069 Observe client

**Testing**: Go unit and integration tests, official SDK Streamable HTTP client tests, transport parity tests, race tests, security matrices, documentation checks, and canonical `scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux daemon installations with local TCP networking

**Project Type**: Go daemon, local IPC API, and CLI

**Performance Goals**: Start or reject a listener within one second; reject invalid requests before MCP parsing; shut down within five seconds; retain S069's ten-second Observe deadline and bounded result sizes

**Constraints**: Default off; IPv4 numeric loopback only; one listener per daemon; no persistence; no plaintext credential retention; no CORS wildcard; no tools or mutation; 1 MiB request body; five-second graceful shutdown; no configuration, GUI, remote JSON, or future remote MCP coupling

**Scale/Scope**: One runtime manager, four local IPC operations, four CLI commands, one SDK HTTP transport, existing five-resource Observe surface, security and lifecycle matrices, and focused documentation

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: Listener state, request policy, secret handling, and SDK adaptation are isolated in `internal/mcphttp`; daemon API and CLI expose narrow lifecycle operations; scheduler core remains unchanged.
- **Testing**: Failure-first tests cover every lifecycle transition, security matrix, protocol revision, transport parity case, port conflict, cancellation, request bound, and concurrency behavior.
- **UX consistency**: CLI verbs mirror other administration surfaces, status is non-secret, newly issued credentials are clearly one-time values, and errors distinguish validation, conflict, and runtime failures.
- **Performance**: Bind is synchronous and atomic, rejected requests do no MCP work, Streamable HTTP is stateless, request bodies are capped, and shutdown has a five-second ceiling.
- **Security and truth**: Numeric loopback, exact Host and Origin policy, digest-only credential retention, constant-time digest comparison, no persistence, no CORS grant, and exact S069 Observe reuse make the boundary fail closed.
- **Dependency governance**: No new dependency is added; the already approved official MCP SDK provides the current Streamable HTTP transport and retains its own localhost protection.
- **Review workflow**: Work remains on `codex/070-localhost-mcp-http`; the user explicitly authorized push, PR publication, verified review-fix pushes, and one manually triggered second Codex review round.
- **Pinned artifacts**: No pinned workflow, toolchain, packaging, or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/070-localhost-mcp-http/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── localhost-mcp-http.md
├── checklists/
│   ├── localhost-security.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
cmd/goschedd/
├── main.go
└── main_test.go

internal/
├── api/
│   ├── client/
│   │   ├── mcp_http.go
│   │   └── mcp_http_test.go
│   └── server/
│       ├── mcp_http.go
│       ├── mcp_http_test.go
│       └── server.go
├── cli/
│   ├── mcp.go
│   └── mcp_test.go
├── mcphttp/
│   ├── handler.go
│   ├── manager.go
│   └── manager_test.go
└── mcpobserve/
    └── server.go

docs/
└── mcp.md
```

**Structure Decision**: Keep the transport and runtime secret boundary in a new `internal/mcphttp` package. Define the lifecycle interface and wire contract in the daemon API package so the API does not import its concrete manager. Let the manager depend on the existing IPC client and `mcpobserve.NewServer`, preserving exact Observe parity without a direct store path. Construct and shut down the manager in `cmd/goschedd`.

## Complexity Tracking

No constitutional violations require justification.
