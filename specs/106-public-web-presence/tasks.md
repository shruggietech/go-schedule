# Tasks: Public website and project discovery

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, data model and contracts. All six issues remain individually traceable.

## Phase 1: Setup

- [x] T001 Initialize `site/package.json`, exact lockfile, Next/Fumadocs configuration and generated-file ignores. (FR-011)
- [x] T002 Record the published baseline, dependency/license choices and approved brand consumers in `site/data/` and `brand/repository-consumers.json`. (FR-002/007/011)

## Phase 2: Foundation

- [x] T003 Write failing document/route/link transformation tests in `site/tests/content.test.mjs`. (FR-003/004)
- [x] T004 Implement canonical Markdown preparation, asset copying and route inventory in `site/scripts/content.mjs`. (FR-003/011)

## Phase 3: User Story 1, discovery and first use

- [x] T005 [US1] Build shared chrome, appearance and brand styling in `site/components/` and `site/app/`. (FR-001/002)
- [x] T006 [US1] Build the product home and static 404 in `site/app/page.tsx` and `site/app/not-found.tsx`. (FR-001)
- [x] T007 [US1] Rewrite `README.md` and restore #9's final attribution while retaining release/policy contracts. (FR-005)

## Phase 4: User Story 2, documentation continuity

- [x] T008 [US2] Implement Fumadocs routes, navigation, static search, code and Mermaid components in `site/app/docs/` and `site/app/static.json/`. (FR-003/004)
- [x] T009 [US2] Implement legacy aliases, export link audit and metadata in `site/scripts/postbuild.mjs` and `site/app/`. (FR-003)
- [x] T010 [US2] Replace Jekyll-specific checks and add exact site build/browser CI in `scripts/docs-check.sh` and `.github/workflows/ci.yml`. (FR-003/004/011)

## Phase 5: User Story 3, competitive capability and release review

- [x] T011 [US3] Write failing structural/release identity tests in `scripts/comparison-check/main_test.go`. (FR-006/007/008)
- [x] T012 [US3] Add all seven sourced products, capability cells and opportunity links in `site/data/comparisons/`. (FR-006/007/009)
- [x] T013 [US3] Implement offline structural/release validation in `scripts/comparison-check/main.go`. (FR-008)
- [x] T014 [US3] Render the accessible dedicated comparison page in `site/app/compare/page.tsx`. (FR-006/007/009)
- [x] T015 [US3] Wire the shared guard into `.github/workflows/release.yml`, both promotion workflows and automation regression fixtures. (FR-008)

## Phase 6: User Story 4, canonical domain preparation

- [x] T016 [US4] Obtain owner Cloudflare authentication for read-only inventory and record redacted findings in `specs/106-public-web-presence/delivery.md`. (FR-010/012)
- [x] T017 [US4] Prepare idempotent DNS/Pages actions, rollback and authentication handoff in `docs/website.md`; mutate no public settings. (FR-010/012)
- [x] T018 [US4] Add `.github/workflows/pages.yml` with static artifact deployment and explicit latest-public-release snapshot selection. (FR-007/010/011)

## Phase 7: Integration and local review

- [x] T019 Add production browser tests and exported internal-link/fragment checks in `site/tests/`; exercise search, appearance, keyboard and mobile. (FR-001/002/003/004)
- [x] T020 Update `CHANGELOG.md`, `CONTRIBUTING.md`, `docs/website.md` and S106 delivery traceability, run full CI parity and record actual outcomes. (FR-011/012, SC-001/005)
- [x] T021 Commit verified local S106 changes on the review branch and present the exact pending publication/domain actions without executing them. (FR-012)

## Dependencies and execution order

T001/T002 precede T003/T004; shared preparation precedes US1/US2 rendering. US3 regression tests/data/guard may proceed alongside independent content work, but render integration consumes the shell. T016 needs owner authentication; pending authentication does not block T005-T015. US4 is deliberately a prepared local handoff at this boundary; public #278 completion remains pending cutover. T019/T020 cover integrated behavior before T021. No push, PR, deployment, tag, release or DNS mutation is an unchecked implementation task.

## Parallel opportunities and strategy

README editing and competitor evidence gathering affect different files from the documentation compiler. Source research used separate agents as required by the planning command. Execute shared-file changes sequentially. Complete all four local story increments and the domain handoff before review; a home-only MVP does not satisfy the authorized bundle.
