# Contract: Per-task DST scheduling policy

## Wire values

Task create, task update, schedule preview, and task responses use:

```json
{
  "time_basis": "wall_clock",
  "dst_gap_policy": "next_valid",
  "dst_overlap_policy": "first"
}
```

Omitted create/preview values use the defaults shown. Omitted update values leave stored choices unchanged.
Task detail and schedule preview also return `policy_summary`, a human-readable
description of the effective basis and whether transition choices are active.
Task detail exposes `schedule.elapsed_epoch` for elapsed schedules. It is the
persisted UTC phase and does not change when only the task timezone changes.

## Validation errors

- `time_basis`: `wall_clock`, `elapsed`, or `utc`
- `dst_gap_policy`: `next_valid` or `skip`
- `dst_overlap_policy`: `first`, `both`, or `last`
- incompatible elapsed recurrence: validation error on `time_basis`

Failures do not create a schedule or mutate a task.

## CLI

- `gosched task add|edit --time-basis wall_clock|elapsed|utc`
- `gosched task add|edit --dst-gap next_valid|skip`
- `gosched task add|edit --dst-overlap first|both|last`
- `gosched task show` names all effective values.

## Desktop

Advanced Settings contains friendly selectors for Time basis, Spring gap, and
Fall overlap. Stored wire values remain the enum strings above. Live preview
includes the current selections, and Save remains disabled when elapsed mode is
incompatible with the entered recurrence.

## Recurrence examples

For New York spring-forward in 2026, wall-clock every six hours anchored at 09:00 retains local 09:00/15:00/21:00/03:00 and therefore includes a five-hour absolute gap. Elapsed retains six-hour absolute gaps and shifts local display. UTC keeps recurrence fields fixed in UTC.

For New York fall-back at local 01:30, `first`, `both`, and `last` select the earlier instant, both ordered instants, and the later instant respectively.
