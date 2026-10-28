# S124 delivery

## Delivered source

The owner squash-merged [PR #396](https://github.com/shruggietech/go-schedule/pull/396) on 2026-10-28T02:17:16Z at b4f1e281e66003e08ee2a7abd35e3e25468fb6bd. Its tree equals reviewed source 5c7dce3cb8427c8f23367b50c5ee40d8d0139346. The remote branch was deleted and the obsolete local branch removed after that equality was established.

S124 delivers v1.7.0 source metadata, tag-specific notes, upgrade/rollback guidance, the reviewed eight-product comparison and the prepared-version validator correction [#395](https://github.com/shruggietech/go-schedule/issues/395). #395 closed on merge. S124's accepted base was a2e4653021b6007135550ed8f84f1421f9072db0; that historical identity is not replaced by the later squash commit.

## Separate acceptance stages

The original assembled-source Windows/WSL evaluation passed 1763 applicable cases, with zero failed/blocked/findings and three not applicable; all 1109 messages received semantic review. Its actual binaries report 0.0.1-dev. The complete corrected canonical aggregate passed all eight gates in 885.179 seconds. [verification.md](verification.md) retains the original source/binary bindings and both earlier failed aggregates.

All 22 PR checks and two Codex rounds passed on reviewed source 5c7dce3cb8427c8f23367b50c5ee40d8d0139346. The automatic round completed 2026-10-27T23:48:12.750187Z; the second and final round completed 2026-10-27T23:57:45.482750Z. The second request is [6456619066](https://github.com/shruggietech/go-schedule/pull/396#issuecomment-6456619066) and clean receipt is [6456648567](https://github.com/shruggietech/go-schedule/pull/396#issuecomment-6456648567). These reviewed-source results do not turn a later main failure into success.

On exact merged main b4f1e28, [CI 40620037771](https://github.com/shruggietech/go-schedule/actions/runs/40620037771) failed the production browser Search focus assertion. The other 22 check runs succeeded, including [CodeQL 40620039916](https://github.com/shruggietech/go-schedule/actions/runs/40620039916) and [website deployment 40620037720](https://github.com/shruggietech/go-schedule/actions/runs/40620037720). This independent failure is [#397](https://github.com/shruggietech/go-schedule/issues/397), with original receipt retained; b4f1e28 is not an accepted release source.

Pages deployment 7375406654 succeeded. Trusted public HTTPS checks covered apex/www, documentation, comparison, Windows installation guidance, current/legacy CLI QA links and search index. The merged Windows guide is live; public comparison remains v1.6.0.

## Housekeeping and remaining outcomes

Post-merge housekeeping removed 51 inventoried S124 scratch helper/draft files while preserving QA reports, logs, reviews and publication receipts. All 137 frozen evidence/binary inputs remained byte-identical. The final post-merge receipts are retained under ignored .test-temp/postmerge-s124.

[#394](https://github.com/shruggietech/go-schedule/issues/394) stays open for actual release delivery. Public v1.6.0 remains release 423908120 at approved immutable source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. No v1.7.0 tag, draft release or promotion was created by S124. [S125's continuation](../125-reliability-delivery-decision/handoff.md) carries the #397 correction and the existing release sequence through separately authorized owner boundaries.

Original [#355](https://github.com/shruggietech/go-schedule/issues/355) and [#356](https://github.com/shruggietech/go-schedule/issues/356) causes remain unconfirmed. S124's passing evaluation neither diagnoses nor closes them.
