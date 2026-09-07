# Quickstart: Desktop Settings, Information, and Recovery

## Focused backend verification

```bash
go test ./desktop/settings ./desktop
go test -race ./desktop/settings ./desktop
```

## Focused frontend verification

```bash
cd desktop/frontend
npm test -- --run
npm run build
npx playwright test e2e/settings.spec.ts
```

## Native desktop verification

```bash
cd desktop
wails build
```

## Canonical repository gate

```bash
sh scripts/verify.sh all
```

## Manual smoke test

1. Start once with each valid legacy appearance and confirm it is migrated and applied.
2. Start with missing, invalid, and malformed legacy preferences and confirm system appearance plus a non-blocking explanation.
3. Change appearance, restart, and confirm the new Wails preference remains authoritative.
4. Restore defaults and confirm only desktop preferences change.
5. Inspect every storage record, verify ownership and removal language, and copy each available path.
6. Stop the daemon and confirm local preferences and About remain usable while daemon storage records show unavailable.
7. Activate source and documentation links and confirm the operating-system browser receives the fixed HTTPS destinations.
8. Open Connections, exercise each diagnosis fixture and Retry, and confirm status remains inline without focus-stealing dialogs.
9. Repeat primary Settings and Connections flows by keyboard at 80, 100, 150, and 200 percent zoom.
