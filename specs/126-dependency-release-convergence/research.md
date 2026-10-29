# S126 Research and Decisions

## Dependency convergence

Decision: integrate #401 API runtime 1.8.0 and YAML 3.0.5, #402 Vite 8.3.5 and #403 Next 16.4.3, Fumadocs core/UI 16.17.8 and MDX 15.4.10. Regenerate complete manifests/locks and perform clean restores. The installed planning command requested three read-only research agents for graph compatibility, storage selection and release/QA continuity.

Rationale: the individual proposals passed, but the assembled graph needs its own acceptance. Desktop replaces the root module locally; regenerate both Go graphs. Existing Vite/plugin/Vitest and Fumadocs/React/Next peers support these selections, with Shiki retained at 4.5.0. Do not change renderer behavior or add features just because a dependency exposes them.

Alternatives considered: merge three proposals independently (less coherent candidate acceptance); copy only root metadata (can miss secondary module requirements); broader unrequested dependency upgrades (excluded).

## Conditional storage decision

Decision: defer #400/#404 and retain modernc.org/sqlite v1.60.1 with modernc.org/libc v1.77.1.

Rationale: [official changelog](https://gitlab.com/cznic/sqlite/-/raw/master/CHANGELOG.md) describes v1.61.0 prepared-statement reparse regression and a v1.61.1 correction. At kickoff the [registry latest](https://proxy.golang.org/modernc.org/sqlite/@latest) still selects v1.61.0, origin 4df1ebea0f1cf0cd6ec3245b26f5fec4234ce738, published 2026-10-21T11:57:52Z; [v1.61.1 metadata](https://proxy.golang.org/modernc.org/sqlite/@v/v1.61.1.info) returns 404 unknown revision. The driver's published v1.61.0 module requires libc v1.79.0; retained v1.60.1 requires v1.77.1. A dated master changelog does not establish a published version.

Alternatives considered: adopt known-regressed v1.61.0 after fixing its graph (rejected); use unpublished upstream code (rejected); delay all ready work (unnecessary). This is a version-selection hold, not a demonstrated project runtime failure or #355/#356 cause. #404 remains open and proposal-specific.

## Release and evidence continuity

Decision: reuse the existing S125 release continuation and existing unattended QA coordinator. Fresh candidate QA is warranted for runtime dependencies; an unchanged historical causality campaign is not. Capture/review/finalize new reports and preserve all original artifacts and binaries.

Rationale: source PR publication is explicitly authorized; merge, exact immutable tag/draft and concrete-asset promotion are separate owner actions. No actual v1.7.0 asset identities are available before staging. Preserve public v1.6.0 selection, original causal uncertainty and standing native-testing waiver attribution.

Alternatives considered: another release framework or new preparation-only cycle (unnecessary); treating green source as actual publication (incorrect). Refresh pending content-bound comparison hashes without changing 80 capability cells or the ten decisions absent changed capabilities.
