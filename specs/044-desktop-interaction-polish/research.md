# Research: Desktop Interaction and Appearance Polish

## Decision 1: Use translucent semantic interaction overlays

**Decision**: Keep the stock button renderer and change hover, pressed, and
focus roles from opaque replacement colors to intentionally translucent
overlays whose composites pass the required contrast matrix.

**Rationale**: Fyne 2.8.1's button renderer blends state colors over the base
importance fill while retaining the importance-specific foreground. Opaque
overlays erase the base, exactly reproducing the S043 active-hover defect.
Fyne's theme guidance also favors semantic global roles for consistent widgets.

**Alternatives considered**:

- Build a complete custom button renderer: rejected because it would duplicate
  focus, disabled, tap animation, alignment, icon, and accessibility behavior.
- Patch navigation and Save individually: rejected because the observed defect
  is shared and would remain in other importance states.
- Change foreground-on-primary globally: rejected because one foreground cannot
  compensate for unrelated opaque hover, focus, and pressed replacement fills.

**Sources**:

- [Fyne theme and customisation](https://docs.fyne.io/faq/theme/)
- [Fyne button API](https://docs.fyne.io/api/v2/widget/button/)
- Pinned Fyne 2.8.1 `widget/button.go` renderer source in the module cache

## Decision 2: Default to System and package Inter plus Ubuntu

**Decision**: Make System the clean/reset/invalid default; retain explicit Brand
and Monospace preferences; add Inter and Ubuntu as unmodified packaged
regular/bold families.

**Rationale**: The operator's S043 native A/B test found System immediately
sharp on the affected QHD ultrawide display. Keeping Brand preserves explicit
user intent and compatibility, while Inter is designed for screen readability
and Ubuntu satisfies the requested familiar family choice. Packaging exact
resources avoids depending on host font installation and keeps selections
portable.

**Alternatives considered**:

- Remove Brand: rejected because a stored explicit choice would silently change
  and retaining it has negligible UI cost.
- Enumerate installed fonts: rejected because platform discovery, missing-font
  recovery, licensing ambiguity, and an unbounded chooser exceed the issue.
- Add many families: rejected because each asset increases binary size, license
  inventory, visual verification, and maintenance.

**Sources**:

- [Inter upstream project](https://github.com/rsms/inter)
- [Inter 4.1 distribution](https://rsms.me/inter/download/)
- [Ubuntu Sans upstream project](https://github.com/canonical/Ubuntu-Sans-fonts)
- [Ubuntu Font Licence 1.0](https://ubuntu.com/legal/font-licence)

**Pinned asset provenance**:

- Inter regular/bold and `OFL-Inter.txt` come unmodified from official Inter
  release `v4.1` (`e3a3d4c57d5ecc01453a575621882a384c1995a3`). SHA-256:
  `40d692fce188e4471e2b3cba937be967878f631ad3ebbbdcd587687c7ebe0c82`,
  `288316099b1e0a47a4716d159098005eef7c0066921f34e3200393dbdb01947f`,
  and `262481e844521b326f5ecd053e59b98c8b2da78c8ee1bdbb6e8174305e54935a`.
- Ubuntu Sans regular/bold and `LICENCE-Ubuntu.txt` come unmodified from
  Canonical commit `9554af00fb9d438a12c916df8451c10dcedc9b7e`. SHA-256:
  `74f238be44ac5e2ad41021f0b4acc5ccc66f585d06c36b22931319d9751d50ea`,
  `185c0fcde30b8b75c543793aeaa27927cc1a9970dac89b13e18acd1b26a3bbb7`,
  and `2f0015108d68627bd788d313f529c21ff4da2c2c42a5e1f3883acc83480f9002`.

## Decision 3: Keep storage details visible in compact two-line rows

**Decision**: Replace cards with a header plus aligned two-line rows. The exact
path and Copy occupy the primary line; scope, existence, and removal behavior
wrap on the secondary line.

**Rationale**: This removes repeated card chrome and reduces vertical/lateral
surface area while keeping critical uninstall claims continuously available to
keyboard and pointer users. Wrapping vertically satisfies the explicit ban on
horizontal scrollbars, so a hover-only tooltip is unnecessary.

**Alternatives considered**:

- A tooltip for all metadata: rejected because hover-only discovery is poor for
  keyboard and touch users and hides important wipe behavior.
- A framework table: rejected because variable-height wrapped paths and details
  conflict with fixed row sizing and would overtake S045's dedicated table work.
- Truncate exact paths: rejected because storage diagnosis requires the full
  value to remain readable/selectable.

## Decision 4: Exclude the current value from selector menus

**Decision**: Keep the selected text in the closed selector and replace its
options with only the other valid values after every selection.

**Rationale**: The stock Select explicitly stores displayed `Selected` apart
from `Options`, and `SetOptions` does not clear `Selected`. This directly meets
the user's complaint without a custom popup or confusing duplicate current
entry.

**Alternatives considered**:

- Keep every value and add “(current)”: rejected because the current entry still
  looks actionable.
- Disable the current popup item: rejected because the stock API exposes no
  per-item disabled state.
- Replace selectors with radio groups: rejected because it consumes additional
  space on an already long Options page.

**Sources**:

- [Fyne Select API](https://docs.fyne.io/api/v2/widget/select/)
- [Fyne choice widgets](https://docs.fyne.io/widget/choices/)

## Decision 5: Scale only discrete scroll deltas

**Decision**: Wrap each application-owned vertical scroll and multiply event
deltas at or above 25 logical pixels by a persisted 1x-4x factor; leave smaller
precision deltas unchanged and delegate all movement to the stock scroll.

**Rationale**: Fyne 2.8.1 exposes only a two-axis delta, not an input-device
identifier. Its pinned Windows/Linux GLFW driver converts one conventional
detent to 25 logical pixels, while precision devices can deliver smaller
fractional deltas. The threshold is the narrowest available distinction and is
isolated behind directly tested normalization logic.

**Alternatives considered**:

- Multiply every delta: rejected because it would amplify precision touchpad
  movement and violate the issue.
- Adjust offset from `OnScrolled`: rejected because that callback also covers
  scrollbar movement and explicitly warns against updating Offset there.
- Fork or patch Fyne's driver: rejected because the requested behavior is an app
  preference and a dependency fork is disproportionate.

**Sources**:

- [Fyne Scrollable API](https://docs.fyne.io/api/v2/fyne/scrollable/)
- [Fyne vertical scroll guidance](https://docs.fyne.io/container/scroll/)
- Pinned Fyne 2.8.1 `internal/driver/glfw/window.go` and
  `internal/widget/scroller.go` source in the module cache

## Decision 6: Use the existing shell and semantic danger treatment for Exit

**Decision**: Add a full-height rail/content separator and use the framework's
danger importance for Exit, leaving its current bottom-right layout and
one-shot callback unchanged.

**Rationale**: S043 confirmed placement and shutdown behavior. A semantic danger
button colors label and themed glyph together and provides a recognizable
destructive command without treating Exit as selected navigation.

**Alternatives considered**:

- Add a custom red icon only: rejected because label and glyph state could drift.
- Rebuild navigation again: rejected because current sizing, order, placement,
  and behavior are accepted and tested.

## Decision 7: Separate local proof from native release evidence

**Decision**: Complete deterministic tests and canonical verification locally,
and retain native DPI, physical wheel/touchpad, and final glyph-rasterization
checks as exact-candidate release evidence.

**Rationale**: Headless rendering can prove semantics and geometry but cannot
honestly prove how Windows rasterizes text or how physical hardware feels. The
prior native A/B report authorizes the System default but does not prove the
post-change installer.

**Alternatives considered**:

- Represent headless screenshots as native proof: rejected as false evidence.
- Expand S044 into release-candidate qualification: rejected because #94 and
  the release-readiness epic already own that gate.
