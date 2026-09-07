# Preferences, Diagnostics, and Trust Checklist

**Purpose**: Test requirement coverage for migration, storage truth, native actions, and recovery

**Created**: 2026-09-07

**Feature**: [spec.md](../spec.md)

## Preference Transition

- [x] Every durable legacy preference has a migrate, default, or retire decision.
- [x] New-install and damaged-file defaults are defined.
- [x] One-time migration authority and repeat-start behavior are defined.
- [x] Persistence validation, atomic replacement, and write failure behavior are defined.
- [x] The transition is disclosed in Settings and user documentation.

## Storage Truth

- [x] Application, daemon, desktop, executable, documentation, and platform records are enumerated.
- [x] Ownership, scope, existence, normal removal, and explicit wipe are required for every row.
- [x] Daemon paths cannot be guessed when runtime metadata is unavailable.
- [x] External data cannot be claimed for removal.
- [x] Copy actions require a current backend-resolved path.

## Native Action Boundaries

- [x] Arbitrary frontend URLs and clipboard text are rejected by design.
- [x] Product destinations are fixed backend-defined HTTPS links.
- [x] Clipboard and browser failures produce concise actionable status.
- [x] Duplicate pending actions are suppressed where repeated activation matters.

## Connection Recovery

- [x] Every supported diagnosis has explicit inline guidance.
- [x] Retry behavior reuses the existing connection manager.
- [x] Retry preserves route and focus and suppresses duplicate work.
- [x] Status changes do not trigger repeated modal interruption.
- [x] Local settings remain usable offline.

## Accessibility and Scale

- [x] Every state is conveyed without color alone.
- [x] Controls and announcements are keyboard and screen-reader accessible.
- [x] Zoom coverage spans 80 through 200 percent.
- [x] Missing, unavailable, failure, and success states are explicit.

## Notes

- Every item maps to an acceptance scenario, functional requirement, edge case, or measurable outcome.
