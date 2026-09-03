# Research: Desktop UX and Options

## Decision 1: Replace leading AppTabs with a package-owned shell

**Decision**: Use ordinary Fyne buttons for selectable destinations inside a
fixed-minimum-width leading rail, a stack for destination content, and a
separate bottom-aligned Exit button.

**Rationale**: `container.AppTabs` supports leading placement but does not offer
a bottom command slot or an API for balanced rail width. Exit is an application
command, not a content destination, so representing it as a tab would produce
incorrect selection and keyboard semantics.

**Alternatives considered**:

- Pad tab captions with spaces: rejected because font and DPI changes make the
  width unstable and Exit would still be a selected tab.
- Nest a second control under AppTabs: rejected because the two rails would not
  share width, selection styling, or focus order.
- Add a general navigation framework: rejected as disproportionate for seven
  destinations and one command.

## Decision 2: Persist validated identifiers through Fyne preferences

**Decision**: Store `dark`, `light`, or `system` and `brand`, `system`, or
`monospace` under namespaced preference keys. Normalize invalid values to Dark
and Brand before applying them.

**Rationale**: Fyne preferences are already scoped by the established
`tech.shruggie.goschedule` application identifier and are appropriate for two
small, non-sensitive, per-user choices. Persisting identifiers rather than
resources keeps upgrades and malformed state safe.

**Alternatives considered**:

- Add appearance to daemon configuration: rejected because appearance is
  user-specific desktop state, while the daemon configuration is machine-wide.
- Add a new preferences file: rejected because it duplicates framework
  lifecycle and platform location handling.
- Store serialized color or font resources: rejected because it creates
  compatibility and validation problems with no user benefit.

## Decision 3: Configure one immutable theme per selection

**Decision**: Extend `brandTheme` with normalized appearance and font fields.
Dark and light force their respective palettes; Follow system uses the variant
supplied by Fyne. Replacing the theme applies changes across current and future
controls.

**Rationale**: Fyne's Settings API owns application-wide theme replacement and
refresh. Immutable theme values avoid mutable shared styling state and permit
direct tests of every combination. Delegating symbols and unhandled roles to
the default theme preserves framework behavior.

**Alternatives considered**:

- Rebuild all application content after each choice: rejected because theme
  replacement already refreshes widgets and rebuilding risks losing view state.
- Maintain separate dark and light widget trees: rejected because it duplicates
  interface logic.
- Support arbitrary font files: rejected because untrusted font parsing,
  licensing, missing-file behavior, and chooser UX require a separate feature.

## Decision 4: Use stable task IDs for double activation

**Decision**: Bind the current task ID into each rendered row and implement
Fyne's single- and double-tap interfaces on the row. Forward single taps to
stable-ID list selection, and resolve the ID against the latest model snapshot
immediately before editing.

**Rationale**: `widget.List` selection exposes an index, but live refreshes can
reorder or remove tasks. An index captured during rendering is not an identity.
The stable ID is already the daemon/API identity and lets stale rows fail closed.
Explicit single-tap forwarding is required because Fyne directs the gesture to
the deepest tap-capable row child once it implements double activation.

**Alternatives considered**:

- Reuse the last selected index: rejected because it can edit a different task
  after a refresh.
- Open on every selection: rejected because it breaks single-click selection and
  keyboard browsing.
- Disable refresh while editing: rejected because it undermines the existing
  live-event contract.

## Decision 5: Present a bounded storage inventory

**Decision**: Resolve only known application-owned locations: machine data
root, database, logs, runtime lock/state, Fyne preference root/file, executable
directory, installed documentation when discoverable, and Windows cleanup
evidence when meaningful. Show the machine configuration path and the per-user
Fyne root separately from their files. Probe those exact paths with `os.Stat`
and no traversal. Running executable paths use ownership-neutral lifecycle copy
because a development binary is not necessarily installer-owned.

**Rationale**: Users need transparent support and uninstall information, not a
filesystem browser. Scope, existence, and preserve/wipe language prevent a path
from being mistaken for a deletion promise. Injected platform and filesystem
inputs make the contract deterministic in headless tests.

**Alternatives considered**:

- Recursively discover related-looking files: rejected because it is slow,
  privacy-invasive, and cannot establish ownership.
- Show hard-coded Windows paths on every platform: rejected because the values
  would be false outside Windows and can vary by environment.
- Claim custom daemon configuration paths are wipe-owned: rejected because the
  GUI cannot currently discover or authorize deletion of arbitrary overlays.

## Decision 6: Separate automated and native visual evidence

**Decision**: Headless tests verify semantic text configuration and layout
invariants. The exact staged candidate in #94 supplies Windows 11 standard and
scaled-DPI visual acceptance.

**Rationale**: Fyne's headless renderer can prove alignment, wrapping, resource
selection, and geometry, but it cannot prove native glyph rasterization on the
operator's Windows display. Treating screenshots or widget properties as native
sharpness proof would overstate evidence.

**Alternatives considered**:

- Close visual issues from headless tests alone: rejected because the reported
  defect is native display rendering.
- Add GUI automation to hosted Windows Server: rejected because #94 already
  defines the clean Windows 11 attended release-candidate gate and environment.
