# S110 Verification

## Subject and scope

Baseline main: 95c05692e41bdb3630aa0b3c1bf946e5d443c2af. Review branch: codex/110-reliability-release-repairs. Results below are local engineering observations. At the local implementation boundary, no branch push, hosted review, live promotion, tag or release had been performed. Subsequent owner-authorized publication is tracked in [review.md](review.md).

## Failing-before evidence

- #262: TestWaiverPromotionBindsTrustedHelperBeforeTagCheckout fails against the unchanged workflow (trusted=-1, build=-1, tag=897). The contract also refuses independent tag-resolving release paths. Evidence: .test-temp/s110-release-before.log.
- #271: The unchanged runtimeState/fake-clock recovery reproduction loses the replacement (zero dispatches, one expected). The overlay and promoted regression are retained in .test-temp/s110-watcher-research/{result,tracked-before}.log.
- #335: Distinct valid hexadecimal session IDs, 1e followed by 30 nines and 2e followed by 30 nines, both become text Inf through the original SQLite CLI adapters. The second heartbeat fails with UNIQUE constraint failed. The six-writer regression retains only one of six records. Exact-value assertions also fail for leading zeroes, arithmetic-looking text, apostrophes and newline/dot-command text. Both actual Windows twins fail in .test-temp/s110-binding-before.log.
- Retry diagnostics: Four missing-final-error assertions fail before correction, while attempt counts and waits already meet their contract. Evidence: .test-temp/s110-heartbeat-research/retry-before.log.

The original hosted #335 run discarded writer diagnostics, so its exact historical cause remains unconfirmed. Two hundred additional isolated repetitions of the unchanged fixture passed during research. These repetitions are not correction evidence. S110 corrects the separately demonstrated binding data loss and retains actual final SQLite errors.

## Focused corrected checks

| Subject | Environment and result | Evidence under .test-temp/ |
| --- | --- | --- |
| Complete watcher package | Windows race detector, pass (6.670s); includes existing native atomic replacement checks and nine deterministic groups with 24 recovery scenarios | s110-watcher-research/watcher-race.log |
| Exact bindings and concurrent IDs | Actual Windows PowerShell and Git Bash twins with native SQLite 3.47.2, race pass (15.606s), including concurrent first use and initialized databases | s110-binding-after.log |
| Exact bindings | Native Ubuntu-24.04 WSL, Go 1.26.4 and SQLite 3.45.1, race pass (2.597s) | s110-binding-linux.log |
| Retry/error contracts | Both actual Windows twins, all eight transient-success, locked/busy exhaustion and non-contention scenarios pass (3.401s) | s110-heartbeat-research/retry-after.log |
| Release helper and workflow contract | Windows race passes (2.559s and 1.058s) | s110-release-focused.log |
| Complete watcher and release helper packages | Native Linux race passes (6.407s and 1.927s) | s110-native-linux-core.log |
| Entire real Bash payload suite | Native Linux race pass (10.949s), including binding, concurrency and retry checks | s110-native-linux-scripts.log |

Watcher checks cover full/root/repeated recovery, pending debounce and stability, inactive timer behavior, previous dispatches and late hints, later ordinary writes, changed/disabled/deleted definitions, explicit reload, absent/create behavior, health revisions and rejected dispatches. Release fixtures cover paginated identity selection, ambiguity, changed identity/inventory/bytes, interrupted checksum upload, published latest-only recovery, redirect credential removal, upload origin/ID validation, finite timeout and redacted permission diagnostics. Payload checks preserve the original simultaneous success and record-count assertions; retry limits and waits remain unchanged.

## Canonical aggregate

The foreground command sh scripts/verify.sh all completed on 2026-10-21 with exit code 0. Session 82297 retained the process until completion. The log is .test-temp/verify-s110.log (SHA256 f8b8b97aabe6b1dfacea8e7efcbe736c0bbb25c2ea9231faa70eb0bca78725b1); its terminal receipt is .test-temp/s110-aggregate-exit.json. All eight required gates passed:

| Gate | Result |
| --- | --- |
| format | PASS, Go formatting and GitHub prose/punctuation checker |
| vet | PASS |
| lint | PASS, zero issues |
| race | PASS, complete root-module race suite |
| gui | PASS, desktop race tests, native Wails build, 28 frontend files / 160 tests and production frontend build |
| coverage | PASS, engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1% |
| docs | PASS, 22 canonical documents, links, metadata and static export |
| automation | PASS, repository policy and complete negative-fixture suite |

Only delivery and lifecycle records were finalized after this run. Final formatting, lifecycle and automation-policy checks are recorded below; no implementation source was changed after the aggregate.

Earlier attempts are retained separately. The first stopped at an import-order lint failure, corrected before the final run (.test-temp/verify-s110-import-order.log). The second was interrupted by a tool runtime limit during automation fixtures and has no terminal aggregate result (.test-temp/verify-s110-tool-interrupted.log). Its owned orphan fixture process tree and temporary directory were removed before the replacement run. Neither attempt is counted as aggregate success.

## Local prerequisites and boundaries

Windows used the existing Go toolchain, GCC, PowerShell 7.6.6 and Git Bash. Portable SQLite 3.47.2 was fetched into ignored scratch from the existing repository manifest and verified with SHA256 8c7fffbf4eec1f43e63153cca6deb018e4360d4d6b0d99bbdd2a541c53b7fa1c before extraction. No global installation or host service was changed. Isolated tool environments restore ordinary executable suffixes because the inherited PATHEXT=.CPL prevented PowerShell native executable invocation. Empty initial output is not counted as a passing check. Payload children are hidden and PowerShell is non-interactive.

Native Linux PowerShell was unavailable; its twin ran natively on Windows. The original macOS observer degradation trigger is not established by the platform-independent state-loss reproduction. Native macOS was not exercised locally; its subsequent successful hosted result is recorded in [review.md](review.md). Release mutations were modeled with HTTP fixtures; no live release promotion was attempted.

## Final record audit

After finalizing lifecycle and delivery records, go run ./scripts/github-format, sh scripts/spec-lifecycle-check.sh . and sh scripts/automation-check.sh . each passed with exit code 0. The lifecycle checker reports 108 consistent specifications. All 24 S110 tasks are complete. The 36 changed text files pass fatal UTF-8 decoding, BOM, mojibake and em-dash checks; git diff --check passes. The automation fixture scratch directories were removed by the successful suite. No implementation source changed after the aggregate.
