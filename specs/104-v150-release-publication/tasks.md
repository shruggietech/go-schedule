# Tasks: S104 publish v1.5.0

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/publication.md](contracts/publication.md), and [quickstart.md](quickstart.md).

## Phase 1: Source preparation

- [x] T001 Record source, tag, workflow, and waiver constraints in `specs/104-v150-release-publication/` (FR-001 through FR-005).
- [x] T002 [US1] Correct README health example in `README.md` to the intended v1.5.0 source identity (FR-001).
- [x] T003 [US2] Date the v1.5.0 section in `CHANGELOG.md` and fix the matching link in `.github/release-notes/v1.5.0.md` (FR-003).
- [x] T004 Check source metadata, formatting, and local CI parity for `README.md`, `CHANGELOG.md`, and `.github/release-notes/v1.5.0.md` (FR-001, FR-003). All eight local gates passed on 2026-09-23.

## Phase 2: Draft and candidate

- [ ] T005 [US1] Tag exact reviewed main as `v1.5.0` and complete the draft-staging run in `.github/workflows/release.yml` (FR-001, FR-002).
- [ ] T006 [US1] Audit all draft assets and validate the exact Windows MSI against `windows-candidate-manifest.json` (FR-002).

## Phase 3: Public release

- [ ] T007 [US2] Record the standing waiver's application and accurate public limitations on issue #185 and release notes (FR-003, FR-004).
- [ ] T008 [US1] Promote with the applicable full-evidence or expressly authorized waiver path, create `SHA256SUMS.txt`, and verify fresh public downloads (FR-002, FR-004, FR-005).
- [ ] T009 [US1] Reconcile issue #185, v1.5.0 milestone, and project status after publication (FR-005).

## Dependencies

T004 follows T002-T003. T005-T006 require the reviewed main merge. T007 must precede T008. T009 follows confirmed public release. No phase claims completion merely because a draft exists. The PR, review, merge, and branch cleanup are workflow evidence recorded in the specification's Delivery field, not implementation tasks.
