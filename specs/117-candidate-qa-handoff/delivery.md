# S117 delivery and release handoff

## Delivered source scope

S117 bundled #354, #355 and #356 with independently demonstrated #375 and #376. It restores bounded human schedule descriptions through the authenticated observation API and MCP resources, and preserves valid non-object QA observations during diagnostic attachment. Each correction has separate failing-before/passing-after evidence in [verification.md](verification.md).

| Issue | Individual outcome | Current disposition |
| --- | --- | --- |
| #354 | Complete candidate evaluation, accurate known-issue disposition and pending source guidance | Open for public release/discovery |
| #355 | Native original-evidence/runtime/prerequisite follow-up; original cause unconfirmed | Open, needs reproduction |
| #356 | Full Windows/WSL observations and phase/request follow-up; original cause unconfirmed | Open, needs reproduction |
| #375 | Bounded scheduled-task descriptions restored with truthful truncation | Closed by merged PR #377 on 2026-10-25 |
| #376 | Object/array/scalar/null observation retention and malformed-state diagnosis | Closed by merged PR #377 on 2026-10-25 |

## Actual source delivery

[PR #377](https://github.com/shruggietech/go-schedule/pull/377) reviewed head 7327a7c401b6140d0b6c02d2769ea3082d74a123 passed all 22 checks, including [CI 40197369910](https://github.com/shruggietech/go-schedule/actions/runs/40197369910) and [CodeQL 40197369870](https://github.com/shruggietech/go-schedule/actions/runs/40197369870), attempt 1. Exactly two Codex rounds completed without findings. The requested security-focused second round was recorded as ManualCodeReview with no dedicated security status; it is not represented as a separate security certification. See the [second-round receipt](https://github.com/shruggietech/go-schedule/pull/377#issuecomment-6398229982) and [final unmerged handoff](https://github.com/shruggietech/go-schedule/pull/377#issuecomment-6398304444).

The owner squash-merged at 2026-10-25T17:48:41Z as 9ad13821c825d01efcc038b6b1634b187a854758, with the same tree as the reviewed head. All 23 exact merged-main checks passed: [CI 40218254702](https://github.com/shruggietech/go-schedule/actions/runs/40218254702), [CodeQL 40218254649](https://github.com/shruggietech/go-schedule/actions/runs/40218254649) and [Website 40218254658](https://github.com/shruggietech/go-schedule/actions/runs/40218254658), all attempt 1. Pages deployment 7316855350 identifies that exact merge; normal HTTPS checks passed the home, comparison, documentation, CLI reference, search index and linked stylesheet routes. These are source/website receipts, not release-package identities.

## Preserved candidate and release identities

The original capture remains bound to dirty planning-only 11485a8a6ce97fe8e6966e1cc3cdb6a55b9d5651, source digest 98514d388059638cf08b00ee5c93e82fca24d4437ccca426172f822200082d18 and development version 0.0.1-dev. Its final report remains failed for six captured communication records; later corrections and successful checks do not relabel it.

Pending source version is v1.6.0; public selection is v1.5.1. Selected changelog/notes hashes and capability decisions remain as recorded in verification.md. The schema-v24 stopped-daemon complete-state backup and matching-binary rollback instructions are unchanged.

## S118 reconciliation and current handoff

On 2026-10-25, S118 corrected the missed source lifecycle bookkeeping after the owner merge. Historical publication tasks T016/T017 now appear as workflow receipts; the source reconciliation portion of T018 is explicitly attributed to S118. S117 implementation is complete independently of later release publication.

The current pathway and owner boundaries are in [S118 handoff](../118-release-records/handoff.md). No S117 source receipt establishes a future tag, draft, asset or public v1.6.0 release. #354/#355/#356 retain the separate open outcomes above.
