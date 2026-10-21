# [bug] Observe created notification channel identities in human-output QA cases

Publication state: local draft awaiting the S109 pre-publication review. Proposed existing taxonomy: bug, area: testing, area: cli, effort: S, platform: windows, platform: linux. Milestone: 15, CLI completeness and unattended application QA. Parent: #277. No assignee, deadline or release claim.

## Outcome and bounded scope

The unattended evaluator recognizes the independently persisted identity created by a human-mode notification channel add command and carries it into dependent channel/policy cases. It retains the existing stdout-identity, ambiguity and persisted-state assertions.

## Authority and source

Discovered during the owner-authorized full S109 #311-#322 bundle, while validating #314's now-working human webhook authoring. The existing mapping in [internal/clievaluate/human.go at merged baseline c6d50a1](https://github.com/shruggietech/go-schedule/blob/c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd/internal/clievaluate/human.go) omits notification channels. The CLI operation succeeds, but the evaluator looks for its created identity under an empty snapshot entity key. The evaluator's prior recovery of failed channel authoring obscured this independently actionable defect.

Affected pending source: c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd with dirty=true and exact source digest 96d22c4425a44e11ea4e53b35bc93009ceb06a2e41d4140fa425f10d1c99e118. Diagnostic run a1524141-b52e-4e1d-85dd-86f57bd86efb is retained under .test-temp/s109-run-01. This is a pending-source observation, not a claim about a tagged release. Windows demonstrates the defect; the mapping is shared with Linux.

## Reproduction and observed effect

1. Use the existing S108 scoped Windows/selected-WSL capture with S109's #314 endpoint correction and fresh receipt-bound inputs.
2. Inspect notification.channel-create.human. The actual CLI exits 0 and prints the newly persisted channel ID.
3. The evaluator reports failed required-output and persisted-state checks because humanResponse cannot select the channels snapshot entity. Dependent channel/policy cases become fixture_prerequisite blocked results.

Representative captured stdout: .test-temp/s109-run-01/attempt-001/windows/cases/notification.channel-create.human/stdout.txt, 125 bytes, SHA256 e7879a9c9fc231f03d9bb956539aba681121d599e0f64416826ef6862dc54653. The artifact is fixture-only and retained locally with the canonical diagnostic report.

Expected: select exactly one newly persisted channel identity, confirm it appears in actual stdout, save that identity and continue dependent cases.

Actual: the mapping covers tasks/groups/chains/actors/pairings/triggers/watchers but omits channels. The failure is in QA interpretation, not in the successful CLI mutation.

## Functional acceptance criteria

- [ ] Human channel creation selects its identity from the independently observed channels snapshot.
- [ ] Unchanged or ambiguous creation state remains refused, and stdout must still identify the selected object.
- [ ] Dependent human channel and notification policy cases receive the actual created identity.

## Dependencies and implementation

Depends on merged #289/#290 for the evaluator. #314 makes this successful path reachable in the normal full suite. The correction is implemented locally in S109 using a shared command-to-snapshot identity helper and a regression covering a new channel alongside an unrelated new task, unchanged state and ambiguous channel creation. Publish this independent bug record with S109, then link the implementation; do not reopen delivered #289/#290 or weaken the product assertions.
