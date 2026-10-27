# S124 source and release handoff

## Selected version and current public release

S124 prepares v1.7.0 for [#394](https://github.com/shruggietech/go-schedule/issues/394). The initial v1.6.1 proposal was conditional on classification. Documented additive maintainer incident/replay CLI selections require a minor increment under the declared SemVer policy; no breaking API, configuration or persisted-schema change was found.

The current public release remains [v1.6.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.6.0), approved source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b, release ID 423908120. Public website selection and immutable v1.6.0 notes/comparison/review records remain unchanged. New v1.7.0 tag URLs are prepared release references and do not establish that a tag or assets exist.

## Source proposal

Accepted base: a2e4653021b6007135550ed8f84f1421f9072db0, owner-merged S123 PR #393. S124 branch: codex/124-reliability-release-preparation. The source proposal includes release notes, selected-version metadata, upgrade/rollback guidance, all 80 comparison cells, ten content-bound decisions and reconciled S123 delivery. Current qualification and publication receipts are recorded separately in [verification.md](verification.md).

The new assembled-source QA uses actual development-version binaries (0.0.1-dev), not fabricated v1.7.0 packages. Source hashes, six binary identities, generated run/corpus/report bindings and semantic review describe the evaluated source. Later exact-head CI and Codex receipts describe the published source PR. Neither substitutes for actual future staging identities.

## Known issues and standing waiver

[#355](https://github.com/shruggietech/go-schedule/issues/355) and [#356](https://github.com/shruggietech/go-schedule/issues/356) remain open with needs: reproduction. Their historical runtime/prerequisite and intermittent inspection causes remain unconfirmed. Separately demonstrated defects have specific delivered fixes; successful current runs do not prove historical causes. The reconstructed historical whole-tree digest differs from the original, replay used a current controller and the old daemon lacked the private incident collector. No unchanged investigation is repeated in this release-preparation slice.

Attended native observations remain untested under the standing native-testing waiver; automated CI and artifact checks remain required. The standing authorization category is **not tested, release authorized by maintainer**. Preserve actual failed or timed-out outcomes separately. Do not use the fully attended promotion path with a waiver. This disclosure grants no new release authorization.

## Subsequent owner boundaries

1. Review and merge the source PR after exact-head checks and no more than two Codex rounds. S124 stops before this step.
2. Verify the exact accepted main and its terminal successful CI. Obtain explicit authorization for the immutable v1.7.0 tag at that commit and draft staging. Only then may existing release tooling create the tag/draft.
3. Download and verify all eight actual staged assets: four CLI/daemon archives, two desktop archives, Windows MSI and windows-candidate-manifest.json. Record the real tag/commit, Release workflow run and attempt, asset IDs/sizes/hashes, MSI ProductVersion/ProductCode and manifest binding. No such identities exist in this source-preparation handoff.
4. Present that concrete inventory for separate owner promotion approval. The existing standing-waiver promotion workflow rechecks exact identities, adds SHA256SUMS.txt and publishes the same assets as latest. Nine public assets and eight checksum entries are expected. Preserve immutable bytes and the explicit untested disclosure.
5. Reconcile actual publication, public website selection and delivery state on #394 only when those functional outcomes occur. Keep independent #355/#356 causes open.

## Upgrade and rollback

Follow the [Windows](../../docs/INSTALL-windows.md), [Linux](../../docs/INSTALL-linux.md) or [macOS](../../docs/INSTALL-macos.md) guide. Stop the daemon that owns the effective configuration/data directory, preserve its entire database and SQLite sidecars, configuration and protected credentials, and retain matching older binaries or package. v1.7.0 retains schema v24 introduced in v1.6.0. Rollback restores complete pre-upgrade state with its matching older binaries; do not point older binaries at migrated state. Remote exposure, certificate/credential identity and MSI upgrade identity are unchanged.
