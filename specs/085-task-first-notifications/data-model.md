# Data Model: Task-First Notifications

S085 adds no persisted entity. It extends the secret-free desktop read projection and adds transient presentation state.

## Configured scope summary

- **Scope type**: Task or group.
- **Identity**: Existing stable scope identifier.
- **Name and context**: Existing human-readable scope values.
- **Assignment source**: Direct task, inherited group, direct group, or none.
- **Outcome coverage**: Success, failure, or both.
- **Destination counts**: Configured destinations and the subset currently enabled.
- **Completeness**: Workspace-level flag indicating whether all requested summaries were resolved.
- **Security boundary**: Contains no endpoint, authorization, credential, payload, or secret-derived value.

## Notification overview

- **Enabled destinations**: Count of currently enabled channels.
- **Configured tasks and groups**: Counts and named summaries from the complete coverage projection.
- **Recent outcomes**: At most five newest delivery summaries.
- **Overall state**: Setup needed, disabled, healthy, in progress, retrying, or failed.
- **Guidance**: Plain-language meaning and one relevant next action or an explicit no-action-needed statement.

## Recent result

- **Identity**: Existing delivery ID.
- **Context**: Task name for task outcomes or Test notification for transport tests.
- **Destination**: Secret-free channel name.
- **State**: Queued, sending, retrying, successful, or failed.
- **Time**: Existing created timestamp.
- **Guidance**: Deterministic state meaning and next action.

## Advanced section state

- **Destinations open**: Local disclosure state, initially false.
- **Assignment rules open**: Local disclosure state, initially false.
- **Delivery diagnostics open**: Local disclosure state, initially false.
- **Transition**: Summary activation toggles exactly its section; channel or policy actions remain within that section.
