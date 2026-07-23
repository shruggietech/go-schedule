# Feature Specification: Run engine benchmarks in CI and enforce the p99 dispatch-latency budget

**Feature Branch**: `009-ci-bench-p99` (trunk-based: committed directly to `main`)

**Created**: 2026-07-23

**Status**: Draft

**Input**: Close GitHub issue #14 — run the engine benchmarks in CI and enforce
the constitution's p99 dispatch-latency budget.

## Overview

The project's Performance principle sets a hard budget — **p99 dispatch latency
< 100 ms under nominal load** — and requires continuous verification. That budget
is written down and believed but never measured: the engine has benchmarks that
no pipeline runs, and the one dispatch benchmark reports a mean, not the p99 the
budget is stated in. This feature closes the gap between the stated budget and a
committed, automatically-checked measurement, and makes the existing benchmarks a
live signal rather than dead files.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A latency regression fails the build (Priority: P1)

A maintainer changes the dispatch path in a way that inflates scheduling
overhead. Before this feature the regression is invisible until an operator
notices late runs in production. After it, an automated check that measures the
p99 dispatch latency against the budget fails, so the regression is caught at
commit/CI time instead of in the field.

**Why this priority**: This is the feature's core promise — turning a documented
budget into an enforced one. Without it, nothing else here matters.

**Independent Test**: Run the p99 latency check against the current engine; it
passes with the observed p99 far under the budget. Artificially inflate the
dispatch overhead in a scratch experiment and confirm the same check fails. The
check is deterministic (does not depend on real `time.Sleep` for correctness of
the assertion) and self-contained (no external services).

**Acceptance Scenarios**:

1. **Given** the engine at its current dispatch performance, **When** the p99
   latency check runs, **Then** it reports the observed p99 and passes because it
   is within the documented budget.
2. **Given** a change that pushes p99 dispatch latency over the budget, **When**
   the check runs, **Then** it fails and names the observed p99 versus the budget.
3. **Given** the check runs on a shared/loaded CI runner, **When** it executes
   repeatedly, **Then** it does not flake — it passes every time the engine is
   within budget, because the assertion carries generous headroom.

---

### User Story 2 - Benchmark numbers are visible on every CI run (Priority: P2)

A maintainer reviewing a change wants to see how the engine's dispatch and
next-run computation are performing, without running anything locally. After this
feature, each CI run executes the engine benchmarks and publishes their output as
a retrievable artifact, so the numbers are reviewable rather than un-run.

**Why this priority**: Converts the two existing benchmark functions from
"documentation with a compiler check" into a live per-run signal, and provides
the raw trend a maintainer can eyeball for a slowdown that is still within budget.

**Independent Test**: Trigger a CI run and confirm a benchmark job executes the
engine benchmarks and produces a downloadable artifact containing their output.

**Acceptance Scenarios**:

1. **Given** a push or pull request, **When** CI runs, **Then** a benchmark job
   executes the engine benchmarks and completes successfully.
2. **Given** the benchmark job has run, **When** a maintainer inspects the run,
   **Then** the benchmark output is available as a build artifact.

---

### User Story 3 - The goroutine-leak guarantee stays enforced under the race detector (Priority: P3)

A maintainer needs confidence that the engine's goroutine-termination guarantee
remains checked under the race detector. This story is confirmation, not new
work: the existing leak test already runs under the race detector in CI, and this
feature records that fact so issue #14's fourth item is closed by observation.

**Why this priority**: It is a verification/documentation item with no code
change; it rounds out issue #14 without inventing redundant machinery.

**Independent Test**: Confirm the goroutine-leak test is included in the
race-detector test selection (it is not among the excluded packages) and runs
under `-race`.

**Acceptance Scenarios**:

1. **Given** the CI race-detector test selection, **When** it runs, **Then** the
   goroutine-leak test executes under the race detector (it is not excluded).

---

### Edge Cases

- **Loaded CI runner / timer noise**: The p99 assertion must not flake under
  scheduler contention or coarse timer granularity. Mitigated by asserting an
  absolute budget with large headroom (real overhead is microscopic relative to
  the 100 ms ceiling), not a tight relative delta.
- **Sample-count sensitivity**: The p99 must be computed over enough samples to
  be meaningful yet run fast enough for the normal test suite; the sample count
  is a fixed, documented value.
- **Environment without a C toolchain**: The p99 check and the benchmark run must
  not require cgo, so they run on any machine and in the cgo-free CI jobs.
