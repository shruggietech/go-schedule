# Tasks: Repository Brand System Integration

**Input**: Design documents from `specs/032-brand-system-integration/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/brand-integrity.md`, `quickstart.md`

**Tests**: Required by the specification, constitution, and autopilot protocol. Integrity behavior is implemented test-first, and every existing CI-parity gate remains mandatory.

**Organization**: Tasks are grouped by user story. Publication bookkeeping is excluded; autopilot halts after the verified local commit.

## Phase 1: Setup and Spec-Kit Gates

**Purpose**: Establish lifecycle state, source provenance, and a reviewed executable contract.

- [x] T001 Record S032 as Draft in `specs/README.md`, activate it in `.specify/feature.json` and `CLAUDE.md`, and validate `specs/032-brand-system-integration/checklists/requirements.md`
- [x] T002 Validate `specs/032-brand-system-integration/checklists/brand-system.md` as requirements-quality coverage for canonical assets, documentation, accessibility, and consumer synchronization
- [x] T003 Run the blocking Spec-Kit analysis over `specs/032-brand-system-integration/spec.md`, `plan.md`, and `tasks.md`, resolving all critical or high findings before implementation
- [x] T004 Capture the approved import source, standalone verification result, exact exclusion policy, and initial 3 MiB size check in `specs/032-brand-system-integration/verification.md`

---

## Phase 2: Foundational Canonical Import

**Purpose**: Bring the approved standalone kit into the repository before any consumer wiring.

**Critical**: Consumer work cannot begin until the canonical source and mapping schema exist.

- [x] T005 Import the complete approved kit into `brand/`, excluding ZIP archives, backup archives, `__pycache__`, temporary environments, and rendered audit previews
- [x] T006 Add repository ownership, canonical-versus-copy rules, regeneration prerequisites, and ordinary validation guidance in `brand/REPOSITORY.md`
- [x] T007 Define versioned exact-copy relationships for documentation, fonts, favicons, GUI, and Windows consumers in `brand/repository-consumers.json`
- [x] T008 Confirm `brand/manifest.json`, `brand/VERIFY.md`, and `brand/brand-guide.pdf` retain the approved standalone evidence and are not rewritten as repository-control substitutes

**Checkpoint**: The complete standalone brand system is a discoverable repository source of truth.

---

## Phase 3: User Story 1 - Obtain the Complete Canonical Kit (Priority: P1)

**Goal**: Make the imported kit self-validating in a normal checkout without optional graphics dependencies.

**Independent Test**: `go test ./scripts/brand-check` and `go run ./scripts/brand-check` validate every canonical hash, SVG, text file, and declared consumer with an actionable failure for deliberate drift.

### Tests for User Story 1

- [x] T009 [P] [US1] Add failing valid-repository, missing-artifact, hash-mismatch, malformed-manifest, and path-escape cases in `scripts/brand-check/manifest_test.go`
- [x] T010 [P] [US1] Add failing UTF-8 BOM, invalid UTF-8, mojibake, missing SVG title, live SVG text, and font-family cases in `scripts/brand-check/text_test.go`
- [x] T011 [P] [US1] Add failing missing-consumer, changed-consumer, duplicate-target, malformed-map, and required-control-file cases in `scripts/brand-check/consumers_test.go`

### Implementation for User Story 1

- [x] T012 [US1] Implement safe manifest parsing, bounded paths, streaming length/SHA-256 validation, and accumulated actionable failures in `scripts/brand-check/main.go`
- [x] T013 [US1] Implement UTF-8/BOM/mojibake and portable accessible-SVG checks in `scripts/brand-check/main.go`
- [x] T014 [US1] Implement consumer-map validation, byte comparison, summary counts, and command exit contract in `scripts/brand-check/main.go`
- [x] T015 [US1] Document focused validation and standalone-versus-repository evidence in `brand/REPOSITORY.md` and `specs/032-brand-system-integration/contracts/brand-integrity.md`

**Checkpoint**: Canonical inventory and every required copy can be verified offline with the existing Go toolchain.

---

## Phase 4: User Story 2 - Discover and Download the Brand (Priority: P1)

**Goal**: Publish the brand system through a useful, accessible documentation page.

**Independent Test**: `sh scripts/verify.sh docs` validates front matter and every local download; browser inspection confirms readable desktop and narrow layouts with the required previews and selection guidance.

### Tests for User Story 2

- [x] T016 [P] [US2] Extend documentation contract coverage for the Brand system page, required download categories, and approved source paths in `test/scripts/docs-policy-check_test.sh`
- [x] T017 [US2] Add brand-page policy rules and actionable failures in `scripts/docs-policy-check.sh`

### Implementation for User Story 2

