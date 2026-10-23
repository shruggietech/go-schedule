# Implementation Plan: S113 reliable unattended QA and actionable delivery diagnostics

**Branch**: `codex/113-qa-diagnostics` | **Date**: 2026-10-23 | **Spec**: [spec.md](spec.md)

**Input**: Full approved #355-#358 and PR #360 bundle, based on main 11c2a9a46d4610fe8a67903cea246050224bd606.

## Summary

Repair concrete evaluator lifecycle and diagnostic loss, investigate the separate intermittent runtime/IPC symptoms, and correct SMTP success/error semantics with compatible persisted history. Incorporate the bounded Fumadocs patch without unrelated dependency changes. Source fixes and issue closure remain separate: no passing repetition establishes a historical root cause.

## Technical Context

**Language/Version**: Go 1.26, PowerShell 7, Bash, TypeScript/React and the existing Wails module.
**Primary Dependencies**: Existing typed local/remote HTTP clients, contained process helpers and SQLite store; Fumadocs core/UI 16.17.7 and MDX 15.4.8 from PR #360. No new runtime dependency.
**Storage**: Forward-only SQLite migration 24, diagnostic payloads inside existing hash-bound v1 evidence artifacts, private receipt-bound evidence artifacts.
**Testing**: Regression-first Go race tests, both script twins, real owned SMTP/IPC adapters, desktop components/browser checks, canonical scripts/verify.sh all and a fresh owned Windows/WSL evaluation.
**Target Platform**: Native Windows and Linux/WSL; hosted macOS remains covered by existing CI.
**Project Type**: System daemon, thin CLI/desktop clients and maintainer evaluator with a static documentation website.
**Performance Goals**: Existing 15-second case and selected total capture budgets remain intact; bounded shutdown/report reserve is separate. Avoid repeated physical writes when exact sanitized bytes are already persisted; benchmark that evidence path.
**Constraints**: Existing 25-percent CPU, memory, process, channel and artifact limits are retained. Every child is hidden/noninteractive and owned. No machine-wide runtime replacement, installed-service mutation, account session operation, tag or release.
**Scale/Scope**: Existing complete catalog, two environments and all message reviews; four defect records plus one existing dependency PR.

## Constitution Check

- Code quality: reuse containment and typed clients, bounded optional observations, explicit errors and joined owned cleanup. Pass pre/post design.
- Testing: write regressions before each behavioral fix; preserve six concurrent writers, secret/auth assertions, forward migrations, clock/DST/recovery coverage and core cross-package coverage >=80 percent. Pass pre/post design.
- UX: neutral diagnostics and actual errors have explicit labels, original failures survive, terminal recovery instructions name exact original inputs. Pass pre/post design.
- Performance: retain limits; measure unchanged-artifact finalization and request phases before attributing stalls or changing transport. Pass pre/post design.
- Autopilot: execute installed specify, clarify, checklist, plan, tasks, analyze, implement in order. Current owner request already authorizes PR publication, so no redundant pre-push question. Two Codex invocations maximum; stop before merge. Pass pre/post design.

## Project Structure

### Documentation

`specs/113-qa-diagnostics/` contains spec.md, plan.md, research.md, data-model.md, contracts/diagnostics.md, quickstart.md, tasks.md, analysis.md and delivery/verification records. The selected feature and CLAUDE.md plan reference must both remain 113.

### Source Code

- cmd/gosched-qa/main.go: validated timing-policy entry point, capture/resume budget handling and actionable interruption output.
- internal/clievaluate/: bounded lifecycle reserve, latest checkpoint recovery, unchanged-artifact optimization, partial observations and optional timing metadata with strict validation/rendering.
- internal/api/client/: context-scoped request observation limited to safe method/route/category/timing metadata, without headers or bodies.
- scripts/cli-qa.ps1: consume the same validated timing policy, preserve output on timeout and identify exact recovery inputs.
- test/scripts/lib/Sqlite.ps1 and sqlite.sh, test/integration/testscripts*.go: numeric candidate validation, bounded candidate/process diagnostics, explicit-path semantics and unchanged concurrent writer outcomes.
- internal/domain/notification.go, internal/store/notifications*.go and store.go: additive diagnostic field and migration24.
- api/openapi/remote-v1.yaml, generated internal/remote/openapi.gen.go, internal/cli/notification.go, desktop/notifications/ and frontend notification models/details: compatible neutral/error presentation.
- site/package.json and package-lock.json: only PR #360's existing three patches.
- docs/, CHANGELOG.md, pending v1.6.0 notes/comparison review and specs index: schema24 upgrade and source/publication reconciliation.

