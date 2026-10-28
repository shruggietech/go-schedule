# S125 Verification

## Source and scope

Accepted main: b4f1e281e66003e08ee2a7abd35e3e25468fb6bd, owner-merged S124 PR #396. Its independent merged-main failure is [#397](https://github.com/shruggietech/go-schedule/issues/397), [job 122591575501](https://github.com/shruggietech/go-schedule/actions/runs/40620037771/job/122591575501). The original log and prior application QA remain retained. This proposal changes the browser fixture, CI failure retention and release/decision documentation; Shell, daemon, scheduler and application QA implementation are unchanged.

## Before/fix browser evidence

The controlled original sequence pauses animation frames, opens Search, focuses the query and then runs the queued route frame. It failed as intended with the input inactive (one case, exit 1, 30.4 seconds). The corrected ordering completes the route/main handoff before focusing and filling the query; one case passed, exit 0, 2.5 seconds. No arbitrary sleep, retry, timeout reduction or assertion removal was introduced.

The strengthened existing browser suite passed all 35 cases, one worker, 28.8 seconds. The shipped fixture retains keyboard focus, compact viewport, 200 percent zoom, filter visibility and horizontal-overflow assertions, and additionally checks the exact archive value. Shell SHA-256 remains b83e754acd3847066255f1e3e1c91a5f65b04daf20172cf21a321aeaa7330e8f.

The private receipt .test-temp/s125/focus-regression-receipt.json binds original/fixed source sequences, logs and three preserved original failure artifacts, including error context and trace. The comparator is controlled causal evidence for the fixture sequence, not a new product focus policy or diagnosis of #355/#356.

## Release and documentation checks

Existing source metadata and comparison guards passed against the pending source content. Supplied b4f1e28 is historical identity context for that local guard exercise, not approval of a tag or an accepted release source. Content-bound v1.7.0 review hashes were refreshed after the pending notes/changelog correction; the 80 capability cells and public v1.6.0 selection are unchanged.

The corrected docs gate passed all seven website tests and the 24-page production export/link audit. Headless visual review of the exported assessment passed at 1440px dark and 390px light viewports, with zero document overflow; the existing brand/navigation and readable content layout remain intact.

Initial GitHub formatting and git diff --check passed. All generated source is UTF-8 without BOM; the final factual-finalization checks cover encoding/corruption, lifecycle and publication formatting separately from the aggregate. Native issue comments record the seven issue dispositions with no unsupported closure.

## Canonical acceptance

The first complete foreground aggregate exited 1 after 558.275 seconds at the docs gate: the new distributed assessment lacked required YAML front matter. All preceding gates passed, including all six core coverage thresholds (engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.6%, catchup 88.9%, logbus 91.1%). All 2463 source paths and 151 frozen inputs remained unchanged. The original receipt and log are retained. Independent introduced-source defect #398 records the metadata correction; no documentation guard is weakened. Corrected-source acceptance follows below after completion.

## Authorization and limits

Push/PR and at most two Codex rounds are explicitly authorized. Merge, tagging and public promotion remain excluded. No actual S125 merged-main candidate or v1.7.0 staged asset identities exist at this checkpoint. The original #355/#356 causes remain unconfirmed; no unchanged historical campaign is repeated.

## Corrected complete aggregate

The second complete foreground aggregate passed all eight gates on 2026-10-28: 933.338 seconds, exit 0, no signal. All 2463 frozen source paths and all 158 frozen evidence inputs were unchanged. The frozen set includes the 137 prior QA/binary inputs, original browser failure/context/trace, first failed aggregate and corrected docs/visual receipts. Receipts: .test-temp/s125/canonical-inputs-02.json and canonical-receipt-02.json; complete output: canonical-all-02.log. The first failed aggregate remains distinct and retained.

Factual finalization updates only spec.md, tasks.md, verification.md, delivery.md and the specs/README.md inventory after that aggregate. No application, fixture, workflow, release content or demand requirement changes in this finalization stage. Formatting, UTF-8/no-BOM/corruption, lifecycle, reviewed release-content guards and unchanged frozen inputs are checked separately before commit. Hosted exact-head CI and external review are additional publication evidence, not a claim of merged-main or public release acceptance.
