# Implementation Plan: Publish v1.5.0

**Branch**: `codex/104-v150-release` | **Date**: 2026-09-23 | **Spec**: [spec.md](spec.md)
**Input**: S104 and release issue #185.

## Summary

Correct the reviewed source-owned v1.5.0 metadata, merge it through the normal PR path, then tag exact green main and stage the standard platform assets as a draft. Verify source, workflow run, candidate manifest, downloaded bytes, and checksum inventory before publishing. The maintainer clarified that the native-testing waiver is standing policy, so use the reviewed standing-waiver promotion workflow without pretending that the full-attended-evidence promotion workflow passed.

## Technical Context

**Language/Version**: Go 1.26, GitHub Actions, Markdown release metadata.
**Primary Dependencies**: Existing release and promotion workflows, GitHub CLI, Windows candidate verifier.
**Storage**: Git tag, GitHub draft/public release, release assets, GitHub issue and project records.
**Testing**: Local CI parity for reviewed source changes, exact-main CI, release-staging jobs, artifact identity and checksum checks.
**Target Platform**: Linux amd64/arm64, macOS amd64/arm64 daemon and CLI, Linux amd64 and macOS arm64 desktop, Windows amd64 MSI.
**Project Type**: Release operations and source metadata.
**Performance Goals**: No runtime performance change.
**Constraints**: No direct main push, no fabricated qualification evidence, no moved tag after candidate verification, no public release with missing assets or failed automated checks.
**Scale/Scope**: One v1.5.0 release, issue #185, and its milestone. SMTP and clustered execution remain future work.

## Constitution Check

- Principle I: No application behavior changes. Reviewed source and artifact identity stay exact.
- Principle II: CI parity and hosted CI remain mandatory under the current constitution. The standing waiver applies to attended native checks, not automated CI or artifact identity.
- Principle III: Release documentation must be accurate and consistent across README, changelog, notes, and latest-release pointer.
- Principle IV: No performance-sensitive path changes.
- Principle V: Spec-kit artifacts and analysis precede implementation. Reviewed source changes use a PR; release publication is separately authorized by this S104 user request, subject to the mandatory pre-push halt.
- Governance: Close #185 only for completed publication, not merely test execution; file any discovered product defect separately.

## Research and Design

- [research.md](research.md) records the source-metadata mismatch, release workflow contract, and historical waiver precedent.
- [data-model.md](data-model.md) defines the source, candidate, public release, and disposition identities.
- [contracts/publication.md](contracts/publication.md) defines the release invariants.
- [quickstart.md](quickstart.md) orders the operator steps.

## Project Structure

```text
README.md
CHANGELOG.md
.github/release-notes/v1.5.0.md
.github/workflows/release.yml
.github/workflows/promote-release.yml
specs/104-v150-release-publication/
```

## Delivery Sequence

1. Review and merge source metadata corrections through a normal PR, with CI and external review handling.
2. Confirm exact merge-commit CI and no public v1.5.0 release or tag.
3. Create an annotated v1.5.0 tag on exact reviewed main and let the release workflow stage a draft.
4. Audit all draft assets, candidate manifest, release notes, and immutable identity. Record the standing waiver's application and every untested native observation honestly.
5. Generate and verify all-asset checksums, publish through the explicit waiver path, then verify the public release by fresh download.
6. Reconcile #185, milestone, and delivery project only after publication is confirmed.

## Complexity / Deviation

The existing attended promotion workflow remains unchanged. The maintainer clarified in S104 that native testing is waived for all releases. Do not fabricate its required evidence archive or label untested observations as passing; use the reviewed standing-waiver promotion workflow with the same source, staging, candidate, asset, and checksum checks. This is a standing operator disposition for attended native checks, not a change to automated CI or a claim of full qualification.
