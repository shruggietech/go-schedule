---
title: Distributed scheduling decision
nav_order: 8.7
---

# Distributed scheduling: evidence and decision

go-schedule manages independent daemons. Remote control, All Systems, cross-daemon search and portable bundles do not currently provide shared execution ownership or automatic host failover. Clustered execution remains an open product proposal under [#20](https://github.com/shruggietech/go-schedule/issues/20).

The S125 assessment on 2026-10-28 found insufficient recorded real-use evidence to activate the architecture and prototype stages. This is a reversible hold on those stages, not a finding that clustered execution has no value or a decision to abandon it. No coordination backend or clustered execution semantics have been selected.

## Evidence register

| Source | What it establishes | What remains unknown |
|---|---|---|
| [Owner proposal on #20](https://github.com/shruggietech/go-schedule/issues/20#issuecomment-5073941594), 2026-07-24 | Product interest in node activity reporting, master/worker job claims, node reliability and optional result callbacks. The fleet-key suggestion is a design idea, not approved security architecture. | Actual workload, affected operator, machine/task counts, recovery target, failure tradeoffs and why available alternatives are insufficient. |
| [Connected control center completion](https://github.com/shruggietech/go-schedule/issues/174#issuecomment-5787231012) and [v1.5.0 publication](https://github.com/shruggietech/go-schedule/issues/185#issuecomment-5796659033), 2026-09-23 | Independent-daemon overview, target-specific search/actions and portable configuration shipped. | Adoption, operating experience and unmet shared-ownership requirements. Publication alone does not establish those findings. |
| [Demand issue #186](https://github.com/shruggietech/go-schedule/issues/186) and [decision issue #189](https://github.com/shruggietech/go-schedule/issues/189), inspected 2026-10-28 | The existing sequence requires actual scenarios, failure semantics and a go/no-go decision. Neither issue contained a real-use report at this assessment. | The workload and operator evidence required by #186. |
| Public issue search and relevant delivery discussions, inspected 2026-10-28 | Open/closed searches for cluster, distributed, failover and ownership, plus discussions on #20/#146/#174/#183/#184/#185/#282, returned roadmap, design and delivery material rather than a concrete post-release unmet workload. Repository Discussions were disabled. | Private use and unreported needs were not assessed. Search results cannot establish that there are no users or no demand. |

## Available approaches and their limits

| Need | Available approach | Limit relevant to the decision |
|---|---|---|
| Operate several machines from one interface | [Remote profiles and All Systems](remote-access.md#client-connection-profiles), cross-daemon search and deliberately selected target actions. | Each daemon remains authoritative for its own work; a successful action on one target does not create shared ownership across targets. |
| Reuse automation definitions on another machine | [Bundle export, validation, comparison, preview and apply](cli.md#bundle), with explicit target identity and target-local bindings. | Bundles move intent. They do not elect a leader, claim runs, synchronize runtime state or perform automatic failover. |
| Resume a local schedule after downtime | The task [catch-up policy](cli.md) can run one catch-up invocation or skip the missed window. | Recovery occurs on that daemon. Another machine does not acquire ownership while it is unavailable. |
| Obtain an established scheduler's failover capability | Review the project's [release-bound comparison](https://go-schedule.io/compare/) and the vendor's supported topology, edition and recovery behavior for the actual workload. | A documented failover scheduler is not proof that in-flight jobs, arbitrary side effects or every topology recover automatically. Product selection needs the same workload and failure-tradeoff inputs. |

The [prepared v1.7.0 comparison source](https://github.com/shruggietech/go-schedule/blob/b4f1e281e66003e08ee2a7abd35e3e25468fb6bd/site/data/comparisons/v1.7.0.json), reviewed 2026-10-27, records the official high-availability references for JAMS, ActiveBatch, Control-M, Azure Logic Apps and Resolve, and leaves unsupported conclusions about ScriptRunner and OpCon unknown. Reuse that bounded review when evaluating alternatives; it is not evidence that those systems fail to meet a particular user's needs. The prepared version is separate from the current public release.

## Supply a concrete workload

Add a report to [#186](https://github.com/shruggietech/go-schedule/issues/186), or provide a private report with a publishable summary if its operational details are sensitive. Do not include credentials, personal data, private addresses or secret-bearing command lines. Record unknown values explicitly instead of inventing estimates or treating the form as an eligibility test for using go-schedule.

Use this outline for each actual scenario:

1. **Source and operator:** report date, contact or attributable role, and whether this is observed use, a planned deployment or a design idea. Identify which details may be published.
2. **Job and value:** what the job does, who needs its result, the failure or limitation experienced, and a concrete example of the impact.
3. **Current deployment:** operating systems, host count, task count, schedule frequency, typical/maximum execution time and required network reachability. State unknowns.
4. **Ownership need:** whether work must move between machines, share a queue, survive loss of its current host or stay on one named machine. Identify local files, devices, accounts and other placement constraints.
5. **Recovery expectations:** acceptable start delay, outage duration and recovery time; whether running work can be interrupted, restarted or resumed; how an uncertain result is resolved.
6. **Side effects:** whether replay is idempotent, how duplicate execution is detected, and which cases tolerate duplicate, missed, delayed or uncertain runs. Explain any irreversible external effect.
7. **Failure cases:** required behavior for process crash, host loss, partition, clock disagreement and recovery of an old node. Distinguish observed failures from requirements not yet exercised.
8. **Alternatives tried:** independent remote daemons, portable bundles, local catch-up, external orchestration or another scheduler; what each accomplished and the specific remaining gap.
9. **Operating cost:** available operator skills, deployment/backup/upgrade expectations and acceptable dependency or support burden. Identify any unmet constraint instead of assigning an invented budget.
10. **Requested outcome:** the smallest useful capability and the evidence that would demonstrate its value. Attach sanitized incident or job references where available.

## Decision sequence

1. [#186](https://github.com/shruggietech/go-schedule/issues/186) assesses the supplied actual scenarios, bounds, alternatives and maintenance cost. It must explicitly recommend continued exploration before [#187](https://github.com/shruggietech/go-schedule/issues/187) proceeds.
2. #187 defines ownership, placement, partitions, clocks, duplicate/missed runs, recovery, uncertainty and single-node behavior. Maintainer review must authorize the prototype before technology selection.
3. [#188](https://github.com/shruggietech/go-schedule/issues/188) evaluates an isolated, unshipped prototype using established coordination and security components against the approved semantics. It compares a simpler alternative and records operational cost.
4. [#189](https://github.com/shruggietech/go-schedule/issues/189) records the findings-based go/no-go result for #20. It may document why the process stopped before the later stages. A go result creates a separate implementation roadmap; a stop result records alternatives and the disposition of unstarted work.

S125 records a hold because the first gate lacks concrete demand inputs. #186, #189 and #20 remain incomplete; #187 and #188 remain gated and unperformed. Reopen the assessment when an actual scenario supplies enough information to compare the ownership/recovery gap, scale, failure tradeoffs and operating cost. A new owner product decision can also revise the gate explicitly. No permanent product restriction follows from this assessment, and no release depends on inventing demand or repeating the same search.
