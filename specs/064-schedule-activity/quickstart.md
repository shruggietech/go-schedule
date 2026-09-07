# Quickstart: Operational Schedule and Activity

## Focused backend verification

```bash
go test ./desktop/operations ./desktop
go test -race ./desktop/operations ./desktop
```

## Focused frontend verification

```bash
cd desktop/frontend
npm test -- --run
npm run build
npx playwright test e2e/operations.spec.ts
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

1. Start the local scheduler daemon and Wails development application.
2. Open Schedule, change each range, switch agenda and calendar views, and inspect predicted and recorded occurrences.
3. Leave one occurrence selected while a relevant live event arrives and confirm focus and selection remain stable.
4. Open Activity, combine text, type, severity, and outcome filters, then inspect run, log, and alert details.
5. Acknowledge one alert and confirm unrelated alerts remain unchanged.
6. Clear a filtered view and confirm visible alerts are acknowledged, records are not deleted, and later activity appears.
7. Stop the daemon and confirm the last complete snapshots remain visible as read-only context.
