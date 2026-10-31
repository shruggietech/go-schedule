# S129 Release and Review Contract

## Current source boundary

One official PR targets main from codex/129-prelaunch-release and uses Refs #412. Complete local foreground verification, exact-head CI, all received findings and at most two Codex requests including automatic PR-open review. A terminal service receipt or correctly attributable post-request bot thumbs-up can establish a clean result; pending/eyes-only receipts cannot. Report the actual review mode and any absent separate security receipt. Stop before owner merge; no recurring monitor.

## Candidate source identity

Prepare v2.0.0. Root Go module is github.com/shruggietech/go-schedule/v2 and nested locally built desktop module is github.com/shruggietech/go-schedule/v2/desktop. The desktop requires root v2.0.0 with its existing local replacement. CLI command names, HTTP /api/v1, persisted current formats and development version fallback are unchanged. Existing release linker symbols use the /v2 buildinfo path.

README distinguishes its prepared-source health example from actual available v1.7.0 downloads. The new comparison/review are bound to v2.0.0; site/data/published.json and immutable historical records stay unchanged until actual promotion.

## Later exact-source and asset boundary

After owner integration, obtain clean synchronized main and successful exact-main CI, CodeQL and website results. Run existing source and comparison guards against the actual main SHA. Obtain separate approval naming that SHA and immutable v2.0.0 tag before creating/pushing it. Existing Release workflow stages exactly these eight assets:

- go-schedule_v2.0.0_darwin_amd64.tar.gz
- go-schedule_v2.0.0_darwin_arm64.tar.gz
- go-schedule_v2.0.0_linux_amd64.tar.gz
- go-schedule_v2.0.0_linux_arm64.tar.gz
- go-schedule-desktop_v2.0.0_darwin_arm64.tar.gz
- go-schedule-desktop_v2.0.0_linux_amd64.tar.gz
- go-schedule_v2.0.0_windows_amd64.msi
- windows-candidate-manifest.json

Record the actual Release run/attempt and seven job conclusions, immutable source/tag, release ID, all asset IDs/sizes/digests, archive contents, embedded versions and MSI ProductVersion/ProductCode against the manifest. No installed-product or attended native acceptance is inferred.

Present the actual inventory for separate concrete promotion approval. Then use existing promote-release-waived.yml from reviewed main with the exact tag. It must publish the same assets and add/reuse exact SHA256SUMS.txt, yielding nine files. Never rebuild after acceptance, move a tag or replace assets. Existing refusal and idempotency behavior governs retries.

## Public completion

Verify actual public/latest release ID/tag/source and unchanged original asset IDs/bytes; verify checksum inventory, Pages exact-source workflow/deployment and HTTPS home/install/comparison/download selection. Reconcile current discovery wording through normal PR workflow within S129 continuation. Keep #355/#356, #400/#404 and cluster work open/deferred. Attended observations remain untested under the standing waiver.
