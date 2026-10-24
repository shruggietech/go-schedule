# S115 source handoff

## Scope and current publication state

Owner-authorized source work covers #354, #355, #356, 18e6fd8 and independently demonstrated #366 and preparation correction #367 and timeout-fixture correction #368, plus message corrections #370/#371. Carry-in reconciliation is commit dde3e7b on the S115 review branch; original 574d071 remains archived. Pending source is v1.6.0, while public selection remains v1.5.1. No v1.6.0 tag, staged package hash, release ID or publication is supplied by this source work.

## Pre-tag source operation

Run `go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.6.0`. The same command runs in Release after Go setup and before comparison. It reads local files, rejects inconsistent badge/health/heading/notes metadata and requires a calendar-valid release date. Success does not verify a real Git tag, exact-main CI, comparison content or package identity; those retain their existing separate checks.

The dated-heading defect is #366. Actual pending notes remain linked to the unique 2026-10-22 heading. Changing that date or release scope requires coordinated notes-anchor and comparison content-digest updates before tag selection.

## Native outcomes and issue traceability

See [observations.md](observations.md) for actual source/runtime-bound results. #355 retains the original independent CLR crash and Bash prerequisite refusal, and #356 retains the original owned IPC/observation failures. Passing current observations do not establish those causes. #366 closes only on its delivered validator correction after merge. #354 stays open for the actual public-release/discovery outcome.

## Later immutable staging and promotion

After source review, owner merge and terminal successful exact-main push CI, obtain separate owner authorization for one exact v1.6.0 tag/source. Inspect existing tag/release state first; never replace or move a tag. The existing tag-only Release workflow stages four daemon/CLI archives (Linux amd64/arm64, Darwin amd64/arm64), two desktop archives (Linux amd64, Darwin arm64), Windows amd64 MSI and windows-candidate-manifest.json. Retain actual run/attempt, authenticated release ID, full allowlist, MSI identity/hash and immutable source resolution. A local build is not a substitute staged candidate.

Present actual staged identities for separately authorized standing-waiver promotion through the reviewed workflow. Keep the same release ID, candidate bytes and exact checksum inventory across interrupted retries. Missing/conflicting assets or identities stop the operation; do not rebuild, delete or substitute a candidate to bypass a mismatch. The standing waiver and required untested disclosure remain unchanged.

For schema-v24 upgrade, stop the daemon and retain a complete application-state backup before switching binaries. Rollback restores that complete backup together with matching older binaries; do not point older binaries at upgraded storage. Existing platform installation guidance remains authoritative.

After actual publication, verify release ID/tag/source, latest pointer, complete downloads and checksums, then reconcile the source-owned published selector and website discovery through review. Preserve historical comparison bytes. This S115 kickoff performs none of those merge/tag/staging/promotion operations.

## Frozen evaluation and later correction identities

PR #369 carries the source work. The full owned report is permanently bound to clean fc39261ed8a874aaad11dd83f7afcccaf06fb10d and content SHA256 74da617216a745dc038c21a727161fe5f4b4f3c825b84283d096b85d562f4208. It is a validated failed outcome: execution 1,760 pass/3 fail, effective 1,744 pass/19 fail/0 blocked/3 platform exclusions, 1,109 exact-bound message reviews and three findings. All four runtime cleanups completed. Three metadata deadlines remain under #356; #370/#371 have subsequent source corrections with separate verification. Never relabel that catalog as a pass or as the later PR head.

Merge may close only the independently delivered #366/#367/#368/#370/#371 outcomes. #354 remains open until public release/discovery is delivered; #355/#356 remain open with original causes unconfirmed. Actual exact-head CI/review receipts belong to the PR timeline. No source correction supplies a staged package or publication identity.
