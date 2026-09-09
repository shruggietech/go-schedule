# Quickstart: Authenticated Localhost MCP

## Focused verification

```bash
go test ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli ./cmd/goschedd
go test -race ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli ./cmd/goschedd
```

## Manual lifecycle smoke

With the daemon running, choose an unused local port:

```bash
gosched mcp http status
gosched mcp http enable --port 43123
gosched mcp http status
gosched mcp http rotate
gosched mcp http disable
```

Copy the one-time credential from enable or rotate into the local MCP client's Bearer authorization configuration. Restarting the daemon intentionally disables the endpoint and invalidates the copied credential.

## Canonical verification

Run the repository's complete foreground gate:

```bash
sh scripts/verify.sh all
```

Record exact focused and canonical evidence in `verification.md`, run spec-kit analysis, scan changed content for publication and encoding defects, and publish only from a clean review branch.
