# Tasks: Authenticated Localhost MCP

**Input**: Design documents from `/specs/070-localhost-mcp-http/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/localhost-mcp-http.md`

**Tests**: Required by FR-020 and SC-006. Security and lifecycle tests are failure-first and must pass before their implementation task is complete.

## Phase 1: Specification and contracts

- [x] T001 Confirm issue #163 traceability and complete specification, clarification, checklists, plan, research, data model, contract, quickstart, and task artifacts in `specs/070-localhost-mcp-http/`
- [x] T002 [P] Define API status, enable, and one-time credential response types plus a narrow lifecycle manager interface in `internal/api/server/mcp_http.go`
- [x] T003 [P] Define failure-first origin normalization, bearer parsing, digest verification, Host policy, and request-bound tests in `internal/mcphttp/manager_test.go`

## Phase 2: Foundational runtime manager

- [x] T004 Implement the authenticated handler boundary and official stateless Streamable HTTP adapter in `internal/mcphttp/handler.go`
- [x] T005 Add failure-first manager tests for default-off state, atomic bind, port conflict, repeated enable, rotation, idempotent disable, unexpected serve exit, shutdown, and concurrent lifecycle calls in `internal/mcphttp/manager_test.go`
- [x] T006 Implement the concurrency-safe runtime manager with numeric loopback binding, random credential issuance, digest-only retention, safe status, and bounded shutdown in `internal/mcphttp/manager.go`

## Phase 3: User Story 1 - Control temporary access (Priority: P1)

**Goal**: Let a local operator explicitly enable, inspect, rotate, and disable one runtime-only endpoint through protected IPC.

**Independent Test**: Exercise all four operations against a fresh daemon and prove restart/default-off behavior plus independent local IPC health.

- [x] T007 [US1] Add protected IPC lifecycle routes with validation, conflict, and safe-error mapping in `internal/api/server/mcp_http.go` and `internal/api/server/server.go`
- [x] T008 [US1] Add typed local client methods and path/error contract tests in `internal/api/client/mcp_http.go` and `internal/api/client/mcp_http_test.go`
- [x] T009 [US1] Wire manager construction, API injection, serve-failure logging, credential revocation, and bounded shutdown into `cmd/goschedd/main.go` with lifecycle tests
- [x] T010 [US1] Add `gosched mcp http status|enable|rotate|disable` including JSON and one-time-secret output tests in `internal/cli/mcp.go` and `internal/cli/mcp_test.go`

## Phase 4: User Story 2 - Observe over official HTTP (Priority: P1)

**Goal**: Deliver the complete S069 Observe surface through current stateless Streamable HTTP with no schema or authority drift.

**Independent Test**: Initialize official SDK clients for both supported revisions, compare discovery and resource results across HTTP and stdio-backed server sessions, and prove zero tools.

- [x] T011 [US2] Construct `mcpobserve.NewServer` for Streamable HTTP requests using the existing daemon IPC client and supported build version
- [x] T012 [US2] Add official SDK compatibility, discovery, zero-tool, resource-read, parity, cancellation, oversized-body, and graceful-shutdown integration tests

## Phase 5: User Story 3 - Reject hostile local requests (Priority: P1)

**Goal**: Fail closed against DNS rebinding, hostile browser origins, and missing or stale credentials before MCP parsing.

**Independent Test**: Run complete Host, Origin, authorization, rotation, disablement, and restart matrices against a real loopback listener.

- [x] T013 [US3] Complete exact Host and origin security matrices including aliases, alternate loopback text, opaque and malformed origins, and forwarded headers
- [x] T014 [US3] Complete authorization matrices including duplicate, malformed, whitespace, suffix, stale, rotated, revoked, and post-restart credentials
- [x] T015 [US3] Prove rejected requests never dispatch MCP, never return permissive CORS headers, never echo secrets, and do not alter local IPC, GUI, CLI, or stdio authority

## Phase 6: Documentation and verification

- [x] T016 Document enablement, one-time credential handling, browser origins, client setup, rotation, disablement, restart behavior, and troubleshooting in `docs/mcp.md`
- [x] T017 Update Unreleased change notes and architectural rationale in `CHANGELOG.md`
- [x] T018 Run focused manager, API, CLI, SDK, security, lifecycle, and race verification from `quickstart.md`
- [x] T019 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T020 Run canonical `scripts/verify.sh all` in the foreground and record objective evidence in `specs/070-localhost-mcp-http/verification.md`
- [x] T021 Scan changed files for mojibake, credentials, Unicode em dashes, hard-wrapped Markdown prose, and publication formatting defects
- [x] T022 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 fixes the public and security contract.
- Phase 2 establishes the fail-closed manager and handler before any control route can expose them.
- User Story 1 wires lifecycle authority through protected IPC.
- User Story 2 adds protocol parity over the active listener.
- User Story 3 completes adversarial verification across the working endpoint.
- Phase 6 follows all required stories.

## Implementation Strategy

Write origin, Host, credential, and lifecycle tests first; implement the manager and authenticated handler; then expose the four protected control operations and CLI. Reuse the exact S069 Observe server for HTTP, prove compatibility and parity through the official SDK, complete the adversarial matrices, run spec-kit analysis, and finish with the canonical repository gate.
