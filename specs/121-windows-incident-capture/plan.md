# Implementation Plan: Windows incident capture and causal resolution

**Branch**: `codex/121-windows-incident-capture` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)

**Input**: S121 #355 and #356, accepted base 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7.

## Summary

Add explicitly selected incident capture to existing disposable script/API fixtures and unattended evaluation. Use CoreCLR child-only Mini dump configuration for the owned PowerShell process. Use a private file rendezvous independent of HTTP/storage for an owned daemon's bounded goroutine/pool/active-request snapshot. Preserve all native, inspection and cleanup outcomes. Diagnose and correct actual new defects; leave unconfirmed historical causes explicitly unresolved.

## Technical Context

**Language/Version**: Go 1.26, existing PowerShell/Bash harnesses.

**Primary Dependencies**: Standard library, existing UUID and Windows process containment. No new third-party runtime or implicit EULA.

**Storage**: Exclusive private incident subdirectories, versioned JSON and bounded artifact inventory. No persisted application-schema migration.

**Testing**: Owned native IPC memory/WAL and held-pool fixtures, intentional crash child, malformed/escape/oversize/cancel cases, targeted race checks and canonical eight gates.

**Target Platform**: Native Windows and existing nonroot WSL Linux. Ordinary daemon behavior remains disabled by default.

**Project Type**: Maintainer QA and daemon diagnostics.

**Performance Goals**: Disabled path adds no collector/polling. Enabled daemon polls at 100ms; at most eight incident snapshots per owned daemon, active requests capped at64, stack buffer64KiB. Watchdog selects five seconds or half the remaining caller budget, whichever is smaller, and never changes case cancellation.

**Constraints**: Each request/response exchange has a two-second independent collection deadline, request <=1KiB, structured snapshot <=128KiB, retained Mini dump <=32MiB and metadata/log <=64KiB each. Runtime dump generation remains within the existing 120-second owned-process watchdog; retained-byte limits do not claim an OS file quota. Collector failures and truncation are explicit.

**Scale/Scope**: Reuse S120 containment, cancellation, runtime baselines and existing strict report envelope. One incident per case at most, eight per daemon. No profiler HTTP listener, installed-service operation, account-session change, global debugger registration, dependency-baseline substitution or blind mutation retry.

## Constitution Check

Before research and after design: PASS. Simplicity: two small shared diagnostic primitives and existing containment; no replacement QA framework. Testing: regressions precede implementation, access/clock/recovery assertions preserved. UX: plain actionable collection statuses and unchanged primary outcomes. Performance: finite buffers/count/time; no scheduler hot-path changes. Autonomy: source push/PR/two-round review are explicitly authorized; merge/tag/promotion are not.

## Project Structure

- internal/incident/: fixed incident format, exclusive artifact handling, child crash configuration, private daemon rendezvous and bounded active requests.
- internal/api/server/: optional request tracking outside authorization, exposing no new route.
- internal/store/: read-only connection-pool statistics only.
- cmd/goschedd/: opt-in private responder tied to daemon context.
- internal/clievaluate/: explicit capture option propagated through coordinator/worker configuration and resume identity, live watchdog and independent incident summary in case state.
- cmd/gosched-qa/ and scripts/cli-qa.ps1: --incidents/-Incidents on capture/resume, additive maintainer option with existing policy forwarding.
- test/integration/: CoreCLR pre-payload identity/configuration and owned FailFast regression.
- docs/cli-qa.md, docs/test-scripts.md: enablement and incident interpretation.
- specs/120-windows-causal-diagnostics/: historical merged delivery references only.

## Decisions and alternatives

Use built-in CoreCLR createdump instead of ProcDump because the supported runtime already contains it and no external EULA is needed. Record loaded runtime identity inside the same owned PowerShell child before the original harness body; do not wrap and alter its native exit semantics. If no dump exists after a nonzero exit, report absent/unavailable, not collected.

Use an independent private file side channel instead of a profiler endpoint or a query through the blocked API. The evaluator creates the exact directory; daemon startup validates absolute canonical directory identity and uses exclusive response creation. Requests bind session nonce, request ID and expected PID. Only fixed route patterns and safe method/stage categories are tracked; bodies, queries, credentials and arbitrary IDs are excluded. Raw stacks remain local private evidence and only their bounded manifest enters report state.

Capture while the request is still alive, before cancellation can release the S120 SQL queue. A canceled collection never changes command or inspection success. Explicit --incidents is recorded in worker and resume identity so resumed environments cannot silently change diagnostics.

Keep the original #355 and #356 evidence immutable. A passing controlled or full-catalog follow-up cannot establish a historical cause. Newly demonstrated bugs receive separate issue references. Reconcile S120 source delivery without changing v1.6.0 immutable release identity.

## Verification and delivery

Spec-Kit tasks/analyze remain blocking. Run targeted native and WSL race experiments, retained/current PowerShell comparisons and one new source-bound Windows/WSL catalog after source changes. Run sh scripts/verify.sh all in the foreground, preserving first failures and input manifests. Finalize Implemented metadata with objective review-branch delivery before the final lifecycle check. Push/open the official PR under current authorization, handle every finding in at most two Codex rounds, require terminal successful exact-head CI, then stop before merge/tag/promotion.

## Complexity Tracking

No constitution violation. Raw native dumps remain private and may require native-capable analysis; missing historical same-process state remains a disclosed causal limit.

## Pre-publication design completion

Unexpected failures before the live threshold receive one best-effort two-second post-case snapshot, explicitly distinguished from live evidence; expected invalid-input refusals do not. Capture directory/ID fields are bound to private creation identity so changing a display field cannot redirect writes or native dump configuration. Snapshot UTC time and cumulative wait units support interpretation. These complete FR-002/FR-006/FR-007 and do not alter command checks or supplied budgets.

## Implementation discovery (2026-10-27)

Independent #385 corrects the demonstrated desktop search test race with Shell route focus. The change is limited to explicit readiness and intended-query assertions in desktop/frontend/src/App.test.tsx; product navigation, target identity checks and timeout policies are unchanged. The failed foreground aggregate is retained, and final verification follows this correction. The completed CLI catalog remains attributable to its frozen input source; this later GUI-only test correction does not change the evaluated CLI or daemon implementation.
