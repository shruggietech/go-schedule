# S114 verification record

## Identity and boundaries

Base/main is 7f7c5ec2c5252bb4b116b8d3b365438f30196779; branch codex/114-release-readiness. Specification, clarification, requirements checklists, planning, tasks and blocking analyze preceded implementation. Private observations currently concern this base plus the identified dirty S114 tree, not a committed release or historical capture. No merge, tag, staging workflow or promotion has been performed.

Original S112 log SHA256 remains 4eded1dd882f92b7f09aad62265065c46631212d79d97ce7b12922a4b2d3cacf. The original CLR access violation and separate prerequisite refusal remain unconfirmed. S113's full catalog stays bound to ffebda5c3230ecfeb6eab5c542aca1a12eeee621, not this slice.

## Demonstrated regressions

- release-notes-before.log: native exit 1. Actual prepared v1.6.0 notes fail production read-only resolution because untested disclosure is absent. After correcting notes, all release-promotion race tests pass (release-notes-after.log), preserving no-mutation assertions and existing identity/checksum refusal cases. New independent bug: #363.
- launcher-before.log: native exit 1. The controlled wrapper exits 73 with its own explanation, but the helper fails again on an absent capture file before retaining that explanation. launcher-after.log passes with original status and explanation. The owner explicitly approved correcting this default failure output; successful output and writer assertions remain unchanged. New independent bug: #364.
- pipe-before.log: native exit 1. The real hidden candidate exits after spawning an owned child retaining redirected handles. The previous helper remains blocked until the separate 20-second test guard, demonstrating that process-only timeout does not bound redirected reads. With a shared process/stream deadline, probes-after.log passes the inherited-handle regression and every existing candidate case for both twins. The inherited-handle subtest includes fixture compilation; the actual probe-return assertion remains inside the unchanged ten-second policy plus startup allowance. Its child receipt is cleaned by the test. This establishes the current helper defect only, not the original CLR cause.

