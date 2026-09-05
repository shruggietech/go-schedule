# Tasks: Lightweight Pull-Request Workflow

**Task Reconciliation**: PR #45 merged and its review branch was removed; the 2026-08-30 lifecycle audit records the completed publication bookkeeping.

## Phase 1: Scope correction

- [X] T001 Record the operator's decision that PRs exist for AI review, not branch-protection enforcement.
- [X] T002 Remove the earlier branch-protection payload, governance checker, fixtures, and verification composition.
- [X] T003 Revise the spec, research, process model, contract, and plan to match the lightweight scope.

## Phase 2: Synchronized guidance

- [X] T004 Amend constitution v3.0.0 to require the review branch and PR while explicitly leaving merge enforcement to maintainer judgment.
- [X] T005 Update `CLAUDE.md`, `CONTRIBUTING.md`, and `docs/build-autopilot.md` with the same lightweight lifecycle, and synchronize the root `README.md` contributor entry point.
- [X] T006 Simplify `.github/PULL_REQUEST_TEMPLATE.md` to the canonical local aggregate, issue closing, a review summary, and a short checklist.
- [X] T007 Correct `CHANGELOG.md` to record the proportionality decision and remove the superseded protection design.

## Phase 3: Verification and publication

- [X] T008 Run the full aggregate in the foreground and audit scope, encoding, stale prose, and working-tree stability.
- [X] T009 Replace the local Feature 013 commit with the corrected verified slice.
- [X] T010 Halt before pushing `codex/013-pr-governance` and opening the PR.
- [X] T011 After authorization, push and open a PR with `Closes #23`.
- [X] T012 Consider third-party AI review comments, respond under each, and leave final review and merge to the maintainer.
- [x] T013 After merge, remove the merged branch and synchronize local `main`.
