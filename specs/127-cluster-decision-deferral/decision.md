# S127 Owner-Directed Clustered-Execution Deferral

**Decision date**: 2026-10-31
**Authority**: Owner response to the S127 product-direction clarification: "Explicitly defer clustered execution."
**Disposition**: Explicit reversible deferral. No final go/no-go, research completion, architecture approval or clustered capability is claimed.

## What changed

[S125](../125-reliability-delivery-decision/decision.md) recorded an evidence-specific hold because actual workload inputs were missing. S127 records the owner's subsequent product direction to defer the clustered-execution effort. The lack of recorded demand was not proof of no demand; the owner's new direction is the authority for this deferral.

The shipped product continues to manage independent local and explicitly enrolled remote daemons, including All Systems, portable definitions, drift reporting, notifications and optional controlled MCP access. Shared run ownership and automatic cross-host failover remain unimplemented proposals. Deferral introduces no product-usage restriction.

## Issue-level disposition

| Issue | S127 disposition | Functional outcome remaining |
|---|---|---|
| [#186](https://github.com/shruggietech/go-schedule/issues/186) | Open, explicitly deferred | Actual attributable scenarios, bounded scale/failure tradeoffs and supported proceed/stop recommendation |
| [#187](https://github.com/shruggietech/go-schedule/issues/187) | Open, deferred and prerequisite-gated | Positive #186 recommendation, failure semantics and maintainer approval before prototype technology selection |
| [#188](https://github.com/shruggietech/go-schedule/issues/188) | Open, deferred and prerequisite-gated | Approved semantics, isolated unshipped prototype, alternatives and operating cost |
| [#189](https://github.com/shruggietech/go-schedule/issues/189) | Open, explicitly deferred | Findings-based final go/no-go or separately authorized stop |
| [#20](https://github.com/shruggietech/go-schedule/issues/20) | Open, explicitly deferred | Supported final direction and child dispositions |
| [#146](https://github.com/shruggietech/go-schedule/issues/146) | Open | Remaining #20 outcome; release #394 remains delivered |

No issue in this table closes on S127 merge. Included issues are traceability, not a claim of functional completion. No new label, milestone, assignee, priority or deadline is required.

## Reactivation

A later owner revisit or explicit scope revision is required before work resumes. A concrete actual or planned deployment report can inform that revisit through the [existing intake](../../docs/distributed-scheduling.md#supply-a-concrete-workload); describe its provenance honestly and assess it under #186's actual-use criteria rather than silently calling a plan observed use. Record known ownership/recovery gaps, host/task bounds, acceptable duplicate/missed/uncertain outcomes and alternatives. Unknowns remain unknown.

Unless the owner explicitly revises the gates, retain the sequence: #186 recommends continued exploration, #187 defines semantics and obtains maintainer approval, #188 evaluates an isolated unshipped prototype using established components, and #189 records the final decision. A go result needs its own implementation roadmap; a stop result needs an explicit disposition and alternatives. Neither is preselected here.

There is no automatic resumption, scheduled check, invented release date or repeated unchanged evidence campaign. Deferral does not permanently rule out clustered execution.

## Independent work

Public v1.7.0 delivery is complete through release 428427523 and owner-merged PR #406 at 6483fc81bbaa75ec0bddeb0fdf17fc4cc66be426. S127 does not mutate its tag, assets, reviewed comparison data or release notes. Original #355/#356 causes remain unresolved; deferred #400/#404 remain separate dependency work. No runtime, storage, API or workflow behavior changes.
