# Research: Desktop Visual and Shell Foundations

## Compact shared controls

**Decision**: Use a 32 CSS-pixel ordinary control height, compact vertical padding, at least 24 by 24 CSS-pixel interactive targets, an 8-pixel spacing unit between adjacent actions, and distinct text, fill, border, focus-ring, and pressed-state changes for each semantic variant.

**Rationale**: The observed 41.6 CSS-pixel minimum button height can be reduced by about 23 percent to a conventional compact desktop control while satisfying the issue's target floor. Multiple visual channels preserve meaning and meet WCAG contrast expectations.

**Alternatives considered**: Keep current sizing and only add hover, which does not address density; reduce below 24 pixels, which weakens target accessibility; introduce Fluent UI React, which adds a large dependency and migration boundary for a focused repair.

## Viewport and navigation ownership

**Decision**: Make the application grid exactly own the dynamic viewport, give the rail and workspace a zero-minimum-height constraint, and put overflow on the main page region only. Preserve all navigation at narrow sizes through a compact horizontal destination strip while keeping Exit visible.

**Rationale**: Microsoft NavigationView guidance treats navigation and header as persistent app structure. A single deliberate page scroller avoids the current body-scroll failure and makes 800 by 600 behavior testable.

**Alternatives considered**: Sticky rail content inside body scroll, which still permits document and footer drift; hide Exit at narrow sizes, which directly violates #230; create independent scroll regions inside each page, which multiplies nested scrolling.

## Feedback lifetime and placement

**Decision**: Treat routine action completion as a single five-second toast positioned above the footer with safe margins. Replace the current toast when a new message arrives, pause dismissal on hover or focus, and use dismissible message bars or field errors for actionable failures.

**Rationale**: Fluent toast guidance reserves transient overlays for short-lived action feedback, while Microsoft InfoBar guidance supports persistent actionable conditions. One active toast prevents the concatenated string observed in the candidate.

**Alternatives considered**: Stack multiple toasts, which increases overlap risk in the constrained desktop viewport; keep status text in the footer, which caused clipping; make every error a toast, which can remove recovery information too early.

## Follow-system identity

**Decision**: Track the resolved system palette with a media-query listener and expose it as application state. Keep the saved preference as `system`, but apply resolved theme data consistently to surfaces and brand-mark treatment.

**Rationale**: The saved intent and current rendered palette are different concepts. Explicit resolution makes live operating-system changes and native evidence deterministic while retaining the user's preference.

**Alternatives considered**: Resolve only at startup, which becomes stale after an OS theme change; add separate dark and light artwork during this repair, which expands brand review scope unnecessarily.

## Verification depth

**Decision**: Use deterministic component and Playwright coverage for all variants, timeouts, focus behavior, scroll ownership, and reflow, followed by a short native Windows observation against one exact development build.

**Rationale**: The defects are primarily shared frontend contracts. Automation should carry the repeatable matrix, while native observation confirms Wails viewport and Windows theme integration without repeating the unrelated 47-observation release ritual.

**Alternatives considered**: Source-only review, which misses native viewport behavior; a full release qualification rerun, which is disproportionate and would mutate the release workflow.
