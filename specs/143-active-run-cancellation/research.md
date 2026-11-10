# S143 Research Decisions

## Process ownership


## Decision

Keep process ownership in a small production executor adapter rather than using the QA launcher. The fixture contains validated platform techniques but also file-output, absolute-path and optional resource-limit policy that must not enter ordinary task execution. No tracked source was edited by this research.

## Recommended shape

- Common executor code builds the existing command semantics, retains capBuffer and uses explicit owned input/output pipe endpoints. Platform start returns one native root wait capability, one owned tree-stop capability and ordinary handle release. One owner waits/reaps the root; cancellation feedback may be surfaced before that root actually exits.
- Unix keeps exec.Cmd and merges Setpgid=true into its existing SysProcAttr after run-as preparation, preserving Credential. Each execution has its own group, including concurrent executions of one task. Reuse internal/processgroup termination/query behavior under the remaining common cleanup deadline. Do not select processes by name or rediscover unrelated host processes.
- Windows uses native CreateProcess with STARTUPINFOEX and only explicitly inherited standard handles. Reuse the existing fixture's CREATE_SUSPENDED -> AssignProcessToJobObject -> ResumeThread creation sequence. It retains exact process/thread/job handles and is already exercised in this repository. Atomic PROC_THREAD_ATTRIBUTE_JOB_LIST is also officially supported on Windows 10+/Server 2016+ and removes assignment/resume stages, but changing the proven local pattern is optional.
- Windows launches retain CREATE_NO_WINDOW, SW_HIDE and Unicode environment/argument handling. No memory/CPU/process-count/priority limits or breakaway restrictions beyond ordinary job membership are added. Per-task Windows run-as remains rejected exactly as before.
- Go1.26 SysProcAttr has no JobList attribute and its StartProcess closes the primary thread handle. Using exec.Cmd suspended start plus Toolhelp/OpenThread therefore adds thread enumeration, PID reopen and identity checking. The small native adapter is preferable to that workaround.
- Preserve command resolution and working-directory semantics from exec.Cmd, including relative executable paths under Dir, PATH resolution, .exe extension handling, empty/quoted/Unicode arguments, task environment overrides, stdin and output caps. This is a required adapter regression surface, not permission to simplify inputs.

## Normal completion and explicit cancellation

Do not set a global WaitDelay, because it also caps ordinary inherited-output draining after root exit. Do not let ordinary handle release kill background descendants: omit JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE or explicitly clear it before normal release. A natural root completion and intentional background descendants must retain the current ordinary completion behavior.

Only an explicit run cancellation starts the single original five-second budget. Terminate the exact job/group, concurrently observe the native root result, and drain owned I/O. Windows success requires ActiveProcesses==0; Unix group cleanup requires ESRCH, with permission ambiguity treated as unconfirmed. Native root exit information remains separate from cancellation outcome. All phase allowances consume the original deadline rather than stacking two five-second waits.

If termination/collection cannot be confirmed by that deadline, publish cancellation_failed feedback and close the owned parent pipe endpoints to end copy goroutines, but retain the native reaper and exact run ownership until the actual process exits. A kernel-uninterruptible process cannot honestly be declared terminated or have its worker capacity freed merely because a timer elapsed. The cancellation acknowledgement remains a request, not proof of termination.

A natural result already observed wins over a late request. A root that exited while a descendant retains the output pipe is not fully collected; explicit cancellation may end owned descendants/drain collection while retaining the genuine root exit code. The engine must synchronize request/terminal retirement so repeated stale Run IDs cannot address a later process.

## Files

Production: internal/executor/executor.go, new internal/executor/process.go, process_windows.go and process_unix.go, with a small shared cleanup-deadline/result contract; internal/processgroup/terminate_unix.go only if extracting an existing wait/query primitive is needed. Keep internal/executor/runas_unix.go semantics unchanged and preserve platform.HideConsole behavior. Do not make executor import internal/clifixture.

Tests: internal/executor/cancellation_test.go and platform-specific process tests, plus existing executor_test.go/runas_unix_test.go. Engine tests exercise the resulting exact-run handles separately.

## Meaningful regressions

