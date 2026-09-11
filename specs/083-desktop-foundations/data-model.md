# Data Model: Desktop Visual and Shell Foundations

S083 introduces no persisted entity or schema. The following transient view models define the interaction contract.

## Resolved appearance

- **Saved preference**: `light`, `dark`, or `system`.
- **Resolved palette**: `light` or `dark` based on the saved preference and current operating-system media preference.
- **Transition**: A saved explicit preference resolves immediately; `system` re-resolves whenever the operating-system preference changes.

## Feedback item

- **Message**: Non-empty user-facing text.
- **Tone**: Success, information, warning, or error.
- **Lifetime**: Transient or persistent.
- **Announcement**: Polite for routine information and success; assertive for actionable errors.
- **Dismissal state**: Active, paused by hover, paused by focus, explicitly dismissed, automatically dismissed, or replaced.
- **Identity**: A changing render key derived from the latest message, ensuring repeated outcomes restart the lifetime rather than concatenate.

## Control variant

- **Role**: Primary, secondary, subtle, affirmative, or destructive.
- **Availability**: Enabled, disabled, or pending.
- **Interaction state**: Rest, hover, focus-visible, pressed, or selected where applicable.
- **Accessible identity**: Visible label plus native disabled and busy semantics when relevant.

## Application shell

- **Persistent regions**: Brand, destinations, Exit, target context, connection action, and Appearance.
- **Scrollable region**: One active route page.
- **Viewport modes**: Standard rail and compact horizontal navigation.
