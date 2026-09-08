# Contract: Desktop Notifications Bridge

## Workspace

### `NotificationWorkspace() NotificationResult`

Returns one complete secret-free snapshot containing current channel summaries, task and group policy scopes, and up to 200 newest redacted deliveries. Any constituent read failure returns unavailable without a partial workspace.

## Channels

### `SaveNotificationChannel(draft) NotificationResult`

Creates or updates a channel. Create requires name and HTTPS endpoint. Update changes protected values only when their explicit replacement flag is true. A blank authorization replacement clears it; a blank endpoint replacement is rejected. Success returns a complete refreshed workspace and never echoes protected draft values.

### `SetNotificationChannelEnabled(id, enabled) NotificationResult`

Enables or disables one exact non-empty channel and returns a complete refreshed workspace.

### `TestNotificationChannel(id) NotificationResult`

Queues one transport-equivalent test for an enabled channel and returns a complete refreshed workspace whose history labels the record as Test.

### `DeleteNotificationChannel(id) NotificationResult`

Removes one exact channel after frontend confirmation and returns a complete refreshed workspace. Terminal redacted history may remain according to the server lifecycle.

## Policies

### `NotificationPolicy(scopeType, scopeId) NotificationResult`

Returns direct assignments for one selected task or group. Task results additionally contain the server-selected effective source and assignments. Group results identify the selected group as direct configuration and explain conditional descendant inheritance.

### `SaveNotificationPolicy(draft) NotificationResult`

Atomically replaces the complete direct assignment list for one selected task or group. Every assignment needs a channel and at least one outcome. Success returns refreshed policy detail.

## Shared guarantees

- Every daemon call has a bounded timeout.
- Result types contain no endpoint, authorization, or payload field.
- Errors use a bounded safe vocabulary and never reflect arbitrary daemon content.
- Reads preserve only complete snapshots, and stale frontend completions are discarded.
- Relevant domain events are debounced.
- Mutation controls are disabled while disconnected or pending.
- Duplicate frontend activation is suppressed.
