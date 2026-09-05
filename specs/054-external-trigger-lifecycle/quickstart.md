# Quickstart: External Triggers

## Create and invoke a trigger

1. Create a runnable task and activate it.
2. Open Triggers in the desktop application or run `gosched trigger create --name "Build hook" --task <task-id>`.
3. Copy the returned key or the complete fire command.
4. Run `gosched trigger fire <key>` from another local process.
5. Confirm one run appears with source `external_trigger` and the trigger identifier.

## Exercise lifecycle controls

1. Disable the trigger and confirm firing returns `trigger_disabled` without creating a run.
2. Re-enable it and confirm firing succeeds.
3. Rotate the key and confirm the old key returns `trigger_unknown` while the new key succeeds.
4. Delete the trigger and confirm its historical run provenance remains visible.

## Verify the implementation

Run `sh scripts/verify.sh all` and confirm every mandatory gate passes. Also review logs, event payloads, history output, and ordinary trigger responses to ensure no raw key is present.
