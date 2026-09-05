# UI Contract: Structured Desktop Data Tables

## Shared contract

1. Each view exposes a fixed header followed by a vertically virtualized body.
2. Header and row cells use the same ordered column definitions and exact allocated widths.
3. The total allocated width never exceeds the available body width; there is no horizontal scrolling.
4. Cells are single-line and ellipsized. Complete values are available through keyboard-accessible, non-mutating disclosure.
5. Empty data leaves headers visible and creates zero data rows.
6. Missing values use view-specific fallback text. Unknown enums are readable and neutral, never silently mapped to a known state.
7. Alternating rows are subtle and structural. Hover, focus, and selection remain distinguishable on either row surface.
8. Semantic status uses both a glyph and a normalized text label. Color supplements those cues.
9. Normal text and essential indicators meet 4.5:1 and 3:1 contrast respectively across supported appearances and overlapping row states.

## Tasks contract

Header order:

```text
Task | Enabled | Lifecycle | Time zone | Group
```

- `Enabled` describes whether scheduling is currently allowed.
- `Lifecycle` describes the stored lifecycle (`Active`, `Completed`, or `Disabled`).
- No status value is enclosed in unexplained square brackets.
- Single click and keyboard navigation select the whole task row.
- Toolbar actions resolve the selected stable task ID against the latest snapshot.
- Double activation resolves and edits the bound stable task ID.
- Selection survives reorder/update by identity and clears when that identity disappears.

## Schedule List contract

Header order:

```text
When | Task | Event | Outcome
```

| Source state | Event | Outcome label | Glyph | Semantic role |
| --- | --- | --- | --- | --- |
| Future | SCHEDULED | N/A Not available | ▷ | Informational |
| Past success | COMPLETED | SUCCESS | ✓ | Success |
| Past failure | COMPLETED | FAILURE | ✗ | Error |
| Past skipped | COMPLETED | SKIPPED | ↷ | Disabled/secondary |
| Past caught up | COMPLETED | CAUGHT UP | ↻ | Informational |
| Past queued | COMPLETED | QUEUED | ⋯ | Warning |
| Past missing | COMPLETED | N/A Not available | • | Neutral |
| Unknown | COMPLETED | normalized source value | ? | Neutral |

- Existing ascending chronological order, live refresh, range selection, and Calendar switching remain unchanged.
- Past rows use the Calendar response's optional stored run ID as their stable identity. Equal-time run records must not depend on query order or row ordinal.
- Future computed rows, which do not yet have a run, use task identity plus scheduled timestamp as their deterministic fallback.
- Selecting a list row exposes all complete values in a read-only disclosure.

## Activity contract

Header order:

```text
When | Severity | Source | Summary
```

| Source severity | Label | Glyph | Semantic role |
| --- | --- | --- | --- |
| info or empty | INFO | • | Informational |
| warning | WARNING | ⚠ | Warning |
| error | ERROR | ✗ | Error |
| other | normalized uppercase source value or UNKNOWN | ? | Neutral |

- Existing newest-first ordering, severity filters, clear/acknowledge behavior, and live refresh remain unchanged.
- Activating a row opens the existing complete Activity detail and clears transient visual selection.

## Responsive priorities

| View | Protected first | Flexible expansion | Protected last |
| --- | --- | --- | --- |
| Tasks | Task, Enabled, Lifecycle | Task, Group | Time zone, Group |
| Schedule | When, Task, Event | Task | Outcome |
| Activity | When, Severity | Summary | Source |

All columns remain present at supported application widths. Primary/flexible text receives any surplus after protected shares are allocated.
