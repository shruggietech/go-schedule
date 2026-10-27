# S122 bounded investigation

## Source and limits

Accepted production baseline is merged S121 6f1de579bed866eb2bb328a25a0ad165caaa278a. New replay/observation tests and Spec-Kit records precede the baseline campaign; production manager and observer corrections follow it. Original failed and S121 finalized artifacts are frozen in .test-temp/s122/evidence-before.json.

Before execution, prepare source-bound Windows/Linux inputs and retain the build receipt. Campaign execution uses one shared thirty-minute wall-clock deadline, ten repetitions per normal selected condition, one deliberate prerequisite fault-matrix pass and sequential groups. Stop on the first unexpected condition failure for analysis. Preparation, controlled development regressions, implementation and later independent verification are separately identified.

## Planned and observed matrix

| Group | Selected existing test/condition | Repetitions | Observed outcome |
| --- | --- | ---: | --- |
| Current PowerShell | Six writers fresh/existing plus metadata/writer implicit/explicit same-SQLite | 10 | 160 attributable passes |
| Retained PowerShell 7.6.6 | Same conditions, child PATH selects retained portable runtime | 10 | 160 attributable passes |
| Current Bash | Same writer/selection contrasts | 10 | 160 attributable passes |
| Prerequisite fault matrix | Candidate refusal, launcher and inherited-pipe diagnostics | 1 | All three expected-fault test families passed |
| Native IPC component | Five original targets plus group prerequisites in memory and WAL | 10 | 140 passes, twenty removed runtimes |
| Packaged daemon WAL | Same seven unchanged catalog steps using source-bound real daemon/CLI | 10 | Seventy passes, ten removed runtimes |

Each normal runtime repetition contains both schema modes and all four operation/selection contrasts. IPC targets are service.restart, mcp.http.enable, mcp.http.rotate, mcp.bad-port and group.remove-child. group.create/group.child establish prerequisites. Preserve expected bad-port exit 2, recording service adapter, all state assertions, MCP resource/confirmation/retry checks and stale-token rejection. Component mode explicitly records actual selected database separately from logical owned runtime paths; it is not a packaged in-memory daemon claim.

## Runtime/prerequisite outcomes

The runtime matrix passed all ten repetitions on current PowerShell 7.6.9/.NET 10.0.14, retained portable 7.6.6/.NET 10.0.12 and Bash. It retained 480 observations: 360 concurrent-writer outcomes and 120 metadata/writer selection contrasts. The deliberate prerequisite fault matrix passed once. No spontaneous original CLR crash or historical prerequisite refusal recurred. The separate controlled #388 development regression already demonstrates existing-file launch failure being misclassified as command/absent with PID zero. Its native Windows launch error and exit -1 remain retained. This is not a spontaneous CLR recurrence.

## IPC/inspection outcomes

The memory/WAL component campaign passed 140 catalog observations across twenty fresh sessions; the packaged-daemon WAL campaign passed seventy observations across ten fresh sessions. These include all 150 original-target executions and sixty group prerequisites. Every original expected exit/state/protocol check remained active, including stale-token rejection and bad-port exit 2. All thirty runtimes were removed. No historical inspection stall recurred. The separate controlled #389 development regressions demonstrate canceled lifecycle mutations and a queued rotation waiting for a held Enable, then replacing the credential after cancellation. These are component-specific findings, not original-capture causality.

## Causal findings and corrections

#388 is a separately demonstrated script observer classification defect under #355. #389 is a separately demonstrated MCP lifecycle cancellation defect under #356. Keep their first failures and actual corrections distinct from original history. Original CLR/prerequisite and intermittent IPC causes remain unconfirmed until a matching recurrence supports a conclusion.

## Cleanup and unresolved outcomes

The campaign records exact ownership cleanup and retained private evidence for each fresh session. No installed runtime, service, account session, public endpoint, tag or promotion change is authorized. Original #355/#356 stay open if their causal outcomes remain incomplete.

### Completed campaign receipt

Shared campaign execution ran from 2026-10-27T17:41:38.714Z through 2026-10-27T17:47:49.040Z, 370.326 seconds, under the original 18:11:38.714Z deadline. The initial compiled-test launch used the wrong working directory and stopped before a writer ran; its failed receipt remains in campaign-01. The corrected launcher used normal package working directories without resetting that deadline. campaign-02 passed all six planned groups with zero source drift and all thirty original/finalized evidence artifacts unchanged. Prepared application/test source digest: 81198382ca942bc8c6fa666442836726870382fd8dd31d7b370fa1d8ae824089, base 6f1de579bed866eb2bb328a25a0ad165caaa278a, dirty=true. Private exact binaries, test executable hashes, all 690 observations, logs and cleanup records are retained under .test-temp/s122/.

The campaign precedes #388/#389 corrections. Native Windows development regressions separately failed before those fixes and passed afterward, including cancellation during Enable preparation and preserved committed Disable cleanup. The source corrections do not explain the original historical causes; #355/#356 remain open. Later full-QA and canonical verification are separately attributed in verification.md.

The private investigation-matrix.json expands every runtime/IPC observation into its condition, repetition and writer/case identity, with 690 rows and exact log references. All ten repetitions are present for each selected normal observation key. Original prepared inputs and raw campaign receipts remain unchanged.
