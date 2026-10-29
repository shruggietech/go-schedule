# S126 Source Delivery

## Delivered in this source proposal

The selected #401-#403 dependency sets are coherent across both Go modules and both npm projects. Pending v1.7.0 notes/changelog disclose selected and deferred versions. S125's actual owner merge and exact-main acceptance are reconciled, and current Windows guidance points to the [same-slice release continuation](handoff.md).

| Record | Functional outcome | Disposition |
| --- | --- | --- |
| PR #401 | API runtime 1.8.0 and YAML 3.0.5, complete intended Go graphs | Selected equivalent source delivered here; supersede proposal after owner integration |
| PR #402 | Vite 8.3.5 and selected lock graph | Selected equivalent source delivered here; supersede proposal after owner integration |
| PR #403 | Coupled Next/Fumadocs update and complete site lock graph | Selected equivalent source delivered here; supersede proposal after owner integration |
| PR #400 / #404 | Suitable published SQLite update and complete secondary graph | Deferred, no suitable corrected driver published at selection; remain open |
| #394 | Dependency/source convergence and concrete existing release continuation | Partial, actual tag/assets/publication/discovery remain open |
| #355 / #356 | Original Windows historical causes | Unresolved, no causal closure inferred from this candidate |

## Verification and owner boundary

[verification.md](verification.md) records fresh Windows/WSL evaluation with 1763 passes, all 1109 messages reviewed, no findings, all four owned runtimes removed, both read-only Go graphs, clean restores, 35 desktop and 10 site browser cases, and the complete eight-gate foreground aggregate. Frozen source and historical evidence retain zero drift.

The source implementation is complete; publishing/review is its authorized protocol rather than an unfinished implementation task. The final PR must still receive exact-head CI and at most two Codex rounds before owner review. Leave it unmerged. Continue after owner merge within S126 through actual main acceptance, separately approved immutable tag/draft, real asset inspection, separate concrete-asset promotion approval and public discovery reconciliation. No new release framework or another preparation-only design slice is required. Public selection remains v1.6.0.
