# Data Model: Desktop Interaction and Appearance Polish

S044 adds no daemon database entity or migration. It extends bounded per-user GUI preferences and in-memory presentation state.

## Desktop preferences

| Field | Allowed values | Default | Invalid/missing behavior | Persistence |
| --- | --- | --- | --- | --- |
| Mode | `dark`, `light`, `system` | `dark` | `dark` | Per-user application preferences |
| Font | `system`, `brand`, `inter`, `ubuntu`, `monospace` | `system` | `system` | Per-user application preferences |
| Scroll sensitivity | `1.0` through `4.0`, step `0.5` | `2.0` | `2.0` | Per-user application preferences |

Each field normalizes independently. Applying any value constructs one complete preference snapshot, persists it, and refreshes affected application-owned views. Restore defaults writes Dark, System, and 2x together from one user action; the preference store does not promise a multi-key transaction.

## Font family

| Choice | Regular face | Bold face | Monospace semantic face | Symbol face |
| --- | --- | --- | --- | --- |
| System | Framework default regular | Framework default bold | Framework default monospace | Framework default symbol |
| Geist (brand) | Geist Regular | Space Grotesk Bold | Geist Mono | Framework default symbol |
| Inter | Inter Regular | Inter Bold | Geist Mono | Framework default symbol |
| Ubuntu | Ubuntu Sans Regular | Ubuntu Sans Bold | Geist Mono | Framework default symbol |
| Monospace | Geist Mono | Geist Mono | Geist Mono | Framework default symbol |

## Interaction state style

| Field | Meaning |
| --- | --- |
| Palette | Effective dark or light colors after explicit/follow-system selection |
| Importance | Low, medium, high, danger, success, or warning base treatment |
| State | Rest, hover, selected, selected-plus-hover, focus, pressed, or disabled |
| Foreground | Semantic label and themed-glyph color |
| Base background | Importance-specific persistent surface |
| Overlay | Translucent state color blended over the base surface |
| Non-color cue | Persistent fill, focus shape, label/glyph, or disabled behavior |

Every normal-text foreground/composite-background pair must reach 4.5:1. Every essential focus/boundary indicator against its adjacent surface must reach 3:1.

## Storage row view

| Field | Meaning |
| --- | --- |
| Location | Existing truthful `storageLocation` value |
| Category | Stable leading column text |
| Path | Exact selectable available path, or unavailability message |
| Detail | Scope, existence, software-only removal, and explicit-wipe behavior |
| Copy | Far-right action, enabled only for an available path |
| Muted | Whole-row unavailable presentation using disabled semantics |

Storage data and ownership do not change. The row is a projection of the existing location; it performs no path creation, traversal, or deletion.

## Choice selector state

| Field | Meaning |
| --- | --- |
| Current | Closed-control value and effective persisted selection |
| Alternatives | Every allowed value except Current, in canonical order |

Transition:

```text
Current A + alternatives [B, C]
  -> choose B
Current B + alternatives [A, C]
  -> invalid/empty external preference
Default + alternatives excluding Default
```

## Sensitive vertical scroll

| Field | Meaning |
| --- | --- |
| Content | One application-owned vertically scrollable object |
| Sensitivity source | Function returning the latest normalized preference |
| Incoming delta | Framework logical-pixel scroll vector |
| Discrete threshold | Absolute vertical delta of 25 logical pixels |
| Effective delta | Incoming value unchanged below threshold, otherwise multiplied by sensitivity |

The copied event is delegated once to the stock scroller. Keyboard, drag, scrollbar, and programmatic offset changes never transition through the custom delta function.

## Navigation shell extension

The S042 destination model is unchanged. One full-height separator becomes the boundary between rail and content. Exit retains no destination ID or selected state and changes only from medium to danger importance.
