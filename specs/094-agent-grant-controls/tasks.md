# Tasks: Agent Grant Controls

## Phase 1: Specification and contracts

- [x] T001 Define the S094 grant inventory, creation, monotonic lifecycle, audit, accessibility, transport isolation, and epic acceptance requirements in specs/094-agent-grant-controls/spec.md
- [x] T002 Record the persistence, projection, native handoff, edit, revoke, audit, and verification decisions in specs/094-agent-grant-controls/plan.md, specs/094-agent-grant-controls/research.md, specs/094-agent-grant-controls/data-model.md, specs/094-agent-grant-controls/contracts/agent-access.md, and specs/094-agent-grant-controls/quickstart.md

## Phase 2: Foundational pairing expiry

- [x] T003 [P] Add failing domain, store migration, and enrollment tests for fixed finite and deliberate non-expiring grant deadlines in internal/domain/remote_access_test.go, internal/store/migration_v19_test.go, internal/store/remote_access_test.go, and internal/enrollment/service_test.go (FR-005, FR-006)
- [x] T004 Add GrantExpiresAt to pairing contracts, schema migration v19, persistence scans, and atomic exchange propagation in internal/domain/remote_access.go, internal/store/store.go, internal/store/remote_access.go, and internal/enrollment/service.go (FR-005, FR-006)
- [x] T005 Update typed API requests, clients, CLI duration options, and documentation compatibility tests for grant expiry in internal/api/server/pairing.go, internal/api/client/access.go, internal/cli/access.go, and docs/cli.md (FR-005, FR-006)

## Phase 3: User Story 1, understand current agent access

**Goal**: Present a secret-free, transport-aware MCP inventory for this daemon.

**Independent Test**: Seed persistent, localhost runtime, stdio runtime, expired, and revoked MCP actors, then verify every required grant fact and independent transport state.

- [x] T006 [P] [US1] Add failing Go projection tests for daemon identity, MCP filtering, transport correlation, safe timestamps, authority copy, and secret exclusion in desktop/agentaccess/service_test.go (FR-001, FR-002, FR-003, FR-004, FR-014)
- [x] T007 [P] [US1] Add failing frontend tests for MCP Off, grant rows, authority distinctions, expiry, state, and responsive secret-free rendering in desktop/frontend/src/agentaccess/AgentAccessPage.test.tsx (FR-001, FR-002, FR-003, FR-004, FR-014)
- [x] T008 [US1] Extend the Agent Access backend and models to load manifest, actors, credentials, audit summaries, and localhost status into a bounded safe workspace in desktop/agentaccess/local.go, desktop/agentaccess/model.go, and desktop/agentaccess/service.go (FR-001, FR-002, FR-003, FR-004, FR-014)
- [x] T009 [US1] Extend frontend models and the Agent Access page with transport status cards and accessible grant inventory cards in desktop/frontend/src/agentaccess/model.ts and desktop/frontend/src/agentaccess/AgentAccessPage.tsx (FR-001, FR-002, FR-003, FR-004, FR-014)

## Phase 4: User Story 2, grant constrained remote MCP access

**Goal**: Create a named, duration-bounded remote MCP pairing without returning its phrase to React.

**Independent Test**: Create each authority and duration choice, exchange one pairing, and prove native copy, cancellation on copy failure, and secret-free bridge results.

- [x] T010 [P] [US2] Add failing Go tests for duration validation, enrollment bundle copy, phrase clearing, clipboard rollback, and safe results in desktop/agentaccess/service_test.go and desktop/app_test.go (FR-005, FR-006, FR-007, FR-008)
- [x] T011 [P] [US2] Add failing frontend store and page tests for the grant dialog, deliberate non-expiring choice, pending state, outcome announcement, and absence of phrase fields in desktop/frontend/src/agentaccess/store.test.ts and desktop/frontend/src/agentaccess/AgentAccessPage.test.tsx (FR-005, FR-007, FR-008, FR-015)
- [x] T012 [US2] Implement the native create-grant lifecycle and Wails facade in desktop/agentaccess/model.go, desktop/agentaccess/service.go, desktop/agentaccess/local.go, and desktop/app.go (FR-005, FR-006, FR-007, FR-008)
- [x] T013 [US2] Implement the typed bridge, store mutation, and accessible create-grant dialog in desktop/frontend/src/agentaccess/model.ts, desktop/frontend/src/agentaccess/bridge.ts, desktop/frontend/src/agentaccess/store.ts, and desktop/frontend/src/agentaccess/AgentAccessPage.tsx (FR-005, FR-007, FR-008, FR-015)

