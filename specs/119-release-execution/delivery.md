# S119 release and source delivery

## Authorities and identities

The owner authorized source push/one official PR and at most two Codex rounds, stopping before owner merge. Separately, the owner approved immutable v1.6.0 tagging and draft staging at 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. On 2026-10-27 the owner separately approved the actual staged candidate presented in #354 comment 6409559198: release 423908120, source 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b, staging run 40278046423 attempt 1 and all eight original assets.

The approved tag was created/pushed and resolves to that exact commit. [Release staging run 40278046423](https://github.com/shruggietech/go-schedule/actions/runs/40278046423), attempt 1, completed successfully through the existing tag-only pathway, with all seven jobs successful. Draft release ID 423908120 has all eight expected original assets. Every asset was downloaded, hashed and matched to its GitHub digest; the Windows candidate manifest passed existing exact-source validation. See [candidate.md](candidate.md) for complete identities. [Guarded promotion run 40450478494](https://github.com/shruggietech/go-schedule/actions/runs/40450478494), attempt 1, completed successfully from the same reviewed main source. The same release became public at 2026-10-27T04:35:10Z and is latest. All eight original asset IDs, sizes and SHA-256 values remain unchanged. SHA256SUMS.txt (asset 680449833, 844 bytes, SHA-256 42db590df8671d04baa1bcd4119ce51feed36c282daf1106ccf4023750e8950c) contains eight entries independently matched to the downloaded staged bytes.

## Individual issue outcomes

| Issue | Current S119 outcome | Remaining outcome |
| --- | --- | --- |
| #354 | Exact approved release published and latest; discovery reconciliation prepared in this source branch | Source discovery PR merge remains owner-controlled; closes with completed source reconciliation |
| #355 | Original runtime/prerequisite evidence and unconfirmed causes remain independent | Demonstrated project correction or supported upstream disposition; remains open |
| #356 | Original IPC/inspection evidence and unconfirmed stall cause remain independent | Demonstrated causal correction; remains open |
| #379 | Reproduced obsolete comparison-version browser assertion and corrected it to actual prepared metadata; both selection paths passed | Source PR delivery, closes with the merged correction |

No installed-system, account-session or DNS operation is part of this slice. Original failed captures retain their original outcomes. [verification.md](verification.md) records actual source/candidate checks and their limits.

## Immutable release and reconciled source records

The immutable tagged comparison changelog digest remains 26fa24a30c70cc5f401213457382c9f57e10bff8a7e0e85e9b925ad7712aa269. Current-main source corrects the publication introduction and refreshes its reviewed digest to f784d30ea8b3e12f03a07154e817d0ed8c588d49a22f481979da0fa4a2e332e1 with review date 2026-10-27. The comparison dataset (ab2bb5bf2de20133735fa4cda541e6d48e557e4632edbd8d95dd8ca14af0cb5b), original release notes (4062f9e664f32298f9f4de8dbeacf4c0b317e0e003c83b185290e18a1c8874b1), decisions and tag source remain unchanged. The v1.6.0 changelog heading/anchor and dated historical decisions are preserved.

Website deployment [40450554565](https://github.com/shruggietech/go-schedule/actions/runs/40450554565) completed successfully after promotion on the accepted main source. Its release-driven selection is separate from this unmerged source-discovery reconciliation.

## Source implementation acceptance

All eleven implementation tasks are complete on review branch codex/119-release-execution. All eight foreground gates passed in 703.326s with all 2,360 frozen inputs unchanged, followed by factual lifecycle finalization. Source publication and at most two Codex rounds are owner-authorized; actual PR/check/review receipts will be recorded on the official PR. The branch remains separate from the immutable release tag and owner merge.
