# S125 release continuation

## Current source and publication state

**S126 continuation:** S125 is owner-merged in PR #399 at 34f59e7828ad87da7dac587c40f35b6ac2eb3407, with exact-main CI 40648849800, CodeQL 40648849803 and website 40648849797 successful. #397/#398 closed. The [S126 continuation](../126-dependency-release-convergence/handoff.md) selects the next combined candidate after compatible dependency convergence; the historical S124 failure below remains preserved.

This is the continuation of [#394](https://github.com/shruggietech/go-schedule/issues/394), using the existing release workflows. S124 source preparation is owner-merged in [PR #396](https://github.com/shruggietech/go-schedule/pull/396), commit b4f1e281e66003e08ee2a7abd35e3e25468fb6bd. Its exact-main [CI 40620037771](https://github.com/shruggietech/go-schedule/actions/runs/40620037771) failed the browser focus contract tracked by [#397](https://github.com/shruggietech/go-schedule/issues/397). Successful S124 PR checks, CodeQL and website deployment do not make that red main source eligible for tagging.

S125 corrects the browser fixture ordering and retains bounded failure diagnostics, includes that change in pending v1.7.0 notes/changelog, and refreshes their content-bound comparison-review hashes. The 80 capability cells and their 2026-10-27 vendor evidence remain unchanged; this content rebind advertises no new scheduler or clustered-execution capability. Current S125 qualification and exact PR identities belong in [verification.md](verification.md) and [delivery.md](delivery.md), separate from the original S124 application QA.

Public release remains [v1.6.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.6.0), ID 423908120, approved source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. site/data/published.json stays at v1.6.0. v1.7.0 is the prepared version, not an assertion that tagged assets exist.

| Identity | Current disposition |
| --- | --- |
| Accepted S125 baseline | 34f59e7828ad87da7dac587c40f35b6ac2eb3407, exact-main CI/CodeQL/website successful; S126 final release source remains unresolved until its owner merge and acceptance |
| Immutable v1.7.0 tag | Not created by this source work; requires exact-source approval |
| v1.7.0 release ID and staging run/attempt | Unresolved until authorized draft staging actually creates them |
| Asset IDs, sizes, hashes, MSI ProductVersion/ProductCode and candidate manifest | Unresolved until the real staged bytes are downloaded and inspected |
| Public promotion and latest-release selection | Not authorized by S125 source publication |

## 1. Finish source review and owner merge

Complete all applicable checks on the final PR head and address every review comment within the authorized maximum of two Codex rounds. Present that exact source to the owner for final review. S125 stops before merge, tagging, draft staging or promotion. The later owner merge is an explicit boundary, not an action inferred from green checks.

After the owner merges, synchronize the clean checkout and obtain the actual new main identity:

```sh
git fetch --prune origin
git switch main
git merge --ff-only origin/main
RELEASE_COMMIT=$(git rev-parse HEAD)
gh run list --repo shruggietech/go-schedule --commit "$RELEASE_COMMIT" --limit 20 --json databaseId,name,event,headSha,status,conclusion,url
```

Require terminal successful exact-source main CI and inspect its individual checks; also verify CodeQL and website deployment for that same SHA. Do not substitute reviewed branch CI, a different main commit, a green rerun of unrelated source or website success for source acceptance. A genuine new failure receives its own bug record and correction. Preserve the first failure.

Run the existing source guards from that clean accepted checkout:

```sh
go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.7.0
go run ./scripts/comparison-check -tag v1.7.0 -commit "$RELEASE_COMMIT"
```

The source guard validates metadata, not tag/package existence. The comparison guard validates reviewed content against the supplied identity; it does not prove a remote tag exists. On Windows, use the established hidden noninteractive launcher for console tools.

## 2. Obtain exact tag and draft-staging approval

Present the owner with the resolved 40-character main SHA, terminal source acceptance, selected v1.7.0 metadata and unresolved known-issue disclosure. Ask for approval of that exact immutable tag source and the existing draft-building workflow. Do not preselect b4f1e28 or the unmerged S125 branch.

Check for an existing v1.7.0 tag and authenticated draft/public release before creating anything. An existing tag must resolve to the approved source; never move it to hide a changed candidate. An unexpected existing release or changed source requires inspection rather than automatic replacement. Once, and only once, exact-source authorization is given, create/push the approved v1.7.0 tag. The tag push triggers [.github/workflows/release.yml](../../.github/workflows/release.yml).

The Release workflow requires successful exact-commit push/main CI, current source metadata, reviewed comparison and absent/draft release state before staging. Its seven jobs build the same source into four daemon/CLI archives, two desktop archives, one Windows MSI and one candidate manifest. Watch the actual workflow run and every staging job to terminal success; preserve run and attempt identities.

## 3. Inspect the concrete draft

Use authenticated release listing or the existing release-promotion helper to select the draft by exact tag and retain its release ID. Do not mistake a published-tag lookup for authenticated draft discovery. The existing helper's resolve/download operations bind one selected release and asset IDs:

```sh
go run ./scripts/release-promotion resolve --repository shruggietech/go-schedule --tag v1.7.0 --state .test-temp/v1.7.0-release/release-selection.json
go run ./scripts/release-promotion download --state .test-temp/v1.7.0-release/release-selection.json --dir .test-temp/v1.7.0-release/staged-assets
```

Run these commands from the accepted repository checkout after creating a new private .test-temp/v1.7.0-release directory; stop to inspect an existing directory rather than overwriting prior receipts. Supply authenticated GH_TOKEN through the normal credential route; do not print credentials or save them in source. These read/download commands neither upload checksums nor publish the release. Preserve the selection receipt with the actual asset inventory.

Require exactly these eight initial staged assets:

- go-schedule_v1.7.0_darwin_amd64.tar.gz
- go-schedule_v1.7.0_darwin_arm64.tar.gz
- go-schedule_v1.7.0_linux_amd64.tar.gz
- go-schedule_v1.7.0_linux_arm64.tar.gz
- go-schedule-desktop_v1.7.0_darwin_arm64.tar.gz
- go-schedule-desktop_v1.7.0_linux_amd64.tar.gz
- go-schedule_v1.7.0_windows_amd64.msi
- windows-candidate-manifest.json

Record the actual repository/tag/commit, release ID, Release workflow run and attempt, every asset ID/name/size/SHA-256, and MSI ProductVersion/ProductCode. Use the existing candidate verifier against the downloaded MSI and manifest:

```sh
go run ./scripts/windows-release-gate verify-candidate --candidate-manifest .test-temp/v1.7.0-release/staged-assets/windows-candidate-manifest.json --artifact .test-temp/v1.7.0-release/staged-assets/go-schedule_v1.7.0_windows_amd64.msi --repository shruggietech/go-schedule --tag v1.7.0 --commit "$RELEASE_COMMIT"
```

The manifest must identify the same immutable source and successful Release staging attempt, including the Windows job. Recheck release disclosure, source guards and actual inventory before presenting the assets. No development-version S124/S125 binary or local-demo MSI substitutes for a workflow-staged candidate.

## 4. Obtain concrete-asset promotion approval

Present the downloaded inventory and exact identities to the owner. Request separate approval to publish those actual staged bytes. The standing native-testing waiver remains **not tested, release authorized by maintainer**; attended native observations remain untested. Preserve failed or timed-out automated outcomes separately. The waiver does not excuse unsuccessful CI, changed assets or mismatched identities.

After that approval, dispatch the existing standing-waiver workflow from reviewed main:

```sh
gh workflow run promote-release-waived.yml --repo shruggietech/go-schedule --ref main -f tag=v1.7.0
```

Do not use promote-release.yml under this waiver: that separate path requires genuine fully attended evidence. The standing-waiver workflow builds its helper from its own immutable workflow source, resolves the selected release ID, verifies all eight assets and the exact successful staging attempt, checks the candidate manifest/comparison, generates SHA256SUMS.txt, re-downloads the final set and rechecks the immutable tag before publishing the same release as latest.

Expect nine public assets with eight checksum entries. A same-tag recovery reuses identical checksums and selected bytes; the helper refuses changed identities, unexpected/incomplete assets and conflicting checksums. Do not rebuild or replace assets while claiming the original candidate was promoted.

## 5. Reconcile actual delivery and public discovery

After successful promotion, verify the actual release is public/latest with the approved tag/source, unchanged original asset IDs/hashes and matching final checksums. Watch the promotion-triggered Pages deployment and verify public HTTPS, documentation, downloads and the comparison's selected version. Pages resolves the latest published GitHub release; the candidate snapshot's prior merge alone does not switch public selection.

Reconcile site/data/published.json, current installation/release wording and source delivery records only after publication actually occurs. Any necessary source reconciliation follows the ordinary PR workflow as operational follow-through in this same delivery plan, without another preparation-only design slice. Close #394 when its release and discovery outcomes are delivered, then reconcile parent #146. Keep #355/#356 open while their original causes remain unconfirmed; their closure is not a release prerequisite.

## Recovery and remaining decisions

Keep the complete stopped-state database/SQLite sidecars, configuration and protected credentials together with the matching pre-upgrade binaries. v1.7.0 retains schema v24 introduced in v1.6.0. Rollback restores the complete prior state with its matching older binaries, never older binaries against migrated state. Platform guides remain authoritative.

The distributed-scheduling assessment is a separate lane of S125. Its hold/proceed/stop recommendation does not add clustered execution to this release, activate #187/#188 without their gates, or authorize release mutation. See [decision.md](decision.md) for issue-level disposition.
