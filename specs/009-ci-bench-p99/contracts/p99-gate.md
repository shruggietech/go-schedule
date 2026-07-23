# Contract: p99 latency gate + CI benchmark job

This feature exposes two internal "contracts": a committed test that behaves as a
gate, and a CI job that behaves as a signal. Neither is a user-facing API.

## Contract 1 — `TestDispatchLatencyP99` (package `engine`)

**Location**: `internal/engine/latency_test.go`

**Precondition**: an `Engine` built with the real clock, a `timingRunner`, a
worker pool, and one `OverlapAllowConcurrent` task backed by an in-memory store
(mirrors `BenchmarkDispatch`).

**Behavior**:
- Dispatches `N = 2000` runs serially through `Engine.dispatch`, waiting on the
  `SetOnRun` completion signal between dispatches.
- Collects `N` `DispatchLatencySample` values (`executionStart − scheduledFor`).
- Sorts ascending; computes `p99 = samples[int(0.99*N)]`.
- **Asserts**: `p99 < engine.DispatchLatencyBudget`.

**Postcondition / outputs**:
- On pass: the test logs the observed p99 (via `t.Logf`) for visibility; exit
  success.
- On fail: `t.Fatalf` naming the observed p99 and the budget, e.g.
  `p99 dispatch latency 142ms exceeds budget 100ms`.

**Guarantees**:
- Deterministic pass/fail semantics: correctness of the assertion does not depend
  on real `time.Sleep`; it depends only on measured elapsed latency versus an
  absolute ceiling with large headroom.
- cgo-free; runs in the standard `go test` suite (locally and in the CI `test`
  job, including under `-race`).
- No engine production behavior is invoked differently than in normal dispatch;
  no production code is modified except the added budget constant.

## Contract 2 — CI `bench` job (`.github/workflows/ci.yml`)

**Trigger**: same as the workflow (`push` and `pull_request` to `main`).

**Runner/env**: `ubuntu-latest`, `CGO_ENABLED: "0"`, Go from `go.mod` via
`actions/setup-go@v5` (`cache: true`), `actions/checkout@v4`.

**Behavior**:
- Runs `go test -run '^$' -bench . -benchmem ./internal/engine/...`, capturing
  stdout to a file (e.g. `bench.txt`).
- Uploads that file as a build artifact via `actions/upload-artifact@v4`.

**Guarantees**:
- Executes both `BenchmarkDispatch` and `BenchmarkNextRun` (the current engine
  benchmark set) on every run.
- Runs benchmarks only — `-run '^$'` matches no tests, so the functional suite is
  not re-executed.
- Does not gate the build on the benchmark numbers (informational artifact); the
  enforced gate is Contract 1.

## Non-goals (explicit)

- No benchstat comparison or stored baseline gate.
- No latency logging/metrics added to production `engine` code paths.
- No change to dispatch behavior, worker-pool sizing, or scheduling.
