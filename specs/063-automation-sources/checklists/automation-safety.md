# Automation Source Safety and UX Checklist

**Purpose**: Test the completeness, clarity, and consistency of the automation-source requirements before implementation

**Created**: 2026-09-07

**Feature**: [spec.md](../spec.md)

## Shared Mental Model

- [x] Every source type has a defined initiating source or selection rule.
- [x] Every source type has a defined target-task presentation.
- [x] Disabled, missing, degraded, and ready states require textual explanations.
- [x] Search, filtering, and large-collection behavior have measurable expectations.
- [x] Deleted sources and targets have explicit recovery behavior.

## Lifecycle Completeness

- [x] Completion-chain create, inspect, edit, retarget, and delete behavior is required.
- [x] External-trigger create, inspect, edit, toggle, rotate, reveal, copy, fire, retarget, and delete behavior is required.
- [x] Trigger Set create, inspect, toggle, rotate, reveal, copy, retarget, and delete behavior is required.
- [x] Filesystem-watcher create, inspect, edit, toggle, retarget, and delete behavior is required.
- [x] Unsupported lifecycle actions are not invented for source types that lack them.
- [x] Trigger Set ordering and atomic mutations are preserved.

## Secret Boundary

- [x] Ordinary list and workspace contracts explicitly exclude raw keys.
- [x] Create, reveal, and rotate are the only secret-producing actions.
- [x] Fire consumes a key only inside the backend boundary.
- [x] Secret dialogs have explicit clearing conditions.
- [x] Logs, events, diagnostics, screenshots, and errors are included in the no-secret requirement.
- [x] Copy feedback is scoped to an explicit user action.

## Reliability and Concurrency

- [x] Stale and deleted entity updates have distinct expected outcomes.
- [x] Duplicate mutation suppression is measurable.
- [x] Successful mutation refresh behavior is defined.
- [x] Failed mutation snapshot preservation is defined.
- [x] Offline read-only continuity is defined.
- [x] Partial list failures cannot masquerade as a complete snapshot.
- [x] Reconnect ordering cannot replace newer state with an older response.

## Accessibility and Platform Behavior

- [x] Every interaction is required to work with a keyboard.
- [x] Status is not communicated by color alone.
- [x] Dialog focus management and screen-reader semantics are required.
- [x] Long and platform-specific watcher paths are preserved.
- [x] Watcher health remains daemon-owned and accurately presented.

## Verification

- [x] Go contract and service tests are required.
- [x] React lifecycle and accessibility tests are required.
- [x] Secret-leak regression tests are required.
- [x] Large-collection and focus-retention tests are required.
- [x] Canonical repository verification is required before publication.

## Notes

- All checklist items are supported by explicit requirements, acceptance scenarios, edge cases, or measurable outcomes in the specification.
