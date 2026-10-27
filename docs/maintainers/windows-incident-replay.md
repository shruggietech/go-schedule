# Bounded Windows incident replay

Use this recipe to investigate #355 and #356 with the existing owned fixtures and collectors. It is a targeted diagnostic campaign; the complete unattended CLI evaluation and its seven-section report remain the integration workflow after application changes.

## Prerequisites and identity

Use the existing Windows user, Go/C toolchains, PowerShell, Bash and SQLite. Use the retained hash-verified portable PowerShell directory for the older-runtime condition without changing the installed runtime. Use a hidden noninteractive launcher with redirected I/O and CREATE_NO_WINDOW for every console process. Keep all runtime/service/session operations confined to the existing disposable fixtures. Private native dumps and stacks stay local.

Freeze the source revision/content digest, dirty state and original evidence hashes. Before the campaign clock starts, prepare six binaries with the existing coordinator and compile the two race-enabled test executables. Supply an existing Go executable matching go.mod; do not silently install another toolchain. Use new ignored input/evidence directories.

```text
Go run ./cmd/gosched-qa build --source SOURCE --bin-dir NEW_INPUTS --go EXISTING_GO
Go test -race -c -o PRIVATE_INTEGRATION_TEST ./test/integration
Go test -race -c -o PRIVATE_IPC_TEST ./internal/clievaluate
```

In these command templates, Go names the selected existing Go executable; replace capitalized paths with the selected absolute paths. Hash both compiled test executables and preserve build-receipt.json. Direct compiled tests require the same working directory that go test normally supplies: SOURCE/test/integration for script tests and SOURCE/internal/clievaluate for IPC tests. Record any setup failure separately and retain its original receipt.

## One shared campaign deadline

Set one UTC deadline thirty minutes after campaign execution starts. Run condition groups sequentially. For each compiled test invocation, pass -test.v=true, -test.failfast=true, -test.count=COUNT and -test.timeout=REMAINING, where REMAINING is recalculated from that original deadline, not another thirty-minute allowance. The hidden parent also bounds its wait by the remaining deadline. Stop after an unexpected failure for analysis. Record omitted conditions as incomplete; do not turn failed writers into successful repetitions. The fixture's six writers remain concurrent.

Set GOSCHEDULE_PROBE_DIAGNOSTICS=1 and GOSCHEDULE_INCIDENT_DIR to a new canonical private directory per runtime group. For the retained PowerShell group only, prepend its verified portable directory to the test process PATH. Restore the original PATH for subsequent conditions. Do not change global environment or installed files.

| Group | Executable and -test.run filter | Count |
| --- | --- | ---: |
| Current PowerShell | Integration test: ^(TestScriptsBindingsKeepConcurrentSessionIdentities\|TestScriptsRuntimeAndSelectionContrasts)$/^powershell$ | 10 |
| Retained PowerShell | Same filter with retained portable child PATH | 10 |
| Current Bash | Integration test: ^(TestScriptsBindingsKeepConcurrentSessionIdentities\|TestScriptsRuntimeAndSelectionContrasts)$/^posix$ | 10 |
| Deliberate prerequisite failures | Integration test: ^(TestScriptsSQLiteCandidateFailureDiagnostics\|TestScriptsSQLiteLauncherFailureDiagnostics\|TestScriptsPowerShellProbeBoundsInheritedPipes)$ | 1 |
| Memory/WAL component replay | IPC test: ^TestWindowsIncidentReplay$/^component_(memory\|file_wal)$ | 10 |
| Packaged-daemon WAL replay | IPC test: ^TestWindowsIncidentReplay$/^packaged_file_wal$ | 10 |

The table escapes regex alternation pipes for Markdown. Pass actual pipes to the executable. Runtime repetitions retain both fresh/existing schemas and metadata/writer implicit/explicit selection of the same SQLite candidate. Deliberate prerequisite timeouts run once rather than multiplying them by ten.

