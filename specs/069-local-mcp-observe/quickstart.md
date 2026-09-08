# Quickstart: Local Observe-Only MCP

## Focused development verification

```bash
go test -race ./internal/mcpobserve ./internal/cli ./cmd/gosched
```

## Build and manual protocol smoke test

```bash
go build -o ./dist/gosched ./cmd/gosched
./dist/gosched mcp serve
```

The second command expects newline-delimited MCP JSON-RPC on stdin and writes protocol responses to stdout. Launch it through an MCP host for normal use. Closing stdin must end the process.

## Publication checks

```bash
go run ./scripts/github-format
sh scripts/verify.sh all
```

Run the canonical verification in the foreground. Record exact evidence in `verification.md`, then scan changed text for Unicode em dashes, mojibake, accidental hard wrapping, secret canaries, and untrusted content that lacks a trust label.
