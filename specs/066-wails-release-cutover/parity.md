# Desktop Parity and Intentional Deviations

This inventory closes the #147 functional boundary. Every materially supported Fyne capability has one Wails disposition and evidence. No entry is blocked.

| Legacy capability | Disposition | Wails destination or retirement | Evidence | User impact |
| --- | --- | --- | --- | --- |
| Application shell and primary navigation | Wails | Persistent rail, responsive compact navigation, appearance selector, and separate Exit command | S061 shell tests and `desktop/frontend/e2e/shell.spec.ts` | Fresh layout with the same reachable work areas |
| Local daemon startup, health, retry, and shutdown | Wails | Transport-neutral local connection manager with bounded automatic and manual recovery | S061 connection tests and S065 Connections tests | Inline recovery replaces repeated modal errors |
| Tasks search, detail, create, edit, state, delete, and Run now | Wails | Tasks workspace and editor | S062 task tests and `desktop/frontend/e2e/tasks.spec.ts` | Advanced options are progressively disclosed |
| Platform-native example command | Wails | Shared read-only example insertion in task command entry | S062 tests for #192 | Same safe example, clearer one-shot insertion |
| Nested groups and inherited state | Wails | Groups panel integrated with Tasks | S062 group service and frontend tests | Tasks and Groups are one authoring flow |
| Chains | Wails | Automation workspace Chains section | S063 automation service and browser tests | Connected source presentation replaces isolated tabs |
| External triggers and trigger sets | Wails | Automation workspace Triggers section | S063 trigger tests | Trigger relationships are visible in one workspace |
| Filesystem watchers and health | Wails | Automation workspace Watchers section | S063 watcher tests | Health and target context are presented inline |
| Schedule list and calendar context | Wails | Schedule workspace with bounded horizon and selection detail | S064 operations tests and browser checks | Operational emphasis replaces Fyne table mechanics |
| Activity, logs, run detail, and alerts | Wails | Activity workspace with severity/source filters, detail, and refresh | S064 activity tests | Alerts are unified with authoritative run activity |
| Options and appearance | Wails | Settings workspace with System, Light, and Dark | S065 preference and settings tests | System is the safe default; choices persist atomically |
| Info, version, storage paths, copy, and product links | Wails | Settings information and storage sections with allowlisted native actions | S065 settings bridge and frontend tests | Information is consolidated and remains useful offline |
| Connection diagnostics | Wails | Connections workspace with stable focused retry control | S065 review regression tests | Recovery is visible without modal repetition |
| Application exit | Wails | Separate Exit rail command using orderly native shutdown | S061 application tests | No selectable pseudo-tab |
| Keyboard navigation, focus, scaling, contrast, and reduced motion | Wails | Semantic React controls and browser/native contracts | S061 through S065 unit and Playwright suites | Stronger documented accessibility contract |
| Fixed headers, whole-row selection, and full-value disclosure | Wails | Responsive structured lists and detail disclosures | S062 and S064 component/browser tests | No horizontal-scrolling table dependency |
| Interface font selection | Retired | Bundled responsive typography owns the supported type system | S065 specification and `docs/options.md` | Arbitrary framework font selection is no longer offered |
| Scroll-sensitivity multiplier | Retired | Native webview and operating-system scrolling remain authoritative | S065 specification and `docs/options.md` | Platform scrolling settings apply directly |
| Fyne canvas-size evidence hook | Retired | Browser zoom, viewport, native build, and platform package evidence replace toolkit canvas units | S060 experience contract plus production Playwright and native CI | Evidence describes current webview and native surfaces |

## Intentional deviations

### Connected workspaces replace a tab-for-widget translation

Tasks and Groups, automation sources, Schedule and Activity, and Settings and Connections are organized around user outcomes instead of preserving the former widget tree. This reduces navigation churn and keeps target and recovery context visible. The feature-level browser and service tests prove every underlying operation remains reachable.

### Inline recovery replaces modal repetition

Connection failures stay inside the affected workspace with diagnosis and a focus-stable retry action. This avoids repeated dialogs during daemon recovery while retaining actionable error text and manual retry.

### Framework-specific font and scrolling controls are retired

Wails uses the bundled brand typography and native webview scrolling. Recreating Fyne-only controls would imply behavior the replacement runtime does not own. Valid appearance intent still migrates once; unsupported or malformed legacy values choose System safely.

### Current evidence replaces toolkit-specific evidence

Fyne canvas units and headless widget rendering are no longer release signals. The production contract uses semantic frontend tests, Chromium accessibility and zoom checks, native Wails builds, stable executable inspection, and platform package lifecycle tests tied to one candidate revision.