For IPC groups, set GOSCHEDULE_INCIDENT_REPLAY=1 and GOSCHEDULE_REPLAY_BIN to the prepared input directory. The test validates all six binary hashes against the strict receipt before selecting its native binaries. The original five targets are service.restart, mcp.http.enable, mcp.http.rotate, mcp.bad-port and group.remove-child. Unchanged catalog group.create/group.child establish prerequisites; enable precedes rotation. The service target records adapter intent. Expected bad-port exit 2, state assertions, MCP resource/confirmation/retry checks and previous-credential rejection remain active.

Component modes use actual API/store/MCP implementations over a real owned native endpoint. They report the actual memory/WAL store separately from logical owned runtime paths and identify the test host process. They are not a packaged-daemon memory mode. Packaged mode uses the actual file-backed daemon. Seven steps per fresh session stay within the existing eight-snapshot allowance without raising it.

## Observation and cleanup

Opt-in script logs emit incident_observation JSON, including test/condition/writer, native exit, combined result error and independent collection/cleanup errors, bounded output, exact executable/harness provenance and incident summary. IPC logs emit incident_replay JSON, including source receipt, actual storage selection, unchanged case/check outcomes and private evidence references. Expected negative cases remain expected. An absent incident on a fast successful case means the watchdog did not need to trigger, not a failed case or proof that an earlier stall was impossible.

Treat command, protocol, inspection, collection and cleanup independently. Capture taken after cancellation is not live evidence of an earlier wait. Original controlled failures and deliberately invalid prerequisites are not spontaneous recurrences. A controlled FailFast only demonstrates the collector; analyze a matching actual fault before assigning a project or upstream cause.

MCP cleanup runs while the native API is available for revocation, followed by incident cancellation/join, API shutdown/join, store close and owned runtime removal. Never terminate a synthetic process handle containing the test host PID. Keep retained evidence separately from disposable runtime state and list every remaining resource. Later independent verification is recorded separately from the original campaign clock.

## Fixed investigation report

Use these six sections in this order. Preserve the canonical full-QA JSON/Markdown separately and do not edit its generated report.

1. Source and limits: baseline/current source digest, dirty state, exact build receipt and test binary hashes, interpreter/module identities, selected tools, shared start/deadline, attempts, preparation versus campaign versus independent verification timing.
2. Planned and observed matrix: one row per condition/repetition with track, baseline, actual storage/execution kind, expected/actual outcome, original exit/errors, completed/failed/unavailable/interrupted state and exact incident/log references/hashes. Count all planned and observed rows; omitted rows remain explicit.
3. Runtime/prerequisite outcomes: distinguish interpreter launch, native runtime, prerequisite discovery and writer behavior; retain six simultaneous writers, literal identities and exact records.
4. IPC/inspection outcomes: retain original command, protocol, inspection, collection, partial-state and cleanup outcomes; distinguish handler progress, SQL-pool wait and transport observations.
5. Causal findings and corrections: hypothesis, matching supporting and contradictory observations, supported conclusion or uncertainty, separately tracked defect/fix or evidence-matched external disposition, affected exact source and preserved regression.
6. Cleanup and unresolved outcomes: owned process/runtime inventory, completed removal, retained private artifacts, remaining resources and issue-level disposition. Passing repeats and unrelated fixes do not close historical causes.

The optional machine-readable sidecar has kind go-schedule.incident-replay-campaign and version 1, with baseline/buildReceipt/testBinaries, started/deadline/finished, conditions, source_drift, originals_changed and overall status. Each condition identifies planned_repetitions, pattern, status, timestamps, original native exit/launch error, log/stderr identities, exact observations and cleanup counts. It supplements the full-QA schema; it does not replace or relax it.

## Supported cancellation and launch outcomes

MCP lifecycle admission honors cancellation while another lifecycle operation remains active. Cancellation observed before commit prevents the mutation and releases uncommitted listener/session preparation. Once listener revocation commits, bounded cleanup completes independently of later caller cancellation. Existing permissions, confirmation and atomic credential replacement remain authoritative.

An existing interpreter file that cannot start is a launch failure with no child PID and unavailable runtime evidence. Preserve the original launch error and selected executable identity. A started process with no dump has a different, explicitly absent collection result; do not invent a crash.0.dmp observation.
