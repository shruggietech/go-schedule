# Production desktop application

`desktop/` is the production Wails control-center module introduced by S061 and promoted to the sole maintained desktop application by S066. It remains a separate Go module so native WebView dependencies do not enter the cgo-free daemon and CLI module.

Current installers and release workflows build this module while preserving the established `gosched-gui` executable and application-bundle identity.

## Build and test

Requirements are Go 1.25, Node.js 24, npm, the platform WebView development prerequisites required by Wails v2.14.0, and Chromium for the optional browser contract.

```text
cd desktop/frontend
npm ci
npm test
npm run build
npm run test:e2e

cd ..
go test -race ./...
go run github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 build -clean
```

Linux native builds also use the `webkit2_41` build tag and require GTK 3 plus WebKitGTK 4.1 development packages. Generated frontend bindings, dependency directories, browser reports, and native build outputs are ignored.

## Trust boundary

The Go connection manager owns health negotiation, automatic and manual retry, event streaming, generation rejection, and shutdown. `NewLocalBackend` uses the existing Unix socket or Windows named pipe client for This computer. It does not start a listener and it does not add authentication, remote targets, or network transport.

React receives only the documented connection snapshot, action result, and `desktop:event` payload. Endpoint names, credentials, raw transport errors, backend values, and configuration paths stay on the Go side. Browser-only tests use an unavailable adapter and never claim a daemon connection.

## Dependency and asset attribution

The Go module pins Wails v2.14.0 and uses the root module through a local replacement. The frontend pins React 19.1.0, React DOM 19.1.0, TypeScript 5.6.3, Vite 7.3.6, Vitest 5.0.0, Playwright 1.63.0, Testing Library 16.3.3, axe-core 4.13.0, and the exact transitive dependency lock in `frontend/package-lock.json`. Dependency licenses remain available through their upstream packages and module caches and are audited with `go-licenses` and `npm audit` in CI.

The local Geist and Space Grotesk fonts use their Open Font License files under `brand/fonts/licenses/`. The mark and application icon are canonical repository assets registered in `brand/repository-consumers.json`. The runtime makes no CDN, telemetry, analytics, remote font, remote image, script, or stylesheet request.

## Feature contract

Desktop workflows use `connection.Backend` on the Go side and `DesktopBridge` plus the shared `components/` catalog on the React side. They do not call Wails-generated transport bindings directly for daemon data.
