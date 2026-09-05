# Research: Persisted Adjustable Columns

## Decision 1: Persist versioned normalized proportions

**Decision**: Store one JSON value per view containing a schema version, ordered column identities, and normalized positive proportions.

**Rationale**: Proportions preserve intent across window size, DPI, and font changes. Identities and a version reject stale schemas atomically.

**Alternatives considered**: Raw pixels do not adapt; separate float keys permit partial state; daemon storage gives presentation state the wrong owner.

## Decision 2: Transfer width between adjacent columns

**Decision**: A boundary change increases one adjacent proportion and decreases the other equally, clamped to logical minimums when space permits.

**Rationale**: The boundary tracks the pointer, unrelated columns remain stable, and total width is conserved. Narrow views retain proportional compression.

**Alternatives considered**: Rebalancing all columns causes surprising motion; unbounded shrink harms usability; horizontal scroll is explicitly out of scope.

## Decision 3: Use one focusable boundary control

**Decision**: Each header boundary is a narrow custom widget with a separator, resize cursor, accessible label, drag events, focus, and arrow keys.

**Rationale**: A real focusable object is discoverable and operable without a pointer, while one implementation prevents accessibility drift.

**Alternatives considered**: Invisible hit targets are undiscoverable; separate keyboard settings duplicate behavior; replacing the list would regress its selection, virtualization, disclosure, and responsive contracts.

## Decision 4: Persist completed interactions, refresh live

**Decision**: Apply drag changes live but persist at drag end. Persist each accepted keyboard step immediately.

**Rationale**: This provides feedback without preference write amplification.

**Alternatives considered**: Per-event writes are unnecessary; completion-only refresh makes dragging feel disconnected.

## Decision 5: Reset locally in each affected view

**Decision**: Add **Reset columns** to both affected view toolbars.

**Rationale**: The action is discoverable in context and preserves isolation.

**Alternatives considered**: A global reset obscures scope; context-menu-only access is less discoverable and less keyboard-friendly.
