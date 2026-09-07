# Operational Truth and UX Checklist

**Purpose**: Test the completeness, accuracy, and accessibility of Schedule and Activity requirements

**Created**: 2026-09-07

**Feature**: [spec.md](../spec.md)

## Operational Truth

- [x] Predictions cannot be mistaken for persisted runs.
- [x] Runs, logs, and alerts retain visible record types.
- [x] Every supported run state has a textual representation.
- [x] Missing task, run, output, exit status, and log-path data have explicit unavailable behavior.
- [x] Daemon scheduling, retention, and acknowledgement remain authoritative.

## Navigation and Inspection

- [x] Agenda, calendar, and range behavior are defined.
- [x] Calendar counts and selected-day inspection are accessible.
- [x] Activity search and type, severity, and outcome filters are defined.
- [x] Run provenance and retained output requirements are complete.
- [x] Empty, loading, stale-selection, and offline states are explicit.

## Live Updates

- [x] Relevant event kinds are identified.
- [x] Burst debouncing and stale-response rejection are required.
- [x] Range, view, filters, focus, and stable identity are preserved.
- [x] Failed constituent reads cannot publish partial state.
- [x] Last complete snapshots remain read-only during disconnection.

## Safe Actions

- [x] Individual acknowledgement changes only one alert.
- [x] Clear View acknowledges only visible alerts.
- [x] Clear View is explicitly local and non-destructive.
- [x] Later activity remains visible after clearing.
- [x] Duplicate pending mutations are suppressed.

## Accessibility and Scale

- [x] Status never relies on color alone.
- [x] Tables, filters, days, detail, and actions are keyboard accessible.
- [x] Screen-reader announcements cover result and mutation changes.
- [x] At least 100 rows per workspace are required in automated coverage.
- [x] The interaction target is measurable at two seconds.

## Notes

- Every item maps to an acceptance scenario, functional requirement, edge case, or success criterion.
