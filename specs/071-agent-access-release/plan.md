# Implementation Plan: Agent Access Controls and MCP Release Gates

**Branch**: `codex/071-agent-access-controls` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/071-agent-access-release/spec.md`

## Summary

Complete issue #164 by extending the existing runtime-only localhost MCP status with one bounded client label and aggregate successful-access evidence, adding a desktop Agent Access service and workspace that performs fail-closed clipboard handoff for enable and rotation, publishing precise Codex and generic host guidance, and consolidating official-SDK conformance plus real built-command smoke coverage. The slice reuses S069 resource registration and S070 listener authority without persistence, multiple grants, remote access, or mutation tools.

## Technical Context

**Language/Version**: Go 1.25.0; TypeScript 5.9 and React 19 in the Wails frontend

**Primary Dependencies**: Go standard library, existing `github.com/modelcontextprotocol/go-sdk` v1.7.0, Wails v2.14.0, Cobra, Vitest, Testing Library, Playwright, axe-core

**Storage**: No new durable storage. Client name, request count, last access time, endpoint, origins, fingerprint, and digest remain daemon-memory state.

**Testing**: Go `testing` with race detector, official MCP SDK transports, Vitest and Testing Library, Playwright plus axe-core, canonical `scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop and CLI installations

**Project Type**: Cross-platform daemon, CLI, local API, and Wails desktop application

**Performance Goals**: Successful HTTP authorization adds one bounded lock-protected counter and timestamp update; status and desktop operations remain constant-space and complete within existing two-second local-operation deadlines

**Constraints**: No credential enters React state or status; no persistence or multiple client grants; no stdio listener state; no tools; no remote transport; no workflow change unless existing platform jobs cannot exercise the built-command smoke; Windows child processes remain hidden

**Scale/Scope**: One active localhost listener, one named runtime client, one saturating 64-bit request count, five resources, four continuation templates, two protocol revisions, and one new desktop route

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **I. Code Quality**: PASS. A dedicated `desktop/agentaccess` boundary keeps the Wails facade small; manager ownership remains explicit; additive API fields carry documented intent; no new dependency is introduced.
- **II. Testing Standards**: PASS. Manager, API, client, desktop service, React, browser accessibility, conformance, hostile-content, and built-command subprocess behaviors receive regression tests before implementation and run under the canonical race gate.
- **III. User Experience Consistency**: PASS. Agent Access uses existing workspace/result patterns, authoritative daemon status, duplicate-safe operations, actionable errors, and one-time native clipboard handoff. CLI and JSON additions are optional and additive.
- **IV. Performance Requirements**: PASS. The added access accounting is O(1), bounded, and outside scheduling dispatch. No scheduler hot path changes, so a new benchmark would not provide useful evidence.
- **V. Autonomous Build-Phase Execution**: PASS. S071 is traceable to #164, uses this review branch, runs the complete spec-kit sequence and analysis gate, updates the Unreleased changelog, and will publish only under the user's explicit authorization in this kickoff.
- **Engineering constraints**: PASS. No new dependency, schema, configuration, public listener, secret log, or platform exception is introduced.

## Project Structure

### Documentation (this feature)

```text
specs/071-agent-access-release/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── agent-access.md
├── checklists/
│   ├── requirements.md
│   └── ux-security-conformance.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/api/server/mcp_http.go
internal/api/client/mcp_http.go
internal/mcphttp/
internal/mcpobserve/
internal/cli/mcp.go

desktop/
├── agentaccess/
│   ├── model.go
│   ├── local.go
│   ├── service.go
│   └── service_test.go
├── app.go
├── app_test.go
├── main.go
└── frontend/
    ├── src/agentaccess/
    │   ├── model.ts
    │   ├── bridge.ts
    │   ├── store.ts
    │   ├── store.test.ts
    │   ├── AgentAccessPage.tsx
    │   └── AgentAccessPage.test.tsx
    ├── src/App.tsx
    ├── src/components/Shell.tsx
    ├── src/connection/model.ts
    └── e2e/agent-access.spec.ts

test/integration/mcp_packaged_test.go
docs/mcp.md
docs/cli.md
docs/architecture.md
CHANGELOG.md
```

**Structure Decision**: Extend S070's manager and additive local API types at their existing boundaries. Add a peer desktop domain package rather than placing daemon lifecycle calls in Settings or React. Keep conformance in `internal/mcpobserve` and package-shaped subprocess proof in the existing root integration suite so every hosted operating-system race job exercises the real command.

## Design Decisions

### Decision 1: Name the single runtime client instead of introducing grants

S070 has one listener and one credential. S071 attaches one validated display name to that credential generation and exposes it in non-secret status. A durable client registry would pull #181's grants, expiry, revocation, and audit model into this release and conflict with S070's restart-off contract.

### Decision 2: Treat disable as revoke and rotation as replacement

The Agent Access workspace calls the existing disable operation for revocation because there is no useful listener state without its sole credential. Rotation preserves the named client and endpoint policy while clearing per-credential access evidence.

### Decision 3: Keep plaintext credentials behind the desktop backend

The desktop service consumes enable or rotation results, copies the credential through the existing native clipboard boundary, clears the response reference, and returns only non-secret workspace state. Clipboard failure triggers disablement before failure is reported. Sending the secret into React for a copy button would create unnecessary retention in browser state and test snapshots.

### Decision 4: Record aggregate successful access only

The manager records one UTC timestamp and saturating count after current Host, Origin, and Bearer checks succeed and before SDK dispatch. It records no request URI, body, peer address, failed attempt, or task data. This satisfies recent evidence without creating an audit system or privacy-sensitive history.

### Decision 5: Use existing cross-platform jobs for built-command smoke

The root integration test builds `gosched` into a package-shaped temporary directory, launches it with the official SDK command transport, applies the project hidden-process helper, and checks initialization, discovery, zero tools, and clean shutdown. The existing Windows, macOS, and Linux race matrix already runs `./...`, so no pinned workflow edit is needed.

## Post-Design Constitution Check

PASS. The design remains additive, runtime-only, test-first, cross-platform, constant-space, and dependency-neutral. It does not modify scheduler timing, persisted schemas, release workflows, or any authorization surface beyond safe status metadata and existing lifecycle controls. Complexity Tracking is therefore unnecessary.
