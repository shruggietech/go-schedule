# Responsive Workflows Interaction Contract

## Task dialog

- Create task and Edit open one modal TaskEditor without changing underlying Tasks page height or scroll position.
- Initial focus enters the first task field, Tab and Shift+Tab remain contained, Escape and Cancel close without mutation, and closure restores focus to the exact available invoker.
- The dialog body owns overflow while the heading and action group remain visible.
- Common fields precede a collapsed Advanced settings disclosure.
- Preview and Save remain distinct actions; pending state prevents duplicate activation of only the active task operation.
- Validation keeps the draft open, scrolls the invalid control into view, and focuses it.

## Form grid

- FormGrid places each Field label above its control.
- Auto columns use available width and stack to one column before any label, control, or action collides.
- DOM order is visual order and therefore keyboard order.

## Description list

- DescriptionList preserves native `dt` and `dd` semantics.
- Each item remains a label-value pair at wide widths and stacks without overlap at constrained widths.
- Literal values wrap anywhere when required and remain selectable.

## Status label and card section

- StatusLabel includes readable text and never depends on color alone.
- CardSection groups one heading and its related description, status, content, and actions with shared panel spacing.

## Path display and copy action

- PathDisplay renders the complete literal path in compact monospace text with a readable fallback when unavailable.
- A copyable record exposes one action whose accessible label identifies the path record.
- Only the matching action may enter pending or copied presentation.
- Different record operations may overlap; each result is associated by record identity and cannot relabel another action.

## Progressive disclosure

- Inactive localhost HTTP configuration and ordinary remote pairing are collapsed initially and expose explicit summary controls.
- Active localhost status and lifecycle actions remain visible.
- Selecting Repair opens the pairing disclosure for the selected profile.

## Responsive boundary

- Affected pages produce no document-level horizontal overflow at 1280 by 800, 800 by 600, or 200 percent zoom.
- Focused controls remain visible within the active page or dialog scroll container.
- Light, dark, and follow-system appearances retain readable text, boundaries, and focus indicators.
