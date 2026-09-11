# Research: Responsive Task and Administration Workflows

## Focused task surface

**Decision**: Reuse the shared modal Dialog for both create and edit, place TaskEditor content inside its scrollable body, and supply Preview and Save as dialog actions with Cancel as the safe close action.

**Rationale**: One workflow prevents create and edit drift, preserves the underlying Tasks geometry, and inherits the focus containment and return behavior delivered in S083. A bounded dialog with an internal scroller fits 800 by 600 without turning the page into the editor.

**Alternatives considered**: A separate route, which loses list context; a side sheet, which has no existing shared accessibility contract; retaining the inline card with compact CSS, which does not satisfy #231.

## Administration composition

**Decision**: Add lightweight semantic React wrappers for responsive form grids, definition lists, section cards, status labels, and literal paths, backed by one CSS contract.

**Rationale**: Agent Access, Connections, Pairing, and Settings repeat the same markup problems. Shared composition makes alignment and reflow testable without importing a full design system or coupling domain pages.

**Alternatives considered**: Page-specific selectors, which would repeat the defect source; Fluent UI React, which is disproportionate for four existing pages and adds a dependency; a schema-driven form engine, which over-abstracts static forms.

## Disclosure hierarchy

**Decision**: Keep safety status, diagnosis, and common actions visible. Collapse inactive localhost configuration and ordinary remote pairing using native Details behavior; automatically open pairing only when a repair target is already selected.

**Rationale**: This meets progressive-disclosure requirements without hiding an active recovery action. Native Details provides keyboard activation and an exposed expanded state with minimal code.

**Alternatives considered**: Always-visible advanced forms, which caused the current density failure; a new nested navigation system, which would overlap #232 and increase routing scope; always-collapsed repair, which hides the action the user just requested.

## Settings action isolation

**Decision**: Track in-flight settings operations by stable action keys and maintain record-scoped copy outcome state. Each Copy path button derives its pending and completion label only from its record key.

**Rationale**: A global pending boolean causes every button to restyle and cannot represent overlapping results correctly. A bounded set of action keys preserves concurrency and stable identity while keeping state local to the settings hook.

**Alternatives considered**: Keep global serialization and conceal it visually, which creates ignored clicks; add a global state library, which is unnecessary; place asynchronous bridge calls in every card, which fragments error handling.

## Long operational values

**Decision**: Show complete paths in selectable monospace blocks and apply overflow wrapping to other operational identifiers and endpoints.

**Rationale**: Operators need the exact value for diagnosis. Wrapping preserves complete access and eliminates horizontal overflow at constrained widths.

**Alternatives considered**: Ellipsis-only truncation, which hides required data; horizontal card scrolling, which introduces two-dimensional navigation; smaller proportional text, which does not clearly identify literals.

## Verification depth

**Decision**: Combine focused Vitest interaction tests, Playwright workflow and reflow checks, the existing native Wails production build, and canonical eight-gate verification.

**Rationale**: Component tests cover deterministic focus and state identity; browser checks prove CSS geometry, zoom, keyboard, and accessibility; the native build proves the exact frontend remains packageable without repeating release promotion.

**Alternatives considered**: Screenshot-only review, which cannot prove focus or state isolation; a full attended release rerun inside this code slice, which belongs after all remediation children finish.
