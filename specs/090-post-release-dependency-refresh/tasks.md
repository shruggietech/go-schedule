# Tasks: Post-release dependency refresh

**Input**: Design documents from `specs/090-post-release-dependency-refresh/`

## Phase 1: Specification and Traceability

- [x] T001 Create #243 and inventory #240, #241, and #242.
- [x] T002 Author the S090 specification, plan, research, data model, contract, quickstart, checklists, tasks, and verification record.
- [x] T003 Record S090 in `.specify/feature.json`, `CLAUDE.md`, and `specs/README.md`, then complete cross-artifact analysis.

## Phase 2: Go Baseline and Graphs

- [x] T004 [US1] Advance root and desktop modules plus operational guidance to Go 1.26.
- [x] T005 [US1] Apply go-sdk 1.8.0, x/crypto 0.57.0, x/sys 0.48.0, and x/time 0.16.0 and regenerate the root graph.
- [x] T006 [US1] Reconcile and verify the desktop Go graph against the updated root replacement.

## Phase 3: Frontend Graph

- [x] T007 [US1] Apply React 19.3.0, matching types, Node types 26.5.1, and Vite 8.3.0.
- [x] T008 [US1] Regenerate the lockfile and prove clean Node 26 restoration and high-severity audit.

## Phase 4: Compatibility Verification

- [x] T009 [US2] Run focused MCP, remote security, authorization, audit, and Windows platform tests.
- [x] T010 [US2] Run desktop Go race, frontend unit, type-check, production build, browser, native Wails, and packaging checks.
- [x] T011 [US2] Resolve only compatibility failures caused by the selected dependency baseline.

## Phase 5: Evidence and Publication

- [x] T012 [US3] Complete the version inventory, graph-cleanliness evidence, changelog decision, and S090 verification record.
- [x] T013 [US3] Run all eight canonical gates, GitHub formatting, UTF-8 and mojibake checks, and diff integrity checks.
- [ ] T014 [US3] Commit, push, and publish the official S090 pull request with complete issue and source-PR traceability.
- [ ] T015 [US3] Comment on and close #240, #241, and #242 as superseded by the official replacement.
- [ ] T016 [US3] Address every first-round review finding, optionally request one second Codex round, and address every resulting finding.
- [ ] T017 [US3] Confirm all required latest-head checks are green and report merge readiness to the maintainer.

## Dependencies and Execution Order

- Specification and analysis precede dependency changes.
- The Go baseline decision precedes Go graph regeneration.
- Root graph resolution precedes desktop graph reconciliation.
- Each ecosystem must restore cleanly before focused behavior checks.
- Canonical verification follows all focused checks and precedes publication.
- Source pull-request closure follows replacement publication; merge readiness follows final review and CI.
