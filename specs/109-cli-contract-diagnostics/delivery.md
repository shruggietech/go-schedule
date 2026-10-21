# S109 local delivery and publication handoff

S109 delivers the full #311-#322 bundle under parent #277 and milestone 15 on codex/109-cli-contract-diagnostics. The synchronized baseline is merged main c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd (PR #310).

## Delivered behavior

Shared CLI invocation contracts classify parser and operand failures before target resolution, keep operational failures distinct, and preserve generated help/completion. Alerts keep structured stdout clean. Daemon targeting is independent of webhook destination via --daemon-endpoint, with legacy targeting retained and explicit conflicts refused. Bundle preview refuses duplicate and unknown watcher bindings. Cron import keeps the complete report and returns usage failure for invalid lines.

Credential, pairing and service help describes actual effects. Human lists label each column. Diagnostics identify relevant objects, fields, files and useful corrections while retaining causes. Credential rotation assembles actual actor and daemon context without authentication side effects, with a transaction retaining the old credential on response-assembly failure. Watcher API, CLI and desktop responses distinguish saved definitions from observations and pending reconciliation. IPC compatibility warnings explain admin_group and restart configuration.

The specification, plan, contracts, 29 completed tasks, issue traceability and detailed verification live in this feature directory. The refreshed inventory has 152 initialized command nodes and 844 scenarios per environment. An additional evaluator channel-identity mapping defect found during real capture is corrected locally and has a separate [issue publication draft](findings/human-channel-identity.md).

## Decisions and deviations

Existing Cobra contracts, typed causes and additive selectors avoid a second parser or a flag-breaking change. Exact durable watcher revisions distinguish coalesced edits even when timestamps or health are unchanged, without a reconciliation wait. The existing rotation response shape is populated accurately, with no persisted schema, new eligibility restriction or dependency. The existing scoped S108 QA runner and fixed report schema remain authoritative; assertions were retained. No functional issue outcome is intentionally deferred.

## Verification

All eight foreground local gates passed; the final Windows/WSL canonical evaluation passed 1,685 of 1,688 planned combinations with three manifest-defined exclusions, zero failures/blocked cases, 1,049 complete contextual reviews and no remaining owned resources. Focused native Linux race checks and Linux-target lint passed. Finalize/validate/render passed and repeated Markdown rendering was byte-identical.

The evaluation uses the receipt-bound dirty-source identity, not a claim that a later documentation-only delivery commit was its captured subject. See [verification.md](verification.md) for exact source/binary identities, artifact hashes, earlier corrective runs and platform boundaries. Hosted CI and publication review follow authorization.

## Publication boundary

Stopped before publication as requested. Nothing has been pushed, no S109 PR or new defect issue has been published, and no merge, tag or release is performed. Existing #311-#322 remain open until their functional corrections merge. No authentication or other owner involvement was required during local implementation.

The exact review-branch commands awaiting authorization are:

```powershell
git push -u origin codex/109-cli-contract-diagnostics
gh pr create --repo shruggietech/go-schedule --base main --head codex/109-cli-contract-diagnostics --title "feat(109): repair CLI contracts and diagnostics" --body-file specs/109-cli-contract-diagnostics/pull-request.md --milestone "CLI completeness and unattended application QA"
```

At publication, create the separately drafted evaluator defect with the verified existing labels and milestone, add its real issue identity to the PR closing references and parent #277 child checklist, and retain the issue-level functional outcomes. The PR draft already includes closing references for all 12 authorized product issues.
