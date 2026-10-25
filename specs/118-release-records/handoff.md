# S118 current v1.6.0 release handoff

Updated: 2026-10-25. S118 reconciles completed source delivery and the remaining #354-#356 outcomes. It authorizes one reviewed source PR, not merge, tagging, staging or promotion. Public package selection remains v1.5.1.

## Delivered source and remaining issues

Accepted baseline is PR #377 squash 9ad13821c825d01efcc038b6b1634b187a854758, tree-identical to reviewed head 7327a7c401b6140d0b6c02d2769ea3082d74a123. Its exact merged-main CI, CodeQL and Website runs passed all 23 checks, with Pages deployment 7316855350. See [S117 delivery](../117-candidate-qa-handoff/delivery.md) for linked run/review receipts. S118 adds source records after that baseline, so a future release must select the actual later reviewed main commit rather than automatically tagging 9ad1382.

| Issue | Delivered outcome | Remaining functional outcome |
| --- | --- | --- |
| [#354](https://github.com/shruggietech/go-schedule/issues/354) | SMTP, CLI, unattended evaluation and source corrections are merged; #375/#376 are closed by PR #377; pending v1.6.0 guidance is coherent | Actual public release and public discovery reconciliation; keep open |
| [#355](https://github.com/shruggietech/go-schedule/issues/355) | Original failure evidence, bounded prerequisite/launcher diagnostics and preserved six-writer assertions; focused native contracts passed in S117 | Original PowerShell CLR access violation and separate historical prerequisite refusal causes remain unconfirmed; keep open |
| [#356](https://github.com/shruggietech/go-schedule/issues/356) | Partial observations, phase/request attribution, cancellation/recovery evidence and owner-directed supported completion allowances | Original intermittent owned Windows IPC/inspection stall cause remains unconfirmed; keep open |

The original S117 report remains failed (1,757 effective passes, six communication failures, three exclusions; all 1,763 execution cases passed). It is bound to dirty planning-only 11485a8 and development version 0.0.1-dev, not the corrected PR or a release package. The S115 failed report also remains unchanged. Exact report/input hashes and later corrective acceptance remain in [S117 verification](../117-candidate-qa-handoff/verification.md). Passing later runs do not establish either original cause or close #355/#356.

## Reviewed pending content

The selected changelog heading remains `[1.6.0] - 2026-10-22`, anchor `160---2026-10-22`. Its date records source preparation. Notes already include #375/#376 and unresolved #355/#356. Current content bindings are:

| Record | SHA-256 |
| --- | --- |
| Comparison dataset | ab2bb5bf2de20133735fa4cda541e6d48e557e4632edbd8d95dd8ca14af0cb5b |
| Selected v1.6.0 changelog section | 26fa24a30c70cc5f401213457382c9f57e10bff8a7e0e85e9b925ad7712aa269 |
| v1.6.0 release notes | 4062f9e664f32298f9f4de8dbeacf4c0b317e0e003c83b185290e18a1c8874b1 |

S118 does not change those selected bytes or capability decisions. Source-owned `site/data/published.json` continues to select v1.5.1. Existing platform instructions retain stopped-daemon complete-state backup before forward-only schema-v24 migration and restoring matching older binaries plus the complete earlier state for rollback. See [Windows installation](../../docs/INSTALL-windows.md), [Linux installation](../../docs/INSTALL-linux.md) and [macOS installation](../../docs/INSTALL-macos.md).

## Later source selection and immutable tag

After the owner merges the source PR, fetch/fast-forward main and verify terminal push-event CI for that exact merged SHA, plus CodeQL and website deployment. PR-head success does not substitute for exact-main acceptance. Inspect whether v1.6.0 or a matching release already exists; never move an existing tag or replace published bytes.

Present the actual final main SHA for explicit tag/staging authorization. Recheck the selected notes, changelog anchor, comparison bindings and backup/rollback guidance on that source. Any later source change requires selecting and auditing the new source. Read-only release-source/comparison checks validate source content and supplied identity syntax; they do not establish an actual tag, staging run, release ID or asset.

With separate authorization, create/push only the approved immutable tag on that exact SHA. The existing tag-only [Release workflow](../../.github/workflows/release.yml) requires successful exact-commit main CI and read-only source/comparison checks before draft staging. Do not substitute a dispatch build or local demonstration package.

## Actual draft staging identities

Retain the resolved tag commit, terminal successful Release run and attempt, authenticated release ID, asset IDs and original bytes. The expected staged inventory is four daemon/CLI archives (Linux amd64/arm64, Darwin amd64/arm64), two desktop archives (Linux amd64, Darwin arm64), Windows amd64 MSI and windows-candidate-manifest.json.

The Windows manifest binds repository, tag commit, workflow run/attempt, MSI filename, byte size, SHA-256 and product identity. Retain the actual successful Windows attempt and complete expected allowlist. Do not rebuild a nominally equivalent MSI, reuse an earlier manifest or invent future asset IDs/digests. Source readiness alone supplies none of these identities.

## Separately approved standing-waiver promotion

Present the actual source/tag/run/attempt/release-ID/asset identities for separate owner publication approval. Use [Promote Release (standing native-testing waiver)](../../.github/workflows/promote-release-waived.yml) from reviewed main with the exact staged tag. The helper is built from immutable workflow source before tagged checkout. Authenticated paginated lookup selects one release ID retained through all asset/checksum/publication operations.

Retain the exact successful staging attempt, candidate manifest, complete asset bytes and SHA256SUMS.txt. Identical existing checksums may be reused after an interrupted upload; differing, changed, incomplete or unexpected assets stop promotion. Preserve the standing disclosure: attended native observations are not tested, release authorized by maintainer. Existing source/automated/artifact contracts remain required; no installed-system, workstation restart or session operation is requested. See [Windows waiver and recovery procedure](../../test/windows/README.md).

For interruption recovery, inspect the actual failure and rerun only the same reviewed workflow and tag. A complete already-published release can recover its latest pointer without changing content. Authentication or permission failure requires the actual access fix; do not substitute a different release, delete candidate bytes or bypass identity checks.

## Public discovery outcome

After authorized promotion, verify the actual public release ID, tag/source, latest pointer, complete downloads and matching checksums. Through a reviewed source PR, reconcile `site/data/published.json` and pending-publication prose to the actual release while preserving historical snapshots. Verify terminal website deployment and public HTTPS discovery, installation and comparison pages. Close #354 when its public release/discovery outcome is delivered. Keep independently unresolved #355/#356 open; qualification archives and passing repeats are not closure gates.

No future staging or publication identity is claimed by S118. Actual S118 exact-head CI and review receipts belong in its PR handoff; no more than two Codex rounds are authorized.
