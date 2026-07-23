# Tasks: Run engine benchmarks in CI and enforce the p99 dispatch-latency budget

**Feature**: 009-ci-bench-p99 | **Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

Design inputs: [research.md](research.md), [data-model.md](data-model.md),
[contracts/p99-gate.md](contracts/p99-gate.md), [quickstart.md](quickstart.md).

Test discipline: TDD. US1's test is written to reference the budget constant; the
constant (Phase 2) is the blocking prerequisite that makes it compile and pass.
The negative check in Phase 6 proves the gate genuinely fails on a breach.

---

## Phase 1: Setup

- [X] T001 Confirm the working tree is clean on `main` and feature 009 context is active (`git status`; `.specify/feature.json` points at `specs/009-ci-bench-p99`). No code changes.

## Phase 2: Foundational (blocking prerequisite for US1)

- [X] T002 Add the exported budget constant `const DispatchLatencyBudget = 100 * time.Millisecond` to `internal/engine/engine.go`, placed near the worker-pool/dispatch code, with a doc comment citing constitution Principle IV (p99 dispatch latency < 100 ms; the budget lives next to the code it governs). Ensure `time` is imported (it already is). No behavior change.

## Phase 3: User Story 1 — A latency regression fails the build (Priority: P1) 🎯 MVP

**Goal**: A committed, deterministic check asserts the p99 dispatch latency against the budget and fails when it is exceeded.

**Independent test**: `go test -run TestDispatchLatencyP99 ./internal/engine/...` passes with the observed p99 far under 100 ms; artificially breaching the budget (Phase 6 negative check) makes it fail.

- [X] T003 [US1] Create `internal/engine/latency_test.go` (package `engine`) with a `timingRunner` implementing the engine's `Runner` interface: capture `start := time.Now()` at the top of `Run`, send `start.Sub(scheduledFor)` on a buffered `chan time.Duration`, and return a success `domain.Run` with `StartedAt`/`EndedAt` set to `start` (mirroring `noopRunner` in `engine_bench_test.go`).
- [X] T004 [US1] In the same file, implement `TestDispatchLatencyP99`: open an in-memory store, create one `OverlapAllowConcurrent` recurring task, build `eng := New(st, clock.NewReal(), timingRunner{...}, testLogger(), 8)`, set `eng.runCtx = context.Background()`, and signal completion via `eng.SetOnRun`. Dispatch `N = 2000` runs serially (wait on the completion channel each iteration), drain the `N` latency samples, `sort.Slice` ascending, compute `p99 := samples[int(0.99*float64(N))]`, `t.Logf` the observed p99, and `t.Fatalf` (naming observed p99 vs `DispatchLatencyBudget`) if `p99 >= DispatchLatencyBudget`.
- [X] T005 [US1] Run `go test -run TestDispatchLatencyP99 -v ./internal/engine/...` and confirm it passes and logs a p99 in the microsecond-to-low-millisecond range (validates the assertion, sample count, and non-flakiness at spec altitude).

**Checkpoint**: US1 is independently complete — the budget is enforced by a committed test that runs locally and in the existing race job.

## Phase 4: User Story 2 — Benchmark numbers visible on every CI run (Priority: P2)

**Goal**: CI executes the engine benchmarks and publishes their output as an artifact.

**Independent test**: A CI run shows a `bench` job that runs the benchmarks and produces a downloadable artifact.

- [X] T006 [US2] Add a `bench` job to `.github/workflows/ci.yml` matching house style: `runs-on: ubuntu-latest`, `env: CGO_ENABLED: "0"`, steps `actions/checkout@v4` → `actions/setup-go@v5` (`go-version-file: go.mod`, `cache: true`) → run `go test -run '^$' -bench . -benchmem ./internal/engine/... | tee bench.txt` → `actions/upload-artifact@v4` uploading `bench.txt`. Add a comment noting the job is an informational signal; the enforced gate is `TestDispatchLatencyP99`.
- [X] T007 [US2] Validate the exact command locally: `go test -run '^$' -bench . -benchmem ./internal/engine/...` runs `BenchmarkDispatch` and `BenchmarkNextRun` and runs no test bodies.

**Checkpoint**: US2 complete — benchmarks are a live per-run signal (observable on CI after the authorized push).

## Phase 5: User Story 3 — Goroutine-leak guarantee stays under -race (Priority: P3)

**Goal**: Confirm (no new code) the leak test runs under the race detector in CI.

**Independent test**: `test/integration` is not in the CI race exclusion list, so `TestEngine_NoGoroutineLeak` runs under `-race`.

- [X] T008 [US3] Verify `test/integration` is included in the CI race command in `.github/workflows/ci.yml` (`go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')` — only the two GUI packages are excluded). Record this as verified in the feature notes; make no code change.

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T009 Negative check (scratch, revert after): temporarily set `DispatchLatencyBudget` to a tiny value, run `go test -run TestDispatchLatencyP99 ./internal/engine/...`, confirm it FAILS naming observed p99 vs budget, then revert the constant. Proves the gate is real (SC-003).
- [X] T010 Update `CHANGELOG.md` `[Unreleased]`: an **Added** line for the p99 dispatch-latency test and the CI `bench` job (closes #14); a **Changed** line for the pinned `.github/workflows/ci.yml` bench job; and a dated **Decisions** entry recording the absolute-p99-budget-over-benchstat-delta choice with its rationale (constitution Principle IV permits the alternative with recorded justification).
- [X] T011 Run full CI parity in the foreground, watched to completion: `gofmt -l internal cmd test`; `go vet ./...`; `golangci-lint run ./...`; `CGO_ENABLED=1 go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')`; `go test ./gui/...`; `sh scripts/coverage-gate.sh`. If no C toolchain is present, state the `-race` gate cannot run locally and rely on CI; do not report it as passing.

---

## Dependencies

- T002 (budget constant) blocks T003–T005 (the test references it).
- US1 (T003–T005) is the MVP and is independent of US2/US3.
- US2 (T006–T007) is independent of US1; touches only `ci.yml`.
- US3 (T008) is verification-only, independent.
- Phase 6 (T009–T011) runs after US1/US2/US3; T009 depends on the test (T004);
  T010 records the ci.yml decision (after T006); T011 is the final gate.

## Parallel opportunities

- After T002: US1 (T003–T005) and US2 (T006–T007) touch disjoint files
  (`internal/engine/latency_test.go` vs `.github/workflows/ci.yml`) and can
  proceed in parallel. T008 (US3) is independent of both.

## MVP scope

User Story 1 alone (T001–T005) delivers the core promise: the p99 budget becomes
enforced by a committed test. US2 (visibility) and US3 (confirmation) round out
issue #14.

## Format validation

All tasks use `- [ ] Txxx [P?] [Story?] description with file path`; setup,
foundational, and polish tasks carry no story label; US phase tasks carry [US1]/
[US2]/[US3]. File paths are explicit.
