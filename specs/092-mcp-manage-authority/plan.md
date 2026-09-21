# Implementation Plan: MCP Manage Authority

**Branch**: `codex/092-mcp-manage-authority` | **Date**: 2026-09-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/092-mcp-manage-authority/spec.md`

## Summary

Add explicit local MCP Manage sessions and six definition-family tools. Each tool accepts one action and one bounded typed definition, validates exact daemon and request identity, optionally requires caller confirmation, deduplicates retries, and delegates to the existing local API client. Results expose identity and outcome only, while shared authorization and intent-first audit middleware remain authoritative.

## Technical Context

**Language/Version**: Go 1.26.0 and TypeScript 5

**Primary Dependencies**: official MCP Go SDK 1.8.0, Cobra, existing local API client/server, React desktop frontend

**Storage**: Existing SQLite definitions and audit records; Manage session policy and deduplication remain memory-only

**Testing**: Go unit, API integration, MCP conformance, desktop frontend, race, coverage, documentation, and automation gates

**Target Platform**: Supported Windows, macOS, and Linux local daemon environments

**Project Type**: Go daemon, CLI, local IPC API, localhost MCP transport, desktop client

**Performance Goals**: One API mutation per tool call; bounded 256-entry retry cache per Manage session

**Constraints**: Observe remains default; no remote MCP, bulk mutation, secret reveal, permission administration, enrollment, direct storage mutation, task environment, or stdin values

**Scale/Scope**: Six tools spanning five definition families plus notification-assignment replacement

## Constitution Check

- Existing API validation, authorization, readiness, event publication, and audit paths are reused.
- Capability ordering remains fail-closed and discovery never exposes Manage tools below Manage.
- Secret-bearing fields are absent from the MCP schema and all results are allowlisted.
- One mutation per request preserves clear atomic and uncertain-outcome semantics.
- Tests cover authority boundaries, confirmation, hostile content, secret redaction, retries, revocation, and transport ambiguity.
- The implementation is bounded to issue #179 and does not pull remote transport or durable grant administration into S092.

All gates pass before implementation. Re-check after design: six family tools are simpler and easier to audit than separate tools for every action, while typed family payloads preserve validation clarity.

## Project Structure

### Documentation

```text
specs/092-mcp-manage-authority/
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
├── cli/                 # explicit Manage and confirmation flags
├── mcphttp/             # listener-scoped Manage executor and policy
├── mcpmanage/           # typed tools, redacted results, validation, deduplication
└── mcpsession/          # permit explicit Manage runtime actors

desktop/
├── agentaccess/         # non-secret Manage projection and policy
└── frontend/            # permission and confirmation controls

docs/mcp.md              # tool, safety, failure, and confirmation guidance
```

**Structure Decision**: Add `mcpmanage` beside `mcpoperate`, reuse the session registry and API client, and keep the transport responsible only for selecting authority and connection policy.

## Complexity Tracking

| Addition | Why Needed | Simpler Alternative Rejected Because |
| --- | --- | --- |
| Six family tools with action discriminators | Keeps discovery bounded while retaining typed domain contracts | One untyped universal mutation tool would make schemas, review, and security boundaries ambiguous; one tool per action would unnecessarily triple discovery surface |
| Optional connection confirmation policy | Supports attended high-risk deployments without breaking unattended automation | Mandatory prompts would make agent automation unusable, while no policy would omit an explicit issue requirement |
