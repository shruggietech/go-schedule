# Quickstart: Validate Desktop Visual and Shell Foundations

## Prerequisites

- Node.js 26 and the repository's pinned frontend dependencies.
- Go 1.25 toolchain and Wails prerequisites for the native build.
- Chromium installed through the existing Playwright setup.

## Focused automated checks

```sh
cd desktop/frontend
npm test -- --run src/components/components.test.tsx src/App.test.tsx src/accessibility.test.tsx src/tasks/TaskActions.test.tsx
npm run build
npx playwright test e2e/shell.spec.ts
```

Expected result: shared variants, dialog focus and spacing hooks, transient feedback timing, follow-system resolution, route scrolling, 800 by 600 reflow, 200 percent zoom, reduced motion, and axe checks pass.

## Native Windows observation

1. Build and launch the development Wails desktop through the repository GUI gate or documented native build command.
2. Record the exact commit and executable used.
3. In light, dark, and follow-system appearances, observe Connection details, Tasks actions, a confirmation dialog, and at least one routine success message.
4. Confirm visible hover and keyboard focus, compact semantic actions, theme-consistent identity, fixed navigation and Exit while a long route scrolls, and automatic removal of routine feedback.
5. Record the focused observation in `verification.md`; do not stage, alter, or promote a GitHub release.

## Canonical verification

```sh
sh scripts/verify.sh all
```

Expected result: format, vet, lint, race, GUI, coverage, documentation, and automation gates all pass without an S083 exclusion.
