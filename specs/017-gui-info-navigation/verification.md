# Verification Evidence: GUI Navigation and Information

**Date**: 2026-08-28 **Branch**: `codex/017-gui-info-navigation`

## Test-first evidence

### Baseline

- Existing navigation: Tasks, Schedule, Groups, Activity.
- Activity owns a retained tab item whose text changes for its bounded alert count.
- The full application mark is already embedded as `appIcon` in `gui/icon.go`.
- The exact process version is already available as `buildinfo.Version`.
- `go test ./gui -run 'TestUI_(BuildsAllTabs|ActivityBadgeReflectsUnacked)' -count=1` passed before test or production changes.

## Focused verification

### User Story 1 red phase

- Changed the expected navigation prefix to Tasks, Groups, Schedule, Activity and asserted that the retained Activity item occupies index 3.
- `go test ./gui -run 'TestUI_(BuildsAllTabs|ActivityBadgeReflectsUnacked)' -count=1` failed as expected: tab 1 was Schedule instead of Groups.

### User Story 1 green phase

- Reordered the existing management views without reconstructing Activity.
- The same focused command passed with Activity retained at index 3.

### User Story 2 red phase

- Added the exact five-item contract plus isolated mark, long-version, attribution, and canonical-link assertions.
- `go test ./gui -run 'TestUI_(BuildsAllTabs|ActivityBadgeReflectsUnacked|Info)' -count=1` failed as expected because `buildInfoContent` did not yet exist.

### User Story 2 green phase

- Implemented a local, vertically scrollable Info hierarchy using the existing full mark, exact supplied version, descriptive attribution, and standard Fyne hyperlinks, then appended Info after the retained Activity item.
- `go test ./gui -run 'TestUI_(BuildsAllTabs|ActivityBadgeReflectsUnacked|Info)' -count=1` passed.

## Repository verification

### First full run

- Format, vet, lint, and race passed.
- The GUI gate exposed one stale existing schedule-position assertion that still expected Schedule at index 1. This is a legitimate consequence of the specified reorder; the assertion was updated to protect Schedule at index 2.

### Final full run

`scripts/verify.sh all` passed in the foreground:

1. format: passed;
2. vet: passed;
3. lint: passed with 0 issues;
4. race: passed;
5. gui: passed;
6. coverage: passed (engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%, catchup 87.5%, logbus 91.1%);
7. docs: passed for 11 pages, links, front matter, fences, and theme contract;
8. automation: passed for approved actions, CodeQL contract, and the eight-gate manifest.

The final focused Info, navigation, Activity badge, and Schedule-position suite also passed.

## Artifact audits

- `git diff --check` passed.
- All 18 final changed files decoded as strict UTF-8, had no BOM, and contained no common mojibake signature.
- Requirements checklist: 16/16 complete.
- UX checklist: 20/20 complete.
- No pinned artifact or architecture decision changed, so no dated decision record was required.
- Final Spec-Kit analysis covered 17/17 requirements and outcomes with 19/19 valid unique task IDs, no ambiguity or duplication, and no constitution conflict.

## Issue disposition

- #29: complete and eligible for `Closes #29`; Tasks and Groups are adjacent, followed by Schedule and the current Activity terminology.
- #32: complete and eligible for `Closes #32`; Info includes the local mark, exact build version, attribution, and all three canonical links.
