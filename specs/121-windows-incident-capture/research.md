# S121 research

## Runtime capture

Decision: use CoreCLR's child-only DOTNET_DbgEnableMiniDump=1, Mini type1 and exact unique dump/log destination. Rationale: Windows support uses MiniDumpWriteDump and the runtime's createdump child, retained by existing Job Object containment; no new dependency or license acceptance. Alternatives: ProcDump requires an unaccepted local EULA; global WER/debugger registration affects unrelated processes; a managed-only debugger cannot explain native CLR frames. Same-process pre-payload identity records loaded CoreCLR metadata, not neighboring-file inference. Process-local SetErrorMode in the owned child suppresses failure dialogs. Retained size and process watchdog are explicit, with no claim of a hard disk quota.

References: [CoreCLR design](https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md), [.NET crash configuration](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/collect-dumps-crash), [SetErrorMode](https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-seterrormode), [dotnet-dump native-frame limitation](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump).

## Owned IPC capture

Decision: private file request/response plus bounded runtime.Stack, sql.DB.Stats and an optional active-handler tracker. Rationale: independent of HTTP authorization and the single SQL connection; no new publicly reachable endpoint. Alternatives: same-API diagnostics fail at the same stall; always-on profiling perturbs execution; SIGQUIT/process termination changes the observed result. A live diagnostic watchdog captures before cancellation releases the pool queue. Every collector is stopped/awaited before workspace cleanup.

References: [Go diagnostics](https://go.dev/doc/diagnostics), [runtime.Stack](https://pkg.go.dev/runtime#Stack), [sql.DB.Stats](https://pkg.go.dev/database/sql#DB.Stats).

## Existing seams and research authority

The installed /speckit-plan command explicitly requires research agents. Two read-only researchers inspected actual code and official sources. Both finished without source edits or publication. Reuse internal/clifixture.StartPiped, Workspace.Contains, executionDiagnostics, case State artifacts and the S120 held-pool test seam. Store statistics are metadata only, never a new SQL capability. No unknown remains.
