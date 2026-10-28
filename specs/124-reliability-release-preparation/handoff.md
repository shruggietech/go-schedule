# S124 source and release handoff

## Selected version and current public release

S124 prepares v1.7.0 for [#394](https://github.com/shruggietech/go-schedule/issues/394). The initial v1.6.1 proposal was conditional on classification. Documented additive maintainer incident/replay CLI selections require a minor increment under the declared SemVer policy; no breaking API, configuration or persisted-schema change was found.

The current public release remains [v1.6.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.6.0), approved source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b, release ID 423908120. Public website selection and immutable v1.6.0 notes/comparison/review records remain unchanged. New v1.7.0 tag URLs are prepared release references and do not establish that a tag or assets exist.

## Source proposal and owner merge

Accepted base: a2e4653021b6007135550ed8f84f1421f9072db0, owner-merged S123 PR #393. S124 branch: codex/124-reliability-release-preparation. The source proposal includes release notes, selected-version metadata, upgrade/rollback guidance, all 80 comparison cells, ten content-bound decisions and reconciled S123 delivery. Current qualification and publication receipts are recorded separately in [verification.md](verification.md). The owner merged [PR #396](https://github.com/shruggietech/go-schedule/pull/396) on 2026-10-28T02:17:16Z at b4f1e281e66003e08ee2a7abd35e3e25468fb6bd, tree-identical to reviewed head 5c7dce3cb8427c8f23367b50c5ee40d8d0139346. [delivery.md](delivery.md) records the actual merge and distinct main acceptance.

The new assembled-source QA uses actual development-version binaries (0.0.1-dev), not fabricated v1.7.0 packages. Source hashes, six binary identities, generated run/corpus/report bindings and semantic review describe the evaluated source. Later exact-head CI and Codex receipts describe the published source PR. Neither substitutes for actual future staging identities.

## Known issues and standing waiver

[#355](https://github.com/shruggietech/go-schedule/issues/355) and [#356](https://github.com/shruggietech/go-schedule/issues/356) remain open with needs: reproduction. Their historical runtime/prerequisite and intermittent inspection causes remain unconfirmed. Separately demonstrated defects have specific delivered fixes; successful current runs do not prove historical causes. The reconstructed historical whole-tree digest differs from the original, replay used a current controller and the old daemon lacked the private incident collector. No unchanged investigation is repeated in this release-preparation slice.

Attended native observations remain untested under the standing native-testing waiver; automated CI and artifact checks remain required. The standing authorization category is **not tested, release authorized by maintainer**. Preserve actual failed or timed-out outcomes separately. Do not use the fully attended promotion path with a waiver. This disclosure grants no new release authorization.

## Subsequent owner boundaries

The source review/merge boundary for S124 is complete. Its exact-main [CI 40620037771](https://github.com/shruggietech/go-schedule/actions/runs/40620037771) then failed the production browser Search focus assertion, tracked independently as [#397](https://github.com/shruggietech/go-schedule/issues/397). The other 22 check runs succeeded, including CodeQL and Pages. Do not tag this red b4f1e28 source, or treat the green reviewed-head checks as a successful main receipt.

[S125 release continuation](../125-reliability-delivery-decision/handoff.md) now carries the bounded correction, pending v1.7.0 content binding and concrete existing-workflow commands. Its future owner-merged source must receive successful exact-main acceptance before the owner approves an immutable v1.7.0 tag/draft. Actual assets then require their own inventory and separate promotion approval. This handoff selects no future source SHA, draft ID or asset identity.

No tag, draft release or promotion was performed in S124 or its post-merge housekeeping. Public v1.6.0 remains unchanged; #394 stays open for actual delivery and #355/#356 causes remain unresolved. Do not repeat another unchanged historical investigation as a release gate.

## Upgrade and rollback

Follow the [Windows](../../docs/INSTALL-windows.md), [Linux](../../docs/INSTALL-linux.md) or [macOS](../../docs/INSTALL-macos.md) guide. Stop the daemon that owns the effective configuration/data directory, preserve its entire database and SQLite sidecars, configuration and protected credentials, and retain matching older binaries or package. v1.7.0 retains schema v24 introduced in v1.6.0. Rollback restores complete pre-upgrade state with its matching older binaries; do not point older binaries at migrated state. Remote exposure, certificate/credential identity and MSI upgrade identity are unchanged.
