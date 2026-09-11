# Quickstart: Validate Task-First Notifications

## Prerequisites

- Node.js 26 and repository-pinned frontend dependencies.
- Go 1.25 toolchain and Wails prerequisites for the native build.
- Chromium installed through the existing Playwright setup.

## Focused automated checks

```sh
cd desktop
go test -race ./notifications
cd frontend
npm test -- --run src/notifications/NotificationsPage.test.tsx src/notifications/store.test.ts
npm run build
npx playwright test e2e/notifications.spec.ts
```

Expected result: configured scope coverage is accurate and secret-free; initial overview, state guidance, recent results, collapsed advanced sections, preserved channel and policy operations, diagnostics, keyboard use, 800 by 600 reflow, 200 percent zoom, and accessibility checks pass.

## Native Windows observation

1. Build the production Wails desktop through the canonical GUI gate and record the exact commit and executable.
2. Open Notifications with no channels and verify the first viewport explains notifications are off and offers setup.
3. Open Notifications with enabled configured destinations and recent successful, retrying, and failed results and verify status, affected task context, and next action are understandable before specialist setup.
4. Expand Destinations, Assignment rules, and Delivery diagnostics one at a time by keyboard and exercise existing create or edit, policy selection, filtering, and selected delivery detail.
5. Repeat at 800 by 600, 200 percent zoom, and light, dark, and follow-system appearances.
6. Confirm no horizontal page scrolling, clipped controls, inaccessible disclosure, secret display, or release mutation.

## Canonical verification

```sh
sh scripts/verify.sh all
```

Expected result: format, vet, lint, race, GUI, coverage, documentation, and automation gates all pass without an S085 exclusion.
