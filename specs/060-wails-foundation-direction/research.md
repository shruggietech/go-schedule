# Research: Wails Foundation and Experience Direction

## Decision 1: Pin Wails v2.14.0 stable

**Decision**: Use Wails v2.14.0, the latest non-prerelease GitHub release observed on 2026-09-07. Do not adopt Wails v3 until a stable release exists and a later slice repeats the platform proof.

**Rationale**: The production replacement needs a supported release rather than a beta. Wails v2 is identified upstream as stable, supports Windows, macOS, and Linux, uses platform rendering engines, exposes Go/TypeScript bindings, native dialogs, clipboard, and unified events, and does not bundle a browser runtime. Sources: [Wails releases](https://github.com/wailsapp/wails/releases), [Wails repository](https://github.com/wailsapp/wails), [installation](https://wails.io/docs/gettingstarted/installation/), [runtime](https://wails.io/docs/reference/runtime/intro/), and [events](https://wails.io/docs/reference/runtime/events/).

**Alternatives considered**: Wails v3.0.0-beta.17 was rejected for production because upstream labels it prerelease software. Deferring all work for v3 GA was rejected because it would leave issues #149 and #150 blocked on an announced but not yet stable release. Electron was rejected because the issue and specification prohibit a bundled browser runtime.

## Decision 2: Use React, TypeScript, Vite, and authored CSS without a UI component framework

**Decision**: Start from the official Wails `react-ts` template baseline: React 19.1.0, React DOM 19.1.0, TypeScript 5.6.3, and `@vitejs/plugin-react` 5.0.0. Pin Vite 7.3.6 instead of the template's vulnerable 7.0.0 baseline. Use semantic HTML and project-authored CSS custom properties derived from the canonical brand tokens. Add no component or CSS framework.

**Rationale**: React has the strongest maintainership and accessibility-testing ecosystem among the built-in Wails templates, TypeScript makes the Go/frontend boundary explicit, and Vite is already the Wails template build path. Avoiding a component framework keeps design ownership local and prevents another visual system from competing with the approved go-schedule identity. The baseline stays on the official stable Wails template architecture, with bounded patch-level compatibility and security corrections. Sources: [Wails project creation](https://wails.io/docs/gettingstarted/firstproject/) and the [v2.14.0 React TypeScript template](https://github.com/wailsapp/wails/tree/v2.14.0/v2/pkg/templates/templates/react-ts).

**Installation finding**: The initial exact Vite 7.0.0 template pin was covered by high-severity development-server advisories. S060 therefore advances within the same major to Vite 7.3.6. jsdom 30.0.1 also requires Node 24.15.0 or newer, while the supported workspace Node 24 line includes 24.11.0, so the proof pins compatible jsdom 29.0.0. Both deviations preserve the selected architecture and produce a clean package audit.

**Alternatives considered**: Vanilla TypeScript minimizes dependencies but would require hand-building state, composition, and test conventions across the large migration. Preact is smaller but adds compatibility judgment without meaningful value at this scale. Svelte and Lit are viable, but their smaller project familiarity and testing ecosystems do not offset React's maintainability. Tailwind, shadcn, Material UI, and similar systems were rejected because they add build and design abstraction before the product primitives are settled.

## Decision 3: Pin Node.js 24 LTS and commit npm's lockfile

**Decision**: Set the frontend engine floor and hosted jobs to Node.js 24, use npm, pin direct dependency versions, and commit `package-lock.json`.

**Rationale**: Node 24 is an active LTS line, meets the current Vitest floor, is available on all GitHub-hosted operating systems through `actions/setup-node`, and matches the verified local toolchain. A lockfile makes the proof reproducible and reviewable. Automatic action caching is disabled because this job needs no secret-bearing or cross-trust cache. Sources: [Node.js release schedule](https://nodejs.org/en/about/previous-releases), [setup-node](https://github.com/actions/setup-node), and [Vitest requirements](https://vitest.dev/guide/).

**Alternatives considered**: Node 22 LTS remains supported but gives less runway. Node 26 is Current rather than LTS. Relying on runner-preinstalled Node was rejected because it makes the proof nondeterministic.

## Decision 4: Use layered Go, component, accessibility, and browser evidence

**Decision**: Test Go boundaries with the standard library under the race detector, React behavior with Vitest and Testing Library, DOM accessibility with axe-core, and real layout/keyboard behavior with Playwright Chromium at 1440 by 900, 900 by 650, and 200 percent zoom. Build the Wails proof on all three hosted desktop operating systems.

**Rationale**: No single layer can prove both native compilation and accessible interaction. Go tests make daemon and native substitutions deterministic; component tests make states fast to exercise; browser tests provide actual layout, focus, overflow, and media-query behavior; Wails builds expose platform dependency failures. Vitest supports the selected Vite line and Node 24. Sources: [Vitest guide](https://vitest.dev/guide/), [Playwright](https://playwright.dev/docs/intro), and [axe-core](https://github.com/dequelabs/axe-core).

**Alternatives considered**: Snapshot-only tests were rejected because they cannot prove keyboard or accessibility behavior. Full native UI automation was deferred to #157 because GitHub runners do not provide trustworthy attended window observation. Browser tests on all three platforms were rejected as redundant; the web layer is platform-neutral, while the native Wails build is the platform-sensitive evidence.

## Decision 5: Keep the proof real at boundaries and disposable in architecture

**Decision**: The committed proof uses the existing protected IPC client for health, task listing, and event streaming in its real entry point. Tests inject fakes. Native action calls pass through a small interface whose Wails implementation opens an informational native dialog. The proof lives in a nested module and is excluded from release packaging.

**Rationale**: This demonstrates the real local integration path without changing production or needing a live daemon in deterministic tests. The interface seam is the minimum needed for headless evidence and later transport-neutral work in #152. Isolation prevents experimental dependencies and generated bindings from becoming accidental production commitments.

**Alternatives considered**: Fixture-only Go methods would not exercise the existing daemon client. Importing Wails into the root module would expand every production dependency and verification path prematurely. Running a daemon inside frontend tests would make tests slower and platform-sensitive.

### Bounded capability evaluation

| Capability | Evidence and disposition |
| --- | --- |
| Local IPC bridge | The real entry point constructs the existing protected local client; deterministic tests cover health and list mapping without adding a listener. |
| Live daemon events | One startup-owned stream converts daemon events into bounded Wails events and is canceled and joined at shutdown under the race detector. |
| Window lifecycle | Wails owns startup and shutdown callbacks; duplicate startup is ignored and clean cancellation is tested. |
| Native dialogs | A replaceable adapter opens a real Wails informational dialog, while tests cover success and safe failure text. |
| Clipboard | Wails v2 exposes native clipboard operations through its runtime. The proof does not add a redundant clipboard wrapper because the required representative native boundary is the tested dialog adapter; clipboard behavior remains a production-shell decision in #151. |
| File dialogs | Wails v2 exposes native file-selection dialogs through the same context-bound runtime. No current scheduler workflow requires a file picker, so adding one to the proof would create a purposeless product contract; #151 must adopt it only for a concrete workflow. |
| Accessibility hooks | React renders semantic HTML into the platform webview. Vitest, axe-core, and Chromium checks cover names, landmarks, focus return, contrast, keyboard use, zoom, reflow, and reduced motion. |
| Installer integration | The local Windows proof builds to a 15,859,200-byte non-compressed executable with the approved icon. S060 deliberately excludes that output from current WiX and release inputs; final installer replacement and attended runtime checks remain in #157. |
| Dependency footprint | The locked frontend audit reports 204 total development and runtime packages with zero known vulnerabilities; the production proof has only React and React DOM as direct runtime packages. The nested Go module directly adds Wails and reuses the root module. |
| Release cadence | Stable Wails releases are evaluated deliberately rather than floated. Exact pins remain until a stable update repeats unit, race, browser, and three-platform native builds. Prereleases require a separate architecture decision. |

## Decision 6: Adopt the Calm Operations control-center direction

**Decision**: Use a persistent slim navigation rail, a target-and-health bar, one primary page region, and a contextual inspector. Favor quiet surfaces, sentence case, information-rich tables, strong empty states, and one obvious primary action. The direction is dark-first but fully supports light and system appearances.

**Rationale**: The structure keeps `This computer` visible without turning connection selection into dashboard clutter, scales to later multi-daemon work, and gives dense scheduler information a clear hierarchy. It reuses approved colors and type while materially departing from the current tabbed Fyne widget structure.

**Alternatives considered**: A dashboard-first grid was rejected because decorative metrics obscure operational decisions. A direct web recreation of Fyne was rejected by issue #150. A command-palette-first interface was rejected as the sole navigation model because discoverability and assistive-technology use matter for first-time users.

## Decision 7: Treat WCAG 2.2 AA as the minimum measurable contract

**Decision**: Require semantic landmarks, labeled controls, logical focus order, no keyboard traps, 4.5:1 normal-text contrast, 3:1 large-text and meaningful-graphic contrast, 3:1 focus indication, 24 by 24 CSS-pixel minimum targets or spacing exceptions, 200 percent zoom usability, and reduced-motion behavior.

**Rationale**: These thresholds translate the roadmap's accessibility intent into repeatable evidence and align with the current WCAG recommendation. Brand tokens already contain measured accessible light and dark roles. Source: [WCAG 2.2](https://www.w3.org/TR/WCAG22/).

**Alternatives considered**: Subjective accessibility review alone was rejected. WCAG AAA everywhere was rejected as disproportionate for the initial direction, though individual tokens can exceed it. Native screen-reader certification remains an attended #157 gate.

## Decision 8: Graduate contracts, not prototype shortcuts

**Decision**: #151 may reuse the selected dependency pins, tokens, semantic structure, tests, and component contracts. #152 may reuse the boundary vocabulary and cancellation rules. Fixture data, prototype routing, experimental module layout, and fake-native fallback are discarded unless their later issue independently justifies them.

**Rationale**: A proof is useful when it answers questions, but harmful when temporary shortcuts silently become architecture. An explicit graduation ledger makes the boundary reviewable.

**Alternatives considered**: Treating the proof as the first production shell would overrun both issues and skip #151's lifecycle and shared-component acceptance criteria. Throwing away all code would discard verified contracts and create needless rework.
