# Quickstart: Agent Access Controls and MCP Release Gates

## Prerequisites

- Build the daemon, CLI, and Wails desktop from the S071 branch.
- Start the daemon under an identity that can use the protected local IPC endpoint.
- Keep any one-time credential out of command history, logs, screenshots, and issue text.

## Desktop Off state

1. Open Agent Access.
2. Confirm stdio is described as an on-demand command that opens no listener.
3. Confirm localhost HTTP reports Off and offers client name, port, and optional origin inputs.
4. Confirm Observe is Available while Operate and Manage are visibly unavailable future capabilities.

## Named localhost client lifecycle

1. Enter a client name, free port, and any required numeric-loopback browser origins.
2. Enable localhost access and confirm the message says the credential was copied once without rendering its value.
3. Configure a local Streamable HTTP client with the displayed endpoint and copied bearer credential.
4. Perform one Observe request, refresh Agent Access, and confirm the request count and last-access time appear.
5. Rotate access, confirm a replacement credential is copied and prior evidence resets, then prove the old credential fails.
6. Choose `Revoke and turn off` and confirm the listener reports Off and the prior credential fails.

## Supported host setup

Follow [Local MCP access](../../docs/mcp.md) for Codex stdio, generic command-based stdio, or generic Streamable HTTP configuration. Remove stdio configuration in the host itself; revoke HTTP from Agent Access or `gosched mcp http disable`.

## Focused validation

```bash
go test -race ./internal/mcpobserve ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli ./test/integration
cd desktop && go test -race ./...
cd desktop/frontend && npm test && npm run build && npm run test:e2e -- agent-access.spec.ts
```

Expected outcomes are exact resource and template counts, zero tools, safe schema and hostile-content results, real built-command initialization, clean shutdown, fail-closed clipboard rollback, and accessible desktop states.

## Canonical validation

```bash
sh scripts/verify.sh all
```

All eight gates must pass in the foreground before publication. Signed and exact-candidate installed-artifact qualification remains part of #190.
