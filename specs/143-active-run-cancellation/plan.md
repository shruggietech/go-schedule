# Implementation Plan: Active-run cancellation

**Branch**: `codex/143-active-run-cancellation` | **Date**: 2026-11-09 | **Spec**: [spec](spec.md)

## Summary

Deliver #465 through exact Run ID cancellation across executor, engine, local/remote API, CLI, controlled MCP and desktop. Correct separately reproduced #466 desktop API negotiation, #467 concurrent accounting and #468 exact-record refresh in the same end-to-end source increment.

## Technical Context

**Language/Version**: Go 1.26.0, TypeScript and React through the maintained Wails 2.17.1 frontend.

**Dependencies**: Existing standard library, golang.org/x/sys platform APIs, Cobra, existing MCP library and current frontend tooling. No new dependency.

**Storage**: Existing current SQLite schema 24. Outcome TEXT supports canceled and cancellation_failed without changing table shape. Active cancellation status/error are transient observation fields, while terminal outcome and bounded diagnostic output persist through existing run history.

**Testing**: Failing-first focused Go race regressions; owned hidden Windows and Linux/WSL processes, existing macOS CI, frontend component/headless browser coverage and the full eight-gate canonical aggregate.

**Platform**: Windows, Linux, macOS; independent local/identity-pinned remote daemons.

**Performance**: O(1) active handle/count operations per run; storage fallback only after active lookup misses. Preserve p99 dispatch <100ms and existing benchmark budget. No unbounded cancellation history cache.

**Constraints**: Five-second original cleanup allowance; acknowledgement does not claim termination. Unreaped native work stays owned/active. Hidden Windows console launch, exact command/run-as/stdin/output behavior, no operator focus, sanitized public evidence and two-total Codex review cap.

## Constitution Check

Pre-research and post-design checks pass principles I-V. Owned goroutine/handle lifetimes are explicit; safety-critical behavior receives race-enabled regressions before source changes; interfaces share one contract; clock use remains injected in engine; active operations remain bounded. No historical migration, automatic reset, dependency, workflow, tag or release change is required. User explicitly authorizes branch/official PR publication and in-scope fixes; final merge remains with the owner.

## Design

- Shared domain observation/result values plus a small internal/runcontrol context-cause/failure-feedback contract join engine and executor without coupling either to API or QA fixture policy.
- Executor owns process and pipe handles. Unix adds a process group while retaining exec.Cmd run-as behavior. Windows uses the proven native suspended-create/job-assignment/resume technique with only selected inherited streams and no QA resource limits. BREAKAWAY_OK preserves explicit application-requested Windows detachment without silent breakaway; intentionally escaped work remains outside the owned target. Normal completion retains ordinary background-child semantics; cancellation alone starts the cleanup allowance and closes owned pipes if collection fails.
- Each active Run ID registers its own cancellation handle before execution. Completion freezes the observed native result before persistence; repeated requests return its actual requested/already-finished state. A metadata-only context-aware store lookup handles completed IDs without reacquiring SQL under an uncancelable request.
- Cancellation failure feedback updates the active observation promptly. A reaper remains responsible for unreaped work; capacity is not falsely released. Final canceled/cancellation_failed outcomes suppress completion chains and outbound failure/recovery policy streak changes; run_finished events retain truthful outcomes.
- Per-task in-flight counts fix #467 so cancellation/completion of one concurrent execution does not release another execution's overlap accounting.
- Existing authority/audit/remote retry policy gains runs.cancel as Operate and RetryNever. CLI gains exact-run cancellation and active ID discovery. MCP adds bounded active-run discovery and gains runs_cancel using its current target/request-id deduplication and revocation checks.
- Desktop cancellation captures the rendered daemon identity, uses existing mutation/current-state guards and a synchronous submission lock, refreshes exact active/history state and reconciles at 250ms for a bounded eight seconds. Expiry remains unconfirmed with Refresh available.
- #466 removes product-major coupling and requires exact advertised v1 local/remote API support alongside existing health/version/identity/access checks.

## Project Structure

Production changes remain in internal/domain, internal/runcontrol, internal/executor, internal/engine, internal/store, internal/api, internal/authorization, internal/remote, internal/cli, internal/mcpoperate and their current consumers; desktop/operations, desktop/connection, App and frontend operations bindings/views. Existing OpenAPI/docs/catalogue contracts are updated coherently.

## Execution and parallel ownership

The installed speckit-plan workflow dispatched three read-only research agents for process ownership, engine/API semantics and client/desktop integration. Implementation follows the task dependency order; independent file-owned executor, engine/API and interface tasks may proceed in parallel after shared contracts and the blocking analysis gate. Integration, shared-model changes, source freezes, canonical verification and publication remain sequential.

## Decisions

Use OS job/group ownership rather than process-name/PID scans or QA launchers. Keep distinct cancellation outcomes in existing run storage instead of adding task-schema/reset work. Preserve unreaped native work instead of reporting timer expiry as termination. Negotiate advertised API version rather than extending the obsolete product-major allowlist. Skip optional extension hooks; maintain the required CLAUDE plan reference directly. The checklist helper initially required plan.md before its documented pre-plan position; setup-plan copied its scaffold and the checklist prerequisite readback then passed.
