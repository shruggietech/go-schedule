# Tasks: S129 Prelaunch Baseline Publication

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/release.md and quickstart.md.
**Branch**: codex/129-prelaunch-release. Native outcome: #412, milestone 19.

Current authorization covers T001-T017 and the source publication workflow below. T018-T022 remain this same slice's later owner-controlled continuation and MUST NOT execute without their stated approval. Source readiness is not public delivery. Tests validate engineering compliance; #412 closure depends on the actual public functional outcome.

## Phase 1: Setup

- [X] T001 Confirm branch/spec/selector identity and installed hooks in .specify/feature.json and specs/129-prelaunch-release/plan.md.
- [X] T002 Preserve 215 historical inputs and original private receipts in .test-temp/s129/protected-inputs.json; leave historical release/comparison files unchanged.

## Phase 2: Foundational gate

- [X] T003 Run installed read-only speckit-analyze after task generation; record its coverage and any resolved findings in specs/129-prelaunch-release/analysis.md before source implementation.

## Phase 3: US1 - Coherent candidate identity (P1)

**Independent test**: Actual source guard selects v2.0.0, module graphs resolve locally, a disposable tag-injected CLI reports 2.0.0, and published selection stays v1.7.0.

- [X] T004 [US1] Align root go.mod, desktop/go.mod, active Go imports, release.yml linker symbols, .golangci.yml local prefix and test/scripts/coverage-profile_test.sh fixture with the /v2 namespace; preserve dependencies and repository URLs.
- [X] T005 [US1] Prepare v2.0.0 CHANGELOG.md section, .github/release-notes/v2.0.0.md and explicitly unpublished README.md health example; record major-version and pinned-artifact decisions.
- [X] T006 [US1] Update current candidate wording in docs/install.md, docs/INSTALL-windows.md and test/windows/README.md while retaining actual public v1.7.0 availability and historical records.
- [X] T007 [US1] Check both read-only module graphs/checksums, actual release-source-check and disposable version-injected build against the design in quickstart.md; retain exact receipts in .test-temp/s129.

## Phase 4: US2 - Accepted-byte release continuation (P1)

**Independent test**: handoff.md names all eight staged assets, separate approval boundaries, unchanged-byte promotion and later public readback without claiming any release exists.

- [X] T008 [P] [US2] Write specs/129-prelaunch-release/handoff.md using existing workflows and contracts/release.md, including exact future source/attempt/asset inspection and separate approvals.
- [X] T009 [US2] Record current source-only delivery, standing-waiver untested scope and later pending outcomes in specs/129-prelaunch-release/delivery.md and verification.md.

## Phase 5: US3 - Truthful capabilities and deferred outcomes (P2)

**Independent test**: Structural and content-bound comparison checks pass for v2.0.0, public selection stays v1.7.0, and unchanged historical files/215 inputs have zero drift.

- [X] T010 [US3] Add site/data/comparisons/v2.0.0.json with reviewed official-source qualifiers, current candidate links and explicit owner-deferred clustered-execution opportunity; preserve every historical snapshot.
- [X] T011 [US3] Generate site/data/reviews/v2.0.0.json with exact formatted dataset/section/notes digests and ten explicit decisions; distinguish prepared/public behavior in docs/website.md.
- [X] T012 [US3] Reconcile #412 with selected candidate and source/public disposition using formatted GitHub prose, keeping #355/#356/#400/#404 and cluster records unresolved/deferred.

## Phase 6: Source identity correction and verification

- [X] T013 [US1] Reproduce and correct #413 in scripts/release-source-check/main.go/main_test.go and README.md: require the full injected tag, retain every existing negative fixture and preserve red/green receipts in .test-temp/s129.

- [X] T014 Run actual github-format, UTF-8/no-BOM/mojibake and staged whitespace checks across changed source; finalize candidate review hashes after relevant formatting in site/data/reviews/v2.0.0.json.
- [X] T015 Run all eight foreground scripts/verify.sh all gates to terminal completion with quiescent source and zero protected-input drift; retain first failures in .test-temp/s129.
- [X] T016 Run clean desktop/site restores, component/export and existing headless browser cases, all test-script ShellCheck and PowerShell compliance; redirect new legacy-named site captures into .test-temp/s129/browser-captures.
- [X] T017 Recheck all current lifecycle/integrity guards and candidate source/comparison identity; record source-versus-artifact evidence and limits in specs/129-prelaunch-release/verification.md.

## Source publication workflow evidence

Publication bookkeeping is evidence rather than open implementation work, following specs/README.md.

- W01: Commit with project convention, automatically push codex/129-prelaunch-release and open one official PR targeting main with Refs #412; retain and attach the PR receipt.
- W02: Obtain terminal successful exact-head CI and actual review receipts, address every finding and resolve threads within at most two Codex requests total; retain requests/results in .test-temp/s129/publication-ledger.json.
- W03: Publish final exact-head source CI/review evidence in the PR and concise handoff comment, then notify the owner for final review/merge with later release boundaries explicit in specs/129-prelaunch-release/handoff.md.

## Phase 7: Same-slice later owner-controlled delivery

- [X] T018 [US2] After owner integration, sync clean main, verify squash tree and fresh exact-main CI/CodeQL/website, and bind source/comparison checks to that actual SHA in specs/129-prelaunch-release/delivery.md.
- [X] T019 [US2] Obtain exact v2.0.0 tag/draft approval for accepted main, inspect absence/conflicts, then create/push the immutable tag to existing Release workflow; record actual run/attempt in delivery.md.
- [X] T020 [US2] Inspect all eight actual staged assets, GitHub digests, archive versions, MSI/manifest identities and all seven job conclusions, retaining the actual frozen inventory in .test-temp/s131 as the S131 execution continuation.
- [X] T021 [US2] Obtain separate concrete-asset promotion approval, then dispatch existing standing-waiver promotion for the same candidate; record actual nine-file result without rebuilding accepted bytes in delivery.md.
- [X] T022 [US3] Verify public/latest, original asset IDs/bytes, checksums, Pages/HTTPS/discovery, implement current wording/default reconciliation in S131, retaining #412 closure on final owner source integration as separate workflow evidence; preserve deferred outcomes.

## Dependencies and execution order

Setup -> analyze -> US1 source identity. US2 handoff can be written independently after analyze. US3 snapshot requires US1's selected metadata, and its review digests follow final formatting. Complete all source changes before the aggregate verification snapshot. Supplemental verification and final factual records precede commit/push. CI/review precede owner handoff. T018-T022 require the stated later owner actions and approvals, in order.

## Parallel opportunities

US1 module identity and US2 handoff use separate files after analyze. US3 official-source review can run read-only during metadata preparation. The installed plan required three research agents; implementation will be integrated serially to keep the broad mechanical namespace change and verification source stable. Do not run mutations concurrently with tests or hash snapshots.

## Implementation strategy

Deliver one coherent source increment, validate it completely and publish one official PR under existing authorization. Preserve the same S129 path to actual release after owner approvals. No additional preparation slice, new release framework, historical conversion or unattended publication is introduced.