1. Hidden helper launches a child that retains output; cancel the owning execution, assert root/child exit, partial output and bounded collection, and keep a sibling execution alive.
2. Repeat with two concurrent executions of the same task to prove Run ID isolation.
3. Hold the process at a startup handshake and cancel before launch; no payload side effect and no native exit code is invented.
4. Release a natural completion before cancellation and assert its original success or failure remains. Repeat competing completion/cancel transitions under race.
5. Inject termination failure, nonempty ownership query and pipe drain deadline; assert explicit failure and retained live ownership rather than a fake canceled terminal result.
6. A naturally finished root creates a deliberately background child with detached output. Assert ordinary completion does not kill that child; clean the owned test fixture separately.
7. Preserve existing literal args, stdin, environment overrides, output cap and run-as tests; add relative executable/working-directory and Windows environment/PATHEXT cases for any native adapter.
8. Verify no inheritable handles leak to the child except selected standard streams; setup failures close all created handles and never resume an uncontained suspended child.

## Limits

OS process groups cannot contain deliberately escaped Unix sessions, and termination cannot reverse external side effects. Windows parent/nested-job constraints may cause ownership setup to fail; report a start/ownership failure without running uncontained work. These are explicit platform observations, not new product eligibility checks.

## Primary references

- Go os/exec cancellation, pipes and WaitDelay: https://pkg.go.dev/os/exec
- Microsoft process creation flags and CREATE_SUSPENDED: https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
- Microsoft suspended thread execution/ResumeThread: https://learn.microsoft.com/en-us/windows/win32/procthread/suspending-thread-execution
- Microsoft process attribute list, including JOB_LIST: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute
- Local primary Go1.26 source: src/syscall/exec_windows.go and src/os/exec_windows.go.


## Engine and API


Read-only research against the S143 base (main 3d772b99c7684062c80cebb410d4dcf3d0d489f2). No tracked implementation changes.

## Recommended engine contract

Add per-Run-ID cancellation controls alongside the current active map. Construct context.WithCancelCause for each acquired-worker execution, register the handle before onRunStarted and Runner.Run, and preserve existing UUID/correlation binding. Use a shared explicit-cancellation cause sentinel in a low-level domain/control boundary, not an engine import from executor. A request context must authorize/admit the mutation but must not own the continuing execution context.

CancelRun(ctx, runID) should return domain.RunCancellation with run_id, task_id, status requested or already_finished, and outcome when finished. A requested response acknowledges initiation only. Mark a control finishing and freeze the returned Runner result under engine.mu before persistence, so a cancel during database work cannot change an already observed terminal result. Keep the handle through persistence. Release registry state and context resources exactly once. Actual executor observation determines canceled versus natural result; engine must not overwrite a natural result solely because a cancellation context became done.

Use a context-aware metadata-only store lookup for a completed/missing fallback: id, task_id, outcome, started_at, ended_at are enough. Current Store.GetRun uses QueryRow without the request context and also loads output; introducing that path into cancellation would repeat the known SQL-pool cancellation defect. QueryRowContext or a focused GetRunResultContext avoids both issues.

A failed five-second cleanup attempt is not proof that execution ended. If the OS process remains unreaped, keep the worker and active ownership, expose an explicit transient cancellation failure/unconfirmed state and continue the owned reaper lifecycle. A cancellation_failed terminal value is appropriate only with an actual terminal observation, not as a substitute for a live process.

## Persistence and downstream effects

The current runs.outcome column is unrestricted TEXT, and scans return RunOutcome strings directly. New canceled and cancellation_failed terminal enums do not require schema shape/version changes. Run.CancellationRequested or a compact cancellation-status field may be transient ActiveRuns presentation metadata without schema persistence, provided terminal outcome remains durable.

Store.RecordRunAndCreateDeliveries (chains.go) already dispatches downstream chains only for success/failure, including the any selector. Distinct cancellation outcomes naturally create no downstream work, while incoming completion deliveries can still be resolved to the actual target run. Add explicit regressions for all three chain selectors and cancellation_failed so future broadening cannot accidentally fan out work.

Outbound webhook/SMTP assignments have success/failure, start-failure, duration, recovery and streak/reminder conditions. They have no OnCanceled or generic completion subscription. Keep deliberate cancellation outside those conditions, preserving failure streak/recovery state. The current engine onRun -> broker.PublishRun -> run_finished stream can carry the actual canceled outcome and update clients truthfully. Adding a new configurable outbound cancellation subscription would be separate schema/UI scope.

## API and transport files

