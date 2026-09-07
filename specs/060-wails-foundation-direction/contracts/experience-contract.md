# Calm Operations Experience Contract

## Composition

The control center uses four stable regions in top-to-bottom reading order:

1. A narrow application rail contains product identity and primary destinations.
2. A target bar names `This computer`, connection state, platform, and the target-switch action.
3. The page header states the current job-to-be-done and exposes one primary action.
4. The content region uses a list plus contextual inspector when detail is useful, collapsing to one region at compact width.

The interface does not begin with a dashboard of decorative metrics. Summary values appear only when they support a decision or route to the underlying work.

## Representative views

| View | Required visible decisions |
| --- | --- |
| Tasks | Active target, task state, schedule, next run, last result, create action, selected task context |
| Task editor | New versus edit mode, active target, command and schedule hierarchy, validation, save consequence |
| Schedule | Active target, temporal range, grouped upcoming work, task relationship |
| Activity | Active target, run outcome, time, duration, exit status, output disclosure |
| Target switcher | Current target, available targets, connection condition, mutation consequence |
| Compact | Target identity, destination access, primary action, and one-column detail flow |

## Required states

New, empty, loading, connected, disconnected, degraded, destructive, validation-error, and success conditions each have:

- a visible text label;
- a semantic icon or shape where appropriate;
- concise explanation;
- a primary recovery or continuation action when one exists;
- no meaning conveyed by color alone.

## Tokens

The proof consumes the canonical brand roles: Night, Panel, Raised, Line, Text, Muted, Paper, Ink, Interval, Anchor, Hold, and Stop. Light surfaces use the approved accessible light variants. Space Grotesk is display-only, Geist is body and control text, and Geist Mono is limited to commands, identifiers, schedules, timestamps, and exit data.

Spacing uses the 4, 8, 12, 16, 24, 32, and 48 pixel scale. Radii use 4 or 8 pixels except pills. Elevation is quiet and never substitutes for a border or heading. Motion uses 120, 180, or 240 milliseconds on the approved easing curve.

## Accessibility

- Landmarks include application navigation, target context, main content, and contextual detail where present.
- One H1 identifies every page; headings do not skip levels.
- Every control has an accessible name and visible label unless an established icon-only control has equivalent accessible text and tooltip.
- Keyboard focus follows visual reading order, remains visible at all times, and never becomes trapped.
- Escape closes dismissible overlays; focus returns to the invoking control.
- Normal text reaches 4.5:1 contrast; large text, meaningful graphics, and focus indicators reach 3:1.
- Targets reach 24 by 24 CSS pixels or have equivalent spacing under WCAG 2.2 exceptions.
- At 200 percent zoom and at 900 by 650 CSS pixels, primary workflows reflow without horizontal page scrolling.
- Status changes use an appropriate live region without repeatedly announcing unchanged content.
- `prefers-reduced-motion: reduce` removes nonessential transitions and animated movement.

## Review heuristic

For each representative view, a reviewer should identify within 10 seconds:

1. Which daemon is active?
2. What page or task is this?
3. What needs attention?
4. What is the primary action?
5. What happens next?

Failure on any question is a hierarchy defect, not a request for more decorative content.
