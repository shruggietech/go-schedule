# Implementation Plan: S119 v1.6.0 release execution

**Branch**: `codex/119-release-execution` | **Date**: 2026-10-26 | **Spec**: [spec.md](spec.md)

## Summary

Execute the existing exact-source release pathway, retain actual candidate/public identities and prepare a reviewed discovery reconciliation PR. Original #355/#356 causes and failed captures remain separate and open. Tag/staging source is owner-approved 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b; S119 source changes never silently replace it.

## Technical Context

**Language/Version**: Existing Go 1.26 release validators, Node 26 website tooling, Markdown/JSON records.

**Primary Dependencies**: Installed Spec-Kit 0.11.3, GitHub Actions release/promotion/Pages, current Next.js/Fumadocs and Playwright.

**Storage**: Authenticated GitHub release/assets; versioned source discovery and immutable comparison data; ignored local receipts.

**Testing**: Existing source/comparison/format guards, candidate-manifest validation, actual browser regression and all eight canonical foreground gates before source commit.

**Target Platform**: Existing Linux/macOS archives and Windows amd64 MSI; static website.

**Project Type**: Release execution and source discovery reconciliation.

**Performance Goals**: Existing declared completion/resource allowances; no new application performance contract.

**Constraints**: Exact tag authorization received; public promotion awaits concrete-asset approval. No installed-system/session/DNS changes. Stop source PR unmerged and at most two Codex rounds.

**Scale/Scope**: Eight staged assets, nine promoted assets including checksums, canonical current-status paragraphs, selected published version and one demonstrated stale browser assertion.

## Constitution Check

Pre-research PASS: Existing architecture, security, migrations and test assertions are retained. No new dependency or qualification policy. Exact tag/staging authority is explicit, public promotion separate.

Post-design PASS: Reuse existing release and website mechanisms. Any demonstrated browser correction requires failing-before/passing-after acceptance and its own issue. Pinned Windows documentation status change receives a dated changelog decision. All eight gates precede every source commit. PR publication and release operations are recorded as external delivery receipts rather than unchecked implementation tasks. No complexity exception.

## Project Structure

- `specs/119-release-execution/`: spec, design, source tasks, delivery and actual verification/publication identities.
- `site/data/published.json`: actual public default after promotion; immutable tagged comparison/notes/original review stay unchanged; only the current main changelog review digest/date is refreshed for reconciled publication prose.
- `README.md`, current `docs/*.md`, `CHANGELOG.md`: current release status and preserved preparation history.
- `site/e2e/site.spec.ts`: demonstrated release-aware browser assertion correction.
- `specs/README.md`, `CLAUDE.md`, `.specify/feature.json`: aligned slice lifecycle.
- `.test-temp/s119/`: private source/candidate/download/hash/CI/review receipts, no credentials.

## Execution strategy

Complete blocking analysis, then existing read-only source preflight and issue kickoff. Execute the approved immutable tag/staging externally and retain complete receipts. Request promotion approval against the actual complete candidate. After authorized public promotion, reconcile source discovery with the real published state; reproduce/correct the release-specific browser defect and track it independently. Run targeted browser acceptance then freeze finished source for all eight gates. Finalize factual task/lifecycle metadata, commit, push and open the authorized PR. Monitor exact-head CI and every review, maximum two Codex rounds; stop unmerged.
