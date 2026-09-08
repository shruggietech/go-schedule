# Tasks: Local Observe-Only MCP

**Input**: Design documents from `/specs/069-local-mcp-observe/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/observe-resources.md`

**Tests**: Required by FR-021 and SC-006. Tests are authored before or with the implementation they constrain and must pass before a task is complete.

## Phase 1: Setup and contracts

- [x] T001 Confirm issues #161 and #162 traceability and complete specification, clarification, plan, checklist, contract, and task artifacts in `specs/069-local-mcp-observe/`
- [x] T002 Add the official MCP Go SDK at the researched stable version and inspect every new direct and transitive dependency license for compatibility
- [x] T003 [P] Add dedicated Observe envelope, task, schedule, run, alert, page, trust, and safe-error models in `internal/mcpobserve/model.go`
- [x] T004 [P] Add UTF-8-safe text bounds and versioned opaque cursor helpers with failure-first tests in `internal/mcpobserve/bounds.go`, `internal/mcpobserve/cursor.go`, and matching tests

## Phase 2: Foundational safe adapter

- [x] T005 Define a narrow read-client interface and adapt the existing daemon IPC client without direct persistence or fallback transport in `internal/mcpobserve/server.go`
- [x] T006 Add allowlisted mapping, deterministic ordering, page construction, trust labels, output truncation, and bounded errors in `internal/mcpobserve/resources.go`
- [x] T007 Add failure-first tests with canaries for commands, arguments, environment, stdin, paths, run-as identities, trigger keys, notification credentials, raw schedules, hostile text, malformed cursors, denied access, unavailable daemon, timeout, and cancellation

## Phase 3: User Story 1 - Connect a local MCP host safely (Priority: P1)

**Goal**: Deliver an explicit local stdio endpoint with official SDK negotiation and no network or mutation capability.

**Independent Test**: Launch the command as a hidden subprocess, initialize, discover resources and templates, read health, close stdin, and verify protocol-only stdout plus prompt termination.

- [x] T008 [US1] Register the five static resources and four continuation templates with the official SDK while registering zero tools in `internal/mcpobserve/server.go`
- [x] T009 [US1] Add `gosched mcp serve` with context propagation, stderr diagnostics, and no ordinary stdout output in `internal/cli/mcp.go` and `internal/cli/cli.go`
- [x] T010 [US1] Add SDK protocol tests for both supported revisions, discovery, zero tools, protocol framing, cancellation, and clean shutdown
- [x] T011 [US1] Add cross-platform subprocess coverage with `CREATE_NO_WINDOW` on Windows for stdio framing, stdout cleanliness, host disconnect, exit status, and structural absence of a network listener

## Phase 4: User Story 2 - Inspect scheduler state safely (Priority: P1)

**Goal**: Make all approved observations useful, bounded, paginated, and structurally free of protected execution data.

**Independent Test**: Read every resource against hostile secret-bearing fixtures and follow every continuation page without duplication.

- [x] T012 [US2] Implement daemon health and active task resources with dedicated safe response types
- [x] T013 [US2] Implement upcoming schedule projections without raw recurrence or execution configuration
- [x] T014 [US2] Implement recent alert and run resources with bounded untrusted messages and output excerpts
- [x] T015 [US2] Add end-to-end resource tests for empty, exact-limit, multi-page, changing, long, invalid UTF-8, active-run, and source-truncated fixtures

## Phase 5: User Story 3 - Preserve local authorization (Priority: P2)

**Goal**: Keep existing operating-system IPC access authoritative and expose only bounded recovery guidance.

**Independent Test**: Exercise denied, absent, timed-out, malformed, and canceled daemon calls and prove no fallback, leaked endpoint, or state change.

- [x] T016 [US3] Complete bounded public error mapping and safe stderr diagnostics without endpoint or secret disclosure
- [x] T017 [US3] Prove every Observe read uses only the injected IPC client, respects cancellation and deadlines, and performs no mutation through failure-first tests

## Phase 6: Documentation and verification

- [x] T018 Document Codex clean-install configuration, resource inventory, trust boundary, bounds, daemon prerequisite, and disabled future permission classes in `docs/mcp.md` and link it from the documentation index
- [x] T019 Update Unreleased change notes and architecture rationale in `CHANGELOG.md`
- [x] T020 Run focused SDK, adapter, CLI, subprocess, race, and dependency-license verification from `quickstart.md`
- [x] T021 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T022 Run canonical `scripts/verify.sh all` in the foreground and record objective evidence in `specs/069-local-mcp-observe/verification.md`
- [x] T023 Scan changed files for mojibake, secret canaries, Unicode em dashes, hard-wrapped Markdown prose, and publication formatting defects
- [x] T024 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 establishes protocol and safe-data contracts.
- Phase 2 blocks all resource handlers and proves the secret boundary first.
- User Story 1 establishes host lifecycle and protocol discovery.
- User Story 2 fills the discovered surface and is independently testable through the injected client.
- User Story 3 hardens errors and authorization behavior across the completed surface.
- Phase 6 follows all required stories.

## Implementation Strategy

Build failure-first bound, cursor, redaction, and protocol tests, then implement the safe adapter and explicit CLI entry point. Complete each resource against hostile fixtures, document Codex setup and future permission gates, run the analysis pass, and use the canonical repository gate as the final definition of green.