- **A pre-existing budget breach**: If the engine were already over budget, the
  check would fail immediately on introduction — an acceptable and correct
  outcome, surfacing a real defect rather than masking it.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide an automated, committed check that measures
  the p99 of dispatch latency (the interval from a run's scheduled time to the
  start of its execution, excluding command execution time) over a fixed number
  of dispatched runs and asserts it against the documented budget.
- **FR-002**: The dispatch-latency budget MUST be a single named value recorded
  next to the engine code it governs, so the budget and the code cannot drift
  apart, and the check MUST assert against that same value.
- **FR-003**: The p99 check MUST fail when the observed p99 exceeds the budget and
  MUST report the observed p99 alongside the budget in its failure message.
- **FR-004**: The p99 check MUST be deterministic and non-flaky on shared CI
  hardware: it MUST NOT depend on real sleeps for the correctness of its
  assertion and MUST carry enough headroom to pass reliably whenever the engine
  is within budget.
- **FR-005**: The p99 check MUST run as part of the standard automated test suite
  (so it executes locally and in the existing test pipeline) and MUST NOT require
  a C toolchain.
- **FR-006**: Continuous integration MUST execute the engine benchmarks on each
  run and MUST publish their output as a retrievable build artifact.
- **FR-007**: The benchmark execution in CI MUST run the benchmarks only (not the
  functional test suite a second time) and MUST NOT require a C toolchain.
- **FR-008**: The feature MUST confirm and record that the existing
  goroutine-leak test runs under the race detector in CI, and MUST NOT exclude it
  from the race-detector selection.
- **FR-009**: The enforced regression gate MUST be the absolute p99 budget
  assertion; the feature MUST record the decision to enforce the absolute budget
  rather than a relative benchmark-delta gate, with the rationale, in the
  project's changelog decisions.
- **FR-010**: Any change to a pinned process artifact (e.g. a CI workflow) made by
  this feature MUST be recorded as a dated decision in the changelog.
- **FR-011**: The feature MUST NOT alter the engine's runtime dispatch behavior or
  add latency logging/metrics to production execution paths; it is measurement and
  CI wiring only.

### Key Entities

- **Dispatch latency sample**: One measurement of the interval between a run's
  scheduled time and the moment its execution starts, for a single dispatched
  run, with command execution excluded.
- **Dispatch-latency budget**: The single documented ceiling (p99 < 100 ms) that
  the check asserts against, recorded next to the engine code.
- **Benchmark artifact**: The published output of the engine benchmarks from a CI
  run, retained for review.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A committed, automatically-run check asserts the p99 dispatch
  latency against the 100 ms budget; before this feature, zero such checks exist.
- **SC-002**: On the current engine, the check reports an observed p99 at least an
  order of magnitude under the budget and passes.
- **SC-003**: A dispatch-path change that pushes p99 over the budget causes the
  check to fail (demonstrable in a scratch experiment), so the budget is enforced
  rather than merely documented.
- **SC-004**: Every CI run executes the engine benchmarks and produces a
  downloadable artifact of their output; before this feature, no CI run executes
  any benchmark.
- **SC-005**: The p99 check passes on repeated runs with no flakiness observed,
  including under the race detector.
- **SC-006**: The decision to enforce an absolute budget over a relative
  benchmark-delta gate is recorded in the changelog with its rationale.

## Assumptions

- The relevant "nominal load" for the budget is the engine's scheduling overhead
  in isolation (command execution excluded), matching how the existing dispatch
  benchmark is constructed and how the budget is phrased ("scheduled time →
  execution start").
- Asserting the absolute p99 budget satisfies the constitution's
  "benchmark regression checks for performance-sensitive packages" requirement: a
  change that regresses dispatch latency past the budget fails the check. The
  constitution's relative-10 %-delta clause explicitly permits an alternative
  "with explicit, recorded justification"; this feature records that justification
  rather than adopting a noise-prone relative gate.
- A fixed sample count in the low thousands yields a meaningful p99 while keeping
  the check well under a second, so it fits the normal test suite.
- The existing benchmark functions (dispatch and next-run computation) are the
  correct set to run in CI; no new benchmark functions are required by this
  feature beyond the p99 measurement.
- Publishing benchmark output as an artifact (rather than gating on a stored
  baseline) is sufficient for the trend-signal goal, given the absolute budget is
  the enforced gate.
