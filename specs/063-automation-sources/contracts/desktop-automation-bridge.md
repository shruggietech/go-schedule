# Contract: Desktop Automation Bridge

## Secret-free queries

### `AutomationWorkspace() OperationResult`

Loads task choices and all four source collections as one complete, secret-free snapshot. Any failed constituent read returns `unavailable` without a partial workspace.

## Completion chains

### `SaveChain(ChainDraft) OperationResult`

Creates or updates a chain after local required-field validation and stale-write detection. A successful result includes the refreshed workspace.

### `DeleteChain(id) OperationResult`

Deletes a chain once and returns the refreshed workspace.

## External triggers

### `SaveTrigger(TriggerDraft) SecretResult | OperationResult`

Create returns an ephemeral secret result. Update returns an ordinary result and refreshed workspace.

### `SetTriggerEnabled(id, enabled) OperationResult`

Toggles one standalone trigger and returns the refreshed workspace.

### `RevealTrigger(id) SecretResult`

Returns the current key and command only after explicit invocation.

### `RotateTrigger(id) SecretResult`

Rotates the key and returns the replacement only after explicit invocation.

### `FireTrigger(id) OperationResult`

Resolves the key and fires inside Go. The returned result contains no key or command.

### `DeleteTrigger(id) OperationResult`

Deletes one standalone trigger and returns the refreshed workspace.

## Trigger Sets

### `CreateTriggerSet(TriggerSetDraft) SecretResult`

Creates the set atomically and returns ordered ephemeral member secrets.

### `RetargetTriggerSet(id, targetTaskId, originalUpdatedAt, overwriteStale) OperationResult`

Detects stale state, retargets every member atomically, and returns the refreshed workspace.

### `SetTriggerSetEnabled(id, enabled) OperationResult`

Toggles every member atomically and returns the refreshed workspace.

### `RevealTriggerSet(id) SecretResult`

Returns ordered current member secrets after explicit invocation.

### `RotateTriggerSet(id) SecretResult`

Rotates every member atomically and returns ordered replacements after explicit invocation.

### `DeleteTriggerSet(id) OperationResult`

Deletes the set and its members atomically and returns the refreshed workspace.

## Filesystem watchers

### `SaveFilesystemWatcher(WatcherDraft) OperationResult`

Creates or updates a watcher after local validation and stale-write detection, then returns daemon health in the refreshed workspace.

### `SetFilesystemWatcherEnabled(id, enabled) OperationResult`

Toggles one watcher and returns daemon-owned health in the refreshed workspace.

### `DeleteFilesystemWatcher(id) OperationResult`

Deletes one watcher and returns the refreshed workspace.

## Shared guarantees

- Ordinary results and workspace models never contain raw keys.
- Status errors map only recognized fields and never reflect backend payloads containing secrets.
- Every backend call has a bounded timeout.
- Mutating frontend controls are disabled while disconnected or pending.
- Successful mutations replace state only with a complete refreshed workspace.
- Stale async loads cannot replace newer state.
