# Quickstart: Trigger Sets

## Prerequisites

Run a local daemon from the S055 branch and create one command-ready task. Record its task ID.

## Create and Copy

Create a three-member set with `gosched trigger set create --name "Deploy callers" --task <task-id> --count 3`. Confirm three complete commands appear in ascending position with no blank lines. Repeat with `--json` and confirm set ID, positions, trigger IDs, keys, and commands are structured and ordered.

## Exercise Individual Isolation

Use `gosched trigger list` to select one member. Rename, disable, enable, rotate, and delete it through the existing individual commands. Confirm sibling identities, positions, keys, target, and enabled state remain unchanged. Attempt an individual target update and confirm the error directs you to the set retarget command.

## Exercise Set Lifecycle

Use `show`, `retarget`, `disable`, `enable`, `reveal`, and `rotate` under `gosched trigger set`. Confirm ordinary output is redacted, secret output is ordered, old rotated keys are unknown, and every member changes together. Delete the set and confirm all former keys are unknown.

## Desktop Validation

Open Triggers, create a set, inspect Set and Position values, copy all commands, and exercise each set-level action through a selected member. Confirm broad confirmations name the set and count, live refresh shows no intermediate state, and light and dark modes keep every control readable.

## Regression and Verification

Confirm standalone triggers remain ungrouped and fully functional after migration. Run `sh scripts/verify.sh all` and require all eight gates to pass. Inspect ordinary API responses, events, logs, Activity, history, and errors for zero raw-key occurrences.
