# Tasks: Dependency Consolidation

**Input**: Design documents from `specs/073-dependency-consolidation/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/dependency-baseline.md`, `quickstart.md`

**Tests**: Dependency consolidation is test-first. Existing product tests remain authoritative, and runtime-baseline automation assertions must fail before workflow implementation.

**Organization**: Tasks are grouped by independently testable maintainer and user stories and execute in chronological phase order.

## Phase 1: Setup and Traceability

**Purpose**: Establish the complete source-update inventory and executable specification.

- [x] T001 Validate the S073 specification and requirements-quality checklists in `specs/073-dependency-consolidation/`
- [x] T002 [P] Record the active S073 plan in `.specify/feature.json`, `CLAUDE.md`, and `specs/README.md`
- [x] T003 Run `/speckit-analyze`, resolve every actionable cross-artifact finding, and record the disposition in `specs/073-dependency-consolidation/verification.md`

---

## Phase 2: Foundational Compatibility Contracts

**Purpose**: Protect the runtime and dependency-graph decisions before changing manifests.

- [x] T004 Add failing Node 26 CI and release baseline assertions and negative fixtures in `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`
- [x] T005 Record the selected versions, companion updates, graph relationships, and clean-restoration contract in `specs/073-dependency-consolidation/data-model.md` and `specs/073-dependency-consolidation/contracts/dependency-baseline.md`

**Checkpoint**: The repository rejects a Node 24 workflow baseline and the complete selected dependency set is traceable.

---

## Phase 3: User Story 1 - Maintain One Coherent Dependency Baseline (Priority: P1)

**Goal**: Resolve all eight source updates into clean root, desktop, and frontend dependency graphs.

**Independent Test**: Both Go modules tidy and verify without residual changes, while the frontend completes a clean install with no peer or engine bypass and reports every selected direct version.

- [x] T006 [US1] Apply fsnotify 1.10.1 and modernc.org/sqlite 1.58.0, then regenerate `go.mod` and `go.sum`
- [x] T007 [US1] Apply Wails 2.15.0, reconcile the desktop graph, and update current build pins and guidance in `desktop/go.mod`, `desktop/go.sum`, `scripts/verify.sh`, `desktop/README.md`, `internal/cli/gui.go`, and desktop packaging contract tests
- [x] T008 [US1] Apply the React, testing, jsdom, Vite plugin, Node type, Node runtime, and Vite companion baseline in `desktop/frontend/package.json`
- [x] T009 [US1] Regenerate and clean-install `desktop/frontend/package-lock.json` without force, legacy-peer, or engine bypass
- [x] T010 [US1] Run both Go module tidy and verification commands plus the direct-version inventory from `specs/073-dependency-consolidation/quickstart.md`

**Checkpoint**: All three dependency graphs are internally consistent and reproducible.

---

## Phase 4: User Story 2 - Preserve Supported Product Behavior (Priority: P1)

**Goal**: Preserve all shipped scheduling, storage, watcher, desktop, accessibility, and packaging behavior under the consolidated baseline.

**Independent Test**: Focused storage, watcher, integration, desktop Go, frontend unit, type-check, production bundle, browser, native, and installer evidence passes without weakened assertions.

- [x] T011 [US2] Update the Node 26 runtime and Wails 2.15.0 build baselines in `.github/workflows/ci.yml` and `.github/workflows/release.yml`
- [x] T012 [US2] Update and pass the repository automation contract in `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`
- [x] T013 [US2] Run focused root race, desktop race, frontend unit, type-check, production bundle, browser, native Wails, and Windows installer checks from `specs/073-dependency-consolidation/quickstart.md`
- [x] T014 [US2] Resolve only compatibility failures caused by the selected dependency baseline in the affected source, test, generated, build, or packaging files

**Checkpoint**: Existing supported behavior passes the complete affected compatibility surface.

---

## Phase 5: User Story 3 - Retain Review and Supersession Traceability (Priority: P2)

**Goal**: Make every version, compatibility decision, source pull request, and final-head result independently auditable.

**Independent Test**: The S073 verification record and pull request plan account for #201 through #208 and #215, with no premature closure and no unexplained change.

- [x] T015 [US3] Record selected versions, clean-restoration evidence, focused results, and source-update disposition in `specs/073-dependency-consolidation/verification.md`
- [x] T016 [US3] Add the dependency consolidation and dated pinned Node-baseline decision to `CHANGELOG.md`
- [x] T017 [US3] Reconcile the S073 lifecycle and delivery inventory in `specs/073-dependency-consolidation/spec.md` and `specs/README.md`

**Checkpoint**: The review branch is complete and traceable while all source records remain open until merge.

---

## Phase 6: Canonical Verification and Publication Readiness

**Purpose**: Establish exact local evidence and a publication-safe branch.

- [x] T018 Run the canonical foreground `scripts/verify.sh all` gate and record all eight results in `specs/073-dependency-consolidation/verification.md`
- [x] T019 Audit UTF-8 without BOM, mojibake, GitHub publication formatting, diff integrity, branch scope, issue and pull-request traceability, completed task markers, clean dependency restoration, and the explicit no-release boundary
- [x] T020 Commit the implemented S073 review branch with the required attribution trailer

---

## Dependencies and Execution Order

- Phase 1 fixes the executable scope and completes analysis before implementation.
- Phase 2 protects the Node baseline and dependency contract before manifest changes.
- User Story 1 resolves all graphs and blocks behavior verification.
- User Story 2 depends on User Story 1 because clean manifests are required for build and test evidence.
- User Story 3 records exact results after the graphs and compatibility surface are complete.
- Phase 6 follows all three stories and is the final local publication gate.

## Parallel Opportunities

- T002 can proceed independently of initial checklist validation.
- Root Go, desktop Go, and frontend version research are independent, but their committed integrity files are reconciled sequentially to avoid graph ambiguity.
- Focused root, desktop, and frontend diagnostics can run independently after clean restoration, while canonical verification remains one foreground sequence.

## Implementation Strategy

1. Protect the Node 26 runtime decision with a failing automation contract.
2. Recreate each requested update on current `main`, add only the Node and Vite companions required for valid engines and peers, and regenerate native integrity files.
3. Run focused compatibility evidence and repair only regressions caused by the selected baseline.
4. Complete traceability, analysis reconciliation, and canonical verification.
5. Publish one official replacement pull request, retain all source pull requests until merge, and process no more than two external review rounds.

## Format Validation

All tasks use a checkbox, sequential ID, optional parallel marker, required user-story label within story phases, concrete action, and exact file path.
