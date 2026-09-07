# Verification: Wails Foundation and Experience Direction

## Local implementation evidence

| Gate | Result |
| --- | --- |
| Nested Go boundary and lifecycle | `go test -race ./...` passed on Windows with connected, disconnected, degraded, native-action, event, duplicate-start, cancellation, and shutdown coverage. |
| Frontend component and accessibility | `npm run test` passed 10 tests across four files. |
| Type safety and production bundle | `npm run build` passed with TypeScript 5.6.3 and Vite 7.3.6. |
| Browser experience contract | `npm run test:e2e` passed five Chromium tests at 1440 by 900, 900 by 650, light and dark appearances, required states, keyboard flow, axe serious and critical violations, screenshots, horizontal overflow, and 200 percent zoom. |
| Dependency audit | `npm audit --audit-level=high` reported zero vulnerabilities from the committed lockfile. |
| Native Windows proof | `go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean` generated the typed bindings, frontend bundle, approved application icon, and a 15,859,200-byte Windows executable. |
| Brand integrity | The canonical automation gate ran `go run ./scripts/brand-check` and accepted every registered proof font and icon copy. |
| Workflow policy | `test/scripts/automation-check_test.sh automation` accepted the exact Node and Wails baseline and rejected stale Node action and incomplete platform-matrix fixtures. |
| Publication format | `go run ./scripts/github-format` reported no Unicode em dash or hard-wrapped Markdown defects. |
| Encoding and diff integrity | UTF-8, mojibake search, `git diff --check`, ignored-output boundaries, and repository status are audited before each publication commit. |

## Canonical repository verification

`scripts/verify.sh all` passed in the foreground on 2026-09-07 after supplying the native Windows Go and gofmt executable paths to the WSL shell environment.

| Gate | Result |
| --- | --- |
| format | Passed, including GitHub publication format. |
| vet | Passed. |
| lint | Passed with zero issues. |
| race | Passed across the complete selected package set, including the 41.307-second integration suite. |
| gui | Passed for `gui` and `gui/viewmodel`. |
| coverage | Passed: engine 81.9%, schedule 89.2%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1%. |
| docs | Passed all 15 pages, links, front matter, fences, theme, product policy, and fixtures. |
| automation | Passed canonical automation plus its fixture regression suite. |

## Hosted evidence pending publication

The official pull request must run the exact proof on Windows, macOS, and Linux and run the Chromium accessibility contract on Ubuntu. Those hosted results, CI URLs, and Codex review dispositions will be appended before the slice advances from `In Progress` to `Implemented`.

## Issue traceability

| Issue | Acceptance evidence |
| --- | --- |
| #149 | `research.md` selects exact stable versions, compares bounded alternatives, evaluates IPC, events, lifecycle, native APIs, accessibility, platform requirements, installer isolation, dependency footprint, release cadence, and upgrade policy; the nested proof supplies deterministic boundary tests and native builds. |
| #150 | `experience-contract.md` and the React prototype establish the Calm Operations direction across Tasks, task editing, Schedule, Activity, target switching, compact reflow, nine state classes, three appearance modes, reusable tokens, focus, dialog, keyboard, contrast, zoom, target, live-region, and reduced-motion rules. |
| #147 | Remains open. S060 resolves only children #149 and #150; #151 through #157 remain actionable. |

## Evidence boundary

Local Windows compilation does not claim hosted macOS or Linux compilation, attended native observation, installer replacement, screen-reader certification, production shell completion, or transport-layer completion. The proof is excluded from current Fyne binaries and release inputs.
