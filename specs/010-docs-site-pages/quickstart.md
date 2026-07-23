# Quickstart: validate the documentation site + consolidation

How to prove feature 010 works end to end, and the one-time operator go-live.

## Prerequisites

- Repository checked out on `main`.
- POSIX shell (`sh`) available. No Ruby/Jekyll needed locally — the site build
  runs on GitHub Pages.

## 1. Documentation integrity (the automated gate)

```sh
sh scripts/docs-check.sh
```

Expected: exit `0` with a success summary naming the number of pages checked.
Contract and failure scenarios: [contracts/docs-check.md](contracts/docs-check.md).

Spot-check the negative paths (revert each edit afterward):
- Delete a `nav_order:` line from any `docs/*.md` → the check fails naming that
  file.
- Add `[x](../nope.md)` to a `docs/*.md` → the check fails on the escape/broken
  link.
- Add `[x](does-not-exist.md)` inside `docs/` → the check fails on the broken
  link.

## 2. Front matter & nav (manual inspection)

- Every `docs/*.md` opens with a `---` front-matter block carrying `title` and
  `nav_order` (see [data-model.md](data-model.md) for the ordering table).
- `docs/README.md` is `title: Home, nav_order: 1`; there is no `docs/index.md`.
- `docs/install.md` has `has_children: true`; the three `INSTALL-*.md` guides
  set `parent: Installation`.
- Viewing any `docs/*.md` in the GitHub repo file view still renders cleanly
  (front matter is hidden by GitHub's Markdown view).

## 3. No relative link escapes `docs/`

```sh
grep -rEn '\]\(\.\./' docs/*.md || echo "no ../ escapes — good"
```

Expected: `no ../ escapes — good`. Every reference outside `docs/` is an absolute
`https://github.com/shruggietech/go-schedule/...` URL.

## 4. Inbound references point at the site

- `README.md` has a prominent Documentation link to
  `https://shruggietech.github.io/go-schedule/`, and the quick-start version
  reads `0.7.0` (matching the release badge).
- `.github/ISSUE_TEMPLATE/config.yml` contact links point at the site (root and
  `test-scripts.html`), not raw repository file URLs.
- `CONTRIBUTING.md` has a Documentation section stating `docs/` is the single
  source of truth and subdirectory READMEs are pointers.

## 5. Full CI parity (before the pre-push halt)

Run the project's CI-parity gates in the foreground, watched to completion (see
the `go-schedule-verify` skill / CLAUDE.md). No Go changed, so all Go gates are
expected green; the `-race` gate is CI-only on a machine without a C toolchain
and is reported as such. The new `docs` job mirrors `sh scripts/docs-check.sh`.

## 6. Operator go-live (after push + release authorization)

Two settings changes, both outside what the repository content can do itself:

1. **Repository → Settings → Pages → Build and deployment →**
   **Source: "Deploy from a branch" → Branch: `main`, Folder: `/docs` → Save.**
2. **Repository main page → About (gear) → Website:**
   **`https://shruggietech.github.io/go-schedule/`.**

Then confirm on the published site: the home page renders, search returns
results, the sidebar lists every page, the three install guides sit under one
Installation section, and no in-page link 404s.
