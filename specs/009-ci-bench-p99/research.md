# Phase 0 Research: p99 dispatch-latency gate + CI benchmarks

All Technical Context unknowns are resolved below. There were no open
`NEEDS CLARIFICATION` markers; the items here record the design decisions and
their rationale.

## D1 — What to measure, and how

**Decision**: Measure per-dispatch latency as `executionStart − scheduledFor`,
where `scheduledFor` is the timestamp handed to `Engine.dispatch` and
`executionStart` is `time.Now()` captured at the top of the runner's `Run`. A
`timingRunner` reports each sample on a buffered channel; command execution is
excluded (the runner does no work).

**Rationale**: This is exactly the constitution's phrasing — "job dispatch
latency (scheduled time → execution start)". It isolates scheduling overhead (the
queue → goroutine → semaphore → runner-start path in `engine.launch`) from command
execution, matching how `BenchmarkDispatch` already isolates it. The
`domain.Run` model already carries `ScheduledFor` and `StartedAt` with this exact
meaning, and the real executor already stamps `StartedAt = time.Now().UTC()`.

**Alternatives considered**:
- *Read latency from persisted runs after the fact.* Rejected: adds store
  round-trips into the measurement and couples the test to persistence timing.
- *Add latency instrumentation to production `launch`.* Rejected: out of scope
  (FR-011) — the feature must not add logging/metrics to production paths. A test
  runner captures the same number without touching production code.

## D2 — Assert an absolute budget, not a relative delta

**Decision**: The enforced gate is `p99 < DispatchLatencyBudget` (100 ms). No
benchstat percentage-delta gate against a stored `bench.txt` baseline is added.
The benchmarks are instead run in CI and their output published as an artifact
for trend inspection.

**Rationale**: The constitution's Performance principle budgets an absolute p99
(< 100 ms); a change that regresses past it fails the check, which *is* a
benchmark regression check. A relative 10 % gate on shared CI runners fires on
scheduler noise, and per issue #14 "a gate that fires on noise gets disabled, and
a disabled gate is worse than none." Principle IV explicitly permits an
alternative to the 10 % clause "with explicit, recorded justification"; this is
that justification, recorded in `CHANGELOG.md` Decisions. The published artifact
preserves the raw numbers for anyone wanting to eyeball a within-budget slowdown.

**Alternatives considered**:
- *benchstat vs committed baseline, fail on >X %.* Rejected for noise/maintenance
  (baseline must be regenerated on legitimate changes; flakes erode trust).
- *Assert on the mean (what `testing.B` reports).* Rejected: measures a different
  property than the budgeted p99; the tail is what hurts a scheduler.

## D3 — Sample count and non-flakiness

**Decision**: Fixed sample count **N = 2000**, dispatched serially (wait on the
completion signal each iteration). p99 index = `int(0.99 * N)` on the sorted
samples (= index 1980, i.e. the 20th-largest of 2000). Assert `p99 < 100 ms`.

**Rationale**: 2000 samples give a stable p99 (the 99th percentile is averaged
over ~20 tail samples) while each dispatch is microseconds, so the whole test runs
in well under a second even under `-race`. Serial dispatch (one in flight at a
time) measures pure scheduling overhead without self-induced queue contention,
which keeps the number representative of "nominal load." The 100 ms ceiling sits
~4–5 orders of magnitude above the real overhead, so runner load or timer
granularity cannot produce a false failure.

**Alternatives considered**:
- *Larger N (e.g. 100k).* Rejected: slower with no material stability gain at this
  headroom.
- *Concurrent dispatch to stress the pool.* Rejected for this gate: it would
  measure contention latency, not the budgeted dispatch overhead, and make the
  number load-dependent and flakier. The pool concurrency is already exercised by
  `BenchmarkDispatch` and the leak test.

## D4 — Where the budget constant lives

**Decision**: `const DispatchLatencyBudget = 100 * time.Millisecond` exported from
`internal/engine/engine.go`, near the worker-pool/dispatch code, with a doc
comment citing constitution Principle IV. The test references the same constant.

**Rationale**: The constitution requires "the budget MUST live next to the code it
governs." Defining it in the production engine file (not the test) documents the
budget beside the dispatch path and guarantees the test asserts against the same
value a maintainer reads. Exported so the test (same package, but exporting also
lets future callers/tools read it) and documentation can reference it.

**Alternatives considered**:
- *Constant in the test file.* Rejected: the budget would not live next to the
  governed code, and a reader of `engine.go` would not see it.

## D5 — CI job shape

**Decision**: Add a `bench` job to `.github/workflows/ci.yml` matching the
existing jobs' house style: `actions/checkout@v4`; `actions/setup-go@v5` with
`go-version-file: go.mod` and `cache: true`; `env: CGO_ENABLED: "0"` (the engine
package is cgo-free); run
`go test -run '^$' -bench . -benchmem ./internal/engine/...`; capture stdout to a
file and upload it via `actions/upload-artifact@v4`. `-run '^$'` runs no tests
(benchmarks only), so the functional suite is not re-run.

**Rationale**: Mirrors `lint`/`coverage`/`build`, which all read the Go version
from `go.mod`, set `CGO_ENABLED` per job, and cache modules. `-benchmem` surfaces
allocations, which the constitution's hot-path guidance cares about. Publishing an
artifact makes the numbers reviewable without gating on them.

**Alternatives considered**:
- *Fold benchmarks into the existing `test` job.* Rejected: that job is the race
  matrix across three OSes; benchmarks belong in a single, clearly-named cgo-free
  job with an artifact, and running benches under `-race` would distort timings.
- *A new `Makefile`/script gate.* Rejected: the Makefile already has a `bench`
  target for local use; the absolute-budget gate is the p99 *test*, which needs no
  separate script. Avoids touching a second pinned artifact.

## D6 — Leak test under -race

**Decision**: No change. `test/integration/leak_test.go` is already selected by
the CI race command (`go test -race $(go list ./... | grep -vE
'/cmd/gosched-gui|/gui$')`) — only the two GUI packages are excluded — so
`TestEngine_NoGoroutineLeak` already runs under `-race`. Recorded as verified;
issue #14 item 4 is closed by observation.

**Rationale**: Adding redundant wiring would be churn for a guarantee that already
holds. The spec (FR-008, US3) frames this as confirmation, not new code.

## Outcome

No unresolved unknowns. The design uses only the standard library and existing
engine test scaffolding, adds one exported constant to production code (no
behavior change), and touches one pinned artifact (`ci.yml`) recorded as a dated
changelog decision.
