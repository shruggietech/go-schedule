# Implementation Plan: Correlated parallel task joins

**Branch**: `codex/144-parallel-task-joins` | **Date**: 2026-11-10 | **Spec**: [spec](spec.md)

## Summary

Deliver #470 as reusable parallel blocks: one successful entry Run creates immutable named branch invocations and one joined target admission after every exact invocation succeeds. Reuse the existing durable completion delivery queue and engine. Correct confirmed shared-path #471 queued eligibility and #472 Automation uncertainty in the same source increment.

## Technical Context

**Language/Version**: Go 1.26.0, existing Wails 2.17.1 desktop with TypeScript/React.

**Primary Dependencies**: Existing standard library, SQLite driver, Cobra and MCP SDK only. No new service, scheduler, canvas or dependency.

**Storage**: Current database schema 25, initialized directly; immutable join snapshots and explicit delivery-role unique indexes. Existing current-only refusal and explicit reset guidance remain. Portable v3 adds join definitions; maintained v1/v2 portable readers retain their existing behavior, reject new join fields under old schema declarations, and do not migrate any persisted store.

**Testing**: Failing-first store/engine/API/interface regressions; deterministic race tests, valid-current reopen, owned Windows/Linux packaged integration, headless browser coverage, full eight-gate foreground CI parity and exact-head hosted Windows/Linux/macOS CI.

**Target Platform**: Windows, Linux and macOS; local and explicitly enrolled independent remote daemons.

**Project Type**: Existing daemon/CLI/API/MCP plus desktop modules.

**Performance Goals**: Existing p99 dispatch <100ms and ten-percent benchmark regression budget. Indexed entry-definition lookup; transactional branch contribution and target decision; no polling scheduler or ambient history scan.

**Constraints**: All surfaces complete; same-task branch references remain distinct. No operator focus or native GUI interaction. Exact ownership, authority, audit, retry/uncertainty and output limits preserved. Two total Codex PR reviews including automatic; final merge remains with owner.

## Constitution Check

Pre-research and post-design checks pass principles I-V. Ownership and state transitions stay in current store/engine. Current schema initialization/refusal/reopen is retained without historical migrations or automatic deletion. Source changes receive meaningful tests first and all mandatory gates remain. The kickoff already authorizes push/official PR/CI fixes; no additional pre-push confirmation is required. No merge, tag, staging or public promotion is authorized.

## Design and decisions

- A parallel block owns explicit branch invocations rather than consuming ambient chain/task results. This preserves exact instance identity even when two slots reference the same task. Existing independent chains remain configured separate work.
- Use ParallelJoin, ParallelJoinBranch, ParallelJoinInstance and ParallelJoinInvocation domain records with exact delivery/run correlation. Current instance bindings and names are snapshotted; edits change future instances only.
- Extend CompletionDelivery with chain/join_branch/join_target roles and exact instance/invocation IDs. Partial unique indexes separately protect ordinary edges, branch invocations and target admissions. Reuse RecordRunAndCreateDeliveries, Claim/Resolve/Recover transactions. A repeated identical committed run is idempotent; conflicting identity reuse is rejected.
- Durable instance states are waiting, blocked, target_pending, completed and resolved. target_pending means a target decision exists, not native completion. Delivery pending/claimed/completed/resolved and live runs identify actual dispatch/running phases. Completed includes the actual target outcome, including failure/cancellation.
- Retain live correlation in Run.DeliveryID, JoinInstanceID, JoinBranchID and JoinRole, with no duplicate JoinID. Merge active runs by DeliveryID into bounded API observations, sampling active before persistence; never persist a fabricated running/terminal record.
- Runtime snapshot identifiers do not cascade away with live task/definition deletion. Explicit invalidation resolves pending work; native work already started keeps its owned result and cannot revive a resolved instance. One branch failure does not invent sibling fail-fast cancellation.
- Add typed task.Sources (Completion, Trigger, Watcher, Join) rather than extending the existing positional boolean list. Join branch/target destinations and outstanding immutable snapshots count toward readiness; edits cannot indirectly disable old in-flight bindings. Explicit owner disablement still applies at admission.
- Validate the combined ordinary-chain/join graph in the same mutation transaction, including projected portable replacements. Reuse one graph algorithm; no separate validators that disagree.
- #471 uses one fresh automatic-delivery eligibility helper at initial dispatch, queue release and worker acquisition, preserving manual-run semantics and first-pending queue-one collapse.
- Full remote Automation requires explicit remote chain CRUD as well as new join routes; a partially loaded joins-only remote workspace is not selected. New and existing Automation mutations preserve #472 uncertainty, exact expected daemon identity, current Manage authority, drafts and authoritative same-target refresh fences.
- Reuse current Automation/Activity UI and command patterns. New top-level join create/list/show/update/rm/instances/instance commands and joins_manage MCP family expose the whole capability. Bounded MCP definition/instance resources plus exact and nested-branch pagination avoid loading unbounded history or silently omitting prerequisites.

## Project Structure

Shared domain contracts: internal/domain/joins.go and domain.go. Durable implementation: internal/store joins files, chains/current schema/CRUD/readiness/portable helpers and internal/task/readiness.go. Engine/API/remote/bundle work stays in internal/engine, internal/api/server, internal/authorization, internal/remote, internal/bundle and current daemon event wiring. Interfaces stay in internal/api/client, internal/cli, MCP modules and desktop Automation/operations/frontend. Root integrates OpenAPI/generated bindings, command catalogue, documentation, packaged validation and benchmark/freeze evidence.

## Parallel ownership and execution

The installed speckit-plan command explicitly dispatched three research agents for store, engine/API/bundle and interfaces. Reports are consolidated in research.md. After the blocking analysis gate and shared domain foundation, independent file-owned [P] tasks may run together under speckit-implement; tasks touching shared files stay sequential. Root owns shared contracts, OpenAPI/catalogue/docs and final integration. Agents do not request formal Codex reviews or publish.

## Tooling and retained boundaries

The checklist prerequisite requires a plan scaffold before checklist generation; setup-plan copied that scaffold before the checklist helper, without performing design out of order. Optional context extension hooks are skipped; the required CLAUDE plan marker is updated directly. All execution uses verified hidden launchers and redirected noninteractive I/O. Historical/S143 evidence is hashed for preservation; no failed receipt is rewritten or used to qualify changed source.
