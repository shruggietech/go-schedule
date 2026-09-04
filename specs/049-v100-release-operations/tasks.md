# Tasks: v1.0.0 Release Operations

**Input**: Design documents from `specs/049-v100-release-operations/`

**Tests**: Required by the constitution and release-security scope. Renderer,
atomicity, mapping, and CLI contract tests must fail before implementation.

**Organization**: Tasks execute chronologically. `[P]` records file-level
independence but does not authorize concurrent agent execution in this session.

## Phase 1: Setup and Specification

- [x] T001 Confirm clean synchronized `main`, create `codex/049-v100-release-operations`, activate `specs/049-v100-release-operations` in `.specify/feature.json`, and add its Draft inventory row in `specs/README.md`
- [x] T002 Inventory issue #122, the ten open v1.0.0 readiness issues, merged PRs #121/#123, absent v1.0.0 tag/release state, and current release workflows in `specs/049-v100-release-operations/spec.md`
- [x] T003 Record the S049 merge-boundary deviation, offline-generator, production-validation, packet-format, atomicity, issue-mapping, rendering-safety, compatibility, and documentation decisions in `specs/049-v100-release-operations/research.md`
- [x] T004 Complete and validate `specs/049-v100-release-operations/checklists/requirements.md` and `specs/049-v100-release-operations/checklists/release-operations.md`
- [x] T005 Complete the plan, data model, packet contract, quickstart, and active-plan pointer in `specs/049-v100-release-operations/plan.md`, `specs/049-v100-release-operations/data-model.md`, `specs/049-v100-release-operations/contracts/disposition-packet.md`, `specs/049-v100-release-operations/quickstart.md`, and `CLAUDE.md`

---

## Phase 2: Spec Kit Analysis Gate and Lifecycle

- [x] T006 Run `/speckit-analyze`, correct the #98/#96 traceability error and task-ordering defect, resolve every critical/high finding and requirement/task coverage gap, and record metrics in `specs/049-v100-release-operations/verification.md`
- [x] T007 Advance S049 to In Progress in `specs/049-v100-release-operations/spec.md` and `specs/README.md` before implementation begins

**Checkpoint**: Requirements, design, and tasks have complete coverage with no
critical or high finding.

---

## Phase 3: Foundational Red Tests

- [x] T008 [P] Add failing exact issue-set, observation-mapping, deterministic-output, Markdown-escaping, index-schema, and referenced-environment tests in `internal/releasegate/disposition_test.go`
- [x] T009 [P] Add failing absent-target, existing-target, parent-link, write-failure cleanup, and atomic-commit tests in `internal/releasegate/disposition_test.go`
- [x] T010 [P] Add failing `render-dispositions` success, required-option, formal-class, candidate-identity mutation-matrix, destination-conflict, and stream/exit tests in `scripts/windows-release-gate/main_test.go`, reusing the existing per-field validator tests in `internal/releasegate/validate_test.go`
- [x] T011 [P] Add failing command/help, exact output inventory, no-network, and production-validator reuse assertions in `test/integration/windows_release_gate_contract_test.go`
- [x] T012 Run the focused tests, prove they fail only because disposition models, rendering, atomic writing, and the new CLI path are absent, and record red evidence in `specs/049-v100-release-operations/verification.md`

**Checkpoint**: No implementation exists, and every new trust boundary has a
failing test.

---

## Phase 4: User Story 1 - Generate Evidence-Backed Dispositions (Priority: P1)

**Goal**: Produce ten complete deterministic issue records from already
validated formal evidence.

**Independent Test**: An in-memory passing 47-observation record yields the
fixed file inventory, exact canonical mappings, complete traceability, escaped
Markdown, and identical bytes across two renders.

- [x] T013 [US1] Define copied canonical issue mappings and packet index models in `internal/releasegate/disposition.go`
- [x] T014 [US1] Implement deterministic candidate, workflow, evidence-archive, environment, observation, metrics, and attachment rendering in `internal/releasegate/disposition.go`
- [x] T015 [US1] Implement Markdown table escaping and mention suppression for every evidence-controlled string in `internal/releasegate/disposition.go`
- [x] T016 [US1] Implement deterministic UTF-8-without-BOM `packet.json` rendering with exact issue/file/observation inventory in `internal/releasegate/disposition.go`
- [x] T017 [US1] Run the mapping, rendering, escaping, schema, environment, and determinism tests and record the result in `specs/049-v100-release-operations/verification.md`

**Checkpoint**: Complete packet bytes can be rendered safely in memory without
filesystem or network side effects.

---

