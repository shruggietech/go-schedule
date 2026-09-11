# Data Model: Responsive Task and Administration Workflows

S084 adds no persisted schema or backend entity. The following transient view models define presentation and interaction state.

## Focused task workflow

- **Mode**: Create or edit.
- **Draft**: The existing complete task draft and original revision metadata.
- **Invoker**: The exact Create task or Edit element that opened the workflow.
- **Action state**: Idle, previewing, or saving.
- **Result**: Validation failure, stale conflict, exact command preview, schedule preview, or accepted save.
- **Disclosure**: Common fields visible and advanced fields collapsed or expanded.
- **Transition**: Invocation opens the modal; Cancel or Escape closes without mutation; accepted Save closes after publishing the result; validation retains the modal and focuses the invalid field.

## Administration section

- **Purpose**: Status, details, configuration, or actions.
- **Heading**: Eyebrow, title, and optional summary.
- **Status**: A text label that remains meaningful without color.
- **Content**: Responsive form grid, definition list, prose, or action group.
- **Disclosure**: Visible or deliberately collapsed based on whether an active task exists.

## Definition item

- **Label**: Stable human-readable term.
- **Value**: Text or semantic content.
- **Literal flag**: Whether the value uses code typography and aggressive wrapping.
- **Availability**: Present, absent, unavailable, or not yet observed.

## Path presentation

- **Record identity**: Existing storage record identifier.
- **Label**: Human-readable storage label or contextual path purpose.
- **Literal path**: Complete selectable value or an explicit unavailable fallback.
- **Copy state**: Idle, pending, copied, or failed for this record only.
- **Transition**: Copy activation changes only the matching record state; completion is applied by stable record identity even when other actions overlap.

## Disclosure state

- **Task advanced settings**: Collapsed whenever create or edit opens.
- **Inactive localhost configuration**: Collapsed on page entry.
- **Ordinary remote pairing**: Collapsed on page entry.
- **Repair pairing**: Expanded when a specific saved profile is selected for repair.
