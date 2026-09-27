# Implementation Plan: Adopt official brand 2.0.0

**Branch**: `codex/105-brand-kit-update` | **Date**: 2026-09-26 | **Spec**: [spec.md](spec.md)
**Planning record**: [#267](https://github.com/shruggietech/go-schedule/issues/267), Brand 2.0 adoption milestone.
**Input**: Owner request, official hosted go-schedule kit, BrandBuilder v2.4.0 formal release, and S105 research.

## Summary

Pin the verified complete formal archive in `brand/` and map its current assets into repository presentation paths. Add a standard-library Go import command with preflight checks and update the existing offline brand check to enforce archive identity, complete inventory, and copied assets. Update documentation and packaging inputs for the simplified mark, then run repository CI parity before committing and pushing the review branch.

## Technical Context

**Language/Version**: Go 1.26 for importer and verifier; existing Wails frontend and documentation assets.
**Primary Dependencies**: Go standard library, existing `scripts/brand-check`, official release archive and `SHA256SUMS`.
**Storage**: Versioned `brand/` archive, `brand/source.json` pin, `brand/repository-consumers.json` mapping, repository consumer copies.
**Testing**: Importer negative tests, brand-check negative tests, existing docs and automation checks, complete `scripts/verify.sh all` foreground run.
**Target Platform**: Existing Windows, macOS, Linux, web documentation, and Wails desktop surfaces.
**Project Type**: Brand asset import and repository build pipeline.
**Performance Goals**: No runtime scheduling path change; brand import remains a maintainer operation.
**Constraints**: Exact official archive bytes, no path traversal or symlink writes, no unpinned download, no unsupported platform claims, no rewrite of upstream generated prose.

## Constitution Check

- Principle I: Keep the importer and verifier small, explicit, and standard-library based; give error context and avoid unchecked filesystem operations.
- Principle II: Add failure-path tests for archive checksum, package identity, path handling, manifest inventory, and consumer integrity before running CI parity.
- Principle III: Align docs, repository artwork, frontend, and package inputs with the current identity. A site-relative manifest must keep icon URLs valid under the documentation base path.
- Principle IV: No scheduler hot path changes or benchmark requirement.
- Principle V: Run specification, clarification, checklist, plan, tasks, analysis, implementation, and verification in order. The owner explicitly requested a push before the approval halt, overriding the default pre-push boundary for this slice. Review branch and PR integration remain required; PR publication and merge await approval.
- Governance: One traceable outcome issue closes only when the functional import and update route are delivered; test results remain engineering evidence.

## Project Structure

```text
brand/                             Exact official ZIP plus repository-owned control and mapped presentation files
brand/source.json                 Formal release pin
brand/repository-consumers.json   Exact source-to-target map
scripts/brand-import/             Repeatable import command and tests
scripts/brand-check/              Offline identity, inventory, and consumer validation
docs/assets/                      Public brand and favicon copies
desktop/build/, desktop/frontend/ Desktop presentation copies
specs/105-brand-kit-update/       Slice artifacts
```

## Delivery Sequence

1. Record the formal release pin and a reviewed source-to-consumer migration map.
2. Write importer and checker failure-path tests, then implement archive preflight and exact selected-asset import.
3. Import the 2.0.0 archive, update copied assets and package paths, and remove stale presentation files.
4. Update repository and public brand guidance, including the distinction between official assets and repository-specific derivatives.
5. Run the analysis gate, resolve findings, run CI parity, format authored content, commit, and push the review branch.

## Decisions

- **Complete archive versus selected assets**: Keep the complete upstream ZIP so the repository retains an offline reference, manifest, migration contract, and exact recovery path. Extract only selected presentation assets to avoid publishing upstream text that conflicts with repository format rules.
- **Exact bytes versus local formatting**: Keep the ZIP byte-for-byte; repository-authored control files and published prose remain subject to `github-format`.
- **Packaging compatibility paths**: Keep repository packaging paths where build tooling requires them, but map each byte-identical copy to the new kit and document any local derivative. This avoids a large unrelated packaging rewrite.
- **Documentation manifest**: Derive site-relative icon URLs from the official web manifest for GitHub Pages under `/go-schedule`; record the transform in the consumer contract and test it.

## Complexity / Deviation

The official archive includes an em dash and width-wrapped prose, so its bytes remain inside the ZIP while repository-authored source files comply with formatting rules. Source asset validation does not establish attended native appearance. The user-authorized push precedes this slice's approval halt; pull-request publication and merge remain pending.
