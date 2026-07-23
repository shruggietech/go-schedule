# Implementation Plan: Run engine benchmarks in CI and enforce the p99 dispatch-latency budget

**Branch**: `009-ci-bench-p99` (trunk-based: committed directly to `main`) | **Date**: 2026-07-23 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/009-ci-bench-p99/spec.md`

## Summary

Close the loop between the constitution's documented p99 dispatch-latency budget
(< 100 ms) and a committed, automatically-checked measurement, and make the
existing but un-run engine benchmarks a live CI signal. Two deliverables: (1) a
committed `TestDispatchLatencyP99` in `internal/engine` that measures per-dispatch
latency (scheduled-time → execution-start, command execution excluded) over a
fixed sample count, computes the p99, and asserts it against a named
`DispatchLatencyBudget` constant defined next to the engine dispatch code; (2) a
`bench` job in CI that runs the engine benchmarks and publishes their output as an
artifact. The enforced regression gate is the absolute p99 assertion; a
noise-prone relative benchstat delta gate is deliberately not adopted (recorded
decision). No engine runtime behavior changes and no latency logging is added to
production paths.

## Technical Context

**Language/Version**: Go 1.25.0 (per `go.mod`; `GOTOOLCHAIN=auto`).

**Primary Dependencies**: Standard library only for the new test (`testing`,
`sort`, `time`, `context`). Reuses `internal/clock`, `internal/domain`,
`internal/schedule`, `internal/store`. CI uses `actions/checkout@v4`,
`actions/setup-go@v5`, `actions/upload-artifact@v4`.

**Storage**: In-memory SQLite (`store.Open(":memory:")`) in the test; N/A in
production for this feature.

**Testing**: `go test` (standard suite). The p99 test is a normal `Test`, so it
runs in the existing `test` (race) job and locally; it is cgo-free.

**Target Platform**: Linux/macOS/Windows (the p99 test runs on all three via the
existing test matrix; the `bench` job runs on `ubuntu-latest`).

**Project Type**: Single Go module (daemon + CLI + GUI); this feature touches
`internal/engine` and CI configuration only.

**Performance Goals**: The property under test is the budget itself — p99 dispatch
latency < 100 ms. Observed overhead is expected to be microseconds.

**Constraints**: The test must be non-flaky on shared CI hardware (absolute budget
with ~4–5 orders of magnitude headroom); it must not depend on real `time.Sleep`
for the correctness of its assertion; it must be cgo-free; it must run well under
~1 s even under `-race`.

**Scale/Scope**: Fixed sample count of **2000** dispatched runs per test
invocation (see research.md for the sizing rationale). Two files changed in
`internal/engine`, one CI workflow, one changelog.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality** — PASS. New code is a single test file plus one exported
  constant with a doc comment citing the principle. `gofmt`/`vet`/lint apply. The
  test's goroutine usage is bounded (it drives the engine's existing worker pool
  and waits on completion; no new long-lived goroutines).
- **II. Testing Standards (NON-NEGOTIABLE)** — PASS and directly advanced. This
  feature *is* a test that enforces a budgeted property. The p99 assertion does
  not depend on real sleeps for correctness (it measures real elapsed latency but
  asserts an absolute ceiling with vast headroom); it runs under `-race` in the
  existing job. No safety-critical surface is weakened.
- **III. UX Consistency** — N/A (no user-facing interface changes). The CI
  artifact and test failure message follow existing conventions.
- **IV. Performance Requirements** — PASS and directly advanced. This is the
  feature that makes the p99 budget measured and enforced, the budget lives next
  to the code it governs (the `DispatchLatencyBudget` constant), and the
  benchmarks run in CI. The decision to enforce the absolute budget rather than a
  relative-10 % delta is recorded with justification (constitution permits the
  alternative "with explicit, recorded justification").
- **V. Autonomous Build-Phase Execution** — PASS. Runs under autopilot for open
  issue #14; one pre-push halt; `/speckit-analyze` gate honored; `ci.yml` change
  recorded as a dated changelog decision.

**Clock discipline note**: `internal/engine` production code must take time through
the injected `Clock`. This feature adds no `time.Now()` to production code — the
`time.Now()` call lives only in the test's `timingRunner`, measuring real elapsed
latency, exactly as the existing `BenchmarkDispatch` already does. No violation.

No gate violations. Complexity Tracking is empty.

## Project Structure

### Documentation (this feature)

```text
specs/009-ci-bench-p99/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (conceptual entities)
├── quickstart.md        # Phase 1 output (validation guide)
├── contracts/
│   └── p99-gate.md      # The test + CI-job behavioral contract
└── checklists/
    ├── requirements.md  # Spec-quality checklist (from /speckit-specify)
    └── performance.md   # Domain checklist (from /speckit-checklist)
```

### Source Code (repository root)

```text
internal/engine/
├── engine.go            # + DispatchLatencyBudget constant (near the pool/dispatch code)
├── engine_bench_test.go # (unchanged) BenchmarkDispatch, BenchmarkNextRun — run by the new CI job
└── latency_test.go      # NEW: timingRunner + TestDispatchLatencyP99

.github/workflows/
└── ci.yml               # + `bench` job (pinned artifact; dated CHANGELOG decision)

CHANGELOG.md             # [Unreleased] Added line + dated Decisions entry
```

**Structure Decision**: Single Go module. The change is localized to
`internal/engine` (one production constant, one new test file) and one CI job.
No new packages, no new production dependencies, no schema or migration.

## Complexity Tracking

*No constitution violations; no entries.*
