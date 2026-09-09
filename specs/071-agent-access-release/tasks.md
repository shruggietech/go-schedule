# Tasks: Agent Access Controls and MCP Release Gates

**Input**: Design documents from `specs/071-agent-access-release/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/agent-access.md`, `quickstart.md`

**Tests**: Behavioral changes are test-first per constitution principle II; security, lifecycle, subprocess, and accessibility coverage is mandatory.

**Organization**: Tasks are grouped by independently testable user story and execute in chronological phase order.

## Phase 1: Setup and Contract Alignment

**Purpose**: Establish S071 traceability and additive contracts before behavior changes.

- [x] T001 Validate the S071 specification and both requirements-quality checklists in `specs/071-agent-access-release/`
- [x] T002 [P] Record the active S071 plan in `.specify/feature.json`, `CLAUDE.md`, and `specs/README.md`
- [x] T003 [P] Define additive named-client and access-evidence contracts in `specs/071-agent-access-release/contracts/agent-access.md` and `specs/071-agent-access-release/data-model.md`

---

## Phase 2: Foundational Runtime Contract

**Purpose**: Extend the S070 manager and protected API with safe additive runtime metadata used by every desktop story.

- [x] T004 Write failing client-name validation, compatibility-default, access-accounting, saturation, rotation-reset, and lifecycle-clear tests in `internal/mcphttp/manager_test.go`
- [x] T005 [P] Write failing additive serialization and API lifecycle tests in `internal/api/server/mcp_http_test.go` and `internal/api/client/mcp_http_test.go`
- [x] T006 Implement optional client name plus non-secret access evidence fields in `internal/api/server/mcp_http.go`
- [x] T007 Implement validated runtime client identity and successful-request accounting in `internal/mcphttp/handler.go` and `internal/mcphttp/manager.go`
- [x] T008 Update optional CLI client-name input and additive output coverage in `internal/cli/mcp.go` and `internal/cli/mcp_test.go`
- [x] T009 Run focused race verification for `internal/mcphttp`, `internal/api/server`, `internal/api/client`, and `internal/cli`

**Checkpoint**: The daemon owns one named runtime client and safe access evidence without changing credential or listener authority.

---

## Phase 3: User Story 1 - Understand and Control Local Agent Access (Priority: P1)

**Goal**: Deliver one truthful, accessible desktop workspace for stdio availability and the complete named localhost client lifecycle.

**Independent Test**: Exercise disconnected, Off, Active, used, rotated, clipboard-failure, and revoked states through the desktop service and React workspace without any plaintext credential entering frontend state.

- [x] T010 [P] [US1] Write failing desktop service tests for workspace projection, validation, timeouts, credential copy, rollback, rotation, revocation, and duplicate suppression in `desktop/agentaccess/service_test.go`
- [x] T011 [P] [US1] Write failing React bridge and store tests for safe result handling, stale loads, and duplicate actions in `desktop/frontend/src/agentaccess/bridge.test.ts` and `desktop/frontend/src/agentaccess/store.test.ts`
- [x] T012 [P] [US1] Write failing Agent Access component tests for Off, Active, evidence, authority, validation, and secret absence in `desktop/frontend/src/agentaccess/AgentAccessPage.test.tsx`
- [x] T013 [US1] Implement the desktop Agent Access model, local API adapter, serialized service, clipboard handoff, and fail-closed rollback in `desktop/agentaccess/model.go`, `desktop/agentaccess/local.go`, and `desktop/agentaccess/service.go`
- [x] T014 [US1] Bind Agent Access workspace and lifecycle methods through `desktop/app.go`, `desktop/main.go`, and `desktop/app_test.go`
- [x] T015 [US1] Implement typed frontend model, native bridge, and stale-safe store in `desktop/frontend/src/agentaccess/model.ts`, `desktop/frontend/src/agentaccess/bridge.ts`, and `desktop/frontend/src/agentaccess/store.ts`
- [x] T016 [US1] Implement the Agent Access page without credential rendering in `desktop/frontend/src/agentaccess/AgentAccessPage.tsx` and `desktop/frontend/src/styles.css`
- [x] T017 [US1] Add Agent Access routing and truthful navigation copy in `desktop/frontend/src/App.tsx`, `desktop/frontend/src/components/Shell.tsx`, and `desktop/frontend/src/connection/model.ts`
- [x] T018 [US1] Run focused desktop Go race tests plus frontend component and build validation for `desktop/agentaccess` and `desktop/frontend/src/agentaccess`

**Checkpoint**: A local user can identify, rotate, and revoke the only active HTTP client while stdio remains accurately described as on demand.

