# Implementation Plan: Correlated parallel recovery and cleanup

**Branch**: `codex/145-parallel-recovery-cleanup` | **Date**: 2026-11-10 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #477 as one complete extension of existing parallel instances: settle every exact branch, select the success or configured recovery follow-up, and admit optional cleanup only after that follow-up settles. Reuse durable completion deliveries, current eligibility checks and the existing engine event loop. Preserve all-success-only definitions when optional fields are absent.

## Technical Context

**Language/Version**: Go 1.26.0; existing Wails 2.17.1 and TypeScript/React desktop.

**Primary Dependencies**: Existing standard library, modernc SQLite, Cobra and official MCP SDK. No new package, service or orchestration framework.

**Storage**: Direct current schema 26 initialization, immutable recovery/cleanup snapshots, unique per-role deliveries and explicit branch settlement. Schema 25 is refused non-mutatingly under the current-only prelaunch policy. Portable v4 represents new routing configuration; unchanged base joins remain v3, and maintained readers retain their existing contracts.

**Testing**: Failing-first store/engine/API/interface tests, concurrency and current reopen, actual owned Windows/Linux packaged workflows, headless Chromium, paired benchmarks and complete foreground `sh scripts/verify.sh all`; exact-head hosted CI follows publication.

**Target Platform**: Windows, Linux and macOS; protected local and enrolled independent remote daemons.

**Project Type**: Existing daemon, CLI/API/MCP and Wails desktop modules.

**Performance Goals**: Existing nominal dispatch p99 below 100ms; unchanged accepted-main benchmarks paired with final candidate. Any regression above ten percent requires recorded evidence and justification. No polling loop, ambient history scan or scheduling-result cache.

**Constraints**: No operator focus or visible/interactive child processes. No historical migration, release or cluster implementation. At most two total Codex PR reviews, including automatic; stop for owner final review/merge.

## Constitution Check

Pre-research principles I-V pass: no new dependency, state remains daemon-owned and transactional, tests precede behavior changes, actual run ownership is authoritative, current initialization/refusal persists and measured performance/gates are retained. Post-design checks also pass with complete all-interface propagation and no unresolved contract decision. The owner kickoff already authorizes push/official PR and in-scope CI/review fixes, so no additional publication confirmation is needed. Merge/tag/packages/releases are outside scope.

## Design and decisions

- Add optional RecoveryTaskID/RecoveryOn/CleanupTaskID; conditions failure/canceled/skipped/not_started. Nil omitted conditions default to failure when recovery exists; explicit empty conditions deliberately match nothing and survive transport/export. Nonempty conditions without recovery are rejected.
- Preserve Target* as success-only evidence. Add separate Recovery*/Cleanup* identities/outcomes/live cancellation metadata and SelectedRoute pending/success/recovery/none. Snapshot configured references and conditions per entry instance.
- Branch settlement is explicit success/failure/canceled/skipped/not_started. Native cancellation_failed maps to failure only after the actual executor owner returns; live cancellation-request failure is not termination. Queued/caught-up markers never complete an incoming delivery.
- Store advancement occurs transactionally after actual native result or documented pre-native no-start resolution. Invalidations remain distinct direct SQL paths with a global resolved guard; a resolved delivery alone cannot prove an invalidated native command stopped. Unique role indexes and exact binding retain idempotence.
- Select recovery when any configured cause is present after all branches settle; otherwise retain unselected reason. Cleanup can follow an unselected route or definitive no-start, but waits for a selected actual native run. Cleanup has its own terminal result and cannot rewrite upstream outcomes.
- Engine uses one buffered coalesced completionWake consumed in its existing event loop; drain until zero claimed. This handles continuation creation during no-start resolution and removes recursive record/overlap-skip draining without a new goroutine or polling. Startup still drains directly.
- Graph validation models ordinary chain edges plus join entry to every owned destination. Finite intra-instance role sequencing is not a new automatic task activation edge, so reusing a task across finite stages does not become a false self-cycle. Every stage can activate chains/new joins, so entry-to-all destinations still detects potential automatic cascades. The same validator handles projected atomic portable replacement.
- APIs keep current endpoints and complete Manage definition reads. Add optional request fields; RecoveryOn is a pointer in create/PATCH DTOs to preserve explicit empty versus omitted. Domain observations retain plain arrays without omitempty. Authoring uses full-text hydration and sparse CAS.
- Advertise parallel-join-recovery alongside parallel-joins. Older compatible Automation workspaces remain usable; gate only new options and capability-requiring portable applications. Remote authority, audit, retry policy and current daemon fences remain unchanged.
- Portable v4 carries all new references/conditions and remaps/drifts/validates them. Base v3 export omits absent new fields so older strict readers remain compatible. No instance/run history is exported.
- CLI adds --recovery, --recovery-on and --cleanup with explicit empty-ID clearing and none for an empty cause list. MCP management and bounded observation maps propagate every new field; desktop uses four cause controls, routing-order help and distinct exact stage inspection.

## Project Structure

Shared contracts: internal/domain/joins.go. Persistence/task readiness: internal/store and internal/task, including current_schema.sql, join decision/readback, invalidation, graph and portable readiness helpers. Engine/API/remote/bundle: internal/engine, internal/api/server, internal/remote operations, internal/bundle and relevant manifest/event wiring. Interfaces: internal/api/client, internal/cli, MCP modules, desktop app/automation/operations and frontend. Root owns domain, OpenAPI/generated contract, CLI catalogue, docs, packaged integration and performance/final qualification.

## Parallel ownership and execution

The installed speckit-plan command mandates research agents; store, engine/API/portable and interface agents performed independent read-only research and coordinated contracts. Their results are consolidated in research.md. After blocking analyze passes and root writes the shared domain foundation, distinct file-owned implementation tasks may proceed in parallel. Root retains shared/generated files and does not compete for agent-owned source. Heavy checks use bounded GOMAXPROCS=2 and GOFLAGS=-p=2; complete aggregate is one watched foreground run after source settles.

## Tooling and evidence

The installed checklist prerequisite expects plan.md even though the project sequence puts the requirements checklist before planning. Its original failure is retained; pre-plan path discovery and requirements review were followed by the real setup-plan and successful full prerequisite validation. No gate was reported passing while absent. Optional context hooks are skipped at specify, and the installed agent-context updater runs after plan. Retain 14,794 protected historical/prior-slice inputs and original failures; all new receipts use .test-temp/s145 and do not rewrite old evidence.

## Complexity Tracking

Principle IV exception is recorded in verification.md: existing isolated last-branch decision 477,188 ns to 680,140 ns (+42.53%, +0.203ms) after removing redundant reads; full 100-target completion improves 25.27 percent and nominal dispatch is within the ten-percent regression budget. Durable settlement/routing evidence is retained instead of omitted or hidden behind a cross-call cache. Inputs/thresholds are preserved with one identical live benchmark harness on both implementations. This deviation is surfaced on the PR.
