# Tasks: Windows version identity

## Phase 1: Setup

- [x] T001 Resolve #458 scope, installed commands and source authority in spec.md and research.md (FR-005, FR-006).

## Phase 2: Foundation

- [x] T002 Write regression-first version, configuration restore and lock tests in scripts/desktop-build/main_test.go (FR-002, FR-003, FR-004).

## Phase 3: US1 release identity

- [x] T003 [US1] Implement selected-version builder in scripts/desktop-build/main.go and wire release.yml (FR-001, FR-002).
- [x] T004 [US1] Add compiled PE/MSI mismatch refusal in build/windows/verify_version.ps1 and existing CI/release packaging paths (FR-001, FR-004).

## Phase 4: US2 development identity

- [x] T005 [US2] Wire Makefile and scripts/verify.sh with explicit development defaults and document local version selection (FR-002, FR-003).

## Phase 5: Verification and delivery records

- [x] T006 Inspect actual compiled Windows identity and mismatch refusal without candidate execution (FR-001, FR-004, FR-006).
- [x] T007 Run all eight foreground gates and record actual outcomes in verification.md (FR-001 through FR-006).
- [x] T008 Complete lifecycle, changelog, source and protected-evidence integrity records (FR-005, FR-006).

## Dependencies and strategy

T001 -> T002 -> T003 -> T004 -> T005 -> T006 -> T007 -> T008. US1 is the build identity MVP; US2 adds local/default paths. Same build files require sequential edits. Research may run independently, but source writes remain sequential. No parallel source task is required. Hosted checks and final owner merge are publication workflow, not unchecked implementation tasks.
