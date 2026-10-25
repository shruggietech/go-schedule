# S115 observations

## Planning and carried source

Spec-Kit specify, clarify, requirement checklists, plan, tasks and blocking read-only analysis completed before implementation. Both checklists passed (16 and 12 items); analysis found zero critical/high issues and full requirement coverage across 15 tasks. The optional context hook was inspected and its plan-pointer effect applied directly without unwrapped Python launches. Carry-in 18e6fd8 is incorporated as dde3e7b; original 574d071 remains archived.

## Release-source defect

The unchanged production inline preflight from main 7bce610 accepted a local isolated heading ending in not-a-date when its notes anchor matched, native exit 0. This demonstrated #366 without changing any release. The shared command regression matrix then failed before its implementation because the command was absent, and passed under native Windows race after implementation. Tests cover actual candidate metadata, valid/whitespace/CRLF fixtures and missing, duplicate, malformed or mismatched source inputs. The actual v1.6.0 metadata independently passes the command.

## Native Windows follow-up

One original six-writer diagnostic run passed for both PowerShell and Bash and both fresh/existing schemas, retaining exactly six distinct rows. PowerShell reported 7.6.6, .NET 10.0.12, Windows 10.0.26300/X64; Bash reported 5.3.15 and MSYS 3.6.10. This is a current observation, not a correction of original #355 causes.

The initial memory/file-backed IPC experiment passed under native Windows race. A subsequent source refinement releases the two snapshot requests at create-request receipt so reads and writes proceed together, preserving cancellation and original command/failed-observation attribution. Final-source results are recorded after completion.

## Initial canonical failure and preparation correction

verify-all-01 failed the unchanged ten-minute integration-package timeout. The active TestV13PackageDefaultsRemainOptIn had run for 52 seconds and was waiting on its daemon build, not a CLI invocation. The raw stack remains retained. No whole-suite or original #355/#356 causal correction is inferred from a later pass.

Independent #367 now reuses each identical CLI/daemon build once within a test process, copies immutable bytes to the original independent fixture paths and uses existing hidden process-tree containment with a two-minute preparation bound. Existing runtime assertions, global timeout and QA limits remain unchanged. Concurrent reuse and sticky-failure regressions failed before the helper existed, then passed under race with the original packaged CLI/MCP/v1.3/v1.4 scenarios (27.356 seconds). This focused timing is not a measured whole-suite speedup.

The original two-mutation-worker IPC case is retained separately from the new simultaneous complete-snapshot workload. Final canonical acceptance covers both memory and file-backed cases and preserves every original request/cancellation/mutation assertion.

## Second aggregate and timeout-fixture correction

verify-all-02 failed only the injected QA-wrapper capture timeout fixture, which assumed both child outputs had been emitted before its 150 ms deadline. The integration package completed in 473.051 seconds, within the unchanged global limit. This remains a failed aggregate, separate from the initial timeout and the original native causes.

Independent #368 establishes child-output readiness at the existing wait boundary in a private test copy before the same 150 ms policy is enforced. All original stdout/stderr, recovery and budget/resume assertions remain, with an added exact 150 ms assertion. Production scripts/cli-qa.ps1 is unchanged. Native race checks pass both capture and resume plus the preparation-cache regressions (9.923 seconds). New package builds explicitly disable Git/GCM prompts while inheriting existing build settings.

Final native Linux IPC coverage passes both memory and file-backed cases, with the original two mutation workers retained alongside the new concurrent complete-snapshot workload. The corresponding final-source Windows evaluator package passed under race in verify-all-02. #355/#356 remain unknown-cause observations.

## Accepted canonical aggregate

verify-all-03 returned native exit 0 after all source corrections. All eight gates passed: format, vet, lint (zero issues), root race, GUI, coverage, docs and automation. Root integration completed in 424.985 seconds under the unchanged global timeout. Native Wails build, desktop tests and all 167 frontend tests passed. Core coverage is engine 83.3%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%. Docs exported 22 canonical pages, and both positive automation policy and mutation fixtures completed successfully.

