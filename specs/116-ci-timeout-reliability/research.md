# Research: S116

## Increased completion policy

Decision: apply the complete [timeout table](contracts/timeouts.md) across all six workflows and their local/hosted runners. Rationale: the owner explicitly directs increased allowances and nested old caps can terminate work before outer jobs. Alternatives: only fixing the single macOS test leaves other CI systems unchanged; inventing an outer value above GitHub maximum cannot work. Keep already-unlimited runners unlimited.

## Source of #372

Decision: raise frontend whole-test/hook/teardown to 60s, DOM lookup to 10s and remove the isolated 5s override. The actual main receipt shows a 5s whole-test expiry, not an identity mismatch assertion failure. Preserve production 6500ms reconnect and identity code. Add actual delayed-work verification and strengthen exact-target assertions without retrying navigation. Root cause beyond the demonstrated old cap remains an inference.

## Native evaluation

Decision: use the existing supported maxima 120s/case and 2 h execution as defaults, keeping exact supplied/resumed budgets authoritative. Scale completion reserves together so evidence freeze 120s fits hard worker 180s, aggregate 240s and outer drainage 60s. Retain original failed report and use targeted real IPC/CLI metadata replays plus six-writer diagnostics before considering a full catalog. Alternatives: modifying production API limits or serializing writers lacks evidence and changes unrelated contracts.

## Release handoff

Decision: include delivered S115/S116 source changes in the selected pending v1.6.0 changelog/notes, refresh only governed content digests and retain dataset bytes/public v1.5.1 selection. Record actual S115 failed-main versus passed-PR evidence separately. No merge/tag/staging/promotion authorization exists.

## Planning research provenance

Spec-Kit required three read-only research agents: CI timeout inventory, native QA and search/release seams. All completed without edits, GitHub mutations or test launches. All unknowns are resolved sufficiently for implementation; original runtime/IPC causes remain explicitly unknown findings, not design placeholders.
