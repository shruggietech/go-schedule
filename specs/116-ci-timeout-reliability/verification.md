# S116 verification record

Date: 2026-10-25. Baseline: 9e744e434f07a8d22041e0d4d6318867410900b6. These are source-development receipts, not a published release or installed-product qualification.

## Chronological evidence

1. Spec-Kit setup, requirements checklists and blocking cross-artifact analysis completed before implementation: 10 requirements, 17 tasks, complete coverage and no findings. The initial checklist prerequisite correctly failed without a plan; setup-plan supplied it and the retry passed.
2. The controlled six-second frontend probe failed under the original 5000ms watchdog (5016ms). After the shared runner increase it completed in 6009ms; the actual App and Search suites passed alongside it (21 tests, including that one temporary probe). The probe was deleted. Production reconnect and exact daemon identity behavior remain unchanged.
3. The new timeout guard rejected the old workflow/nested policies before implementation and passed after implementation. A later Makefile audit added 60-minute test/benchmark parity and mutation coverage for Make variable commands. Budget expectations rejected the old defaults/reserves before implementation and passed afterward.
4. Focused race packages passed: clievaluate 84.998s, ci-timeout-check 1.106s, release-promotion 2.112s and windows-qualification-session 8.630s. The guest tests use nondestructive fixtures and host refusal; they do not install a guest or service.
5. The first metadata fixture run failed because its expected request label was “complete” instead of the actual “completed” and it assumed six metadata definitions instead of the actual three. The fixture was corrected to the existing manifest contract. Both the failed receipt and successful retry are retained.
6. Windows six-writer follow-up passed for PowerShell and Bash, each with fresh and existing schemas. All six literal session identities, rows and exits remain asserted; test 5.13s, package 6.185s.
7. Linux follow-up passed the POSIX six-writer checks and memory/WAL IPC workloads. PowerShell was unavailable on that Linux environment and its twin explicitly skipped. The IPC workload retains 32 snapshots, 320 requests, both mutation workers, 16 exact create/delete pairs, all 13 state keys, cancellation, partial-result and recovery assertions. IPC 6.13s. The actual CLI exercised help.audit, help.profile.default.reset and help.profile.default.set with all four ownership requests, metadata_only, command_ran, exit/output and process_cleanup assertions. Controlled completion beyond the former 15-second limit took 16.175162966s and passed under the selected 120-second default. An explicit 100ms budget still failed at the first manifest guard without running the command. Package 59.499s.
8. Windows memory/WAL IPC and the same three actual metadata cases passed. IPC 5.76s; controlled beyond-old-limit completion 16.1441179s; explicit 100ms failure retained; package 26.297s. Owned server/store/temp cleanup is asserted by the fixture.
9. Source release and comparison guards passed with pending tag v1.6.0, baseline source 9e744e4 and refreshed governed content digests. The published selector remains v1.5.1. No tag, package, staging or publication identity is inferred from these checks.

## Runtime and source attribution

Windows: Windows 10.0.26300 x64, Go 1.26.0, PowerShell 7.6.6 on .NET 10.0.12, Bash 5.3.15(2)-release on MSYS 3.6.10, SQLite 3.47.2 64-bit. No runtime was installed or replaced. Linux: existing nonroot Ubuntu-24.04 WSL2, kernel 6.18.33.2-microsoft-standard-WSL2, Go 1.26.4, Bash 5.2.21 and SQLite 3.45.1.

Both native metadata runs used CPU 25%, memory 1GiB, Windows processes 32, Linux tasks 128, case 120000ms, total 7200000ms, channels 262144 bytes and artifacts 67108864 bytes. Explicit and resumed budgets remain authoritative. These targeted tests are not a new full catalog/coordinator/cgroup qualification.

| Receipt | Source and binary identity | Log SHA-256 |
| --- | --- | --- |
| Windows IPC/metadata | Baseline 9e744e4, dirty source SHA-256 bcafe90859130829d33a289968f88ccfa0919101d4be2755643e6bdbffc89963; CLI 2caba57515f04dbbce79860c297aebab9d631e8a95b679e6e54a588bbcf65f58 | 75a6dcdc915c0d9f8cef282dad111ad8af50284ebe58c9168aab9893b16e6fc9 |
| Linux writer/IPC/metadata | Baseline 9e744e4, dirty source SHA-256 3755727fba98ac96e638f1ef19d1ff73f784dbab6410f7793bd55edf858e67d0; CLI 1096312ea17b4dce0d244d5bcfe35a22e067bf4e22d1c742e26e9628eea44345 | 483f70dbfc91901f68ed42abcf57b86bd9cddd6e1207d572df23e9872dd388bb |
| Windows six writers | PowerShell and Bash, both schema states, six exact identities | 9ba88c8713494f15f9fb5528febafad948ba2d7e8108a0016e49d0b8a04d8a62 |

Platform source/binary hashes are separately attributed and do not identify the eventual PR commit. Native logs, exit receipts, runtime executable hashes and the evidence summary are retained privately under .test-temp/s116/.

## Preserved earlier outcomes

S115 exact merged-main CI 40177500270 failed on macOS at the 5000ms frontend watchdog; its log SHA-256 is b71b82a8e98ac0cc7c0d5e866574c3958506d313d9371ca0173995f8f8a8d80f. CodeQL 40177500250 and website 40177500252 passed on that same head. S116 does not relabel that failed run.

The S115 full catalog report remains unchanged: 1760 execution passes, three Windows deadline failures, 1744 effective passes and 19 effective failures, with findings #356, #370 and #371. Report JSON SHA-256 e378151781ac1a4433eb790e68108d455a92e17e50559035946c99d17896169b, Markdown 106dde5e3bfcec68c3e3048b597c0f4997912ea21c25853369959e6ff27e73df, corpus e8abd02d91b9a1ab8dc4584ae014d61e9ca68182f99e63f80966b60fd1c6ed49. #370/#371 were fixed separately in S115. Passing S116 targeted runs and longer allowances do not establish the original #355 CLR/prerequisite or #356 IPC causes; both remain open.

## Canonical suite

scripts/verify.sh all completed successfully in canonical-all-01 with native exit 0 on 2026-10-25. All eight gates passed: format, vet, lint (zero issues), race, GUI, coverage, docs and automation including negative fixtures. The race evaluator ran freshly in 86.823s and integration in 227.274s; timeout-policy and promotion packages reused their successful test cache. Native desktop build and all 167 frontend tests across 28 files passed. Uncached coverage passed engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%. Website tests and static export passed with 23 canonical documents and internal link/asset audit.

Application, GUI, coverage and documentation gates ran on native Windows through the verified hidden launcher; portable automation used the existing private WSL policy-shell bridge. The bridge is not repository source. All 760 frozen application/test/workflow/runner input hashes remained unchanged during and after execution. Canonical log SHA-256 f2d33371e006c2e25adf2f5c39a5752a317dcf569681a137c0265a76c67d67c7. The 67 changed files passed UTF-8/no-BOM/mojibake checks; final publication formatting and lifecycle consistency were checked again after metadata finalization. Final source preflight, comparison content binding and lifecycle consistency across 114 specifications passed with native exit 0. This is a warm dependency/build-cache development run, not cold-cache qualification or a replacement for hosted exact-head CI.
