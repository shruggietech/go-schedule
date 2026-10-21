# Implementation Plan: Unattended CLI evaluation and canonical QA reports

**Branch**: `codex/108-unattended-cli-qa` | **Date**: 2026-10-21 | **Spec**: [spec.md](spec.md)

**Input**: Full #289/#290 bundle under parent #277, baseline main 57461fc.

## Summary

Deliver a repository-owned Go coordinator with one shared versioned contract for scenario inventory, execution events, message corpus, review records and reports. Reuse actual Cobra handlers, the real packaged daemon/CLI and S107 owned fixtures. Native Windows and the selected WSL environment execute sequentially. The executing agent performs semantic communication review between capture and finalization; no model-service integration or fabricated verdicts.

## Technical Context

**Language/Version**: Existing Go 1.26 module and pinned toolchain conventions.

**Primary Dependencies**: Standard library, existing Cobra, jsonschema/v6, x/sys and product packages. No new dependency.

**Storage**: Exclusive private owned runtime roots; bounded redacted artifacts through confined os.Root operations; atomic JSON checkpoints and append-only events.

**Testing**: Contract/adversarial unit tests, actual command and packaged-daemon scenarios, deterministic clock fixtures, native Windows/WSL resource-bound worker integration, canonical foreground sh scripts/verify.sh all.

**Target Platform**: Windows 11 native plus explicitly selected existing WSL2 distro/user; shared contracts compile on the existing Linux/macOS/Windows matrix. Unsupported execution capabilities report blocked.

**Project Type**: Maintainer CLI, separate from installed gosched product commands.

**Performance Goals**: Sequential environments, worker_pool_size 1, GOMAXPROCS 2, GOFLAGS=-p=2, per-case and total deadlines, bounded output/artifacts. Default environment budget: 25% CPU (Windows system denominator, Linux one-core denominator), 1 GiB committed/cgroup memory, Windows 32 processes, Linux 128 tasks/threads; configurable recorded bounded values. Total coordinator deadline 30 minutes with cleanup/reporting reserve, case default 15 seconds, 256 KiB per captured channel and 64 MiB total artifacts.

**Constraints**: Unelevated hidden noninteractive child processes, exact ownership before mutations, native user state untouched, redaction before artifact persistence, independent watchdog and process-tree cleanup. Missing controls cannot be reclassified as pass.

**Scale/Scope**: Complete actual command tree, generated commands/aliases/flags and all feature/negative classes in FR-003/004. Behavioral coverage is distinct from help/metadata coverage.

## Constitution Check

Pre-research and post-design: PASS. I: idiomatic small modules and handled lifecycle/errors. II: behavioral and adversarial tests precede implementation, race/cancellation/security gates retained. III: strict CLI/report states and actionable diagnostics. IV: enforced budgets and bounded I/O; no scheduler hot-path change. V: full installed Spec-Kit sequence, review branch, local verification/commit and one pre-publication halt. No exception or constitution change.

## Project Structure

```text
cmd/gosched-qa/main.go
internal/cli/inventory.go
internal/clifixture/process*.go
internal/clievaluate/
  model.go, schema.go, decode.go, artifacts.go, validate.go
  inventory.go, manifest.go, capture.go, scenarios*.go
  environment*.go, coordinator.go, resume.go
  corpus.go, review.go, render.go
  schemas/v1/*.schema.json
  *_test.go
docs/cli-qa.md
docs/maintainers/cli-qa-run-v1.md
docs/maintainers/cli-qa-triage-v1.md
specs/108-unattended-cli-qa/
```

**Structure Decision**: A focused internal maintainer package and thin standalone coordinator command; no second scheduler, alternate production credential backend or generic test framework.

## Execution and dependency order

1. Freeze concrete Go/JSON schemas, event and corpus contracts; implement semantic validators and redacting bounded artifact sink.
2. Export metadata from the actual initialized Cobra tree without contacting installed services. Bind the checked-in inventory and behavioral manifest; distinguish metadata probes from functional outcomes.
3. Extend S107 launchers with optional environment resource controls and memory capture. Preserve existing callers. Windows suspended-before-assignment jobs add CPU/memory/process limits and low priority. Linux workers run inside a unique transient user-systemd cgroup with CPUQuota/MemoryMax/TasksMax and independent runtime termination; individual cases retain process-group ownership.
4. Implement real packaged execution, scoped host-facing handlers and deterministic engine fixtures for the complete feature groups. Each case asserts relevant independent state/history/audit/log evidence; expected rejection remains pass only if its complete contract holds.
5. Capture/checkpoint every planned case, freeze redacted evidence/corpus, support identity-checked interruption/resume and cross-environment aggregation.
6. Validate executing-agent reviews, derive truthful outcomes, render the seven-section report and document execution/issue-triage prompts.
7. Run the actual Windows/selected-WSL routine and communication review, file independently observed defects, run canonical local gates and commit locally before halt.

## Decisions

1. Reuse existing jsonschema/v6 with Draft2020, asserted formats and embedded-only resources. Strict types plus semantic validation handle invariants JSON Schema cannot express.
2. Keep execution_status separate from effective pending_review status so missing semantic review does not erase functional evidence.
3. Evidence manifest hashes captured artifacts, events and corpus; derived report JSON/Markdown do not recursively hash themselves.
4. Capture bounded streams in memory, discover fixture-only dynamic secrets there, then sanitize all persistent output. Never persist raw secret registries.
5. Use actual CLI metadata export rather than a parallel hand-maintained command tree. Help/error probes do not satisfy functional lifecycle coverage.
6. Reuse source-level deterministic fixtures where packaged wall-clock execution would be slow or imprecise, labeling the exact boundary and subject identity.
7. Use owned Linux cgroups, not RLIMIT_AS (incompatible with Go's address reservations), RLIMIT_RSS (ineffective) or UID-wide RLIMIT_NPROC. No global WSL changes.
8. The inspected existing Ubuntu-24.04 user manager delegates cpu/memory/pids. Verify actual controls inside each worker before cases; unavailable environments remain blocked without installation or elevation.
9. Current logs CLI has severity/limit/metadata only. Stream cancellation/reconnect belongs to the actual API/client fixture boundary; do not invent CLI tail/follow flags.
10. Product defects, including usage exit/extra-argument candidates, retain failing scenario results and separate issue records. Completing the evaluator is independent of making the evaluated application defect-free.
11. Resume must match subject binaries, manifest, runner, budgets, environments, artifact hashes and inactive prior ownership. In-progress cases are not reusable; completed equivalent cases retain evidence.

## Complexity Tracking

No constitution exception. OS-specific boundaries are demonstrated technical necessities for the requested Windows/WSL host and resource contract. All other components share the same simple manifest/event/report model.
