# S110 publication and bounded review record

## 2026-10-21: Authorized publication

Owner authorization covers pushing S110, opening PR #338, monitoring CI and completing up to two Codex review rounds, resolving findings autonomously and stopping before merge. The initial local implementation and first published head is 19df5010f697a9a9c4ff2cdedf1e8ab95f68f9b5. Authentication worked without owner involvement.

PR #338 covers #262, #271 and #335 in full. Existing issue milestones remain authoritative, and #335 remains under parent #277. The PR uses existing bug, build, daemon, packaging, testing and docs labels. Closing references are in its description; no issue was manually closed.

## Review and hosted verification

The PR was created as a draft while publication metadata was reconciled. Marking it ready automatically started round 1; no duplicate manual request was sent. Automatic and manual rounds count toward the same maximum of two. No merge or auto-merge is authorized.

Hosted CI and Codex review are separate evidence. The final handoff requires the actual final-head checks to reach terminal success; local checks, workflow aggregate state or earlier reviewed heads cannot substitute for those conclusions. The initial pending snapshot is superseded by the exact-head results below.

## 2026-10-21: Codex round 1 and hosted verification

Draft-ready round 1 reviewed 8badadee838c9d10f620162deee20720b254245e, starting at 2026-10-21T20:11:24.016843Z and completing at 2026-10-21T20:14:59.976364Z. [The bot summary](https://github.com/shruggietech/go-schedule/pull/338#issuecomment-6325502228) reported has_any_findings=false and has_unresolved_findings=false. The inline-thread inventory is empty. No source correction was needed.

All 22 checks on that same head reached terminal success, including native Linux/Windows/macOS race checks and qualification, desktop builds, packaging, docs and CodeQL. Exact runs: [CI 39668115276](https://github.com/shruggietech/go-schedule/actions/runs/39668115276) and [CodeQL 39668115324](https://github.com/shruggietech/go-schedule/actions/runs/39668115324). This is hosted evidence for 8badade, separate from the earlier local aggregate.

## 2026-10-21: Final publication-record head

The next commit updates only S110 delivery, review, verification and PR text. It changes no implementation source. The second and final allowed Codex round targets that publication-record head; actual request, completion and final-head CI identities are recorded in the final PR handoff comment. Never request a third round. No earlier review or check result is substituted for final-head CI conclusions. The PR remains open and unmerged, with auto-merge disabled and all three issues open pending merge.

## 2026-10-21: Codex round 2 and subsequent Windows CI correction

[The second and final request](https://github.com/shruggietech/go-schedule/pull/338#issuecomment-6325693641) targeted cb614525cfe650330a0a2847568069a11122ae2c. The bot summary completed at 2026-10-21T20:22:56.905517Z with has_any_findings=false and has_unresolved_findings=false. The inline-thread inventory is empty. Exactly two Codex rounds are complete; no third is requested, including for the CI correction below.

The final Windows race job on cb61452 failed TestReadHelperRejectsDirectEntryAndUntrustedLivePeer at context_test.go:123 with use of closed network connection after 2.01s. The other 21 checks succeeded. [Issue #339](https://github.com/shruggietech/go-schedule/issues/339) preserves the failed [CI job](https://github.com/shruggietech/go-schedule/actions/runs/39669525921/job/119463385548), independent of #262/#271/#335. The former all-green result on 8badade does not supersede this failure.

A forced legal schedule reproduces the failure against unchanged source: the helper correctly rejects its forged parent and closes before Windows Accept completes. go-winio discards the closed client, leaving Accept waiting until the fixture context closes its listener. The correction establishes both real native endpoints before invoking the existing peer checks. A private established-connection handler is extracted without changing kernel PID, actual parent, executable image, framing, cancellation or credential access checks. The direct-entry refusal and zero-backend-call assertions remain.

Correction checks: deterministic before case fails (2.164s); 100 native Windows race repetitions pass (43.694s); full Windows clientsecret package passes race (7.229s); full native Linux clientsecret package passes race (15.368s); format, vet and lint pass with zero lint issues. The fixture's existing two-second context remains unchanged. See [the defect record](findings/windows-peer-fixture.md) and [verification](verification.md). This source correction occurred after the two clean Codex rounds and is not represented as reviewed by a third round. The correction head receives fresh hosted checks, whose exact final outcome belongs in the PR handoff.

## 2026-10-21: Subsequent Linux reader-fixture correction

Head 5be817f033b9125e70d6db6842c9631c49514887 completed 21 successful hosted checks, including the corrected Windows race job, but failed TestScriptsReaderQueries/posix in [the Linux race job](https://github.com/shruggietech/go-schedule/actions/runs/39671902546/job/119471665372). The jitter reader returned usage exit 2 because no interval could be inferred. [Issue #340](https://github.com/shruggietech/go-schedule/issues/340), linked under #277, tracks this independent fixture outcome. The hosted seed timestamps were not retained, so their exact values or a particular clock cause are not established.

Twenty unchanged native Linux repetitions passed (33.792s), which is not a diagnosis. A controlled legal two-row seed with equal timestamps reproduces the old fixture's failure (1.681s). The correction retains the actual recorder loop and requires both rows, then assigns the reader's own consistent one-second timeline. Positive calls still omit an interval and must disclose inferred 1s; gaps joins the canned-query inventory. Equal-time and single-beat cases retain usage exit 2 and interval guidance, and explicit interval overrides must work. Failing queries now retain actual seed rows. No production reader or timestamp behavior changes.

Correction results: three native Windows race repetitions across both twins pass (80.008s); 20 native Linux repetitions pass (41.808s); the entire native Linux script suite passes race (12.688s). Formatting, vet and lint also pass. See [the issue record](findings/reader-cadence.md) and [verification](verification.md). The next correction head receives independent hosted checks, whose final exact-head results belong in the PR handoff. Both Codex rounds remain complete; no third review is requested.