**Structure Decision**: Extend existing modules and contracts; do not add a second process launcher, reporting framework or duplicate documentation source.

## Decisions and Implementation Order

1. Shared additive diagnostic contracts first. Timing/request/completeness metadata lives inside existing hash-bound state evidence artifacts, preserving the strict non-omitted v1 report envelope and retained report serialization. Complete-state assertions require complete observations even if partial maps happen to match.
2. #357: derive the outer capture/resume deadline from Go-validated policy, not duplicated PowerShell defaults. Keep the selected total deadline, allow bounded worker shutdown and aggregate finalization after execution cancellation, preserve the newest durable per-case outcomes, and emit explicit original-input recovery guidance if final reporting cannot finish. Checkpoint completed cleanup before fallible freezing. Normal worker freeze remains inside the selected execution deadline, with a separate 30-second reserve after cancellation rather than an unrelated 30-second normal-run cap. Record stages. Unchanged bytes can skip replacement only after regular-file/bounds checks and exact sanitized comparison; changed redaction still replaces and updates hashes.
3. #356: retain partial snapshot/protocol details and command exit/output separately; record guard/setup/before/command/protocol/after/cleanup timings and bounded request identities. Reproduce sustained owned requests and compare unchanged/current dependency baselines before selecting a transport correction. Do not guess a dependency/host cause or mask a deadline with retries or larger limits.
4. #355: validate numeric SQLite versions in both twins, retain bounded original candidate rejection detail and process/runtime metadata, and strengthen fixture prerequisite checks. Exercise existing six-writer cases unchanged. The original CLR signature has no exact matching proven upstream cause; preserve it and report any still-unconfirmed investigation as incomplete rather than claiming an upstream fix.
5. #358: add diagnostic to delivery domain/storage/API/CLI/desktop; migration24 copies old SMTP last_error to diagnostic, clearing last_error only for persisted succeeded SMTP rows. New attempts store neutral diagnostic always; last_error is empty only when all recipients accepted. Preserve webhook history, recipient state machine, immutable queued credentials, terminal credential erasure, retry limits and uncertain-ACK warning. Regenerate existing typed OpenAPI/bindings using project tools.
6. #360: review the actual two-file patch and official Windows cross-drive/UNC/dot-folder fixes, apply exact versions/lockfile, then validate the real site. No unrelated dependency churn.
7. Update docs and schema24 rollback guidance, dated architectural decisions, S112 merge records and v1.6.0 pending release content binding. Canonical eight gates, targeted regressions, full owned QA and complete message review precede commit/publication. New defects get separate accurate records; #354 stays open until actual publication.

## Research and Delegation

/speckit-plan explicitly requests research agents. Three read-only assignments covered #355, #356/#357 and #358; findings are consolidated in research.md. Implementation remains governed by tasks.md and the blocking analyze gate.

## Tooling Note

The generic checklist prerequisite requires plan.md despite this repository's required checklist-before-plan ordering. Its failure was retained; the checklist used already validated absolute feature paths. setup-plan then ran successfully. Re-run full prerequisites for tasks/analyze. PowerShell helper BRANCH is empty in this environment, so direct git verifies codex/113-qa-diagnostics independently; the actual feature selector resolves 113 correctly.

## Complexity Tracking

No constitution violation is requested. Additive report evidence and a forward-only schema column are required to retain accurate history; destructive rewriting, retrying crashed writers and weakened assertions were rejected.

## PR review follow-up (2026-10-23)

Round 1 on PR #361 identified an interruption window between late-redacted channel replacement and refreshed report/progress metadata. Record exact old/new artifact identities before each changed-file replacement, then reconcile only verified recorded identities during recovery. This preserves recovery at every interruption boundary without accepting arbitrary conflicting metadata, changing case outcomes or rerunning commands. The review thread is https://github.com/shruggietech/go-schedule/pull/361#discussion_r4359185167.
