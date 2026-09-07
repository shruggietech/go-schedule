# Quickstart: Wails Foundation and Experience Direction

## Prerequisites

- Go 1.25.0
- Node.js 24 LTS and npm
- Windows: WebView2 runtime for attended launch
- macOS: Xcode command-line tools
- Linux: GCC, GTK3 development files, and WebKitGTK 4.1 development files

The canonical proof commands run from `experiments/wails-foundation` and do not modify scheduler data.

## Deterministic boundary tests

```bash
go test -race ./...
```

Expected: health and list mapping, disconnected behavior, event conversion, native action mapping, and cancellation ownership all pass without a running daemon or interactive desktop.

## Frontend component and accessibility tests

```bash
cd frontend
npm ci
npm test
npm run build
```

Expected: strict type-check, component states, keyboard interaction, semantic structure, axe checks, and offline-reference checks pass. The production bundle is created using repository-hosted assets only.

## Browser experience contract

```bash
npx playwright install chromium
npm run test:e2e
```

Expected: Tasks, task editing, Schedule, Activity, target switching, and compact views pass at the documented viewport and zoom settings with no horizontal page overflow, keyboard trap, serious accessibility finding, or hidden target identity.

## Native build proof

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean
```

On current Ubuntu runners, append `-tags webkit2_41`. The hosted CI matrix performs this build on Windows, macOS, and Linux. A successful build is compile and linkage evidence, not an attended installed-app observation.

## Optional local review

```bash
cd frontend
npm run dev
```

Open the printed loopback URL. Use the page navigation, state selector, appearance toggle, target switcher, keyboard focus order, task editor, and compact browser sizing to review the [experience contract](contracts/experience-contract.md).

## Repository verification

From the repository root:

```bash
sh scripts/verify.sh all
```

Expected: all eight canonical gates pass and the production Fyne desktop, installer, and release configuration remain unchanged except for the additive proof CI contract.

## Acceptance boundary

S060 closes #149 and #150 when the proof, hosted platform builds, documented direction, and maintainer review are complete. It does not close #151 or #152, publish a Wails application, replace Fyne, or claim attended native qualification.
