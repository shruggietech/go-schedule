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

## Remaining current-candidate evaluation

The next operation builds six exact binaries from one clean frozen source checkpoint and performs a fresh full owned Windows/WSL capture with actual semantic reviews, then finalize/validate/render. This work is still pending and is not replaced by the canonical gate pass or any historical full-catalog report. Full private logs and reproduction fixtures remain under ignored .test-temp/s115. Hosted CI and Codex receipts belong to their actual published heads.
