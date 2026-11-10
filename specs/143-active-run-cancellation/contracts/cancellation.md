# S143 control contracts

## Shared Go values

domain.RunCancellation has RunID, TaskID, Status, Outcome, CancellationStatus and CancellationError fields, with snake_case JSON and optional result/status fields. Engine implements CancelRun(context.Context, string) (domain.RunCancellation, error). API client exposes the same result. Domain Run gains transient CancellationStatus and CancellationError. Shared internal/runcontrol defines ErrRequested, ErrNotActive, CleanupTimeout=5s, WithFailureReporter(ctx, func(error)) and ReportFailure(ctx,error).

## Local and remote JSON

POST /v1/runs/{id}/cancel and its existing /api/v1 remote projection accept no command/PID override. 202 returns requested; 200 returns already_finished with actual outcome; 404 identifies an unknown Run ID; 409 identifies known nonactive queued/skipped work; existing authentication/identity/authority failures retain their established errors. Operation runs.cancel requires Operate, audits run target and uses remote RetryNever. Store fallback reads bounded result metadata with request context.

## CLI and MCP

CLI: gosched runs cancel RUN_ID supports checked human output and global JSON, normal stderr/exit conventions and explicit requested versus already-finished wording. Existing gosched runs active provides ID discovery through GET /v1/runs/active. MCP goschedule://runs/active exposes paged exact Run IDs and authored cancellation metadata without raw output/execution inputs; active-read failures never substitute history. MCP runs_cancel takes daemon_id, run_id and request_id and shares current bounded deduplication, actor/revocation checks and mutation-uncertainty rules.

## Desktop

App.CancelRun(runID, expectedDaemonID) uses the operations service/client and existing reconcileUncertain behavior. Captured identity is the actual selected daemon installation ID, not merely a profile alias. Activity/Run details offers a keyboard-accessible Cancel run action only for its authoritative active selection with operating authority. A ref-backed in-flight lock prevents duplicate same-render submissions. The UI refreshes the selected exact record, preferring persisted completion over an earlier active snapshot. After acknowledgement, 250ms polling is bounded to eight seconds and is canceled on unmount or target-generation change; timeout leaves termination unconfirmed and offers Refresh.

## Completion and platform semantics

Canceled/cancellation_failed never dispatch completion chains or affect outbound success/failure/duration/recovery policy streaks. Existing run_finished events carry their actual outcome. Successful explicit cancellation confirms owned job/group emptiness, native reap and collection. Cleanup failure is observable within five seconds, but unreaped work remains active and owns its capacity until actual completion. Ordinary task completion retains background-child, args/env/stdin/run-as and output-cap behavior. Explicit Windows breakaway, deliberately escaped Unix sessions and already-performed external effects are outside OS-enforced cancellation.
