# S120 attributable investigation

## Source and preserved originals

S120 begins at accepted S119 merge dac2e97895353435ad30515402b9f99d5e5d4ebc with the separately identified pending S120 tree. The private evidence-before.json records 26 unchanged original S112/S115/S117/S118 artifacts, including frozen failed reports and executables. The original writer log remains 11068 bytes, SHA-256 4eded1dd882f92b7f09aad62265065c46631212d79d97ce7b12922a4b2d3cacf. The S112 capture remains interrupted/blocked with original per-case failures. Neither later success nor S120 component corrections overwrite those outcomes.

## Individual dispositions

| Original category | Controlled observation | Supported disposition |
| --- | --- | --- |
| #355 PowerShell CLR 0xc0000005 / 0x80131506 | Original metadata-resolution frames retained; current metadata-only and writer contrasts passed | Original cause unconfirmed; no dump or matching demonstrated upstream fix; #355 remains open |
| #355 Bash SQLite prerequisite refusal | Current implicit discovery and explicit selection of exactly the same candidate passed; existing malformed/old/timeout/launcher distinctions retained | Historical refusal unconfirmed; passing current probes do not explain its original cause |
| #356 service.restart credentials inspection | Real native credentials request cancellation plus memory/WAL sustained snapshots passed | Separate reproduced pool defect corrected under #382; original intent-only service case cause unconfirmed |
| #356 mcp.http.enable/rotate runtime-info/tasks inspection and group cleanup | Real native runtime-info/tasks/DELETE request shapes exercised against a held pool, with cancellation before release | Separate SQL cancellation correction under #382; original MCP/protocol outcomes remain independently attributable |
| #356 mcp.bad-port/group.remove-child canceled execution | Existing partial/canceled observation guarantees retained with server-side queued cancellation | Original intermittent cause unconfirmed; no blind mutation retry or replacement capture |

## Demonstrated #381 correction

On unchanged executeBinding, an owned parent exited 0 or 37 while a descendant retained stdout/stderr. Both calls reached the independent four-second guard and completed only after the fixture released the descendant. This demonstrates unbounded collection, independently of any CLR failure. The corrected contained runner preserves native 0/37, bounds execution at 120 seconds, allows two seconds to drain streams, caps each retained stream at one MiB and separates collection/cleanup errors. Exit 0 plus incomplete collection remains a failed fixture. The overflow regression retains native 37 and exactly one MiB with an explicit collection error. Existing six-writer fresh/existing-schema, distinct-identity and exact-row assertions remain unchanged.

## Demonstrated #382 correction

A real typed native named-pipe ListTasks request queued behind the test-owned single SQL connection (WaitCount=1, InUse=1). Canceling the client left the unchanged server handler blocked until the pool was released, in both memory and file-WAL stores. A preceding sandbox access-denied run never reached the pool and supplies no causal evidence. After explicit request-context propagation, handlers finish before pool release for tasks, runtime-info, credentials and group DELETE on both stores. Twelve direct affected SQL operations also cancel while queued, retaining the group and uncertain audit intent. Existing context-free store callers delegate through Background.

Audit intent remains durable before mutation. After actual handler execution, terminal audit persistence receives an independent 30-second context through context.WithoutCancel, so disconnect cannot erase an already committed outcome. The real audit regression commits a group, cancels the caller and retains its succeeded audit target. A failed/canceled terminal SQL write retains uncertain intent; existing fail-closed authorization and audit checks remain active. This is a bounded affected-path correction, not a claim that every store operation or original transport stall was diagnosed.

## Runtime and environment contrasts

Windows Go race checks passed with PowerShell 7.6.9/.NET 10.0.14, Windows 10.0.26300 X64, Bash 5.3.15 and SQLite 3.47.2. The current PowerShell executable SHA-256 is 2a6769a3330dc1084d0154bfde7d164ab750ca7f700f313f9d417c1be78b2233; Bash is 9eedc85f00d4b090f50ce4bede70802619212baa90c15e57893413bca0de9fad. These differ from retained historical PowerShell inventory and cannot be retroactively assigned to the original crashing process. Both metadata-only and writer conditions passed implicit/explicit same-candidate selection without retries. Metadata conditions created no database; writers retained the exact literal. Opt-in receipts include exact executable/harness hashes and root PID, excluding environment values and arbitrary payloads.

Existing nonroot Ubuntu 24.04 WSL repeated the controlled pool, audit, inherited-stream and overflow regressions with Go 1.26.0, Bash 5.2.21 and SQLite 3.45.1; Bash contrasts and six-writer assertions passed. PowerShell was unavailable and explicitly skipped. Existing sustained complete snapshots, concurrent writes and cancellation recovery passed on both native environments. No full unchanged catalog was substituted for these discriminating experiments.

## Evidence and remaining work

Private S120 receipts retain regressions-before.log (original inherited-stream failures plus separately attributed sandbox refusal), pool-before-native.log (actual pool regression failure), regressions-after-01.log, targeted-native-01.log (includes the initial new audit-test duplicate-route setup failure), contrasts-audit-02.log (corrected setup and contrasts), targeted-native-final.log and targeted-linux-01.log. The new audit test's route setup was corrected; its transient setup panic is not reported as a product defect. Canonical acceptance binds final source separately in [verification.md](verification.md).

#355 and #356 retain their unmet historical causal outcomes. If the CLR failure recurs, preserve the exact owned PID/executable/runtime, native exit, completed phases and an available native-capable dump. [The nearest runtime lead](https://github.com/dotnet/runtime/issues/129458) differs in concurrency/runtime context and establishes no matching supported fix. No runtime replacement, global debugger registration, installed-service operation, merge, tag or public promotion occurs in S120.
