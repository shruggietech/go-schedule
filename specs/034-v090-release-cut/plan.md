# Implementation Plan: v0.9.0 Release Cut

**Branch**: `codex/034-v090-release-cut` | **Date**: 2026-08-30 |
**Spec**: [spec.md](spec.md)

## Summary

Prepare the reviewed v0.9.0 boundary, replace exhaustive generated GitHub
release copy with tag-specific concise highlights, enforce that presentation
contract offline, and preserve the existing package and checksum workflow. The
implementation ends at the normal pre-publication halt. Pull-request
publication and the later v0.9.0 tag remain separately authorized workflow
steps, with issue #79 open until the hosted release is fully audited.

## Technical Context

**Language/Version**: GitHub Actions YAML, POSIX shell, Markdown, Go 1.25.0
verification suite

**Primary Dependencies**: `softprops/action-gh-release@v3`, GitHub Releases,
existing Spec-Kit lifecycle and offline automation checker

**Storage**: Versioned repository documents plus hosted tag, release, and assets

**Testing**: Shell fixtures, workflow contract checks, eight canonical gates,
hosted pull-request CI, and post-tag GitHub readback

**Target Platform**: GitHub release workflow producing Windows, Linux, and macOS
artifacts

**Project Type**: Release automation and documentation slice

**Performance Goals**: No runtime performance change; release notes remain
scannable as four to six highlights

**Constraints**: UTF-8 without BOM, LF, concise public notes, complete detailed
changelog, no generated notes, no manual install session, no tag before explicit
authorization

**Scale/Scope**: One minor release, one detailed changelog boundary, one
tag-specific highlight file, existing package matrix, one tracked issue

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | The workflow uses one explicit notes source and fails when the tag-specific file is absent. |
| II. Testing Standards | PASS | Negative fixtures precede the release-policy checker change; all canonical gates remain required. |
| III. UX Consistency | PASS | Public notes are concise, version-specific, and link to one authoritative detailed record. |
| IV. Performance | PASS | No runtime path changes; release-copy size is deliberately bounded. |
| V. Autonomous Execution | PASS | S034 follows the complete command sequence and halts before branch publication; tagging remains separately authorized. |

### Post-design re-check

The design changes a pinned release workflow, so the decision is recorded in
the v0.9.0 changelog. Existing artifact jobs, supported platforms, checksums,
and quality gates remain intact. No constitution exception is required.

One tooling-sequence deviation is recorded: the checklist prerequisite script
requires `plan.md`, while project autopilot requires checklist before plan. The
requirements checklist was therefore produced directly from the completed spec.
This does not skip or weaken the checklist gate and does not modify Spec-Kit.

## Architecture and Decision Log

### Select release notes by exact tag

Store curated notes at `.github/release-notes/<tag>.md` and configure the
release action with a dynamic `body_path`. A missing tag-specific file then
fails publication instead of reusing stale notes. A hard-coded v0.9.0 path was
rejected because it would publish the wrong copy on the next release.

### Disable generated release notes

Set `generate_release_notes` to false. Generated commit and pull-request lists
conflict with the maintainer's highlights-only requirement and duplicate the
full changelog. Retaining generated notes and trying to trim them afterward was
rejected because it makes the public contract dependent on hosted formatting.

### Keep detail in the changelog

Move the current Unreleased history under a dated v0.9.0 heading and create a
new empty Unreleased section. The highlight file links to the versioned
changelog at the v0.9.0 tag. Duplicating detailed entries into release notes was
rejected because the two records would drift.

### Guard the release-copy contract offline

Extend the existing automation checker and its fixtures rather than create a
ninth verification gate. The checker requires a dynamic notes path, rejects
generated notes, and validates each notes file for a Highlights heading, four
to six bullets, one versioned changelog link, and prohibited exhaustive-copy
headings.

### Keep publication distinct from implementation

The reviewed branch prepares and verifies the release mechanism. Publishing the
branch, merging its pull request, and creating the version tag are workflow
evidence rather than implementation tasks. Issue #79 remains open until the tag,
assets, checksums, notes, and README synchronization are verified. This follows
the project lifecycle contract while honoring the unconditional tag guardrail.

## Project Structure

```text
.github/
├── release-notes/
│   └── v0.9.0.md
└── workflows/
    └── release.yml
scripts/
└── automation-check.sh
test/scripts/
└── automation-check_test.sh
specs/034-v090-release-cut/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── publication.md
├── checklists/
└── tasks.md
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend the existing release workflow, offline automation
contract, and documentation surfaces. No application package or dependency is
introduced.

## Complexity Tracking

No constitution violations require justification.

During full verification, the canonical lint gate exposed a release blocker:
the pinned golangci-lint v2.1.6 executable was built with Go 1.24 and refused to
analyze the Go 1.25 module. S034 therefore includes a scope-proportional pin
refresh to v2.12.0, the latest upstream release whose declared build baseline
remains Go 1.25, in the verifier and current contributor instructions. This
deviation keeps the existing gate and configuration intact while restoring the
release baseline; historical specifications retain their original evidence.