- [x] T018 [US2] Replace `docs/assets/brand/`, `docs/assets/favicons/`, and `docs/assets/fonts/` with declared approved consumer copies from `brand/`
- [x] T019 [US2] Publish identity, variant selection, palette, type, accessibility, attribution, misuse, and downloads in `docs/brand.md`
- [x] T020 [US2] Align the docs logo and social metadata with canonical mapped assets in `docs/_config.yml`, `docs/_includes/head_custom.html`, and `docs/_sass/custom/custom.scss`
- [x] T021 [US2] Validate the page source and preview presentation at desktop and narrow widths, recording evidence in `specs/032-brand-system-integration/verification.md`

**Checkpoint**: Visitors can discover, understand, and download the approved kit from the documentation site.

---

## Phase 5: User Story 3 - Keep Product Surfaces Synchronized (Priority: P2)

**Goal**: Consume canonical identity across application and release surfaces while mechanically rejecting drift.

**Independent Test**: The brand verifier, GUI tests, Windows installer verifier, release-contract checks, and automation gate all pass; a deliberate mapped-copy mutation fails with exact paths.

### Tests for User Story 3

- [x] T022 [P] [US3] Update embedded-icon expectations and size/purpose coverage in `gui/app_test.go` and `gui/info_test.go`
- [x] T023 [P] [US3] Update canonical Windows icon source assertions in `build/windows/verify_wxs.ps1`
- [x] T024 [P] [US3] Add offline release-contract assertions for canonical Windows/macOS inputs and Linux desktop assets in `scripts/automation-check.sh`

### Implementation for User Story 3

- [x] T025 [US3] Synchronize `gui/assets/icon.png`, `gui/assets/icon-window.png`, bundled GUI fonts/licenses, and `cmd/gosched-gui/icon.ico` from declared canonical sources
- [x] T026 [US3] Point Windows executable and MSI identity at the canonical ICO in `.github/workflows/release.yml` and `build/windows/goschedule.wxs`
- [x] T027 [US3] Replace macOS release-time icon resampling with the canonical ICNS in `.github/workflows/release.yml`
- [x] T028 [US3] Add the canonical `.desktop` entry and hicolor icon tree to Linux desktop archives in `.github/workflows/release.yml`
- [x] T029 [US3] Run `go run ./scripts/brand-check` from `scripts/automation-check.sh` and keep the existing eight-gate manifest unchanged

**Checkpoint**: Every current identity surface is canonical or declared, and routine automation detects drift.

---

## Phase 6: Documentation, Verification, and Delivery

**Purpose**: Close the two brand issues with review-ready evidence and no claim against deferred Windows observation work.

- [x] T030 [P] Update canonical-kit discovery and maintainer workflow in `README.md`, `CONTRIBUTING.md`, and `docs/README.md`
- [x] T031 [P] Add the feature and dated pinned-artifact/brand-architecture decision in `CHANGELOG.md`
- [x] T032 Advance S032 to In Progress and run focused kit, docs, GUI, Windows packaging, and deliberate-drift validation from `specs/032-brand-system-integration/quickstart.md`
- [x] T033 Run `sh scripts/verify.sh all` in the foreground and capture format, vet, lint, race, GUI, coverage, docs, and automation evidence in `specs/032-brand-system-integration/verification.md`
- [x] T034 Re-run Spec-Kit analysis with zero critical/high findings and validate both S032 checklists against the delivered result
- [x] T035 Confirm issue #10 and #34 closing eligibility, explicitly preserve #33, scan UTF-8/BOM/mojibake, audit the complete diff and imported size, and record final evidence in `specs/032-brand-system-integration/verification.md`
- [x] T036 Mark all tasks complete, advance S032 to Implemented in `specs/032-brand-system-integration/spec.md` and `specs/README.md`, and commit locally as `feat(032): integrate the complete brand system`

---

## Dependencies and Execution Order

- Phase 1 gates all implementation.
- Phase 2 establishes the authoritative kit and mapping schema and blocks all user stories.
- User Story 1 supplies the integrity command required by User Stories 2 and 3.
- User Story 2 and User Story 3 can proceed independently after User Story 1.
- Phase 6 depends on every story checkpoint.
- Within each story, failing tests precede implementation.

## Parallel Opportunities

- T009, T010, and T011 cover independent verifier failure classes in separate test files.
- T016 and T017 have separate test/policy files but remain test-before-implementation ordered.
- T022, T023, and T024 cover independent GUI, Windows, and release-contract surfaces.
- T030 and T031 touch independent documentation and changelog files.
- This autopilot run uses one working agent because repository instructions do not authorize subagent delegation; `[P]` records dependency structure rather than active delegation.

## Implementation Strategy

Deliver the complete brand integration as one review unit. The canonical import and integrity contract land first, followed by the public documentation surface and all product/release consumers. No partial publication checkpoint is introduced because splitting canonical assets, documentation, and drift enforcement would leave temporary competing sources and create additional PR ceremony without useful review isolation.
