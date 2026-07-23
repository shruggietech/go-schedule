# Feature Specification: Documentation site on GitHub Pages + README consolidation

**Feature Branch**: `010-docs-site-pages`

**Created**: 2026-07-23

**Status**: Draft

**Input**: User description: "Publish docs/ as a GitHub Pages documentation site and consolidate in-repo READMEs against it so nothing drifts (closes #11)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Read the documentation as a searchable, navigable site (Priority: P1)

A person evaluating or installing go-schedule opens the project's documentation
in a browser and finds a landing page, a sidebar listing every document, and a
search box. They type a question ("how do I check whether the daemon is
running?"), jump straight to the answer, and move between related pages with
next/previous and cross-links — without cloning the repository or guessing which
of eight files holds the answer.

**Why this priority**: This is the whole point of the feature and of issue #11.
Without a published, navigable, searchable site there is no user-facing value;
everything else is plumbing in support of this.

**Independent Test**: Point the Pages source at the `docs/` folder, load the
published URL, confirm the landing page renders, the sidebar lists every page,
search returns results, install guides are grouped together, and no in-page link
leads to a 404. Deliverable and demonstrable on its own.

**Acceptance Scenarios**:

1. **Given** the published site, **When** a reader loads the site root, **Then**
   they see the documentation index as the home page with a navigation sidebar
   listing every document.
2. **Given** the published site, **When** a reader searches for a term that
   appears in a page body, **Then** the matching page is offered as a result.
3. **Given** the published site, **When** a reader opens any documentation page
   and follows any in-page link, **Then** the link resolves to a real
   destination (another site page, or an absolute repository URL for content
   that lives outside the documentation set).
4. **Given** the published site, **When** a reader looks for installation help,
   **Then** the three per-platform install guides appear grouped under a single
   Installation section.

---

### User Story 2 - Documentation stays the single source of truth and cannot silently drift (Priority: P2)

A maintainer edits documentation confident that there is exactly one place each
fact lives. In-repo READMEs in subdirectories are thin pointers to the
documentation, not copies that rot. If an edit breaks an internal link, moves a
document a pointer depends on, or omits the metadata a page needs to appear on
the site, an automated check fails and names the problem before it ships.

**Why this priority**: The operator explicitly asked for consolidation so that
"any required in-place directory README's are kept in sync." A published site
that quietly drifts out of step with the repo is worse than no site. This
depends on Story 1 existing but is independently valuable and independently
testable.

**Independent Test**: Run the documentation check script against the repo; it
passes on a consistent tree. Introduce a broken link or a missing pointer
target, re-run, and confirm the check fails and names the offending file. The
check runs as its own CI job.

**Acceptance Scenarios**:

1. **Given** a consistent documentation tree, **When** the documentation check
   runs, **Then** it exits successfully.
2. **Given** a documentation page with a relative link to a file that does not
   exist, **When** the check runs, **Then** it fails and identifies the page and
   the broken link.
3. **Given** a pointer README whose target document has been removed, **When**
   the check runs, **Then** it fails and identifies the stale pointer.
4. **Given** a documentation page missing its required front matter, **When** the
   check runs, **Then** it fails and identifies the page.

---

### User Story 3 - Every reference to the docs points at the site, not raw files (Priority: P3)

Someone who lands on the repository README, or opens the "new issue" chooser,
follows a documentation link and arrives at the published site rather than a raw
file listing. Support answers become durable site URLs.

**Why this priority**: Valuable polish that makes the site discoverable and the
references durable, but the site delivers its core value (Stories 1 and 2) even
before every inbound link is repointed.

**Independent Test**: Inspect the root README and the issue-form contact
configuration; confirm their documentation links target the site URL. Confirm
the README no longer shows a stale version number.

**Acceptance Scenarios**:

1. **Given** the repository README, **When** a reader looks for documentation,
   **Then** a prominent link takes them to the published site.
2. **Given** the "new issue" chooser, **When** a reader opens the installation-
   help or verify-your-install contact link, **Then** it opens the corresponding
   site page rather than a raw repository file URL.
3. **Given** the repository README, **When** a reader reads the quick-start
   walkthrough, **Then** the version shown matches the current release badge.

---

### Edge Cases

- **A documentation page links to content outside the documentation set**
  (repository README, CONTRIBUTING, the constitution, a spec, a test script).
  On a site rooted at the documentation folder those targets are not published,
  so such links must resolve to absolute repository URLs rather than site-
  relative paths that would 404.
- **A page is added without navigation metadata.** It would appear unordered or
  invisibly in the sidebar; the documentation check must catch a page that lacks
  the required front matter.
- **The theme requires a newer site builder than the hosting provides.** The
  chosen theme version must be compatible with the builder the branch-based
  hosting actually runs, or the site fails to build; the compatible version is a
  fixed, recorded choice.
- **A reader browses a documentation file directly in the repository** (not on
  the site). Added navigation metadata must not degrade that in-repo reading
  experience.
- **The documentation check encounters an external (http/https) link.** It must
  not fail the build on network conditions; external links are out of scope for
  the on-disk check.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The documentation set MUST be publishable as a website served
  directly from the repository's documentation folder on the default branch, so
  that the documentation source and the served content are one and the same and
  no built artifacts are committed.
- **FR-002**: The site MUST provide full-text search across the documentation
  set, a navigation sidebar spanning all pages, and next/previous movement
  between pages.
- **FR-003**: The site's home page MUST be the existing documentation index, and
  the documentation folder's index file MUST continue to render correctly when
  browsed directly in the repository.
- **FR-004**: Every documentation page MUST carry the metadata the site needs to
  place it in navigation (a human title and an explicit ordering), and the three
  per-platform install guides MUST be grouped together under a single
  Installation section.
- **FR-005**: No documentation page may contain a repository-relative link that
  escapes the documentation folder; every such reference MUST be an absolute
  repository URL so it resolves both on the site and from the repository.
- **FR-006**: Links between documentation pages MUST continue to work on the
  published site without manual per-link rewriting.
- **FR-007**: An automated documentation check MUST fail when a documentation
  page has a broken on-disk link, when a pointer README references a missing
  document, or when a documentation page lacks its required navigation metadata;
  it MUST ignore external (http/https) links and MUST run without a network
  dependency.
- **FR-008**: The documentation check MUST run in continuous integration as its
  own job on every push and pull request, on the same triggers as the existing
  checks.
- **FR-009**: In-repo subdirectory READMEs that would otherwise duplicate
  documentation MUST remain thin pointers to the documentation; the contributor
  guide MUST state that the documentation folder is the single source of truth,
  that subdirectory READMEs are pointers, and how the site is served.
- **FR-010**: The repository README and the issue-form contact links MUST point
  at the published site rather than raw repository file URLs, and the README's
  quick-start version reference MUST match the current release badge.
- **FR-011**: Changes to pinned artifacts required by this feature MUST each be
  recorded as a dated decision in the changelog, and the changelog MUST carry an
  unreleased feature entry for the site.
- **FR-012**: The feature MUST NOT modify the scheduler's Go code; the existing
  code quality, test, race, coverage, and benchmark gates MUST remain green.

### Key Entities *(include if feature involves data)*

- **Documentation page**: one Markdown file in the documentation folder. Carries
  a title and a navigation order; may declare a parent section (the install
  guides declare the Installation parent). Its body is the source of truth for
  the topic it covers.
- **Site configuration**: the single file that turns the documentation folder
  into a themed, searchable, navigable site (theme, search, plugins, site
  identity). Not committed as built output — it configures a build the host
  performs.
- **Pointer README**: a subdirectory README whose role is to redirect a reader
  to the canonical documentation page rather than to hold content.
- **Documentation check**: the script and its CI job that enforce link
  integrity, pointer validity, and page metadata across the documentation set.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A reader can reach any documentation topic from the site's home
  page in at most two clicks (sidebar section → page) or one search.
- **SC-002**: 100% of in-page documentation links resolve — zero 404s across the
  published site.
- **SC-003**: Every documentation page appears in the site navigation in a
  deliberate order, and all three install guides appear under one Installation
  section.
- **SC-004**: A broken internal link, a stale pointer, or a page missing its
  navigation metadata is caught automatically and fails the build, rather than
  reaching the published site.
- **SC-005**: Every inbound documentation reference the project controls (README,
  issue-form contact links) resolves to a site URL, and the README shows a
  single consistent version.
- **SC-006**: The publication requires exactly one repository settings change by
  the operator (point Pages at the documentation folder) with no ongoing deploy
  workflow to maintain.

## Assumptions

- The hosting provider serves a site directly from a documentation folder on the
  default branch and builds it with its bundled site generator; the chosen theme
  is pinned to a version compatible with that generator (a newer theme line that
  would require a separate build pipeline is deliberately not used).
- The documentation set is user-facing documentation, not a mirror of the
  repository: the repository README and the specification tree stay where they
  are and are not migrated onto the site.
- There is essentially no literal content duplication between the documentation
  and the subdirectory READMEs today (the one subdirectory README is already a
  pointer), so a pointer-plus-check convention is sufficient and no
  content-copying/generation machinery is warranted.
- Versioned documentation, a custom domain, and branding assets (favicon, social
  preview image) are out of scope for this first publish; the branding assets
  depend on a separate branding-package effort.
- The operator performs the one-time hosting settings change and sets the
  repository "Website" field after the change ships; both are outside what the
  repository content can do for itself.
- The site is served from the default branch on every push, so the published
  documentation tracks the default branch (the versioning question is deferred
  rather than answered here).
