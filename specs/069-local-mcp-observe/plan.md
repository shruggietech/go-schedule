# Implementation Plan: Local Observe-Only MCP

**Branch**: `codex/069-local-mcp-observe` | **Date**: 2026-09-08 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/069-local-mcp-observe/spec.md`

## Summary

Add an `internal/mcpobserve` adapter that registers five bounded JSON resources and their continuation templates with the official MCP Go SDK, maps existing daemon reads into dedicated allowlisted types, and runs only when `gosched mcp serve` is launched. The adapter communicates with the daemon through the existing Unix-socket or Windows-named-pipe client, exposes no tools, treats user-controlled content as untrusted data, and keeps protocol stdout isolated from diagnostics.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, existing local daemon client, official `github.com/modelcontextprotocol/go-sdk/mcp` v1.7.0

**Storage**: Existing daemon persistence through read-only local IPC methods; no direct database access and no migration

**Testing**: Go unit, MCP in-memory protocol, stdio subprocess, race, dependency-license inspection, documentation, and canonical `scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux CLI installations used by local MCP hosts including Codex

**Project Type**: Go daemon and CLI with a local protocol adapter

**Performance Goals**: Discover resources immediately; complete an individual Observe read within 10 seconds; cap collection pages at 100 records, text at 2 KiB, and output excerpts at 8 KiB

**Constraints**: No TCP listener; no tools or mutation; no secret-bearing daemon type may be serialized; stdout remains protocol-only; cancellation propagates; Windows subprocess tests use `CREATE_NO_WINDOW`; no MCP logic enters scheduler core

**Scale/Scope**: Five static first-page resources, four continuation templates, one CLI subcommand, one adapter package, one official SDK dependency, Codex setup documentation, and focused plus canonical verification

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: Protocol registration, cursor validation, safe mapping, bounds, and errors are isolated in one adapter package. Existing domain and scheduler logic remain unchanged.
- **Testing**: Failure-first tests cover every resource, prohibited field, limit, cursor, protocol revision, cancellation, denial, and subprocess lifecycle before completion.
- **UX consistency**: Resources use stable names, JSON media types, UTC timestamps, a shared envelope, and actionable bounded errors.
- **Performance**: Existing aggregate reads are reused, pages are capped at 100, each request has a 10-second deadline, and no background polling is introduced.
- **Security and truth**: Dedicated allowlisted types make redaction structural. Existing IPC access is authoritative. User-controlled strings are bounded and labeled untrusted. No tool is registered.
- **Dependency governance**: The official SDK is the smallest compliant implementation path, is pinned through Go modules, and its new direct and transitive licenses are inspected before review.
- **Review workflow**: Work remains on `codex/069-local-mcp-observe`; the user explicitly authorized push, PR publication, review fixes, and at most one manually triggered second Codex review round.
- **Pinned artifacts**: No pinned workflow, toolchain, packaging, or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/069-local-mcp-observe/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── observe-resources.md
├── checklists/
│   ├── mcp-trust.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/
├── cli/
│   ├── cli.go
│   ├── mcp.go
│   └── mcp_test.go
└── mcpobserve/
    ├── bounds.go
    ├── cursor.go
    ├── model.go
    ├── resources.go
    ├── server.go
    └── *_test.go

docs/
└── mcp.md
```

**Structure Decision**: Place MCP-specific contracts and handlers in `internal/mcpobserve`, with a narrow read-client interface adapting the existing daemon client. Add only a Cobra entry point in `internal/cli`. This creates a replaceable protocol edge and prevents MCP concepts from contaminating domain, store, executor, or scheduler packages.

## Complexity Tracking

No constitutional violations require justification.
