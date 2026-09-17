# Implementation Plan: MCP Operate Authority

**Branch**: `codex/091-mcp-operate-authority` | **Date**: 2026-09-17 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/091-mcp-operate-authority/spec.md`

## Summary

Add explicit local MCP Operate sessions and exactly three task tools. A memory-only session registry creates an MCP actor, resolves a secret on each local API request, and revokes that actor when the connection ends. Tool handlers verify daemon and task identity, deduplicate retries, use the shared task client, and rely on the existing authorization and audit middleware for enforcement and evidence.

## Technical Context

**Language/Version**: Go 1.26.0

**Primary Dependencies**: official MCP Go SDK 1.8.0, Cobra, existing local API client/server, SQLite audit store

**Storage**: Existing SQLite actor and audit records; session credentials and deduplication records remain in memory

**Testing**: Go unit, contract, conformance, integration, race, coverage, documentation, and automation gates

**Target Platform**: Supported Windows, macOS, and Linux local daemon environments

**Project Type**: Go daemon, CLI, local IPC API, optional localhost MCP transport

**Performance Goals**: Tool dispatch adds bounded identity validation and lookup; retry cache is capped and expires entries

**Constraints**: Observe remains default; no remote MCP, Manage tools, persistent grants, secret disclosure, or direct store mutation from tool handlers

**Scale/Scope**: One runtime session per stdio process or localhost listener, three task tools, bounded deduplication state

## Constitution Check

- Existing local API validation, authorization, and audit paths are reused.
- Authority remains fail-closed and capability based.
- Runtime secrets are shown once, stored only as digests, and never logged.
- Tool results are allowlisted and exclude execution inputs.
- The implementation is bounded to issue #178 and does not pull durable grant administration from #181.
- Verification covers adversarial authorization, retry, revocation, transport uncertainty, and hostile content.

All gates pass before implementation. Re-check after design: the session boundary adds only the minimum identity bridge required to avoid falsely auditing MCP mutations as the local operating-system actor.

## Project Structure

### Documentation

```text
specs/091-mcp-operate-authority/
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
├── api/client/          # runtime session client and task calls
├── api/server/          # local session lifecycle and actor resolution
├── authorization/       # Operate mapping for task toggles
├── cli/                 # explicit stdio and HTTP permission selection
├── mcphttp/             # listener-scoped authority and revocation
├── mcpobserve/          # unchanged Observe resources
├── mcpoperate/          # typed tools, target validation, result contract, deduplication
└── mcpsession/          # memory-only token registry and actor lifecycle

desktop/agentaccess/     # honest authority projection and permission selection
docs/mcp.md              # setup, boundaries, outcomes, and revocation guidance
```

**Structure Decision**: Keep Observe resources intact, isolate mutation tools in `mcpoperate`, and isolate ephemeral identity from both transport and tool logic in `mcpsession`.

## Complexity Tracking

| Addition | Why Needed | Simpler Alternative Rejected Because |
| --- | --- | --- |
| Runtime MCP session registry | Tool mutations require real MCP actor attribution and live revocation | Reusing the local actor would produce false audit attribution and could not revoke an existing MCP connection independently |
| Bounded request deduplicator | Run-now retries can create duplicate executions | Blind retry or relying on HTTP semantics cannot make a non-idempotent scheduler request safe |
