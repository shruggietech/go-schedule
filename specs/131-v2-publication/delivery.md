# S131 Delivery

## Actual public release

[v2.0.0](https://github.com/shruggietech/go-schedule/releases/tag/v2.0.0), release **430039348**, became public/latest on **2026-10-31T16:16:04Z** from immutable source **75eae1ae11ad40db38341dc7d63f32876cba6483**. The owner separately approved exact tag/draft staging and promotion of the inspected concrete inventory. Existing [Release run 41218081205](https://github.com/shruggietech/go-schedule/actions/runs/41218081205), attempt 1, passed all seven jobs. [Promotion 41218776827](https://github.com/shruggietech/go-schedule/actions/runs/41218776827), attempt 1, completed successfully through the existing standing-waiver workflow. No accepted asset was rebuilt.

[The exact candidate inventory](https://github.com/shruggietech/go-schedule/issues/412#issuecomment-6545923386) records all eight original asset IDs, sizes and hashes. Independent public re-download verified those same eight originals plus SHA256SUMS.txt (asset 695846859), for nine files/eight checksum entries. All original IDs/digests are unchanged. MSI ProductVersion 2.0.0 and ProductCode {FD066B26-9F5F-4A16-A3FA-BAC7ECDB75BD} agree with the manifest and successful Windows staging attempt. Six archive inventories and all 15 embedded release versions were inspected. Available embedded source revisions match; Wails GUI omits VCS metadata and existing packaging output marks some other binaries modified, so no clean-worktree binary claim is made.

## Discovery and current-source reconciliation

Publication-triggered [website 41218860788](https://github.com/shruggietech/go-schedule/actions/runs/41218860788), attempt 1, succeeded from the accepted source; Pages deployment 7461113991 serves public v2.0.0 comparison selection. HTTPS home, comparison, Windows installation and unauthenticated manifest/checksum download URLs returned 200. The deployed Windows guide still contains its pre-publication v1.7.0 status sentence until the S131 source PR merges; source fallback and current README/platform/install/website/maintainer wording are reconciled in this branch. Public package delivery and final current-guide integration remain distinct.

[#412](https://github.com/shruggietech/go-schedule/issues/412) and milestone 19 remain open for that final source integration. The official S131 PR will close #412 on owner merge only after the source reconciliation is ready. The inherited S129 release continuation is executed here with actual receipts; no new release framework or preparation-only outcome is introduced.

## Verification and review boundary

All eight foreground gates passed in 4534.339 seconds with zero drift across 2557 source paths and all 221 retained inputs; ten headless site browser cases and seven rendered-publication checks passed. See [verification](verification.md). Exact final head/check/review receipts belong in the private .test-temp/s131/publication-ledger.json and PR description/comment, avoiding a source-only evidence commit that would change the reviewed head. Automatic PR-open review consumes request one; at most one manual second request is authorized. The PR stays unmerged for the owner.

Attended native observations remain untested under the standing waiver; no active-workstation installation or lifecycle acceptance is claimed. Original #355/#356 causes remain unconfirmed. Clustered execution remains owner-deferred. Prior public v1.7.0 release metadata/asset IDs/digests and all 221 retained inputs plus three newly tagged source records remain unchanged. No scheduled monitor is created.
