# Quickstart: Validate v1.3 Release Qualification

## Prerequisites

- Go toolchain selected from `go.mod`
- C compiler available for race detection
- POSIX shell available for repository verification
- Network access only when the Go module cache is not already populated

## Focused package-shaped journey

```bash
go test -race ./test/integration -run '^TestV13PackageDefaultsRemainOptIn$' -count=1
```

Expected result: the current daemon builds in an isolated candidate directory, starts twice against the same state, remains reachable through local IPC, and reports no notification channels, deliveries, or MCP HTTP listener metadata.

## Focused notification evidence

```bash
go test -race ./internal/notification ./internal/store ./internal/api/server ./internal/api/client ./internal/cli
go test -race ./test/integration -run '^TestWebhookNotificationEndToEndPreservesRunOutcome$' -count=1
```

Expected result: webhook delivery, retries, persistence, migration, redaction, and source-run isolation pass without enabling notifications by default.

## Focused MCP evidence

```bash
go test -race ./internal/mcpobserve ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli
go test -race ./test/integration -run '^TestPackagedMCPCommandDiscovery$' -count=1
```

Expected result: official SDK clients observe the bounded resource contract, no tools appear, authorization and request-boundary tests pass, and process lifecycle tests terminate cleanly.

## Workflow contract

```bash
sh scripts/verify.sh automation
```

Expected result: repository automation verifies the named v1.3 qualification matrix, its three operating systems, and its required test invocation.

## Canonical pre-publication verification

```bash
sh scripts/verify.sh all
```

Expected result: all eight canonical gates pass in order. Record the exact outcome in `specs/072-v13-release-qualification/verification.md`; do not claim a tag or public release artifact.