- internal/api/server/runs.go: new handler/source contract, explicit response classification. Prefer 202 requested, 200 already_finished, 404 unknown, 409 existing non-started/nonactive queued marker. Do not return captured task output in the mutation acknowledgement.
- internal/api/server/server.go: POST /v1/runs/{id}/cancel route.
- internal/authorization/catalog.go: runs.cancel, CapabilityOperate, AuditMutation, TargetKind run. Existing authorization/audit wrapper provides actor validation, expected daemon identity, denied/succeeded/failed audit and independent completion after mutation.
- internal/remote/operations.go: /api/v1/runs/{id}/cancel -> local path, operation runs.cancel, RetryNever. Regenerate internal/remote/openapi.gen.go from api/openapi/remote-v1.yaml using the pinned generator rather than hand editing.
- internal/api/server/manifest.go: add run-cancellation capability for UI/discovery as needed; avoid advertising actual support if the injected scheduler lacks the cancellation source.
- internal/api/client/methods.go and observation.go: CancelRun method and cancellation route token normalization. The latter currently has a fixed path-token allowlist.
- API/remote tests should cover local and enrolled-remote targets, observer denial, revoked credentials/sessions, expected-daemon mismatch, mutation audit failure before request, post-mutation audit completion, canceled request pool wait and identical response binding.

## Consumer mappings

- desktop/operations/service.go normalizeOutcome currently maps unknown run outcomes to unavailable; add both cancellation outcomes.
- desktop/frontend/src/operations/ActivityPage.tsx filter list and labels; SchedulePage.tsx stateLabel; inspect shared state badges/tone choices.
- desktop/popups/service.go only triggers success/failure conditions, so deliberate cancellation should not become a false success/failure popup.
- internal/store/system_summary.go and search.go currently count/select failure only. Decide deliberately whether cancellation_failed is an actionable failure-type observation; canceled should not inflate failure counts. If included, preserve its actual distinct outcome in results.
- Internal completion and bundle chain selectors stay success/failure/any; adding cancellation to RunOutcome does not add a chain selector.
- internal/releasegate/validate.go has a historical eight-state qualification set. Preserve historical packets/fixtures; S143-specific tests can assess new states without rewriting old source evidence or claiming new native observations.

## Confirmed pre-existing concurrency bug

The current per-task running map is bool. With two allow_concurrent executions active, finish of either sets it false and immediately releases a queued run even if another original is active. Use an in-flight count (including reserved/waiting-worker executions) so queued work is released only when the last execution finishes. This is needed to preserve exact-run cancellation and policy changes while other instances remain active.

A private source-preserving Go overlay reproduces the defect: start two same-task allow_concurrent runs with existing blockingRunner, change the incoming task policy to queue_one, queue a third and release only one original. The queued third starts prematurely. The fixture drains all goroutines and changes no source file.

Command: go test -race -overlay .test-temp/s143/concurrent-overlay.json -run '^TestS143Overlay_QueuedRunWaitsForAllConcurrentExecutions$' -count=1 ./internal/engine

Observed exit 1 in 0.113 seconds: queued run started while another concurrent execution was still active. Private files: concurrent-overlay-test.go, concurrent-overlay.json, concurrent-overlay-red.log. File a separate source bug linked to #465; the same engine correction can be bundled in S143.


## Control interfaces


## Recommended shared contracts

Add a transport-neutral cancellation acknowledgement (Run ID, task ID, status requested/already_finished, terminal outcome when already finished) separate from Run. A202 response means only requested, not canceled. Existing API client do accepts all2xx and already wraps lost remote mutation responses as MutationUncertainError; reuse it. Use url.PathEscape for Run ID. CLI runs cancel RUN_ID uses the existing reqCtx, clientFor and global JSON flag; print canonical acknowledgement JSON or a checked human write explaining that history/active observation establishes final termination. Do not replay after output-write failure.

## File map

