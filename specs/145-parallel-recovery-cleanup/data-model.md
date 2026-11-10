# S145 data model

## Definition

Extend ParallelJoin with RecoveryTaskID, bounded observed RecoveryTaskName, RecoveryOn []JoinRecoveryCondition and CleanupTaskID/bounded CleanupTaskName. Causes are failure, canceled, skipped, not_started. Nil omitted policy defaults to failure when recovery exists; explicit [] matches none. Unknown causes and nonempty causes without recovery are invalid; repeated known causes normalize to a stable canonical set. Optional IDs can be cleared; references and combined causal graph are validated atomically.

Current schema 26 initializes directly. Definitions keep current revision/CAS and complete name fidelity. Definition changes affect future snapshots only.

## Instance and exact stages

Snapshot all optional configured references/conditions. SelectedRoute is empty pending, success, recovery or none. Existing Target* remains success-only. Recovery* and Cleanup* have independent DeliveryID, DeliveryState, RunID, Outcome, StartedAt, EndedAt, CancellationStatus and CancellationError. Each stage also exposes Settlement/SettledAt and bounded Reason projected from its delivery resolution.

ParallelJoinInvocation adds Settlement/SettledAt. Settlement is empty pending or success/failure/canceled/skipped/not_started. Actual StartFailed=true maps not_started; terminal cancellation_failed maps failure. Actual outcomes/IDs remain unchanged. Live failed cancellation does not settle.

## Durable lifecycle

Definitions without recovery/cleanup preserve established success-only states. Extended instances wait for every branch settlement before choosing a route. Success admits join_target, matching unsuccessful causes admit join_recovery, and no matching route records none. One partial unique index spans target OR recovery per instance. After the selected stage settles, join_cleanup is admitted once when configured; no selected stage advances directly to cleanup after branch settlement. No-start stage resolution advances without a fabricated native run. Cleanup completion resolves the final lifecycle while preserving upstream outcomes.

Additional lifecycle states recovery_pending and cleanup_pending express durable decisions, not native completion. Global resolved remains an invalidation tombstone. Late native results retain evidence and never revive invalidated automatic admissions. Pending/claimed/recovered work remains unsettled until a real terminal/no-start fact.

## Activation, graph and portable state

Current configured and outstanding snapshotted recovery/cleanup references participate in readiness, current eligibility and deletion handling. Replacing a definition cannot indirectly disable a task still required by its old instance. Task deletion removes invalid definitions and resolves future admissions according to current policy.

The causal task graph combines chains with join entry-to-all-destination edges. Finite same-task stage reuse is allowed; any potential chain/join cascade back to its entry remains rejected. Atomic replacement validates the final projected graph.

Portable v4 stores only configured references and conditions, including explicit []; remap every task reference and include all configuration in fingerprints/drift. Base-only joins remain portable v3. No settlement, selected route, deliveries, native IDs or runtime history is portable.
