# Implementation Plan: Documentation site on GitHub Pages + README consolidation

**Branch**: `010-docs-site-pages` (trunk-based — committed onto `main`) | **Date**: 2026-07-23 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/010-docs-site-pages/spec.md`

## Summary

Publish the existing `docs/` Markdown set as a searchable, navigable GitHub
Pages site served branch-based from the `/docs` folder on `main`, using the
just-the-docs remote theme (pinned to a Jekyll-3.9-compatible tag). Establish
`docs/` as the single source of truth: subdirectory READMEs stay thin pointers,
inbound references (README, issue templates) point at the site, and a POSIX-sh
CI link-check prevents drift. No Go code changes; no committed build artifacts;
one operator settings change to go live.

## Technical Context

**Language/Version**: No application language change. New artifacts are YAML
(`docs/_config.yml`), Markdown front matter (`docs/*.md`), and POSIX `sh`
(`scripts/docs-check.sh`). Site build: GitHub Pages' bundled Jekyll 3.9.x
(github-pages gem) — run by the host, not this repo.

**Primary Dependencies**: just-the-docs remote theme pinned to `@v0.4.2`;
GitHub-Pages-allowlisted plugins `jekyll-remote-theme`, `jekyll-relative-links`,
`jekyll-seo-tag`, `jekyll-sitemap`. No Gemfile is committed (branch-based Pages
resolves the remote theme + allowlisted plugins itself).

**Storage**: N/A (static docs).

**Testing**: `scripts/docs-check.sh` (front-matter + on-disk link + pointer
integrity), run locally and as a new `docs` CI job. Existing Go gates unchanged.

**Target Platform**: A static website on GitHub Pages + the repository file view.

**Project Type**: Documentation site over an existing Go project.

**Performance Goals**: N/A — a static site. The link-check must run in seconds
with no network dependency.

**Constraints**: Source == served content (no committed HTML). Branch-based
Pages only (no deploy workflow). Theme must build under Jekyll 3.9 (libsass).
`docs/INSTALL-windows.md` and `.github/workflows/ci.yml` are pinned artifacts.

**Scale/Scope**: 9 documentation pages, ~11 cross-directory links to rewrite,
one new config, one new script, one new CI job, four wire-in edits.

## Constitution Check

*GATE: must pass before Phase 0 and re-checked after Phase 1.*

- **I. Code Quality** — No Go code changes. `scripts/docs-check.sh` is POSIX sh
  (not covered by gofmt/go vet/golangci-lint, which scope `internal cmd test`
  and `./...` Go packages). It will be kept POSIX-clean, `set -eu`, quoted. PASS.
- **II. Testing Standards (NON-NEGOTIABLE)** — No behavioral Go change, so no Go
  tests are added or weakened; the safety-critical surfaces (clock injection,
  DST, migrations, catch-up, goroutine termination, IPC access) are untouched.
  The feature's own verification is the docs-check plus full CI parity. PASS.
- **III. UX Consistency** — Not a CLI/API change. The site improves
  documentation UX; front matter does not alter in-repo rendering. PASS.
- **IV. Performance** — No hot path touched; dispatch-latency budget untouched.
  PASS.
- **V. Autonomous Build-Phase Execution** — Runs under autopilot; this feature
  traces to open issue #11. Pinned artifacts (`ci.yml`, `docs/INSTALL-windows.md`)
  are changed within scope and recorded as dated CHANGELOG decisions, surfaced
  at the single pre-push halt. The `/speckit-analyze` gate is not skipped. PASS.

**Result: PASS.** No violations; Complexity Tracking not required.

## Project Structure

### Documentation (this feature)

```text
specs/010-docs-site-pages/
├── plan.md              # This file
├── research.md          # Phase 0 — decisions (theme pin, grouping, checker rules)
├── data-model.md        # Phase 1 — page/config/pointer/check entities + front-matter schema
├── quickstart.md        # Phase 1 — how to validate locally + operator go-live steps
├── contracts/
│   └── docs-check.md    # Phase 1 — the link-check contract (inputs, rules, exit codes)
├── checklists/
│   ├── requirements.md  # from /speckit-specify
│   └── docs.md          # from /speckit-checklist
└── tasks.md             # Phase 2 — from /speckit-tasks
```

### Source changes (repository root)

```text
docs/
├── _config.yml          # NEW — just-the-docs remote theme, search, nav, plugins
├── README.md            # front matter (title: Home, nav_order: 1); ../ links → absolute
├── INSTALL-windows.md   # PINNED — front matter (parent: Installation); dated CHANGELOG
├── INSTALL-linux.md     # front matter (parent: Installation)
├── INSTALL-macos.md     # front matter (parent: Installation)
├── install.md           # NEW — tiny Installation section index (has_children)
├── cli.md               # front matter
├── gui-fields.md        # front matter; ../ link → absolute
├── cron.md              # front matter
├── test-scripts.md      # front matter; ../ links → absolute
└── build-autopilot.md   # front matter

scripts/
└── docs-check.sh        # NEW — POSIX-sh front-matter + link + pointer checker

.github/
├── workflows/ci.yml     # PINNED — NEW `docs` job; dated CHANGELOG
└── ISSUE_TEMPLATE/config.yml  # contact links → site URLs

README.md                # Documentation link to site; fix 0.6.0 → 0.7.0 drift
CONTRIBUTING.md          # NEW "Documentation" section (SSOT + pointer convention)
CHANGELOG.md             # Unreleased feature line + dated pinned-artifact decisions
```

**Structure Decision**: This is a documentation + tooling feature; there is no
application source tree to add. The layout above is the concrete change set.

## Key decisions (recorded per constitution principle V)

- **Serving: branch-based `/docs` + just-the-docs, not Hugo/MkDocs+Actions.**
  Chosen so the `docs/` Markdown is both the reviewable source and the served
  content, with zero deploy workflow and exactly one operator settings change —
  matching the operator's intent to "serve from the docs directory." Hugo and
  MkDocs (issue #11 alternatives) add a build pipeline and a second toolchain.
- **Theme pin `just-the-docs@v0.4.2`.** The last release before just-the-docs's
  practical Jekyll-4/dart-sass requirement; builds under GitHub Pages' libsass
  Jekyll 3.9. Migration path if the newest theme is wanted later: switch to a
  GitHub Actions build (a separate future feature). See research.md.
- **Link-checker strips `#fragments`, does not validate anchors.** Anchor
  validation is brittle (heading-slug rules) for little gain; file-existence +
  no-escape + front-matter presence catch the real drift. See research.md.
- **Installation grouping via a small `docs/install.md` parent page** with
  `has_children: true`; the three guides set `parent: Installation`. A real
  section index page reads better than nav-only nesting and gives the section a
  landing target. See research.md.

## Complexity Tracking

No Constitution Check violations. Section intentionally empty.
