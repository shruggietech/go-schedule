# S110 local delivery and publication handoff

S110 delivers #262, #271 and #335 on codex/110-reliability-release-repairs. The synchronized baseline is merged main 95c05692e41bdb3630aa0b3c1bf946e5d443c2af (S109 PR #332). A final fetch confirmed the same remote main before local delivery.

## Delivered behavior

| Issue | Functional result |
| --- | --- |
| #262 | Authenticated release listing selects exactly one ID; every later asset/checksum/publication operation uses that selection. Interrupted checksum attempts reuse identical bytes, differing contents are refused, and published releases can repair only the latest pointer. The helper is built from the exact workflow revision before tagged checkout. Operator guidance names reviewed main and the refusal conditions. |
| #271 | Observer recovery preserves unchanged enabled definition revisions, prior file identity, pending debounce/stability and last attempted dispatch. Recovery reconciles replacements, suspends inactive candidate timers and suppresses duplicate late hints while preserving later ordinary writes and definition lifecycle behavior. |
| #335 | Bash and PowerShell SQLite adapters preserve exact UTF-8 text, including numeric-looking session IDs. Concurrent first-use and initialized writes retain distinct records. Final SQLite errors survive bounded retry exhaustion, with the existing attempts, waits and simultaneous-writer assertions retained. |

## Decisions and deviations

A small standard-library Go helper with one atomic temporary selection file keeps release identity consistent between workflow steps without a new dependency. The standing-waiver workflow and its policy contract change under a dated CHANGELOG decision; the full attended-evidence workflow remains unchanged. Watcher recovery retains revision-scoped state instead of reusing the explicit Reload reset. Script adapters encode literal bytes rather than changing session-ID generation or speculatively increasing retries. All 24 required tasks are resolved and no functional outcome is intentionally deferred.

The exact cause of the earlier hosted heartbeat failure remains unconfirmed because that run discarded writer output. The scientific-notation ID collision is independently demonstrated and corrected in both adapters; repeated unchanged passes are not represented as a diagnosis. The macOS observer's original degradation trigger is likewise distinct from the reproduced recovery state loss. These boundaries do not change the issue acceptance criteria or conceal unimplemented work.

## Verification

All eight local canonical gates passed with aggregate exit code 0: format, vet, lint, race, GUI, coverage, docs and automation. The GUI suite passed 160 tests; the docs export checked 22 canonical documents. Focused native Windows and native Linux race checks cover all three corrections. Final delivery-only records receive formatting, lifecycle and automation-policy checks. See [verification.md](verification.md) for exact log identity, failing-before evidence, individual results and earlier interrupted attempts.

Hosted CI and Codex review have not run for S110. Native macOS confirmation will use the existing hosted workflow after publication. Live release promotion was modeled with HTTP fixtures and was not performed. No tag, release or public configuration was changed.

## Publication and merge boundary

The initial pre-publication halt completed at local implementation commit 19df5010f697a9a9c4ff2cdedf1e8ab95f68f9b5. On 2026-10-21 the owner authorized pushing S110, opening the PR, monitoring CI and completing up to two Codex review rounds with autonomous corrections, stopping before merge.

[PR #338](https://github.com/shruggietech/go-schedule/pull/338) is open for the complete #262, #271 and #335 bundle. All three issues remain open pending merge with their existing milestones; #335 retains parent #277. Existing bug and functional-area labels identify the PR. Its cross-milestone scope does not create a new release milestone.

The [publication and review record](review.md) tracks exact revisions, review rounds and hosted results. The PR was initially created as a draft while this metadata was reconciled. No merge, auto-merge, tag, release or live promotion is authorized. Local verification claims remain scoped to the implementation checks in [verification.md](verification.md); hosted checks and review are recorded separately.
