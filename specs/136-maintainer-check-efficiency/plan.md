# Implementation Plan: Maintainer check efficiency

**Branch**: `codex/136-maintainer-check-efficiency` | **Date**: 2026-11-08 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #433 and children #434/#435 together using a native checker with an evaluation-scoped file cache, YAML-aware workflow checks and in-process lifecycle evaluation. Preserve the existing shell interfaces and all meaningful fixture decisions.

## Technical Context

**Language/Version**: Repository Go module toolchain and POSIX wrappers.

**Primary Dependencies**: Existing go.yaml.in/yaml/v3; standard library and internal/platform hidden-process support. No added dependency.

**Storage**: Read-only repository files and isolated test fixtures. No persisted cache.

**Testing**: Existing automation/lifecycle/failure-propagation corpus, Go race tests and benchmarks, full foreground scripts/verify.sh all, hosted cross-platform CI.

**Target Platform**: Windows, Linux, macOS.

**Project Type**: Internal maintainer CLI, not an application feature/release.

**Performance Goals**: Reduce isolated native Windows checker durations against retained 31.541-second lifecycle and 44.701-second automation measurements; remove utility fan-out without promising total canonical savings.

**Constraints**: Same policy/gates; hidden noninteractive child processes; no timeouts, dependency or product changes; immutable historical inputs.

**Scale/Scope**: 133 existing specifications plus S136, all workflow files and existing synthetic cases.

## Constitution Check

Pre-design and post-design: PASS. I: idiomatic Go, explicit failures and no new concurrency. II: existing corpus retained, discriminating parser/lifecycle regressions and full race/coverage gates. III: existing entry points, status and useful diagnostics. IV: unchanged-source timing before edits, isolated after timing and benchmarks. V: installed specify/clarify/checklist/plan/tasks/analyze/implement procedures, review branch, explicit push/PR authority and owner merge boundary. No principle waiver is needed.

## Research and decisions

- Choose native Go over shell micro-optimization or another embedded interpreter because measured Windows costs and visible per-rule/per-spec process loops justify reusable in-process work.
- Parse workflow/dependency structure using YAML nodes, retain source lines, reject malformed/duplicate keys and inspect only executable action fields (#434). Use required literal shell-command checks only for shell payload contracts. Do not implement a general shell interpreter.
- Cache each read/parse only inside one evaluation. A later fixture mutation must be observed by a fresh evaluation.
- Lifecycle evaluation stays in-process. Replace the synthetic lifecycle stub with real minimal inventory/spec/task data. Quote invalid synthetic YAML run scalars and repair the overbroad desktop-job deletion fixture while preserving its intended rejection.
- Keep bounded existing brand and gate-manifest helper execution, with native process errors retained and Windows CREATE_NO_WINDOW. Never execute shell payloads from workflows. Reuse one compiled native checker across the fixture corpus; wrapper invocations build through the actual source module rather than the target fixture module.
- Optional after-specify agent-context hook is omitted because the core plan step updates the existing context reference. The research agent completed read-only policy/fixture analysis; implementation and the blocking analyze gate remain the primary agent's work. No Codex review request has occurred.

## Project Structure

- `internal/maintainercheck/`: evaluation cache, structured workflow policy, lifecycle policy, regression tests and benchmarks.
- `scripts/maintainer-check/`: CLI/status/diagnostic adapter.
- `scripts/automation-check.sh`, `scripts/spec-lifecycle-check.sh`: thin compatible wrappers.
- `test/scripts/*check_test.sh`: existing corpus with one compiled checker and valid isolated fixtures.
- `specs/136-maintainer-check-efficiency/`: specification, plan, research, model, CLI contract, checklists, tasks, analysis and verification records.

## Verification and delivery

Retain old source/timing under ignored .test-temp/s136. Capture original failure results separately. Run focused regressions before/after, original fixture decisions, native Windows before/after timing, race/benchmarks and all eight foreground canonical gates. Final source freeze binds evidence; factual completion records are separately checked. Publish only the authorized official PR, watch checks/reviews and resolve findings within two total requests including automatic review. Leave auto-merge disabled and stop before merge/tag/release.

## Exact candidate hosted evidence

CI and CodeQL retain their existing pull-request/main triggers and jobs but explicitly select the PR head for checkout (fallback github.sha for non-PR events). This avoids qualifying a temporary merge checkout as the handoff candidate. The dated pinned-process decision is in CHANGELOG.md.

CI qualification ArtifactOrigin labels use GOSCHED_CI_SOURCE_SHA, defined from the same selected checkout expression. The source-config regression keeps that metadata consistent with CI/CodeQL checkout. #436 preserves exact dependency membership and unique Makefile pin selection; no functional requirement or policy is weakened.

Review correction #438 preserves roots as literal filesystem paths during specification/release-note discovery. Windows wrappers use existing shell builtins for native path output because implicit MSYS conversion leaves bracketed arguments unchanged. This retains the original quoted-root contract and adds no new dependency. The final review cap is exhausted after two requests; further verified corrections do not request another round.
