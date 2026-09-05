# Cron Fidelity Checklist: Cron Parity Closure

**Purpose**: Validate that the requirements define faithful behavior and explicit boundaries before implementation **Created**: 2026-08-29 **Feature**: [spec.md](../spec.md)

## Import Context

- [x] Assignment ordering and independent job snapshots are specified
- [x] `CRON_TZ` and process `TZ` have distinct contracts
- [x] Shell selection and original command preservation are explicit
- [x] Percent splitting, escaping, newline conversion, and stdin are defined
- [x] User and system layouts are independently selectable
- [x] Unix and Quartz file layouts are explicit rather than guessed
- [x] Mail directives have a visible deferred outcome

## Expression Fidelity

- [x] Five-field default seconds and Unix weekday numbering are specified
- [x] Six-field seconds, Quartz weekday numbering, and `?` placement are specified
- [x] Invalid year, day, and modifier combinations have non-approximating outcomes
- [x] Compilation, description, source retention, execution, and export must agree
- [x] Five-field export compatibility and six-field selection are explicit

## Lifecycle and Boundaries

- [x] Standard-input persistence, clearing, restart, and execution are specified
- [x] Lossy operational export is explicitly refused
- [x] Windows run-as remains a visible platform validation boundary
- [x] Every issue #22 row requires one documented disposition
- [x] Public copy must match the tested subset
- [x] Migration, performance, encoding, and canonical verification criteria are measurable

## Notes

- All 18 items pass. The checklist rejects structural guessing and lossy output as violations even when they would reduce implementation work.
