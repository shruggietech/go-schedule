# Research: Early-crash incident evidence

## 2026-11-08 preflight and owner scope decision

Retained S120-S123 comparisons had no matching original recurrence. Read-only Windows event inspection newly recovered Application .NET Runtime record 229042 at 2026-10-22T17:43:35.8766937Z, reporting pwsh.exe PID 479476, .NET 10.0.12 and CoreCLR 10.0.1226.41902. The original writer log lacks a PID, so temporal association is not definitive writer attribution. No original dump/WER artifact or matching PowerShell/Security process record was recovered. Private raw evidence remains under .test-temp/s138/preflight-2026-11-08T21-24-20-081Z.

The owner explicitly approved a bounded diagnostic improvement while #355/#356 remain open. This resolves the material scope question; no unsupported original-cause correction is promised.

## Adapter choice

Decision: use built-in Windows wevtutil via existing contained, hidden process support, then decode and validate a bounded XML response locally. Read-only research confirmed /uni:true emits UTF-16LE with a BOM, and the original 1023 event plus controlled FailFast 1025 share application/runtime identity lines. Alternatives: embedding another PowerShell query (depends on the runtime being diagnosed), adding a native Event Log wrapper (more surface), or collecting broad event exports (unnecessary unrelated data).

## Causal limits and primary sources

The original failure occurs in metadata/type resolution before authored script code, so prelude identity can be absent. Existing dump collection remains useful and unchanged; event readback supplements it. [dotnet/runtime #129458](https://github.com/dotnet/runtime/issues/129458) concerns a different concurrent attribute-loading report on 10.0.5, with later reported success on 10.0.9, and is not an established resolution for the recovered 10.0.12 event. [dotnet/runtime #134928](https://github.com/dotnet/runtime/issues/134928) concerns collectible-assembly/thread-static lifetime; no matching dump or assembly-lifetime evidence connects it to this failure. Neither source establishes original cause or a reason to replace an installed runtime.

## Attribution and dependency refinement

The adapter belongs in the existing integration fixture boundary because internal/clifixture reaches the API server, which already imports incident. A direct incident-to-clifixture import produced a retained compile-time cycle; keeping the generic receipt/filter in incident and injecting the platform adapter preserves the existing dependency graph. Native Windows GetProcessTimes is read from the exact completed process handle, with creation/exit times exposed only after completion. This avoids using a wider wall-clock interval that could include a reused PID. Unsupported/failed native timing remains explicit and does not alter the child outcome.

The skill-directed read-only agent confirmed .NET Runtime event 1025 for retained FailFast observations, and [Microsoft documents inherited process error mode](https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-seterrormode). The native test suppresses WER UI only in its own parent process and restores that mode. [wevtutil](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/wevtutil) and the [event schema](https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-elements) define the built-in readback and identity fields. No test child was crashed by the research agent.

## Development fixture and qualification corrections

The initial self-hosted race-test query helper exited 66 before its role marker and did not qualify overflow/cancellation coverage. A tiny standard-library non-race helper now models the native utility within unchanged query limits; the parent parser/adapter/ownership tests still run under race. All helper compilation also uses the existing contained launcher, bounded wait and cleanup, with bounded diagnostic readback. An initial aggregate was intentionally interrupted before this containment correction; it is retained as interrupted, not passed. The replacement qualification uses distinct input/log/result paths.

One development retry reused its foreground log filename. The first failure's tool-output facts are retained in a labeled transcription with explicit provenance, without claiming original-file byte identity. All historical protected evidence remained unchanged.
