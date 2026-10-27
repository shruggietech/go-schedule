# S120 research and decisions

## Existing evidence and completed work

Two installed-plan research agents independently audited runtime/prerequisite and IPC evidence without modifying source or running tests. Original S112 writer log remains at the private artifact s112/verify-all-corrected.log, 11068 bytes, SHA-256 4eded1dd882f92b7f09aad62265065c46631212d79d97ce7b12922a4b2d3cacf. Its PowerShell writer 2 exited 3221225477 (0xc0000005), CLR 0x80131506, with GetTypeCoreIgnoreCase, TypeNameResolver and ParamBlockAst.UsesCmdletBinding frames. Separate Bash writer 0 exited 2 for unusable SQLite discovery. Each retained five of six records. The original text does not establish a crash dump or exact same-process runtime provenance.

Original IPC capture cb3c46fc-4134-4d90-9971-029284d04e6c remains interrupted/blocked, with original per-case failures under private s112/qa-run-01-interrupted/. Its identity SHA-256 is 3434f2f9c0ee68cac261150aa18daccaf60cff0d69c63bc8940d978579415632; build receipt SHA-256 e1288542124089a17ed233c9bef1a6e5b02a23a46b37c1d1a14e5aed5b377d29. service.restart command exit 0 preceded credentials inspection failure; MCP enable/rotate also had command exit 0, with runtime-info/tasks inspection or group DELETE failures. This is separate from S115 metadata deadline failures and S117 later passing executions.

Reuse S113-S117 phase/probe diagnostics, partial snapshots, memory/WAL overlap and cancellation fixtures, six concurrent writers and supported completion budgets. Do not repeat a full unchanged catalog or infer historical fixes from new passing experiments.

## Decision: discriminate present collection and cancellation gaps

Static hypothesis A: integration scriptCommand/executeBinding use CombinedOutput without a post-exit stream bound; requireSqlite has a context but no WaitDelay. A parent that exits while an owned descendant retains stdout/stderr can discriminate root completion from collection completion. Reproduce before correction. Prefer the existing contained process capability where a new runner needs tree cleanup; preserve actual native exit independently from collection/cleanup errors. Add outer opt-in provenance because a metadata-resolution crash can precede the existing interpreter-entry phase. Do not replace the existing phase protocol or successful default output.

Static hypothesis B: the store has one connection, while local actor resolution, authorization and affected handlers call context-free SQL methods. Hold the owned connection, observe a real typed native IPC waiter, cancel the client and require server-handler completion before releasing the pool. This differs from the existing pre-handler canceled-credentials fixture. Reproduce before correction. If demonstrated, use explicit Context method variants and request-context propagation on the affected read/cleanup/auth paths, preserving Background wrappers for legacy callers. No mutable store context, copied server, generic persistence rewrite or speculative dependency rollback.

Audit intent must remain durable before mutation; canceled requests must not be relabelled succeeded. Complete actual mutation evidence with an independently bounded completion context where needed, retaining uncertain intent if terminal persistence fails. Existing authorization, denied audit and identity mismatch tests remain authoritative.

## Alternatives considered

- Repeating the whole catalog: rejected because it does not distinguish the remaining causes.
- Replacing PowerShell/.NET or go-winio: rejected without a matching signature or causal experiment.
- Public profiler endpoint or global dump registration: rejected because owned fixture diagnostics suffice.
- Broad SQL/store facade conversion: rejected as disproportionate; explicit variants preserve old call contracts.
- New process containment: rejected in favor of internal/clifixture when process-tree control is required.

## Upstream interpretation and primary references

[Go 1.26 os/exec](https://pkg.go.dev/os/exec@go1.26.0) documents that context cancellation alone does not bound inherited output handles when WaitDelay is zero. [Go database cancellation](https://go.dev/doc/database/cancel-operations) describes propagation of request context into SQL operations. These support the proposed discriminating experiments, not historical attribution.

[dotnet/runtime #129458](https://github.com/dotnet/runtime/issues/129458) has a nearby type-name/attribute-resolution signature but describes concurrent threads in one WPF process, unlike the separate PowerShell writers here, and supplies no demonstrated matching fix. Treat it as a non-matching lead. [Microsoft Windows dump guidance](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/debug-windows-dumps) distinguishes native-capable analysis from managed-only dotnet-dump. If a crash actually recurs, use an available owned-PID collector; no new installed tooling or global configuration is presumed.

## Workflow decisions

The real checklist prerequisite helper requires plan.md. setup-plan copied its resolved template before checklist execution; substantive planning followed completed clarification/checklists. No mandatory extension hooks exist. The optional context script can spawn an unverified Python launcher, so it is skipped; the installed plan command's required CLAUDE marker update is performed directly. Source publication is explicitly authorized by this kickoff; merge/tag/promotion are excluded. No unresolved technical choice remains before task generation; test outcomes determine which conditional corrections are warranted.
