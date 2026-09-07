# Wails foundation proof

This directory is the disposable, non-shipping evidence for Spec-Kit slice S060 and GitHub issues #149 and #150. It proves the selected desktop foundation and the Calm Operations experience direction without changing the current Fyne application, installer, release artifacts, or scheduler data.

## Selected baseline

| Concern | Version | License |
| --- | --- | --- |
| Wails | 2.14.0 stable | MIT |
| React and React DOM | 19.1.0 | MIT |
| TypeScript | 5.6.3 | Apache-2.0 |
| Node.js type declarations | 24.13.3 | MIT |
| Vite and React plugin | 7.3.6 and 5.0.0 | MIT |
| Vitest | 5.0.0 | MIT |
| Testing Library | 16.3.3 and user-event 14.6.7 | MIT |
| jsdom | 29.0.0 | MIT |
| axe-core and Playwright adapter | 4.13.0 | MPL-2.0 |
| Playwright | 1.63.0 | Apache-2.0 |
| Node.js | 24 LTS | MIT |

Transitive dependency licenses remain recorded in the committed npm and Go lock data. This proof is not distributed, but all direct licenses are compatible with repository use.

## Evidence boundaries

- The real entry point uses the existing protected local IPC client for health, task listing, and event streaming.
- Go tests replace daemon, event, and native-dialog boundaries deterministically.
- The browser preview uses local fixtures because it runs outside Wails; it does not open a network listener beyond Vite's loopback-only development server.
- Wails builds on Windows, macOS, and Linux prove native compilation and linkage. They do not claim installed or attended screen-reader qualification.
- All application fonts, marks, scripts, and styles are local repository assets. No CDN, remote font, telemetry, or analytics is used.
- Generated platform build directories, `build/bin`, `frontend/dist`, `node_modules`, Playwright reports, and test results are ignored and never become release inputs.

## Prerequisites

- Go 1.25.0
- Node.js 24 LTS and npm
- Windows WebView2 for attended launch
- macOS Xcode command-line tools
- Linux GCC, GTK3 development files, and WebKitGTK 4.1 development files

See [the S060 quickstart](../../specs/060-wails-foundation-direction/quickstart.md) for commands and expected evidence.

## Graduation ledger

The exact dependency baselines, brand-derived tokens, semantic structure, state vocabulary, and accessibility tests are candidates for #151. Target identity, snapshot, event, cancellation, and safe-error vocabulary are candidates for #152. Fixture data, prototype route switching, the state toolbar, the nested module, and the browser-preview native fallback are intentionally disposable.

No proof shortcut becomes production architecture merely because it exists here. A later issue must adopt it explicitly and satisfy that issue's acceptance criteria.
