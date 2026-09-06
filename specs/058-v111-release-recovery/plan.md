# Implementation Plan: v1.1.1 Release Recovery

**Branch**: `codex/058-v111-release-recovery` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/058-v111-release-recovery/spec.md`, corrective PR [#143](https://github.com/shruggietech/go-schedule/pull/143), and release issue [#140](https://github.com/shruggietech/go-schedule/issues/140).

## Summary

Preserve the existing unpublished v1.1.0 tag as immutable history, prepare v1.1.1 as the first public release containing the accumulated feature set and watcher correction, update the existing release coordinator rather than duplicating it, and require successful main-branch CI for the exact reviewed release commit before staging.

## Technical Context

**Language/Version**: Go 1.25 plus Markdown, YAML, PowerShell, and POSIX shell release surfaces

**Primary Dependencies**: GitHub Actions, GitHub Releases, WiX 6.0.2, and existing release-gate tooling

**Storage**: Git objects, GitHub release metadata, workflow artifacts, and evidence archives

**Testing**: `go run ./scripts/github-format`, release integration tests, `scripts/verify.sh all`, hosted CI, candidate validation, evidence validation, and checksum audit

**Target Platform**: GitHub-hosted Linux, macOS, and Windows builders plus attended Windows installation

**Project Type**: Cross-platform daemon, CLI, and desktop release operation

**Performance Goals**: Not applicable; release correctness and provenance dominate elapsed time

**Constraints**: Immutable historical tag, no artifact rebuild during promotion, exact-commit main CI, no public release before qualification, no hard-wrapped GitHub prose, and no Unicode em dash

**Scale/Scope**: One corrective version boundary, four release highlights, seven distributable packages, one candidate manifest, one attended evidence archive, and one checksum inventory

## Constitution Check

- Independent reasoning: PASS. Moving v1.1.0 would corrupt immutable history, so recovery advances to v1.1.1.
- Scope proportionality: PASS. Product code remains unchanged because PR #143 already delivered and verified the required corrections.
- Encoding and output integrity: PASS. Repository text remains UTF-8 without BOM and publication formatting is mechanically checked.
- Chronological planning: PASS. Recovery preparation precedes review, old-draft retirement, tagging, staging, qualification, promotion, and audit.
- GitHub planning: PASS. Existing issue #140 and milestone #4 remain authoritative and are revised instead of duplicated.

## Project Structure

```text
.github/release-notes/v1.1.1.md
.specify/feature.json
CHANGELOG.md
README.md
specs/058-v111-release-recovery/
├── checklists/requirements.md
├── contracts/publication.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
specs/README.md
```

**Structure Decision**: Release recovery changes only versioned identity, release copy, planning records, and Spec Kit evidence. Existing build and promotion workflows remain unchanged because PR #143 already added the exact-commit CI preflight.

## Complexity Tracking

No constitution violation requires justification.
