# S108 Research

## Shared contracts and validation

Decision: Reuse existing jsonschema/v6 v6.0.3 with strict embedded Draft 2020-12 schemas and explicit format assertions, complemented by semantic validators. Rationale: already direct in go.mod; supports the required standard without a new dependency. Alternatives: handwritten structural checks risk schema drift; reflection-only schemas cannot establish cross-record truthfulness. References: [compiler format assertions](https://pkg.go.dev/github.com/santhosh-tekuri/jsonschema/v6#Compiler.AssertFormat), [JSON Schema validation](https://json-schema.org/draft/2020-12/json-schema-validation).

Decision: Confine artifact operations with os.Root, portable canonical relative paths, regular-file checks and hashing the opened handle. Rationale: prevent traversal/link races and validate exact evidence. Do not copy ordinary Lstat-then-Open hashing from unrelated release code. Reference: [Go Root](https://pkg.go.dev/os#Root).

Decision: Freeze a context-bearing sanitized corpus and let the executing agent produce structured reviews. Preserve execution outcomes separately from pending semantic review. Rationale: the requested layperson review needs interpretation without an external model service. Alternatives: automatic good-wording verdicts would fabricate evidence; adding an AI API would add unrequested credentials.

## Resource and ownership boundaries

Decision: Extend Windows suspended-before-assignment jobs with JobMemoryLimit, ActiveProcessLimit, hard CPU rate and BELOW_NORMAL_PRIORITY_CLASS. Retain hidden redirected handles, nested containment and exact drainage. References: [CPU rate](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_cpu_rate_control_information), [extended limits](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_extended_limit_information), [nested jobs](https://learn.microsoft.com/en-us/windows/win32/procthread/nested-jobs).

Decision: Use unique transient systemd-run --user workers inside the selected existing WSL distro/user, applying CPUQuota, MemoryMax, MemorySwapMax=0, TasksMax, Nice, RuntimeMaxSec, KillMode=control-group and noninteractive flags. The user manager and cpu/memory/pids delegation were inspected read-only in Ubuntu-24.04. Enforcement still requires actual owned-worker execution. References: [WSL commands](https://learn.microsoft.com/en-us/windows/wsl/basic-commands), [systemd-run 255](https://www.freedesktop.org/software/systemd/man/255/systemd-run.html), [resource controls](https://www.freedesktop.org/software/systemd/man/255/systemd.resource-control.html), [service lifecycle](https://www.freedesktop.org/software/systemd/man/255/systemd.service.html).

Rationale: Windows job bounds do not control Linux VM processes; cgroups bind the owned Linux tree independently. Linux TasksMax counts threads, unlike Windows process limits. Reject RLIMIT_AS for Go's reserved address space, RLIMIT_RSS as ineffective and RLIMIT_NPROC as UID-scoped. GOMEMLIMIT supplements OS controls but cannot replace them. References: [Linux resource limits](https://www.man7.org/linux/man-pages/man2/getrlimit.2.html), [Go GC guide](https://go.dev/doc/gc-guide).

## Scenario coverage

Read-only source inventory identified 27 declared roots and 119 runnable application paths, including runs/alerts parents; generated Cobra entries must be counted by the runtime exporter. Include systems search independently, watcher alias watch, inherited selectors and shadowed endpoint flags. Existing S107 scoped Dependencies and owned packaged fixtures cover the required seams. Functional scenarios must cover actual behavior independently of metadata probes.

Reuse fake-clock scheduling/catch-up/overlap, watcher, webhook, authority/TLS and MCP fixture patterns. The runtime evaluator must record its actual execution boundary and assert state/history/audit/logs rather than count an unrelated go test process as application coverage. Unknown-command/flag exit classification and absent arity validators are source candidates requiring separate reproduction and issue tracking if demonstrated.

The installed speckit-plan explicitly dispatched three read-only research agents for resource controls, scenarios and report contracts. Findings above resolve all planning unknowns; no host configuration was changed and no human decision remains.
