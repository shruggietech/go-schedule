# Quickstart: Validate Responsive Task and Administration Workflows

## Prerequisites

- Node.js 26 and the repository's pinned frontend dependencies.
- Go 1.25 toolchain and Wails prerequisites for the native build.
- Chromium installed through the existing Playwright setup.

## Focused automated checks

```sh
cd desktop/frontend
npm test -- --run src/components/components.test.tsx src/tasks/TasksPage.test.tsx src/tasks/TaskEditor.test.tsx src/tasks/TaskActions.test.tsx src/agentaccess/AgentAccessPage.test.tsx src/settings/ConnectionsPage.test.tsx src/settings/SettingsPage.test.tsx src/settings/store.test.ts src/remotepairing/PairingForm.test.tsx
npm run build
npx playwright test e2e/shell.spec.ts
```

Expected result: modal task create and edit, focus containment and restoration, task validation, safe example insertion, initial disclosures, responsive administration grids, complete long values, record-scoped copy state, supported appearances, keyboard navigation, and axe checks pass.

## Native Windows observation

1. Build the production Wails desktop through the canonical GUI gate and record the exact commit and executable.
2. At the minimum supported window, open Create task, use Insert example, expand Advanced settings, cancel with Escape, reopen, and save or cancel through the action row.
3. Select a task, open Edit, close it, and exercise Delete confirmation with keyboard navigation.
4. Open Agent Access, inspect basic status, and expand optional localhost configuration.
5. Open Connections, inspect diagnosis and saved profiles, then expand ordinary pairing and one repair flow.
6. Open Settings in light, dark, and follow-system appearances, inspect long paths, and activate one Copy path action while observing unrelated controls.
7. Confirm no collision, clipping, horizontal page scrolling, hidden focus, interface-wide button flicker, secret display, or release mutation.

## Canonical verification

```sh
sh scripts/verify.sh all
```

Expected result: format, vet, lint, race, GUI, coverage, documentation, and automation gates all pass without an S084 exclusion.
