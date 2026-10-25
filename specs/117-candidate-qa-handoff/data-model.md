# S117 evidence relationships

| Entity | Identity and relationship | Invariant |
|---|---|---|
| Candidate | Commit, dirty flag, source digest, six binary digests, build receipt | Frozen through capture and finalization |
| Run | Run ID, manifest, budget, runtime identities, evidence manifest | Every applicable case and exclusion accounted for |
| Case | Environment, command/scenario/mode, command exit, inspection, diagnostics, cleanup | Independent outcomes remain attributable |
| Message review | Run/corpus digest, exact message ID/text/channel/cases, seven rubrics | Actual contextual review, no inherited verdict |
| Finding | Dedup key, affected revision/platform/cases, hashed evidence, functional outcome | Existing issue reused where equivalent |
| Correction | Separate source identity, failing-before and passing-after regression | Never relabels original capture |
| Handoff | Pending source claims, final commit/check/review identities | No implied merge, tag, package or public release |

State progression: identified inputs, capture, pending review, finalized pass/fail/blocked, immutable report, separate corrective acceptance, source PR handoff. An integrity error is not a finalized evaluation outcome. An interrupted capture follows only the existing identity-preserving recovery pathway.
