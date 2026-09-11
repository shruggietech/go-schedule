# Desktop Foundations Interaction Contract

## Shared button

- Supported semantic variants are primary, secondary, subtle, affirmative, and destructive.
- Every variant provides visibly different rest, hover, focus-visible, pressed, disabled, and pending states without removing its accessible name.
- Ordinary actions use the compact control dimension; specialized large targets must opt in explicitly.
- Pending actions expose `aria-busy` and remain unavailable to repeated activation.

## Dialog

- A dialog moves initial focus to its first enabled action, contains Tab and Shift+Tab focus, closes on Escape, and restores focus to its exact invoker.
- Dialog actions are grouped in one `.dialog-actions` region with at least 8 CSS pixels between actions.
- The caller supplies semantic actions; the shared dialog supplies a safe Close action.

## Persistent notice

- A persistent notice uses `status` for information and success or `alert` for warning and error.
- A notice with a dismissal callback exposes a compact action named `Dismiss <title>`.
- Dismissing removes the notice from layout and the accessibility tree.

## Transient toast

- An empty message produces no visible toast.
- A non-empty message produces exactly one overlay toast with a polite atomic live region and a dismiss action.
- A new message replaces the current toast and restarts its five-second lifetime.
- Hover or focus within the toast pauses dismissal; leaving or blurring resumes a fresh readable interval.
- Explicit dismissal removes the toast immediately.

## Application shell

- The `.app` region owns the dynamic viewport and the document itself does not scroll.
- The rail and workspace chrome remain visible while `#main-content` is the only vertical route scroller.
- At or below 900 CSS pixels, destinations may scroll horizontally inside navigation, but brand, Exit, target context, and Appearance remain available.
- `data-appearance` retains saved intent and `data-resolved-appearance` exposes the active light or dark palette.
