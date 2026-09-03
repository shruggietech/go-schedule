# Implementation Plan: Desktop Interaction and Appearance Polish

**Branch**: `codex/044-desktop-interaction-polish` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/044-desktop-interaction-polish/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Correct the systemic active/hover contrast defect by supplying the translucent
semantic overlays the Fyne button renderer expects, make System the safe font
default while adding packaged Inter and Ubuntu alternatives, compact Options
storage into aligned rows, omit the current selection from selector menus,
finish the navigation boundary and Exit glyph, and route every application-owned
vertical scroll through one bounded discrete-wheel multiplier. Preserve S042's
storage truth, navigation order, editor behavior, and one-shot shutdown.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Fyne 2.8.1, Fyne preferences/settings/theme APIs,
packaged Inter 4.1 and Ubuntu Sans font resources

**Storage**: Existing Fyne per-user preferences only; no daemon database or
configuration migration

**Testing**: Go unit tests, Fyne headless GUI tests, focused race tests, WCAG
contrast calculations, repository canonical `scripts/verify.sh all`

**Target Platform**: Cross-platform desktop with native Windows 11 acceptance
at 100% and scaled DPI

**Project Type**: Desktop application over the existing local daemon/API

**Performance Goals**: One conventional Windows wheel detent moves a vertical
view by 50 logical pixels at the default 2x setting; preference application and
theme replacement remain synchronous and visibly immediate

**Constraints**: No horizontal Options scrolling, no recursive storage access,
no daemon schema change, no locally installed font enumeration, no input-source
field in Fyne scroll events, no weakening of existing shutdown or task editor
semantics

**Scale/Scope**: Six linked v1 GUI issues, five application-owned vertical
scroll containers, five curated font choices, three color modes, and the
existing seven destinations plus one Exit command

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Code Quality**: PASS. Shared preference, palette, row-layout, selector, and
  scroll helpers keep behavior centralized and bounded. No new goroutine or
  mutable global lifecycle is introduced.
- **Testing Standards**: PASS. Each behavioral correction starts with a focused
  regression test. Pure normalization, contrast, geometry, and delta logic stay
  deterministic; the complete GUI and canonical gates remain required.
- **UX Consistency**: PASS. One semantic state palette covers shared controls,
  one font catalog serves the complete app, one storage-row contract replaces
  cards, and one scroll preference applies to all owned vertical scroll views.
- **Performance**: PASS. The work adds constant-time color calculations,
  bounded preference lookup, and one delta multiplication per wheel event. No
  scheduling hot path changes.
- **Autonomous Build-Phase Execution**: PASS. S044 follows specify, clarify,
  checklist, plan, tasks, analyze, implement, verify, and commit on a review
  branch. The operator's kickoff explicitly pre-authorizes the branch push and
  PR publication after the local pre-publication report.
- **Engineering Constraints**: PASS. Font assets are redistributed unmodified
  with their upstream licenses. Persisted identifiers are forward-compatible;
  unknown values fail safely. No secrets or external inputs are added.
- **Post-design re-check**: PASS. The design adds no constitution exception and
  does not modify a pinned process artifact.

### Spec-kit procedure note

The installed `/speckit-checklist` prerequisite invocation requires `plan.md`
even though both the project autopilot protocol and the command order place
checklist before plan. S044 retained the required order and used the feature
path already resolved by `/speckit-clarify` to produce `checklists/ux.md`. This
explicitly avoids copying the contradictory precheck into project practice.

## Technical Design

### Preserve semantic button blending

Fyne 2.8.1 composes a button's base importance color with semantic hover,
pressed, or focus colors. S042 supplied opaque hover and pressed colors, which
replaced the primary fill while leaving the foreground-on-primary color in
place. S044 uses translucent overlays for those roles and verifies the actual
composites for medium, low, high, danger, success, and warning importance in
both palettes. Selected navigation keeps its persistent filled shape, so hover
cannot erase selected identity.

### Migrate only the default font

