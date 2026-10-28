# Tasks: S125 Reliability Delivery and Scheduling Decision

**Input**: [spec.md](spec.md), [plan.md](plan.md), research, data model and delivery contract.
**Tests**: Required by the constitution for the demonstrated browser sequencing correction.

## Phase 1: Setup

- [x] T001 Synchronize and verify the intended repository, review branch and feature selector in .specify/feature.json.
- [x] T002 Complete specify/clarify/checklist/plan with approved boundaries in specs/125-reliability-delivery-decision/spec.md and plan.md.

## Phase 2: Foundations

- [x] T003 Preserve the original merged-main failure and assemble research findings in specs/125-reliability-delivery-decision/research.md.
- [x] T004 Run the read-only coverage/constitution analysis of spec.md, plan.md and tasks.md, then retain its receipt in specs/125-reliability-delivery-decision/analysis.md.

## Phase 3: US1 Browser reliability (Priority: P1)

**Goal**: Correct the demonstrated fixture race and preserve useful failed-run evidence.
**Independent test**: Controlled late route focus exposes the original ordering and the corrected query remains intact; full browser contracts retain their assertions.

- [x] T005 [US1] Add a deterministic delayed-focus regression before the ordering correction in desktop/frontend/e2e/cross-daemon-search.spec.ts and preserve the failing receipt.
- [x] T006 [US1] Synchronize route readiness and assert exact query value in desktop/frontend/e2e/cross-daemon-search.spec.ts (FR-001).
- [x] T007 [P] [US1] Retain bounded browser failure diagnostics in .github/workflows/ci.yml without changing exercise failure status (FR-002/FR-009).
- [x] T008 [US1] Run the corrected regression and full headless browser suite, recording original and corrected identities in specs/125-reliability-delivery-decision/verification.md.

## Phase 4: US2 Release continuation (Priority: P1)

**Goal**: Finish the authorized source outcome and make actual release continuation concrete, with no invented source/asset identity.
**Independent test**: Existing metadata/comparison checks accept pending source while published discovery stays v1.6.0.

- [x] T009 [US2] Reconcile owner-merged S124, its failed main receipt and #397 in specs/124-reliability-release-preparation/delivery.md and handoff.md (FR-003).
- [x] T010 [US2] Include the bounded #397 correction and dated pinned-artifact rationale in CHANGELOG.md and .github/release-notes/v1.7.0.md.
- [x] T011 [US2] Rebind content review digests in site/data/reviews/v1.7.0.json without unsupported capability changes.
- [x] T012 [US2] Write the complete guarded release continuation in specs/125-reliability-delivery-decision/handoff.md and link current maintainer guides in docs/INSTALL-windows.md and test/windows/README.md (FR-004).

## Phase 5: US3 Distributed decision (Priority: P2)

**Goal**: Deliver the evidence register and gated disposition without fabricating demand or completing gated work.
**Independent review**: Every demand claim has a source or explicit missing input; both downstream gates remain closed on insufficient evidence.

- [x] T013 [P] [US3] Document real owner intent, observed evidence, alternatives, exact intake and reopening criteria in docs/distributed-scheduling.md (FR-005).
- [x] T014 [US3] Record continue/hold/stop alternatives and explicit #186/#189/#20 incomplete outcomes and #187/#188 gate results in specs/125-reliability-delivery-decision/decision.md (FR-006/FR-007).
- [x] T015 [US3] Link the decision document from docs/remote-access.md and reconcile native issue records without unsupported closing references (FR-008).

## Phase 6: Verification and Delivery

- [x] T016 Reconcile source lifecycle and inventories in specs/README.md, specs/124-reliability-release-preparation/ and specs/125-reliability-delivery-decision/.
- [x] T017 Run GitHub formatting, UTF-8/no-BOM/corruption and publication-language checks on the final source proposal; retain receipts in specs/125-reliability-delivery-decision/verification.md.
- [x] T018 Freeze retained original QA/evidence inputs and run the complete foreground canonical aggregate, recording exact source identity and unchanged evidence in specs/125-reliability-delivery-decision/verification.md.
- [x] T019 Complete delivery.md and all issue-level dispositions in specs/125-reliability-delivery-decision/delivery.md; commit the coherent source proposal after verification.

## Dependencies and Execution Order

T001-T004 precede implementation. T005 precedes T006; T008 follows T006/T007. US2 and US3 proceed independently once analysis passes; T011 follows T010. T016-T019 follow all source work. Spec-Kit research agents may implement disjoint file groups after the gate; shared changelog and release review files stay with one owner. MVP is US1, but deliver the full authorized compatible bundle.

## Publication Protocol (not feature tasks)

The owner already authorizes branch push and one PR, exact-head CI monitoring and at most two Codex rounds. Review/publication bookkeeping does not keep implementation tasks open. Stop before merging, tagging or promotion. #394 and conditional #186/#189/#20 remain open while their actual functional outcomes are incomplete. #187/#188 do not become prototype work unless their gates pass.
