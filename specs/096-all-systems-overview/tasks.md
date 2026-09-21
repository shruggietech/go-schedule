# Tasks: All Systems Operational Overview

**Input**: Design documents from `/specs/096-all-systems-overview/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/system-summary.md

## Phase 1: Shared contract and bounded daemon summary

**Purpose**: Establish the safe additive contract used by local and remote observations.

- [x] T001 [P] Define safe system-summary domain models in `internal/domain/system_summary.go`.
- [x] T002 [P] Add deterministic bounded store-summary tests in `internal/store/system_summary_test.go`.
- [x] T003 Implement count and representative-row queries in `internal/store/system_summary.go`.
- [x] T004 [P] Add server authorization, method, response, redaction, and error tests in `internal/api/server/system_summary_test.go` and authorization catalog tests.
- [x] T005 Add `GET /v1/system-summary` to `internal/api/server/server.go`, its handler, and `internal/authorization/catalog.go` with Observe authority.
- [x] T006 [P] Add client decoding and compatibility tests in `internal/api/client/system_summary_test.go`.
- [x] T007 Implement the client method in `internal/api/client/system_summary.go`.

**Checkpoint**: One daemon returns one bounded, read-only summary through either transport.

---

## Phase 2: User Story 1 - Understand every daemon at a glance (Priority: P1)

**Goal**: Refresh This computer and every saved profile independently with explicit state and identity.

**Independent Test**: Load local, reachable, unreachable, duplicate-label, and slow registrations and confirm one correctly classified observation per current registration.

- [x] T008 [P] [US1] Add registration, observation, snapshot, and failure models in `desktop/systems/model.go`.
- [x] T009 [P] [US1] Add deterministic service tests for membership, four-worker fan-out, five-second target contexts, partial results, stale fallback, cancellation, and removed profiles in `desktop/systems/service_test.go`.
- [x] T010 [US1] Implement independent target construction, refresh generations, concurrency bounds, typed failure mapping, and session cache in `desktop/systems/service.go`.
- [x] T011 [US1] Wire the systems service through `desktop/main.go` and `desktop/app.go` without changing the selected connection during refresh.
- [x] T012 [P] [US1] Add frontend models and the native bridge adapter in `desktop/frontend/src/connection/model.ts` and `desktop/frontend/src/connection/bridge.ts`.
- [x] T013 [P] [US1] Add generation replacement and current-registration coverage in `desktop/systems/service_test.go` and `desktop/frontend/src/systems/SystemsPage.test.tsx`.
- [x] T014 [US1] Implement generation-safe overview state in `desktop/frontend/src/systems/SystemsPage.tsx`.

**Checkpoint**: Every registration appears once with independent current, stale, or failure state.

---

## Phase 3: User Story 2 - Triage operational attention (Priority: P1)

**Goal**: Compare bounded operational counts, representatives, attention state, filters, and stable sorting.

**Independent Test**: Seed mixed summaries and isolate each attention category through text, state, and attention filters while changing sort order.

- [x] T015 [P] [US2] Add component tests for counts, representatives, filters, same-name identity, stable sorts, stale labeling, and empty state in `desktop/frontend/src/systems/SystemsPage.test.tsx`.
- [x] T016 [US2] Implement the All Systems page in `desktop/frontend/src/systems/SystemsPage.tsx`.
- [x] T017 [US2] Add All Systems to the route model, shell navigation, and page composition in `desktop/frontend/src/connection/model.ts`, `desktop/frontend/src/components/Shell.tsx`, and `desktop/frontend/src/App.tsx`.

**Checkpoint**: Operators can locate every daemon and isolate systems requiring attention.

---

## Phase 4: User Story 3 - Drill into one safe target (Priority: P1)

**Goal**: Select exactly one registration before navigating to relevant source context.

**Independent Test**: Drill into same-named local and remote registrations and confirm navigation occurs only after the exact selection succeeds.

- [x] T018 [P] [US3] Add drill-down tests for exact profile selection, preserved task and record context, and selection failure in `desktop/frontend/src/systems/SystemsPage.test.tsx` and `desktop/frontend/src/App.test.tsx`.
- [x] T019 [US3] Implement registration-key drill-down through the existing connection selection bridge in `desktop/frontend/src/App.tsx`.
- [x] T020 [US3] Surface optional source context in relevant Tasks, Schedule, Activity, and Notifications destinations without inferring target identity from display names.

**Checkpoint**: Every action opens the correct daemon or leaves the overview unchanged with recovery guidance.

---

## Phase 5: User Story 4 - Refresh and navigate accessibly (Priority: P2)

**Goal**: Keep refresh, filter, sort, inspection, and navigation accessible at supported density and viewport constraints.

**Independent Test**: Operate the page by keyboard at 800 by 600 and 200 percent zoom with reduced motion and 100 profiles, then run focused accessibility automation.

- [x] T021 [P] [US4] Add keyboard, live-region, responsive-overflow, and accessible-name coverage in `desktop/frontend/src/systems/SystemsPage.test.tsx` and the desktop Playwright suite.
- [x] T022 [US4] Add compact responsive overview styles, visible focus, status shapes, reduced-motion behavior, and narrow-layout stacking in `desktop/frontend/src/styles.css`.
- [x] T023 [US4] Ensure refresh announcements preserve focus and expose partial, stale, and failed result counts in `desktop/frontend/src/systems/SystemsPage.tsx`.

**Checkpoint**: The overview remains fully operable and understandable across supported input and layout modes.

---

## Phase 6: Documentation and integration quality

- [x] T024 [P] Document the additive endpoint and authority boundary in `docs/api.md` and `docs/remote-access.md`.
- [x] T025 [P] Add the S096 product change and architectural rationale to `CHANGELOG.md` and the specification index.
- [x] T026 Run focused Go and frontend tests, race tests, builds, formatting, and repository CI-parity checks; record exact commands and results in `specs/096-all-systems-overview/verification.md`.
- [x] T027 Run `go run ./scripts/github-format`, inspect the diff for encoding and scope integrity, and update every completed task checkbox.

---

## Dependencies and execution order

- T001 through T007 establish the daemon contract and block desktop orchestration.
- T008 through T014 establish registration refresh and block the presentation stories.
- T015 through T017 deliver comparison and can overlap with T018 test preparation after frontend models exist.
- T018 through T020 depend on exact observation identities and the existing connection selection path.
- T021 through T023 depend on the rendered page but can be refined alongside presentation implementation.
- T024 and T025 can proceed after the contract stabilizes. T026 and T027 are final integration tasks.

## Scope boundaries

- S096 includes issue #182 only and leaves #183, #184, #185, #176, and #177 open.
- No bulk mutation, shared ownership, failover, reconciliation, portable bundle, SMTP, or native notification delivery work belongs in this slice.
- Tests, CI, and review remain engineering gates but are not GitHub issue closure criteria.
