# Data Model: Explicit DST Scheduling Intent

## Task additions

| Field | Values | Default | Applies when |
| --- | --- | --- | --- |
| `time_basis` | `wall_clock`, `elapsed`, `utc` | `wall_clock` | recurring schedules |
| `dst_gap_policy` | `next_valid`, `skip` | `next_valid` | wall-clock recurrence |
| `dst_overlap_policy` | `first`, `both`, `last` | `first` | wall-clock recurrence |

The fields belong to Task, remain independent of Schedule replacement, and are returned in the task JSON object. Empty in-memory values normalize to defaults before evaluation or persistence.

An elapsed recurring Schedule also carries nullable `elapsed_epoch`, the UTC
instant from which fixed durations are counted. It is bound once in the
authoring timezone and then persists independently of the task's presentation
timezone. Wall-clock and UTC schedules leave it empty.

## SchedulePolicy value

The evaluator receives a value containing missing-date policy, time basis, DST gap policy, and DST overlap policy. The value exposes effective defaults and is immutable during one evaluation.

## Validation rules

- Unknown enum values are rejected at preview/create/update boundaries.
- `elapsed` requires a recurring fixed-duration shape with one occurrence per base period.
- `wall_clock` and `utc` accept every existing recurring shape.
- One-off and event schedules retain fields without applying them.
- Schedule replacement does not reset any task-level policy.

## Schema transition

Migration v7 adds the three non-null task columns with defaults and nullable
`schedules.elapsed_epoch`. All pre-v7 rows therefore read as `wall_clock`,
`next_valid`, and `first`, with no elapsed epoch. No reverse migration is
provided.

## Occurrence transition

1. Parse the authoritative schedule.
2. Select recurrence intent under the chosen basis.
3. For wall-clock schedules, resolve missing date and calendar adjustment.
4. Resolve wall time using gap and overlap policies.
5. Sort concrete UTC candidates, discard candidates not strictly after the cursor, and suppress duplicates.
6. Return the earliest remaining instant.
