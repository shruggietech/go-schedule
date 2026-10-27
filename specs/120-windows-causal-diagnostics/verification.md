# S120 verification

## Targeted evidence

The inherited-stream regression failed before correction for native exits 0 and 37. The real native IPC pool-wait regression failed before correction in both memory and file-WAL stores. After correction, targeted native Windows race checks passed across store/API/script packages, all twelve queued SQL cancellation boundaries, all four original request shapes, terminal audit persistence after a committed mutation, bounded inherited streams, output overflow and unchanged concurrent writer contracts. Separate corrected audit/contrast checks passed. Existing memory/WAL snapshots and concurrent mutations passed on Windows and Linux.

Existing nonroot WSL Go 1.26.0 race checks passed controlled component and Bash regressions; Linux PowerShell was unavailable. Windows runtime identity is PowerShell 7.6.9/.NET 10.0.14, distinct from the original inventory. Full attributable dispositions and first failed outcomes are in [investigation.md](investigation.md). Original evidence is immutable.

## Canonical acceptance

All eight foreground gates passed on 2026-10-27 from 2026-10-27T07:19:04.259Z to 2026-10-27T07:32:14.519Z, native exit 0, 790.26 seconds. The input manifest has 2380 entries (SHA-256 1ce1cd8132db3066ffc72d4672e514cb14efdf64178484e9f0486d370337c2db); before/after drift is empty and all 26 original evidence hashes remain unchanged. Gates: format, vet, lint, race, gui, coverage, docs and automation. Native GUI build, 167 frontend tests and production bundle passed; all six core coverage thresholds passed (store 80.6%). Documentation exports 23 canonical pages and automation validates 22 hosted jobs plus nested completion allowances. The first 25.508-second aggregate stopped at a new test's import grouping; that receipt is preserved separately and the grouping was corrected before this complete successful run. Only factual lifecycle, task and verification metadata are finalized after the frozen-input result.

## Hosted delivery boundary

PR #383 is open under the owner's source-publication authorization. Both permitted Codex rounds completed without findings on implementation cb48f4813d87fde27c245b53e807604e609d0562. The metadata-only correction consumes no further review request. Exact-head hosted checks follow; stop before merge, tagging or promotion.

## Post-publication lifecycle correction (#384)

Initial implementation cb48f4813d87fde27c245b53e807604e609d0562 received two completed Codex rounds without findings. Hosted CI 40467905081 correctly rejected its finalized Delivery field, because the pre-finalization aggregate exercised In Progress metadata and finalization omitted an objective PR/branch reference. The unchanged checker requires that reference for Implemented. The correction binds the actual PR #383 and review branch; application and test source are unchanged. The original failed hosted log is retained separately.

The repeated full foreground aggregate passed all eight gates in 534.88 seconds, from 2026-10-27T07:40:52.725Z to 2026-10-27T07:49:47.605Z, native exit 0. All 2380 frozen inputs and 26 original artifacts remained unchanged; manifest SHA-256 05798dbd918663108693b2e9f2e1e81991848330b145f0b38080720b4c5c51a2. This run exercised the corrected Implemented lifecycle metadata. After recording this factual receipt, final formatting and the unchanged automation/lifecycle gate run again before the correction commit. Both completed Codex rounds reviewed implementation cb48f4813d87fde27c245b53e807604e609d0562; the correction changes delivery metadata only and requests no additional review.

## Post-merge reconciliation (2026-10-27)

The owner merged PR #383 as 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7 and deleted its remote branch. Main CI [40498355544](https://github.com/shruggietech/go-schedule/actions/runs/40498355544) passed all 20 jobs; CodeQL [40498355556](https://github.com/shruggietech/go-schedule/actions/runs/40498355556) and website [40498355692](https://github.com/shruggietech/go-schedule/actions/runs/40498355692) passed on that merge SHA. All 23 check runs were terminal/successful, Pages deployment 7357551490 succeeded, and five live canonical HTTPS pages returned 200 with the new documentation visible. #381/#382/#384 closed; original #355/#356 causes remain unconfirmed and open. Public v1.6.0 release 423908120 and immutable tag source are unchanged.
