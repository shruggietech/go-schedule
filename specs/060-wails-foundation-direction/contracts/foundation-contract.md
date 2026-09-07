# Foundation Proof Contract

## Selected baseline

| Concern | Contract |
| --- | --- |
| Desktop runtime | Wails v2.14.0 stable |
| Frontend | React 19.1.0 and React DOM 19.1.0 |
| Language | TypeScript 5.6.3 in strict mode |
| Build | Vite 7.3.6 on Node.js 24 LTS with npm lockfile |
| Component tests | Vitest, Testing Library, jsdom, and axe-core |
| Browser tests | Playwright Chromium on Ubuntu |
| Native platforms | Wails build on Windows, macOS, and Linux |

All direct package versions are exact in manifests. The npm lockfile and Go sum file are committed. No application asset may resolve from an HTTP or HTTPS URL at runtime.

## Go to frontend boundary

The proof exposes these typed operations:

```text
Snapshot() -> ProofSnapshot
ShowAbout() -> NativeActionResult
proof:event -> ProofEvent
```

`Snapshot` uses the existing protected local IPC client in the real application and an injected deterministic daemon in tests. It requests daemon health and a task list. Transport endpoints, credentials, environment values, and raw errors never cross the frontend boundary.

`ShowAbout` calls a replaceable native-action adapter. The real adapter opens a Wails native informational dialog after the DOM is ready. Tests use a fake and verify invocation and result mapping.

The application begins one cancelable event stream after startup. Every received daemon event is converted to a non-sensitive `ProofEvent` and emitted through Wails. Shutdown cancels the stream and waits for it to terminate; no goroutine remains owned by the proof.

## Platform evidence

| Platform | Automated requirement | Explicit evidence boundary |
| --- | --- | --- |
| Windows | Go tests, frontend tests, exact Wails build | Build proves native binding and WebView2 loader compilation, not an attended installed-window observation |
| macOS | Go tests, frontend tests, exact Wails build | Build proves Cocoa/WebKit integration against the runner SDK, not notarization or an attended installed-window observation |
| Linux | Go tests, frontend tests, Wails build with WebKitGTK 4.1 | Build proves GTK/WebKit linkage on the runner distribution, not every supported distribution |

## Offline and release isolation

- Fonts, icons, styles, scripts, and fixtures are repository files included by Vite.
- No CDN, telemetry, analytics, remote font, runtime package download, network listener, or browser bundle is permitted.
- The proof directory is not referenced by the production `cmd/gosched-gui`, current installer inputs, or release workflow.
- CI may build proof artifacts but must not publish or upload them as release assets.

## Graduation ledger

| Element | Destination | Disposition |
| --- | --- | --- |
| Stable dependency pins and upgrade policy | #151 | Candidate to retain |
| Tokens, semantic layout, state components, and accessibility tests | #151 | Candidate to retain |
| Daemon identity, snapshot, event, cancellation, and error vocabulary | #152 | Candidate to refine and retain |
| Fixture data and prototype state toolbar | None | Discard |
| Nested module and proof routing | None | Discard |
| Fake native action fallback | Test support only | Reassess in #151 |
