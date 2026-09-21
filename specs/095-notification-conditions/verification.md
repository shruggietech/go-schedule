# S095 Verification: Actionable Notification Conditions

## Outcome

S095 completes issue #175 with durable, explainable notification conditions for consecutive failures, process start failure, duration thresholds, recovery, bounded reminders, success quiet periods, and opt-in daemon healthy-presence heartbeats.

## Acceptance evidence

- Schema v20 preserves existing assignments with a failure threshold of one and adds durable per-task and per-channel condition state.
- The run transaction evaluates conditions after history is recorded, so notification creation cannot change the task outcome or block scheduler dispatch.
- One primary problem is selected in deterministic order: failure to start, consecutive failure, then duration exceeded.
- The first problem transition notifies once, identical repeats are suppressed, an elapsed reminder can notify again, and optional recovery fires when the active problem clears.
- Success quiet periods suppress routine repeat volume without hiding the first problem or recovery transition.
- Policy fingerprints reset stale counters and active conditions when an assignment changes. Durable state survives daemon restart and evaluates only runs present in durable history.
- Enabled channels may opt into healthy-presence heartbeats with a receiver-visible next expected deadline. Zero disables heartbeat work, and a stopped daemon makes no false claim that it can send its own outage.
- API, CLI, MCP Manage-compatible assignment input, desktop channel and policy controls, webhook payloads, and redacted delivery history carry the complete condition configuration and safe explanation.
- Desktop controls use labeled native inputs, remain inside the existing responsive policy layout, and preserve the progressive-disclosure workflow.

## Commands and results

- `go test ./...` at the repository root: passed, including integration tests.
- `go test ./...` in `desktop`: passed.
- `npm run build` in `desktop/frontend`: TypeScript and Vite production build passed.
- `npm test -- --run` in `desktop/frontend`: 24 files and 135 tests passed.
- `npm run test:e2e -- --grep "notifications"` in `desktop/frontend`: 2 Playwright workflows passed with accessibility coverage.
- `go test -race ./internal/store ./internal/notification ./internal/api/server`: passed.
- `go run ./scripts/github-format`: passed with no em dashes or hard-wrapped Markdown prose.
- `scripts/verify.sh all` through WSL Bash with explicit Windows Go tool paths: format, vet, lint, full race, desktop race, native Windows Wails production build, frontend tests and build, coverage, documentation, and automation gates passed.

## Compatibility and privacy

Existing assignment JSON remains valid and gains threshold-one failure semantics by default. Webhook additions are optional fields under the existing schema identifier. No endpoint, authorization value, task command, arguments, environment, stdin, output, or other protected task content enters condition summaries, desktop projections, or terminal delivery evidence.

## Traceability

- Closes #175.
- Advances parent #19 while leaving sibling transport work independent.