## Phase 5: User Story 2 - Fail Closed and Commit Atomically (Priority: P1)

**Goal**: Reuse the full production validator and make the packet visible only
as one complete absent-to-present directory transition.

**Independent Test**: Valid inputs atomically commit one packet; every invalid
candidate/evidence/archive/destination variant emits an actionable failure and
leaves no target or staging residue.

- [x] T018 [US2] Implement absent-target, real-parent, private sibling staging, restrictive-file-mode, cleanup, and atomic-rename behavior in `internal/releasegate/disposition.go`
- [x] T019 [US2] Add `render-dispositions` parsing with mandatory bundle, manifest, MSI, repository, tag, commit, and output options in `scripts/windows-release-gate/main.go`
- [x] T020 [US2] Route `render-dispositions` through strict decoding, bounded extraction, production evidence validation, archive inventory validation, manifest comparison, and expected identity validation in `scripts/windows-release-gate/main.go`
- [x] T021 [US2] Add actionable success/failure reporting and document the new operation in CLI usage in `scripts/windows-release-gate/main.go`
- [x] T022 [US2] Run atomicity, validation-mutation, CLI, race, and integration-contract tests and record the result in `specs/049-v100-release-operations/verification.md`

**Checkpoint**: A report cannot exist unless the exact formal candidate passes
the same production checks used by promotion.

---

## Phase 6: User Story 3 - Execute and Audit the Release Ritual (Priority: P2)

**Goal**: Give the maintainer one chronological path from reviewed S049 merge to
public v1.0.0 without weakening the S048 contract.

**Independent Test**: Documentation and policy tests prove the tag target,
payload cardinalities, packet step, issue boundary, promotion order, and final
audit are explicit while the pull request performs no remote mutation.

- [x] T023 [P] [US3] Update the tag authorization and issue reconciliation steps for the reviewed S049 merge boundary and disposition packet in `specs/048-v100-release-cut/contracts/publication.md`
- [x] T024 [P] [US3] Add formal packet generation, review, and non-authority instructions to `test/windows/README.md`
- [x] T025 [US3] Add positive and negative documentation assertions for the command and S049 boundary in `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`
- [x] T026 [US3] Add the S049 release-operations entry and dated boundary decision to the not-yet-published `CHANGELOG.md` v1.0.0 section while keeping the future Unreleased section empty
- [x] T027 [US3] Run documentation, automation-policy, lifecycle, and packet-contract checks and record the result in `specs/049-v100-release-operations/verification.md`

**Checkpoint**: The reviewed repository states exactly how the later authorized
tag, evidence, issues, promotion, and audit proceed.

---

## Phase 7: Verification and Local Delivery

- [x] T028 [P] Run focused package, CLI, integration, automation-fixture, full-suite, and race tests from `specs/049-v100-release-operations/quickstart.md`
- [x] T029 Run the canonical foreground `scripts/verify.sh all` aggregate and record all eight gate results and coverage in `specs/049-v100-release-operations/verification.md`
- [x] T030 Run strict UTF-8 without BOM, mojibake, unresolved-placeholder, task-format, packet-inventory, documentation-link, and `git diff --check` audits and record results in `specs/049-v100-release-operations/verification.md`
- [x] T031 Mark all implementation tasks complete, advance S049 to Implemented with objective local evidence in `specs/049-v100-release-operations/spec.md` and `specs/README.md`, and create the local `feat(049): add v1.0.0 release operations packet` commit with the required co-author trailer

---

## Dependencies and Execution Order

```text
specification -> analyze -> In Progress -> red tests -> deterministic renderer
              -> atomic CLI -> release runbook -> canonical verification -> commit
```

- T008/T009, T010, and T011 are independent red-test surfaces but precede all
  implementation.
- US1 supplies bytes and mappings consumed by US2's atomic writer and command.
- US3 updates documentation after the command contract is stable.
- Push, pull request, hosted CI, review responses, merge, tag staging, issue
  mutation, release promotion, and milestone closure are delivery or release
  operations, not unchecked implementation tasks.

## Parallel Opportunities

- T008/T009, T010, and T011 touch independent test files.
- T023 and T024 touch independent operator documents.
- T028's focused checks are conceptually independent but run sequentially on
  this Windows host to keep process execution headless and observable.

## Implementation Strategy

US1 is the MVP: deterministic issue-specific evidence without copy/omission
risk. US2 makes that output trustworthy and durable by enforcing production
validation and atomicity. US3 updates the reviewed release boundary and makes
the packet operational. The slice stops at a reviewed, mergeable pull request;
tag and release mutations require later explicit authorization.
