# UI Contract: Structured Column Adjustment

## Boundary discovery

- Every adjacent pair in Schedule list and Activity has a visible separator and
  horizontal-resize pointer.
- Keyboard focus reaches each separator in header order.
- Its accessible label is `Resize <left> and <right> columns`.

## Pointer contract

- Horizontal drag transfers width only between the named adjacent columns.
- Layout updates during drag; the preference commits at drag end.
- Minimum constraints stop motion cleanly without overlap.

## Keyboard contract

- Left Arrow moves the boundary one 10-unit logical step left.
- Right Arrow moves it one 10-unit logical step right.
- Unsupported keys do nothing; each accepted step persists immediately.

## Reset and persistence contract

- Each affected list toolbar includes **Reset columns**.
- Reset applies and stores only that view's defaults immediately.
- Schedule and Activity have distinct stable preference identities.
- Stored proportions are normalized, not raw pixels.
- Any incompatible value falls back atomically to defaults.
- Narrow windows compress proportionally without horizontal scrolling.
