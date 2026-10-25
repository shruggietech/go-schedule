# S117 verification record

Date: 2026-10-25. Accepted main baseline: 11485a8a6ce97fe8e6966e1cc3cdb6a55b9d5651. This is candidate/source evidence; it does not identify a tag, staged package or public v1.6.0 release.

## Chronological execution

1. Completed the installed Spec-Kit specify, clarify, checklist, plan, tasks and blocking analysis sequence before implementation. Both 16-item checklists passed; all 12 requirements were mapped to 18 tasks with no blocking analysis findings. The feature selector, branch and current plan agree on S117.
2. Verified existing headless tools and nonroot Ubuntu-24.04 WSL controls without provisioning. Windows Go 1.26.0, PowerShell 7.6.6/.NET 10.0.12, Bash 5.3.15 and SQLite 3.47.2 were available. The existing Linux environment is Ubuntu 24.04.3, kernel 6.18.33.2-microsoft-standard-WSL2; the compiled Linux evaluator identifies Go 1.26.0.
3. Froze accepted main plus the finished planning-only working tree, built six immutable binaries and ran one full existing-coordinator capture without budget overrides. Source stayed unchanged through contextual review, finalization, validation and rendering. No full replay was used to seek a preferred outcome.
4. Capture 76051e90-e129-48b0-8dad-63646360ad97 started at 2026-10-25T12:33:58.743243Z and finished at 2026-10-25T12:39:12.7972597Z. All 1,763 executable cases passed (881 Windows, 882 WSL), with three explicit platform exclusions. CPU, memory, process/task, output and ownership controls were verified. All four owned runtimes were removed, with no remaining resources.
5. Reviewed all 1,109 required fresh messages against their current command/case contexts using the seven fixed rubric dimensions. Exact shared text was read once with every distinct context and changed value inspected. There are 43 explicitly silent cases, with no fabricated message reviews. Six stdout records expose one missing-summary product defect (#375); the other 1,103 message records pass.
6. Finalize, validate and render each returned native outcome 2 for a valid finalized failed report: 1,757 effective passes, six communication failures, zero blocked, three exclusions and one finding. The outer hidden launcher maps the child failure to shell exit 1; native outcome 2 is retained in each receipt. The report is not malformed and is not relabelled passed after corrections.
7. The real persisted/authenticated API-to-MCP regression failed before correction for recurring, one-off, startup and missing-date descriptions. Corrected SQL loads only the byte-bounded human description with a source-side truncation flag; the API composes policy wording and UTF-8 bounds. The typed Observe session and both resource projections preserve descriptions and exclude execution-input canaries. ASCII and Unicode truncation plus manual-only behavior are covered.
8. Supplemental inspection of the frozen authority-state artifacts found discarded valid JSON arrays. A focused regression reproduced array loss and a panic on JSON null before correction (#376). Diagnostic attachment now preserves object paths, wraps non-object observations, retains exact JSON numbers and retains malformed raw state with its explicit decode error. This supplemental finding does not rewrite the original finalized report.
9. Corrective race acceptance passed store, API server, MCP observer and evaluator diagnostic contracts. Native Windows follow-up passed the six script-binding/prerequisite tests in 88.974s: both PowerShell/Bash, fresh/existing schemas, six simultaneous writers, exact identities/rows, opt-in diagnostics, candidate failures, launcher errors and inherited-pipe bounds. No original CLR/prerequisite failure was reproduced or causally corrected.

## Frozen candidate identity

Commit: 11485a8a6ce97fe8e6966e1cc3cdb6a55b9d5651. Dirty: true, planning-only source snapshot. Source SHA-256: 98514d388059638cf08b00ee5c93e82fca24d4437ccca426172f822200082d18. Actual binary version: gosched version 0.0.1-dev. These inputs predate the corrective source and final PR.

| Role | Windows amd64 SHA-256 | Linux amd64 SHA-256 |
| --- | --- | --- |
| CLI | 61b385964f033de142f804f181773fcfdbd383a75fed19fcb0d13c4fcd6c51c5 | cd7cfeecaaa4c6de73ab8437e3099fd32d341cfadcaa5331ebb2eed2bf2730ec |
| Daemon | 7c3713bef21860010d1d114a277ebe7c0b4c57a3b1050808dbca7c6f9f91a5f3 | 95b6f78202cd9b6ed5e4d8048ed81f428345424fe24a4b8d40e7db6026e9feca |
| Evaluator | 501613d66a8c73752a93c14890d474574b024a4e044c8f6d591063c55b403b29 | 4f8523ea90df8d475f1da849693f0650ed1e7aaeb69fa6f2635945e8e0a3e9bb |

Supported selection: 120000ms/case, 7200000ms total, 25% CPU, 1GiB memory, 32 Windows processes and 128 Linux tasks. Native Jobs and user cgroups were verified. Explicit/resume budget identity, intentional expiry assertions and installed-system isolation remain unchanged.

| Immutable record | SHA-256 |
| --- | --- |
| report.json | c825e902cb934c2e3bcf7dcdffc7e4f34f3e71e81de894636e0a453f98d11c7b |
| report.md | 91eb8029766c12d12c4439d8903ff498450a87e26c8199e61221e2700a703e2b |
| corpus.json | 3d5128a68399f0f0babd017b90442ee402bb17f5cb58429aadde16ec1ec18686 |
| reviews-01.json | e6774a2869dcb13ee74a836e83d50a35b0f2351ef35bf3e56dc6e2685ae01550 |

Private originals remain under .test-temp/s117/qa-run-01, qa-inputs-01 and reviews-01.json, with the hash receipt qa-frozen-01.json. Raw machine-specific artifacts are not published.

## Independent original causes

#355 remains open: the original CLR access violation and separate prerequisite refusal have no demonstrated cause or supported upstream resolution. Current native successes preserve all writer assertions and do not establish a fix. The original failed log SHA-256 remains 4eded1dd882f92b7f09aad62265065c46631212d79d97ce7b12922a4b2d3cacf.

#356 remains open: fresh Windows service.restart (102ms), mcp.http.enable (229ms), mcp.http.rotate (226ms), mcp.bad-port (150ms) and group.remove-child (144ms) passed with complete before/after observations and cleanup. The three previously failing S115 metadata cases passed: help.audit 117ms, help.profile.default.reset 114ms and help.profile.default.set 125ms. All their recorded ownership/API requests completed. These observations and supported larger allowances do not identify the original intermittent stall cause. The service adapter recorded intent without operating the installed service.

The original S115 full report remains failed and unchanged: report JSON e378151781ac1a4433eb790e68108d455a92e17e50559035946c99d17896169b, Markdown 106dde5e3bfcec68c3e3048b597c0f4997912ea21c25853369959e6ff27e73df and corpus e8abd02d91b9a1ab8dc4584ae014d61e9ca68182f99e63f80966b60fd1c6ed49. Its separately corrected communication findings and original deadlines are not replaced by S117 results.

## Canonical source acceptance

The foreground scripts/verify.sh all run completed with native exit 0 from 2026-10-25T13:06:34.841Z to 2026-10-25T13:21:30.253Z (895.412 seconds). All eight gates passed: format, vet, lint (zero issues), race, GUI, coverage, docs and automation including negative fixtures. All 2334 frozen repository input hashes remained unchanged during the run. Evaluator race acceptance ran freshly in 79.619s, store in 100.945s and integration in 219.437s. Native desktop build, all 167 frontend tests across 28 files and 23-document website export passed. Uncached core coverage passed engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%.

Application, GUI, coverage and documentation gates ran on native Windows through the verified hidden launcher; portable automation used the existing private WSL policy-shell bridge. This is a warm dependency/build-cache development run, not cold-cache qualification. Canonical log SHA-256: 3bba02d08f4b2fec44245b29abf15406da1775c046920abf98051a38a8578cbf.

Read-only release-source and comparison guards passed for pending v1.6.0. The comparison check received the actual current baseline SHA and validates content bindings; it does not resolve or establish a tag. The exact changelog/notes hashes are 26fa24a30c70cc5f401213457382c9f57e10bff8a7e0e85e9b925ad7712aa269 and 4062f9e664f32298f9f4de8dbeacf4c0b317e0e003c83b185290e18a1c8874b1. Formatting and UTF-8/no-BOM/mojibake audits passed and are repeated after metadata finalization. Only factual verification/task metadata is finalized after the canonical run; no executable input changes.

Hosted exact-head checks and external reviews were subsequent delivery events, separate from the frozen candidate identity. Actual receipts follow chronologically below.

## Subsequent source delivery and owner merge (reconciled by S118)

The final reviewed PR #377 head was 7327a7c401b6140d0b6c02d2769ea3082d74a123. Its 22 checks passed, and two Codex rounds completed without findings. The owner merged at 2026-10-25T17:48:41Z as 9ad13821c825d01efcc038b6b1634b187a854758; both commits have the same tree. Exact merged-main CI 40218254702, CodeQL 40218254649 and Website 40218254658 all completed successfully at attempt 1, covering 23 checks. Pages deployment 7316855350 and HTTPS route checks passed for the merge. See [delivery.md](delivery.md) for linked runs and actual review classification.

These subsequent receipts preserve every original failed capture, candidate source/binary hash and later corrective acceptance above. S118 reconciles lifecycle metadata rather than rerunning or reassigning that evidence. Tag, staging and public v1.6.0 package identities remain absent from this source delivery.
