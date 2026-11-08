# Implementation Plan: Early-crash incident evidence

**Branch**: `codex/138-early-crash-diagnostics` | **Date**: 2026-11-08 | **Spec**: [spec.md](spec.md)

## Summary

Extend the existing opt-in owned-script incident capture with a separate Windows runtime-event receipt for failed PowerShell children, including failures before the identity prelude executes. Original #355/#356 causes remain unresolved. Owner scope approval supersedes the provisional causal-resolution objective.

## Technical Context

**Language/Version**: Existing Go module/toolchain, test-only PowerShell fixture.

**Primary Dependencies**: Standard library, existing internal/clifixture containment (the platform adapter stays in the integration fixture boundary), existing golang.org/x/sys/windows and built-in Windows wevtutil. No new dependency.

**Storage**: Existing exclusively owned incident directory, one bounded runtime-events.json artifact. Raw event XML and system/user identifiers are not exported.

**Testing**: Failing-first parser/collector cases, exact PID/time/provider/application rejection, bounded query faults, actual hidden Windows pre-identity FailFast, Go race/lint and all eight foreground canonical gates.

**Target Platform**: Windows event readback; unchanged Linux/macOS behavior.

**Project Type**: Internal maintainer diagnostics, no public application or release change.

**Performance Goals**: Query only failed owned PowerShell children when capture is enabled; no successful-writer query overhead.

**Constraints**: Three-second query context plus existing five-second contained cleanup and one-second stream drain, at most sixteen retained events (query one extra to detect truncation), 256-KiB raw output per channel, existing 128-KiB metadata limit, bounded stream drain and contained cleanup. No global event/dump configuration, elevation or installed-service/runtime change.

**Scale/Scope**: One separately tracked diagnostic gap under #355, plus historical evidence record relevant to #355/#356. No broad replay or new QA framework.

## Constitution Check

Pre-design and post-design: PASS. I: existing process-tree ownership and explicit cleanup, no orphaned goroutines. II: failing-first discriminating evidence tests and preserved concurrency/security suites. III: native results and evidence availability remain separate. IV: fixed collection limits and no successful-operation query cost. V: installed Spec-Kit sequence, approved narrowed scope, authorized publication and owner merge boundary. No waiver.

## Research and design decisions

- Use built-in wevtutil with direct arguments through existing hidden, noninteractive contained process support. Prefer this small read-only adapter over a new native event API wrapper, a PowerShell runtime dependency or another collector framework.
- Decode BOM-marked UTF-16LE explicitly, then parse XML and re-filter locally. Query filtering alone cannot establish attribution. Retain only record/time/PID/event/provider/application and CoreCLR/.NET version fields; omit arbitrary Message/Stack, machine and account data.
- Accept only .NET Runtime 1023 (fatal runtime error) and 1025 (controlled FailFast), with exact owned PID and closed native creation/exit interval from the exact owned process handle and case-insensitive executable basename. Readback of older events outside that window cannot qualify the child.
- Keep the platform adapter at the existing integration-fixture boundary: importing clifixture from incident creates an indirect client/server dependency cycle. Preserve the package graph instead of refactoring unrelated application dependencies.
- Keep dump Collection.Status unchanged. An additional artifact contains its own collected/absent/unavailable/failed/truncated disposition; an artifact file's existence is not a claimed native dump or causal diagnosis.
- Run only after a nonzero native PowerShell result, never for disabled capture, success or launch failure. The event query receives independent collection allowance and cannot change primary elapsed/native/stream/cancellation results.
- Query-bound probes use one minimal non-race utility fixture built once through the same contained process launcher; parent tests retain race checks. This matches the small native utility budget without weakening collection limits.
- A controlled pre-identity crash suppresses inherited Windows error UI from its owned test parent, without altering system settings. A missing event log remains explicit; native proof on this host is inspected separately.
- Optional after-specify/after-plan context hooks are omitted because the core plan context update records this plan. Skill-directed read-only research verifies the Windows adapter and retained event schemas; no Codex review has been requested.

## Project Structure

- internal/incident/runtime_events*.go: bounded receipt/filter; test/integration/script_runtime_events*test.go hosts the existing-launcher platform adapter.
- internal/incident/runtime_events_test.go: discriminating parsing, attribution and failure cases.
- internal/clifixture/process*.go: native handle-derived creation/exit times with explicit unavailable state; test/integration/testscripts_observer_test.go integrates failed-child collection.
- test/integration/*incident*test.go: native pre-identity probe and unchanged outcome checks.
- docs/test-scripts.md: opt-in evidence boundaries.
- specs/138-early-crash-diagnostics/: design, investigation, tasks and verification.

## Verification and delivery

Preserve all prior inputs, retain first failures, run focused tests and actual hidden Windows controlled probe, then all eight canonical foreground gates. Freeze source and re-audit evidence. Publish only the official PR, monitor exact-head CI and findings within two total Codex requests including automatic review, and stop for owner final review/merge. Use Refs for #355/#356, with no causal closure claim.
