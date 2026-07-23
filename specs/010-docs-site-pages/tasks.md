# Tasks: Documentation site on GitHub Pages + README consolidation

**Feature**: 010-docs-site-pages | **Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

Organized by user story. No Go code changes; the docs-check script is this
feature's automated test. Tests are the contract scenarios in
[contracts/docs-check.md](contracts/docs-check.md), exercised in Polish.

**[P]** = parallelizable (distinct file, no dependency on an incomplete task).

## Phase 1: Setup (shared site scaffold)

- [x] T001 Create `docs/_config.yml` with `remote_theme: just-the-docs/just-the-docs@v0.4.2`, `title`, `description`, `url: https://shruggietech.github.io`, `baseurl: /go-schedule`, `search_enabled: true`, `aux_links` to the GitHub repo, and `plugins: [jekyll-remote-theme, jekyll-relative-links, jekyll-seo-tag, jekyll-sitemap]` (per [research.md](research.md) D1/D2, [data-model.md](data-model.md)).

## Phase 2: Foundational

No additional blocking prerequisites — `docs/_config.yml` (Setup) is the only
shared dependency. User-story phases follow directly.

## Phase 3: User Story 1 — Searchable, navigable site (Priority: P1)

**Goal**: Every `docs/` page carries nav metadata, install guides are grouped,
and no in-page link escapes the served root — so the published site renders with
search, a spanning sidebar, and next/previous, with zero 404s.

**Independent test**: With Pages pointed at `/docs`, the home page renders, the
sidebar lists every page, search returns results, the three install guides sit
under one Installation section, and every in-page link resolves.

- [x] T002 [P] [US1] Add front matter (`title: Home`, `nav_order: 1`) to `docs/README.md` and rewrite its six `../` links (`../README.md`, `../CONTRIBUTING.md`, `../SECURITY.md`, `../CHANGELOG.md`, `../.specify/memory/constitution.md`, `../specs/001-task-scheduler/spec.md`) to absolute `https://github.com/shruggietech/go-schedule/blob/main/...` URLs.
- [x] T003 [P] [US1] Create `docs/install.md` — Installation section index (`title: Installation`, `nav_order: 2`, `has_children: true`) with a short intro linking the three platform guides.
- [x] T004 [P] [US1] Add front matter (`title: Windows`, `parent: Installation`, `nav_order: 1`) to `docs/INSTALL-windows.md` (PINNED artifact — the dated CHANGELOG entry is T017).
- [x] T005 [P] [US1] Add front matter (`title: Linux`, `parent: Installation`, `nav_order: 2`) to `docs/INSTALL-linux.md`.
- [x] T006 [P] [US1] Add front matter (`title: macOS`, `parent: Installation`, `nav_order: 3`) to `docs/INSTALL-macos.md`.
- [x] T007 [P] [US1] Add front matter (`title: CLI reference`, `nav_order: 3`) to `docs/cli.md`.
- [x] T008 [P] [US1] Add front matter (`title: GUI field reference`, `nav_order: 4`) to `docs/gui-fields.md` and rewrite its `../specs/001-task-scheduler/contracts/cli.md` link to an absolute repo URL.
- [x] T009 [P] [US1] Add front matter (`title: Cron interoperability`, `nav_order: 5`) to `docs/cron.md`.
- [x] T010 [P] [US1] Add front matter (`title: Maintainer test scripts`, `nav_order: 6`) to `docs/test-scripts.md` and rewrite its four `../` links (`../.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1`, `../specs/006-maintainer-test-scripts/`, `../specs/006-maintainer-test-scripts/data-model.md`, `../test/scripts/lib/sqlite-manifest.json`) to absolute repo URLs (`blob/main` for files, `tree/main` for the directory).
- [x] T011 [P] [US1] Add front matter (`title: Build-phase autopilot`, `nav_order: 7`) to `docs/build-autopilot.md`.

**Checkpoint**: `grep -rEn '\]\(\.\./' docs/*.md` returns nothing; every page has front matter.

## Phase 4: User Story 2 — Single source of truth, no silent drift (Priority: P2)

**Goal**: A fast, network-free gate fails on broken internal links, stale
pointers, or missing front matter; the SSOT + pointer convention is written down.

**Independent test**: `sh scripts/docs-check.sh` exits 0 on the consistent tree;
introducing a broken link / missing pointer target / missing front matter makes
it exit non-zero naming the offender (contract scenarios T2–T7).

