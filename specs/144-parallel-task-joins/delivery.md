# S144 delivery checkpoint

## Source outcome

Implemented 144-parallel-task-joins on review branch codex/144-parallel-task-joins in PR #473, based on accepted main 28cd3ba83457384cc16bfbed3ed5f37241c22098. Current source accounts for all 19 functional requirements, five success criteria, 32 completed implementation tasks and 22 checklist items. The final corrected complete local eight-gate run passed. Exact-head hosted CI follows the final push; owner final review and merge remain pending.

- #470: entry-owned named parallel invocations, exact all-success target admission, durable snapshots/current reopen, mixed graph validation and equivalent authoring/observation across CLI, local/remote API, controlled MCP, desktop Automation/Activity and portable v3.
- #471: queued and capacity-delayed automatic deliveries reload current ownership/task/group/command eligibility at their actual admission boundaries, with explicit resolution and preserved manual/overlap behavior.
- #472: ordinary and secret Automation mutations retain uncertain outcomes, current selected-daemon identity, in-flight guards and uncertainty across failed authoritative refreshes. Remote Automation includes existing chain CRUD.
- #474: complete branch replacements discard readonly observation metadata, allowing removal, addition and reordering while preserving slot IDs and revision conflicts.
- #475: selected-target capabilities control initial and post-mutation join reads, preserving compatible Automation workspaces and visible advertised failures.
- #476: full writable configuration pages, audited Manage reads, complete editor hydration and sparse revision-checked saves preserve long Unicode names and exact branch identity. Bounded observations remain distinct; join-only strict JSON budgets match the existing 2 MiB configuration budget.

The official PR closes these issue outcomes on integration. Checks support functional delivery and are not independent issue-closure gates. Deferred cluster work remains deferred, and public v2.0.3 and historical release/comparison bytes are unchanged.

## Verification and performance

The final corrected complete watched foreground eight-gate aggregate passed in 910.891 seconds with zero drift across 2,829 source files and 14,758 protected historical/prior-candidate inputs. Root packaged integration passed in 275.971 seconds; desktop race/build, TypeScript and all 187 frontend tests passed. Full headless Chromium coverage passed 38 cases. Actual owned Windows and Linux/WSL packaged workflows cover correlated entries, repeated-task slots, cancellation and portable export. Final Linux Store/API race and packaged workflow receipts are separately retained. [Verification](verification.md) records exact boundaries and original failed observations.

Final-correction paired measurements preserve accepted-main inputs and thresholds: ordinary dispatch median improves 28.35%; the full 100-target completion path is 25.001ms versus 22.015ms (+13.57%, +2.986ms). The recorded Principle IV exception covers fresh automatic-admission checks after redundant reads and SQL preparation were removed. The unchanged 2000-sample nominal p99 dispatch limit remains 100ms and passes.

## Publication checkpoint and stop boundary

The official non-draft PR is #473. Both permitted Codex reviews are consumed: round one reviewed 3d4ce0227108d3b505e303411797d91641965bf6 and found #474/#475; round two reviewed a52f34b83dc8c8b4610f5742e30af77d43c9bce7 and found #476. All findings have corrections. The final #476 delta is qualified with focused/full local tests and fresh exact-head CI; it has no third Codex review. Publication-time review receipts, final SHA, terminal hosted checks and thread state are retained privately and recorded on the PR. Literal review triggers are omitted from evidence.

Stop with the PR ready for owner final review and merge. Do not merge, tag, stage packages or publish a release.

## Concrete limits

Schema 25 follows the owner-directed current-only prelaunch policy: fresh/current schema 25 stores work; schema 24 is refused without conversion or deletion. Actual native Wails window interaction is unverified. WSL host mirrored networking fell back to None, so actual Linux execution does not establish a fresh cross-host network result. No operator focus, service installation or host network mutation was used. External-command effects retain their existing at-least-once crash boundary.
