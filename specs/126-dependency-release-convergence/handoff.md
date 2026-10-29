# S126 Release Continuation

## Selected source and current state

S125 is owner-merged at 34f59e7828ad87da7dac587c40f35b6ac2eb3407. Exact-main CI 40648849800, CodeQL 40648849803 and website 40648849797 all passed, with Pages deployment 7379445600. #397/#398 closed. S126 converges selected #401-#403 dependencies; #400/#404 remain deferred because a suitable corrected SQLite release is not published. Public release remains v1.6.0, ID 423908120, approved source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b.

| Identity | Disposition |
| --- | --- |
| Final S126 merged-main release source | Unresolved until owner merge and successful exact-main acceptance |
| v1.7.0 immutable tag | Not created; requires approval of the exact accepted main SHA |
| Release ID, staging run/attempt, eight asset IDs/hashes and MSI identity | Unresolved until approved staging actually creates them |
| Public promotion | Requires separate approval of the real selected staged bytes |
| Public discovery | Remains v1.6.0 until actual publication and reconciliation |

## Same-slice delivery sequence

1. Finish S126 source PR and review at its actual final head, at most two Codex rounds. Stop for the owner's merge ritual. The current authorization does not merge, tag, stage or promote.
2. After owner merge, synchronize clean main and require successful exact-commit main CI, CodeQL and website checks. Select the actual new source SHA, not the earlier accepted S125 baseline or unmerged proposal. Run existing release-source-check and comparison-check guards.
3. Obtain explicit approval of that exact immutable v1.7.0 tag and draft staging. Check remote tag and authenticated release inventory first; never move an existing tag to conceal a changed candidate. Push the approved tag only after that approval and monitor existing release.yml through all seven jobs.
4. Resolve and download the actual draft with the existing release-promotion helper. Inspect all eight expected assets, their IDs/names/sizes/SHA-256, Release run/attempt, MSI ProductVersion/ProductCode and windows-candidate-manifest.json against the same immutable source. No local QA binary substitutes for a staged asset.
5. Present the concrete inventory for separate promotion approval. The standing waiver remains **not tested, release authorized by maintainer** for attended native observations; automated and artifact gates remain required. Dispatch existing promote-release-waived.yml only after approval, publishing the same selected bytes and adding SHA256SUMS.txt. Expect nine public assets and eight checksum entries.
6. Verify public/latest identity, unchanged original assets, matching checksums and promotion-triggered Pages/HTTPS. Reconcile site/data/published.json and actual delivery records through the normal source review path as same-slice follow-through. Close #394 only when release/discovery outcomes are delivered; update #146 accordingly. Supersede #401-#403 only after their selected equivalent source actually integrates. Leave #400/#404 open while deferred.

The [existing detailed operational commands and expected asset inventory](../125-reliability-delivery-decision/handoff.md) remain authoritative for these unchanged workflows. This continuation updates source selection and issue disposition rather than introducing another release mechanism or preparation-only slice.

## Preserved product and evidence contracts

Schema remains v24. Preserve complete stopped-state database/sidecars, configuration and protected credentials with matching pre-upgrade binaries; rollback restores that complete state with matching older binaries. No cluster capability, new binder adoption or highlighter migration is claimed. #355/#356 original causes remain unconfirmed; no unchanged historical campaign or causal closure is implied by this candidate. Historical reports, binaries and first failures remain frozen.
