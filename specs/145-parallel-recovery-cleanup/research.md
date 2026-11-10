# S145 research

## Settlement evidence

Decision: persist explicit success/failure/canceled/skipped/not_started settlement and time per exact branch/stage. Native results own settlement only when the runner returns; explicit pre-native ResolveCompletionDelivery writes not_started. Direct invalidation resolves the instance without pretending that native work ended.

Rationale: an early blocked route, queued marker, failed cancellation request or resolved invalidated delivery does not establish termination. Actual StartFailed=true is not_started while retaining its diagnostic RunID/Outcome; terminal cancellation_failed maps failure after the native owner returns.

Alternatives: inferring completion from delivery_state or adding a start-write to every native dispatch. The first loses ownership truth; the second is unnecessary because no-start and invalidation paths are already separate.

## Finite stages and causal graph

Decision: use success/recovery/none selection plus separate success Target, Recovery and Cleanup evidence. Extend current schema to 26 with direct initialization/refusal, exact immutable snapshots and unique follow-up/cleanup admissions.

Rationale: success and recovery must never both admit for one instance. Cleanup remains independent from upstream success and cannot overwrite outcomes. Existing no-option joins retain their established behavior.

Graph correction: ordinary chain source-to-target edges plus join entry-to-every-owned-destination edges represent potential automatic cascades. Intra-instance finite stage ordering does not create new join instances, so modeling it as global task edges falsely rejects task reuse. Any stage may independently trigger chains/new joins, so the entry-to-all graph still exposes true cascaded cycles. Reuse the same validator for atomic replacement and portable planning.

## Event-driven progress

Decision: one buffered completionWake channel is consumed in the existing engine select. Successful result/no-start commits coalesce a wake; drain stops only at zero claimed deliveries, with an initial startup drain.

Rationale: a no-start resolution can itself create recovery/cleanup with no future native result to wake the queue. A short-batch early return can strand such work. This also removes recursive record/skip/drain nesting.

Alternatives: recursive drains, polling or another resident goroutine. They add stack/lifecycle/latency costs without improving this already event-driven architecture.

## Writable intent and interface compatibility

Decision: RecoveryTaskID, RecoveryOn and CleanupTaskID are optional additive configuration. Domain RecoveryOn arrays preserve explicit empty values; request create/PATCH use pointer arrays and optional IDs. Omitted recovery conditions default to failure when recovery is configured; [] matches no cause. Clearing recovery without a condition patch clears its old policy.

Rationale: omitempty on a plain writable array destroys [] intent, while unconfigured null fields sent to older strict servers break compatible base joins. Complete Manage authoring and sparse revision/target-fenced edits remain authoritative. Observe stays bounded.

Add parallel-join-recovery alongside parallel-joins; gate new options only. Portable v4 carries new options and their full remapped references; base-only joins export v3 with absent fields omitted. Plan/apply compatibility retains required new capabilities even after items become generic join records.

## Research execution and constitution

The installed plan command dispatched three read-only research agents for store/task, engine/API/portable and interfaces. They inspected actual source, coordinated exact contracts and retained private reports under .test-temp/s145. No implementation preceded the blocking analysis gate.

Post-design constitution check passes I-V: established libraries only, no new service, owned lifecycles and clock discipline retained, failing-first behavior tests and all mandatory gates planned, measured scheduling budgets unchanged, no historical migrations, no scope substitution. The owner kickoff supplies publication authorization; merge/release remain outside scope.

Known duplicate recovery causes normalize to one stable ordered set, consistent with Store validation, rather than adding a needless rejection for repeated equivalent intent. Unknown causes remain invalid. Valid Store-backed unconfigured responses always expose recovery_on as [], so the response schema is a required array.
