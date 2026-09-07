# Quickstart: Production Wails Shell and Local Connection

## Prerequisites

- Go 1.25.0
- Node.js 24 and npm
- Platform Wails prerequisites documented by S060
- A clean checkout on the S061 review branch

## Connection tests

```powershell
Set-Location desktop
go test -race ./...
```

Expected: safe mapping, compatibility, all state transitions, retry cadence, manual retry coalescing, stale generation rejection, event degradation and recovery, shutdown cancellation, and 100 lifecycle cycles pass without race findings.

## Frontend tests

```powershell
Set-Location desktop/frontend
npm ci
npm audit --audit-level=high
npm run test
npm run build
npm run test:e2e
```

Expected: every shared primitive, shell route, appearance, connection state, keyboard path, dialog focus rule, offline asset rule, compact layout, 200 percent zoom case, and axe serious or critical rule passes.

## Native build

```powershell
Set-Location desktop
go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean
```

On Linux, add `-tags webkit2_41`. Expected: a production desktop foundation binary is generated without changing current release inputs.

## Real local connection

Start `goschedd` through the normal development or installed service path, then run the Wails foundation:

```powershell
Set-Location desktop
go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 dev
```

Expected: This computer becomes connected within two seconds, the sanitized daemon version and capabilities appear, one event stream is active, and stopping the daemon produces degraded or unavailable guidance followed by bounded recovery. No network listener is created by the desktop connection; the development asset server is a Wails tooling exception and is not present in a production build.

## Canonical repository verification

From the repository root:

```bash
sh scripts/verify.sh all
```

Expected: format, vet, lint, race, GUI, coverage, docs, and automation pass in order. Hosted CI additionally builds `desktop/` on Windows, macOS, and Linux and exercises its browser accessibility contract.

## Evidence boundary

S061 proves a production-intent shell and local connection foundation. It does not ship Wails, migrate feature workflows or preferences, change installers, retire Fyne, or claim attended native qualification.
