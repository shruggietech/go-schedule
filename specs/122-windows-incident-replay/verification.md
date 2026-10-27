# S122 verification

## Identity and original evidence

Accepted baseline is owner-merged S121 6f1de579bed866eb2bb328a25a0ad165caaa278a. The bounded baseline campaign uses exact prepared source digest 81198382ca942bc8c6fa666442836726870382fd8dd31d7b370fa1d8ae824089, dirty=true, with source/test binary hashes in its private receipt. Its before/after source manifests have zero drift, and all thirty original/finalized evidence artifacts remain unchanged. The initial wrong-working-directory launch stopped before writer setup and is retained separately; the original campaign deadline was not reset. One component fixture setup validation preceded the campaign and is separately retained.

## Demonstrated regressions and corrections

#388 failed before correction: existing invalid executable, native exit -1, PID zero, phase command, collection absent and fabricated crash.0.dmp inventory. It passes after truthful launch/unavailable classification, retaining the original launch error and executable hash.

#389 failed before correction for already canceled enable/rotate/disable and queued rotation. The queued test waited through its thirty-second observation bound and then replaced the credential after cancellation. A second first-failure regression showed cancellation during session creation still publishing a listener and retaining an uncommitted session. The corrections pass pre-admission, queued and pre-commit cancellation regressions. Independent committed Disable cleanup remains passing. A test context acknowledges entry to the cancelable select, making the queued case discriminating rather than scheduler-dependent.

## Bounded baseline investigation

All six groups passed under one shared deadline in 370.326 seconds. There are 480 runtime observations (360 six-writer outcomes and 120 metadata/writer-selection contrasts), one expected prerequisite-fault matrix and 210 IPC catalog observations. The original five targets passed 150 times across memory/WAL component and packaged-daemon WAL modes; sixty group prerequisites also passed. All thirty owned IPC runtimes were removed. The raw receipts and expanded 690-row matrix remain private. No original spontaneous CLR crash, historical prerequisite refusal or intermittent inspection stall recurred; original #355/#356 causes remain unconfirmed.

## Targeted platform contracts

Native Windows race checks passed the corrected launch and lifecycle tests, then the complete mcphttp, incident, API server/client and store packages. Native nonroot Linux on existing Ubuntu-24.04 passed complete mcphttp/incident race checks and the launch-incident regressions. Ordinary permissions, confirmation, SDK discovery, stale-credential rejection and ownership assertions remain active. No installed runtime/service or account session changed.

## Integrated evaluation and canonical acceptance

Fresh corrected-source Windows/WSL evaluation e35a91fa-0c76-453f-b369-5bfc78d66735 finalized, validated and rendered with exit 0: 1,766 planned, 1,763 passed, zero failed/blocked/findings and three not applicable. All 1,109 distinct messages received seven-criterion semantic review (7,763 rubric entries). Frozen source digest is 5609908ca5bc7c40c3f909b23cde4e46a77459d0a6130b9fd85362f68921605e; corpus digest is 76196c62b88ed3b1fab418f35e4f3c2a583b603e69b6b2854aee7a48bdc78157. The original coordinator, all six prepared inputs, evidence manifest and generated reports remain under .test-temp/s122/qa-run-01 and qa-inputs-01.

After finalization, only a test replay diagnostic field was renamed from original_error to result_error to identify its combined error honestly, and its recipe wording/table escaping plus factual records were corrected. No CLI or daemon implementation changed after the evaluation freeze. The bounded campaign retains its original empty original_error fields and immutable receipt. The complete foreground canonical aggregate follows these corrections. Its before/after source and original evidence identities will be recorded when complete.

### First canonical attempt and documentation correction

