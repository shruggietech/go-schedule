# S109 local delivery and publication handoff

S109 delivers the full #311-#322 bundle under parent #277 and milestone 15 on codex/109-cli-contract-diagnostics. The synchronized baseline is merged main c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd (PR #310).

## Delivered behavior

Shared CLI invocation contracts classify parser and operand failures before target resolution, keep operational failures distinct, and preserve generated help/completion. Alerts keep structured stdout clean. Daemon targeting is independent of webhook destination via --daemon-endpoint, with legacy targeting retained and explicit conflicts refused. Bundle preview refuses duplicate and unknown watcher bindings. Cron import keeps the complete report and returns usage failure for invalid lines.

Credential, pairing and service help describes actual effects. Human lists label each column. Diagnostics identify relevant objects, fields, files and useful corrections while retaining causes. Credential rotation assembles actual actor and daemon context without authentication side effects, with a transaction retaining the old credential on response-assembly failure. Watcher API, CLI and desktop responses distinguish saved definitions from observations and pending reconciliation. IPC compatibility warnings explain admin_group and restart configuration.

The specification, plan, contracts, 29 completed tasks, issue traceability and detailed verification live in this feature directory. The refreshed inventory has 152 initialized command nodes and 844 scenarios per environment. An additional evaluator channel-identity mapping defect found during real capture is corrected locally and has a separate [published issue record](findings/human-channel-identity.md).

## Decisions and deviations

Existing Cobra contracts, typed causes and additive selectors avoid a second parser or a flag-breaking change. Exact durable watcher revisions distinguish coalesced edits even when timestamps or health are unchanged, without a reconciliation wait. The existing rotation response shape is populated accurately, with no persisted schema, new eligibility restriction or dependency. The existing scoped S108 QA runner and fixed report schema remain authoritative; assertions were retained. No functional issue outcome is intentionally deferred.

## Verification

All eight foreground local gates passed; the final Windows/WSL canonical evaluation passed 1,685 of 1,688 planned combinations with three manifest-defined exclusions, zero failures/blocked cases, 1,049 complete contextual reviews and no remaining owned resources. Focused native Linux race checks and Linux-target lint passed. Finalize/validate/render passed and repeated Markdown rendering was byte-identical.

The evaluation uses the receipt-bound dirty-source identity, not a claim that a later documentation-only delivery commit was its captured subject. See [verification.md](verification.md) for exact source/binary identities, artifact hashes, earlier corrective runs and platform boundaries. Hosted CI and publication review follow authorization.

## Publication and merge boundary

The initial pre-publication halt completed at local implementation commit 34ab5cb544ac6bb2671dde30238163e213120d4d. On 2026-10-21 the owner authorized pushing S109, opening the PR, publishing the additional defect, monitoring CI and completing up to two Codex review rounds with autonomous corrections.

[PR #332](https://github.com/shruggietech/go-schedule/pull/332) is published for the complete #311-#322 bundle plus independently tracked evaluator defect [#331](https://github.com/shruggietech/go-schedule/issues/331). All issue records remain open pending merge. The [publication and review record](review.md) records exact reviewed revisions, corrections and hosted checks. Stop before merging; no merge, tag or release is authorized.

Two Codex rounds are complete. The final reviewed correction head is 4db6501b4ec3d7256609525288e52a4b37287e21, and round 2 has no findings. The two first-round threads are resolved with the documentation correction and exact-output evidence for the unreproduced version report. CI-discovered macOS fixture correction #333 and documentation correction #334 are included with the original bundle and #331. The final exact-head hosted-check result is recorded in the PR handoff after the publication-record commit.
