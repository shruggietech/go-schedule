# S116 source handoff

## Delivered scope and issue disposition

S116 delivers #372 and #373 together: increased shared frontend and cross-CI completion allowances, explicit 360-minute job ceilings, 60-minute Go package watchdogs, enforced nested-runner policy and 120-second/2-hour QA defaults with aligned cleanup reserves. Exact target, error, cleanup, resource, concurrency and intentional-expiry assertions remain. The owner explicitly superseded older no-increase guidance on 2026-10-25. Existing jobs already inherited GitHub's six-hour maximum; making that ceiling explicit does not increase it above the platform maximum.

#354, #355 and #356 are included with source-bound native follow-up, known-issue attribution and candidate release-content reconciliation. They remain open for their individual remaining functional outcomes. The original intermittent CLR/prerequisite and IPC causes are unconfirmed; successful repeats and a larger budget do not prove those causes fixed. #364 remains closed after its separately approved launcher explanation fix.

## Source, CI and publication identities

S115 PR #369 was squash-merged as 9e744e434f07a8d22041e0d4d6318867410900b6. The reviewed source tree matched that merge, but merged-main CI 40177500270 failed on the macOS frontend watchdog (#372). Exact-main CodeQL 40177500250 and website 40177500252 succeeded; website deployment 7310749121 belongs to that source. This failure remains part of the record, independently of later S116 results.

The selected pending release is v1.6.0. S115 and S116 delivered candidate scope now appears in CHANGELOG.md and .github/release-notes/v1.6.0.md, with refreshed comparison-review content digests. Capability decisions/dataset remain unchanged. The published release selector remains v1.5.1 (release 397490139). Source checks do not establish a v1.6.0 tag, package, staging or installed-product identity. No such publication operation was performed.

## Remaining release pathway

1. Publish the owner-authorized S116 review branch and official PR. Audit terminal checks against its exact head and resolve all warranted individual review threads, with no more than two Codex rounds.
2. Stop unmerged for the owner's final review. After owner merge, separately verify that exact merged-main CI and website deployment; PR-head success does not replace those identities.
3. Continue #355/#356 causal diagnosis as evidence warrants. Preserve original S115 reports and native runtime identities; file distinct defects when observations demonstrate distinct failures.
4. #354 public release requires the separately authorized tag/package/staging/promotion pathway and its actual identities. Standing owner-approved disclosure remains explicit in candidate notes. This slice does not tag, merge, stage, promote or change the public release selector.

## Completed local source acceptance

All eight canonical gates passed with native exit 0, 760 frozen source/test input hashes unchanged, 167 frontend tests passed and all six uncached coverage thresholds satisfied. The 23-document website export and timeout-policy/negative automation fixtures passed. This accepts local source implementation; hosted exact-head CI, reviews and later merged-main acceptance remain independently recorded.

## Verification reference

See [verification.md](verification.md) for original failure attribution, before/after experiments, targeted native Windows/Linux outcomes, explicit-budget behavior and canonical-suite results. Private scratch evidence is retained outside tracked source; no full-catalog rerun or installed-service action is claimed.
