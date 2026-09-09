# Tasks: v1.3 Notifications and Local Agent Access Qualification

**Input**: Design documents from `specs/072-v13-release-qualification/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/qualification-gate.md`, `quickstart.md`

**Tests**: Release qualification is test-first per constitution principle II. Existing detailed suites remain authoritative and the new package-shaped journey composes them into explicit three-platform evidence.

**Organization**: Tasks are grouped by independently testable user story and execute in chronological phase order.

## Phase 1: Setup and Traceability

**Purpose**: Establish the approved qualification boundary before implementation.

- [x] T001 Validate the S072 specification and requirements-quality checklists in `specs/072-v13-release-qualification/`
- [x] T002 [P] Record the active S072 plan in `.specify/feature.json`, `CLAUDE.md`, and `specs/README.md`
- [x] T003 [P] Define the qualification evidence model and three-platform gate contract in `specs/072-v13-release-qualification/data-model.md` and `specs/072-v13-release-qualification/contracts/qualification-gate.md`

---

## Phase 2: Foundational Qualification Contract

**Purpose**: Protect the named hosted gate before relying on it as release evidence.

- [x] T004 Write failing repository automation assertions for the named job, three operating systems, and package-shaped invocation in `scripts/automation-check.sh`
- [x] T005 Add the fail-fast-disabled v1.3 qualification matrix and focused race commands in `.github/workflows/ci.yml`
- [x] T006 Run the automation gate and confirm the pinned workflow contract passes through `scripts/verify.sh`

**Checkpoint**: Every supported platform has a visible, non-optional qualification job whose wiring is protected against drift.

---

## Phase 3: User Story 1 - Preserve Safe Installation Defaults (Priority: P1)

**Goal**: Prove a built daemon remains local, offline, and opt-in on first start and restart.

**Independent Test**: Build and launch `goschedd` twice against isolated retained state, then inspect health, notification state, and MCP HTTP status through local IPC.

- [x] T007 [US1] Write the package-shaped fresh and retained default-state contract test in `test/integration/v13_release_qualification_test.go`
- [x] T008 [US1] Implement hidden cross-platform daemon lifecycle, bounded health readiness, local task creation and listing, and retained-state assertions in `test/integration/v13_release_qualification_test.go`
- [x] T009 [US1] Run the package-shaped integration test under the race detector using `specs/072-v13-release-qualification/quickstart.md`

**Checkpoint**: The candidate starts twice with protected local health access, zero notification records, and no MCP HTTP listener or metadata.

---

## Phase 4: User Story 2 - Qualify Webhook Delivery Without Affecting Task Truth (Priority: P1)

**Goal**: Bind the existing detailed webhook and migration evidence into the named release gate.

**Independent Test**: Run the focused notification, store migration, and end-to-end webhook suites under the race detector and verify source outcomes and secrets remain isolated.

- [x] T010 [US2] Map success, retry, restart, disabled-channel, redaction, migration, and source-outcome tests to issue #190 in `specs/072-v13-release-qualification/verification.md`
- [x] T011 [US2] Run the focused race-checked notification qualification commands from `specs/072-v13-release-qualification/quickstart.md`

**Checkpoint**: Existing webhook behavior is explicitly qualified without changing scheduling truth or default network behavior.

---

## Phase 5: User Story 3 - Qualify Local Observe Access and Truthful Release Claims (Priority: P1)

**Goal**: Bind official-client stdio and authenticated localhost Observe evidence into the same release boundary without implying mutation or remote authority.

**Independent Test**: Run the focused MCP Observe, HTTP, and package-shaped official-SDK suites under the race detector and audit release-facing capability boundaries.

- [x] T012 [US3] Map protocol, resource, template, zero-tool, authorization, lifecycle, redaction, and hostile-content tests to issue #190 in `specs/072-v13-release-qualification/verification.md`
- [x] T013 [US3] Run the focused race-checked MCP qualification commands, update the shipped-versus-future boundary in `docs/notifications.md`, and audit capability claims from `specs/072-v13-release-qualification/quickstart.md`

**Checkpoint**: Both local transports retain the same bounded Observe authority and future remote or mutation work remains unclaimed.

---

## Phase 6: Polish and Publication Readiness

**Purpose**: Close traceability, decisions, verification, and repository-wide quality gates.

- [x] T014 Update the Unreleased qualification entry and dated pinned-workflow decision in `CHANGELOG.md`
- [x] T015 Reconcile issue #190 acceptance evidence and S072 lifecycle state in `specs/072-v13-release-qualification/spec.md`, `specs/072-v13-release-qualification/verification.md`, and `specs/README.md`
- [x] T016 Run `/speckit-analyze`, resolve every actionable finding, and record the disposition in `specs/072-v13-release-qualification/verification.md`
- [x] T017 Run the canonical foreground `scripts/verify.sh all` gate and record all eight results in `specs/072-v13-release-qualification/verification.md`
- [x] T018 Audit UTF-8 without BOM, mojibake, GitHub publication formatting, diff integrity, branch scope, issue traceability, completed task markers, and the explicit no-release boundary

---

## Dependencies and Execution Order

- Phase 1 establishes the slice boundary and completes before analysis or implementation.
- Phase 2 protects the hosted qualification contract and blocks release-evidence claims.
- User Story 1 depends on Phase 2 because its package-shaped invocation is the matrix's primary platform evidence.
- User Stories 2 and 3 reuse existing independently testable suites and can be validated after Phase 2, but their final evidence depends on User Story 1 completing the matrix contract.
- Phase 6 follows all three user stories.

## Parallel Opportunities

- T002 and T003 affect independent context and contract files.
- After Phase 2, the detailed evidence inventory for T010 and T012 can be assembled independently.
- The focused notification and MCP commands in T011 and T013 exercise separate packages.

## Implementation Strategy

1. Protect the named three-platform workflow contract first.
2. Add the smallest package-shaped journey that proves fresh and retained opt-in defaults end to end.
3. Reuse and map the detailed S067 through S071 suites instead of duplicating them.
4. Run focused evidence, analyze the artifacts, then run the canonical eight-gate verification.
5. Publish only the review branch and pull request; leave tag and release creation to the maintainer's later release ritual.

## Format Validation

All tasks use a checkbox, sequential ID, optional parallel marker, required user-story label within story phases, concrete action, and exact file path.
