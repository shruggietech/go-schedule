# Research: Documentation site on GitHub Pages + README consolidation

Phase 0 decisions. Each resolves a design unknown for feature 010.

## D1 — Serving mechanism: branch-based `/docs` + just-the-docs

**Decision**: Serve the site from the `/docs` folder on `main` via GitHub Pages'
branch-based publishing (Pages source = "Deploy from a branch → `main` →
`/docs`"). GitHub's bundled Jekyll (github-pages gem) builds it. Theme is
just-the-docs applied with `remote_theme`.

**Rationale**: Keeps the `docs/` Markdown as *both* the reviewable repository
source and the served content — no built artifacts committed, no deploy
workflow to maintain, and exactly one operator settings change to go live. This
matches the operator's stated intent to "serve from the docs directory."
just-the-docs supplies the three capabilities issue #11 asks for out of the box:
client-side search, a spanning sidebar, and next/previous.

**Alternatives considered**:
- *Hugo + `actions/deploy-pages`* (issue #11's tentative default): a Go-toolchain
  build, but it deploys an artifact rather than serving `/docs`, adds a pinned
  workflow, and commits Hugo config/theme. More moving parts for no gain here.
- *MkDocs Material + Actions*: best-in-class docs UX (incl. `mike` versioning)
  but adds a Python toolchain and an Actions deploy. Heaviest; diverges most
  from the Go-only toolchain.
- *Plain GitHub Markdown (status quo)*: free, but caps discoverability — the
  exact problem #11 exists to solve.

## D2 — just-the-docs version pin for branch-based Pages

**Decision**: `remote_theme: just-the-docs/just-the-docs@v0.4.2`.

**Rationale**: GitHub Pages' branch-based build runs the github-pages gem
(Jekyll 3.9.x with the libsass sass-converter). While just-the-docs gemspecs
still declare `jekyll >= 3.8.5` even at 0.10.x, the theme's SCSS migrated to
dart-sass module syntax (`@use`) which libsass cannot compile; the theme's own
MIGRATION notes flag from v0.4.0 that "future releases … may require the use of
Jekyll 4." v0.4.2 is the newest release of the last line that builds cleanly on
the classic github-pages workflow, and it already provides search + nav +
next/prev. `remote_theme` (via `jekyll-remote-theme`) bypasses the Pages theme
allowlist, so the pin is honored.

**Migration path** (out of scope): to adopt the newest just-the-docs (newer
search, callouts, native dark mode toggle) later, switch to a GitHub Actions
build with Jekyll 4 + `actions/deploy-pages`. That is a separate future feature;
it changes the Pages source from "branch" to "Actions."

**Alternatives considered**: `@v0.3.3` (older, fewer features, no benefit over
0.4.2 on this toolchain); newest `@v0.10.x` (fails to build under libsass).

## D3 — Link-checker scope and rules

**Decision**: `scripts/docs-check.sh` validates, for every `docs/*.md`:
1. YAML front matter is present and contains `title:` and `nav_order:` (the
   home page and section index may use `nav_order`; all pages carry `title`).
2. Every Markdown link that is **not** `http(s):` and **not** a pure `#fragment`
   is checked: strip any trailing `#fragment`, resolve the path relative to the
   file's directory, and assert the target exists on disk.
3. No link resolves to a path **outside** `docs/` (no surviving `../` escape).
4. The pointer README `test/scripts/README.md` — its `docs/*.md` target exists.

It exits non-zero listing each offending `file: link`. No network; pure `sh` +
coreutils (`grep`, `sed`), mirroring `scripts/coverage-gate.sh`'s style.

**Rationale**: File-existence + no-escape + front-matter presence catch the drift
that actually breaks the published site (404s, unlisted pages, links that escape
the served root). Anchor (`#fragment`) validation is deliberately excluded: it
requires replicating Jekyll/kramdown heading-slug rules and is brittle for
little value. Skipping `http(s)` keeps the check network-free and deterministic
in CI.

**Alternatives considered**: a third-party link checker (e.g. `lychee`) — adds a
binary/action dependency and network calls; rejected for a 9-page doc set. A Go
test — heavier than a shell script for a docs concern and would pull docs into
the Go test graph.

## D4 — Installation grouping

**Decision**: Add a small `docs/install.md` section page with
`title: Installation`, `nav_order: 2`, `has_children: true`. The three guides
(`INSTALL-windows/linux/macos.md`) set `parent: Installation` with their own
`nav_order` (1/2/3) for platform ordering.

**Rationale**: just-the-docs renders a real, collapsible "Installation" section
with a landing page a reader (and the README) can link to, which reads better
than nav-only nesting and gives the grouping a stable URL. The section page is
a few lines pointing to the three guides — no content duplication.

**Alternatives considered**: nav-only nesting without a section page (no landing
target; the parent label isn't clickable to anything useful); flattening the
three guides at top level (fails the "grouped under one Installation section"
requirement, FR-004 / SC-003).

## D5 — Home page and in-repo rendering

**Decision**: Keep `docs/README.md` as the site home (`title: Home`,
`nav_order: 1`); do **not** add a duplicate `docs/index.md`. GitHub Pages serves
a folder's `README.md` as its directory index when no `index.*` is present.

**Rationale**: One index file, serving both the repo folder view and the site
root. YAML front matter is rendered invisibly by GitHub's Markdown view, so the
in-repo reading experience is unchanged (edge case in spec).

## D6 — Page URL shape (for repointed inbound links)

**Decision**: Use default Jekyll permalinks (no `permalink: pretty`), so
`docs/test-scripts.md` publishes at
`https://shruggietech.github.io/go-schedule/test-scripts.html`. The
issue-template "verify your install" link points there; the "installation help"
link points at the site root `https://shruggietech.github.io/go-schedule/`.

**Rationale**: Deterministic output paths with no permalink configuration to
reason about; `jekyll-relative-links` rewrites in-page `.md` links to the same
`.html` targets, so on-site and inbound links agree.

## D7 — Cross-directory link rewrite targets

The `../`-escaping links to rewrite to
`https://github.com/shruggietech/go-schedule/blob/main/<path>` (verified present
by exploration):

- `docs/README.md`: `../README.md`, `../CONTRIBUTING.md`, `../SECURITY.md`,
  `../CHANGELOG.md`, `../.specify/memory/constitution.md`,
  `../specs/001-task-scheduler/spec.md`
- `docs/gui-fields.md`: `../specs/001-task-scheduler/contracts/cli.md`
- `docs/test-scripts.md`: `../.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1`,
  `../specs/006-maintainer-test-scripts/`, `../specs/006-maintainer-test-scripts/data-model.md`,
  `../test/scripts/lib/sqlite-manifest.json`

Note the `blob/main` form works for files and the `tree/main` form for
directories; the `specs/006-maintainer-test-scripts/` directory link uses
`tree/main`. Within-`docs/` `.md` links stay relative.

## Open items

None. All Phase 0 unknowns resolved; no `[NEEDS CLARIFICATION]` remains.

## Sources

- [just-the-docs (theme)](https://github.com/just-the-docs/just-the-docs)
- [just-the-docs MIGRATION notes](https://github.com/just-the-docs/just-the-docs/blob/main/MIGRATION.md)
- [Using just-the-docs as a remote theme (example)](https://github.com/pmarsceill/jtd-remote)
- [RubyGems: just-the-docs versions](https://rubygems.org/gems/just-the-docs)