Application, GUI, coverage and docs execution was native Windows. Only portable automation policy and its unchanged driver used the recorded private SH bridge to Ubuntu-24.04/WSL, matching the prior accepted workstation route. No bridge is committed. All sixteen recorded source/test hashes stayed fixed through this accepted aggregate. Temporary root coverage and mutation-fixture directories were removed by their normal cleanup. Earlier failed aggregates remain failed, independently attributed receipts.

## Runtime identity and upstream disposition

Fresh native SQLite inventory reports 3.47.2 (64-bit). Current interpreter/runtime SHA256 identities are retained with the private native log: pwsh.exe bfb46af89433268872ddb43d1ca7a3f433452ee91ed356a9786940f90118e285; System.Management.Automation.dll ed8d39663e05435aab59f1c7fb0ad20292e22b0cb12cf47aad4b62b31e4d4126; coreclr.dll 20ea763d5ad8be8656a3123c4ad5afdb4c714a009bfc63bb8f4906d6678130db; bash.exe 9eedc85f00d4b090f50ce4bede70802619212baa90c15e57893413bca0de9fad; sqlite3.exe 5e33fe030fb7c3bbe2bca1f70f21a406716961aefdfb1bc030d7c65b7db055e9.

[Official PowerShell documentation](https://github.com/MicrosoftDocs/PowerShell-Docs/blob/main/reference/docs-conceptual/whats-new/What-s-New-in-PowerShell-76.md) identifies 7.6.6 with .NET 10.0.12, matching the observed inventory. Reviewed upstream reports describe different .NET 7 analyzer/build and older Windows CET-startup contexts; they do not establish this compiler/parameter-metadata crash cause. The [7.6.7 release notes](https://github.com/PowerShell/PowerShell/releases/tag/v7.6.7) likewise do not establish a matching correction. No installed runtime was replaced, and #355/#356 remain open with their original causes unconfirmed.

## Completed frozen candidate evaluation

Fresh owned run ada3d85c-aee8-42c8-bdfc-fdbf5f846a76 captured clean source fc39261ed8a874aaad11dd83f7afcccaf06fb10d with content SHA256 74da617216a745dc038c21a727161fe5f4b4f3c825b84283d096b85d562f4208 from 2026-10-24T20:26:37.9253392Z through 2026-10-24T21:17:26.7692275Z. The selected Windows 11 and existing Ubuntu-24.04/nonroot-user environments used the complete two-hour budget, original 15-second case deadline and unchanged resource/ownership controls. No catalog retry was used to replace a failure.

Execution completed all 1,763 applicable cases: 1,760 passed and three Windows metadata cases failed. There were three explicit platform exclusions. All 1,109 distinct messages received exact command/scenario/channel/occurrence-bound seven-part reviews: 1,094 passed and 15 failed; 44 no-message occurrences remained explicitly silent. The final effective result is fail, with 1,744 passes, 19 failures and zero blocked cases across three findings. All four owned process trees were drained, runtime directories removed and cleanup remaining lists empty. Finalize, validate and render each returned native outcome code 2, identifying a valid failed evaluation.

| Windows case | Original exit | Elapsed ms | Original stdout |
| --- | --- | --- | --- |
| help.audit | 0 | 17009 | complete audit help |
| help.profile.default.reset | 1 | 15376 | empty |
| help.profile.default.set | 0 | 18617 | complete profile default set help |

All three state diagnostics say context deadline exceeded, command_ran true, metadata_only true and successful case-process cleanup. Their guard/setup manifest and runtime-info requests completed. Metadata-only cases intentionally do not collect full before/after snapshots, so their false before_complete/after_complete flags do not establish a failed full-state observation. These are new attributable #356 observations; a common cause remains unconfirmed and no product missing-help defect is inferred from the timing failure.

| Receipt binary | SHA256 |
| --- | --- |
| windows cli | f497aa5e3f7bfd92f44d88f991d24281abcb05db40305f7e5c8453bbd4892d67 |
| windows daemon | 76791b19949ba7284f634539fd6327d5817464365640d1067b0afe635e6ab8b8 |
| windows runner | a73fb445e16a4a128901a4642b5576037e4f73a580cbf5830aab1aa7dc2e8fb4 |
| linux cli | 4a3094643619fa42a0df95e157f045fa73a21284f56227c74e4fb9856d2428a0 |
| linux daemon | f61d2d4b8ea0c6c8040b802826ae706337a402eb0e35f9805f0d85f0cdbd3ecf |
| linux runner | 6e47ff87b163d41eb5420f5e7b49f8c21ad48b6e4d7f5135f9a22885d7a9a92d |

Private complete inputs, report and evidence remain under .test-temp/s115/qa-inputs-01 and .test-temp/s115/qa-run-01; the review submission is .test-temp/s115/reviews-01.json. The generated report retains the fixed seven-section human outline and full defect schema. Report JSON SHA256 e378151781ac1a4433eb790e68108d455a92e17e50559035946c99d17896169b; report Markdown SHA256 106dde5e3bfcec68c3e3048b597c0f4997912ea21c25853369959e6ff27e73df; corpus SHA256 e8abd02d91b9a1ab8dc4584ae014d61e9ca68182f99e63f80966b60fd1c6ed49. No generated frozen artifact is hand-edited or relabelled after source changes.

## Independent communication corrections

The actual message review demonstrated #370: default PowerShell help told readers to persist output from a pipeline that evaluates Register-ArgumentCompleter. The scoped correction tells readers to add the invocation itself as a line in $PROFILE, preserves current-session loading, default/no-description generated bytes and other shell help, and edits no owner profile.

The same review demonstrated #371: authorized Manage/Operate MCP sessions inherited whole-session Observe-only/read-only identity and instructions despite advertising and accepting authorized mutations. The scoped correction uses permission-neutral server identity and explicitly authorized-tool wording, preserving resource-specific Observe descriptions, the untrusted-content instruction and all permission/confirmation/retry behavior. Real SDK discovery regressions retain exact tool counts.

Incremental blocking Spec-Kit analysis maps both discoveries to FR-005 and T018 with zero critical/high findings. Regressions reproduced both problems before correction. The first test-baseline experiment also lacked the dummy child Cobra requires to initialize default completions; that independent setup was corrected before the clean pre-fix regression receipt. The original failed full evaluation remains bound to fc39261; later source and focused/canonical evidence are separate identities, not a replacement full catalog.

The initial post-correction native race run passed CLI, Observe/Operate/Manage, localhost MCP and remote MCP packages. Its evaluator inventory check rejected a manual JSON encoder formatting change, so that run remains failed. Preserve Go canonical HTML escaping in the inventory; the intended semantic change is only the PowerShell Long field. Final corrected-inventory acceptance is recorded separately.

## Final communication-correction acceptance

Corrected-inventory native race checks returned exit 0 in 3.675 seconds for TestCompleteManifestAndDriftAreSeparateFromHelp and TestEveryCatalogFlagExistsInActualTree. The inventory diff changes only PowerShell Long guidance and preserves Go canonical escaping.

verify-all-04 returned native exit 0 with all eight canonical gates passing on the final communication-correction source. All 22 recorded source/test hashes remained fixed through completion. CLI and evaluator race packages ran freshly in 130.372 and 114.232 seconds; the root race integration result was cached, so it is not a fresh packaged integration execution. The separate uncached coverage gate passed at engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%. Native Wails build, 167 frontend tests, the desktop production build, 22-page docs export, positive automation policy and negative mutation fixtures all passed.

Application, GUI, coverage and docs execution used native Windows. Portable automation policy and its unchanged fixture driver used the private Ubuntu-24.04/WSL SH bridge; it is not committed. Normal cleanup removed the temporary coverage profile and mutation directories. Canonical log SHA256 ff9d6f00e9aac75bda728f4d9d31118778adbd434b82ec2f46e5830a6f4c88a8. This aggregate does not rerun or replace the frozen full catalog, whose failed result and three findings remain recorded above.
