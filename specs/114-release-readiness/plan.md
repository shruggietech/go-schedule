# Implementation Plan: S114 release readiness

**Branch**: `codex/114-release-readiness` | **Date**: 2026-10-23 | **Spec**: [spec.md](spec.md)

**Input**: #354, #355 and #356 on main 7f7c5ec2c5252bb4b116b8d3b365438f30196779.

## Summary

Complete the source preparation and bounded Windows investigations authorized for S114. Correct the demonstrated release-note disclosure mismatch through existing promotion validation; improve demonstrably lost launcher diagnostics and concurrent failure attribution; exercise real owned IPC cancellation and mixed operations. Document original causes independently and leave unconfirmed bugs open. Prepare the exact later release procedure without staging, tagging or promotion.

## Technical Context

**Language/Version**: Existing Go 1.26 module, PowerShell 7 and Bash maintainer payloads, Markdown and JSON release metadata.

**Primary Dependencies**: Existing SQLite CLI, go-winio 0.6.4, x/sys 0.49.0, net/http and existing release-promotion helper. No new dependency or system runtime replacement.

**Storage**: Existing schema-v24 scheduler state and owned disposable script fixture databases; no new production schema.

**Testing**: Regression-before-fix for prepared notes and Bash launcher stderr; native Windows six-writer fixtures; real typed owned IPC/store requests and cancellation; native Linux comparison; all eight canonical gates.

**Target Platform**: Native Windows and Linux/WSL; shared release guidance includes macOS.

**Project Type**: Scheduler daemon, CLI and desktop with maintainer scripts and exact-source release workflows.

**Performance Goals**: Preserve 15-second case limits, six simultaneous writers and existing process/output/resource caps. Diagnostic payloads remain bounded.

**Constraints**: Foreground suites through verified hidden launcher. No login/logout, installed service, DNS, machine-wide runtime change, merge, tag, staging workflow or promotion.

**Scale/Scope**: Three original issues with independent outcomes, one source PR, at most two Codex rounds. Any demonstrated independent new bug receives a separate linked record.

## Constitution Check

Before research and after design: pass. Review branch, ordered installed Spec-Kit commands and blocking analyze gate are required. Behavioral changes receive failing-before/passing-after regressions and race verification. Existing migration, credential, access, concurrency and release assertions remain intact. The owner explicitly authorizes push/PR and verified in-scope review fixes; no merge/tag/release authority is inferred. Issue closure requires delivered functional outcomes, not tests or qualification receipts.

## Project Structure

```text
specs/114-release-readiness/
  spec.md
  plan.md
  research.md
  data-model.md
  contracts/diagnosis-release.md
  quickstart.md
  tasks.md
  delivery.md
test/integration/testscripts_binding_test.go
test/integration/testscripts_probe_test.go
test/scripts/lib/sqlite.sh
internal/clievaluate/ipc_diagnostics_test.go
scripts/release-promotion/promotion_test.go
.github/release-notes/v1.6.0.md
docs/INSTALL-{windows,linux,macos}.md
site/data/reviews/v1.6.0.json
```

**Structure Decision**: Reuse current helpers, real fixture boundaries and promotion validator. Do not add a new release service or speculative transport replacement.

## Phase 0: Research and decisions

- Preserve original S112 CLR access violation and separate SQLite refusal; suppressed historical candidate output prevents retrospective diagnosis. Current S113 observations do not establish their causes.
- Distinguish the independently proved release notes mismatch from unchanged original runtime/IPC causes. Production promotion requires explicit untested disclosure, but prepared notes omit it.
- Bash probe captures launcher stderr into a file and removes it unread. Retain a capped sanitized launch explanation separately from original candidate stderr/status.
- Add per-writer phase attribution and all-six results/rows on failure, retaining actual process exits and exact assertions. Do not retry or serialize.
- Exercise actual OS-IPC cancellation, a fresh succeeding request and mixed reads/mutations. Prior/current dependency runs are diagnostic comparisons only; do not speculate that a dependency caused the historical event.
- Investigate the PowerShell inherited-pipe deadline hypothesis separately. Change it only if a deterministic bounded reproduction proves an in-scope defect; otherwise record a specific follow-up rather than a speculative fix.
- Update only current upgrade/runbook references; preserve historical S112 schema-v23 records and public v1.5.1 bytes.
- Keep #354 open for actual later publication and #355/#356 open if causes remain unsupported. Source work, reviewed PR, tag and public artifacts are distinct.

## Phase 1: Design

Existing external JSON/report schemas remain unchanged. Diagnostic stderr and test failure logs gain bounded context. Release notes disclose native observations untested under the standing waiver and link unresolved known issues. All platforms receive stopped-daemon full-state backup and matching old-state/old-binary rollback guidance. The delivery runbook names exact later tag/source/run/attempt/release/candidate/asset/checksum identities and owner boundaries.

## Verification and delivery

Run meaningful regressions before correction. Preserve source, input, runtime and log identities for native Windows and Linux diagnostics. Run foreground `sh scripts/verify.sh all` with its eight gates. Full-catalog S113 evidence remains bound to ffebda5, with new fixture/helper regressions separately identified; a new full capture is necessary only if application/evaluator executable behavior changes. Format all GitHub/repository Markdown and check encoding before the authorized push. Open one official PR, audit actual reviews and exact-head CI, request no more than one additional Codex round after the automatic review, resolve every thread and stop for owner merge.

## Complexity Tracking

No constitution exceptions. Existing source notes require an explicit disclosure correction, not validator weakening. Original causal uncertainty is reported as incomplete bug outcomes, not hidden by the slice lifecycle.

## Output compatibility audit

The issue contract preserves existing outputs and native outcomes. Expanded per-writer phases and all-six failure detail are opt-in through GOSCHEDULE_PROBE_DIAGNOSTICS=1. No additional default output or file-size cap is introduced. The proven PowerShell stream deadline retains the existing timeout wording. The owner explicitly approved corrected default launcher-failure output: preserve the launcher status/explanation instead of the secondary missing-file failure. This exception applies only to that demonstrated failure boundary; successful output and existing writer assertions remain unchanged. Original before-fix logs are preserved.
