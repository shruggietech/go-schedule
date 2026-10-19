# Implementation Plan: Public website and project discovery

**Branch**: `codex/106-public-web-presence` | **Date**: 2026-10-19 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #278, #279, #280, #281, #9 and #282 as one static public presence. Build a Next.js export with Fumadocs from canonical `docs/` Markdown; add the branded landing and sourced comparison; replace Jekyll-specific checks; enforce comparison review in staging and both promotion paths. Prepare Cloudflare/Pages changes without live cutover before owner review.

## Technical Context

**Language/Version**: TypeScript/JavaScript, Node 26 and existing Go release tooling.

**Primary Dependencies**: Next 16.4.2 (official registry confirmed), React 19.2.8, Fumadocs core/UI 16.17.6 and MDX 15.4.7 (current compatible set validated by production export), Tailwind/PostCSS 4.2.4, TypeScript 5.9.3. Roles are static compilation, docs navigation/search/highlighting and shared accessible components; retain exact lockfile and licenses.

**Storage**: Canonical Markdown and versioned JSON snapshots/reviews; browser appearance preference only.

**Testing**: Node content/export contract tests, Go release-review regression tests, production browser checks and canonical `sh scripts/verify.sh all`.

**Target Platform**: Static GitHub Pages at go-schedule.io, desktop/mobile browsers, no runtime server.

**Project Type**: Documentation/marketing site and release tooling.

**Performance Goals**: Build documents once; static search; defer Mermaid until needed.

**Constraints**: UTF-8 without BOM, no em dashes/hard-wrapped Markdown, official brand, no application behavior changes, no publication before review.

**Scale/Scope**: Existing top-level docs and all assets; seven competitors plus go-schedule; six child issues.

## Constitution Check

- I: Small content/release helpers reject ambiguous routes and data; no custom docs renderer.
- II: Contract tests precede implementation; preserve all existing safety-critical gates and core coverage.
- III: Actionable build/review errors and coherent reader navigation/source identity.
- IV: No scheduler hot-path changes; static search and deferred diagrams bound browser work.
- V: Spec Kit sequence and analysis before implementation; local review branch only.
- Post-design: PASS. Live #278 outcome remains pending publication, not falsely completed from local tests.

## Project Structure

```text
site/app/                 # home, docs, comparison, static search
site/components/          # shared chrome and diagram renderer
site/lib/                 # source loader and metadata
site/scripts/             # preparation and export audit
site/tests/               # content/browser contracts
site/data/                # comparison snapshots/reviews
scripts/comparison-check/ # structural/release guard and tests
docs/website.md           # authoring/release/domain handoff
.github/workflows/pages.yml
specs/106-public-web-presence/
```

## Design direction

Approved palette: background #071014, panel #0D171C, raised #13232A, blue #58A6FF, accessible blue #1160C6 and white; orange for deliberate emphasis. Space Grotesk headlines, Geist body and Geist Mono commands use local licensed fonts. The memorable element is a readable task-to-trigger-to-result workflow beside a real CLI example. Keep supporting content quiet, vertical and left aligned with generous reading widths. Light mode uses white surfaces/deep blue. No invented metrics or repeated decorative cards. This is specific to scheduling and the approved cursor mark.

## Dependencies and delivery

The content pipeline precedes home/comparison rendering. README and research are independent. One offline Go command enforces comparison review in existing staging preflight and before the first mutation in both promotion workflows. Preserve release-state dependencies; do not absorb unrelated #262. Cloudflare read-only inventory needs owner authentication. DNS/Pages mutations, deployment, branch push and PR stay at the requested review boundary. Optional agent-context hooks are skipped; the required CLAUDE plan reference is updated directly.