The first full foreground aggregate ran from 18:22:31.749Z to 18:32:29.565Z for 597.816 seconds and returned exit 1. Format, vet, lint, race, GUI and all six coverage thresholds passed. Website preparation rejected the new relative maintainer recipe link in cli-qa.md because repository-only maintainer files are outside canonical website routes; automation had not run. Separately filed #390 tracks this pending-source documentation defect. Both new recipe links now use the existing absolute repository-link convention. No application implementation changed. The first receipt/log/stderr and before/after source manifests remain under .test-temp/s122/canonical-attempt-01*, with zero source drift and all thirty original/finalized historical artifacts unchanged. A complete corrected aggregate follows.

### Corrected canonical acceptance

The complete corrected foreground sh scripts/verify.sh all passed all eight gates from 2026-10-27T18:34:52.795Z through 2026-10-27T18:50:51.116Z, 958.321 seconds, native exit 0, 2419 frozen inputs, zero source drift and all thirty original/finalized historical artifacts unchanged. Log SHA-256: 5e56225532077473ecd3c015d3574aef02228d6b470a0b343371c765a0b53698. Stderr SHA-256: 1f4eeb2047516abe32d4f39a06ace851db8f8e830caecccb3741250dcb4e0204. Core coverage: engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.6%, catchup 88.9%, logbus 91.1%. All 167 frontend tests, native desktop build, 23 canonical documentation pages, website static export and automation negative fixtures passed. The first exit-1 documentation attempt remains separately retained.

Only factual completion, task/lifecycle and delivery records follow this frozen aggregate. Final GitHub prose and lifecycle validation run after those metadata edits; no application or test implementation changes follow the successful aggregate. Hosted CI/reviews remain separately attributable to their exact published head.

### Finalized evaluation identity

| Artifact | SHA-256 |
| --- | --- |
| report.json | 953c4589d8e49de2b00460c125e1194ffcf6be0380be439a18f633fc55ae11f9 |
| report.md | 9c410899d27e6f28c729f5cc5afb504ea0c02ddea8cd4f98900e8a1c4881c148 |
| corpus.json | 76196c62b88ed3b1fab418f35e4f3c2a583b603e69b6b2854aee7a48bdc78157 |
| manifest.json | c6723dd5a8c1650ca756b008862b6f80691dd38128ae808690512b00a02f1e9c |
| evidence-manifest.json | 29e333cfebed2ac20b659eb1e7db8b9b8de53af00f46a9250bdf5e34d53b0383 |
| capture-identity.json | f95b125d22177e21b3d5e336c4cb0abd1b1738d3d88152c546408d7595705f5d |

| Prepared input | SHA-256 |
| --- | --- |
| windows/cli | 12b2c57e156bf81602da68489b4aaba201fe7544cb3df49e36c5081fa7c0b7a6 |
| windows/daemon | 772148e234ec4b9b77e550ec87d2513117794c8fc315d3cfa87017a0a42a2e0c |
| windows/runner | cb3303da2b34935ef45b66df01e068641c57d32c44f998c1197cab52c3c11c11 |
| linux/cli | 7eb2ed124903250d83fc2e4856e8feca26e8fc4a28630aaa3406b88a3fad2e2c |
| linux/daemon | 544c46aa6795a22da4ccef1f85cd46e53a2de2250268a7f8f1eb5d320f6fe70a |
| linux/runner | a57f81c26225b0fd0beb52fcab01bd3c5a4a38c2a9b10da12a3c4cc16ab6cae9 |

All four complete-catalog owned runtimes were removed, all cleanup actions passed and both environments report zero remaining resources. Report capture ran from 18:02:07.3862083Z to 18:07:33.2689516Z. Capture-pending-review was not treated as a final pass; all 1,109 current messages were read in occurrence context through factored paragraphs/JSON leaves, then exact full excerpts and complete cases were submitted and independently checked by the original coordinator.

## Hosted CI discovery and correction #392

