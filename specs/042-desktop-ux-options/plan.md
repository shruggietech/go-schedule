# Implementation Plan: Desktop UX and Options

**Branch**: `codex/042-desktop-ux-options` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/042-desktop-ux-options/spec.md`

## Summary

Replace the fixed leading `AppTabs` shell with a package-owned navigation rail,
add a persistent Options view for bounded appearance choices and transparent
application-storage locations, route all exits through one idempotent shutdown
path, and bind task-row double-clicks to stable task identities. Extend the
existing theme rather than introducing a second styling system, remove wrapping
from the two affected Info labels, and keep native DPI acceptance in the exact
candidate gate defined by #94.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Fyne 2.8.1; existing embedded Geist, Geist Mono, and Space Grotesk resources; existing daemon client and view model

**Storage**: Fyne per-user preferences for two non-sensitive appearance identifiers; read-only presentation of existing configuration, runtime, executable, documentation, and maintenance paths; no schema migration

**Testing**: Go unit tests with Fyne's headless test driver, race tests for non-widget state, existing Windows evidence contracts, and canonical eight-gate verification

**Target Platform**: Windows desktop is the reported platform; production code and headless contracts remain cross-platform for Windows, macOS, and Linux

**Project Type**: Go desktop application with a local daemon backend

**Performance Goals**: Appearance changes and view selection visibly apply during the initiating interaction; no filesystem traversal; no new background polling; no scheduler-runtime impact

**Constraints**: Preserve 1280 by 800 preferred sizing and 800 by 600 clamp; UTF-8 without BOM; no arbitrary external fonts; no path mutation; no new process launcher; no release, merge, or tag; native sharpness claims require attended exact-candidate evidence

**Scale/Scope**: Five linked GUI issues, one navigation shell, one Options view, three appearance modes, three font choices, ten bounded storage categories, one task-row interaction, and one shutdown coordinator

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | UI state is decomposed into appearance, storage, navigation, row activation, and shutdown types with bounded responsibilities. The shell replaces an incapable widget rather than layering positional workarounds over it. |
| II. Testing Standards | PASS | Each behavior starts with a failing headless or pure-state test. Stable identity, malformed preferences, path classification, geometry, and one-shot shutdown have explicit regression cases. The canonical gate remains mandatory. |
| III. UX Consistency | PASS | Existing brand defaults remain intact; light/system/font options update the full app; navigation order, selected state, copy controls, and exit semantics are explicit. Native DPI acceptance stays attached to #94. |
| IV. Performance | PASS | Storage resolution probes only declared paths, appearance state is constant-sized, and navigation swaps already-built views. No filesystem scan or new polling loop is introduced. |
| V. Autonomous Execution | PASS WITH RECORDED OPERATOR OVERRIDE | Full Spec Kit and analyze gates remain mandatory. The operator explicitly authorized automatic S042 push and PR plus at most one additional Codex review round; merge, tag, and release remain outside authority. |

### Post-design re-check

All principles remain satisfied. No dependency, executable, database schema,
privilege boundary, or deletion boundary changes. A bounded read-only local API
response is added so Options cannot misrepresent a custom daemon configuration.
`CHANGELOG.md` and the exact-candidate validation instructions are pinned release
artifacts whose updates are required to record the behavior and evidence boundary.

## Architecture and Decision Log

### Own the leading navigation layout

Introduce a `navigationShell` that owns ordinary destination buttons, a content
stack, and a separately laid-out Exit command. The rail derives a stable minimum
width from the longest supported label plus symmetric theme padding, while a
bottom border region keeps Exit right-aligned and outside selected navigation.
Activity badge changes update the existing destination label through the shell.

`container.AppTabs` is removed because it exposes neither a bottom command slot
nor sufficient rail-width control. Padding around its tabs would still leave
Exit as a fake destination and would not provide the required geometry contract.

### Treat appearance as validated per-user state

Define bounded `appearanceMode` and `fontChoice` values. Load them through Fyne
preferences, normalize unsupported strings to Dark and Brand, and apply an
immutable `brandTheme` configured with both choices. `SetTheme` refreshes current
windows, and later controls inherit the same settings. Reset writes and applies
both defaults through one user action.

Dark and light palettes remain go-schedule branded. Follow-system honors the
variant passed by Fyne. System font delegates to the framework default; Brand
uses the current embedded faces; Monospace uses the bundled Geist Mono for all
non-symbol text. Arbitrary font files are excluded because they expand parsing,
licensing, accessibility, and persistence risk beyond the reported need.

### Resolve storage locations from explicit inputs only

Build daemon storage rows from absolute effective paths returned by a read-only
local runtime-information endpoint, then combine them with the Fyne app storage
root, running executable directory, and documented platform-specific maintenance
evidence. Each row carries scope, existence, and platform-accurate preserve/wipe
copy. The resolver receives injectable stat and platform inputs for deterministic
tests. It probes exact declared paths only and never walks a profile or filesystem.

Documentation is shown only when a known installed documentation path can be
derived and exists. Custom external configuration is presented from daemon
metadata but is never represented as owned or wiped.

### Bind task gestures to identity and guard editor ownership

Use a reusable task-row widget that implements Fyne's single- and double-tap
contracts and stores the rendered task ID. Because Fyne targets the deepest
tap-capable child, the single-tap callback explicitly forwards stable-ID
selection to `widget.List`. On double activation, resolve that ID against the
current task snapshot before calling the same detail lookup and editor path used
by the toolbar. A one-at-a-time editor guard is released by the dialog close
callback, including Save, Cancel, and dismissal.

### Centralize shutdown

Route title-bar close, connection-card Exit, and navigation Exit through a
single `requestClose` method guarded by `sync.Once`. The method cancels the run
context, clears the intercept, and closes the window once. This removes the
current duplicated close sequence and prevents rapid competing requests from
repeating lifecycle work.

### Keep visual evidence honest

Headless tests pin theme selection, font resources, label wrapping/alignment,
navigation geometry, and control semantics. They do not claim native Windows
text sharpness or DPI rendering. The #94 exact-candidate runbook remains the
required standard-DPI and scaled-DPI visual gate for #101, #104, and #105.

## External Research

- [Fyne widget package 2.8.1](https://pkg.go.dev/fyne.io/fyne/v2/widget) documents selectable labels and the list/widget contracts used for headless interaction.
- [Fyne theme package](https://docs.fyne.io/api/v2/theme/pkg/) defines semantic colors, sizes, variants, and default-theme delegation.
- [Fyne settings API](https://docs.fyne.io/api/v2/fyne/settings/) defines application-wide theme replacement and change notification.
- [Fyne preferences](https://docs.fyne.io/explore/preferences/) defines per-application persisted preferences through the established application identifier.
- [Fyne custom themes](https://docs.fyne.io/extend/custom-theme/) recommends semantic theme overrides with default delegation for unhandled roles.

## Project Structure

```text
gui/
├── app.go
├── app_test.go
├── appearance.go
├── appearance_test.go
├── navigation.go
├── navigation_test.go
├── options.go
├── options_test.go
├── storage_locations.go
├── storage_locations_test.go
├── tasks.go
├── widgets.go
├── widgets_test.go
├── theme.go
├── theme_test.go
├── info.go
└── info_test.go
internal/winuninstall/
├── ledger.go
├── ledger_test.go
└── platform_windows.go
test/windows/
└── README.md
docs/
└── options.md
specs/042-desktop-ux-options/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/desktop-ux.md
├── checklists/
├── tasks.md
└── verification.md
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Keep state and widgets in the existing `gui` package so
they can reuse the established Fyne headless harness and package-private
builders. Add focused files rather than enlarging `app.go`, and update only the
existing exact-candidate runbook and release-history surfaces outside that
package.

## Complexity Tracking

No constitution violation or exceptional complexity is introduced.
