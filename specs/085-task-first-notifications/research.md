# Research: Task-First Notifications

## Initial information hierarchy

**Decision**: Put current status, configured coverage, recent task outcomes, and the next action in the first visible region, before transport or policy administration.

**Rationale**: Users enter a destination named Notifications to understand notification state and results. Configuration is a secondary task and should not obscure the destination's primary answer.

**Alternatives considered**: Retaining the two-column setup grid with an overview banner, which still makes specialist controls dominant; a separate dashboard route, which fragments one small capability; nested sidebar destinations, which adds navigation complexity for three related tasks.

## Progressive disclosure

**Decision**: Use three existing native Disclosure instances for Destinations, Assignment rules, and Delivery diagnostics, all collapsed initially.

**Rationale**: Native details and summary semantics provide keyboard activation, exposed state, and simple responsive flow. Each label describes a user task rather than an implementation object.

**Alternatives considered**: Tabs, which hide content through custom selection state and complicate narrow-height use; accordions that permit only one section, which obstruct cross-reference; custom popovers, which are unsuitable for complete workflows.

## Coverage projection

**Decision**: Extend the desktop-only workspace with secret-free task and group coverage derived from existing direct and effective assignment methods. Preserve the rest of the snapshot if one coverage lookup fails and mark coverage incomplete.

**Rationale**: The initial surface cannot truthfully identify configured tasks and groups from channels and delivery history alone. Existing methods are authoritative and require no daemon API or persistence change. Explicit incompleteness avoids turning optional explanatory projection into total page failure.

**Alternatives considered**: Inferring configuration from deliveries, which misses never-fired assignments; loading every policy from the browser, which duplicates orchestration and produces incremental visual churn; adding a daemon bulk endpoint, which expands public contracts beyond this UI slice.

## State guidance

**Decision**: Define deterministic, color-independent copy for setup needed, disabled, healthy, queued, sending, retrying, successful, and failed states. Failed and disabled states point to a relevant disclosure; healthy and successful states explicitly say no action is needed.

**Rationale**: Plain language and one next step reduce specialist interpretation and provide stable assertions for tests and assistive technology.

**Alternatives considered**: Badge-only states, which rely on terminology and color; arbitrary free-form messages, which drift; tooltips alone, which hide essential status.

## Recent versus diagnostic evidence

**Decision**: Show the five newest results on the initial surface and retain the existing full 200-record table, filters, and selected detail under diagnostics.

**Rationale**: Five results answer recency without overwhelming an 800 by 600 viewport. Detailed history remains available without losing filtering or redacted correlation evidence.

**Alternatives considered**: Showing all results initially, which recreates excessive page length; task outcomes only, which conceals channel-test evidence; pagination, which adds state and contract complexity without changing the existing bounded source.

## Verification depth

**Decision**: Combine Go projection tests, Vitest interaction tests, Playwright comprehension, reflow and accessibility checks, native Wails build, and canonical verification.

**Rationale**: Service tests establish accurate and secret-free summaries, component tests establish disclosure and interaction state, and browser checks establish actual geometry, focus, and accessibility.

**Alternatives considered**: Screenshot review alone, which cannot prove keyboard state or semantics; attended release qualification in this slice, which belongs after all #228 children are merged.
