# Quickstart: Post-release dependency refresh

## Graph Restoration

```bash
go mod tidy
go mod verify
(cd desktop && go mod tidy && go mod verify)
(cd desktop/frontend && npm ci && npm audit --audit-level=high)
git diff --exit-code -- go.mod go.sum desktop/go.mod desktop/go.sum desktop/frontend/package.json desktop/frontend/package-lock.json
```

## Focused Verification

```bash
go test -race ./internal/mcphttp ./internal/mcpobserve ./internal/remote ./internal/remoteenroll ./internal/api/client ./internal/api/server
(cd desktop && go test -race ./...)
(cd desktop/frontend && npm test && npm run build && npm run test:e2e)
go test ./test/integration -run 'TestProductionDesktopUsesStableWailsIdentity|TestReleaseWorkflowBuildsProductionWailsPayload|TestWindowsInstallerGUIResourceContract' -count=1
```

## Canonical Verification

```bash
sh scripts/verify.sh all
go run ./scripts/github-format
git diff --check
```

## Version Inventory

```bash
go list -m github.com/modelcontextprotocol/go-sdk golang.org/x/crypto golang.org/x/sys golang.org/x/time
(cd desktop && go list -m golang.org/x/crypto golang.org/x/sys)
(cd desktop/frontend && npm ls --depth=0 react react-dom @types/react @types/react-dom @types/node vite)
```
