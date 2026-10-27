# S120 source delivery

Owner-authorized scope is #355/#356, with demonstrated child defects [#381](https://github.com/shruggietech/go-schedule/issues/381) and [#382](https://github.com/shruggietech/go-schedule/issues/382), plus the separately demonstrated source-record defect [#384](https://github.com/shruggietech/go-schedule/issues/384). Implemented corrections address inherited output collection and lost request cancellation at affected SQL boundaries. The historical CLR/prerequisite and intermittent IPC causes remain unconfirmed; #355/#356 remain open. The source PR will close only the delivered child defects on merge.

S119 PR #380 and exact-main CI/CodeQL/Website/Pages delivery are reconciled as historical metadata. Public v1.6.0 source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b and release 423908120 remain unchanged.

See [spec.md](spec.md), [investigation.md](investigation.md) and [verification.md](verification.md). Source publication and at most two Codex rounds are authorized. The final PR will retain its exact checked head and review receipts; merge/tag/promotion remain excluded.

## Completed local source acceptance

All 19 implementation tasks are complete. All eight canonical foreground gates passed in 790.26 seconds with 2380 frozen inputs unchanged, followed only by factual lifecycle finalization and formatting. All 26 original evidence hashes remain unchanged. The source is ready for the already authorized branch push and official PR; exact hosted CI and capped review receipts belong on that PR.

## Authorized publication and lifecycle correction

[PR #383](https://github.com/shruggietech/go-schedule/pull/383) was opened on implementation cb48f4813d87fde27c245b53e807604e609d0562. Both authorized Codex rounds completed without findings on that implementation. The second security-focused request was recorded by the bot as Manual Code Review; no separate dedicated security verdict is claimed. Hosted CI correctly found #384: the finalized Implemented Delivery field lacked the objective reference required by the unchanged lifecycle checker. This follow-up adds the actual PR/review branch and retains the original failed hosted receipt. It changes only source delivery metadata. No third review is requested. #384 is separate from the original #356 cause.

The #384 metadata correction passed a second complete frozen-input aggregate in 534.88 seconds, with all 2380 inputs and all 26 originals unchanged. The final post-receipt formatting and lifecycle recheck precede its correction commit. The original implementation's successful checks and both no-finding Codex receipts remain independently attributed to cb48f4813d87fde27c245b53e807604e609d0562.

## Post-merge reconciliation (2026-10-27)

The owner merged PR #383 as 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7 and deleted its remote branch. Main CI [40498355544](https://github.com/shruggietech/go-schedule/actions/runs/40498355544) passed all 20 jobs; CodeQL [40498355556](https://github.com/shruggietech/go-schedule/actions/runs/40498355556) and website [40498355692](https://github.com/shruggietech/go-schedule/actions/runs/40498355692) passed on that merge SHA. All 23 check runs were terminal/successful, Pages deployment 7357551490 succeeded, and five live canonical HTTPS pages returned 200 with the new documentation visible. #381/#382/#384 closed; original #355/#356 causes remain unconfirmed and open. Public v1.6.0 release 423908120 and immutable tag source are unchanged.
