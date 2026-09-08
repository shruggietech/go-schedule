# Implementation Plan: Dependable Webhook Notifications

**Branch**: `codex/067-webhook-notifications` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/067-webhook-notifications/spec.md`

## Summary

Deliver issues #158 and #159 as one end-to-end increment: durable webhook channels, task and group assignments with nearest-scope replacement, atomic run-to-delivery creation, a separate bounded delivery runtime, redacted history, local API and CLI management, and receiver documentation. Extend the existing SQLite, engine callback, local HTTP API, and Cobra structures with standard-library HTTP and JSON only. Protect write-only endpoint and authorization values with Windows DPAPI or daemon-only Unix permissions, clear delivery secret snapshots at terminal state, and never expose them through ordinary representations.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library (`net/http`, `net/url`, `encoding/json`), existing `modernc.org/sqlite`, `github.com/google/uuid`, Cobra, `golang.org/x/sys/windows` DPAPI bindings, existing local IPC client/server

**Storage**: Forward-only SQLite migration v15 with notification channels, assignments, and deliveries

**Testing**: Go unit and integration tests, `httptest`, deterministic injected notification clock/delays, installed-daemon scripts, `go test -race`, canonical `scripts/verify.sh all`

**Target Platform**: Windows amd64, Linux amd64, macOS arm64; daemon transport remains local IPC

**Project Type**: Long-running system daemon with shared local API and CLI client

**Performance Goals**: Notification network waits consume zero task-worker slots; pending delivery lookup and retention remain bounded; existing dispatch p99 budget remains below 100 ms

**Constraints**: Three attempts maximum, 5-second request timeout, 1 and 2-second backoff, 1,000 terminal history records, no redirects, no task-secret payload fields, no new third-party dependency, no custom cryptography, clean goroutine shutdown and restart recovery

**Scale/Scope**: One webhook channel kind, success/failure conditions, task and nested-group scopes, local API and CLI management; desktop notification management remains #160

## Constitution Check

- **Code Quality**: PASS. Notification state, persistence, transport, API, and CLI responsibilities remain separate. Exported contracts receive intent-focused documentation and all goroutines are owned by one runtime lifecycle.
- **Testing Standards**: PASS. Tests precede behavior changes, deterministic clocks and injected senders avoid sleep-based assertions, restart and concurrent isolation receive integration coverage, and the race gate remains mandatory.
- **User Experience Consistency**: PASS. CLI follows noun-verb command families with human and JSON output, API errors use the existing envelope, timestamps remain RFC 3339, and delivery logs use stable channel, delivery, task, and run identifiers.
- **Performance Requirements**: PASS. Outbound delivery uses a dedicated bounded worker pool and never acquires the engine task semaphore. Queries, retry count, diagnostic body size, and retained history are bounded.
- **Autonomous Build-Phase Execution**: PASS. S067 has a review branch, spec-kit artifacts, issue traceability, analysis gate, local CI parity, and an authorized PR publication path. Tags and releases remain out of scope.
- **Engineering Constraints**: PASS. The standard library satisfies HTTP, URL, and JSON needs. SQLite already supplies durable transactions. Windows uses the already-pinned system DPAPI binding, while Unix uses daemon-only directory and database permissions; neither path introduces home-grown cryptography.

## Project Structure

### Documentation (this feature)

```text
specs/067-webhook-notifications/
├── checklists/
│   ├── requirements.md
│   └── security-delivery.md
├── contracts/
│   ├── local-api.md
│   └── webhook-v1.schema.json
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
cmd/goschedd/main.go
docs/
├── api.md
├── cli.md
└── notifications.md
internal/
├── api/
│   ├── client/
│   └── server/
├── cli/
├── domain/
├── engine/
├── events/
├── notification/
└── store/
test/
├── integration/
└── scripts/
CHANGELOG.md
```

**Structure Decision**: Add one `internal/notification` package for outbound validation, payload construction, retries, and runtime ownership. Keep durable queries and transactions in `internal/store`, plain entities in `internal/domain`, local interface handlers in the established API and CLI packages, and use the engine's post-commit callback only to wake notification work. This avoids coupling untrusted network transport to scheduler execution or introducing another service or broker.

## Delivery Phases

1. Define migration, domain, redaction, payload, and local API contracts with failing tests.
2. Implement channel and assignment persistence, effective-policy resolution, and atomic notification creation beside run persistence.
3. Implement the independently bounded webhook dispatcher with retry, recovery, secret clearing, pruning, and task-worker isolation.
4. Add local API, shared client, CLI lifecycle and evidence operations, plus installed-daemon coverage.
5. Publish receiver, security, retry, duplicate, and precedence documentation; run focused and canonical verification.

## Post-Design Constitution Re-check

PASS. Phase 1 design preserves the simple root-module architecture, adds no dependency, uses one durable database transaction at the run boundary, isolates network work from task workers, bounds every retry and retained collection, and makes the platform-specific secret-protection boundary explicit. No constitutional deviation is required.

## Complexity Tracking

No constitution violations require justification.
