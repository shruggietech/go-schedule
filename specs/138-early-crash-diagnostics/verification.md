# Verification: Early-crash incident evidence

## 2026-11-08 local implementation qualification

The installed specify, clarify, checklist, plan, tasks, analyze and implement workflows ran in order. The owner approved the diagnostic scope after historical inspection could not establish original causality. All 24 checklist items passed; analyze mapped seven functional requirements and four success criteria to eight implementation tasks with no blocking conflict. No workflow waiver remains.

| Gate | Result | Evidence |
| --- | --- | --- |
| format | PASS | Go formatting and GitHub prose policy |
| vet | PASS | Canonical root vet |
| lint | PASS | Zero issues |
| race | PASS | Root packages and Windows integration tests, including native crash and bounded query roles |
| gui | PASS | Desktop race, native Windows Wails build, 168 frontend tests and production bundle |
| coverage | PASS | engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.9%, catchup 88.9%, logbus 91.1% |
| docs | PASS | Policy/architecture fixtures, seven site tests and 25-page static export |
| automation | PASS | Timeout/policy/coverage guards, 96 policy cases, lifecycle corpus, 30 public adapters and bounded driver mutation |

The full feature-implementation foreground aggregate started at 2026-11-08T22:20:06.729Z, finished at 2026-11-08T22:33:32.411Z, took 805.682 seconds and exited zero. The same run first executed the exact-source native pre-identity probe successfully. Its before/after audit found zero changes across 2666 frozen source files and zero drift across 757 protected evidence files. Private receipts are canonical-02-result.json, canonical-02-jobs.foreground.log and native-final-02-proof.json under .test-temp/s138.

## Final native proof and artifact readback

The controlled owned PowerShell child PID 42776 exited with native 0x80131623. Its exact process-handle creation/exit interval was 2026-11-08T22:20:09.5629345Z through 2026-11-08T22:20:10.0652662Z. No runtime.json prelude identity existed. The event receipt independently reported collected with one .NET Runtime event 1025, record 233654, timestamp 2026-11-08T22:20:09.7465621Z, application pwsh.exe, CoreCLR 10.0.1426.52326 and .NET 10.0.14. The query took 57 milliseconds. This is an intentional FailFast demonstration, not a recurrence of #355 or #356.

Every collected outer artifact was read back and its byte count/SHA-256 recomputed: runtime-events.json is 689 bytes, dcb14d9e1519c27820140031c6c25e6737526d66d140e84a62216db062fe1471; crash.42776.dmp is 6397094 bytes, 4eb7abbb7308440703aa0bb1ee7df176f7a14bfad1d4648d958f8e172cb3ddce. The native dump was collected on this host; no forced no-dump native scenario is claimed. Receipt independence from dump status is covered separately by the collector regression.

Fifteen attribution cases cover event IDs, case-insensitive application identity, PID/provider/application/time rejection, boundary timestamps and ambiguous/malformed identity. Eight bounded/fault cases cover unavailable queries, timeout, output overflow/truncation, malformed UTF-16, excessive event count and wrong envelope. Native utility roles prove overflow, bounded wait and exit 37 actually execute within the query allowance. Disabled capture, launch failure, success and non-Windows paths do not request runtime events. Arbitrary message/stack data is excluded from the retained receipt.

## Retained development observations

The failing-first parser API compilation, initial indirect import-cycle failure and native missing-artifact failure remain diagnostic observations, followed by corrected focused passes. An early race-instrumented query helper exceeded the small native utility allowance before its role marker; the corrected minimal utility is compiled through the contained launcher while the parent tests retain race instrumentation.

The first full qualifier was deliberately interrupted after identifying an uncontained helper compiler path. Its owned process identity, interruption receipt and partial log remain private and are not counted as passing qualification. Final attempt 02 includes the corrected contained compiler path and supersedes that attempt for source qualification.

An initial development query-bound foreground log was accidentally overwritten by a retry. query-bounds-log-provenance.json records that loss explicitly, and query-bounds-first-failure-transcript.txt is a labeled recovery from tool output rather than an original byte-identical log. Historical inputs were unaffected. Subsequent foreground paths are unique; the final protected-input audit passed.

## 2026-11-08 hosted fixture correction

Initial hosted CI on 6cd1e9dbbd682f3f27ef18adcf167aeb56ea5b24 passed 21 of 22 checks but rejected TestRuntimeEventReceiptIsSeparateArtifact on macOS. The fixture passed t.TempDir() through the platform temporary-root symlink, unlike the package's established canonicalTempDir(t) convention. The collector correctly refused that linked parent. [#443](https://github.com/shruggietech/go-schedule/issues/443) separately tracks this defect against the pending source; the same PR reuses the existing canonical helper for that single setup line. Production code and native query behavior are unchanged.

The failed macOS job is 128885096293 in run 42531051328. Its complete decoded log is retained privately as ci-42531051328-attempt-1-macos-128885096293.log (28621 UTF-8 bytes, SHA-256 bfa021bf6ff9d2eb208c7426aa8cd2a904cc9a45cc300f17ad20d76536431748). The local incident package race suite passed after the correction, including existing caller-link rejection and platform-alias coverage. Formatting and lifecycle checks passed. Full hosted checks on the correction head and the remaining review allowance are recorded on the PR; the earlier full local aggregate is evidence for the feature implementation before this test-only fixture correction.

## Completion records and evidence boundaries

Only factual completion records (specification status/delivery, completed tasks, inventory and this verification/delivery record) are finalized after the full aggregate. Formatting, lifecycle consistency, UTF-8/mojibake and source/evidence drift are checked again before commit. Hosted checks, candidate SHA and provider review are subsequent evidence maintained on the official PR. No Codex review was requested during implementation; the two-request cap includes automatic review.

The recovered October event remains a temporal association because the original writer PID and dump are absent. #355 and #356 remain open. No original cause, shared cause, upstream fix, original-operation success, installed-product qualification or release action is claimed.
