# User Experience Requirements Checklist: Agent Grant Controls

**Purpose**: Validate that S094's administrative and accessibility experience is unambiguous before planning
**Created**: 2026-09-20
**Feature**: [spec.md](../spec.md)

## Information Clarity

- [x] The interface answers who, daemon, authority, transport, expiry, and state from one view
- [x] Manage authority uses distinct high-impact wording and non-color-only treatment
- [x] MCP Off and stdio on-demand availability are not conflated
- [x] Non-expiring access is an explicit, deliberate duration choice

## Interaction Quality

- [x] Grant creation requires the minimum user choices and provides practical duration presets
- [x] Narrow, expire, revoke, and recent-audit actions have bounded outcomes
- [x] Destructive revocation requires a clear confirmation without repeated confirmation ceremonies
- [x] Async success and failure outcomes are announced

## Accessibility

- [x] Keyboard operation is required for every grant workflow
- [x] Controls and state have programmatic names
- [x] Dialog focus entry, containment, dismissal, and return are required
- [x] Authority and state remain distinguishable without relying on color

## Notes

- All UX requirements are specific and traceable to user stories or functional requirements.
