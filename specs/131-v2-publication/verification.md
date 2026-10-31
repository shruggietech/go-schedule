# S131 Verification

## Exact release-source selection (2026-10-31)

Accepted source 75eae1ae11ad40db38341dc7d63f32876cba6483 is owner-merged S130. Fresh readback confirms all 24 exact-main checks and CI 41203356073, CodeQL 41203355921 and website 41203356094 succeeded on attempt 1. Existing source/comparison guards pass for v2.0.0 at this SHA. No local/remote v2.0.0 tag or authenticated draft exists at selection; public latest is v1.7.0. No new issues arrived since the preceding assessment.

Spec-Kit specify, clarify, checklist, plan, tasks and blocking analysis completed. Both requirements-quality checklists pass 12/12. Actual staged-asset, promotion, public-discovery and final source-canonical verification remain pending. These source checks do not establish tag/package acceptance. Attended native observations remain untested under the standing waiver. Private exact receipts are under .test-temp/s131.

## Approved staging

The owner explicitly approved v2.0.0 tag/draft staging at accepted source. The immutable tag is pushed; Release run 41218081205 attempt 1 is in progress. Exact-source CI, metadata/comparison and absent/draft-state staging gates succeeded. Complete candidate inventory, actual staged-byte inspection and public promotion remain pending.

## Actual staged and published bytes

Release run 41218081205 attempt 1 completed all seven jobs successfully. Draft 430039348 contained all eight expected assets; the existing authenticated identity helper downloaded those exact asset IDs. Independent SHA256 checks matched every GitHub digest. Six archive inventories and 15 embedded v2.0.0 binary identities passed; all available VCS revisions matched the source. Two Wails GUI binaries omit VCS revision; eight other binaries report modified checkout metadata from existing staging output, retained explicitly rather than described as clean builds. macOS bundle version is 2.0.0. Read-only compiled MSI inspection, independent ProductVersion/ProductCode extraction and windows-release-gate candidate verification passed, bound to source/run/attempt and downloaded MSI hash.

After separate concrete-candidate approval, promotion 41218776827 attempt 1 succeeded. setup-go reported a nonfatal existing-file cache restore warning; all actual identity/checksum/publication steps passed. The same release became public/latest at 2026-10-31T16:16:04Z. Independent final download verified all nine files/eight checksum entries and unchanged original asset IDs/bytes. SHA256SUMS.txt asset 695846859 has SHA256 b9eea87a4c1f33fbab0290248839864c9d57e6205fe3dbedc2bbbd61b88705f8. No artifact was rebuilt or replaced. Prior v1.7.0 release metadata and asset identities/digests remained unchanged.

Website 41218860788 attempt 1 succeeded; deployment 7461113991 selects public v2.0.0 comparison. HTTPS home, comparison, Windows guide and unauthenticated manifest/checksum URLs returned 200. The live guide pre-publication status sentence awaits this reconciliation PR merge, while local current guide/fallback wording is corrected. Source/CI, actual package bytes, public discovery and attended native behavior remain separate evidence classes.

## Final reconciliation-source verification

All eight foreground gates passed in 4534.339 seconds with zero drift across 2557 source paths and all 221 retained inputs; ten headless site browser cases and seven rendered-publication checks passed. Core coverage: engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.9%, catchup 88.9%, logbus 91.1%. All 167 frontend tests, the native desktop build and 24-document export passed. The canonical aggregate includes the existing release/security/identity negative fixtures and CI timeout/lifecycle checks; no gate or assertion was weakened.

Evaluated source manifest SHA256: 84243d59835b0e5104157ba5117d103dcc77de157871c0eb1f27684df9c12a24. The aggregate evaluated accepted base 75eae1ae11ad40db38341dc7d63f32876cba6483 plus this slice working tree. It is independent of the already published tag/package receipts. Later changes only finalize factual specification/status records; the final source audit allows only those named records and keeps all evaluated current guide, fallback, application, dependency, workflow and tagged release bytes unchanged. Final encoding/format/lifecycle and immutable-record checks follow the factual finalization.

The ten existing browser cases retain their assertions; an ignored owned copy changes only capture output paths and its relative selected-comparison fixture path. Six new captures live under .test-temp/s131/browser-captures, preserving earlier captures byte-for-byte. Seven one-off rendered checks confirm v2.0.0 fallback, comparison and current platform/install wording. All 15 staged binary metadata files independently select SQLite 1.61.1/libc 1.79.0. Attended native observations remain untested under the standing waiver.
