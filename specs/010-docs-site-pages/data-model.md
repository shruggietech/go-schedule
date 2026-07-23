# Data Model: Documentation site + README consolidation

The "entities" here are documents and configuration, not persisted records.

## Documentation page

One Markdown file in `docs/`. On the site it becomes a navigable, searchable
page; in the repo it stays a readable file.

**Front-matter schema** (YAML block at the very top, delimited by `---`):

| Field         | Required | Applies to                    | Meaning |
|---------------|----------|-------------------------------|---------|
| `title`       | yes      | every page                    | Human page title shown in nav and tab. |
| `nav_order`   | yes*     | every top-level page + children | Integer ordering within its nav level. |
| `parent`      | when nested | the three install guides    | `Installation` — places the page under the section. |
| `has_children`| yes      | `install.md` only             | `true` — marks the Installation section parent. |

\* The home page carries `nav_order: 1`; the section index and all leaf pages
carry an explicit `nav_order`. The link-check requires `title` on every page and
`nav_order` on every page.

**Nav layout & ordering** (top level unless nested):

| File                    | title                         | nav_order | parent |
|-------------------------|-------------------------------|-----------|--------|
| `README.md`             | Home                          | 1         | —      |
| `install.md`            | Installation                  | 2 (has_children) | — |
| `INSTALL-windows.md`    | Windows                       | 1         | Installation |
| `INSTALL-linux.md`      | Linux                         | 2         | Installation |
| `INSTALL-macos.md`      | macOS                         | 3         | Installation |
| `cli.md`                | CLI reference                 | 3         | —      |
| `gui-fields.md`         | GUI field reference           | 4         | —      |
| `cron.md`               | Cron interoperability         | 5         | —      |
| `test-scripts.md`       | Maintainer test scripts       | 6         | —      |
| `build-autopilot.md`    | Build-phase autopilot         | 7         | —      |

**Invariant**: adding a page without valid front matter fails the docs-check
(FR-007). Front matter does not alter GitHub's in-repo Markdown rendering.

## Site configuration (`docs/_config.yml`)

The single file that turns `docs/` into a themed, searchable site. Not built
output — it configures the host's build.

Key settings:
- `title`, `description` — site identity (feeds `jekyll-seo-tag`).
- `url: https://shruggietech.github.io`, `baseurl: /go-schedule` — project-pages
  base path.
- `remote_theme: just-the-docs/just-the-docs@v0.4.2` — pinned theme (see D2).
- `search_enabled: true` — client-side search.
- `aux_links` — a persistent link back to the GitHub repository.
- `plugins: [jekyll-remote-theme, jekyll-relative-links, jekyll-seo-tag, jekyll-sitemap]`
  — all on the GitHub Pages allowlist.
- `exclude` — any non-page files that must not be published (e.g. `_config.yml`
  is auto-excluded; add others only if needed).

## Pointer README

A subdirectory README whose role is to redirect to the canonical doc, not to
hold content.

- `test/scripts/README.md` — already a pointer to `docs/test-scripts.md`;
  unchanged except that the docs-check now guards its target's existence.
- **Convention** (documented in CONTRIBUTING): `docs/` is the single source of
  truth; a subdirectory README is a signpost to the relevant `docs/` page and
  must not duplicate documentation content.

## Documentation check (`scripts/docs-check.sh` + CI job)

Enforces the invariants above across the doc set. Full contract in
[contracts/docs-check.md](contracts/docs-check.md).

- **Inputs**: the `docs/` tree + the known pointer README(s). No arguments, no
  network.
- **Rules**: front-matter presence (`title` + `nav_order`); every non-`http(s)`,
  non-pure-`#fragment` link resolves on disk; no link escapes `docs/`; pointer
  targets exist.
- **Output**: `0` when clean; non-zero listing `file: offending-link` otherwise.
- **CI**: a `docs` job in `.github/workflows/ci.yml` runs `sh scripts/docs-check.sh`
  on push/PR to `main`, in parallel with the existing jobs.

## Inbound references (repointed, not new entities)

- `README.md` — a prominent Documentation link to the site; quick-start version
  reference corrected `0.6.0` → `0.7.0` (the release badge is authoritative and
  auto-bumped by `release.yml`; its line is not edited here).
- `.github/ISSUE_TEMPLATE/config.yml` — two contact links point at site URLs
  (root and `test-scripts.html`) instead of raw repository file URLs.
