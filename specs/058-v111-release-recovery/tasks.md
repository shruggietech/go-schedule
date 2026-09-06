# Tasks: v1.1.1 Release Recovery

**Input**: Design documents from `specs/058-v111-release-recovery/`

**Tests**: Release identity, copy, historical preservation, workflow contracts, repository gates, candidate integrity, and public checksums are mandatory.

## Phase 1: Setup and specification

- [x] T001 Verify synchronized main, the original v1.1.0 tag commit, unpublished v1.1.0 draft state, absent v1.1.1 state, and successful post-merge CI for PR #143.
- [x] T002 Create branch `codex/058-v111-release-recovery` and initialize the Spec Kit artifact set.
- [x] T003 Complete specification, clarification, requirements checklist, research, data model, plan, publication contract, and quickstart.
- [x] T004 Analyze requirement coverage, constitution compliance, chronological ordering, and scope boundaries before implementation.

## Phase 2: Corrective release preparation

- [x] T005 Revise issue #140 and milestone #4 to make v1.1.1 the authoritative release target without duplicating planning records.
- [x] T006 Add a dated v1.1.1 changelog section for the watcher and exact-commit CI corrections while preserving the v1.1.0 historical section.
- [x] T007 Update the README health example to 1.1.1.
- [x] T008 Add exactly four concise highlights and one final tagged changelog link in `.github/release-notes/v1.1.1.md`.
- [x] T009 Run release-copy, lifecycle, formatting, focused integration, and complete eight-gate verification.
- [x] T010 Record exact verification evidence and mark S058 implementation complete.

## Phase 3: Publication operations

Publication operations are tracked by issue #140 rather than represented as repository implementation tasks. They proceed after the reviewed preparation merge in this order: revalidate and preserve the v1.1.0 tag, retire only its unpublished draft, create the immutable v1.1.1 tag, require successful main CI for the exact commit, stage draft artifacts, qualify the exact candidate, promote without rebuild, perform the final checksum and identity audit, then close the issue and milestone.

## Dependencies and execution order

T001 through T004 establish release authority and an internally consistent recovery design. T005 synchronizes authoritative hosted planning. T006 through T008 prepare the source identity and public copy. T009 validates the result before T010 records completion. Public release operations begin only after the preparation pull request is reviewed, green, and merged.
