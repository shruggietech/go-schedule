# Research: Structured Desktop Data Tables

## Decision 1: Use structured virtualized rows with a fixed header

**Decision**: Compose a persistent header and a `widget.List` whose reusable row widget lays out multiple labeled cells.

**Rationale**: The pinned Fyne 2.8.1 `widget.List` already provides virtualized row reuse, full-row selection, hover, focus, keyboard navigation, and vertical scrolling. Placing the header outside the scrolling body makes it unconditionally fixed. It also preserves the existing task double-tap forwarding pattern and Activity selection behavior.

**Alternatives considered**:

- `widget.Table`: provides sticky headers and cell virtualization, but its built-in hover, keyboard highlight, and selection operate on one cell, not the full record row, and its content model permits horizontal scrolling. Correcting those behaviors would require fighting private renderer state.
- Static grid inside a vertical scroll: easy to style but eagerly creates every cell and abandons virtualization for long Activity histories.
- Three view-specific row layouts: initially small but duplicates column sizing, truncation, fallback, row state, and theme behavior and invites drift between issues #112 and #113.

## Decision 2: Fit columns to available width

**Decision**: Use one shared width allocator driven by per-column minimum shares and expansion weights. The allocator returns widths whose sum, including gaps, equals the available row width. Cells use single-line ellipsis.

**Rationale**: This satisfies the explicit no-horizontal-scroll constraint at any supported width and guarantees header/body alignment because both use the same column definitions. Minimum shares protect identity and semantic columns; flexible text columns receive remaining space.

**Alternatives considered**:

- Content-measured fixed widths: long values force overflow or steal space unpredictably.
- User-resizable columns: introduces persistent preferences and horizontal overflow outside the approved scope.
- Breakpoint-based column hiding: can work, but permanently visible compact columns provide more information at the application's existing minimum width and avoid breakpoint-specific accessibility behavior.

## Decision 3: Expose complete values through view-appropriate disclosure

**Decision**: Ellipsized Tasks and Schedule cells expose a selectable read-only full-row summary below their lists. Activity continues to open its existing detail dialog, which already contains complete values and diagnostic context.

**Rationale**: Fyne has no stable native tooltip contract across desktop/touch/keyboard input. A visible disclosure works with mouse and keyboard, remains selectable, and does not require hover. It avoids changing task data or opening the edit modal merely to read a clipped value.

**Alternatives considered**:

- Hover-only tooltip: excludes keyboard and touch users and would require a custom overlay lifecycle.
- Multi-line table rows: destroys scan alignment and makes 100-row navigation inefficient.
- Horizontal scrolling: explicitly prohibited by both issues.

## Decision 4: Normalize presentation before widget binding

**Decision**: Map source records into immutable row models with stable identity, ordered cells, full summary, glyph, and semantic importance.

**Rationale**: Pure mapping makes every known/unknown value testable, keeps widget update callbacks simple, and ensures the header/cell contract cannot vary by row. It also separates user-facing `Enabled` from lifecycle `Active`, eliminating the original apparent contradiction without changing domain state.

**Alternatives considered**:

- Format directly in list callbacks: repeats existing concatenated-string coupling and makes fallbacks and semantics harder to test.
- Change domain values: would alter API/persistence contracts for a presentation problem.

## Decision 5: Reuse semantic importance and add one quiet row surface

**Decision**: Use Fyne label importance mapped through the existing S044 theme roles for informational, warning, error, success, and disabled text/glyphs. Add only one theme-aware alternate-row background role if the default background does not provide visible striping.

**Rationale**: S044 already proves semantic palette contrast and state overlays. Reusing these roles keeps Schedule and Activity consistent and avoids the harsh bespoke colors explicitly rejected by the user. Alternation is structural, not semantic, so it stays deliberately subtle.

**Alternatives considered**:

- Hard-coded per-widget colors: do not follow appearance changes and duplicate palette knowledge.
- Bright filled status badges: over-amplify contrast and compete with selection/hover.
- Color-only status: fails the acceptance requirement; every semantic cell retains a label and distinct glyph.

## Decision 6: Preserve ordering and scope boundaries

**Decision**: Keep Tasks in backend snapshot order, Schedule ascending by time, and Activity newest-first. Do not add sorting, new source fields, or richer error capture.

**Rationale**: The issues ask for understandable presentation and expressly coordinate, but do not authorize new sorting or backend contracts. Existing controls and behavior remain independently verifiable.

**Alternatives considered**:

- Clickable sorting headers: valuable future work but expands keyboard, persistence, selection, and test scope.
- Enrich Activity rows from captured process output: belongs to #102.
- Redesign command display or editing: belongs to #110.

## Decision 7: Preserve the mandated spec-kit order despite a tooling defect

**Decision**: Generate the UX requirements checklist from the feature path returned by `check-prerequisites.ps1 -Json -PathsOnly`, then run plan setup.

**Rationale**: The installed checklist prerequisite rejects a missing plan even though both the command skill and autopilot protocol require checklist before plan. Using the already validated feature path preserves the governing order without modifying unrelated project tooling.

**Alternatives considered**:

- Run plan before checklist: violates the explicit autopilot sequence.
- Patch the spec-kit prerequisite inside S045: unrelated workflow-tooling scope and review noise.
- Skip checklist: violates autopilot completion requirements.
