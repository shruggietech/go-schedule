# Quickstart: Desktop Notification Management

## Prerequisites

- Go 1.25 toolchain
- Node.js 24 and npm
- Wails v2 and platform WebView build prerequisites
- A running local daemon using an isolated configuration and database
- A local HTTPS-capable webhook test receiver for attended channel testing

## Focused backend validation

```powershell
Set-Location desktop
go test -race ./notifications ./...
```

Expected: secret-boundary, mapping, validation, mutation refresh, task/group policy, inheritance, and delivery-state tests pass with the race detector.

## Focused frontend validation

```powershell
Set-Location desktop/frontend
npm test -- src/notifications src/App.test.tsx src/accessibility.test.tsx
npm run build
npx playwright test e2e/notifications.spec.ts
```

Expected: channel, policy, history, stale-response, duplicate-action, accessibility, keyboard, 200-record, and zoom scenarios pass, then the production bundle compiles.

## Native validation

```powershell
Set-Location desktop
go run github.com/wailsapp/wails/v2/cmd/wails@v2.11.0 build -clean
```

Expected: the native desktop executable builds with the Notifications route bound.

## Attended workflow

1. Create a channel against the controlled HTTPS receiver and confirm endpoint and authorization inputs clear after save.
2. Send a test and verify both receiver payload and desktop history identify it as a test.
3. Disable the channel, observe the disabled explanation, repair endpoint or authorization through explicit replacement, and re-enable it.
4. Assign failure delivery to a group, inspect a child task's inherited policy, add a direct task override, then clear the override and observe inheritance return.
5. Generate success, retry, sending, and failure evidence and inspect each history state without exposing protected values.
6. Complete the entire workflow by keyboard and inspect at 80, 100, 150, and 200 percent zoom.

## Canonical repository verification

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

Expected: format, vet, lint, race, gui, coverage, docs, and automation gates pass in order.
