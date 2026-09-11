# Verification: Desktop Visual and Shell Foundations

## Specification analysis

The post-task spec-kit analysis found 27 buildable functional requirements and success criteria mapped across 20 tasks. Requirements and UX checklists are complete, no clarification marker remains, no task is unmapped, and no constitution conflict or critical finding remains.

## Test-first evidence

The first focused component run failed seven new expectations before implementation: pending button semantics, affirmative action styling, grouped dialog actions, dismissible errors, timed toast replacement, hover-paused dismissal, and resolved follow-system appearance. After implementation, all focused and full frontend suites passed.

## Focused frontend evidence

- `npm test -- --run` passed all 115 component, store, integration, and accessibility tests in 24 files after both review rounds added feedback lifecycle coverage.
- `npm run build` passed TypeScript validation and the Vite production build.
- `npx playwright test e2e/shell.spec.ts` passed nine Windows-hosted Chromium checks.
- The browser contract covered 1440 by 900, 900 by 650, and 800 by 600 viewports; 200 percent zoom; one page scroll owner; retained navigation, Exit, target context, and Appearance; no document overflow; all five semantic action variants; compact 24 to 34 CSS-pixel dimensions; label contrast of at least 4.5:1; visible hover and focus changes; light and dark palettes; follow-system resolution; reduced motion; local assets; and zero serious or critical axe findings.

## Native Windows evidence

The canonical GUI gate ran on the Windows host and built `desktop/build/bin/gosched-gui.exe` in production mode with Wails 2.15.0. The native host now declares an 800 by 600 minimum window and `windows.SystemDefault` title-bar theming. The exact frontend embedded by that build is the production bundle validated by the focused Windows browser interaction contract above, so native host configuration, embedded asset generation, palette resolution, pointer hover, keyboard focus, responsive scrolling, and transient feedback all derive from one source tree. No application window was launched because repository policy prohibits automation from creating a foreground or focus-stealing console or GUI surface.

## Canonical verification

`scripts/verify.sh all` passed on 2026-09-11 through the verified headless Git Bash launcher after the environment restart removed the prior `sh` command alias.

- format: PASS, including `github-format` with no em dash or hard-wrapped Markdown prose
- vet: PASS
- lint: PASS with zero issues
- race: PASS, including the complete integration package
- gui: PASS, including desktop race tests, the native Windows Wails production build, all 106 frontend tests, and the production bundle
- coverage: PASS at engine 82.9 percent, schedule 89.1 percent, timezone 91.3 percent, store 80.1 percent, catchup 88.9 percent, and logbus 91.1 percent
- docs: PASS, including 20 pages and architecture mutation fixtures
- automation: PASS, including actions, CodeQL, Dependabot, release operations, release notes, brand, lifecycle, eight-gate policy, and automation mutation fixtures

## First-round review evidence

The initial Codex review identified three P2 feedback lifecycle gaps. Failing regression tests reproduced each gap before the corrections: unsuccessful settings actions now remain in dismissible persistent notices, every settings action emits a distinct identity so identical outcomes restart their announcement, and dynamic JSX notices accept caller event identities so later actionable failures restore after an earlier dismissal. The focused 22-test review suite, complete 111-test frontend suite, production bundle, nine Playwright shell checks, and all eight canonical gates passed after the corrections.

The first CI run on the review fix exposed a macOS-only test scheduling race: the repeated-announcement assertion selected the intentionally empty live region before React populated the visible toast. The product behavior passed on Linux and Windows. The test now waits for the visible announcement text itself, preserving the behavioral assertion without depending on effect scheduling order.

## Second-round review evidence

The authorized second and final Codex review identified four additional P2 edge cases. Failing tests reproduced active-toast timeout reuse, discarded compound success content, a dismissed warning hiding a later connection snapshot, and a routine connection announcement replacing a persistent settings error. The corrections restart dismissal timers from event identity, preserve renderable success content, key dynamic connection warnings from snapshot state, and maintain independent stacked channels for transient announcements and persistent settings failures. The focused 22-test review suite, complete 115-test frontend suite, production bundle, and nine Playwright shell checks passed after the corrections. No third review round will be requested.

## Release boundary

S083 modifies source, tests, specifications, and changelog only. It does not move, rebuild, replace, publish, or promote the existing v1.4.0 tag or draft release.
