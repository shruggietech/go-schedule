# Shell and Component Contract

## Shell

The Shell owns landmarks, skip navigation, active route, persistent target identity, page title, appearance selector, connection notice placement, notification region, and orderly exit. Route content is supplied as children and cannot replace those landmarks.

## Primitive Catalog

| Primitive | Required contract |
| --- | --- |
| Button | Native button semantics; primary, secondary, quiet, and destructive variants; disabled and busy states |
| Link | Native link semantics for navigation; external behavior is announced when applicable |
| StatusBadge | Text, shape, and color communicate a closed state; status does not rely on color alone |
| Notice | Named title and detail; optional relevant action; polite status unless urgent |
| StatePanel | Named state, concise guidance, and optional recovery or continuation action |
| Field | Programmatic label, optional help, error association, required and disabled semantics |
| DataTable | Caption, headers, keyboard-reachable row actions, and an explicit empty state |
| Disclosure | Native expanded state and keyboard behavior; content remains in reading order |
| Dialog | Modal name, initial focus, Escape and Close behavior, focus containment, and exact-invoker return |
| ToastRegion | Polite live region with stable message identity and bounded retention |

## Interaction Rules

- Page changes move focus to the main content heading without losing persistent target context.
- Dialog close by button or Escape returns focus to the exact connected invoking element when it still exists.
- Busy buttons retain their accessible name and expose busy state while preventing duplicate activation.
- Error text is associated with its field and does not disappear on focus.
- Disabled controls remain visually and semantically distinct from unavailable application state.
- Reduced-motion mode removes nonessential transitions.
- Compact and zoomed layouts reflow vertically and never create page-level horizontal overflow.

## Styling Rules

Tokens define color, typography, spacing, radius, elevation, focus ring, and motion. Components consume tokens rather than hard-coded visual values. System appearance uses the operating-system preference, and explicit light or dark selection overrides it for the session.

## Boundary Rules

Components are presentational. They receive safe models and callbacks, never import generated Wails bindings, connection transport, fixture data, or backend errors. The frontend connection adapter is the only consumer of generated bridge and runtime event APIs.