The PowerShell correction uses the [asynchronous StreamReader contract](https://learn.microsoft.com/en-us/dotnet/api/system.io.streamreader.readasync). Existing timeout wording, candidate output limits, explicit-selection refusal and implicit fallback remain intact. Expanded writer phases and all-six failure evidence are opt-in through GOSCHEDULE_PROBE_DIAGNOSTICS=1; default output receives direct coverage.

## Native investigation

writers-windows.log passes two repetitions of each native six-writer twin, with both fresh and existing schemas and exact identity/count assertions. It records PowerShell 7.6.6, .NET 10.0.12, Windows 10.0.26300/X64 and Git Bash 5.3.15/MSYS 3.6.10. These passing observations do not explain the original failure.

ipc-current-windows.log and ipc-prior-windows.log pass real owned Windows IPC/store requests: 32 complete snapshots, cancellation receipt and retained earlier collections, a successful fresh snapshot, sixteen exact group creation/deletion pairs without retries, and the existing incomplete-command/inspection distinction. The prior variant is explicitly a diagnostic-only current-source dependency overlay selecting go-winio 0.6.2 and x/sys 0.48.0; all other source is current. It is not a full prior-release qualification.

native-linux-current-02.log passes real native Linux IPC and POSIX six-writer/candidate cases on Ubuntu-24.04/WSL2. PowerShell is absent there, so its twin and inherited-pipe case are explicitly skipped; native Windows exercised them. Dependency input uses the existing local file proxy and retained go.sum verification, with GOSUMDB disabled only for this cache-backed diagnostic command. The first native-linux-current.log attempt failed parsing launcher JSON before WSL or tests began; it is retained as an invocation error, not an application failure.

## Remaining outcomes

The original #355/#356 causes remain causally unconfirmed, so both issues remain open. #363/#364 track the demonstrated source corrections independently. #354 remains open for the later public release. This slice changes maintainer probes/fixtures and release guidance, not application/evaluator executable behavior; no new full-catalog capture is substituted for S113's immutable report.

Canonical verification, final exact source identity and hosted CI/review receipts are recorded below when completed.

## Output compatibility correction

The first S114 stream-deadline attempt discarded stderr already emitted before cancellation. partial-stderr-before.log records that intermediate failing correction, not a baseline product failure. The final bounded incremental read loop retains completed chunks; probe-contract-after.log passes original timeout stderr, inherited-handle timeout, approved launcher failure and explicit default/opt-in writer contracts. Existing timeout wording, candidate limits and writer assertions remain intact. No extra default output or file-size cap was added.

ipc-prior-linux.log passes the same real owned IPC cancellation/mixed-mutation regression on native Linux with the attributable prior dependency overlay. Both current/prior Windows and Linux observations pass; historical causation is still unconfirmed.

## Native Linux timeout-stderr discovery

Final native-linux-final.log retained native outcome 1: default/opt-in writer behavior and other POSIX candidate cases passed, but emitted stderr was absent after timeout. The Coreutils group terminates stdio-buffered head before its small pending output reaches disk. This is a demonstrated current helper defect under #364, not an inferred historical cause. The initial canonical run was stopped after that independent failure became known and is not accepted verification. The correction uses byte-sized bounded dd writes with the same 4097-byte detection cap, preserving original emitted chunks without changing candidate or case deadlines.

The second canonical attempt passed format/vet/lint and entered race checks, then was stopped to make the launcher regression explicitly select default diagnostics=0 after the owner exception. It is not accepted aggregate verification. The final aggregate below must cover that strengthened default-path regression.

## Canonical attempt 03 and release-copy policy correction

Attempt 03 returned native exit 1 at the final automation gate. Format/vet/lint, root race (integration 456.845 seconds), native Wails build, all 167 frontend tests, six core coverage thresholds and the 22-page docs export passed. Automation correctly rejected a separate native-observation section because release notes require exactly one Highlights heading, four to six consecutive bullets and one changelog link. No policy or assertion is relaxed. The required waiver disclosure now lives in the fifth highlight, with a sixth concise known-issue bullet; content receipts are refreshed. The corrected gate, actual-note resolver and new full aggregate must pass before publication.

## Portable policy execution environment

The corrected actual-note resolver passes. Native Windows positive automation policy also passes, but its unchanged negative-fixture batch remained pending for more than ten minutes and was stopped without a success claim. A private windowless shell bridge delegates only scripts/automation-check.sh and its unchanged negative-fixture harness to native Ubuntu-24.04/WSL, matching hosted POSIX policy execution. All application race, GUI, coverage and documentation gates remain native Windows. The repository driver, manifest, policy rules, assertions and fixtures are unchanged. The subsequent complete aggregate explicitly uses this recorded SH override, rather than claiming the unfinished Windows negative batch passed. No bridge is committed or shipped.

## Accepted local aggregate

verify-all-04.log returned native exit 0 after all source corrections: format, vet, lint, root race, desktop race/native build and 167 frontend tests, coverage, documentation and automation. Core statement coverage is engine 83.3%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%. The documentation gate exported 22 canonical pages. Both the positive policy and negative automation fixtures returned success.

Application, GUI, coverage and docs execution was native Windows. Only the unchanged portable automation policy and negative fixtures used the explicit private SH bridge to Ubuntu-24.04/WSL. The driver and all eight gates, policy rules and fixture assertions remain unchanged. The earlier unfinished native Windows negative batch is not accepted as a pass. The private bridge is not part of the source deliverable.

The recorded eleven input hashes were unchanged throughout the accepted aggregate. Native-linux-stderr-after.log independently passes the final POSIX timeout-stderr and default/opt-in contracts; release-notes-final.log passes the final six-highlight notes through production read-only selection. No installed runtime, public selector, historical comparison, release tag or release asset was modified. Final lifecycle and publication receipts are documented separately from the tested executable/test inputs.

## Initial hosted review and metadata correction

Codex round one completed on 15e2578554a5ebc1b4b9b0ea7eec2f2113620575 with no findings and no review threads. Hosted CI run 39985108430 found a lifecycle error introduced after the accepted local aggregate: the specification and inventory were marked Implemented while two delivery tasks were still open, and the delivery paragraph lacked the required implementation reference. Lint itself reported zero issues. The correction keeps status In Progress while CI/reviews remain pending, records PR #365 and the actual source commit, and checks off only the now-completed publication task. The lifecycle gate and its assertions remain unchanged.

## Accepted hosted checkpoint

All 22 applicable PR checks are terminal COMPLETED/SUCCESS on d11646111be52851c5b32a468050b508c63954c4. CI run [39985842927](https://github.com/shruggietech/go-schedule/actions/runs/39985842927) and CodeQL run [39985842835](https://github.com/shruggietech/go-schedule/actions/runs/39985842835) provide that source-bound receipt; Windows race completed successfully after 10 minutes 8 seconds. The initial automatic Codex review completed with no findings on 15e2578554a5ebc1b4b9b0ea7eec2f2113620575 ([review summary](https://github.com/shruggietech/go-schedule/pull/365#issuecomment-6369390824)); there are zero review threads. Application/test source is unchanged between these two commits.

The final lifecycle/evidence checkpoint updates only spec.md, tasks.md, verification.md and specs/README.md. Its own CI and the optional second/final Codex review are monitored separately on the actual final head, with terminal receipts retained in PR #365. This paragraph does not claim those future results. The PR remains open, unmerged and without auto-merge. Its closing references are only #363 and #364; #354, #355 and #356 remain open.
