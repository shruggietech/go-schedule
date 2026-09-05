# Data Model: Trigger-Ready Task Authoring

## Task

Existing task identity and execution fields remain. S053 changes these invariants:

- `name`: stored string may be empty; user-facing fallback is `unnamed`.
- `command`: stored string may be empty; empty means command readiness is false.
- `schedule_id`: nullable reference; null means no time or startup activation source.
- `enabled`: false whenever command readiness or automatic activation readiness is false.
- `state`: existing active, completed, and disabled lifecycle remains distinct from readiness and enabled state.

## Schedule

Existing schedule rows remain immutable timing definitions. S053 does not add a draft or placeholder kind. A task may reference no schedule.

## Command Readiness

Derived value:

- `runnable`: true when the stored executable command is nonempty and valid at the persisted command boundary.
- `reason`: empty when runnable; otherwise a stable command-specific explanation.

This value controls manual Run now availability but is not persisted.

## Automatic Activation Readiness

Derived value:

- `ready`: true when the task is runnable, nonterminal, and has at least one valid automatic source.
- `sources`: zero or more of `schedule`, `startup`, `completion`, and later `trigger`.
- `reason`: stable explanation when no valid automatic source exists.

During S053 a nonempty valid schedule reference covers recurring, one-off, and startup kinds. At least one incoming completion chain adds `completion` independently of schedule presence.

## Effective Eligibility

Derived precedence for automatic dispatch and display:

1. Not runnable when command readiness is false.
2. Terminal when lifecycle is not active.
3. Manual only when command-ready but no automatic source exists.
4. Disabled when automatic activation-ready but local enabled is false.
5. Blocked when enabled but an ancestor group is disabled or invalid.
6. Ready when all preceding gates pass.

Manual Run now consults command readiness only and therefore remains available in states 3 through 6 unless another existing execution error applies.

## Group

No schema change. A group with no child groups or member tasks is valid. Group name validation and hierarchy semantics remain unchanged.

## Navigation Section

Presentation-only metadata:

- `task-definition`: Tasks, Groups, Chains, and later Triggers.
- `operation`: Schedule, Activity, Options, and Info.

Destination IDs remain stable. Exit remains outside both sections.

## State Transitions

- Create with missing command or automatic source: persist disabled.
- Add valid command to an unscheduled draft: becomes manual-only and stays disabled.
- Add first schedule or incoming completion chain to a runnable task: becomes activation-ready and stays disabled until explicitly enabled.
- Remove command from enabled task: atomically disable and become not runnable.
- Remove final automatic source from enabled task: atomically disable and become manual-only when the command remains valid.
- Enable request for a not-runnable, manual-only, or terminal task: reject without mutation.
- Manual Run now for missing command: reject before dispatch.
- Manual Run now for runnable disabled, blocked, or manual-only task: preserve existing explicit operator behavior and dispatch.
