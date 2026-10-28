# S125 Research

## Focus sequencing

Decision: wait for main-content route focus before focusing and entering the query, mirroring the previously delivered #385 unit-fixture synchronization. Shell.selectRoute queues requestAnimationFrame; the browser fixture currently races that handoff. Preserve the original b4f1e28 failed receipt, not just later passing repetitions.

Alternatives considered: changing product focus, increasing timeouts, blind input retries. None is supported by the observed fixture failure. Use controlled animation-frame ordering for a meaningful regression; Playwright supports [clock control](https://playwright.dev/docs/clock), including requestAnimationFrame. Retain generated failure context/traces using the existing output directory and artifact action.

## Release delivery

Decision: reuse S124 prepared v1.7.0 metadata and existing release/promotion guards. The current red main b4f1e281e66003e08ee2a7abd35e3e25468fb6bd cannot be selected as a release source. The S125 merge and assets do not yet exist. Existing workflow stages eight assets in seven jobs; the standing-waiver promotion adds SHA256SUMS.txt, rechecks immutable source and concrete asset identities, and publishes the same release.

Alternatives considered: new release tooling or another full application/historical campaign. Existing infrastructure already covers the unchanged application; source correction and a concrete continuation handoff are sufficient for this authorized pre-merge boundary. Public selection remains v1.6.0 until actual approved promotion.

## Distributed demand

Decision: reversible hold on architecture/prototype activation, retaining all incomplete outcomes. [Owner intent](https://github.com/shruggietech/go-schedule/issues/20#issuecomment-5073941594) describes activity endpoints, coordination and reliability priorities, but not a concrete user/job, bounded host/task counts, availability target or duplicate/missed tolerance. #186/#189 have no actual-use reports. The [#174 delivery record](https://github.com/shruggietech/go-schedule/issues/174#issuecomment-5787231012) leaves demand to #186; [#185 publication](https://github.com/shruggietech/go-schedule/issues/185#issuecomment-5796659033) proves shipping, not subsequent unmet demand.

The public issue search covered open/closed cluster, distributed, failover and ownership records plus relevant roadmap comments. Discussions are disabled. This limited search cannot prove absence of users or need. Preserve the distinction.

Alternatives considered: inventing scenarios and proceeding, permanently cancelling the proposal, or blocking all independent release work. A documented hold plus exact missing-input intake preserves ambition and the existing gates with the least unsupported commitment.

## Existing alternatives

Use the implemented independent remote operations, All Systems/search, portable bundles and local downtime catch-up contracts. None silently delivers shared execution ownership. Reuse the reviewed 2026-10-27 highavailability comparison references instead of adding unsupported competitor claims. See the delivered demand document for limitations and reopening criteria.
