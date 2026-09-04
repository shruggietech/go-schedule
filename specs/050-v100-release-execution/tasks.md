# Tasks: v1.0.0 Release Execution and Audit

**Input**: Design documents from `specs/050-v100-release-execution/`

**Tests**: Formal candidate identity, native evidence, release integrity, and
repository verification are mandatory because S050 publishes the first stable
release and mutates public project state.

## Phase 1: Setup

**Purpose**: Establish the immutable boundary and the isolated audit branch.

- [x] T001 Verify clean synchronized `main`, explicit authorization, absent local/remote release state, create and push annotated `v1.0.0` at `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`, and record the boundary in `specs/050-v100-release-execution/spec.md`
- [x] T002 Create `codex/050-v100-release-execution`, activate `.specify/feature.json`, and add the Draft inventory entry in `specs/README.md`
- [x] T003 Complete specification, clarification, research, data model, release-audit contract, quickstart, and requirement-quality checklists under `specs/050-v100-release-execution/`

---

## Phase 2: Foundational Release State

**Purpose**: Accept one exact staged candidate before native qualification.

**Critical**: No installation, evidence collection, issue mutation, or promotion
may start until this phase passes.

- [x] T004 Require Release run `33838072246` to complete successfully with tag-push event and exact head commit, and record its jobs in `specs/050-v100-release-execution/verification.md`
- [x] T005 Require one draft v1.0.0 release with exactly seven packages plus `windows-candidate-manifest.json`, recording release ID and inventory in `specs/050-v100-release-execution/verification.md`
- [x] T006 Create the absent fixed-volume workspace `A:/_tmp/go-schedule-v1.0.0-s050`, download all eight draft assets, and record safe hashes and sizes in `specs/050-v100-release-execution/verification.md`
- [x] T007 Run production `verify-candidate` against the independently downloaded manifest and MSI, and record complete candidate identity in `specs/050-v100-release-execution/verification.md`

**Checkpoint**: One immutable staged MSI is eligible for formal Windows work.

---

## Phase 3: User Story 1 - Qualify One Immutable Candidate (Priority: P1)

**Goal**: Produce one independently valid 47-observation formal archive from
the exact staged MSI.

**Disposition**: Waived by the explicit publication exception. Formal
qualification remains incomplete.

**Independent Test**: The production bundle validator accepts the finalized ZIP
only with the separately downloaded manifest and MSI.

- [x] T008 [US1] Initialize the formal attended workspace at `A:/_tmp/go-schedule-v1.0.0-s050/attended` from the exact candidate identity
- [x] T009 [US1] Inventory generated observations, fragments, metrics, attachments, and environment requirements in `A:/_tmp/go-schedule-v1.0.0-s050/attended`
- [x] T010 [US1] WAIVED by the explicit publication exception: do not claim unperformed access, failure, task, lifecycle, profile, or reinstall observations
- [x] T011 [US1] WAIVED by the explicit publication exception: do not claim the eleven unperformed formal desktop observations
- [x] T012 [US1] Resolve fragment review by preserving all placeholders as unavailable and refusing to import or relabel them as passing evidence
- [x] T013 [US1] WAIVED by the explicit publication exception: do not finalize an incomplete workspace or create a misleading evidence archive
- [x] T014 [US1] WAIVED because no formal archive exists; preserve the production validator without invoking it on fabricated input
- [x] T015 [US1] Record the zero imported formal observations, absent archive, and explicit maintainer waiver in `specs/050-v100-release-execution/verification.md`

**Checkpoint**: Formal candidate evidence remains incomplete; no archive or
passing observation claim was created.

---

## Phase 4: User Story 2 - Reconcile Readiness Issues Individually (Priority: P1)

**Goal**: Apply evidence only to issues whose complete acceptance criteria pass.

**Disposition**: No evidence disposition was available, so no issue was closed.

**Independent Test**: The draft has nine exact assets, the packet has ten
records plus one index, and live issue/coordinator states match individual
review decisions.

