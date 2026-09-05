# Implementation Plan: v1.0.0 Release Cut

**Branch**: `codex/048-v100-release-cut` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/048-v100-release-cut/spec.md`

## Summary

Prepare the reviewed v1.0.0 boundary by cutting the complete post-v0.9.1 changelog, adding concise tag-specific highlights, synchronizing the two README version lines, and documenting the exact post-merge staging, 47-observation Windows qualification, issue reconciliation, and promotion ritual. The implementation ends with a green preparation PR. Tag and release mutations remain separately authorized after merge.

## Technical Context

**Language/Version**: Markdown, GitHub Actions YAML contract inspection, PowerShell 7.6-compatible operator commands, Go 1.25 verification suite

**Primary Dependencies**: Existing Release and Promote Release workflows, GitHub Releases, S047 release-gate CLI and attended collector

**Storage**: Versioned repository documents plus future hosted draft assets and attended evidence archive

**Testing**: Existing release-note policy fixtures, release workflow integration contracts, S047 synthetic evidence suite, all eight canonical gates, hosted PR CI

**Target Platform**: Repository preparation on all supported hosts; post-merge candidate staging on GitHub Actions; attended qualification on Windows 11

**Project Type**: Release preparation and operational contract

**Performance Goals**: No runtime path changes; release notes remain four to six highlights; validation remains within existing evidence/archive bounds

**Constraints**: UTF-8 without BOM; no tag or release before reviewed merge and separate authority; no product changes unless a release blocker is reproduced; no copied local-demo results; no Post-v1 scope

**Scale/Scope**: 33 current Unreleased entries, 12 merge commits since v0.9.1, one version boundary, one release-note file, two README strings, nine staged pre-checksum assets, 47 observations, ten remaining v1 issues, one release issue

## Constitution Check

*GATE: Passed before research and rechecked after design.*

| Principle | Requirement | Design disposition |
| --- | --- | --- |
| I. Code Quality | Release artifacts and contracts are explicit and maintainable | PASS. S048 reuses the existing workflow architecture and makes one bounded README metadata-preflight correction rather than introducing another path. |
| II. Testing Standards | Behavioral changes require tests and canonical gates | PASS. The preflight correction has a negative regression fixture; existing release-copy, workflow, evidence, race, coverage, and canonical checks remain mandatory. |
| III. UX Consistency | Public release copy and version identity agree | PASS. The badge reflects published state while the changelog, health example, notes, tag, artifacts, and issue state consistently describe the v1.0.0 boundary. |
| IV. Performance | No unexplained runtime regression | PASS. No runtime or hot-path source changes; no benchmark change is warranted. |
| V. Autonomous Execution | Use Spec Kit, PR-first integration, and explicit publication authority | PASS. S048 runs end to end through the green PR and leaves tag/release mutations behind the required post-merge authorization. |

No constitutional deviation is required. The scope correction from “qualify before PR” to “prepare PR, then qualify after merge” is required by Principle V and by the candidate's immutable reviewed-source identity.

## Technical Design

### Release boundary

Move the current contents between `## [Unreleased]` and `## [0.9.1]` beneath `## [1.0.0] - 2026-09-03`, leaving a new empty Unreleased heading. Preserve the 33 top-level entries byte-for-byte apart from their enclosing boundary. Add comparison links for `v1.0.0...HEAD` and `v0.9.1...v1.0.0`.

### Public identity

Replace the hard-coded README release badge with Shields.io's GitHub release badge, whose source is the latest published release behind the existing `releases/latest` link. Change the single health example from 0.9.1 to 1.0.0. The tagged workflow validates the stable badge contract and requires the health example to match `GITHUB_REF_NAME`, so preparation can be reviewed without advertising an unpublished version.

### Curated release notes

Add `.github/release-notes/v1.0.0.md` with five outcome-oriented highlights: cross-platform scheduling breadth, dependable task execution, Windows setup and uninstall control, desktop usability/accessibility, and exact-candidate release assurance. Link once to the tagged v1.0.0 changelog. Keep generated notes off.

### Post-merge contract

Document a chronological, fail-closed procedure:

1. verify the reviewed merge is clean synchronized `main` and v1.0.0 is absent;
2. obtain explicit tag authorization, create and push the annotated tag;
3. wait for the Release workflow and audit the draft's seven build packages plus candidate manifest;
4. initialize an attended workspace from the exact staged MSI/manifest;
5. complete and finalize all 47 observations without copying local-demo passes;
6. upload the evidence archive and reconcile ten issues individually;
7. dispatch Promote Release and audit nine checksummed payloads plus `SHA256SUMS.txt`, public/latest identity, issue #122, and milestone state.

The contract names stop conditions at every irreversible boundary.

## Project Structure

```text
.github/release-notes/
└── v1.0.0.md

specs/048-v100-release-cut/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── publication.md
└── checklists/
    ├── requirements.md
    └── release-contract.md

README.md
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Reuse the established S034 release-cut surfaces and the S040/S047 gate. Amend the existing Release metadata preflight and automation fixture for published-state badge correctness; no new workflow, validator, dependency, or product component is justified.

## Implementation Sequence

1. Record baseline identity, history cardinality/hash, tag/release absence, milestone state, and expected assets.
2. Add the v1.0.0 release-note file and confirm the existing policy accepts it.
3. Cut the changelog, update the health example, and make the release badge publication-aware.
4. Write the post-merge publication contract, data model, and quickstart.
5. Run Spec Kit analysis, focused release contracts, and all eight gates.
6. Commit, push the authorized review branch, open the PR with `Refs #122`, and resolve up to two Codex review rounds plus hosted CI.
7. Stop for the maintainer's merge ritual. Do not tag or publish.

## Post-Design Constitution Recheck

All five principles remain PASS. The final design changes one release metadata preflight but no pinned action or runtime code, and it preserves every existing staging, evidence, asset, checksum, and promotion boundary. Release publication remains impossible from the preparation merge alone.

## Complexity Tracking

No constitution violation requires justification.
