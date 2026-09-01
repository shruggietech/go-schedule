# Contract: Startup Window

## Inputs

- Selected monitor work-area width and height in physical pixels.
- Positive effective display scale.

## Output

- Restored Fyne size in logical units.

## Rules

1. Unknown work area or non-positive scale returns 1280x800.
2. Known dimensions are divided by scale.
3. Each logical dimension is multiplied by 0.9.
4. Width is `min(1280, capped width)` and height is `min(800, capped height)`.
5. No lower clamp may exceed the cap.
6. The Windows adapter obtains `rcWork` for the monitor nearest the launch pointer.
7. The caller centers the restored window and does not force maximize or fullscreen state.
