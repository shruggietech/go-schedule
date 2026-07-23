# Quickstart / Validation: p99 latency gate + CI benchmarks

Runnable checks that prove the feature works end to end. Run from the repo root.

## Prerequisites

- Go toolchain matching `go.mod` (1.25.x; `GOTOOLCHAIN=auto` upgrades inside the
  repo). No C toolchain required for these checks.

## 1. The p99 gate passes and reports a sane number

```bash
go test -run TestDispatchLatencyP99 -v ./internal/engine/...
```

**Expected**: `PASS`, with a logged line reporting the observed p99 (a small
number — microseconds to low single-digit milliseconds), far under the 100 ms
budget.

## 2. The benchmarks run locally (mirrors the CI bench job)

```bash
go test -run '^$' -bench . -benchmem ./internal/engine/...
```

**Expected**: `BenchmarkDispatch` and `BenchmarkNextRun` both execute and print
`ns/op` and `B/op` lines; no test bodies run.

## 3. The gate actually fails when the budget is breached (negative check)

This is a scratch experiment, not a committed change. Temporarily lower the budget
in `internal/engine/engine.go`, e.g. `DispatchLatencyBudget = 1 * time.Nanosecond`,
then:

```bash
go test -run TestDispatchLatencyP99 ./internal/engine/...
```

**Expected**: `FAIL`, with a message naming the observed p99 versus the (tiny)
budget — demonstrating the gate is real. **Revert the constant afterward.**

## 4. Full CI parity (must be green before the pre-push halt)

```bash
gofmt -l internal cmd test
```
```bash
go vet ./...
```
```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```
```bash
CGO_ENABLED=1 go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')
```
```bash
go test ./gui/...
```
```bash
sh scripts/coverage-gate.sh
```

**Expected**: `gofmt` prints nothing; vet and lint clean; the race run (which
includes `TestDispatchLatencyP99` and the goroutine-leak test) passes; GUI tests
pass; coverage gate stays ≥ 80 % on the six core packages. If no C toolchain is
present, the `-race` run cannot execute locally — say so explicitly and rely on CI
for that gate.

## 5. CI benchmark artifact (observable only after the authorized push)

After the push, the `bench` job in `.github/workflows/ci.yml` runs the engine
benchmarks and uploads their output as a downloadable build artifact on the run.

## Reference

- Test + CI-job contract: [contracts/p99-gate.md](contracts/p99-gate.md)
- Conceptual entities: [data-model.md](data-model.md)
- Design decisions: [research.md](research.md)