## Phase 5: User Story 3, narrow, expire, revoke, and audit grants

**Goal**: Apply monotonic grant changes immediately and present bounded shared audit evidence.

**Independent Test**: Narrow, expire, and revoke grants while applicable connections remain open, then inspect a 25-event secret-free audit view and complete every interaction by keyboard.

- [x] T014 [P] [US3] Add failing Go tests for monotonic edit validation, revocation, bounded actor audit, and next-request authorization across local and remote MCP in desktop/agentaccess/service_test.go, internal/api/server/access_test.go, internal/mcphttp/manager_test.go, and internal/remotemcp/handler_test.go (FR-009, FR-010, FR-011, FR-012, FR-013, FR-016)
- [x] T015 [P] [US3] Add failing frontend and accessibility tests for edit options, expiry constraints, revoke confirmation, audit disclosure, focus return, announcements, and non-color-only state in desktop/frontend/src/agentaccess/store.test.ts, desktop/frontend/src/agentaccess/AgentAccessPage.test.tsx, and desktop/frontend/src/accessibility.test.tsx (FR-009, FR-010, FR-012, FR-013, FR-015)
- [x] T016 [US3] Implement server-side monotonic actor change validation and typed grant edit, revoke, and audit service operations in internal/store/access.go, internal/api/server/access.go, desktop/agentaccess/local.go, desktop/agentaccess/model.go, and desktop/agentaccess/service.go (FR-009, FR-010, FR-011, FR-012, FR-013, FR-016)
- [x] T017 [US3] Implement Wails, bridge, store, edit and revoke dialogs, and recent-action disclosure in desktop/app.go, desktop/frontend/src/agentaccess/model.ts, desktop/frontend/src/agentaccess/bridge.ts, desktop/frontend/src/agentaccess/store.ts, and desktop/frontend/src/agentaccess/AgentAccessPage.tsx (FR-009, FR-010, FR-012, FR-013, FR-015)

## Phase 6: Polish and cross-cutting verification

- [x] T018 Add compact grant, authority, transport, state, and audit styling for supported widths, themes, focus, and reduced motion in desktop/frontend/src/styles.css (FR-004, FR-015)
- [x] T019 Add Playwright coverage for the complete keyboard grant lifecycle and independent transport states in desktop/frontend/e2e/agent-access.spec.ts (FR-001, FR-005, FR-009, FR-012, FR-014, FR-015)
- [x] T020 Update administrator guidance, security boundaries, changelog, and specification inventory in docs/mcp.md, docs/remote-access.md, docs/access-control.md, CHANGELOG.md, and specs/README.md (FR-003, FR-004, FR-014, FR-017)
- [x] T021 Run focused Go, frontend, accessibility, Playwright, race, MCP conformance, hostile-content, redaction, formatting, and full scripts/verify.sh all gates; record issue #181 and epic #148 evidence in specs/094-agent-grant-controls/verification.md (SC-001 through SC-006)
- [x] T022 Mark every completed task and transition the S094 specification to Implemented in specs/094-agent-grant-controls/tasks.md and specs/094-agent-grant-controls/spec.md

## Dependencies

- Phase 2 blocks all user stories because exchanged grants need fixed expiry semantics.
- User Story 1 establishes the shared workspace projection used by User Stories 2 and 3.
- User Story 2 and User Story 3 can be reviewed independently after User Story 1, but their shared service and page files execute sequentially.
- Phase 6 begins after all three user stories pass their independent tests.

## Parallel Opportunities

- T003 spans independent test files and can be developed in parallel before T004.
- T006 and T007 are independent backend and frontend failing-test tasks.
- T010 and T011 are independent backend and frontend failing-test tasks.
- T014 and T015 are independent backend and frontend failing-test tasks.

## Implementation Strategy

Deliver the persistent expiry foundation first, then the read-only inventory, then grant creation, then monotonic lifecycle and audit controls. Each story is independently demonstrable, while the final slice is accepted only after issue #181 and parent #148 pass together.
