# Phase 1 Data Model: p99 dispatch-latency gate

This feature is measurement and CI wiring; it introduces **no persisted schema,
no migration, and no new domain type**. The "entities" below are conceptual, used
by the test and CI only.

## Conceptual entities

### DispatchLatencySample
- **Represents**: one measurement of dispatch overhead for a single run.
- **Value**: `time.Duration` = `executionStart − scheduledFor`.
  - `scheduledFor`: the `time.Time` passed to `Engine.dispatch`.
  - `executionStart`: `time.Now()` captured at the top of the runner's `Run`.
- **Lifecycle**: produced by `timingRunner.Run`, sent on a buffered channel,
  drained by the test after all dispatches complete. Not persisted.
- **Constraints**: non-negative in practice; command execution excluded.

### DispatchLatencyBudget
- **Represents**: the single documented ceiling the p99 is asserted against.
- **Value**: `const DispatchLatencyBudget = 100 * time.Millisecond` in
  `internal/engine/engine.go`.
- **Relationships**: cited by the doc comment to constitution Principle IV;
  referenced by `TestDispatchLatencyP99`.
- **Rules**: changing it changes the enforced gate; it must remain the same value
  the constitution documents unless the constitution is amended.

### BenchmarkArtifact
- **Represents**: the published output of the engine benchmarks from one CI run.
- **Value**: the captured stdout of
  `go test -run '^$' -bench . -benchmem ./internal/engine/...`.
- **Lifecycle**: produced by the `bench` CI job, uploaded via
  `actions/upload-artifact@v4`, retained per the repo/Actions default retention.
- **Rules**: informational only; not gated on.

## Existing types reused (unchanged)

- `domain.Run` — already carries `ScheduledFor time.Time` and `StartedAt
  *time.Time` with the exact "scheduled → start" meaning the sample uses.
- `domain.Task`, `domain.Schedule` — the test creates one recurring task exactly
  as `BenchmarkDispatch` does.
- `engine.Engine` — the test drives its existing `dispatch`/`launch`/`SetOnRun`
  surface; no new fields or methods on the engine.

## State transitions

None. No entity in this feature has persisted state or a lifecycle beyond a single
test invocation or CI run.