- internal/cli/runs.go: add child command to newRunsCmd; runs_test.go plus operational/JSON command fixtures for accepted, already-finished, missing, denied, malformed arity and output-write failure. Existing CLI generic exit mapping handles validation2, ordinary errors1 and JSON diagnostics.
- internal/api/client/methods.go: CancelRun through POST /v1/runs/{id}/cancel and response decode. runs_test.go, remote_test.go, mcp_sessions_test.go. Observe path sanitizer in observation.go needs cancel literal allowlist so operation diagnostics do not expose path segments or collapse incorrectly.
- internal/authorization/catalog.go: runs.cancel, Operate, run target, AuditMutation. internal/remote/operations.go: matching /api/v1/runs/{id}/cancel, RetryNever, inherited catalog authority/audit. api/openapi/remote-v1.yaml needs documented route/ack schemas and regenerated remote bindings using existing go generate target.
- internal/mcpoperate/operate.go: dedicated public CancelInput daemon_id/run_id/request_id and runs_cancel tool; shared existing session request-ID cache/fingerprint must cover cancellation together with task operations (do not create an independent cache allowing conflicting request-ID reuse across operations). Add run_id/status to safe result, not task command/output. Recheck revocation on cached result using existing VerifyAccess. Existing operationClient fakes in mcpoperate, mcphttp and remotemcp need new method if interface extended. Update exact3-tool discovery assertions to4 in mcpoperate/operate_test.go, mcphttp/manager_test.go and remotemcp/handler_test.go. Manage may include Operate tools; inspect assertions as well.
- desktop/operations/local.go: Backend/LocalBackend cancel forwarding. service.go: cancellation request result, safe status classification and ActivityRecord active/persisted handoff; model.go: cancellation acknowledgement data. service_test.go contains stubBackend and handoff coverage.
- desktop/app.go: exported CancelRun facade, preserve a.reconcileUncertain(result.Outcome). desktop/app_test.go fake facade backend grows accordingly. Wails frontend bindings are ignored, not tracked; production build generates them.
- desktop/frontend/src/operations/model.ts, bridge.ts, store.ts, ActivityPage.tsx: expose cancelRun, button in selected active Run details, pending-by-run acknowledgement distinct from terminal state, canceled/cleanup-failed filter/state/severity display. App.tsx must pass selected daemon identity along with canOperate, not merely target display name. styles.css state-running/queued/warning currently use dashed border.
- desktop/frontend/src/operations/ActivityPage.test.tsx, store.test.ts, desktop/frontend/e2e/operations.spec.ts: same-task isolation, exact RunID wiring, keyboard disabled observer/stale states, duplicate click guard, no success claim from acknowledgement, eventual terminal result, target-switch suppression and retained-workspace uncertain outcome. Use headless only.
- docs/cli.md, docs/api.md, docs/mcp.md and docs/gui-fields.md: command, endpoint, fourth Operate tool, identity/ack/final-result meaning and cancellation chain/notification semantics.

## UI state and target safety

Current useActivity mutating guard uses React pending state, which does not prevent two submissions in one render turn; use a useRef lock like notifications/store.ts. Defer background refresh while mutation pending, then refresh, and reject stale sequence/generation responses. Preserve workspace on errors. Exact-record mode currently disables subscription and loads exact record only on mount/selection; Refresh incorrectly calls recent ActivityWorkspace. For selected run cancellation, explicitly reconcile exact identity instead of relying only on events. Snapshot ListActiveRuns before GetRun (GET exact is persisted-only), use persisted completion if found, else use earlier active snapshot only on a genuine not_found. Poll cancellation every250ms for up to8s, with timer cleanup on unmount/target change; on expiry show unconfirmed termination and Refresh, not canceled.

Capture daemon identity from the rendered connection snapshot with the cancel request and bind with client.ExpectDaemon through the Wails/API path. That prevents a rapid target switch from dispatching an old Run ID to a new daemon, including clones with copied run IDs. Backend facade should also require connected/current manager snapshot and existing Operate permission for the selected target. Uncertain remote mutation calls a.reconcileUncertain and keeps cancellation disabled until authoritative refresh.

## Confirmed integration blocker: product2 rejected

Both local.go and remote.go connection Health use compatibleVersion that accepts only product majors0/1. Existing local_test.go intentionally labels2.0.0 incompatible, so publishedv2.x desktops cannot negotiate their own daemon. Recommended proportional correction is manifest API-version negotiation: local requires exact v1 in LocalAPIVersions; remote exact v1 in RemoteAPIVersions (currently only nonempty). Preserve health status, nonempty consistent product-version metadata, installation/display identity, remote VerifyIdentity and VerifyAccess. Testv2.0.3 + v1 success, missing/unsupported API rejection, health/manifest version mismatch and current identity/access failure cases. Future product majors with supported APIv1 can be compatible; future unsupported API-only manifests remain rejected. Source is direct demonstration; parent will reproduce with private overlay before filing separate bug.

## Analysis correction

The initial analysis inventory exposed missing explicit FR-013/014 entries for the two already-scoped prerequisite fixes, although their tasks were present. Added those requirements and repeated the blocking cross-artifact review before source implementation. The corrected inventory is fourteen unique functional requirements, five success criteria and twenty-six tasks, with all mapped and no remaining ambiguity, duplication, unmapped task or constitution conflict.
