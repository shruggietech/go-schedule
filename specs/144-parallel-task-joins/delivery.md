# S144 delivery checkpoint

## Source outcome

Implemented 144-parallel-task-joins on codex/144-parallel-task-joins, based on accepted main 28cd3ba83457384cc16bfbed3ed5f37241c22098. All 19 functional requirements, five success criteria, 28 implementation tasks and 22 checklist items are accounted for. Source is implemented; final owner review/merge and a future release remain distinct.

- #470: entry-owned named parallel invocations, exact all-success target admission, durable snapshots/current reopen, mixed graph validation and equivalent authoring/observation across CLI, local/remote API, controlled MCP, desktop Automation/Activity and portable v3 are implemented.
- #471: queued and capacity-delayed automatic deliveries reload current ownership/task/group/command eligibility at their real admission boundaries, with explicit resolution and preserved manual/overlap behavior.
- #472: ordinary and secret Automation mutations retain uncertain outcomes, current selected-daemon identity, in-flight guards and uncertainty across failed authoritative refreshes. The remote Automation workspace includes existing chain CRUD.

The official PR closes these issue outcomes on integration; checks support the functional delivery and are not independent issue-closure gates. No deferred cluster issue is resumed. Public v2.0.3 and historical release/comparison bytes remain unchanged.

## Verification and performance

The complete watched foreground eight-gate aggregate passed in 746.055 seconds with zero drift across 2,820 source paths and 11,608 protected historical inputs. Actual Windows and Linux/WSL packaged native workflows and portable export passed, alongside 184 frontend tests and 37 complete headless browser cases. [Verification](verification.md) records exact boundaries, original failed observations and the measured Principle IV exception: the 100-target fan-out path takes 27.164ms versus 23.346ms (+16.35%) after redundant SQL reads/preparation were removed; ordinary dispatch improves 4.91% and the unchanged nominal p99 limit passes.

## Publication checkpoint and stop boundary

The owner explicitly authorizes push, official non-draft PR creation, CI monitoring and in-scope findings resolution. The total Codex review cap is two including automatic review. At this source checkpoint no PR review has been requested; publication-time review counts, final candidate SHA, hosted check runs and any resolution follow-ups are retained privately and recorded on the official PR. Do not reproduce literal review triggers in evidence.

Stop with the official PR ready for the owner's final review and merge. Do not merge, tag, stage packages or publish a release. No new permission request is needed to perform the already authorized publication workflow.

## Concrete limits

Schema 25 follows the owner-directed current-only prelaunch policy: fresh/current schema 25 stores work, schema 24 is refused without conversion/deletion. Actual native Wails window interaction remains unverified. WSL host mirrored networking fell back to None, so native Linux execution does not establish a fresh cross-host network result. No operator desktop focus, service installation or host network mutation was used. External-command effects retain their existing at-least-once crash boundary.
