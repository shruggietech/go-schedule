# Tasks: S105 official brand kit adoption

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/import.md](contracts/import.md), and [quickstart.md](quickstart.md).

## Phase 1: Release identity and source mapping

- [x] T001 Record the formal release pin in `brand/source.json` and document the exact archive provenance in `brand/REPOSITORY.md` (FR-001, FR-007).
- [x] T002 Define current source-to-consumer paths and any documented transform in `brand/repository-consumers.json`, including docs, frontend, Linux, macOS, and Windows destinations (FR-004, FR-006).

## Phase 2: Maintainer import pipeline, User Story 2

**Independent test**: A verified release import succeeds; altered digest, malicious path, or mismatched manifest fails before changing the installed kit.

- [x] T003 [US2] Add failure-path tests for archive and consumer validation in `internal/brandrelease/` and `scripts/brand-import/` (FR-003, FR-004).
- [x] T004 [US2] Implement `scripts/brand-import/` with pin, release-checksum, ZIP, manifest, bundle, path, and symlink preflight, then selected-asset synchronization and stale-copy cleanup (FR-001 through FR-004).
- [x] T005 [US2] Update `scripts/brand-check/` to verify the stored ZIP and every declared consumer, including the docs manifest transform (FR-002, FR-004).

## Phase 3: Current presentation, User Story 1

**Independent test**: Repository, docs, desktop, and packaging source assets use the simplified mark and the approved current social image.

- [x] T006 [US1] Import the pinned archive and mapped assets, remove the old unpacked kit and stale declared copies, and retain repository-owned Linux desktop metadata (FR-002, FR-006).
- [x] T007 [US1] Update docs, README, desktop and packaging references, contracts, and existing brand-specific tests for the current asset names and roles (FR-005, FR-006).
- [x] T008 [US1] Review the docs and desktop token usage against the official migration contract; correct current-identity drift without expanding platform support (FR-005, FR-006).

## Phase 4: Documentation and verification

- [x] T009 Update `CHANGELOG.md`, `brand/REPOSITORY.md`, `docs/brand.md`, and the S105 quickstart with the repeatable process and source-versus-native-verification boundary (FR-005, FR-007).
- [x] T010 Run S105 analysis and resolve all critical or high findings; execute `go run ./scripts/github-format` and `sh scripts/verify.sh all` in the foreground, then record results (FR-001 through FR-007).

## Dependencies

T002 follows T001. T003 precedes T004 and T005. T006 requires T004 and T005. T007 and T008 follow T006. T009 follows the final migration behavior. T010 follows implementation and documentation. Push, pull-request review, and merge are publication bookkeeping outside the feature task list; this slice stops after the user-authorized push for approval.
