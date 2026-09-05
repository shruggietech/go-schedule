# Contract: Desktop Interaction and Appearance Polish

## Interactive-state contract

1. The theme supplies dark and light semantic surfaces plus translucent hover, pressed, and focus overlays compatible with stock control blending.
2. Normal text reaches 4.5:1 against its final composite background for every supported importance and interaction state.
3. Essential focus, boundary, and glyph indicators reach 3:1 against adjacent surfaces.
4. A selected destination remains persistently filled while hovered or focused; its corresponding content and label continue to identify selection without relying on the transient state color alone.
5. Disabled controls remain labeled, visibly muted as a complete control or row, and non-actionable.

## Appearance contract

1. Missing, empty, unknown, or reset font state resolves to System.
2. Explicit `brand`, `system`, `inter`, `ubuntu`, and `monospace` identifiers restore unchanged.
3. The user-facing order is System, Geist (brand), Inter, Ubuntu, Monospace.
4. System delegates all text faces to the framework. Inter and Ubuntu use their packaged regular/bold faces, proportional families retain Geist Mono for explicitly monospace text, and symbols always delegate.
5. A font or color-mode change applies to existing and future views/dialogs in the same application instance.
6. Info continues to show `Version <buildinfo.Version>` and the unchanged attribution, centered and unwrapped.

## Selector contract

1. A closed appearance selector displays its current value.
2. Opening it offers every valid alternative exactly once and omits the current value.
3. Choosing an alternative applies it, makes it current, and rebuilds the menu so the former value becomes an alternative.
4. Pointer and keyboard choice paths produce the same persisted result.

## Storage presentation contract

1. A header labels Category, Location and removal details, and Action.
2. Each known location occupies one compact aligned row with Copy at the far right and no horizontal scroll container.
3. Long paths and details wrap vertically and keep the exact path selectable.
4. Available rows copy the exact existing resolved string.
5. Unavailable rows use disabled semantic treatment across category, path, details, and action; selection and Copy are disabled.
6. Existing scope, existence, ownership, and removal wording remains truthful and visible.

## Navigation contract

1. The rail keeps its content-derived minimum of at least 168 logical pixels, symmetric horizontal insets, current destination order, and Activity badge.
2. A full-height semantic separator defines the rail's trailing boundary.
3. Exit remains bottom-right anchored below the destination separator, carries semantic danger label/glyph treatment, and never becomes selected content.
4. Exit and title-bar close share the existing exactly-once orderly shutdown.

## Scroll contract

1. Sensitivity accepts 1x through 4x in 0.5 steps, defaults to 2x, persists per user, normalizes invalid state to 2x, and resets with other Options defaults.
2. Every direct application-owned vertical scroll container uses the same live preference source.
3. Absolute deltas below 25 logical pixels pass through unchanged. Deltas at or above 25 are multiplied once by sensitivity and keep their direction.
4. Horizontal delta is never introduced; horizontal containers and scrollbars are not added.
5. Keyboard, precision-delta, drag, scrollbar, and programmatic scrolling keep stock behavior.

## Evidence contract

1. Headless tests cover theme composites, contrast thresholds, preference migration/defaults, font resources, selector alternatives, row states and geometry, navigation boundary/Exit semantics, and scroll delta calculation.
2. The canonical eight-gate aggregate passes before publication.
3. Native after-change claims require the exact Windows candidate at 100% and a scaled-DPI setting. Physical precision-touchpad evidence is recorded as unavailable when the attended machine lacks one.
