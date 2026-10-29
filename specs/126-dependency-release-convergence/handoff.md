# S126 Published Release and Reconciliation

## Current publication

The owner merged source PR #405 at df54fc7b2d49d3b71125a049c42bfd4225503b1d, tree-identical to proposal 78ae1cab2cf5e420ae4ddb52a644b95be0a8e6d0. After successful exact-main acceptance, the owner separately approved the immutable v1.7.0 tag/draft and then the concrete staged assets. Release [v1.7.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.7.0), ID 428427523, became public/latest at 2026-10-29T21:54:39Z.

| Identity | Accepted value |
| --- | --- |
| Immutable source | df54fc7b2d49d3b71125a049c42bfd4225503b1d |
| Annotated tag object | df8de602bdb1f99fcd567910440871eb5c78b43d |
| Exact-main CI | [40951385827](https://github.com/shruggietech/go-schedule/actions/runs/40951385827), success, attempt 1 |
| Exact-main CodeQL | [40951385833](https://github.com/shruggietech/go-schedule/actions/runs/40951385833), success, attempt 1 |
| Main website | [40951385847](https://github.com/shruggietech/go-schedule/actions/runs/40951385847), success, Pages 7423721911 |
| Draft staging | [Release 40960084352](https://github.com/shruggietech/go-schedule/actions/runs/40960084352), seven successful jobs, attempt 1 |
| Windows staging job | 123751628992 |
| Promotion | [40961590136](https://github.com/shruggietech/go-schedule/actions/runs/40961590136), success, attempt 1 |
| Publication-triggered website | [40961693487](https://github.com/shruggietech/go-schedule/actions/runs/40961693487), success, Pages 7425164305 |
| MSI ProductVersion | 1.7.0 |
| MSI ProductCode | {46A1517E-8ECD-4F7F-8D5E-F77D2AC2BAD5} |
| MSI SHA-256 | 671476fa3023d55de68c15b71043c2784ecb672d8886a4685b63e4bd3f8dd8b9 |
| SHA256SUMS.txt | Asset 690023419, 844 bytes, SHA-256 513f72f21e44e79ece2539fa9b72f64a9e5efc1abe6f2d4614f63f4aaea7638a |

The [concrete approved inventory](https://github.com/shruggietech/go-schedule/issues/394#issuecomment-6509164486) records all eight original asset IDs, sizes and SHA-256 values. Promotion preserved every original ID and byte and added only SHA256SUMS.txt. Nine public files were downloaded independently; all match GitHub digests, and the checksum file contains exactly eight matching entries. The Windows manifest binds this source, Release run/attempt and MSI hash. Read-only MSI database inspection independently confirms ProductVersion and ProductCode.

The six archives contain fifteen Go executables identifying v1.7.0. Thirteen embed the approved revision; the two Wails GUI executables omit a VCS stamp and are bound by the successful exact-source staging workflow. No independent embedded revision or attended native acceptance is claimed for those two binaries. Attended native observations remain **not tested, release authorized by maintainer** under the standing waiver; automated staging and artifact gates passed.

## Source and review history

S125 baseline 34f59e7828ad87da7dac587c40f35b6ac2eb3407 was accepted before S126. Source proposal #405 selected equivalent #401-#403 dependency blobs; Dependabot closed those proposals and deleted their branches after integration. Retained SQLite 1.60.1/libc 1.77.1 and deferred #400/#404 remain distinct from shipped dependency updates.

Both authorized Codex request slots were consumed on proposal head 78ae1ca. The first automatic request retained a pending receipt; the final manual request is comment 6500037248. At owner merge, the summary still reported running and no completed review or findings were published. The owner subsequently merged; no clean completion is inferred. The owner deleted the scheduled review monitor. No third review request was made.

[verification.md](verification.md) preserves the earlier assembled-source QA and canonical acceptance separately from exact-main CI, actual staging and public downloads. Development-version binaries are not substituted for released packages. Historical reports and failures remain unchanged.

## Same-slice publication reconciliation

Public HTTPS at go-schedule.io and the www redirect work. Documentation, the Windows guide and static search return successfully, and the live comparison selects v1.7.0 after publication-triggered Pages deployment. This follow-through aligns site/data/published.json and current source prose with those delivered facts. It changes no immutable tag, published asset, release note, comparison snapshot or comparison-review binding.

The source reconciliation follows the normal review-branch/PR boundary. #394 can close when this remaining source discovery outcome is integrated; #146 remains open for the independently conditional distributed-scheduling decision. #355/#356 original causes remain unconfirmed and open; #400/#404 remain deferred. No clustered capability is added or claimed.

## Recovery and operational references

v1.7.0 retains schema v24 introduced in v1.6.0. Preserve the complete stopped-state database/sidecars, configuration and protected credentials together with matching pre-upgrade binaries. Rollback restores that complete prior state with its matching older binaries. Platform installation guides remain authoritative.

The [S125 operational procedure](../125-reliability-delivery-decision/handoff.md) retains the existing staging, selection, download, candidate and promotion commands. This candidate has already been published; do not recreate its tag, rebuild its assets or dispatch promotion as a new release. Further releases require their own exact-source and concrete-asset approvals.