An empty, missing, unknown, or reset font value resolves to System. Explicit
`brand`, `monospace`, `inter`, or `ubuntu` identifiers survive. Brand remains
available as `Geist (brand)` to avoid silently invalidating an existing choice;
System, Inter, Ubuntu, and Monospace provide the requested alternatives. Inter
and Ubuntu ship as unmodified regular/bold static faces with their licenses.
Symbol roles delegate to Fyne and monospace-semantic roles retain the packaged
Geist Mono resource when a proportional family is selected.

### Replace storage cards with compact rows

Use a header and one two-line aligned row per `storageLocation`: category at the
left, a flexible path/detail column in the middle, and Copy at the far right.
Long content wraps vertically. Unavailable rows use disabled semantic text and
background treatment across the entire row and expose no selectable path or
enabled Copy action. All existing path, scope, existence, and wipe text remains
visible, so a tooltip is unnecessary and keyboard users lose no information.

### Treat selector menus as actions

The closed selector displays the current value while its menu contains only
other valid values. Fyne permits this by preserving `Selected` while
`SetOptions` replaces the menu entries; after each change the menu is rebuilt
around the new current value. Arrow-key selection remains deterministic because
an absent selected index advances to the first or last alternative.

### Amplify discrete wheel input at the owned scroll boundary

Introduce a `sensitiveVScroll` that embeds the framework scroll, rebinds its
base widget to the outer type, and overrides only `Scrolled`. The wrapper copies
the event, multiplies Windows/Linux discrete deltas whose magnitude is at least
the pinned driver's 25-logical-pixel detent by the normalized preference, and
delegates to the stock scroll. Smaller deltas, used by precision devices,
remain unchanged. Keyboard, drag, scrollbar, and programmatic offset paths never
enter `Scrolled` and therefore remain stock behavior. Replacing each direct
application `NewVScroll` prevents nested multiplication.

Sensitivity ranges from 1x through 4x in 0.5 increments, defaults to 2x, and is
stored separately from appearance so Restore defaults can reset all three
choices deliberately. A label presents the active multiplier next to the
slider.

### Finish the navigation rail without restructuring it

Place a full-height semantic separator between the rail and content using the
existing border shell, retain the 168-pixel content-derived minimum, and color
Exit through danger importance so both its glyph and label use the semantic
error treatment. Exit remains outside destination selection and retains the
one-shot close coordinator.

### Keep visual evidence honest

Headless tests pin contrast math, theme resources, state composites, label
invariants, row geometry, selector alternatives, navigation layout, preference
normalization, and scroll deltas. The operator's S043 A/B observation is valid
evidence for choosing System as the default. Native after-change sharpness,
palette state inspection, physical wheel feel, and scaled-DPI behavior remain
release-candidate evidence and are not fabricated by headless tests.

## Project Structure

### Documentation (this feature)

```text
specs/044-desktop-interaction-polish/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── desktop-interaction.md
├── checklists/
│   ├── requirements.md
│   └── ux.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
gui/
├── appearance.go
├── appearance_test.go
├── scroll.go
├── scroll_test.go
├── theme.go
├── theme_test.go
├── options.go
├── options_test.go
├── navigation.go
├── navigation_test.go
├── info.go
├── info_test.go
├── widgets.go
├── widgets_test.go
└── assets/fonts/
    ├── Inter-Regular.ttf
    ├── Inter-Bold.ttf
    ├── OFL-Inter.txt
    ├── UbuntuSans-Regular.ttf
    ├── UbuntuSans-Bold.ttf
    └── LICENCE-Ubuntu.txt
docs/
└── options.md
test/windows/
└── README.md
CHANGELOG.md
specs/README.md
```

**Structure Decision**: Extend the existing `gui` package and S042 files so the
headless harness can exercise every state. Add one focused scroll module and
licensed font resources; do not introduce a new package or alter daemon-facing
interfaces.

## Complexity Tracking

No constitution violation or exceptional complexity is introduced.
