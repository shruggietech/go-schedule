# Tasks: S135 v2.0.1 publication closeout

## Phase 1: Setup

- [x] T001 Record accepted source, exact-main CI/CodeQL/Pages and public release baseline in .test-temp/s135/baseline.json.
- [x] T002 Complete requirements, research and design in specs/135-v201-publication-closeout/ and update .specify/feature.json and CLAUDE.md.

## Phase 2: Foundation

- [x] T003 Run the blocking cross-artifact gate and retain its result in specs/135-v201-publication-closeout/analysis.md before implementation.

## Phase 3: US1, obtain the accepted maintenance release

**Independent outcome**: Complete public v2.0.1 packages identify accepted source and retain the native-testing disclosure.

- [x] T004 [US1] After explicit exact-source approval, create the immutable accepted-main tag and stage through .github/workflows/release.yml, retaining actual run/attempt in .test-temp/s135/.
- [x] T005 [US1] Download and inspect the selected draft assets using existing scripts/windows-release-gate, test/windows/inspect-installer.ps1 and fresh .test-temp/s135/ identities; preserve failed observations.
- [x] T006 [US1] After separate actual-asset approval, publish the selected bytes through .github/workflows/promote-release-waived.yml and record the actual release identity in specs/135-v201-publication-closeout/delivery.md.
- [x] T007 [US1] Independently read back public/latest, all nine downloaded assets, eight checksum entries and Pages/HTTPS in .test-temp/s135/, recording factual results in specs/135-v201-publication-closeout/verification.md.

## Phase 4: US2, discover the release consistently

**Independent outcome**: Current source discovery and availability guidance identify the actual public release; deployed guide status is separately reported.

- [x] T008 [US2] After publication, reconcile README.md, docs/first-task.md, docs/gui-fields.md, docs/install.md, docs/INSTALL-windows.md, docs/website.md and test/windows/README.md, preserving historical baseline claims.
- [x] T009 [US2] Select the existing v2.0.1 comparison in site/data/published.json, preserving the original live-selector baseline and historical comparison/review records under .test-temp/s135/.
- [x] T010 [US2] Update the existing availability assertion in site/e2e/site.spec.ts and exercise production documentation/comparison/browser contracts with fresh .test-temp/s135/ outputs.

## Phase 5: US3, review a bounded closeout

**Independent outcome**: Durable delivery and issue traceability distinguish public release, source integration and deployed discovery.

- [x] T011 [US3] Record actual package/source outcomes, limitations, #419 traceability and the two-request review ledger in specs/135-v201-publication-closeout/delivery.md and verification.md, keeping PR workflow bookkeeping outside feature tasks.

## Phase 6: Polish and validation

- [x] T012 Run all eight canonical gates from scripts/verify.sh all in the foreground through fresh hidden S135 wrappers, plus exact-source release/comparison and preservation checks, recording results in specs/135-v201-publication-closeout/verification.md.
- [x] T013 Update CHANGELOG.md Unreleased with discovery and the dated pinned-guide decision; finalize specs/135-v201-publication-closeout/spec.md and specs/README.md lifecycle evidence and run GitHub format/encoding/lifecycle checks.

## Dependencies and execution strategy

T001-T003 precede implementation. US1 staging, actual-asset selection and promotion are sequential and require their stated decisions. US2 source changes follow actual publication; discovery can be reviewed independently once the public release exists. US3 records the outcomes of both stories. T012 follows the final source edits, and T013 finalizes factual records without changing qualified behavior. Official PR push/review/CI monitoring and final owner merge are workflow evidence, not unresolved feature tasks.

Independent read-only source/package metadata observations may overlap within US1 after a concrete selection. Once publication exists, independent guide edits in T008 and the fallback edit in T009 may be prepared in parallel, then T010 validates them together. US3 evidence summaries can be assembled alongside finished story observations. No parallel worker may write protected historical outputs.

The first useful increment is US1's actual public package delivery; US2 completes the parent discovery outcome. No preparation-only substitute completes #419.