Both authorized Codex rounds completed without findings on initial PR #391 head 8b0dda27e0569f73bddfc5cd82676c59e9819540. CodeQL run 40564151195 succeeded. CI run 40564151092 coverage job 122400789413 failed TestOwnedMetadataUsesSelectedCompletionAllowance before percentage evaluation with evaluated source changed or could not be read while hashing. The hosted message does not name the file. Source inspection and a separately executed regression establish that the canonical driver places its mutable profile in the unignored root source inventory: the first regression failed with coverage profile is an unignored mutable source input: ./goschedule-cover.ZQ3c4s. The original hosted log and regression remain private under .test-temp/s122/.

#392 tracks this separate CI artifact-boundary defect. The canonical profile now lives under existing ignored .test-temp while retaining a repository-relative native Windows-compatible path. The regression exercises the actual driver, ignored-source boundary, cleanup on success/failure and original exit-7/error preservation. Windows passes after correction. identifySource, legitimate source inclusion/change rejection, coverage thresholds and deadlines are unchanged. No application implementation changed; no fresh full-catalog claim is inferred for the later build-script source. The prior application QA remains exact to its recorded frozen input; the complete canonical rerun follows this CI-only correction. The two-round Codex cap is exhausted, so no third request will be made.

### Post-hosted-correction canonical acceptance

The actual artifact-boundary regression passed on native Windows and existing nonroot Linux. The complete corrected foreground aggregate then passed all eight gates from 2026-10-27T19:04:10.513Z through 2026-10-27T19:19:44.592Z, 934.079 seconds, native exit 0, 2420 frozen inputs, zero source drift and all thirty original/finalized historical artifacts unchanged. Log SHA-256: 853dd217fefc3d134862c8862e3cf1f399d5577ac5d56b19306d4e6f98961ba3. Stderr SHA-256: 12a61d180ce7637b377368570625db1957cd3590d8ada06f16e50d891d9d3274. All 167 frontend tests, native desktop build, six core coverage thresholds, 23 canonical pages/static export, the new profile-isolation success/failure contract and the existing automation negative fixtures passed. The prior successful pre-hosted receipt remains under canonical-pre-hosted*, and the first documentation failure remains under canonical-attempt-01*. No implementation changes follow this latest aggregate; only factual completion/lifecycle/delivery records do.

Both Codex rounds remain explicitly bound to initial head 8b0dda27e0569f73bddfc5cd82676c59e9819540 (automatic completion 18:54:40.102699Z; manual round-two completion 18:56:32.541803Z, receipt comment 6450901492), with no findings or unresolved threads. The later #392 correction changes scripts/verify.sh, its new shell regression and factual records only; it is not represented as another externally reviewed application head. Final published-head CI is retained on PR #391 and in the private ledger. No third review, merge, tag or promotion is authorized.

## Owner merge and main acceptance (reconciled by S123)

The owner merged PR #391 on 2026-10-27 as e8cefc1c8059364a0ee84fe22ac740842885d53d and deleted the remote branch. The squash tree exactly matches final source head 47181947868271325294ddf900fdf39f7cc8ddf1. All 23 merged-main checks succeeded: [CI 40575544962](https://github.com/shruggietech/go-schedule/actions/runs/40575544962), [CodeQL 40575544974](https://github.com/shruggietech/go-schedule/actions/runs/40575544974) and [Public website 40575545135](https://github.com/shruggietech/go-schedule/actions/runs/40575545135), attempt 1. Pages deployment 7369221764 succeeded on that same main commit; six public routes passed certificate-validated HTTPS. [Post-merge receipt](https://github.com/shruggietech/go-schedule/pull/391#issuecomment-6452630491) records clean synchronized main, branch/helper cleanup and unchanged retained evidence.

#388/#389/#390/#392 are closed by the merge. Both clean Codex rounds remain attributable to initial application head 8b0dda27e0569f73bddfc5cd82676c59e9819540; later 4718194 changes only CI scripts and factual records for separately demonstrated #392. No third review was requested. Original #355/#356 causes remain unconfirmed; merge and source acceptance do not close them. Earlier pre-publication statements above describe their historical stage, not the current PR state. No tag or promotion occurred.
