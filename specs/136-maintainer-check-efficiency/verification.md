# S136 verification

## 2026-11-08 unchanged baseline

Accepted main 5e4f4f4ca373ca31b607f7d456fb4d3d3339a638: native Windows lifecycle 31.541 seconds; automation 44.701 seconds, both successful. Automation includes lifecycle. Original source and timing receipts are retained privately under .test-temp/s136.

## Scope and preservation

Parent #433 includes parser children #434/#435. Original #355/#356 remain unresolved; clustered scheduling stays deferred. No release or installed service changes. Prior S135 frozen inventories and housekeeping receipt will be audited before and after candidate verification. New source qualification and factual completion edits will be identified separately.

## Implementation and qualification

Focused native race tests, both command/root modes, action-field/alias/malformed/duplicate cases, publication dependency/order cases, child failure retention, the lifecycle corpus, all 93 automation assertion cases and aggregate failure propagation passed. The first native test run failed to compile before the new implementation existed; that scaffold observation is not a reproduced semantic bug. The original comment action scan failed with the unaudited-action finding, and the original disconnected-publication fixture returned zero; those real behavioral observations are retained separately.

On the same Windows host with resolved Go modules, wrapper-inclusive lifecycle improved from 31.541 to 0.496 seconds and automation from 44.701 to 0.624 seconds. The new inventory contains 134 specifications instead of 133. Automation includes lifecycle. These single before/after observations do not attribute the entire historical S135 aggregate or promise a future aggregate duration. Native evaluation benchmarks: 133-spec lifecycle 15.277113 ms/op, 840228 B/op, 6853 allocations; small workflow evaluation 31.432 microseconds/op, 12012 B/op, 99 allocations. Benchmarks exclude wrapper compilation and brand/manifest helpers.

All 577 selected original/historical inputs match their preserved sizes and hashes before canonical qualification. The full eight-gate foreground aggregate is pending on the frozen implementation; final factual completion records will be checked separately.

## First canonical attempt

The first foreground aggregate stopped at lint after 11.132 seconds for local import grouping and a tagged-switch style finding. No later gate ran, source drift was zero and all 577 frozen inputs were unchanged. Its failed result is retained. These are new-code style findings, not a passed canonical run. The source grouping/switch corrections follow that observation. CI and CodeQL now explicitly check out the PR head for exact-candidate evidence without changing triggers, permissions or gate scope.

## Final corrected implementation qualification

Canonical attempt 02 passed all eight gates in 663.266 seconds with 2631 source files, zero source drift and all 577 preserved inputs unchanged. During that frozen run, separate isolated negative probes demonstrated #436 in pending native source: near-match dependency labels, update types and module patterns were accepted. The original failing probes remain retained. The correction uses exact decoded sequence membership and rejects repeated Makefile pins; its permanent native race regressions pass. Three added Makefile cases increase the existing 93-case automation corpus to 96 without removing a case.

Canonical attempt 03 passed all eight gates in 482.232 seconds on 2632 frozen source files, zero drift and all 577 preserved inputs unchanged. Format, vet, lint, root race, Wails Go/native build/frontend, coverage, docs and automation are successful. All 168 frontend and seven website tests pass; 25 documents/30 routes export. Core coverage remains engine 83.6, schedule 89.1, timezone 91.3, store 80.9, catchup 88.9 and logbus 91.1 percent. Lifecycle corpus and aggregate failure propagation also pass. Wrapper smoke observations cover default and caller-relative roots, spaces and source-module selection.

The final wrapper-inclusive native measurements are lifecycle 0.444 seconds and automation 0.623 seconds, compared with unchanged-source 31.541 and 44.701 seconds. The new inventory includes S136 (134 vs 133 specifications). These are individual observed runs with overlapping scopes, not a statistical or whole-suite promise. Earlier optimized observations and the historical S135 aggregate remain separately identified.

After implementation qualification, factual status/delivery/traceability records and CI artifact-origin metadata are finalized. A source-config regression requires CI/CodeQL checkout refs and Windows qualification origin labels to select the same PR head, with github.sha fallback for non-PR events. The initial new guard failed before the label was completed; the corrected configuration and new guard are checked by focused native race, formatting, lint, lifecycle and full automation verification. No application implementation changes follow canonical attempt 03. Hosted checks will independently qualify the official PR head.

Released-MSI installation and attended native GUI remain untested; S136 changes maintainer tooling, not those product surfaces. No new application QA catalog replay, installed service change, tag, staging or promotion is performed.

## Final record/configuration checks

Focused native race/provenance tests, publication formatting, all 134 lifecycle entries, lint, the full automation gate with all 96 assertion cases and failure propagation passed after record/configuration finalization. Final benchmarks on this source: 133-spec lifecycle 14.563420 ms/op, 835742 B/op, 6853 allocations; small workflow evaluation 30.927 microseconds/op, 12012 B/op, 99 allocations. Earlier benchmark observations remain historical. Post-canonical differences are the candidate-origin configuration/guard and factual status, traceability and delivery records; the application and native policy implementation are unchanged.
