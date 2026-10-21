# Maintenance contracts

## Standing-waiver promotion

resolve selects exactly one authenticated release across all pages and freezes ID/content identity. download lists and fetches by frozen release/asset IDs, refusing unsafe/duplicate/incomplete assets; the initial inventory and digests constrain the final download. ensure-checksum reuses identical checksum bytes, permits absent checksum upload only on the still-matching draft, and refuses differing bytes or incomplete published recovery. publish revalidates identity/content and changes only draft=false plus make_latest=true for a draft or make_latest=true for a published release. Already-latest is a no-op. Latest confirmation matches both ID and tag. Content mutation, asset deletion/overwrite and independent tag re-resolution are prohibited by this correction.

The workflow retains successful staging-attempt, candidate, eight-asset allowlist, comparison-before-mutation, complete checksums and immutable-tag checks. The helper is built from exact workflow source before tagged checkout. Full attended-evidence promotion is unchanged.

## Watcher and payload contracts

Observer recovery preserves unchanged enabled revision's prior meaningful file state, pending replacement and dispatch identity; old definitions do not dispatch. Reconciliation and duplicate suppression remain bounded and fake-clock driven. Health/revision output stays compatible.

Script bindings preserve exact UTF-8 text and explicit null. Existing numeric conversions, exit0/1/2 categories and bounded retries remain. Concurrent success means independently persisted records. Actual failed-writer output is retained; passing repetitions never establish the historical intermittent cause.
