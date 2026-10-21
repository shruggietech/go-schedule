# S109 publication and bounded review record

## 2026-10-21: Authorized publication

Owner authorization covers pushing S109, opening PR #332, publishing additional evaluator defect #331, monitoring CI and completing up to two Codex review rounds, resolving findings autonomously and stopping before merge. The initial local implementation is 34ab5cb544ac6bb2671dde30238163e213120d4d; first published head is 01e6d8efb0b0a3da42cf77712d58ce9d77e9eac7. No authentication was needed.

PR #332 covers the full #311-#322 bundle and #331, all under parent #277 and milestone 15. The additional issue uses the existing taxonomy, functional acceptance criteria, explicit dependencies and source-pinned reproduction. Open and closed issue searches found no duplicate. All issues remain open pending merge.

## Review and hosted verification

The PR was created in draft while publication metadata was reconciled, then marked ready on a fixed head for round 1. Hosted CI and Codex review are separate evidence. Exactly two rounds completed; no third round, merge or auto-merge is authorized.

## 2026-10-21: Codex round 1

Marking the draft ready automatically started round 1 at 2026-10-21T13:20:16Z on 0dcc0de9e4d95bf593c2dfce19d58239863da8cd. The bot summary is https://github.com/shruggietech/go-schedule/pull/332#issuecomment-6317104578. No duplicate manual request was sent. CI run 39607085730 and CodeQL run 39607085887 target the same revision. That initial snapshot recorded review and checks in progress; their terminal outcomes follow.

Round 1 completed at 2026-10-21T13:24:48Z with two findings. The documentation contradiction is confirmed and tracked as #334. The version finding is not reproducible on the reviewed implementation: a new exact-output regression passed before any execution-code change for --version, -v, --version=true and --profile=missing --version. Parsed pflag values persist when Cobra reparses the empty argument list, so no parser workaround is introduced. The regression remains as explicit coverage.

CI independently found a new macOS targeting-test socket path above the platform limit (#333). The fixture now creates its macOS socket in an exclusively owned short private directory with test cleanup; actual IPC/profile/TLS assertions remain. This is distinct from #308's delivered clifixture workspace correction. The documentation now consistently distinguishes invalid-input exit 2, unsupported well-formed lines and operational exit 1.

Review-fix validation: the complete Windows CLI race package passed (16.177s), native Linux Go 1.26.4 CLI race package passed (8.104s), and the full docs gate passed including site tests, 22 canonical documents and static export. An initial native-test launcher attempt interpreted the test-selection separator through a shell and failed before it could establish a valid result; direct WSL execution replaced that launcher and produced the recorded successful run. No product execution code changed in this correction set. The initial head completed 19 successful hosted checks and three macOS jobs failed at the same overlong CLI socket fixture; the fixed head receives fresh hosted checks.

## 2026-10-21: Codex round 2 and final handoff

[The round-2 request](https://github.com/shruggietech/go-schedule/pull/332#issuecomment-6317446854) targeted 4db6501b4ec3d7256609525288e52a4b37287e21. [The bot summary](https://github.com/shruggietech/go-schedule/pull/332#issuecomment-6317104578) records completion at 2026-10-21T13:39:22.796896Z with has_any_findings=false and has_unresolved_findings=false. Both round-1 inline threads are resolved after the confirmed docs correction and evidence-based version disposition. The original receipt-bound Windows root.version artifact also contains exactly gosched version 0.0.1-dev followed by a newline, corroborating the original successful behavior.

Exactly two Codex rounds were used: automatic draft-ready round 1 on 0dcc0de and manual round 2 on 4db6501. No third round was requested. The final publication-record commit changes documentation only after that clean review. Current-head hosted checks are monitored to terminal status independently, and their exact final SHA/run links and readiness are recorded in the PR handoff comment rather than pretending the earlier runtime capture or review covered a different source identity.

At the delivery-record update, 17 of 22 checks on 4db6501 had passed, including the previously failing native macOS remote qualification and CodeQL; five remained running and none had failed. This is an interim historical snapshot, not the final CI conclusion. The final metadata revision receives its own complete check run before handoff. PR #332 remains open and unmerged with auto-merge disabled; #311-#322, #331, #333 and #334 remain open until their functional corrections merge.