- [x] T016 [US2] WAIVED because no formal archive exists; leave the release without an evidence asset and record the exact public inventory
- [x] T017 [US2] WAIVED because rendering dispositions without a valid archive is correctly rejected by design
- [x] T018 [US2] WAIVED because no disposition packet exists; do not invent packet files or hashes
- [x] T019 [US2] Audit the live evidence-dependent leaf set and preserve every issue as open because the required formal records are absent
- [x] T020 [US2] Apply zero leaf closures and record that conservative disposition in `specs/050-v100-release-execution/verification.md`
- [x] T021 [US2] Preserve #96 as open because its evidence-dependent child criteria are unresolved
- [x] T022 [US2] Record the absence of disposition files, zero closure comments, final issue states, and the publication-exception audit URL in `specs/050-v100-release-execution/verification.md`

**Checkpoint**: Evidence-dependent readiness issues and #122 remain open with
their unresolved criteria visible.

---

## Phase 5: User Story 3 - Publish and Audit v1.0.0 (Priority: P2)

**Goal**: Publish the existing draft under the explicit exception and audit the
public/project identities that can be proven.

**Independent Test**: All nine public assets verify independently, the latest-
release pointer is correct, and issue/milestone state accurately remains open.

- [x] T023 [US3] Reverify remote tag immutability, the exact eight-asset staged state, absent formal evidence, and all readiness issue states before the exception publication
- [x] T024 [US3] WAIVED by direct maintainer instruction: publish the existing draft without dispatching `Promote Release`, and record that exception
- [x] T025 [US3] Generate `SHA256SUMS.txt` for the eight unchanged staged payloads, verify it locally, upload it, and publish the existing release without rebuilding assets
- [x] T026 [US3] Download all nine public assets into `A:/_tmp/go-schedule-v1.0.0-s050/public-audit` and independently verify names, non-empty bytes, and all eight payload checksums
- [x] T027 [US3] Audit public/latest release, tag/commit, manifest/packages, release notes, README, and changelog identity and record results in `specs/050-v100-release-execution/verification.md`
- [x] T028 [US3] Add the exception-aware release audit to #122 and preserve #122 plus milestone `v1.0.0 - Release readiness` as open because their formal criteria did not pass

**Checkpoint**: v1.0.0 is public under exception and byte-audited; formal
qualification remains incomplete and is accurately represented.

---

## Phase 6: Audit Record and Repository Verification

**Purpose**: Deliver the post-release evidence record without changing candidate bytes.

- [x] T029 Update `specs/050-v100-release-execution/spec.md` and `specs/README.md` to Implemented with exact public delivery evidence
- [x] T030 Add the post-release audit note to the Unreleased section of `CHANGELOG.md` without modifying the historical v1.0.0 boundary
- [x] T031 Run `/speckit-analyze`, resolve ordinary consistency findings, and record the intentional publication-exception nonconformance in `specs/050-v100-release-execution/verification.md`
- [x] T032 Run focused release validators plus the foreground canonical `scripts/verify.sh all` suite and record all eight gates in `specs/050-v100-release-execution/verification.md`
- [x] T033 Run strict UTF-8 without BOM, mojibake, placeholder, Markdown-link, task-format, spec-lifecycle, runtime-diff, and `git diff --check` audits
- [x] T034 Mark all tasks complete and create the local S050 audit commit with the required conventional message and co-author trailer

## Dependencies and Execution Order

- Phase 1 is complete.
- Phase 2 blocks every formal observation.
- User Story 1 blocks User Story 2.
- User Story 2 blocks User Story 3.
- Public audit blocks final repository verification and completion claims.
- Pull-request publication and review are workflow evidence after T034, not open feature tasks.

## Parallel Opportunities

- The three platform GUI staging jobs run in parallel inside the Release workflow.
- After download, asset hashing and manifest inspection may run in parallel, but
  candidate validation must finish before installation.
- Issue acceptance reviews may be prepared in parallel after packet generation,
  but leaf mutations and coordinator reconciliation remain ordered.
- Public asset metadata and checksum audits may run in parallel after promotion.

## Implementation Strategy

The minimum viable increment is a valid formal evidence archive (User Story 1),
but it has no publication value by itself. The complete slice is necessarily
sequential: candidate, evidence, dispositions, promotion, public audit, then the
reviewed repository record.