- [x] T012 [US2] Create `scripts/docs-check.sh` (POSIX `sh`, `set -eu`, style of `scripts/coverage-gate.sh`) implementing the contract in [contracts/docs-check.md](contracts/docs-check.md): front-matter `title`+`nav_order` presence; on-disk link resolution (strip `#fragment`, skip `http(s)`, skip pure `#`); no `../` escape from `docs/`; pointer README `test/scripts/README.md` target exists. Non-zero lists `file: reason: link`.
- [x] T013 [US2] Add a `docs` job to `.github/workflows/ci.yml` (ubuntu, `actions/checkout@v4`, `run: sh scripts/docs-check.sh`, no `needs`, existing push/PR-to-`main` triggers) (PINNED artifact — dated CHANGELOG entry is T017).
- [x] T014 [P] [US2] Add a "Documentation" section to `CONTRIBUTING.md`: `docs/` is the single source of truth, subdirectory READMEs are pointers (not copies), the site is served branch-based from `/docs`, and `sh scripts/docs-check.sh` runs before pushing.

**Checkpoint**: `sh scripts/docs-check.sh` is green; CI `docs` job defined.

## Phase 5: User Story 3 — Every inbound reference points at the site (Priority: P3)

**Goal**: README and the issue-form contact links resolve to the site; the
README shows one consistent version.

**Independent test**: README has a site Documentation link and reads `0.7.0` in
quick-start; the issue chooser's two contact links open site URLs.

- [x] T015 [P] [US3] In `README.md`, add a prominent Documentation link to `https://shruggietech.github.io/go-schedule/` near the top, and fix the quick-start version drift `version 0.6.0` → `version 0.7.0` (line ~125). Do NOT edit the release badge line (auto-bumped by `release.yml`).
- [x] T016 [P] [US3] Repoint the two `contact_links` in `.github/ISSUE_TEMPLATE/config.yml` from raw `tree/main/docs` and `blob/.../test-scripts.md` URLs to `https://shruggietech.github.io/go-schedule/` and `https://shruggietech.github.io/go-schedule/test-scripts.html`.

## Phase 6: Polish & cross-cutting

- [x] T017 Update `CHANGELOG.md` `[Unreleased]`: an `### Added` feature line for the docs site + consolidation (referencing `closes #11`); `### Decisions` entries dated `2026-07-23` (branch-based `/docs` + just-the-docs@v0.4.2 vs Hugo/MkDocs+Actions; link-checker strips `#fragments`); and two `### Changed` "**Pinned artifact — `<path>` (2026-07-23).**" entries for `.github/workflows/ci.yml` (new `docs` job) and `docs/INSTALL-windows.md` (front matter added).
- [x] T018 Run `sh scripts/docs-check.sh` locally (expect T1 = exit 0), then exercise contract negative scenarios T2–T7 with temporary edits reverted afterward, confirming each fails and names the offender.
- [x] T019 Run full CI parity in the foreground per the `go-schedule-verify` skill: `gofmt -l internal cmd test`, `go vet ./...`, golangci-lint, the `-race` test set, `go test ./gui/...`, `sh scripts/coverage-gate.sh`. No Go changed → expect green; report the `-race` gate honestly (CI-only if no local C toolchain).

## Dependencies & execution order

- **T001 (Setup)** first — the site scaffold. Front-matter tasks (T002–T011) may
  proceed in parallel with T001 (distinct files).
- **US1 (T002–T011)** should complete before **US2's T012** run is expected to
  pass, because the docs-check validates US1's front matter and rewritten links.
  (T012 authoring is independent; its *green run* depends on US1.)
- **US3 (T015–T016)** is independent of US1/US2 and can run any time after Setup.
- **Polish**: T017 after T004 + T013 (the pinned changes exist). T018 after all
  `docs/` + script changes. T019 last, immediately before the pre-push halt.

## Parallel opportunities

- T002–T011 are all `[P]` (nine distinct `docs/*.md` files + one new file).
- T014 (`CONTRIBUTING.md`), T015 (`README.md`), T016 (`config.yml`) are `[P]`
  across distinct files and independent of the `docs/` edits.

## MVP scope

**User Story 1** alone (T001–T011) is a shippable MVP: a real, searchable,
navigable published site. US2 adds the anti-drift guarantee; US3 adds
discoverability polish.

## Format validation

All tasks use `- [ ] Txxx [P?] [US?] description + file path`. Setup and Polish
tasks carry no story label by design; US1/US2/US3 tasks are labelled.
