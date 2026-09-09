# Quickstart: Dependency Consolidation Verification

Run from the repository root unless a step changes directories.

## Inspect the Selected Baseline

```sh
go list -m github.com/fsnotify/fsnotify modernc.org/sqlite
(cd desktop && go list -m github.com/wailsapp/wails/v2)
(cd desktop/frontend && npm ls --depth=0)
```

Expected: the direct versions match `data-model.md`, Node runtime and types use major 26, and Vite React plugin 6.1.1 resolves with Vite 8.2.2.

## Verify Reproducible Dependency State

```sh
go mod tidy
go mod verify
(cd desktop && go mod tidy && go mod verify)
(cd desktop/frontend && npm ci)
git diff --exit-code -- go.mod go.sum desktop/go.mod desktop/go.sum desktop/frontend/package.json desktop/frontend/package-lock.json
```

Expected: both Go graphs verify, npm resolves without peer or engine bypass, and restoration produces no manifest drift.

## Run Focused Compatibility Evidence

```sh
go test -race ./internal/store ./internal/watcher ./test/integration
(cd desktop && go test -race ./...)
(cd desktop/frontend && npm test)
(cd desktop/frontend && npm run build)
(cd desktop/frontend && npm run test:e2e)
(cd desktop && go run github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 build -clean)
go test ./test/integration -run 'TestProductionDesktopUsesStableWailsIdentity|TestReleaseWorkflowBuildsProductionWailsPayload|TestWindowsInstallerGUIResourceContract' -count=1
```

Expected: storage, watchers, integration, desktop bridge, frontend unit, type-check, production bundle, browser accessibility and responsive, native Wails, and installer contract evidence passes.

## Run Canonical Verification

```sh
sh scripts/verify.sh all
```

Expected: format, vet, lint, race, GUI, coverage, docs, and automation pass in order.

## Final Audit

```sh
go run ./scripts/github-format
git diff --check
git status --short
```

Expected: publication formatting and diff integrity pass, with only intentional S073 changes present. Pull requests #201 through #208 and issue #215 remain open until merge.
