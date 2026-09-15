# Tasks: Corrected candidate qualification

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, data model, and session contract.

## Setup and foundation

- [x] T001 Complete specification, clarification, requirements-quality checklists, and plan in specs/086-release-candidate-refresh/.
- [x] T002 Analyze artifact coverage and constitutional boundaries before implementation.

## US1 - Prepared guest sessions

- [x] T003 [US1] Write failing preparation tests for hashes, identity, occupied output, overlap, and XML escaping in scripts/windows-qualification-session/main_test.go.
- [x] T004 [US1] Implement create-only preparation, immutable packaged inputs, and independent scenario launch files in scripts/windows-qualification-session/main.go.
- [x] T005 [US1] Implement guest-only input verification/local staging in test/windows/Start-QualificationGuest.ps1.

## US2 - Observable waits

- [x] T006 [US2] Add nondestructive child-process fixtures and Windows regression tests for success, failure, timeout, and hidden launch in scripts/windows-qualification-session/guest_windows_test.go.
- [x] T007 [US2] Implement serialized operations, flushed MSI logs, progress, deadline state, and stop-on-failure in test/windows/Start-QualificationGuest.ps1.

## US3 - Honest attended evidence

- [x] T008 [US3] Add tests ensuring preparation never emits passing attended evidence and refuses guest execution on the development host.
- [x] T009 [US3] Document the remaining native matrix and #229-#233 walkthrough in test/windows/README.md and generated session instructions.

## US4 - Safe refresh

- [x] T010 [US4] Document immutable candidate provenance and explicit old-tag/draft replacement, staging, and promotion boundaries in specs/086-release-candidate-refresh/quickstart.md.

## Validation and local handoff

- [x] T011 Review changes for process visibility, encoding, paths, evidence integrity, and scope; run PowerShell compliance checks.
- [x] T012 Run focused tests and full foreground sh scripts/verify.sh all; record truthful results and update local implementation status.
- [x] T013 Commit the verified local slice and prepare the single pre-publication authorization handoff.

## Dependencies and strategy

T001-T002 precede code. Write T003 before T004-T005, T006 before T007, and T008 before T009. T010 can be authored independently. T011-T013 follow all stories. Installation, hosted publication, review rounds, merge, tag replacement, and attended release qualification are delivery gates rather than local implementation task checkboxes. #226 and #228 remain open.
