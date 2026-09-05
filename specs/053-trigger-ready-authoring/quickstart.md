# Quickstart: Trigger-Ready Task Authoring

## Prerequisites

- Go 1.25.0 toolchain with the repository's C compiler prerequisites.
- A disposable daemon data directory for CLI and restart scenarios.
- Fyne headless test environment for GUI scenarios.

## Scenario 1 - Completely Blank Draft

1. Create a task with no name, command, schedule, or one-off time through the API and desktop editor.
2. Confirm the task is stored disabled, displays as `unnamed`, reports not runnable, and returns `schedule: null`.
3. Restart the store or daemon and confirm the same values and stable identity remain.
4. Confirm Enable and Run now return a command-specific validation result without launching a process.

## Scenario 2 - Manual-Only Task

1. Create a task with a known harmless command and no automatic source.
2. Confirm it reports manual only and remains disabled.
3. Invoke Run now and confirm exactly one manual run reaches history.
4. Confirm scheduler recomputation and calendar listing create no occurrence and no missing-schedule error.

## Scenario 3 - Completion Source Transition

1. Create a runnable manual-only target and a fully configured source task.
2. Add a completion chain from source to target, then enable the target.
3. Confirm the target reports ready and retains no time schedule.
4. Delete the final incoming chain directly and through source-task deletion, then confirm the target becomes disabled and manual only in the same committed transition.

## Scenario 4 - Clear Existing Configuration

1. Create and enable a scheduled task.
2. Clear its command and confirm the task persists disabled and not runnable while its schedule remains.
3. Restore the command, clear the schedule, and confirm the task persists disabled and manual only.
4. Submit malformed replacement input and confirm no task or schedule field changes.

## Scenario 5 - Nameless Identity

1. Create at least ten nameless tasks with distinct commands or stable IDs.
2. Confirm every task label reads `unnamed` in Tasks, Groups, dialogs, CLI output, and detail summaries.
3. Edit, run, and delete selected records and confirm actions affect only their stable identities.

## Scenario 6 - Navigation and Group Keyboard Flow

1. Open light and dark modes at supported default, narrow, and wide window sizes.
2. Confirm Tasks, Groups, and Chains form the first section; Schedule, Activity, Options, and Info form the second; Exit remains separate at the bottom.
3. Confirm selected, hover, focus, badge, and keyboard behavior remains readable and ordered.
4. Open New Group, enter a valid name, press Enter, and confirm exactly one group is created.
5. Repeat with blank input, key repeat, composition, and immediate click-after-key cases and confirm zero or one creation as specified.

## Verification

```sh
sh scripts/verify.sh all
```

The aggregate must pass format, vet, lint, race, GUI, coverage, docs, and automation gates in the foreground.
