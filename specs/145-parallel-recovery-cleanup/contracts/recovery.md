# S145 shared contracts

## Shared domain

Add JoinRecoveryCondition constants JoinRecoveryOnFailure/Canceled/Skipped/NotStarted; JoinSettlement constants JoinSettlementPending/Success/Failure/Canceled/Skipped/NotStarted; JoinSelectedRoute constants JoinRoutePending/Success/Recovery/None. Pending values are empty strings. Add DeliveryRoleJoinRecovery/Cleanup and JoinInstanceRecoveryPending/CleanupPending.

Definition fields: RecoveryTaskID/RecoveryTaskName, RecoveryOn []JoinRecoveryCondition with json recovery_on (no omitempty), CleanupTaskID/CleanupTaskName. Instance snapshots these configuration fields and adds SelectedRoute. Keep Target* success meanings. Add RecoveryDeliveryID/DeliveryState/RunID/Outcome/StartedAt/EndedAt/CancellationStatus/CancellationError and Cleanup* mirrors. Target/Recovery/Cleanup also expose Settlement, SettledAt and Reason. Invocation adds Settlement/SettledAt. All new task names/reasons are bounded in Observe; authoring excludes display/runtime fields.

## Store and engine

Existing exported store signatures remain unchanged. Internal advanceJoinInstance(tx,id,now) runs transactionally after owned actual terminal results and explicit pre-native no-start resolutions. Invalidation remains a separate path and a resolved instance never advances. Native StartFailed maps not_started; actual terminal cancellation_failed maps failure. Stable unique indexes protect a single success OR recovery follow-up plus one cleanup.

All new roles retain GetCompletionAdmission initial/queue/worker freshness, exact snapshot task/delivery binding and current eligibility. Extend activation/deletion/readiness and current-only fixtures. Engine completionJoinRole returns branch/target/recovery/cleanup. One buffered completionWake is consumed by the existing Start select; no recursive native-result drain, new goroutine or polling; drain until zero claimed.

## API, capabilities and complete authoring

Existing join endpoints, full Manage definition pages and Observe instance pages remain. JoinCreateRequest optional IDs use string omitempty; PATCH optional IDs use *string omitempty. Both use RecoveryOn *[]JoinRecoveryCondition omitempty to distinguish absent/default from explicit []. Empty supplied optional IDs clear fields. Clearing recovery with omitted RecoveryOn clears old causes. Strict input, 2 MiB body/MIME, authority, identity, audit, retry classification and revision CAS remain.

Advertise parallel-join-recovery in addition to parallel-joins. Existing compatible base workspaces and base create/PATCH packets omit new unconfigured fields. New options require the additional capability. Merge active observations into the matching exact stage only; persisted terminal evidence wins. State filters include recovery_pending and cleanup_pending.

## CLI, MCP and desktop

Extend join create/update with --recovery TASK_ID, --recovery-on failure,canceled,skipped,not_started and --cleanup TASK_ID. --recovery-on none means []; none mixed with other values is invalid. Update only supplied fields; --recovery= and --cleanup= clear. Human/JSON definition/instance output names selected/unselected routes and exact separate stages, reasons and settlement.

MCP joins_manage adds pointer-preserved fields without a new tool. Bounded allowlisted projections expose configuration, selected route and exact stage/settlement metadata, excluding commands/output.

JoinDraft adds recoveryTaskId/recoveryOn/cleanupTaskId. Full hydration and fresh-full sparse CAS preserve names, [] and references. New recovery defaults to [failure]; no boxes checked visibly means no matching causes. Four labeled cause controls and a vertical semantic routing preview explain actual settlement and cleanup order. Only new options are capability-gated; current target/revision/uncertainty/submission guards remain. Activity uses branch/target/recovery/cleanup-specific wording and retains exact-run cancellation.

## Portable and final integration

Portable v4 fields: recovery_task_portable_id, recovery_on and cleanup_task_portable_id. Optional pointer/custom serialization preserves [] and omits absent fields in base-only v3. Require parallel-join-recovery only when content needs it, and retain that requirement through preview/apply target revalidation. Remap/fingerprint/drift every new field; validate one final causal graph atomically and exclude runtime state.

Root owns internal/domain/joins.go, api/openapi/remote-v1.yaml and generated internal/remote/openapi.gen.go, CLI catalogue, docs and packaged/performance/final evidence. Store agent owns internal/store/internal/task. Engine agent owns internal/engine (except root-owned new route benchmark), API/server, remote operation regressions and internal/bundle. Interface agent owns client/CLI/MCP/desktop. Shared changes are coordinated before edits.