---

## Phase 4: User Story 2 - Configure Supported Hosts Safely (Priority: P1)

**Goal**: Publish one precise guide for Codex, generic stdio, and generic Streamable HTTP without unsupported authority claims.

**Independent Test**: Follow each configuration path from a package-shaped command layout and audit all authority, credential, privacy, and removal language.

- [x] T019 [P] [US2] Expand Codex, generic stdio, generic HTTP, privacy, hostile-content, removal, rotation, and troubleshooting guidance in `docs/mcp.md`
- [x] T020 [P] [US2] Synchronize Agent Access and optional client-name CLI contracts in `docs/cli.md` and `docs/architecture.md`
- [x] T021 [US2] Add the fixed Agent Access guide link to the trusted desktop backend contract in `desktop/agentaccess/service.go` and its tests in `desktop/agentaccess/service_test.go`
- [x] T022 [US2] Run documentation policy, link, formatting, and unsupported-claim audits through `scripts/docs-check.sh` and `scripts/github-format`

**Checkpoint**: Supported local host setup and removal are reproducible and the shipped authority boundary is explicit.

---

## Phase 5: User Story 3 - Qualify the Complete Local MCP Contract (Priority: P2)

**Goal**: Consolidate deterministic MCP conformance, hostile-content, package-shaped execution, and desktop accessibility evidence.

**Independent Test**: Run official clients for both revisions, every resource and template, schema and hostile-content invariants, the built CLI subprocess, and browser accessibility at every supported zoom.

- [x] T023 [P] [US3] Write the complete schema, safe-error, trust, redaction, pagination, revision, resource, template, and zero-tool conformance matrix in `internal/mcpobserve/conformance_test.go`
- [x] T024 [P] [US3] Write hostile-content invariance coverage for discovery, instructions, permission, and tool authority in `internal/mcpobserve/conformance_test.go`
- [x] T025 [P] [US3] Add a real built-command official-SDK smoke with hidden Windows launch in `test/integration/mcp_packaged_test.go`
- [x] T026 [P] [US3] Add keyboard, WCAG 2.2 AA, and 80 through 200 percent zoom coverage in `desktop/frontend/e2e/agent-access.spec.ts`
- [x] T027 [US3] Run focused MCP race, package-shaped subprocess, browser accessibility, and production frontend validation

**Checkpoint**: The complete local Observe surface and desktop control path have deterministic release evidence on supported platforms.

---

## Phase 6: Polish and Publication Readiness

**Purpose**: Close traceability, documentation, verification, and repository-wide quality gates.

- [x] T028 Update the Unreleased feature and dated architecture decision records in `CHANGELOG.md`
- [x] T029 Reconcile issue #164 acceptance evidence and S071 status in `specs/071-agent-access-release/spec.md`, `specs/071-agent-access-release/verification.md`, and `specs/README.md`
- [x] T030 Run `/speckit-analyze`, resolve every actionable finding, and record the disposition in `specs/071-agent-access-release/verification.md`
- [x] T031 Run the canonical foreground `scripts/verify.sh all` gate and record all eight results in `specs/071-agent-access-release/verification.md`
- [x] T032 Audit UTF-8 without BOM, mojibake, secret canaries, GitHub publication formatting, diff integrity, branch scope, and completed task markers

---

## Dependencies and Execution Order

- Phase 1 completes before analysis and implementation.
- Phase 2 blocks every desktop lifecycle action because the workspace projects its additive status contract.
- User Story 1 depends on Phase 2 and is the interactive MVP.
- User Story 2 can begin after the Agent Access model names its user-facing terms; its documentation is independently reviewable.
- User Story 3 depends only on the stable runtime contract, then validates User Stories 1 and 2 as a release increment.
- Phase 6 follows all three user stories.

## Parallel Opportunities

- T002 and T003 affect independent context and contract files.
- T004 and T005 establish separate manager and API failure-first suites.
- T010, T011, and T012 cover independent Go service, frontend state, and presentation contracts.
- T019 and T020 update separate documentation surfaces.
- T023 through T026 cover independent conformance, subprocess, and browser surfaces.

## Implementation Strategy

1. Complete and race-check additive runtime metadata before desktop work.
2. Deliver the desktop workspace end to end, including credential rollback, before documentation claims it.
3. Publish setup guidance and then validate it through the real built-command smoke.
4. Consolidate conformance and accessibility evidence, run analysis, then run the canonical repository gate.

## Format Validation

All tasks use checkbox, sequential ID, optional parallel marker, required user-story label within story phases, concrete action, and exact file path.
