# Research: GUI Activity Clarity

## Decision 1: Canonical user-facing name

**Decision**: Use `Activity` for the mixed destination and its badge.

**Rationale**: The view merges daemon log records and scheduler alerts, so `Logs` is narrower than the content. `Activity` describes both without changing their underlying domain types.

**Alternatives considered**: `Events` could be confused with the internal live event stream; `History` suggests durable completeness that the current-view clear action does not provide.

## Decision 2: Badge boundary

**Decision**: Show exact counts through 99 and `99+` at 100 or above.

**Rationale**: This is the exact boundary requested by issue #26 and keeps tab width stable. Treating negative values as zero makes the formatter total even though the view model should not produce them.

## Decision 3: Clear presentation

**Decision**: Use `Clear View`, Fyne's `ContentClearIcon`, and persistent help text.

**Rationale**: The action changes the current view and acknowledges visible alerts but does not delete persisted records. The selected wording and icon do not imply deletion. Persistent copy works without hover and uses an icon already provided by the pinned Fyne dependency.

**Alternatives considered**: `Dismiss All` plus a delete icon retains the current ambiguity. A custom tooltip adds a new interaction pattern and excludes non-hover environments. A confirmation dialog adds ceremony to a reversible view action.

## Decision 4: Internal naming

**Decision**: Keep internal `logsTab`, `buildLogsTab`, and `logEntry` names.

**Rationale**: Those identifiers describe existing implementation structures and are invisible to users. Renaming them would enlarge the diff without improving the requested experience.
